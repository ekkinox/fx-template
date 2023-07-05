package fxhttpserver

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type options struct {
	Debug            bool
	Banner           bool
	Logger           echo.Logger
	Binder           echo.Binder
	JsonSerializer   echo.JSONSerializer
	HttpErrorHandler echo.HTTPErrorHandler
}

var defaultHttpServerOptions = options{
	Debug:            false,
	Banner:           false,
	Logger:           log.New("default"),
	Binder:           &echo.DefaultBinder{},
	JsonSerializer:   &echo.DefaultJSONSerializer{},
	HttpErrorHandler: nil,
}

type HttpServerOption func(o *options)

func WithDebug(d bool) HttpServerOption {
	return func(o *options) {
		o.Debug = d
	}
}

func WithBanner(b bool) HttpServerOption {
	return func(o *options) {
		o.Banner = b
	}
}

func WithLogger(l echo.Logger) HttpServerOption {
	return func(o *options) {
		o.Logger = l
	}
}

func WithBinder(b echo.Binder) HttpServerOption {
	return func(o *options) {
		o.Binder = b
	}
}

func WithJsonSerializer(s echo.JSONSerializer) HttpServerOption {
	return func(o *options) {
		o.JsonSerializer = s
	}
}

func WithHttpErrorHandler(h echo.HTTPErrorHandler) HttpServerOption {
	return func(o *options) {
		o.HttpErrorHandler = h
	}
}
