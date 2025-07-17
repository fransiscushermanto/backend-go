package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/fransiscushermanto/backend/internal/utils"
	_ "github.com/lib/pq"
)

type Database struct {
	*sql.DB
}

func NewDatabase(dataSourceName string) (*Database, error) {
	db, err := sql.Open("postgres", dataSourceName)

	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	ctx, cancel := utils.ContextWithTimeout(5 * time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	utils.Log().Info().Msg("Database connection established")
	return &Database{DB: db}, nil
}
