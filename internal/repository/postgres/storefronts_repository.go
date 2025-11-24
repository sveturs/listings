package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/sveturs/listings/internal/domain"
)

// CreateStorefront creates a new storefront
func (r *Repository) CreateStorefront(ctx context.Context, storefront *domain.Storefront) error {
	query := `
		INSERT INTO storefronts (
			user_id, slug, name, description, logo_url, banner_url, theme,
			phone, email, website,
			address, city, postal_code, country, latitude, longitude, formatted_address,
			geo_strategy, default_privacy_level, address_verified,
			settings, seo_meta,
			is_active, is_verified, rating, reviews_count, products_count, sales_count, views_count,
			subscription_plan, subscription_expires_at, commission_rate, subscription_id, is_subscription_active,
			ai_agent_enabled, ai_agent_config, live_shopping_enabled, group_buying_enabled,
			followers_count
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
			$21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39
		) RETURNING id, created_at, updated_at`

	return r.db.QueryRowContext(ctx, query,
		storefront.UserID, storefront.Slug, storefront.Name, storefront.Description,
		storefront.LogoURL, storefront.BannerURL, storefront.Theme,
		storefront.Phone, storefront.Email, storefront.Website,
		storefront.Address, storefront.City, storefront.PostalCode, storefront.Country,
		storefront.Latitude, storefront.Longitude, storefront.FormattedAddress,
		storefront.GeoStrategy, storefront.DefaultPrivacyLevel, storefront.AddressVerified,
		storefront.Settings, storefront.SeoMeta,
		storefront.IsActive, storefront.IsVerified, storefront.Rating, storefront.ReviewsCount,
		storefront.ProductsCount, storefront.SalesCount, storefront.ViewsCount,
		storefront.SubscriptionPlan, storefront.SubscriptionExpiresAt, storefront.CommissionRate,
		storefront.SubscriptionID, storefront.IsSubscriptionActive,
		storefront.AIAgentEnabled, storefront.AIAgentConfig, storefront.LiveShoppingEnabled,
		storefront.GroupBuyingEnabled, storefront.FollowersCount,
	).Scan(&storefront.ID, &storefront.CreatedAt, &storefront.UpdatedAt)
}

// GetStorefrontByID retrieves storefront by ID with related entities
func (r *Repository) GetStorefrontByID(ctx context.Context, id int64, includes *domain.Includes) (*domain.Storefront, error) {
	var storefront domain.Storefront
	query := `SELECT * FROM storefronts WHERE id = $1 AND deleted_at IS NULL`

	if err := r.db.GetContext(ctx, &storefront, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("storefront not found")
		}
		return nil, fmt.Errorf("failed to get storefront: %w", err)
	}

	if includes != nil {
		if err := r.loadRelatedEntities(ctx, &storefront, includes); err != nil {
			return nil, err
		}
	}

	return &storefront, nil
}

// GetStorefrontBySlug retrieves storefront by slug with related entities
func (r *Repository) GetStorefrontBySlug(ctx context.Context, slug string, includes *domain.Includes) (*domain.Storefront, error) {
	var storefront domain.Storefront
	query := `SELECT * FROM storefronts WHERE slug = $1 AND deleted_at IS NULL`

	if err := r.db.GetContext(ctx, &storefront, query, slug); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("storefront not found")
		}
		return nil, fmt.Errorf("failed to get storefront: %w", err)
	}

	if includes != nil {
		if err := r.loadRelatedEntities(ctx, &storefront, includes); err != nil {
			return nil, err
		}
	}

	return &storefront, nil
}

