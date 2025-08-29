package processor

import (
	"context"

	"github.com/ipavlov93/universe-demo/product-eventbus-pkg/event"
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
				msg := envelope.Message

				eventPayload, err := event.ParsePayload(msg.Headers.EventType, msg.Payload)
				if err != nil {
					m.lg.Error("failed to parse message payload", zap.Error(err))
				}

				m.lg.Info("Message logged",
					zap.String("message_id", msg.Headers.MessageID),
					zap.String("message_type", msg.Headers.EventType),
					zap.String("producer", msg.Headers.Producer),
					zap.Any("event", eventPayload),
				)
				receiptHandles = append(receiptHandles, envelope.ReceiptHandle)
			}
			out <- receiptHandles
		}
	}
}
