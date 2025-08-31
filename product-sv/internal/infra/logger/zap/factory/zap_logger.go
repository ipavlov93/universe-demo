package factory

import (
	"io"
	"time"

	"github.com/ipavlov93/universe-demo/universe-pkg/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewZapLogger constructs logger.
// logger.ZapLogger write logs to the given io.Writer (zapcore.WriteSyncer) with RFC3339 timestamps.
func NewZapLogger(w io.Writer, minLevel zapcore.Level) *logger.ZapLogger {
	cfg := zap.NewProductionEncoderConfig()

	cfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format(time.RFC3339))
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg),
		zapcore.AddSync(w),
		minLevel,
	)

	options := []zap.Option{
		zap.AddCaller(),
	}

	return logger.NewWithCore(core, options...)
}
