// Package postgres implements PostgreSQL repository layer for listings microservice.
package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"github.com/vondi-global/listings/internal/domain"
)

// CartRepository defines operations for cart management
type CartRepository interface {
	// Cart operations
	Create(ctx context.Context, cart *domain.Cart) error
	GetByID(ctx context.Context, cartID int64) (*domain.Cart, error)
	GetByUserAndStorefront(ctx context.Context, userID, storefrontID int64) (*domain.Cart, error)
	GetBySessionAndStorefront(ctx context.Context, sessionID string, storefrontID int64) (*domain.Cart, error)
	GetUserCarts(ctx context.Context, userID int64) ([]*domain.Cart, error)
	Update(ctx context.Context, cart *domain.Cart) error
	Delete(ctx context.Context, cartID int64) error

	// Cart item operations
	AddItem(ctx context.Context, item *domain.CartItem) error
	UpdateItem(ctx context.Context, item *domain.CartItem) error
	RemoveItem(ctx context.Context, cartID, itemID int64) error
	ClearItems(ctx context.Context, cartID int64) error
	GetItemsByCartID(ctx context.Context, cartID int64) ([]*domain.CartItem, error)
	GetCartItemByID(ctx context.Context, cartItemID int64) (*domain.CartItem, error)

	// Transaction support
	WithTx(tx pgx.Tx) CartRepository
}

// cartRepository implements CartRepository using PostgreSQL
type cartRepository struct {
	db     dbOrTx
	logger zerolog.Logger
}

// NewCartRepository creates a new cart repository
func NewCartRepository(pool *pgxpool.Pool, logger zerolog.Logger) CartRepository {
	return &cartRepository{
		db:     pool,
		logger: logger.With().Str("component", "cart_repository").Logger(),
	}
}

// WithTx returns a new repository instance using the provided transaction
func (r *cartRepository) WithTx(tx pgx.Tx) CartRepository {
	return &cartRepository{
		db:     tx,
		logger: r.logger,
	}
}

