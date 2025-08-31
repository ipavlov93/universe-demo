package factory

import (
	"io"

	logfactory "github.com/ipavlov93/universe-demo/notification-sv/internal/infra/logger/std"
	processor "github.com/ipavlov93/universe-demo/notification-sv/internal/service/message-logger"
)

func NewMessageLogger(w io.Writer) (*processor.MessageLogger, error) {
	msgLogger := logfactory.NewWriterLogger(w)

	return processor.NewMessageLogger(msgLogger), nil
}
