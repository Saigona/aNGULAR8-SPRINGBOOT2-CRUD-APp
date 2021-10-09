package imagescore

import (
	"context"
	"image"
	"image/gif"
)

type GifScorer struct {
	uncompressedImageSizeCache
}

var _ ImageScorer = &GifScorer{}

func NewGifScorer() *GifScorer { return &GifScorer{} }

func (ps *GifScorer) ScoreImage(ctx conte