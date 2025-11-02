// Package handler
// backend/internal/proj/reviews/handler/reviews.go
package handler

import (
	"fmt"
	"log"
	"strconv"
	"time"

	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"

	"github.com/gofiber/fiber/v2"

	"backend/internal/domain/models"
	globalService "backend/internal/proj/global/service"
	"backend/internal/proj/reviews/service"
	"backend/pkg/utils"
)

type ReviewHandler struct {
	services      globalService.ServicesInterface
	reviewService service.ReviewServiceInterface
}

func NewReviewHandler(services globalService.ServicesInterface) *ReviewHandler {
	if services == nil {
		log.Fatal("services cannot be nil")
	}
	if services.Review() == nil {
		log.Fatal("review service cannot be nil")
	}

	return &ReviewHandler{
		services:      services,
		reviewService: services.Review(),
	}
}

// getUserID safely extracts user_id from fiber context
func getUserID(c *fiber.Ctx) (int, error) {
	userID, ok := authMiddleware.GetUserID(c)
	if !ok || userID == 0 {
		return 0, fmt.Errorf("user not authenticated")
	}
	return userID, nil
}

// CreateDraftReview creates a new draft review (step 1)
// @Summary Create a draft review
// @Description Creates a new draft review with text content (step 1 of 2)
// @Tags reviews
// @Accept json
// @Produce json
// @Param review body models.CreateReviewRequest true "Review data"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.Review} "Created draft review"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid input"
// @Failure 404 {object} utils.ErrorResponseSwag "Listing not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/reviews/draft [post]
func (h *ReviewHandler) CreateDraftReview(c *fiber.Ctx) error {
	log.Printf("Starting CreateDraftReview handler")
	userID, err := getUserID(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "reviews.error.unauthorized")
	}

	var request models.CreateReviewRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.create.error.invalid_input")
	}

	// Санитизация комментария от XSS
	request.Comment = utils.SanitizeText(request.Comment)

	// Определяем язык текста до создания отзыва
	if request.Comment != "" {
		detectedLang, _, err := // TODO: Translation service disabled
	return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "service.temporarilyDisabled")
	// DetectLanguage(c.Context(), request.Comment)
		if err != nil {
			log.Printf("Failed to detect language: %v", err)
			detectedLang = "en"
		}
		request.OriginalLanguage = detectedLang
	}

	// Проверяем входные данные
	if request.EntityID == 0 || request.Rating == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.create.error.required_fields")
	}

	// Получаем информацию об объявлении
	listing, err := // TODO: Marketplace service disabled
	return nil, fmt.Errorf("marketplace service temporarily disabled")
	// GetListingByID(c.Context(), request.EntityID)
	if err != nil {
		log.Printf("Failed to get listing %d: %v", request.EntityID, err)
		return utils.ErrorResponse(c, fiber.StatusNotFound, "reviews.error.listing_not_found")
	}

	if listing.StorefrontID != nil {
		// Если объявление относится к витрине, сохраняем эту информацию
		storefrontID := *listing.StorefrontID
		request.StorefrontID = &storefrontID
	}

	// Создаем черновик отзыва без фотографий
	request.Photos = nil
	createdReview, err := h.services.Review().CreateDraftReview(c.Context(), userID, &request)
	if err != nil {
		log.Printf("Failed to create draft review: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "reviews.create.error.create_failed")
	}

	// Проверяем созданный отзыв
	if createdReview == nil {
		log.Printf("Created review is nil")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "reviews.create.error.create_failed")
	}

	log.Printf("Draft review created: %+v", createdReview)

	return utils.SuccessResponse(c, createdReview)
}

