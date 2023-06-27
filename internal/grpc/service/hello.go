package service

import (
	"context"

	"github.com/ekkinox/fx-template/modules/fxlogger"
	"github.com/ekkinox/fx-template/proto"
)

type HelloServer struct {
	proto.UnimplementedHelloServiceServer
	logger *fxlogger.Logger
}

func NewHelloServer(logger *fxlogger.Logger) *HelloServer {
	return &HelloServer{
		logger: logger,
	}
}

func (s *HelloServer) SayHello(ctx context.Context, in *proto.HelloRequest) (*proto.HelloResponse, error) {

	s.logger.Info().Msgf("called SayHello with %s", in.Name)

	return &proto.HelloResponse{Message: "Hello, " + in.Name}, nil
}
