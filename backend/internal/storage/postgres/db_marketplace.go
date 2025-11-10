// backend/internal/storage/postgres/db_marketplace.go
package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"backend/internal/domain/models"
	"backend/internal/domain/search"

	"github.com/jackc/pgx/v5"
	authservice "github.com/sveturs/auth/pkg/service"
)

func (db *Database) GetListingImageByID(ctx context.Context, imageID int) (*models.MarketplaceImage, error) {
	return db.grpcClient.GetListingImage(ctx, int64(imageID))
}

func (db *Database) DeleteListingImage(ctx context.Context, imageID int) error {
	return db.grpcClient.DeleteListingImage(ctx, int64(imageID))
}

// IndexListing indexes listing in OpenSearch via gRPC microservice
// DEV MODE: grpcClient must be initialized (fail-fast if nil)
func (db *Database) IndexListing(ctx context.Context, listing *models.MarketplaceListing) error {
	return db.grpcClient.IndexListing(ctx, listing)
}

// DeleteListingIndex removes listing from OpenSearch via gRPC microservice
// DEV MODE: grpcClient must be initialized (fail-fast if nil)
func (db *Database) DeleteListingIndex(ctx context.Context, id string) error {
	return db.grpcClient.DeleteListingIndex(ctx, id)
}

// SuggestListings provides autocomplete via gRPC microservice
// DEV MODE: grpcClient must be initialized (fail-fast if nil)
func (db *Database) SuggestListings(ctx context.Context, prefix string, size int) ([]string, error) {
	return db.grpcClient.SuggestListings(ctx, prefix, size)
}

// ReindexAllListings triggers full reindexing via gRPC microservice
// DEV MODE: grpcClient must be initialized (fail-fast if nil)
func (db *Database) ReindexAllListings(ctx context.Context) error {
	// Get all active listings that need reindexing
	listings, err := db.GetMarketplaceListingsForReindex(ctx, 1000)
	if err != nil {
		return fmt.Errorf("failed to get listings for reindex: %w", err)
	}

	log.Printf("Reindexing %d listings via gRPC...", len(listings))

	// Index each listing
	successCount := 0
	errorCount := 0
	for _, listing := range listings {
		if err := db.grpcClient.IndexListing(ctx, listing); err != nil {
			log.Printf("Failed to index listing %d: %v", listing.ID, err)
			errorCount++
		} else {
			successCount++
		}
	}

	log.Printf("Reindexing complete: %d successful, %d errors", successCount, errorCount)

	// Reset reindex flag for successfully indexed listings
	if successCount > 0 {
		listingIDs := make([]int, 0, successCount)
		for _, listing := range listings {
			listingIDs = append(listingIDs, listing.ID)
		}
		if err := db.ResetMarketplaceListingsReindexFlag(ctx, listingIDs); err != nil {
			log.Printf("Warning: failed to reset reindex flags: %v", err)
		}
	}

	if errorCount > 0 {
		return fmt.Errorf("reindexing completed with %d errors", errorCount)
	}

	return nil
}

// GetCategoryAttributes получает атрибуты для указанной категории
func (db *Database) GetCategoryAttributes(ctx context.Context, categoryID int) ([]models.CategoryAttribute, error) {
	return db.GetCategoryAttributesImpl(ctx, categoryID)
}

// SaveListingAttributes сохраняет значения атрибутов для объявления
func (db *Database) SaveListingAttributes(ctx context.Context, listingID int, attributes []models.ListingAttributeValue) error {
	return db.SaveListingAttributesImpl(ctx, listingID, attributes)
}

// GetAttributeRanges возвращает диапазоны значений для фильтруемых атрибутов
func (db *Database) GetAttributeRanges(ctx context.Context, categoryID int) (map[string]map[string]interface{}, error) {
	return db.GetAttributeRangesImpl(ctx, categoryID)
}

// GetListingAttributes получает значения атрибутов для объявления
func (db *Database) GetListingAttributes(ctx context.Context, listingID int) ([]models.ListingAttributeValue, error) {
	return db.GetListingAttributesImpl(ctx, listingID)
}

// GetSession - DEPRECATED: Sessions are now managed via JWT tokens in auth-service
// This method is kept for backward compatibility but should not be used in new code

func (db *Database) GetFavoritedUsers(ctx context.Context, listingID int) ([]int, error) {
	userIDStrs, err := db.grpcClient.GetFavoritedUsers(ctx, int64(listingID))
	if err != nil {
		return nil, err
	}

	// Convert string user IDs to int
	userIDs := make([]int, len(userIDStrs))
	for i, idStr := range userIDStrs {
		var id int
		_, _ = fmt.Sscanf(idStr, "%d", &id)
		userIDs[i] = id
	}

	return userIDs, nil
}

