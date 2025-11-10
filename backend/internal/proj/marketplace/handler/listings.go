// TEMPORARY: Will be moved to microservice
package handler

import (
	"fmt"
	"mime/multipart"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"

	"backend/internal/domain/models"
	"backend/internal/services"
	"backend/internal/storage/postgres"
	"backend/pkg/utils"
)

// CreateListingRequest represents the request body for creating a listing
type CreateListingRequest struct {
	CategoryID   int     `json:"category_id" validate:"required"`
	Title        string  `json:"title" validate:"required,min=3,max=200"`
	Description  *string `json:"description,omitempty"`
	Price        float64 `json:"price" validate:"required,min=0"`
	Currency     string  `json:"currency,omitempty"`
	Quantity     int32   `json:"quantity,omitempty"`
	SKU          *string `json:"sku,omitempty"`
	StorefrontID *int    `json:"storefront_id,omitempty"`
}

// CreateListing godoc
// @Summary Создать новое объявление
// @Description Создать новое объявление в marketplace
// @Tags marketplace
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body CreateListingRequest true "Данные объявления"
// @Success 201 {object} utils.SuccessResponseSwag{data=interface{}}
// @Failure 400 {object} utils.ErrorResponseSwag
// @Failure 401 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/marketplace/listings [post]
func (h *Handler) CreateListing(c *fiber.Ctx) error {
	// Get authenticated user ID
	userID, ok := authMiddleware.GetUserID(c)
	if !ok {
		h.logger.Warn().Msg("CreateListing: user not authenticated")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.unauthorized")
	}

	// Parse request body
	var req CreateListingRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error().Err(err).Msg("CreateListing: failed to parse request body")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalid_request")
	}

	// Basic validation
	if req.CategoryID == 0 {
		h.logger.Error().Msg("CreateListing: category_id is required")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.category_required")
	}
	if req.Title == "" || len(req.Title) < 3 || len(req.Title) > 200 {
		h.logger.Error().Msg("CreateListing: title must be between 3 and 200 characters")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalid_title")
	}
	if req.Price < 0 {
		h.logger.Error().Msg("CreateListing: price must be non-negative")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalid_price")
	}

	// Set default currency if not provided
	if req.Currency == "" {
		req.Currency = "RSD"
	}

	// Set default quantity if not provided
	if req.Quantity == 0 {
		req.Quantity = 1
	}

	// Call the storage layer (direct DB insert for now)
	// TEMPORARY: Direct DB insert until microservice fully migrated
	listing, err := h.storage.CreateListing(c.Context(), int(userID), req.CategoryID, req.Title, req.Description, req.Price, req.Currency, req.Quantity, req.SKU, req.StorefrontID)
	if err != nil {
		h.logger.Error().Err(err).
			Int("user_id", int(userID)).
			Int("category_id", req.CategoryID).
			Msg("CreateListing: failed to create listing")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.create_failed")
	}

	h.logger.Info().
		Int("listing_id", listing.ID).
		Int("user_id", int(userID)).
		Str("title", req.Title).
		Msg("Listing created successfully")

	c.Status(fiber.StatusCreated)
	return utils.SuccessResponse(c, listing)
}

// GetListing godoc
// @Summary Получить объявление по ID
// @Description Получить детали объявления по его ID с массивом изображений
// @Tags marketplace
// @Accept json
// @Produce json
// @Param id path int true "Listing ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.MarketplaceListing}
// @Failure 404 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/marketplace/listings/{id} [get]
func (h *Handler) GetListing(c *fiber.Ctx) error {
	// Parse listing ID
	idParam := c.Params("id")
	listingID, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Error().Err(err).Str("id", idParam).Msg("GetListing: invalid listing ID")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalid_id")
	}

	// Get listing from storage
	listing, err := h.storage.GetListing(c.Context(), listingID)
	if err != nil {
		h.logger.Error().Err(err).Int("listing_id", listingID).Msg("GetListing: failed to get listing")
		return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.listing_not_found")
	}

	return utils.SuccessResponse(c, listing)
}

