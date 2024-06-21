package slogger

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

// TODO 2024-06-20 Sam Borkent: add test cases.

const (
	traceIDKey = "traceId"
	spanIDKey  = "spanId"
)

type OtelHandler struct {
	slog.Handler
}

// Extract OpenTelemetry tracing data from context and add to logs if available.
func (t OtelHandler) Handle(ctx context.Context, record slog.Record) error {
	// Only trace logs that received non-empty contexts
	if ctx != context.Background() {
		// Extract traced span context from context
		spanCtx := trace.SpanContextFromContext(ctx)

		if spanCtx.HasTraceID() {
			record.AddAttrs(slog.String(traceIDKey, spanCtx.TraceID().String()))
		}

		if spanCtx.HasSpanID() {
			record.AddAttrs(slog.String(spanIDKey, spanCtx.SpanID().String()))
		}
	}

	return t.Handler.Handle(ctx, record)
}