// CreateListing creates a new listing via gRPC microservice
// DEV MODE: grpcClient must be initialized (fail-fast if nil)
func (db *Database) CreateListing(ctx context.Context, listing *models.MarketplaceListing) (int, error) {
	return db.grpcClient.CreateListing(ctx, listing)
}

// GetListings retrieves listings with filters via gRPC microservice
// DEV MODE: grpcClient must be initialized (fail-fast if nil)
func (db *Database) GetListings(ctx context.Context, filters map[string]string, limit int, offset int) ([]models.MarketplaceListing, int64, error) {
	return db.grpcClient.GetListings(ctx, filters, limit, offset)
}

// GetListingByID retrieves a single listing by ID via gRPC microservice
// DEV MODE: grpcClient must be initialized (fail-fast if nil)
func (db *Database) GetListingByID(ctx context.Context, id int) (*models.MarketplaceListing, error) {
	return db.grpcClient.GetListingByID(ctx, id)
}

// GetMarketplaceListingsForReindex возвращает listings с needs_reindex=true для индексации в OpenSearch
func (db *Database) GetMarketplaceListingsForReindex(ctx context.Context, limit int) ([]*models.MarketplaceListing, error) {
	if limit <= 0 {
		limit = 1000 // default limit
	}

	return db.grpcClient.GetMarketplaceListingsForReindex(ctx, limit)
}

// ResetMarketplaceListingsReindexFlag сбрасывает флаг needs_reindex для listings после успешной индексации
func (db *Database) ResetMarketplaceListingsReindexFlag(ctx context.Context, listingIDs []int) error {
	if len(listingIDs) == 0 {
		return nil
	}

	// Convert []int to []int64
	listingIDs64 := make([]int64, len(listingIDs))
	for i, id := range listingIDs {
		listingIDs64[i] = int64(id)
	}

	return db.grpcClient.ResetMarketplaceListingsReindexFlag(ctx, listingIDs64)
}

// GetListingBySlug retrieves a single listing by slug via gRPC microservice
// DEV MODE: grpcClient must be initialized (fail-fast if nil)
func (db *Database) GetListingBySlug(ctx context.Context, slug string) (*models.MarketplaceListing, error) {
	return db.grpcClient.GetListingBySlug(ctx, slug)
}

// IsSlugUnique checks if slug is unique in c2c_listings table
// Implementation in marketplace_slugs.go
// Note: Direct database query, not using gRPC microservice
// Requires 'slug' column to exist in c2c_listings table

// GenerateUniqueSlug generates unique slug by appending suffix if needed
// Implementation in marketplace_slugs.go
// Note: Direct database query, not using gRPC microservice
// Requires 'slug' column to exist in c2c_listings table

// UpdateListing updates an existing listing via gRPC microservice
// DEV MODE: grpcClient must be initialized (fail-fast if nil)
func (db *Database) UpdateListing(ctx context.Context, listing *models.MarketplaceListing) error {
	return db.grpcClient.UpdateListing(ctx, listing)
}

// DeleteListing soft-deletes a listing (user ownership check) via gRPC microservice
// DEV MODE: grpcClient must be initialized (fail-fast if nil)
func (db *Database) DeleteListing(ctx context.Context, id int, userID int) error {
	return db.grpcClient.DeleteListing(ctx, id, userID)
}

// DeleteListingAdmin hard-deletes a listing (admin, no ownership check) via gRPC microservice
// DEV MODE: grpcClient must be initialized (fail-fast if nil)
func (db *Database) DeleteListingAdmin(ctx context.Context, id int) error {
	return db.grpcClient.DeleteListingAdmin(ctx, id)
}

// GetCategories - implementation moved to marketplace_categories.go
// GetAllCategories - implementation moved to marketplace_categories.go
// GetPopularCategories - implementation moved to marketplace_categories.go
// GetCategoryByID - implementation moved to marketplace_categories.go
// GetCategoryTree - implementation moved to marketplace_categories.go

// AddListingImage - implementation moved to marketplace_images.go
// GetListingImages - implementation moved to marketplace_images.go

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

// SearchListings performs OpenSearch query via gRPC microservice
// DEV MODE: grpcClient must be initialized (fail-fast if nil)
func (db *Database) SearchListings(ctx context.Context, params *search.SearchParams) (*search.SearchResult, error) {
	return db.grpcClient.SearchListings(ctx, params)
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
	var settings models.UserPrivacySettings

	query := `
		SELECT
			user_id,
			allow_contact_requests,
			allow_messages_from_contacts_only,
			settings,
			created_at,
			updated_at
		FROM user_privacy_settings
		WHERE user_id = $1
	`

	err := db.pool.QueryRow(ctx, query, userID).Scan(
		&settings.UserID,
		&settings.AllowContactRequests,
		&settings.AllowMessagesFromContactsOnly,
		&settings.Settings,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {
			// Создаем дефолтные настройки если их нет
			insertQuery := `
				INSERT INTO user_privacy_settings (user_id, allow_contact_requests, allow_messages_from_contacts_only, settings)
				VALUES ($1, true, false, '{}'::jsonb)
				RETURNING user_id, allow_contact_requests, allow_messages_from_contacts_only, settings, created_at, updated_at
			`
			err = db.pool.QueryRow(ctx, insertQuery, userID).Scan(
				&settings.UserID,
				&settings.AllowContactRequests,
				&settings.AllowMessagesFromContactsOnly,
				&settings.Settings,
				&settings.CreatedAt,
				&settings.UpdatedAt,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to create default privacy settings: %w", err)
			}
			return &settings, nil
		}
		return nil, fmt.Errorf("failed to get privacy settings: %w", err)
	}

	return &settings, nil
}