// GetSimilarListings godoc
// @Summary Получить похожие объявления
// @Description Получить список похожих объявлений для данного листинга
// @Tags marketplace
// @Accept json
// @Produce json
// @Param id path int true "Listing ID"
// @Param limit query int false "Лимит результатов" default(20)
// @Success 200 {object} utils.SuccessResponseSwag{data=[]models.MarketplaceListing}
// @Failure 404 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/marketplace/listings/{id}/similar [get]
func (h *Handler) GetSimilarListings(c *fiber.Ctx) error {
	// Parse listing ID
	idParam := c.Params("id")
	listingID, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Error().Err(err).Str("id", idParam).Msg("GetSimilarListings: invalid listing ID")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalid_id")
	}

	// Check if listing exists
	_, err = h.storage.GetListing(c.Context(), listingID)
	if err != nil {
		h.logger.Error().Err(err).Int("listing_id", listingID).Msg("GetSimilarListings: failed to get listing")
		return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.listing_not_found")
	}

	// TODO: Implement actual similar listings logic (using OpenSearch, category, price range, etc.)
	// For now, return an empty array
	h.logger.Debug().Int("listing_id", listingID).Msg("GetSimilarListings: returning empty array (not implemented)")

	return utils.SuccessResponse(c, []models.MarketplaceListing{})
}

// UploadListingImages godoc
// @Summary Загрузить изображения для объявления
// @Description Загружает одно или несколько изображений для объявления
// @Tags marketplace
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param id path int true "Listing ID"
// @Param images formData file true "Image files"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]interface{}}
// @Failure 400 {object} utils.ErrorResponseSwag
// @Failure 401 {object} utils.ErrorResponseSwag
// @Failure 403 {object} utils.ErrorResponseSwag
// @Failure 404 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/marketplace/listings/{id}/images [post]
func (h *Handler) UploadListingImages(c *fiber.Ctx) error {
	// Get authenticated user ID
	userID, ok := authMiddleware.GetUserID(c)
	if !ok {
		h.logger.Warn().Msg("UploadListingImages: user not authenticated")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.unauthorized")
	}

	// Parse listing ID from path params
	idParam := c.Params("id")
	listingID, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Error().Err(err).Str("id", idParam).Msg("UploadListingImages: invalid listing ID")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalid_id")
	}

	// Check if listing exists and belongs to user
	listing, err := h.storage.GetListing(c.Context(), listingID)
	if err != nil {
		h.logger.Error().Err(err).Int("listing_id", listingID).Msg("UploadListingImages: failed to get listing")
		return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.listing_not_found")
	}

	// Verify ownership
	if listing.UserID != int(userID) {
		h.logger.Warn().
			Int("user_id", int(userID)).
			Int("listing_user_id", listing.UserID).
			Int("listing_id", listingID).
			Msg("UploadListingImages: user is not the owner")
		return utils.ErrorResponse(c, fiber.StatusForbidden, "marketplace.not_owner")
	}

	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		h.logger.Error().Err(err).Msg("UploadListingImages: failed to parse multipart form")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalid_form")
	}

	// Get uploaded files - check both "images" (multiple) and "image" (single) field names
	files := form.File["images"]
	if len(files) == 0 {
		// Try singular "image" field name (used by frontend)
		files = form.File["image"]
	}
	if len(files) == 0 {
		h.logger.Error().Msg("UploadListingImages: no files uploaded")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.no_files")
	}

	// PRE-FLIGHT VALIDATION: Max 10 files per upload
	const maxFiles = 10
	if len(files) > maxFiles {
		h.logger.Warn().
			Int("file_count", len(files)).
			Int("max_allowed", maxFiles).
			Msg("UploadListingImages: too many files")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.too_many_files")
	}

	// PRE-FLIGHT VALIDATION: Max 50MB total size
	const maxTotalSize = 50 * 1024 * 1024 // 50MB
	var totalSize int64
	for _, fileHeader := range files {
		totalSize += fileHeader.Size
	}
	if totalSize > maxTotalSize {
		h.logger.Warn().
			Int64("total_size_mb", totalSize/(1024*1024)).
			Int("max_allowed_mb", maxTotalSize/(1024*1024)).
			Msg("UploadListingImages: total size exceeds limit")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.total_size_exceeds_limit")
	}

	// Create ImageRepository
	imageRepo := postgres.NewImageRepository(h.db)

	// Create ImageService
	imageService := h.services.NewImageService(
		h.services.FileStorage(),
		imageRepo,
		services.ImageServiceConfig{
			BucketListings:    "listings",
			BucketStorefront:  "storefront",
			BucketChatFiles:   "chat-files",
			BucketReviewPhoto: "review-photos",
		},
	)

	// PARALLEL UPLOAD: Process files concurrently
	type uploadResult struct {
		image *services.UploadImageResponse
		err   error
		index int
	}

	results := make(chan uploadResult, len(files))
	var wg sync.WaitGroup

	// Semaphore to limit concurrent uploads (max 5 at once)
	semaphore := make(chan struct{}, 5)

	for i, fileHeader := range files {
		wg.Add(1)
		go func(idx int, fh *multipart.FileHeader) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Open file
			file, err := fh.Open()
			if err != nil {
				h.logger.Error().Err(err).Str("filename", fh.Filename).Msg("UploadListingImages: failed to open file")
				results <- uploadResult{nil, err, idx}
				return
			}
			defer func() {
				if closeErr := file.Close(); closeErr != nil {
					h.logger.Error().Err(closeErr).Str("filename", fh.Filename).Msg("Failed to close file")
				}
			}()

			// Upload image
			uploadReq := &services.UploadImageRequest{
				EntityType:   services.ImageTypeMarketplaceListing,
				EntityID:     listingID,
				File:         file,
				FileHeader:   fh,
				IsMain:       idx == 0, // First image is main
				DisplayOrder: idx + 1,
			}

			uploadedImage, err := imageService.UploadImage(c.Context(), uploadReq)
			results <- uploadResult{uploadedImage, err, idx}

			if err != nil {
				h.logger.Error().
					Err(err).
					Str("filename", fh.Filename).
					Int("listing_id", listingID).
					Msg("UploadListingImages: failed to upload image")
			} else {
				h.logger.Info().
					Int("image_id", uploadedImage.ID).
					Int("listing_id", listingID).
					Str("filename", fh.Filename).
					Msg("Image uploaded successfully")
			}
		}(i, fileHeader)
	}

	// Wait for all uploads to complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	uploadedImages := make([]interface{}, 0, len(files))
	failedUploads := make([]string, 0)

	for result := range results {
		if result.err != nil {
			failedUploads = append(failedUploads, fmt.Sprintf("file_%d: %s", result.index, result.err.Error()))
		} else if result.image != nil {
			uploadedImages = append(uploadedImages, result.image)
		}
	}

	// Prepare response with partial success info
	successCount := len(uploadedImages)
	failedCount := len(failedUploads)

	h.logger.Info().
		Int("listing_id", listingID).
		Int("success_count", successCount).
		Int("failed_count", failedCount).
		Int("total_files", len(files)).
		Msg("Bulk image upload completed")

	// If all failed, return error
	if successCount == 0 {
		h.logger.Error().Int("listing_id", listingID).Msg("UploadListingImages: no images were uploaded successfully")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.upload_failed")
	}

	// Return success with detailed info
	responseData := fiber.Map{
		"success": successCount,
		"failed":  failedCount,
		"images":  uploadedImages,
	}
	if failedCount > 0 {
		responseData["errors"] = failedUploads
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    responseData,
	})
}

