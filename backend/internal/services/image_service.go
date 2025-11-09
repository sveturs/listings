package services

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif" // Регистрация GIF декодера
	"image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	// Регистрация PNG декодера

	"backend/internal/domain/models"
	"backend/internal/storage/filestorage"
	"backend/internal/storage/interfaces"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
	_ "golang.org/x/image/webp" // Регистрация WebP декодера
)

// ImageType определяет тип изображения для выбора правильного бакета
type ImageType string

const (
	ImageTypeMarketplaceListing ImageType = "marketplace_listing"
	ImageTypeStorefrontProduct  ImageType = "storefront_product"
	ImageTypeStorefrontLogo     ImageType = "storefront_logo"
	ImageTypeStorefrontBanner   ImageType = "storefront_banner"
	ImageTypeChatFile           ImageType = "chat_file"
	ImageTypeReviewPhoto        ImageType = "review_photo"
)

// ImageService - единый сервис для работы с изображениями
type ImageService struct {
	fileStorage       filestorage.FileStorageInterface
	repo              interfaces.ImageRepositoryInterface
	bucketListings    string // Bucket для marketplace listings
	bucketStorefront  string // Bucket для storefront products
	bucketChatFiles   string // Bucket для chat files
	bucketReviewPhoto string // Bucket для review photos
}

// GetRepo возвращает репозиторий (для использования в handler)
func (s *ImageService) GetRepo() interfaces.ImageRepositoryInterface {
	return s.repo
}

// ImageServiceConfig содержит конфигурацию для ImageService
type ImageServiceConfig struct {
	BucketListings    string
	BucketStorefront  string
	BucketChatFiles   string
	BucketReviewPhoto string
}

// NewImageService создает новый ImageService
func NewImageService(fileStorage filestorage.FileStorageInterface, repo interfaces.ImageRepositoryInterface, cfg ImageServiceConfig) *ImageService {
	return &ImageService{
		fileStorage:       fileStorage,
		repo:              repo,
		bucketListings:    cfg.BucketListings,
		bucketStorefront:  cfg.BucketStorefront,
		bucketChatFiles:   cfg.BucketChatFiles,
		bucketReviewPhoto: cfg.BucketReviewPhoto,
	}
}

// UploadImageRequest - запрос для загрузки изображения
type UploadImageRequest struct {
	EntityType   ImageType             // Тип сущности (marketplace_listing, storefront_product, etc.)
	EntityID     int                   // ID сущности
	File         multipart.File        // Файл изображения
	FileHeader   *multipart.FileHeader // Заголовок файла
	IsMain       bool                  // Является ли главным изображением
	DisplayOrder int                   // Порядок отображения
}

// UploadImageResponse - ответ при загрузке изображения
type UploadImageResponse struct {
	ID           int    `json:"id"`
	ImageURL     string `json:"image_url"`
	ThumbnailURL string `json:"thumbnail_url"`
	PublicURL    string `json:"public_url"`
	IsMain       bool   `json:"is_main"`
	DisplayOrder int    `json:"display_order"`
}

