package postgres

import (
	"context"
	"fmt"
)

// GetFavoritedUsers retrieves list of user IDs who favorited a listing
func (r *Repository) GetFavoritedUsers(ctx context.Context, listingID int64) ([]int64, error) {
	query := `
		SELECT DISTINCT user_id
		FROM c2c_favorites
		WHERE listing_id = $1
		ORDER BY user_id ASC
	`

	rows, err := r.db.QueryxContext(ctx, query, listingID)
	if err != nil {
		r.logger.Error().Err(err).Int64("listing_id", listingID).Msg("failed to query favorited users")
		return nil, fmt.Errorf("failed to query favorited users: %w", err)
	}
	defer rows.Close()

	var userIDs []int64
	for rows.Next() {
		var userID int64
		if err := rows.Scan(&userID); err != nil {
			r.logger.Error().Err(err).Msg("failed to scan user ID")
			return nil, fmt.Errorf("failed to scan user ID: %w", err)
		}
		userIDs = append(userIDs, userID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return userIDs, nil
}

// AddToFavorites adds a listing to user's favorites
func (r *Repository) AddToFavorites(ctx context.Context, userID, listingID int64) error {
	query := `
		INSERT INTO c2c_favorites (user_id, listing_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, listing_id) DO NOTHING
	`

	_, err := r.db.ExecContext(ctx, query, userID, listingID)
	if err != nil {
		r.logger.Error().Err(err).Int64("user_id", userID).Int64("listing_id", listingID).Msg("failed to add to favorites")
		return fmt.Errorf("failed to add to favorites: %w", err)
	}

	r.logger.Debug().Int64("user_id", userID).Int64("listing_id", listingID).Msg("added to favorites")
	return nil
}

// RemoveFromFavorites removes a listing from user's favorites
func (r *Repository) RemoveFromFavorites(ctx context.Context, userID, listingID int64) error {
	query := `
		DELETE FROM c2c_favorites
		WHERE user_id = $1 AND listing_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, userID, listingID)
	if err != nil {
		r.logger.Error().Err(err).Int64("user_id", userID).Int64("listing_id", listingID).Msg("failed to remove from favorites")
		return fmt.Errorf("failed to remove from favorites: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	r.logger.Debug().Int64("user_id", userID).Int64("listing_id", listingID).Int64("rows_affected", rowsAffected).Msg("removed from favorites")
	return nil
}

// GetUserFavorites retrieves list of listing IDs favorited by a user
func (r *Repository) GetUserFavorites(ctx context.Context, userID int64) ([]int64, error) {
	query := `
		SELECT listing_id
		FROM c2c_favorites
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryxContext(ctx, query, userID)
	if err != nil {
		r.logger.Error().Err(err).Int64("user_id", userID).Msg("failed to query user favorites")
		return nil, fmt.Errorf("failed to query user favorites: %w", err)
	}
	defer rows.Close()

	var listingIDs []int64
	for rows.Next() {
		var listingID int64
		if err := rows.Scan(&listingID); err != nil {
			r.logger.Error().Err(err).Msg("failed to scan listing ID")
			return nil, fmt.Errorf("failed to scan listing ID: %w", err)
		}
		listingIDs = append(listingIDs, listingID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return listingIDs, nil
}

// IsFavorite checks if a listing is in user's favorites
func (r *Repository) IsFavorite(ctx context.Context, userID, listingID int64) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM c2c_favorites
			WHERE user_id = $1 AND listing_id = $2
		)
	`

	var exists bool
	err := r.db.QueryRowxContext(ctx, query, userID, listingID).Scan(&exists)
	if err != nil {
		r.logger.Error().Err(err).Int64("user_id", userID).Int64("listing_id", listingID).Msg("failed to check favorite status")
		return false, fmt.Errorf("failed to check favorite status: %w", err)
	}

	return exists, nil
}
