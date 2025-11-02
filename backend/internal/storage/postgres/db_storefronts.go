// backend/internal/storage/postgres/db_storefronts.go
package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"backend/internal/domain/models"

	"github.com/jackc/pgx/v5"
)

// stringPtr creates a pointer to a string
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// SearchStorefrontsOpenSearch removed - use microservice
// func (db *Database) SearchStorefrontsOpenSearch(ctx context.Context, params *storefrontOpenSearch.StorefrontSearchParams) (*storefrontOpenSearch.StorefrontSearchResult, error) {
//	return nil, fmt.Errorf("marketplace service removed - use microservice")
//
}

// IndexStorefront индексирует витрину в OpenSearch
func (db *Database) IndexStorefront(ctx context.Context, storefront *models.Storefront) error {
	if db.osStorefrontRepo == nil {
		return fmt.Errorf("OpenSearch для витрин не настроен")
	}
	return db.osStorefrontRepo.Index(ctx, storefront)
}

// DeleteStorefrontIndex удаляет витрину из индекса OpenSearch
func (db *Database) DeleteStorefrontIndex(ctx context.Context, storefrontID int) error {
	if db.osStorefrontRepo == nil {
		return fmt.Errorf("OpenSearch для витрин не настроен")
	}
	return db.osStorefrontRepo.Delete(ctx, storefrontID)
}

// ReindexAllStorefronts переиндексирует все витрины
func (db *Database) ReindexAllStorefronts(ctx context.Context) error {
	if db.osStorefrontRepo == nil {
		return fmt.Errorf("OpenSearch для витрин не настроен")
	}
	return db.osStorefrontRepo.ReindexAll(ctx)
}

// ReindexAllProducts переиндексирует все товары витрин
func (db *Database) ReindexAllProducts(ctx context.Context) error {
	if db.productSearchRepo == nil {
		return fmt.Errorf("OpenSearch для товаров витрин не настроен")
	}

	// Получаем все активные товары витрин из БД
	query := `
		SELECT
			p.id, p.storefront_id, p.name, p.description, p.price, p.currency,
			p.category_id, p.sku, p.barcode, p.stock_quantity, p.stock_status,
			p.is_active, p.attributes, p.view_count, p.sold_count,
			p.created_at, p.updated_at,
			p.has_individual_location, p.individual_address, p.individual_latitude,
			p.individual_longitude, p.location_privacy, p.show_on_map, p.has_variants,
			c.id, c.name, c.slug, c.icon, c.parent_id,
			s.name as storefront_name
		FROM b2c_products p
		LEFT JOIN c2c_categories c ON p.category_id = c.id
		LEFT JOIN b2c_stores s ON p.storefront_id = s.id
		WHERE p.is_active = true
		ORDER BY p.id
	`

	rows, err := db.pool.Query(ctx, query)
	if err != nil {
		return fmt.Errorf("ошибка получения товаров: %w", err)
	}
	defer rows.Close()

	var products []*models.StorefrontProduct
	for rows.Next() {
		product := &models.StorefrontProduct{}
		var category models.MarketplaceCategory
		var storefrontName sql.NullString
		var categoryID sql.NullInt32
		var categoryName, categorySlug, categoryIcon sql.NullString
		var categoryParentID sql.NullInt32
		var attributesJSON []byte

		err := rows.Scan(
			&product.ID, &product.StorefrontID, &product.Name, &product.Description,
			&product.Price, &product.Currency, &categoryID, &product.SKU, &product.Barcode,
			&product.StockQuantity, &product.StockStatus, &product.IsActive,
			&attributesJSON, &product.ViewCount, &product.SoldCount,
			&product.CreatedAt, &product.UpdatedAt,
			&product.HasIndividualLocation, &product.IndividualAddress,
			&product.IndividualLatitude, &product.IndividualLongitude,
			&product.LocationPrivacy, &product.ShowOnMap, &product.HasVariants,
			&category.ID, &categoryName, &categorySlug, &categoryIcon, &categoryParentID,
			&storefrontName,
		)
		if err != nil {
			log.Printf("Ошибка сканирования товара: %v", err)
			continue
		}

		// Десериализуем attributes из JSON
		if attributesJSON != nil {
			if unmarshalErr := json.Unmarshal(attributesJSON, &product.Attributes); unmarshalErr != nil {
				log.Printf("Ошибка десериализации атрибутов товара %d: %v", product.ID, unmarshalErr)
			}
		}

		// Заполняем категорию если есть
		if categoryID.Valid {
			category.ID = int(categoryID.Int32)
			product.CategoryID = category.ID
			if categoryName.Valid {
				category.Name = categoryName.String
			}
			if categorySlug.Valid {
				category.Slug = categorySlug.String
			}
			if categoryIcon.Valid {
				category.Icon = &categoryIcon.String
			}
			if categoryParentID.Valid {
				parentID := int(categoryParentID.Int32)
				category.ParentID = &parentID
			}
			product.Category = &category
		}

		// Загружаем варианты если есть
		if product.HasVariants {
			variants, err := db.GetProductVariants(ctx, product.ID)
			if err == nil {
				// Конвертируем []*StorefrontProductVariant в []StorefrontProductVariant
				product.Variants = make([]models.StorefrontProductVariant, len(variants))
				for i, v := range variants {
					product.Variants[i] = *v
				}
			}
		}

		// TODO: Реализовать GetProductImages когда понадобится
		// images, err := db.GetProductImages(ctx, []int{product.ID})
		// if err == nil {
		// 	product.Images = images
		// }

		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("ошибка итерации товаров: %w", err)
	}

	log.Printf("Получено %d товаров для индексации", len(products))

	// Индексируем все товары пакетами
	const batchSize = 100
	for i := 0; i < len(products); i += batchSize {
		end := i + batchSize
		if end > len(products) {
			end = len(products)
		}
		batch := products[i:end]

		if err := db.productSearchRepo.BulkIndexProducts(ctx, batch); err != nil {
			return fmt.Errorf("ошибка индексации пакета товаров: %w", err)
		}
		log.Printf("Проиндексировано %d/%d товаров", end, len(products))
	}

	log.Printf("Успешно проиндексировано %d товаров витрин", len(products))
	return nil
}

