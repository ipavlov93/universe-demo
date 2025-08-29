package product

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ipavlov93/universe-demo/universe-pkg/logger"

	"github.com/ipavlov93/universe-demo/product-sv/internal/controller"
	"github.com/ipavlov93/universe-demo/product-sv/internal/controller/product/dto"
	"github.com/ipavlov93/universe-demo/product-sv/internal/dto/smodel"
	apperror "github.com/ipavlov93/universe-demo/product-sv/internal/error"
	mapper "github.com/ipavlov93/universe-demo/product-sv/internal/mapper/product/smodel"
	"github.com/ipavlov93/universe-demo/product-sv/internal/service"
)

type ProductController struct {
	productSrvFacade service.Facade
	lg               logger.Logger
}

func NewController(productSrvFacade service.Facade, lg logger.Logger) *ProductController {
	return &ProductController{
		productSrvFacade: productSrvFacade,
		lg:               lg,
	}
}

func (c *ProductController) GetProductHandler(w http.ResponseWriter, r *http.Request) {
	productIDStr, err := controller.GetIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid productID format", http.StatusBadRequest)
		return
	}

	product, err := c.productSrvFacade.GetProductByID(r.Context(), productID)
	if err != nil {
		if errors.Is(err, apperror.ErrProductNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(product); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (c *ProductController) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	var productDTO smodel.Product

	if err := json.NewDecoder(r.Body).Decode(&productDTO); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	var productDTO dto.ProductDelete

	if err := json.NewDecoder(r.Body).Decode(&productDTO); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.productSrvFacade.DeleteProduct(r.Context(), productDTO.ProductID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
