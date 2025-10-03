package service_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"backend/internal/domain/models"
	"backend/internal/domain/search"
	"backend/internal/proj/marketplace/service"
	"backend/internal/proj/storefronts/storage/opensearch"
	"backend/internal/storage"
	"backend/internal/storage/filestorage"
	"backend/internal/types"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var ErrNotImplemented = errors.New("not implemented in mock")

type CategoryIntegrationTestSuite struct {
	suite.Suite
	ctx          context.Context
	pgContainer  *postgres.PostgresContainer
	db           *sql.DB
	storage      storage.Storage
	service      *service.MarketplaceService
	redisCache   service.CacheInterface
	testCategory *models.MarketplaceCategory
}

func (suite *CategoryIntegrationTestSuite) SetupSuite() {
	ctx := context.Background()
	suite.ctx = ctx

	// Запускаем PostgreSQL контейнер
	pgContainer, err := postgres.Run(ctx, "postgres:16",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Minute)),
	)
	require.NoError(suite.T(), err)
	suite.pgContainer = pgContainer

	// Получаем строку подключения
	connectionString, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(suite.T(), err)

	// Подключаемся к БД
	db, err := sql.Open("pgx", connectionString)
	require.NoError(suite.T(), err)
	suite.db = db

	// Создаем схему БД
	err = suite.createSchema()
	require.NoError(suite.T(), err)

	// Создаем storage
	suite.storage = &testStorage{db: suite.db}

	// Создаем Redis mock кеш
	suite.redisCache = &memoryCache{data: make(map[string][]byte)}

	// Создаем сервис
	translationService := &dummyTranslationService{}
	suite.service = service.NewMarketplaceService(suite.storage, translationService, nil, suite.redisCache).(*service.MarketplaceService)
}

func (suite *CategoryIntegrationTestSuite) TearDownSuite() {
	if suite.db != nil {
		if err := suite.db.Close(); err != nil {
			log.Printf("Ошибка при закрытии подключения к базе данных: %v", err)
		}
	}
	if suite.pgContainer != nil {
		if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
			log.Printf("Ошибка при остановке контейнера: %v", err)
		}
	}
}

func (suite *CategoryIntegrationTestSuite) SetupTest() {
	// Очищаем данные перед каждым тестом
	suite.cleanupData()

	// Создаем тестовую категорию
	suite.testCategory = &models.MarketplaceCategory{
		Name:        "Тестовая категория",
		Slug:        "test-category",
		Description: "Описание тестовой категории",
		IsActive:    true,
	}
}

func (suite *CategoryIntegrationTestSuite) cleanupData() {
	queries := []string{
		"DELETE FROM translations",
		"DELETE FROM category_attribute_mapping",
		"DELETE FROM category_attributes",
		"DELETE FROM marketplace_listings",
		"DELETE FROM marketplace_categories",
		"DELETE FROM rating_cache",
	}

	for _, query := range queries {
		_, err := suite.db.Exec(query)
		if err != nil {
			log.Printf("Ошибка при очистке таблицы: %v", err)
		}
	}

	// Обновляем материализованное представление
	_, _ = suite.db.Exec("REFRESH MATERIALIZED VIEW category_listing_count_view")
}

func (suite *CategoryIntegrationTestSuite) createSchema() error {
	schema := `
	-- Основные таблицы
	CREATE TABLE IF NOT EXISTS marketplace_categories (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		slug VARCHAR(255) NOT NULL UNIQUE,
		parent_id INTEGER REFERENCES marketplace_categories(id),
		icon VARCHAR(255),
		has_custom_ui BOOLEAN DEFAULT FALSE,
		custom_ui_component VARCHAR(255),
		description TEXT,
		is_active BOOLEAN DEFAULT TRUE,
		seo_title VARCHAR(255),
		seo_description TEXT,
		seo_keywords TEXT,
		sort_order INTEGER DEFAULT 0,
		level INTEGER DEFAULT 0,
		count INTEGER DEFAULT 0,
		external_id VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS translations (
		id SERIAL PRIMARY KEY,
		entity_type VARCHAR(50) NOT NULL,
		entity_id INTEGER NOT NULL,
		language VARCHAR(10) NOT NULL,
		field_name VARCHAR(100) NOT NULL,
		translated_text TEXT NOT NULL,
		is_machine_translated BOOLEAN DEFAULT FALSE,
		is_verified BOOLEAN DEFAULT FALSE,
		provider VARCHAR(50),
		metadata JSONB,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(entity_type, entity_id, language, field_name)
	);

	CREATE TABLE IF NOT EXISTS category_attributes (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		display_name VARCHAR(255) NOT NULL,
		type VARCHAR(50) NOT NULL,
		is_required BOOLEAN DEFAULT FALSE,
		is_filterable BOOLEAN DEFAULT TRUE,
		validation_rules JSONB,
		options JSONB,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS category_attribute_mapping (
		id SERIAL PRIMARY KEY,
		category_id INTEGER NOT NULL REFERENCES marketplace_categories(id),
		attribute_id INTEGER NOT NULL REFERENCES category_attributes(id),
		is_required BOOLEAN DEFAULT FALSE,
		is_enabled BOOLEAN DEFAULT TRUE,
		sort_order INTEGER DEFAULT 0,
		custom_component VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(category_id, attribute_id)
	);

	CREATE TABLE IF NOT EXISTS marketplace_listings (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		price DECIMAL(15,2) NOT NULL,
		category_id INTEGER REFERENCES marketplace_categories(id),
		user_id INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Материализованное представление для счетчиков
	CREATE MATERIALIZED VIEW IF NOT EXISTS category_listing_count_view AS
	WITH RECURSIVE category_tree AS (
		SELECT id, parent_id FROM marketplace_categories
	),
	listing_counts AS (
		SELECT category_id, COUNT(*) as count 
		FROM marketplace_listings 
		GROUP BY category_id
	)
	SELECT 
		c.id,
		COALESCE(lc.count, 0) as listing_count
	FROM category_tree c
	LEFT JOIN listing_counts lc ON c.id = lc.category_id;

	-- Таблица для кеша рейтингов
	CREATE TABLE IF NOT EXISTS rating_cache (
		entity_type VARCHAR(50) NOT NULL,
		entity_id INTEGER NOT NULL,
		average_rating DECIMAL(3,2),
		total_reviews INTEGER DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (entity_type, entity_id)
	);
	`

	_, err := suite.db.Exec(schema)
	return err
}