// Create creates a new cart
func (r *cartRepository) Create(ctx context.Context, cart *domain.Cart) error {
	if err := cart.Validate(); err != nil {
		return fmt.Errorf("invalid cart: %w", err)
	}

	query := `
		INSERT INTO shopping_carts (user_id, session_id, storefront_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		cart.UserID,
		cart.SessionID,
		cart.StorefrontID,
	).Scan(&cart.ID, &cart.CreatedAt, &cart.UpdatedAt)

	if err != nil {
		r.logger.Error().Err(err).Msg("failed to create cart")
		return fmt.Errorf("failed to create cart: %w", err)
	}

	r.logger.Info().Int64("cart_id", cart.ID).Msg("cart created")
	return nil
}

// GetByID retrieves a cart by its ID
func (r *cartRepository) GetByID(ctx context.Context, cartID int64) (*domain.Cart, error) {
	query := `
		SELECT id, user_id, session_id, storefront_id, created_at, updated_at
		FROM shopping_carts
		WHERE id = $1
	`

	var cart domain.Cart
	var userID sql.NullInt64
	var sessionID sql.NullString

	err := r.db.QueryRow(ctx, query, cartID).Scan(
		&cart.ID,
		&userID,
		&sessionID,
		&cart.StorefrontID,
		&cart.CreatedAt,
		&cart.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("cart not found")
		}
		r.logger.Error().Err(err).Int64("cart_id", cartID).Msg("failed to get cart by ID")
		return nil, fmt.Errorf("failed to get cart by ID: %w", err)
	}

	// Handle nullable fields
	if userID.Valid {
		cart.UserID = &userID.Int64
	}
	if sessionID.Valid {
		cart.SessionID = &sessionID.String
	}

	// Load cart items
	items, err := r.GetItemsByCartID(ctx, cart.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load cart items: %w", err)
	}
	cart.Items = items

	return &cart, nil
}

// GetByUserAndStorefront retrieves a cart by user ID and storefront ID
func (r *cartRepository) GetByUserAndStorefront(ctx context.Context, userID, storefrontID int64) (*domain.Cart, error) {
	query := `
		SELECT id, user_id, session_id, storefront_id, created_at, updated_at
		FROM shopping_carts
		WHERE user_id = $1 AND storefront_id = $2
	`

	var cart domain.Cart
	var userIDNullable sql.NullInt64
	var sessionID sql.NullString

	err := r.db.QueryRow(ctx, query, userID, storefrontID).Scan(
		&cart.ID,
		&userIDNullable,
		&sessionID,
		&cart.StorefrontID,
		&cart.CreatedAt,
		&cart.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("cart not found")
		}
		r.logger.Error().Err(err).Int64("user_id", userID).Int64("storefront_id", storefrontID).Msg("failed to get cart")
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Handle nullable fields
	if userIDNullable.Valid {
		cart.UserID = &userIDNullable.Int64
	}
	if sessionID.Valid {
		cart.SessionID = &sessionID.String
	}

	// Load cart items
	items, err := r.GetItemsByCartID(ctx, cart.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load cart items: %w", err)
	}
	cart.Items = items

	return &cart, nil
}

// GetBySessionAndStorefront retrieves a cart by session ID and storefront ID
func (r *cartRepository) GetBySessionAndStorefront(ctx context.Context, sessionID string, storefrontID int64) (*domain.Cart, error) {
	query := `
		SELECT id, user_id, session_id, storefront_id, created_at, updated_at
		FROM shopping_carts
		WHERE session_id = $1 AND storefront_id = $2
	`

	var cart domain.Cart
	var userID sql.NullInt64
	var sessionIDNullable sql.NullString

	err := r.db.QueryRow(ctx, query, sessionID, storefrontID).Scan(
		&cart.ID,
		&userID,
		&sessionIDNullable,
		&cart.StorefrontID,
		&cart.CreatedAt,
		&cart.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("cart not found")
		}
		r.logger.Error().Err(err).Str("session_id", sessionID).Int64("storefront_id", storefrontID).Msg("failed to get cart")
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Handle nullable fields
	if userID.Valid {
		cart.UserID = &userID.Int64
	}
	if sessionIDNullable.Valid {
		cart.SessionID = &sessionIDNullable.String
	}

	// Load cart items
	items, err := r.GetItemsByCartID(ctx, cart.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load cart items: %w", err)
	}
	cart.Items = items

	return &cart, nil
}

// GetUserCarts retrieves all carts for a user (across all storefronts)
func (r *cartRepository) GetUserCarts(ctx context.Context, userID int64) ([]*domain.Cart, error) {
	query := `
		SELECT id, user_id, session_id, storefront_id, created_at, updated_at
		FROM shopping_carts
		WHERE user_id = $1
		ORDER BY updated_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		r.logger.Error().Err(err).Int64("user_id", userID).Msg("failed to get user carts")
		return nil, fmt.Errorf("failed to get user carts: %w", err)
	}
	defer rows.Close()

	var carts []*domain.Cart
	for rows.Next() {
		var cart domain.Cart
		var userIDNullable sql.NullInt64
		var sessionID sql.NullString

		err := rows.Scan(
			&cart.ID,
			&userIDNullable,
			&sessionID,
			&cart.StorefrontID,
			&cart.CreatedAt,
			&cart.UpdatedAt,
		)
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to scan cart")
			return nil, fmt.Errorf("failed to scan cart: %w", err)
		}

		// Handle nullable fields
		if userIDNullable.Valid {
			cart.UserID = &userIDNullable.Int64
		}
		if sessionID.Valid {
			cart.SessionID = &sessionID.String
		}

		// Load cart items
		items, err := r.GetItemsByCartID(ctx, cart.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to load cart items: %w", err)
		}
		cart.Items = items

		carts = append(carts, &cart)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error().Err(err).Msg("error iterating cart rows")
		return nil, fmt.Errorf("error iterating cart rows: %w", err)
	}

	return carts, nil
}

// Update updates an existing cart
func (r *cartRepository) Update(ctx context.Context, cart *domain.Cart) error {
	if err := cart.Validate(); err != nil {
		return fmt.Errorf("invalid cart: %w", err)
	}

	query := `
		UPDATE shopping_carts
		SET user_id = $1, session_id = $2, storefront_id = $3
		WHERE id = $4
		RETURNING updated_at
	`

	err := r.db.QueryRow(ctx, query,
		cart.UserID,
		cart.SessionID,
		cart.StorefrontID,
		cart.ID,
	).Scan(&cart.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("cart not found")
		}
		r.logger.Error().Err(err).Int64("cart_id", cart.ID).Msg("failed to update cart")
		return fmt.Errorf("failed to update cart: %w", err)
	}

	r.logger.Info().Int64("cart_id", cart.ID).Msg("cart updated")
	return nil
}

