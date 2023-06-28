package middleware

import (
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/labstack/echo/v4"
)

type GroupMiddleware struct {
	config *fxconfig.Config
}

func NewGroupMiddleware(config *fxconfig.Config) *GroupMiddleware {
	return &GroupMiddleware{
		config: config,
	}
}

func (m *GroupMiddleware) Handle() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Logger().Infof("GROUP middleware for app: %s", m.config.AppName())

			return next(c)
		}
	}
}
