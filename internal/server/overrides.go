package server

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

func RegisterOverrides() fx.Option {
	return fx.Options(
		// test echo decoration
		fx.Decorate(
			func(e *echo.Echo) *echo.Echo {
				e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
					return func(c echo.Context) error {
						c.Logger().Info("override example middleware")
						return next(c)
					}
				})

				return e
			},
		),
	)
}
