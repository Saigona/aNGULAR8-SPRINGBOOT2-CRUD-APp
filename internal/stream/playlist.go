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
	maxPol