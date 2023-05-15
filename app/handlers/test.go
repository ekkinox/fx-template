package handlers

import (
	"fmt"
	"github.com/ekkinox/fx-template/modules/fxhttpserver"
	"net/http"

	"github.com/ekkinox/fx-template/app/services"
	"github.com/ekkinox/fx-template/modules/fxlogger"

	"github.com/labstack/echo/v4"
)

type TestHandler struct {
	service *services.TestService
}

func NewHelloHandler(logger *fxlogger.Logger, service *services.TestService) *TestHandler {
	return &TestHandler{
		service: service,
	}
}

func (*TestHandler) Method() string {
	return "GET"
}

func (*TestHandler) Path() string {
	return "/test/:name"
}

func (h *TestHandler) Handler() echo.HandlerFunc {
	return func(c echo.Context) error {

		name := c.Param("name")

		c.Logger().Infof("called %s with name=%s", c.Path(), name)
		fxhttpserver.GetCtxLogger(c).Info().Msg("fxhttpserver.GetLogger")

		_, span := fxhttpserver.GetCtxTracer(c).Start(c.Request().Context(), "some span")
		defer span.End()

		output, err := h.service.Test(c)
		if err != nil {
			c.Logger().Errorf("handler failed=%s", err)
			return err
		}

		return c.String(
			http.StatusOK,
			fmt.Sprintf("Service output: %s.", output),
		)
	}
}

func (h *TestHandler) Middlewares() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				c.Logger().Info("custom handler middleware adding custom response header")
				c.Response().Header().Set("x-custom-header", "/test/:name")
				return next(c)
			}
		},
	}
}
