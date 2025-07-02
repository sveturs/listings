package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"backend/internal/domain/models"
)

// CartRepositoryInterface определяет интерфейс для работы с корзинами
type CartRepositoryInterface interface {
	GetByID(ctx context.Context, cartID int64) (*models.ShoppingCart, error)
	GetByUser(ctx context.Context, userID int, storefrontID int) (*models.ShoppingCart, error)
	GetBySession(ctx context.Context, sessionID string, storefrontID int) (*models.ShoppingCart, error)
	Create(ctx context.Context, cart *models.ShoppingCart) (*models.ShoppingCart, error)
	Update(ctx context.Context, cart *models.ShoppingCart) error
	Delete(ctx context.Context, cartID int64) error
	Clear(ctx context.Context, cartID int64) error

	AddItem(ctx context.Context, item *models.ShoppingCartItem) (*models.ShoppingCartItem, error)
	UpdateItem(ctx context.Context, item *models.ShoppingCartItem) error
	RemoveItem(ctx context.Context, cartID int64, productID int64, variantID *int64) error
	GetItems(ctx context.Context, cartID int64) ([]models.ShoppingCartItem, error)

	CleanupExpiredCarts(ctx context.Context, olderThanDays int) error
}

// CartRepository реализует интерфейс для работы с корзинами
type CartRepository struct {
	db *sqlx.DB
}

// NewCartRepository создает новый репозиторий корзин
func NewCartRepository(db *sqlx.DB) *CartRepository {
	return &CartRepository{db: db}
}

// GetByID получает корзину по ID
func (r *CartRepository) GetByID(ctx context.Context, cartID int64) (*models.ShoppingCart, error) {
	query := `
		SELECT id, user_id, storefront_id, session_id, created_at, updated_at
		FROM shopping_carts 
		WHERE id = $1`

	var cart models.ShoppingCart
	err := r.db.GetContext(ctx, &cart, query, cartID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("cart not found")
		}
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Получаем позиции корзины
	items, err := r.GetItems(ctx, cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart items: %w", err)
	}
	cart.Items = items

	return &cart, nil
}

// GetByUser получает корзину пользователя для определенной витрины
func (r *CartRepository) GetByUser(ctx context.Context, userID int, storefrontID int) (*models.ShoppingCart, error) {
	query := `
		SELECT id, user_id, storefront_id, session_id, created_at, updated_at
		FROM shopping_carts 
		WHERE user_id = $1 AND storefront_id = $2`

	var cart models.ShoppingCart
	err := r.db.GetContext(ctx, &cart, query, userID, storefrontID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("cart not found")
		}
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Получаем позиции корзины
	items, err := r.GetItems(ctx, cart.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart items: %w", err)
	}
	cart.Items = items

	return &cart, nil
}

// GetBySession получает корзину по session ID для неавторизованных пользователей
func (r *CartRepository) GetBySession(ctx context.Context, sessionID string, storefrontID int) (*models.ShoppingCart, error) {
	query := `
		SELECT id, user_id, storefront_id, session_id, created_at, updated_at
		FROM shopping_carts 
		WHERE session_id = $1 AND storefront_id = $2`

	var cart models.ShoppingCart
	err := r.db.GetContext(ctx, &cart, query, sessionID, storefrontID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("cart not found")
		}
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Получаем позиции корзины
	items, err := r.GetItems(ctx, cart.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart items: %w", err)
	}
	cart.Items = items

	return &cart, nil
}

// Create создает новую корзину
func (r *CartRepository) Create(ctx context.Context, cart *models.ShoppingCart) (*models.ShoppingCart, error) {
	query := `
		INSERT INTO shopping_carts (user_id, storefront_id, session_id)
		VALUES (:user_id, :storefront_id, :session_id)
		RETURNING id, created_at, updated_at`

	rows, err := r.db.NamedQueryContext(ctx, query, cart)
	if err != nil {
		return nil, fmt.Errorf("failed to create cart: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&cart.ID, &cart.CreatedAt, &cart.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan created cart: %w", err)
		}
	}

	return cart, nil
}

// Update обновляет корзину
func (r *CartRepository) Update(ctx context.Context, cart *models.ShoppingCart) error {
	query := `
		UPDATE shopping_carts SET
			user_id = :user_id,
			session_id = :session_id,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = :id`

	_, err := r.db.NamedExecContext(ctx, query, cart)
	if err != nil {
		return fmt.Errorf("failed to update cart: %w", err)
	}

	return nil
}