// UploadImage загружает изображение в MinIO и сохраняет в БД
func (s *ImageService) UploadImage(ctx context.Context, req *UploadImageRequest) (*UploadImageResponse, error) {
	// Валидация файла
	if err := s.validateImageFile(req.FileHeader); err != nil {
		return nil, fmt.Errorf("invalid image file: %w", err)
	}

	// Определение пути для сохранения
	filePath := s.generateFilePath(req.EntityType, req.EntityID, req.FileHeader.Filename)

	// Чтение файла
	fileBytes, err := io.ReadAll(req.File)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Создание миниатюры
	thumbnailBytes, err := s.createThumbnail(fileBytes, req.FileHeader.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create thumbnail: %w", err)
	}

	// Определение бакета
	bucket := s.getBucketForImageType(req.EntityType)

	// Загрузка основного изображения в соответствующий bucket
	// Всегда используем UploadToCustomBucket, так как теперь у нас есть имена всех buckets
	imageURL, uploadErr := s.fileStorage.UploadToCustomBucket(ctx, bucket, filePath, bytes.NewReader(fileBytes), req.FileHeader.Size, req.FileHeader.Header.Get("Content-Type"))
	if uploadErr != nil {
		return nil, fmt.Errorf("failed to upload image: %w", uploadErr)
	}

	// Загрузка миниатюры в тот же bucket
	thumbnailPath := s.generateThumbnailPath(filePath)
	thumbnailURL, uploadErr := s.fileStorage.UploadToCustomBucket(ctx, bucket, thumbnailPath, bytes.NewReader(thumbnailBytes), int64(len(thumbnailBytes)), "image/jpeg")
	if uploadErr != nil {
		return nil, fmt.Errorf("failed to upload thumbnail: %w", uploadErr)
	}

	// Создание записи в БД
	imageRecord := s.createImageRecord(req, imageURL, thumbnailURL, bucket, filePath)

	// Сохранение в БД
	savedImage, err := s.repo.CreateImage(ctx, imageRecord)
	if err != nil {
		// Откат: удаление файлов из MinIO
		if deleteErr := s.fileStorage.DeleteFile(ctx, filePath); deleteErr != nil {
			// Логируем ошибку удаления, но не прерываем выполнение
			_ = deleteErr // Explicitly ignore error
		}
		if deleteErr := s.fileStorage.DeleteFile(ctx, thumbnailPath); deleteErr != nil {
			// Логируем ошибку удаления, но не прерываем выполнение
			_ = deleteErr // Explicitly ignore error
		}
		return nil, fmt.Errorf("failed to save image to database: %w", err)
	}

	return &UploadImageResponse{
		ID:           savedImage.GetID(),
		ImageURL:     imageURL,
		ThumbnailURL: thumbnailURL,
		PublicURL:    imageURL,
		IsMain:       req.IsMain,
		DisplayOrder: req.DisplayOrder,
	}, nil
}

// DeleteImage удаляет изображение из MinIO и БД
func (s *ImageService) DeleteImage(ctx context.Context, imageID int, entityType ImageType) error {
	// Получение информации об изображении
	imageRecord, err := s.repo.GetImageByID(ctx, imageID)
	if err != nil {
		return fmt.Errorf("failed to get image: %w", err)
	}

	// Получаем пути файлов из URL
	imageURL := imageRecord.GetImageURL()
	thumbnailURL := imageRecord.GetThumbnailURL()

	// Определяем bucket для типа изображения
	bucket := s.getBucketForImageType(entityType)

	// Удаляем префикс /{bucket}/ из URL для получения пути в MinIO
	bucketPrefix := "/" + bucket + "/"
	imagePath := strings.TrimPrefix(imageURL, bucketPrefix)
	// Также удаляем старый префикс для обратной совместимости
	if imagePath == imageURL {
		imagePath = strings.TrimPrefix(imageURL, "/listings/")
	}
	thumbnailPath := strings.TrimPrefix(thumbnailURL, bucketPrefix)
	if thumbnailPath == thumbnailURL {
		thumbnailPath = strings.TrimPrefix(thumbnailURL, "/listings/")
	}

	// Удаление основного изображения из MinIO
	if imagePath != "" && imagePath != imageURL {
		err = s.fileStorage.DeleteFileFromCustomBucket(ctx, bucket, imagePath)
		if err != nil {
			return fmt.Errorf("failed to delete image from storage: %w", err)
		}
	}

	// Удаление миниатюры из MinIO
	if thumbnailPath != "" && thumbnailPath != thumbnailURL {
		err = s.fileStorage.DeleteFileFromCustomBucket(ctx, bucket, thumbnailPath)
		if err != nil {
			// Логируем ошибку, но не останавливаем процесс
			fmt.Printf("Warning: failed to delete thumbnail: %v\n", err)
		}
	}

	// Удаление из БД
	if err := s.repo.DeleteImage(ctx, imageID); err != nil {
		return fmt.Errorf("failed to delete image from database: %w", err)
	}

	return nil
}

