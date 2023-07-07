package fxlogger

import (
	"io"
	"os"

	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

var FxLoggerModule = fx.Module("logger",
	fx.Provide(
		NewDefaultLoggerFactory,
		NewFxLogger,
		fx.Annotate(
			GetTestLogBufferInstance,
			fx.ResultTags(`name:"test-log-buffer"`),
		),
	),
)

type FxLoggerParam struct {
	fx.In
	Factory       LoggerFactory
	Config        *fxconfig.Config
	TestLogBuffer *TestLogBuffer `name:"test-log-buffer"`
}

func NewFxLogger(p FxLoggerParam) (*Logger, error) {

	// level
	level := FetchLogLevel(p.Config.GetString("modules.logger.level"))
	if p.Config.AppDebug() {
		level = zerolog.DebugLevel
	}

	// output writer
	var outputWriter io.Writer
	if p.Config.AppEnv() == fxconfig.Test {
		outputWriter = p.TestLogBuffer
	} else {
		outputWriter = os.Stdout
	}

	// logger
	return p.Factory.Create(
		WithName(p.Config.AppName()),
		WithLevel(level),
		WithOutputWriter(outputWriter),
	)
}