// Delete deletes a cart by ID
func (r *cartRepository) Delete(ctx context.Context, cartID int64) error {
	query := `DELETE FROM shopping_carts WHERE id = $1`

	result, err := r.db.Exec(ctx, query, cartID)
	if err != nil {
		r.logger.Error().Err(err).Int64("cart_id", cartID).Msg("failed to delete cart")
		return fmt.Errorf("failed to delete cart: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("cart not found")
	}

	r.logger.Info().Int64("cart_id", cartID).Msg("cart deleted")
	return nil
}

// AddItem adds a new item to the cart
func (r *cartRepository) AddItem(ctx context.Context, item *domain.CartItem) error {
	if err := item.Validate(); err != nil {
		return fmt.Errorf("invalid cart item: %w", err)
	}

	query := `
		INSERT INTO cart_items (cart_id, listing_id, variant_id, quantity, price_snapshot)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (cart_id, listing_id, COALESCE(variant_id, 0))
		DO UPDATE SET
			quantity = cart_items.quantity + EXCLUDED.quantity,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		item.CartID,
		item.ListingID,
		item.VariantID,
		item.Quantity,
		item.PriceSnapshot,
	).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)

	if err != nil {
		r.logger.Error().Err(err).Int64("cart_id", item.CartID).Msg("failed to add cart item")
		return fmt.Errorf("failed to add cart item: %w", err)
	}

	r.logger.Info().Int64("cart_item_id", item.ID).Int64("cart_id", item.CartID).Msg("cart item added")
	return nil
}

// UpdateItem updates an existing cart item
func (r *cartRepository) UpdateItem(ctx context.Context, item *domain.CartItem) error {
	if err := item.Validate(); err != nil {
		return fmt.Errorf("invalid cart item: %w", err)
	}

	query := `
		UPDATE cart_items
		SET quantity = $1, price_snapshot = $2
		WHERE id = $3 AND cart_id = $4
		RETURNING updated_at
	`

	err := r.db.QueryRow(ctx, query,
		item.Quantity,
		item.PriceSnapshot,
		item.ID,
		item.CartID,
	).Scan(&item.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("cart item not found")
		}
		r.logger.Error().Err(err).Int64("cart_item_id", item.ID).Msg("failed to update cart item")
		return fmt.Errorf("failed to update cart item: %w", err)
	}

	r.logger.Info().Int64("cart_item_id", item.ID).Msg("cart item updated")
	return nil
}

// RemoveItem removes an item from the cart
func (r *cartRepository) RemoveItem(ctx context.Context, cartID, itemID int64) error {
	query := `DELETE FROM cart_items WHERE id = $1 AND cart_id = $2`

	result, err := r.db.Exec(ctx, query, itemID, cartID)
	if err != nil {
		r.logger.Error().Err(err).Int64("cart_item_id", itemID).Msg("failed to remove cart item")
		return fmt.Errorf("failed to remove cart item: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("cart item not found")
	}

	r.logger.Info().Int64("cart_item_id", itemID).Int64("cart_id", cartID).Msg("cart item removed")
	return nil
}

// ClearItems removes all items from a cart
func (r *cartRepository) ClearItems(ctx context.Context, cartID int64) error {
	query := `DELETE FROM cart_items WHERE cart_id = $1`

	result, err := r.db.Exec(ctx, query, cartID)
	if err != nil {
		r.logger.Error().Err(err).Int64("cart_id", cartID).Msg("failed to clear cart items")
		return fmt.Errorf("failed to clear cart items: %w", err)
	}

	r.logger.Info().Int64("cart_id", cartID).Int64("items_removed", result.RowsAffected()).Msg("cart items cleared")
	return nil
}

// GetCartItemByID retrieves a single cart item by its ID
func (r *cartRepository) GetCartItemByID(ctx context.Context, cartItemID int64) (*domain.CartItem, error) {
	query := `
		SELECT ci.id, ci.cart_id, ci.listing_id, ci.variant_id, ci.quantity, ci.price_snapshot,
		       ci.created_at, ci.updated_at,
		       l.title as listing_name,
		       COALESCE(lv.image_url, (SELECT url FROM listing_images WHERE listing_id = l.id AND is_primary = true LIMIT 1)) as listing_image,
		       lv.attributes as variant_data,
		       CASE
		           WHEN lv.id IS NOT NULL THEN lv.stock
		           ELSE l.quantity
		       END as available_stock,
		       CASE
		           WHEN lv.id IS NOT NULL AND lv.price IS NOT NULL THEN lv.price
		           ELSE l.price
		       END as current_price
		FROM cart_items ci
		JOIN listings l ON ci.listing_id = l.id
		LEFT JOIN listing_variants lv ON ci.variant_id = lv.id
		WHERE ci.id = $1
	`

	var item domain.CartItem
	var variantID sql.NullInt64
	var listingName, listingImage sql.NullString
	var variantDataJSON []byte
	var availableStock sql.NullInt32
	var currentPrice sql.NullFloat64

	err := r.db.QueryRow(ctx, query, cartItemID).Scan(
		&item.ID,
		&item.CartID,
		&item.ListingID,
		&variantID,
		&item.Quantity,
		&item.PriceSnapshot,
		&item.CreatedAt,
		&item.UpdatedAt,
		&listingName,
		&listingImage,
		&variantDataJSON,
		&availableStock,
		&currentPrice,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("cart item not found")
		}
		r.logger.Error().Err(err).Int64("cart_item_id", cartItemID).Msg("failed to get cart item by ID")
		return nil, fmt.Errorf("failed to get cart item by ID: %w", err)
	}

	// Handle nullable fields
	if variantID.Valid {
		item.VariantID = &variantID.Int64
	}
	if listingName.Valid {
		item.ListingName = &listingName.String
	}
	if listingImage.Valid {
		item.ListingImage = &listingImage.String
	}
	if availableStock.Valid {
		item.AvailableStock = &availableStock.Int32
	}
	if currentPrice.Valid {
		item.CurrentPrice = &currentPrice.Float64
	}

	// Parse variant data JSON
	if len(variantDataJSON) > 0 {
		if err := json.Unmarshal(variantDataJSON, &item.VariantData); err != nil {
			r.logger.Error().Err(err).Msg("failed to unmarshal variant data")
		}
	}

	return &item, nil
}

// GetItemsByCartID retrieves all items for a cart
func (r *cartRepository) GetItemsByCartID(ctx context.Context, cartID int64) ([]*domain.CartItem, error) {
	query := `
		SELECT ci.id, ci.cart_id, ci.listing_id, ci.variant_id, ci.quantity, ci.price_snapshot,
		       ci.created_at, ci.updated_at,
		       l.title as listing_name,
		       COALESCE(lv.image_url, (SELECT url FROM listing_images WHERE listing_id = l.id AND is_primary = true LIMIT 1)) as listing_image,
		       lv.attributes as variant_data,
		       CASE
		           WHEN lv.id IS NOT NULL THEN lv.stock
		           ELSE l.quantity
		       END as available_stock,
		       CASE
		           WHEN lv.id IS NOT NULL AND lv.price IS NOT NULL THEN lv.price
		           ELSE l.price
		       END as current_price
		FROM cart_items ci
		JOIN listings l ON ci.listing_id = l.id
		LEFT JOIN listing_variants lv ON ci.variant_id = lv.id
		WHERE ci.cart_id = $1
		ORDER BY ci.created_at ASC
	`

	rows, err := r.db.Query(ctx, query, cartID)
	if err != nil {
		r.logger.Error().Err(err).Int64("cart_id", cartID).Msg("failed to get cart items")
		return nil, fmt.Errorf("failed to get cart items: %w", err)
	}
	defer rows.Close()

	var items []*domain.CartItem
	for rows.Next() {
		var item domain.CartItem
		var variantID sql.NullInt64
		var listingName, listingImage sql.NullString
		var variantDataJSON []byte
		var availableStock sql.NullInt32
		var currentPrice sql.NullFloat64

		err := rows.Scan(
			&item.ID,
			&item.CartID,
			&item.ListingID,
			&variantID,
			&item.Quantity,
			&item.PriceSnapshot,
			&item.CreatedAt,
			&item.UpdatedAt,
			&listingName,
			&listingImage,
			&variantDataJSON,
			&availableStock,
			&currentPrice,
		)
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to scan cart item")
			return nil, fmt.Errorf("failed to scan cart item: %w", err)
		}

		// Handle nullable fields
		if variantID.Valid {
			item.VariantID = &variantID.Int64
		}
		if listingName.Valid {
			item.ListingName = &listingName.String
		}
		if listingImage.Valid {
			item.ListingImage = &listingImage.String
		}
		if availableStock.Valid {
			item.AvailableStock = &availableStock.Int32
		}
		if currentPrice.Valid {
			item.CurrentPrice = &currentPrice.Float64
		}

		// Parse variant data JSON
		if len(variantDataJSON) > 0 {
			if err := json.Unmarshal(variantDataJSON, &item.VariantData); err != nil {
				r.logger.Error().Err(err).Msg("failed to unmarshal variant data")
			}
		}

		items = append(items, &item)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error().Err(err).Msg("error iterating cart item rows")
		return nil, fmt.Errorf("error iterating cart item rows: %w", err)
	}

	return items, nil
}
