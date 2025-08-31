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

// Process starts to log messages in a separate goroutine.
// It respects context cancellation (e.g., via <-ctx.Done()) and wait group by design.
// Notice: actual logs format is different from JSON.
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

			var receiptHandles []string
			for _, envelope := range envelopes {
				if envelope == nil {
					continue
				}

				// todo: add error sending to separate channel
				err := m.ParseAndLogMessage(envelope.Message)
				if err != nil {
					m.lg.Log(err.Error())
				}

				receiptHandles = append(receiptHandles, envelope.ReceiptHandle)
			}
			out <- receiptHandles
		}
	}
}

func (m *MessageLogger) ParseAndLogMessage(msg *msgpkg.Message) error {
	if msg == nil {
		return nil
	}

	eventPayload, err := event.ParsePayload(msg.Headers.EventType, msg.Payload)
	if err != nil {
		return err
	}

	m.logMessage(msg, eventPayload)
	return nil
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
