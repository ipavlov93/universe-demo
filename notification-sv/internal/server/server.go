package server

import (
	"log"
	"net"
	"net/http"

	"github.com/ipavlov93/universe-demo/universe-pkg/logger"
	"go.uber.org/zap"
)

func RunServer(socket string, lg logger.Logger) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	listener, err := net.Listen("tcp", socket)
	if err != nil {
		lg.Fatal("", zap.Error(err))
	}
	lg.Info("notification-sv has started listening on", zap.String("socket", listener.Addr().String()))

	if err = http.Serve(listener, mux); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"notification-sv":{"status":"healthy"}}`))
}
