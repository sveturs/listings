package cache

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRedisCacheRealIntegration тестирует работу с реальным Redis
func TestRedisCacheRealIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Создаем логгер
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Подключаемся к реальному Redis
	cache, err := NewRedisCache("localhost:6379", "", 0, 10, logger)
	if err != nil {
		t.Skipf("Skipping test: Redis not available: %v", err)
	}
	defer func() {
		if err := cache.Close(); err != nil {
			t.Logf("Failed to close cache: %v", err)
		}
	}()

	ctx := context.Background()

	t.Run("Test Category Caching", func(t *testing.T) {
		// Симулируем реальные данные категорий
		type Category struct {
			ID          int               `json:"id"`
			Name        string            `json:"name"`
			Description string            `json:"description"`
			ParentID    *int              `json:"parent_id"`
			IsActive    bool              `json:"is_active"`
			Attributes  map[string]string `json:"attributes"`
		}

		categories := []Category{
			{
				ID:          1,
				Name:        "Electronics",
				Description: "Electronic devices and accessories",
				ParentID:    nil,
				IsActive:    true,
				Attributes: map[string]string{
					"icon": "electronics-icon",
					"slug": "electronics",
				},
			},
			{
				ID:          2,
				Name:        "Smartphones",
				Description: "Mobile phones and tablets",
				ParentID:    func() *int { i := 1; return &i }(),
				IsActive:    true,
				Attributes: map[string]string{
					"icon": "phone-icon",
					"slug": "smartphones",
				},
			},
		}

		// Тестируем кеширование для разных языков
		locales := []string{"en", "ru", "sr"}

		for _, locale := range locales {
			t.Run(fmt.Sprintf("Locale_%s", locale), func(t *testing.T) {
				// Формируем ключ кеша
				cacheKey := BuildCategoriesKey(locale)

				// Сохраняем в кеш
				err := cache.Set(ctx, cacheKey, categories, 1*time.Hour)
				require.NoError(t, err)

				// Получаем из кеша
				var cachedCategories []Category
				err = cache.Get(ctx, cacheKey, &cachedCategories)
				require.NoError(t, err)

				// Проверяем данные
				assert.Len(t, cachedCategories, 2)
				assert.Equal(t, categories[0].Name, cachedCategories[0].Name)
				assert.Equal(t, categories[1].Name, cachedCategories[1].Name)

				// Проверяем что ключ существует
				exists, err := cache.Exists(ctx, cacheKey)
				require.NoError(t, err)
				assert.True(t, exists)

				t.Logf("Successfully cached and retrieved categories for locale: %s", locale)
			})
		}

		// Очищаем все категории
		pattern := BuildAllCategoriesInvalidationPattern()
		err = cache.DeletePattern(ctx, pattern)
		require.NoError(t, err)

		// Проверяем что все ключи удалены
		for _, locale := range locales {
			key := BuildCategoriesKey(locale)
			exists, err := cache.Exists(ctx, key)
			require.NoError(t, err)
			assert.False(t, exists)
		}
	})

	t.Run("Test GetOrSet Functionality", func(t *testing.T) {
		key := "test:getorset:categories"
		loadCount := 0

		// Функция загрузки данных
		loader := func() (interface{}, error) {
			loadCount++
			t.Log("Loader called, simulating database query...")

			// Симуляция загрузки из БД
			return map[string]interface{}{
				"id":     100,
				"name":   "Test Category",
				"loaded": time.Now().Format(time.RFC3339),
			}, nil
		}

		// Первый вызов - должен вызвать loader
		var result1 map[string]interface{}
		err := cache.GetOrSet(ctx, key, &result1, 30*time.Second, loader)
		require.NoError(t, err)
		assert.Equal(t, 1, loadCount)
		assert.Equal(t, 100, int(result1["id"].(float64)))

		// Второй вызов - должен получить из кеша
		var result2 map[string]interface{}
		err = cache.GetOrSet(ctx, key, &result2, 30*time.Second, loader)
		require.NoError(t, err)
		assert.Equal(t, 1, loadCount)                         // loader не должен вызываться
		assert.Equal(t, result1["loaded"], result2["loaded"]) // время должно быть одинаковое

		// Очищаем ключ
		err = cache.Delete(ctx, key)
		require.NoError(t, err)
	})

	t.Run("Test Cache Performance", func(t *testing.T) {
		// Массив тестовых данных
		testData := make([]map[string]interface{}, 100)
		for i := 0; i < 100; i++ {
			testData[i] = map[string]interface{}{
				"id":          i,
				"name":        fmt.Sprintf("Item %d", i),
				"description": fmt.Sprintf("Description for item %d", i),
				"price":       float64(i) * 10.99,
				"active":      i%2 == 0,
			}
		}

		// Измеряем время записи
		startWrite := time.Now()
		for i, data := range testData {
			key := fmt.Sprintf("perf:test:%d", i)
			err := cache.Set(ctx, key, data, 5*time.Minute)
			require.NoError(t, err)
		}
		writeTime := time.Since(startWrite)

		// Измеряем время чтения
		startRead := time.Now()
		for i := range testData {
			key := fmt.Sprintf("perf:test:%d", i)
			var result map[string]interface{}
			err := cache.Get(ctx, key, &result)
			require.NoError(t, err)
		}
		readTime := time.Since(startRead)

		t.Logf("Performance results:")
		t.Logf("  Write 100 items: %v", writeTime)
		t.Logf("  Read 100 items: %v", readTime)
		t.Logf("  Avg write time: %v", writeTime/100)
		t.Logf("  Avg read time: %v", readTime/100)

		// Очищаем тестовые данные
		err = cache.DeletePattern(ctx, "perf:test:*")
		require.NoError(t, err)
	})

	t.Run("Test Attribute Groups Caching", func(t *testing.T) {
		type AttributeGroup struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			DisplayRank int    `json:"display_rank"`
			Attributes  []struct {
				Key  string `json:"key"`
				Name string `json:"name"`
				Type string `json:"type"`
			} `json:"attributes"`
		}

		groups := []AttributeGroup{
			{
				ID:          1,
				Name:        "Basic Info",
				DisplayRank: 1,
				Attributes: []struct {
					Key  string `json:"key"`
					Name string `json:"name"`
					Type string `json:"type"`
				}{
					{Key: "brand", Name: "Brand", Type: "text"},
					{Key: "model", Name: "Model", Type: "text"},
				},
			},
			{
				ID:          2,
				Name:        "Technical Specs",
				DisplayRank: 2,
				Attributes: []struct {
					Key  string `json:"key"`
					Name string `json:"name"`
					Type string `json:"type"`
				}{
					{Key: "screen_size", Name: "Screen Size", Type: "number"},
					{Key: "ram", Name: "RAM", Type: "select"},
				},
			},
		}

		categoryID := int64(5)
		locale := "en"
		key := BuildAttributeGroupsKey(categoryID, locale)

		// Сохраняем группы атрибутов
		err := cache.Set(ctx, key, groups, 6*time.Hour)
		require.NoError(t, err)

		// Получаем из кеша
		var cachedGroups []AttributeGroup
		err = cache.Get(ctx, key, &cachedGroups)
		require.NoError(t, err)

		// Проверяем данные
		assert.Len(t, cachedGroups, 2)
		assert.Equal(t, "Basic Info", cachedGroups[0].Name)
		assert.Len(t, cachedGroups[0].Attributes, 2)
		assert.Equal(t, "Technical Specs", cachedGroups[1].Name)
		assert.Len(t, cachedGroups[1].Attributes, 2)

		// Очищаем
		err = cache.Delete(ctx, key)
		require.NoError(t, err)
	})

	t.Run("Test Cache Invalidation Patterns", func(t *testing.T) {
		// Создаем разные типы ключей для категории
		categoryID := int64(10)

		// Атрибуты категории для разных локалей
		attrKeyEn := BuildCategoryAttributesKey(categoryID, "en")
		attrKeyRu := BuildCategoryAttributesKey(categoryID, "ru")

		// Группы атрибутов
		groupKeyEn := BuildAttributeGroupsKey(categoryID, "en")
		groupKeyRu := BuildAttributeGroupsKey(categoryID, "ru")

		// Сам ключ категории
		categoryKey := BuildCategoryKey(categoryID)

		// Сохраняем все ключи
		testValue := map[string]string{"test": "data"}
		require.NoError(t, cache.Set(ctx, attrKeyEn, testValue, time.Hour))
		require.NoError(t, cache.Set(ctx, attrKeyRu, testValue, time.Hour))
		require.NoError(t, cache.Set(ctx, groupKeyEn, testValue, time.Hour))
		require.NoError(t, cache.Set(ctx, groupKeyRu, testValue, time.Hour))
		require.NoError(t, cache.Set(ctx, categoryKey, testValue, time.Hour))

		// Проверяем что все существуют
		for _, key := range []string{attrKeyEn, attrKeyRu, groupKeyEn, groupKeyRu, categoryKey} {
			exists, err := cache.Exists(ctx, key)
			require.NoError(t, err)
			assert.True(t, exists, "Key should exist: %s", key)
		}

		// Инвалидируем все ключи категории
		pattern := BuildCategoryInvalidationPattern(categoryID)
		err = cache.DeletePattern(ctx, pattern)
		require.NoError(t, err)

		// Проверяем что все удалены
		for _, key := range []string{attrKeyEn, attrKeyRu, groupKeyEn, groupKeyRu} {
			exists, err := cache.Exists(ctx, key)
			require.NoError(t, err)
			assert.False(t, exists, "Key should be deleted: %s", key)
		}
	})
}

