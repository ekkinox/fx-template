package fxlogger

import (
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"go.uber.org/fx"
)

var LoggerModule = fx.Module("logger",
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