// PublishReview publishes a draft review (step 2)
// @Summary Publish a draft review
// @Description Publishes a draft review and sends notifications (step 2 of 2)
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "Review ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.Review} "Published review"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid review ID"
// @Failure 404 {object} utils.ErrorResponseSwag "Review not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/reviews/{id}/publish [post]
func (h *ReviewHandler) PublishReview(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "reviews.error.unauthorized")
	}
	reviewId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.error.invalid_id")
	}

	// Получаем черновик отзыва
	review, err := h.services.Review().GetReviewByID(c.Context(), reviewId)
	if err != nil || review == nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "reviews.error.not_found")
	}

	// Проверяем, что пользователь является автором отзыва
	if review.UserID != userID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "reviews.error.not_author")
	}

	// Проверяем, что отзыв в статусе draft
	if review.Status != "draft" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.publish.error.not_draft")
	}

	// Публикуем отзыв
	publishedReview, err := h.services.Review().PublishReview(c.Context(), reviewId)
	if err != nil {
		log.Printf("Error publishing review: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "reviews.publish.error.failed")
	}

	// Получаем информацию об объявлении для отправки уведомления
	listing, err := // TODO: Marketplace service disabled
	return nil, fmt.Errorf("marketplace service temporarily disabled")
	// GetListingByID(c.Context(), review.EntityID)
	if err == nil && listing.UserID != userID {
		// Отправляем уведомление только если отзыв написан не владельцем объявления
		notificationText := fmt.Sprintf(
			"reviews.notification.new_review\nreviews.notification.rating: %d/5\nreviews.notification.listing: %s\n\n%s",
			review.Rating,
			listing.Title,
			review.Comment,
		)

		if err := h.services.Notification().SendNotification(
			c.Context(),
			listing.UserID,
			models.NotificationTypeNewReview,
			notificationText,
			listing.ID,
		); err != nil {
			log.Printf("Error sending notification: %v", err)
			// Не возвращаем ошибку, так как отзыв уже опубликован
		}
	}

	return utils.SuccessResponse(c, publishedReview)
}

// GetReviews returns filtered list of reviews
// @Summary Get reviews list
// @Description Returns paginated list of reviews with filters
// @Tags reviews
// @Accept json
// @Produce json
// @Param entity_type query string false "Entity type filter"
// @Param entity_id query int false "Entity ID filter"
// @Param user_id query int false "User ID filter"
// @Param min_rating query int false "Minimum rating filter"
// @Param max_rating query int false "Maximum rating filter"
// @Param status query string false "Status filter" default(published)
// @Param sort_by query string false "Sort by field" default(date)
// @Param sort_order query string false "Sort order" default(desc)
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} utils.SuccessResponseSwag{data=ReviewsListResponse} "List of reviews"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/reviews [get]
func (h *ReviewHandler) GetReviews(c *fiber.Ctx) error {
	filter := models.ReviewsFilter{
		EntityType: c.Query("entity_type"),
		EntityID:   utils.StringToInt(c.Query("entity_id"), 0),
		UserID:     utils.StringToInt(c.Query("user_id"), 0),
		MinRating:  utils.StringToInt(c.Query("min_rating"), 0),
		MaxRating:  utils.StringToInt(c.Query("max_rating"), 5),
		Status:     c.Query("status", "published"),
		SortBy:     c.Query("sort_by", "date"),
		SortOrder:  c.Query("sort_order", "desc"),
		Page:       utils.StringToInt(c.Query("page"), 1),
		Limit:      utils.StringToInt(c.Query("limit"), 20),
	}

	reviews, total, err := h.services.Review().GetReviews(c.Context(), filter)
	if err != nil {
		log.Printf("Error getting reviews: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "reviews.list.error.fetch_failed")
	}

	return utils.SuccessResponse(c, ReviewsListResponse{
		Success: true,
		Data:    reviews,
		Meta: ReviewsMeta{
			Total: int(total),
			Page:  filter.Page,
			Limit: filter.Limit,
		},
	})
}

// VoteForReview adds a vote to a review
// @Summary Vote for a review
// @Description Adds a helpful or unhelpful vote to a review
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "Review ID"
// @Param vote body VoteRequest true "Vote data"
// @Success 200 {object} utils.SuccessResponseSwag{data=ReviewMessageResponse} "Vote recorded successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 404 {object} utils.ErrorResponseSwag "Review not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/reviews/{id}/vote [post]
func (h *ReviewHandler) VoteForReview(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "reviews.error.unauthorized")
	}
	reviewID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.error.invalid_id")
	}

	var request struct {
		VoteType string `json:"vote_type"`
	}

	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.error.invalid_request")
	}

	// Получаем информацию об отзыве
	review, err := h.services.Review().GetReviewByID(c.Context(), reviewID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "reviews.error.not_found")
	}

	// Получаем информацию об объявлении
	listing, err := // TODO: Marketplace service disabled
	return nil, fmt.Errorf("marketplace service temporarily disabled")
	// GetListingByID(c.Context(), review.EntityID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "reviews.error.listing_not_found")
	}

	err = h.services.Review().VoteForReview(c.Context(), userID, reviewID, request.VoteType)
	if err != nil {
		log.Printf("Error voting for review: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "reviews.vote.error.failed")
	}

	notificationText := fmt.Sprintf(
		"reviews.notification.vote: %s\nreviews.notification.listing: %s",
		request.VoteType,
		listing.Title,
	)
	if err := h.services.Notification().SendNotification(
		c.Context(),
		review.UserID,
		models.NotificationTypeReviewVote,
		notificationText,
		listing.ID,
	); err != nil {
		log.Printf("Error sending notification: %v", err)
	}

	return utils.SuccessResponse(c, ReviewMessageResponse{
		Success: true,
		Message: "reviews.vote.success.recorded",
	})
}

