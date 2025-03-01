// internal/handlers/reviews.go

package handler

import (
	"backend/internal/domain/models"
	globalService "backend/internal/proj/global/service"
	"backend/internal/proj/reviews/service"
	"backend/pkg/utils"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
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

// internal/handlers/reviews.go

func (h *ReviewHandler) CreateReview(c *fiber.Ctx) error {
	log.Printf("Starting CreateReview handler")
	userID := c.Locals("user_id").(int)

	var request models.CreateReviewRequest
	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid input")
	}

	// Определяем язык текста до создания отзыва
	if request.Comment != "" {
		detectedLang, _, err := h.services.Translation().DetectLanguage(c.Context(), request.Comment)
		if err != nil {
			log.Printf("Failed to detect language: %v", err)
			detectedLang = "en"
		}
		request.OriginalLanguage = detectedLang
	}

	// Проверяем входные данные
	if request.EntityID == 0 || request.Rating == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Required fields missing")
	}

	// Получаем информацию об объявлении
	listing, err := h.services.Marketplace().GetListingByID(c.Context(), request.EntityID)
	if err != nil {
		log.Printf("Failed to get listing %d: %v", request.EntityID, err)
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Listing not found")
	}

	if listing.StorefrontID != nil {
		// Если объявление относится к витрине, сохраняем эту информацию
		storefrontID := *listing.StorefrontID
		request.StorefrontID = &storefrontID
	}
	// Создаем отзыв через сервис
	createdReview, err := h.services.Review().CreateReview(c.Context(), userID, &request)
	if err != nil {
		log.Printf("Failed to create review: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error creating review")
	}

	log.Printf("Parsed request: %+v", request)

	// Проверяем созданный отзыв
	if createdReview == nil {
		log.Printf("Created review is nil")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error creating review")
	}

	log.Printf("Review created: %+v", createdReview)

	// Отправляем уведомление только если отзыв написан не владельцем объявления
	if listing.UserID != userID {
		notificationText := fmt.Sprintf(
			"Новый отзыв\nРейтинг: %d/5\nОбъявление: %s\n\n%s",
			request.Rating,
			listing.Title,
			request.Comment,
		)

		// Отправляем уведомление
		if err := h.services.Notification().SendNotification(
			c.Context(),
			listing.UserID,
			models.NotificationTypeNewReview,
			notificationText,
			listing.ID,
		); err != nil {
			log.Printf("Error sending notification: %v", err)
			// Не возвращаем ошибку, так как отзыв уже создан
		}
	}

	return utils.SuccessResponse(c, createdReview)
}

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
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error getting reviews")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"data": reviews,
		"meta": fiber.Map{
			"total": total,
			"page":  filter.Page,
			"limit": filter.Limit,
		},
	})
}

func (h *ReviewHandler) VoteForReview(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	reviewID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid review ID")
	}

	var request struct {
		VoteType string `json:"vote_type"`
	}

	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Получаем информацию об отзыве
	review, err := h.services.Review().GetReviewByID(c.Context(), reviewID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Review not found")
	}

	// Получаем информацию об объявлении
	listing, err := h.services.Marketplace().GetListingByID(c.Context(), review.EntityID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Listing not found")
	}

	err = h.services.Review().VoteForReview(c.Context(), userID, reviewID, request.VoteType)
	if err != nil {
		log.Printf("Error voting for review: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error voting for review")
	}

	notificationText := fmt.Sprintf(
		"Ваш отзыв оценили как %s\nОбъявление: %s",
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

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Vote recorded successfully",
	})
}
func (h *ReviewHandler) AddResponse(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	reviewID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid review ID")
	}

	var request struct {
		Response string `json:"response"`
	}

	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Получаем информацию об отзыве
	review, err := h.services.Review().GetReviewByID(c.Context(), reviewID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Review not found")
	}

	// Получаем информацию об объявлении
	listing, err := h.services.Marketplace().GetListingByID(c.Context(), review.EntityID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Listing not found")
	}

	err = h.services.Review().AddResponse(c.Context(), userID, reviewID, request.Response)
	if err != nil {
		log.Printf("Error adding response: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error adding response")
	}

	notificationText := fmt.Sprintf(
		"Получен ответ на ваш отзыв\nОбъявление: %s\n\n%s",
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

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Response added successfully",
	})
}

// internal/handlers/reviews.go

func (h *ReviewHandler) GetReviewByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid review ID")
	}

	review, err := h.services.Review().GetReviewByID(c.Context(), id)
	if err != nil {
		log.Printf("Error getting review: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error getting review")
	}

	return utils.SuccessResponse(c, review)
}

