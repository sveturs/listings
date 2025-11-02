// backend/internal/storage/postgres/db_marketplace.go
package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"backend/internal/domain/models"
	"backend/internal/domain/search"

	"github.com/jackc/pgx/v5"
	authservice "github.com/sveturs/auth/pkg/http/service"
)

func (db *Database) GetListingImageByID(ctx context.Context, imageID int) (*models.MarketplaceImage, error) {
	var image models.MarketplaceImage
	var storageBucket, publicURL sql.NullString

	err := db.pool.QueryRow(ctx, `
		SELECT id, listing_id, file_path, file_name, file_size, content_type, is_main,
		       storage_type, storage_bucket, public_url, created_at
		FROM c2c_images
		WHERE id = $1
	`, imageID).Scan(
		&image.ID, &image.ListingID, &image.FilePath, &image.FileName, &image.FileSize,
		&image.ContentType, &image.IsMain, &image.StorageType, &storageBucket,
		&publicURL, &image.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("image not found")
		}
		return nil, err
	}

	// Handle nullable fields
	if storageBucket.Valid {
		image.StorageBucket = storageBucket.String
	}
	if publicURL.Valid {
		image.PublicURL = publicURL.String
	}

	return &image, nil
}

func (db *Database) DeleteListingImage(ctx context.Context, imageID int) error {
	_, err := db.pool.Exec(ctx, `
		DELETE FROM c2c_images
		WHERE id = $1
	`, imageID)

	return err
}

// IndexListing - TODO: OpenSearch integration disabled during refactoring
func (db *Database) IndexListing(ctx context.Context, listing *models.MarketplaceListing) error {
	log.Println("IndexListing: OpenSearch disabled during refactoring")
	return nil
}

// DeleteListingIndex - TODO: OpenSearch integration disabled during refactoring
func (db *Database) DeleteListingIndex(ctx context.Context, id string) error {
	log.Println("DeleteListingIndex: OpenSearch disabled during refactoring")
	return nil
}

// SuggestListings - TODO: OpenSearch integration disabled during refactoring
func (db *Database) SuggestListings(ctx context.Context, prefix string, size int) ([]string, error) {
	log.Println("SuggestListings: OpenSearch disabled during refactoring")
	return []string{}, nil
}

// ReindexAllListings - TODO: OpenSearch integration disabled during refactoring
func (db *Database) ReindexAllListings(ctx context.Context) error {
	log.Println("ReindexAllListings: OpenSearch disabled during refactoring")
	return nil
}

// GetCategoryAttributes получает атрибуты для указанной категории
// TODO: Реализация временно отключена, требуется рефакторинг
func (db *Database) GetCategoryAttributes(ctx context.Context, categoryID int) ([]models.CategoryAttribute, error) {
	return []models.CategoryAttribute{}, nil
}

// SaveListingAttributes сохраняет значения атрибутов для объявления
// TODO: Реализация временно отключена, требуется рефакторинг
func (db *Database) SaveListingAttributes(ctx context.Context, listingID int, attributes []models.ListingAttributeValue) error {
	return nil
}

// GetAttributeRanges - TODO: Реализация временно отключена, требуется рефакторинг
func (db *Database) GetAttributeRanges(ctx context.Context, categoryID int) (map[string]map[string]interface{}, error) {
	return make(map[string]map[string]interface{}), nil
}

// GetListingAttributes получает значения атрибутов для объявления
// TODO: Реализация временно отключена, требуется рефакторинг
func (db *Database) GetListingAttributes(ctx context.Context, listingID int) ([]models.ListingAttributeValue, error) {
	return []models.ListingAttributeValue{}, nil
}

// GetSession - DEPRECATED: Sessions are now managed via JWT tokens in auth-service
// This method is kept for backward compatibility but should not be used in new code

