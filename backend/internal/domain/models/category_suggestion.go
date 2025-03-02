package models

// CategorySuggestion представляет предложение категории для поиска
type CategorySuggestion struct {
    ID           int    `json:"id"`
    Name         string `json:"name"`
    ListingCount int    `json:"listing_count"`
}