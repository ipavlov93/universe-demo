package factory

import (
	"io"

	logfactory "github.com/ipavlov93/universe-demo/notification-sv/internal/infra/logger/factory"
	processor "github.com/ipavlov93/universe-demo/notification-sv/internal/service/message-logger"
)

func NewMessageLogger(w io.Writer, minLevel string) (*processor.MessageLogger, error) {
	msgLogger, err := logfactory.NewLogger(w, minLevel)
	if err != nil {
		return nil, err
	}

	return processor.NewMessageLogger(msgLogger), nil
}