func (db *Database) GetFavoritedUsers(ctx context.Context, listingID int) ([]int, error) {
	query := `
        SELECT user_id
        FROM c2c_favorites
        WHERE listing_id = $1
    `
	rows, err := db.pool.Query(ctx, query, listingID)
	if err != nil {
		return nil, fmt.Errorf("error querying favorited users: %w", err)
	}
	defer rows.Close()

	var userIDs []int
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return nil, fmt.Errorf("error scanning user ID: %w", err)
		}
		userIDs = append(userIDs, userID)
	}

	return userIDs, nil
}

// CreateListing - методы для работы с listings находятся в отдельных файлах
// TODO: Консолидировать все методы listings в один файл после рефакторинга
func (db *Database) CreateListing(ctx context.Context, listing *models.MarketplaceListing) (int, error) {
	return 0, fmt.Errorf("CreateListing: method implementation removed, needs refactoring")
}

func (db *Database) GetListings(ctx context.Context, filters map[string]string, limit int, offset int) ([]models.MarketplaceListing, int64, error) {
	return nil, 0, fmt.Errorf("GetListings: method implementation removed, needs refactoring")
}

func (db *Database) GetListingByID(ctx context.Context, id int) (*models.MarketplaceListing, error) {
	return nil, fmt.Errorf("GetListingByID: method implementation removed, needs refactoring")
}

// GetMarketplaceListingsForReindex возвращает listings с needs_reindex=true для индексации в OpenSearch
func (db *Database) GetMarketplaceListingsForReindex(ctx context.Context, limit int) ([]*models.MarketplaceListing, error) {
	if limit <= 0 {
		limit = 1000 // default limit
	}

	query := `
		SELECT
			id, user_id, category_id, title, description, price,
			condition, status, location, latitude, longitude,
			address_city, address_country, views_count, show_on_map,
			original_language, created_at, updated_at,
			storefront_id, external_id, metadata,
			address_multilingual
		FROM c2c_listings
		WHERE needs_reindex = true
		  AND status = 'active'
		ORDER BY id
		LIMIT $1
	`

	rows, err := db.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query marketplace listings for reindex: %w", err)
	}
	defer rows.Close()

	listings := make([]*models.MarketplaceListing, 0)
	for rows.Next() {
		listing := &models.MarketplaceListing{}
		var addressMultilingual []byte
		var metadataJSON []byte
		var city sql.NullString
		var country sql.NullString

		err := rows.Scan(
			&listing.ID,
			&listing.UserID,
			&listing.CategoryID,
			&listing.Title,
			&listing.Description,
			&listing.Price,
			&listing.Condition,
			&listing.Status,
			&listing.Location,
			&listing.Latitude,
			&listing.Longitude,
			&city,
			&country,
			&listing.ViewsCount,
			&listing.ShowOnMap,
			&listing.OriginalLanguage,
			&listing.CreatedAt,
			&listing.UpdatedAt,
			&listing.StorefrontID,
			&listing.ExternalID,
			&metadataJSON,
			&addressMultilingual,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan marketplace listing: %w", err)
		}

		// Handle nullable fields
		if city.Valid {
			listing.City = city.String
		}
		if country.Valid {
			listing.Country = country.String
		}

		// Parse metadata JSON
		if len(metadataJSON) > 0 {
			if err := json.Unmarshal(metadataJSON, &listing.Metadata); err != nil {
				log.Printf("Warning: failed to unmarshal metadata for listing %d: %v", listing.ID, err)
			}
		}

		// Parse address_multilingual JSON
		if len(addressMultilingual) > 0 {
			if err := json.Unmarshal(addressMultilingual, &listing.AddressMultilingual); err != nil {
				log.Printf("Warning: failed to unmarshal address_multilingual for listing %d: %v", listing.ID, err)
			}
		}

		listings = append(listings, listing)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating marketplace listings: %w", err)
	}

	return listings, nil
}

// ResetMarketplaceListingsReindexFlag сбрасывает флаг needs_reindex для listings после успешной индексации
func (db *Database) ResetMarketplaceListingsReindexFlag(ctx context.Context, listingIDs []int) error {
	if len(listingIDs) == 0 {
		return nil
	}

	query := `
		UPDATE c2c_listings
		SET needs_reindex = false
		WHERE id = ANY($1)
	`

	_, err := db.pool.Exec(ctx, query, listingIDs)
	if err != nil {
		return fmt.Errorf("failed to reset needs_reindex flag: %w", err)
	}

	return nil
}

