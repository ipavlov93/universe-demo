package publisher_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/ipavlov93/universe-demo/product-eventbus-pkg/event"
	"github.com/ipavlov93/universe-demo/product-eventbus-pkg/message"
	apperror "github.com/ipavlov93/universe-demo/product-sv/internal/error"
	"github.com/ipavlov93/universe-demo/product-sv/internal/mocks"
	"github.com/ipavlov93/universe-demo/product-sv/internal/service/publisher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewPublisherSQS(t *testing.T) {
	mockAdapter := mocks.NewMockAdapter(t)
	mockSqsClient := mocks.NewMockSQSClientAPI(t)

	ctx := context.Background()
	queueName := "test-queue"

	t.Run("should return a publisher when the queue URL is found", func(t *testing.T) {
		queueURL := "http://sqs.aws.com/123/test-queue"

		mockSqsClient.EXPECT().GetQueueUrl(
			mock.Anything,
			&sqs.GetQueueUrlInput{QueueName: aws.String(queueName)},
		).Return(&sqs.GetQueueUrlOutput{QueueUrl: aws.String(queueURL)}, nil).Once()

		mockAdapter.EXPECT().Client().Return(mockSqsClient).Once()

		// ACT
		publisherSQS, err := publisher.NewPublisherSQS(ctx, mockAdapter, queueName)

		assert.NotNil(t, publisherSQS)
		assert.Nil(t, err)
		assert.IsType(t, &publisher.ServiceSQS{}, publisherSQS)

		mockAdapter.AssertExpectations(t)
		mockSqsClient.AssertExpectations(t)
	})

	t.Run("should return an error when getting the queue URL fails", func(t *testing.T) {
		expectedError := &types.QueueDoesNotExist{}
		mockSqsClient.EXPECT().GetQueueUrl(
			mock.Anything,
			&sqs.GetQueueUrlInput{QueueName: aws.String(queueName)},
		).Return(nil, expectedError).Once()

		mockAdapter.EXPECT().Client().Return(mockSqsClient).Once()

		// ACT
		publisherSQS, err := publisher.NewPublisherSQS(ctx, mockAdapter, queueName)

		assert.Nil(t, publisherSQS)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, apperror.ErrMessageBroker))

		mockAdapter.AssertExpectations(t)
		mockSqsClient.AssertExpectations(t)
	})
}

func TestPublish(t *testing.T) {
	mockAdapter := mocks.NewMockAdapter(t)
	mockSqsClient := mocks.NewMockSQSClientAPI(t)
	queueURL := "http://sqs.aws.com/123/test-queue"
	producer := "test producer"

	ctx := context.Background()

	t.Run("should send a message successfully", func(t *testing.T) {
		productCreatedEvent := event.ProductCreatedEvent{
			ID:   10,
			Name: "Test Product",
		}
		headers := message.NewHeaders(productCreatedEvent.Type(), producer)
		msg, err := message.New(headers, productCreatedEvent)
		require.NoError(t, err)

		mockAdapter.EXPECT().Client().Return(mockSqsClient).Once()
		mockSqsClient.EXPECT().GetQueueUrl(
			mock.Anything,
			mock.Anything,
		).Return(&sqs.GetQueueUrlOutput{QueueUrl: aws.String(queueURL)}, nil).Once()

		mockAdapter.EXPECT().Client().Return(mockSqsClient).Once()
		mockSqsClient.EXPECT().SendMessage(
			mock.Anything,
			mock.AnythingOfType("*sqs.SendMessageInput"),
		).Return(&sqs.SendMessageOutput{}, nil).Once()

		publisherSQS, err := publisher.NewPublisherSQS(ctx, mockAdapter, queueURL)
		require.NoError(t, err)

		// ACT
		err = publisherSQS.Publish(ctx, msg)

		assert.NoError(t, err)
		mockAdapter.AssertExpectations(t)
		mockSqsClient.AssertExpectations(t)
	})

	t.Run("should return an error when sending the message fails", func(t *testing.T) {
		productCreatedEvent := event.ProductCreatedEvent{
			ID:   10,
			Name: "Test Product",
		}
		headers := message.NewHeaders(productCreatedEvent.Type(), producer)
		msg, err := message.New(headers, productCreatedEvent)
		expectedError := apperror.ErrMessageBroker

		mockAdapter.EXPECT().Client().Return(mockSqsClient).Once()
		mockSqsClient.EXPECT().GetQueueUrl(
			mock.Anything,
			mock.Anything,
		).Return(&sqs.GetQueueUrlOutput{QueueUrl: aws.String(queueURL)}, nil).Once()

		mockAdapter.EXPECT().Client().Return(mockSqsClient).Once()
		mockSqsClient.EXPECT().SendMessage(
			mock.Anything,
			mock.AnythingOfType("*sqs.SendMessageInput"),
		).Return(nil, expectedError).Once()

		publisherSQS, err := publisher.NewPublisherSQS(ctx, mockAdapter, queueURL)
		require.NoError(t, err)

		// ACT
		err = publisherSQS.Publish(ctx, msg)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, apperror.ErrMessageBroker))
		assert.Contains(t, err.Error(), "failed to publish message to SQS")

		mockAdapter.AssertExpectations(t)
		mockSqsClient.AssertExpectations(t)
	})
}