// Тест создания категории с переводами
func (suite *CategoryIntegrationTestSuite) TestCreateCategoryWithTranslations() {
	// Подготовка
	category := &models.MarketplaceCategory{
		Name:           "Electronics",
		Slug:           "electronics",
		Description:    "Electronic devices and accessories",
		IsActive:       true,
		SEOTitle:       "Buy Electronics Online",
		SEODescription: "Best deals on electronic devices",
		SEOKeywords:    "electronics, gadgets, devices",
		Translations: map[string]string{
			"ru": "Электроника",
			"sr": "Електроника",
		},
	}

	// Выполнение
	id, err := suite.service.CreateCategory(suite.ctx, category)

	// Проверки
	assert.NoError(suite.T(), err)
	assert.Greater(suite.T(), id, 0)

	// Проверяем, что категория создана
	var count int
	err = suite.db.QueryRow("SELECT COUNT(*) FROM marketplace_categories WHERE id = $1", id).Scan(&count)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, count)

	// Проверяем переводы
	var translationCount int
	err = suite.db.QueryRow("SELECT COUNT(*) FROM translations WHERE entity_type = 'category' AND entity_id = $1", id).Scan(&translationCount)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, translationCount)

	// Проверяем конкретный перевод
	var translatedText string
	err = suite.db.QueryRow(
		"SELECT translated_text FROM translations WHERE entity_type = 'category' AND entity_id = $1 AND language = 'ru'",
		id,
	).Scan(&translatedText)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Электроника", translatedText)
}

// Тест обновления категории с инвалидацией кеша
func (suite *CategoryIntegrationTestSuite) TestUpdateCategoryWithCacheInvalidation() {
	// Создаем категорию
	id, err := suite.service.CreateCategory(suite.ctx, suite.testCategory)
	require.NoError(suite.T(), err)

	// Загружаем категорию в кеш
	category, err := suite.service.GetCategoryByID(suite.ctx, id)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Тестовая категория", category.Name)

	// Обновляем категорию
	category.Name = "Обновленная категория"
	category.SEOTitle = "SEO заголовок"
	category.SEODescription = "SEO описание"
	category.SEOKeywords = "ключевые, слова"
	category.Translations = map[string]string{
		"en": "Updated Category",
		"ru": "Обновленная категория",
	}

	err = suite.service.UpdateCategory(suite.ctx, category)
	assert.NoError(suite.T(), err)

	// Проверяем, что данные обновились в БД
	var name, seoTitle string
	err = suite.db.QueryRow(
		"SELECT name, seo_title FROM marketplace_categories WHERE id = $1",
		id,
	).Scan(&name, &seoTitle)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Обновленная категория", name)
	assert.Equal(suite.T(), "SEO заголовок", seoTitle)

	// Проверяем, что кеш был инвалидирован и данные обновились
	updatedCategory, err := suite.service.GetCategoryByID(suite.ctx, id)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Обновленная категория", updatedCategory.Name)
	assert.Equal(suite.T(), "SEO заголовок", updatedCategory.SEOTitle)
}

// Тест удаления категории с проверкой подкатегорий
func (suite *CategoryIntegrationTestSuite) TestDeleteCategoryWithSubcategories() {
	// Создаем родительскую категорию
	parentID, err := suite.service.CreateCategory(suite.ctx, &models.MarketplaceCategory{
		Name:     "Родительская категория",
		Slug:     "parent-category",
		IsActive: true,
	})
	require.NoError(suite.T(), err)

	// Создаем дочернюю категорию
	childCategory := &models.MarketplaceCategory{
		Name:     "Дочерняя категория",
		Slug:     "child-category",
		ParentID: &parentID,
		IsActive: true,
	}
	childID, err := suite.service.CreateCategory(suite.ctx, childCategory)
	require.NoError(suite.T(), err)
	assert.Greater(suite.T(), childID, 0)

	// Пытаемся удалить родительскую категорию (должна быть ошибка)
	err = suite.service.DeleteCategory(suite.ctx, parentID)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "дочерних категорий")

	// Удаляем сначала дочернюю категорию
	err = suite.service.DeleteCategory(suite.ctx, childID)
	assert.NoError(suite.T(), err)

	// Теперь можем удалить родительскую
	err = suite.service.DeleteCategory(suite.ctx, parentID)
	assert.NoError(suite.T(), err)

	// Проверяем, что категории удалены
	var count int
	err = suite.db.QueryRow("SELECT COUNT(*) FROM marketplace_categories WHERE id IN ($1, $2)", parentID, childID).Scan(&count)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 0, count)
}

// Тест получения категорий с кешированием
func (suite *CategoryIntegrationTestSuite) TestGetCategoriesWithCaching() {
	// Создаем несколько категорий
	categories := []models.MarketplaceCategory{
		{Name: "Категория 1", Slug: "cat-1", IsActive: true},
		{Name: "Категория 2", Slug: "cat-2", IsActive: true},
		{Name: "Категория 3", Slug: "cat-3", IsActive: false},
	}

	for _, cat := range categories {
		_, err := suite.service.CreateCategory(suite.ctx, &cat)
		require.NoError(suite.T(), err)
	}

	// Первый запрос - загружаем из БД
	result1, err := suite.service.GetCategories(suite.ctx)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result1, 2) // Только активные

	// Второй запрос - должен взять из кеша
	result2, err := suite.service.GetCategories(suite.ctx)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result2, 2)

	// Проверяем, что это те же данные
	assert.Equal(suite.T(), result1[0].ID, result2[0].ID)
	assert.Equal(suite.T(), result1[0].Name, result2[0].Name)
}

