package worker

import (
	"context"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/WIZARDISHUNGRY/hls-await/internal/logger"
	"github.com/WIZARDISHUNGRY/hls-await/internal/segment"
	"github.com/WIZARDISHUNGRY/hls-await/pkg/unixmsg"
	"github.com/pkg/errors"
)

type Parent struct {
	once         sync.Once
	mutex        sync.RWMutex
	cmd          *exec.Cmd
	listener     *net.UnixListener
	client       *rpc.Client
	conn, connFD *net.UnixConn
	launchCount  int
	lastLaunch   time.Time
	context      context.Context
}

func (p *Parent) Start(ctx context.Context) error {

	var retErr error
	p.once.Do(func() {
		retErr = p.spawnChild(ctx)
		if retErr == nil {
			go p.loop(ctx)
		}
	})
	return retErr
}

func (p *Parent) closeChild(ctx context.Context) error {
	// PRE: must own write mutex
	if p.client != nil {
		p.client.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
	if p.connFD != nil {
		p.connFD.Close()
	}
	if p.listener != nil {
		p.listener.Close()
	}

	p.nicelyKill(ctx, p.cmd) // kick the process

	return nil
}

func (p *Parent) Restart(ctx context.Context) {
	p.mutex.RLock()
	cmd := p.cmd
	p.mutex.RUnlock()
	p.nicelyKill(ctx, cmd)
}

func (p *Parent) nicelyKill(ctx context.Context, cmd *exec.Cmd) {
	log := logger.Entry(ctx)
	// PRE: must own write mutex
	if cmd != nil && cmd.Process != nil {
		log.Info("Signaling child to exit")
		cmd.Process.Signal(syscall.SIGTERM)
		time.Sleep(3 * time.Second)
		cmd.Process.Kill()
	}
}

func (p *Parent) spawnChild(ctx context.Context) (err error) {
	log := logger.Entry(ctx)

	defer func() {
		l := log.WithField("count", p.launchCount)
		if !p.lastLaunch.IsZero() {
			l = l.WithField("lifetime", time.Now().Sub(p.lastLaunch))
		}
		p.lastLaunch = time.Now()
		l.WithError(err).Infof("spawnChild")
	}()
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.closeChi