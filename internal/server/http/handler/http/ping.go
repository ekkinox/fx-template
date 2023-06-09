package http

import (
	"io"
	"net/http"

	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxhttpserver"

	"github.com/labstack/echo/v4"
)

type PingHandler struct {
	config *fxconfig.Config
}

func NewPingHandler(config *fxconfig.Config) *PingHandler {
	return &PingHandler{
		config: config,
	}
}

func (h *PingHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {

		client := fxhttpserver.CtxHttpClient(c)

		req, err := http.NewRequest("GET", h.config.GetString("config.ping.url"), nil)
		if err != nil {
			c.Logger().Errorf("cannot create request: %v", err)
		}

		req.Header.Set("custom", "value")

		res, err := client.Do(req)
		if err != nil {
			c.Logger().Errorf("cannot request target: %v", err)
			return err
		}

		body, err := io.ReadAll(res.Body)
		err = res.Body.Close()
		if err != nil {
			c.Logger().Errorf("cannot close response body: %v", err)
			return err
		}

		return c.String(res.StatusCode, string(body))
	}
}
