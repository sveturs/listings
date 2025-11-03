// backend/internal/storage/postgres/marketplace_slugs_test.go
package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestIsSlugUnique_ColumnDoesNotExist проверяет поведение когда колонка slug не существует
func TestIsSlugUnique_ColumnDoesNotExist(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := setupTestDatabase(t)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Пытаемся проверить уникальность slug
	isUnique, err := db.IsSlugUnique(ctx, "test-slug", 0)

	// Ожидаем ошибку, так как колонка slug не существует в текущей схеме
	assert.Error(t, err)
	assert.False(t, isUnique)
	assert.Contains(t, err.Error(), "slug column does not exist")
}

// TestGenerateUniqueSlug_ColumnDoesNotExist проверяет поведение когда колонка slug не существует
func TestGenerateUniqueSlug_ColumnDoesNotExist(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := setupTestDatabase(t)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Пытаемся сгенерировать уникальный slug
	slug, err := db.GenerateUniqueSlug(ctx, "test-slug", 0)

	// Ожидаем ошибку, так как колонка slug не существует в текущей схеме
	assert.Error(t, err)
	assert.Empty(t, slug)
	assert.Contains(t, err.Error(), "slug column does not exist")
}

// NOTE: Следующие тесты будут работать только после добавления колонки slug в таблицу c2c_listings
// Они закомментированы до момента создания соответствующей миграции

