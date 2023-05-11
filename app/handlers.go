package app

import (
	"github.com/ekkinox/fx-template/app/handlers"
	"github.com/ekkinox/fx-template/modules/fxhttpserver"
	"go.uber.org/fx"
)

func RegisterHandlers() fx.Option {
	return fx.Provide(
		fxhttpserver.AsHttpServerHandler(handlers.NewHelloHandler),
	)
}
