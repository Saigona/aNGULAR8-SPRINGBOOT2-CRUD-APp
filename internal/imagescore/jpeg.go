package imagescore

import (
	"context"
	"image"
	"image/jpeg"
)

type JpegScorer struct {
	uncompressedImageSizeC