// Тест сохранения и проверки SEO полей
func (suite *CategoryIntegrationTestSuite) TestSEOFieldsSaving() {
	// Создаем категорию с SEO полями
	category := &models.MarketplaceCategory{
		Name:           "SEO Test Category",
		Slug:           "seo-test",
		Description:    "Category for SEO testing",
		IsActive:       true,
		SEOTitle:       "Best SEO Test Category | Buy Online",
		SEODescription: "Find the best products in our SEO test category. Great prices and fast delivery.",
		SEOKeywords:    "seo, test, category, online, shopping",
	}

	id, err := suite.service.CreateCategory(suite.ctx, category)
	require.NoError(suite.T(), err)

	// Загружаем категорию и проверяем SEO поля
	loaded, err := suite.service.GetCategoryByID(suite.ctx, id)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), category.SEOTitle, loaded.SEOTitle)
	assert.Equal(suite.T(), category.SEODescription, loaded.SEODescription)
	assert.Equal(suite.T(), category.SEOKeywords, loaded.SEOKeywords)

	// Обновляем SEO поля
	loaded.SEOTitle = "Updated SEO Title"
	loaded.SEODescription = "Updated SEO description with more keywords"
	loaded.SEOKeywords = "updated, seo, keywords"

	err = suite.service.UpdateCategory(suite.ctx, loaded)
	assert.NoError(suite.T(), err)

	// Проверяем обновление
	updated, err := suite.service.GetCategoryByID(suite.ctx, id)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Updated SEO Title", updated.SEOTitle)
	assert.Equal(suite.T(), "Updated SEO description with more keywords", updated.SEODescription)
	assert.Equal(suite.T(), "updated, seo, keywords", updated.SEOKeywords)
}

// Тест проверки удаления категории с объявлениями
func (suite *CategoryIntegrationTestSuite) TestDeleteCategoryWithListings() {
	// Создаем категорию
	categoryID, err := suite.service.CreateCategory(suite.ctx, &models.MarketplaceCategory{
		Name:     "Категория с объявлениями",
		Slug:     "category-with-listings",
		IsActive: true,
	})
	require.NoError(suite.T(), err)

	// Создаем объявление в этой категории
	_, err = suite.db.Exec(
		"INSERT INTO marketplace_listings (title, description, price, category_id, user_id) VALUES ($1, $2, $3, $4, $5)",
		"Тестовое объявление", "Описание", 100.00, categoryID, 1,
	)
	require.NoError(suite.T(), err)

	// Пытаемся удалить категорию (должна быть ошибка)
	err = suite.service.DeleteCategory(suite.ctx, categoryID)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "объявлений")

	// Удаляем объявление
	_, err = suite.db.Exec("DELETE FROM marketplace_listings WHERE category_id = $1", categoryID)
	require.NoError(suite.T(), err)

	// Теперь можем удалить категорию
	err = suite.service.DeleteCategory(suite.ctx, categoryID)
	assert.NoError(suite.T(), err)
}

// Тест перемещения категории в иерархии
func (suite *CategoryIntegrationTestSuite) TestMoveCategory() {
	// Создаем структуру категорий
	// Родитель 1
	parent1ID, err := suite.service.CreateCategory(suite.ctx, &models.MarketplaceCategory{
		Name: "Родитель 1", Slug: "parent-1", IsActive: true,
	})
	require.NoError(suite.T(), err)

	// Родитель 2
	parent2ID, err := suite.service.CreateCategory(suite.ctx, &models.MarketplaceCategory{
		Name: "Родитель 2", Slug: "parent-2", IsActive: true,
	})
	require.NoError(suite.T(), err)

	// Дочерняя категория под Родителем 1
	childID, err := suite.service.CreateCategory(suite.ctx, &models.MarketplaceCategory{
		Name: "Дочерняя", Slug: "child", ParentID: &parent1ID, IsActive: true,
	})
	require.NoError(suite.T(), err)

	// Перемещаем дочернюю категорию под Родителя 2
	err = suite.service.MoveCategory(suite.ctx, childID, parent2ID)
	assert.NoError(suite.T(), err)

	// Проверяем, что категория переместилась
	var parentID int
	err = suite.db.QueryRow("SELECT parent_id FROM marketplace_categories WHERE id = $1", childID).Scan(&parentID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), parent2ID, parentID)
}

// Тест автоматического создания переводов
func (suite *CategoryIntegrationTestSuite) TestAutomaticTranslationCreation() {
	// Создаем мок сервиса перевода с предсказуемыми результатами
	mockTranslationService := &mockTranslationService{
		translations: map[string]string{
			"en": "Furniture",
			"ru": "Мебель",
			"sr": "Намештај",
		},
	}
	// Создаем новый сервис с моком
	mockService := service.NewMarketplaceService(suite.storage, mockTranslationService, nil, suite.redisCache).(*service.MarketplaceService)

	// Создаем категорию без явных переводов
	category := &models.MarketplaceCategory{
		Name:        "Мебель",
		Slug:        "furniture",
		Description: "Мебель для дома и офиса",
		IsActive:    true,
	}

	id, err := mockService.CreateCategory(suite.ctx, category)
	require.NoError(suite.T(), err)

	// Проверяем, что автоматические переводы созданы
	var translations []struct {
		Language       string
		TranslatedText string
		IsMachine      bool
	}

	rows, err := suite.db.Query(
		`SELECT language, translated_text, is_machine_translated 
		 FROM translations 
		 WHERE entity_type = 'category' AND entity_id = $1 AND field_name = 'name'
		 ORDER BY language`,
		id,
	)
	require.NoError(suite.T(), err)
	defer func() {
		if err := rows.Close(); err != nil {
			// Логируем ошибку закрытия rows в тесте
			_ = err // Explicitly ignore error
		}
	}()

	for rows.Next() {
		var t struct {
			Language       string
			TranslatedText string
			IsMachine      bool
		}
		err := rows.Scan(&t.Language, &t.TranslatedText, &t.IsMachine)
		require.NoError(suite.T(), err)
		translations = append(translations, t)
	}

	// Check for iteration errors
	require.NoError(suite.T(), rows.Err())

	// Проверяем, что создались переводы для en и sr
	assert.Len(suite.T(), translations, 2)

	// Проверяем английский перевод
	enTranslation := translations[0]
	assert.Equal(suite.T(), "en", enTranslation.Language)
	assert.Equal(suite.T(), "Furniture", enTranslation.TranslatedText)
	assert.True(suite.T(), enTranslation.IsMachine)

	// Проверяем сербский перевод
	srTranslation := translations[1]
	assert.Equal(suite.T(), "sr", srTranslation.Language)
	assert.Equal(suite.T(), "Намештај", srTranslation.TranslatedText)
	assert.True(suite.T(), srTranslation.IsMachine)
}

