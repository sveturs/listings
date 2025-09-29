package postgres

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
)

// HardDeleteDebug - временная функция для отладки удаления
func (r *storefrontRepo) HardDeleteDebug(ctx context.Context, id int) error {
	logger := log.With().Str("component", "storefront_repo_debug").Logger()

	// Начинаем транзакцию для атомарного удаления
	logger.Info().Int("storefrontID", id).Msg("Starting transaction")
	tx, err := r.db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		logger.Info().Msg("Rolling back transaction")
		_ = tx.Rollback(ctx)
	}()

	// Простое удаление только самой витрины
	logger.Info().Int("storefrontID", id).Msg("Deleting storefront")
	result, err := tx.Exec(ctx, "DELETE FROM storefronts WHERE id = $1", id)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to delete storefront")
		return fmt.Errorf("failed to delete storefront: %w", err)
	}

	logger.Info().Int64("rowsAffected", result.RowsAffected()).Msg("Storefront deleted")
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	// Коммитим транзакцию
	logger.Info().Msg("Committing transaction")
	if err := tx.Commit(ctx); err != nil {
		logger.Error().Err(err).Msg("Failed to commit transaction")
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	logger.Info().Int("storefrontID", id).Msg("Storefront deleted successfully")
	return nil
}