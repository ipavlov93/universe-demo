package product_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ipavlov93/universe-demo/product-sv/internal/domain"
	apperror "github.com/ipavlov93/universe-demo/product-sv/internal/error"
	"github.com/ipavlov93/universe-demo/product-sv/internal/mocks"
	"github.com/ipavlov93/universe-demo/product-sv/internal/service/product"
)

func TestGetProductByID(t *testing.T) {
	mockRepo := mocks.NewMockProductRepository(t)
	s := product.NewProductService(mockRepo)
	ctx := context.Background()

	t.Run("should return the product when the ID is valid", func(t *testing.T) {
		expectedProduct := domain.Product{ID: 10, Name: "Test Product"}
		mockRepo.EXPECT().GetProductByID(mock.Anything, int64(10)).Return(expectedProduct, nil).Once()

		result, err := s.GetProductByID(ctx, 10)

		assert.Nil(t, err)
		assert.Equal(t, expectedProduct, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return ErrInvalidArgument when the ID is invalid", func(t *testing.T) {
		result, err := s.GetProductByID(ctx, 0)

		assert.Equal(t, domain.Product{}, result)
		assert.Equal(t, apperror.ErrInvalidArgument, err)
		mockRepo.AssertNotCalled(t, "GetProductByID")
	})

	t.Run("should return an error when the repository fails", func(t *testing.T) {
		expectedError := errors.New("database connection failed")
		mockRepo.EXPECT().GetProductByID(mock.Anything, int64(11)).Return(domain.Product{}, expectedError).Once()

		result, err := s.GetProductByID(ctx, 11)

		assert.Equal(t, domain.Product{}, result)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestCreateProduct(t *testing.T) {
	mockRepo := mocks.NewMockProductRepository(t)
	s := product.NewProductService(mockRepo)
	ctx := context.Background()

	t.Run("should create the product when it is valid", func(t *testing.T) {
		validProduct := domain.Product{Name: "New Product"}
		expectedID := int64(100)
		mockRepo.EXPECT().CreateProduct(mock.Anything, validProduct).Return(expectedID, nil).Once()

		id, err := s.CreateProduct(ctx, validProduct)

		assert.Nil(t, err)
		assert.Equal(t, expectedID, id)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return ErrInvalidArgument when the product is invalid", func(t *testing.T) {
		invalidProduct := domain.Product{Name: ""}

		id, err := s.CreateProduct(ctx, invalidProduct)

		assert.Equal(t, int64(0), id)
		assert.Equal(t, apperror.ErrInvalidArgument, err)
		mockRepo.AssertNotCalled(t, "CreateProduct")
	})

	t.Run("should return an error when the repository fails", func(t *testing.T) {
		validProduct := domain.Product{Name: "Valid Product"}
		expectedError := errors.New("insertion failed")
		mockRepo.EXPECT().CreateProduct(mock.Anything, validProduct).Return(int64(0), expectedError).Once()

		id, err := s.CreateProduct(ctx, validProduct)

		assert.Equal(t, int64(0), id)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteProduct(t *testing.T) {
	mockRepo := mocks.NewMockProductRepository(t)
	s := product.NewProductService(mockRepo)
	ctx := context.Background()

	t.Run("should delete the product when it exists", func(t *testing.T) {
		productID := int64(20)
		existingProduct := domain.Product{ID: productID}
		mockRepo.EXPECT().GetProductByID(mock.Anything, productID).Return(existingProduct, nil).Once()
		mockRepo.EXPECT().DeleteProductByID(mock.Anything, productID).Return(nil).Once()

		err := s.DeleteProduct(ctx, productID)

		assert.Nil(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return an error when the ID is invalid", func(t *testing.T) {
		err := s.DeleteProduct(ctx, 0)

		assert.Equal(t, apperror.ErrInvalidArgument, err)
		mockRepo.AssertNotCalled(t, "GetProductByID")
		mockRepo.AssertNotCalled(t, "DeleteProductByID")
	})

	t.Run("should not return an error when the product is not found", func(t *testing.T) {
		productID := int64(21)
		mockRepo.EXPECT().GetProductByID(mock.Anything, productID).Return(domain.Product{}, apperror.ErrProductNotFound).Once()

		err := s.DeleteProduct(ctx, productID)

		assert.Nil(t, err)
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "DeleteProductByID")
	})

	t.Run("should return an error when GetProductByID fails with a different error", func(t *testing.T) {
		productID := int64(22)
		expectedError := errors.New("unexpected database error")
		mockRepo.EXPECT().GetProductByID(mock.Anything, productID).Return(domain.Product{}, expectedError).Once()

		err := s.DeleteProduct(ctx, productID)

		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "DeleteProductByID")
	})

	t.Run("should return an error when DeleteProductByID fails", func(t *testing.T) {
		productID := int64(23)
		existingProduct := domain.Product{ID: productID}
		expectedError := errors.New("deletion failed")
		mockRepo.EXPECT().GetProductByID(mock.Anything, productID).Return(existingProduct, nil).Once()
		mockRepo.EXPECT().DeleteProductByID(mock.Anything, productID).Return(expectedError).Once()

		err := s.DeleteProduct(ctx, productID)

		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}
