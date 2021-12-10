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
			BufferPo