package slogger_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/samborkent/slogger"
)

func TestWithLogLevel(t *testing.T) {
	t.Parallel()

	t.Run("debug", func(t *testing.T) {
		t.Parallel()

		log, err := slogger.NewWithOptions(slogger.Options{
			LogLevel: slog.LevelDebug,
		})
		if err != nil {
			t.Errorf("initializing logger: " + err.Error())
			return
		}

		ctx := context.Background()

		if !log.Enabled(ctx, slog.LevelDebug) {
			t.Error("debug should be enabled")
		}

		if !log.Enabled(ctx, slog.LevelInfo) {
			t.Error("info should be enabled")
		}
	})

	t.Run("info", func(t *testing.T) {
		t.Parallel()

		log, err := slogger.NewWithOptions(slogger.Options{
			LogLevel: slog.LevelInfo,
		})
		if err != nil {
			t.Errorf("initializing logger: " + err.Error())
			return
		}

		ctx := context.Background()

		if !log.Enabled(ctx, slog.LevelInfo) {
			t.Error("info should be enabled")
		}

		if !log.Enabled(ctx, slog.LevelWarn) {
			t.Error("warn should be enabled")
		}

		if log.Enabled(ctx, slog.LevelDebug) {
			t.Error("debug should be disabled")
		}
	})
}
