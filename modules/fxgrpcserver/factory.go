package fxgrpcserver

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcServerFactory interface {
	Create(options ...GrpcServerOption) (*grpc.Server, error)
}

type DefaultGrpcServerFactory struct{}

func NewDefaultGrpcServerFactory() GrpcServerFactory {
	return &DefaultGrpcServerFactory{}
}

func (f *DefaultGrpcServerFactory) Create(options ...GrpcServerOption) (*grpc.Server, error) {

	appliedOpts := defaultGrpcServerOptions
	for _, applyOpt := range options {
		applyOpt(&appliedOpts)
	}

	grpcServer := grpc.NewServer(appliedOpts.ServerOptions...)

	if appliedOpts.Reflection {
		reflection.Register(grpcServer)
	}

	return grpcServer, nil
}