// Запуск тестов
func TestCategoryIntegrationTestSuite(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Пропуск интеграционных тестов")
	}
	suite.Run(t, new(CategoryIntegrationTestSuite))
}

// ErrTestMethodNotImplemented возвращается когда тестовый метод не реализован
var ErrTestMethodNotImplemented = errors.New("test method not implemented")

// testStorage - простая реализация storage.Storage для тестов
type testStorage struct {
	db *sql.DB
}

func (ts *testStorage) QueryRow(ctx context.Context, query string, args ...interface{}) storage.Row {
	return ts.db.QueryRowContext(ctx, query, args...)
}

func (ts *testStorage) Query(ctx context.Context, query string, args ...interface{}) (storage.Rows, error) {
	return ts.db.QueryContext(ctx, query, args...) //nolint:rowserrcheck // rows.Err() checked by caller
}

func (ts *testStorage) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return ts.db.ExecContext(ctx, query, args...)
}

func (ts *testStorage) BeginTx(ctx context.Context, opts *sql.TxOptions) (storage.Transaction, error) {
	tx, err := ts.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &testTransaction{tx: tx}, nil
}

func (ts *testStorage) GetCategoryByID(ctx context.Context, id int) (*models.MarketplaceCategory, error) {
	var category models.MarketplaceCategory
	err := ts.QueryRow(ctx, `
		SELECT id, name, slug, parent_id, icon, has_custom_ui, custom_ui_component, 
		       description, is_active, seo_title, seo_description, seo_keywords, sort_order
		FROM marketplace_categories WHERE id = $1
	`, id).Scan(
		&category.ID, &category.Name, &category.Slug, &category.ParentID,
		&category.Icon, &category.HasCustomUI, &category.CustomUIComponent,
		&category.Description, &category.IsActive, &category.SEOTitle,
		&category.SEODescription, &category.SEOKeywords, &category.SortOrder,
	)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (ts *testStorage) GetCategories(ctx context.Context) ([]models.MarketplaceCategory, error) {
	rows, err := ts.Query(ctx, `
		SELECT id, name, slug, parent_id, icon, has_custom_ui, custom_ui_component, 
		       description, is_active, seo_title, seo_description, seo_keywords, sort_order
		FROM marketplace_categories WHERE is_active = true
		ORDER BY sort_order, name
	`)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			// Логируем ошибку закрытия rows в тесте
			_ = err // Explicitly ignore error
		}
	}()

	var categories []models.MarketplaceCategory
	for rows.Next() {
		var category models.MarketplaceCategory
		err := rows.Scan(
			&category.ID, &category.Name, &category.Slug, &category.ParentID,
			&category.Icon, &category.HasCustomUI, &category.CustomUIComponent,
			&category.Description, &category.IsActive, &category.SEOTitle,
			&category.SEODescription, &category.SEOKeywords, &category.SortOrder,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

// Остальные методы storage.Storage - заглушки для компиляции
func (ts *testStorage) GetOrCreateGoogleUser(ctx context.Context, user *models.User) (*models.User, error) {
	return nil, ErrTestMethodNotImplemented
}

func (ts *testStorage) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return nil, ErrTestMethodNotImplemented
}

func (ts *testStorage) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	return nil, ErrTestMethodNotImplemented
}

func (ts *testStorage) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	return nil, ErrTestMethodNotImplemented
}
func (ts *testStorage) UpdateUser(ctx context.Context, user *models.User) error { return nil }
func (ts *testStorage) GetUserProfile(ctx context.Context, id int) (*models.UserProfile, error) {
	return nil, ErrTestMethodNotImplemented
}

func (ts *testStorage) UpdateUserProfile(ctx context.Context, id int, update *models.UserProfileUpdate) error {
	return nil
}
func (ts *testStorage) UpdateLastSeen(ctx context.Context, id int) error { return nil }
func (ts *testStorage) GetFavoritedUsers(ctx context.Context, listingID int) ([]int, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) GetSession(ctx context.Context, token string) (*types.SessionData, error) {
	return nil, ErrTestMethodNotImplemented
}

func (ts *testStorage) CreateRefreshToken(ctx context.Context, token *models.RefreshToken) error {
	return nil
}

func (ts *testStorage) GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	return nil, ErrTestMethodNotImplemented
}

func (ts *testStorage) GetRefreshTokenByID(ctx context.Context, id int) (*models.RefreshToken, error) {
	return nil, ErrTestMethodNotImplemented
}

func (ts *testStorage) GetUserRefreshTokens(ctx context.Context, userID int) ([]*models.RefreshToken, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) UpdateRefreshToken(ctx context.Context, token *models.RefreshToken) error {
	return nil
}
func (ts *testStorage) RevokeRefreshToken(ctx context.Context, tokenID int) error { return nil }
func (ts *testStorage) RevokeRefreshTokenByValue(ctx context.Context, tokenValue string) error {
	return nil
}
func (ts *testStorage) RevokeUserRefreshTokens(ctx context.Context, userID int) error { return nil }
func (ts *testStorage) DeleteExpiredRefreshTokens(ctx context.Context) (int64, error) { return 0, nil }
func (ts *testStorage) CountActiveUserTokens(ctx context.Context, userID int) (int, error) {
	return 0, nil
}

