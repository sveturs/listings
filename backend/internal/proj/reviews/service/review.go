// internal/services/review.go

package service

import (
	"context"
	"fmt"
	"log"

	"backend/internal/domain/models"
	"backend/internal/storage"

	"github.com/jackc/pgx/v5"
)

const (
	// Entity types
	entityTypeListing    = "listing"
	entityTypeUser       = "user"
	entityTypeStorefront = "storefront"

	// Rating trends
	trendStable = "stable"
)

type ReviewService struct {
	storage storage.Storage
}

func NewReviewService(storage storage.Storage) ReviewServiceInterface { //  ReviewService -> ReviewServiceInterface
	if storage == nil {
		log.Fatal("storage cannot be nil")
	}
	return &ReviewService{
		storage: storage,
	}
}

// UpdateEntityRatingInSearch обновляет рейтинг объекта в поисковом индексе
// UpdateEntityRatingInSearch обновляет рейтинг объекта в поисковом индексе
func (s *ReviewService) UpdateEntityRatingInSearch(ctx context.Context, entityType string, entityID int, avgRating float64) error {
	// Обновляем только для листингов
	if entityType != entityTypeListing {
		return nil
	}

	// Получаем текущий объект листинга из базы данных
	listing, err := s.storage.GetListingByID(ctx, entityID)
	if err != nil {
		return fmt.Errorf("ошибка получения объявления для обновления рейтинга: %w", err)
	}

	// Получаем статистику отзывов
	var reviewCount int
	var averageRating float64

	// Запрашиваем данные из базы
	err = s.storage.QueryRow(ctx, `
        SELECT COUNT(*), COALESCE(AVG(rating), 0)
        FROM reviews
        WHERE entity_type = $1 AND entity_id = $2 AND status = 'published'
    `, entityType, entityID).Scan(&reviewCount, &averageRating)
	if err != nil {
		// Если не удалось получить статистику, используем переданное значение
		averageRating = avgRating
		log.Printf("Ошибка получения статистики отзывов: %v, используем переданное значение рейтинга: %f", err, avgRating)
	}

	// Обновляем рейтинг в объекте листинга
	listing.AverageRating = averageRating
	listing.ReviewCount = reviewCount

	// Переиндексируем объявление с обновленным рейтингом
	err = s.storage.IndexListing(ctx, listing)
	if err != nil {
		return fmt.Errorf("ошибка переиндексации объявления с обновленным рейтингом: %w", err)
	}

	log.Printf("Успешно обновлен рейтинг в индексе для %s ID=%d: %.2f (%d отзывов)",
		entityType, entityID, averageRating, reviewCount)
	return nil
}

// CreateDraftReview создает черновик отзыва (этап 1)
func (s *ReviewService) CreateDraftReview(ctx context.Context, userId int, req *models.CreateReviewRequest) (*models.Review, error) {
	review := &models.Review{
		UserID:           userId,
		EntityType:       req.EntityType,
		EntityID:         req.EntityID,
		Rating:           req.Rating,
		Comment:          req.Comment,
		Pros:             req.Pros,
		Cons:             req.Cons,
		Photos:           nil,     // Фотографии добавятся позже
		Status:           "draft", // Статус черновика
		OriginalLanguage: req.OriginalLanguage,
	}

	// Заполняем entity_origin_type и entity_origin_id для агрегации рейтингов
	err := s.setEntityOrigin(ctx, review)
	if err != nil {
		log.Printf("Ошибка определения origin для отзыва: %v", err)
		// Не блокируем создание отзыва, но логируем ошибку
	}

	// Проверяем, является ли покупка верифицированной
	review.IsVerifiedPurchase = s.checkVerifiedPurchase(ctx, userId, req.EntityType, req.EntityID)

	// Создаем черновик отзыва в базе
	createdReview, err := s.storage.CreateReview(ctx, review)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания черновика отзыва: %w", err)
	}

	log.Printf("Created draft review with ID %d", createdReview.ID)
	return createdReview, nil
}

