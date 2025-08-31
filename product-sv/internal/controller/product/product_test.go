package product

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ipavlov93/universe-demo/product-sv/internal/domain"
	"github.com/ipavlov93/universe-demo/product-sv/internal/dto/smodel"
	apperror "github.com/ipavlov93/universe-demo/product-sv/internal/error"
	mapper "github.com/ipavlov93/universe-demo/product-sv/internal/mapper/product/smodel"
	"github.com/ipavlov93/universe-demo/product-sv/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestGetProductHandler(t *testing.T) {
	mockFacade := mocks.NewMockFacade(t)
	logger := zaptest.NewLogger(t)
	controller := NewController(mockFacade, logger)

	router := http.NewServeMux()
	router.HandleFunc("/v1/products/{id}", controller.GetProductHandler)

	t.Run("should return 200 OK and the product on success", func(t *testing.T) {
		expectedProduct := domain.Product{ID: 10, Name: "Test Product"}
		expectedProductDTO := mapper.ProductToProductDto(expectedProduct)
		expectedBody, err := json.Marshal(expectedProductDTO)
		require.NoError(t, err)

		req := httptest.NewRequest("GET", "/v1/products/10", nil)

		rec := httptest.NewRecorder()

		mockFacade.EXPECT().GetProductByID(mock.Anything, int64(10)).Return(expectedProduct, nil).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
		assert.Equal(t, string(expectedBody)+"\n", rec.Body.String())
		mockFacade.AssertCalled(t, "GetProductByID", mock.Anything, int64(10))
	})

	t.Run("should return 400 Bad Request for an invalid ID format", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/v1/products/abc", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "bad request argument(s)\n", rec.Body.String())
		mockFacade.AssertNotCalled(t, "GetProductByID")
	})

	t.Run("should return 404 Not Found if product doesn't exist", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/v1/products/11", nil)
		rec := httptest.NewRecorder()

		mockFacade.EXPECT().GetProductByID(mock.Anything, int64(11)).Return(domain.Product{}, apperror.ErrProductNotFound).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, apperror.ErrProductNotFound.Error()+"\n", rec.Body.String())
	})

	t.Run("should return 500 Internal Server Error", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/v1/products/12", nil)
		rec := httptest.NewRecorder()

		expectedError := errors.New("database connection failed")
		mockFacade.EXPECT().GetProductByID(mock.Anything, int64(12)).Return(domain.Product{}, expectedError).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, expectedError.Error()+"\n", rec.Body.String())
	})
}

func TestCreateProductHandler(t *testing.T) {
	mockFacade := mocks.NewMockFacade(t)
	logger := zaptest.NewLogger(t)
	controller := NewController(mockFacade, logger)

	t.Run("should return 201 Created and the product ID on success", func(t *testing.T) {
		productDTO := smodel.Product{Name: "New Product"}
		expectedID := int64(123)

		requestBody, _ := json.Marshal(productDTO)
		req := httptest.NewRequest("POST", "/v1/products", bytes.NewBuffer(requestBody))
		rec := httptest.NewRecorder()

		mockFacade.EXPECT().CreateProduct(mock.Anything, mock.AnythingOfType("domain.Product")).Return(expectedID, nil).Once()

		controller.CreateProductHandler(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, fmt.Sprintf(`{"product_id":"%d"}`, expectedID), rec.Body.String())
		mockFacade.AssertCalled(t, "CreateProduct", mock.Anything, mock.AnythingOfType("domain.Product"))
	})

	t.Run("should return 400 Bad Request for an invalid JSON body", func(t *testing.T) {
		invalidJSON := `{"name": 123}` // Name should be a string
		req := httptest.NewRequest("POST", "/v1/products", bytes.NewBufferString(invalidJSON))
		rec := httptest.NewRecorder()

		controller.CreateProductHandler(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, apperror.ErrInvalidArgument.Error()+"\n", rec.Body.String())
		mockFacade.AssertNotCalled(t, "CreateProduct")
	})

	t.Run("should return 500 Internal Server Error for facade failure", func(t *testing.T) {
		productDTO := smodel.Product{Name: "New Product"}
		expectedError := errors.New("database connection failed")

		requestBody, _ := json.Marshal(productDTO)
		req := httptest.NewRequest("POST", "/v1/products", bytes.NewBuffer(requestBody))
		rec := httptest.NewRecorder()

		mockFacade.EXPECT().CreateProduct(mock.Anything, mock.AnythingOfType("domain.Product")).Return(int64(0), expectedError).Once()

		controller.CreateProductHandler(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, expectedError.Error()+"\n", rec.Body.String())
	})
}

func TestDeleteProductHandler(t *testing.T) {
	mockFacade := mocks.NewMockFacade(t)
	logger := zaptest.NewLogger(t)
	controller := NewController(mockFacade, logger)

	router := http.NewServeMux()
	router.HandleFunc("/v1/products/{id}", controller.DeleteProductHandler)

	t.Run("should return 200 OK on successful deletion", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/v1/products/456", nil)
		rec := httptest.NewRecorder()

		mockFacade.EXPECT().DeleteProduct(mock.Anything, int64(456)).Return(nil).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockFacade.AssertCalled(t, "DeleteProduct", mock.Anything, int64(456))
	})

	t.Run("should return 400 Bad Request for invalid JSON body", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/v1/products/abc", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "bad request argument(s)\n", rec.Body.String())
		mockFacade.AssertNotCalled(t, "DeleteProduct")
	})

	t.Run("should return 500 Internal Server Error for facade failure", func(t *testing.T) {
		expectedError := errors.New("deletion failed")

		req := httptest.NewRequest("DELETE", "/v1/products/457", nil)
		rec := httptest.NewRecorder()

		mockFacade.EXPECT().DeleteProduct(mock.Anything, int64(457)).Return(expectedError).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, "", rec.Body.String())
	})
}
