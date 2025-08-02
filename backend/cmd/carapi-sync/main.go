package main

import (
	"backend/internal/config"
	"backend/internal/integrations/carapi"
	"backend/internal/logger"
	"backend/internal/storage/postgres"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// SyncService handles all CarAPI synchronization
type SyncService struct {
	db       *sqlx.DB
	carAPI   *carapi.Client
	redis    *redis.Client
	logger   logger.Logger
}

func main() {
	// Инициализация логгера
	log := logger.New()
	log.Info().Msg("Starting CarAPI sync...")

	// Загрузка конфигурации
	cfg, err := config.Parse()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse config")
	}

	// Подключение к базе данных
	db, err := sqlx.Connect("postgres", cfg.Database.ConnectionString())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	// Подключение к Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	defer redisClient.Close()

	// Создаем клиент CarAPI
	// Токен нужно установить в переменной окружения CARAPI_TOKEN
	token := os.Getenv("CARAPI_TOKEN")
	if token == "" {
		log.Fatal().Msg("CARAPI_TOKEN environment variable is required")
	}

	carAPIClient := carapi.NewClient(token, redisClient)

	// Создаем сервис синхронизации
	syncService := &SyncService{
		db:       db,
		carAPI:   carAPIClient,
		redis:    redisClient,
		logger:   log,
	}

	// Выполняем полную синхронизацию
	ctx := context.Background()
	
	// 1. Синхронизируем все марки
	log.Info().Msg("Starting makes sync...")
	if err := syncService.SyncAllMakes(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to sync makes")
	}

	// 2. Синхронизируем модели для всех марок
	log.Info().Msg("Starting models sync...")
	if err := syncService.SyncAllModels(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to sync models")
	}

	// 3. Синхронизируем комплектации (trims) для популярных моделей
	log.Info().Msg("Starting trims sync...")
	if err := syncService.SyncPopularTrims(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to sync trims")
	}

	log.Info().Msg("CarAPI sync completed!")
}

// SyncAllMakes синхронизирует все марки автомобилей
func (s *SyncService) SyncAllMakes(ctx context.Context) error {
	// Получаем все марки из API
	apiMakes, err := s.carAPI.GetMakes(ctx)
	if err != nil {
		return fmt.Errorf("get makes from API: %w", err)
	}

	s.logger.Info().Int("count", len(apiMakes)).Msg("Fetched makes from CarAPI")

	// Начинаем транзакцию
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Обновляем или создаем марки
	for _, apiMake := range apiMakes {
		// Проверяем существует ли марка
		var existingID int
		err := tx.GetContext(ctx, &existingID, 
			"SELECT id FROM car_makes WHERE name = $1", 
			apiMake.Name)
		
		if err == nil {
			// Обновляем существующую марку
			_, err = tx.ExecContext(ctx, `
				UPDATE car_makes 
				SET external_id = $2, 
				    last_sync_at = NOW(),
				    metadata = $3
				WHERE id = $1`,
				existingID, 
				fmt.Sprintf("carapi_%d", apiMake.ID),
				map[string]interface{}{"carapi_id": apiMake.ID})
			
			if err != nil {
				s.logger.Error().Err(err).Str("make", apiMake.Name).Msg("Failed to update make")
			} else {
				s.logger.Debug().Str("make", apiMake.Name).Msg("Updated make")
			}
		} else {
			// Создаем новую марку
			slug := generateSlug(apiMake.Name)
			_, err = tx.ExecContext(ctx, `
				INSERT INTO car_makes (name, slug, external_id, last_sync_at, metadata)
				VALUES ($1, $2, $3, NOW(), $4)`,
				apiMake.Name, 
				slug,
				fmt.Sprintf("carapi_%d", apiMake.ID),
				map[string]interface{}{"carapi_id": apiMake.ID})
			
			if err != nil {
				s.logger.Error().Err(err).Str("make", apiMake.Name).Msg("Failed to create make")
			} else {
				s.logger.Info().Str("make", apiMake.Name).Msg("Created new make")
			}
		}

		// Логируем синхронизацию
		s.logSync(ctx, tx, "makes", "make", existingID, fmt.Sprintf("carapi_%d", apiMake.ID), "sync")
	}

	// Коммитим транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	// Обновляем счетчик использования API
	s.updateAPIUsage(ctx, len(apiMakes))

	return nil
}

