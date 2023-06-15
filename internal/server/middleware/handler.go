package middleware

import (
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/labstack/echo/v4"
)

type HandlerMiddleware struct {
	config *fxconfig.Config
}

func NewHandlerMiddleware(config *fxconfig.Config) *HandlerMiddleware {
	return &HandlerMiddleware{
		config: config,
	}
}

func (m *HandlerMiddleware) Handle() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Logger().Infof("HANDLER middleware for app: %s", m.config.AppName())

			return next(c)
		}
	}
}