// AddResponse adds a response to a review
// @Summary Add response to review
// @Description Adds a response from the listing owner to a review
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "Review ID"
// @Param response body ResponseRequest true "Response data"
// @Success 200 {object} utils.SuccessResponseSwag{data=ReviewMessageResponse} "Response added successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 404 {object} utils.ErrorResponseSwag "Review not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/reviews/{id}/response [post]
func (h *ReviewHandler) AddResponse(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "reviews.error.unauthorized")
	}
	reviewID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.error.invalid_id")
	}

	var request struct {
		Response string `json:"response"`
	}

	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.error.invalid_request")
	}

	// Санитизация ответа от XSS
	request.Response = utils.SanitizeText(request.Response)

	// Получаем информацию об отзыве
	review, err := h.services.Review().GetReviewByID(c.Context(), reviewID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "reviews.error.not_found")
	}

	// Получаем информацию об объявлении
	listing, err := // TODO: Marketplace service disabled
	return nil, fmt.Errorf("marketplace service temporarily disabled")
	// GetListingByID(c.Context(), review.EntityID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "reviews.error.listing_not_found")
	}

	err = h.services.Review().AddResponse(c.Context(), userID, reviewID, request.Response)
	if err != nil {
		log.Printf("Error adding response: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "reviews.response.error.add_failed")
	}

	notificationText := fmt.Sprintf(
		"reviews.notification.response\nreviews.notification.listing: %s\n\n%s",
		listing.Title,
		request.Response,
	)
	if err := h.services.Notification().SendNotification(
		c.Context(),
		review.UserID,
		models.NotificationTypeReviewResponse,
		notificationText,
		listing.ID,
	); err != nil {
		log.Printf("Error sending notification: %v", err)
	}

	return utils.SuccessResponse(c, ReviewMessageResponse{
		Success: true,
		Message: "reviews.response.success.added",
	})
}

// GetReviewByID returns a review by ID
// @Summary Get review by ID
// @Description Returns a single review by its ID
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "Review ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.Review} "Review details"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid review ID"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/reviews/{id} [get]
func (h *ReviewHandler) GetReviewByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.error.invalid_id")
	}

	review, err := h.services.Review().GetReviewByID(c.Context(), id)
	if err != nil {
		log.Printf("Error getting review: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "reviews.get.error.failed")
	}

	return utils.SuccessResponse(c, review)
}

// UpdateReview updates an existing review
// @Summary Update review
// @Description Updates an existing review by its author
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "Review ID"
// @Param review body models.Review true "Updated review data"
// @Success 200 {object} utils.SuccessResponseSwag{data=ReviewMessageResponse} "Review updated successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/reviews/{id} [put]
func (h *ReviewHandler) UpdateReview(c *fiber.Ctx) error {
	userId, err := getUserID(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "reviews.error.unauthorized")
	}
	reviewId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.error.invalid_id")
	}

	var review models.Review
	if err := c.BodyParser(&review); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.error.invalid_request")
	}

	// Санитизация полей от XSS
	review.Comment = utils.SanitizeText(review.Comment)

	err = h.services.Review().UpdateReview(c.Context(), userId, reviewId, &review)
	if err != nil {
		log.Printf("Error updating review: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "reviews.update.error.failed")
	}

	return utils.SuccessResponse(c, ReviewMessageResponse{
		Success: true,
		Message: "reviews.update.success.updated",
	})
}

