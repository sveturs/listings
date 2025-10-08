package handler

import (
	"context"
	"fmt"
	"testing"

	"backend/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/stretchr/testify/require"
)

// TestDatabaseConfig - конфигурация тестовой БД
type TestDatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

// GetTestDBConfig - получить конфигурацию тестовой БД
func GetTestDBConfig() *TestDatabaseConfig {
	return &TestDatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "mX3g1XGhMRUZEX3l",
		Name:     "svetubd_test",
	}
}

// SetupTestDatabase - создать и настроить тестовую БД
func SetupTestDatabase(t *testing.T) (*pgxpool.Pool, func()) {
	cfg := GetTestDBConfig()

	// Подключаемся к основной БД для создания тестовой
	mainDSN := fmt.Sprintf("postgres://%s:%s@%s:%d/postgres?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port)

	mainDB, err := pgxpool.New(context.Background(), mainDSN)
	require.NoError(t, err)

	// Создаем тестовую БД если не существует
	ctx := context.Background()
	var exists bool
	err = mainDB.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", cfg.Name).Scan(&exists)
	require.NoError(t, err)

	if !exists {
		_, err = mainDB.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", cfg.Name))
		require.NoError(t, err)
	}
	mainDB.Close()

	// Подключаемся к тестовой БД
	testDSN := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	testDB, err := pgxpool.New(context.Background(), testDSN)
	require.NoError(t, err)

	// Применяем миграции (если нужно)
	err = applyTestMigrations(testDB)
	require.NoError(t, err)

	// Функция очистки
	cleanup := func() {
		cleanupTestData(testDB)
		testDB.Close()
	}

	return testDB, cleanup
}

// applyTestMigrations - применить миграции к тестовой БД
func applyTestMigrations(db *pgxpool.Pool) error {
	ctx := context.Background()

	// Создаем таблицы если не существуют
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS unified_attributes (
			id SERIAL PRIMARY KEY,
			code VARCHAR(100) UNIQUE NOT NULL,
			name VARCHAR(255) NOT NULL,
			attribute_type VARCHAR(50) NOT NULL,
			options JSONB,
			validation_rules JSONB,
			purpose VARCHAR(50) DEFAULT 'regular',
			is_required BOOLEAN DEFAULT false,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		`CREATE TABLE IF NOT EXISTS unified_category_attributes (
			category_id INTEGER NOT NULL,
			attribute_id INTEGER NOT NULL REFERENCES unified_attributes(id) ON DELETE CASCADE,
			is_enabled BOOLEAN DEFAULT true,
			is_required BOOLEAN DEFAULT false,
			is_filter BOOLEAN DEFAULT false,
			sort_order INTEGER DEFAULT 0,
			group_id INTEGER,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (category_id, attribute_id)
		)`,

		`CREATE TABLE IF NOT EXISTS unified_attribute_values (
			id SERIAL PRIMARY KEY,
			entity_type VARCHAR(50) NOT NULL,
			entity_id INTEGER NOT NULL,
			attribute_id INTEGER NOT NULL REFERENCES unified_attributes(id) ON DELETE CASCADE,
			value JSONB NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(entity_type, entity_id, attribute_id)
		)`,

		`CREATE INDEX IF NOT EXISTS idx_unified_attribute_values_entity 
		ON unified_attribute_values(entity_type, entity_id)`,

		`CREATE INDEX IF NOT EXISTS idx_unified_category_attributes_category 
		ON unified_category_attributes(category_id)`,

		`CREATE INDEX IF NOT EXISTS idx_unified_attributes_code 
		ON unified_attributes(code)`,
	}

	for _, migration := range migrations {
		_, err := db.Exec(ctx, migration)
		if err != nil {
			return fmt.Errorf("failed to apply migration: %w", err)
		}
	}

	return nil
}

// cleanupTestData - очистить тестовые данные
func cleanupTestData(db *pgxpool.Pool) {
	ctx := context.Background()

	// Очищаем в правильном порядке из-за foreign keys
	tables := []string{
		"unified_attribute_values",
		"unified_category_attributes",
		"unified_attributes",
	}

	for _, table := range tables {
		_, _ = db.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
	}
}

// CreateTestConfig - создать тестовую конфигурацию
func CreateTestConfig() *config.Config {
	return &config.Config{
		DatabaseURL: "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5432/svetubd_test?sslmode=disable",
		Port:        "3000",
		FeatureFlags: &config.FeatureFlags{
			UseUnifiedAttributes:      true,
			UnifiedAttributesFallback: true,
			UnifiedAttributesPercent:  100,
		},
	}
}

// MockAuthContext - создать контекст с авторизацией для тестов
func MockAuthContext(userID int, isAdmin bool) map[string]interface{} {
	return map[string]interface{}{
		"user_id":  userID,
		"is_admin": isAdmin,
	}
}
