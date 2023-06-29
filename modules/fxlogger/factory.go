package fxlogger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type LoggerFactory interface {
	Create(options ...LoggerOption) (*Logger, error)
}

type DefaultLoggerFactory struct{}

func NewDefaultLoggerFactory() LoggerFactory {
	return &DefaultLoggerFactory{}
}

func (f *DefaultLoggerFactory) Create(options ...LoggerOption) (*Logger, error) {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	appliedOpts := defaultLoggerOptions
	for _, applyOpt := range options {
		applyOpt(&appliedOpts)
	}

	logger := log.
		Output(appliedOpts.OutputWriter).
		With().
		Str("service", appliedOpts.Name).
		Logger().
		Level(appliedOpts.Level)

	return &Logger{&logger}, nil
}
