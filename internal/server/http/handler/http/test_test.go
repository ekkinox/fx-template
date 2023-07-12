package http_test

import (
	"net/http/httptest"
	"testing"

	"github.com/ekkinox/fx-template/internal/server"
	"github.com/ekkinox/fx-template/modules/fxhttpserver"
	"github.com/ekkinox/fx-template/modules/fxhttpserver/fxhttpservertest"
	"github.com/ekkinox/fx-template/modules/fxlogger/fxloggertest"
	"github.com/ekkinox/fx-template/modules/fxtracer/fxtracertest"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/fx"
)

func TestHttpEndpoint(t *testing.T) {

	// test env vars
	t.Setenv("APP_CONFIG_PATH", "../../../../../configs")
	t.Setenv("PUBSUB_PROJECT_ID", "test-project")
	t.Setenv("PUBSUB_EMULATOR_HOST", "fake")

	// preparation
	var httpServer *echo.Echo
	server.ServerBoostrapper.BoostrapAndRunTestApp(
		t,
		fxhttpserver.StartFxHttpServer(),
		fx.Populate(&httpServer),
	)

	// http request
	req := httptest.NewRequest("GET", "/test", nil)
	rec := fxhttpservertest.RecordResponse(httpServer, req)

	// http response assertion
	fxhttpservertest.AssertRecordedResponseCode(t, rec, 200)
	fxhttpservertest.AssertRecordedResponseBody(t, rec, `{"value":"testtest"}`)

	// log assertion
	fxloggertest.AssertHasLogRecord(t, map[string]interface{}{
		"level":   "info",
		"message": "test log",
	})

	// trace assertion
	fxtracertest.AssertHasTraceSpan(t, "test span", attribute.String("test attribute name", "test attribute value"))
}
