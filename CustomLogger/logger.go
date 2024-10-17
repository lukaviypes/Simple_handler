package logger

import (
	"log/slog"
	"os"
)

func New() *slog.Logger {
	logPath := os.Getenv("LOG_PATH")
	if logPath == "" {
		logPath = "logs.txt"
	}

	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err.Error())
	}

	return slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{}))
}
