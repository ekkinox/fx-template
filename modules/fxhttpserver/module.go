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

var HttpServerModule = fx.Module("http-server",
	fx.Provide(
		NewEcho,
		NewHTTPServer,
	),
	fx.Invoke(func(*http.Server) {}),
)

type HTTPServerParam struct {
	fx.In
	Config    *fxconfig.Config
	Echo      *echo.Echo
	LifeCycle fx.Lifecycle
}

func NewHTTPServer(p HTTPServerParam) *http.Server {

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
