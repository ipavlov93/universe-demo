package factory

import (
	"context"
	"time"

	"github.com/ipavlov93/universe-demo/universe-pkg/logger"
	"go.uber.org/zap"

	"github.com/ipavlov93/universe-demo/product-sv/internal/config"
	productctrl "github.com/ipavlov93/universe-demo/product-sv/internal/controller/product"
	"github.com/ipavlov93/universe-demo/product-sv/internal/infra/database"
	adapterfactory "github.com/ipavlov93/universe-demo/product-sv/internal/infra/sqs/adapter/factory"
	"github.com/ipavlov93/universe-demo/product-sv/internal/repository/postgres"
	"github.com/ipavlov93/universe-demo/product-sv/internal/service/facade"
	productsrv "github.com/ipavlov93/universe-demo/product-sv/internal/service/product"
	promservice "github.com/ipavlov93/universe-demo/product-sv/internal/service/prometheus"
	"github.com/ipavlov93/universe-demo/product-sv/internal/service/publisher"
)

func NewProductController(parentCtx context.Context, appConfig config.Config, lg logger.Logger) *productctrl.ProductController {
	ctx, cancel := context.WithTimeout(parentCtx, 3*time.Second)
	defer cancel()

	sqsAdapter, err := adapterfactory.NewInsecureAdapter(parentCtx, appConfig.LocalStackCfg)
	if err != nil {
		lg.Fatal("failed to create AdapterSQS", zap.Error(err))
	}
	sqsPublisher, err := publisher.NewPublisherSQS(ctx, *sqsAdapter, appConfig.LocalStackCfg.Queue)
	if err != nil {
		lg.Fatal("failed to create PublisherSQS", zap.Error(err))
	}

	pgAdapter, err := database.NewWithConfig(appConfig.PostgresCfg)
	if err != nil {
		lg.Fatal("failed to connect to Postgres DB", zap.Error(err))
	}
	promService := promservice.PromService{}
	promService.RegisterPrometheusMetrics()

	productRepository := postgres.NewUserRepository(pgAdapter.GetConnection())
	productService := productsrv.NewProductService(productRepository)
	productSrvFacade := facade.NewServiceFacade(productService, sqsPublisher, promService, lg)
	return productctrl.NewController(productSrvFacade, lg)
}
