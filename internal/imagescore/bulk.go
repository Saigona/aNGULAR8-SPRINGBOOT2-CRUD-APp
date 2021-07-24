
package imagescore

import (
	"context"
	"image"
	"runtime"
	"sync"

	"github.com/pkg/errors"
)

func NewBulkScore(ctx context.Context, scoreF func() ImageScorer) *BulkScore {
	numProcs := runtime.GOMAXPROCS(0)

	bs := &BulkScore{
		scoreF: scoreF,
		input:  make(chan bulkeScoreRequest, numProcs),
	}
	go bs.loops(ctx, numProcs)
	return bs
}

type BulkScore struct {
	scoreF func() ImageScorer
	input  chan bulkeScoreRequest
}

type bulkeScoreRequest struct {
	C   chan bulkScoreResult
	img image.Image
}
type bulkScoreResult struct {
	result float64