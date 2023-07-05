package fxgrpcserver

import (
	"google.golang.org/grpc"
)

type options struct {
	ServerOptions []grpc.ServerOption
	Reflection    bool
}

var defaultGrpcServerOptions = options{
	ServerOptions: []grpc.ServerOption{},
	Reflection:    false,
}

type GrpcServerOption func(o *options)

func WithServerOptions(s ...grpc.ServerOption) GrpcServerOption {
	return func(o *options) {
		o.ServerOptions = s
	}
}

func WithReflection(r bool) GrpcServerOption {
	return func(o *options) {
		o.Reflection = r
	}
}
