package imagescore

import (
	"context"
	"image"
	"sync"
	"testing"
)

//go:generate sh -c "go test ./... -r