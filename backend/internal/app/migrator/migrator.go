package migrator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"backend/internal/logger"
)

type Migrator struct {
	migrationPath string
	dsn           string
}

// NewMigrator creates a new Migrator instance
func NewMigrator(migrationPath, dsn string) *Migrator {
	if !strings.HasPrefix(migrationPath, "file://") {
		migrationPath = "file://" + migrationPath
	}
	return &Migrator{
		migrationPath: migrationPath,
		dsn:           dsn,
	}
}

func (m *Migrator) Up() error {
	logger.Info().Str("migrationPath", m.migrationPath).Str("dsn", m.dsn).Msg("migrating...")

	migrator, err := migrate.New(m.migrationPath, m.dsn)
	if err != nil {
		return fmt.Errorf("error creating migrate instance: %w", err)
	}
	defer migrator.Close()

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("error running migrations up: %w", err)
	}

	logger.Info().Msg("Migrations up completed successfully")
	return nil
}

// Down reverts all migrations
// TODO: I want only 1 down migration
func (m *Migrator) Down() error {
	logger.Info().Str("migrationPath", m.migrationPath).Str("dsn", m.dsn).Msg("migrating...")

	migrator, err := migrate.New(m.migrationPath, m.dsn)
	if err != nil {
		return fmt.Errorf("error creating migrate instance: %w", err)
	}
	defer migrator.Close()

	if err := migrator.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("error running migrations down: %w", err)
	}

	logger.Info().Msg("Migrations down completed successfully")
	return nil
}
