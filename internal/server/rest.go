package server

import (
	"net/http"
	"time"

	"github.com/ITarako/coshkey_tree/internal/config"
	"github.com/ITarako/coshkey_tree/internal/service/tree"
)

func createRestServer(cfg *config.Config, addr string, service tree.Service) *http.Server {
	restServer := &http.Server{
		Addr:         addr,
		Handler:      NewHandler(service),
		ReadTimeout:  time.Duration(cfg.Rest.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Rest.WriteTimeout) * time.Second,
	}

	return restServer
}