// UpdateStorefront updates storefront fields
func (r *Repository) UpdateStorefront(ctx context.Context, id int64, update *domain.StorefrontUpdate) error {
	sets := []string{}
	args := []interface{}{}
	argIdx := 1

	if update.Name != nil {
		sets = append(sets, fmt.Sprintf("name = $%d", argIdx))
		args = append(args, *update.Name)
		argIdx++
	}
	if update.Description != nil {
		sets = append(sets, fmt.Sprintf("description = $%d", argIdx))
		args = append(args, *update.Description)
		argIdx++
	}
	if update.IsActive != nil {
		sets = append(sets, fmt.Sprintf("is_active = $%d", argIdx))
		args = append(args, *update.IsActive)
		argIdx++
	}
	if update.LogoURL != nil {
		sets = append(sets, fmt.Sprintf("logo_url = $%d", argIdx))
		args = append(args, *update.LogoURL)
		argIdx++
	}
	if update.BannerURL != nil {
		sets = append(sets, fmt.Sprintf("banner_url = $%d", argIdx))
		args = append(args, *update.BannerURL)
		argIdx++
	}
	if update.Theme != nil {
		sets = append(sets, fmt.Sprintf("theme = $%d", argIdx))
		args = append(args, update.Theme)
		argIdx++
	}
	if update.Phone != nil {
		sets = append(sets, fmt.Sprintf("phone = $%d", argIdx))
		args = append(args, *update.Phone)
		argIdx++
	}
	if update.Email != nil {
		sets = append(sets, fmt.Sprintf("email = $%d", argIdx))
		args = append(args, *update.Email)
		argIdx++
	}
	if update.Website != nil {
		sets = append(sets, fmt.Sprintf("website = $%d", argIdx))
		args = append(args, *update.Website)
		argIdx++
	}
	if update.Address != nil {
		sets = append(sets, fmt.Sprintf("address = $%d", argIdx))
		args = append(args, *update.Address)
		argIdx++
	}
	if update.City != nil {
		sets = append(sets, fmt.Sprintf("city = $%d", argIdx))
		args = append(args, *update.City)
		argIdx++
	}
	if update.PostalCode != nil {
		sets = append(sets, fmt.Sprintf("postal_code = $%d", argIdx))
		args = append(args, *update.PostalCode)
		argIdx++
	}
	if update.Country != nil {
		sets = append(sets, fmt.Sprintf("country = $%d", argIdx))
		args = append(args, *update.Country)
		argIdx++
	}
	if update.Latitude != nil {
		sets = append(sets, fmt.Sprintf("latitude = $%d", argIdx))
		args = append(args, *update.Latitude)
		argIdx++
	}
	if update.Longitude != nil {
		sets = append(sets, fmt.Sprintf("longitude = $%d", argIdx))
		args = append(args, *update.Longitude)
		argIdx++
	}
	if update.FormattedAddress != nil {
		sets = append(sets, fmt.Sprintf("formatted_address = $%d", argIdx))
		args = append(args, *update.FormattedAddress)
		argIdx++
	}
	if update.Settings != nil {
		sets = append(sets, fmt.Sprintf("settings = $%d", argIdx))
		args = append(args, update.Settings)
		argIdx++
	}
	if update.SeoMeta != nil {
		sets = append(sets, fmt.Sprintf("seo_meta = $%d", argIdx))
		args = append(args, update.SeoMeta)
		argIdx++
	}
	if update.AIAgentEnabled != nil {
		sets = append(sets, fmt.Sprintf("ai_agent_enabled = $%d", argIdx))
		args = append(args, *update.AIAgentEnabled)
		argIdx++
	}
	if update.LiveShoppingEnabled != nil {
		sets = append(sets, fmt.Sprintf("live_shopping_enabled = $%d", argIdx))
		args = append(args, *update.LiveShoppingEnabled)
		argIdx++
	}
	if update.GroupBuyingEnabled != nil {
		sets = append(sets, fmt.Sprintf("group_buying_enabled = $%d", argIdx))
		args = append(args, *update.GroupBuyingEnabled)
		argIdx++
	}

	if len(sets) == 0 {
		return fmt.Errorf("no fields to update")
	}

	sets = append(sets, fmt.Sprintf("updated_at = $%d", argIdx))
	args = append(args, time.Now())
	argIdx++

	args = append(args, id)

	query := fmt.Sprintf("UPDATE storefronts SET %s WHERE id = $%d AND deleted_at IS NULL", strings.Join(sets, ", "), argIdx)

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update storefront: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("storefront not found")
	}

	return nil
}