func (db *Database) GetB2CProductImages(ctx context.Context, productID int) ([]models.MarketplaceImage, error) {
	return db.marketplaceDB.GetB2CProductImages(ctx, productID)
}

func (db *Database) AddToFavorites(ctx context.Context, userID int, listingID int) error {
	return db.marketplaceDB.AddToFavorites(ctx, userID, listingID)
}

func (db *Database) RemoveFromFavorites(ctx context.Context, userID int, listingID int) error {
	return db.marketplaceDB.RemoveFromFavorites(ctx, userID, listingID)
}

func (db *Database) GetUserFavorites(ctx context.Context, userID int) ([]models.MarketplaceListing, error) {
	return db.marketplaceDB.GetUserFavorites(ctx, userID)
}

// Storefront favorites
func (db *Database) AddStorefrontToFavorites(ctx context.Context, userID int, productID int) error {
	return db.marketplaceDB.AddStorefrontToFavorites(ctx, userID, productID)
}

func (db *Database) RemoveStorefrontFromFavorites(ctx context.Context, userID int, productID int) error {
	return db.marketplaceDB.RemoveStorefrontFromFavorites(ctx, userID, productID)
}

func (db *Database) GetUserStorefrontFavorites(ctx context.Context, userID int) ([]models.MarketplaceListing, error) {
	return db.marketplaceDB.GetUserStorefrontFavorites(ctx, userID)
}