// GetStats returns review statistics
// @Summary Get review statistics
// @Description Returns review statistics for a specific entity
// @Tags reviews
// @Accept json
// @Produce json
// @Param entity_type query string false "Entity type"
// @Param entity_id query int false "Entity ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.ReviewStats} "Review statistics"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/reviews/stats [get]
func (h *ReviewHandler) GetStats(c *fiber.Ctx) error {
	entityType := c.Query("entity_type")
	entityID := utils.StringToInt(c.Query("entity_id"), 0)

	stats, err := h.services.Review().GetReviewStats(c.Context(), entityType, entityID)
	if err != nil {
		log.Printf("Error getting review stats: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "reviews.stats.error.fetch_failed")
	}

	return utils.SuccessResponse(c, stats)
}

// UploadPhotos uploads photos for a review
// @Summary Upload review photos
// @Description Uploads photos for an existing review
// @Tags reviews
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Review ID"
// @Param photos formData file true "Photos to upload (max 5, max 5MB each, formats: jpg/png/webp)"
// @Success 200 {object} utils.SuccessResponseSwag{data=PhotosResponse} "Photos uploaded successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/reviews/{id}/photos [post]
func (h *ReviewHandler) UploadPhotos(c *fiber.Ctx) error {
	userId, err := getUserID(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "reviews.error.unauthorized")
	}
	reviewId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.error.invalid_id")
	}

	// Получаем существующий отзыв для проверки авторства и количества фото
	review, err := h.services.Review().GetReviewByID(c.Context(), reviewId)
	if err != nil || review == nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "reviews.error.not_found")
	}

	// Проверяем, что пользователь является автором отзыва
	if review.UserID != userId {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "reviews.error.not_author")
	}

	// Получаем загруженные файлы
	form, err := c.MultipartForm()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.photos.error.parse_form")
	}

	files := form.File["photos"]
	if len(files) == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.photos.error.no_files")
	}

	// Проверяем общее количество фото (существующие + новые)
	existingPhotosCount := len(review.Photos)
	totalPhotos := existingPhotosCount + len(files)
	if totalPhotos > 5 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.photos.error.total_max_files")
	}

	// Разрешенные форматы
	allowedFormats := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/webp": true,
	}

	photoUrls := make([]string, 0)
	for _, file := range files {
		// Проверка размера файла (максимум 5MB)
		if file.Size > 5*1024*1024 {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.photos.error.file_too_large")
		}

		// Проверка типа файла
		contentType := file.Header.Get("Content-Type")
		if !allowedFormats[contentType] {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.photos.error.invalid_type")
		}

		// Генерируем уникальное имя файла
		filename := fmt.Sprintf("review_%d_%d_%s", reviewId, time.Now().UnixNano(), file.Filename)
		objectKey := "reviews/" + filename

		// Открываем файл для чтения
		src, err := file.Open()
		if err != nil {
			log.Printf("Error opening file: %v", err)
			continue
		}
		defer func() {
			if err := src.Close(); err != nil {
				log.Printf("Error closing file: %v", err)
			}
		}()

		// Загружаем файл в MinIO и получаем полный URL
		imageURL, err := h.services.Storage().FileStorage().UploadFile(c.Context(), objectKey, src, file.Size, contentType)
		if err != nil {
			log.Printf("Error uploading file to MinIO: %v", err)
			continue
		}

		photoUrls = append(photoUrls, imageURL)
	}

	// Обновляем отзыв с новыми фотографиями
	err = h.services.Review().UpdateReviewPhotos(c.Context(), reviewId, photoUrls)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "reviews.photos.error.update_failed")
	}

	return utils.SuccessResponse(c, PhotosResponse{
		Success: true,
		Message: "reviews.photos.success.uploaded",
		Photos:  photoUrls,
	})
}

