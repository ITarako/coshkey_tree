package server

import (
	"net/http"
	"time"

	"coshkey_tree/internal/config"
	"coshkey_tree/internal/server/middleware"
	"coshkey_tree/internal/service/tree"
)

func createRestServer(cfg *config.Config, addr string, service tree.Service) *http.Server {
	restServer := &http.Server{
		Addr:         addr,
		Handler:      middleware.AuthToken(NewHandler(service), cfg.Project.CoshkeyToken),
		ReadTimeout:  time.Duration(cfg.Rest.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Rest.WriteTimeout) * time.Second,
	}

	return restServer
}
