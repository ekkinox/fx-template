package middleware

import (
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/labstack/echo/v4"
)

type TestMiddleware struct {
	config *fxconfig.Config
}

func NewTestMiddleware(config *fxconfig.Config) *TestMiddleware {
	return &TestMiddleware{
		config: config,
	}
}

func (m *TestMiddleware) Handle() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Logger().Infof("test middleware for app: %s", m.config.AppName())

			return next(c)
		}
	}
}
