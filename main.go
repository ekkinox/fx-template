package main

import (
	"net/http"

	"github.com/ekkinox/fx-template/handler"
	"github.com/ekkinox/fx-template/server"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.Provide(
			server.NewHTTPServer,
			fx.Annotate(
				server.NewServeMux,
				fx.ParamTags(`group:"routes"`),
			),
			server.AsRoute(handler.NewEchoHandler),
			server.AsRoute(handler.NewHelloHandler),
			zap.NewExample,
		),
		fx.Invoke(func(*http.Server) {}),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
	).Run()
}
