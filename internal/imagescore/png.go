package imagescore

import (
	"context"
	"image"
	"image/png"
	"sync"
)

type PngScorer struct {
	enc png.Encoder
	un