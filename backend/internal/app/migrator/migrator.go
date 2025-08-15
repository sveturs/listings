package migrator

import (
	"errors"
	"fmt"
	"os"
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

// MigrationConfig содержит конфигурацию для мигратора
type MigrationConfig struct {
	DatabaseURL    string
	MigrationsPath string
	Direction      string
	TargetVersion  string
	WithFixtures   bool
	OnlyFixtures   bool
	FixturesPath   string
	FixturesTable  string
}

// LoadConfigFromEnv загружает конфигурацию из переменных окружения
func LoadConfigFromEnv() *MigrationConfig {
	// Получаем настройки миграций из переменных окружения
	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		migrationsPath = "migrations" // default path
	}

	// Получаем настройки фикстур из переменных окружения
	fixturesPath := os.Getenv("MIGRATIONS_FIXTURES_PATH")
	if fixturesPath == "" {
		fixturesPath = "fixtures" // default path
	}

	fixturesTable := os.Getenv("MIGRATIONS_FIXTURES_TABLE")
	if fixturesTable == "" {
		fixturesTable = "schema_fixtures" // default table name
	}

	direction := os.Getenv("MIGRATION_DIRECTION")
	if direction == "" {
		direction = "up" // default direction
	}

	targetVersion := os.Getenv("MIGRATION_TARGET")

	return &MigrationConfig{
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		MigrationsPath: migrationsPath,
		Direction:      direction,
		TargetVersion:  targetVersion,
		FixturesPath:   fixturesPath,
		FixturesTable:  fixturesTable,
	}
}

// migratorInfo содержит информацию о мигрatore
type migratorInfo struct {
	migr    *Migrator
	name    string
	errName string
}

// RunMigrations запускает миграции в соответствии с конфигурацией
func RunMigrations(config *MigrationConfig) error {
	if config.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is not set")
	}

	logger.Info().
		Str("direction", config.Direction).
		Str("migrationsPath", config.MigrationsPath).
		Str("fixturesPath", config.FixturesPath).
		Str("fixturesTable", config.FixturesTable).
		Str("target", config.TargetVersion).
		Bool("withFixtures", config.WithFixtures).
		Bool("onlyFixtures", config.OnlyFixtures).
		Msg("Starting migration")

	// Проверка конфликта флагов
	if config.WithFixtures && config.OnlyFixtures {
		return fmt.Errorf("cannot use --with-fixtures and --only-fixtures together")
	}

	// Создаем список мигрatorов в зависимости от флагов
	var migrators []migratorInfo

	// Добавляем схемные миграции, если не используется --only-fixtures
	if !config.OnlyFixtures {
		migrators = append(migrators, migratorInfo{
			migr:    NewMigrator(config.MigrationsPath, config.DatabaseURL),
			name:    "migrations",
			errName: "migrations",
		})
	}

	// Добавляем фикстуры если requested (--with-fixtures или --only-fixtures)
	if config.WithFixtures || config.OnlyFixtures {
		logger.Info().Msg("Running fixtures")
		// Проверяем есть ли уже query параметры в URL
		var fixturesDSN string
		if strings.Contains(config.DatabaseURL, "?") {
			fixturesDSN = config.DatabaseURL + "&x-migrations-table=" + config.FixturesTable
		} else {
			fixturesDSN = config.DatabaseURL + "?x-migrations-table=" + config.FixturesTable
		}
		migrators = append(migrators, migratorInfo{
			migr:    NewMigrator(config.FixturesPath, fixturesDSN),
			name:    "fixtures",
			errName: "fixtures",
		})
	}

	// Выполняем миграции для всех мигрatorов
	for _, m := range migrators {
		switch config.Direction {
		case "down":
			if config.TargetVersion != "" {
				// Down to specific version
				if err := m.migr.DownTo(config.TargetVersion); err != nil {
					return fmt.Errorf("error running %s down to version %s: %w", m.errName, config.TargetVersion, err)
				}
			} else {
				// Down all
				if err := m.migr.Down(); err != nil {
					return fmt.Errorf("error running %s down: %w", m.errName, err)
				}
			}
		case "up":
			if config.TargetVersion != "" {
				// Up to specific version
				if err := m.migr.UpTo(config.TargetVersion); err != nil {
					return fmt.Errorf("error running %s up to version %s: %w", m.errName, config.TargetVersion, err)
				}
			} else {
				// Up all pending
				if err := m.migr.Up(); err != nil {
					return fmt.Errorf("error running %s up: %w", m.errName, err)
				}
			}
		default:
			return fmt.Errorf("unknown migration direction: %s", config.Direction)
		}
	}

	return nil
}

// RunMigrationsSchema запускает только схемные миграции
func RunMigrationsSchema(dbURL string) error {
	config := LoadConfigFromEnv()
	config.DatabaseURL = dbURL
	config.WithFixtures = false
	config.OnlyFixtures = false

	logger.Info().
		Str("mode", "schema").
		Str("migrationsPath", config.MigrationsPath).
		Msg("Running schema migrations on API startup")

	return RunMigrations(config)
}

// RunMigrationsFull запускает миграции и фикстуры
func RunMigrationsFull(dbURL string) error {
	config := LoadConfigFromEnv()
	config.DatabaseURL = dbURL
	config.WithFixtures = true
	config.OnlyFixtures = false

	logger.Info().
		Str("mode", "full").
		Str("migrationsPath", config.MigrationsPath).
		Str("fixturesPath", config.FixturesPath).
		Str("fixturesTable", config.FixturesTable).
		Msg("Running full migrations on API startup")

	return RunMigrations(config)
}

// RunMigrationsOnlyFixtures запускает только фикстуры
func RunMigrationsOnlyFixtures(dbURL string) error {
	config := LoadConfigFromEnv()
	config.DatabaseURL = dbURL
	config.WithFixtures = false
	config.OnlyFixtures = true

	logger.Info().
		Str("mode", "fixtures-only").
		Str("fixturesPath", config.FixturesPath).
		Str("fixturesTable", config.FixturesTable).
		Msg("Running fixtures only")

	return RunMigrations(config)
}
