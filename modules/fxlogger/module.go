package fxlogger

import (
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

var LoggerModule = fx.Module("logger",
	fx.WithLogger(func(logger *Logger) fxevent.Logger {
		return &FxEventLogger{Logger: logger}
	}),
	fx.Provide(
		NewLogger,
	),
)

type LoggerParam struct {
	fx.In
	Config *fxconfig.Config
}

type LoggerResult struct {
	fx.Out
	Logger *Logger
}

func NewLogger(p LoggerParam) LoggerResult {
	return LoggerResult{
		Logger: newLogger(p.Config.AppName),
	}
}
