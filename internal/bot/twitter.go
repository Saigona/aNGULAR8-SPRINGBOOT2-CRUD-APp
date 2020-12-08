package bot

import (
	"context"
	"image"
	"path/filepath"
	"runtime"
	"time"

	"github.com/WIZARDISHUNGRY/hls-await/internal/logger"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
)

const (
	TWITTER_CONSUMER_KEY    = "TWITTER_CONSUMER_KEY"
	TWITTER_CONSUMER_SECRET = "TWITTER_CONSUMER_SECRET"
	TWITTER_ACCESS_TOKEN    = "TWITTER_ACCESS_TOKEN"
	TWITTER_ACCESS_SECRET   = "TWITTER_ACCESS_SECRET"

	updateIntervalMinutes = 10
	updateInterval        = updateIntervalMinutes * time.Minute
	minUpdateInterval     = 20 * time.Second // used for tweeting quickly after a manual restart
	numImages             = 1                // per post, think 4 is max for twitter
	maxQueuedIntervals    = 4
	maxQueuedImages       = 25 * updateIntervalMinutes * 60 * maxQueuedIntervals * ImageFraction // about 4 updateIntervals at 25fps x the image fraction
	maxQueuedImagesMult   = 1.5
	replyWindow           = 3 * updateInterval
	ImageFraction         = (4 / FPS) // this is the proportion of images that make it from the decoder to here
	FPS                   = 25.0      // assume source is 25 fps PAL
	postTimeout           = 1 * time.Minute
)

var (
	_, b, _, _ = runtime.Caller(0)

	// root folder of this project. TODO: does not work with trimpath or builds
	root = filepath.Join(filepath.Dir(b), "../..")
)

func newClient() *twitter.Client {
	pa