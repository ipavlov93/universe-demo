package service

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/ipavlov93/universe-demo/product-sv/internal/domain"
)

type ProductRepository interface {
	GetProductByID(ctx context.Context, ID int64) (obj domain.Product, err error)
	CreateProduct(ctx context.Context, product domain.Product) (productID int64, err error)
	DeleteProductByID(ctx context.Context, ID int64) error
}

type Adapter interface {
	Client() *sqs.Client
}
