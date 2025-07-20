package db

import (
	"context"
	"embed"
	"fmt"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/meh-hackathon/meh/apperror"
	"github.com/meh-hackathon/meh/logger"


	"github.com/golang-migrate/migrate/v4"
	migrate_postgres "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

var (
	ErrConnectionFailed = apperror.Define("db:connection_failed", "Database connection failed")
	ErrTransaction      = apperror.Define("db:transaction_begin", "Failed to begin database transaction")
	ErrReset            = apperror.Define("db:reset", "Failed to reset database to a clean state")
	ErrMigrationFailed  = apperror.Define("db:migrate", "Failed to migrate database")
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

func Init(host string, port string, user string, password string, options ...dbOption) error {
	start := time.Now()
	if strings.TrimSpace(host) == "" {
		return ErrConnectionFailed.WithMessage("Database host must be provided").WithOrigin()
	}
	if strings.TrimSpace(port) == "" {
		return ErrConnectionFailed.WithMessage("Database port must be provided").WithOrigin()
	}
	if strings.TrimSpace(user) == "" {
		return ErrConnectionFailed.WithMessage("Database user must be provided").WithOrigin()
	}
	if strings.TrimSpace(password) == "" {
		return ErrConnectionFailed.WithMessage("Database password must be provided").WithOrigin()
	}

	cfg := &config{
		maxConns:        25,
		maxIdleConns:    5,
		connMaxLifetime: 5 * time.Minute,
	}

	for _, opt := range options {
		opt(cfg)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/meh", user, password, host, port)
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

//go:embed migrations/*.sql
var migrations embed.FS

// createMigrator creates a new migrate instance
func createMigrator() (*migrate.Migrate, error) {
	driver, err := migrate_postgres.WithInstance(DB.DB, &migrate_postgres.Config{})
	if err != nil {
		return nil, ErrMigrationFailed.WithMessage("failed to create migration database instance").WithOrigin().WithCause(err)
	}

	d, err := iofs.New(migrations, "migrations")
	if err != nil {
		return nil, ErrMigrationFailed.WithMessage("failed to read migrations").WithOrigin().WithCause(err)
	}

	m, err := migrate.NewWithInstance("iofs", d, "postgres", driver)
	if err != nil {
		return nil, ErrMigrationFailed.WithMessage("failed to create migration instance").WithOrigin().WithCause(err)
	}

	return m, nil
}

// MigrateUp migrates the database to the latest version
func MigrateUp() error {
	m, err := createMigrator()
	if err != nil {
		return err
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return ErrMigrationFailed.WithMessage("failed to get current version").WithOrigin().WithCause(err)
	}

	logger.Info("Database version before migration", "version", version, "dirty", dirty)

	if err == migrate.ErrNilVersion {
		// No migrations applied yet
		logger.Info("No migrations applied yet")
		version = 0
	}

	// Apply all pending migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return ErrMigrationFailed.WithMessage("failed to run migrations").WithOrigin().WithCause(err)
	}

	newVersion, dirty, _ := m.Version()
	logger.Info("Database version after migration", "version", newVersion, "dirty", dirty)
	return nil
}

// MigrateDown migrates the database down to a specific version
// If targetVersion is 0, it will revert all migrations
func MigrateDown(targetVersion uint) error {
	m, err := createMigrator()
	if err != nil {
		return err
	}

	currentVersion, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return ErrMigrationFailed.WithMessage("failed to get current version").WithOrigin().WithCause(err)
	}

	if err == migrate.ErrNilVersion {
		logger.Info("No migrations to revert - database is at version 0")
		return nil
	}

	logger.Info("Database version before downgrade", "version", currentVersion, "dirty", dirty, "target", targetVersion)

	if targetVersion >= currentVersion {
		return ErrMigrationFailed.WithMessage(fmt.Sprintf("target version %d is not lower than current version %d", targetVersion, currentVersion)).WithOrigin()
	}

	if targetVersion == 0 {
		// Migrate all the way down
		err = m.Down()
	} else {
		// Migrate to specific version
		err = m.Migrate(targetVersion)
	}

	if err != nil && err != migrate.ErrNoChange {
		return ErrMigrationFailed.WithMessage("failed to run downgrade migrations").WithOrigin().WithCause(err)
	}

	newVersion, dirty, _ := m.Version()
	if err == migrate.ErrNilVersion {
		newVersion = 0
	}
	logger.Info("Database version after downgrade", "version", newVersion, "dirty", dirty)
	return nil
}
