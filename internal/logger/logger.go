package logger

import (
	"log/slog"
	"os"
)

func SetupLogger(silent, verbose bool) *slog.Logger {
	var level slog.Level

	switch {
	case silent:
		level = slog.LevelError
	case verbose:
		level = slog.LevelDebug
	default:
		level = slog.LevelInfo
	}

	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	})

	return slog.New(handler)
}
