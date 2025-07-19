package interfaces

import (
	"context"

	"backend/internal/domain/models"
)

// ImageRepositoryInterface - интерфейс для работы с изображениями
type ImageRepositoryInterface interface {
	// CreateImage создает новое изображение
	CreateImage(ctx context.Context, image models.ImageInterface) (models.ImageInterface, error)

	// GetImageByID получает изображение по ID
	GetImageByID(ctx context.Context, imageID int) (models.ImageInterface, error)

	// GetImagesByEntity получает все изображения для сущности
	GetImagesByEntity(ctx context.Context, entityType string, entityID int) ([]models.ImageInterface, error)

	// DeleteImage удаляет изображение
	DeleteImage(ctx context.Context, imageID int) error

	// UnsetMainImages сбрасывает флаг главного изображения для всех изображений сущности
	UnsetMainImages(ctx context.Context, entityType string, entityID int) error

	// SetMainImage устанавливает изображение как главное
	SetMainImage(ctx context.Context, imageID int, isMain bool) error

	// UpdateDisplayOrder обновляет порядок отображения изображений
	UpdateDisplayOrder(ctx context.Context, imageID int, displayOrder int) error
}
