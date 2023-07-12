package fxhttpservertest

import (
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
)

func RecordResponse(httpServer *echo.Echo, request *http.Request) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, request)

	return rec
}
