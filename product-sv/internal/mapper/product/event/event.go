package event

import (
	"time"

	"github.com/ipavlov93/universe-demo/product-eventbus-pkg/event"

	"github.com/ipavlov93/universe-demo/product-sv/internal/domain"
)

func ProductCreatedEvent(p domain.Product) event.ProductCreatedEvent {
	return event.ProductCreatedEvent{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
	}
}

func ProductDeletedEvent(p domain.Product, deletedAt time.Time) event.ProductDeletedEvent {
	return event.ProductDeletedEvent{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		DeletedAt:   deletedAt,
	}
}
