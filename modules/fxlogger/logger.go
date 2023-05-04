package fxlogger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type options struct {
	OutputWriter io.Writer
	Level        zerolog.Level
}

var defaultOptions = options{
	OutputWriter: os.Stdout,
	Level:        zerolog.InfoLevel,
}

type LoggerOption func(o *options)

func WithOutputWriter(w io.Writer) LoggerOption {
	return func(o *options) {
		o.OutputWriter = w
	}
}

func WithLevel(l zerolog.Level) LoggerOption {
	return func(o *options) {
		o.Level = l
	}
}

type Logger struct {
	*zerolog.Logger
}

func newLogger(appName string, opts ...LoggerOption) *Logger {

	appliedOpts := defaultOptions

	for _, applyOpt := range opts {
		applyOpt(&appliedOpts)
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	l := log.
		Output(appliedOpts.OutputWriter).With().
		Str("service", appName).
		Logger().
		Level(appliedOpts.Level)

	return &Logger{&l}
}
