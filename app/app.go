package app

import (
	"github.com/ekkinox/fx-template/app/handlers"
	"github.com/ekkinox/fx-template/app/services"
	"github.com/ekkinox/fx-template/modules/fxhttpserver"

	"go.uber.org/fx"
)

var AppModule = fx.Module("app",
	fx.Provide(
		//routes
		fxhttpserver.AsHttpServerHandler(handlers.NewHelloHandler),
		//services
		services.NewTestService,
	),
)
