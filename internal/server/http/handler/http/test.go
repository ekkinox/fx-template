package http

import (
	"net/http"

	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxhttpserver"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

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

		c.Logger().Info("test log")

		_, span := fxhttpserver.CtxTracer(c).Start(
			c.Request().Context(),
			"test span",
			trace.WithAttributes(attribute.String("test attribute name", "test attribute value")),
		)
		defer span.End()

		return c.JSON(http.StatusOK, echo.Map{
			"value": h.config.GetString("config.test.value"),
		})
	}
}
