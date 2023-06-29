package fxlogger

import (
	"github.com/rs/zerolog"
)

const (
	Service = "service"
)

type Logger struct {
	*zerolog.Logger
}

func FromLogger(logger zerolog.Logger) *Logger {
	return &Logger{&logger}
}
