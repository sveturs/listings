// internal/handlers/reviews.go

package handler

import (
	"backend/internal/domain/models"
 	"backend/pkg/utils"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
    "backend/internal/proj/reviews/service"
    globalService "backend/internal/proj/global/service"


	"github.com/gofiber/fiber/v2"
)

type ReviewHandler struct {
    services        globalService.ServicesInterface
	reviewService   service.ReviewServiceInterface
}

func NewReviewHandler(services globalService.ServicesInterface) *ReviewHandler {
	return &ReviewHandler{
		services:      services,
		reviewService: services.Review(),
	}
}

// internal/handlers/reviews.go

func (h *ReviewHandler) CreateReview(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(int)
	var request models.CreateReviewRequest

	if err := c.BodyParser(&request); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	log.Printf("Creating review: %+v", request)
	review, err := h.reviewService.CreateReview(c.Context(), userId, &request)
	if err != nil {
		log.Printf("Error creating review: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error creating review")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Review created successfully",
		"id":      review.ID,
	})
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
	userId := c.Locals("user_id").(int)
	reviewId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid review ID")
	}

	var request struct {
		VoteType string `json:"vote_type"`
	}

	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	err = h.services.Review().VoteForReview(c.Context(), userId, reviewId, request.VoteType)
	if err != nil {
		log.Printf("Error voting for review: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error voting for review")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Vote recorded successfully",
	})
}

func (h *ReviewHandler) AddResponse(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(int)
	reviewId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid review ID")
	}

	var request struct {
		Response string `json:"response"`
	}

	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	err = h.services.Review().AddResponse(c.Context(), userId, reviewId, request.Response)
	if err != nil {
		log.Printf("Error adding response: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error adding response")
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