func (ts *testStorage) GetAllUsers(ctx context.Context, limit, offset int) ([]*models.UserProfile, int, error) {
	return nil, 0, nil
}
func (ts *testStorage) UpdateUserStatus(ctx context.Context, id int, status string) error { return nil }
func (ts *testStorage) DeleteUser(ctx context.Context, id int) error                      { return nil }
func (ts *testStorage) IsUserAdmin(ctx context.Context, email string) (bool, error) {
	return false, nil
}

func (ts *testStorage) GetAllAdmins(ctx context.Context) ([]*models.AdminUser, error) {
	return nil, ErrNotImplemented
}
func (ts *testStorage) AddAdmin(ctx context.Context, admin *models.AdminUser) error { return nil }
func (ts *testStorage) RemoveAdmin(ctx context.Context, email string) error         { return nil }
func (ts *testStorage) DeleteListingAdmin(ctx context.Context, listingID int) error { return nil }
func (ts *testStorage) CreateReview(ctx context.Context, review *models.Review) (*models.Review, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) GetReviews(ctx context.Context, filter models.ReviewsFilter) ([]models.Review, int64, error) {
	return nil, 0, nil
}

func (ts *testStorage) GetReviewByID(ctx context.Context, id int) (*models.Review, error) {
	return nil, ErrNotImplemented
}
func (ts *testStorage) UpdateReview(ctx context.Context, review *models.Review) error { return nil }
func (ts *testStorage) UpdateReviewStatus(ctx context.Context, reviewId int, status string) error {
	return nil
}
func (ts *testStorage) DeleteReview(ctx context.Context, id int) error { return nil }
func (ts *testStorage) AddReviewResponse(ctx context.Context, response *models.ReviewResponse) error {
	return nil
}
func (ts *testStorage) AddReviewVote(ctx context.Context, vote *models.ReviewVote) error { return nil }
func (ts *testStorage) GetReviewVotes(ctx context.Context, reviewId int) (helpful int, notHelpful int, err error) {
	return 0, 0, nil
}

func (ts *testStorage) GetUserReviewVote(ctx context.Context, userId int, reviewId int) (string, error) {
	return "", nil
}

func (ts *testStorage) GetEntityRating(ctx context.Context, entityType string, entityId int) (float64, error) {
	return 0, nil
}

func (ts *testStorage) GetUserReviews(ctx context.Context, userID int, filter models.ReviewsFilter) ([]models.Review, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) GetStorefrontReviews(ctx context.Context, storefrontID int, filter models.ReviewsFilter) ([]models.Review, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) GetUserRatingSummary(ctx context.Context, userID int) (*models.UserRatingSummary, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) GetStorefrontRatingSummary(ctx context.Context, storefrontID int) (*models.StorefrontRatingSummary, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) GetNotificationSettings(ctx context.Context, userID int) ([]models.NotificationSettings, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) UpdateNotificationSettings(ctx context.Context, settings *models.NotificationSettings) error {
	return nil
}

func (ts *testStorage) SaveTelegramConnection(ctx context.Context, userID int, chatID string, username string) error {
	return nil
}

func (ts *testStorage) GetTelegramConnection(ctx context.Context, userID int) (*models.TelegramConnection, error) {
	return nil, ErrNotImplemented
}
func (ts *testStorage) DeleteTelegramConnection(ctx context.Context, userID int) error { return nil }
func (ts *testStorage) CreateNotification(ctx context.Context, notification *models.Notification) error {
	return nil
}

func (ts *testStorage) GetUserNotifications(ctx context.Context, userID int, limit, offset int) ([]models.Notification, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) MarkNotificationAsRead(ctx context.Context, userID int, notificationID int) error {
	return nil
}

func (ts *testStorage) DeleteNotification(ctx context.Context, userID int, notificationID int) error {
	return nil
}

func (ts *testStorage) CreateListing(ctx context.Context, listing *models.MarketplaceListing) (int, error) {
	return 0, nil
}

func (ts *testStorage) GetListings(ctx context.Context, filters map[string]string, limit int, offset int) ([]models.MarketplaceListing, int64, error) {
	return nil, 0, nil
}

func (ts *testStorage) GetListingByID(ctx context.Context, id int) (*models.MarketplaceListing, error) {
	return nil, ErrNotImplemented
}
func (ts *testStorage) IncrementViewsCount(ctx context.Context, id int) error { return nil }
func (ts *testStorage) UpdateListing(ctx context.Context, listing *models.MarketplaceListing) error {
	return nil
}
func (ts *testStorage) DeleteListing(ctx context.Context, id int, userID int) error { return nil }
func (ts *testStorage) GetCategoryTree(ctx context.Context) ([]models.CategoryTreeNode, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) AddListingImage(ctx context.Context, image *models.MarketplaceImage) (int, error) {
	return 0, nil
}

func (ts *testStorage) GetListingImages(ctx context.Context, listingID string) ([]models.MarketplaceImage, error) {
	return nil, ErrNotImplemented
}
func (ts *testStorage) FileStorage() filestorage.FileStorageInterface { return nil }
func (ts *testStorage) GetListingImageByID(ctx context.Context, imageID int) (*models.MarketplaceImage, error) {
	return nil, ErrNotImplemented
}
func (ts *testStorage) DeleteListingImage(ctx context.Context, imageID int) error { return nil }
func (ts *testStorage) GetAttributeOptionTranslations(ctx context.Context, attributeName, optionValue string) (map[string]string, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) GetAttributeRanges(ctx context.Context, categoryID int) (map[string]map[string]interface{}, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) GetCategoryAttributes(ctx context.Context, categoryID int) ([]models.CategoryAttribute, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) GetAllCategories(ctx context.Context) ([]models.MarketplaceCategory, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) AddToFavorites(ctx context.Context, userID int, listingID int) error {
	return nil
}

func (ts *testStorage) RemoveFromFavorites(ctx context.Context, userID int, listingID int) error {
	return nil
}

