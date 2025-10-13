// backend/internal/proj/c2c/service/marketplace_images.go
package service

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"backend/internal/domain/models"
	"backend/internal/logger"
)

// ProcessImage обрабатывает загружаемое изображение и генерирует уникальное имя
func (s *MarketplaceService) ProcessImage(file *multipart.FileHeader) (string, error) {
	// Получаем расширение файла
	ext := filepath.Ext(file.Filename)
	if ext == "" {
		// Если расширение отсутствует, определяем его по MIME-типу
		switch file.Header.Get("Content-Type") {
		case "image/jpeg", "image/jpg":
			ext = ".jpg"
		case "image/png":
			ext = ".png"
		case "image/gif":
			ext = ".gif"
		case "image/webp":
			ext = ".webp"
		default:
			ext = ".jpg" // По умолчанию используем .jpg
		}
	}

	// Генерируем уникальное имя файла с расширением
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	return fileName, nil
}

// UploadImage загружает изображение в хранилище и сохраняет информацию в БД
func (s *MarketplaceService) UploadImage(ctx context.Context, file *multipart.FileHeader, listingID int, isMain bool) (*models.MarketplaceImage, error) {
	// Get file name
	fileName, err := s.ProcessImage(file)
	if err != nil {
		return nil, err
	}

	// Create object path - ensure no duplicate 'listings/' prefix
	objectName := fmt.Sprintf("%d/%s", listingID, fileName)

	// Open the file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer func() {
		if err := src.Close(); err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}()

	// Use FileStorage to upload
	fileStorage := s.storage.FileStorage()
	if fileStorage == nil {
		return nil, fmt.Errorf("file storage service not initialized")
	}

	// Upload to storage
	publicURL, err := fileStorage.UploadFile(ctx, objectName, src, file.Size, file.Header.Get("Content-Type"))
	if err != nil {
		return nil, fmt.Errorf("error uploading file: %w", err)
	}
	log.Printf("UploadImage: Изображение загружено в MinIO. objectName=%s, publicURL=%s", objectName, publicURL)

	// Create image information
	image := &models.MarketplaceImage{
		ListingID:     listingID,
		FilePath:      objectName,
		FileName:      file.Filename,
		FileSize:      int(file.Size),
		ContentType:   file.Header.Get("Content-Type"),
		IsMain:        isMain,
		StorageType:   "minio", // Явно указываем тип хранилища!
		StorageBucket: "listings",
		PublicURL:     publicURL,
		ImageURL:      publicURL, // Заполняем ImageURL для API
	}
	log.Printf("UploadImage: Сохраняем информацию об изображении: ListingID=%d, FilePath=%s, StorageType=%s, PublicURL=%s",
		image.ListingID, image.FilePath, image.StorageType, image.PublicURL)
	// Save image information to database
	imageID, err := s.storage.AddListingImage(ctx, image)
	if err != nil {
		// Если не удалось сохранить информацию, удаляем файл
		if err := fileStorage.DeleteFile(ctx, objectName); err != nil {
			logger.Error().Err(err).Str("objectName", objectName).Msg("Failed to delete file from storage")
		}
		return nil, fmt.Errorf("error saving image information: %w", err)
	}
	log.Printf("UploadImage: Изображение успешно сохранено в базе данных с ID=%d", imageID)

	image.ID = imageID
	return image, nil
}

// DeleteImage удаляет изображение из хранилища и БД
func (s *MarketplaceService) DeleteImage(ctx context.Context, imageID int) error {
	// Получаем информацию об изображении
	image, err := s.storage.GetListingImageByID(ctx, imageID)
	if err != nil {
		return fmt.Errorf("ошибка получения информации об изображении: %w", err)
	}

	// Используем FileStorage для удаления файла
	fileStorage := s.storage.FileStorage()
	if fileStorage == nil {
		return fmt.Errorf("сервис файлового хранилища не инициализирован")
	}

	// Удаляем файл из хранилища
	err = fileStorage.DeleteFile(ctx, image.FilePath)
	if err != nil {
		log.Printf("Ошибка удаления файла из хранилища: %v", err)
		// Продолжаем выполнение для удаления записи из базы данных
	}

	// Удаляем информацию об изображении из базы данных
	err = s.storage.DeleteListingImage(ctx, imageID)
	if err != nil {
		return fmt.Errorf("ошибка удаления информации об изображении: %w", err)
	}

	return nil
}