// GetImagesByEntity получает все изображения для сущности
func (s *ImageService) GetImagesByEntity(ctx context.Context, entityType ImageType, entityID int) ([]models.ImageInterface, error) {
	return s.repo.GetImagesByEntity(ctx, string(entityType), entityID)
}

// SetMainImage устанавливает изображение как главное
func (s *ImageService) SetMainImage(ctx context.Context, imageID int, entityType ImageType, entityID int) error {
	// Сброс всех изображений как не главных
	if err := s.repo.UnsetMainImages(ctx, string(entityType), entityID); err != nil {
		return fmt.Errorf("failed to unset main images: %w", err)
	}

	// Установка изображения как главного
	if err := s.repo.SetMainImage(ctx, imageID, true); err != nil {
		return fmt.Errorf("failed to set main image: %w", err)
	}

	return nil
}

// validateImageFile валидирует загружаемый файл
func (s *ImageService) validateImageFile(fileHeader *multipart.FileHeader) error {
	// Проверка размера файла (максимум 10MB)
	const maxFileSize = 10 * 1024 * 1024
	if fileHeader.Size > maxFileSize {
		return fmt.Errorf("file size exceeds maximum limit of 10MB")
	}

	// Проверка расширения файла
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}

	isAllowed := false
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		return fmt.Errorf("unsupported file format: %s", ext)
	}

	return nil
}

// generateFilePath генерирует путь для сохранения файла
// Для кастомных buckets путь не должен включать название bucket (оно уже в UploadToCustomBucket)
func (s *ImageService) generateFilePath(entityType ImageType, entityID int, filename string) string {
	ext := filepath.Ext(filename)
	uniqueID := uuid.New().String()
	timestamp := time.Now().Format("20060102150405")

	switch entityType {
	case ImageTypeMarketplaceListing:
		return fmt.Sprintf("listings/%d/%s_%s%s", entityID, timestamp, uniqueID, ext)
	case ImageTypeStorefrontProduct:
		// Для товаров витрин путь без префикса bucket (bucket передается отдельно)
		return fmt.Sprintf("%d/%s_%s%s", entityID, timestamp, uniqueID, ext)
	case ImageTypeStorefrontLogo:
		return fmt.Sprintf("%d/logo_%s_%s%s", entityID, timestamp, uniqueID, ext)
	case ImageTypeStorefrontBanner:
		return fmt.Sprintf("%d/banner_%s_%s%s", entityID, timestamp, uniqueID, ext)
	case ImageTypeChatFile:
		return fmt.Sprintf("%d/%s_%s%s", entityID, timestamp, uniqueID, ext)
	case ImageTypeReviewPhoto:
		return fmt.Sprintf("%d/%s_%s%s", entityID, timestamp, uniqueID, ext)
	default:
		return fmt.Sprintf("misc/%s_%s%s", timestamp, uniqueID, ext)
	}
}

// generateThumbnailPath генерирует путь для миниатюры
func (s *ImageService) generateThumbnailPath(originalPath string) string {
	dir := filepath.Dir(originalPath)
	filename := filepath.Base(originalPath)
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)

	return fmt.Sprintf("%s/thumb_%s.jpg", dir, nameWithoutExt)
}

// getBucketForImageType возвращает бакет для типа изображения
func (s *ImageService) getBucketForImageType(entityType ImageType) string {
	switch entityType {
	case ImageTypeMarketplaceListing:
		return s.bucketListings
	case ImageTypeStorefrontProduct:
		return s.bucketStorefront
	case ImageTypeStorefrontLogo, ImageTypeStorefrontBanner:
		return s.bucketStorefront // Используем тот же bucket для storefront images
	case ImageTypeChatFile:
		return s.bucketChatFiles
	case ImageTypeReviewPhoto:
		return s.bucketReviewPhoto
	default:
		return s.bucketListings // Fallback на listings bucket
	}
}

