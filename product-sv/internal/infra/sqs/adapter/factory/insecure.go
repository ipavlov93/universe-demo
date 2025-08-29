package factory

import (
	"context"

	"github.com/ipavlov93/universe-demo/product-sv/internal/config"
	"github.com/ipavlov93/universe-demo/product-sv/internal/infra/sqs/adapter"
)

func NewInsecureAdapter(ctx context.Context, cfg config.LocalStackConfig) (*adapter.Adapter, error) {
	return adapter.NewAdapter(ctx, cfg, "http", "")
}
