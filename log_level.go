package slogger

import (
	"errors"
	"log/slog"
	"strings"
)

var ErrLogLevelUnknown = errors.New("log level unknown")

func SlogLevel(logLevel string) (slog.Level, error) {
	switch strings.ToUpper(logLevel) {
	case slog.LevelDebug.String():
		return slog.LevelDebug, nil
	case slog.LevelInfo.String():
		return slog.LevelInfo, nil
	case slog.LevelWarn.String():
		return slog.LevelWarn, nil
	case slog.LevelError.String():
		return slog.LevelError, nil
	default:
		return 0, ErrLogLevelUnknown
	}
}
