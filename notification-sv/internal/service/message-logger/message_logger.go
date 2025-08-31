package processor

import (
	"context"
	"fmt"

	"github.com/ipavlov93/universe-demo/notification-sv/internal/message"
	"github.com/ipavlov93/universe-demo/product-eventbus-pkg/event"
	msgpkg "github.com/ipavlov93/universe-demo/product-eventbus-pkg/message"
)

type Logger interface {
	Log(msg string)
}

type MessageLogger struct {
	lg Logger
}

func NewMessageLogger(lg Logger) *MessageLogger {
	return &MessageLogger{
		lg: lg,
	}
}

// Process starts long-running process of logging messages.
// Process will end after one of following conditions:
// 1. input channel is closed.
// 2. ctx is done.
func (m *MessageLogger) Process(
	ctx context.Context,
	input <-chan []*message.Envelope,
	out chan<- []string,
) {
	defer func() {
		close(out)
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case envelopes, ok := <-input:
			if !ok {
				return
			}

			m.process(envelopes, out)
		}
	}
}

func (m *MessageLogger) process(envelopes []*message.Envelope, out chan<- []string) {
	var receiptHandles []string
	for _, envelope := range envelopes {
		if envelope == nil {
			continue
		}

		// todo: add error sending to separate channel
		eventPayload, err := event.ParsePayload(envelope.Message.Headers.EventType, envelope.Message.Payload)
		if err != nil {
			m.lg.Log(err.Error())
		}

		m.logMessage(envelope.Message, eventPayload)

		receiptHandles = append(receiptHandles, envelope.ReceiptHandle)
	}
	out <- receiptHandles
}

func (m *MessageLogger) logMessage(msg *msgpkg.Message, eventPayload any) {
	if msg == nil {
		return
	}

	logRow := fmt.Sprintf(
		"Message logged | message_id=%s | event_type=%s | producer=%s | event=%+v",
		msg.Headers.MessageID,
		msg.Headers.EventType,
		msg.Headers.Producer,
		eventPayload,
	)

	m.lg.Log(logRow)
	return
}
