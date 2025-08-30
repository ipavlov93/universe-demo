package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/ipavlov93/universe-demo/product-eventbus-pkg/event"
	msgpkg "github.com/ipavlov93/universe-demo/product-eventbus-pkg/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"github.com/ipavlov93/universe-demo/notification-sv/internal/message"
)

func TestMessageLogger_Run_LogsMessages(t *testing.T) {
	// Create an observer to capture logs
	core, recorded := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	msgLogger := &MessageLogger{lg: logger}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	input := make(chan []*message.Envelope, 1)
	out := make(chan []string, 1)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		msgLogger.Process(ctx, input, out)
		wg.Done()
	}()

	expectedFormat := `{"id":1,"name":"Test Product","created_at":"%s"}`
	expectedJSON := fmt.Sprintf(expectedFormat, time.Now().Format(time.RFC3339))

	producer := "TestMessageLogger"
	headers := msgpkg.NewHeaders(event.TypeProductCreated, producer)

	msg, err := msgpkg.New(headers, json.RawMessage(expectedJSON))
	require.NoError(t, err)

	input <- []*message.Envelope{{Message: msg, ReceiptHandle: ""}}

	time.Sleep(50 * time.Millisecond)

	cancel()
	wg.Wait()

	logs := recorded.All()

	// Headers
	assert.Equal(t, msg.Headers.MessageID, logs[0].ContextMap()["message_id"])
	assert.Equal(t, msg.Headers.EventType, logs[0].ContextMap()["event_type"])
	assert.Equal(t, producer, logs[0].ContextMap()["producer"])
}
