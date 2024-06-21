package slogger

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

type Option func(*config) error

// Add trace ID and span ID attributes to context-aware log methods with Elastic APM traced context.
func WithAPMTracing() Option {
	return func(c *config) error {
		c.tracingType = tracingTypeElasticAPM
		return nil
	}
}

// Add trace ID and span ID attributes to context-aware log methods with OpenTelemtery traced context.
func WithOtelTracing() Option {
	return func(c *config) error {
		c.tracingType = tracingTypeOpenTelemetry
		return nil
	}
}

// Specify log level as string.
func WithLogLevel(logLevel string) Option {
	return func(c *config) error {
		slogLevel, err := stringToSlogLevel(logLevel)
		if err != nil {
			return fmt.Errorf("WithLogLevel: %w", err)
		}

		c.logLevel = slogLevel

		return nil
	}
}

// Specify log level as slog.Level.
func WithSlogLevel(logLevel slog.Level) Option {
	return func(c *config) error {
		c.logLevel = logLevel
		return nil
	}
}

const defaultLogLevel = slog.LevelInfo

var defaultConfig config = config{
	logLevel:    defaultLogLevel,
	tracingType: tracingTypeDisabled,
}

type tracingType string

const (
	tracingTypeDisabled      = "disabled"
	tracingTypeElasticAPM    = "apm"
	tracingTypeOpenTelemetry = "otel"
)

type config struct {
	logLevel    slog.Level
	tracingType tracingType
}

var errLogLevelUnknown = errors.New("log level unknown")

func stringToSlogLevel(logLevel string) (slog.Level, error) {
	switch strings.ToLower(logLevel) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, errLogLevelUnknown
	}
}
