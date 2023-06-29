package fxlogger

import (
	"os"

	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

var FxLoggerModule = fx.Module("logger",
	fx.Provide(
		NewDefaultLoggerFactory,
		NewFxLogger,
	),
)

type FxLoggerParam struct {
	fx.In
	Config  *fxconfig.Config
	Factory LoggerFactory
}

func NewFxLogger(p FxLoggerParam) (*Logger, error) {

	level := zerolog.InfoLevel
	if p.Config.AppDebug() {
		level = zerolog.DebugLevel
	}

	return p.Factory.Create(
		WithName(p.Config.AppName()),
		WithLevel(level),
		WithOutputWriter(os.Stdout),
	)
}
