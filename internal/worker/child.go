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
	"github.com/WIZARDISHUNGRY/hls-await/inte