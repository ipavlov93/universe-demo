package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ipavlov93/universe-demo/universe-pkg/logger"
	"github.com/jmoiron/sqlx"

	"github.com/ipavlov93/universe-demo/product-sv/internal/domain"
	"github.com/ipavlov93/universe-demo/product-sv/internal/dto/dmodel"
	errs "github.com/ipavlov93/universe-demo/product-sv/internal/error"
	mapper "github.com/ipavlov93/universe-demo/product-sv/internal/mapper/product/dmodel"
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

func (repo *UserRepositoryPostgres) GetUsersTotalCount(ctx context.Context) (int64, error) {
	var count int64
	err := sqlx.GetContext(ctx, repo.dbDriver, &count,
		`SELECT count(*) FROM products`)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		errInfo := fmt.Sprintf("repository.GetProductByID: %v", err)
		return 0, errs.ErrDB.WithReason(errInfo)
	}

	return count, nil
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
			errorInfo = fmt.Sprintf("repository.GetProductByID: user not found for ID=%d", productID)
			return domain.Product{}, errs.ErrProductNotFound.WithReason(errorInfo)
		}
		errorInfo = fmt.Sprintf("repository.GetProductByID: %v", err)
		return domain.Product{}, errs.ErrDB.WithReason(errorInfo)
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
		errorInfo := fmt.Sprintf("repository.CreateProduct: %v", err)

		if len(err.Error()) > 50 && err.Error()[:50] == pqDuplicateErr {
			return 0, errs.ErrProductExists.WithReason(errorInfo)
		}
		return 0, errs.ErrDB.WithReason(errorInfo)
	}
	return productID, nil
}

func (repo *UserRepositoryPostgres) DeleteProductByID(ctx context.Context, productID int64) error {
	err := repo.dbDriver.QueryRowxContext(
		ctx,
		`DELETE FROM products WHERE id = $1`,
		productID,
	)
	if err != nil {
		errorInfo := fmt.Sprintf("repository.DeleteProductByID: %v", err)
		return errs.ErrDB.WithReason(errorInfo)
	}
	return nil
}