func (db *Database) GetListingBySlug(ctx context.Context, slug string) (*models.MarketplaceListing, error) {
	return nil, fmt.Errorf("method removed during refactoring, needs reimplementation") // OLD: db.marketplaceDB.GetListingBySlug(ctx, slug)
}

func (db *Database) IsSlugUnique(ctx context.Context, slug string, excludeID int) (bool, error) {
	return false, fmt.Errorf("method removed during refactoring, needs reimplementation") // OLD: db.marketplaceDB.IsSlugUnique(ctx, slug, excludeID)
}

func (db *Database) GenerateUniqueSlug(ctx context.Context, baseSlug string, excludeID int) (string, error) {
	return "", fmt.Errorf("method removed during refactoring, needs reimplementation") // OLD: db.marketplaceDB.GenerateUniqueSlug(ctx, baseSlug, excludeID)
}

func (db *Database) UpdateListing(ctx context.Context, listing *models.MarketplaceListing) error {
	return fmt.Errorf("method removed during refactoring, needs reimplementation") // OLD: db.marketplaceDB.UpdateListing(ctx, listing)
}

func (db *Database) DeleteListing(ctx context.Context, id int, userID int) error {
	return fmt.Errorf("method removed during refactoring, needs reimplementation") // OLD: db.marketplaceDB.DeleteListing(ctx, id, userID)
}

func (db *Database) DeleteListingAdmin(ctx context.Context, id int) error {
	return fmt.Errorf("method removed during refactoring, needs reimplementation") // OLD: db.marketplaceDB.DeleteListingAdmin(ctx, id)
}

func (db *Database) GetCategories(ctx context.Context) ([]models.MarketplaceCategory, error) {
	return nil, fmt.Errorf("method removed during refactoring, needs reimplementation") // OLD: db.marketplaceDB.GetCategories(ctx)
}

func (db *Database) GetAllCategories(ctx context.Context) ([]models.MarketplaceCategory, error) {
	return nil, fmt.Errorf("method removed during refactoring, needs reimplementation") // OLD: db.marketplaceDB.GetAllCategories(ctx)
}

func (db *Database) GetPopularCategories(ctx context.Context, limit int) ([]models.MarketplaceCategory, error) {
	return nil, fmt.Errorf("method removed during refactoring, needs reimplementation") // OLD: db.marketplaceDB.GetPopularCategories(ctx, limit)
}

func (db *Database) GetCategoryByID(ctx context.Context, id int) (*models.MarketplaceCategory, error) {
	return nil, fmt.Errorf("method removed during refactoring, needs reimplementation") // OLD: db.marketplaceDB.GetCategoryByID(ctx, id)
}

func (db *Database) GetCategoryTree(ctx context.Context) ([]models.CategoryTreeNode, error) {
	return nil, fmt.Errorf("method removed during refactoring, needs reimplementation") // OLD: db.marketplaceDB.GetCategoryTree(ctx)
}

func (db *Database) AddListingImage(ctx context.Context, image *models.MarketplaceImage) (int, error) {
	return 0, fmt.Errorf("method removed during refactoring, needs reimplementation") // OLD: db.marketplaceDB.AddListingImage(ctx, image)
}

func (db *Database) GetListingImages(ctx context.Context, listingID string) ([]models.MarketplaceImage, error) {
	return nil, fmt.Errorf("method removed during refactoring, needs reimplementation") // OLD: db.marketplaceDB.GetListingImages(ctx, listingID)
}

// Chat and messages methods
func (db *Database) ArchiveChat(ctx context.Context, chatID int, userID int) error {
	return fmt.Errorf("method removed during refactoring, needs reimplementation") // OLD: db.marketplaceDB.ArchiveChat(ctx, chatID, userID)
}

