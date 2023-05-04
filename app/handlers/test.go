package handlers

import (
	"fmt"
	"github.com/ekkinox/fx-template/app/services"
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"github.com/labstack/echo/v4"
	"net/http"
)

type TestHandler struct {
	logger  *fxlogger.Logger
	service *services.TestService
}

func NewHelloHandler(logger *fxlogger.Logger, service *services.TestService) *TestHandler {
	return &TestHandler{
		logger:  logger,
		service: service,
	}
}

func (*TestHandler) Method() string {
	return "GET"
}

func (*TestHandler) Path() string {
	return "/test/:name"
}

func (h *TestHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {

		name := c.Param("name")

		h.logger.Info().Msgf("called %s with name %s", h.Path(), name)

		return c.String(
			http.StatusOK,
			fmt.Sprintf("Test hello world for %s. Service output: %s.", name, h.service.Test()),
		)
	}
}