// Delete удаляет корзину
func (r *CartRepository) Delete(ctx context.Context, cartID int64) error {
	query := `DELETE FROM shopping_carts WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, cartID)
	if err != nil {
		return fmt.Errorf("failed to delete cart: %w", err)
	}
	return nil
}

// Clear очищает корзину (удаляет все позиции)
func (r *CartRepository) Clear(ctx context.Context, cartID int64) error {
	query := `DELETE FROM shopping_cart_items WHERE cart_id = $1`
	_, err := r.db.ExecContext(ctx, query, cartID)
	if err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}
	return nil
}

// AddItem добавляет позицию в корзину
func (r *CartRepository) AddItem(ctx context.Context, item *models.ShoppingCartItem) (*models.ShoppingCartItem, error) {
	// Сначала проверяем, есть ли уже такая позиция
	existingQuery := `
		SELECT id, quantity, price_per_unit 
		FROM shopping_cart_items 
		WHERE cart_id = $1 AND product_id = $2 AND COALESCE(variant_id, 0) = COALESCE($3, 0)`

	var existingItem struct {
		ID           int64       `db:"id"`
		Quantity     int         `db:"quantity"`
		PricePerUnit interface{} `db:"price_per_unit"`
	}

	err := r.db.GetContext(ctx, &existingItem, existingQuery, item.CartID, item.ProductID, item.VariantID)
	if err == nil {
		// Позиция уже существует, обновляем количество
		item.ID = existingItem.ID
		item.Quantity += existingItem.Quantity
		item.UpdateTotalPrice()
		err := r.UpdateItem(ctx, item)
		if err != nil {
			return nil, err
		}
		return item, nil
	} else if err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to check existing item: %w", err)
	}

	// Позиции нет, создаем новую
	item.UpdateTotalPrice()

	query := `
		INSERT INTO shopping_cart_items (
			cart_id, product_id, variant_id, quantity, 
			price_per_unit, total_price
		) VALUES (
			:cart_id, :product_id, :variant_id, :quantity,
			:price_per_unit, :total_price
		) RETURNING id, created_at, updated_at`

	rows, err := r.db.NamedQueryContext(ctx, query, item)
	if err != nil {
		return nil, fmt.Errorf("failed to add cart item: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan created cart item: %w", err)
		}
	}

	return item, nil
}

// UpdateItem обновляет позицию в корзине
func (r *CartRepository) UpdateItem(ctx context.Context, item *models.ShoppingCartItem) error {
	item.UpdateTotalPrice()

	query := `
		UPDATE shopping_cart_items SET
			quantity = :quantity,
			price_per_unit = :price_per_unit,
			total_price = :total_price,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = :id`

	result, err := r.db.NamedExecContext(ctx, query, item)
	if err != nil {
		return fmt.Errorf("failed to update cart item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("cart item not found")
	}

	return nil
}

// RemoveItem удаляет позицию из корзины
func (r *CartRepository) RemoveItem(ctx context.Context, cartID int64, productID int64, variantID *int64) error {
	query := `
		DELETE FROM shopping_cart_items 
		WHERE cart_id = $1 AND product_id = $2 AND COALESCE(variant_id, 0) = COALESCE($3, 0)`

	result, err := r.db.ExecContext(ctx, query, cartID, productID, variantID)
	if err != nil {
		return fmt.Errorf("failed to remove cart item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("cart item not found")
	}

	return nil
}

// GetItems получает все позиции корзины
func (r *CartRepository) GetItems(ctx context.Context, cartID int64) ([]models.ShoppingCartItem, error) {
	query := `
		SELECT id, cart_id, product_id, variant_id, quantity,
			   price_per_unit, total_price, created_at, updated_at
		FROM shopping_cart_items 
		WHERE cart_id = $1
		ORDER BY created_at ASC`

	var items []models.ShoppingCartItem
	err := r.db.SelectContext(ctx, &items, query, cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart items: %w", err)
	}

	return items, nil
}

// CleanupExpiredCarts удаляет старые корзины
func (r *CartRepository) CleanupExpiredCarts(ctx context.Context, olderThanDays int) error {
	query := `
		DELETE FROM shopping_carts 
		WHERE created_at < CURRENT_TIMESTAMP - INTERVAL '%d days'
		AND user_id IS NULL` // Удаляем только корзины неавторизованных пользователей

	_, err := r.db.ExecContext(ctx, fmt.Sprintf(query, olderThanDays))
	if err != nil {
		return fmt.Errorf("failed to cleanup expired carts: %w", err)
	}

	return nil
}