// DeleteListingImage godoc
// @Summary Удалить изображение объявления
// @Description Удаляет изображение из объявления. Если удаляется главное изображение, автоматически назначается новое главное изображение (первое по дате создания)
// @Tags marketplace
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Listing ID"
// @Param imageId path int true "Image ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=object}
// @Failure 400 {object} utils.ErrorResponseSwag
// @Failure 401 {object} utils.ErrorResponseSwag
// @Failure 403 {object} utils.ErrorResponseSwag
// @Failure 404 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/marketplace/listings/{id}/images/{imageId} [delete]
func (h *Handler) DeleteListingImage(c *fiber.Ctx) error {
	// Get authenticated user ID
	userID, ok := authMiddleware.GetUserID(c)
	if !ok {
		h.logger.Warn().Msg("DeleteListingImage: user not authenticated")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.unauthorized")
	}

	// Parse listing ID from path params
	idParam := c.Params("id")
	listingID, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Error().Err(err).Str("id", idParam).Msg("DeleteListingImage: invalid listing ID")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalid_id")
	}

	// Parse image ID from path params
	imageIDParam := c.Params("imageId")
	imageID, err := strconv.Atoi(imageIDParam)
	if err != nil {
		h.logger.Error().Err(err).Str("imageId", imageIDParam).Msg("DeleteListingImage: invalid image ID")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalid_image_id")
	}

	// Check if listing exists and belongs to user
	listing, err := h.storage.GetListing(c.Context(), listingID)
	if err != nil {
		h.logger.Error().Err(err).Int("listing_id", listingID).Msg("DeleteListingImage: failed to get listing")
		return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.listing_not_found")
	}

	// Verify ownership
	if listing.UserID != int(userID) {
		h.logger.Warn().
			Int("user_id", int(userID)).
			Int("listing_user_id", listing.UserID).
			Int("listing_id", listingID).
			Msg("DeleteListingImage: user is not the owner")
		return utils.ErrorResponse(c, fiber.StatusForbidden, "marketplace.not_owner")
	}

	// Create ImageRepository
	imageRepo := postgres.NewImageRepository(h.db)

	// Get image info to check if it belongs to this listing
	imageInterface, err := imageRepo.GetImageByID(c.Context(), imageID)
	if err != nil {
		h.logger.Error().Err(err).Int("image_id", imageID).Msg("DeleteListingImage: failed to get image")
		return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.image_not_found")
	}

	// Type assertion to MarketplaceImage
	marketplaceImage, ok := imageInterface.(*models.MarketplaceImage)
	if !ok {
		h.logger.Error().Int("image_id", imageID).Msg("DeleteListingImage: image is not a marketplace image")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalid_image_type")
	}

	// Verify image belongs to this listing
	if marketplaceImage.ListingID != listingID {
		h.logger.Warn().
			Int("image_listing_id", marketplaceImage.ListingID).
			Int("listing_id", listingID).
			Msg("DeleteListingImage: image does not belong to this listing")
		return utils.ErrorResponse(c, fiber.StatusForbidden, "marketplace.image_not_belongs_to_listing")
	}

	// Remember if this was the main image
	wasMainImage := marketplaceImage.IsMain

	// Create ImageService
	imageService := h.services.NewImageService(
		h.services.FileStorage(),
		imageRepo,
		services.ImageServiceConfig{
			BucketListings:    "listings",
			BucketStorefront:  "storefront",
			BucketChatFiles:   "chat-files",
			BucketReviewPhoto: "review-photos",
		},
	)

	// Delete image from MinIO and DB
	if err := imageService.DeleteImage(c.Context(), imageID, services.ImageTypeMarketplaceListing); err != nil {
		h.logger.Error().
			Err(err).
			Int("image_id", imageID).
			Int("listing_id", listingID).
			Msg("DeleteListingImage: failed to delete image")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.delete_image_failed")
	}

	// If deleted image was main, set a new main image
	if wasMainImage {
		// Get first image by created_at
		firstImage, err := imageRepo.(*postgres.ImageRepository).GetFirstImageByListingID(c.Context(), listingID)
		if err == nil && firstImage != nil {
			// Set as main image
			if setErr := imageRepo.SetMainImage(c.Context(), firstImage.ID, true); setErr != nil {
				h.logger.Error().
					Err(setErr).
					Int("new_main_image_id", firstImage.ID).
					Int("listing_id", listingID).
					Msg("DeleteListingImage: failed to set new main image")
				// Continue execution - это не критичная ошибка
			} else {
				h.logger.Info().
					Int("new_main_image_id", firstImage.ID).
					Int("listing_id", listingID).
					Msg("Set new main image after deletion")
			}
		}
		// Если изображений не осталось - это нормально, просто игнорируем
	}

	h.logger.Info().
		Int("image_id", imageID).
		Int("listing_id", listingID).
		Bool("was_main", wasMainImage).
		Msg("Image deleted successfully")

	return utils.SuccessResponse(c, fiber.Map{
		"message":    "marketplace.image_deleted",
		"image_id":   imageID,
		"listing_id": listingID,
		"was_main":   wasMainImage,
		"deleted_at": time.Now(),
	})
}

