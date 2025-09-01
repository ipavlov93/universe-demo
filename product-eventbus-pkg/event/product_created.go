package event

import (
	"encoding/json"
	"fmt"
	"time"
)

type ProductCreatedEvent struct {
	ID          int64
	Name        string
	Description string
	CreatedAt   time.Time
}

type productCreated struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at"`
}

func (p *ProductCreatedEvent) Type() Type {
	return TypeProductCreated
}

// UnmarshalJSON deserializes JSON data with time.RFC3339 format support for timestamps.
func (p *ProductCreatedEvent) UnmarshalJSON(data []byte) error {
	var aux productCreated
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var createdAt time.Time
	if aux.CreatedAt != "" {
		timestamp, err := time.Parse(timeFormatRFC3339, aux.CreatedAt)
		if err != nil {
			return fmt.Errorf("invalid created_at: %v", err)
		}
		createdAt = timestamp
	}
	p.setEvent(&aux, createdAt)
	return nil
}

// MarshalJSON serialises instance to JSON with time.RFC3339 format for timestamps.
func (p *ProductCreatedEvent) MarshalJSON() ([]byte, error) {
	aux := &struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description,omitempty"`
		DeletedAt   string `json:"deleted_at"`
	}{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		DeletedAt:   p.CreatedAt.Format(time.RFC3339),
	}

	return json.Marshal(aux)
}

func (p *ProductCreatedEvent) setEvent(pc *productCreated, createdAt time.Time) {
	if p == nil {
		return
	}

	p.ID = pc.ID
	p.Name = pc.Name
	p.Description = pc.Description
	p.CreatedAt = createdAt
}
