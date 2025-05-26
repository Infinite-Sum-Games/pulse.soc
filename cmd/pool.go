package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DBPool *pgxpool.Pool

const (
	defaultMaxConns          = int32(100)
	defaultMinConns          = int32(5)
	defaultMaxConnLifetime   = time.Hour
	defaultMaxConnIdleTime   = time.Minute * 30
	defaultHealthCheckPeriod = time.Minute
	defaultConnectTimeout    = time.Second * 5
)

func InitDB() (*pgxpool.Pool, error) {
	connStr := AppConfig.DatabaseURL

	dbConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse connection string: %w", err)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	pool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to create connection pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pool.Acquire(ctx) // Acquire a connection to test
	if err != nil {
		return nil, fmt.Errorf("Failed to acquire connection from pool: %w", err)
	}
	defer conn.Release()

	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("Database connection test failed: %w", err)
	}
	return pool, nil
}
