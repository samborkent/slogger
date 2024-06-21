package slogger_test

import (
	"context"
	"testing"

	"github.com/samborkent/slogger"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// TODO 2024-06-20 Sam Borkent: extend and finish tests.

// func TestHandle(t *testing.T) {
// 	// Setup stderr capture to capture log output.
// 	testStderr := os.Stderr

// 	reader, writer, err := os.Pipe()
// 	if err != nil {
// 		t.Errorf("creating os i/o pipeline: " + err.Error())
// 		return
// 	}

// 	os.Stderr = writer

// 	wg := new(sync.WaitGroup)

// 	wg.Add(1)
// 	t.Run("untraced", func(t *testing.T) {
// 		t.Parallel()

// 		// ctx := context.Background()
// 		log := slogger.New()

// 		log.Info("test")
// 		wg.Done()
// 	})

// 	wg.Wait()

// 	_ = writer.Close()

// 	logOutput, _ := io.ReadAll(reader)
// }

func TestLogger(t *testing.T) {
	log, err := slogger.New()
	if err != nil {
		t.Errorf("initializing logger: " + err.Error())
		return
	}

	ctx := context.Background()

	log.InfoContext(ctx, "TEST")

	tracerProvider := sdktrace.NewTracerProvider()
	tracer := tracerProvider.Tracer("test")

	ctx, span := tracer.Start(ctx, "TestHandle")
	defer span.End()

	log.InfoContext(ctx, "TEST - traced logger")
}
