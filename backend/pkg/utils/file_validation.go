package utils

import (
	"fmt"
	"mime"
	"path/filepath"
	"strings"
)

// AllowedFileTypes определяет разрешенные типы файлов
var AllowedFileTypes = map[string][]string{
	"image": {
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
	},
	"video": {
		"video/mp4",
		"video/mpeg",
		"video/quicktime",
		"video/x-msvideo",
		"video/webm",
	},
	"document": {
		"application/pdf",
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/vnd.ms-excel",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"text/plain",
	},
}

// MaxFileSizes определяет максимальные размеры файлов в байтах
var MaxFileSizes = map[string]int64{
	"image":    10 * 1024 * 1024,  // 10MB
	"video":    100 * 1024 * 1024, // 100MB
	"document": 20 * 1024 * 1024,  // 20MB
}

// ValidateFileType проверяет, является ли тип файла разрешенным
func ValidateFileType(contentType string) (string, error) {
	contentType = strings.ToLower(contentType)

	// Пытаемся определить основной тип из MIME
	baseMime, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", fmt.Errorf("invalid content type: %s", contentType)
	}

	// Проверяем каждую категорию
	for category, allowedTypes := range AllowedFileTypes {
		for _, allowed := range allowedTypes {
			if baseMime == allowed {
				return category, nil
			}
		}
	}

	return "", fmt.Errorf("file type not allowed: %s", contentType)
}

// ValidateFileSize проверяет размер файла
func ValidateFileSize(fileType string, size int64) error {
	maxSize, exists := MaxFileSizes[fileType]
	if !exists {
		return fmt.Errorf("unknown file type: %s", fileType)
	}

	if size > maxSize {
		return fmt.Errorf("file too large: %d bytes (max: %d bytes)", size, maxSize)
	}

	if size <= 0 {
		return fmt.Errorf("invalid file size: %d", size)
	}

	return nil
}

// ValidateFileName проверяет и санитизирует имя файла
func ValidateFileName(filename string) (string, error) {
	if filename == "" {
		return "", fmt.Errorf("filename cannot be empty")
	}

	// Удаляем потенциально опасные символы
	filename = filepath.Base(filename) // Убираем путь
	filename = strings.TrimSpace(filename)

	// Проверяем расширение
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		return "", fmt.Errorf("file must have an extension")
	}

	// Список опасных расширений
	dangerousExtensions := []string{
		".exe", ".bat", ".cmd", ".com", ".scr", ".vbs", ".js",
		".jar", ".app", ".deb", ".rpm", ".dmg", ".pkg", ".run",
	}

	for _, dangerous := range dangerousExtensions {
		if ext == dangerous {
			return "", fmt.Errorf("file extension not allowed: %s", ext)
		}
	}

	// Заменяем недопустимые символы
	replacer := strings.NewReplacer(
		"<", "_",
		">", "_",
		":", "_",
		"\"", "_",
		"/", "_",
		"\\", "_",
		"|", "_",
		"?", "_",
		"*", "_",
	)

	sanitized := replacer.Replace(filename)

	// Ограничиваем длину имени файла
	if len(sanitized) > 255 {
		// Сохраняем расширение
		extLen := len(ext)
		maxBaseName := 255 - extLen
		baseName := sanitized[:len(sanitized)-extLen]
		if len(baseName) > maxBaseName {
			baseName = baseName[:maxBaseName]
		}
		sanitized = baseName + ext
	}

	return sanitized, nil
}

// GetFileTypeFromMIME определяет категорию файла по MIME типу
func GetFileTypeFromMIME(contentType string) string {
	contentType = strings.ToLower(contentType)

	switch {
	case strings.HasPrefix(contentType, "image/"):
		return "image"
	case strings.HasPrefix(contentType, "video/"):
		return "video"
	case strings.HasPrefix(contentType, "application/pdf") ||
		strings.HasPrefix(contentType, "application/msword") ||
		strings.HasPrefix(contentType, "application/vnd.") ||
		strings.HasPrefix(contentType, "text/"):
		return "document"
	default:
		return "unknown"
	}
}
