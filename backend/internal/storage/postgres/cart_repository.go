package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

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
	UpdateItemQuantity(ctx context.Context, cartID int64, productID int64, variantID *int64, quantity int) error
	RemoveItem(ctx context.Context, cartID int64, productID int64, variantID *int64) error
	GetItems(ctx context.Context, cartID int64) ([]models.ShoppingCartItem, error)

	CleanupExpiredCarts(ctx context.Context, olderThanDays int) error
	GetAllUserCarts(ctx context.Context, userID int) ([]*models.ShoppingCart, error)
}

// cartRepository реализует интерфейс для работы с корзинами
type cartRepository struct {
	pool *pgxpool.Pool
}

// NewCartRepository создает новый репозиторий корзин
func NewCartRepository(pool *pgxpool.Pool) CartRepositoryInterface {
	return &cartRepository{pool: pool}
}

// GetByID получает корзину по ID
func (r *cartRepository) GetByID(ctx context.Context, cartID int64) (*models.ShoppingCart, error) {
	query := `
		SELECT id, user_id, storefront_id, session_id, created_at, updated_at
		FROM shopping_carts 
		WHERE id = $1`

	var cart models.ShoppingCart
	var userID sql.NullInt32
	var sessionID sql.NullString

	err := r.pool.QueryRow(ctx, query, cartID).Scan(
		&cart.ID,
		&userID,
		&cart.StorefrontID,
		&sessionID,
		&cart.CreatedAt,
		&cart.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("cart not found")
		}
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Обработка NULL значений
	if userID.Valid {
		userIDInt := int(userID.Int32)
		cart.UserID = &userIDInt
	}
	if sessionID.Valid {
		cart.SessionID = &sessionID.String
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
func (r *cartRepository) GetByUser(ctx context.Context, userID int, storefrontID int) (*models.ShoppingCart, error) {
	query := `
		SELECT c.id, c.user_id, c.storefront_id, c.session_id, c.created_at, c.updated_at,
			   s.name, s.slug
		FROM shopping_carts c
		LEFT JOIN storefronts s ON s.id = c.storefront_id
		WHERE c.user_id = $1 AND c.storefront_id = $2`

	var cart models.ShoppingCart
	var dbUserID sql.NullInt32
	var sessionID sql.NullString
	var storefrontName, storefrontSlug sql.NullString

	err := r.pool.QueryRow(ctx, query, userID, storefrontID).Scan(
		&cart.ID,
		&dbUserID,
		&cart.StorefrontID,
		&sessionID,
		&cart.CreatedAt,
		&cart.UpdatedAt,
		&storefrontName,
		&storefrontSlug,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("cart not found")
		}
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Обработка NULL значений
	if dbUserID.Valid {
		userIDInt := int(dbUserID.Int32)
		cart.UserID = &userIDInt
	}
	if sessionID.Valid {
		cart.SessionID = &sessionID.String
	}

	// Заполняем информацию о витрине
	if storefrontName.Valid && storefrontSlug.Valid {
		cart.Storefront = &models.Storefront{
			ID:   cart.StorefrontID,
			Name: storefrontName.String,
			Slug: storefrontSlug.String,
		}
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
func (r *cartRepository) GetBySession(ctx context.Context, sessionID string, storefrontID int) (*models.ShoppingCart, error) {
	query := `
		SELECT c.id, c.user_id, c.storefront_id, c.session_id, c.created_at, c.updated_at,
			   s.name, s.slug
		FROM shopping_carts c
		LEFT JOIN storefronts s ON s.id = c.storefront_id
		WHERE c.session_id = $1 AND c.storefront_id = $2`

	var cart models.ShoppingCart
	var userID sql.NullInt32
	var dbSessionID sql.NullString
	var storefrontName, storefrontSlug sql.NullString

	err := r.pool.QueryRow(ctx, query, sessionID, storefrontID).Scan(
		&cart.ID,
		&userID,
		&cart.StorefrontID,
		&dbSessionID,
		&cart.CreatedAt,
		&cart.UpdatedAt,
		&storefrontName,
		&storefrontSlug,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("cart not found")
		}
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Обработка NULL значений
	if userID.Valid {
		userIDInt := int(userID.Int32)
		cart.UserID = &userIDInt
	}
	if dbSessionID.Valid {
		cart.SessionID = &dbSessionID.String
	}

	// Заполняем информацию о витрине
	if storefrontName.Valid && storefrontSlug.Valid {
		cart.Storefront = &models.Storefront{
			ID:   cart.StorefrontID,
			Name: storefrontName.String,
			Slug: storefrontSlug.String,
		}
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
func (r *cartRepository) Create(ctx context.Context, cart *models.ShoppingCart) (*models.ShoppingCart, error) {
	query := `
		INSERT INTO shopping_carts (user_id, storefront_id, session_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`

	err := r.pool.QueryRow(ctx, query, cart.UserID, cart.StorefrontID, cart.SessionID).Scan(
		&cart.ID,
		&cart.CreatedAt,
		&cart.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create cart: %w", err)
	}

	return cart, nil
}

// Update обновляет корзину
func (r *cartRepository) Update(ctx context.Context, cart *models.ShoppingCart) error {
	query := `
		UPDATE shopping_carts SET
			user_id = $1,
			session_id = $2,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $3`

	_, err := r.pool.Exec(ctx, query, cart.UserID, cart.SessionID, cart.ID)
	if err != nil {
		return fmt.Errorf("failed to update cart: %w", err)
	}

	return nil
}

// Delete удаляет корзину
func (r *cartRepository) Delete(ctx context.Context, cartID int64) error {
	query := `DELETE FROM shopping_carts WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, cartID)
	if err != nil {
		return fmt.Errorf("failed to delete cart: %w", err)
	}
	return nil
}

// Clear очищает корзину (удаляет все позиции)
func (r *cartRepository) Clear(ctx context.Context, cartID int64) error {
	query := `DELETE FROM shopping_cart_items WHERE cart_id = $1`
	_, err := r.pool.Exec(ctx, query, cartID)
	if err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}
	return nil
}

// AddItem добавляет позицию в корзину
func (r *cartRepository) AddItem(ctx context.Context, item *models.ShoppingCartItem) (*models.ShoppingCartItem, error) {
	// Сначала проверяем, есть ли уже такая позиция
	existingQuery := `
		SELECT id, quantity, price_per_unit 
		FROM shopping_cart_items 
		WHERE cart_id = $1 AND product_id = $2 AND COALESCE(variant_id, 0) = COALESCE($3, 0)`

	var existingID int64
	var existingQuantity int
	var existingPrice float64

	err := r.pool.QueryRow(ctx, existingQuery, item.CartID, item.ProductID, item.VariantID).Scan(
		&existingID,
		&existingQuantity,
		&existingPrice,
	)
	if err == nil {
		// Позиция уже существует, обновляем количество
		item.ID = existingID
		item.Quantity += existingQuantity
		item.UpdateTotalPrice()
		err := r.UpdateItem(ctx, item)
		if err != nil {
			return nil, err
		}
		return item, nil
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("failed to check existing item: %w", err)
	}

	// Позиции нет, создаем новую
	item.UpdateTotalPrice()

	query := `
		INSERT INTO shopping_cart_items (
			cart_id, product_id, variant_id, quantity, 
			price_per_unit, total_price
		) VALUES (
			$1, $2, $3, $4, $5, $6
		) RETURNING id, created_at, updated_at`

	priceFloat, _ := item.PricePerUnit.Float64()
	totalFloat, _ := item.TotalPrice.Float64()

	err = r.pool.QueryRow(ctx, query,
		item.CartID,
		item.ProductID,
		item.VariantID,
		item.Quantity,
		priceFloat,
		totalFloat,
	).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to add cart item: %w", err)
	}

	return item, nil
}

// UpdateItem обновляет позицию в корзине
func (r *cartRepository) UpdateItem(ctx context.Context, item *models.ShoppingCartItem) error {
	item.UpdateTotalPrice()

	query := `
		UPDATE shopping_cart_items SET
			quantity = $1,
			price_per_unit = $2,
			total_price = $3,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $4`

	priceFloat, _ := item.PricePerUnit.Float64()
	totalFloat, _ := item.TotalPrice.Float64()

	result, err := r.pool.Exec(ctx, query,
		item.Quantity,
		priceFloat,
		totalFloat,
		item.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update cart item: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("cart item not found")
	}

	return nil
}

// UpdateItemQuantity обновляет количество товара в корзине
func (r *cartRepository) UpdateItemQuantity(ctx context.Context, cartID int64, productID int64, variantID *int64, quantity int) error {
	// Получаем текущую позицию
	query := `
		SELECT id, price_per_unit 
		FROM shopping_cart_items 
		WHERE cart_id = $1 AND product_id = $2 AND COALESCE(variant_id, 0) = COALESCE($3, 0)`

	var itemID int64
	var pricePerUnit float64

	err := r.pool.QueryRow(ctx, query, cartID, productID, variantID).Scan(&itemID, &pricePerUnit)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("cart item not found")
		}
		return fmt.Errorf("failed to get cart item: %w", err)
	}

	// Обновляем количество и пересчитываем общую стоимость
	totalPrice := pricePerUnit * float64(quantity)

	updateQuery := `
		UPDATE shopping_cart_items SET
			quantity = $1,
			total_price = $2,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $3`

	result, err := r.pool.Exec(ctx, updateQuery, quantity, totalPrice, itemID)
	if err != nil {
		return fmt.Errorf("failed to update cart item quantity: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("cart item not found")
	}

	return nil
}

// RemoveItem удаляет позицию из корзины
func (r *cartRepository) RemoveItem(ctx context.Context, cartID int64, productID int64, variantID *int64) error {
	query := `
		DELETE FROM shopping_cart_items 
		WHERE cart_id = $1 AND product_id = $2 AND COALESCE(variant_id, 0) = COALESCE($3, 0)`

	result, err := r.pool.Exec(ctx, query, cartID, productID, variantID)
	if err != nil {
		return fmt.Errorf("failed to remove cart item: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("cart item not found")
	}

	return nil
}

// GetItems получает все позиции корзины с полной информацией о товарах
func (r *cartRepository) GetItems(ctx context.Context, cartID int64) ([]models.ShoppingCartItem, error) {
	query := `
		SELECT 
			ci.id, ci.cart_id, ci.product_id, ci.variant_id, ci.quantity,
			ci.price_per_unit, ci.total_price, ci.created_at, ci.updated_at,
			p.name, p.description, p.price, p.sku, p.stock_quantity,
			p.category_id, p.is_active,
			v.variant_attributes, v.sku, v.price, v.stock_quantity
		FROM shopping_cart_items ci
		LEFT JOIN storefront_products p ON ci.product_id = p.id
		LEFT JOIN storefront_product_variants v ON ci.variant_id = v.id
		WHERE ci.cart_id = $1
		ORDER BY ci.created_at ASC`

	rows, err := r.pool.Query(ctx, query, cartID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart items: %w", err)
	}
	defer rows.Close()

	var items []models.ShoppingCartItem
	for rows.Next() {
		var item models.ShoppingCartItem
		var variantID sql.NullInt64
		var pricePerUnit, totalPrice float64

		// Product fields
		var productName, productDesc, productSKU sql.NullString
		var productPrice sql.NullFloat64
		var productStock, productCategoryID sql.NullInt64
		var productActive sql.NullBool

		// Variant fields
		var variantAttributes, variantSKU sql.NullString
		var variantPrice sql.NullFloat64
		var variantStock sql.NullInt64

		err := rows.Scan(
			&item.ID,
			&item.CartID,
			&item.ProductID,
			&variantID,
			&item.Quantity,
			&pricePerUnit,
			&totalPrice,
			&item.CreatedAt,
			&item.UpdatedAt,
			// Product fields
			&productName,
			&productDesc,
			&productPrice,
			&productSKU,
			&productStock,
			&productCategoryID,
			&productActive,
			// Variant fields
			&variantAttributes,
			&variantSKU,
			&variantPrice,
			&variantStock,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan cart item: %w", err)
		}

		// Обработка NULL значений
		if variantID.Valid {
			item.VariantID = &variantID.Int64
		}

		// Конвертируем цены из float64 в decimal
		item.PricePerUnit = decimal.NewFromFloat(pricePerUnit)
		item.TotalPrice = decimal.NewFromFloat(totalPrice)

		// Заполняем информацию о продукте
		if productName.Valid {
			var productSKUPtr *string
			if productSKU.Valid {
				productSKUPtr = &productSKU.String
			}

			item.Product = &models.StorefrontProduct{
				ID:          int(item.ProductID),
				Name:        productName.String,
				Description: productDesc.String,
				SKU:         productSKUPtr,
				IsActive:    productActive.Bool,
			}

			if productPrice.Valid {
				item.Product.Price = productPrice.Float64
			}
			if productStock.Valid {
				item.Product.StockQuantity = int(productStock.Int64)
			}
			if productCategoryID.Valid {
				item.Product.CategoryID = int(productCategoryID.Int64)
			}

			// Получаем изображения товара
			imagesQuery := `
				SELECT id, image_url, display_order 
				FROM storefront_product_images 
				WHERE storefront_product_id = $1 
				ORDER BY display_order ASC, id ASC
				LIMIT 1`

			var imageID int
			var imageURL string
			var displayOrder int

			err := r.pool.QueryRow(ctx, imagesQuery, item.ProductID).Scan(&imageID, &imageURL, &displayOrder)
			if err == nil {
				item.Product.Images = []models.StorefrontProductImage{
					{
						ID:                  imageID,
						StorefrontProductID: int(item.ProductID),
						ImageURL:            imageURL,
						DisplayOrder:        displayOrder,
					},
				}
			}
		}

		// Заполняем информацию о варианте
		if variantID.Valid {
			var variantSKUPtr *string
			if variantSKU.Valid {
				variantSKUPtr = &variantSKU.String
			}

			// Создаем имя варианта из атрибутов
			variantName := ""
			if variantAttributes.Valid {
				// Пытаемся извлечь атрибуты для создания имени
				// Например: {"color": "red", "size": "L"} -> "Red - L"
				variantName = variantAttributes.String // Пока используем raw JSON как имя
			}

			item.Variant = &models.StorefrontProductVariant{
				ID:                  int(*item.VariantID),
				StorefrontProductID: int(item.ProductID),
				Name:                variantName,
				SKU:                 variantSKUPtr,
			}

			if variantPrice.Valid {
				item.Variant.Price = variantPrice.Float64
			}
			if variantStock.Valid {
				item.Variant.StockQuantity = int(variantStock.Int64)
			}
		}

		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating cart items: %w", err)
	}

	return items, nil
}

// CleanupExpiredCarts удаляет старые корзины
func (r *cartRepository) CleanupExpiredCarts(ctx context.Context, olderThanDays int) error {
	query := fmt.Sprintf(`
		DELETE FROM shopping_carts 
		WHERE created_at < CURRENT_TIMESTAMP - INTERVAL '%d days'
		AND user_id IS NULL`, olderThanDays) // Удаляем только корзины неавторизованных пользователей

	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to cleanup expired carts: %w", err)
	}

	return nil
}

// GetAllUserCarts получает все корзины пользователя
func (r *cartRepository) GetAllUserCarts(ctx context.Context, userID int) ([]*models.ShoppingCart, error) {
	query := `
		SELECT c.id, c.user_id, c.storefront_id, c.session_id, c.created_at, c.updated_at,
			   s.name as storefront_name
		FROM shopping_carts c
		JOIN storefronts s ON s.id = c.storefront_id
		WHERE c.user_id = $1
		ORDER BY c.updated_at DESC`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user carts: %w", err)
	}
	defer rows.Close()

	var carts []*models.ShoppingCart
	for rows.Next() {
		cart := &models.ShoppingCart{}
		var storefrontName string

		err := rows.Scan(
			&cart.ID,
			&cart.UserID,
			&cart.StorefrontID,
			&cart.SessionID,
			&cart.CreatedAt,
			&cart.UpdatedAt,
			&storefrontName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan cart: %w", err)
		}

		// Загружаем товары в корзине
		items, err := r.GetItems(ctx, cart.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get cart items: %w", err)
		}
		cart.Items = items

		carts = append(carts, cart)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return carts, nil
}
