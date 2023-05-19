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

func RegisterHandlers() fx.Option {
	return fx.Options(
		fxhttpserver.RegisterHandler(
			"GET",
			"/bar",
			func(c echo.Context) error {
				c.Logger().Info("from bar")
				return c.JSON(http.StatusOK, "ok")
			},
			middleware.CORS(),
			middlewares.NewTest2Middleware,
		),
		fxhttpserver.RegisterHandler(
			"GET",
			"/foo",
			handlers.NewFooHandler,
			middlewares.NewTest2Middleware,
			middleware.Gzip(),
			middlewares.NewTest1Middleware,
		),
	)
}
