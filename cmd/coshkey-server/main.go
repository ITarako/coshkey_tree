package main

import (
	"flag"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"

	"github.com/ITarako/coshkey_tree/internal/config"
	"github.com/ITarako/coshkey_tree/internal/database"
	"github.com/ITarako/coshkey_tree/internal/server"
)

func main() {
	if err := config.ReadConfigYML("config.yml"); err != nil {
		log.Fatal().Err(err).Msg("Failed init configuration")
	}
	cfg := config.GetConfigInstance()

	flag.Parse()

	log.Info().
		Str("version", cfg.Project.Version).
		Str("commitHash", cfg.Project.CommitHash).
		Bool("debug", cfg.Project.Debug).
		Str("environment", cfg.Project.Environment).
		Msgf("Starting service: %s", cfg.Project.Name)

	if cfg.Project.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SslMode,
	)

	db, err := database.NewPostgres(dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed init postgres")
	}
	defer db.Close()

	if err := server.NewRestServer(db).Start(&cfg); err != nil {
		log.Error().Err(err).Msg("Failed creating rest server")

		return
	}
}