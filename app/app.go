package app

import (
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"go.uber.org/fx"
)

var App = fx.New(
	RegisterModules(),
	RegisterHandlers(),
	RegisterServices(),
	RegisterOverrides(),
	fx.WithLogger(fxlogger.FxEventLogger),
)
