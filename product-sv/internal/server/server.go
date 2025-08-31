package server

import (
	"net/http"

	productctrl "github.com/ipavlov93/universe-demo/product-sv/internal/controller/product"
)

func ConfigureRoutes(mux *http.ServeMux, productController *productctrl.ProductController) http.Handler {
	mux.HandleFunc("/v1/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			productController.CreateProductHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/v1/products/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			productController.GetProductHandler(w, r)
		case http.MethodDelete:
			productController.DeleteProductHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	return mux
}
