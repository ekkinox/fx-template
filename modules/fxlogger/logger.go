package fxlogger

import (
	"github.com/rs/zerolog"
)

type Logger struct {
	*zerolog.Logger
}

func FromLogger(logger zerolog.Logger) *Logger {
	return &Logger{&logger}
}
