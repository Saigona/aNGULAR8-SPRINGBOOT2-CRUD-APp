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

func (p *Par