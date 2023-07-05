package fxhttpserver

import (
	"github.com/labstack/echo/v4"
)

type HttpServerFactory interface {
	Create(options ...HttpServerOption) (*echo.Echo, error)
}

type DefaultHttpServerFactory struct{}

func NewDefaultHttpServerFactory() HttpServerFactory {
	return &DefaultHttpServerFactory{}
}

func (f *DefaultHttpServerFactory) Create(options ...HttpServerOption) (*echo.Echo, error) {

	appliedOpts := defaultHttpServerOptions
	for _, applyOpt := range options {
		applyOpt(&appliedOpts)
	}

	httpServer := echo.New()

	httpServer.Debug = appliedOpts.Debug
	httpServer.HideBanner = !appliedOpts.Banner

	httpServer.Logger = appliedOpts.Logger
	httpServer.Binder = appliedOpts.Binder
	httpServer.JSONSerializer = appliedOpts.JsonSerializer

	if appliedOpts.HttpErrorHandler != nil {
		httpServer.HTTPErrorHandler = appliedOpts.HttpErrorHandler
	}

	return httpServer, nil
}
