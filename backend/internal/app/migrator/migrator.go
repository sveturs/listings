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
	defer func() {
		if sourceErr, dbErr := migrator.Close(); sourceErr != nil || dbErr != nil {
			logger.Warn().Interface("sourceErr", sourceErr).Interface("dbErr", dbErr).Msg("Failed to close migrator")
		}
	}()

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("error running migrations up: %w", err)
	}

	logger.Info().Msg("Migrations up completed successfully")
	return nil
}

// Down reverts all migrations
func (m *Migrator) Down() error {
	logger.Info().Str("migrationPath", m.migrationPath).Msg("Reverting all migrations...")

	migrator, err := migrate.New(m.migrationPath, m.dsn)
	if err != nil {
		return fmt.Errorf("error creating migrate instance: %w", err)
	}
	defer func() {
		if sourceErr, dbErr := migrator.Close(); sourceErr != nil || dbErr != nil {
			logger.Warn().Interface("sourceErr", sourceErr).Interface("dbErr", dbErr).Msg("Failed to close migrator")
		}
	}()

	if err := migrator.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("error running migrations down: %w", err)
	}

	logger.Info().Msg("All migrations reverted successfully")
	return nil
}

// UpTo migrates up to a specific version
func (m *Migrator) UpTo(version string) error {
	logger.Info().
		Str("migrationPath", m.migrationPath).
		Str("targetVersion", version).
		Msg("Migrating up to specific version...")

	migrator, err := migrate.New(m.migrationPath, m.dsn)
	if err != nil {
		return fmt.Errorf("error creating migrate instance: %w", err)
	}
	defer func() {
		if sourceErr, dbErr := migrator.Close(); sourceErr != nil || dbErr != nil {
			logger.Warn().Interface("sourceErr", sourceErr).Interface("dbErr", dbErr).Msg("Failed to close migrator")
		}
	}()

	// Convert version string to uint
	targetVersion, err := parseVersion(version)
	if err != nil {
		return fmt.Errorf("error parsing target version: %w", err)
	}

	if err := migrator.Migrate(targetVersion); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("error migrating to version %d: %w", targetVersion, err)
	}

	logger.Info().Uint("version", targetVersion).Msg("Migrated to specific version successfully")
	return nil
}

// DownTo migrates down to a specific version
func (m *Migrator) DownTo(version string) error {
	logger.Info().
		Str("migrationPath", m.migrationPath).
		Str("targetVersion", version).
		Msg("Reverting to specific version...")

	migrator, err := migrate.New(m.migrationPath, m.dsn)
	if err != nil {
		return fmt.Errorf("error creating migrate instance: %w", err)
	}
	defer func() {
		if sourceErr, dbErr := migrator.Close(); sourceErr != nil || dbErr != nil {
			logger.Warn().Interface("sourceErr", sourceErr).Interface("dbErr", dbErr).Msg("Failed to close migrator")
		}
	}()

	// Convert version string to uint
	targetVersion, err := parseVersion(version)
	if err != nil {
		return fmt.Errorf("error parsing target version: %w", err)
	}

	if err := migrator.Migrate(targetVersion); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("error migrating to version %d: %w", targetVersion, err)
	}

	logger.Info().Uint("version", targetVersion).Msg("Reverted to specific version successfully")
	return nil
}

// parseVersion converts version string to uint
func parseVersion(version string) (uint, error) {
	// Remove leading zeros if any
	version = strings.TrimLeft(version, "0")
	if version == "" {
		version = "0"
	}

	var v uint
	_, err := fmt.Sscanf(version, "%d", &v)
	if err != nil {
		return 0, fmt.Errorf("invalid version format: %s", version)
	}
	return v, nil
}
