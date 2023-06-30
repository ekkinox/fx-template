package fxgrpcserver

import "google.golang.org/grpc"

type resolvedGrpcService struct {
	description    *grpc.ServiceDesc
	implementation any
}

func newResolvedGrpcService(description *grpc.ServiceDesc, implementation any) GrpcService {
	return &resolvedGrpcService{
		description:    description,
		implementation: implementation,
	}
}

func (r *resolvedGrpcService) Description() *grpc.ServiceDesc {
	return r.description
}

func (r *resolvedGrpcService) Implementation() any {
	return r.implementation
}