// PublishReview публикует черновик отзыва (этап 2)
func (s *ReviewService) PublishReview(ctx context.Context, reviewId int) (*models.Review, error) {
	// Получаем отзыв
	review, err := s.storage.GetReviewByID(ctx, reviewId)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения отзыва: %w", err)
	}

	// Проверяем, что отзыв в статусе draft
	if review.Status != "draft" {
		return nil, fmt.Errorf("отзыв не является черновиком")
	}

	// Обновляем статус на published
	err = s.storage.UpdateReviewStatus(ctx, reviewId, "published")
	if err != nil {
		return nil, fmt.Errorf("ошибка публикации отзыва: %w", err)
	}

	// Обновляем рейтинг в поисковом индексе
	err = s.UpdateEntityRatingInSearch(ctx, review.EntityType, review.EntityID, float64(review.Rating))
	if err != nil {
		log.Printf("Ошибка обновления рейтинга в поиске: %v", err)
		// Не блокируем публикацию отзыва
	}

	// Получаем обновленный отзыв
	publishedReview, err := s.storage.GetReviewByID(ctx, reviewId)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения опубликованного отзыва: %w", err)
	}

	log.Printf("Published review with ID %d", reviewId)
	return publishedReview, nil
}

// Исправленная версия метода CreateReview (legacy, одношаговое создание)
func (s *ReviewService) CreateReview(ctx context.Context, userId int, req *models.CreateReviewRequest) (*models.Review, error) {
	review := &models.Review{
		UserID:           userId,
		EntityType:       req.EntityType,
		EntityID:         req.EntityID,
		Rating:           req.Rating,
		Comment:          req.Comment,
		Pros:             req.Pros,
		Cons:             req.Cons,
		Photos:           req.Photos,
		Status:           "published",
		OriginalLanguage: req.OriginalLanguage,
	}

	// Заполняем entity_origin_type и entity_origin_id для агрегации рейтингов
	err := s.setEntityOrigin(ctx, review)
	if err != nil {
		log.Printf("Ошибка определения origin для отзыва: %v", err)
		// Не блокируем создание отзыва, но логируем ошибку
	}

	// Проверяем, является ли покупка верифицированной
	review.IsVerifiedPurchase = s.checkVerifiedPurchase(ctx, userId, req.EntityType, req.EntityID)

	// Создаем отзыв в базе
	createdReview, err := s.storage.CreateReview(ctx, review)
	if err != nil {
		return nil, err
	}

	// Используем ID из возвращенного отзыва
	// или просто возвращаем сам createdReview, если он уже содержит ID

	// После успешного создания отзыва
	// Обновляем рейтинг в поисковом индексе
	err = s.UpdateEntityRatingInSearch(ctx, review.EntityType, review.EntityID, float64(review.Rating))
	if err != nil {
		// Логируем ошибку, но не возвращаем ее, чтобы не блокировать создание отзыва
		log.Printf("Ошибка обновления рейтинга в индексе: %v", err)
	}

	return createdReview, nil
}

func (s *ReviewService) GetReviews(ctx context.Context, filter models.ReviewsFilter) ([]models.Review, int64, error) {
	return s.storage.GetReviews(ctx, filter)
}

func (s *ReviewService) GetReviewByID(ctx context.Context, id int) (*models.Review, error) {
	return s.storage.GetReviewByID(ctx, id)
}

func (s *ReviewService) UpdateReview(ctx context.Context, userId int, reviewId int, review *models.Review) error {
	// Проверяем, принадлежит ли отзыв пользователю
	existingReview, err := s.storage.GetReviewByID(ctx, reviewId)
	if err != nil {
		return err
	}

	if existingReview.UserID != userId {
		return fmt.Errorf("unauthorized to update this review")
	}

	review.ID = reviewId
	return s.storage.UpdateReview(ctx, review)
}