func (ts *testStorage) GetUserFavorites(ctx context.Context, userID int) ([]models.MarketplaceListing, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) GetPriceHistory(ctx context.Context, listingID int) ([]models.PriceHistoryEntry, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) AddPriceHistoryEntry(ctx context.Context, entry *models.PriceHistoryEntry) error {
	return nil
}
func (ts *testStorage) ClosePriceHistoryEntry(ctx context.Context, listingID int) error { return nil }
func (ts *testStorage) CheckPriceManipulation(ctx context.Context, listingID int) (bool, error) {
	return false, nil
}

func (ts *testStorage) SaveListingAttributes(ctx context.Context, listingID int, attributes []models.ListingAttributeValue) error {
	return nil
}

func (ts *testStorage) GetListingAttributes(ctx context.Context, listingID int) ([]models.ListingAttributeValue, error) {
	return nil, ErrNotImplemented
}
func (ts *testStorage) SynchronizeDiscountMetadata(ctx context.Context) error { return nil }
func (ts *testStorage) GetUserBalance(ctx context.Context, userID int) (*models.UserBalance, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) GetUserTransactions(ctx context.Context, userID int, limit, offset int) ([]models.BalanceTransaction, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) CreateTransaction(ctx context.Context, transaction *models.BalanceTransaction) (int, error) {
	return 0, nil
}

func (ts *testStorage) GetActivePaymentMethods(ctx context.Context) ([]models.PaymentMethod, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) UpdateBalance(ctx context.Context, userID int, amount float64) error {
	return nil
}

func (ts *testStorage) CreateMessage(ctx context.Context, msg *models.MarketplaceMessage) error {
	return nil
}

func (ts *testStorage) GetMessages(ctx context.Context, listingID int, userID int, offset int, limit int) ([]models.MarketplaceMessage, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) GetChats(ctx context.Context, userID int) ([]models.MarketplaceChat, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) GetChat(ctx context.Context, chatID int, userID int) (*models.MarketplaceChat, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) MarkMessagesAsRead(ctx context.Context, messageIDs []int, userID int) error {
	return nil
}
func (ts *testStorage) ArchiveChat(ctx context.Context, chatID int, userID int) error { return nil }
func (ts *testStorage) GetUnreadMessagesCount(ctx context.Context, userID int) (int, error) {
	return 0, nil
}

func (ts *testStorage) CreateChatAttachment(ctx context.Context, attachment *models.ChatAttachment) error {
	return nil
}

func (ts *testStorage) GetChatAttachment(ctx context.Context, attachmentID int) (*models.ChatAttachment, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) GetMessageAttachments(ctx context.Context, messageID int) ([]*models.ChatAttachment, error) {
	return nil, ErrNotImplemented
}
func (ts *testStorage) DeleteChatAttachment(ctx context.Context, attachmentID int) error { return nil }
func (ts *testStorage) UpdateMessageAttachmentsCount(ctx context.Context, messageID int, count int) error {
	return nil
}

func (ts *testStorage) GetMessageByID(ctx context.Context, messageID int) (*models.MarketplaceMessage, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) GetChatActivityStats(ctx context.Context, buyerID int, sellerID int, listingID int) (*models.ChatActivityStats, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) GetUserAggregatedRating(ctx context.Context, userID int) (*models.UserAggregatedRating, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) GetStorefrontAggregatedRating(ctx context.Context, storefrontID int) (*models.StorefrontAggregatedRating, error) {
	return nil, ErrNotImplemented
}
func (ts *testStorage) RefreshRatingViews(ctx context.Context) error { return nil }
func (ts *testStorage) CreateReviewConfirmation(ctx context.Context, confirmation *models.ReviewConfirmation) error {
	return nil
}

func (ts *testStorage) GetReviewConfirmation(ctx context.Context, reviewID int) (*models.ReviewConfirmation, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) CreateReviewDispute(ctx context.Context, dispute *models.ReviewDispute) error {
	return nil
}

func (ts *testStorage) GetReviewDispute(ctx context.Context, reviewID int) (*models.ReviewDispute, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) UpdateReviewDispute(ctx context.Context, dispute *models.ReviewDispute) error {
	return nil
}

func (ts *testStorage) CanUserReviewEntity(ctx context.Context, userID int, entityType string, entityID int) (*models.CanReviewResponse, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) CreateStorefront(ctx context.Context, userID int, dto *models.StorefrontCreateDTO) (*models.Storefront, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) GetUserStorefronts(ctx context.Context, userID int) ([]models.Storefront, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) GetStorefrontByID(ctx context.Context, id int) (*models.Storefront, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) UpdateStorefront(ctx context.Context, storefront *models.Storefront) error {
	return nil
}
func (ts *testStorage) DeleteStorefront(ctx context.Context, id int) error { return nil }
func (ts *testStorage) Storefront() interface{}                            { return nil }
func (ts *testStorage) Cart() interface{}                                  { return nil }
func (ts *testStorage) Order() interface{}                                 { return nil }
func (ts *testStorage) Inventory() interface{}                             { return nil }
func (ts *testStorage) MarketplaceOrder() interface{}                      { return nil }
func (ts *testStorage) StorefrontProductSearch() interface{}               { return nil }
func (ts *testStorage) SearchListings(ctx context.Context, params *search.SearchParams) (*search.SearchResult, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) SearchListingsOpenSearch(ctx context.Context, params *search.SearchParams) (*search.SearchResult, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) SuggestListings(ctx context.Context, prefix string, size int) ([]string, error) {
	return nil, ErrNotImplemented
}
func (ts *testStorage) ReindexAllListings(ctx context.Context) error { return nil }
func (ts *testStorage) IndexListing(ctx context.Context, listing *models.MarketplaceListing) error {
	return nil
}
func (ts *testStorage) DeleteListingIndex(ctx context.Context, id string) error { return nil }
func (ts *testStorage) PrepareIndex(ctx context.Context) error                  { return nil }
func (ts *testStorage) SearchStorefrontsOpenSearch(ctx context.Context, params *opensearch.StorefrontSearchParams) (*opensearch.StorefrontSearchResult, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) IndexStorefront(ctx context.Context, storefront *models.Storefront) error {
	return nil
}
func (ts *testStorage) DeleteStorefrontIndex(ctx context.Context, storefrontID int) error { return nil }
func (ts *testStorage) ReindexAllStorefronts(ctx context.Context) error                   { return nil }
func (ts *testStorage) GetTranslationsForEntity(ctx context.Context, entityType string, entityID int) ([]models.Translation, error) {
	return nil, ErrNotImplemented
}
func (ts *testStorage) AddContact(ctx context.Context, contact *models.UserContact) error { return nil }
func (ts *testStorage) UpdateContactStatus(ctx context.Context, userID, contactUserID int, status, notes string) error {
	return nil
}

