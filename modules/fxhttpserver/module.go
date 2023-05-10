package fxhttpserver

import (
	"context"
	"fmt"
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"net"
	"net/http"
)

var FxHttpServerModule = fx.Module("http-server",
	fx.Provide(
		NewEcho,
		NewFxHttpServer,
	),
	fx.Invoke(func(*http.Server) {}),
)

type FxHttpServerParam struct {
	fx.In
	Config    *fxconfig.Config
	Echo      *echo.Echo
	LifeCycle fx.Lifecycle
}

func NewFxHttpServer(p FxHttpServerParam) *http.Server {

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", p.Config.AppPort),
		Handler: p.Echo,
	}

	p.LifeCycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", server.Addr)
			if err != nil {
				return err
			}
			go server.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})

	return server
}
