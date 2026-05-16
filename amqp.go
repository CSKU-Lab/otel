package otel

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// InjectTraceHeaders extracts W3C TraceContext from ctx and returns it
// as an AMQP-compatible header map for use in queue.Derivery.Headers.
func InjectTraceHeaders(ctx context.Context) map[string]interface{} {
	carrier := make(propagation.MapCarrier)
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	headers := make(map[string]interface{}, len(carrier))
	for k, v := range carrier {
		headers[k] = v
	}
	return headers
}

// ExtractTraceContext reconstructs a context with the remote span from
// AMQP message headers produced by InjectTraceHeaders.
func ExtractTraceContext(headers map[string]interface{}) context.Context {
	carrier := make(propagation.MapCarrier, len(headers))
	for k, v := range headers {
		if s, ok := v.(string); ok {
			carrier[k] = s
		}
	}
	return otel.GetTextMapPropagator().Extract(context.Background(), carrier)
}
