package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Test3Middleware struct{}

func NewTest3Middleware() *Test3Middleware {
	return &Test3Middleware{}
}

func (m *Test3Middleware) Handle() echo.MiddlewareFunc {
	return middleware.CORS()
}
