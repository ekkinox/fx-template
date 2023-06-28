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
	checker *fxhealthchecker.Checker
	logger  *fxlogger.Logger
}

func NewGrpcHealthCheckServer(checker *fxhealthchecker.Checker, logger *fxlogger.Logger) *GrpcHealthCheckServer {
	return &GrpcHealthCheckServer{
		checker: checker,
		logger:  logger,
	}
}

func (s *GrpcHealthCheckServer) Check(ctx context.Context, in *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {

	s.logger.Info().Msgf("health check invoked with %v", in)

	result := s.checker.Run(ctx)

	hcStatus := grpc_health_v1.HealthCheckResponse_SERVING
	if !result.Success {
		hcStatus = grpc_health_v1.HealthCheckResponse_NOT_SERVING
	}

	return &grpc_health_v1.HealthCheckResponse{Status: hcStatus}, nil
}

func (s *GrpcHealthCheckServer) Watch(in *grpc_health_v1.HealthCheckRequest, watchServer grpc_health_v1.Health_WatchServer) error {

	s.logger.Info().Msgf("health watch invoked with %v", in)

	return status.Error(codes.Unimplemented, "watch is not implemented")
}
