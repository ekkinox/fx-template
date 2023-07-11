package fxgrpcserver

import (
	"context"
	"fmt"
	"net"

	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxhealthchecker"
	"github.com/ekkinox/fx-template/modules/fxlogger"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

const DefaultPort = 50051

var FxGrpcServerModule = fx.Module(
	"grpc-server",
	fx.Provide(
		NewDefaultGrpcServerFactory,
		NewFxGrpcServiceRegistry,
		NewFxGrpcServer,
	),
)

type FxGrpcServerParam struct {
	fx.In
	LifeCycle     fx.Lifecycle
	Factory       GrpcServerFactory
	Registry      *GrpcServiceRegistry
	Config        *fxconfig.Config
	Logger        *fxlogger.Logger
	HealthChecker *fxhealthchecker.HealthChecker
}

func StartFxGrpcServer() fx.Option {
	return fx.Invoke(func(*GrpcServiceRegistry, *grpc.Server) {})
}

func NewFxGrpcServer(p FxGrpcServerParam) (*grpc.Server, error) {

	// interceptors
	grpcPanicRecoveryHandler := NewGrpcPanicRecoveryHandler(p.Config, p.Logger)

	loggerInterceptor := NewLoggerInterceptor(p.Logger)

	unaryInterceptors := []grpc.UnaryServerInterceptor{
		loggerInterceptor.UnaryInterceptor(),
		recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler.Handle())),
	}
	streamInterceptors := []grpc.StreamServerInterceptor{
		loggerInterceptor.StreamInterceptor(),
		recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler.Handle())),
	}

	if p.Config.GetBool("modules.grpc.tracer.enabled") {
		unaryInterceptors = append(unaryInterceptors, otelgrpc.UnaryServerInterceptor())
		streamInterceptors = append(streamInterceptors, otelgrpc.StreamServerInterceptor())
	}

	// server
	grpcServer, err := p.Factory.Create(
		WithServerOptions(
			middleware.WithUnaryServerChain(unaryInterceptors...),
			middleware.WithStreamServerChain(streamInterceptors...),
		),
		WithReflection(p.Config.GetBool("modules.grpc.server.reflection")),
	)
	if err != nil {
		p.Logger.Error().Err(err).Msg("failed to create grpc server")

		return nil, err
	}

	// health check registration
	grpcServer.RegisterService(
		&grpc_health_v1.Health_ServiceDesc,
		NewGrpcHealthCheckServer(p.HealthChecker, p.Logger),
	)

	// services registration
	grpcServices, err := p.Registry.ResolveGrpcServices()
	if err != nil {
		p.Logger.Error().Err(err).Msg("failed to resolve grpc services")

		return nil, err
	}

	for _, service := range grpcServices {
		grpcServer.RegisterService(service.Description(), service.Implementation())
	}

	// lifecycles
	p.LifeCycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {

			port := p.Config.GetInt("modules.grpc.server.port")
			if port == 0 {
				port = DefaultPort
			}

			go func() {
				lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
				if err != nil {
					p.Logger.Error().Err(err).Msgf("failed to listen on %d for grpc server", port)
				}

				if err = grpcServer.Serve(lis); err != nil {
					p.Logger.Error().Err(err).Msg("failed to serve grpc server")
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