// MigrateImagesToMinio мигрирует изображения из локального хранилища в MinIO
func (s *MarketplaceService) MigrateImagesToMinio(ctx context.Context) error {
	// Этот метод будем вызывать вручную при необходимости миграции

	// Получаем все изображения с типом хранилища 'local'
	query := `
		SELECT id, listing_id, file_path, file_name, file_size, content_type, is_main, created_at
		FROM c2c_images
		WHERE storage_type = 'local' OR storage_type IS NULL
	`

	rows, err := s.storage.Query(ctx, query)
	if err != nil {
		return fmt.Errorf("ошибка получения изображений: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			// Логирование ошибки закрытия rows
			_ = err // Explicitly ignore error
		}
	}()

	var count int
	for rows.Next() {
		var image models.MarketplaceImage
		err := rows.Scan(
			&image.ID, &image.ListingID, &image.FilePath, &image.FileName,
			&image.FileSize, &image.ContentType, &image.IsMain, &image.CreatedAt,
		)
		if err != nil {
			log.Printf("Ошибка сканирования данных изображения: %v", err)
			continue
		}

		// Пропускаем, если путь к файлу пустой
		if image.FilePath == "" {
			continue
		}

		// Исключаем уже мигрированные изображения
		if strings.HasPrefix(image.FilePath, "listings/") {
			continue
		}

		// Создаем новый путь для изображения в MinIO
		newPath := fmt.Sprintf("listings/%d/%s", image.ListingID, filepath.Base(image.FilePath))

		// Security check: validate file path
		if strings.Contains(image.FilePath, "..") {
			log.Printf("Skipping image with invalid path: %s", image.FilePath)
			continue
		}

		// Открываем исходный файл
		localPath := fmt.Sprintf("./uploads/%s", image.FilePath)
		file, err := os.Open(localPath) // #nosec G304 -- path validated above
		if err != nil {
			log.Printf("Ошибка открытия файла %s: %v", localPath, err)
			continue
		}

		// Получаем размер файла
		fileInfo, err := file.Stat()
		if err != nil {
			log.Printf("Ошибка получения информации о файле %s: %v", localPath, err)
			if closeErr := file.Close(); closeErr != nil {
				log.Printf("Ошибка закрытия файла %s: %v", localPath, closeErr)
			}
			continue
		}

		// Загружаем файл в MinIO
		fileStorage := s.storage.FileStorage()
		if fileStorage == nil {
			if closeErr := file.Close(); closeErr != nil {
				log.Printf("Ошибка закрытия файла %s: %v", localPath, closeErr)
			}
			return fmt.Errorf("сервис файлового хранилища не инициализирован")
		}

		publicURL, err := fileStorage.UploadFile(ctx, newPath, file, fileInfo.Size(), image.ContentType)
		if closeErr := file.Close(); closeErr != nil {
			log.Printf("Ошибка закрытия файла %s: %v", localPath, closeErr)
		}
		if err != nil {
			log.Printf("Ошибка загрузки файла %s в MinIO: %v", localPath, err)
			continue
		}

		// Обновляем информацию об изображении в базе данных
		_, err = s.storage.Exec(ctx, `
			UPDATE c2c_images
			SET file_path = $1, storage_type = 'minio', storage_bucket = 'listings', public_url = $2
			WHERE id = $3
		`, newPath, publicURL, image.ID)
		if err != nil {
			log.Printf("Ошибка обновления информации об изображении %d: %v", image.ID, err)
			continue
		}

		count++
		log.Printf("Успешно мигрировано изображение %d для объявления %d", image.ID, image.ListingID)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("ошибка итерации по изображениям: %w", err)
	}

	log.Printf("Миграция завершена. Всего мигрировано %d изображений", count)

	return nil
}

// AddListingImage добавляет изображение к объявлению
func (s *MarketplaceService) AddListingImage(ctx context.Context, image *models.MarketplaceImage) (int, error) {
	return s.storage.AddListingImage(ctx, image)
}
