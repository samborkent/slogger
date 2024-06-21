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

		log, err := slogger.New(slogger.WithLogLevel("debug"))
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

		log, err := slogger.New(slogger.WithLogLevel("info"))
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

func TestWithSlogLevel(t *testing.T) {
	t.Parallel()

	t.Run("debug", func(t *testing.T) {
		t.Parallel()

		log, err := slogger.New(slogger.WithSlogLevel(slog.LevelDebug))
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

		log, err := slogger.New(slogger.WithSlogLevel(slog.LevelInfo))
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