// DeleteStorefront soft-deletes or hard-deletes a storefront
func (r *Repository) DeleteStorefront(ctx context.Context, id int64, hardDelete bool) error {
	if hardDelete {
		query := `DELETE FROM storefronts WHERE id = $1`
		result, err := r.db.ExecContext(ctx, query, id)
		if err != nil {
			return fmt.Errorf("failed to delete storefront: %w", err)
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get rows affected: %w", err)
		}
		if rows == 0 {
			return fmt.Errorf("storefront not found")
		}
		return nil
	}

	// Soft delete
	query := `UPDATE storefronts SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to soft delete storefront: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("storefront not found")
	}
	return nil
}

// ListStorefronts lists storefronts with filters and pagination
func (r *Repository) ListStorefronts(ctx context.Context, filter *domain.ListStorefrontsFilter) ([]domain.Storefront, int, error) {
	where := []string{"deleted_at IS NULL"}
	args := []interface{}{}
	argIdx := 1

	if filter.UserID != nil {
		where = append(where, fmt.Sprintf("user_id = $%d", argIdx))
		args = append(args, *filter.UserID)
		argIdx++
	}
	if filter.IsActive != nil {
		where = append(where, fmt.Sprintf("is_active = $%d", argIdx))
		args = append(args, *filter.IsActive)
		argIdx++
	}
	if filter.IsVerified != nil {
		where = append(where, fmt.Sprintf("is_verified = $%d", argIdx))
		args = append(args, *filter.IsVerified)
		argIdx++
	}
	if filter.City != nil {
		where = append(where, fmt.Sprintf("city = $%d", argIdx))
		args = append(args, *filter.City)
		argIdx++
	}
	if filter.Country != nil {
		where = append(where, fmt.Sprintf("country = $%d", argIdx))
		args = append(args, *filter.Country)
		argIdx++
	}
	if filter.MinRating != nil {
		where = append(where, fmt.Sprintf("rating >= $%d", argIdx))
		args = append(args, *filter.MinRating)
		argIdx++
	}
	if filter.Latitude != nil && filter.Longitude != nil && filter.RadiusKm != nil {
		where = append(where, fmt.Sprintf(`
			earth_distance(
				ll_to_earth(latitude, longitude),
				ll_to_earth($%d, $%d)
			) <= $%d * 1000`, argIdx, argIdx+1, argIdx+2))
		args = append(args, *filter.Latitude, *filter.Longitude, *filter.RadiusKm)
		argIdx += 3
	}
	if filter.HasAIAgent != nil && *filter.HasAIAgent {
		where = append(where, "ai_agent_enabled = true")
	}
	if filter.HasLiveShopping != nil && *filter.HasLiveShopping {
		where = append(where, "live_shopping_enabled = true")
	}
	if filter.HasGroupBuying != nil && *filter.HasGroupBuying {
		where = append(where, "group_buying_enabled = true")
	}
	if filter.Search != nil && *filter.Search != "" {
		where = append(where, fmt.Sprintf("(name ILIKE $%d OR description ILIKE $%d)", argIdx, argIdx))
		searchPattern := "%" + *filter.Search + "%"
		args = append(args, searchPattern)
		argIdx++
	}

	whereClause := strings.Join(where, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM storefronts WHERE %s", whereClause)
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("failed to count storefronts: %w", err)
	}

	orderBy := "created_at DESC"
	if filter.SortBy != "" {
		direction := "DESC"
		if filter.SortOrder == "asc" {
			direction = "ASC"
		}
		switch filter.SortBy {
		case "rating":
			orderBy = fmt.Sprintf("rating %s", direction)
		case "created_at":
			orderBy = fmt.Sprintf("created_at %s", direction)
		case "products_count":
			orderBy = fmt.Sprintf("products_count %s", direction)
		case "views_count":
			orderBy = fmt.Sprintf("views_count %s", direction)
		}
	}

	limit := int32(10)
	if filter.Limit > 0 {
		limit = filter.Limit
	}
	if limit > 100 {
		limit = 100
	}

	page := int32(1)
	if filter.Page > 0 {
		page = filter.Page
	}
	offset := (page - 1) * limit

	query := fmt.Sprintf(`
		SELECT * FROM storefronts
		WHERE %s
		ORDER BY %s
		LIMIT %d OFFSET %d`,
		whereClause, orderBy, limit, offset)

	var storefronts []domain.Storefront
	if err := r.db.SelectContext(ctx, &storefronts, query, args...); err != nil {
		return nil, 0, fmt.Errorf("failed to list storefronts: %w", err)
	}

	return storefronts, total, nil
}

// AddStaff adds a staff member to storefront
func (r *Repository) AddStaff(ctx context.Context, staff *domain.StorefrontStaff) error {
	query := `
		INSERT INTO storefront_staff (storefront_id, user_id, role, permissions, actions_count)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRowContext(ctx, query,
		staff.StorefrontID, staff.UserID, staff.Role, staff.Permissions, staff.ActionsCount,
	).Scan(&staff.ID, &staff.CreatedAt, &staff.UpdatedAt)
}

