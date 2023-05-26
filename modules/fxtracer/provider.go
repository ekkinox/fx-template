package fxtracer

import (
	"context"
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"time"

	"github.com/ekkinox/fx-template/modules/fxconfig"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewTracerProvider(config *fxconfig.Config, logger *fxlogger.Logger) (*trace.TracerProvider, error) {

	bgCtx := context.Background()

	res, err := resource.New(
		bgCtx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(config.AppName()),
		),
	)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create tracing resource")
		return nil, err
	}

	var bsp trace.SpanProcessor
	if config.GetBool("TRACING_ENABLED") {

		dialCtx, cancel := context.WithTimeout(bgCtx, 5*time.Second)
		defer cancel()

		conn, err := grpc.DialContext(
			dialCtx,
			config.GetString("TRACING_COLLECTOR"),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		)
		if err != nil {
			logger.Error().Err(err).Msg("failed to create gRPC connection to tracing collector")
			return nil, err
		}

		traceExporter, err := otlptracegrpc.New(dialCtx, otlptracegrpc.WithGRPCConn(conn))
		if err != nil {
			logger.Error().Err(err).Msg("failed to create gRPC tracing exporter")
			return nil, err
		}

		bsp = trace.NewBatchSpanProcessor(traceExporter)
	} else {
		//	=> "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
		//exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
		//if err != nil {
		//	return nil, err
		//}
		//bsp = trace.NewBatchSpanProcessor(exporter)

		bsp = trace.NewBatchSpanProcessor(tracetest.NewNoopExporter())
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(res),
		trace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	logger.Debug().Msg("tracer is ready")

	return tracerProvider, nil
}