// createThumbnail создает миниатюру изображения
func (s *ImageService) createThumbnail(imageBytes []byte, filename string) ([]byte, error) {
	// Декодирование изображения
	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Создание миниатюры (200x200 пикселей)
	// Используем Bilinear вместо Lanczos3 для ~3x speed improvement
	// Незначительная потеря качества acceptable для thumbnails
	thumbnail := resize.Thumbnail(200, 200, img, resize.Bilinear)

	// Кодирование в JPEG
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, thumbnail, &jpeg.Options{Quality: 85}); err != nil {
		return nil, fmt.Errorf("failed to encode thumbnail: %w", err)
	}

	return buf.Bytes(), nil
}

// createImageRecord создает запись изображения для БД
func (s *ImageService) createImageRecord(req *UploadImageRequest, imageURL, thumbnailURL, bucket, filePath string) models.ImageInterface {
	switch req.EntityType {
	case ImageTypeMarketplaceListing:
		return &models.MarketplaceImage{
			ListingID:     req.EntityID,
			ImageURL:      imageURL,
			ThumbnailURL:  thumbnailURL,
			FilePath:      filePath,
			StorageType:   "minio",
			StorageBucket: bucket,
			PublicURL:     imageURL,
			IsMain:        req.IsMain,
			DisplayOrder:  req.DisplayOrder,
			FileName:      req.FileHeader.Filename,
			ContentType:   req.FileHeader.Header.Get("Content-Type"),
		}
	case ImageTypeStorefrontProduct:
		return &models.StorefrontProductImage{
			StorefrontProductID: req.EntityID,
			ImageURL:            imageURL,
			ThumbnailURL:        thumbnailURL,
			DisplayOrder:        req.DisplayOrder,
			IsDefault:           req.IsMain,
		}
	case ImageTypeStorefrontLogo:
		// Для логотипа витрины
		return &models.StorefrontProductImage{
			StorefrontProductID: req.EntityID,
			ImageURL:            imageURL,
			ThumbnailURL:        thumbnailURL,
			DisplayOrder:        req.DisplayOrder,
			IsDefault:           req.IsMain,
		}
	case ImageTypeStorefrontBanner:
		// Для баннера витрины
		return &models.StorefrontProductImage{
			StorefrontProductID: req.EntityID,
			ImageURL:            imageURL,
			ThumbnailURL:        thumbnailURL,
			DisplayOrder:        req.DisplayOrder,
			IsDefault:           req.IsMain,
		}
	case ImageTypeChatFile:
		// Для файлов чата
		return &models.StorefrontProductImage{
			StorefrontProductID: req.EntityID,
			ImageURL:            imageURL,
			ThumbnailURL:        thumbnailURL,
			DisplayOrder:        req.DisplayOrder,
			IsDefault:           req.IsMain,
			FilePath:            filePath,
			FileName:            req.FileHeader.Filename,
			FileSize:            int(req.FileHeader.Size),
			ContentType:         req.FileHeader.Header.Get("Content-Type"),
			StorageType:         "minio",
			StorageBucket:       bucket,
			PublicURL:           imageURL,
		}
	case ImageTypeReviewPhoto:
		// Для фото отзывов
		return &models.StorefrontProductImage{
			StorefrontProductID: req.EntityID,
			ImageURL:            imageURL,
			ThumbnailURL:        thumbnailURL,
			DisplayOrder:        req.DisplayOrder,
			IsDefault:           req.IsMain,
			FilePath:            filePath,
			FileName:            req.FileHeader.Filename,
			FileSize:            int(req.FileHeader.Size),
			ContentType:         req.FileHeader.Header.Get("Content-Type"),
			StorageType:         "minio",
			StorageBucket:       bucket,
			PublicURL:           imageURL,
		}
	default:
		// Для других типов можно добавить обработку
		return &models.StorefrontProductImage{
			StorefrontProductID: req.EntityID,
			ImageURL:            imageURL,
			ThumbnailURL:        thumbnailURL,
			DisplayOrder:        req.DisplayOrder,
			IsDefault:           req.IsMain,
		}
	}
}
