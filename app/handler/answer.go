package handler

import (
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AnswerHandler struct {
	config *fxconfig.Config
}

func NewAnswerHandler(config *fxconfig.Config) *AnswerHandler {
	return &AnswerHandler{
		config: config,
	}
}

func (h *AnswerHandler) Handle() echo.HandlerFunc {

	status := http.StatusOK
	if h.config.GetBool("config.answer.should_fail") {
		status = http.StatusInternalServerError
	}

	return func(c echo.Context) error {
		return c.JSON(status, c.Request().Header)
	}
}
