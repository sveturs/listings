// internal/services/review.go

package service

import (
    "context"
    "backend/internal/domain/models"
    "backend/internal/storage"
	"fmt"
    "log"
)

type ReviewService struct {
    storage storage.Storage
}

func NewReviewService(storage storage.Storage) ReviewServiceInterface { //  ReviewService -> ReviewServiceInterface
    if storage == nil {
        log.Fatal("storage cannot be nil")
    }
    return &ReviewService{
        storage: storage,
    }
}

func (s *ReviewService) CreateReview(ctx context.Context, userId int, req *models.CreateReviewRequest) (*models.Review, error) {
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
        OriginalLanguage: req.OriginalLanguage, 
    }
    
    // Проверяем, является ли покупка верифицированной
    review.IsVerifiedPurchase = s.checkVerifiedPurchase(ctx, userId, req.EntityType, req.EntityID)
    
    // Сохраняем отзыв
    createdReview, err := s.storage.CreateReview(ctx, review)
    if err != nil {
        return nil, err
    }
    
    return createdReview, nil
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
func (s *ReviewService) GetReviewStats(ctx context.Context, entityType string, entityId int) (*models.ReviewStats, error) {
    stats := &models.ReviewStats{
        RatingDistribution: make(map[int]int),
    }

    // Получаем общую статистику
    err := s.storage.QueryRow(ctx, `
        SELECT 
            COUNT(*) as total,
            COALESCE(AVG(rating), 0) as avg_rating,
            COUNT(*) FILTER (WHERE is_verified_purchase) as verified,
            COUNT(*) FILTER (WHERE array_length(photos, 1) > 0) as with_photos
        FROM reviews
        WHERE entity_type = $1 
        AND entity_id = $2
        AND status = 'published'
    `, entityType, entityId).Scan(
        &stats.TotalReviews,
        &stats.AverageRating,
        &stats.VerifiedReviews,
        &stats.PhotoReviews,
    )
    if err != nil {
        return nil, err
    }

    // Получаем распределение оценок
    rows, err := s.storage.Query(ctx, `
        SELECT rating, COUNT(*)
        FROM reviews
        WHERE entity_type = $1 
        AND entity_id = $2
        AND status = 'published'
        GROUP BY rating
    `, entityType, entityId)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var rating, count int
        if err := rows.Scan(&rating, &count); err != nil {
            return nil, err
        }
        stats.RatingDistribution[rating] = count
    }

    return stats, nil
}

func (s *ReviewService) UpdateReviewPhotos(ctx context.Context, reviewId int, photoUrls []string) error {
    // Получаем текущий отзыв
    review, err := s.storage.GetReviewByID(ctx, reviewId)
    if err != nil {
        return err
    }

    // Обновляем массив фотографий
    review.Photos = photoUrls

    // Сохраняем изменения
    return s.storage.UpdateReview(ctx, review)
}
