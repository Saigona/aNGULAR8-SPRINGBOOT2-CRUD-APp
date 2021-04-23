package filter

import (
	"context"
	"image"
	"testing"

	"github.com/WIZARDISHUNGRY/hls-await/internal/corpus"
	"github.com/WIZARDISHUNGRY/hls-await/internal/logger"
	"github.com/sirupsen/logrus"
)

//go:generate sh -c "go test ./... -run '^$' -benchmem -bench .