func (db *Database) UpdateMessageTranslations(ctx context.Context, messageID int, translations map[string]string) error {
	return fmt.Errorf("method removed during refactoring, needs reimplementation") // OLD: db.marketplaceDB.UpdateMessageTranslations(ctx, messageID, translations)
}

func (db *Database) CreateMessage(ctx context.Context, msg *models.MarketplaceMessage) error {
	return fmt.Errorf("method removed during refactoring, needs reimplementation") // OLD: db.marketplaceDB.CreateMessage(ctx, msg)
}

func (db *Database) GetMessages(ctx context.Context, listingID int, userID int, offset int, limit int) ([]models.MarketplaceMessage, error) {
	return nil, fmt.Errorf("method removed during refactoring, needs reimplementation") // OLD: db.marketplaceDB.GetMessages(ctx, listingID, userID, offset, limit)
}

func (db *Database) GetChats(ctx context.Context, userID int) ([]models.MarketplaceChat, error) {
	return nil, fmt.Errorf("method removed during refactoring, needs reimplementation") // OLD: db.marketplaceDB.GetChats(ctx, userID)
}

func (db *Database) GetChat(ctx context.Context, chatID int, userID int) (*models.MarketplaceChat, error) {
	return nil, fmt.Errorf("method removed during refactoring, needs reimplementation") // OLD: db.marketplaceDB.GetChat(ctx, chatID, userID)
}

func (db *Database) MarkMessagesAsRead(ctx context.Context, messageIDs []int, userID int) error {
	return fmt.Errorf("method removed during refactoring, needs reimplementation") // OLD: db.marketplaceDB.MarkMessagesAsRead(ctx, messageIDs, userID)
}

