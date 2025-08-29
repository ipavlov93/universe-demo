package consumer

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	msgpkg "github.com/ipavlov93/universe-demo/product-eventbus-pkg/message"
	"github.com/ipavlov93/universe-demo/universe-pkg/logger"
	"go.uber.org/zap"

	"github.com/ipavlov93/universe-demo/notification-sv/internal/infra/sqs/adapter"
	"github.com/ipavlov93/universe-demo/notification-sv/internal/message"
	"github.com/ipavlov93/universe-demo/notification-sv/internal/service"
)

type ServiceSQS struct {
	adapter  adapter.Adapter
	queueURL string
	lg       logger.Logger
}

func NewConsumerSQS(
	ctx context.Context,
	adapter adapter.Adapter,
	queueName string,
	lg logger.Logger,
) (service.Consumer, error) {
	sqsClient := adapter.Client()
	result, err := sqsClient.GetQueueUrl(
		ctx,
		&sqs.GetQueueUrlInput{
			QueueName: aws.String(queueName),
		})
	if err != nil {
		return nil, fmt.Errorf("failed to get queue URL for %s: %w", queueName, err)
	}

	return &ServiceSQS{
		adapter:  adapter,
		queueURL: *result.QueueUrl,
		lg:       lg,
	}, nil
}

// Subscribe starts a long-running process to consume messages from the queue.
func (s *ServiceSQS) Subscribe(ctx context.Context, out chan<- []*message.Envelope) {
	sqsClient := s.adapter.Client()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			result, err := sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(s.queueURL),
				MaxNumberOfMessages: 10,
				WaitTimeSeconds:     10,
			})
			if err != nil {
				s.lg.Error("failed to receive message", zap.Error(err))
				continue
			}

			if len(result.Messages) == 0 {
				continue
			}

			messages, err := deserializeMessages(result.Messages)
			if err != nil {
				s.lg.Error("failed to deserialize messages", zap.Error(err))
				continue
			}

			out <- messages
		}
	}
}

// Acknowledge starts a long-running process to delete messages from the queue.
func (s *ServiceSQS) Acknowledge(parentCtx context.Context, out <-chan []string) {
	for {
		select {
		case <-parentCtx.Done():
			return
		case receiptHandles, ok := <-out:
			if !ok {
				return
			}
			err := s.deleteBatch(parentCtx, receiptHandles)
			if err != nil {
				s.lg.Error("failed to delete messages", zap.Error(err))
			}
		}
	}
}

func (s *ServiceSQS) deleteBatch(ctx context.Context, receiptHandle []string) error {
	sqsClient := s.adapter.Client()

	_, err := sqsClient.DeleteMessageBatch(ctx, &sqs.DeleteMessageBatchInput{
		QueueUrl: aws.String(s.queueURL),
		Entries:  deleteReceiptHandles(receiptHandle),
	})
	return err
}

func deleteReceiptHandles(receiptHandles []string) []types.DeleteMessageBatchRequestEntry {
	result := make([]types.DeleteMessageBatchRequestEntry, 0, len(receiptHandles))
	for i, receiptHandle := range receiptHandles {
		id := fmt.Sprintf("%d", i)
		result = append(result, types.DeleteMessageBatchRequestEntry{
			Id:            aws.String(id),
			ReceiptHandle: aws.String(receiptHandle),
		})
	}
	return result
}

func deserializeMessages(messages []types.Message) ([]*message.Envelope, error) {
	result := make([]*message.Envelope, 0)
	for _, msg := range messages {
		env := message.Envelope{
			Message:       &msgpkg.Message{},
			ReceiptHandle: *msg.ReceiptHandle,
		}

		if err := env.Message.DecodeJSON([]byte(*msg.Body)); err != nil {
			return nil, err
		}

		result = append(result, &env)
	}
	return result, nil
}
