package processor

import (
	"context"

	"github.com/ipavlov93/universe-demo/product-eventbus-pkg/event"
	msgpkg "github.com/ipavlov93/universe-demo/product-eventbus-pkg/message"
	"github.com/ipavlov93/universe-demo/universe-pkg/logger"
	"go.uber.org/zap"

	"github.com/ipavlov93/universe-demo/notification-sv/internal/message"
)

type MessageLogger struct {
	lg logger.Logger
}

func NewMessageLogger(lg logger.Logger) *MessageLogger {
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
		m.lg.Sync()
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
				m.ParseAndLogMessage(envelope.Message)

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

	m.lg.Info("Message logged",
		zap.String("message_id", msg.Headers.MessageID),
		zap.String("event_type", msg.Headers.EventType),
		zap.String("producer", msg.Headers.Producer),
		zap.Any("event", eventPayload),
	)
}
