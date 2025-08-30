package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ipavlov93/universe-demo/notification-sv/internal/config"
	logfactory "github.com/ipavlov93/universe-demo/notification-sv/internal/infra/logger/factory"
	adapterfactory "github.com/ipavlov93/universe-demo/notification-sv/internal/infra/sqs/adapter/factory"
	"github.com/ipavlov93/universe-demo/notification-sv/internal/service/consumer"
	msglogfactory "github.com/ipavlov93/universe-demo/notification-sv/internal/service/message-logger/factory"
	"github.com/ipavlov93/universe-demo/notification-sv/internal/service/worker"
	"go.uber.org/zap"
)

func main() {
	appConfig := config.LoadConfigEnv()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	parentCtx, parentCancel := context.WithCancel(context.Background())

	appLogger := logfactory.NewAppLogger(os.Stdout, appConfig.MinLogLevel)
	defer appLogger.Sync()

	messageLogger, err := msglogfactory.NewMessageLogger(os.Stdout, appConfig.MinLogLevel)
	if err != nil {
		appLogger.Fatal("failed to create MessageLogger", zap.Error(err))
	}

	childCtx, cancel := context.WithTimeout(parentCtx, 3*time.Second)
	defer cancel()

	// insecure is used only for local development
	sqsAdapter, err := adapterfactory.NewInsecureAdapter(childCtx, appConfig.LocalStackCfg)
	if err != nil {
		appLogger.Fatal("failed to create AdapterSQS", zap.Error(err))
	}
	sqsConsumer, err := consumer.NewConsumerSQS(childCtx, *sqsAdapter, appConfig.LocalStackCfg.Queue, appLogger)
	if err != nil {
		appLogger.Fatal("failed to create ConsumerServiceSQS", zap.Error(err))
	}

	var wg sync.WaitGroup
	workerService := worker.NewWorkerService(sqsConsumer, messageLogger, appConfig.WorkersBufferSize)
	workerService.RunWorkers(parentCtx, &wg)

	appLogger.Info("notification-sv consumer has subscribed to message queues")

	<-signalCh
	parentCancel()
	wg.Wait()

	appLogger.Info("notification-sv shut down")
}
