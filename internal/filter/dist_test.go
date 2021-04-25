package filter

import (
	"context"
	"image"
	"testing"

	"github.com/WIZARDISHUNGRY/hls-await/internal/corpus"
	"github.com/WIZARDISHUNGRY/hls-await/internal/logger"
	"github.com/sirupsen/logrus"
)

//go:generate sh -c "go test ./... -run '^$' -benchmem -bench . | tee benchresult.txt"
//go:generate sh -c "git show :./benchresult.txt | go run golang.org/x/perf/cmd/benchstat -delta-test none -geomean /dev/stdin benchresult.tx