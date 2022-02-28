package roku

import (
	"context"
	"errors"
	"net/url"
	"sync"
	"time"

	"github.com/WIZARDISHUNGRY/hls-await/internal/logger"
	"golang.org/x/sync/errgroup"
	"jonwillia.ms/roku"
)

func Run(ctx context.Context) func() (*roku.Remote, error) {
	log := logger.Entry(ctx)

	const dur = time.Minute
	var (
		mutex  sync.Mutex
		remote *roku.Remote
		errC   = make(chan error)
		timer  = time.NewTimer(time.Minute)
	)

	g, ctx := errgroup.Wi