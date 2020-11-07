
package bot

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"time"

	"github.com/WIZARDISHUNGRY/hls-await/internal/logger"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

const (
	maxAge = 2 * updateIntervalMinutes * time.Minute // Don't post old images
	minAge = 90 * time.Second
)

func (b *Bot) maybeDoPost(ctx context.Context, srcImages []imageRecord) ([]imageRecord, error) {
	log := logger.Entry(ctx)

	const mimeType = "image/png"

	if len(srcImages) < numImages {
		log.WithField("num_images", len(srcImages)).Info("not enough images to post")
		return srcImages, nil
	}

	// Don't post old images
	firstGood := -1
	for i, img := range srcImages {
		if firstGood == -1 && time.Now().Sub(img.time) < maxAge {
			firstGood = i
		}
	}

	if firstGood < 0 {
		log.WithField("num_images", len(srcImages)).Info("discarding all images due to age")
		return nil, nil
	}
	if firstGood != 0 {
		log.WithField("num_images", len(srcImages)-firstGood).Info("discarding some images due to age")
		srcImages = srcImages[firstGood:]
	}

	if len(srcImages) < numImages {
		log.WithField("num_images", len(srcImages)).Info("not enough images to post")
		return srcImages, nil
	}