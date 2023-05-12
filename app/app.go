package app

import (
	"github.com/ekkinox/fx-template/modules/fxlogger"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

var App = fx.New(
	// resources
	RegisterModules(),
	RegisterHandlers(),
	RegisterServices(),
	// fx logger
	fx.WithLogger(func(log *fxlogger.Logger) fxevent.Logger {
		return log
	}),
)
