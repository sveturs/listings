// internal/handlers/reviews.go

package handlers

import (
    "backend/internal/domain/models"
    "backend/internal/services"
    "backend/pkg/utils"
    "github.com/gofiber/fiber/v2"
    "strconv"
    "log"
)

type ReviewHandler struct {
    services services.ServicesInterface
}

func NewReviewHandler(services services.ServicesInterface) *ReviewHandler {
    return &ReviewHandler{
        services: services,
    }
}

func (h *ReviewHandler) CreateReview(c *fiber.Ctx) error {
    userId := c.Locals("user_id").(int)
    var request models.CreateReviewRequest

    if err := c.BodyParser(&request); err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
    }

    err := h.services.Review().CreateReview(c.Context(), userId, &request)
    if err != nil {
        log.Printf("Error creating review: %v", err)
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error creating review")
    }

    return utils.SuccessResponse(c, fiber.Map{
        "message": "Review created successfully",
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