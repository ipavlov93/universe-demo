package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ipavlov93/universe-demo/product-sv/internal/config"
	"github.com/ipavlov93/universe-demo/product-sv/internal/controller/product/factory"
	logfactory "github.com/ipavlov93/universe-demo/product-sv/internal/infra/logger/zap/factory"
	"github.com/ipavlov93/universe-demo/product-sv/internal/server"
)

func main() {
	appConfig := config.LoadConfigEnv()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	appLogger := logfactory.NewAppLogger(os.Stdout, appConfig.MinLogLevel)
	defer appLogger.Sync()

	parentCtx, parentCancel := context.WithCancel(context.Background())
	productCtrl := factory.NewProductController(parentCtx, appConfig, appLogger)

	httpHandler := server.ConfigureRoutes(productCtrl)

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		<-signalCh
		parentCancel()
	}()

	server.Listen(parentCtx, fmt.Sprintf(":%d", appConfig.ServerPort), httpHandler, appLogger)
}
