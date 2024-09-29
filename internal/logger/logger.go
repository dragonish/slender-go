package logger

import (
	"io"
	"log/slog"
	"os"
	"slender/internal/model"
	"strings"

	goLogger "github.com/donnie4w/go-logger/logger"
)

// instance is default logger instance.
var instance *slog.Logger

// programLevel set log lowest level, Info by default.
var programLevel = new(slog.LevelVar)

func init() {
	logFile, err := goLogger.NewLogger().SetRollingDaily(model.LOGS_DIR, model.LOG_FILENAME)
	if err != nil {
		slog.New(slog.NewTextHandler(os.Stdout, nil)).Error("unable to create log file", "error", err)
		os.Exit(1)
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	h := slog.NewTextHandler(multiWriter, &slog.HandlerOptions{Level: programLevel})
	instance = slog.New(h)

	// Set as default logger.
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
