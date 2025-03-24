package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"
	"users-app/migrations"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func Connect(ctx context.Context, dsn string, maxConns int32) (*pgxpool.Pool, error) {
	const connectTimeout = time.Second * 5

	dbCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	dbCfg.ConnConfig.ConnectTimeout = connectTimeout
	dbCfg.MaxConns = maxConns

	pool, err := pgxpool.NewWithConfig(ctx, dbCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}

func UpMigrations(pool *pgxpool.Pool) error {
	db := stdlib.OpenDBFromPool(pool)

	fs := migrations.FS
	goose.SetBaseFS(fs)
	goose.SetLogger(goose.NopLogger())

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(db, "."); err != nil && !errors.Is(err, goose.ErrNoNextVersion) {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