func (ts *testStorage) GetContact(ctx context.Context, userID, contactUserID int) (*models.UserContact, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) GetUserContacts(ctx context.Context, userID int, status string, page, limit int) ([]models.UserContact, int, error) {
	return nil, 0, nil
}

func (ts *testStorage) RemoveContact(ctx context.Context, userID, contactUserID int) error {
	return nil
}

func (ts *testStorage) GetUserPrivacySettings(ctx context.Context, userID int) (*models.UserPrivacySettings, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) UpdateUserPrivacySettings(ctx context.Context, userID int, settings *models.UpdatePrivacySettingsRequest) error {
	return nil
}

func (ts *testStorage) CanAddContact(ctx context.Context, userID, targetUserID int) (bool, error) {
	return false, nil
}

func (ts *testStorage) ExpandSearchQuery(ctx context.Context, query string, language string) (string, error) {
	return query, nil
}

func (ts *testStorage) SearchCategoriesFuzzy(ctx context.Context, searchTerm string, language string, similarityThreshold float64) ([]interface{}, error) {
	return []interface{}{}, nil
}

func (ts *testStorage) SearchCategories(ctx context.Context, query string, limit int) ([]models.MarketplaceCategory, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) GetPopularSearchQueries(ctx context.Context, query string, limit int) ([]interface{}, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) SaveSearchQuery(ctx context.Context, query, normalized string, resultsCount int, language string) error {
	return nil
}

func (ts *testStorage) SearchListingsAdvanced(ctx context.Context, params interface{}) (interface{}, error) {
	return nil, ErrNotImplemented
}
func (ts *testStorage) Close()                         {}
func (ts *testStorage) Ping(ctx context.Context) error { return nil }
func (ts *testStorage) GetCarMakeBySlug(ctx context.Context, slug string) (*models.CarMake, error) {
	return nil, ErrNotImplemented
}

// Marketplace specific methods
func (ts *testStorage) GenerateUniqueSlug(ctx context.Context, baseSlug string, excludeID int) (string, error) {
	return baseSlug, nil
}

func (ts *testStorage) GetListingBySlug(ctx context.Context, slug string) (*models.MarketplaceListing, error) {
	return nil, sql.ErrNoRows
}

func (ts *testStorage) IsSlugUnique(ctx context.Context, slug string, excludeID int) (bool, error) {
	return true, nil
}

// Car-related methods
func (ts *testStorage) GetCarMakes(ctx context.Context, country string, isDomestic bool, isMotorcycle bool, activeOnly bool) ([]models.CarMake, error) {
	return []models.CarMake{}, nil
}

func (ts *testStorage) GetCarModelsByMake(ctx context.Context, makeSlug string, activeOnly bool) ([]models.CarModel, error) {
	return []models.CarModel{}, nil
}

func (ts *testStorage) GetCarGenerationsByModel(ctx context.Context, modelID int, activeOnly bool) ([]models.CarGeneration, error) {
	return []models.CarGeneration{}, nil
}

func (ts *testStorage) SearchCarMakes(ctx context.Context, query string, limit int) ([]models.CarMake, error) {
	return []models.CarMake{}, nil
}

// Additional mock methods to satisfy storage.Storage interface
func (ts *testStorage) GetAllRoles(ctx context.Context) ([]*models.Role, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) UpdateUserRole(ctx context.Context, id int, roleID int) error {
	return ErrNotImplemented
}

func (ts *testStorage) GetAllUsersWithSort(ctx context.Context, limit, offset int, sortBy, sortOrder, statusFilter string) ([]*models.UserProfile, int, error) {
	return nil, 0, ErrNotImplemented
}

func (ts *testStorage) GetStorefrontOwnerByProductID(ctx context.Context, productID int) (int, error) {
	return 0, ErrNotImplemented
}

func (ts *testStorage) GetPrivacySettings(ctx context.Context, userID int) (*models.UserPrivacySettings, error) {
	return nil, ErrNotImplemented
}

func (ts *testStorage) UpdatePrivacySettings(ctx context.Context, userID int, settings *models.UpdatePrivacySettingsRequest) error {
	return ErrNotImplemented
}

func (ts *testStorage) GetPopularCategories(ctx context.Context, limit int) ([]models.MarketplaceCategory, error) {
	return nil, ErrNotImplemented
}

// Listing Variants methods
func (ts *testStorage) CreateListingVariants(ctx context.Context, listingID int, variants []models.MarketplaceListingVariant) error {
	return nil
}

func (ts *testStorage) GetListingVariants(ctx context.Context, listingID int) ([]models.MarketplaceListingVariant, error) {
	return []models.MarketplaceListingVariant{}, nil
}

func (ts *testStorage) UpdateListingVariant(ctx context.Context, variant *models.MarketplaceListingVariant) error {
	return nil
}

func (ts *testStorage) DeleteListingVariant(ctx context.Context, variantID int) error {
	return nil
}

func (ts *testStorage) GetIncomingContactRequests(ctx context.Context, userID int, page, limit int) ([]models.UserContact, int, error) {
	return []models.UserContact{}, 0, nil
}

func (ts *testStorage) AddStorefrontToFavorites(ctx context.Context, userID int, productID int) error {
	return nil
}

func (ts *testStorage) RemoveStorefrontFromFavorites(ctx context.Context, userID int, productID int) error {
	return nil
}

