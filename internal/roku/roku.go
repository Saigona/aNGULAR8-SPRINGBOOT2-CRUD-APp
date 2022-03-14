package roku

import (
	"context"
	"errors"
	"net/url"
	"sync"
	"time"

	"github.com/WIZARDISHUNGRY/hls-await/internal/logger"
	"golang.org/x/sync/errgroup"
	"jonwillia.ms/roku"
)

func Run(ctx context.Context) func() (*roku.Remote, error) {
	log := logger.Entry(ctx)

	const dur = time.Minute
	var (
		mutex  sync.Mutex
		remote *roku.Remote
		errC   = make(chan error)
		timer  = time.NewTimer(time.Minute)
	)

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
	LOOP:
		for ctx.Err() == nil {

			devs, err := roku.FindRokuDevices()
			switch {
			case len(devs) == 0:
				fallthrough
			case err != nil:
				log.WithError(err).Warn("roku.FindRokuDevices")
				time.Sleep(10 * time.Second) // TODO not abortable
				continue LOOP
			}
			dev := devs[0]
			log.Infof("found roku %s : %s", dev.Addr, dev.Name)
			r, err := roku.NewRemote(dev.Addr)
			if err != nil {
				log.WithError(err).Warn("roku.NewRemote")
				time.Sleep(10 * time.Second) // TODO not abortable
				continue LOOP
			}
			mutex.Lock()
			remo