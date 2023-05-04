package handler

import (
	"fmt"
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"github.com/labstack/echo/v4"
	"net/http"
)

type HelloHandler struct {
	config *fxconfig.Config
	logger *fxlogger.Logger
}

func NewHelloHandler(config *fxconfig.Config, logger *fxlogger.Logger) *HelloHandler {
	return &HelloHandler{
		config: config,
		logger: logger,
	}
}

func (*HelloHandler) Method() string {
	return "GET"
}

func (*HelloHandler) Path() string {
	return "/hello/:name"
}

func (h *HelloHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {

		return c.String(http.StatusOK, fmt.Sprintf("Hello %s world from %s.", c.Param("name"), h.config.AppName))
	}
}
