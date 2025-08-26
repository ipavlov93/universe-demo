package zap

import (
	"io"

	"github.com/ipavlov93/universe-demo/universe-pkg/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New logs to the given io.Writer (zapcore.WriteSyncer) using JSON encoding with RFC3339 timestamps.
func New(w io.Writer, minLevel zapcore.Level) *logger.ZapLogger {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.RFC3339TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg),
		zapcore.AddSync(w),
		minLevel,
	)

	return logger.NewWithCore(core)
}
