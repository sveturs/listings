package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Структуры для парсинга JSON
type CarAPIResponse struct {
	Collection struct {
		Count int `json:"count"`
		Total int `json:"total"`
	} `json:"collection"`
	Data json.RawMessage `json:"data"`
}

type CarAPIMake struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CarAPIModel struct {
	ID     int    `json:"id"`
	MakeID int    `json:"make_id"`
	Make   string `json:"make"`
	Name   string `json:"name"`
}

// Структуры БД
type DBMake struct {
	ID           int    `db:"id"`
	Name         string `db:"name"`
	Slug         string `db:"slug"`
	ExternalID   string `db:"external_id"`
	IsDomestic   bool   `db:"is_domestic"`
}

func main() {
	// Подключение к БД
	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Путь к данным
	dataDir := "/data/hostel-booking-system/backend/data/carapi-final-20250802-205650"

	// Импорт марок
	log.Println("=== Importing Makes ===")
	if err := importMakes(db, dataDir); err != nil {
		log.Printf("Error importing makes: %v", err)
	}

	// Импорт моделей
	log.Println("\n=== Importing Models ===")
	if err := importModels(db, dataDir); err != nil {
		log.Printf("Error importing models: %v", err)
	}

	// Добавление сербских марок
	log.Println("\n=== Adding Serbian Makes ===")
	if err := addSerbianMakes(db); err != nil {
		log.Printf("Error adding Serbian makes: %v", err)
	}

	// Статистика
	showStatistics(db)
}

