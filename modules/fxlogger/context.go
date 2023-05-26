package fxlogger

import (
	"context"
	"github.com/rs/zerolog"
)

func CtxLogger(ctx context.Context) *Logger {
	return &Logger{zerolog.Ctx(ctx)}
}
