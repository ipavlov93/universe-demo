package message

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ProductCreated struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func TestMessage_Decode(t *testing.T) {
	t.Run("Valid JSON should be decoded correctly", func(t *testing.T) {
		// ARRANGE
		payloadBytes := `{}`

		jsonString := fmt.Sprintf(
			`{"headers":{"event_type":"product.created","producer":"test-producer","message_id":"abc-123"},"payload":%s}`,
			string(payloadBytes))

		var decodedMessage Message
		dataToDecode := []byte(jsonString)

		// ACT
		err := decodedMessage.DecodeJSON(dataToDecode)
		require.NoError(t, err)

		// ASSERT
		assert.Equal(t, "product.created", decodedMessage.Headers.EventType)
		assert.Equal(t, "test-producer", decodedMessage.Headers.Producer)
		assert.Equal(t, "abc-123", decodedMessage.Headers.MessageID)
	})

	t.Run("Invalid JSON headers should return an error", func(t *testing.T) {
		// ARRANGE
		invalidJSON := []byte(`{"headers":"should be object"`)

		// ACT
		var decodedMessage Message
		err := decodedMessage.DecodeJSON(invalidJSON)

		// ASSERT
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unmarshal message")
	})
	t.Run("Invalid JSON payload should return an error", func(t *testing.T) {
		// ARRANGE
		invalidJSON := []byte(`{"payload":{}`)

		// ACT
		var decodedMessage Message
		err := decodedMessage.DecodeJSON(invalidJSON)

		// ASSERT
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unmarshal message")
	})
}
