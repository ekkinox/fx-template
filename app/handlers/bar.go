package handlers

import (
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"net/http"

	"github.com/labstack/echo/v4"
)

type BarHandler struct {
	config *fxconfig.Config
}

func NewBarHandler(config *fxconfig.Config) (*BarHandler, error) {
	return &BarHandler{
		config: config,
	}, nil
}

func (h *BarHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {

		c.Logger().Infof("called %s", c.Path())

		return c.JSON(
			http.StatusOK,
			h.config.GetStringMap("config"),
		)
	}
}
