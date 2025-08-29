package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/ipavlov93/universe-demo/notification-sv/internal/config"
	logfactory "github.com/ipavlov93/universe-demo/notification-sv/internal/infra/logger/factory"
	adapterfactory "github.com/ipavlov93/universe-demo/notification-sv/internal/infra/sqs/adapter/factory"
	"github.com/ipavlov93/universe-demo/notification-sv/internal/service/consumer"
	msgprocessor "github.com/ipavlov93/universe-demo/notification-sv/internal/service/message-processor"
	"github.com/ipavlov93/universe-demo/notification-sv/internal/service/worker"
)

func main() {
	appConfig := config.LoadConfigEnv()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	lg, err := logfactory.NewAppLoggerOrDefault(appConfig.MinLogLevel)
	if err != nil {
		lg = zap.NewNop()
	}

	defer lg.Sync()

	parentCtx, parentCancel := context.WithCancel(context.Background())
	ctx, cancel := context.WithTimeout(parentCtx, 3*time.Second)
	defer cancel()

	messageLogger := msgprocessor.NewMessageLogger(lg)

	sqsAdapter, err := adapterfactory.NewInsecureAdapter(parentCtx, appConfig.LocalStackCfg)
	if err != nil {
		lg.Fatal("failed to create AdapterSQS", zap.Error(err))
	}
	sqsConsumer, err := consumer.NewConsumerSQS(ctx, *sqsAdapter, appConfig.LocalStackCfg.Queue, lg)
	if err != nil {
		lg.Fatal("failed to create ConsumerServiceSQS", zap.Error(err))
	}

	var wg sync.WaitGroup
	worker.RunWorkers(parentCtx, sqsConsumer, messageLogger, &wg)

	<-signalCh
	parentCancel()
	wg.Wait()

	lg.Info("Notification service shut down")
}
