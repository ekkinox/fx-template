package http_test

import (
	"context"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ekkinox/fx-template/internal/server"
	h "github.com/ekkinox/fx-template/internal/server/http"
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxhealthchecker"
	"github.com/ekkinox/fx-template/modules/fxhttpserver"
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"github.com/ekkinox/fx-template/modules/fxlogger/fxloggertest"
	"github.com/ekkinox/fx-template/modules/fxtracer"
	"github.com/ekkinox/fx-template/modules/fxtracer/fxtracertest"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestHttpEndpoint(t *testing.T) {

	ctx := context.Background()
	var httpServer *echo.Echo

	os.Setenv("APP_ENV", "test")
	os.Setenv("APP_CONFIG_PATH", "../../../../../configs")
	os.Setenv("PUBSUB_PROJECT_ID", "test-project")
	os.Setenv("PUBSUB_EMULATOR_HOST", "pubsub:8085")

	fxtest.New(
		t,
		fx.NopLogger,
		//fx.WithLogger(fxlogger.FxEventLogger),
		// core
		fxconfig.FxConfigModule,
		fxlogger.FxLoggerModule,
		fxtracer.FxTracerModule,
		fxhealthchecker.FxHealthCheckerModule,
		// common
		server.RegisterModules(ctx),
		server.RegisterServices(ctx),
		server.RegisterOverrides(ctx),
		// http server
		fxhttpserver.FxHttpServerModule,
		h.RegisterHandlers(),
		//options
		fx.Options(
			fx.Populate(&httpServer),
		),
	).RequireStart().RequireStop()

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	// response assertion
	assert.Equal(t, 200, rec.Code)
	assert.Contains(t, rec.Body.String(), "test")

	//log assertion
	fxloggertest.AssertHasLogRecord(t, map[string]interface{}{
		"level":   "info",
		"message": "test log",
	})

	//trace assertion
	fxtracertest.AssertHasTraceSpan(t, "test span", attribute.String("test attribute name", "test attribute value"))
}