// SyncAllModels синхронизирует модели для всех марок
func (s *SyncService) SyncAllModels(ctx context.Context) error {
	// Получаем все марки из БД
	var makes []struct {
		ID         int    `db:"id"`
		Name       string `db:"name"`
		ExternalID string `db:"external_id"`
	}
	
	err := s.db.SelectContext(ctx, &makes, 
		"SELECT id, name, external_id FROM car_makes WHERE external_id IS NOT NULL")
	if err != nil {
		return fmt.Errorf("get makes from DB: %w", err)
	}

	s.logger.Info().Int("count", len(makes)).Msg("Found makes to sync models")

	for _, make := range makes {
		// Извлекаем CarAPI ID из external_id
		var carapiID int
		fmt.Sscanf(make.ExternalID, "carapi_%d", &carapiID)
		
		// Получаем модели для марки
		apiModels, err := s.carAPI.GetModels(ctx, carapiID)
		if err != nil {
			s.logger.Error().Err(err).Str("make", make.Name).Msg("Failed to get models")
			continue
		}

		s.logger.Info().
			Str("make", make.Name).
			Int("models", len(apiModels)).
			Msg("Fetched models for make")

		// Сохраняем модели
		tx, _ := s.db.BeginTxx(ctx, nil)
		
		for _, apiModel := range apiModels {
			// Проверяем существует ли модель
			var existingID int
			err := tx.GetContext(ctx, &existingID,
				"SELECT id FROM car_models WHERE make_id = $1 AND name = $2",
				make.ID, apiModel.Name)
			
			if err == nil {
				// Обновляем
				_, err = tx.ExecContext(ctx, `
					UPDATE car_models 
					SET external_id = $2, 
					    last_sync_at = NOW(),
					    metadata = metadata || $3
					WHERE id = $1`,
					existingID,
					fmt.Sprintf("carapi_%d", apiModel.ID),
					map[string]interface{}{
						"carapi_id": apiModel.ID,
						"make_id": apiModel.MakeID,
					})
			} else {
				// Создаем
				slug := generateSlug(apiModel.Name)
				_, err = tx.ExecContext(ctx, `
					INSERT INTO car_models (make_id, name, slug, external_id, last_sync_at, metadata)
					VALUES ($1, $2, $3, $4, NOW(), $5)`,
					make.ID,
					apiModel.Name,
					slug,
					fmt.Sprintf("carapi_%d", apiModel.ID),
					map[string]interface{}{
						"carapi_id": apiModel.ID,
						"make_id": apiModel.MakeID,
					})
					
				if err == nil {
					s.logger.Info().
						Str("make", make.Name).
						Str("model", apiModel.Name).
						Msg("Created new model")
				}
			}
		}
		
		tx.Commit()
		
		// Обновляем счетчик API
		s.updateAPIUsage(ctx, len(apiModels))
		
		// Небольшая задержка чтобы не превысить rate limit
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}

// SyncPopularTrims синхронизирует комплектации для популярных моделей
func (s *SyncService) SyncPopularTrims(ctx context.Context) error {
	// Получаем популярные модели (например, последних 5 лет)
	currentYear := time.Now().Year()
	startYear := currentYear - 5

	// Получаем модели с external_id
	var models []struct {
		ID         int    `db:"id"`
		MakeID     int    `db:"make_id"`
		Name       string `db:"name"`
		MakeName   string `db:"make_name"`
		ExternalID string `db:"external_id"`
	}
	
	err := s.db.SelectContext(ctx, &models, `
		SELECT 
			m.id, 
			m.make_id,
			m.name, 
			mk.name as make_name,
			m.external_id
		FROM car_models m
		JOIN car_makes mk ON mk.id = m.make_id
		WHERE m.external_id IS NOT NULL
		AND mk.name IN ('Volkswagen', 'Toyota', 'Škoda', 'Fiat', 'Ford', 'Opel', 'Peugeot', 'Renault', 'Hyundai', 'Kia')
		ORDER BY mk.popularity_rs DESC, m.name
		LIMIT 100`)
	
	if err != nil {
		return fmt.Errorf("get popular models: %w", err)
	}

	s.logger.Info().Int("count", len(models)).Msg("Found popular models to sync trims")

	for _, model := range models {
		// Извлекаем CarAPI ID
		var carapiModelID int
		fmt.Sscanf(model.ExternalID, "carapi_%d", &carapiModelID)
		
		// Синхронизируем комплектации для последних лет
		for year := startYear; year <= currentYear; year++ {
			// Получаем комплектации для модели и года
			apiTrims, err := s.carAPI.GetTrims(ctx, carapiModelID, year)
			if err != nil {
				s.logger.Debug().
					Err(err).
					Str("model", model.Name).
					Int("year", year).
					Msg("No trims found for model/year")
				continue
			}

			if len(apiTrims) == 0 {
				continue
			}

			s.logger.Info().
				Str("make", model.MakeName).
				Str("model", model.Name).
				Int("year", year).
				Int("trims", len(apiTrims)).
				Msg("Found trims")

			// Сохраняем комплектации
			tx, _ := s.db.BeginTxx(ctx, nil)
			
			for _, apiTrim := range apiTrims {
				slug := generateSlug(fmt.Sprintf("%s-%d", apiTrim.Name, year))
				
				// Сохраняем полные данные в metadata
				metadata, _ := json.Marshal(apiTrim)
				
				_, err = tx.ExecContext(ctx, `
					INSERT INTO car_trims (
						model_id, year, name, slug, 
						external_id, carapi_trim_id, 
						metadata, last_sync_at
					) VALUES ($1, $2, $3, $4, $5, $6, $7::jsonb, NOW())
					ON CONFLICT (model_id, year, slug) DO UPDATE
					SET metadata = EXCLUDED.metadata,
					    last_sync_at = NOW()`,
					model.ID,
					year,
					apiTrim.Name,
					slug,
					fmt.Sprintf("carapi_%d", apiTrim.ID),
					apiTrim.ID,
					string(metadata))
				
				if err != nil {
					s.logger.Error().
						Err(err).
						Str("trim", apiTrim.Name).
						Msg("Failed to save trim")
				}
			}
			
			tx.Commit()
			
			// Обновляем счетчик API
			s.updateAPIUsage(ctx, 1)
			
			// Задержка для rate limit
			time.Sleep(time.Second)
		}
	}

	return nil
}

// Helper функции

func (s *SyncService) logSync(ctx context.Context, tx *sqlx.Tx, syncType, entityType string, entityID int, externalID, action string) {
	tx.ExecContext(ctx, `
		INSERT INTO car_sync_log (sync_type, entity_type, entity_id, external_id, action)
		VALUES ($1, $2, $3, $4, $5)`,
		syncType, entityType, entityID, externalID, action)
}

func (s *SyncService) updateAPIUsage(ctx context.Context, requests int) {
	s.db.ExecContext(ctx, `
		INSERT INTO carapi_usage (date, requests_count, last_request_at)
		VALUES (CURRENT_DATE, $1, NOW())
		ON CONFLICT (date) DO UPDATE
		SET requests_count = carapi_usage.requests_count + $1,
		    last_request_at = NOW()`,
		requests)
}

func generateSlug(name string) string {
	// Простая генерация slug - в реальности нужна более сложная логика
	return strings.ToLower(strings.ReplaceAll(name, " ", "-"))
}