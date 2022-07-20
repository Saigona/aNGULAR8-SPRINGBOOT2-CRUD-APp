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

ty