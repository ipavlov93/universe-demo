package event

import (
	"fmt"
	"time"
)

type Type string

const timeFormatRFC3339 = time.RFC3339

const (
	TypeProductCreated Type = "product.created"
	TypeProductDeleted Type = "product.deleted"
)

var validTypes = map[Type]struct{}{
	TypeProductCreated: {},
	TypeProductDeleted: {},
}

func (t Type) Valid() bool {
	if _, ok := validTypes[t]; !ok {
		return false
	}
	return true
}

func (t Type) String() string {
	_, ok := validTypes[t]
	if !ok {
		return fmt.Sprintf("Type(%s)", string(t))
	}
	return string(t)
}
