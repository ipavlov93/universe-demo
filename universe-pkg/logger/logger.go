package logger

import (
	"go.uber.org/zap"
)

// Logger usage makes it possible swap logger in the future.
type Logger interface {
	// Sync flushing any buffered log entries.
	// Applications should take care to call Sync before exiting.
	Sync() error
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}
