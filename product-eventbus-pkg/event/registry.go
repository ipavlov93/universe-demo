package event

import (
	"encoding/json"
	"fmt"
)

type Payload interface {
	json.Marshaler
	json.Unmarshaler
}

var messageRegistry = map[Type]func() Payload{
	TypeProductCreated: func() Payload { return &ProductCreatedEvent{} },
	TypeProductDeleted: func() Payload { return &ProductDeletedEvent{} },
}

func ParsePayload(msgType string, payload []byte) (json.Unmarshaler, error) {
	eventType, ok := eventTypeFound(msgType)
	if !ok {
		return nil, fmt.Errorf("unknown message type %s", msgType)
	}

	constructor, ok := messageRegistry[eventType]
	if !ok {
		return nil, fmt.Errorf("unsupported message type: %s", msgType)
	}

	v := constructor()
	if err := v.UnmarshalJSON(payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s: %v", msgType, err)
	}

	return v, nil
}

func eventTypeFound(t string) (Type, bool) {
	eventType := Type(t)
	return eventType, eventType.Valid()
}
