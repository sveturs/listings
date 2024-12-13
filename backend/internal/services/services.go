//backend/internal/services/services.go
package services

import (
    "backend/internal/config"
    "backend/internal/storage"
    "context"
    "backend/internal/domain/models"

)

type ServicesInterface interface {
    Auth() AuthServiceInterface
    Room() RoomServiceInterface
    Booking() BookingServiceInterface
    User() UserServiceInterface
    Car() CarServiceInterface
    Config() *config.Config
    Marketplace() MarketplaceServiceInterface
    Review() ReviewServiceInterface  
}

type Services struct {
    auth    AuthServiceInterface
    room    RoomServiceInterface
    booking BookingServiceInterface
    user    UserServiceInterface
    car     CarServiceInterface
    marketplace MarketplaceServiceInterface 
    review      ReviewServiceInterface  
    config  *config.Config
}

func NewServices(storage storage.Storage, cfg *config.Config) *Services {
    return &Services{
        auth:    NewAuthService(cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.GoogleRedirectURL, storage),
        room:    NewRoomService(storage),
        booking: NewBookingService(storage),
        user:    NewUserService(storage),
        car:     NewCarService(storage),
        marketplace: NewMarketplaceService(storage),
        review:      NewReviewService(storage),
        config:  cfg,
    }
}
type ReviewServiceInterface interface {
    CreateReview(ctx context.Context, userId int, review *models.CreateReviewRequest) error
    GetReviews(ctx context.Context, filter models.ReviewsFilter) ([]models.Review, int64, error)
    GetReviewByID(ctx context.Context, id int) (*models.Review, error)
    UpdateReview(ctx context.Context, userId int, reviewId int, review *models.Review) error
    DeleteReview(ctx context.Context, userId int, reviewId int) error
    VoteForReview(ctx context.Context, userId int, reviewId int, voteType string) error
    AddResponse(ctx context.Context, userId int, reviewId int, response string) error
    GetEntityRating(ctx context.Context, entityType string, entityId int) (float64, error)
    GetReviewStats(ctx context.Context, entityType string, entityId int) (*models.ReviewStats, error)
    UpdateReviewPhotos(ctx context.Context, reviewId int, photoUrls []string) error

}
func (s *Services) Auth() AuthServiceInterface { 
    return s.auth
}
func (s *Services) Car() CarServiceInterface {
    return s.car
}
func (s *Services) Room() RoomServiceInterface { 
    return s.room
}

func (s *Services) Booking() BookingServiceInterface { 
    return s.booking
}

func (s *Services) User() UserServiceInterface { 
    return s.user
}

func (s *Services) Config() *config.Config { 
    return s.config
}
func (s *Services) Marketplace() MarketplaceServiceInterface {
    return s.marketplace
}
func (s *Services) Review() ReviewServiceInterface {
    return s.review
}
// Проверяем, что Services реализует ServicesInterface
var _ ServicesInterface = (*Services)(nil)