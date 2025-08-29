package sqsclient

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	sqsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func NewClientSQS(
	ctx context.Context,
	sessionKey string,
	baseURL string,
	defaultRegion string,
	accessKeyID string,
	accessKeySecret string,
) (*sqs.Client, error) {
	config, err := sqsconfig.LoadDefaultConfig(ctx,
		sqsconfig.WithDefaultRegion(defaultRegion),
		sqsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKeyID, accessKeySecret, sessionKey),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load SQS config: %w", err)
	}

	client := sqs.NewFromConfig(config, func(o *sqs.Options) {
		o.BaseEndpoint = aws.String(baseURL)
	})

	return client, nil
}
