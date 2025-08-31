package error

import (
	"net/http"

	errorpkg "github.com/ipavlov93/universe-demo/universe-pkg/error"
)

var (
	ErrDB            = errorpkg.New("PRODUCT_SV_DB_ERROR", "database error")
	ErrMessageBroker = errorpkg.New("PRODUCT_SV_MESSAGE_BROKER_ERROR", "message broker error")

	ErrInvalidArgument = errorpkg.New("PRODUCT_SV_BAD_REQUEST", "bad request argument(s)").WithCode(http.StatusBadRequest)
	ErrNotFound        = errorpkg.New("PRODUCT_SV_NOT_FOUND", "not found").WithCode(http.StatusNotFound)

	ErrProductNotFound = errorpkg.New("PRODUCT_SV_PRODUCT_NOT_FOUND", "product not found").WithCode(http.StatusNotFound)
	ErrProductExists   = errorpkg.New("PRODUCT_SV_PRODUCT_EXISTS", "product already exists").WithCode(http.StatusConflict)
)
