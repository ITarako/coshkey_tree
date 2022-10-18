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

	restAddr := fmt.Sprintf("%s:%v", cfg.Rest.Host, cfg.Rest.Port)
	restServer := createRestServer(restAddr)

	go func() {
		log.Info().Msgf("Rest server is running on %s", restAddr)
		if err := restServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("Failed running rest server")
			cancel()
		}
	}()

	isReady := &atomic.Value{}
	isReady.Store(false)

	//grpcServer := grpc.NewServer(
	//	grpc.KeepaliveParams(keepalive.ServerParameters{
	//		MaxConnectionIdle: time.Duration(cfg.Grpc.MaxConnectionIdle) * time.Minute,
	//		Timeout:           time.Duration(cfg.Grpc.Timeout) * time.Second,
	//		MaxConnectionAge:  time.Duration(cfg.Grpc.MaxConnectionAge) * time.Minute,
	//		Time:              time.Duration(cfg.Grpc.Timeout) * time.Minute,
	//	}),
	//	grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
	//		grpc_ctxtags.UnaryServerInterceptor(),
	//		grpc_prometheus.UnaryServerInterceptor,
	//		grpc_opentracing.UnaryServerInterceptor(),
	//		grpcrecovery.UnaryServerInterceptor(),
	//	)),
	//)

	//r := repo.NewRepo(s.db, s.batchSize)
	//
	//pb.RegisterOmpTemplateApiServiceServer(grpcServer, api.NewTemplateAPI(r))
	//grpc_prometheus.EnableHandlingTimeHistogram()
	//grpc_prometheus.Register(grpcServer)
	//
	//go func() {
	//	log.Info().Msgf("GRPC Server is listening on: %s", grpcAddr)
	//	if err := grpcServer.Serve(l); err != nil {
	//		log.Fatal().Err(err).Msg("Failed running gRPC server")
	//	}
	//}()

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
