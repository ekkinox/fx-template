package fxhttpserver

import (
	"github.com/ekkinox/fx-template/modules/fxhttpclient"
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"github.com/ekkinox/fx-template/modules/fxtracer"
	"github.com/labstack/echo/v4"
)

func CtxLogger(ctx echo.Context) *fxlogger.Logger {
	return fxlogger.CtxLogger(ctx.Request().Context())
}

func CtxTracer(ctx echo.Context) *fxtracer.Tracer {
	return fxtracer.CtxTracer(ctx.Request().Context())
}

func CtxHttpClient(ctx echo.Context, opts ...fxhttpclient.HttpClientOption) *fxhttpclient.HttpClient {
	return fxhttpclient.NewCtxHttpClient(ctx.Request().Context(), opts...)
}
