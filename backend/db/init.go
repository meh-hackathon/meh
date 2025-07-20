package db

import (
	"context"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/meh-hackathon/meh/apperror"
	"github.com/meh-hackathon/meh/logger"
)

var (
	ErrConnectionFailed = apperror.Define("db:connection_failed", "Database connection failed")
	ErrTransaction      = apperror.Define("db:transaction_begin", "Failed to begin database transaction")
	ErrReset            = apperror.Define("db:reset", "Failed to reset database to a clean state")
)

var DB *sqlx.DB

type dbOption func(*config)

type config struct {
	maxConns        int
	maxIdleConns    int
	connMaxLifetime time.Duration
}

func WithMaxConns(n int) dbOption {
	return func(c *config) { c.maxConns = n }
}

func WithMaxIdleConns(n int) dbOption {
	return func(c *config) { c.maxIdleConns = n }
}

func WithConnMaxLifetime(d time.Duration) dbOption {
	return func(c *config) { c.connMaxLifetime = d }
}

func Init(dsn string, options ...dbOption) error {
	start := time.Now()
	if strings.TrimSpace(dsn) == "" {
		return ErrConnectionFailed.WithMessage("Database DSN must be provided").WithOrigin()
	}

	cfg := &config{
		maxConns:        25,
		maxIdleConns:    5,
		connMaxLifetime: 5 * time.Minute,
	}

	for _, opt := range options {
		opt(cfg)
	}

	var err error
	DB, err = sqlx.Open("pgx", dsn)
	if err != nil {
		return ErrConnectionFailed.WithMessage("Failed to open database connection").WithOrigin().WithCause(err)
	}

	DB.SetMaxOpenConns(cfg.maxConns)
	DB.SetMaxIdleConns(cfg.maxIdleConns)
	DB.SetConnMaxLifetime(cfg.connMaxLifetime)

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()

	if err := DB.PingContext(pingCtx); err != nil {
		DB.Close()
		return ErrConnectionFailed.WithMessage("Failed to ping database").WithOrigin().WithCause(err)
	}

	logger.Debug("Database connection established", "took", time.Since(start))
	return nil
}
