package stream

import (
	"bufio"
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/grafov/m3u8"
)

const (
	minPollDuration = time.Second
	maxPollDuration = time.Minute
)

func (s *Stream) doPlaylist(ctx context.Context, u *url.URL) (*m3u8.MediaPlaylist, error) {
	resp, err := s.httpGet(ctx, u.String())
	if err != nil {
		re