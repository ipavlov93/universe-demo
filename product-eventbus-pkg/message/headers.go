package message

import (
	"errors"
	"time"

	"github.com/ipavlov93/universe-demo/product-eventbus-pkg/event"
	"github.com/ipavlov93/universe-demo/product-eventbus-pkg/message/id"
)

type Headers struct {
	MessageID        string    `json:"message_id"`
	CorrelationID    string    `json:"correlation_id,omitempty"`
	CausationID      string    `json:"causation_id,omitempty"`
	Type             string    `json:"type"`
	Version          int       `json:"version,omitempty"`
	MessageCreatedAt time.Time `json:"message_created_at"`
	Producer         string    `json:"producer"`
}

func (h *Headers) Validate() error {
	if h.MessageID == "" {
		return errors.New("header messageID is required")
	}
	if h.Type == "" {
		return errors.New("header type is required")
	}
	if h.MessageCreatedAt.IsZero() {
		return errors.New("header messageCreatedAt is required")
	}
	return nil
}

func NewHeaders(msgType event.Type, producer string) *Headers {
	return &Headers{
		MessageID:        id.NewMessageID(),
		Type:             msgType.String(),
		Producer:         producer,
		MessageCreatedAt: time.Now().UTC(),
	}
}

func (h *Headers) SetDefaultHeaders() {
	if h.MessageID == "" {
		h.MessageID = id.NewMessageID()
	}
	if h.MessageCreatedAt.IsZero() {
		h.MessageCreatedAt = time.Now().UTC()
	}
}
