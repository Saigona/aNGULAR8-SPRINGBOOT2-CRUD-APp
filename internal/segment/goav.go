
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

var pngEncoder = &png.Encoder{
	CompressionLevel: png.NoCompression,
}

var _ Handler = &GoAV{}

var onceAvcodecRegisterAll sync.Once

func (goav *GoAV) HandleSegment(req *Request, resp *Response) (err error) {

	defer func() {

	}()

	ctx := goav.Context
	log := logger.Entry(ctx)

	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}
			if err != nil {
				log.WithError(err).Error("HandleSegment panic handler")
			}
		}
		if goav.DoneCB != nil {
			goav.DoneCB(err)
		}
	}()

	defer func() { fractionImages(ctx, resp, err) }()

	onceAvcodecRegisterAll.Do(func() {
		avcodec.AvcodecRegisterAll() // only instantiate if we build a GoAV
		// go func() { // put this in here to test ffmpeg crashes
		// 	time.Sleep(10 * time.Second)
		// 	panic("test restarting")
		// }()
	})

	fd := req.FD
	if goav.RecvUnixMsg {
		var ok bool
		fd, ok = <-goav.FDs
		if !ok {
			return fmt.Errorf("fd channel closed")
		}
	}

	if fd <= 0 {
		return fmt.Errorf("fd is weird %d", fd)
	}

	defer func() {
		f := os.NewFile(fd, "whatever")
		if f != nil {
			f.Close()
		}
	}()

	// file := fmt.Sprintf("pipe:%d", fd)
	file := fmt.Sprintf("/proc/self/fd/%d", fd) // This is a Linuxism, but it works. Otherwise we get like 10% of images (30 instead of 248)

	pFormatContext := avformat.AvformatAllocContext()
	if pFormatContext == nil {
		return errors.New("pFormatContext == nil")
	}

	// if avformat.AvformatOpenInput(&pFormatContext, "x.ts", nil, nil) != 0 {
	if e := avformat.AvformatOpenInput(&pFormatContext, file, nil, nil); e != 0 {
		return errors.Wrap(goavError(e), "couldn't open file.")
	}
	defer avformat.AvformatCloseInput(pFormatContext)

	// Retrieve stream information (TODO: is this needeed)
	// var dict *avutil.Dictionary
	// if e := pFormatContext.AvformatFindStreamInfo(&dict); e < 0 {
	if e := pFormatContext.AvformatFindStreamInfo(nil); e < 0 {
		return errors.Wrap(goavError(e), "couldn't find stream information.")
	}
	// defer avutil.AvDictFree(&dict)

	// Dump information about file onto standard error
	if goav.VerboseDecoder {
		pFormatContext.AvDumpFormat(0, file, 0)