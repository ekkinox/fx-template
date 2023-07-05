package grpcservice

import (
	"context"

	"github.com/ekkinox/fx-template/modules/fxgrpcserver"
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

	tracer := fxgrpcserver.CtxTracer(ctx)
	ctx, span := tracer.Start("ping")
	defer span.End()

	logger := fxgrpcserver.CtxLogger(ctx)
	logger.Info().Msg("test 123")
	logger.Info().Msgf("called Ping with %s", in.Message)

	return &ping.PingResponse{Message: "Your message was: " + in.Message}, nil
}
