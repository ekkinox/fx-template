package fxlogger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct {
	*zerolog.Logger
}

func NewLogger(opts ...Option) *Logger {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	appliedOpts := defaultOptions
	for _, applyOpt := range opts {
		applyOpt(&appliedOpts)
	}

	logger := log.
		Output(appliedOpts.OutputWriter).
		With().
		Str("service", appliedOpts.Name).
		Logger().
		Level(appliedOpts.Level)

	return &Logger{&logger}
}

func FromLogger(logger zerolog.Logger) *Logger {
	return &Logger{&logger}
}
