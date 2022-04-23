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

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Referer", "https://kcnawatch.org/korea-central-tv-livestream/") // TODO flag
	req.Header.Set("Accept", "*/*")
	// req.Header.Set("Cookie", "