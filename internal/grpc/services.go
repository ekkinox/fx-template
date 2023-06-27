package grpc

import (
	"github.com/ekkinox/fx-template/internal/grpc/service"
	"github.com/ekkinox/fx-template/modules/fxgrpcserver"
	"github.com/ekkinox/fx-template/proto"
	"go.uber.org/fx"
)

func RegisterGrpcServices() fx.Option {
	return fx.Options(
		fxgrpcserver.AsGrpcService(&proto.HelloService_ServiceDesc, service.NewHelloServer),
		fxgrpcserver.AsGrpcService(&proto.GoodbyeService_ServiceDesc, service.NewGoodbyeServer),
	)
}
