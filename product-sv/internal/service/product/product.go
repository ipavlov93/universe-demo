package product

import (
	"context"
	"errors"

	"github.com/ipavlov93/universe-demo/product-sv/internal/domain"
	apperror "github.com/ipavlov93/universe-demo/product-sv/internal/error"
	"github.com/ipavlov93/universe-demo/product-sv/internal/service"
)

type Service struct {
	repo service.ProductRepository
}

func NewProductService(repo service.ProductRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetProductByID(ctx context.Context, productID int64) (obj domain.Product, err error) {
	if productID < 1 {
		return domain.Product{}, apperror.ErrInvalidArgument
	}

	return s.repo.GetProductByID(ctx, productID)
}

func (s *Service) CreateProduct(ctx context.Context, product domain.Product) (int64, error) {
	if !product.Valid() {
		return 0, apperror.ErrInvalidArgument
	}

	return s.repo.CreateProduct(ctx, product)
}

func (s *Service) DeleteProduct(ctx context.Context, productID int64) error {
	if productID < 1 {
		return apperror.ErrInvalidArgument
	}

	_, err := s.repo.GetProductByID(ctx, productID)
	if err != nil {
		if errors.Is(err, apperror.ErrProductNotFound) {
			return nil
		}
		return err
	}

	return s.repo.DeleteProductByID(ctx, productID)
}
