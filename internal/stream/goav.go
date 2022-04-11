package stream

import (
	"context"
	"time"

	"github.com/WIZARDISHUNGRY/hls-await/internal/logger"
	"github.com/WIZARDISHUNGRY/hls-await/internal/segment"
	"github.com/sirupsen/logrus"
)

const workerMaxDuration = 10 * time.Second // if the worker appears to be stalled

func (s *Stream) ProcessSegment(ctx context.Context, request *segment.Request) error {
	log 