package event

import (
	"encoding/json"
	"fmt"
	"time"
)

type ProductDeletedEvent struct {
	ID          string
	Name        string
	Description string
	DeletedAt   time.Time
}

type productDeleted struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	DeletedAt   string `json:"deleted_at"`
}

func (p *ProductDeletedEvent) Type() Type {
	return TypeProductDeleted
}

// UnmarshalJSON deserializes JSON data with time.RFC3339 format support for timestamps.
func (p *ProductDeletedEvent) UnmarshalJSON(data []byte) error {
	var aux productDeleted
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var deletedAt time.Time
	if aux.DeletedAt != "" {
		timestamp, err := time.Parse(timeFormatRFC3339, aux.DeletedAt)
		if err != nil {
			return fmt.Errorf("invalid deleted_at: %v", err)
		}
		deletedAt = timestamp
	}
	p.setEvent(&aux, deletedAt)
	return nil
}

// MarshalJSON serialises instance to JSON with time.RFC3339 format for timestamps.
func (p *ProductDeletedEvent) MarshalJSON() ([]byte, error) {
	aux := &struct {
		DeletedAt string `json:"deleted_at"`
		*productDeleted
	}{
		DeletedAt:      p.DeletedAt.Format(timeFormatRFC3339),
		productDeleted: convertProductDeleted(p),
	}

	return json.Marshal(aux)
}

func (p *ProductDeletedEvent) setEvent(pd *productDeleted, deletedAt time.Time) {
	if p == nil {
		return
	}

	p.ID = pd.ID
	p.Name = pd.Name
	p.Description = pd.Description
	p.DeletedAt = deletedAt
}

func convertProductDeleted(p *ProductDeletedEvent) *productDeleted {
	return &productDeleted{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		DeletedAt:   p.DeletedAt.String(),
	}
}
