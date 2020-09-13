
package main

import (
	"context"
	"flag"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"github.com/WIZARDISHUNGRY/hls-await/internal/bot"
	"github.com/WIZARDISHUNGRY/hls-await/internal/logger"
	"github.com/WIZARDISHUNGRY/hls-await/internal/roku"
	"github.com/WIZARDISHUNGRY/hls-await/internal/stream"
	"github.com/WIZARDISHUNGRY/hls-await/internal/worker"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

const streamURL = "https://tv.nknews.org/tvhls/stream.m3u8"

var (
	currentStream *stream.Stream
)

var (
	_, b, _, _ = runtime.Caller(0)

	// root folder of this project for trimming frames
	root = filepath.Join(filepath.Dir(b), "../..")
)

func main() {
	logr := logrus.New()
	logr.ReportCaller = true
	logr.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			lastIdx := strings.LastIndexByte(f.Function, '/')
			if lastIdx == -1 {
				lastIdx = 0
			} else {
				lastIdx++
			}
			fxn := f.Function[lastIdx:]
			lastIdx = strings.IndexByte(fxn, '(')
			if lastIdx == -1 {
				lastIdx = 0
			}
			fxn = fxn[lastIdx:]
			return fxn,
				""
			//	f.File[len(root)+1:]
		},
	}
	log := logr.WithFields(nil)
	flag.Parse()

	args := flag.Args()