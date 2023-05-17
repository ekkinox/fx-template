package handlers

import (
	"fmt"
	"github.com/ekkinox/fx-template/modules/fxhttpserver"
	"net/http"

	"github.com/ekkinox/fx-template/app/services"
	"github.com/labstack/echo/v4"
)

type FooHandler struct {
	service *services.TestService
}

func NewFooHandler(service *services.TestService) *FooHandler {
	return &FooHandler{
		service: service,
	}
}

func (h *FooHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {

		// example logging
		c.Logger().Infof("called %s", c.Path())
		fxhttpserver.GetCtxLogger(c).Info().Msg("fxhttpserver.GetLogger")

		// example tracing
		_, span := fxhttpserver.GetCtxTracer(c).Start("some span")
		defer span.End()

		_, span2 := fxhttpserver.GetCtxTracer(c).Start("some span 2")
		defer span2.End()

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
