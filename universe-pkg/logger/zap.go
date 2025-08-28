package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapLogger is wrapper around underlying zap.Logger.
type ZapLogger struct {
	logger *zap.Logger
}

// Sync flushing any buffered log entries.
// Applications should take care to call Sync before exiting.
func (z *ZapLogger) Sync() error { return z.logger.Sync() }

func (z *ZapLogger) Info(msg string, fields ...zap.Field)  { z.logger.Info(msg, fields...) }
func (z *ZapLogger) Warn(msg string, fields ...zap.Field)  { z.logger.Warn(msg, fields...) }
func (z *ZapLogger) Debug(msg string, fields ...zap.Field) { z.logger.Debug(msg, fields...) }
func (z *ZapLogger) Error(msg string, fields ...zap.Field) { z.logger.Error(msg, fields...) }
func (z *ZapLogger) Fatal(msg string, fields ...zap.Field) { z.logger.Fatal(msg, fields...) }

func NewWithCore(core zapcore.Core, options ...zap.Option) *ZapLogger {
	options = append(options, zap.AddCallerSkip(1))

	return &ZapLogger{
		logger: zap.New(core, options...),
	}
}
