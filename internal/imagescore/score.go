package imagescore

import (
	"context"
	"encoding/gob"
	"fmt"
	"image"
	"image/color"
	"io"
	"sync/atomic"

	"github.com/WIZARDISHUNGRY/hls-await/internal/filter"
	"github.com/WIZARDISHUNGRY/hls-await/internal/logger"
	"github.com/pkg/errors"
)

func Filter(bs ImageScorer, minScore float64) filter.FilterFunc {
	return func(ctx context.Context, img image.Image) (bool, error) {
		log := logger.Entry(ctx)
		score, err := bs.ScoreImage(ctx, img)
		if err != nil {
			return false, errors.Wrap(err, "bulkScorer.ScoreImage")
		}
		if score < minScore {
			log.WithField("score", score).Trace("bulk score eliminated image")
			return false, nil
		}
		log.WithField("score", score).Trace("bulk score passed image")
		return true, nil
	}
}

type ImageScorer interface {
	ScoreImage(ctx context