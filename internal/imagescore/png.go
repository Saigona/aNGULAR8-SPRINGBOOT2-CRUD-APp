package imagescore

import (
	"context"
	"image"
	"image/png"
	"sync"
)

type PngScorer struct {
	enc png.Encoder
	uncompressedImageSizeCache
}

var _ ImageScorer = &PngScorer{}

func NewPngScorer() *PngScorer {
	return &PngScorer{
		enc: png.Encoder{
			CompressionLevel: png.BestSpeed,
			BufferPool:       &singleThreadBufferPool{},
		},
	}
}

func (ps *PngScorer) ScoreImage(ctx context.Context, img image.Image) (float64, error) {
	buf := &discardCounter{}

	err := p