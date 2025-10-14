// backend/internal/storage/postgres/db_reviews.go
package postgres

import (
	"context"
	"errors"

	"backend/internal/domain/models"

	"github.com/jackc/pgx/v5"
)

func (db *Database) CreateReview(ctx context.Context, review *models.Review) (*models.Review, error) {
	return db.reviewDB.CreateReview(ctx, review)
}

func (db *Database) GetReviews(ctx context.Context, filter models.ReviewsFilter) ([]models.Review, int64, error) {
	return db.reviewDB.GetReviews(ctx, filter)
}

func (db *Database) GetReviewByID(ctx context.Context, id int) (*models.Review, error) {
	return db.reviewDB.GetReviewByID(ctx, id)
}

func (db *Database) UpdateReview(ctx context.Context, review *models.Review) error {
	return db.reviewDB.UpdateReview(ctx, review)
}

func (db *Database) UpdateReviewStatus(ctx context.Context, reviewId int, status string) error {
	return db.reviewDB.UpdateReviewStatus(ctx, reviewId, status)
}

func (db *Database) DeleteReview(ctx context.Context, id int) error {
	return db.reviewDB.DeleteReview(ctx, id)
}

func (db *Database) AddReviewResponse(ctx context.Context, response *models.ReviewResponse) error {
	return db.reviewDB.AddReviewResponse(ctx, response)
}

func (db *Database) AddReviewVote(ctx context.Context, vote *models.ReviewVote) error {
	return db.reviewDB.AddReviewVote(ctx, vote)
}

func (db *Database) GetReviewVotes(ctx context.Context, reviewId int) (helpful int, notHelpful int, err error) {
	return db.reviewDB.GetReviewVotes(ctx, reviewId)
}

func (db *Database) GetUserReviewVote(ctx context.Context, userId int, reviewId int) (string, error) {
	return db.reviewDB.GetUserReviewVote(ctx, userId, reviewId)
}

func (db *Database) GetEntityRating(ctx context.Context, entityType string, entityId int) (float64, error) {
	return db.reviewDB.GetEntityRating(ctx, entityType, entityId)
}

func (db *Database) GetUserReviews(ctx context.Context, userID int, filter models.ReviewsFilter) ([]models.Review, error) {
	return db.reviewDB.GetUserReviews(ctx, userID, filter)
}

func (db *Database) GetStorefrontReviews(ctx context.Context, storefrontID int, filter models.ReviewsFilter) ([]models.Review, error) {
	return db.reviewDB.GetStorefrontReviews(ctx, storefrontID, filter)
}

func (db *Database) GetUserRatingSummary(ctx context.Context, userID int) (*models.UserRatingSummary, error) {
	return db.reviewDB.GetUserRatingSummary(ctx, userID)
}

func (db *Database) GetStorefrontRatingSummary(ctx context.Context, storefrontID int) (*models.StorefrontRatingSummary, error) {
	return db.reviewDB.GetStorefrontRatingSummary(ctx, storefrontID)
}

func (db *Database) GetUserAggregatedRating(ctx context.Context, userID int) (*models.UserAggregatedRating, error) {
	rating := &models.UserAggregatedRating{UserID: userID}

	query := `
		SELECT
			total_reviews, average_rating, direct_reviews, listing_reviews,
			storefront_reviews, verified_reviews, rating_1, rating_2, rating_3,
			rating_4, rating_5, recent_rating, recent_reviews, last_review_at
		FROM user_ratings
		WHERE user_id = $1
	`

	row := db.pool.QueryRow(ctx, query, userID)
	err := row.Scan(
		&rating.TotalReviews, &rating.AverageRating, &rating.DirectReviews,
		&rating.ListingReviews, &rating.StorefrontReviews, &rating.VerifiedReviews,
		&rating.Rating1, &rating.Rating2, &rating.Rating3,
		&rating.Rating4, &rating.Rating5, &rating.RecentRating,
		&rating.RecentReviews, &rating.LastReviewAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		// Если нет данных в материализованном представлении, возвращаем пустой рейтинг
		return rating, nil
	}

	return rating, err
}

// GetStorefrontAggregatedRating получает агрегированный рейтинг магазина
func (db *Database) GetStorefrontAggregatedRating(ctx context.Context, storefrontID int) (*models.StorefrontAggregatedRating, error) {
	rating := &models.StorefrontAggregatedRating{StorefrontID: storefrontID}

	query := `
		SELECT
			total_reviews, average_rating, direct_reviews, listing_reviews,
			verified_reviews, rating_1, rating_2, rating_3, rating_4, rating_5,
			recent_rating, recent_reviews, last_review_at, owner_id
		FROM b2c_rating_summary
		WHERE storefront_id = $1
	`

	row := db.pool.QueryRow(ctx, query, storefrontID)
	err := row.Scan(
		&rating.TotalReviews, &rating.AverageRating, &rating.DirectReviews,
		&rating.ListingReviews, &rating.VerifiedReviews,
		&rating.Rating1, &rating.Rating2, &rating.Rating3,
		&rating.Rating4, &rating.Rating5, &rating.RecentRating,
		&rating.RecentReviews, &rating.LastReviewAt, &rating.OwnerID,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return rating, nil
	}

	return rating, err
}

// RefreshRatingViews обновляет материализованные представления
func (db *Database) RefreshRatingViews(ctx context.Context) error {
	_, err := db.pool.Exec(ctx, "SELECT rebuild_all_ratings()")
	return err
}

