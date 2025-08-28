package event_test

import (
	"testing"
	"time"

	"github.com/ipavlov93/universe-demo/product-eventbus-pkg/event"
	"github.com/stretchr/testify/require"
)

func TestParsePayload(t *testing.T) {
	createdAt := time.Now()

	messagePayload := newProductCreated(
		"Test ID",
		"Test Product",
		"Test Description",
		createdAt,
	)

	// Prepare test message
	payload, err := messagePayload.MarshalJSON()
	require.NoError(t, err)

	// ACT
	m, err := event.ParsePayload(messagePayload.Type().String(), payload)
	require.NoError(t, err)

	t.Log(m)
}

func newProductCreated(
	id string,
	name string,
	description string,
	createdAt time.Time,
) event.ProductCreatedEvent {
	return event.ProductCreatedEvent{
		ID:          id,
		Name:        name,
		Description: description,
		CreatedAt:   createdAt,
	}
}
