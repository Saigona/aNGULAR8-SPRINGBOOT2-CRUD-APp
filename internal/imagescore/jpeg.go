package imagescore

import (
	"context"
	"image"
	"image/jpeg"
)

type JpegScorer struct {
	uncompressedImageSizeCache
}

var _ ImageScorer = &JpegScorer{}

func NewJpegScorer() *JpegScorer { return &JpegScorer{} }

func (js *JpegScorer) ScoreImage(ctx context.Context, img image.Imag