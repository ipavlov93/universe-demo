package message

import (
	msgpkg "github.com/ipavlov93/universe-demo/product-eventbus-pkg/message"
)

type Envelope struct {
	ReceiptHandle string
	Message       *msgpkg.Message
}

func NewEnvelope(receiptHandle string, msg *msgpkg.Message) *Envelope {
	return &Envelope{
		ReceiptHandle: receiptHandle,
		Message:       msg,
	}
}
