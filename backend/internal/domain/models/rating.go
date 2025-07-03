package models

import "time"

// AggregatedRating представляет агрегированный рейтинг сущности
type AggregatedRating struct {
	EntityType         string          `json:"entity_type"`
	EntityID           int             `json:"entity_id"`
	Average            float64         `json:"average"`
	TotalReviews       int             `json:"total_reviews"`
	Distribution       map[int]int     `json:"distribution"`
	Breakdown          RatingBreakdown `json:"breakdown"`
	VerifiedPercentage int             `json:"verified_percentage"`
	RecentTrend        string          `json:"recent_trend"` // up, down, stable
	LastReviewAt       *time.Time      `json:"last_review_at,omitempty"`
	RecentRating       *float64        `json:"recent_rating,omitempty"`
	RecentReviews      int             `json:"recent_reviews"`
}

// RatingBreakdown разбивка рейтинга по источникам
type RatingBreakdown struct {
	Direct      BreakdownItem `json:"direct"`      // Прямые отзывы
	Listings    BreakdownItem `json:"listings"`    // Через товары
	Storefronts BreakdownItem `json:"storefronts"` // Через магазины
}

// BreakdownItem элемент разбивки рейтинга
type BreakdownItem struct {
	Count   int     `json:"count"`
	Average float64 `json:"average"`
}

// UserAggregatedRating рейтинг пользователя из материализованного представления
type UserAggregatedRating struct {
	UserID            int        `json:"user_id"`
	TotalReviews      int        `json:"total_reviews"`
	AverageRating     float64    `json:"average_rating"`
	DirectReviews     int        `json:"direct_reviews"`
	ListingReviews    int        `json:"listing_reviews"`
	StorefrontReviews int        `json:"storefront_reviews"`
	VerifiedReviews   int        `json:"verified_reviews"`
	Rating1           int        `json:"rating_1"`
	Rating2           int        `json:"rating_2"`
	Rating3           int        `json:"rating_3"`
	Rating4           int        `json:"rating_4"`
	Rating5           int        `json:"rating_5"`
	RecentRating      *float64   `json:"recent_rating"`
	RecentReviews     int        `json:"recent_reviews"`
	LastReviewAt      *time.Time `json:"last_review_at"`
}

// StorefrontAggregatedRating рейтинг магазина из материализованного представления
type StorefrontAggregatedRating struct {
	StorefrontID    int        `json:"storefront_id"`
	TotalReviews    int        `json:"total_reviews"`
	AverageRating   float64    `json:"average_rating"`
	DirectReviews   int        `json:"direct_reviews"`
	ListingReviews  int        `json:"listing_reviews"`
	VerifiedReviews int        `json:"verified_reviews"`
	Rating1         int        `json:"rating_1"`
	Rating2         int        `json:"rating_2"`
	Rating3         int        `json:"rating_3"`
	Rating4         int        `json:"rating_4"`
	Rating5         int        `json:"rating_5"`
	RecentRating    *float64   `json:"recent_rating"`
	RecentReviews   int        `json:"recent_reviews"`
	LastReviewAt    *time.Time `json:"last_review_at"`
	OwnerID         int        `json:"owner_id"`
}

// CanReviewResponse ответ на проверку возможности оставить отзыв
type CanReviewResponse struct {
	CanReview         bool   `json:"can_review"`
	Reason            string `json:"reason,omitempty"`
	HasExistingReview bool   `json:"has_existing_review"`
	ExistingReviewID  *int   `json:"existing_review_id,omitempty"`
}

// ReviewConfirmation подтверждение отзыва продавцом
type ReviewConfirmation struct {
	ID                 int       `json:"id"`
	ReviewID           int       `json:"review_id"`
	ConfirmedBy        int       `json:"confirmed_by"`
	ConfirmationStatus string    `json:"confirmation_status"` // confirmed, disputed
	ConfirmedAt        time.Time `json:"confirmed_at"`
	Notes              *string   `json:"notes,omitempty"`
}

// ReviewDispute спор по отзыву
type ReviewDispute struct {
	ID                 int        `json:"id"`
	ReviewID           int        `json:"review_id"`
	DisputedBy         int        `json:"disputed_by"`
	DisputeReason      string     `json:"dispute_reason"`
	DisputeDescription string     `json:"dispute_description"`
	Status             string     `json:"status"`
	AdminID            *int       `json:"admin_id,omitempty"`
	AdminNotes         *string    `json:"admin_notes,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	ResolvedAt         *time.Time `json:"resolved_at,omitempty"`
}

// CreateReviewConfirmationRequest запрос на подтверждение отзыва
type CreateReviewConfirmationRequest struct {
	Status string  `json:"status" validate:"required,oneof=confirmed disputed"`
	Notes  *string `json:"notes,omitempty"`
}

// CreateReviewDisputeRequest запрос на создание спора
type CreateReviewDisputeRequest struct {
	Reason      string `json:"reason" validate:"required,oneof=not_a_customer false_information deal_cancelled spam other"`
	Description string `json:"description" validate:"required,min=10,max=1000"`
}