func (ts *testStorage) GetUserStorefrontFavorites(ctx context.Context, userID int) ([]models.MarketplaceListing, error) {
	return []models.MarketplaceListing{}, nil
}

// GetCarListingsCount - недостающий метод для подсчета автомобильных объявлений
func (ts *testStorage) GetCarListingsCount(ctx context.Context) (int, error) {
	var count int
	err := ts.QueryRow(ctx, `SELECT COUNT(*) FROM marketplace_listings WHERE category_id = 1301`).Scan(&count)
	return count, err
}

// GetTotalCarModelsCount - недостающий метод для подсчета моделей автомобилей
func (ts *testStorage) GetTotalCarModelsCount(ctx context.Context) (int, error) {
	// В тестах возвращаем фиксированное значение
	return 0, nil
}

// testTransaction - простая реализация транзакции для тестов
type testTransaction struct {
	tx *sql.Tx
}

func (tt *testTransaction) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return tt.tx.ExecContext(ctx, query, args...)
}

func (tt *testTransaction) QueryRow(ctx context.Context, query string, args ...interface{}) storage.Row {
	return tt.tx.QueryRowContext(ctx, query, args...)
}

func (tt *testTransaction) Query(ctx context.Context, query string, args ...interface{}) (storage.Rows, error) {
	return tt.tx.QueryContext(ctx, query, args...) //nolint:rowserrcheck // rows.Err() checked by caller
}

func (tt *testTransaction) Commit() error {
	return tt.tx.Commit()
}

func (tt *testTransaction) Rollback() error {
	return tt.tx.Rollback()
}

// dummyTranslationService - заглушка для сервиса перевода
type dummyTranslationService struct{}

func (d *dummyTranslationService) TranslateText(ctx context.Context, text, sourceLanguage, targetLanguage string) (string, error) {
	return text, nil
}

func (d *dummyTranslationService) SaveTranslation(ctx context.Context, entityType string, entityID int, language, fieldName, translatedText string, metadata map[string]any) error {
	return nil
}

func (d *dummyTranslationService) DetectLanguage(ctx context.Context, text string) (string, float64, error) {
	return "en", 1.0, nil
}

func (d *dummyTranslationService) Translate(ctx context.Context, text string, sourceLanguage string, targetLanguage string) (string, error) {
	return text, nil
}

func (d *dummyTranslationService) ModerateText(ctx context.Context, text string, language string) (string, error) {
	return text, nil
}

func (d *dummyTranslationService) TranslateToAllLanguages(ctx context.Context, text string) (map[string]string, error) {
	return map[string]string{"en": text, "ru": text, "sr": text}, nil
}

func (d *dummyTranslationService) TranslateEntityFields(ctx context.Context, sourceLanguage string, targetLanguages []string, fields map[string]string) (map[string]map[string]string, error) {
	return nil, ErrNotImplemented
}

func (d *dummyTranslationService) TranslateWithContext(ctx context.Context, text string, sourceLanguage string, targetLanguage string, context string, fieldName string) (string, error) {
	return text, nil
}

func (d *dummyTranslationService) TranslateWithToneModeration(ctx context.Context, text string, sourceLanguage string, targetLanguage string, moderateTone bool) (string, error) {
	return text, nil
}

// memoryCache - простой вариант кеша в памяти
type memoryCache struct {
	data map[string][]byte
}

func (m *memoryCache) Get(ctx context.Context, key string, dest interface{}) error {
	return fmt.Errorf("not found")
}

func (m *memoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return nil
}

func (m *memoryCache) Delete(ctx context.Context, keys ...string) error {
	return nil
}

func (m *memoryCache) DeletePattern(ctx context.Context, pattern string) error {
	return nil
}

func (m *memoryCache) GetOrSet(ctx context.Context, key string, dest interface{}, ttl time.Duration, loader func() (interface{}, error)) error {
	result, err := loader()
	if err != nil {
		return err
	}
	// Просто копируем результат в dest
	if result != nil {
		switch d := dest.(type) {
		case *[]models.MarketplaceCategory:
			if categories, ok := result.([]models.MarketplaceCategory); ok {
				*d = categories
			}
		case *[]*models.AttributeGroup:
			if groups, ok := result.([]*models.AttributeGroup); ok {
				*d = groups
			}
		}
	}
	return nil
}

// mockTranslationService - мок сервиса перевода для тестов
type mockTranslationService struct {
	translations map[string]string
}

func (m *mockTranslationService) TranslateText(ctx context.Context, text, source, target string) (string, error) {
	if translated, ok := m.translations[target]; ok {
		return translated, nil
	}
	return text, nil
}

func (m *mockTranslationService) SaveTranslation(ctx context.Context, entityType string, entityID int, language, fieldName, translatedText string, metadata map[string]any) error {
	return nil
}

func (m *mockTranslationService) DetectLanguage(ctx context.Context, text string) (string, float64, error) {
	return "en", 1.0, nil
}

func (m *mockTranslationService) Translate(ctx context.Context, text string, sourceLanguage string, targetLanguage string) (string, error) {
	return text, nil
}

func (m *mockTranslationService) ModerateText(ctx context.Context, text string, language string) (string, error) {
	return text, nil
}

func (m *mockTranslationService) TranslateToAllLanguages(ctx context.Context, text string) (map[string]string, error) {
	return m.translations, nil
}

func (m *mockTranslationService) TranslateEntityFields(ctx context.Context, sourceLanguage string, targetLanguages []string, fields map[string]string) (map[string]map[string]string, error) {
	return nil, ErrNotImplemented
}

func (m *mockTranslationService) TranslateWithContext(ctx context.Context, text string, sourceLanguage string, targetLanguage string, context string, fieldName string) (string, error) {
	if translated, ok := m.translations[targetLanguage]; ok {
		return translated, nil
	}
	return text, nil
}

func (m *mockTranslationService) TranslateWithToneModeration(ctx context.Context, text string, sourceLanguage string, targetLanguage string, moderateTone bool) (string, error) {
	if translated, ok := m.translations[targetLanguage]; ok {
		return translated, nil
	}
	return text, nil
}
