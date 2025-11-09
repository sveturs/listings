package service

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"

	"backend/internal/config"
	"backend/internal/domain/models"
	distributedlock "backend/pkg/distributed_lock"

	pb "github.com/sveturs/listings/api/proto/listings/v1"
)

// CreateOrderWithTx создает новый заказ с использованием gRPC Listings Service для управления стоками
func (s *OrderService) CreateOrderWithTx(ctx context.Context, db *sqlx.DB, req *models.CreateOrderRequest, userID int) (*models.StorefrontOrder, error) {
	s.logger.Info("Creating order with gRPC Listings Service (user_id: %d, storefront_id: %d)", userID, req.StorefrontID)

	// Validate items
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("orders.empty_items")
	}

	// Prepare stock items for gRPC
	productIDs := make([]string, 0, len(req.Items))
	stockItems := make([]*pb.StockItem, 0, len(req.Items))

	for _, item := range req.Items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("orders.invalid_quantity")
		}

		productIDs = append(productIDs, fmt.Sprintf("%d", item.ProductID))

		stockItem := &pb.StockItem{
			ProductId: int64(item.ProductID),
			Quantity:  int32(item.Quantity),
		}

		if item.VariantID != nil {
			variantID := int64(*item.VariantID)
			stockItem.VariantId = &variantID
		}

		stockItems = append(stockItems, stockItem)
	}

	// Sort product IDs to prevent deadlocks
	sort.Strings(productIDs)
	lockKey := fmt.Sprintf("order:lock:%s", strings.Join(productIDs, ","))

	// 1. Acquire distributed lock
	lock := distributedlock.NewRedisLock(s.redisClient, lockKey, 30*time.Second)
	acquired, err := lock.TryLock(ctx)
	if err != nil {
		s.logger.Error("Failed to acquire lock (key: %s): %v", lockKey, err)
		return nil, fmt.Errorf("orders.lock_failed")
	}
	if !acquired {
		s.logger.Warn("Lock already held (key: %s)", lockKey)
		return nil, fmt.Errorf("orders.lock_busy")
	}
	defer func() {
		if unlockErr := lock.Unlock(ctx); unlockErr != nil {
			s.logger.Error("Failed to unlock (key: %s): %v", lockKey, unlockErr)
		}
	}()

	// 2. Check stock availability first
	availResp, err := s.listingsClient.CheckStockAvailability(ctx, stockItems)
	if err != nil {
		s.logger.Error("Failed to check stock availability: %v", err)
		return nil, fmt.Errorf("orders.stock_check_failed")
	}
	if !availResp.AllAvailable {
		s.logger.Warn("Insufficient stock for some items")
		return nil, fmt.Errorf("orders.insufficient_stock")
	}

	// 3. Generate order ID early (needed for DecrementStock)
	orderID := uuid.New().String()

	// 4. Decrement stock via gRPC
	decrementResp, err := s.listingsClient.DecrementStock(ctx, stockItems, orderID)
	if err != nil {
		s.logger.Error("Failed to call DecrementStock: %v", err)
		return nil, fmt.Errorf("orders.stock_service_error")
	}
	if !decrementResp.Success {
		errorMsg := "unknown error"
		if decrementResp.Error != nil {
			errorMsg = *decrementResp.Error
		}
		s.logger.Warn("Stock decrement failed: %s", errorMsg)
		return nil, fmt.Errorf("orders.insufficient_stock")
	}

	// 5. Create order in database (SQL transaction)
	tx, err := db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		s.logger.Error("Failed to begin transaction: %v", err)
		// Rollback stock
		if rollbackErr := s.listingsClient.RollbackStock(ctx, stockItems, orderID); rollbackErr != nil {
			s.logger.Error("Failed to rollback stock after tx error: %v", rollbackErr)
		}
		return nil, fmt.Errorf("orders.transaction_failed")
	}
	defer func() {
		_ = tx.Rollback() // ignore error - transaction might already be committed
	}()

	// Get storefront
	storefront, err := s.getStorefrontTx(ctx, tx, req.StorefrontID)
	if err != nil {
		s.logger.Error("Failed to get storefront: %v", err)
		if rollbackErr := s.listingsClient.RollbackStock(ctx, stockItems, orderID); rollbackErr != nil {
			s.logger.Error("Failed to rollback stock after storefront error: %v", rollbackErr)
		}
		return nil, fmt.Errorf("orders.storefront_not_found")
	}

	if !storefront.IsActive {
		s.logger.Warn("Storefront is not active (storefront_id: %d)", req.StorefrontID)
		if rollbackErr := s.listingsClient.RollbackStock(ctx, stockItems, orderID); rollbackErr != nil {
			s.logger.Error("Failed to rollback stock after inactive storefront: %v", rollbackErr)
		}
		return nil, fmt.Errorf("orders.storefront_inactive")
	}

	// Prepare order struct
	order := s.prepareOrderStruct(req, userID, storefront)

	// Insert order with generated UUID
	createdOrder, err := s.createOrderRecordWithUUID(ctx, tx, order, orderID)
	if err != nil {
		s.logger.Error("Failed to create order record: %v", err)
		// Rollback stock
		if rollbackErr := s.listingsClient.RollbackStock(ctx, stockItems, orderID); rollbackErr != nil {
			s.logger.Error("Failed to rollback stock after order error: %v", rollbackErr)
		}
		return nil, fmt.Errorf("orders.creation_failed")
	}

	// Clear cart if needed
	if req.CartID != nil {
		if err := s.clearCartTx(ctx, tx, *req.CartID); err != nil {
			s.logger.Error("Failed to clear cart %d: %v", *req.CartID, err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		s.logger.Error("Failed to commit transaction: %v", err)
		// Rollback stock
		if rollbackErr := s.listingsClient.RollbackStock(ctx, stockItems, orderID); rollbackErr != nil {
			s.logger.Error("Failed to rollback stock after commit error: %v", rollbackErr)
		}
		return nil, fmt.Errorf("orders.commit_failed")
	}

	s.logger.Info("Order created successfully with stock decremented (order_id: %s, items_count: %d)", orderID, len(stockItems))

	return createdOrder, nil
}

// generateOrderNumberForOrder генерирует уникальный номер заказа
// Дублируется из order_service_tx.go для независимости модулей
func generateOrderNumberForOrder(userID int, storefrontID int) string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("ORD-%d-%d-%d", storefrontID, userID, timestamp)
}

