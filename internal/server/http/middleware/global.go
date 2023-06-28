package middleware

import (
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/labstack/echo/v4"
)

type GlobalMiddleware struct {
	config *fxconfig.Config
}

func NewGlobalMiddleware(config *fxconfig.Config) *GlobalMiddleware {
	return &GlobalMiddleware{
		config: config,
	}
}

func (m *GlobalMiddleware) Handle() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Logger().Infof("GLOBAL middleware for app: %s", m.config.AppName())

			return next(c)
		}
	}
}
