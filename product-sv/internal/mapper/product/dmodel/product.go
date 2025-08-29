package dmodel

import (
	"database/sql"

	"github.com/ipavlov93/universe-demo/product-sv/internal/domain"
	"github.com/ipavlov93/universe-demo/product-sv/internal/dto/dmodel"
)

func ProductToProductDto(product domain.Product) dmodel.Product {
	return dmodel.Product{
		ID:   product.ID,
		Name: product.Name,
		Description: func(desc string) sql.NullString {
			return sql.NullString{
				String: desc,
				Valid:  len(desc) > 0,
			}
		}(product.Description),
	}
}

func ProductDtoToProduct(product dmodel.Product) domain.Product {
	return domain.Product{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description.String,
		CreatedAt:   product.CreatedAt,
	}
}