// BenchmarkRedisCacheOperations бенчмарк для реального Redis
func BenchmarkRedisCacheOperations(b *testing.B) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	cache, err := NewRedisCache("localhost:6379", "", 0, 10, logger)
	if err != nil {
		b.Skipf("Skipping benchmark: Redis not available: %v", err)
	}
	defer func() {
		if err := cache.Close(); err != nil {
			b.Logf("Failed to close cache: %v", err)
		}
	}()

	ctx := context.Background()

	b.Run("Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("bench:set:%d", i%1000)
			value := map[string]interface{}{
				"id":    i,
				"name":  fmt.Sprintf("Item %d", i),
				"price": float64(i) * 10.5,
			}
			_ = cache.Set(ctx, key, value, 5*time.Minute)
		}
	})

	b.Run("Get", func(b *testing.B) {
		// Подготовка данных
		key := "bench:get:key"
		value := map[string]interface{}{
			"id":    1,
			"name":  "Benchmark Item",
			"price": 99.99,
		}
		_ = cache.Set(ctx, key, value, 5*time.Minute)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var result map[string]interface{}
			_ = cache.Get(ctx, key, &result)
		}
	})

	b.Run("GetOrSet", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("bench:getorset:%d", i%100)
			var result map[string]interface{}
			_ = cache.GetOrSet(ctx, key, &result, 5*time.Minute, func() (interface{}, error) {
				return map[string]interface{}{
					"id":   i,
					"name": fmt.Sprintf("Generated %d", i),
				}, nil
			})
		}
	})
}
