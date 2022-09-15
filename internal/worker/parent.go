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
	"github.com/WIZARDISHUNGRY/hls-await/internal/seg