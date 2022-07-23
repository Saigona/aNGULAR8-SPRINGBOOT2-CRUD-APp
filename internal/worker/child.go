package worker

import (
	"context"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/WIZARDISHUNGRY/hls-await/internal/logger"
	"github.com/WIZARDISHUNGRY/hls-await/internal/segment"
	"github.com/WIZARDISHUNGRY/hls-await/pkg/unixmsg"
	"github.com/pkg/errors"
)

const (
	durWaitBeforeStopTheWorld = 2 * time.Second
	maxConsecutivePanics      = 2
)

type Child struct {
	once      sync.Once
	memstatsC chan error
	MemQuota  int
}

func (c *Child) Start(ctx context.Context) error {
	var retErr error
	c.once.Do(func() { // This should block and then error out
		retErr = c.runWorker(ctx)
	})
	ret