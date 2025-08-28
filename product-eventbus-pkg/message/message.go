package message

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	Headers Headers         `json:"headers"`
	Payload json.RawMessage `json:"payload"`
}

// EncodeJSON serializes the Message instance to a JSON byte slice.
func (e *Message) EncodeJSON() ([]byte, error) {
	return json.Marshal(e)
}

// DecodeJSON deserializes a JSON byte slice into the Message.
func (e *Message) DecodeJSON(data []byte) error {
	if err := json.Unmarshal(data, e); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}
	return nil
}

func New[T any](headers *Headers, payload T) (*Message, error) {
	if headers == nil {
		return nil, fmt.Errorf("headers cannot be nil")
	}

	payloadJSON, err := marshalPayload(payload)
	if err != nil {
		return nil, err
	}

	return &Message{
		Headers: *headers,
		Payload: payloadJSON,
	}, nil
}

func marshalPayload[T any](payload T) (json.RawMessage, error) {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return payloadJSON, nil
}
