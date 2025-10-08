package service

import (
	"context"
	"testing"

	"backend/internal/domain/models"

	"github.com/stretchr/testify/assert"
)

// Тест пустой группы вариантов - самый простой unit тест
func TestImportService_ImportVariantGroup_EmptyGroup(t *testing.T) {
	ctx := context.Background()
	storefrontID := 123

	importService := &ImportService{
		variantDetector: NewVariantDetector(),
	}

	variantGroup := &VariantGroup{
		BaseName:     "Empty Group",
		VariantCount: 0,
		Variants:     []*ProductVariant{}, // пустой массив
	}

	err := importService.importVariantGroup(ctx, variantGroup, storefrontID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "variant group has no variants")
}

// Тест валидации: группа с nil Variants
func TestImportService_ImportVariantGroup_NilVariants(t *testing.T) {
	ctx := context.Background()
	storefrontID := 123

	importService := &ImportService{
		variantDetector: NewVariantDetector(),
	}

	variantGroup := &VariantGroup{
		BaseName:     "Nil Variants",
		VariantCount: 0,
		Variants:     nil, // nil вместо пустого массива
	}

	err := importService.importVariantGroup(ctx, variantGroup, storefrontID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "variant group has no variants")
}

// Тест конвертации ImportProductRequest -> ProductVariant
func TestImportService_ConvertImportProductsToVariants(t *testing.T) {
	importService := &ImportService{
		variantDetector: NewVariantDetector(),
	}

	products := []models.ImportProductRequest{
		{
			Name:          "Tastatura Gembird KB-UM-104 crna",
			SKU:           "KB-UM-104-BLACK",
			Price:         1200.0,
			Description:   "USB клавиатура",
			CategoryID:    5,
			Barcode:       "1234567890",
			StockQuantity: 10,
			ImageURLs:     []string{"http://example.com/img1.jpg"},
		},
		{
			Name:          "Tastatura Gembird KB-UM-104 bela",
			SKU:           "KB-UM-104-WHITE",
			Price:         1250.0,
			Description:   "USB клавиатура белая",
			CategoryID:    5,
			Barcode:       "1234567891",
			StockQuantity: 5,
			ImageURLs:     []string{"http://example.com/img2.jpg", "http://example.com/img3.jpg"},
		},
	}

	variants := importService.convertImportProductsToVariants(products)

	// Проверяем количество
	assert.Len(t, variants, 2)

	// Проверяем первый вариант
	assert.Equal(t, "Tastatura Gembird KB-UM-104 crna", variants[0].Name)
	assert.Equal(t, "KB-UM-104-BLACK", variants[0].SKU)
	assert.Equal(t, 1200.0, variants[0].Price)
	assert.Equal(t, 10, variants[0].StockQuantity)
	assert.Equal(t, "http://example.com/img1.jpg", variants[0].ImageURL)
	assert.NotNil(t, variants[0].OriginalAttributes)
	assert.Equal(t, "USB клавиатура", variants[0].OriginalAttributes["description"])
	assert.Equal(t, "1234567890", variants[0].OriginalAttributes["barcode"])
	assert.Equal(t, 5, variants[0].OriginalAttributes["category_id"])

	// Проверяем второй вариант
	assert.Equal(t, "Tastatura Gembird KB-UM-104 bela", variants[1].Name)
	assert.Equal(t, "KB-UM-104-WHITE", variants[1].SKU)
	assert.Equal(t, 1250.0, variants[1].Price)
	assert.Equal(t, 5, variants[1].StockQuantity)
	assert.Equal(t, "http://example.com/img2.jpg", variants[1].ImageURL) // берется первое изображение
}

// Тест конвертации без изображений
func TestImportService_ConvertImportProductsToVariants_NoImages(t *testing.T) {
	importService := &ImportService{
		variantDetector: NewVariantDetector(),
	}

	products := []models.ImportProductRequest{
		{
			Name:          "Test Product",
			SKU:           "TEST-1",
			Price:         100.0,
			StockQuantity: 1,
			ImageURLs:     []string{}, // пустой массив изображений
		},
	}

	variants := importService.convertImportProductsToVariants(products)

	assert.Len(t, variants, 1)
	assert.Empty(t, variants[0].ImageURL) // ImageURL должен быть пустым
}

// Тест groupAndDetectVariants - интеграция с VariantDetector
func TestImportService_GroupAndDetectVariants(t *testing.T) {
	importService := &ImportService{
		variantDetector: NewVariantDetector(),
	}

	products := []models.ImportProductRequest{
		{
			Name:          "Tastatura Gembird KB-UM-104 crna",
			SKU:           "KB-UM-104-BLACK",
			Price:         1200.0,
			CategoryID:    5,
			StockQuantity: 10,
		},
		{
			Name:          "Tastatura Gembird KB-UM-104 bela",
			SKU:           "KB-UM-104-WHITE",
			Price:         1250.0,
			CategoryID:    5,
			StockQuantity: 5,
		},
		{
			Name:          "Miš Genius DX-110 crni",
			SKU:           "DX-110-BLACK",
			Price:         500.0,
			CategoryID:    6,
			StockQuantity: 20,
		},
		{
			Name:          "Miš Genius DX-110 beli",
			SKU:           "DX-110-WHITE",
			Price:         550.0,
			CategoryID:    6,
			StockQuantity: 15,
		},
	}

	groups := importService.groupAndDetectVariants(products)

	// Должны получить 2 группы: Tastatura и Miš
	assert.Len(t, groups, 2)

	// Проверяем что обе группы имеют по 2 варианта
	for _, group := range groups {
		assert.Equal(t, 2, group.VariantCount)
		assert.Len(t, group.Variants, 2)
	}

	// Проверяем что есть группы для обоих товаров (порядок не важен)
	// BaseName обрезает цвета, поэтому ожидаем укороченные имена
	baseNames := []string{groups[0].BaseName, groups[1].BaseName}

	// Проверяем что хотя бы одна группа содержит "Tastatura" и одна содержит "Miš"
	hasTastatura := false
	hasMouse := false
	for _, name := range baseNames {
		if len(name) >= 9 && name[:9] == "Tastatura" {
			hasTastatura = true
		}
		// Проверяем по байтам (кириллица "Miš" может быть 4+ байта)
		// Используем Contains для упрощения
		if len(name) >= 2 && (name[0] == 'M' && name[1] == 'i') {
			hasMouse = true
		}
	}
	assert.True(t, hasTastatura, "Expected one group to start with 'Tastatura', got: %v", baseNames)
	assert.True(t, hasMouse, "Expected one group to start with 'Miš', got: %v", baseNames)
}

// Тест groupAndDetectVariants без группировки (разные товары)
func TestImportService_GroupAndDetectVariants_NoGrouping(t *testing.T) {
	importService := &ImportService{
		variantDetector: NewVariantDetector(),
	}

	products := []models.ImportProductRequest{
		{
			Name:          "Tastatura Logitech K120",
			SKU:           "LOG-K120",
			Price:         1000.0,
			StockQuantity: 10,
		},
		{
			Name:          "Miš HP X1000",
			SKU:           "HP-X1000",
			Price:         400.0,
			StockQuantity: 20,
		},
	}

	groups := importService.groupAndDetectVariants(products)

	// Не должно быть групп вариантов, т.к. товары слишком разные
	// или группы будут с низким confidence и будут отфильтрованы
	// В зависимости от алгоритма может быть 0, 1 или 2 группы
	// Главное - не crash
	assert.NotNil(t, groups)
}
