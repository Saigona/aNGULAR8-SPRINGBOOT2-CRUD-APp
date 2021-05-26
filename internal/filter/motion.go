package filter

import (
	"context"
	"image"
	"sync"

	"github.com/WIZARDISHUNGRY/hls-await/internal/logger"
	"github.com/corona10/goimagehash"
	"github.com/pkg/errors"
)

// Motion returns a filter function that rejects images that fall under a threshold when
// comparing the ExtPerceptionHash against the previous image hash.
func Motion(dim, minDis