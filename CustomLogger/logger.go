package logger

import (
	"log/slog"
	"net/http"
	"os"
)

type Logmid struct {
	Logger slog.Logger
	Status int
	Err    error
}

func New() *Logmid {
	logPath := os.Getenv("LOG_PATH")
	if logPath == "" {
		logPath = "logs.txt"
	}

	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err.Error())
	}

	return &Logmid{
		Logger: *slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{})),
		Status: http.StatusOK,
	}
}
