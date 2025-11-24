package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/lib/pq"
)

// CategoryAttribute представляет связь категории с атрибутом
type CategoryAttribute struct {
	ID                      int
	CategoryID              int
	AttributeID             int
	IsEnabled               bool
	IsRequired              bool
	SortOrder               int
	CategorySpecificOptions *json.RawMessage
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

// MigrationStats статистика миграции
type MigrationStats struct {
	TotalRecords    int
	MigratedRecords int
	SkippedRecords  int
	FailedRecords   int
	StartTime       time.Time
	EndTime         time.Time
}

// Config конфигурация подключений
type Config struct {
	SourceDSN      string
	DestinationDSN string
	DryRun         bool
	BatchSize      int
	Verbose        bool
}

func main() {
	config := parseFlags()

	if err := runMigration(config); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}

func parseFlags() *Config {
	config := &Config{}

	flag.StringVar(&config.SourceDSN, "source",
		"postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5433/svetubd?sslmode=disable",
		"Source database DSN (монолит)")
	flag.StringVar(&config.DestinationDSN, "dest",
		"postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable",
		"Destination database DSN (микросервис)")
	flag.BoolVar(&config.DryRun, "dry-run", false,
		"Dry run mode (не вносить изменения)")
	flag.IntVar(&config.BatchSize, "batch-size", 100,
		"Размер батча для вставки")
	flag.BoolVar(&config.Verbose, "verbose", false,
		"Подробный вывод")

	flag.Parse()
	return config
}

func runMigration(config *Config) error {
	ctx := context.Background()

	log.Printf("🚀 Начало миграции category_attributes")
	log.Printf("📊 Режим: %s", getModeString(config.DryRun))
	log.Printf("📦 Размер батча: %d", config.BatchSize)

	// Подключение к базам данных
	sourceDB, err := sql.Open("postgres", config.SourceDSN)
	if err != nil {
		return fmt.Errorf("failed to connect to source DB: %w", err)
	}
	defer sourceDB.Close()

	destDB, err := sql.Open("postgres", config.DestinationDSN)
	if err != nil {
		return fmt.Errorf("failed to connect to destination DB: %w", err)
	}
	defer destDB.Close()

	// Проверка подключений
	if err := sourceDB.PingContext(ctx); err != nil {
		return fmt.Errorf("source DB ping failed: %w", err)
	}
	if err := destDB.PingContext(ctx); err != nil {
		return fmt.Errorf("destination DB ping failed: %w", err)
	}

	log.Printf("✅ Подключение к базам данных успешно")

	// Получение данных из source
	categoryAttrs, err := fetchCategoryAttributes(ctx, sourceDB)
	if err != nil {
		return fmt.Errorf("failed to fetch category attributes: %w", err)
	}

	log.Printf("📥 Получено %d записей из монолита", len(categoryAttrs))

	// Валидация и фильтрация
	validAttrs, invalidCount, err := validateCategoryAttributes(ctx, destDB, categoryAttrs)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if invalidCount > 0 {
		log.Printf("⚠️  Пропущено %d невалидных записей", invalidCount)
	}

	log.Printf("✅ Валидно %d записей для миграции", len(validAttrs))

	// Статистика
	stats := &MigrationStats{
		TotalRecords:   len(categoryAttrs),
		SkippedRecords: invalidCount,
		StartTime:      time.Now(),
	}

	// Миграция данных
	if !config.DryRun {
		log.Printf("💾 Начало вставки данных...")
		if err := migrateCategoryAttributes(ctx, destDB, validAttrs, config, stats); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	} else {
		log.Printf("🔍 DRY RUN: Данные НЕ будут вставлены")
		stats.MigratedRecords = len(validAttrs)
	}

	stats.EndTime = time.Now()

	// Вывод статистики
	printStats(stats)

	return nil
}

func fetchCategoryAttributes(ctx context.Context, db *sql.DB) ([]*CategoryAttribute, error) {
	query := `
		SELECT
			id,
			category_id,
			attribute_id,
			is_enabled,
			is_required,
			sort_order,
			category_specific_options,
			created_at,
			updated_at
		FROM unified_category_attributes
		ORDER BY category_id, sort_order
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categoryAttrs []*CategoryAttribute

	for rows.Next() {
		ca := &CategoryAttribute{}
		var options sql.NullString

		err := rows.Scan(
			&ca.ID,
			&ca.CategoryID,
			&ca.AttributeID,
			&ca.IsEnabled,
			&ca.IsRequired,
			&ca.SortOrder,
			&options,
			&ca.CreatedAt,
			&ca.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if options.Valid && options.String != "" {
			raw := json.RawMessage(options.String)
			ca.CategorySpecificOptions = &raw
		}

		categoryAttrs = append(categoryAttrs, ca)
	}

	return categoryAttrs, rows.Err()
}

func validateCategoryAttributes(ctx context.Context, db *sql.DB, categoryAttrs []*CategoryAttribute) ([]*CategoryAttribute, int, error) {
	// Проверка существования категорий
	categoryIDs := make(map[int]bool)
	for _, ca := range categoryAttrs {
		categoryIDs[ca.CategoryID] = true
	}

	existingCategories, err := checkCategoriesExist(ctx, db, categoryIDs)
	if err != nil {
		return nil, 0, err
	}

	// Проверка существования атрибутов
	attributeIDs := make(map[int]bool)
	for _, ca := range categoryAttrs {
		attributeIDs[ca.AttributeID] = true
	}

	existingAttributes, err := checkAttributesExist(ctx, db, attributeIDs)
	if err != nil {
		return nil, 0, err
	}

	// Фильтрация валидных записей
	var validAttrs []*CategoryAttribute
	invalidCount := 0

	for _, ca := range categoryAttrs {
		if !existingCategories[ca.CategoryID] {
			log.Printf("⚠️  Категория %d не существует, пропускаем запись", ca.CategoryID)
			invalidCount++
			continue
		}

		if !existingAttributes[ca.AttributeID] {
			log.Printf("⚠️  Атрибут %d не существует, пропускаем запись", ca.AttributeID)
			invalidCount++
			continue
		}

		validAttrs = append(validAttrs, ca)
	}

	return validAttrs, invalidCount, nil
}

func checkCategoriesExist(ctx context.Context, db *sql.DB, categoryIDs map[int]bool) (map[int]bool, error) {
	if len(categoryIDs) == 0 {
		return make(map[int]bool), nil
	}

	// Преобразуем map в slice для запроса
	ids := make([]int, 0, len(categoryIDs))
	for id := range categoryIDs {
		ids = append(ids, id)
	}

	query := `SELECT id FROM categories WHERE id = ANY($1)`
	rows, err := db.QueryContext(ctx, query, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	existingIDs := make(map[int]bool)
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		existingIDs[id] = true
	}

	return existingIDs, rows.Err()
}

func checkAttributesExist(ctx context.Context, db *sql.DB, attributeIDs map[int]bool) (map[int]bool, error) {
	if len(attributeIDs) == 0 {
		return make(map[int]bool), nil
	}

	// Преобразуем map в slice для запроса
	ids := make([]int, 0, len(attributeIDs))
	for id := range attributeIDs {
		ids = append(ids, id)
	}

	query := `SELECT id FROM attributes WHERE id = ANY($1)`
	rows, err := db.QueryContext(ctx, query, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	existingIDs := make(map[int]bool)
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		existingIDs[id] = true
	}

	return existingIDs, rows.Err()
}

func migrateCategoryAttributes(ctx context.Context, db *sql.DB, categoryAttrs []*CategoryAttribute, config *Config, stats *MigrationStats) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Подготовка statement для вставки
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO category_attributes (
			category_id,
			attribute_id,
			is_enabled,
			is_required,
			is_searchable,
			is_filterable,
			sort_order,
			category_specific_options,
			is_active,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (category_id, attribute_id) DO UPDATE SET
			is_enabled = EXCLUDED.is_enabled,
			is_required = EXCLUDED.is_required,
			sort_order = EXCLUDED.sort_order,
			category_specific_options = EXCLUDED.category_specific_options,
			updated_at = EXCLUDED.updated_at
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Вставка данных батчами
	for i, ca := range categoryAttrs {
		var optionsJSON *string
		if ca.CategorySpecificOptions != nil {
			str := string(*ca.CategorySpecificOptions)
			optionsJSON = &str
		}

		_, err := stmt.ExecContext(ctx,
			ca.CategoryID,
			ca.AttributeID,
			ca.IsEnabled,
			ca.IsRequired,
			true, // is_searchable - по умолчанию true
			true, // is_filterable - по умолчанию true
			ca.SortOrder,
			optionsJSON,
			ca.IsEnabled, // is_active = is_enabled
			ca.CreatedAt,
			ca.UpdatedAt,
		)
		if err != nil {
			stats.FailedRecords++
			log.Printf("❌ Ошибка вставки записи %d (category_id=%d, attribute_id=%d): %v",
				i+1, ca.CategoryID, ca.AttributeID, err)
			continue
		}

		stats.MigratedRecords++

		if config.Verbose && (i+1)%100 == 0 {
			log.Printf("📊 Обработано %d/%d записей", i+1, len(categoryAttrs))
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func printStats(stats *MigrationStats) {
	duration := stats.EndTime.Sub(stats.StartTime)

	separator := strings.Repeat("═", 60)
	fmt.Println("\n" + separator)
	fmt.Println("📊 СТАТИСТИКА МИГРАЦИИ")
	fmt.Println(separator)
	fmt.Printf("📥 Всего записей в источнике:    %d\n", stats.TotalRecords)
	fmt.Printf("✅ Успешно мигрировано:          %d\n", stats.MigratedRecords)
	fmt.Printf("⚠️  Пропущено (невалидные):      %d\n", stats.SkippedRecords)
	fmt.Printf("❌ Ошибки при вставке:           %d\n", stats.FailedRecords)
	fmt.Printf("⏱️  Время выполнения:            %s\n", duration.Round(time.Millisecond))
	fmt.Println(separator)

	if stats.FailedRecords > 0 {
		fmt.Println("⚠️  ВНИМАНИЕ: Некоторые записи не были мигрированы!")
		os.Exit(1)
	}

	if stats.MigratedRecords == stats.TotalRecords-stats.SkippedRecords {
		fmt.Println("✅ Миграция завершена успешно!")
	}
}

func getModeString(dryRun bool) string {
	if dryRun {
		return "🔍 DRY RUN (без изменений)"
	}
	return "💾 PRODUCTION (с записью в БД)"
}
