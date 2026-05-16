package otel

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// WithTrace enriches a zap logger with trace_id and span_id fields
// from the active span in ctx. Returns logger unchanged if no valid span.
func WithTrace(ctx context.Context, l *zap.SugaredLogger) *zap.SugaredLogger {
	sc := trace.SpanFromContext(ctx).SpanContext()
	if !sc.IsValid() {
		return l
	}
	return l.With(
		"trace_id", sc.TraceID().String(),
		"span_id", sc.SpanID().String(),
	)
}
