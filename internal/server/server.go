package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"

	"github.com/ITarako/coshkey_tree/internal/config"
	"github.com/ITarako/coshkey_tree/internal/service/tree"
)

type RestServer struct {
	db        *sqlx.DB
	batchSize uint
}

func NewRestServer(db *sqlx.DB, batchSize uint) *RestServer {
	return &RestServer{
		db:        db,
		batchSize: batchSize,
	}
}

func (s *RestServer) Start(cfg *config.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r := tree.NewRepository(s.db, s.batchSize)
	service := tree.NewService(r)

	restAddr := fmt.Sprintf("%s:%v", cfg.Rest.Host, cfg.Rest.Port)
	restServer := createRestServer(cfg, restAddr, service)

	go func() {
		log.Info().Msgf("Rest server is running on %s", restAddr)
		if err := restServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("Failed running rest server")
			cancel()
		}
	}()

	isReady := &atomic.Value{}
	isReady.Store(false)

	go func() {
		time.Sleep(2 * time.Second)
		isReady.Store(true)
		log.Info().Msg("The service is ready to accept requests")
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		log.Info().Msgf("signal.Notify: %v", v)
	case done := <-ctx.Done():
		log.Info().Msgf("ctx.Done: %v", done)
	}

	isReady.Store(false)

	if err := restServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("restServer.Shutdown")
	} else {
		log.Info().Msg("restServer shut down correctly")
	}

	return nil
}