func importMakes(db *sqlx.DB, dataDir string) error {
	// Читаем файл с марками
	data, err := ioutil.ReadFile(filepath.Join(dataDir, "makes", "all_makes.json"))
	if err != nil {
		return fmt.Errorf("read makes file: %w", err)
	}

	var response CarAPIResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}

	var makes []CarAPIMake
	if err := json.Unmarshal(response.Data, &makes); err != nil {
		return fmt.Errorf("unmarshal makes: %w", err)
	}

	log.Printf("Found %d makes to import", len(makes))

	// Начинаем транзакцию
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	imported := 0
	updated := 0
	skipped := 0

	for _, make := range makes {
		// Проверяем существует ли марка по external_id или имени
		var existingID int
		err := tx.Get(&existingID, `
			SELECT id FROM car_makes 
			WHERE external_id = $1 OR LOWER(name) = LOWER($2)
			LIMIT 1`,
			fmt.Sprintf("carapi_%d", make.ID), make.Name)

		if err == nil {
			// Обновляем существующую
			_, err = tx.Exec(`
				UPDATE car_makes 
				SET external_id = $2, 
				    last_sync_at = NOW(),
				    metadata = jsonb_build_object('carapi_id', $3::int)
				WHERE id = $1`,
				existingID, fmt.Sprintf("carapi_%d", make.ID), make.ID)
			
			if err != nil {
				log.Printf("Failed to update make %s: %v", make.Name, err)
			} else {
				updated++
			}
		} else {
			// Создаем новую
			slug := generateSlug(make.Name)
			
			// Проверяем уникальность slug
			var slugCount int
			tx.Get(&slugCount, "SELECT COUNT(*) FROM car_makes WHERE slug = $1", slug)
			if slugCount > 0 {
				slug = fmt.Sprintf("%s-%d", slug, make.ID)
			}

			_, err = tx.Exec(`
				INSERT INTO car_makes (name, slug, external_id, last_sync_at, metadata, is_domestic, popularity_rs)
				VALUES ($1, $2, $3, NOW(), jsonb_build_object('carapi_id', $4::int), false, 0)`,
				make.Name, slug, fmt.Sprintf("carapi_%d", make.ID), make.ID)

			if err != nil {
				log.Printf("Failed to insert make %s: %v", make.Name, err)
				skipped++
			} else {
				imported++
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	log.Printf("Makes: imported=%d, updated=%d, skipped=%d", imported, updated, skipped)
	return nil
}

func importModels(db *sqlx.DB, dataDir string) error {
	// Получаем все марки из БД для маппинга
	makeMap := make(map[string]int) // external_id -> db_id
	
	rows, err := db.Query("SELECT id, external_id FROM car_makes WHERE external_id IS NOT NULL")
	if err != nil {
		return fmt.Errorf("get makes: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var externalID string
		if err := rows.Scan(&id, &externalID); err == nil {
			makeMap[externalID] = id
		}
	}

	log.Printf("Loaded %d makes from database", len(makeMap))

	// Читаем все файлы с моделями
	modelsDir := filepath.Join(dataDir, "models")
	files, err := ioutil.ReadDir(modelsDir)
	if err != nil {
		return fmt.Errorf("read models dir: %w", err)
	}

	totalImported := 0
	totalUpdated := 0
	totalSkipped := 0

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		// Читаем файл
		data, err := ioutil.ReadFile(filepath.Join(modelsDir, file.Name()))
		if err != nil {
			log.Printf("Failed to read %s: %v", file.Name(), err)
			continue
		}

		var response CarAPIResponse
		if err := json.Unmarshal(data, &response); err != nil {
			log.Printf("Failed to unmarshal %s: %v", file.Name(), err)
			continue
		}

		var models []CarAPIModel
		if err := json.Unmarshal(response.Data, &models); err != nil {
			log.Printf("Failed to unmarshal models from %s: %v", file.Name(), err)
			continue
		}

		// Импортируем модели
		tx, _ := db.Beginx()
		
		for _, model := range models {
			// Находим make_id в нашей БД
			makeExternalID := fmt.Sprintf("carapi_%d", model.MakeID)
			dbMakeID, ok := makeMap[makeExternalID]
			if !ok {
				// log.Printf("Make not found for model %s (make_id=%d)", model.Name, model.MakeID)
				totalSkipped++
				continue
			}

			// Проверяем существует ли модель
			var existingID int
			err := tx.Get(&existingID, `
				SELECT id FROM car_models 
				WHERE make_id = $1 AND LOWER(name) = LOWER($2)
				LIMIT 1`,
				dbMakeID, model.Name)

			if err == nil {
				// Обновляем
				_, err = tx.Exec(`
					UPDATE car_models 
					SET external_id = $2, 
					    last_sync_at = NOW(),
					    metadata = COALESCE(metadata, '{}'::jsonb) || jsonb_build_object('carapi_id', $3::int)
					WHERE id = $1`,
					existingID, 
					fmt.Sprintf("carapi_%d", model.ID),
					model.ID)
				
				if err == nil {
					totalUpdated++
				}
			} else {
				// Создаем новую
				slug := generateSlug(model.Name)
				
				// Проверяем уникальность slug для этой марки
				var slugCount int
				tx.Get(&slugCount, "SELECT COUNT(*) FROM car_models WHERE make_id = $1 AND slug = $2", dbMakeID, slug)
				if slugCount > 0 {
					slug = fmt.Sprintf("%s-%d", slug, model.ID)
				}

				_, err = tx.Exec(`
					INSERT INTO car_models (make_id, name, slug, external_id, last_sync_at, metadata)
					VALUES ($1, $2, $3, $4, NOW(), jsonb_build_object('carapi_id', $5::int))`,
					dbMakeID, model.Name, slug,
					fmt.Sprintf("carapi_%d", model.ID),
					model.ID)

				if err == nil {
					totalImported++
				} else {
					log.Printf("Failed to insert model %s: %v", model.Name, err)
					totalSkipped++
				}
			}
		}
		
		tx.Commit()
	}

	log.Printf("Models: imported=%d, updated=%d, skipped=%d", totalImported, totalUpdated, totalSkipped)
	return nil
}

func addSerbianMakes(db *sqlx.DB) error {
	serbianMakes := []struct {
		Name         string
		Country      string
		IsDomestic   bool
		PopularityRS int
	}{
		{"Zastava", "Serbia", true, 100},
		{"Yugo", "Serbia", true, 90},
		{"FAP", "Serbia", true, 70},
		{"IMT", "Serbia", true, 60},
		{"IMK", "Serbia", true, 50},
		{"IDA-Opel", "Serbia", true, 40},
	}

	added := 0
	for _, make := range serbianMakes {
		// Проверяем существует ли
		var exists bool
		err := db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM car_makes WHERE LOWER(name) = LOWER($1))", make.Name)
		if err != nil || exists {
			continue
		}

		slug := generateSlug(make.Name)
		_, err = db.Exec(`
			INSERT INTO car_makes (name, slug, country, is_domestic, popularity_rs, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, NOW(), NOW())`,
			make.Name, slug, make.Country, make.IsDomestic, make.PopularityRS)

		if err != nil {
			log.Printf("Failed to add Serbian make %s: %v", make.Name, err)
		} else {
			log.Printf("Added Serbian make: %s", make.Name)
			added++
		}
	}

	log.Printf("Added %d Serbian makes", added)
	return nil
}

func showStatistics(db *sqlx.DB) {
	log.Println("\n=== Final Statistics ===")

	var stats struct {
		TotalMakes       int `db:"total_makes"`
		DomesticMakes    int `db:"domestic_makes"`
		MakesWithExtID   int `db:"makes_with_ext_id"`
		TotalModels      int `db:"total_models"`
		ModelsWithExtID  int `db:"models_with_ext_id"`
	}

	db.Get(&stats, `
		SELECT 
			(SELECT COUNT(*) FROM car_makes) as total_makes,
			(SELECT COUNT(*) FROM car_makes WHERE is_domestic = true) as domestic_makes,
			(SELECT COUNT(*) FROM car_makes WHERE external_id IS NOT NULL) as makes_with_ext_id,
			(SELECT COUNT(*) FROM car_models) as total_models,
			(SELECT COUNT(*) FROM car_models WHERE external_id IS NOT NULL) as models_with_ext_id
	`)

	log.Printf("Total makes: %d (domestic: %d, from CarAPI: %d)", 
		stats.TotalMakes, stats.DomesticMakes, stats.MakesWithExtID)
	log.Printf("Total models: %d (from CarAPI: %d)", 
		stats.TotalModels, stats.ModelsWithExtID)

	// Топ марок по количеству моделей
	log.Println("\nTop makes by model count:")
	rows, _ := db.Query(`
		SELECT m.name, COUNT(md.id) as model_count
		FROM car_makes m
		LEFT JOIN car_models md ON md.make_id = m.id
		GROUP BY m.id, m.name
		ORDER BY model_count DESC
		LIMIT 10
	`)
	defer rows.Close()

	for rows.Next() {
		var name string
		var count int
		if rows.Scan(&name, &count) == nil {
			log.Printf("  %s: %d models", name, count)
		}
	}
}

func generateSlug(name string) string {
	// Простая генерация slug
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, ".", "")
	slug = strings.ReplaceAll(slug, "&", "and")
	slug = strings.ReplaceAll(slug, "/", "-")
	slug = strings.ReplaceAll(slug, "'", "")
	return slug
}