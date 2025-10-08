// backend/internal/proj/c2c/handler/images.go
package handler

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"

	"backend/internal/domain/models"
	"backend/internal/logger"
	"backend/internal/proj/c2c/service"
	globalService "backend/internal/proj/global/service"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// ImagesHandler handles requests related to listing images
type ImagesHandler struct {
	services           globalService.ServicesInterface
	marketplaceService service.MarketplaceServiceInterface
}

// NewImagesHandler creates a new images handler
func NewImagesHandler(services globalService.ServicesInterface) *ImagesHandler {
	return &ImagesHandler{
		services:           services,
		marketplaceService: services.Marketplace(),
	}
}

// UploadImages uploads images for a listing
// @Summary Upload listing images
// @Description Uploads multiple images for a marketplace listing
// @Tags marketplace-images
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Listing ID"
// @Param file formData file true "Image files to upload"
// @Param main_image_index formData int false "Index of the main image"
// @Success 200 {object} utils.SuccessResponseSwag{data=ImagesUploadResponse} "Images uploaded successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidData or marketplace.invalidId"
// @Failure 401 {object} utils.ErrorResponseSwag "marketplace.authRequired"
// @Failure 403 {object} utils.ErrorResponseSwag "marketplace.forbidden"
// @Failure 404 {object} utils.ErrorResponseSwag "marketplace.notFound"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.uploadError"
// @Security BearerAuth
// @Router /api/v1/marketplace/listings/{id}/images [post]
func (h *ImagesHandler) UploadImages(c *fiber.Ctx) error {
	// Добавим явные логи для отладки
	logger.Info().Str("method", c.Method()).Str("path", c.Path()).Msg("Starting image upload")

	// Получаем ID пользователя из контекста
	userID, ok := authMiddleware.GetUserID(c)
	if !ok {
		logger.Warn().Msg("User ID not found in context")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "marketplace.authRequired")
	}
	logger.Info().Int("userId", userID).Msg("Authenticated user")

	// Проверяем, пришли ли какие-то файлы
	form, err := c.MultipartForm()
	if err != nil {
		logger.Error().Err(err).Msg("Error getting MultipartForm")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.fileError")
	}

	files := form.File["file"]
	if len(files) == 0 {
		// Попробуем альтернативное имя поля
		files = form.File["files"]
		if len(files) == 0 {
			// Проверим все поля
			logger.Info().Interface("keys", getMapKeys(form.File)).Msg("Searching files in form.File")
			for key, values := range form.File {
				logger.Info().Str("field", key).Int("filesCount", len(values)).Msg("Field contains files")
				if len(values) > 0 {
					files = values
					break
				}
			}

			if len(files) == 0 {
				logger.Warn().Msg("No files found in request")
				return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.noFiles")
			}
		}
	}

	logger.Info().Int("filesCount", len(files)).Msg("Files found for upload")

	// Получаем ID объявления из параметров
	listingIDStr := c.FormValue("listing_id")
	if listingIDStr == "" {
		// Попробуем из параметров URL
		listingIDStr = c.Params("id")
		if listingIDStr == "" {
			logger.Error().Msg("Error: listing ID not specified")
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.missingListingId")
		}
	}

	listingID, err := strconv.Atoi(listingIDStr)
	if err != nil {
		logger.Error().Err(err).Str("listingIdStr", listingIDStr).Msg("Error converting listing ID")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidId")
	}
	logger.Info().Int("listingId", listingID).Msg("Listing ID for image upload")

	// Получаем информацию об объявлении для проверки владельца
	listing, err := h.marketplaceService.GetListingByID(c.Context(), listingID)
	if err != nil {
		logger.Error().Err(err).Int("listingId", listingID).Msg("Error getting listing")
		return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.notFound")
	}

	// Проверяем, владеет ли пользователь объявлением
	if listing.UserID != userID {
		logger.Warn().Int("userId", userID).Int("listingId", listingID).Int("ownerId", listing.UserID).Msg("Access denied: user does not own listing")

		return utils.ErrorResponse(c, fiber.StatusForbidden, "marketplace.forbidden")
	}

	// Определяем главное изображение
	mainImageIndex := 0
	if mainImageIndexStr := c.FormValue("main_image_index"); mainImageIndexStr != "" {
		if idx, err := strconv.Atoi(mainImageIndexStr); err == nil {
			mainImageIndex = idx
			logger.Info().Int("mainImageIndex", mainImageIndex).Msg("Main image index set")
		}
	}

	var uploadedImages []models.MarketplaceImage
	for i, file := range files {
		contentType := file.Header.Get("Content-Type")
		logger.Info().Int("fileIndex", i).Str("filename", file.Filename).Int64("size", file.Size).Str("contentType", contentType).Msg("Processing file")

		// Проверка типа файла и размера
		if file.Size > 10*1024*1024 {
			logger.Warn().Int64("size", file.Size).Msg("File too large")
			continue // Пропускаем слишком большие файлы
		}

		// Уже определена выше
		if !strings.HasPrefix(contentType, "image/") {
			logger.Warn().Str("contentType", contentType).Msg("Unsupported file type")
			continue // Пропускаем файлы не-изображения
		}

		// Загружаем изображение
		isMain := i == mainImageIndex
		logger.Info().Int("listingId", listingID).Bool("isMain", isMain).Msg("Uploading image")
		image, err := h.marketplaceService.UploadImage(c.Context(), file, listingID, isMain)
		if err != nil {
			logger.Error().Err(err).Msg("Error uploading image")
			continue
		}

		logger.Info().Int("imageId", image.ID).Str("filePath", image.FilePath).Msg("Image successfully uploaded")

		uploadedImages = append(uploadedImages, *image)
	}

	// Переиндексируем объявление с загруженными изображениями
	logger.Info().Int("listingId", listingID).Msg("Reindexing listing with new images")
	fullListing, err := h.marketplaceService.GetListingByID(c.Context(), listingID)
	if err != nil {
		logger.Error().Err(err).Msg("Error getting full listing info for reindexing")
	} else {
		if err := h.marketplaceService.Storage().IndexListing(c.Context(), fullListing); err != nil {
			logger.Error().Err(err).Msg("Error reindexing listing after image upload")
		} else {
			logger.Info().Int("listingId", listingID).Int("imagesCount", len(fullListing.Images)).Msg("Successfully reindexed listing with images")
		}
	}

	logger.Info().Int("uploadedCount", len(uploadedImages)).Int("totalFiles", len(files)).Msg("Image upload completed")

	// Возвращаем успешный результат с информацией о загруженных изображениях
	response := ImagesUploadResponse{
		Success: true,
		Message: "marketplace.imagesUploaded",
		Images:  uploadedImages,
		Count:   len(uploadedImages),
	}
	return utils.SuccessResponse(c, response)
}

