// backend/internal/storage/postgres/auto_properties.go
package postgres

import (
	"backend/internal/domain/models"
	"context"
	"database/sql"
	"fmt"
	"log"
)

// CreateAutoProperties создает новую запись с автомобильными свойствами
func (db *Database) CreateAutoProperties(ctx context.Context, props *models.AutoProperties) error {
	query := `
		INSERT INTO auto_properties (
			listing_id, brand, model, year, mileage, fuel_type, transmission,
			engine_capacity, power, color, body_type, drive_type,
			number_of_doors, number_of_seats, additional_features
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		)
	`

	_, err := db.pool.Exec(ctx, query,
		props.ListingID, props.Brand, props.Model, props.Year, props.Mileage,
		props.FuelType, props.Transmission, props.EngineCapacity, props.Power,
		props.Color, props.BodyType, props.DriveType, props.NumberOfDoors,
		props.NumberOfSeats, props.AdditionalFeatures,
	)

	return err
}

// GetAutoPropertiesByListingID получает автомобильные свойства по ID объявления
func (db *Database) GetAutoPropertiesByListingID(ctx context.Context, listingID int) (*models.AutoProperties, error) {
	query := `
		SELECT 
			listing_id, brand, model, year, mileage, fuel_type, transmission,
			engine_capacity, power, color, body_type, drive_type,
			number_of_doors, number_of_seats, additional_features,
			created_at, updated_at
		FROM auto_properties
		WHERE listing_id = $1
	`

	props := &models.AutoProperties{}
	err := db.pool.QueryRow(ctx, query, listingID).Scan(
		&props.ListingID, &props.Brand, &props.Model, &props.Year, &props.Mileage,
		&props.FuelType, &props.Transmission, &props.EngineCapacity, &props.Power,
		&props.Color, &props.BodyType, &props.DriveType, &props.NumberOfDoors,
		&props.NumberOfSeats, &props.AdditionalFeatures, &props.CreatedAt, &props.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return props, nil
}

// UpdateAutoProperties обновляет автомобильные свойства
func (db *Database) UpdateAutoProperties(ctx context.Context, props *models.AutoProperties) error {
	query := `
		UPDATE auto_properties
		SET 
			brand = $2, model = $3, year = $4, mileage = $5,
			fuel_type = $6, transmission = $7, engine_capacity = $8,
			power = $9, color = $10, body_type = $11, drive_type = $12,
			number_of_doors = $13, number_of_seats = $14, additional_features = $15
		WHERE listing_id = $1
	`

	result, err := db.pool.Exec(ctx, query,
		props.ListingID, props.Brand, props.Model, props.Year, props.Mileage,
		props.FuelType, props.Transmission, props.EngineCapacity, props.Power,
		props.Color, props.BodyType, props.DriveType, props.NumberOfDoors,
		props.NumberOfSeats, props.AdditionalFeatures,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		// Если запись не найдена, создаем новую
		return db.CreateAutoProperties(ctx, props)
	}

	return nil
}

// DeleteAutoProperties удаляет автомобильные свойства
func (db *Database) DeleteAutoProperties(ctx context.Context, listingID int) error {
	query := `DELETE FROM auto_properties WHERE listing_id = $1`
	_, err := db.pool.Exec(ctx, query, listingID)
	return err
}

// GetAutoListings получает список автомобильных объявлений с фильтрацией
func (db *Database) GetAutoListings(ctx context.Context, filters map[string]string, limit int, offset int) ([]models.AutoListing, int64, error) {
	// Проверка и создание индекса OpenSearch при необходимости
	if db.osMarketplaceRepo != nil {
		if err := db.osMarketplaceRepo.PrepareIndex(ctx); err != nil {
			log.Printf("Ошибка подготовки индекса OpenSearch: %v", err)
			// Продолжаем выполнение - можем работать и без OpenSearch
		}
	}

	// Базовая часть запроса
	query := `
        SELECT 
            l.id, l.user_id, l.category_id, l.title, l.description,
            l.price, l.condition, l.status, l.location, l.latitude,
            l.longitude, l.address_city, l.address_country, l.views_count,
            l.created_at, l.updated_at, l.show_on_map, l.original_language,
            ap.brand, ap.model, ap.year, ap.mileage, ap.fuel_type, 
            ap.transmission, ap.engine_capacity, ap.power, ap.color, 
            ap.body_type, ap.drive_type, ap.number_of_doors, 
            ap.number_of_seats, ap.additional_features,
            u.name as user_name, u.email as user_email, 
            u.created_at as user_created_at, u.picture_url as user_picture_url,
            c.name as category_name, c.slug as category_slug,
            COUNT(*) OVER() as total_count
        FROM marketplace_listings l
        JOIN users u ON l.user_id = u.id
        JOIN marketplace_categories c ON l.category_id = c.id
        LEFT JOIN auto_properties ap ON l.id = ap.listing_id
        WHERE 1=1
    `

	// Добавляем условие, что это автомобильная категория
	query += `
        AND l.category_id IN (
            SELECT id FROM marketplace_categories 
            WHERE name LIKE '%автомобиль%' 
            OR id = 2000 
            OR parent_id = 2000
            OR parent_id IN (SELECT id FROM marketplace_categories WHERE parent_id = 2000)
        )
    `

	args := []interface{}{}
	argIndex := 1 // Начинаем счетчик аргументов с 1

	// Добавляем фильтры
	for key, value := range filters {
		if value == "" {
			continue
		}

		switch key {
		case "brand":
			query += fmt.Sprintf(" AND ap.brand = $%d", argIndex)
			args = append(args, value)
			argIndex++
		case "model":
			query += fmt.Sprintf(" AND ap.model = $%d", argIndex)
			args = append(args, value)
			argIndex++
		case "year_from":
			query += fmt.Sprintf(" AND ap.year >= $%d", argIndex)
			args = append(args, value)
			argIndex++
		case "year_to":
			query += fmt.Sprintf(" AND ap.year <= $%d", argIndex)
			args = append(args, value)
			argIndex++
		case "fuel_type":
			query += fmt.Sprintf(" AND ap.fuel_type = $%d", argIndex)
			args = append(args, value)
			argIndex++
		case "transmission":
			query += fmt.Sprintf(" AND ap.transmission = $%d", argIndex)
			args = append(args, value)
			argIndex++
		case "body_type":
			query += fmt.Sprintf(" AND ap.body_type = $%d", argIndex)
			args = append(args, value)
			argIndex++
		case "drive_type":
			query += fmt.Sprintf(" AND ap.drive_type = $%d", argIndex)
			args = append(args, value)
			argIndex++
		case "min_price":
			query += fmt.Sprintf(" AND l.price >= $%d", argIndex)
			args = append(args, value)
			argIndex++
		case "max_price":
			query += fmt.Sprintf(" AND l.price <= $%d", argIndex)
			args = append(args, value)
			argIndex++
		case "condition":
			query += fmt.Sprintf(" AND l.condition = $%d", argIndex)
			args = append(args, value)
			argIndex++
		case "city":
			query += fmt.Sprintf(" AND l.address_city = $%d", argIndex)
			args = append(args, value)
			argIndex++
		case "country":
			query += fmt.Sprintf(" AND l.address_country = $%d", argIndex)
			args = append(args, value)
			argIndex++
		case "query":
			query += fmt.Sprintf(" AND (LOWER(l.title) LIKE LOWER($%d) OR LOWER(l.description) LIKE LOWER($%d))", argIndex, argIndex)
			args = append(args, "%"+value+"%")
			argIndex++
		}
	}

	// Сортировка - добавляем только ОДИН раз ORDER BY и параметры пагинации
	query += " ORDER BY l.created_at DESC"

	// Добавляем параметры пагинации
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	// Логирование для отладки
	log.Printf("Выполнение SQL запроса:\nQuery: %s\nArgs: %v", query, args)

	// Выполняем запрос
	rows, err := db.pool.Query(ctx, query, args...)
	if err != nil {
		log.Printf("Ошибка выполнения запроса: %v", err)
		return nil, 0, fmt.Errorf("error querying auto listings: %w", err)
	}
	defer rows.Close()

	// Остальной код без изменений...
	var listings []models.AutoListing
	var totalCount int64

	for rows.Next() {
		var listing models.AutoListing
		listing.User = &models.User{}
		listing.Category = &models.MarketplaceCategory{}
		listing.AutoProperties = &models.AutoProperties{}

		// Объявляем временные переменные для NULL-значений
		var brand, model, fuelType, transmission, color, bodyType, driveType, additionalFeatures sql.NullString
		var year, mileage, power, numberOfDoors, numberOfSeats sql.NullInt32
		var engineCapacity sql.NullFloat64

		err := rows.Scan(
			&listing.ID, &listing.UserID, &listing.CategoryID, &listing.Title,
			&listing.Description, &listing.Price, &listing.Condition, &listing.Status,
			&listing.Location, &listing.Latitude, &listing.Longitude, &listing.City,
			&listing.Country, &listing.ViewsCount, &listing.CreatedAt, &listing.UpdatedAt,
			&listing.ShowOnMap, &listing.OriginalLanguage,
			&brand, &model, &year, &mileage, &fuelType,
			&transmission, &engineCapacity, &power, &color,
			&bodyType, &driveType, &numberOfDoors,
			&numberOfSeats, &additionalFeatures,
			&listing.User.Name, &listing.User.Email, &listing.User.CreatedAt,
			&listing.User.PictureURL, &listing.Category.Name, &listing.Category.Slug,
			&totalCount,
		)

		if err != nil {
			log.Printf("Ошибка сканирования результатов: %v", err)
			return nil, 0, fmt.Errorf("error scanning auto listing: %w", err)
		}

		// Копируем значения из временных переменных в структуру с проверкой на NULL
		if brand.Valid {
			listing.AutoProperties.Brand = brand.String
		}
		if model.Valid {
			listing.AutoProperties.Model = model.String
		}
		if year.Valid {
			listing.AutoProperties.Year = int(year.Int32)
		}
		if mileage.Valid {
			listing.AutoProperties.Mileage = int(mileage.Int32)
		}
		if fuelType.Valid {
			listing.AutoProperties.FuelType = fuelType.String
		}
		if transmission.Valid {
			listing.AutoProperties.Transmission = transmission.String
		}
		if engineCapacity.Valid {
			listing.AutoProperties.EngineCapacity = float64(engineCapacity.Float64)
		}
		if power.Valid {
			listing.AutoProperties.Power = int(power.Int32)
		}
		if color.Valid {
			listing.AutoProperties.Color = color.String
		}
		if bodyType.Valid {
			listing.AutoProperties.BodyType = bodyType.String
		}
		if driveType.Valid {
			listing.AutoProperties.DriveType = driveType.String
		}
		if numberOfDoors.Valid {
			listing.AutoProperties.NumberOfDoors = int(numberOfDoors.Int32)
		}
		if numberOfSeats.Valid {
			listing.AutoProperties.NumberOfSeats = int(numberOfSeats.Int32)
		}
		if additionalFeatures.Valid {
			listing.AutoProperties.AdditionalFeatures = additionalFeatures.String
		}

		// Получаем изображения напрямую через метод Database
		images, err := db.GetListingImages(ctx, fmt.Sprintf("%d", listing.ID))
		if err != nil {
			log.Printf("Ошибка получения изображений для объявления %d: %v", listing.ID, err)
			listing.Images = []models.MarketplaceImage{}
		} else {
			listing.Images = images
		}

		listings = append(listings, listing)
	}

	log.Printf("Найдено %d автомобильных объявлений", len(listings))
	return listings, totalCount, nil
}

func (db *Database) GetAutoListingByID(ctx context.Context, id int) (*models.AutoListing, error) {
	// GetAutoListingByID получает информацию об автомобильном объявлении по ID
	// Получаем базовое объявление
	baseListing, err := db.GetListingByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Проверяем, что это автомобильное объявление
	isAuto := false
	categoryID := baseListing.CategoryID

	// Проверка категории
	row := db.pool.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM marketplace_categories 
			WHERE (id = $1 AND parent_id = 2000) 
			OR id = 2000
			OR (id = $1 AND parent_id IN (SELECT id FROM marketplace_categories WHERE parent_id = 2000))
		)
	`, categoryID)

	err = row.Scan(&isAuto)
	if err != nil {
		return nil, err
	}

	// Если это не автомобильное объявление, возвращаем как обычное
	if !isAuto {
		return &models.AutoListing{
			MarketplaceListing: *baseListing,
			AutoProperties:     nil,
		}, nil
	}

	// Получаем автомобильные свойства
	autoProps, err := db.GetAutoPropertiesByListingID(ctx, id)
	if err != nil {
		// Если свойства не найдены, возвращаем только базовое объявление
		return &models.AutoListing{
			MarketplaceListing: *baseListing,
			AutoProperties:     nil,
		}, nil
	}

	// Возвращаем полное автомобильное объявление
	return &models.AutoListing{
		MarketplaceListing: *baseListing,
		AutoProperties:     autoProps,
	}, nil
}

// SearchAutoListings выполняет расширенный поиск автомобилей
func (db *Database) SearchAutoListings(ctx context.Context, autoFilters *models.AutoFilter, baseFilters map[string]string, limit int, offset int) ([]models.AutoListing, int64, error) {
	// Преобразуем расширенные фильтры в обычные
	filters := make(map[string]string)

	// Сначала копируем базовые фильтры
	for key, value := range baseFilters {
		filters[key] = value
	}

	// Добавляем автомобильные фильтры
	if autoFilters.Brand != "" {
		filters["brand"] = autoFilters.Brand
	}
	if autoFilters.Model != "" {
		filters["model"] = autoFilters.Model
	}
	if autoFilters.YearFrom > 0 {
		filters["year_from"] = fmt.Sprintf("%d", autoFilters.YearFrom)
	}
	if autoFilters.YearTo > 0 {
		filters["year_to"] = fmt.Sprintf("%d", autoFilters.YearTo)
	}
	if autoFilters.FuelType != "" {
		filters["fuel_type"] = autoFilters.FuelType
	}
	if autoFilters.Transmission != "" {
		filters["transmission"] = autoFilters.Transmission
	}
	if autoFilters.BodyType != "" {
		filters["body_type"] = autoFilters.BodyType
	}
	if autoFilters.DriveType != "" {
		filters["drive_type"] = autoFilters.DriveType
	}

	// Выполняем поиск
	return db.GetAutoListings(ctx, filters, limit, offset)
}

// Обновим метод ReindexAllListings для включения автосвойств в индекс
// Заменяем существующий метод ReindexAllListings на этот
func (db *Database) ReindexAllListings(ctx context.Context) error {
	// Проверяем, что opensearch клиент доступен
	if db.osMarketplaceRepo == nil {
		return fmt.Errorf("OpenSearch не настроен")
	}

	// Базовая переиндексация через стандартный метод
	if err := db.osMarketplaceRepo.ReindexAll(ctx); err != nil {
		return fmt.Errorf("ошибка базовой переиндексации: %w", err)
	}

	// Дополнительно обрабатываем автомобильные объявления
	query := `
        SELECT id FROM marketplace_listings
        WHERE category_id IN (
            SELECT id FROM marketplace_categories 
            WHERE parent_id = 2000 
            OR id = 2000 
            OR parent_id IN (SELECT id FROM marketplace_categories WHERE parent_id = 2000)
        )
    `

	rows, err := db.pool.Query(ctx, query)
	if err != nil {
		return fmt.Errorf("error querying auto listings for reindex: %w", err)
	}
	defer rows.Close()

	var autoListingIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Printf("Error scanning auto listing ID: %v", err)
			continue
		}
		autoListingIDs = append(autoListingIDs, id)
	}

	log.Printf("Found %d auto listings to update with auto properties", len(autoListingIDs))

	// Обновляем каждое автомобильное объявление
	for _, id := range autoListingIDs {
		autoListing, err := db.GetAutoListingByID(ctx, id)
		if err != nil {
			log.Printf("Error getting auto listing %d: %v", id, err)
			continue
		}

		if autoListing.AutoProperties != nil {
			// Переиндексируем объявление с дополнительными автосвойствами
			if err := db.IndexListing(ctx, &autoListing.MarketplaceListing); err != nil {
				log.Printf("Error indexing auto listing %d: %v", id, err)
			}
		}
	}

	return nil
}
