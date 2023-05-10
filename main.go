package main

import (
	"github.com/ekkinox/fx-template/app"
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxhttpserver"
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"go.uber.org/fx/fxevent"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		// core modules
		fxconfig.FxConfigModule,
		fxlogger.FxLoggerModule,
		fxhttpserver.FxHttpServerModule,
		// app module
		app.AppModule,
		// logger
		fx.WithLogger(func(log *fxlogger.Logger) fxevent.Logger {
			return log
		}),
	).Run()
}
