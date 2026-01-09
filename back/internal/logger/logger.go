package logger

import (
	"log/slog"
	"os"
)

func New() *slog.Logger {
	// Structured JSON logs with levels; production-friendly
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}
