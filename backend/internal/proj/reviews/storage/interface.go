package storage

import (
    "context"
    "backend/internal/domain/models"
)

type ReviewRepository interface {
    CreateReview(ctx context.Context, review *models.Review) (*models.Review, error)
    GetReviews(ctx context.Context, filter models.ReviewsFilter) ([]models.Review, int64, error)
    GetReviewByID(ctx context.Context, id int) (*models.Review, error)
    UpdateReview(ctx context.Context, review *models.Review) error
    DeleteReview(ctx context.Context, id int) error
    AddReviewResponse(ctx context.Context, response *models.ReviewResponse) error
    AddReviewVote(ctx context.Context, vote *models.ReviewVote) error 
    GetReviewVotes(ctx context.Context, reviewId int) (helpful int, notHelpful int, err error)
    GetUserReviewVote(ctx context.Context, userId int, reviewId int) (string, error)
    GetEntityRating(ctx context.Context, entityType string, entityId int) (float64, error)
    UpdateReviewVotes(ctx context.Context, reviewId int) error
    GetReviewStats(ctx context.Context, entityType string, entityId int) (*models.ReviewStats, error)
}