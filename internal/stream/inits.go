package stream

import (
	"github.com/WIZARDISHUNGRY/hls-await/internal/bot"
	"github.com/WIZARDISHUNGRY/hls-await/internal/worker"
	"jonwillia.ms/roku"
)

func InitWorker() worker.Worker {
	if someFlags.Worker {
		return &worker.Child{
			MemQuota: someFlags.WorkerMemQuota,
		}
	}
	if !someFlags.Privsep {
		return &worker.InProcess{}
	}
	return &worker.Parent{}
}

func WithWorker(w worker.Worker) StreamOption {
	return func(s *Stream) error {
		s.worker = w
		return nil
	}
}

func WithBot(b 