package config

import (
	"backend/internal/config"
	"backend/internal/middleware"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// Handler представляет обработчик для конфигурации
type Handler struct {
	cfg *config.Config
}

// NewHandler создает новый экземпляр обработчика конфигурации
func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		cfg: cfg,
	}
}

// GetPrefix возвращает префикс для API конфигурации
func (h *Handler) GetPrefix() string {
	return "/api/v1"
}

// RegisterRoutes регистрирует маршруты для конфигурации
func (h *Handler) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
	// Регистрируем публичный endpoint для получения конфигурации storage
	// Важно: регистрируем напрямую в app, чтобы избежать middleware авторизации
	app.Get("/api/v1/config/storage", h.GetStorageConfig)

	return nil
}

// StorageConfig представляет конфигурацию хранилища для frontend
type StorageConfig struct {
	Provider           string `json:"provider"`
	BaseURL            string `json:"base_url"`
	ListingsBucket     string `json:"listings_bucket"`
	ChatFilesBucket    string `json:"chat_files_bucket"`
	StorefrontBucket   string `json:"storefront_bucket"`
	ReviewPhotosBucket string `json:"review_photos_bucket"`
}

// GetStorageConfig возвращает конфигурацию хранилища для frontend
// @Summary Get storage configuration
// @Description Returns storage configuration including bucket names for frontend
// @Tags config
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=StorageConfig} "Storage configuration"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/config/storage [get]
func (h *Handler) GetStorageConfig(c *fiber.Ctx) error {
	log.Info().Msg("Getting storage configuration")

	storageConfig := StorageConfig{
		Provider:           h.cfg.FileStorage.Provider,
		BaseURL:            h.cfg.FileStorage.PublicBaseURL,
		ListingsBucket:     h.cfg.FileStorage.MinioBucketName,
		ChatFilesBucket:    h.cfg.FileStorage.MinioChatBucket,
		StorefrontBucket:   h.cfg.FileStorage.MinioStorefrontBucket,
		ReviewPhotosBucket: h.cfg.FileStorage.MinioReviewPhotosBucket,
	}

	return utils.SuccessResponse(c, storageConfig)
}
