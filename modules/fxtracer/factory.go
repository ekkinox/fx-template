package fxtracer

import (
	"context"
	"time"

	"github.com/ekkinox/fx-template/modules/fxlogger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TracerProviderFactory interface {
	Create(options ...TracerProviderOption) (*trace.TracerProvider, error)
}

type DefaultTracerProviderFactory struct {
	logger *fxlogger.Logger
}

func NewDefaultTracerProviderFactory(logger *fxlogger.Logger) TracerProviderFactory {
	return &DefaultTracerProviderFactory{
		logger: logger,
	}
}

func (f *DefaultTracerProviderFactory) Create(options ...TracerProviderOption) (*trace.TracerProvider, error) {

	appliedOptions := defaultTracerProviderOptions
	for _, opt := range options {
		opt(&appliedOptions)
	}

	ctx := context.Background()

	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(appliedOptions.Name),
		),
	)
	if err != nil {
		f.logger.Error().Err(err).Msg("failed to create resource for tracer provider")

		return nil, err
	}

	spanExporter, err := f.createSpanExporter(ctx, appliedOptions)
	if err != nil {
		f.logger.Error().Err(err).Msg("failed to create span exporter for tracer provider")

		return nil, err
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(res),
		trace.WithSpanProcessor(trace.NewBatchSpanProcessor(spanExporter)),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return tracerProvider, nil
}

func (f *DefaultTracerProviderFactory) createSpanExporter(ctx context.Context, opts options) (trace.SpanExporter, error) {
	switch opts.Exporter {
	case Memory:
		return tracetest.NewInMemoryExporter(), nil
	case Stdout:
		exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			return nil, err
		}

		return exporter, nil
	case OtlpGrpc:
		dialCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		conn, err := grpc.DialContext(
			dialCtx,
			opts.Collector,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		)
		if err != nil {
			f.logger.Error().Err(err).Msg("failed to create gRPC connection for otlp-grpc span exporter")

			return nil, err
		}

		exporter, err := otlptracegrpc.New(dialCtx, otlptracegrpc.WithGRPCConn(conn))
		if err != nil {
			f.logger.Error().Err(err).Msg("failed to create otlp-grpc span exporter")

			return nil, err
		}

		return exporter, nil
	default:
		return tracetest.NewNoopExporter(), nil
	}
}