func (h *ReviewHandler) UpdateReview(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(int)
	reviewId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid review ID")
	}

	var review models.Review
	if err := c.BodyParser(&review); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	err = h.services.Review().UpdateReview(c.Context(), userId, reviewId, &review)
	if err != nil {
		log.Printf("Error updating review: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error updating review")
	}

	return utils.SuccessResponse(c, fiber.Map{"message": "Review updated successfully"})
}

func (h *ReviewHandler) GetStats(c *fiber.Ctx) error {
	entityType := c.Query("entity_type")
	entityID := utils.StringToInt(c.Query("entity_id"), 0)

	stats, err := h.services.Review().GetReviewStats(c.Context(), entityType, entityID)
	if err != nil {
		log.Printf("Error getting review stats: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error getting review stats")
	}

	return utils.SuccessResponse(c, stats)
}

func (h *ReviewHandler) UploadPhotos(c *fiber.Ctx) error {
	reviewId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid review ID")
	}

	// Получаем загруженные файлы
	form, err := c.MultipartForm()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Error parsing form")
	}

	files := form.File["photos"]
	if len(files) == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "No files uploaded")
	}

	if len(files) > 10 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Maximum 10 photos allowed")
	}

	photoUrls := make([]string, 0)
	for _, file := range files {
		// Проверка типа файла
		if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Only images are allowed")
		}

		// Генерируем уникальное имя файла
		filename := fmt.Sprintf("review_%d_%d_%s", reviewId, time.Now().UnixNano(), file.Filename)

		// Сохраняем файл
		if err := c.SaveFile(file, "./uploads/"+filename); err != nil {
			log.Printf("Error saving file: %v", err)
			continue
		}

		photoUrls = append(photoUrls, filename)
	}

	// Обновляем отзыв с новыми фотографиями
	err = h.services.Review().UpdateReviewPhotos(c.Context(), reviewId, photoUrls)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error updating review photos")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Photos uploaded successfully",
		"photos":  photoUrls,
	})
}
func (h *ReviewHandler) DeleteReview(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(int)
	reviewId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный ID отзыва")
	}

	err = h.services.Review().DeleteReview(c.Context(), userId, reviewId)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка при удалении отзыва")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Отзыв успешно удален",
	})
}

func (h *ReviewHandler) GetEntityRating(c *fiber.Ctx) error {
	entityType := c.Params("type")
	entityId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный ID сущности")
	}

	rating, err := h.services.Review().GetEntityRating(c.Context(), entityType, entityId)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка при получении рейтинга")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"rating": rating,
	})
}

func (h *ReviewHandler) GetEntityStats(c *fiber.Ctx) error {
	entityType := c.Params("type")
	entityId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный ID сущности")
	}

	stats, err := h.services.Review().GetReviewStats(c.Context(), entityType, entityId)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка при получении статистики")
	}

	return utils.SuccessResponse(c, stats)
}

// GetUserReviews получает все отзывы для пользователя
func (h *ReviewHandler) GetUserReviews(c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный ID пользователя")
	}

	reviews, err := h.services.Review().GetUserReviews(c.Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка при получении отзывов пользователя")
	}

	return utils.SuccessResponse(c, reviews)
}

// GetStorefrontReviews получает все отзывы для витрины
func (h *ReviewHandler) GetStorefrontReviews(c *fiber.Ctx) error {
	storefrontID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный ID витрины")
	}

	reviews, err := h.services.Review().GetStorefrontReviews(c.Context(), storefrontID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка при получении отзывов витрины")
	}

	return utils.SuccessResponse(c, reviews)
}

// GetUserRatingSummary получает сводные данные о рейтинге пользователя
func (h *ReviewHandler) GetUserRatingSummary(c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный ID пользователя")
	}

	summary, err := h.services.Review().GetUserRatingSummary(c.Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка при получении сводки рейтинга пользователя")
	}

	return utils.SuccessResponse(c, summary)
}

// GetStorefrontRatingSummary получает сводные данные о рейтинге витрины
func (h *ReviewHandler) GetStorefrontRatingSummary(c *fiber.Ctx) error {
	storefrontID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный ID витрины")
	}

	summary, err := h.services.Review().GetStorefrontRatingSummary(c.Context(), storefrontID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка при получении сводки рейтинга витрины")
	}

	return utils.SuccessResponse(c, summary)
}
