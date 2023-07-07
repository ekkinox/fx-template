package fxlogger

import (
	"github.com/rs/zerolog"
)

const (
	Service = "service"
)

type Logger struct {
	zerolog.Logger
}

func (l *Logger) ToZerolog() *zerolog.Logger {
	return &l.Logger
}

func FromZerolog(logger zerolog.Logger) *Logger {
	return &Logger{logger}
}
