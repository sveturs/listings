// internal/proj/reviews/storage/postgres/review.go
package postgres

import (
	"backend/internal/domain/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
    
	"github.com/jackc/pgx/v5"
)

func (s *Storage) CreateReview(ctx context.Context, review *models.Review) (*models.Review, error) {
    tx, err := s.pool.Begin(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to start transaction: %w", err)
    }
    defer tx.Rollback(ctx)

    // Сначала модерируем оригинальный текст
    moderatedComment := review.Comment
    if review.Comment != "" {
        moderatedComment, err = s.translationService.ModerateText(ctx, review.Comment, review.OriginalLanguage)
        if err != nil {
            return nil, fmt.Errorf("failed to moderate comment: %w", err)
        }
    }

    // Создаем запись отзыва с модерированным текстом
    err = tx.QueryRow(ctx, `
        INSERT INTO reviews (
            user_id, entity_type, entity_id, rating, comment, 
            pros, cons, photos, is_verified_purchase, status,
            original_language
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        RETURNING id, created_at, updated_at
    `,
        review.UserID, review.EntityType, review.EntityID, review.Rating,
        moderatedComment, review.Pros, review.Cons, review.Photos,
        review.IsVerifiedPurchase, review.Status, review.OriginalLanguage,
    ).Scan(&review.ID, &review.CreatedAt, &review.UpdatedAt)

    if err != nil {
        return nil, fmt.Errorf("failed to insert review: %w", err)
    }

    // Сохраняем модерированный текст как оригинальный перевод
    err = s.saveTranslation(ctx, tx, "review", review.ID, review.OriginalLanguage, "comment", moderatedComment, false, true)
    if err != nil {
        return nil, fmt.Errorf("failed to save original translation: %w", err)
    }

    // Создаем переводы на другие языки
    targetLangs := []string{"en", "ru", "sr"}
    for _, lang := range targetLangs {
        if lang == review.OriginalLanguage {
            continue
        }

        translatedText, err := s.translationService.Translate(ctx, moderatedComment, review.OriginalLanguage, lang)
        if err != nil {
            log.Printf("Failed to translate to %s: %v", lang, err)
            continue
        }

        if err := s.saveTranslation(ctx, tx, "review", review.ID, lang, "comment", translatedText, true, false); err != nil {
            log.Printf("Failed to save translation to %s: %v", lang, err)
            continue
        }
    }

    if err = tx.Commit(ctx); err != nil {
        return nil, fmt.Errorf("failed to commit transaction: %w", err)
    }

	// Загружаем все переводы в структуру отзыва
	translations := make(map[string]map[string]string)
	rows, err := s.pool.Query(ctx, `
        SELECT language, field_name, translated_text 
        FROM translations 
        WHERE entity_type = 'review' AND entity_id = $1
    `, review.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load translations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var lang, field, text string
		if err := rows.Scan(&lang, &field, &text); err != nil {
			return nil, fmt.Errorf("failed to scan translation: %w", err)
		}
		if translations[lang] == nil {
			translations[lang] = make(map[string]string)
		}
		translations[lang][field] = text
	}
	review.Translations = translations

	return review, nil
}

