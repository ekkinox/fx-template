package fxgrpcserver

import (
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type GrpcServiceRegistration struct {
	description *grpc.ServiceDesc
	constructor any
}

func NewGrpcServiceRegistration(description *grpc.ServiceDesc, constructor any) *GrpcServiceRegistration {
	return &GrpcServiceRegistration{
		description: description,
		constructor: constructor,
	}
}

func (r *GrpcServiceRegistration) Description() *grpc.ServiceDesc {
	return r.description
}

func (r *GrpcServiceRegistration) Constructor() any {
	return r.constructor
}

func AsGrpcService(description *grpc.ServiceDesc, constructor any) fx.Option {
	return RegisterGrpcService(NewGrpcServiceRegistration(description, constructor))
}

func RegisterGrpcService(grpcServiceRegistration *GrpcServiceRegistration) fx.Option {

	serviceDef := newGrpcServiceDefinition(
		grpcServiceRegistration.Description(),
		getReturnType(grpcServiceRegistration.Constructor()),
	)

	return fx.Options(
		fx.Provide(
			fx.Annotate(
				grpcServiceRegistration.Constructor(),
				fx.As(new(interface{})),
				fx.ResultTags(`group:"grpc-server-services"`),
			),
		),
		fx.Supply(
			fx.Annotate(
				serviceDef,
				fx.As(new(GrpcServiceDefinition)),
				fx.ResultTags(`group:"grpc-server-service-definitions"`),
			),
		),
	)
}
