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

func NewFxLogger(p FxLoggerParam) *Logger {

	lvl := log.INFO
	if p.Config.AppConfig.Debug {
		lvl = log.DEBUG
	}

	return NewLogger(
		os.Stdout,
		WithField("service", p.Config.AppConfig.Name),
		WithLevel(lvl),
	)
}
