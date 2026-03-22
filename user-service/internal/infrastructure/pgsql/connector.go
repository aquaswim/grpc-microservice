package pgsql

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func Connect(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	log.Info().Msg("Connected to PostgreSQL")

	return db, nil
}