func (db *Database) UpdateUserPrivacySettings(ctx context.Context, userID int, settings *models.UpdatePrivacySettingsRequest) error {
	// Сначала убеждаемся что запись существует
	_, err := db.GetUserPrivacySettings(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get existing settings: %w", err)
	}

	// Строим динамический UPDATE query
	updates := []string{}
	args := []interface{}{}
	argPos := 1

	if settings.AllowContactRequests != nil {
		updates = append(updates, fmt.Sprintf("allow_contact_requests = $%d", argPos))
		args = append(args, *settings.AllowContactRequests)
		argPos++
	}

	if settings.AllowMessagesFromContactsOnly != nil {
		updates = append(updates, fmt.Sprintf("allow_messages_from_contacts_only = $%d", argPos))
		args = append(args, *settings.AllowMessagesFromContactsOnly)
		argPos++
	}

	if len(updates) == 0 {
		return nil // Nothing to update
	}

	args = append(args, userID)
	query := fmt.Sprintf(`
		UPDATE user_privacy_settings
		SET %s
		WHERE user_id = $%d
	`, strings.Join(updates, ", "), argPos)

	_, err = db.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update privacy settings: %w", err)
	}

	return nil
}

func (db *Database) GetPrivacySettings(ctx context.Context, userID int) (*models.UserPrivacySettings, error) {
	// Используем реализованный метод GetUserPrivacySettings
	return db.GetUserPrivacySettings(ctx, userID)
}

func (db *Database) UpdatePrivacySettings(ctx context.Context, userID int, settings *models.UpdatePrivacySettingsRequest) error {
	// Используем реализованный метод UpdateUserPrivacySettings
	return db.UpdateUserPrivacySettings(ctx, userID, settings)
}

func (db *Database) UpdateChatSettings(ctx context.Context, userID int, settings *models.ChatUserSettings) error {
	// Сначала убеждаемся что запись существует
	_, err := db.GetUserPrivacySettings(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get existing settings: %w", err)
	}

	// Обновляем JSONB поле settings с chat настройками
	query := `
		UPDATE user_privacy_settings
		SET settings = jsonb_set(
			jsonb_set(
				jsonb_set(
					jsonb_set(
						settings,
						'{auto_translate_chat}',
						to_jsonb($2::boolean)
					),
					'{preferred_language}',
					to_jsonb($3::text)
				),
				'{show_original_language_badge}',
				to_jsonb($4::boolean)
			),
			'{chat_tone_moderation}',
			to_jsonb($5::boolean)
		)
		WHERE user_id = $1
	`

	_, err = db.pool.Exec(ctx, query,
		userID,
		settings.AutoTranslate,
		settings.PreferredLanguage,
		settings.ShowLanguageBadge,
		settings.ModerateTone,
	)
	if err != nil {
		return fmt.Errorf("failed to update chat settings: %w", err)
	}

	return nil
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

	// Convert []models.MarketplaceListingVariant to []*models.MarketplaceListingVariant
	variantsPtrs := make([]*models.MarketplaceListingVariant, len(variants))
	for i := range variants {
		variantsPtrs[i] = &variants[i]
	}

	return db.grpcClient.CreateListingVariants(ctx, int64(listingID), variantsPtrs)
}

func (db *Database) GetListingVariants(ctx context.Context, listingID int) ([]models.MarketplaceListingVariant, error) {
	variantsPtrs, err := db.grpcClient.GetListingVariants(ctx, int64(listingID))
	if err != nil {
		return nil, err
	}

	// Convert []*models.MarketplaceListingVariant to []models.MarketplaceListingVariant
	variants := make([]models.MarketplaceListingVariant, len(variantsPtrs))
	for i, v := range variantsPtrs {
		variants[i] = *v
	}

	return variants, nil
}

func (db *Database) UpdateListingVariant(ctx context.Context, variant *models.MarketplaceListingVariant) error {
	return db.grpcClient.UpdateListingVariant(ctx, variant)
}

func (db *Database) DeleteListingVariant(ctx context.Context, variantID int) error {
	return db.grpcClient.DeleteListingVariant(ctx, int64(variantID))
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
	return db.grpcClient.SynchronizeDiscountMetadata(ctx)
}
