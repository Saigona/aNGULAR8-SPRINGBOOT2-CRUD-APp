package filter

import (
	"context"
	"image"
)

// FilterFunc is a function that evaluates if a function passes a filter. It must be safe for concurrent use.
type FilterFunc func(context.Context, image.Image) (bool, error)

func 