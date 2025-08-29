package factory

import (
	"io"
	"os"

	"github.com/ipavlov93/universe-demo/universe-pkg/logger"
	"go.uber.org/zap"

	zapfactory "github.com/ipavlov93/universe-demo/product-sv/internal/infra/logger/zap"
)

func NewAppLoggerOrDefault(level string) (logger.Logger, error) {
	lg, err := newLogger(os.Stdout, level)
	if err != nil {
		return newDefaultAppLogger()
	}
	return lg, nil
}

func newLogger(w io.Writer, level string) (logger.Logger, error) {
	logLevel, err := logger.ParseLevel(level)
	if err != nil {
		return nil, err
	}

	return zapfactory.New(w, logLevel), nil
}

func newDefaultAppLogger() (logger.Logger, error) {
	lg, err := newLoggerInfo()
	if err != nil {
		return nil, err
	}
	return lg, nil
}

func newLoggerInfo() (logger.Logger, error) {
	return newLogger(os.Stdout, zap.InfoLevel.String())
}
