package fxgrpcserver

import "google.golang.org/grpc"

type GrpcServiceDefinition interface {
	Description() *grpc.ServiceDesc
	Service() string
}

type grpcServiceDefinition struct {
	description *grpc.ServiceDesc
	service     string
}

func newGrpcServiceDefinition(description *grpc.ServiceDesc, service string) *grpcServiceDefinition {
	return &grpcServiceDefinition{
		description: description,
		service:     service,
	}
}

func (d *grpcServiceDefinition) Description() *grpc.ServiceDesc {
	return d.description
}

func (d *grpcServiceDefinition) Service() string {
	return d.service
}
