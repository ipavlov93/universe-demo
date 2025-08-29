package server

import (
	"net/http"

	productctrl "github.com/ipavlov93/universe-demo/product-sv/internal/controller/product"
)

func ConfigureRoutes(productController *productctrl.ProductController) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/product", productController.CreateProductHandler)
	mux.HandleFunc("/product/", func(w http.ResponseWriter, r *http.Request) {
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
