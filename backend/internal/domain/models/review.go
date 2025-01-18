// backend/internal/domain/models/review.go

package models

import "time"

type Review struct {
    ID                int       `json:"id"`
    UserID            int       `json:"user_id"`
    EntityType        string    `json:"entity_type"`
    EntityID          int       `json:"entity_id"`
    Rating            int       `json:"rating"`
    Comment           string    `json:"comment,omitempty"`
    Pros              string    `json:"pros,omitempty"`
    Cons              string    `json:"cons,omitempty"`
    Photos            []string  `json:"photos,omitempty"`
    LikesCount        int       `json:"likes_count"`
    IsVerifiedPurchase bool     `json:"is_verified_purchase"`
    Status            string    `json:"status"`
    CreatedAt         time.Time `json:"created_at"`
    UpdatedAt         time.Time `json:"updated_at"`
	HelpfulVotes      int       `json:"helpful_votes"`       
	NotHelpfulVotes   int       `json:"not_helpful_votes"`   
    // поля для мультиязычности
    OriginalLanguage string                            `json:"original_language"`
    Translations    map[string]map[string]string      `json:"translations,omitempty"`
    
    // Дополнительные поля для отображения
    User              *User     `json:"user,omitempty"`
    Responses         []ReviewResponse `json:"responses,omitempty"`
    VotesCount        struct {
        Helpful     int `json:"helpful"`
        NotHelpful  int `json:"not_helpful"`
    } `json:"votes_count,omitempty"`
    CurrentUserVote   string    `json:"current_user_vote,omitempty"`
}

type ReviewResponse struct {
    ID        int       `json:"id"`
    ReviewID  int       `json:"review_id"`
    UserID    int       `json:"user_id"`
    Response  string    `json:"response"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    User      *User     `json:"user,omitempty"`
}

type ReviewVote struct {
    ReviewID  int       `json:"review_id"`
    UserID    int       `json:"user_id"`
    VoteType  string    `json:"vote_type"`
    CreatedAt time.Time `json:"created_at"`
}

// Запрос на создание отзыва
type CreateReviewRequest struct {
    EntityType  string   `json:"entity_type" validate:"required,oneof=listing room car"`
    EntityID    int      `json:"entity_id" validate:"required"`
    Rating      int      `json:"rating" validate:"required,min=1,max=5"`
    Comment     string   `json:"comment"`
    Pros        string   `json:"pros"`
    Cons        string   `json:"cons"`
    Photos      []string `json:"photos"`
}

// Фильтры для получения отзывов
type ReviewsFilter struct {
    EntityType string `query:"entity_type"`
    EntityID   int    `query:"entity_id"`
    UserID     int    `query:"user_id"`
    MinRating  int    `query:"min_rating"`
    MaxRating  int    `query:"max_rating"`
    Status     string `query:"status"`
    SortBy     string `query:"sort_by"` // rating, date, likes
    SortOrder  string `query:"sort_order"` // asc, desc
    Page       int    `query:"page"`
    Limit      int    `query:"limit"`
}
type ReviewStats struct {
    TotalReviews      int     `json:"total_reviews"`
    AverageRating     float64 `json:"average_rating"`
    VerifiedReviews   int     `json:"verified_reviews"`
    RatingDistribution map[int]int `json:"rating_distribution"` // Распределение оценок: {1: 10, 2: 20, ...}
    PhotoReviews      int     `json:"photo_reviews"`          // Количество отзывов с фото
}