// Вспомогательная функция для получения ключей из map
func getMapKeys(m map[string][]*multipart.FileHeader) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// DeleteImage deletes an image
// @Summary Delete listing image
// @Description Deletes an image from a marketplace listing
// @Tags marketplace-images
// @Accept json
// @Produce json
// @Param id path int true "Image ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=MessageResponse} "Image deleted successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidId"
// @Failure 401 {object} utils.ErrorResponseSwag "marketplace.authRequired"
// @Failure 403 {object} utils.ErrorResponseSwag "marketplace.forbidden"
// @Failure 404 {object} utils.ErrorResponseSwag "marketplace.notFound"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.deleteError"
// @Security BearerAuth
// @Router /api/v1/marketplace/images/{id} [delete]
func (h *ImagesHandler) DeleteImage(c *fiber.Ctx) error {
	// Получаем ID пользователя из контекста
	userID, ok := authMiddleware.GetUserID(c)
	if !ok {
		logger.Warn().Msg("User ID not found in context")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "marketplace.authRequired")
	}

	// Получаем ID изображения из параметров URL
	imageID, err := strconv.Atoi(c.Params("image_id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidId")
	}

	// Получаем информацию об изображении
	image, err := h.services.Storage().GetListingImageByID(c.Context(), imageID)
	if err != nil {
		logger.Error().Err(err).Int("imageId", imageID).Msg("Failed to get image")
		return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.imageNotFound")
	}

	// Получаем информацию об объявлении
	listing, err := h.marketplaceService.GetListingByID(c.Context(), image.ListingID)
	if err != nil {
		logger.Error().Err(err).Int("listingId", image.ListingID).Msg("Failed to get listing")
		return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.notFound")
	}

	// Проверяем владельца объявления
	if listing.UserID != userID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "marketplace.forbidden")
	}

	// Удаляем изображение
	err = h.marketplaceService.DeleteImage(c.Context(), imageID)
	if err != nil {
		logger.Error().Err(err).Int("imageId", imageID).Msg("Failed to delete image")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.deleteError")
	}

	// Возвращаем успешный результат
	response := MessageResponse{
		Message: "marketplace.imageDeleted",
	}
	return utils.SuccessResponse(c, response)
}

