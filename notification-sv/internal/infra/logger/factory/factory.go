package factory

import (
	"io"
	"os"

	"github.com/ipavlov93/universe-demo/universe-pkg/logger"
	"go.uber.org/zap"
)

func NewStdoutLogger(minLevel string) (logger.Logger, error) {
	return NewLogger(os.Stdout, minLevel)
}

func NewStdoutLoggerInfoLevel() (logger.Logger, error) {
	return NewLogger(os.Stdout, zap.InfoLevel.String())
}

func NewLogger(w io.Writer, level string) (logger.Logger, error) {
	minLevel, err := logger.ParseLevel(level)
	if err != nil {
		return nil, err
	}

	return NewZapLogger(w, minLevel), nil
}
