package fxgrpcserver

import (
	"context"
	"time"

	"github.com/ekkinox/fx-template/modules/fxlogger"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type LoggerInterceptor struct {
	logger *fxlogger.Logger
}

func NewLoggerInterceptor(logger *fxlogger.Logger) *LoggerInterceptor {
	return &LoggerInterceptor{
		logger: logger,
	}
}

func (i *LoggerInterceptor) UnaryInterceptor() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		reqId, traceParent := i.extractMetadata(ctx)

		grpcLogger := i.logger.
			With().
			Str("x-request-id", reqId).
			Str("traceparent", traceParent).
			Logger()

		grpcLogger.
			Info().
			Str("grpc-type", "unary").
			Str("grpc-method", info.FullMethod).
			Msg("grpc call start")

		newCtx := grpcLogger.WithContext(ctx)

		now := time.Now()

		resp, err := handler(newCtx, req)

		if err != nil {

			errStatus := status.Convert(err)

			grpcLogger.
				Error().
				Err(err).
				Int32("grpc-code", int32(errStatus.Code())).
				Str("grpc-status", errStatus.Code().String()).
				Str("grpc-duration", time.Since(now).String()).
				Msg("grpc call error")
		} else {
			grpcLogger.
				Info().
				Int32("grpc-code", int32(codes.OK)).
				Str("grpc-status", codes.OK.String()).
				Str("grpc-duration", time.Since(now).String()).
				Msg("grpc call success")
		}

		return resp, err
	}
}

func (i *LoggerInterceptor) StreamInterceptor() func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		ctx := ss.Context()
		reqId, traceParent := i.extractMetadata(ctx)

		grpcLogger := i.logger.
			With().
			Str("x-request-id", reqId).
			Str("traceparent", traceParent).
			Logger()

		grpcLogger.
			Info().
			Str("grpc-type", "server-streaming").
			Str("grpc-method", info.FullMethod).
			Msg("grpc call start")

		newCtx := grpcLogger.WithContext(ctx)

		wrapped := middleware.WrapServerStream(ss)
		wrapped.WrappedContext = newCtx

		now := time.Now()

		err := handler(newCtx, wrapped)

		if err != nil {

			errStatus := status.Convert(err)

			grpcLogger.
				Error().
				Err(err).
				Int32("grpc-code", int32(errStatus.Code())).
				Str("grpc-status", errStatus.Code().String()).
				Str("grpc-duration", time.Since(now).String()).
				Msg("grpc call error")
		} else {
			grpcLogger.
				Info().
				Int32("grpc-code", int32(codes.OK)).
				Str("grpc-status", codes.OK.String()).
				Str("grpc-duration", time.Since(now).String()).
				Msg("grpc call success")
		}

		return err
	}
}

func (i *LoggerInterceptor) extractMetadata(ctx context.Context) (string, string) {
	var reqId string
	var traceParent string

	md, _ := metadata.FromIncomingContext(ctx)

	if requestId, ok := md["x-request-id"]; ok && len(requestId) > 0 {
		reqId = requestId[0]
	} else {
		reqId = generateRequestId()
	}

	if requestTraceParent, ok := md["traceparent"]; ok && len(requestTraceParent) > 0 {
		traceParent = requestTraceParent[0]
	} else {
		traceParent = ""
	}

	return reqId, traceParent
}
