package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/ipavlov93/universe-demo/product-sv/internal/config"
	"github.com/ipavlov93/universe-demo/product-sv/internal/controller/product/factory"
	logfactory "github.com/ipavlov93/universe-demo/product-sv/internal/infra/logger/factory"
	"github.com/ipavlov93/universe-demo/product-sv/internal/server"
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
	productCtrl := factory.NewProductController(parentCtx, appConfig, lg)

	httpHandler := server.ConfigureRoutes(productCtrl)

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		<-signalCh
		parentCancel()
	}()

	server.Listen(parentCtx, fmt.Sprintf(":%d", appConfig.ServerPort), httpHandler, lg)
}
