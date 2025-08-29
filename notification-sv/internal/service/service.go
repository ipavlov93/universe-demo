package service

import (
	"context"

	"github.com/ipavlov93/universe-demo/notification-sv/internal/message"
)

// MessageProcessor represents a component that processes incoming messages.
type MessageProcessor interface {
	Process(ctx context.Context, input <-chan []*message.Envelope, outputCh chan<- []string)
}

type Consumer interface {
	// Subscribe starts reading from the external broker and sending messages to out channel.
	Subscribe(ctx context.Context, out chan<- []*message.Envelope)
	// Acknowledge starts reading from input channel and sending delete request to external broker.
	Acknowledge(parentCtx context.Context, input <-chan []string)
}
