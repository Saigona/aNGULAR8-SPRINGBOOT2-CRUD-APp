package stream

import (
	"context"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"

	"github.com/WIZARDISHUNGRY/hls-await/internal/logger"
)

func (s *Stream) httpGet(ctx context.Context, url string) (*http.Response, error) {
	log := logger.Entry(ctx)