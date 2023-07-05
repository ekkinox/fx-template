package worker

import (
	"context"

	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxhealthchecker"
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"github.com/ekkinox/fx-template/modules/fxtracer"
	"go.uber.org/fx"
)

func BootstrapWorker(ctx context.Context) *fx.App {
	return fx.New(
		// logger
		fx.WithLogger(fxlogger.FxEventLogger),
		// core
		fxconfig.FxConfigModule,
		fxlogger.FxLoggerModule,
		fxtracer.FxTracerModule,
		fxhealthchecker.FxHealthCheckerModule,
		// worker
		RegisterModules(ctx),
		RegisterServices(ctx),
		RegisterOverrides(ctx),
	)
}
