package handler

import (
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxhttpserver"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CallHandler struct {
	config *fxconfig.Config
}

func NewCallHandler(config *fxconfig.Config) *CallHandler {
	return &CallHandler{
		config: config,
	}
}

func (h *CallHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {

		client := fxhttpserver.CtxHttpClient(c)

		res, err := client.Get(h.config.GetString("config.call.url"))
		if err != nil {
			c.Logger().Errorf("cannot call target: %v", err)
			return err
		}

		body, err := io.ReadAll(res.Body)
		err = res.Body.Close()
		if err != nil {
			c.Logger().Errorf("cannot close response body: %v", err)
			return err
		}

		return c.String(http.StatusOK, string(body))
	}
}
