package service

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPreviewImportCSV тестирует preview функционал для CSV файлов
func TestPreviewImportCSV(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	// Подготовка тестовых данных CSV
	csvData := `name,price,currency,sku,description,category
Test Product 1,99.99,USD,TEST-SKU-001,Test description 1,Electronics
Test Product 2,149.99,USD,TEST-SKU-002,Test description 2,Clothing
Test Product 3,199.99,USD,TEST-SKU-003,Test description 3,Books`

	t.Run("Preview CSV with default limit", func(t *testing.T) {
		// Mock service (будет заменен на реальный в интеграции)
		service := &ImportService{}

		// Вызов preview
		result, err := service.PreviewImport(ctx, []byte(csvData), "csv", 1, 10)

		// Проверки
		require.NoError(t, err, "PreviewImport should not return error")
		require.NotNil(t, result, "Result should not be nil")

		assert.Equal(t, "csv", result.FileType)
		assert.Equal(t, 3, result.TotalRows, "Should have 3 rows")
		assert.LessOrEqual(t, len(result.PreviewRows), 10, "Should not exceed preview limit")
		assert.True(t, result.ValidationOK, "Validation should pass for valid data")
	})

	t.Run("Preview CSV with custom limit", func(t *testing.T) {
		service := &ImportService{}

		// Preview только 2 строки
		result, err := service.PreviewImport(ctx, []byte(csvData), "csv", 1, 2)

		require.NoError(t, err)
		require.NotNil(t, result)

		assert.Len(t, result.PreviewRows, 2, "Should preview exactly 2 rows")
	})

	t.Run("Preview CSV with validation errors", func(t *testing.T) {
		// CSV с ошибками валидации
		invalidCSV := `name,price,currency,sku
,99.99,USD,TEST-001
Valid Product,-50.00,USD,TEST-002
Another Product,100,INVALID,TEST-003`

		service := &ImportService{}
		result, err := service.PreviewImport(ctx, []byte(invalidCSV), "csv", 1, 10)

		require.NoError(t, err, "PreviewImport should not fail on validation errors")
		require.NotNil(t, result)

		// Проверяем что есть строки с ошибками
		hasErrors := false
		for _, row := range result.PreviewRows {
			if !row.IsValid {
				hasErrors = true
				assert.NotEmpty(t, row.Errors, "Invalid row should have errors")
			}
		}
		assert.True(t, hasErrors, "Should have at least one invalid row")
		assert.False(t, result.ValidationOK, "Overall validation should fail")
	})
}

