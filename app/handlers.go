package app

import (
	"github.com/ekkinox/fx-template/app/handlers"
	"github.com/ekkinox/fx-template/app/middlewares"
	"github.com/ekkinox/fx-template/modules/fxhttpserver"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
	"net/http"
)

var adminHandlersGroup = fxhttpserver.NewHandlersGroupRegistration(
	"/admin",
	[]*fxhttpserver.HandlerRegistration{
		fxhttpserver.NewHandlerRegistration("GET", "/foo", handlers.NewFooHandler),
		fxhttpserver.NewHandlerRegistration("GET", "/bar", handlers.NewBarHandler),
	},
	middlewares.NewTest1Middleware,
	middlewares.NewTest2Middleware,
	middleware.CORS(),
	middleware.Gzip(),
)
var fooHandler = fxhttpserver.NewHandlerRegistration("GET", "/foo", handlers.NewFooHandler, middlewares.NewTest1Middleware)
var barHandler = fxhttpserver.NewHandlerRegistration("GET", "/bar", handlers.NewBarHandler, middlewares.NewTest2Middleware)

func RegisterHandlers() fx.Option {
	return fx.Options(
		fxhttpserver.RegisterHandlersGroup(adminHandlersGroup),
		fxhttpserver.RegisterHandler(fooHandler),
		fxhttpserver.RegisterHandler(barHandler),
		fxhttpserver.AsHandler("GET", "/baz", func(c echo.Context) error {
			c.Logger().Info("from baz")
			return c.JSON(http.StatusOK, "ok")
		}),
	)
}
