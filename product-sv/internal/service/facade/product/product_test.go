package product_test

import (
	"context"
	"testing"

	"github.com/ipavlov93/universe-demo/product-sv/internal/domain"
	apperror "github.com/ipavlov93/universe-demo/product-sv/internal/error"
	"github.com/ipavlov93/universe-demo/product-sv/internal/mocks"
	"github.com/ipavlov93/universe-demo/product-sv/internal/service/facade/product"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
)

func TestGetProductByID(t *testing.T) {
	mockProductService := mocks.NewMockProductService(t)
	mockPublisher := mocks.NewMockPublisher(t)
	mockPromService := mocks.NewMockPrometheusService(t)
	logger := zaptest.NewLogger(t)

	facade := product.NewServiceFacade(
		mockProductService,
		mockPublisher,
		mockPromService,
		logger,
	)

	t.Run("should return the product by ID when it exists", func(t *testing.T) {
		expectedProduct := domain.Product{ID: 1, Name: "Test Product"}
		mockProductService.EXPECT().GetProductByID(mock.Anything, int64(1)).Return(expectedProduct, nil).Once()

		result, err := facade.GetProductByID(context.Background(), 1)

		assert.Equal(t, expectedProduct, result)
		assert.Nil(t, err)

		mockProductService.AssertExpectations(t)
	})

	t.Run("should return an error if the product is not found", func(t *testing.T) {
		expectedError := apperror.ErrProductNotFound
		mockProductService.EXPECT().GetProductByID(mock.Anything, int64(2)).Return(domain.Product{}, expectedError).Once()

		_, err := facade.GetProductByID(context.Background(), 2)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)

		mockProductService.AssertExpectations(t)
	})
}

func TestCreateProduct(t *testing.T) {
	mockProductService := mocks.NewMockProductService(t)
	mockPublisher := mocks.NewMockPublisher(t)
	mockPromService := mocks.NewMockPrometheusService(t)
	logger := zaptest.NewLogger(t)

	facade := product.NewServiceFacade(
		mockProductService,
		mockPublisher,
		mockPromService,
		logger,
	)

	t.Run("should create a product and publish an event on success", func(t *testing.T) {
		product := domain.Product{Name: "New Product"}
		mockPromService.EXPECT().IncProductsCreated().Return().Once()
		mockProductService.EXPECT().CreateProduct(mock.Anything, product).Return(int64(10), nil).Once()
		mockPublisher.EXPECT().Publish(mock.Anything, mock.Anything).Return(nil).Once()

		productID, err := facade.CreateProduct(context.Background(), product)

		assert.Equal(t, int64(10), productID)
		assert.Nil(t, err)

		mockProductService.AssertExpectations(t)
		mockPromService.AssertExpectations(t)
		mockPublisher.AssertExpectations(t)
	})

	t.Run("should return an error if the product creation fails", func(t *testing.T) {
		product := domain.Product{Name: "Failing Product"}
		expectedError := apperror.ErrDB

		mockProductService.EXPECT().CreateProduct(mock.Anything, product).Return(int64(0), expectedError).Once()

		productID, err := facade.CreateProduct(context.Background(), product)

		assert.Equal(t, int64(0), productID)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)

		mockProductService.AssertExpectations(t)
	})

	t.Run("should return the product ID if publishing fails", func(t *testing.T) {
		product := domain.Product{Name: "Product with publish error"}
		productID := int64(20)
		publishError := apperror.ErrMessageBroker

		mockProductService.EXPECT().CreateProduct(mock.Anything, product).Return(productID, nil).Once()
		mockPromService.EXPECT().IncProductsCreated().Return().Once()
		mockPublisher.EXPECT().Publish(mock.Anything, mock.Anything).Return(publishError).Once()

		createdID, err := facade.CreateProduct(context.Background(), product)

		assert.Equal(t, productID, createdID)
		assert.Nil(t, err)

		mockProductService.AssertExpectations(t)
		mockPromService.AssertExpectations(t)
		mockPublisher.AssertExpectations(t)
	})
}

func TestDeleteProduct(t *testing.T) {
	mockProductService := mocks.NewMockProductService(t)
	mockPublisher := mocks.NewMockPublisher(t)
	mockPromService := mocks.NewMockPrometheusService(t)
	logger := zaptest.NewLogger(t)

	facade := product.NewServiceFacade(
		mockProductService,
		mockPublisher,
		mockPromService,
		logger,
	)

	t.Run("should delete a product and publish an event on success", func(t *testing.T) {
		productID := int64(100)
		productObj := domain.Product{ID: productID, Name: "Product to be deleted"}

		mockPromService.EXPECT().IncProductsDeleted().Return().Once()
		mockProductService.EXPECT().GetProductByID(mock.Anything, productID).Return(productObj, nil).Once()
		mockProductService.EXPECT().DeleteProduct(mock.Anything, productID).Return(nil).Once()
		mockPublisher.EXPECT().Publish(mock.Anything, mock.Anything).Return(nil).Once()

		err := facade.DeleteProduct(context.Background(), productID)

		assert.Nil(t, err)

		mockProductService.AssertExpectations(t)
		mockPromService.AssertExpectations(t)
		mockPublisher.AssertExpectations(t)
	})

	t.Run("should return an error if the product to delete is not found", func(t *testing.T) {
		productID := int64(101)
		expectedError := apperror.ErrProductNotFound

		mockProductService.EXPECT().GetProductByID(mock.Anything, productID).Return(domain.Product{}, expectedError).Once()

		err := facade.DeleteProduct(context.Background(), productID)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)

		mockProductService.AssertExpectations(t)
	})

	t.Run("should return an error if the product deletion fails", func(t *testing.T) {
		productID := int64(102)
		productObj := domain.Product{ID: productID, Name: "Product to be deleted"}
		expectedError := apperror.ErrDB

		mockProductService.EXPECT().GetProductByID(mock.Anything, productID).Return(productObj, nil).Once()
		mockProductService.EXPECT().DeleteProduct(mock.Anything, productID).Return(expectedError).Once()

		err := facade.DeleteProduct(context.Background(), productID)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)

		mockProductService.AssertExpectations(t)
	})

	t.Run("should not return an error if publishing fails", func(t *testing.T) {
		productID := int64(103)
		productObj := domain.Product{ID: productID, Name: "Product to be deleted"}
		publishError := apperror.ErrMessageBroker

		mockProductService.EXPECT().GetProductByID(mock.Anything, productID).Return(productObj, nil).Once()
		mockProductService.EXPECT().DeleteProduct(mock.Anything, productID).Return(nil).Once()
		mockPromService.EXPECT().IncProductsDeleted().Return().Once()
		mockPublisher.EXPECT().Publish(mock.Anything, mock.Anything).Return(publishError).Once()

		err := facade.DeleteProduct(context.Background(), productID)

		assert.Nil(t, err)

		mockProductService.AssertExpectations(t)
		mockPromService.AssertExpectations(t)
		mockPublisher.AssertExpectations(t)
	})
}