// UploadPhotosForNewReview uploads photos for a new review being created
// @Summary Upload photos for new review
// @Description Uploads photos that will be attached to a new review during creation
// @Tags reviews
// @Accept multipart/form-data
// @Produce json
// @Param photos formData file true "Photos to upload (max 5, max 5MB each, formats: jpg/png/webp)"
// @Success 200 {object} utils.SuccessResponseSwag{data=PhotosResponse} "Photos uploaded successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/reviews/upload-photos [post]
func (h *ReviewHandler) UploadPhotosForNewReview(c *fiber.Ctx) error {
	userId, err := getUserID(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "reviews.error.unauthorized")
	}

	// Получаем загруженные файлы
	form, err := c.MultipartForm()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.photos.error.parse_form")
	}

	files := form.File["photos"]
	if len(files) == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.photos.error.no_files")
	}

	// Максимум 5 фото
	if len(files) > 5 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.photos.error.max_files")
	}

	// Разрешенные форматы
	allowedFormats := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/webp": true,
	}

	photoUrls := make([]string, 0)
	for _, file := range files {
		// Проверка размера файла (максимум 5MB)
		if file.Size > 5*1024*1024 {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.photos.error.file_too_large")
		}

		// Проверка типа файла
		contentType := file.Header.Get("Content-Type")
		if !allowedFormats[contentType] {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.photos.error.invalid_type")
		}

		// Генерируем уникальное имя файла для временного хранения в review-photos
		filename := fmt.Sprintf("temp_%d_%d_%s", userId, time.Now().UnixNano(), file.Filename)
		tempObjectKey := "temp/" + filename

		// Открываем файл для чтения
		src, err := file.Open()
		if err != nil {
			log.Printf("Error opening file: %v", err)
			continue
		}
		defer func() {
			if err := src.Close(); err != nil {
				log.Printf("Error closing file: %v", err)
			}
		}()

		// Пока используем основной FileStorage, позже добавим review-photos wrapper
		// TODO: Создать отдельный бакет review-photos
		imageURL, err := h.services.Storage().FileStorage().UploadFile(c.Context(), tempObjectKey, src, file.Size, contentType)
		if err != nil {
			log.Printf("Error uploading temp file to MinIO: %v", err)
			continue
		}

		photoUrls = append(photoUrls, imageURL)
	}

	if len(photoUrls) == 0 {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "reviews.photos.error.save_failed")
	}

	return utils.SuccessResponse(c, PhotosResponse{
		Success: true,
		Message: "reviews.photos.success.uploaded_temp",
		Photos:  photoUrls,
	})
}

// DeleteReview deletes a review
// @Summary Delete review
// @Description Deletes a review by its author
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "Review ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=ReviewMessageResponse} "Review deleted successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid review ID"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/reviews/{id} [delete]
func (h *ReviewHandler) DeleteReview(c *fiber.Ctx) error {
	userId, err := getUserID(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "reviews.error.unauthorized")
	}
	reviewId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.error.invalid_id")
	}

	err = h.services.Review().DeleteReview(c.Context(), userId, reviewId)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "reviews.delete.error.failed")
	}

	return utils.SuccessResponse(c, ReviewMessageResponse{
		Success: true,
		Message: "reviews.delete.success.deleted",
	})
}

// GetEntityRating returns entity rating
// @Summary Get entity rating
// @Description Returns average rating for a specific entity
// @Tags reviews
// @Accept json
// @Produce json
// @Param type path string true "Entity type"
// @Param id path int true "Entity ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=RatingResponse} "Entity rating"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid entity ID"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/reviews/rating/{type}/{id} [get]
func (h *ReviewHandler) GetEntityRating(c *fiber.Ctx) error {
	entityType := c.Params("type")
	entityId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.error.invalid_entity_id")
	}

	rating, err := h.services.Review().GetEntityRating(c.Context(), entityType, entityId)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "reviews.rating.error.fetch_failed")
	}

	return utils.SuccessResponse(c, RatingResponse{
		Success: true,
		Rating:  rating,
	})
}

// GetEntityStats returns entity review statistics
// @Summary Get entity review statistics
// @Description Returns detailed review statistics for a specific entity
// @Tags reviews
// @Accept json
// @Produce json
// @Param type path string true "Entity type"
// @Param id path int true "Entity ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.ReviewStats} "Entity review statistics"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid entity ID"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/reviews/stats/{type}/{id} [get]
func (h *ReviewHandler) GetEntityStats(c *fiber.Ctx) error {
	entityType := c.Params("type")
	entityId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.error.invalid_entity_id")
	}

	log.Printf("Getting stats for %s ID=%d", entityType, entityId)

	stats, err := h.services.Review().GetReviewStats(c.Context(), entityType, entityId)
	if err != nil {
		log.Printf("Error getting stats for %s ID=%d: %v", entityType, entityId, err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "reviews.stats.error.fetch_failed")
	}

	log.Printf("Got stats for %s ID=%d: %+v", entityType, entityId, stats)

	return utils.SuccessResponse(c, stats)
}

