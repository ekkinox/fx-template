package app

import (
	"github.com/ekkinox/fx-template/modules/fxhttpserver"
	"go.uber.org/fx"
)

func RegisterModules() fx.Option {
	return fx.Options(
		fxhttpserver.FxHttpServerModule,
	)
}