// prepareOrderStruct подготавливает структуру заказа
func (s *OrderService) prepareOrderStruct(req *models.CreateOrderRequest, userID int, storefront *models.Storefront) *models.StorefrontOrder {
	// Формируем адрес забора из данных витрины
	pickupAddress := map[string]interface{}{
		"street":      storefront.Address,
		"city":        storefront.City,
		"postal_code": storefront.PostalCode,
		"country":     storefront.Country,
		"latitude":    storefront.Latitude,
		"longitude":   storefront.Longitude,
		"name":        storefront.Name, // Название витрины для удобства
		"phone":       storefront.Phone,
		"email":       storefront.Email,
	}

	return &models.StorefrontOrder{
		StorefrontID:    req.StorefrontID,
		CustomerID:      userID,
		UserID:          userID, // Для совместимости
		OrderNumber:     generateOrderNumberForOrder(userID, req.StorefrontID),
		Status:          models.OrderStatusPending,
		Currency:        config.GetGlobalDefaultCurrency(),
		ShippingMethod:  &req.ShippingMethod,
		CustomerNotes:   &req.CustomerNotes,
		EscrowDays:      s.calculateEscrowDays(storefront),
		ShippingAddress: convertToJSONB(req.ShippingAddress),
		BillingAddress:  convertToJSONB(req.BillingAddress),
		PickupAddress:   convertToJSONB(pickupAddress),
		PaymentMethod:   req.PaymentMethod,
		PaymentStatus:   "pending",
		Metadata:        make(map[string]interface{}),
		// Инициализируем финансовые поля нулями
		SubtotalAmount:   decimal.Zero,
		TaxAmount:        decimal.Zero,
		ShippingAmount:   decimal.Zero,
		Discount:         decimal.Zero,
		TotalAmount:      decimal.Zero,
		CommissionAmount: decimal.Zero,
		SellerAmount:     decimal.Zero,
		// Алиасы для совместимости
		Subtotal: decimal.Zero,
		Tax:      decimal.Zero,
		Shipping: decimal.Zero,
		Total:    decimal.Zero,
	}
}

// Вспомогательные методы для работы с транзакциями

