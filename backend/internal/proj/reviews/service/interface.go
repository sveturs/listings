package service

import (
    "context"
    "backend/internal/domain/models"
)

type ReviewServiceInterface interface {
    CreateReview(ctx context.Context, userId int, review *models.CreateReviewRequest) (*models.Review, error)
    GetReviews(ctx context.Context, filter models.ReviewsFilter) ([]models.Review, int64, error)
    GetReviewByID(ctx context.Context, id int) (*models.Review, error)
    UpdateReview(ctx context.Context, userId int, reviewId int, review *models.Review) error
    DeleteReview(ctx context.Context, userId int, reviewId int) error
    VoteForReview(ctx context.Context, userId int, reviewId int, voteType string) error
    AddResponse(ctx context.Context, userId int, reviewId int, response string) error
    GetEntityRating(ctx context.Context, entityType string, entityId int) (float64, error)
    GetReviewStats(ctx context.Context, entityType string, entityId int) (*models.ReviewStats, error)
    UpdateReviewPhotos(ctx context.Context, reviewId int, photoUrls []string) error

    GetUserReviews(ctx context.Context, userID int) ([]models.Review, error)
    GetStorefrontReviews(ctx context.Context, storefrontID int) ([]models.Review, error)
    GetUserRatingSummary(ctx context.Context, userID int) (*models.UserRatingSummary, error)
    GetStorefrontRatingSummary(ctx context.Context, storefrontID int) (*models.StorefrontRatingSummary, error)

}