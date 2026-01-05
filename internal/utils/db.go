package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

type Database struct {
	Pool *pgxpool.Pool
}

func NewDatabase(dataSourceName string) (*Database, error) {
	config, err := pgxpool.ParseConfig(dataSourceName)

	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = 5 * time.Minute
	config.MaxConnIdleTime = 1 * time.Minute

	ctx, cancel := ContextWithTimeout(5 * time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	Log().Info().Msg("Database connection established")
	return &Database{Pool: pool}, nil
}

func (db *Database) WithTransaction(ctx context.Context, fn func(tx pgx.Tx) error) (err error) {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			// A panic occurred, rollback and re-panic.
			tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			// The function returned an error, rollback.
			tx.Rollback(ctx) // err is already set
		} else {
			// The function was successful, commit.
			err = tx.Commit(ctx) // if commit fails, the error will be returned
		}
	}()

	err = fn(tx)
	return err
}

func (db *Database) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}
