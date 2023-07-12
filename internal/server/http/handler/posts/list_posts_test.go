package posts_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/ekkinox/fx-template/internal/model"
	"github.com/ekkinox/fx-template/internal/repository"
	"github.com/ekkinox/fx-template/internal/server"
	"github.com/ekkinox/fx-template/modules/fxgorm"
	"github.com/ekkinox/fx-template/modules/fxhttpserver"
	"github.com/ekkinox/fx-template/modules/fxhttpserver/fxhttpservertest"
	"github.com/ekkinox/fx-template/modules/fxlogger/fxloggertest"
	"github.com/ekkinox/fx-template/modules/fxtracer/fxtracertest"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/fx"
)

func TestListPostsHandler(t *testing.T) {

	// test env vars
	t.Setenv("APP_CONFIG_PATH", "../../../../../configs")
	t.Setenv("PUBSUB_PROJECT_ID", "test-project")
	t.Setenv("PUBSUB_EMULATOR_HOST", "fake")

	// preparation
	var httpServer *echo.Echo
	var repository *repository.PostRepository

	server.ServerBoostrapper.BoostrapAndRunTestApp(
		t,
		fxgorm.StartFxGorm(),
		fxhttpserver.StartFxHttpServer(),
		fx.Populate(&httpServer, &repository),
	)

	// prepare database post
	repository.Create(context.Background(), &model.Post{
		Title:       "test title",
		Description: "test description",
		Likes:       9,
	})

	// http request
	req := httptest.NewRequest("GET", "/posts", nil)
	rec := fxhttpservertest.RecordResponse(httpServer, req)

	// http response assertion
	fxhttpservertest.AssertRecordedResponseCode(t, rec, 200)
	fxhttpservertest.AssertRecordedResponseBody(t, rec, `"title":"test title","description":"test description"`)

	// log assertion
	fxloggertest.AssertHasLogRecord(
		t,
		map[string]interface{}{
			"level":   "info",
			"message": "in list posts handler",
		},
	)

	// trace assertion
	fxtracertest.AssertHasTraceSpan(
		t,
		"gorm.Query",
		attribute.String("db.statement", "SELECT * FROM `posts` WHERE `posts`.`deleted_at` IS NULL"),
	)
}
