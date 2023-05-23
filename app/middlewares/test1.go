package middlewares

import (
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/labstack/echo/v4"
)

type Test1Middleware struct {
	config *fxconfig.Config
}

func NewTest1Middleware(config *fxconfig.Config) *Test1Middleware {
	return &Test1Middleware{
		config: config,
	}
}

func (m *Test1Middleware) Handle() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Logger().Infof("middleware 1 - %s", m.config.AppName())

			return next(c)
		}
	}
}
