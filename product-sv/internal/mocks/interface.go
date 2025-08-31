package mocks

import (
	"github.com/ipavlov93/universe-demo/product-sv/internal/controller"
	"github.com/ipavlov93/universe-demo/product-sv/internal/service"
	"github.com/ipavlov93/universe-demo/product-sv/internal/service/facade"
	"github.com/ipavlov93/universe-demo/product-sv/internal/service/publisher"
)

type ProductServiceFacade interface {
	controller.ProductServiceFacade
}

type ProductService interface {
	facade.ProductService
}

type ProductRepository interface {
	service.ProductRepository
}

type PrometheusService interface {
	facade.PrometheusService
}

type Publisher interface {
	facade.Publisher
}

type Adapter interface {
	publisher.Adapter
}

type SQSClientAPI interface {
	publisher.SQSClientAPI
}
