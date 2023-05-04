package main

import (
	"github.com/ekkinox/fx-template/handler"
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxhttpserver"
	"github.com/ekkinox/fx-template/modules/fxlogger"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func main() {
	fx.New(
		//modules
		fxconfig.ConfigModule,
		fxlogger.LoggerModule,
		fxhttpserver.HttpServerModule,
		//providers
		fx.Provide(
			fxhttpserver.AsHttpServerHandler(handler.NewHelloHandler),
		),
		//logger
		fx.WithLogger(func(logger *fxlogger.Logger) fxevent.Logger {
			return &fxlogger.FxEventLogger{Logger: logger}
		}),
	).Run()
}
