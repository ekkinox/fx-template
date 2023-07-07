package http

import (
	"net/http"

	"github.com/ekkinox/fx-template/modules/fxconfig"

	"github.com/labstack/echo/v4"
)

type TestHandler struct {
	config *fxconfig.Config
}

func NewTestHandler(config *fxconfig.Config) *TestHandler {
	return &TestHandler{
		config: config,
	}
}

func (h *TestHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {

		c.Logger().Info("in test endpoint")

		return c.JSON(http.StatusOK, echo.Map{
			"value": h.config.GetString("config.test.value"),
		})
	}
}
