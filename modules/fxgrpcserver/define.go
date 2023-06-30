package fxgrpcserver

import "google.golang.org/grpc"

type GrpcServiceDefinition interface {
	Description() *grpc.ServiceDesc
	ReturnType() string
}

type grpcServiceDefinition struct {
	description *grpc.ServiceDesc
	returnType  string
}

func newGrpcServiceDefinition(description *grpc.ServiceDesc, returnType string) *grpcServiceDefinition {
	return &grpcServiceDefinition{
		description: description,
		returnType:  returnType,
	}
}

func (d *grpcServiceDefinition) Description() *grpc.ServiceDesc {
	return d.description
}

func (d *grpcServiceDefinition) ReturnType() string {
	return d.returnType
}
