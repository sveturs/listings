package migrator

import (
	"database/sql"
	"fmt"

	"github.com/vondi-global/lib/logger"
	libmigrator "github.com/vondi-global/lib/migrator"

	"github.com/vondi-global/listings/internal/config"
)

// loggerAdapter адаптирует lib/logger к интерфейсу lib/migrator.Logger
type loggerAdapter struct{}

func (l *loggerAdapter) Info() libmigrator.LogEvent {
	return &logEventAdapter{event: logger.Info()}
}

func (l *loggerAdapter) Warn() libmigrator.LogEvent {
	return &logEventAdapter{event: logger.Warn()}
}

func (l *loggerAdapter) Error() libmigrator.LogEvent {
	return &logEventAdapter{event: logger.Error()}
}

// logEventAdapter адаптирует lib/logger.Event к lib/migrator.LogEvent
type logEventAdapter struct {
	event logger.Event
}

func (e *logEventAdapter) Str(key, value string) libmigrator.LogEvent {
	e.event = e.event.Str(key, value)
	return e
}

func (e *logEventAdapter) Int(key string, value int) libmigrator.LogEvent {
	e.event = e.event.Int(key, value)
	return e
}

func (e *logEventAdapter) Err(err error) libmigrator.LogEvent {
	e.event = e.event.Err(err)
	return e
}

func (e *logEventAdapter) Msg(msg string) {
	e.event.Msg(msg)
}

// Migrator - адаптер над lib/migrator для listings сервиса
type Migrator struct {
	db     *sql.DB
	config *config.Config
}

// New создает новый мигратор для listings сервиса
func New(db *sql.DB, cfg *config.Config) *Migrator {
	return &Migrator{
		db:     db,
		config: cfg,
	}
}

// buildConfig конвертирует конфиг приложения в конфиг lib/migrator
func (m *Migrator) buildConfig() libmigrator.Config {
	return libmigrator.Config{
		MigrationsPath:  m.config.DB.MigrationPath,
		MigrationsTable: "schema_migrations",
		FixturesPath:    m.config.Fixtures.Path,
		FixturesTable:   "schema_fixtures",
		Logger:          &loggerAdapter{},
	}
}

// Up применяет все pending миграции
func (m *Migrator) Up() error {
	migrator, err := libmigrator.NewWithConfig(m.db, m.buildConfig())
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	return migrator.Up()
}

// UpWithFixtures применяет миграции и фикстуры
func (m *Migrator) UpWithFixtures() error {
	return libmigrator.RunWithFixtures(m.db, m.buildConfig())
}

// UpOnlyFixtures применяет только фикстуры
func (m *Migrator) UpOnlyFixtures() error {
	return libmigrator.RunOnlyFixtures(m.db, m.buildConfig())
}

// Down откатывает последнюю миграцию
func (m *Migrator) Down() error {
	migrator, err := libmigrator.NewWithConfig(m.db, m.buildConfig())
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	return migrator.Down()
}

// DownOnlyFixtures откатывает последнюю фикстуру
func (m *Migrator) DownOnlyFixtures() error {
	migrator, err := libmigrator.NewForFixtures(m.db, m.buildConfig())
	if err != nil {
		return fmt.Errorf("failed to create fixtures migrator: %w", err)
	}
	defer migrator.Close()

	return migrator.Down()
}

// DownAll откатывает все миграции
func (m *Migrator) DownAll() error {
	migrator, err := libmigrator.NewWithConfig(m.db, m.buildConfig())
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	return migrator.DownAll()
}

// DownAllOnlyFixtures откатывает все фикстуры
func (m *Migrator) DownAllOnlyFixtures() error {
	migrator, err := libmigrator.NewForFixtures(m.db, m.buildConfig())
	if err != nil {
		return fmt.Errorf("failed to create fixtures migrator: %w", err)
	}
	defer migrator.Close()

	return migrator.DownAll()
}

// Status возвращает статус миграций
func (m *Migrator) Status() (libmigrator.MigrationStatus, error) {
	migrator, err := libmigrator.NewWithConfig(m.db, m.buildConfig())
	if err != nil {
		return libmigrator.MigrationStatus{}, fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	return migrator.Status()
}

// FixturesStatus возвращает статус фикстур
func (m *Migrator) FixturesStatus() (libmigrator.MigrationStatus, error) {
	migrator, err := libmigrator.NewForFixtures(m.db, m.buildConfig())
	if err != nil {
		return libmigrator.MigrationStatus{}, fmt.Errorf("failed to create fixtures migrator: %w", err)
	}
	defer migrator.Close()

	return migrator.Status()
}

// Force устанавливает версию принудительно
func (m *Migrator) Force(version int) error {
	migrator, err := libmigrator.NewWithConfig(m.db, m.buildConfig())
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	return migrator.Force(version)
}

// Create создает новую миграцию
func (m *Migrator) Create(name string) error {
	migrator, err := libmigrator.NewWithConfig(m.db, m.buildConfig())
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	return migrator.Create(name)
}

// CreateFixture создает новую фикстуру
func (m *Migrator) CreateFixture(name string) error {
	migrator, err := libmigrator.NewForFixtures(m.db, m.buildConfig())
	if err != nil {
		return fmt.Errorf("failed to create fixtures migrator: %w", err)
	}
	defer migrator.Close()

	return migrator.Create(name)
}
