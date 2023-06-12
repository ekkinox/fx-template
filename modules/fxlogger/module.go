package fxlogger

import (
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
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

	level := zerolog.InfoLevel
	if p.Config.AppDebug() {
		level = zerolog.DebugLevel
	}

	return NewLogger(
		WithName(p.Config.AppName()),
		WithLevel(level),
	)
}
