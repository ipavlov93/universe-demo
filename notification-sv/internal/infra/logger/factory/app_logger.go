package factory

import (
	"io"

	"github.com/ipavlov93/universe-demo/universe-pkg/logger"
	"go.uber.org/zap"
)

// NewAppLogger constructs logger using fall back chain listed bellow and ignores errors
func NewAppLogger(w io.Writer, minLevel string) logger.Logger {
	lg, err := NewLogger(w, minLevel)
	if err == nil {
		return lg
	}

	lg, err = NewStdoutLogger(minLevel)
	if err == nil {
		return lg
	}

	lg, err = NewStdoutLoggerInfoLevel()
	if err == nil {
		return lg
	}

	return zap.NewNop()
}
