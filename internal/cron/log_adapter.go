package cron

import (
	"log/slog"
)

type SlogAdapter struct {
	logger *slog.Logger
}

func (s *SlogAdapter) Info(msg string, keysAndValues ...interface{}) {
	s.logger.Info(msg, keysAndValues...)
}

func (s *SlogAdapter) Error(err error, msg string, keysAndValues ...interface{}) {
	allFields := append(keysAndValues, "error", err)
	s.logger.Error(msg, allFields...)
}
