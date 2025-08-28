package id

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

func NewMessageID() string {
	ts := time.Now().UnixNano()
	randBytes := make([]byte, 4)
	rand.Read(randBytes)
	return fmt.Sprintf("%d-%s", ts, hex.EncodeToString(randBytes))
}
