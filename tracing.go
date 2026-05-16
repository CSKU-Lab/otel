package otel

import (
	"context"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// Init initialises the global TracerProvider and W3C propagator.
// Reads OTEL_EXPORTER_OTLP_ENDPOINT (e.g. http://jaeger:4317) and
// OTEL_SERVICE_NAME from env. The OTLP gRPC exporter handles the URL
// scheme correctly — do not strip http:// from the env var.
// Returns a shutdown function that must be called before process exit.
func Init(ctx context.Context) (shutdown func(context.Context) error, err error) {
	serviceName := os.Getenv("OTEL_SERVICE_NAME")

	// otlptracegrpc reads OTEL_EXPORTER_OTLP_ENDPOINT automatically and
	// handles the http:// scheme (insecure) vs https:// (TLS).
	exporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		return nil, err
	}

	res, _ := resource.New(ctx,
		resource.WithAttributes(semconv.ServiceName(serviceName)),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp.Shutdown, nil
}
