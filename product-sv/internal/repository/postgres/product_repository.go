package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ipavlov93/universe-demo/product-sv/internal/domain"
	"github.com/ipavlov93/universe-demo/product-sv/internal/dto/dmodel"
	errs "github.com/ipavlov93/universe-demo/product-sv/internal/error"
	mapper "github.com/ipavlov93/universe-demo/product-sv/internal/mapper/product/dmodel"
	"github.com/ipavlov93/universe-demo/universe-pkg/logger"
	"github.com/jmoiron/sqlx"
)

type UserRepositoryPostgres struct {
	dbDriver sqlx.ExtContext
	logger   logger.Logger
}

func NewUserRepository(dbDriver sqlx.ExtContext) *UserRepositoryPostgres {
	return &UserRepositoryPostgres{
		dbDriver: dbDriver,
	}
}

// GetProductByID will return errs.ErrProductNotFound if no matching record exists.
func (repo *UserRepositoryPostgres) GetProductByID(ctx context.Context, productID int64) (obj domain.Product, err error) {
	var errorInfo string
	var productDto dmodel.Product

	err = sqlx.GetContext(ctx, repo.dbDriver, &productDto,
		`SELECT * FROM products
				WHERE id = $1`, productID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errorInfo = fmt.Sprintf("product ID=%d not found", productID)
			return domain.Product{}, errs.ErrProductNotFound.WithReason(errorInfo)
		}
		return domain.Product{}, errs.ErrDB.WithReason(err.Error())
	}
	return mapper.ProductDtoToProduct(productDto), nil
}

// CreateProduct will return errs.ErrProductExists if the product already exists.
func (repo *UserRepositoryPostgres) CreateProduct(ctx context.Context, product domain.Product) (productID int64, err error) {
	err = repo.dbDriver.QueryRowxContext(
		ctx,
		`INSERT INTO products (name, description) VALUES ($1, $2) RETURNING id`,
		product.Name, product.Description,
	).Scan(&productID)
	if err != nil {
		if len(err.Error()) > 50 && err.Error()[:50] == pqDuplicateErr {
			return 0, errs.ErrProductExists
		}
		return 0, errs.ErrDB.WithReason(err.Error())
	}
	return productID, nil
}

func (repo *UserRepositoryPostgres) DeleteProductByID(ctx context.Context, productID int64) error {
	_, err := repo.dbDriver.ExecContext(
		ctx,
		`DELETE FROM products WHERE id = $1`,
		productID,
	)
	if err != nil {
		return errs.ErrDB.WithReason(err.Error())
	}
	return nil
}
