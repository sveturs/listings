package main

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
	"github.com/tealeg/xlsx"
)

func main() {
	// Подключение к БД
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:password@localhost:5432/svetubd?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("Failed to close database connection: %v", closeErr)
		}
	}()

	// Импортируем муниципалитеты
	if err := importMunicipalities(db); err != nil {
		log.Printf("Ошибка импорта муниципалитетов: %v", err)
	} else {
		log.Println("✅ Муниципалитеты импортированы")
	}

	// Импортируем населенные пункты
	if err := importPlaces(db); err != nil {
		log.Printf("Ошибка импорта населенных пунктов: %v", err)
	} else {
		log.Println("✅ Населенные пункты импортированы")
	}

	// Импортируем улицы
	if err := importStreets(db); err != nil {
		log.Printf("Ошибка импорта улиц: %v", err)
	} else {
		log.Println("✅ Улицы импортированы")
	}

	// Добавляем начальные настройки BEX
	if err := createDefaultSettings(db); err != nil {
		log.Printf("Ошибка создания настроек: %v", err)
	} else {
		log.Println("✅ Настройки BEX созданы")
	}

	log.Println("🎉 Импорт завершен успешно!")
}

func importMunicipalities(db *sql.DB) error {
	file, err := xlsx.OpenFile("/data/hostel-booking-system/data/bex-reference/Municipalities.xlsx")
	if err != nil {
		return fmt.Errorf("не удалось открыть файл муниципалитетов: %w", err)
	}

	sheet := file.Sheets[0]

	// Начинаем со второй строки (пропускаем заголовок)
	for i := 1; i < len(sheet.Rows); i++ {
		row := sheet.Rows[i]
		if len(row.Cells) < 2 {
			continue
		}

		bexID, _ := row.Cells[0].Int()
		name := row.Cells[1].String()

		if bexID == 0 || name == "" {
			continue
		}

		_, err := db.Exec(`
			INSERT INTO bex_municipalities (bex_id, name, name_cyrillic, is_active, created_at, updated_at)
			VALUES ($1, $2, $3, true, NOW(), NOW())
			ON CONFLICT (bex_id) DO UPDATE
			SET name = EXCLUDED.name,
			    updated_at = NOW()
		`, bexID, name, name)
		if err != nil {
			log.Printf("Ошибка вставки муниципалитета %d: %v", bexID, err)
		}
	}

	return nil
}

func importPlaces(db *sql.DB) error {
	file, err := xlsx.OpenFile("/data/hostel-booking-system/data/bex-reference/Places.xlsx")
	if err != nil {
		return fmt.Errorf("не удалось открыть файл населенных пунктов: %w", err)
	}

	sheet := file.Sheets[0]

	// Начинаем со второй строки (пропускаем заголовок)
	for i := 1; i < len(sheet.Rows); i++ {
		row := sheet.Rows[i]
		if len(row.Cells) < 4 {
			continue
		}

		bexID, _ := row.Cells[0].Int()
		name := row.Cells[1].String()
		postalCode := row.Cells[2].String()
		municipalityID, _ := row.Cells[3].Int()

		if bexID == 0 || name == "" {
			continue
		}

		// Сначала получаем ID муниципалитета из нашей БД
		var munID sql.NullInt64
		err := db.QueryRow("SELECT id FROM bex_municipalities WHERE bex_id = $1", municipalityID).Scan(&munID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			log.Printf("Ошибка поиска муниципалитета %d: %v", municipalityID, err)
			continue
		}

		_, err = db.Exec(`
			INSERT INTO bex_places (bex_id, name, name_cyrillic, postal_code, municipality_id, is_active, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, true, NOW(), NOW())
			ON CONFLICT (bex_id) DO UPDATE
			SET name = EXCLUDED.name,
			    postal_code = EXCLUDED.postal_code,
			    municipality_id = EXCLUDED.municipality_id,
			    updated_at = NOW()
		`, bexID, name, name, postalCode, munID)
		if err != nil {
			log.Printf("Ошибка вставки населенного пункта %d: %v", bexID, err)
		}
	}

	return nil
}