// createOrderRecordWithUUID создает запись заказа в БД с заданным UUID (для интеграции с Listings Service)
func (s *OrderService) createOrderRecordWithUUID(
	ctx context.Context,
	tx *sqlx.Tx,
	order *models.StorefrontOrder,
	orderUUID string,
) (*models.StorefrontOrder, error) {
	// TODO: Временно используем orderID int64, конвертируя UUID -> hash
	// В будущем нужно мигрировать b2c_orders.id на UUID type
	// Пока используем простой подход: генерируем int64 ID из БД
	query := `
		INSERT INTO b2c_orders (
			storefront_id, customer_id, order_number, subtotal_amount, shipping_amount,
			tax_amount, total_amount, commission_amount, seller_amount,
			currency, status, escrow_days, shipping_address, billing_address,
			shipping_method, customer_notes, payment_method, payment_status, metadata, pickup_address
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20
		) RETURNING id, created_at, updated_at`

	var createdOrder models.StorefrontOrder
	err := tx.QueryRowContext(ctx, query,
		order.StorefrontID,
		order.CustomerID,
		order.OrderNumber,
		order.SubtotalAmount,
		order.ShippingAmount,
		order.TaxAmount,
		order.TotalAmount,
		order.CommissionAmount,
		order.SellerAmount,
		order.Currency,
		order.Status,
		order.EscrowDays,
		order.ShippingAddress,
		order.BillingAddress,
		order.ShippingMethod,
		order.CustomerNotes,
		order.PaymentMethod,
		order.PaymentStatus,
		order.Metadata,
		order.PickupAddress,
	).Scan(&createdOrder.ID, &createdOrder.CreatedAt, &createdOrder.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert order: %w", err)
	}

	// Копируем остальные поля
	createdOrder.StorefrontID = order.StorefrontID
	createdOrder.CustomerID = order.CustomerID
	createdOrder.UserID = order.UserID
	createdOrder.OrderNumber = order.OrderNumber
	createdOrder.Status = order.Status
	createdOrder.Currency = order.Currency
	createdOrder.ShippingMethod = order.ShippingMethod
	createdOrder.CustomerNotes = order.CustomerNotes
	createdOrder.EscrowDays = order.EscrowDays
	createdOrder.ShippingAddress = order.ShippingAddress
	createdOrder.BillingAddress = order.BillingAddress
	createdOrder.PickupAddress = order.PickupAddress
	createdOrder.PaymentMethod = order.PaymentMethod
	createdOrder.PaymentStatus = order.PaymentStatus
	createdOrder.Metadata = order.Metadata
	createdOrder.SubtotalAmount = order.SubtotalAmount
	createdOrder.TaxAmount = order.TaxAmount
	createdOrder.ShippingAmount = order.ShippingAmount
	createdOrder.Discount = order.Discount
	createdOrder.TotalAmount = order.TotalAmount
	createdOrder.CommissionAmount = order.CommissionAmount
	createdOrder.SellerAmount = order.SellerAmount
	createdOrder.Items = order.Items

	// Store UUID in metadata for future reference
	if createdOrder.Metadata == nil {
		createdOrder.Metadata = make(map[string]interface{})
	}
	createdOrder.Metadata["order_uuid"] = orderUUID

	return &createdOrder, nil
}

func (s *OrderService) getStorefrontTx(ctx context.Context, tx *sqlx.Tx, storefrontID int) (*models.Storefront, error) {
	var storefront models.Storefront
	query := `SELECT * FROM b2c_stores WHERE id = $1 FOR SHARE`
	err := tx.GetContext(ctx, &storefront, query, storefrontID)
	if err != nil {
		// TODO: Migrate storefronts to separate microservice
		// Graceful degradation: return minimal storefront info if table doesn't exist
		s.logger.Warn("Failed to get storefront (legacy table may be dropped): storefront_id=%d, error=%v", storefrontID, err)
		return &models.Storefront{
			ID:       storefrontID,
			Name:     "Unknown Store",
			IsActive: true,
		}, nil
	}
	return &storefront, nil
}

func (s *OrderService) getOrderItemsTx(ctx context.Context, tx *sqlx.Tx, req *models.CreateOrderRequest, userID int) ([]models.OrderItemRequest, error) {
	var items []models.OrderItemRequest

	if req.CartID != nil {
		// Блокируем корзину для чтения
		var cart models.ShoppingCart
		query := `SELECT * FROM b2c_carts WHERE id = $1 FOR UPDATE`
		err := tx.GetContext(ctx, &cart, query, *req.CartID)
		if err != nil {
			return nil, fmt.Errorf("failed to get cart: %w", err)
		}

		if cart.UserID == nil || *cart.UserID != userID {
			return nil, fmt.Errorf("cart does not belong to user")
		}

		if cart.StorefrontID != req.StorefrontID {
			return nil, fmt.Errorf("cart belongs to different storefront")
		}

		// Получаем позиции корзины
		var cartItems []models.ShoppingCartItem
		itemsQuery := `SELECT * FROM b2c_cart_items WHERE cart_id = $1`
		err = tx.SelectContext(ctx, &cartItems, itemsQuery, cart.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get cart items: %w", err)
		}

		for _, cartItem := range cartItems {
			items = append(items, models.OrderItemRequest{
				ProductID: cartItem.ProductID,
				VariantID: cartItem.VariantID,
				Quantity:  cartItem.Quantity,
			})
		}
	} else {
		items = req.Items
	}

	return items, nil
}

