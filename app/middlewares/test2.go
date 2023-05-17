package middlewares

import (
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/labstack/echo/v4"
)

type Test2Middleware struct {
	config *fxconfig.Config
}

func NewTest2Middleware(config *fxconfig.Config) *Test2Middleware {
	return &Test2Middleware{
		config: config,
	}
}

func (m *Test2Middleware) Handle() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Logger().Infof("middleware 2 - %s", m.config.AppConfig.Name)

			return next(c)
		}
	}
}