// UpdateStaff updates staff member information
func (r *Repository) UpdateStaff(ctx context.Context, id int64, update *domain.StaffUpdate) error {
	sets := []string{}
	args := []interface{}{}
	argIdx := 1

	if update.Role != nil {
		sets = append(sets, fmt.Sprintf("role = $%d", argIdx))
		args = append(args, *update.Role)
		argIdx++
	}
	if update.Permissions != nil {
		sets = append(sets, fmt.Sprintf("permissions = $%d", argIdx))
		args = append(args, update.Permissions)
		argIdx++
	}

	if len(sets) == 0 {
		return fmt.Errorf("no fields to update")
	}

	sets = append(sets, fmt.Sprintf("updated_at = $%d", argIdx))
	args = append(args, time.Now())
	argIdx++

	args = append(args, id)

	query := fmt.Sprintf("UPDATE storefront_staff SET %s WHERE id = $%d", strings.Join(sets, ", "), argIdx)

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update staff: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("staff not found")
	}

	return nil
}

// RemoveStaff removes a staff member from storefront
func (r *Repository) RemoveStaff(ctx context.Context, storefrontID, userID int64) error {
	query := `DELETE FROM storefront_staff WHERE storefront_id = $1 AND user_id = $2`

	result, err := r.db.ExecContext(ctx, query, storefrontID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove staff: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("staff not found")
	}

	return nil
}

// GetStaff retrieves all staff members for a storefront
func (r *Repository) GetStaff(ctx context.Context, storefrontID int64) ([]domain.StorefrontStaff, error) {
	query := `SELECT * FROM storefront_staff WHERE storefront_id = $1 ORDER BY created_at`

	var staff []domain.StorefrontStaff
	if err := r.db.SelectContext(ctx, &staff, query, storefrontID); err != nil {
		return nil, fmt.Errorf("failed to get staff: %w", err)
	}

	return staff, nil
}