// GetUserReviews returns all reviews for a user
// @Summary Get user reviews
// @Description Returns all reviews received by a specific user
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=ReviewsListResponse} "User reviews"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid user ID"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/reviews/user/{id} [get]
func (h *ReviewHandler) GetUserReviews(c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.error.invalid_user_id")
	}

	reviews, err := h.services.Review().GetUserReviews(c.Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "reviews.user.error.fetch_failed")
	}

	return utils.SuccessResponse(c, ReviewsListResponse{
		Success: true,
		Data:    reviews,
		Meta: ReviewsMeta{
			Total: len(reviews),
			Page:  1,
			Limit: len(reviews),
		},
	})
}

// GetStorefrontReviews returns all reviews for a storefront
// @Summary Get storefront reviews
// @Description Returns all reviews for a specific storefront
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "Storefront ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=ReviewsListResponse} "Storefront reviews"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid storefront ID"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/reviews/storefront/{id} [get]
func (h *ReviewHandler) GetStorefrontReviews(c *fiber.Ctx) error {
	storefrontID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.error.invalid_storefront_id")
	}

	reviews, err := h.services.Review().GetStorefrontReviews(c.Context(), storefrontID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "reviews.storefront.error.fetch_failed")
	}

	return utils.SuccessResponse(c, ReviewsListResponse{
		Success: true,
		Data:    reviews,
		Meta: ReviewsMeta{
			Total: len(reviews),
			Page:  1,
			Limit: len(reviews),
		},
	})
}

// GetUserRatingSummary returns user rating summary
// @Summary Get user rating summary
// @Description Returns rating summary for a specific user
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.UserRatingSummary} "User rating summary"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid user ID"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/reviews/user/{id}/rating [get]
func (h *ReviewHandler) GetUserRatingSummary(c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.error.invalid_user_id")
	}

	summary, err := h.services.Review().GetUserRatingSummary(c.Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "reviews.user.error.rating_fetch_failed")
	}

	return utils.SuccessResponse(c, summary)
}

// GetStorefrontRatingSummary returns storefront rating summary
// @Summary Get storefront rating summary
// @Description Returns rating summary for a specific storefront
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "Storefront ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.StorefrontRatingSummary} "Storefront rating summary"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid storefront ID"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/reviews/storefront/{id}/rating [get]
func (h *ReviewHandler) GetStorefrontRatingSummary(c *fiber.Ctx) error {
	storefrontID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.error.invalid_storefront_id")
	}

	summary, err := h.services.Review().GetStorefrontRatingSummary(c.Context(), storefrontID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "reviews.storefront.error.rating_fetch_failed")
	}

	return utils.SuccessResponse(c, summary)
}

// GetUserAggregatedRating возвращает агрегированный рейтинг пользователя
// @Summary Get user aggregated rating
// @Description Returns aggregated rating for a user including breakdown by sources
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.AggregatedRating} "Aggregated rating"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid user ID"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/users/{id}/aggregated-rating [get]
func (h *ReviewHandler) GetUserAggregatedRating(c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.error.invalid_user_id")
	}

	rating, err := h.services.Review().GetUserAggregatedRating(c.Context(), userID)
	if err != nil {
		log.Printf("Error getting user aggregated rating: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "reviews.aggregated_rating.error.fetch_failed")
	}

	return utils.SuccessResponse(c, rating)
}

// GetStorefrontAggregatedRating возвращает агрегированный рейтинг магазина
// @Summary Get storefront aggregated rating
// @Description Returns aggregated rating for a storefront including breakdown by sources
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "Storefront ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.AggregatedRating} "Aggregated rating"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid storefront ID"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/b2c_stores/{id}/aggregated-rating [get]
func (h *ReviewHandler) GetStorefrontAggregatedRating(c *fiber.Ctx) error {
	storefrontID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.error.invalid_storefront_id")
	}

	rating, err := h.services.Review().GetStorefrontAggregatedRating(c.Context(), storefrontID)
	if err != nil {
		log.Printf("Error getting storefront aggregated rating: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "reviews.aggregated_rating.error.fetch_failed")
	}

	return utils.SuccessResponse(c, rating)
}

