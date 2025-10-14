// backend/internal/proj/c2c/storage/postgres/storage_utils.go
package postgres

import (
	"fmt"
	"os"
	"strings"

	"backend/internal/domain/models"
)

// processTranslations обрабатывает сырые переводы из БД в структурированный TranslationMap
func (s *Storage) processTranslations(rawTranslations interface{}) models.TranslationMap {
	translations := make(models.TranslationMap)

	if rawMap, ok := rawTranslations.(map[string]interface{}); ok {
		for key, value := range rawMap {
			parts := strings.Split(key, "_")
			if len(parts) != 2 {
				continue
			}

			lang, field := parts[0], parts[1]
			if translations[lang] == nil {
				translations[lang] = make(map[string]string)
			}

			if strValue, ok := value.(string); ok {
				translations[lang][field] = strValue
			}
		}
	}

	return translations
}

// buildFullImageURL преобразует относительный URL в полный URL с базовым адресом
func buildFullImageURL(relativeURL string) string {
	if relativeURL == "" {
		return ""
	}

	// Если URL уже полный (начинается с http:// или https://), возвращаем как есть
	if strings.HasPrefix(relativeURL, "http://") || strings.HasPrefix(relativeURL, "https://") {
		return relativeURL
	}

	// Получаем базовый URL из переменной окружения
	baseURL := os.Getenv("FILE_STORAGE_PUBLIC_URL")
	if baseURL == "" {
		// Если не задан, используем относительный URL как есть
		return relativeURL
	}

	// Убираем trailing slash из baseURL если есть
	baseURL = strings.TrimRight(baseURL, "/")

	// Обрабатываем URL в зависимости от формата
	switch {
	case strings.HasPrefix(relativeURL, "/listings/"):
		// URL вида /listings/294/... -> нужно заменить на правильный bucket
		bucketName := os.Getenv("MINIO_BUCKET_NAME")
		if bucketName == "" {
			bucketName = "development-listings"
		}
		// Убираем /listings/ и добавляем правильный bucket
		path := strings.TrimPrefix(relativeURL, "/listings/")
		return fmt.Sprintf("%s/%s/%s", baseURL, bucketName, path)
	case strings.HasPrefix(relativeURL, "/"):
		// URL начинается с /, убираем его чтобы не было двойного слэша
		return baseURL + relativeURL
	default:
		// URL без начального слэша
		bucketName := os.Getenv("MINIO_BUCKET_NAME")
		if bucketName == "" {
			bucketName = "development-listings"
		}
		return fmt.Sprintf("%s/%s/%s", baseURL, bucketName, relativeURL)
	}
}
