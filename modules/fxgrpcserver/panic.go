package fxgrpcserver

import (
	"runtime/debug"

	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcPanicRecoveryHandler struct {
	config *fxconfig.Config
	logger *fxlogger.Logger
}

func NewGrpcPanicRecoveryHandler(config *fxconfig.Config, logger *fxlogger.Logger) *GrpcPanicRecoveryHandler {
	return &GrpcPanicRecoveryHandler{
		config: config,
		logger: logger,
	}
}

func (h *GrpcPanicRecoveryHandler) Handle() func(any) error {
	return func(pnc any) error {
		h.logger.Error().Msgf("grpc recovering from panic, panic:%s, stack: %s", pnc, debug.Stack())

		if h.config.AppDebug() {
			return status.Errorf(codes.Internal, "internal grpc server error, panic:%s, stack: %s", pnc, debug.Stack())
		} else {
			return status.Error(codes.Internal, "internal grpc server error")
		}
	}
}
