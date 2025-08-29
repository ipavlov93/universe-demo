package prometheus

import "github.com/prometheus/client_golang/prometheus"

var (
	productsCreatedCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "product_sv_products_created_count",
			Help: "Total number of products created",
		},
	)

	productsDeletedCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "product_sv_products_deleted_count",
			Help: "Total number of products deleted",
		},
	)
)

type PromService struct{}

func (p PromService) RegisterPrometheusMetrics() {
	prometheus.MustRegister(productsCreatedCount)
	prometheus.MustRegister(productsDeletedCount)
}

func (p PromService) IncProductsCreated() {
	productsCreatedCount.Inc()
}

func (p PromService) IncProductsDeleted() {
	productsDeletedCount.Inc()
}