func (db *Database) IncrementViewsCount(ctx context.Context, id int) error {
	// Логируем вызов функции
	fmt.Printf("IncrementViewsCount called for listing %d\n", id)

	// Получаем ID пользователя из контекста
	var userID int
	if uid := ctx.Value("user_id"); uid != nil {
		if uidInt, ok := uid.(int); ok {
			userID = uidInt
		}
		fmt.Printf("User ID from context: %v (type: %T)\n", uid, uid)
	} else {
		fmt.Printf("No user_id in context\n")
	}

	// Для неавторизованных пользователей используем IP-адрес как идентификатор
	var userIdentifier string
	if ip := ctx.Value("ip_address"); ip != nil {
		if ipStr, ok := ip.(string); ok {
			userIdentifier = ipStr
		}
		fmt.Printf("IP address from context: %v (type: %T)\n", ip, ip)
	} else {
		fmt.Printf("No ip_address in context\n")
	}

	// Проверяем, есть ли уже запись о просмотре этого объявления данным пользователем
	// Если userID > 0, проверяем по user_id, иначе по ip_hash
	var viewExists bool
	var err error

	switch {
	case userID > 0:
		// Для авторизованных пользователей проверяем по ID
		err = db.pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM listing_views
				WHERE listing_id = $1 AND user_id = $2
			)
		`, id, userID).Scan(&viewExists)
	case userIdentifier != "":
		// Для неавторизованных пользователей проверяем строго по IP-адресу,
		// убедившись, что user_id IS NULL (чтобы не конфликтовать с ограничением уникальности)
		err = db.pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM listing_views
				WHERE listing_id = $1 AND ip_hash = $2 AND user_id IS NULL
			)
		`, id, userIdentifier).Scan(&viewExists)
	default:
		// Если нет ни ID пользователя, ни IP - считаем, что просмотр уже был (перестраховка)
		return nil
	}

	if err != nil {
		return err
	}

	// Если просмотра ещё не было, увеличиваем счетчик и добавляем запись о просмотре
	fmt.Printf("View exists check for listing %d: %v\n", id, viewExists)
	if !viewExists {
		fmt.Printf("View does not exist, incrementing counter for listing %d\n", id)
		// Начинаем транзакцию
		tx, err := db.pool.Begin(ctx)
		if err != nil {
			return err
		}
		defer func() {
			if err := tx.Rollback(ctx); err != nil {
				// Игнорируем ошибку если транзакция уже была завершена
				_ = err // Explicitly ignore error
			}
		}()

		// Увеличиваем счетчик просмотров
		// Сначала пробуем обновить в c2c_listings
		commandTag, err := tx.Exec(ctx, `
			UPDATE c2c_listings
			SET views_count = views_count + 1
			WHERE id = $1
		`, id)
		if err != nil {
			fmt.Printf("Error updating c2c_listings: %v\n", err)
			return err
		}

		rowsAffected := commandTag.RowsAffected()
		fmt.Printf("c2c_listings rows affected: %d\n", rowsAffected)

		// Если не нашли в c2c_listings, пробуем в storefront_products
		if rowsAffected == 0 {
			fmt.Printf("Trying to update storefront_products for id %d\n", id)
			commandTag, err = tx.Exec(ctx, `
				UPDATE storefront_products
				SET view_count = view_count + 1
				WHERE id = $1
			`, id)
			if err != nil {
				fmt.Printf("Error updating storefront_products: %v\n", err)
				return err
			}

			rowsAffected = commandTag.RowsAffected()
			fmt.Printf("storefront_products rows affected: %d\n", rowsAffected)
			if rowsAffected == 0 {
				// Если ни в одной таблице не нашли товар
				return fmt.Errorf("listing or product with id %d not found", id)
			}
		}

		// Добавляем запись о просмотре
		if userID > 0 {
			_, err = tx.Exec(ctx, `
				INSERT INTO listing_views (listing_id, user_id, view_time)
				VALUES ($1, $2, NOW())
			`, id, userID)
		} else {
			_, err = tx.Exec(ctx, `
				INSERT INTO listing_views (listing_id, ip_hash, view_time, user_id)
				VALUES ($1, $2, NOW(), NULL)
			`, id, userIdentifier)
		}
		if err != nil {
			return err
		}

		// Фиксируем транзакцию
		fmt.Printf("Committing transaction for listing %d view increment\n", id)
		err = tx.Commit(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("View increment committed successfully for listing %d\n", id)

		// После успешного обновления в PostgreSQL синхронизируем данные с OpenSearch
		if db.osMarketplaceRepo != nil && db.osClient != nil {
			// Получаем обновленное значение счетчика просмотров
			var viewsCount int
			// Сначала пробуем получить из c2c_listings
			err = db.pool.QueryRow(ctx, "SELECT views_count FROM c2c_listings WHERE id = $1", id).Scan(&viewsCount)
			if err != nil {
				// Если не нашли в c2c_listings, пробуем из storefront_products
				err = db.pool.QueryRow(ctx, "SELECT view_count FROM b2c_products WHERE id = $1", id).Scan(&viewsCount)
				if err != nil {
					log.Printf("Ошибка при получении обновленного счетчика просмотров: %v", err)
					// Не прерываем выполнение, так как главное - обновить в PostgreSQL
				} else {
					// Обновляем данные в OpenSearch
					go db.updateViewCountInOpenSearch(id, viewsCount) //nolint:contextcheck // фоновое обновление
				}
			} else {
				// Обновляем данные в OpenSearch
				go db.updateViewCountInOpenSearch(id, viewsCount) //nolint:contextcheck // фоновое обновление
			}
		}
	}

	return nil
}

