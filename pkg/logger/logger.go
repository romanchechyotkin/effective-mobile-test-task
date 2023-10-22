package logger

import (
	"log/slog"
	"os"
)

func New(out *os.File) *slog.Logger {
	var level = slog.LevelDebug

	if os.Getenv("ENVIRONMENT") == "prod" {
		level = slog.LevelInfo
	}

	return slog.New(slog.NewTextHandler(out, &slog.HandlerOptions{
		Level: level,
	}))
}

func Error(log *slog.Logger, msg string, err error) {
	log.Error(msg, slog.Any("error", err))
}
