package publisher

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/ipavlov93/universe-demo/product-eventbus-pkg/message"

	"github.com/ipavlov93/universe-demo/product-sv/internal/infra/sqs/adapter"
	"github.com/ipavlov93/universe-demo/product-sv/internal/service"
)

type ServiceSQS struct {
	adapter  adapter.Adapter
	queueURL string
}

func NewPublisherSQS(
	ctx context.Context,
	adapter adapter.Adapter,
	queueName string,
) (service.Publisher, error) {
	result, err := adapter.Client().GetQueueUrl(
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
	}, nil
}

func (s *ServiceSQS) PublishJSON(ctx context.Context, data []byte) error {
	_, err := s.adapter.Client().SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(s.queueURL),
		MessageBody: aws.String(string(data)),
	})
	if err != nil {
		return fmt.Errorf("failed to publish message to SQS: %w", err)
	}
	return nil
}

// Publish serializes and sends a single message to the SQS queue.
func (s *ServiceSQS) Publish(ctx context.Context, msg *message.Message) error {
	body, err := msg.EncodeJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return s.PublishJSON(ctx, body)
}
