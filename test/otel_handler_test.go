package slogger_test

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/samborkent/slogger"
)

// This test and its subtests cannot run in parallel due to stderr capture.
func TestOtelHandler(t *testing.T) {
	t.Run("untraced", func(t *testing.T) {
		// Setup stderr capture to capture log output.
		testStderr := os.Stderr

		reader, writer, err := os.Pipe()
		if err != nil {
			t.Errorf("creating i/o pipeline: " + err.Error())
			return
		}

		os.Stderr = writer

		// Initialize slogger
		log, err := slogger.NewWithOptions(slogger.Options{
			TracingEnabled: true,
		})
		if err != nil {
			t.Errorf("initializing slogger: " + err.Error())
			return
		}

		ctx := context.Background()

		log.InfoContext(ctx, "untraced")

		// Capture log
		_ = writer.Close()
		logOutput, _ := io.ReadAll(reader)
		os.Stderr = testStderr
		_ = reader.Close()

		if strings.Contains(string(logOutput), slogger.TraceIDKey) {
			t.Errorf("log should not contain a trace id")
		}

		if strings.Contains(string(logOutput), slogger.SpanIDKey) {
			t.Errorf("log should not contain a span id")
		}
	})

	t.Run("traced", func(t *testing.T) {
		// Setup stderr capture to capture log output.
		testStderr := os.Stderr

		reader, writer, err := os.Pipe()
		if err != nil {
			t.Errorf("creating i/o pipeline: " + err.Error())
			return
		}

		os.Stderr = writer

		// Initialize slogger
		log, err := slogger.NewWithOptions(slogger.Options{
			TracingEnabled: true,
		})
		if err != nil {
			t.Errorf("initializing slogger: " + err.Error())
			return
		}

		// Setup tracing
		tracerProvider := sdktrace.NewTracerProvider()
		tracer := tracerProvider.Tracer("test")
		ctx, span := tracer.Start(context.Background(), "TestOtelHandler")
		defer span.End()

		log.InfoContext(ctx, "traced")

		// Capture log
		_ = writer.Close()
		logOutput, _ := io.ReadAll(reader)
		os.Stderr = testStderr
		_ = reader.Close()

		if !strings.Contains(string(logOutput), slogger.TraceIDKey) {
			t.Errorf("log should contain a trace id")
		}

		if !strings.Contains(string(logOutput), slogger.SpanIDKey) {
			t.Errorf("log should contain a span id")
		}
	})
}