// NOTE: Removed deprecated methods that are now handled by Listings gRPC Service:
// - lockProductForUpdate (stock locking now via distributed Redis lock + gRPC)
// - createOrderInTransaction (replaced by createOrderRecordWithUUID)
// - createReservationTx (inventory reservations removed - stock managed by Listings Service)
// - updateProductStockTx (stock updates now via gRPC DecrementStock/RollbackStock)

func (s *OrderService) createOrderItemTx(ctx context.Context, tx *sqlx.Tx, item *models.StorefrontOrderItem) error {
	query := `
		INSERT INTO b2c_order_items (
			order_id, product_id, variant_id, product_name, variant_name,
			product_sku, quantity, unit_price, total_price
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)`

	_, err := tx.ExecContext(ctx, query,
		item.OrderID,
		item.ProductID,
		item.VariantID,
		item.ProductName,
		item.VariantName,
		item.ProductSKU,
		item.Quantity,
		item.PricePerUnit,
		item.TotalPrice,
	)

	return err
}

func (s *OrderService) updateOrderTx(ctx context.Context, tx *sqlx.Tx, order *models.StorefrontOrder) error {
	query := `
		UPDATE b2c_orders SET
			subtotal_amount = $1,
			shipping_amount = $2,
			tax_amount = $3,
			total_amount = $4,
			commission_amount = $5,
			seller_amount = $6,
			discount = $7,
			updated_at = NOW()
		WHERE id = $8`

	_, err := tx.ExecContext(ctx, query,
		order.SubtotalAmount,
		order.ShippingAmount,
		order.TaxAmount,
		order.TotalAmount,
		order.CommissionAmount,
		order.SellerAmount,
		order.Discount,
		order.ID,
	)

	return err
}

func (s *OrderService) clearCartTx(ctx context.Context, tx *sqlx.Tx, cartID int64) error {
	// Удаляем позиции корзины
	deleteItemsQuery := `DELETE FROM b2c_cart_items WHERE cart_id = $1`
	_, err := tx.ExecContext(ctx, deleteItemsQuery, cartID)
	if err != nil {
		return err
	}

	// Обновляем корзину
	updateCartQuery := `UPDATE b2c_carts SET updated_at = NOW() WHERE id = $1`
	_, err = tx.ExecContext(ctx, updateCartQuery, cartID)
	return err
}

// updateProductStocksInSearch обновляет остатки товаров в OpenSearch после создания заказа
func (s *OrderService) updateProductStocksInSearch(ctx context.Context, orderItems []models.StorefrontOrderItem) {
	// TODO: OpenSearch integration removed (c2c/b2c deprecated)
	// ProductSearchRepo was removed - stock updates now go directly to PostgreSQL only
	s.logger.Info("ProductSearchRepo removed (deprecated), skipping stock update in OpenSearch")
	_ = ctx
	_ = orderItems
	// Return early - remaining code is disabled

	/* Disabled code - will be restored after OpenSearch refactoring
	// Обновляем остатки для каждого товара в заказе
	for _, item := range orderItems {
		// Получаем актуальную информацию о товаре из БД
		product, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			s.logger.Error("Failed to get product %d for OpenSearch update: %v", item.ProductID, err)
			continue
		}

		// Определяем актуальное количество на складе
		stockQuantity := 0
		if product != nil {
			stockQuantity = product.StockQuantity
		}

		// Если есть вариант, получаем его остатки
		if item.VariantID != nil {
			variant, err := s.productRepo.GetVariantByID(ctx, *item.VariantID)
			if err != nil {
				s.logger.Error("Failed to get variant %d for OpenSearch update: %v", *item.VariantID, err)
			} else if variant != nil {
				stockQuantity = variant.StockQuantity
			}
		}

		// Подготавливаем данные для обновления в OpenSearch
		stockData := map[string]interface{}{
			"stock_quantity": stockQuantity,
			"inventory": map[string]interface{}{
				"quantity":  stockQuantity,
				"in_stock":  stockQuantity > 0,
				"available": stockQuantity,
				"low_stock": stockQuantity > 0 && stockQuantity <= 5,
			},
			"status": "active", // Добавляем статус для совместимости с фильтром поиска
		}

		// Определяем статус наличия
		stockStatus := "in_stock"
		if stockQuantity <= 0 {
			stockStatus = "out_of_stock"
		} else if stockQuantity <= 5 {
			stockStatus = "low_stock"
		}
		stockData["stock_status"] = stockStatus

		// Обновляем данные в OpenSearch
		// TODO: OpenSearch update removed - this code is unreachable (return above)
		// err = s.productSearchRepo.UpdateProductStock(ctx, int(item.ProductID), stockData)
		if false { // Unreachable code - kept for reference
			s.logger.Info("Successfully updated stock for product %d in OpenSearch (new quantity: %d)",
				item.ProductID, stockQuantity)
		}
	}
	*/
}