// updateViewCountInOpenSearch обновляет счетчик просмотров в индексе OpenSearch
func (db *Database) CreateStorefront(ctx context.Context, userID int, dto *models.StorefrontCreateDTO) (*models.Storefront, error) {
	storefront := &models.Storefront{
		UserID:           userID,
		Slug:             dto.Slug,
		Name:             dto.Name,
		Description:      stringPtr(dto.Description),
		LogoURL:          stringPtr(""), // Будет заполнено после загрузки
		BannerURL:        stringPtr(""), // Будет заполнено после загрузки
		Theme:            dto.Theme,
		Phone:            stringPtr(dto.Phone),
		Email:            stringPtr(dto.Email),
		Website:          stringPtr(dto.Website),
		Address:          stringPtr(dto.Location.FullAddress),
		City:             stringPtr(dto.Location.City),
		PostalCode:       stringPtr(dto.Location.PostalCode),
		Country:          stringPtr(dto.Location.Country),
		Latitude:         &dto.Location.BuildingLat,
		Longitude:        &dto.Location.BuildingLng,
		Settings:         dto.Settings,
		SEOMeta:          dto.SEOMeta,
		IsActive:         true,
		SubscriptionPlan: "basic",
		CommissionRate:   0.05, // 5% по умолчанию
	}

	err := db.pool.QueryRow(ctx, `
		INSERT INTO b2c_stores (user_id, slug, name, description, logo_url, banner_url, theme,
			phone, email, website, address, city, postal_code, country, latitude, longitude,
			settings, seo_meta, is_active, subscription_plan, commission_rate)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
		RETURNING id, created_at, updated_at
	`, storefront.UserID, storefront.Slug, storefront.Name, storefront.Description,
		storefront.LogoURL, storefront.BannerURL, storefront.Theme, storefront.Phone,
		storefront.Email, storefront.Website, storefront.Address, storefront.City,
		storefront.PostalCode, storefront.Country, storefront.Latitude, storefront.Longitude,
		storefront.Settings, storefront.SEOMeta, storefront.IsActive, storefront.SubscriptionPlan,
		storefront.CommissionRate).Scan(&storefront.ID, &storefront.CreatedAt, &storefront.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Создаем запись в unified_geo для поддержки GIS поиска витрин
	if storefront.Latitude != nil && storefront.Longitude != nil {
		// Calculate geohash from coordinates
		geohashStr := fmt.Sprintf("%.6f,%.6f", *storefront.Latitude, *storefront.Longitude)

		addressComponents := map[string]interface{}{
			"city":        storefront.City,
			"postal_code": storefront.PostalCode,
			"country":     storefront.Country,
		}
		addressComponentsJSON, _ := json.Marshal(addressComponents)

		_, err = db.pool.Exec(ctx, `
			INSERT INTO unified_geo (
				source_type, source_id, location, geohash,
				formatted_address, address_components,
				geocoding_confidence, address_verified,
				input_method, location_privacy, blur_radius,
				is_precise
			) VALUES (
				'b2c_store', $1, ST_SetSRID(ST_MakePoint($2, $3), 4326), $4,
				$5, $6,
				0.9, true,
				'manual', 'exact', 0,
				true
			)
			ON CONFLICT (source_type, source_id)
			DO UPDATE SET
				location = EXCLUDED.location,
				geohash = EXCLUDED.geohash,
				formatted_address = EXCLUDED.formatted_address,
				address_components = EXCLUDED.address_components,
				updated_at = CURRENT_TIMESTAMP
		`, storefront.ID, *storefront.Longitude, *storefront.Latitude, geohashStr,
			storefront.Address, addressComponentsJSON)

		if err != nil {
			log.Printf("Error creating unified_geo entry for storefront %d: %v", storefront.ID, err)
			// Не прерываем создание витрины из-за ошибки с geo
		} else {
			// Обновляем materialized view после успешного создания geo записи
			_, err = db.pool.Exec(ctx, "SELECT refresh_map_items_cache()")
			if err != nil {
				log.Printf("Error refreshing map_items_cache: %v", err)
			}
		}
	}

	return storefront, nil
}

func (db *Database) GetUserStorefronts(ctx context.Context, userID int) ([]models.Storefront, error) {
	rows, err := db.pool.Query(ctx, `
		SELECT id, user_id, slug, name, description, logo_url, banner_url, theme,
			phone, email, website, address, city, postal_code, country, latitude, longitude,
			settings, seo_meta, is_active, is_verified, verification_date, rating, reviews_count,
			products_count, sales_count, views_count, subscription_plan, subscription_expires_at,
			commission_rate, ai_agent_enabled, ai_agent_config, live_shopping_enabled,
			group_buying_enabled, created_at, updated_at
		FROM b2c_stores
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var b2c_stores []models.Storefront
	for rows.Next() {
		var s models.Storefront
		err := rows.Scan(
			&s.ID, &s.UserID, &s.Slug, &s.Name, &s.Description, &s.LogoURL, &s.BannerURL, &s.Theme,
			&s.Phone, &s.Email, &s.Website, &s.Address, &s.City, &s.PostalCode, &s.Country,
			&s.Latitude, &s.Longitude, &s.Settings, &s.SEOMeta, &s.IsActive, &s.IsVerified,
			&s.VerificationDate, &s.Rating, &s.ReviewsCount, &s.ProductsCount, &s.SalesCount,
			&s.ViewsCount, &s.SubscriptionPlan, &s.SubscriptionExpiresAt, &s.CommissionRate,
			&s.AIAgentEnabled, &s.AIAgentConfig, &s.LiveShoppingEnabled, &s.GroupBuyingEnabled,
			&s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		b2c_stores = append(b2c_stores, s)
	}

	return b2c_stores, nil
}

func (db *Database) GetStorefrontByID(ctx context.Context, id int) (*models.Storefront, error) {
	var s models.Storefront
	var theme, settings, seoMeta, aiConfig json.RawMessage

	err := db.pool.QueryRow(ctx, `
		SELECT id, user_id, slug, name, description, logo_url, banner_url,
			COALESCE(theme, '{}')::jsonb,
			phone, email, website, address, city, postal_code, country, latitude, longitude,
			COALESCE(settings, '{}')::jsonb, COALESCE(seo_meta, '{}')::jsonb,
			is_active, is_verified, verification_date, rating, reviews_count,
			products_count, sales_count, views_count, subscription_plan, subscription_expires_at,
			commission_rate, ai_agent_enabled, COALESCE(ai_agent_config, '{}')::jsonb,
			live_shopping_enabled, group_buying_enabled, created_at, updated_at
		FROM b2c_stores
		WHERE id = $1
	`, id).Scan(
		&s.ID, &s.UserID, &s.Slug, &s.Name, &s.Description, &s.LogoURL, &s.BannerURL, &theme,
		&s.Phone, &s.Email, &s.Website, &s.Address, &s.City, &s.PostalCode, &s.Country,
		&s.Latitude, &s.Longitude, &settings, &seoMeta, &s.IsActive, &s.IsVerified,
		&s.VerificationDate, &s.Rating, &s.ReviewsCount, &s.ProductsCount, &s.SalesCount,
		&s.ViewsCount, &s.SubscriptionPlan, &s.SubscriptionExpiresAt, &s.CommissionRate,
		&s.AIAgentEnabled, &aiConfig, &s.LiveShoppingEnabled, &s.GroupBuyingEnabled,
		&s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrStorefrontNotFound
		}
		return nil, err
	}

	// Конвертируем json.RawMessage в JSONB
	if theme != nil {
		if err := json.Unmarshal(theme, &s.Theme); err != nil {
			// Логируем ошибку, но не прерываем выполнение
			_ = err // Explicitly ignore error
		}
	}
	if settings != nil {
		if err := json.Unmarshal(settings, &s.Settings); err != nil {
			// Логируем ошибку, но не прерываем выполнение
			_ = err // Explicitly ignore error
		}
	}
	if seoMeta != nil {
		if err := json.Unmarshal(seoMeta, &s.SEOMeta); err != nil {
			// Логируем ошибку, но не прерываем выполнение
			_ = err // Explicitly ignore error
		}
	}
	if aiConfig != nil {
		if err := json.Unmarshal(aiConfig, &s.AIAgentConfig); err != nil {
			// Логируем ошибку, но не прерываем выполнение
			_ = err // Explicitly ignore error
		}
	}

	return &s, nil
}

// GetStorefrontOwnerByProductID возвращает ID владельца витрины по ID товара
func (db *Database) GetStorefrontOwnerByProductID(ctx context.Context, productID int) (int, error) {
	var userID int
	err := db.pool.QueryRow(ctx, `
		SELECT s.user_id
		FROM b2c_stores s
		INNER JOIN storefront_products sp ON sp.storefront_id = s.id
		WHERE sp.id = $1
	`, productID).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("storefront product not found: %d", productID)
		}
		return 0, err
	}

	return userID, nil
}

func (db *Database) UpdateStorefront(ctx context.Context, storefront *models.Storefront) error {
	_, err := db.pool.Exec(ctx, `
		UPDATE b2c_stores
		SET name = $2, description = $3, logo_url = $4, banner_url = $5, theme = $6,
			phone = $7, email = $8, website = $9, address = $10, city = $11,
			postal_code = $12, country = $13, latitude = $14, longitude = $15,
			settings = $16, seo_meta = $17, is_active = $18, updated_at = NOW()
		WHERE id = $1
	`, storefront.ID, storefront.Name, storefront.Description, storefront.LogoURL,
		storefront.BannerURL, storefront.Theme, storefront.Phone, storefront.Email,
		storefront.Website, storefront.Address, storefront.City, storefront.PostalCode,
		storefront.Country, storefront.Latitude, storefront.Longitude, storefront.Settings,
		storefront.SEOMeta, storefront.IsActive)
	return err
}

func (db *Database) DeleteStorefront(ctx context.Context, id int) error {
	_, err := db.pool.Exec(ctx, "DELETE FROM b2c_stores WHERE id = $1", id)
	return err
}

// Storefront возвращает репозиторий витрин
func (db *Database) Storefront() interface{} {
	if db.storefrontRepo != nil {
		return db.storefrontRepo
	}
	// Возвращаем новый репозиторий используя текущий экземпляр db
	return NewStorefrontRepository(db)
}

func (db *Database) updateViewCountInOpenSearch(id int, viewsCount int) {
	ctx := context.Background()

	// Получаем объявление из PostgreSQL
	listing, err := db.GetListingByID(ctx, id)
	if err != nil {
		log.Printf("Ошибка при получении листинга для обновления OpenSearch: %v", err)
		return
	}

	// Устанавливаем значение счетчика просмотров
	listing.ViewsCount = viewsCount

	// Индексируем объявление в OpenSearch
	if db.osMarketplaceRepo != nil {
		err = db.osMarketplaceRepo.IndexListing(ctx, listing)
		if err != nil {
			log.Printf("Ошибка при обновлении счетчика просмотров в OpenSearch: %v", err)
		} else {
			log.Printf("Успешно обновлен счетчик просмотров в OpenSearch для объявления %d: %d", id, viewsCount)
		}
	}
}
