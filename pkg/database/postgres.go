package database

import (
	"context"
	"fmt"
	"time"

	"HailowSellerService/pkg/logging"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	Pool   *pgxpool.Pool
	logger *logging.Logger
}

func NewPostgresClient(connectLink string, logger *logging.Logger) (*Postgres, error) {
	config, err := pgxpool.ParseConfig(connectLink)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse config: %v", err)
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnIdleTime = 30 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)

	if err != nil {
		logger.Fatalf("Failed to create postgres pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		logger.Fatalf("Failed to ping postgres pool: %v", err)
	}

	logger.Info("Database initializaed successfully")

	return &Postgres{
		Pool:   pool,
		logger: logger,
	}, nil
}

func (db *Postgres) Close() {
	if db.Pool != nil {
		db.Pool.Close()
		db.logger.Info("Database closed successfully")
	}
}
