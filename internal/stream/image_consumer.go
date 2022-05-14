
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
		case <-ctx.Done():
			return nil
		case <-s.oneShot:
			if s.flags.OneShot {
				oneShot = true
				log.Println("photo time!")
			}
		case img := <-s.imageChan:
			err := s.consumeImage(ctx, sem, filterFunc, img, oneShot, frameCount)
			if err != nil {
				return err
			}
			if oneShot {
				oneShot = false
			}
			frameCount++
		}
	}
}

// consumeImage cannot be spawned in its own goroutine because we must ensure the heap is updated synchonously
func (s *Stream) consumeImage(ctx context.Context,
	sem *semaphore.Weighted,
	filterFunc filter.FilterFunc,
	img image.Image,
	oneShot bool,
	frameCount int,
) error {
	ctx, log := withFrameCount(ctx, frameCount)

	entry := &outputImageEntry{
		counter: frameCount,
		done:    false,
		image:   img,
	}

	s.outputImagesMutex.Lock()
	s.outputImages.Push(entry)
	s.outputImagesMutex.Unlock()

	if err := s.drawImage(ctx, img, oneShot, entry); err != nil {
		log.WithError(err).Warn("drawImage draw")
	}
	err := sem.Acquire(ctx, 1)
	if err != nil {
		return err
	}
	go func() {
		defer sem.Release(1)
		err := s.consumeImageAsync(ctx, filterFunc, img, oneShot, entry)
		if err != nil {
			log.WithError(err).Warn("consumeImageAsync")
		}
	}()
	return nil
}

func (s *Stream) consumeImageAsync(ctx context.Context,
	filterFunc filter.FilterFunc,
	img image.Image,
	oneShot bool,
	entry *outputImageEntry,
) error {

	var sendToBot bool

	defer func() { // When finishing filtering an image try to dequeue all done images
		s.outputImagesMutex.Lock()
		defer s.outputImagesMutex.Unlock()
		entry.done = true // done
		entry.passedFilter = sendToBot

		for {
			if s.outputImages.Len() == 0 {
				return
			}
			entry := s.outputImages.Pop()
			if entry == nil {
				return