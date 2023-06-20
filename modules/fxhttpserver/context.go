package fxhttpserver

import (
	"errors"

	"github.com/ekkinox/fx-template/modules/fxauthenticationcontext"
	"github.com/ekkinox/fx-template/modules/fxhttpclient"
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"github.com/ekkinox/fx-template/modules/fxtracer"
	"github.com/labstack/echo/v4"
)

func CtxLogger(c echo.Context) *fxlogger.Logger {
	return fxlogger.CtxLogger(c.Request().Context())
}

func CtxTracer(c echo.Context) *fxtracer.Tracer {
	return fxtracer.CtxTracer(c.Request().Context())
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
