// internal/storage/postgres/reviews.go

package postgres

import (
    "context"
    "backend/internal/domain/models"
	"fmt"
)

func (db *Database) CreateReview(ctx context.Context, review *models.Review) error {
    return db.pool.QueryRow(ctx, `
        INSERT INTO reviews (
            user_id, entity_type, entity_id, rating, comment, 
            pros, cons, photos, is_verified_purchase, status
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING id, created_at, updated_at
    `,
        review.UserID, review.EntityType, review.EntityID, review.Rating,
        review.Comment, review.Pros, review.Cons, review.Photos,
        review.IsVerifiedPurchase, review.Status,
    ).Scan(&review.ID, &review.CreatedAt, &review.UpdatedAt)
}

func (db *Database) AddReviewResponse(ctx context.Context, response *models.ReviewResponse) error {
    return db.pool.QueryRow(ctx, `
        INSERT INTO review_responses (review_id, user_id, response)
        VALUES ($1, $2, $3)
        RETURNING id, created_at, updated_at
    `,
        response.ReviewID, response.UserID, response.Response,
    ).Scan(&response.ID, &response.CreatedAt, &response.UpdatedAt)
}
// internal/storage/postgres/reviews.go

func (db *Database) GetReviews(ctx context.Context, filter models.ReviewsFilter) ([]models.Review, int64, error) {
    // Базовый запрос
    query := `
        SELECT 
            r.id, r.user_id, r.entity_type, r.entity_id, r.rating, 
            r.comment, r.pros, r.cons, r.photos, r.likes_count,
            r.is_verified_purchase, r.status, r.created_at, r.updated_at,
            u.name as user_name, u.email as user_email, u.picture_url as user_picture,
            COUNT(*) OVER() as total_count
        FROM reviews r
        LEFT JOIN users u ON r.user_id = u.id
        WHERE 1=1
    `
    params := make([]interface{}, 0)
    paramCount := 1

    // Добавляем условия фильтрации
    if filter.EntityType != "" {
        query += fmt.Sprintf(" AND r.entity_type = $%d", paramCount)
        params = append(params, filter.EntityType)
        paramCount++
    }

    if filter.EntityID != 0 {
        query += fmt.Sprintf(" AND r.entity_id = $%d", paramCount)
        params = append(params, filter.EntityID)
        paramCount++
    }

    if filter.UserID != 0 {
        query += fmt.Sprintf(" AND r.user_id = $%d", paramCount)
        params = append(params, filter.UserID)
        paramCount++
    }

    if filter.MinRating > 0 {
        query += fmt.Sprintf(" AND r.rating >= $%d", paramCount)
        params = append(params, filter.MinRating)
        paramCount++
    }

    if filter.MaxRating > 0 {
        query += fmt.Sprintf(" AND r.rating <= $%d", paramCount)
        params = append(params, filter.MaxRating)
        paramCount++
    }

    if filter.Status != "" {
        query += fmt.Sprintf(" AND r.status = $%d", paramCount)
        params = append(params, filter.Status)
        paramCount++
    }

    // Добавляем сортировку
    switch filter.SortBy {
    case "rating":
        query += " ORDER BY r.rating " + filter.SortOrder
    case "likes":
        query += " ORDER BY r.likes_count " + filter.SortOrder
    default:
        query += " ORDER BY r.created_at " + filter.SortOrder
    }

    // Добавляем пагинацию
    query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramCount, paramCount+1)
    params = append(params, filter.Limit, (filter.Page-1)*filter.Limit)

    rows, err := db.pool.Query(ctx, query, params...)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()

    var reviews []models.Review
    var totalCount int64

    for rows.Next() {
        var r models.Review
        r.User = &models.User{}
        err := rows.Scan(
            &r.ID, &r.UserID, &r.EntityType, &r.EntityID, &r.Rating,
            &r.Comment, &r.Pros, &r.Cons, &r.Photos, &r.LikesCount,
            &r.IsVerifiedPurchase, &r.Status, &r.CreatedAt, &r.UpdatedAt,
            &r.User.Name, &r.User.Email, &r.User.PictureURL,
            &totalCount,
        )
        if err != nil {
            return nil, 0, err
        }
        reviews = append(reviews, r)
    }

    return reviews, totalCount, nil
}

