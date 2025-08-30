package factory

import (
	"io"

	"github.com/ipavlov93/universe-demo/universe-pkg/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewZapLogger constructs logger.
// logger.ZapLogger write logs to the given io.Writer (zapcore.WriteSyncer) using JSON encoding with RFC3339 timestamps.
func NewZapLogger(w io.Writer, minLevel zapcore.Level) *logger.ZapLogger {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.RFC3339TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg),
		zapcore.AddSync(w),
		minLevel,
	)

	options := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	}

	return logger.NewWithCore(core, options...)
}
