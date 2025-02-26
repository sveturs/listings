package service

import (
	"backend/internal/domain/models"
	"context"
	"io"
)

type StorefrontServiceInterface interface {
	// Управление витринами
	CreateStorefront(ctx context.Context, userID int, create *models.StorefrontCreate) (*models.Storefront, error)
	GetUserStorefronts(ctx context.Context, userID int) ([]models.Storefront, error)
	GetStorefrontByID(ctx context.Context, id int, userID int) (*models.Storefront, error)
	UpdateStorefront(ctx context.Context, storefront *models.Storefront, userID int) error
	DeleteStorefront(ctx context.Context, id int, userID int) error
	
	// Источники импорта
	CreateImportSource(ctx context.Context, source *models.ImportSourceCreate, userID int) (*models.ImportSource, error)
	UpdateImportSource(ctx context.Context, source *models.ImportSource, userID int) error
	DeleteImportSource(ctx context.Context, id int, userID int) error
	GetImportSources(ctx context.Context, storefrontID int, userID int) ([]models.ImportSource, error)
	GetImportSourceByID(ctx context.Context, id int, userID int) (*models.ImportSource, error)
	
	// Импорт данных
	RunImport(ctx context.Context, sourceID int, userID int) (*models.ImportHistory, error)
	ImportCSV(ctx context.Context, sourceID int, reader io.Reader, userID int) (*models.ImportHistory, error)
	GetImportHistory(ctx context.Context, sourceID int, userID int, limit, offset int) ([]models.ImportHistory, error)
}