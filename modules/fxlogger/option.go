package fxlogger

import (
	"github.com/rs/zerolog"
	"io"
	"os"
)

type options struct {
	Name         string
	Level        zerolog.Level
	OutputWriter io.Writer
}

var defaultOptions = options{
	Name:         "default",
	Level:        zerolog.InfoLevel,
	OutputWriter: os.Stdout,
}

type Option func(o *options)

func WithName(n string) Option {
	return func(o *options) {
		o.Name = n
	}
}

func WithLevel(l zerolog.Level) Option {
	return func(o *options) {
		o.Level = l
	}
}

func WithOutputWriter(w io.Writer) Option {
	return func(o *options) {
		o.OutputWriter = w
	}
}
