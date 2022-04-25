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
	// req.Header.Set("Cookie", " __qca=P0-44019880-1616793366216; _ga=GA1.2.978268718.1616793363; _gid=GA1.2.523786624.1616793363")
	req.Header.Set("Accept-Language", "en-us")
	req.Header.Se