// SetWorkingHours sets working hours for a storefront (replaces existing)
func (r *Repository) SetWorkingHours(ctx context.Context, storefrontID int64, hours []domain.StorefrontHours) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	deleteQuery := `DELETE FROM storefront_hours WHERE storefront_id = $1`
	if _, err := tx.ExecContext(ctx, deleteQuery, storefrontID); err != nil {
		return fmt.Errorf("failed to delete existing hours: %w", err)
	}

	insertQuery := `
		INSERT INTO storefront_hours (
			storefront_id, day_of_week, open_time, close_time, is_closed, special_date, special_note
		) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	for _, hour := range hours {
		if _, err := tx.ExecContext(ctx, insertQuery,
			storefrontID, hour.DayOfWeek, hour.OpenTime, hour.CloseTime,
			hour.IsClosed, hour.SpecialDate, hour.SpecialNote,
		); err != nil {
			return fmt.Errorf("failed to insert hours: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetWorkingHours retrieves working hours for a storefront
func (r *Repository) GetWorkingHours(ctx context.Context, storefrontID int64) ([]domain.StorefrontHours, error) {
	query := `SELECT * FROM storefront_hours WHERE storefront_id = $1 ORDER BY day_of_week, special_date`

	var hours []domain.StorefrontHours
	if err := r.db.SelectContext(ctx, &hours, query, storefrontID); err != nil {
		return nil, fmt.Errorf("failed to get working hours: %w", err)
	}

	return hours, nil
}

// IsOpenNow checks if storefront is currently open
func (r *Repository) IsOpenNow(ctx context.Context, storefrontID int64) (bool, *time.Time, *time.Time, error) {
	now := time.Now()
	currentDay := int32(now.Weekday())
	currentTime := now.Format("15:04:05")
	currentDate := now.Format("2006-01-02")

	query := `
		SELECT open_time, close_time, is_closed
		FROM storefront_hours
		WHERE storefront_id = $1
			AND (
				(day_of_week = $2 AND special_date IS NULL)
				OR special_date = $3
			)
		ORDER BY special_date DESC NULLS LAST
		LIMIT 1`

	var openTime, closeTime *string
	var isClosed bool

	err := r.db.QueryRowContext(ctx, query, storefrontID, currentDay, currentDate).Scan(&openTime, &closeTime, &isClosed)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil, nil, nil
		}
		return false, nil, nil, fmt.Errorf("failed to check if open: %w", err)
	}

	if isClosed {
		return false, nil, nil, nil
	}

	if openTime == nil || closeTime == nil {
		return false, nil, nil, nil
	}

	isOpen := currentTime >= *openTime && currentTime <= *closeTime

	var nextOpenTime, nextCloseTime *time.Time
	if !isOpen {
		nextOpen, _ := time.Parse("15:04:05", *openTime)
		nextOpenTime = &nextOpen
	} else {
		nextClose, _ := time.Parse("15:04:05", *closeTime)
		nextCloseTime = &nextClose
	}

	return isOpen, nextOpenTime, nextCloseTime, nil
}

// SetPaymentMethods sets payment methods for a storefront (replaces existing)
func (r *Repository) SetPaymentMethods(ctx context.Context, storefrontID int64, methods []domain.PaymentMethod) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	deleteQuery := `DELETE FROM storefront_payment_methods WHERE storefront_id = $1`
	if _, err := tx.ExecContext(ctx, deleteQuery, storefrontID); err != nil {
		return fmt.Errorf("failed to delete existing payment methods: %w", err)
	}

	insertQuery := `
		INSERT INTO storefront_payment_methods (
			storefront_id, method_type, is_enabled, provider, settings,
			transaction_fee, min_amount, max_amount
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	for _, method := range methods {
		if _, err := tx.ExecContext(ctx, insertQuery,
			storefrontID, method.MethodType, method.IsEnabled, method.Provider,
			method.Settings, method.TransactionFee, method.MinAmount, method.MaxAmount,
		); err != nil {
			return fmt.Errorf("failed to insert payment method: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetPaymentMethods retrieves payment methods for a storefront
func (r *Repository) GetPaymentMethods(ctx context.Context, storefrontID int64) ([]domain.PaymentMethod, error) {
	query := `SELECT * FROM storefront_payment_methods WHERE storefront_id = $1 ORDER BY created_at`

	var methods []domain.PaymentMethod
	if err := r.db.SelectContext(ctx, &methods, query, storefrontID); err != nil {
		return nil, fmt.Errorf("failed to get payment methods: %w", err)
	}

	return methods, nil
}

// SetDeliveryOptions sets delivery options for a storefront (replaces existing)
func (r *Repository) SetDeliveryOptions(ctx context.Context, storefrontID int64, options []domain.StorefrontDeliveryOption) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	deleteQuery := `DELETE FROM storefront_delivery_options WHERE storefront_id = $1`
	if _, err := tx.ExecContext(ctx, deleteQuery, storefrontID); err != nil {
		return fmt.Errorf("failed to delete existing delivery options: %w", err)
	}

	insertQuery := `
		INSERT INTO storefront_delivery_options (
			storefront_id, name, description, base_price, price_per_km, price_per_kg,
			free_above_amount, min_order_amount, max_weight_kg, max_distance_km,
			estimated_days_min, estimated_days_max, zones, available_days, cutoff_time,
			provider, provider_config, is_active, display_order
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)`

	for _, option := range options {
		if _, err := tx.ExecContext(ctx, insertQuery,
			storefrontID, option.Name, option.Description, option.BasePrice,
			option.PricePerKm, option.PricePerKg, option.FreeAboveAmount,
			option.MinOrderAmount, option.MaxWeightKg, option.MaxDistanceKm,
			option.EstimatedDaysMin, option.EstimatedDaysMax, option.Zones,
			option.AvailableDays, option.CutoffTime, option.Provider,
			option.ProviderConfig, option.IsActive, option.DisplayOrder,
		); err != nil {
			return fmt.Errorf("failed to insert delivery option: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetDeliveryOptions retrieves delivery options for a storefront
func (r *Repository) GetDeliveryOptions(ctx context.Context, storefrontID int64) ([]domain.StorefrontDeliveryOption, error) {
	query := `SELECT * FROM storefront_delivery_options WHERE storefront_id = $1 ORDER BY display_order, created_at`

	var options []domain.StorefrontDeliveryOption
	if err := r.db.SelectContext(ctx, &options, query, storefrontID); err != nil {
		return nil, fmt.Errorf("failed to get delivery options: %w", err)
	}

	return options, nil
}

// GetMapData retrieves storefronts for map display within bounds
func (r *Repository) GetMapData(ctx context.Context, bounds *domain.MapBounds, filter *domain.ListStorefrontsFilter) ([]domain.StorefrontMapData, error) {
	where := []string{
		"latitude IS NOT NULL",
		"longitude IS NOT NULL",
		"deleted_at IS NULL",
		"latitude BETWEEN $1 AND $2",
		"longitude BETWEEN $3 AND $4",
	}
	args := []interface{}{bounds.South, bounds.North, bounds.West, bounds.East}
	argIdx := 5

	if filter != nil {
		if filter.IsActive != nil {
			where = append(where, fmt.Sprintf("is_active = $%d", argIdx))
			args = append(args, *filter.IsActive)
			argIdx++
		}
		if filter.IsVerified != nil {
			where = append(where, fmt.Sprintf("is_verified = $%d", argIdx))
			args = append(args, *filter.IsVerified)
			argIdx++
		}
	}

	query := fmt.Sprintf(`
		SELECT
			id, slug, name, latitude, longitude, rating, logo_url, address, phone,
			false as working_now,
			products_count,
			false as supports_cod,
			false as has_delivery,
			false as has_self_pickup,
			false as accepts_cards
		FROM storefronts
		WHERE %s
		LIMIT 1000`,
		strings.Join(where, " AND "))

	var mapData []domain.StorefrontMapData
	if err := r.db.SelectContext(ctx, &mapData, query, args...); err != nil {
		return nil, fmt.Errorf("failed to get map data: %w", err)
	}

	return mapData, nil
}

// GetStorefrontDashboardStats retrieves dashboard statistics for a storefront
func (r *Repository) GetStorefrontDashboardStats(ctx context.Context, storefrontID int64, from, to *time.Time) (*domain.StorefrontDashboardStats, error) {
	stats := &domain.StorefrontDashboardStats{}

	query := `
		SELECT
			COALESCE(SUM(CASE WHEN is_active = true THEN 1 ELSE 0 END), 0) as total_products,
			COALESCE(SUM(CASE WHEN is_active = true THEN 1 ELSE 0 END), 0) as active_products
		FROM listings
		WHERE storefront_id = $1`

	if err := r.db.QueryRowContext(ctx, query, storefrontID).Scan(
		&stats.TotalProducts, &stats.ActiveProducts,
	); err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get product stats: %w", err)
	}

	stats.ViewsCount = 0
	stats.UniqueVisitors = 0
	stats.ConversionRate = 0.0

	return stats, nil
}

// IsSlugTaken checks if slug is already taken by another storefront
func (r *Repository) IsSlugTaken(ctx context.Context, slug string, excludeID *int64) (bool, error) {
	query := `SELECT COUNT(*) FROM storefronts WHERE slug = $1 AND deleted_at IS NULL`
	args := []interface{}{slug}

	if excludeID != nil {
		query += " AND id != $2"
		args = append(args, *excludeID)
	}

	var count int
	if err := r.db.GetContext(ctx, &count, query, args...); err != nil {
		return false, fmt.Errorf("failed to check slug: %w", err)
	}

	return count > 0, nil
}

// IncrementViewsCount increments views count for a storefront
func (r *Repository) IncrementViewsCount(ctx context.Context, storefrontID int64) error {
	query := `UPDATE storefronts SET views_count = views_count + 1, updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	_, err := r.db.ExecContext(ctx, query, storefrontID)
	if err != nil {
		return fmt.Errorf("failed to increment views count: %w", err)
	}

	return nil
}

// loadRelatedEntities loads related entities based on includes
func (r *Repository) loadRelatedEntities(ctx context.Context, storefront *domain.Storefront, includes *domain.Includes) error {
	if includes.Staff {
		staff, err := r.GetStaff(ctx, storefront.ID)
		if err != nil {
			return fmt.Errorf("failed to load staff: %w", err)
		}
		storefront.Staff = staff
	}

	if includes.Hours {
		hours, err := r.GetWorkingHours(ctx, storefront.ID)
		if err != nil {
			return fmt.Errorf("failed to load hours: %w", err)
		}
		storefront.Hours = hours
	}

	if includes.PaymentMethods {
		methods, err := r.GetPaymentMethods(ctx, storefront.ID)
		if err != nil {
			return fmt.Errorf("failed to load payment methods: %w", err)
		}
		storefront.PaymentMethods = methods
	}

	if includes.DeliveryOptions {
		options, err := r.GetDeliveryOptions(ctx, storefront.ID)
		if err != nil {
			return fmt.Errorf("failed to load delivery options: %w", err)
		}
		storefront.DeliveryOptions = options
	}

	return nil
}
