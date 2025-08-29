package smodel

import (
	"github.com/ipavlov93/universe-demo/product-sv/internal/domain"
	"github.com/ipavlov93/universe-demo/product-sv/internal/dto/smodel"
)

func ProductToProductDto(product domain.Product) smodel.Product {
	return smodel.Product{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		CreatedAt:   product.CreatedAt,
	}
}

func ProductDtoToProduct(product smodel.Product) domain.Product {
	return domain.Product{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
	}
}
