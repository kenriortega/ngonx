package otelify

import (
	"context"
	"log"

	"github.com/kenriortega/ngonx/pkg/errors"
	"github.com/kenriortega/ngonx/pkg/logger"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

// Initializes an OTLP exporter, and configures the corresponding trace and
// metric providers.
func InitProvider(name, version, namEnv, endpoint string) func() {
	ctx := context.Background()

	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
	)
	handleErr(err, "failed to create trace exporter")

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(NewResource(name, version, namEnv)),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return func() {
		// Shutdown will flush any remaining spans and shut down the exporter.
		handleErr(tracerProvider.Shutdown(ctx), "failed to shutdown TracerProvider")
	}
}
func handleErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}

// NewResource returns a resource describing this application.
func NewResource(name, version, namEnv string) *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(name),
			semconv.ServiceVersionKey.String(version),
			attribute.String("environment", namEnv),
		),
	)
	return r
}

func InstrumentedError(span trace.Span, source, traceID string, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
	logger.LogError(
		errors.Errorf("%s: %v", source, err).Error(),
		zap.String("traceID", traceID),
	)
}

func InstrumentedInfo(span trace.Span, source, traceID string) {
	logger.LogInfo(
		source,
		zap.String("traceID", traceID),
	)
	span.SetAttributes(attribute.String(source, "Success"))
}
