package worker

import (
	"context"
	"sync"

	"github.com/ipavlov93/universe-demo/notification-sv/internal/message"
	"github.com/ipavlov93/universe-demo/notification-sv/internal/service"
)

type Service struct {
	consumer         service.Consumer
	messageProcessor service.MessageProcessor
	chanBufferSize   int
}

func NewWorkerService(
	consumer service.Consumer,
	messageProcessor service.MessageProcessor,
	chanBufferSize int,
) Service {
	return Service{
		consumer:         consumer,
		messageProcessor: messageProcessor,
		chanBufferSize:   chanBufferSize,
	}
}

func (s Service) RunWorkers(ctx context.Context, wg *sync.WaitGroup) {
	input := make(chan []*message.Envelope, s.chanBufferSize)
	out := make(chan []string, s.chanBufferSize)

	wg.Add(1)
	go func() {
		s.messageProcessor.Process(ctx, input, out)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		s.consumer.Subscribe(ctx, input)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.consumer.Acknowledge(ctx, out)
	}()
}
