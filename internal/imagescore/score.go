package imagescore

import (
	"context"
	"encoding/gob"
	"fmt"
	"image"
	"image/color"
	"io"
	"sync/atomic"

	"github.com/WIZARDISHUNGRY/hls-await/internal/filter"
	"github.com/WIZARDISHUNGRY/hls-await/internal/logger"
	"github.com/pkg/errors"
)

func Filter(bs ImageScorer, minScore floa