func (s *ReviewService) DeleteReview(ctx context.Context, userId int, reviewId int) error {
	// Проверяем, принадлежит ли отзыв пользователю
	review, err := s.storage.GetReviewByID(ctx, reviewId)
	if err != nil {
		return err
	}

	if review.UserID != userId {
		return fmt.Errorf("unauthorized to delete this review")
	}

	return s.storage.DeleteReview(ctx, reviewId)
}

func (s *ReviewService) VoteForReview(ctx context.Context, userId int, reviewId int, voteType string) error {
	vote := &models.ReviewVote{
		ReviewID: reviewId,
		UserID:   userId,
		VoteType: voteType,
	}
	return s.storage.AddReviewVote(ctx, vote)
}

func (s *ReviewService) AddResponse(ctx context.Context, userId int, reviewId int, responseText string) error {
	response := &models.ReviewResponse{
		ReviewID: reviewId,
		UserID:   userId,
		Response: responseText,
	}
	return s.storage.AddReviewResponse(ctx, response)
}

func (s *ReviewService) GetEntityRating(ctx context.Context, entityType string, entityId int) (float64, error) {
	return s.storage.GetEntityRating(ctx, entityType, entityId)
}

// checkVerifiedPurchase проверяет, совершал ли пользователь покупку
// Использует трехэтапную систему верификации через анализ чата
func (s *ReviewService) checkVerifiedPurchase(ctx context.Context, userId int, entityType string, entityId int) bool {
	// В зависимости от типа сущности проверяем наличие покупки/бронирования
	switch entityType {
	case entityTypeListing:
		// Для маркетплейса проверяем активность в чате
		// Сначала получаем информацию о листинге, чтобы узнать продавца
		listing, err := s.storage.GetListingByID(ctx, entityId)
		if err != nil || listing == nil {
			return false
		}

		// Получаем статистику чата между покупателем и продавцом
		stats, err := s.storage.GetChatActivityStats(ctx, userId, listing.UserID, entityId)
		if err != nil {
			return false
		}

		// Критерии автоматической верификации:
		// 1. Чат существует
		// 2. Минимум 5 сообщений от каждой стороны
		// 3. Чат не старше 30 дней
		// 4. Обе стороны были активны в чате
		if stats.ChatExists &&
			stats.BuyerMessages >= 5 &&
			stats.SellerMessages >= 5 &&
			stats.DaysSinceLastMsg <= 30 &&
			stats.TotalMessages >= 10 {
			return true
		}

		return false
	case "room":
		// Проверяем бронирования комнат
		// TODO: интегрировать с системой бронирований когда она будет реализована
		return false
	case "car":
		// Проверяем аренду автомобилей
		// TODO: интегрировать с системой аренды когда она будет реализована
		return false
	default:
		return false
	}
}

