package service

import (
	"context"

	"github.com/ipavlov93/universe-demo/product-eventbus-pkg/message"

	"github.com/ipavlov93/universe-demo/product-sv/internal/domain"
)

type Facade interface {
	GetProductByID(ctx context.Context, ID int64) (obj domain.Product, err error)
	CreateProduct(ctx context.Context, product domain.Product) (int64, error)
	DeleteProduct(ctx context.Context, ID int64) error
}

type ProductService interface {
	GetProductByID(ctx context.Context, ID int64) (obj domain.Product, err error)
	CreateProduct(ctx context.Context, product domain.Product) (int64, error)
	DeleteProduct(ctx context.Context, productID int64) error
}

type ProductRepository interface {
	GetProductByID(ctx context.Context, ID int64) (obj domain.Product, err error)
	CreateProduct(ctx context.Context, product domain.Product) (productID int64, err error)
	DeleteProductByID(ctx context.Context, ID int64) error
}

type Publisher interface {
	Publish(ctx context.Context, msg *message.Message) error
}

type PrometheusService interface {
	IncProductsCreated()
	IncProductsDeleted()
}
