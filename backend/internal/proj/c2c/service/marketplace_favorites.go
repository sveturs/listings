// backend/internal/proj/c2c/service/marketplace_favorites.go
package service

import (
	"context"
	"log"
	"sort"

	"backend/internal/domain/models"
)

// GetUserFavorites получает список избранных объявлений пользователя
func (s *MarketplaceService) GetUserFavorites(ctx context.Context, userID int) ([]models.MarketplaceListing, error) {
	// Получаем обычные избранные объявления
	regularFavorites, err := s.storage.GetUserFavorites(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Получаем избранные товары витрин
	storefrontFavorites, err := s.storage.GetUserStorefrontFavorites(ctx, userID)
	if err != nil {
		// Если ошибка при получении товаров витрин, просто логируем и возвращаем обычные избранные
		log.Printf("Error getting storefront favorites: %v", err)
		return regularFavorites, nil
	}

	// Объединяем оба списка
	regularFavorites = append(regularFavorites, storefrontFavorites...)

	// Сортируем по времени добавления (новые сначала)
	// Предполагаем, что более новые имеют больший ID
	sort.Slice(regularFavorites, func(i, j int) bool {
		return regularFavorites[i].CreatedAt.After(regularFavorites[j].CreatedAt)
	})

	return regularFavorites, nil
}

// GetFavoritedUsers получает список пользователей, добавивших объявление в избранное
func (s *MarketplaceService) GetFavoritedUsers(ctx context.Context, listingID int) ([]int, error) {
	return s.storage.GetFavoritedUsers(ctx, listingID)
}

// AddToFavorites добавляет объявление в избранное
func (s *MarketplaceService) AddToFavorites(ctx context.Context, userID int, listingID int) error {
	return s.storage.AddToFavorites(ctx, userID, listingID)
}

// RemoveFromFavorites удаляет объявление из избранного
func (s *MarketplaceService) RemoveFromFavorites(ctx context.Context, userID int, listingID int) error {
	return s.storage.RemoveFromFavorites(ctx, userID, listingID)
}

// AddStorefrontToFavorites добавляет товар витрины в избранное
func (s *MarketplaceService) AddStorefrontToFavorites(ctx context.Context, userID int, productID int) error {
	return s.storage.AddStorefrontToFavorites(ctx, userID, productID)
}

// RemoveStorefrontFromFavorites удаляет товар витрины из избранного
func (s *MarketplaceService) RemoveStorefrontFromFavorites(ctx context.Context, userID int, productID int) error {
	return s.storage.RemoveStorefrontFromFavorites(ctx, userID, productID)
}
