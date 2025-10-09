// backend/cmd/diagnose-attributes/main.go
package main

import (
	//	"context"
	"database/sql"
	"fmt"
	"log"

	//	"os"

	"backend/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	// Загрузка .env файла
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load .env file")
	}

	// Загружаем конфигурацию
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Подключаемся к базе данных напрямую
	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Failed to close database connection: %v", err)
		}
	}()

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		log.Printf("Failed to ping database: %v", err)
		return
	}

	// Проверяем категории, которые должны иметь атрибуты
	categoriesToCheck := []int{1100, 2000, 3110, 3310, 3320, 3600, 3810, 3100}

	fmt.Println("=== Диагностика атрибутов категорий ===")

	// Проверка иерархии категорий
	fmt.Println("\n1. Проверка иерархии категорий:")
	for _, catID := range categoriesToCheck {
		checkCategoryHierarchy(db, catID)
	}

	// Проверка наличия связей категорий с атрибутами
	fmt.Println("\n2. Проверка связей категорий с атрибутами:")
	for _, catID := range categoriesToCheck {
		checkCategoryAttributes(db, catID)
	}

	// Выполнение того же запроса, что и в GetCategoryAttributes
	fmt.Println("\n3. Выполнение запроса GetCategoryAttributes для каждой категории:")
	for _, catID := range categoriesToCheck {
		runGetCategoryAttributesQuery(db, catID)
	}
}

func checkCategoryHierarchy(db *sql.DB, categoryID int) {
	var name string
	var parentID sql.NullInt64

	err := db.QueryRow(`
		SELECT name, parent_id 
		FROM c2c_categories 
		WHERE id = $1
	`, categoryID).Scan(&name, &parentID)
	if err != nil {
		fmt.Printf("Категория %d: ОШИБКА - %v\n", categoryID, err)
		return
	}

	if parentID.Valid {
		fmt.Printf("Категория %d: %s, родитель: %d\n", categoryID, name, parentID.Int64)
	} else {
		fmt.Printf("Категория %d: %s, родитель: нет\n", categoryID, name)
	}

	// Проверяем родительскую цепочку
	if parentID.Valid {
		var parents []int
		parent := int(parentID.Int64)

		for parent != 0 {
			parents = append(parents, parent)

			var nextParent sql.NullInt64
			err := db.QueryRow(`
				SELECT parent_id 
				FROM c2c_categories 
				WHERE id = $1
			`, parent).Scan(&nextParent)
			if err != nil {
				fmt.Printf("  Ошибка при получении родителя для %d: %v\n", parent, err)
				break
			}

			if !nextParent.Valid {
				break
			}

			parent = int(nextParent.Int64)
		}

		if len(parents) > 0 {
			fmt.Printf("  Цепочка родителей: %v\n", parents)
		}
	}
}

func checkCategoryAttributes(db *sql.DB, categoryID int) {
	rows, err := db.Query(`
		SELECT a.id, a.name, a.display_name, m.is_enabled
		FROM category_attribute_mapping m
		JOIN category_attributes a ON m.attribute_id = a.id
		WHERE m.category_id = $1
	`, categoryID)
	if err != nil {
		fmt.Printf("Категория %d: ОШИБКА при проверке атрибутов - %v\n", categoryID, err)
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Printf("Warning: failed to close rows: %v\n", err)
		}
	}()

	var count int
	for rows.Next() {
		var id int
		var name, displayName string
		var isEnabled bool

		err := rows.Scan(&id, &name, &displayName, &isEnabled)
		if err != nil {
			fmt.Printf("Категория %d: ОШИБКА при сканировании - %v\n", categoryID, err)
			continue
		}

		if count == 0 {
			fmt.Printf("Категория %d: найдены атрибуты:\n", categoryID)
		}

		fmt.Printf("  - %s (%s, id=%d, enabled=%v)\n", displayName, name, id, isEnabled)
		count++
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		fmt.Printf("Категория %d: ОШИБКА при итерации - %v\n", categoryID, err)
		return
	}

	if count == 0 {
		fmt.Printf("Категория %d: НЕТ атрибутов\n", categoryID)
	} else {
		fmt.Printf("Категория %d: всего атрибутов: %d\n", categoryID, count)
	}
}