/*
// TestIsSlugUnique_NewSlug проверяет уникальность нового slug
func TestIsSlugUnique_NewSlug(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := setupTestDatabase(t)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Создаем тестовый листинг со slug
	_, err := db.pool.Exec(ctx, `
		INSERT INTO c2c_listings (user_id, category_id, title, description, price, slug)
		VALUES (1, 1, 'Test Listing', 'Test Description', 100.00, 'test-listing-1')
	`)
	require.NoError(t, err)

	// Проверяем что существующий slug НЕ уникален
	isUnique, err := db.IsSlugUnique(ctx, "test-listing-1", 0)
	require.NoError(t, err)
	assert.False(t, isUnique, "Existing slug should not be unique")

	// Проверяем что новый slug уникален
	isUnique, err = db.IsSlugUnique(ctx, "new-unique-slug", 0)
	require.NoError(t, err)
	assert.True(t, isUnique, "New slug should be unique")
}

// TestIsSlugUnique_WithExclude проверяет уникальность с исключением ID
func TestIsSlugUnique_WithExclude(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := setupTestDatabase(t)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Создаем два тестовых листинга
	var id1, id2 int
	err := db.pool.QueryRow(ctx, `
		INSERT INTO c2c_listings (user_id, category_id, title, description, price, slug)
		VALUES (1, 1, 'Test Listing 1', 'Test Description', 100.00, 'test-listing-1')
		RETURNING id
	`).Scan(&id1)
	require.NoError(t, err)

	err = db.pool.QueryRow(ctx, `
		INSERT INTO c2c_listings (user_id, category_id, title, description, price, slug)
		VALUES (1, 1, 'Test Listing 2', 'Test Description', 200.00, 'test-listing-2')
		RETURNING id
	`).Scan(&id2)
	require.NoError(t, err)

	// Проверяем что slug уникален если исключаем листинг с таким же slug (для update)
	isUnique, err := db.IsSlugUnique(ctx, "test-listing-1", id1)
	require.NoError(t, err)
	assert.True(t, isUnique, "Slug should be unique when excluding own ID")

	// Проверяем что slug НЕ уникален если исключаем другой листинг
	isUnique, err = db.IsSlugUnique(ctx, "test-listing-1", id2)
	require.NoError(t, err)
	assert.False(t, isUnique, "Slug should not be unique when excluding different ID")
}

// TestGenerateUniqueSlug_BaseSlugUnique проверяет генерацию когда базовый slug уникален
func TestGenerateUniqueSlug_BaseSlugUnique(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := setupTestDatabase(t)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Генерируем slug для нового листинга
	slug, err := db.GenerateUniqueSlug(ctx, "unique-slug", 0)
	require.NoError(t, err)
	assert.Equal(t, "unique-slug", slug, "Should return base slug when it's unique")
}

// TestGenerateUniqueSlug_WithSuffix проверяет генерацию с добавлением суффикса
func TestGenerateUniqueSlug_WithSuffix(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := setupTestDatabase(t)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Создаем тестовые листинги с базовым slug и суффиксами
	baseSlug := "popular-item"
	_, err := db.pool.Exec(ctx, `
		INSERT INTO c2c_listings (user_id, category_id, title, description, price, slug)
		VALUES
			(1, 1, 'Item 1', 'Description', 100.00, $1),
			(1, 1, 'Item 2', 'Description', 100.00, $2)
	`, baseSlug, baseSlug+"-1")
	require.NoError(t, err)

	// Генерируем новый уникальный slug
	slug, err := db.GenerateUniqueSlug(ctx, baseSlug, 0)
	require.NoError(t, err)
	assert.Equal(t, "popular-item-2", slug, "Should generate slug with suffix -2")
}

// TestGenerateUniqueSlug_MaxAttempts проверяет что метод возвращает ошибку после 100 попыток
func TestGenerateUniqueSlug_MaxAttempts(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := setupTestDatabase(t)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Создаем 101 листинг с базовым slug и суффиксами от 0 до 100
	baseSlug := "exhausted-slug"
	tx, err := db.pool.Begin(ctx)
	require.NoError(t, err)
	defer tx.Rollback(ctx)

	// Создаем base slug без суффикса
	_, err = tx.Exec(ctx, `
		INSERT INTO c2c_listings (user_id, category_id, title, description, price, slug)
		VALUES (1, 1, 'Item', 'Description', 100.00, $1)
	`, baseSlug)
	require.NoError(t, err)

	// Создаем slug-1 до slug-100
	for i := 1; i <= 100; i++ {
		slug := baseSlug + "-" + string(rune(i))
		_, err = tx.Exec(ctx, `
			INSERT INTO c2c_listings (user_id, category_id, title, description, price, slug)
			VALUES (1, 1, 'Item', 'Description', 100.00, $1)
		`, slug)
		require.NoError(t, err)
	}

	err = tx.Commit(ctx)
	require.NoError(t, err)

	// Пытаемся сгенерировать slug - должна быть ошибка
	slug, err := db.GenerateUniqueSlug(ctx, baseSlug, 0)
	assert.Error(t, err)
	assert.Empty(t, slug)
	assert.Contains(t, err.Error(), "failed to generate unique slug after 100 attempts")
}

// TestGenerateUniqueSlug_WithExclude проверяет генерацию с исключением ID (для update)
func TestGenerateUniqueSlug_WithExclude(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db := setupTestDatabase(t)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Создаем тестовый листинг
	var listingID int
	err := db.pool.QueryRow(ctx, `
		INSERT INTO c2c_listings (user_id, category_id, title, description, price, slug)
		VALUES (1, 1, 'My Listing', 'Description', 100.00, 'my-listing')
		RETURNING id
	`).Scan(&listingID)
	require.NoError(t, err)

	// Пытаемся сгенерировать тот же slug для update - должен вернуть без суффикса
	slug, err := db.GenerateUniqueSlug(ctx, "my-listing", listingID)
	require.NoError(t, err)
	assert.Equal(t, "my-listing", slug, "Should return base slug when excluding own ID")
}
*/

// setupTestDB создает подключение к тестовой БД
// NOTE: Этот helper требует существования функции NewDatabase или аналогичной
func setupTestDatabase(t *testing.T) *Database {
	// TODO: Implement actual test database setup
	// For now, this is a placeholder that will be implemented when tests are enabled

	// Example implementation:
	// db, err := NewDatabase(testDBConnectionString)
	// require.NoError(t, err)
	// return db

	t.Skip("Test database setup not implemented yet")
	return nil
}