// CreateReviewConfirmation создает подтверждение отзыва
func (db *Database) CreateReviewConfirmation(ctx context.Context, confirmation *models.ReviewConfirmation) error {
	query := `
		INSERT INTO review_confirmations
		(review_id, confirmed_by, confirmation_status, notes)
		VALUES ($1, $2, $3, $4)
		RETURNING id, confirmed_at
	`

	row := db.pool.QueryRow(ctx, query,
		confirmation.ReviewID, confirmation.ConfirmedBy,
		confirmation.ConfirmationStatus, confirmation.Notes,
	)

	return row.Scan(&confirmation.ID, &confirmation.ConfirmedAt)
}

// GetReviewConfirmation получает подтверждение отзыва
func (db *Database) GetReviewConfirmation(ctx context.Context, reviewID int) (*models.ReviewConfirmation, error) {
	confirmation := &models.ReviewConfirmation{}

	query := `
		SELECT id, review_id, confirmed_by, confirmation_status, confirmed_at, notes
		FROM review_confirmations
		WHERE review_id = $1
	`

	row := db.pool.QueryRow(ctx, query, reviewID)
	err := row.Scan(
		&confirmation.ID, &confirmation.ReviewID, &confirmation.ConfirmedBy,
		&confirmation.ConfirmationStatus, &confirmation.ConfirmedAt, &confirmation.Notes,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrReviewConfirmationNotFound
	}

	return confirmation, err
}

// CreateReviewDispute создает спор по отзыву
func (db *Database) CreateReviewDispute(ctx context.Context, dispute *models.ReviewDispute) error {
	query := `
		INSERT INTO review_disputes
		(review_id, disputed_by, dispute_reason, dispute_description, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	row := db.pool.QueryRow(ctx, query,
		dispute.ReviewID, dispute.DisputedBy, dispute.DisputeReason,
		dispute.DisputeDescription, dispute.Status,
	)

	err := row.Scan(&dispute.ID, &dispute.CreatedAt, &dispute.UpdatedAt)
	if err != nil {
		return err
	}

	// Обновляем флаг в таблице reviews
	_, err = db.pool.Exec(ctx,
		"UPDATE reviews SET has_active_dispute = true WHERE id = $1",
		dispute.ReviewID,
	)

	return err
}

// GetReviewDispute получает спор по отзыву
func (db *Database) GetReviewDispute(ctx context.Context, reviewID int) (*models.ReviewDispute, error) {
	dispute := &models.ReviewDispute{}

	query := `
		SELECT id, review_id, disputed_by, dispute_reason, dispute_description,
			   status, admin_id, admin_notes, created_at, updated_at, resolved_at
		FROM review_disputes
		WHERE review_id = $1 AND status NOT IN ('resolved_keep_review', 'resolved_remove_review', 'cancelled')
		ORDER BY created_at DESC
		LIMIT 1
	`

	row := db.pool.QueryRow(ctx, query, reviewID)
	err := row.Scan(
		&dispute.ID, &dispute.ReviewID, &dispute.DisputedBy,
		&dispute.DisputeReason, &dispute.DisputeDescription, &dispute.Status,
		&dispute.AdminID, &dispute.AdminNotes, &dispute.CreatedAt,
		&dispute.UpdatedAt, &dispute.ResolvedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrReviewDisputeNotFound
	}

	return dispute, err
}

// UpdateReviewDispute обновляет спор
func (db *Database) UpdateReviewDispute(ctx context.Context, dispute *models.ReviewDispute) error {
	query := `
		UPDATE review_disputes
		SET status = $2, admin_id = $3, admin_notes = $4,
			updated_at = NOW(), resolved_at = $5
		WHERE id = $1
	`

	_, err := db.pool.Exec(ctx, query,
		dispute.ID, dispute.Status, dispute.AdminID,
		dispute.AdminNotes, dispute.ResolvedAt,
	)
	if err != nil {
		return err
	}

	// Если спор разрешен, обновляем флаг в reviews
	if dispute.Status == "resolved_keep_review" ||
		dispute.Status == "resolved_remove_review" ||
		dispute.Status == "canceled" {
		_, err = db.pool.Exec(ctx,
			"UPDATE reviews SET has_active_dispute = false WHERE id = $1",
			dispute.ReviewID,
		)
	}

	return err
}

// CanUserReviewEntity проверяет может ли пользователь оставить отзыв
func (db *Database) CanUserReviewEntity(ctx context.Context, userID int, entityType string, entityID int) (*models.CanReviewResponse, error) {
	response := &models.CanReviewResponse{
		CanReview: true,
	}

	// Проверяем существующий отзыв
	var existingReviewID *int
	query := `
		SELECT id FROM reviews
		WHERE user_id = $1 AND entity_type = $2 AND entity_id = $3
		AND status != 'deleted'
		LIMIT 1
	`

	err := db.pool.QueryRow(ctx, query, userID, entityType, entityID).Scan(&existingReviewID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	if existingReviewID != nil {
		response.CanReview = false
		response.HasExistingReview = true
		response.ExistingReviewID = existingReviewID
		response.Reason = "Вы уже оставили отзыв на этот объект"
		return response, nil
	}

	// Для отзывов на товары проверяем владельца
	if entityType == "listing" {
		var ownerID int
		err := db.pool.QueryRow(ctx,
			"SELECT user_id FROM c2c_listings WHERE id = $1",
			entityID,
		).Scan(&ownerID)

		if err == nil && ownerID == userID {
			response.CanReview = false
			response.Reason = "Вы не можете оставить отзыв на свой товар"
			return response, nil
		}
	}

	return response, nil
}
