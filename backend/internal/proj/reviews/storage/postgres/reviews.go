// internal/proj/reviews/storage/postgres/review.go
package postgres

import (
    "backend/internal/domain/models"
    "context"
    "fmt"
    "log"
    "database/sql"
)

func (s *Storage) CreateReview(ctx context.Context, review *models.Review) (*models.Review, error) {
    err := s.pool.QueryRow(ctx, `
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

    if err != nil {
        return nil, err
    }

    return review, nil
}

func (s *Storage) GetReviews(ctx context.Context, filter models.ReviewsFilter) ([]models.Review, int64, error) {
    userID, ok := ctx.Value("user_id").(int)
    if !ok {
        userID = 0
    }

    query := `
    SELECT 
        r.id, r.user_id, r.entity_type, r.entity_id, r.rating, 
        r.comment, r.pros, r.cons, r.photos, r.likes_count,
        r.is_verified_purchase, r.status, r.created_at, r.updated_at,
        r.helpful_votes, r.not_helpful_votes,
        u.name as user_name, u.email as user_email, u.picture_url as user_picture,
        COUNT(*) OVER() as total_count,
        (
            SELECT vote_type 
            FROM review_votes 
            WHERE review_id = r.id AND user_id = $1
        ) as current_user_vote
    FROM reviews r
    LEFT JOIN users u ON r.user_id = u.id
    WHERE 1=1`

    params := []interface{}{userID}
    paramCount := 2

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

    rows, err := s.pool.Query(ctx, query, params...)
    if err != nil {
        return nil, 0, fmt.Errorf("error executing query: %w", err)
    }
    defer rows.Close()

    var reviews []models.Review
    var totalCount int64

    for rows.Next() {
        var r models.Review
        r.User = &models.User{}
        var currentUserVote sql.NullString
        err := rows.Scan(
            &r.ID, &r.UserID, &r.EntityType, &r.EntityID, &r.Rating,
            &r.Comment, &r.Pros, &r.Cons, &r.Photos, &r.LikesCount,
            &r.IsVerifiedPurchase, &r.Status, &r.CreatedAt, &r.UpdatedAt,
            &r.HelpfulVotes, &r.NotHelpfulVotes,
            &r.User.Name, &r.User.Email, &r.User.PictureURL,
            &totalCount,
            &currentUserVote,
        )
        if err != nil {
            return nil, 0, fmt.Errorf("error scanning row: %w", err)
        }
        if currentUserVote.Valid {
            r.CurrentUserVote = currentUserVote.String
        }
        
        reviews = append(reviews, r)
    }

    return reviews, totalCount, nil
}

func (s *Storage) GetReviewByID(ctx context.Context, id int) (*models.Review, error) {
    review := &models.Review{
        User: &models.User{},
    }

    err := s.pool.QueryRow(ctx, `
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
    rows, err := s.pool.Query(ctx, `
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

    helpful, notHelpful, err := s.GetReviewVotes(ctx, review.ID)
    if err != nil {
        return nil, err
    }
    review.VotesCount.Helpful = helpful
    review.VotesCount.NotHelpful = notHelpful

    return review, nil
}

func (s *Storage) UpdateReview(ctx context.Context, review *models.Review) error {
    _, err := s.pool.Exec(ctx, `
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

func (s *Storage) DeleteReview(ctx context.Context, id int) error {
    _, err := s.pool.Exec(ctx, `DELETE FROM reviews WHERE id = $1`, id)
    return err
}

func (s *Storage) AddReviewResponse(ctx context.Context, response *models.ReviewResponse) error {
    return s.pool.QueryRow(ctx, `
        INSERT INTO review_responses (review_id, user_id, response)
        VALUES ($1, $2, $3)
        RETURNING id, created_at, updated_at
    `,
        response.ReviewID, response.UserID, response.Response,
    ).Scan(&response.ID, &response.CreatedAt, &response.UpdatedAt)
}

func (s *Storage) AddReviewVote(ctx context.Context, vote *models.ReviewVote) error {
    _, err := s.pool.Exec(ctx, `
        INSERT INTO review_votes (review_id, user_id, vote_type)
        VALUES ($1, $2, $3)
        ON CONFLICT (review_id, user_id) DO UPDATE
        SET vote_type = EXCLUDED.vote_type
    `,
        vote.ReviewID, vote.UserID, vote.VoteType,
    )
    if err != nil {
        return fmt.Errorf("failed to add or update review vote: %w", err)
    }

    // Пересчет голосов
    err = s.UpdateReviewVotes(ctx, vote.ReviewID)
    if err != nil {
        return fmt.Errorf("failed to update review votes for review_id %d: %w", vote.ReviewID, err)
    }

    return nil
}

func (s *Storage) UpdateReviewVotes(ctx context.Context, reviewId int) error {
    _, err := s.pool.Exec(ctx, `
        WITH vote_counts AS (
            SELECT 
                SUM(CASE WHEN vote_type = 'helpful' THEN 1 ELSE 0 END) as helpful_count,
                SUM(CASE WHEN vote_type = 'not_helpful' THEN 1 ELSE 0 END) as not_helpful_count
            FROM review_votes
            WHERE review_id = $1
        )
        UPDATE reviews r
        SET 
            helpful_votes = COALESCE(vote_counts.helpful_count, 0),
            not_helpful_votes = COALESCE(vote_counts.not_helpful_count, 0),
            updated_at = CURRENT_TIMESTAMP
        FROM vote_counts
        WHERE r.id = $1
    `, reviewId)

    if err != nil {
        log.Printf("Failed to update votes for review_id %d: %v", reviewId, err)
        return fmt.Errorf("failed to update votes for review_id %d: %w", reviewId, err)
    }

    return nil
}

func (s *Storage) GetReviewVotes(ctx context.Context, reviewId int) (helpful int, notHelpful int, err error) {
    err = s.pool.QueryRow(ctx, `
        SELECT 
            COUNT(CASE WHEN vote_type = 'helpful' THEN 1 END) as helpful,
            COUNT(CASE WHEN vote_type = 'not_helpful' THEN 1 END) as not_helpful
        FROM review_votes
        WHERE review_id = $1
    `, reviewId).Scan(&helpful, &notHelpful)
    if err != nil {
        return 0, 0, fmt.Errorf("failed to fetch review votes: %w", err)
    }
    return helpful, notHelpful, nil
}

func (s *Storage) GetUserReviewVote(ctx context.Context, userId int, reviewId int) (string, error) {
    var voteType string
    err := s.pool.QueryRow(ctx, `
        SELECT vote_type FROM review_votes
        WHERE user_id = $1 AND review_id = $2
    `, userId, reviewId).Scan(&voteType)
    return voteType, err
}

func (s *Storage) GetEntityRating(ctx context.Context, entityType string, entityId int) (float64, error) {
    var rating float64
    err := s.pool.QueryRow(ctx, `
        SELECT calculate_entity_rating($1, $2)
    `, entityType, entityId).Scan(&rating)
    return rating, err
}

func (s *Storage) GetReviewStats(ctx context.Context, entityType string, entityId int) (*models.ReviewStats, error) {
    stats := &models.ReviewStats{
        RatingDistribution: make(map[int]int),
    }

    err := s.pool.QueryRow(ctx, `
        SELECT 
            COUNT(*) as total,
            COALESCE(AVG(rating), 0) as avg_rating,
            COUNT(*) FILTER (WHERE is_verified_purchase) as verified,
            COUNT(*) FILTER (WHERE array_length(photos, 1) > 0) as with_photos
        FROM reviews
        WHERE entity_type = $1 
        AND entity_id = $2
        AND status = 'published'
    `, entityType, entityId).Scan(
        &stats.TotalReviews,
        &stats.AverageRating,
        &stats.VerifiedReviews,
        &stats.PhotoReviews,
    )
    if err != nil {
        return nil, err
    }

    // Получаем распределение оценок
    rows, err := s.pool.Query(ctx, `
        SELECT rating, COUNT(*)
        FROM reviews
        WHERE entity_type = $1 
        AND entity_id = $2
        AND status = 'published'
        GROUP BY rating
    `, entityType, entityId)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var rating, count int
        if err := rows.Scan(&rating, &count); err != nil {
            return nil, err
        }
        stats.RatingDistribution[rating] = count
    }

    return stats, nil
}