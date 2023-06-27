package service

import (
	"context"

	"github.com/ekkinox/fx-template/modules/fxlogger"
	"github.com/ekkinox/fx-template/proto"
)

type GoodbyeServer struct {
	proto.UnimplementedGoodbyeServiceServer
	logger *fxlogger.Logger
}

func NewGoodbyeServer(logger *fxlogger.Logger) *GoodbyeServer {
	return &GoodbyeServer{
		logger: logger,
	}
}

func (s *GoodbyeServer) SayGoodbye(ctx context.Context, in *proto.GoodbyeRequest) (*proto.GoodbyeResponse, error) {

	s.logger.Info().Msgf("called SayGoodbye with %s", in.Name)

	return &proto.GoodbyeResponse{Message: "Goodbye, " + in.Name}, nil
}
