package stream

import (
	"context"
	"fmt"
	"image"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/WIZARDISHUNGRY/hls-await/internal/bot"
	"github.com/WIZARDISHUNGRY/hls-await/internal/logger"
	my_roku "github.com/WIZARDISHUNGRY/hls-await/internal/roku"
	"github.com/WIZARDISHUNGRY/hls-await/internal/worker"
	"github.com/WIZARDISHUNGRY/hls-await/pkg/heap"
	"github.com/WIZARDISHUNGRY/hls-await/pkg/proxy"
	"github.com/looplab/fsm"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"jonwillia.ms/roku"
)

type StreamOption func(s *Stream) error

func NewStream(opts ...StreamOption) (*Stream, error) {

	s := newStream()

	for _, opt := range opts {
		err := opt(s)
		if err != nil {
			return nil, err
		}
	}

	if !s.flags.Worker {
		target, err := s.url.Parse("/")
		if err != nil {
			return nil, err
		}
		u, err := proxy.NewSingleHostReverseProxy(context.TODO(), target, false)
		if err != nil {
			return nil, errors.Wrap(err, "NewSingleHostReverseProxy")
		}
		u.Path = s.url.Path
		s.url = u
	}

	return s, nil
}

func WithURL(u *url.URL) StreamOption {
	return func(s *Stream) error {
		s.url = u
		return nil
	}
}

type Stream struct {
	rokuCB        func() (*roku.Remote, error)
	url, proxyURL *url.URL

	oneShot    chan struct{}
	imageChan  chan image.Image
	flags      *flags
	segmentMap map[url.URL]struct{