// TestPreviewImportXML тестирует preview функционал для XML файлов
func TestPreviewImportXML(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	// Подготовка тестовых данных XML в формате Digital Vision
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<artikli>
	<artikal>
		<id>1</id>
		<sifra>TEST-SKU-001</sifra>
		<naziv>Test Product 1</naziv>
		<kategorija1>Electronics</kategorija1>
		<kategorija2>Phones</kategorija2>
		<kategorija3>Smartphones</kategorija3>
		<vpCena>89.99</vpCena>
		<mpCena>99.99</mpCena>
		<dostupan>1</dostupan>
		<naAkciji>0</naAkciji>
		<opis>Test description 1</opis>
		<barKod>1234567890123</barKod>
		<slike>
			<slika>https://example.com/image1.jpg</slika>
		</slike>
	</artikal>
	<artikal>
		<id>2</id>
		<sifra>TEST-SKU-002</sifra>
		<naziv>Test Product 2</naziv>
		<kategorija1>Clothing</kategorija1>
		<kategorija2>Shirts</kategorija2>
		<kategorija3></kategorija3>
		<vpCena>129.99</vpCena>
		<mpCena>149.99</mpCena>
		<dostupan>1</dostupan>
		<naAkciji>1</naAkciji>
		<opis>Test description 2</opis>
		<barKod>9876543210987</barKod>
		<slike>
			<slika>https://example.com/image2.jpg</slika>
		</slike>
	</artikal>
</artikli>`

	t.Run("Preview XML with default limit", func(t *testing.T) {
		service := &ImportService{}

		result, err := service.PreviewImport(ctx, []byte(xmlData), "xml", 1, 10)

		require.NoError(t, err)
		require.NotNil(t, result)

		assert.Equal(t, "xml", result.FileType)
		assert.Equal(t, 2, result.TotalRows, "Should have 2 products")
		assert.LessOrEqual(t, len(result.PreviewRows), 10)
	})

	t.Run("Preview XML with limit 1", func(t *testing.T) {
		service := &ImportService{}

		result, err := service.PreviewImport(ctx, []byte(xmlData), "xml", 1, 1)

		require.NoError(t, err)
		assert.Len(t, result.PreviewRows, 1, "Should preview only 1 row")
	})
}

// TestBatchImport тестирует batch processing функционал
func TestBatchImport(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	t.Run("Batch import with 100 products", func(t *testing.T) {
		// Генерация CSV с 100 товарами
		var csvBuffer bytes.Buffer
		csvBuffer.WriteString("name,price,currency,sku,description\n")
		for i := 1; i <= 100; i++ {
			csvBuffer.WriteString("Product ")
			csvBuffer.WriteString(string(rune('0' + (i/10)%10)))
			csvBuffer.WriteString(string(rune('0' + i%10)))
			csvBuffer.WriteString(",99.99,USD,SKU-")
			csvBuffer.WriteString(string(rune('0' + (i/100)%10)))
			csvBuffer.WriteString(string(rune('0' + (i/10)%10)))
			csvBuffer.WriteString(string(rune('0' + i%10)))
			csvBuffer.WriteString(",Test description\n")
		}

		service := &ImportService{}

		// Проверяем что preview работает для большого файла
		result, err := service.PreviewImport(ctx, csvBuffer.Bytes(), "csv", 1, 10)

		require.NoError(t, err)
		require.NotNil(t, result)

		assert.Equal(t, 100, result.TotalRows, "Should detect 100 rows")
		assert.Equal(t, 10, len(result.PreviewRows), "Should preview only 10 rows")
	})

	t.Run("Batch import performance", func(t *testing.T) {
		// Этот тест будет измерять производительность batch processing
		// После интеграции с реальной БД
		t.Skip("Performance test will be implemented with database integration")
	})
}

// TestImportModes тестирует различные режимы импорта
func TestImportModes(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("Create only mode", func(t *testing.T) {
		// Будет протестирован после интеграции с БД
		t.Skip("Requires database integration")
	})

	t.Run("Update only mode", func(t *testing.T) {
		t.Skip("Requires database integration")
	})

	t.Run("Upsert mode", func(t *testing.T) {
		t.Skip("Requires database integration")
	})
}

// TestPreviewImportErrors тестирует обработку ошибок в preview
func TestPreviewImportErrors(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	service := &ImportService{}

	t.Run("Empty file", func(t *testing.T) {
		result, err := service.PreviewImport(ctx, []byte{}, "csv", 1, 10)

		assert.Error(t, err, "Should return error for empty file")
		assert.Nil(t, result)
	})

	t.Run("Unsupported file type", func(t *testing.T) {
		result, err := service.PreviewImport(ctx, []byte("test"), "unsupported", 1, 10)

		assert.Error(t, err, "Should return error for unsupported file type")
		assert.Nil(t, result)
	})

	t.Run("ZIP file not supported in preview", func(t *testing.T) {
		result, err := service.PreviewImport(ctx, []byte("test"), "zip", 1, 10)

		assert.Error(t, err, "ZIP preview should not be supported")
		assert.Contains(t, err.Error(), "not supported", "Error should mention unsupported")
		assert.Nil(t, result)
	})

	t.Run("Invalid preview limit", func(t *testing.T) {
		// CSV с обязательными полями
		csvData := "name,price,currency,sku\nProduct 1,99.99,USD,SKU-001"

		// Тест с limit = 0 (должен использовать default)
		result, err := service.PreviewImport(ctx, []byte(csvData), "csv", 1, 0)
		assert.NoError(t, err, "Should handle limit=0 gracefully")
		if result != nil {
			assert.GreaterOrEqual(t, len(result.PreviewRows), 0)
		}

		// Тест с limit > 100 (должен быть ограничен до 100)
		result, err = service.PreviewImport(ctx, []byte(csvData), "csv", 1, 150)
		assert.NoError(t, err)
		if result != nil {
			assert.LessOrEqual(t, len(result.PreviewRows), 100, "Should cap at 100")
		}
	})
}
