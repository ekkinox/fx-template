package fxhttpserver

import (
	"errors"

	"github.com/ekkinox/fx-template/modules/fxauthenticationcontext"
	"github.com/ekkinox/fx-template/modules/fxhttpclient"
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const ctxTracerKey = "otel-go-contrib-tracer-labstack-echo"

func CtxLogger(c echo.Context) *fxlogger.Logger {
	return fxlogger.CtxLogger(c.Request().Context())
}

func CtxTracer(c echo.Context) trace.Tracer {
	if tracer, ok := c.Request().Context().Value(ctxTracerKey).(trace.Tracer); ok {
		return tracer
	} else {
		return otel.Tracer("fxhttpserver-tracer")
	}
}

func CtxHttpClient(c echo.Context, opts ...fxhttpclient.HttpClientOption) *fxhttpclient.HttpClient {

	opts = append(
		opts,
		fxhttpclient.WithRequestHeadersToForward(c.Request(), fxhttpclient.DefaultHeadersToForward),
	)

	return fxhttpclient.NewCtxHttpClient(c.Request().Context(), opts...)
}

func CtxAuthenticationContext(c echo.Context) (*fxauthenticationcontext.AuthenticationContext, error) {

	if authContext, ok := c.Get(fxauthenticationcontext.AuthenticationContextKey).(*fxauthenticationcontext.AuthenticationContext); ok {
		return authContext, nil
	}

	return nil, errors.New("cannot get authentication context")
}
