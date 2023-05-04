package main

import (
	"github.com/ekkinox/fx-template/app"
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxhttpserver"
	"github.com/ekkinox/fx-template/modules/fxlogger"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		// core modules
		fxconfig.ConfigModule,
		fxlogger.LoggerModule,
		fxhttpserver.HttpServerModule,
		// app module
		app.AppModule,
	).Run()
}