func (s *Storage) GetReviews(ctx context.Context, filter models.ReviewsFilter) ([]models.Review, int64, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		userID = 0
	}

	baseQuery := `
    WITH vote_counts AS (
    SELECT 
        review_id,
        COUNT(*) FILTER (WHERE vote_type = 'helpful') as helpful_votes,
        COUNT(*) FILTER (WHERE vote_type = 'not_helpful') as not_helpful_votes
    FROM review_votes 
    GROUP BY review_id
),
translations_agg AS (
    SELECT 
        entity_id,
        jsonb_object_agg(
            t.language,
            jsonb_build_object(t.field_name, t.translated_text)
        ) as translations
    FROM translations t
    WHERE entity_type = 'review'
    GROUP BY entity_id
)
    SELECT 
        r.id, r.user_id, r.entity_type, r.entity_id, r.rating, 
        r.comment, r.pros, r.cons, r.photos, r.likes_count,
        r.is_verified_purchase, r.status, r.created_at, r.updated_at,
        r.original_language,
        COALESCE(vc.helpful_votes, 0) as helpful_votes,
        COALESCE(vc.not_helpful_votes, 0) as not_helpful_votes,
        u.name as user_name, u.email as user_email, u.picture_url as user_picture,
        COUNT(*) OVER() as total_count,
        COALESCE(ta.translations, '{}'::jsonb) as translations,
        (
            SELECT vote_type 
            FROM review_votes 
            WHERE review_id = r.id AND user_id = $1
        ) as current_user_vote
    FROM reviews r
    LEFT JOIN users u ON r.user_id = u.id
    LEFT JOIN vote_counts vc ON vc.review_id = r.id
    LEFT JOIN translations_agg ta ON ta.entity_id = r.id
    WHERE 1=1`

	params := []interface{}{userID}
	paramCount := 2

	// Добавляем условия фильтрации
	if filter.EntityType != "" {
		baseQuery += fmt.Sprintf(" AND r.entity_type = $%d", paramCount)
		params = append(params, filter.EntityType)
		paramCount++
	}

	if filter.EntityID != 0 {
		baseQuery += fmt.Sprintf(" AND r.entity_id = $%d", paramCount)
		params = append(params, filter.EntityID)
		paramCount++
	}

	if filter.UserID != 0 {
		baseQuery += fmt.Sprintf(" AND r.user_id = $%d", paramCount)
		params = append(params, filter.UserID)
		paramCount++
	}

	if filter.MinRating > 0 {
		baseQuery += fmt.Sprintf(" AND r.rating >= $%d", paramCount)
		params = append(params, filter.MinRating)
		paramCount++
	}

	if filter.MaxRating > 0 {
		baseQuery += fmt.Sprintf(" AND r.rating <= $%d", paramCount)
		params = append(params, filter.MaxRating)
		paramCount++
	}

	if filter.Status != "" {
		baseQuery += fmt.Sprintf(" AND r.status = $%d", paramCount)
		params = append(params, filter.Status)
		paramCount++
	}

	// Сортировка
	switch filter.SortBy {
	case "rating":
		baseQuery += " ORDER BY r.rating " + filter.SortOrder
	case "likes":
		baseQuery += " ORDER BY r.likes_count " + filter.SortOrder
	default:
		baseQuery += " ORDER BY r.created_at " + filter.SortOrder
	}

	// Пагинация
	baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramCount, paramCount+1)
	params = append(params, filter.Limit, (filter.Page-1)*filter.Limit)

	rows, err := s.pool.Query(ctx, baseQuery, params...)
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
		var translationsJSON []byte

		err := rows.Scan(
			&r.ID, &r.UserID, &r.EntityType, &r.EntityID, &r.Rating,
			&r.Comment, &r.Pros, &r.Cons, &r.Photos, &r.LikesCount,
			&r.IsVerifiedPurchase, &r.Status, &r.CreatedAt, &r.UpdatedAt,
			&r.OriginalLanguage,
			&r.HelpfulVotes, &r.NotHelpfulVotes,
			&r.User.Name, &r.User.Email, &r.User.PictureURL,
			&totalCount,
			&translationsJSON,
			&currentUserVote,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning row: %w", err)
		}

		// Парсим переводы из JSON
		if err := json.Unmarshal(translationsJSON, &r.Translations); err != nil {
			log.Printf("Error unmarshaling translations for review %d: %v", r.ID, err)
			r.Translations = make(map[string]map[string]string)
		}

		r.VotesCount = struct {
			Helpful    int `json:"helpful"`
			NotHelpful int `json:"not_helpful"`
		}{
			Helpful:    r.HelpfulVotes,
			NotHelpful: r.NotHelpfulVotes,
		}

		if currentUserVote.Valid {
			r.CurrentUserVote = currentUserVote.String
		}

		reviews = append(reviews, r)
	}

	return reviews, totalCount, nil
}

