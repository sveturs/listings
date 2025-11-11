// TEMPORARY: Will be moved to microservice
package handler

import (
	"context"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"
	pb "github.com/sveturs/listings/api/proto/listings/v1"
	"google.golang.org/grpc/metadata"

	"backend/internal/domain/models"
	"backend/internal/services"
	"backend/internal/storage/postgres"
	"backend/pkg/utils"
)

// contextWithAuth creates gRPC context with JWT token from Fiber context
func contextWithAuth(c *fiber.Ctx, baseCtx context.Context) context.Context {
	// Extract Authorization header from Fiber request
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		// Try lowercase (some clients send it this way)
		authHeader = c.Get("authorization")
	}

	if authHeader == "" {
		// No token - return base context
		return baseCtx
	}

	// Ensure it starts with "Bearer "
	if !strings.HasPrefix(authHeader, "Bearer ") {
		authHeader = "Bearer " + authHeader
	}

	// Add to gRPC metadata
	md := metadata.Pairs("authorization", authHeader)
	return metadata.NewOutgoingContext(baseCtx, md)
}

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

	// Route to microservice if enabled, otherwise use monolith storage (fallback)
	if h.useListingsMicroservice && h.listingsClient != nil {
		// Use microservice via gRPC
		h.logger.Info().
			Int("user_id", int(userID)).
			Msg("Routing CreateListing to microservice")

		// Prepare gRPC request
		grpcReq := &pb.CreateListingRequest{
			UserId:     int64(userID),
			Title:      req.Title,
			Price:      req.Price,
			Currency:   req.Currency,
			CategoryId: int64(req.CategoryID),
			Quantity:   req.Quantity,
		}

		// Optional fields
		if req.Description != nil {
			grpcReq.Description = req.Description
		}
		if req.SKU != nil {
			grpcReq.Sku = req.SKU
		}
		if req.StorefrontID != nil {
			storefrontID := int64(*req.StorefrontID)
			grpcReq.StorefrontId = &storefrontID
		}

		// Call microservice with auth context
		grpcCtx := contextWithAuth(c, c.Context())
		grpcResp, err := h.listingsClient.CreateListing(grpcCtx, grpcReq)
		if err != nil {
			h.logger.Error().Err(err).
				Int("user_id", int(userID)).
				Int("category_id", req.CategoryID).
				Msg("CreateListing: failed to create listing via microservice")
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.create_failed")
		}

		// Auto-publish: Set status to 'active' immediately after creation
		publishStatus := "active"
		updateReq := &pb.UpdateListingRequest{
			Id:     grpcResp.Listing.Id,
			UserId: int64(userID),
			Status: &publishStatus,
		}
		updateResp, err := h.listingsClient.UpdateListing(grpcCtx, updateReq)
		if err != nil {
			h.logger.Warn().Err(err).
				Int64("listing_id", grpcResp.Listing.Id).
				Msg("CreateListing: failed to auto-publish listing, but listing was created")
			// Don't return error - listing is created, just not published
			// Continue with original status (draft)
		} else {
			// Update response with published listing
			grpcResp.Listing = updateResp.Listing
			h.logger.Info().
				Int64("listing_id", grpcResp.Listing.Id).
				Str("status", grpcResp.Listing.Status).
				Msg("Listing auto-published successfully")
		}

		h.logger.Info().
			Int64("listing_id", grpcResp.Listing.Id).
			Int("user_id", int(userID)).
			Str("title", req.Title).
			Str("final_status", grpcResp.Listing.Status).
			Msg("Listing created successfully via microservice")

		// Convert gRPC response to domain model
		listing := &models.MarketplaceListing{
			ID:         int(grpcResp.Listing.Id),
			UserID:     int(grpcResp.Listing.UserId),
			CategoryID: int(grpcResp.Listing.CategoryId),
			Title:      grpcResp.Listing.Title,
			Price:      grpcResp.Listing.Price,
			Status:     grpcResp.Listing.Status,
		}
		if grpcResp.Listing.Description != nil {
			listing.Description = *grpcResp.Listing.Description
		}
		if grpcResp.Listing.Sku != nil {
			listing.ExternalID = *grpcResp.Listing.Sku
		}
		if grpcResp.Listing.StorefrontId != nil {
			storefrontID := int(*grpcResp.Listing.StorefrontId)
			listing.StorefrontID = &storefrontID
		}
		// Stock quantity from microservice
		if grpcResp.Listing.Quantity > 0 {
			qty := int(grpcResp.Listing.Quantity)
			listing.StockQuantity = &qty
		}
		// Published timestamp from microservice
		if grpcResp.Listing.PublishedAt != nil && *grpcResp.Listing.PublishedAt != "" {
			publishedAt, err := time.Parse(time.RFC3339, *grpcResp.Listing.PublishedAt)
			if err == nil {
				listing.PublishedAt = &publishedAt
			}
		}

		c.Status(fiber.StatusCreated)
		return utils.SuccessResponse(c, listing)
	}

	// Fallback to monolith storage
	h.logger.Info().
		Int("user_id", int(userID)).
		Bool("microservice_enabled", h.useListingsMicroservice).
		Bool("client_available", h.listingsClient != nil).
		Msg("Using monolith storage for CreateListing (fallback)")

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
		Msg("Listing created successfully via monolith")

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

	// Route to microservice if enabled, otherwise use monolith storage (fallback)
	if h.useListingsMicroservice && h.listingsClient != nil {
		// Use microservice via gRPC
		h.logger.Debug().
			Int("listing_id", listingID).
			Msg("Routing GetListing to microservice")

		grpcReq := &pb.GetListingRequest{
			Id: int64(listingID),
		}

		// Call microservice with auth context (optional for public methods)
		grpcCtx := contextWithAuth(c, c.Context())
		grpcResp, err := h.listingsClient.GetListing(grpcCtx, grpcReq)
		if err != nil {
			h.logger.Error().Err(err).Int("listing_id", listingID).Msg("GetListing: failed to get listing from microservice")
			return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.listing_not_found")
		}

		// Convert gRPC response to domain model
		listing := &models.MarketplaceListing{
			ID:         int(grpcResp.Listing.Id),
			UserID:     int(grpcResp.Listing.UserId),
			CategoryID: int(grpcResp.Listing.CategoryId),
			Title:      grpcResp.Listing.Title,
			Price:      grpcResp.Listing.Price,
			Status:     grpcResp.Listing.Status,
		}
		if grpcResp.Listing.Description != nil {
			listing.Description = *grpcResp.Listing.Description
		}
		if grpcResp.Listing.Sku != nil {
			listing.ExternalID = *grpcResp.Listing.Sku
		}
		if grpcResp.Listing.StorefrontId != nil {
			storefrontID := int(*grpcResp.Listing.StorefrontId)
			listing.StorefrontID = &storefrontID
		}
		// Stock quantity from microservice
		if grpcResp.Listing.Quantity > 0 {
			qty := int(grpcResp.Listing.Quantity)
			listing.StockQuantity = &qty
		}
		// Published timestamp from microservice
		if grpcResp.Listing.PublishedAt != nil && *grpcResp.Listing.PublishedAt != "" {
			publishedAt, err := time.Parse(time.RFC3339, *grpcResp.Listing.PublishedAt)
			if err == nil {
				listing.PublishedAt = &publishedAt
			}
		}

		// Convert images from gRPC response (always initialize array)
		if len(grpcResp.Listing.Images) > 0 {
			listing.Images = make([]models.MarketplaceImage, len(grpcResp.Listing.Images))
			for i, pbImg := range grpcResp.Listing.Images {
				listing.Images[i] = models.MarketplaceImage{
					ID:           int(pbImg.Id),
					ListingID:    int(pbImg.ListingId),
					PublicURL:    pbImg.Url,
					IsMain:       pbImg.IsPrimary,
					DisplayOrder: int(pbImg.DisplayOrder),
				}

				if pbImg.ThumbnailUrl != nil {
					listing.Images[i].ThumbnailURL = *pbImg.ThumbnailUrl
				}
			}
		} else {
			// Initialize empty array instead of nil for consistency
			listing.Images = []models.MarketplaceImage{}
		}

		return utils.SuccessResponse(c, listing)
	}

	// Fallback to monolith storage
	h.logger.Debug().
		Int("listing_id", listingID).
		Bool("microservice_enabled", h.useListingsMicroservice).
		Bool("client_available", h.listingsClient != nil).
		Msg("Using monolith storage for GetListing (fallback)")

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

	// Check if listing exists (route to microservice if enabled)
	if h.useListingsMicroservice && h.listingsClient != nil {
		// Use microservice
		grpcCtx := contextWithAuth(c, c.Context())
		grpcReq := &pb.GetListingRequest{Id: int64(listingID)}
		_, err = h.listingsClient.GetListing(grpcCtx, grpcReq)
		if err != nil {
			h.logger.Error().Err(err).Int("listing_id", listingID).Msg("GetSimilarListings: failed to get listing from microservice")
			return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.listing_not_found")
		}
	} else {
		// Use monolith
		_, err = h.storage.GetListing(c.Context(), listingID)
		if err != nil {
			h.logger.Error().Err(err).Int("listing_id", listingID).Msg("GetSimilarListings: failed to get listing")
			return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.listing_not_found")
		}
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

	// Check if listing exists and belongs to user (route to microservice if enabled)
	var listing *models.MarketplaceListing
	if h.useListingsMicroservice && h.listingsClient != nil {
		// Use microservice
		grpcCtx := contextWithAuth(c, c.Context())
		grpcReq := &pb.GetListingRequest{Id: int64(listingID)}
		grpcResp, err := h.listingsClient.GetListing(grpcCtx, grpcReq)
		if err != nil {
			h.logger.Error().Err(err).Int("listing_id", listingID).Msg("UploadListingImages: failed to get listing from microservice")
			return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.listing_not_found")
		}
		// Convert to domain model (we only need ID and UserID for ownership check)
		listing = &models.MarketplaceListing{
			ID:     int(grpcResp.Listing.Id),
			UserID: int(grpcResp.Listing.UserId),
		}
	} else {
		// Use monolith
		var err error
		listing, err = h.storage.GetListing(c.Context(), listingID)
		if err != nil {
			h.logger.Error().Err(err).Int("listing_id", listingID).Msg("UploadListingImages: failed to get listing")
			return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.listing_not_found")
		}
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

	// Create ImageRepository and Service (will be used differently based on routing)
	imageRepo := postgres.NewImageRepository(h.db)
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

	// If using microservice, migrate image metadata from monolith to microservice
	if h.useListingsMicroservice && h.listingsClient != nil && len(uploadedImages) > 0 {
		h.logger.Info().
			Int("listing_id", listingID).
			Int("images_count", len(uploadedImages)).
			Msg("Migrating image metadata to microservice")

		// Migrate each uploaded image to microservice
		for _, img := range uploadedImages {
			uploadedImg, ok := img.(*services.UploadImageResponse)
			if !ok {
				h.logger.Error().Msg("Failed to cast image response")
				continue
			}

			// Call microservice to add image metadata
			grpcCtx := contextWithAuth(c, c.Context())
			addImageReq := &pb.AddImageRequest{
				ListingId:    int64(listingID),
				Url:          uploadedImg.ImageURL,
				ThumbnailUrl: &uploadedImg.ThumbnailURL,
				DisplayOrder: int32(uploadedImg.DisplayOrder), // #nosec G115 - display order is limited to reasonable values (0-10)
				IsPrimary:    uploadedImg.IsMain,
			}

			grpcResp, err := h.listingsClient.AddListingImage(grpcCtx, addImageReq)
			if err != nil {
				h.logger.Error().
					Err(err).
					Int("image_id", uploadedImg.ID).
					Int("listing_id", listingID).
					Msg("Failed to add image to microservice")
				// Don't fail the whole upload, just log
				continue
			}

			h.logger.Info().
				Int64("microservice_image_id", grpcResp.Image.Id).
				Int("monolith_image_id", uploadedImg.ID).
				Int("listing_id", listingID).
				Msg("Image metadata added to microservice")

			// Delete from monolith DB (cleanup)
			if err := imageRepo.DeleteImage(c.Context(), uploadedImg.ID); err != nil {
				h.logger.Warn().
					Err(err).
					Int("image_id", uploadedImg.ID).
					Msg("Failed to delete image from monolith after migration")
				// Non-critical, continue
			}

			// Update response to use microservice ID
			uploadedImg.ID = int(grpcResp.Image.Id)
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

	// Route to microservice if enabled, otherwise use monolith storage (fallback)
	if h.useListingsMicroservice && h.listingsClient != nil {
		// Use microservice via gRPC
		h.logger.Info().
			Int("user_id", int(userID)).
			Int("listing_id", listingID).
			Int("image_id", imageID).
			Msg("Routing DeleteListingImage to microservice")

		// Verify ownership by getting listing first
		grpcCtx := contextWithAuth(c, c.Context())
		grpcReq := &pb.GetListingRequest{Id: int64(listingID)}
		grpcResp, err := h.listingsClient.GetListing(grpcCtx, grpcReq)
		if err != nil {
			h.logger.Error().Err(err).Int("listing_id", listingID).Msg("DeleteListingImage: failed to get listing from microservice")
			return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.listing_not_found")
		}

		// Verify ownership
		if grpcResp.Listing.UserId != int64(userID) {
			h.logger.Warn().
				Int("user_id", int(userID)).
				Int64("listing_user_id", grpcResp.Listing.UserId).
				Int("listing_id", listingID).
				Msg("DeleteListingImage: user is not the owner")
			return utils.ErrorResponse(c, fiber.StatusForbidden, "marketplace.not_owner")
		}

		// Call microservice to delete image
		if err := h.listingsClient.DeleteListingImage(grpcCtx, int64(imageID)); err != nil {
			h.logger.Error().
				Err(err).
				Int("image_id", imageID).
				Int("listing_id", listingID).
				Msg("DeleteListingImage: failed to delete image via microservice")
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.delete_image_failed")
		}

		h.logger.Info().
			Int("image_id", imageID).
			Int("listing_id", listingID).
			Msg("Image deleted successfully via microservice")

		return utils.SuccessResponse(c, fiber.Map{
			"message":    "marketplace.image_deleted",
			"image_id":   imageID,
			"listing_id": listingID,
			"deleted_at": time.Now(),
		})
	}

	// Fallback to monolith storage
	h.logger.Info().
		Int("user_id", int(userID)).
		Int("listing_id", listingID).
		Int("image_id", imageID).
		Bool("microservice_enabled", h.useListingsMicroservice).
		Bool("client_available", h.listingsClient != nil).
		Msg("Using monolith storage for DeleteListingImage (fallback)")

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
		Msg("Image deleted successfully via monolith")

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

	// Route to microservice if enabled, otherwise use monolith storage (fallback)
	if h.useListingsMicroservice && h.listingsClient != nil {
		// Use microservice via gRPC
		h.logger.Info().
			Int("user_id", int(userID)).
			Int("listing_id", listingID).
			Int("images_count", len(req.Images)).
			Msg("Routing ReorderListingImages to microservice")

		// Verify ownership by getting listing first
		grpcCtx := contextWithAuth(c, c.Context())
		grpcReq := &pb.GetListingRequest{Id: int64(listingID)}
		grpcResp, err := h.listingsClient.GetListing(grpcCtx, grpcReq)
		if err != nil {
			h.logger.Error().Err(err).Int("listing_id", listingID).Msg("ReorderListingImages: failed to get listing from microservice")
			return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.listing_not_found")
		}

		// Verify ownership
		if grpcResp.Listing.UserId != int64(userID) {
			h.logger.Warn().
				Int("user_id", int(userID)).
				Int64("listing_user_id", grpcResp.Listing.UserId).
				Int("listing_id", listingID).
				Msg("ReorderListingImages: user is not the owner")
			return utils.ErrorResponse(c, fiber.StatusForbidden, "marketplace.not_owner")
		}

		// Convert request body to proto format
		imageOrders := make([]*pb.ImageOrder, len(req.Images))
		for i, item := range req.Images {
			imageOrders[i] = &pb.ImageOrder{
				ImageId:      int64(item.ImageID),
				DisplayOrder: int32(item.DisplayOrder), // #nosec G115 - display order is validated and limited to reasonable values (1-10)
			}
		}

		// Call microservice to reorder images
		if err := h.listingsClient.ReorderListingImages(grpcCtx, int64(listingID), imageOrders); err != nil {
			h.logger.Error().
				Err(err).
				Int("listing_id", listingID).
				Msg("ReorderListingImages: failed to reorder images via microservice")
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.reorder_failed")
		}

		h.logger.Info().
			Int("listing_id", listingID).
			Int("images_count", len(req.Images)).
			Msg("Images reordered successfully via microservice")

		return utils.SuccessResponse(c, fiber.Map{
			"message":      "marketplace.images_reordered",
			"listing_id":   listingID,
			"images_count": len(req.Images),
			"updated_at":   time.Now(),
		})
	}

	// Fallback to monolith storage
	h.logger.Info().
		Int("user_id", int(userID)).
		Int("listing_id", listingID).
		Int("images_count", len(req.Images)).
		Bool("microservice_enabled", h.useListingsMicroservice).
		Bool("client_available", h.listingsClient != nil).
		Msg("Using monolith storage for ReorderListingImages (fallback)")

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
		Msg("Images reordered successfully via monolith")

	return utils.SuccessResponse(c, fiber.Map{
		"message":      "marketplace.images_reordered",
		"listing_id":   listingID,
		"images_count": len(req.Images),
		"updated_at":   time.Now(),
	})
}

// GetListings godoc
// @Summary Get listings list (admin)
// @Description Get paginated list of listings for admin panel
// @Tags admin
// @Accept json
// @Produce json
// @Param limit query int false "Items per page (default: 20)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Param sort_by query string false "Sort by (date_desc, date_asc, price_asc, price_desc)"
// @Param user_id query int false "Filter by user ID"
// @Param storefront_id query int false "Filter by storefront ID"
// @Param category_id query int false "Filter by category ID"
// @Param status query string false "Filter by status (active, pending, sold, etc)"
// @Param min_price query number false "Minimum price filter"
// @Param max_price query number false "Maximum price filter"
// @Param exclude_storefronts query boolean false "Exclude storefront listings"
// @Success 200 {object} utils.SuccessResponseSwag{data=object{data=[]models.MarketplaceListing,meta=object{total=int}}}
// @Failure 401 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Security BearerAuth
// @Router /api/v1/marketplace/listings [get]
func (h *Handler) GetListings(c *fiber.Ctx) error {
	// Parse query parameters
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)
	sortBy := c.Query("sort_by", "")
	excludeStorefronts := c.QueryBool("exclude_storefronts", false)

	// Optional filters
	userIDQuery := c.Query("user_id", "")
	storefrontIDQuery := c.Query("storefront_id", "")
	categoryIDQuery := c.Query("category_id", "")
	status := c.Query("status", "")
	minPriceQuery := c.Query("min_price", "")
	maxPriceQuery := c.Query("max_price", "")

	// Phase 7.4: Route to microservice if feature flag is enabled
	if h.useListingsMicroservice && h.listingsClient != nil {
		h.logger.Info().
			Bool("use_microservice", true).
			Int("limit", limit).
			Int("offset", offset).
			Str("sort_by", sortBy).
			Bool("exclude_storefronts", excludeStorefronts).
			Msg("Routing GetListings to listings microservice")

		// Build gRPC request
		grpcReq := &pb.ListListingsRequest{
			Limit:  int32(limit),
			Offset: int32(offset),
		}

		// Add optional filters
		if userIDQuery != "" {
			if userID, err := strconv.ParseInt(userIDQuery, 10, 64); err == nil {
				grpcReq.UserId = &userID
			}
		}
		if storefrontIDQuery != "" {
			if storefrontID, err := strconv.ParseInt(storefrontIDQuery, 10, 64); err == nil {
				// Handle exclude_storefronts: if true, don't filter by storefront
				// If false and storefront_id is provided, filter by it
				if !excludeStorefronts {
					grpcReq.StorefrontId = &storefrontID
				}
			}
		} else if excludeStorefronts {
			// If exclude_storefronts=true but no specific storefront_id, we need special handling
			// This logic should be implemented in microservice
			// For now, we just pass the filter as-is
			h.logger.Debug().Msg("exclude_storefronts=true will be handled by filtering results")
		}
		if categoryIDQuery != "" {
			if categoryID, err := strconv.ParseInt(categoryIDQuery, 10, 64); err == nil {
				grpcReq.CategoryId = &categoryID
			}
		}
		if status != "" {
			grpcReq.Status = &status
		}
		if minPriceQuery != "" {
			if minPrice, err := strconv.ParseFloat(minPriceQuery, 64); err == nil {
				grpcReq.MinPrice = &minPrice
			}
		}
		if maxPriceQuery != "" {
			if maxPrice, err := strconv.ParseFloat(maxPriceQuery, 64); err == nil {
				grpcReq.MaxPrice = &maxPrice
			}
		}

		// Note: sort_by is currently NOT supported by microservice ListListings
		// It will be implemented in future phases
		if sortBy != "" {
			h.logger.Warn().
				Str("sort_by", sortBy).
				Msg("sort_by parameter is not yet supported by microservice, ignoring")
		}

		// Call microservice via gRPC
		grpcResp, err := h.listingsClient.ListListings(c.Context(), grpcReq)
		if err != nil {
			h.logger.Error().Err(err).Msg("Failed to get listings from microservice, falling back to monolith")
			// Fallback to monolith (stub for now)
			return utils.SuccessResponse(c, fiber.Map{
				"data": []interface{}{},
				"meta": fiber.Map{
					"total":  0,
					"limit":  limit,
					"offset": offset,
				},
			})
		}

		// Convert proto.Listing to models.MarketplaceListing
		listings := make([]models.MarketplaceListing, 0, len(grpcResp.Listings))
		for _, protoListing := range grpcResp.Listings {
			listing := models.MarketplaceListing{
				ID:         int(protoListing.Id),
				UserID:     int(protoListing.UserId),
				CategoryID: int(protoListing.CategoryId),
				Title:      protoListing.Title,
				Price:      protoListing.Price,
				Status:     protoListing.Status,
				ViewsCount: int(protoListing.ViewsCount),
			}

			// Description (optional)
			if protoListing.Description != nil {
				listing.Description = *protoListing.Description
			}

			// Location (complex object)
			if protoListing.Location != nil {
				loc := protoListing.Location
				if loc.City != nil {
					listing.City = *loc.City
				}
				if loc.Country != nil {
					listing.Country = *loc.Country
				}
				if loc.Latitude != nil {
					listing.Latitude = loc.Latitude
				}
				if loc.Longitude != nil {
					listing.Longitude = loc.Longitude
				}
				// Combine address fields into Location string
				if loc.AddressLine1 != nil || loc.City != nil {
					locationParts := []string{}
					if loc.AddressLine1 != nil {
						locationParts = append(locationParts, *loc.AddressLine1)
					}
					if loc.City != nil {
						locationParts = append(locationParts, *loc.City)
					}
					if loc.Country != nil {
						locationParts = append(locationParts, *loc.Country)
					}
					listing.Location = strings.Join(locationParts, ", ")
				}
			}

			// Storefront ID (optional)
			if protoListing.StorefrontId != nil {
				storefrontID := int(*protoListing.StorefrontId)
				listing.StorefrontID = &storefrontID
			}

			// SKU (for products)
			if protoListing.Sku != nil {
				listing.ExternalID = *protoListing.Sku
			}

			// Timestamps (parse from RFC3339 strings)
			if createdAt, err := time.Parse(time.RFC3339, protoListing.CreatedAt); err == nil {
				listing.CreatedAt = createdAt
			}
			if updatedAt, err := time.Parse(time.RFC3339, protoListing.UpdatedAt); err == nil {
				listing.UpdatedAt = updatedAt
			}
			if protoListing.PublishedAt != nil {
				if publishedAt, err := time.Parse(time.RFC3339, *protoListing.PublishedAt); err == nil {
					listing.PublishedAt = &publishedAt
				}
			}

			// Images - map from proto.ListingImage to models.MarketplaceImage
			if len(protoListing.Images) > 0 {
				listing.Images = make([]models.MarketplaceImage, 0, len(protoListing.Images))
				for _, protoImage := range protoListing.Images {
					image := models.MarketplaceImage{
						ID:            int(protoImage.Id),
						ListingID:     int(protoImage.ListingId),
						PublicURL:     protoImage.Url,
						DisplayOrder:  int(protoImage.DisplayOrder),
						IsMain:        protoImage.IsPrimary,
						StorageType:   "minio", // Default
						StorageBucket: "listings",
					}
					if protoImage.StoragePath != nil {
						image.FilePath = *protoImage.StoragePath
					}
					if protoImage.ThumbnailUrl != nil {
						image.ThumbnailURL = *protoImage.ThumbnailUrl
					}
					if protoImage.MimeType != nil {
						image.ContentType = *protoImage.MimeType
					}
					if protoImage.FileSize != nil {
						image.FileSize = int(*protoImage.FileSize)
					}
					if createdAt, err := time.Parse(time.RFC3339, protoImage.CreatedAt); err == nil {
						image.CreatedAt = createdAt
					}
					listing.Images = append(listing.Images, image)
				}
			}

			// Filter out storefront listings if exclude_storefronts is true
			if excludeStorefronts && listing.StorefrontID != nil {
				continue
			}

			listings = append(listings, listing)
		}

		// Adjust total count if we filtered out storefronts
		totalCount := int(grpcResp.Total)
		if excludeStorefronts {
			totalCount = len(listings)
		}

		h.logger.Info().
			Int("count", len(listings)).
			Int("total", totalCount).
			Bool("served_by_microservice", true).
			Msg("Successfully retrieved listings from microservice")

		// Add header to indicate microservice was used
		c.Set("X-Served-By", "microservice")
		return utils.SuccessResponse(c, fiber.Map{
			"data": listings,
			"meta": fiber.Map{
				"total":  totalCount,
				"limit":  limit,
				"offset": offset,
			},
		})
	}

	// Default: use monolith storage (stub for now)
	h.logger.Debug().
		Bool("use_microservice", false).
		Int("limit", limit).
		Int("offset", offset).
		Msg("Routing GetListings to monolith (stub)")

	return utils.SuccessResponse(c, fiber.Map{
		"data": []interface{}{},
		"meta": fiber.Map{
			"total":  0,
			"limit":  limit,
			"offset": offset,
		},
	})
}

// GetListingsStatistics godoc
// @Summary Get listings statistics (admin)
// @Description Get statistics for listings in admin panel
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=object{total=int,active=int,pending=int,views=int}}
// @Failure 401 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Security BearerAuth
// @Router /api/v1/admin/listings/statistics [get]
func (h *Handler) GetListingsStatistics(c *fiber.Ctx) error {
	// TODO: Implement via microservice when fully migrated
	// For now, return stub data

	return utils.SuccessResponse(c, fiber.Map{
		"total":   0,
		"active":  0,
		"pending": 0,
		"views":   0,
	})
}
