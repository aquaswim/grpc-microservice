package pgsql

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func Connect(ctx context.Context, url string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	// Connection pool configuration
	cfg.MinConns = 2
	cfg.MaxConns = 10

	// Health check configuration
	cfg.HealthCheckPeriod = 1 * time.Minute

	// Connection lifecycle configuration
	cfg.MaxConnLifetime = 1 * time.Hour
	cfg.MaxConnIdleTime = 30 * time.Minute

	// After connect hook for connection setup
	cfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		log.Debug().Msg("New database connection established")
		return nil
	}

	// query logger
	cfg.ConnConfig.Tracer = logger{}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	log.
		Info().
		Str("host", cfg.ConnConfig.Host).
		Uint16("port", cfg.ConnConfig.Port).
		Str("dbname", cfg.ConnConfig.Database).
		Msg("Connected to PostgreSQL")

	return pool, nil
}