func (db *Database) GetReviewByID(ctx context.Context, id int) (*models.Review, error) {
    review := &models.Review{
        User: &models.User{},
    }

    err := db.pool.QueryRow(ctx, `
        SELECT 
            r.id, r.user_id, r.entity_type, r.entity_id, r.rating, 
            r.comment, r.pros, r.cons, r.photos, r.likes_count,
            r.is_verified_purchase, r.status, r.created_at, r.updated_at,
            u.name as user_name, u.email as user_email, u.picture_url as user_picture
        FROM reviews r
        LEFT JOIN users u ON r.user_id = u.id
        WHERE r.id = $1
    `, id).Scan(
        &review.ID, &review.UserID, &review.EntityType, &review.EntityID, &review.Rating,
        &review.Comment, &review.Pros, &review.Cons, &review.Photos, &review.LikesCount,
        &review.IsVerifiedPurchase, &review.Status, &review.CreatedAt, &review.UpdatedAt,
        &review.User.Name, &review.User.Email, &review.User.PictureURL,
    )

    if err != nil {
        return nil, err
    }

    // Получаем ответы на отзыв
    rows, err := db.pool.Query(ctx, `
        SELECT 
            rr.id, rr.user_id, rr.response, rr.created_at, rr.updated_at,
            u.name, u.email, u.picture_url
        FROM review_responses rr
        LEFT JOIN users u ON rr.user_id = u.id
        WHERE rr.review_id = $1
        ORDER BY rr.created_at
    `, review.ID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        response := models.ReviewResponse{
            User: &models.User{},
        }
        err := rows.Scan(
            &response.ID, &response.UserID, &response.Response,
            &response.CreatedAt, &response.UpdatedAt,
            &response.User.Name, &response.User.Email, &response.User.PictureURL,
        )
        if err != nil {
            return nil, err
        }
        review.Responses = append(review.Responses, response)
    }

    // Получаем голоса
    helpful, notHelpful, err := db.GetReviewVotes(ctx, review.ID)
    if err != nil {
        return nil, err
    }
    review.VotesCount.Helpful = helpful
    review.VotesCount.NotHelpful = notHelpful

    return review, nil
}

func (db *Database) UpdateReview(ctx context.Context, review *models.Review) error {
    _, err := db.pool.Exec(ctx, `
        UPDATE reviews 
        SET rating = $1, 
            comment = $2, 
            pros = $3, 
            cons = $4, 
            photos = $5, 
            status = $6,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $7
    `, 
        review.Rating, review.Comment, review.Pros, review.Cons, 
        review.Photos, review.Status, review.ID,
    )
    return err
}
func (db *Database) DeleteReview(ctx context.Context, id int) error {
    _, err := db.pool.Exec(ctx, `
        DELETE FROM reviews
        WHERE id = $1
    `, id)
    return err
}
func (db *Database) AddReviewVote(ctx context.Context, vote *models.ReviewVote) error {
    _, err := db.pool.Exec(ctx, `
        INSERT INTO review_votes (review_id, user_id, vote_type)
        VALUES ($1, $2, $3)
        ON CONFLICT (review_id, user_id) DO UPDATE
        SET vote_type = EXCLUDED.vote_type
    `,
        vote.ReviewID, vote.UserID, vote.VoteType,
    )
    return err
}

func (db *Database) GetReviewVotes(ctx context.Context, reviewId int) (helpful int, notHelpful int, err error) {
    err = db.pool.QueryRow(ctx, `
        SELECT 
            COUNT(CASE WHEN vote_type = 'helpful' THEN 1 END) as helpful,
            COUNT(CASE WHEN vote_type = 'not_helpful' THEN 1 END) as not_helpful
        FROM review_votes
        WHERE review_id = $1
    `, reviewId).Scan(&helpful, &notHelpful)
    return
}

func (db *Database) GetUserReviewVote(ctx context.Context, userId int, reviewId int) (string, error) {
    var voteType string
    err := db.pool.QueryRow(ctx, `
        SELECT vote_type FROM review_votes
        WHERE user_id = $1 AND review_id = $2
    `, userId, reviewId).Scan(&voteType)
    return voteType, err
}

func (db *Database) GetEntityRating(ctx context.Context, entityType string, entityId int) (float64, error) {
    var rating float64
    err := db.pool.QueryRow(ctx, `
        SELECT calculate_entity_rating($1, $2)
    `, entityType, entityId).Scan(&rating)
    return rating, err
}