package logger

import (
	"log"
	"log/slog"
	"os"
)

type LogLevel string

const (
	Local LogLevel = "local"
	Dev   LogLevel = "dev"
	Prod  LogLevel = "prod"
)

func SetupLogger(logLevel LogLevel, logFilePath string) *slog.Logger {
	//file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//if err != nil {
	//	log.Fatalf("Failed to open log file: %v", err)
	//}

	var level slog.Level
	switch logLevel {
	case Local:
		level = slog.LevelDebug
	case Dev:
		level = slog.LevelInfo
	case Prod:
		level = slog.LevelError
	default:
		log.Fatalf("Invalid log level: %s. Must be one of 'local', 'dev', or 'prod'.", logLevel)
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	return slog.New(handler)
}
