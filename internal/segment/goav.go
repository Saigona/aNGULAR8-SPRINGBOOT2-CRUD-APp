
package segment

import (
	"context"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"sync"
	"time"
	"unsafe"

	"github.com/WIZARDISHUNGRY/hls-await/internal/bot"
	"github.com/WIZARDISHUNGRY/hls-await/internal/logger"
	"github.com/charlestamz/goav/avcodec"
	"github.com/charlestamz/goav/avformat"
	"github.com/charlestamz/goav/avutil"
	"github.com/charlestamz/goav/swscale"
	old_avutil "github.com/giorgisio/goav/avutil"
	"github.com/pkg/errors"
)

const (
	minImages = 4
)

type GoAV struct {
	Context        context.Context
	VerboseDecoder bool
	RecvUnixMsg    bool // use a secondary unix socket to receive file descriptors in priv sep mode
	RestartOnPanic bool
	FDs            chan uintptr
	DoneCB         func(error)
}
