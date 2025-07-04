// backend/internal/storage/postgres/price_history.go
package postgres

import (
	"context"
	"fmt"
	"time"

	"backend/internal/domain/models"
)

func (db *Database) GetPriceHistory(ctx context.Context, listingID int) ([]models.PriceHistoryEntry, error) {
	query := `
        SELECT 
            id, listing_id, price, effective_from, effective_to, 
            change_source, COALESCE(change_percentage, 0), created_at
        FROM price_history
        WHERE listing_id = $1
        ORDER BY effective_from DESC
    `

	rows, err := db.pool.Query(ctx, query, listingID)
	if err != nil {
		return nil, fmt.Errorf("error querying price history: %w", err)
	}
	defer rows.Close()

	var history []models.PriceHistoryEntry
	for rows.Next() {
		var entry models.PriceHistoryEntry
		if err := rows.Scan(
			&entry.ID,
			&entry.ListingID,
			&entry.Price,
			&entry.EffectiveFrom,
			&entry.EffectiveTo,
			&entry.ChangeSource,
			&entry.ChangePercentage,
			&entry.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning price history: %w", err)
		}
		history = append(history, entry)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating price history: %w", err)
	}

	return history, nil
}

// AddPriceHistoryEntry добавляет запись в историю цен
func (db *Database) AddPriceHistoryEntry(ctx context.Context, entry *models.PriceHistoryEntry) error {
	query := `
        INSERT INTO price_history (
            listing_id, price, effective_from, change_source, change_percentage
        ) VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `

	var id int
	err := db.pool.QueryRow(ctx, query,
		entry.ListingID,
		entry.Price,
		entry.EffectiveFrom,
		entry.ChangeSource,
		entry.ChangePercentage,
	).Scan(&id)
	if err != nil {
		return fmt.Errorf("error inserting price history: %w", err)
	}

	entry.ID = id
	return nil
}

// ClosePriceHistoryEntry закрывает текущую активную запись в истории цен
func (db *Database) ClosePriceHistoryEntry(ctx context.Context, listingID int) error {
	query := `
        UPDATE price_history
        SET effective_to = $1
        WHERE listing_id = $2
        AND effective_to IS NULL
    `

	_, err := db.pool.Exec(ctx, query, time.Now(), listingID)
	if err != nil {
		return fmt.Errorf("error closing price history entry: %w", err)
	}

	return nil
}

// CheckPriceManipulation проверяет, есть ли признаки манипуляций с ценой
func (db *Database) CheckPriceManipulation(ctx context.Context, listingID int) (bool, error) {
	query := `SELECT check_price_manipulation($1)`

	var isSuspicious bool
	err := db.pool.QueryRow(ctx, query, listingID).Scan(&isSuspicious)
	if err != nil {
		return false, fmt.Errorf("error checking price manipulation: %w", err)
	}

	return isSuspicious, nil
}
