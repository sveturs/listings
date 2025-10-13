// backend/internal/proj/c2c/storage/opensearch/repository_mappings.go
package opensearch

import (
	"context"
	"fmt"

	"backend/internal/domain/models"
	"backend/internal/logger"
	osClient "backend/internal/storage/opensearch"
)

func (r *Repository) PrepareIndex(ctx context.Context) error {
	exists, err := r.client.IndexExists(ctx, r.indexName)
	if err != nil {
		return fmt.Errorf("ошибка проверки индекса: %w", err)
	}

	logger.Info().Str("indexName", r.indexName).Bool("exists", exists).Msg("Проверка индекса")

	if !exists {
		logger.Info().Str("indexName", r.indexName).Msg("Создание индекса...")
		if err := r.client.CreateIndex(ctx, r.indexName, osClient.ListingMapping); err != nil {
			return fmt.Errorf("ошибка создания индекса: %w", err)
		}
		logger.Info().Str("indexName", r.indexName).Msg("Индекс успешно создан")

		allListings, _, err := r.storage.GetListings(ctx, map[string]string{}, 1000, 0)
		if err != nil {
			logger.Error().Err(err).Msg("Ошибка получения объявлений")
			return err
		}

		listingPtrs := make([]*models.MarketplaceListing, len(allListings))
		for i := range allListings {
			listingPtrs[i] = &allListings[i]
		}
		logger.Info().Msgf("Запуск переиндексации с обновленной схемой (поддержка метаданных и скидок)")
		if err := r.BulkIndexListings(ctx, listingPtrs); err != nil {
			logger.Error().Err(err).Msg("Ошибка индексации объявлений")
			return err
		}

		logger.Info().Int("listing_count", len(allListings)).Msgf("Успешно проиндексировано объявлений")
	}

	return nil
}