func (s *ReviewService) GetReviewStats(ctx context.Context, entityType string, entityId int) (*models.ReviewStats, error) {
	stats := &models.ReviewStats{
		RatingDistribution: make(map[int]int),
	}

	// Для пользователей и витрин используем материализованные представления
	switch entityType {
	case entityTypeUser:
		// Используем материализованное представление user_ratings
		err := s.storage.QueryRow(ctx, `
			SELECT 
				total_reviews,
				average_rating,
				verified_reviews,
				photo_reviews
			FROM user_ratings
			WHERE user_id = $1
		`, entityId).Scan(
			&stats.TotalReviews,
			&stats.AverageRating,
			&stats.VerifiedReviews,
			&stats.PhotoReviews,
		)
		if err != nil {
			if err == pgx.ErrNoRows {
				// Если нет записи в материализованном представлении, возвращаем пустую статистику
				return stats, nil
			}
			return nil, err
		}

		// Получаем распределение оценок из материализованного представления
		rows, err := s.storage.Query(ctx, `
			SELECT rating, count
			FROM user_rating_distribution
			WHERE user_id = $1
		`, entityId)
		if err != nil {
			return nil, err
		}
		defer func() {
			if err := rows.Close(); err != nil {
				// Логирование ошибки закрытия rows
			}
		}()

		for rows.Next() {
			var rating, count int
			if err := rows.Scan(&rating, &count); err != nil {
				return nil, err
			}
			stats.RatingDistribution[rating] = count
		}

		return stats, nil
	case entityTypeStorefront:
		// Используем материализованное представление storefront_ratings
		err := s.storage.QueryRow(ctx, `
			SELECT 
				total_reviews,
				average_rating,
				verified_reviews,
				photo_reviews
			FROM storefront_ratings
			WHERE storefront_id = $1
		`, entityId).Scan(
			&stats.TotalReviews,
			&stats.AverageRating,
			&stats.VerifiedReviews,
			&stats.PhotoReviews,
		)
		if err != nil {
			if err == pgx.ErrNoRows {
				return stats, nil
			}
			return nil, err
		}

		// Получаем распределение оценок
		rows, err := s.storage.Query(ctx, `
			SELECT rating, count
			FROM storefront_rating_distribution
			WHERE storefront_id = $1
		`, entityId)
		if err != nil {
			return nil, err
		}
		defer func() {
			if err := rows.Close(); err != nil {
				// Логирование ошибки закрытия rows
			}
		}()

		for rows.Next() {
			var rating, count int
			if err := rows.Scan(&rating, &count); err != nil {
				return nil, err
			}
			stats.RatingDistribution[rating] = count
		}

		return stats, nil
	default:
		// Для других типов сущностей (listing, room, car) используем прямой запрос
		err := s.storage.QueryRow(ctx, `
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
		rows, err := s.storage.Query(ctx, `
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
		defer func() {
			if err := rows.Close(); err != nil {
				// Логирование ошибки закрытия rows
			}
		}()

		for rows.Next() {
			var rating, count int
			if err := rows.Scan(&rating, &count); err != nil {
				return nil, err
			}
			stats.RatingDistribution[rating] = count
		}

		return stats, nil
	}
}

func (s *ReviewService) UpdateReviewPhotos(ctx context.Context, reviewId int, photoUrls []string) error {
	// Получаем текущий отзыв
	review, err := s.storage.GetReviewByID(ctx, reviewId)
	if err != nil {
		return err
	}

	// Добавляем новые фото к существующим
	review.Photos = append(review.Photos, photoUrls...)

	// Проверяем общее количество фото
	if len(review.Photos) > 5 {
		return fmt.Errorf("превышено максимальное количество фотографий (5)")
	}

	// Сохраняем изменения
	return s.storage.UpdateReview(ctx, review)
}

// GetUserReviews возвращает все отзывы, связанные с пользователем
func (s *ReviewService) GetUserReviews(ctx context.Context, userID int) ([]models.Review, error) {
	filter := models.ReviewsFilter{
		Status: "published",
		Limit:  1000, // достаточно большое число
		Page:   1,
	}

	// SQL-запрос для получения отзывов выполняется в хранилище
	reviews, err := s.storage.GetUserReviews(ctx, userID, filter)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

// GetStorefrontReviews возвращает все отзывы, связанные с витриной
func (s *ReviewService) GetStorefrontReviews(ctx context.Context, storefrontID int) ([]models.Review, error) {
	filter := models.ReviewsFilter{
		Status: "published",
		Limit:  1000, // достаточно большое число
		Page:   1,
	}

	// SQL-запрос для получения отзывов выполняется в хранилище
	reviews, err := s.storage.GetStorefrontReviews(ctx, storefrontID, filter)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

// GetUserRatingSummary возвращает сводную информацию о рейтинге пользователя
func (s *ReviewService) GetUserRatingSummary(ctx context.Context, userID int) (*models.UserRatingSummary, error) {
	summary, err := s.storage.GetUserRatingSummary(ctx, userID)
	if err != nil {
		return nil, err
	}

	return summary, nil
}

// GetStorefrontRatingSummary возвращает сводную информацию о рейтинге витрины
func (s *ReviewService) GetStorefrontRatingSummary(ctx context.Context, storefrontID int) (*models.StorefrontRatingSummary, error) {
	summary, err := s.storage.GetStorefrontRatingSummary(ctx, storefrontID)
	if err != nil {
		return nil, err
	}

	return summary, nil
}

// setEntityOrigin заполняет поля entity_origin_type и entity_origin_id
// для правильной агрегации рейтингов после удаления объектов
func (s *ReviewService) setEntityOrigin(ctx context.Context, review *models.Review) error {
	switch review.EntityType {
	case entityTypeListing:
		// Для отзывов на товары определяем владельца
		listing, err := s.storage.GetListingByID(ctx, review.EntityID)
		if err != nil {
			return fmt.Errorf("не удалось получить информацию о товаре: %w", err)
		}

		// Если товар принадлежит магазину, origin указывает на магазин
		if listing.StorefrontID != nil && *listing.StorefrontID > 0 {
			review.EntityOriginType = entityTypeStorefront
			review.EntityOriginID = *listing.StorefrontID
		} else {
			// Иначе origin указывает на продавца
			review.EntityOriginType = entityTypeUser
			review.EntityOriginID = listing.UserID
		}

	case entityTypeStorefront:
		// Для отзывов на магазин origin - сам магазин
		review.EntityOriginType = entityTypeStorefront
		review.EntityOriginID = review.EntityID

	case entityTypeUser:
		// Для отзывов на пользователя origin - сам пользователь
		review.EntityOriginType = entityTypeUser
		review.EntityOriginID = review.EntityID

	default:
		return fmt.Errorf("неизвестный тип сущности: %s", review.EntityType)
	}

	return nil
}

// GetUserAggregatedRating получает агрегированный рейтинг пользователя
func (s *ReviewService) GetUserAggregatedRating(ctx context.Context, userID int) (*models.AggregatedRating, error) {
	data, err := s.storage.GetUserAggregatedRating(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Преобразуем в общий формат AggregatedRating
	rating := &models.AggregatedRating{
		EntityType:   entityTypeUser,
		EntityID:     userID,
		Average:      data.AverageRating,
		TotalReviews: data.TotalReviews,
		Distribution: map[int]int{
			1: data.Rating1,
			2: data.Rating2,
			3: data.Rating3,
			4: data.Rating4,
			5: data.Rating5,
		},
		Breakdown: models.RatingBreakdown{
			Direct: models.BreakdownItem{
				Count:   data.DirectReviews,
				Average: 0, // TODO: рассчитать отдельно если нужно
			},
			Listings: models.BreakdownItem{
				Count:   data.ListingReviews,
				Average: 0,
			},
			Storefronts: models.BreakdownItem{
				Count:   data.StorefrontReviews,
				Average: 0,
			},
		},
		VerifiedPercentage: 0,
		RecentRating:       data.RecentRating,
		RecentReviews:      data.RecentReviews,
		LastReviewAt:       data.LastReviewAt,
	}

	// Рассчитываем процент верифицированных
	if data.TotalReviews > 0 {
		rating.VerifiedPercentage = (data.VerifiedReviews * 100) / data.TotalReviews
	}

	// Определяем тренд
	if data.RecentRating != nil && data.RecentReviews >= 5 {
		diff := *data.RecentRating - data.AverageRating
		if diff > 0.2 {
			rating.RecentTrend = "up"
		} else if diff < -0.2 {
			rating.RecentTrend = "down"
		} else {
			rating.RecentTrend = trendStable
		}
	} else {
		rating.RecentTrend = trendStable
	}

	return rating, nil
}

// GetStorefrontAggregatedRating получает агрегированный рейтинг магазина
func (s *ReviewService) GetStorefrontAggregatedRating(ctx context.Context, storefrontID int) (*models.AggregatedRating, error) {
	data, err := s.storage.GetStorefrontAggregatedRating(ctx, storefrontID)
	if err != nil {
		return nil, err
	}

	rating := &models.AggregatedRating{
		EntityType:   entityTypeStorefront,
		EntityID:     storefrontID,
		Average:      data.AverageRating,
		TotalReviews: data.TotalReviews,
		Distribution: map[int]int{
			1: data.Rating1,
			2: data.Rating2,
			3: data.Rating3,
			4: data.Rating4,
			5: data.Rating5,
		},
		Breakdown: models.RatingBreakdown{
			Direct: models.BreakdownItem{
				Count:   data.DirectReviews,
				Average: 0,
			},
			Listings: models.BreakdownItem{
				Count:   data.ListingReviews,
				Average: 0,
			},
		},
		VerifiedPercentage: 0,
		RecentRating:       data.RecentRating,
		RecentReviews:      data.RecentReviews,
		LastReviewAt:       data.LastReviewAt,
	}

	if data.TotalReviews > 0 {
		rating.VerifiedPercentage = (data.VerifiedReviews * 100) / data.TotalReviews
	}

	if data.RecentRating != nil && data.RecentReviews >= 5 {
		diff := *data.RecentRating - data.AverageRating
		if diff > 0.2 {
			rating.RecentTrend = "up"
		} else if diff < -0.2 {
			rating.RecentTrend = "down"
		} else {
			rating.RecentTrend = trendStable
		}
	} else {
		rating.RecentTrend = trendStable
	}

	return rating, nil
}

// ConfirmReview подтверждает отзыв продавцом
func (s *ReviewService) ConfirmReview(ctx context.Context, userID int, reviewID int, req *models.CreateReviewConfirmationRequest) error {
	// Получаем отзыв
	review, err := s.storage.GetReviewByID(ctx, reviewID)
	if err != nil {
		return err
	}

	// Проверяем права - только продавец/владелец может подтвердить
	// TODO: добавить более детальную проверку прав в зависимости от entity_type

	// Проверяем, не было ли уже подтверждения
	existing, err := s.storage.GetReviewConfirmation(ctx, reviewID)
	if err != nil {
		return err
	}
	if existing != nil {
		return fmt.Errorf("отзыв уже имеет подтверждение")
	}

	// Создаем подтверждение
	confirmation := &models.ReviewConfirmation{
		ReviewID:           reviewID,
		ConfirmedBy:        userID,
		ConfirmationStatus: req.Status,
		Notes:              req.Notes,
	}

	err = s.storage.CreateReviewConfirmation(ctx, confirmation)
	if err != nil {
		return err
	}

	// Обновляем флаг в reviews если подтверждено
	if req.Status == "confirmed" {
		review.SellerConfirmed = true
		err = s.storage.UpdateReview(ctx, review)
	}

	return err
}

// DisputeReview создает спор по отзыву
func (s *ReviewService) DisputeReview(ctx context.Context, userID int, reviewID int, req *models.CreateReviewDisputeRequest) error {
	// Получаем отзыв
	_, err := s.storage.GetReviewByID(ctx, reviewID)
	if err != nil {
		return err
	}

	// Проверяем, нет ли активного спора
	existing, err := s.storage.GetReviewDispute(ctx, reviewID)
	if err != nil {
		return err
	}
	if existing != nil {
		return fmt.Errorf("по этому отзыву уже есть активный спор")
	}

	// Создаем спор
	dispute := &models.ReviewDispute{
		ReviewID:           reviewID,
		DisputedBy:         userID,
		DisputeReason:      req.Reason,
		DisputeDescription: req.Description,
		Status:             "pending",
	}

	return s.storage.CreateReviewDispute(ctx, dispute)
}

// CanUserReviewEntity проверяет может ли пользователь оставить отзыв
func (s *ReviewService) CanUserReviewEntity(ctx context.Context, userID int, entityType string, entityID int) (*models.CanReviewResponse, error) {
	return s.storage.CanUserReviewEntity(ctx, userID, entityType, entityID)
}