// ModerateImage checks image for prohibited content
// @Summary Moderate image content
// @Description Checks an uploaded image for prohibited content using AI moderation
// @Tags marketplace-images
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Image file to moderate"
// @Success 200 {object} utils.SuccessResponseSwag{data=ImageModerationResponse} "Moderation results"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidFile or marketplace.fileTooLarge"
// @Failure 401 {object} utils.ErrorResponseSwag "marketplace.authRequired"
// @Failure 403 {object} utils.ErrorResponseSwag "marketplace.adminRequired"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.moderationError"
// @Security BearerAuth
// @Router /api/v1/marketplace/moderate-image [post]
func (h *ImagesHandler) ModerateImage(c *fiber.Ctx) error {
	// Проверяем, является ли пользователь администратором
	userID, ok := authMiddleware.GetUserID(c)
	if !ok {
		logger.Warn().Msg("User ID not found in context")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "marketplace.authRequired")
	}

	// Получаем пользователя по ID для проверки email
	user, err := h.services.User().GetUserByID(c.Context(), userID)
	if err != nil {
		logger.Error().Err(err).Int("userId", userID).Msg("Failed to get user")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.adminCheckError")
	}

	// Проверяем права администратора
	isAdmin, err := h.services.User().IsUserAdmin(c.Context(), user.Email)
	if err != nil || !isAdmin {
		logger.Error().Err(err).Int("userId", userID).Msg("User is not admin")
		return utils.ErrorResponse(c, fiber.StatusForbidden, "marketplace.adminRequired")
	}

	// Получаем файл из запроса
	file, err := c.FormFile("file")
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get file from request")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.fileError")
	}

	// Проверяем размер файла (ограничение 10MB)
	if file.Size > 10*1024*1024 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.fileTooLarge")
	}

	// Проверяем тип файла
	contentType := file.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidFileType")
	}

	// Создаем временный файл для сохранения загруженного изображения
	tempDir := os.TempDir()
	tempFilePath := filepath.Join(tempDir, fmt.Sprintf("moderate_image_%d", time.Now().UnixNano()))
	if err := c.SaveFile(file, tempFilePath); err != nil {
		logger.Error().Err(err).Msg("Failed to save file")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.saveFileError")
	}
	defer func() {
		if err := os.Remove(tempFilePath); err != nil {
			logger.Error().Err(err).Msg("Failed to remove temporary file")
		}
	}() // Удаляем временный файл после завершения

	// Возвращаем заглушку для модерации изображения (полная реализация требует Vision API)
	// TODO: Реализовать проверку с использованием Vision API

	// Возвращаем результаты проверки
	response := ImageModerationResponse{
		Success: true,
		Data: ModerationData{
			Labels:           []string{"image", "photo"},
			ProhibitedLabels: []string{},
			HasProhibited:    false,
		},
	}
	return utils.SuccessResponse(c, response)
}

