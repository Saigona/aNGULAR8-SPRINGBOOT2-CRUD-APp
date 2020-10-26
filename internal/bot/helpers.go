package bot

import (
	"context"
	"time"

	"github.com/WIZARDISHUNGRY/hls-await/internal/logger"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/pkg/errors"
)

var loc = func() *time.Location {
	l, err := time.LoadLocation("Asia/Pyongyang")
	if err != nil {
		panic(err)
	}
	return l
}()

func getLastTweet(c *twitter.Client) (int64, time.T