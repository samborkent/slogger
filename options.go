package slogger

import (
	"log/slog"
)

type Options struct {
	LogLevel       slog.Level
	TracingEnabled bool
}
