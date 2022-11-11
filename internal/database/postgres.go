package database

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var StatementBuilder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

// NewPostgres returns DB
func NewPostgres(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create database connection")

		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Error().Err(err).Msg("Failed ping the database")

		return nil, err
	}

	return db, nil
}
