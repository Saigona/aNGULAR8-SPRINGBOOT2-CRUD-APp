
package stream

import (
	"flag"
	"strings"

	"github.com/sirupsen/logrus"
)

type flags struct {
	URL            string
	DumpHttp       bool
	LogLevel       string
	VerboseDecoder bool
	AnsiArt        int
	Threshold      int
	Flicker        bool
	FastStart      int
	FastResume     bool
	DumpFSM        bool
	OneShot        bool
	Worker         bool
	Privsep        bool
	WorkerMemQuota int
}

func WithFlags() StreamOption {
	return func(s *Stream) error {
		s.flags = someFlags

		return nil
	}
}

func getFlags() *flags {
	f := flags{}
	flag.BoolVar(&f.DumpHttp, "dump-http", false, "dumps http headers")
	flag.BoolVar(&f.VerboseDecoder, "verbose", false, "ffmpeg debuggging info")
	flag.IntVar(&f.AnsiArt, "ansi-art", 0, "output ansi art on modulo frame")
	flag.IntVar(&f.Threshold, "threshold", 8, "need this much to output a warning")
	flag.BoolVar(&f.Flicker, "flicker", false, "reset terminal in ansi mode")