// EnhancePreview creates preview of enhanced image
// @Summary Create image enhancement preview
// @Description Creates a preview of an enhanced image before applying changes
// @Tags marketplace-images
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param image_id formData int true "Image ID to enhance"
// @Param enhancement_type formData string false "Enhancement type" default(quality)
// @Success 200 {object} utils.SuccessResponseSwag{data=EnhancePreviewResponse} "Enhancement preview"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidId"
// @Failure 401 {object} utils.ErrorResponseSwag "marketplace.authRequired"
// @Failure 403 {object} utils.ErrorResponseSwag "marketplace.forbidden"
// @Failure 404 {object} utils.ErrorResponseSwag "marketplace.notFound"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.enhanceError"
// @Security BearerAuth
// @Router /api/v1/marketplace/enhance-preview [post]
func (h *ImagesHandler) EnhancePreview(c *fiber.Ctx) error {
	// Получаем ID пользователя из контекста
	userID, ok := authMiddleware.GetUserID(c)
	if !ok {
		logger.Warn().Msg("User ID not found in context")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "marketplace.authRequired")
	}

	// Получаем ID изображения из параметров
	imageID, err := strconv.Atoi(c.FormValue("image_id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidId")
	}

	// Получаем информацию об изображении
	image, err := h.services.Storage().GetListingImageByID(c.Context(), imageID)
	if err != nil {
		logger.Error().Err(err).Int("imageId", imageID).Msg("Failed to get image")
		return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.imageNotFound")
	}

	// Получаем информацию об объявлении
	listing, err := h.marketplaceService.GetListingByID(c.Context(), image.ListingID)
	if err != nil {
		logger.Error().Err(err).Int("listingId", image.ListingID).Msg("Failed to get listing")
		return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.notFound")
	}

	// Проверяем владельца объявления
	if listing.UserID != userID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "marketplace.forbidden")
	}

	// Получаем тип улучшения
	enhancementType := c.FormValue("enhancement_type")
	if enhancementType == "" {
		enhancementType = "quality" // По умолчанию используется качество
	}

	// Создаем предпросмотр улучшенного изображения
	// Здесь должен быть код для обработки изображения и создания предпросмотра
	// Пока возвращаем заглушку
	response := EnhancePreviewResponse{
		Success: true,
		Data: EnhancePreviewData{
			PreviewURL: fmt.Sprintf("https://example.com/preview/%d/%s", imageID, enhancementType),
		},
	}
	return utils.SuccessResponse(c, response)
}

// EnhanceImages enhances images for a listing
// @Summary Enhance listing images
// @Description Enhances multiple images for a marketplace listing using AI processing
// @Tags marketplace-images
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param listing_id formData int true "Listing ID"
// @Param enhancement_type formData string false "Enhancement type" default(quality)
// @Success 200 {object} utils.SuccessResponseSwag{data=EnhanceImagesResponse} "Enhancement job started"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidId or marketplace.invalidRequest"
// @Failure 401 {object} utils.ErrorResponseSwag "marketplace.authRequired"
// @Failure 403 {object} utils.ErrorResponseSwag "marketplace.forbidden"
// @Failure 404 {object} utils.ErrorResponseSwag "marketplace.notFound or marketplace.imageNotFound"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.enhanceError"
// @Security BearerAuth
// @Router /api/v1/marketplace/enhance-images [post]
func (h *ImagesHandler) EnhanceImages(c *fiber.Ctx) error {
	// Получаем ID пользователя из контекста
	userID, ok := authMiddleware.GetUserID(c)
	if !ok {
		logger.Warn().Msg("User ID not found in context")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "marketplace.authRequired")
	}

	// Получаем ID объявления из параметров
	listingID, err := strconv.Atoi(c.FormValue("listing_id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidId")
	}

	// Получаем информацию об объявлении
	listing, err := h.marketplaceService.GetListingByID(c.Context(), listingID)
	if err != nil {
		logger.Error().Err(err).Int("listingId", listingID).Msg("Failed to get listing")
		return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.notFound")
	}

	// Проверяем владельца объявления
	if listing.UserID != userID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "marketplace.forbidden")
	}

	// Получаем список ID изображений для улучшения
	var imageIDs []int
	if err := c.BodyParser(&imageIDs); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidRequest")
	}

	// Проверяем, что указанные изображения действительно принадлежат данному объявлению
	for _, imageID := range imageIDs {
		image, err := h.services.Storage().GetListingImageByID(c.Context(), imageID)
		if err != nil {
			logger.Error().Err(err).Int("imageId", imageID).Msg("Failed to get image")
			return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.imageNotFound")
		}

		if image.ListingID != listingID {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "marketplace.forbidden")
		}
	}

	// Получаем тип улучшения
	enhancementType := c.FormValue("enhancement_type")
	if enhancementType == "" {
		enhancementType = "quality" // По умолчанию используется качество
	}

	// Запускаем процесс улучшения изображений
	// Здесь должен быть код для обработки изображений
	// Пока возвращаем заглушку
	response := EnhanceImagesResponse{
		Success: true,
		Data: EnhanceImagesData{
			Message: "marketplace.imageEnhancementStarted",
			JobID:   fmt.Sprintf("enhance_%s_%d_%d", enhancementType, listingID, time.Now().Unix()),
		},
	}
	return utils.SuccessResponse(c, response)
}
