package stream

import (
	"github.com/WIZARDISHUNGRY/hls-await/internal/bot"
	"github.com/WIZARDISHUNGRY/hls-await/internal/worker"
	"jonwillia.ms/roku"
)

func InitWorker() worker.Worker {
	if someFlag