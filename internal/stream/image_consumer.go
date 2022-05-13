
package stream

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"os"
	"runtime"
	"sync/atomic"

	"github.com/WIZARDISHUNGRY/hls-await/internal/corpus"
	"github.com/WIZARDISHUNGRY/hls-await/internal/filter"
	"github.com/WIZARDISHUNGRY/hls-await/internal/imagescore"
	"github.com/WIZARDISHUNGRY/hls-await/internal/logger"
	"github.com/eliukblau/pixterm/pkg/ansimage"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
	"golang.org/x/sys/unix"
)

const (
	goimagehashDim = 8    // should be power of 2, color bars show noise at 16
	imagescoreMin  = 0.06 // GZIP specific, calculated from output of TestScoringAlgos
)

func withFrameCount(ctx context.Context, frameCount int) (context.Context, *logrus.Entry) {
	log := logger.Entry(ctx).WithField("frame_count", frameCount)
	logger.WithLogEntry(ctx, log)
	return ctx, log
}

// We picked gzip because it had the best results and good speed + low allocs
var imageScorerAlgo = imagescore.NewGzipScorer

func (s *Stream) consumeImages(ctx context.Context) error {
	log := logger.Entry(ctx)

	c, err := corpus.LoadEmbedded("testpatterns")
	if err != nil {
		return errors.Wrap(err, "corpus.Load")
	}

	bs := imagescore.NewBulkScore(ctx,
		func() imagescore.ImageScorer {
			return imageScorerAlgo()
		},
	)

	filterFunc := filter.Multi(
		filter.Motion(goimagehashDim, s.flags.Threshold),
		filter.DefaultMinDistFromCorpus(c),
		imagescore.Filter(bs, imagescoreMin),
	)

	var frameCount int

	var maxFramesInFlight = runtime.GOMAXPROCS(-1) * 4 // a large number
	sem := semaphore.NewWeighted(int64(maxFramesInFlight))

	oneShot := false
	for {
		select {