package domain

import (
	"fmt"
	"time"
)

type Product struct {
	ID          int64
	Name        string
	Description string
	CreatedAt   time.Time
}

func (p Product) Valid() bool {
	if p.ID < 1 {
		return false
	}
	if p.Name == "" {
		return false
	}
	return true
}

func (p Product) String() string {
	return fmt.Sprintf("ID: %d, Name: %s", p.ID, p.Name)
}
