package worker

import (
	"context"

	"github.com/WIZARDISHUNGRY/hls-await/internal/logger"
	"github.com/WIZARDISHUNGRY/hls-await/internal/segment"
)

type InProcess struct {
}

// startWorker runs in the child process
func (ip *InProcess) Start(ctx context.Context) error {
	return nil
}

func