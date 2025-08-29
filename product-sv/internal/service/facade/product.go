package facade

import (
	"context"
	"time"

	"github.com/ipavlov93/universe-demo/product-eventbus-pkg/event"
	"github.com/ipavlov93/universe-demo/product-eventbus-pkg/message"
	"github.com/ipavlov93/universe-demo/universe-pkg/logger"
	"go.uber.org/zap"

	"github.com/ipavlov93/universe-demo/product-sv/internal/domain"
	eventmapper "github.com/ipavlov93/universe-demo/product-sv/internal/mapper/product/event"
	"github.com/ipavlov93/universe-demo/product-sv/internal/service"
)

const producerService = "product-sv"

type ServiceFacadeImp struct {
	productService service.ProductService
	publisher      service.Publisher
	promService    service.PrometheusService
	lg             logger.Logger
}

func NewServiceFacade(
	productService service.ProductService,
	publisher service.Publisher,
	promService service.PrometheusService,
	lg logger.Logger,
) service.Facade {
	return &ServiceFacadeImp{
		productService: productService,
		publisher:      publisher,
		promService:    promService,
		lg:             lg,
	}
}

func (f *ServiceFacadeImp) GetProductByID(ctx context.Context, ID int64) (obj domain.Product, err error) {
	return f.productService.GetProductByID(ctx, ID)
}

func (f *ServiceFacadeImp) CreateProduct(ctx context.Context, product domain.Product) (int64, error) {
	productID, err := f.productService.CreateProduct(ctx, product)
	if err != nil {
		return 0, err
	}

	go f.promService.IncProductsCreated()

	err = f.publishProductCreatedEvent(ctx, product)
	if err != nil {
		// TODO: add retry strategy
		f.lg.Error("failed to publish message",
			zap.Error(err),
			zap.Int64("product_id", productID),
		)
	}
	return productID, nil
}

func (f *ServiceFacadeImp) DeleteProduct(ctx context.Context, productID int64) error {
	productObj, err := f.productService.GetProductByID(ctx, productID)
	if err != nil {
		return err
	}

	err = f.productService.DeleteProduct(ctx, productObj.ID)
	if err != nil {
		return err
	}

	go f.promService.IncProductsDeleted()

	deletedAt := time.Now()
	err = f.publishProductDeletedEvent(ctx, productObj, deletedAt)
	if err != nil {
		// TODO: add retry strategy
		f.lg.Error("failed to publish message",
			zap.Error(err),
			zap.Int64("product_id", productID),
		)
	}

	return nil
}

func (f *ServiceFacadeImp) publishProductCreatedEvent(
	ctx context.Context,
	product domain.Product,
) error {
	createdEvent := eventmapper.ProductCreatedEvent(product)
	headers := message.NewHeaders(event.TypeProductCreated, producerService)

	msg, err := message.New(headers, createdEvent)
	if err != nil {
		return err
	}

	return f.publisher.Publish(ctx, msg)
}

func (f *ServiceFacadeImp) publishProductDeletedEvent(
	ctx context.Context,
	product domain.Product,
	deletedAt time.Time,
) error {
	deletedEvent := eventmapper.ProductDeletedEvent(product, deletedAt)
	headers := message.NewHeaders(event.TypeProductDeleted, producerService)

	msg, err := message.New(headers, deletedEvent)
	if err != nil {
		return err
	}

	return f.publisher.Publish(ctx, msg)
}