func importStreets(db *sql.DB) error {
	// Улиц очень много, используем CSV для более быстрой обработки
	file, err := os.Open("/data/hostel-booking-system/data/bex-reference/Streets.csv")
	if err != nil {
		// Попробуем Excel если CSV не существует
		return importStreetsFromExcel(db)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Printf("Failed to close file: %v", closeErr)
		}
	}()

	reader := csv.NewReader(file)

	// Пропускаем заголовок
	if _, err := reader.Read(); err != nil {
		return err
	}

	// Batch insert для производительности
	stmt, err := db.Prepare(`
		INSERT INTO bex_streets (bex_id, name, name_cyrillic, place_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, true, NOW(), NOW())
		ON CONFLICT (bex_id) DO UPDATE
		SET name = EXCLUDED.name,
		    place_id = EXCLUDED.place_id,
		    updated_at = NOW()
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	count := 0
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		if len(record) < 3 {
			continue
		}

		bexID, _ := strconv.Atoi(record[0])
		name := strings.TrimSpace(record[1])
		placeID, _ := strconv.Atoi(record[2])

		if bexID == 0 || name == "" {
			continue
		}

		// Получаем ID места из нашей БД
		var pID sql.NullInt64
		err = db.QueryRow("SELECT id FROM bex_places WHERE bex_id = $1", placeID).Scan(&pID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			continue
		}

		_, err = stmt.Exec(bexID, name, name, pID)
		if err != nil {
			log.Printf("Ошибка вставки улицы %d: %v", bexID, err)
		} else {
			count++
			if count%1000 == 0 {
				log.Printf("Импортировано %d улиц...", count)
			}
		}
	}

	log.Printf("Всего импортировано %d улиц", count)
	return nil
}

func importStreetsFromExcel(db *sql.DB) error {
	file, err := xlsx.OpenFile("/data/hostel-booking-system/data/bex-reference/Streets.xlsx")
	if err != nil {
		return fmt.Errorf("не удалось открыть файл улиц: %w", err)
	}

	sheet := file.Sheets[0]

	stmt, err := db.Prepare(`
		INSERT INTO bex_streets (bex_id, name, name_cyrillic, place_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, true, NOW(), NOW())
		ON CONFLICT (bex_id) DO UPDATE
		SET name = EXCLUDED.name,
		    place_id = EXCLUDED.place_id,
		    updated_at = NOW()
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	count := 0
	// Начинаем со второй строки (пропускаем заголовок)
	for i := 1; i < len(sheet.Rows); i++ {
		row := sheet.Rows[i]
		if len(row.Cells) < 3 {
			continue
		}

		bexID, _ := row.Cells[0].Int()
		name := row.Cells[1].String()
		placeID, _ := row.Cells[2].Int()

		if bexID == 0 || name == "" {
			continue
		}

		// Получаем ID места из нашей БД
		var pID sql.NullInt64
		err := db.QueryRow("SELECT id FROM bex_places WHERE bex_id = $1", placeID).Scan(&pID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			continue
		}

		_, err = stmt.Exec(bexID, name, name, pID)
		if err != nil {
			log.Printf("Ошибка вставки улицы %d: %v", bexID, err)
		} else {
			count++
			if count%100 == 0 {
				log.Printf("Импортировано %d улиц...", count)
			}
		}
	}

	log.Printf("Всего импортировано %d улиц", count)
	return nil
}

func createDefaultSettings(db *sql.DB) error {
	// Проверяем, существуют ли уже настройки
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM bex_settings").Scan(&count)
	if err == nil && count > 0 {
		log.Println("Настройки уже существуют")
		return nil
	}

	// Создаем настройки по умолчанию с предоставленными credentials
	_, err = db.Exec(`
		INSERT INTO bex_settings (
			auth_token, client_id, api_endpoint,
			sender_client_id, sender_name, sender_address,
			sender_city, sender_postal_code, sender_phone, sender_email,
			enabled, test_mode, use_address_lookup,
			created_at, updated_at
		) VALUES (
			$1, $2, $3,
			$4, $5, $6,
			$7, $8, $9, $10,
			true, true, true,
			NOW(), NOW()
		)
	`, "d50261-18wo-8539-ee5a-67uu3tu79", "326166", "https://api.bex.rs:62502",
		"326166", "Sve Tu d.o.o.", "Мике Манојловића 53",
		"Нови Сад", "21000", "+381 21 123456", "info@svetu.rs")

	return err
}
