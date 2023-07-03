package fxgrpcserver

import (
	"context"

	"github.com/ekkinox/fx-template/modules/fxhealthchecker"
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

type GrpcHealthCheckServer struct {
	grpc_health_v1.UnimplementedHealthServer
	healthChecker *fxhealthchecker.HealthChecker
	logger        *fxlogger.Logger
}

func NewGrpcHealthCheckServer(healthChecker *fxhealthchecker.HealthChecker, logger *fxlogger.Logger) *GrpcHealthCheckServer {
	return &GrpcHealthCheckServer{
		healthChecker: healthChecker,
		logger:        logger,
	}
}

func (s *GrpcHealthCheckServer) Check(ctx context.Context, in *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {

	CtxLogger(ctx).Info().Msgf("grpc health check called by %s", in.Service)

	result := s.healthChecker.Run(ctx)

	hcStatus := grpc_health_v1.HealthCheckResponse_SERVING
	if !result.Success {
		hcStatus = grpc_health_v1.HealthCheckResponse_NOT_SERVING
	}

	return &grpc_health_v1.HealthCheckResponse{Status: hcStatus}, nil
}

func (s *GrpcHealthCheckServer) Watch(in *grpc_health_v1.HealthCheckRequest, watchServer grpc_health_v1.Health_WatchServer) error {

	CtxLogger(watchServer.Context()).Info().Msgf("grpc health watch called by %s", in.Service)

	return status.Error(codes.Unimplemented, "watch is not implemented")
}
