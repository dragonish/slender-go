package logger

import (
	"log/slog"
	"os"
	"strings"
)

// instance is default logger instance.
var instance *slog.Logger

// programLevel set log lowest level, Info by default.
var programLevel = new(slog.LevelVar)

func init() {
	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel})
	instance = slog.New(h)
	slog.SetDefault(instance)
}

// SetLevel uses level to set the lowest log output level.
//
// Log level optional: "debug" | "info" | "warn" | "error".
func SetLevel(level string) {
	switch strings.ToLower(level) {
	case "debug":
		programLevel.Set(slog.LevelDebug)
	case "info":
		programLevel.Set(slog.LevelInfo)
	case "warn":
		programLevel.Set(slog.LevelWarn)
	case "error":
		programLevel.Set(slog.LevelError)
	}
}
