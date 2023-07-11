package server

import (
	"context"

	"github.com/ekkinox/fx-template/internal/server/grpc"
	"github.com/ekkinox/fx-template/internal/server/http"
	"github.com/ekkinox/fx-template/modules/fxbootstrapper"
	"github.com/ekkinox/fx-template/modules/fxgrpcserver"
	"github.com/ekkinox/fx-template/modules/fxhttpserver"
	"go.uber.org/fx"
)

var ServerBoostrapper = fxbootstrapper.NewBootstrapper().WithOptions(
	// common
	RegisterModules(),
	RegisterServices(),
	RegisterOverrides(),
	// http server
	fxhttpserver.FxHttpServerModule,
	http.RegisterHandlers(),
	// grpc server
	fxgrpcserver.FxGrpcServerModule,
	grpc.RegisterGrpcServices(),
)

func BootstrapServer(ctx context.Context) *fx.App {
	return ServerBoostrapper.BoostrapApp(
		fxhttpserver.StartFxHttpServer(),
		fxgrpcserver.StartFxGrpcServer(),
	)
}