func runGetCategoryAttributesQuery(db *sql.DB, categoryID int) {
	query := `
    WITH category_hierarchy AS (
        WITH RECURSIVE parents AS (
            SELECT id, parent_id
            FROM c2c_categories
            WHERE id = $1
            
            UNION
            
            SELECT c.id, c.parent_id
            FROM c2c_categories c
            INNER JOIN parents p ON c.id = p.parent_id
        )
        SELECT id FROM parents
    )
    SELECT 
        a.id, 
        a.name, 
        a.attribute_type
    FROM category_attribute_mapping m
    JOIN category_attributes a ON m.attribute_id = a.id
    JOIN category_hierarchy h ON m.category_id = h.id
    WHERE m.is_enabled = true
    ORDER BY a.sort_order, a.display_name
    `

	rows, err := db.Query(query, categoryID)
	if err != nil {
		fmt.Printf("Категория %d: ОШИБКА при выполнении запроса - %v\n", categoryID, err)
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Printf("Warning: failed to close rows: %v\n", err)
		}
	}()

	var count int
	for rows.Next() {
		var id int
		var name, attrType string

		err := rows.Scan(&id, &name, &attrType)
		if err != nil {
			fmt.Printf("Категория %d: ОШИБКА при сканировании - %v\n", categoryID, err)
			continue
		}

		if count == 0 {
			fmt.Printf("Категория %d: результаты запроса:\n", categoryID)
		}

		fmt.Printf("  - %s (id=%d, type=%s)\n", name, id, attrType)
		count++
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		fmt.Printf("Категория %d: ОШИБКА при итерации - %v\n", categoryID, err)
		return
	}

	if count == 0 {
		fmt.Printf("Категория %d: запрос НЕ вернул атрибутов\n", categoryID)

		// Диагностика каждого шага запроса
		diagnoseCategoryHierarchy(db, categoryID)
	} else {
		fmt.Printf("Категория %d: запрос вернул %d атрибутов\n", categoryID, count)
	}
}

func diagnoseCategoryHierarchy(db *sql.DB, categoryID int) {
	// Проверяем, что категория существует
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*) FROM c2c_categories WHERE id = $1
	`, categoryID).Scan(&count)
	if err != nil {
		fmt.Printf("  Диагностика: ОШИБКА при проверке существования категории - %v\n", err)
		return
	}

	if count == 0 {
		fmt.Printf("  Диагностика: Категория %d НЕ существует!\n", categoryID)
		return
	}

	// Проверяем построение иерархии
	rows, err := db.Query(`
		WITH RECURSIVE parents AS (
			SELECT id, parent_id, name
			FROM c2c_categories
			WHERE id = $1
			
			UNION
			
			SELECT c.id, c.parent_id, c.name
			FROM c2c_categories c
			INNER JOIN parents p ON c.id = p.parent_id
		)
		SELECT id, name FROM parents
	`, categoryID)
	if err != nil {
		fmt.Printf("  Диагностика: ОШИБКА при построении иерархии - %v\n", err)
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Printf("Warning: failed to close rows: %v\n", err)
		}
	}()

	fmt.Printf("  Диагностика: иерархия категорий:\n")
	var hierarchyCount int
	for rows.Next() {
		var id int
		var name string

		err := rows.Scan(&id, &name)
		if err != nil {
			fmt.Printf("    ОШИБКА при сканировании - %v\n", err)
			continue
		}

		fmt.Printf("    %d - %s\n", id, name)
		hierarchyCount++
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		fmt.Printf("  Диагностика: ОШИБКА при итерации иерархии - %v\n", err)
		return
	}

	if hierarchyCount == 0 {
		fmt.Printf("    Иерархия пуста (возможно, рекурсивный запрос работает некорректно)\n")
	}

	// Проверяем связи с атрибутами для категорий из иерархии
	rows, err = db.Query(`
		WITH RECURSIVE parents AS (
			SELECT id, parent_id
			FROM c2c_categories
			WHERE id = $1
			
			UNION
			
			SELECT c.id, c.parent_id
			FROM c2c_categories c
			INNER JOIN parents p ON c.id = p.parent_id
		)
		SELECT 
			p.id as category_id,
			COUNT(m.attribute_id) as attribute_count
		FROM parents p
		LEFT JOIN category_attribute_mapping m ON p.id = m.category_id
		GROUP BY p.id
	`, categoryID)
	if err != nil {
		fmt.Printf("  Диагностика: ОШИБКА при проверке связей - %v\n", err)
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Printf("Warning: failed to close rows: %v\n", err)
		}
	}()

	fmt.Printf("  Диагностика: связи с атрибутами:\n")
	var hasLinks bool
	for rows.Next() {
		var catID, attrCount int

		err := rows.Scan(&catID, &attrCount)
		if err != nil {
			fmt.Printf("    ОШИБКА при сканировании - %v\n", err)
			continue
		}

		fmt.Printf("    Категория %d: %d атрибутов\n", catID, attrCount)
		hasLinks = hasLinks || (attrCount > 0)
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		fmt.Printf("  Диагностика: ОШИБКА при итерации связей - %v\n", err)
		return
	}

	if !hasLinks {
		fmt.Printf("    ВНИМАНИЕ: Ни одна категория в иерархии не имеет связей с атрибутами!\n")
	}
}
