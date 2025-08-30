package factory

import (
	"io"
	"os"

	"github.com/ipavlov93/universe-demo/universe-pkg/logger"
	"go.uber.org/zap"
)

// NewAppLogger constructs logger using fall back chain listed bellow and ignores errors
func NewAppLogger(w io.Writer, minLevel string) logger.Logger {
	lg, err := newLogger(w, minLevel)
	if err == nil {
		return lg
	}

	lg, err = newStdoutLogger(minLevel)
	if err == nil {
		return lg
	}

	return newStdoutLoggerInfoLevel()
}

func newStdoutLoggerInfoLevel() logger.Logger {
	return NewZapLogger(os.Stdout, zap.InfoLevel)
}

func newStdoutLogger(minLevel string) (logger.Logger, error) {
	return newLogger(os.Stdout, minLevel)
}

func newLogger(w io.Writer, level string) (logger.Logger, error) {
	minLevel, err := logger.ParseLevel(level)
	if err != nil {
		return nil, err
	}

	return NewZapLogger(w, minLevel), nil
}
