package adapter

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqspkg "github.com/ipavlov93/universe-demo/universe-pkg/sqs"

	"github.com/ipavlov93/universe-demo/notification-sv/internal/config"
)

type Adapter struct {
	client *sqs.Client
}

func (a *Adapter) Client() *sqs.Client {
	return a.client
}

func NewAdapter(
	ctx context.Context,
	cfg config.LocalStackConfig,
	scheme string,
	sessionKey string,
) (*Adapter, error) {
	localStackBaseURL := baseURL(scheme, cfg)

	client, err := sqspkg.NewClientSQS(
		ctx, sessionKey, localStackBaseURL,
		cfg.DefaultRegion, cfg.AccessKeyID, cfg.SecretAccessKey,
	)
	if err != nil {
		return nil, err
	}

	return &Adapter{
		client: client,
	}, nil
}

func baseURL(scheme string, cfg config.LocalStackConfig) string {
	return fmt.Sprintf(
		"%s://%s:%d",
		scheme, cfg.Host, cfg.Port,
	)
}