func (s *Storage) GetReviewByID(ctx context.Context, id int) (*models.Review, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		userID = 0
	}

	review := &models.Review{
		User: &models.User{},
	}

	// Используем NullString для всех полей, которые могут быть NULL
	var (
		comment, pros, cons                 sql.NullString
		photos                              []string
		userName, userEmail, userPictureURL sql.NullString
		currentUserVote                     sql.NullString
	)

	err := s.pool.QueryRow(ctx, `
        WITH votes_summary AS (
            SELECT 
                COUNT(*) FILTER (WHERE vote_type = 'helpful') as helpful_count,
                COUNT(*) FILTER (WHERE vote_type = 'not_helpful') as not_helpful_count
            FROM review_votes
            WHERE review_id = $1
        )
        SELECT 
            r.id, r.user_id, r.entity_type, r.entity_id, r.rating,
            r.comment, r.pros, r.cons, r.photos, r.likes_count,
            r.is_verified_purchase, r.status, r.created_at, r.updated_at,
            u.name, u.email, u.picture_url,
            vs.helpful_count, vs.not_helpful_count,
            (
                SELECT vote_type 
                FROM review_votes 
                WHERE review_id = r.id AND user_id = $2
            ) as current_user_vote
        FROM reviews r
        LEFT JOIN users u ON r.user_id = u.id
        CROSS JOIN votes_summary vs
        WHERE r.id = $1
    `, id, userID).Scan(
		&review.ID, &review.UserID, &review.EntityType, &review.EntityID, &review.Rating,
		&comment, &pros, &cons, &photos, &review.LikesCount,
		&review.IsVerifiedPurchase, &review.Status, &review.CreatedAt, &review.UpdatedAt,
		&userName, &userEmail, &userPictureURL,
		&review.HelpfulVotes, &review.NotHelpfulVotes,
		&currentUserVote,
	)

	if err != nil {
		return nil, fmt.Errorf("error getting review: %w", err)
	}

	// Присваиваем значения из NullString только если они валидны
	if comment.Valid {
		review.Comment = comment.String
	}
	if pros.Valid {
		review.Pros = pros.String
	}
	if cons.Valid {
		review.Cons = cons.String
	}
	review.Photos = photos

	// Заполняем информацию о пользователе
	if userName.Valid {
		review.User.Name = userName.String
	}
	if userEmail.Valid {
		review.User.Email = userEmail.String
	}
	if userPictureURL.Valid {
		review.User.PictureURL = userPictureURL.String
	}

	// Устанавливаем current_user_vote только если значение не NULL
	if currentUserVote.Valid {
		review.CurrentUserVote = currentUserVote.String
	}

	// Инициализируем VotesCount
	review.VotesCount = struct {
		Helpful    int `json:"helpful"`
		NotHelpful int `json:"not_helpful"`
	}{
		Helpful:    review.HelpfulVotes,
		NotHelpful: review.NotHelpfulVotes,
	}

	// Загружаем ответы на отзыв
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
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Сначала проверим, существует ли уже такой голос
	var existingVoteType sql.NullString
	err = tx.QueryRow(ctx, `
        SELECT vote_type 
        FROM review_votes 
        WHERE review_id = $1 AND user_id = $2
    `, vote.ReviewID, vote.UserID).Scan(&existingVoteType)

	if err != nil && err != pgx.ErrNoRows {
		return fmt.Errorf("error checking existing vote: %w", err)
	}

	if existingVoteType.Valid && existingVoteType.String == vote.VoteType {
		// Если такой же голос уже есть - удаляем его (снимаем голос)
		_, err = tx.Exec(ctx, `
            DELETE FROM review_votes 
            WHERE review_id = $1 AND user_id = $2
        `, vote.ReviewID, vote.UserID)
	} else {
		// Иначе добавляем/обновляем голос
		_, err = tx.Exec(ctx, `
            INSERT INTO review_votes (review_id, user_id, vote_type)
            VALUES ($1, $2, $3)
            ON CONFLICT (review_id, user_id) 
            DO UPDATE SET vote_type = EXCLUDED.vote_type
        `, vote.ReviewID, vote.UserID, vote.VoteType)
	}

	if err != nil {
		return fmt.Errorf("error updating vote: %w", err)
	}

	// Обновляем счетчики
	_, err = tx.Exec(ctx, `
        UPDATE reviews r
        SET 
            helpful_votes = (
                SELECT COUNT(*) FROM review_votes 
                WHERE review_id = $1 AND vote_type = 'helpful'
            ),
            not_helpful_votes = (
                SELECT COUNT(*) FROM review_votes 
                WHERE review_id = $1 AND vote_type = 'not_helpful'
            )
        WHERE id = $1
    `, vote.ReviewID)

	if err != nil {
		return fmt.Errorf("error updating vote counts: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *Storage) UpdateReviewVotes(ctx context.Context, reviewId int) error {
	// Обновляем через транзакцию для атомарности
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Сначала получаем актуальное количество голосов
	var helpfulCount, notHelpfulCount int
	err = tx.QueryRow(ctx, `
        SELECT 
            COUNT(*) FILTER (WHERE vote_type = 'helpful'),
            COUNT(*) FILTER (WHERE vote_type = 'not_helpful')
        FROM review_votes
        WHERE review_id = $1
    `, reviewId).Scan(&helpfulCount, &notHelpfulCount)
	if err != nil {
		return fmt.Errorf("failed to count votes: %w", err)
	}

	// Обновляем количество голосов в таблице reviews
	_, err = tx.Exec(ctx, `
        UPDATE reviews
        SET 
            helpful_votes = $1,
            not_helpful_votes = $2,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $3
    `, helpfulCount, notHelpfulCount, reviewId)
	if err != nil {
		return fmt.Errorf("failed to update review votes: %w", err)
	}

	// Фиксируем транзакцию
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
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
