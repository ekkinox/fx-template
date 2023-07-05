package grpc

import (
	"github.com/ekkinox/fx-template/internal/server/grpc/grpcservice"
	"github.com/ekkinox/fx-template/modules/fxgrpcserver"
	"github.com/ekkinox/fx-template/proto/ping"
	"github.com/ekkinox/fx-template/proto/posts"
	"go.uber.org/fx"
)

func RegisterGrpcServices() fx.Option {
	return fx.Options(
		fxgrpcserver.AsGrpcService(&ping.PingService_ServiceDesc, grpcservice.NewPingServer),
		fxgrpcserver.AsGrpcService(&posts.PostCrudService_ServiceDesc, grpcservice.NewPostsCrudServer),
	)
}
