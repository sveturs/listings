// internal/services/review.go

package services

import (
    "context"
    "backend/internal/domain/models"
    "backend/internal/storage"
	"fmt"
)

type ReviewService struct {
    storage storage.Storage
}

func NewReviewService(storage storage.Storage) ReviewServiceInterface {
    return &ReviewService{
        storage: storage,
    }
}

func (s *ReviewService) CreateReview(ctx context.Context, userId int, req *models.CreateReviewRequest) error {
    review := &models.Review{
        UserID:     userId,
        EntityType: req.EntityType,
        EntityID:   req.EntityID,
        Rating:     req.Rating,
        Comment:    req.Comment,
        Pros:       req.Pros,
        Cons:       req.Cons,
        Photos:     req.Photos,
        Status:     "published",
    }
    
    // Проверяем, является ли покупка верифицированной
    review.IsVerifiedPurchase = s.checkVerifiedPurchase(ctx, userId, req.EntityType, req.EntityID)
    
    return s.storage.CreateReview(ctx, review)
}

func (s *ReviewService) GetReviews(ctx context.Context, filter models.ReviewsFilter) ([]models.Review, int64, error) {
    return s.storage.GetReviews(ctx, filter)
}

func (s *ReviewService) GetReviewByID(ctx context.Context, id int) (*models.Review, error) {
    return s.storage.GetReviewByID(ctx, id)
}

func (s *ReviewService) UpdateReview(ctx context.Context, userId int, reviewId int, review *models.Review) error {
    // Проверяем, принадлежит ли отзыв пользователю
    existingReview, err := s.storage.GetReviewByID(ctx, reviewId)
    if err != nil {
        return err
    }
    
    if existingReview.UserID != userId {
        return fmt.Errorf("unauthorized to update this review")
    }
    
    review.ID = reviewId
    return s.storage.UpdateReview(ctx, review)
}

func (s *ReviewService) DeleteReview(ctx context.Context, userId int, reviewId int) error {
    // Проверяем, принадлежит ли отзыв пользователю
    review, err := s.storage.GetReviewByID(ctx, reviewId)
    if err != nil {
        return err
    }
    
    if review.UserID != userId {
        return fmt.Errorf("unauthorized to delete this review")
    }
    
    return s.storage.DeleteReview(ctx, reviewId)
}

func (s *ReviewService) VoteForReview(ctx context.Context, userId int, reviewId int, voteType string) error {
    vote := &models.ReviewVote{
        ReviewID: reviewId,
        UserID:   userId,
        VoteType: voteType,
    }
    return s.storage.AddReviewVote(ctx, vote)
}

func (s *ReviewService) AddResponse(ctx context.Context, userId int, reviewId int, responseText string) error {
    response := &models.ReviewResponse{
        ReviewID: reviewId,
        UserID:   userId,
        Response: responseText,
    }
    return s.storage.AddReviewResponse(ctx, response)
}

func (s *ReviewService) GetEntityRating(ctx context.Context, entityType string, entityId int) (float64, error) {
    return s.storage.GetEntityRating(ctx, entityType, entityId)
}

// checkVerifiedPurchase проверяет, совершал ли пользователь покупку
func (s *ReviewService) checkVerifiedPurchase(ctx context.Context, userId int, entityType string, entityId int) bool {
    // В зависимости от типа сущности проверяем наличие покупки/бронирования
    switch entityType {
    case "listing":
        // Проверяем покупки в маркетплейсе
        return true // TODO: реализовать проверку
    case "room":
        // Проверяем бронирования комнат
        return true // TODO: реализовать проверку
    case "car":
        // Проверяем аренду автомобилей
        return true // TODO: реализовать проверку
    default:
        return false
    }
}