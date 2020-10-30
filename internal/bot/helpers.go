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

func getLastTweet(c *twitter.Client) (int64, time.Time, error) {
	u, _, err := c.Accounts.VerifyCredentials(&twitter.AccountVerifyParams{
		IncludeEntities: twitter.Bool(true),
	})
	if err != nil {
		return 0, time.Time{}, errors.Wrap(err, "VerifyCredentials")
	}
	tws, _, err := c.Timelines.UserTimeline(&twitte