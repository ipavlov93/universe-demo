package server

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/ipavlov93/universe-demo/universe-pkg/logger"
	"go.uber.org/zap"
)

// Listen makes blocking http.Serve call
func Listen(ctx context.Context, socket string, handler http.Handler, lg logger.Logger) {
	listener, err := net.Listen("tcp", socket)
	if err != nil {
		lg.Fatal("can't start net.Listen", zap.Error(err))
	}
	lg.Info("product-sv API server has started listening on", zap.String("socket", listener.Addr().String()))

	server := &http.Server{Handler: handler}

	go func() {
		if err = server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			lg.Fatal("can't start http.Serve, server failed to start", zap.Error(err))
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = server.Shutdown(shutdownCtx); err != nil {
		lg.Fatal("server shutdown failed", zap.Error(err))
	}

	lg.Info("API server has been gracefully shut down")
}
