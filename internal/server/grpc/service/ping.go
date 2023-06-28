package service

import (
	"context"

	"github.com/ekkinox/fx-template/modules/fxlogger"
	"github.com/ekkinox/fx-template/proto/ping"
)

type PingServer struct {
	ping.UnimplementedPingServiceServer
	logger *fxlogger.Logger
}

func NewPingServer(logger *fxlogger.Logger) *PingServer {
	return &PingServer{
		logger: logger,
	}
}

func (s *PingServer) Ping(ctx context.Context, in *ping.PingRequest) (*ping.PingResponse, error) {

	s.logger.Info().Msgf("called SayGoodbye with %s", in.Message)

	return &ping.PingResponse{Message: "Goodbye, " + in.Message}, nil
}
