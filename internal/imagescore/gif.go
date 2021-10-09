package imagescore

import (
	"context"
	"image"
	"image/gif"
)

type GifScorer struct {
	uncompressedImageSizeCache
}

var _ ImageScorer = &GifScor