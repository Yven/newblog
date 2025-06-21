package logger

import (
	"fmt"
	"log/slog"
	"newblog/internal/global"
	"os"
	"path/filepath"
	"time"
)

func Init(path string, level slog.Level) {
	today := time.Now().Format("2006-01-02")
	os.MkdirAll(path, os.ModePerm)
	logPath := filepath.Join(path, fmt.Sprintf("%s.log", today))

	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		global.Logger.Error("open log file failed", "error", err)
		panic(err)
	}

	handler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		Level: level,
	})

	global.Logger = slog.New(handler)
}
