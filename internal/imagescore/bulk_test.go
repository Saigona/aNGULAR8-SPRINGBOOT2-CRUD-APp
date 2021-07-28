package imagescore

import (
	"context"
	"image"
	"sync"
	"testing"
)

//go:generate sh -c "go test ./... -run '^$' -benchmem -bench . | tee benchresult.txt"
//go:generate sh -c "git show :./benchresult.txt | go run golang.org/x/perf/cmd/benchstat -delta-test none -geomean /dev/s