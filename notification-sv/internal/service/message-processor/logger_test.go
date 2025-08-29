package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/ipavlov93/universe-demo/product-eventbus-pkg/event"
	"github.com/ipavlov93/universe-demo/product-eventbus-pkg/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestMessageLogger_Run_LogsMessages(t *testing.T) {
	// Create an observer to capture logs
	core, recorded := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	msgLogger := &MessageLogger{lg: logger}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inputChan := make(chan []*message.Message, 1)
	var wg sync.WaitGroup

	msgLogger.Run(ctx, inputChan, &wg)

	expectedFormat := `{"id":1,"name":"Test Product","created_at":"%s"}`
	expectedJSON := fmt.Sprintf(expectedFormat, time.Now().Format(time.RFC3339))

	producer := "TestMessageLogger"
	headers := message.NewHeaders(event.TypeProductCreated, producer)

	env, err := message.New(headers, json.RawMessage(expectedJSON))
	require.NoError(t, err)

	inputChan <- []*message.Message{env}

	time.Sleep(50 * time.Millisecond)

	// Stop goroutine
	cancel()
	wg.Wait()

	logs := recorded.All()

	// Headers
	assert.Equal(t, env.Headers.MessageID, logs[0].ContextMap()["message_id"])
	assert.Equal(t, env.Headers.Type, logs[0].ContextMap()["message_type"])
	assert.Equal(t, producer, logs[0].ContextMap()["producer"])

	//actualJSON, err := json.Marshal(logs[0].ContextMap()["message"])
	//require.NoError(t, err)
	//// Message
	//assert.JSONEq(t, expectedJSON, string(actualJSON))
}
