package logger

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

type ctxKey int

const (
	ctxKeyLog = iota
)

func Entry(ctx context.Context) *logrus.Entry {
	