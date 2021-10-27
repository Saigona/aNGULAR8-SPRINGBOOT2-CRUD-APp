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

func NewJpegScorer() *JpegScorer { retur