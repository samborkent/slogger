/*
slogger is a utility package which wraps a slog.Handler to provide some additional functionality.
The following features are added:
  - Create a handler with standardized config which can be extended through functional options.
  - Automatically enable source for debug logs.
  - Context aware log/slog methods will extract OpenTelemetry trace ID and span ID from context if configured.
  - Set default log/slog logger to enable the use of the global slog methods with the same config.
  - Implements http.Handler to enable changing the log level at runtime.
  - Type safe logging methods for better performance, and enforcing log attribute typing.
*/
package slogger

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"time"
)

// Global log level that can be changed at runtime as the logger is a HTTP handler.
var programLevel = new(slog.LevelVar)

// Wrapper for a slog.Logger that includes some custom config, a HTTP handler,
// and automatically traces context aware logging methods.
type Logger struct {
	*slog.Logger
}

// Create a slog logger with traced JSON handler.
func New(options ...Option) (*Logger, error) {
	// Process options
	config := defaultConfig

	errs := make([]error, 0, len(options))

	for _, option := range options {
		if err := option(&config); err != nil {
			errs = append(errs, err)
		}
	}

	if err := errors.Join(errs...); err != nil {
		return nil, fmt.Errorf("processing options: %w", err)
	}

	// Set global log level.
	programLevel.Set(config.logLevel)

	handlerOptions := &slog.HandlerOptions{
		Level: programLevel,
	}

	// Create handler based on configured tracing type.
	var handler slog.Handler

	switch config.tracingType {
	case tracingTypeOpenTelemetry:
		handler = OtelHandler{
			slog.NewJSONHandler(os.Stderr, handlerOptions),
		}
	case tracingTypeDisabled:
		fallthrough
	default:
		handler = slog.NewJSONHandler(os.Stderr, handlerOptions)
	}

	// Create slog logger
	log := slog.New(handler)

	// Set the handler as default slog.Logger, so the global methods will use the same config.
	slog.SetDefault(log)

	return &Logger{log}, nil
}

// Type safe debug log method.
func (l *Logger) Debug(message string, attributes ...slog.Attr) {
	l.logAttrs(context.Background(), slog.LevelDebug, message, attributes...)
}

// Type safe context-aware debug log method.
func (l *Logger) DebugContext(ctx context.Context, message string, attributes ...slog.Attr) {
	l.logAttrs(ctx, slog.LevelDebug, message, attributes...)
}

// Type safe info log method.
func (l *Logger) Info(message string, attributes ...slog.Attr) {
	l.logAttrs(context.Background(), slog.LevelInfo, message, attributes...)
}

// Type safe context-aware info log method.
func (l *Logger) InfoContext(ctx context.Context, message string, attributes ...slog.Attr) {
	l.logAttrs(ctx, slog.LevelInfo, message, attributes...)
}

// Type safe warn log method.
func (l *Logger) Warn(message string, attributes ...slog.Attr) {
	l.logAttrs(context.Background(), slog.LevelWarn, message, attributes...)
}

// Type safe context-aware warn log method.
func (l *Logger) WarnContext(ctx context.Context, message string, attributes ...slog.Attr) {
	l.logAttrs(ctx, slog.LevelWarn, message, attributes...)
}

// Type safe warn log method.
func (l *Logger) Error(message string, attributes ...slog.Attr) {
	l.logAttrs(context.Background(), slog.LevelError, message, attributes...)
}

// Type safe context-aware warn log method.
func (l *Logger) ErrorContext(ctx context.Context, message string, attributes ...slog.Attr) {
	l.logAttrs(ctx, slog.LevelError, message, attributes...)
}

// Number of stack frames to skip when ketting program counters.
const skipFrames = 3

// Forked from log/slog/logger.go
// logAttrs is like [Logger.log], but for methods that take ...Attr.
func (l *Logger) logAttrs(ctx context.Context, level slog.Level, message string, attributes ...slog.Attr) {
	if !l.Enabled(ctx, level) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(skipFrames, pcs[:])

	record := slog.NewRecord(time.Now(), level, message, pcs[0])

	// Add source information for debug logs
	if level == slog.LevelDebug {
		frames := runtime.CallersFrames([]uintptr{pcs[0]})
		frame, _ := frames.Next()
		source := &slog.Source{
			Function: frame.Function,
			File:     frame.File,
			Line:     frame.Line,
		}
		attributes = append(attributes, slog.Any("source", source))
	}

	record.AddAttrs(attributes...)

	if ctx == nil {
		ctx = context.Background()
	}

	_ = l.Handler().Handle(ctx, record)
}
