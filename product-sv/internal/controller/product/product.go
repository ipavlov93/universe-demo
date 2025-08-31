package product

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ipavlov93/universe-demo/product-sv/internal/controller"
	"github.com/ipavlov93/universe-demo/product-sv/internal/dto/smodel"
	apperror "github.com/ipavlov93/universe-demo/product-sv/internal/error"
	mapper "github.com/ipavlov93/universe-demo/product-sv/internal/mapper/product/smodel"
	errorpkg "github.com/ipavlov93/universe-demo/universe-pkg/error"
	"github.com/ipavlov93/universe-demo/universe-pkg/logger"
)

type ProductController struct {
	productSrvFacade controller.ProductServiceFacade
	lg               logger.Logger
}

func NewController(productSrvFacade controller.ProductServiceFacade, lg logger.Logger) *ProductController {
	return &ProductController{
		productSrvFacade: productSrvFacade,
		lg:               lg,
	}
}

func (c *ProductController) GetProductHandler(w http.ResponseWriter, r *http.Request) {
	productID, err := pathValueID(r, "id")
	if err != nil {
		writeErrorStatusCode(w, err)
		return
	}

	product, err := c.productSrvFacade.GetProductByID(r.Context(), productID)
	if err != nil {
		if errors.Is(err, apperror.ErrProductNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	productDTO := mapper.ProductToProductDto(product)
	err = json.NewEncoder(w).Encode(productDTO)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}

func (c *ProductController) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	var productDTO smodel.Product

	if err := json.NewDecoder(r.Body).Decode(&productDTO); err != nil {
		http.Error(w, apperror.ErrInvalidArgument.Error(), http.StatusBadRequest)
		return
	}

	product := mapper.ProductDtoToProduct(productDTO)

	productID, err := c.productSrvFacade.CreateProduct(r.Context(), product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"product_id":"%d"}`, productID)))
}

func (c *ProductController) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	productID, err := pathValueID(r, "id")
	if err != nil {
		writeErrorStatusCode(w, err)
		return
	}

	err = c.productSrvFacade.DeleteProduct(r.Context(), productID)
	if err != nil {
		writeErrorStatusCode(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func writeErrorStatusCode(w http.ResponseWriter, err error) {
	var appError errorpkg.AppError
	if !errors.As(err, &appError) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	writeAppError(w, appError)
}

func writeAppError(w http.ResponseWriter, appError errorpkg.AppError) {
	code, errorMsg := appError.ToHTTP()
	http.Error(w, errorMsg, code)
}

func pathValueID(r *http.Request, idKey string) (int64, error) {
	productIDStr := r.PathValue(idKey)
	if productIDStr == "" {
		return 0, apperror.ErrNotFound
	}

	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		return 0, apperror.ErrInvalidArgument.WithReason(
			fmt.Sprintf("invalid product id: %s", productIDStr))
	}
	return productID, nil
}
