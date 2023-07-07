package http_test

import (
	"context"
	"net/http/httptest"
	"os"

	"github.com/ekkinox/fx-template/internal/server"
	h "github.com/ekkinox/fx-template/internal/server/http"
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxhealthchecker"
	"github.com/ekkinox/fx-template/modules/fxhttpserver"
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"github.com/ekkinox/fx-template/modules/fxtracer"
	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

var _ = Describe("some http test", func() {

	ctx := context.Background()
	var httpServer *echo.Echo

	BeforeEach(func() {
		os.Setenv("APP_ENV", "test")
		os.Setenv("APP_CONFIG_PATH", "../../../../../configs")
		os.Setenv("PUBSUB_PROJECT_ID", "test-project")
		os.Setenv("PUBSUB_EMULATOR_HOST", "pubsub:8085")

		fxtest.New(
			GinkgoT(),
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
	})

	Describe("test /test", func() {
		Context("calling /test", func() {
			It("response should contain the config value", func() {

				// http call
				req := httptest.NewRequest("GET", "/test", nil)
				rec := httptest.NewRecorder()
				httpServer.ServeHTTP(rec, req)

				// response assertion
				Expect(200).To(Equal(rec.Code))
				Expect(rec.Body.String()).To(ContainSubstring("test"))

				// log assertion
				buf := fxlogger.GetTestLogBufferInstance()
				Expect(buf.HasRecord("info", "in test endpoint")).To(Equal(true))
			})
		})
	})
})
