package http

import (
	"net/http"

	"github.com/ekkinox/fx-template/modules/fxconfig"

	"github.com/labstack/echo/v4"
)

type PongHandler struct {
	config *fxconfig.Config
}

func NewPongHandler(config *fxconfig.Config) *PongHandler {
	return &PongHandler{
		config: config,
	}
}

func (h *PongHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {

		if h.config.GetBool("config.pong.should_fail") {
			return echo.NewHTTPError(http.StatusInternalServerError, "pong configured to fail")
		}

		return c.JSON(http.StatusOK, c.Request().Header)
	}
}
