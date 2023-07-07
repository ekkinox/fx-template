package server

import (
	"context"

	"github.com/ekkinox/fx-template/internal/server/grpc"
	"github.com/ekkinox/fx-template/internal/server/http"
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxgrpcserver"
	"github.com/ekkinox/fx-template/modules/fxhealthchecker"
	"github.com/ekkinox/fx-template/modules/fxhttpserver"
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"github.com/ekkinox/fx-template/modules/fxtracer"
	"go.uber.org/fx"
)

func BootstrapServer(ctx context.Context, options ...fx.Option) *fx.App {
	return fx.New(
		// logger
		fx.WithLogger(fxlogger.FxEventLogger),
		// core
		fxconfig.FxConfigModule,
		fxlogger.FxLoggerModule,
		fxtracer.FxTracerModule,
		fxhealthchecker.FxHealthCheckerModule,
		// common
		RegisterModules(ctx),
		RegisterServices(ctx),
		RegisterOverrides(ctx),
		// http server
		fxhttpserver.FxHttpServerModule,
		http.RegisterHandlers(),
		// grpc server
		fxgrpcserver.FxGrpcServerModule,
		grpc.RegisterGrpcServices(),
		//options
		fx.Options(options...),
	)
}
