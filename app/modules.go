package app

import (
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxhttpserver"
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"github.com/ekkinox/fx-template/modules/fxtracer"

	"go.uber.org/fx"
)

func RegisterModules() fx.Option {
	return fx.Options(
		fxconfig.FxConfigModule,
		fxlogger.FxLoggerModule,
		fxtracer.FxTracerModule,
		fxhttpserver.FxHttpServerModule,
	)
}
