package fxlogger

import (
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/labstack/gommon/log"
	"go.uber.org/fx"
	"os"
)

var FxLoggerModule = fx.Module("logger",
	fx.Provide(
		NewFxLogger,
	),
)

type FxLoggerParam struct {
	fx.In
	Config *fxconfig.Config
}

type FxLoggerResult struct {
	fx.Out
	Logger *Logger
}

func NewFxLogger(p FxLoggerParam) FxLoggerResult {

	lvl := log.INFO
	if p.Config.AppConfig.Debug {
		lvl = log.DEBUG
	}

	return FxLoggerResult{
		Logger: NewLogger(
			os.Stdout,
			WithField("service", p.Config.AppConfig.Name),
			WithLevel(lvl),
		),
	}
}
