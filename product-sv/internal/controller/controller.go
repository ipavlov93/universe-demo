package controller

import (
	"context"

	"github.com/ipavlov93/universe-demo/product-sv/internal/domain"
)

type ProductServiceFacade interface {
	GetProductByID(ctx context.Context, ID int64) (obj domain.Product, err error)
	CreateProduct(ctx context.Context, product domain.Product) (int64, error)
	DeleteProduct(ctx context.Context, ID int64) error
}