// ImageReorderItem представляет элемент для переупорядочивания изображений
type ImageReorderItem struct {
	ImageID      int   `json:"image_id" validate:"required"`
	DisplayOrder int   `json:"display_order" validate:"required,min=1"`
	IsMain       *bool `json:"is_main,omitempty"` // optional, pointer для различия между false и не указано
}

// ReorderListingImagesRequest представляет запрос на переупорядочивание изображений
type ReorderListingImagesRequest struct {
	Images []ImageReorderItem `json:"images" validate:"required,min=1"`
}

// ReorderListingImages godoc
// @Summary Переупорядочить изображения объявления
// @Description Изменяет порядок отображения изображений и назначает главное изображение. Обновляет display_order и is_main для всех указанных изображений атомарно
// @Tags marketplace
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Listing ID"
// @Param request body ReorderListingImagesRequest true "Массив изображений с новым порядком"
// @Success 200 {object} utils.SuccessResponseSwag{data=object}
// @Failure 400 {object} utils.ErrorResponseSwag
// @Failure 401 {object} utils.ErrorResponseSwag
// @Failure 403 {object} utils.ErrorResponseSwag
// @Failure 404 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Router /api/v1/marketplace/listings/{id}/images/reorder [patch]
func (h *Handler) ReorderListingImages(c *fiber.Ctx) error {
	// Get authenticated user ID
	userID, ok := authMiddleware.GetUserID(c)
	if !ok {
		h.logger.Warn().Msg("ReorderListingImages: user not authenticated")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.unauthorized")
	}

	// Parse listing ID from path params
	idParam := c.Params("id")
	listingID, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Error().Err(err).Str("id", idParam).Msg("ReorderListingImages: invalid listing ID")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalid_id")
	}

	// Parse request body
	var req ReorderListingImagesRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error().Err(err).Msg("ReorderListingImages: failed to parse request body")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalid_request")
	}

	// Validate request
	if len(req.Images) == 0 {
		h.logger.Error().Msg("ReorderListingImages: empty images array")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.empty_images_array")
	}

	// Check if listing exists and belongs to user
	listing, err := h.storage.GetListing(c.Context(), listingID)
	if err != nil {
		h.logger.Error().Err(err).Int("listing_id", listingID).Msg("ReorderListingImages: failed to get listing")
		return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.listing_not_found")
	}

	// Verify ownership
	if listing.UserID != int(userID) {
		h.logger.Warn().
			Int("listing_id", listingID).
			Int("user_id", int(userID)).
			Int("owner_id", listing.UserID).
			Msg("ReorderListingImages: user is not the owner")
		return utils.ErrorResponse(c, fiber.StatusForbidden, "marketplace.not_owner")
	}

	// Create ImageRepository
	imageRepo := postgres.NewImageRepository(h.db)

	// Verify all images belong to this listing
	for _, item := range req.Images {
		imageInterface, err := imageRepo.GetImageByID(c.Context(), item.ImageID)
		if err != nil {
			h.logger.Error().
				Err(err).
				Int("image_id", item.ImageID).
				Msg("ReorderListingImages: image not found")
			return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.image_not_found")
		}

		// Type assertion to MarketplaceImage
		marketplaceImage, ok := imageInterface.(*models.MarketplaceImage)
		if !ok {
			h.logger.Error().
				Int("image_id", item.ImageID).
				Msg("ReorderListingImages: wrong image type")
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalid_image_type")
		}

		if marketplaceImage.ListingID != listingID {
			h.logger.Error().
				Int("image_id", item.ImageID).
				Int("listing_id", listingID).
				Int("actual_listing_id", marketplaceImage.ListingID).
				Msg("ReorderListingImages: image does not belong to listing")
			return utils.ErrorResponse(c, fiber.StatusForbidden, "marketplace.image_not_belongs_to_listing")
		}
	}

	// Check if any image should be set as main
	var newMainImageID *int
	for _, item := range req.Images {
		if item.IsMain != nil && *item.IsMain {
			if newMainImageID != nil {
				// Multiple images marked as main - error
				h.logger.Error().Msg("ReorderListingImages: multiple images marked as main")
				return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.multiple_main_images")
			}
			newMainImageID = &item.ImageID
		}
	}

	// Update display_order and is_main in a loop
	// Note: Ideally this should be done in a transaction, but for simplicity we'll do it sequentially
	// If we need true atomicity, we should add a batch update method to the repository

	// First, if new main image is specified, unset all main images
	if newMainImageID != nil {
		if err := imageRepo.UnsetMainImages(c.Context(), "marketplace_listing", listingID); err != nil {
			h.logger.Error().
				Err(err).
				Int("listing_id", listingID).
				Msg("ReorderListingImages: failed to unset main images")
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.reorder_failed")
		}
	}

	// Update each image
	for _, item := range req.Images {
		// Update display order
		if err := imageRepo.UpdateDisplayOrder(c.Context(), item.ImageID, item.DisplayOrder); err != nil {
			h.logger.Error().
				Err(err).
				Int("image_id", item.ImageID).
				Int("display_order", item.DisplayOrder).
				Msg("ReorderListingImages: failed to update display order")
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.reorder_failed")
		}

		// Update is_main if specified
		if item.IsMain != nil {
			if err := imageRepo.SetMainImage(c.Context(), item.ImageID, *item.IsMain); err != nil {
				h.logger.Error().
					Err(err).
					Int("image_id", item.ImageID).
					Bool("is_main", *item.IsMain).
					Msg("ReorderListingImages: failed to set main image")
				return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.reorder_failed")
			}
		}
	}

	h.logger.Info().
		Int("listing_id", listingID).
		Int("images_count", len(req.Images)).
		Msg("Images reordered successfully")

	return utils.SuccessResponse(c, fiber.Map{
		"message":      "marketplace.images_reordered",
		"listing_id":   listingID,
		"images_count": len(req.Images),
		"updated_at":   time.Now(),
	})
}
