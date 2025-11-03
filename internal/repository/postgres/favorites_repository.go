package postgres

import (
	"context"
	"fmt"
)

// GetFavoritedUsers retrieves list of user IDs who favorited a listing
func (r *Repository) GetFavoritedUsers(ctx context.Context, listingID int64) ([]int64, error) {
	query := `
		SELECT DISTINCT user_id
		FROM user_favorites
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
