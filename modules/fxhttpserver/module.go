package fxhttpserver

import (
	"context"
	"fmt"
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxhealthchecker"
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"github.com/ekkinox/fx-template/modules/fxtracer"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.uber.org/fx"
	"net/http"
	"reflect"
	"strings"
)

var FxHttpServerModule = fx.Module("http-server",
	// modules dependencies
	fxconfig.FxConfigModule,
	fxlogger.FxLoggerModule,
	fxtracer.FxTracerModule,
	fxhealthchecker.FxHealthCheckerModule,
	// http server
	fx.Provide(
		NewFxHttpServer,
	),
	fx.Invoke(func(*echo.Echo) {}),
)

type FxHttpServerParam struct {
	fx.In
	LifeCycle    fx.Lifecycle
	Config       *fxconfig.Config
	Logger       *fxlogger.Logger
	HeathChecker *fxhealthchecker.Checker
	Handlers     []Handler    `group:"http-server-handlers"`
	Middlewares  []Middleware `group:"http-server-middlewares"`
	Routes       []Route      `group:"http-server-routes"`
}

func NewFxHttpServer(p FxHttpServerParam) *echo.Echo {
	// echo
	e := echo.New()
	e.HideBanner = true
	e.Debug = p.Config.AppConfig.Debug
	e.Logger = p.Logger

	// middlewares
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(otelecho.Middleware(p.Config.AppConfig.Name))
	e.Use(fxlogger.Middleware(fxlogger.Config{
		Logger:      p.Logger,
		HandleError: true,
	}))

	// handlers groups
	/*	for _, hg := range p.HandlersGroups {
		g := e.Group(hg.Prefix(), hg.Middlewares()...)
		for _, h := range hg.Handlers() {
			g.Add(strings.ToUpper(h.Method()), h.Path(), h.Handler(), h.Middlewares()...)
		}
	}*/

	// handlers
	for _, handler := range p.Handlers {

		handlerType := reflect.TypeOf(handler).String()
		handlerRoute, err := findRouteForHandlerType(p.Routes, handlerType)
		if err != nil {
			p.Logger.Error("cannot register handler")
		}

		fmt.Printf("\nroute: %+v\n", handlerRoute)

		middlewaresRawMap := map[string]echo.MiddlewareFunc{}
		for _, rmt := range handlerRoute.MiddlewareTypes() {
			for _, pm := range p.Middlewares {
				pmt := reflect.TypeOf(pm).String()
				if pmt == rmt {
					middlewaresRawMap[pmt] = pm.Handle()
				}
			}
		}

		var middlewaresCleanSlice []echo.MiddlewareFunc
		for _, vm := range middlewaresRawMap {
			middlewaresCleanSlice = append(middlewaresCleanSlice, vm)
		}
		middlewaresCleanSlice = append(middlewaresCleanSlice, handlerRoute.MiddlewareInstances()...)
		middlewaresCleanSlice = append(middlewaresCleanSlice, middleware.AddTrailingSlash())

		e.Add(
			strings.ToUpper(handlerRoute.Method()),
			handlerRoute.Path(),
			handler.Handle(),
			middlewaresCleanSlice...,
		)

		p.Logger.Infof("registering handler type %s for [%s]%s", handlerType, handlerRoute.Method(), handlerRoute.Path())
	}

	// debugger
	if p.Config.AppConfig.Debug {
		g := e.Group("/_debug")
		// routes
		g.GET("/routes", func(c echo.Context) error {
			return c.JSON(http.StatusOK, e.Routes())
		})
		// version
		g.GET("/version", func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{"version": "0.1.0"})
		})
	}

	// health check
	e.GET("/_health", func(c echo.Context) error {
		r := p.HeathChecker.Run()

		status := http.StatusOK
		if !r.Success {
			status = http.StatusInternalServerError
		}

		return c.JSON(status, r)
	})

	// lifecycles
	p.LifeCycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go e.Start(fmt.Sprintf(":%d", p.Config.AppConfig.Port))
			return nil

		},
		OnStop: func(ctx context.Context) error {
			return e.Shutdown(ctx)
		},
	})

	return e
}
