package publisher

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/ipavlov93/universe-demo/product-eventbus-pkg/message"
	apperror "github.com/ipavlov93/universe-demo/product-sv/internal/error"
	"github.com/ipavlov93/universe-demo/product-sv/internal/service/facade"
)

type Adapter interface {
	Client() SQSClientAPI
}

type SQSClientAPI interface {
	GetQueueUrl(ctx context.Context, params *sqs.GetQueueUrlInput, optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)
	SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
}

type ServiceSQS struct {
	adapter  Adapter
	queueURL string
}

func NewPublisherSQS(
	ctx context.Context,
	adapter Adapter,
	queueName string,
) (facade.Publisher, error) {
	result, err := adapter.Client().GetQueueUrl(
		ctx,
		&sqs.GetQueueUrlInput{
			QueueName: aws.String(queueName),
		})
	if err != nil {
		return nil, apperror.ErrMessageBroker.WithReason(
			fmt.Sprintf("failed to get queue URL for %s: %v", queueName, err))
	}

	return &ServiceSQS{
		adapter:  adapter,
		queueURL: *result.QueueUrl,
	}, nil
}

// Publish serializes and sends a single message to the SQS queue.
func (s *ServiceSQS) Publish(ctx context.Context, msg *message.Message) error {
	body, err := msg.EncodeJSON()
	if err != nil {
		return apperror.ErrMessageBroker.WithReason(
			fmt.Sprintf("failed to marshal message: %v", err))
	}

	return s.publish(ctx, body)
}

func (s *ServiceSQS) publish(ctx context.Context, data []byte) error {
	_, err := s.adapter.Client().SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(s.queueURL),
		MessageBody: aws.String(string(data)),
	})
	if err != nil {
		return apperror.ErrMessageBroker.WithReason(
			fmt.Sprintf("failed to publish message to SQS: %v", err))
	}
	return nil
}
