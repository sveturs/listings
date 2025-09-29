package postgres

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
)

// HardDeleteFixed - исправленная версия функции HardDelete
func (r *storefrontRepo) HardDeleteFixed(ctx context.Context, id int) error {
	logger := log.With().Str("component", "storefront_repo").Logger()

	// Начинаем транзакцию для атомарного удаления
	tx, err := r.db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Удаляем связанные записи в правильном порядке
	// 1. Сначала получаем ID всех товаров витрины
	logger.Debug().Int("storefrontID", id).Msg("Getting product IDs")
	var productIDs []int
	rows, err := tx.Query(ctx, "SELECT id FROM storefront_products WHERE storefront_id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to get product IDs: %w", err)
	}
	for rows.Next() {
		var productID int
		if err := rows.Scan(&productID); err != nil {
			rows.Close()
			return fmt.Errorf("failed to scan product ID: %w", err)
		}
		productIDs = append(productIDs, productID)
	}
	rows.Close()
	logger.Debug().Ints("productIDs", productIDs).Msg("Found products")

	// 2. Удаляем связанные с товарами записи
	if len(productIDs) > 0 {
		// Удаляем товары из избранного
		logger.Debug().Msg("Deleting from storefront_favorites")
		_, err = tx.Exec(ctx, "DELETE FROM storefront_favorites WHERE product_id = ANY($1)", productIDs)
		if err != nil {
			return fmt.Errorf("failed to delete favorites: %w", err)
		}

		// Удаляем изображения товаров
		logger.Debug().Msg("Deleting from storefront_product_images")
		_, err = tx.Exec(ctx, "DELETE FROM storefront_product_images WHERE storefront_product_id = ANY($1)", productIDs)
		if err != nil {
			return fmt.Errorf("failed to delete product images: %w", err)
		}

		// Удаляем варианты товаров
		logger.Debug().Msg("Deleting from storefront_product_variants")
		_, err = tx.Exec(ctx, "DELETE FROM storefront_product_variants WHERE product_id = ANY($1)", productIDs)
		if err != nil {
			return fmt.Errorf("failed to delete product variants: %w", err)
		}
	}

	// 3. Удаляем заказы
	logger.Debug().Msg("Deleting from storefront_orders")
	_, err = tx.Exec(ctx, "DELETE FROM storefront_orders WHERE storefront_id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete storefront orders: %w", err)
	}

	// 4. Удаляем товары
	logger.Debug().Msg("Deleting from storefront_products")
	_, err = tx.Exec(ctx, "DELETE FROM storefront_products WHERE storefront_id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete storefront products: %w", err)
	}

	// 5. Удаляем сотрудников
	logger.Debug().Msg("Deleting from storefront_staff")
	_, err = tx.Exec(ctx, "DELETE FROM storefront_staff WHERE storefront_id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete storefront staff: %w", err)
	}

	// 6. Удаляем витрину из избранного (игнорируем ошибку если таблицы нет)
	logger.Debug().Msg("Trying to delete from user_favorite_storefronts")
	_, _ = tx.Exec(ctx, "DELETE FROM user_favorite_storefronts WHERE storefront_id = $1", id)

	// 7. Удаляем саму витрину
	logger.Debug().Msg("Deleting storefront itself")
	result, err := tx.Exec(ctx, "DELETE FROM storefronts WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete storefront: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	// Коммитим транзакцию
	logger.Debug().Msg("Committing transaction")
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	logger.Info().Int("storefrontID", id).Msg("Storefront hard deleted successfully")
	return nil
}