// CanReview проверяет может ли пользователь оставить отзыв
// @Summary Check if user can review entity
// @Description Checks if the current user can leave a review for the specified entity
// @Tags reviews
// @Accept json
// @Produce json
// @Param type path string true "Entity type (listing, user, storefront)"
// @Param id path int true "Entity ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.CanReviewResponse} "Permission check result"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid parameters"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/reviews/can-review/{type}/{id} [get]
func (h *ReviewHandler) CanReview(c *fiber.Ctx) error {
	userID, ok := authMiddleware.GetUserID(c)
	if !ok || userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "reviews.error.unauthorized")
	}

	entityType := c.Params("type")
	entityID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.error.invalid_entity_id")
	}

	// Валидация типа сущности
	validTypes := map[string]bool{"listing": true, "user": true, "storefront": true}
	if !validTypes[entityType] {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.error.invalid_entity_type")
	}

	response, err := h.services.Review().CanUserReviewEntity(c.Context(), userID, entityType, entityID)
	if err != nil {
		log.Printf("Error checking review permission: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "reviews.permission.error.check_failed")
	}

	return utils.SuccessResponse(c, response)
}

// CanReviewListing проверяет может ли пользователь оставить отзыв на listing
func (h *ReviewHandler) CanReviewListing(c *fiber.Ctx) error {
	userID, ok := authMiddleware.GetUserID(c)
	if !ok || userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "reviews.error.unauthorized")
	}

	entityID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.error.invalid_entity_id")
	}

	response, err := h.services.Review().CanUserReviewEntity(c.Context(), userID, "listing", entityID)
	if err != nil {
		log.Printf("Error checking review permission for listing: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "reviews.permission.error.check_failed")
	}

	return utils.SuccessResponse(c, response)
}

// CanReviewUser проверяет может ли пользователь оставить отзыв на user
func (h *ReviewHandler) CanReviewUser(c *fiber.Ctx) error {
	userID, ok := authMiddleware.GetUserID(c)
	if !ok || userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "reviews.error.unauthorized")
	}

	entityID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.error.invalid_entity_id")
	}

	response, err := h.services.Review().CanUserReviewEntity(c.Context(), userID, "user", entityID)
	if err != nil {
		log.Printf("Error checking review permission for user: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "reviews.permission.error.check_failed")
	}

	return utils.SuccessResponse(c, response)
}

// CanReviewStorefront проверяет может ли пользователь оставить отзыв на storefront
func (h *ReviewHandler) CanReviewStorefront(c *fiber.Ctx) error {
	userID, ok := authMiddleware.GetUserID(c)
	if !ok || userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "reviews.error.unauthorized")
	}

	entityID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.error.invalid_entity_id")
	}

	response, err := h.services.Review().CanUserReviewEntity(c.Context(), userID, "storefront", entityID)
	if err != nil {
		log.Printf("Error checking review permission for storefront: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "reviews.permission.error.check_failed")
	}

	return utils.SuccessResponse(c, response)
}

// ConfirmReview подтверждает отзыв продавцом
// @Summary Confirm review
// @Description Allows seller to confirm or dispute a review
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "Review ID"
// @Param request body models.CreateReviewConfirmationRequest true "Confirmation request"
// @Success 200 {object} utils.SuccessResponseSwag{data=ReviewMessageResponse} "Review confirmed successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 403 {object} utils.ErrorResponseSwag "Not authorized"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/reviews/{id}/confirm [post]
func (h *ReviewHandler) ConfirmReview(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "reviews.error.unauthorized")
	}
	reviewID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.error.invalid_id")
	}

	var req models.CreateReviewConfirmationRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.error.invalid_request")
	}

	err = h.services.Review().ConfirmReview(c.Context(), userID, reviewID, &req)
	if err != nil {
		log.Printf("Error confirming review: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "reviews.confirm.error.failed")
	}

	return utils.SuccessResponse(c, ReviewMessageResponse{
		Success: true,
		Message: "reviews.confirm.success",
	})
}

// DisputeReview создает спор по отзыву
// @Summary Dispute review
// @Description Creates a dispute for a review
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "Review ID"
// @Param request body models.CreateReviewDisputeRequest true "Dispute request"
// @Success 200 {object} utils.SuccessResponseSwag{data=ReviewMessageResponse} "Dispute created successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/reviews/{id}/dispute [post]
func (h *ReviewHandler) DisputeReview(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "reviews.error.unauthorized")
	}
	reviewID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.error.invalid_id")
	}

	var req models.CreateReviewDisputeRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "reviews.error.invalid_request")
	}

	err = h.services.Review().DisputeReview(c.Context(), userID, reviewID, &req)
	if err != nil {
		log.Printf("Error creating dispute: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "reviews.dispute.error.failed")
	}

	return utils.SuccessResponse(c, ReviewMessageResponse{
		Success: true,
		Message: "reviews.dispute.success.created",
	})
}
