package app

import (
	"go.uber.org/fx"
)

var AppModule = fx.Module(
	"app",
	RegisterHandlers(),
	RegisterServices(),
)
