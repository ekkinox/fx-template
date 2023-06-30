package fxgrpcserver

import "google.golang.org/grpc"

type GrpcService interface {
	Description() *grpc.ServiceDesc
	Implementation() any
}
