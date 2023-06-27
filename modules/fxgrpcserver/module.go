package fxgrpcserver

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const DefaultPort = 50051

var FxGrpcServerModule = fx.Module(
	"grpc-server",
	fx.Provide(
		NewFxGrpcServer,
	),
	fx.Invoke(func(*grpc.Server) {}),
)

type FxGrpcServerParam struct {
	fx.In
	LifeCycle               fx.Lifecycle
	Config                  *fxconfig.Config
	Logger                  *fxlogger.Logger
	GrpcServices            []any                   `group:"grpc-server-services"`
	GrpcServicesDefinitions []GrpcServiceDefinition `group:"grpc-server-service-definitions"`
}

func NewFxGrpcServer(p FxGrpcServerParam) (*grpc.Server, error) {

	port := p.Config.GetInt("grpc.server.port")
	if port == 0 {
		port = DefaultPort
	}

	p.Logger.Info().Msgf("************ services definitions: %+#v", p.GrpcServicesDefinitions)
	p.Logger.Info().Msgf("************ services: %+v", p.GrpcServices)

	grpcServer := grpc.NewServer(
	//grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
	//grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)

	if p.Config.GetBool("grpc.server.reflection") {
		reflection.Register(grpcServer)
	}

	for _, def := range p.GrpcServicesDefinitions {

		service, err := lookupRegisteredService(def.Service(), p)
		if err != nil {
			return nil, err
		}

		grpcServer.RegisterService(def.Description(), service)
	}

	p.LifeCycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
				if err != nil {
					p.Logger.Error().Err(err).Msgf("failed to listen on %d for grpc server", port)
				}

				if err = grpcServer.Serve(lis); err != nil {
					p.Logger.Error().Err(err).Msgf("failed to serve grpc server", port)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			grpcServer.GracefulStop()

			return nil
		},
	})

	return grpcServer, nil
}

func lookupRegisteredService(service string, params FxGrpcServerParam) (any, error) {
	for _, s := range params.GrpcServices {
		params.Logger.Info().Msgf("&&&&&&&&&&&& service type: %s", getType(s))
		if getType(s) == service {
			return s, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("cannot find service for type %s", service))
}
