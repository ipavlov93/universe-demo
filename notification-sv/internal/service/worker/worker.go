package worker

import (
	"context"
	"sync"

	"github.com/ipavlov93/universe-demo/notification-sv/internal/message"
	"github.com/ipavlov93/universe-demo/notification-sv/internal/service"
)

func RunWorkers(
	ctx context.Context,
	consumer service.Consumer,
	processor service.MessageProcessor,
	wg *sync.WaitGroup,
) {
	input := make(chan []*message.Envelope, 10)
	out := make(chan []string, 10)

	wg.Add(1)
	go func() {
		processor.Process(ctx, input, out)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		consumer.Subscribe(ctx, input)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer.Acknowledge(ctx, out)
	}()
}