func (db *Database) GetUnreadMessagesCount(ctx context.Context, userID int) (int, error) {
	var count int
	err := db.pool.QueryRow(ctx, `
        SELECT COUNT(*)
        FROM c2c_messages m
        JOIN marketplace_chats c ON m.chat_id = c.id
        WHERE m.receiver_id = $1
        AND NOT m.is_read
        AND NOT c.is_archived
    `, userID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// SearchListings - TODO: OpenSearch integration disabled during refactoring
func (db *Database) SearchListings(ctx context.Context, params *search.SearchParams) (*search.SearchResult, error) {
	log.Println("SearchListings: OpenSearch disabled during refactoring")
	return &search.SearchResult{
		Listings: nil,
		Total:    0,
	}, nil
}

// Добавить этот метод в структуру Database
func (db *Database) GetAttributeOptionTranslations(ctx context.Context, attributeName, optionValue string) (map[string]string, error) {
	query := `
        SELECT option_value, en_translation, sr_translation
        FROM attribute_option_translations
        WHERE attribute_name = $1 AND option_value = $2
    `

	var optValue, enTrans, srTrans string
	err := db.pool.QueryRow(ctx, query, attributeName, optionValue).Scan(
		&optValue, &enTrans, &srTrans,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAttributeTranslationNotFound
		}
		return nil, fmt.Errorf("error getting attribute translations: %w", err)
	}

	translations := map[string]string{
		"en": enTrans,
		"sr": srTrans,
	}

	return translations, nil
}

// User Contacts methods - delegating to marketplace storage
func (db *Database) AddContact(ctx context.Context, contact *models.UserContact) error {
	return fmt.Errorf("method removed during refactoring, needs reimplementation") // OLD: db.marketplaceDB.AddContact(ctx, contact)
}

func (db *Database) UpdateContactStatus(ctx context.Context, userID, contactUserID int, status, notes string) error {
	return fmt.Errorf("method removed during refactoring, needs reimplementation") // OLD: db.marketplaceDB.UpdateContactStatus(ctx, userID, contactUserID, status, notes)
}

func (db *Database) GetContact(ctx context.Context, userID, contactUserID int) (*models.UserContact, error) {
	return nil, fmt.Errorf("method removed during refactoring, needs reimplementation") // OLD: db.marketplaceDB.GetContact(ctx, userID, contactUserID)
}

func (db *Database) GetUserContacts(ctx context.Context, userID int, status string, page, limit int) ([]models.UserContact, int, error) {
	return nil, 0, fmt.Errorf("method removed during refactoring, needs reimplementation") // OLD: db.marketplaceDB.GetUserContacts(ctx, userID, status, page, limit)
}

func (db *Database) GetIncomingContactRequests(ctx context.Context, userID int, page, limit int) ([]models.UserContact, int, error) {
	return nil, 0, fmt.Errorf("method removed during refactoring, needs reimplementation") // OLD: db.marketplaceDB.GetIncomingContactRequests(ctx, userID, page, limit)
}

func (db *Database) RemoveContact(ctx context.Context, userID, contactUserID int) error {
	return fmt.Errorf("method removed during refactoring, needs reimplementation") // OLD: db.marketplaceDB.RemoveContact(ctx, userID, contactUserID)
}

func (db *Database) GetUserPrivacySettings(ctx context.Context, userID int) (*models.UserPrivacySettings, error) {
	return nil, fmt.Errorf("method removed during refactoring, needs reimplementation") // OLD: db.marketplaceDB.GetUserPrivacySettings(ctx, userID)
}

func (db *Database) UpdateUserPrivacySettings(ctx context.Context, userID int, settings *models.UpdatePrivacySettingsRequest) error {
	return fmt.Errorf("method removed during refactoring, needs reimplementation") // OLD: db.marketplaceDB.UpdateUserPrivacySettings(ctx, userID, settings)
}

func (db *Database) GetPrivacySettings(ctx context.Context, userID int) (*models.UserPrivacySettings, error) {
	return nil, fmt.Errorf("method removed during refactoring, needs reimplementation") // OLD: db.marketplaceDB.GetUserPrivacySettings(ctx, userID)
}

func (db *Database) UpdatePrivacySettings(ctx context.Context, userID int, settings *models.UpdatePrivacySettingsRequest) error {
	return fmt.Errorf("method removed during refactoring, needs reimplementation") // OLD: db.marketplaceDB.UpdateUserPrivacySettings(ctx, userID, settings)
}

func (db *Database) UpdateChatSettings(ctx context.Context, userID int, settings *models.ChatUserSettings) error {
	return fmt.Errorf("method removed during refactoring, needs reimplementation") // OLD: db.marketplaceDB.UpdateChatSettings(ctx, userID, settings)
}

func (db *Database) CanAddContact(ctx context.Context, userID, targetUserID int) (bool, error) {
	return false, fmt.Errorf("method removed during refactoring, needs reimplementation") // OLD: db.marketplaceDB.CanAddContact(ctx, userID, targetUserID)
}

func (db *Database) AreContacts(ctx context.Context, userID1, userID2 int) (bool, error) {
	return false, fmt.Errorf("method removed during refactoring, needs reimplementation") // OLD: db.marketplaceDB.AreContacts(ctx, userID1, userID2)
}

// Marketplace listing variants methods
func (db *Database) CreateListingVariants(ctx context.Context, listingID int, variants []models.MarketplaceListingVariant) error {
	if len(variants) == 0 {
		return nil
	}

	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			log.Printf("Failed to rollback transaction: %v", rollbackErr)
		}
	}()

	for _, variant := range variants {
		attributesJSON, err := json.Marshal(variant.Attributes)
		if err != nil {
			return fmt.Errorf("failed to marshal variant attributes: %w", err)
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO marketplace_listing_variants (listing_id, sku, price, stock, attributes, image_url, is_active)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, listingID, variant.SKU, variant.Price, variant.Stock, attributesJSON, variant.ImageURL, variant.IsActive)
		if err != nil {
			return fmt.Errorf("failed to insert variant: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (db *Database) GetListingVariants(ctx context.Context, listingID int) ([]models.MarketplaceListingVariant, error) {
	query := `
		SELECT id, listing_id, sku, price, stock, attributes, image_url, is_active,
		       created_at::text, updated_at::text
		FROM c2c_listing_variants
		WHERE listing_id = $1 AND is_active = true
		ORDER BY id
	`

	rows, err := db.pool.Query(ctx, query, listingID)
	if err != nil {
		return nil, fmt.Errorf("failed to query variants: %w", err)
	}
	defer rows.Close()

	var variants []models.MarketplaceListingVariant
	for rows.Next() {
		var variant models.MarketplaceListingVariant
		var attributesJSON []byte

		err := rows.Scan(
			&variant.ID, &variant.ListingID, &variant.SKU, &variant.Price, &variant.Stock,
			&attributesJSON, &variant.ImageURL, &variant.IsActive,
			&variant.CreatedAt, &variant.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan variant: %w", err)
		}

		if len(attributesJSON) > 0 {
			err = json.Unmarshal(attributesJSON, &variant.Attributes)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal variant attributes: %w", err)
			}
		}

		variants = append(variants, variant)
	}

	return variants, rows.Err()
}

func (db *Database) UpdateListingVariant(ctx context.Context, variant *models.MarketplaceListingVariant) error {
	attributesJSON, err := json.Marshal(variant.Attributes)
	if err != nil {
		return fmt.Errorf("failed to marshal variant attributes: %w", err)
	}

	query := `
		UPDATE marketplace_listing_variants
		SET sku = $1, price = $2, stock = $3, attributes = $4, image_url = $5, is_active = $6
		WHERE id = $7
	`

	result, err := db.pool.Exec(ctx, query,
		variant.SKU, variant.Price, variant.Stock, attributesJSON,
		variant.ImageURL, variant.IsActive, variant.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update variant: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("variant not found")
	}

	return nil
}

func (db *Database) DeleteListingVariant(ctx context.Context, variantID int) error {
	// Soft delete - просто помечаем как неактивный
	query := `UPDATE marketplace_listing_variants SET is_active = false WHERE id = $1`

	result, err := db.pool.Exec(ctx, query, variantID)
	if err != nil {
		return fmt.Errorf("failed to delete variant: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("variant not found")
	}

	return nil
}

// SetMarketplaceUserService устанавливает UserService для marketplace storage
func (db *Database) SetMarketplaceUserService(userService *authservice.UserService) {
}

// GetAllUnifiedAttributes получает все активные unified attributes
func (db *Database) GetAllUnifiedAttributes(ctx context.Context) ([]*models.UnifiedAttribute, error) {
	query := `
		SELECT
			id, code, name, display_name, attribute_type, purpose,
			options, validation_rules, ui_settings,
			is_searchable, is_filterable, is_required,
			is_variant_compatible, affects_stock, affects_price,
			sort_order, is_active, created_at, updated_at,
			legacy_category_attribute_id, legacy_product_variant_attribute_id
		FROM unified_attributes
		WHERE is_active = true
		ORDER BY code`

	rows, err := db.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all unified attributes: %w", err)
	}
	defer rows.Close()

	var attributes []*models.UnifiedAttribute
	for rows.Next() {
		attr := &models.UnifiedAttribute{}
		err := rows.Scan(
			&attr.ID, &attr.Code, &attr.Name, &attr.DisplayName,
			&attr.AttributeType, &attr.Purpose,
			&attr.Options, &attr.ValidationRules, &attr.UISettings,
			&attr.IsSearchable, &attr.IsFilterable, &attr.IsRequired,
			&attr.IsVariantCompatible, &attr.AffectsStock, &attr.AffectsPrice,
			&attr.SortOrder, &attr.IsActive, &attr.CreatedAt, &attr.UpdatedAt,
			&attr.LegacyCategoryAttributeID, &attr.LegacyProductVariantAttributeID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan attribute: %w", err)
		}
		attributes = append(attributes, attr)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return attributes, nil
}

func (db *Database) GetChatActivityStats(ctx context.Context, buyerID int, sellerID int, listingID int) (*models.ChatActivityStats, error) {
	stats := &models.ChatActivityStats{}

	// Проверяем наличие чата и получаем статистику
	query := `
		WITH chat_info AS (
			SELECT
				c.id as chat_id,
				c.created_at as chat_created
			FROM c2c_chats c
			WHERE c.buyer_id = $1
				AND c.seller_id = $2
				AND c.listing_id = $3
			LIMIT 1
		),
		message_stats AS (
			SELECT
				COUNT(*) as total_messages,
				COUNT(*) FILTER (WHERE m.sender_id = $1) as buyer_messages,
				COUNT(*) FILTER (WHERE m.sender_id = $2) as seller_messages,
				MIN(m.created_at) as first_message_date,
				MAX(m.created_at) as last_message_date
			FROM c2c_messages m
			INNER JOIN chat_info ci ON m.chat_id = ci.chat_id
		)
		SELECT
			CASE WHEN ci.chat_id IS NOT NULL THEN true ELSE false END as chat_exists,
			COALESCE(ms.total_messages, 0) as total_messages,
			COALESCE(ms.buyer_messages, 0) as buyer_messages,
			COALESCE(ms.seller_messages, 0) as seller_messages,
			ms.first_message_date,
			ms.last_message_date,
			CASE
				WHEN ms.first_message_date IS NOT NULL
				THEN EXTRACT(DAY FROM NOW() - ms.first_message_date)::int
				ELSE 0
			END as days_since_first_msg,
			CASE
				WHEN ms.last_message_date IS NOT NULL
				THEN EXTRACT(DAY FROM NOW() - ms.last_message_date)::int
				ELSE 0
			END as days_since_last_msg
		FROM chat_info ci
		LEFT JOIN message_stats ms ON true
	`

	row := db.pool.QueryRow(ctx, query, buyerID, sellerID, listingID)

	err := row.Scan(
		&stats.ChatExists,
		&stats.TotalMessages,
		&stats.BuyerMessages,
		&stats.SellerMessages,
		&stats.FirstMessageDate,
		&stats.LastMessageDate,
		&stats.DaysSinceFirstMsg,
		&stats.DaysSinceLastMsg,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		// Чат не существует, возвращаем пустую статистику
		return stats, nil
	}

	return stats, err
}

func (db *Database) SynchronizeDiscountMetadata(ctx context.Context) error {
	// Получаем все объявления с информацией о скидке
	query := `
        SELECT id, price, metadata
        FROM c2c_listings
        WHERE metadata->>'discount' IS NOT NULL
    `

	rows, err := db.pool.Query(ctx, query)
	if err != nil {
		return fmt.Errorf("error querying listings with discounts: %w", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id int
		var price float64
		var metadataJSON []byte

		if err := rows.Scan(&id, &price, &metadataJSON); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		var metadata map[string]interface{}
		if err := json.Unmarshal(metadataJSON, &metadata); err != nil {
			log.Printf("Error unmarshaling metadata: %v", err)
			continue
		}

		// Проверяем и обновляем информацию о скидке
		if discount, ok := metadata["discount"].(map[string]interface{}); ok {
			if prevPrice, ok := discount["previous_price"].(float64); ok && prevPrice > 0 {
				// Пересчитываем актуальный процент скидки
				if prevPrice > price {
					discountPercent := int((prevPrice - price) / prevPrice * 100)
					discount["discount_percent"] = float64(discountPercent)

					// Обновляем метаданные в БД
					metadata["discount"] = discount
					updatedMetadataJSON, err := json.Marshal(metadata)
					if err != nil {
						log.Printf("Error marshaling updated metadata: %v", err)
						continue
					}

					_, err = db.pool.Exec(ctx, `
                        UPDATE c2c_listings
                        SET metadata = $1
                        WHERE id = $2
                    `, updatedMetadataJSON, id)
					if err != nil {
						log.Printf("Error updating metadata for listing %d: %v", id, err)
						continue
					}

					count++
					log.Printf("Updated discount percentage for listing %d: %d%%", id, discountPercent)
				}
			}
		}
	}

	log.Printf("Synchronized discount metadata for %d listings", count)
	return nil
}
