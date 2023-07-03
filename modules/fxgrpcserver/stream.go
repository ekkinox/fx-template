package fxgrpcserver

import (
	"context"

	"google.golang.org/grpc"
)

type WrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func NewWrappedStream(ss grpc.ServerStream, ctx context.Context) *WrappedStream {
	return &WrappedStream{ServerStream: ss, ctx: ctx}
}

func (w *WrappedStream) Context() context.Context {
	return w.ctx
}
