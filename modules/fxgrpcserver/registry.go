package fxgrpcserver

import (
	"errors"
	"fmt"

	"go.uber.org/fx"
)

type GrpcServiceRegistry struct {
	grpcServices            []any
	grpcServicesDefinitions []GrpcServiceDefinition
}

type FxGrpcServiceRegistryParam struct {
	fx.In
	GrpcServices            []any                   `group:"grpc-server-services"`
	GrpcServicesDefinitions []GrpcServiceDefinition `group:"grpc-server-service-definitions"`
}

func NewGrpcServiceRegistry(p FxGrpcServiceRegistryParam) *GrpcServiceRegistry {
	return &GrpcServiceRegistry{
		grpcServices:            p.GrpcServices,
		grpcServicesDefinitions: p.GrpcServicesDefinitions,
	}
}

func (r *GrpcServiceRegistry) ResolveGrpcServices() ([]GrpcService, error) {

	var grpcServices []GrpcService

	for _, def := range r.grpcServicesDefinitions {
		implementation, err := r.lookupRegisteredServiceImplementation(def.ReturnType())
		if err != nil {
			return nil, err
		}

		grpcServices = append(grpcServices, newResolvedGrpcService(def.Description(), implementation))
	}

	return grpcServices, nil
}

func (r *GrpcServiceRegistry) lookupRegisteredServiceImplementation(returnType string) (any, error) {
	for _, implementation := range r.grpcServices {
		if getType(implementation) == returnType {
			return implementation, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("cannot find grpc service implementation for type %s", returnType))
}
