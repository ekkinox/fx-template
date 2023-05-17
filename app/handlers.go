package app

import (
	"github.com/ekkinox/fx-template/app/handlers"
	"github.com/ekkinox/fx-template/app/middlewares"
	"github.com/ekkinox/fx-template/modules/fxhttpserver"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
)

func RegisterHandlers() fx.Option {
	return fx.Options(
		fxhttpserver.RegisterHandler("GET", "/foo", handlers.NewFooHandler, middlewares.NewTest1Middleware, middleware.CORS(), middleware.Gzip()),
		fxhttpserver.RegisterHandler("GET", "/bar", handlers.NewBarHandler, middlewares.NewTest2Middleware, middlewares.NewTest3Middleware),
	)
}
