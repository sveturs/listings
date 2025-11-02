package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"

	"backend/internal/config"
	"backend/internal/domain/models"
)

// CreateOrderWithTx создает новый заказ с использованием транзакций для обеспечения целостности данных
func (s *OrderService) CreateOrderWithTx(ctx context.Context, db *sqlx.DB, req *models.CreateOrderRequest, userID int) (*models.StorefrontOrder, error) {
	s.logger.Info("Creating order with transaction (user_id: %d, storefront_id: %d)", userID, req.StorefrontID)

	var createdOrder *models.StorefrontOrder
	var reservations []*models.InventoryReservation

	// Выполняем все операции в транзакции
	err := func() error {
		tx, err := db.BeginTxx(ctx, &sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
		})
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}
		defer func() {
			_ = tx.Rollback() // ignore error - transaction might already be committed
		}()

		// 1. Проверяем существование витрины
		storefront, err := s.getStorefrontTx(ctx, tx, req.StorefrontID)
		if err != nil {
			return fmt.Errorf("failed to get storefront: %w", err)
		}

		if !storefront.IsActive {
			return fmt.Errorf("storefront is not active")
		}

		// 2. Получаем позиции заказа
		items, err := s.getOrderItemsTx(ctx, tx, req, userID)
		if err != nil {
			return err
		}

		if len(items) == 0 {
			return fmt.Errorf("no items in order")
		}

		// 3. Создаем структуру заказа
		order := s.prepareOrderStruct(req, userID, storefront)

		// 4. Создаем заказ в базе данных
		createdOrder, err = s.createOrderInTransaction(ctx, tx, order)
		if err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		// 5. Обрабатываем позиции заказа и резервируем товары
		var orderItems []models.StorefrontOrderItem
		reservations = make([]*models.InventoryReservation, 0, len(items))

		for _, item := range items {
			// Блокируем товар для чтения с блокировкой (SELECT FOR UPDATE)
			product, variant, err := s.lockProductForUpdate(ctx, tx, item.ProductID, item.VariantID)
			if err != nil {
				return fmt.Errorf("failed to lock product %d: %w", item.ProductID, err)
			}

			// Проверяем активность товара
			if !product.IsActive {
				return fmt.Errorf("product %d is not active", item.ProductID)
			}

			// Определяем цену и количество на складе
			var price decimal.Decimal
			var stockQuantity int
			var variantName *string

			if variant != nil {
				if !variant.IsActive {
					return fmt.Errorf("variant %d is not active", *item.VariantID)
				}
				price = decimal.NewFromFloat(*variant.Price)
				stockQuantity = variant.StockQuantity
				variantName = nil
			} else {
				price = decimal.NewFromFloat(product.Price)
				stockQuantity = product.StockQuantity
			}

			// Проверяем наличие на складе
			if stockQuantity < item.Quantity {
				return fmt.Errorf("insufficient stock for product %d: requested %d, available %d",
					item.ProductID, item.Quantity, stockQuantity)
			}

			// Создаем резервирование в рамках транзакции
			reservation, err := s.createReservationTx(ctx, tx, item.ProductID, item.VariantID, item.Quantity, createdOrder.ID)
			if err != nil {
				return fmt.Errorf("failed to reserve stock for product %d: %w", item.ProductID, err)
			}
			reservations = append(reservations, reservation)

			// Обновляем количество товара на складе
			if err := s.updateProductStockTx(ctx, tx, item.ProductID, item.VariantID, stockQuantity-item.Quantity); err != nil {
				return fmt.Errorf("failed to update stock for product %d: %w", item.ProductID, err)
			}

			// Создаем позицию заказа
			orderItem := models.StorefrontOrderItem{
				OrderID:      createdOrder.ID,
				ProductID:    item.ProductID,
				VariantID:    item.VariantID,
				ProductName:  product.Name,
				ProductSKU:   product.SKU,
				Quantity:     item.Quantity,
				PricePerUnit: price,
				TotalPrice:   price.Mul(decimal.NewFromInt(int64(item.Quantity))),
				VariantName:  variantName,
			}

			// Сохраняем позицию заказа
			if err := s.createOrderItemTx(ctx, tx, &orderItem); err != nil {
				return fmt.Errorf("failed to add order item: %w", err)
			}

			orderItems = append(orderItems, orderItem)
		}

		// 6. Добавляем позиции к заказу
		createdOrder.Items = orderItems

		// 7. Рассчитываем суммы заказа
		s.calculateOrderTotals(ctx, createdOrder, storefront)

		// 8. Обновляем заказ с рассчитанными суммами
		if err := s.updateOrderTx(ctx, tx, createdOrder); err != nil {
			return fmt.Errorf("failed to update order totals: %w", err)
		}

		// 9. Очищаем корзину если заказ был создан из неё
		if req.CartID != nil {
			if err := s.clearCartTx(ctx, tx, *req.CartID); err != nil {
				// Логируем ошибку, но не прерываем транзакцию
				s.logger.Error("Failed to clear cart %d: %v", *req.CartID, err)
			}
		}

		// 10. Фиксируем транзакцию
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		// 11. После успешного коммита обновляем остатки в OpenSearch
		// Это делается после транзакции, чтобы не блокировать основной процесс
		go s.updateProductStocksInSearch(ctx, orderItems)

		s.logger.Info("Order created successfully with transaction (order_id: %d)", createdOrder.ID)
		return nil
	}()
	if err != nil {
		// Если произошла ошибка, освобождаем все резервирования
		for _, reservation := range reservations {
			if releaseErr := s.inventoryMgr.ReleaseReservation(ctx, reservation.ID); releaseErr != nil {
				s.logger.Error("Failed to release reservation %d: %v", reservation.ID, releaseErr)
			}
		}
		return nil, err
	}

	// Сохраняем информацию о резервированиях в метаданных
	createdOrder.Metadata["reservations"] = reservations

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

func (s *OrderService) getStorefrontTx(ctx context.Context, tx *sqlx.Tx, storefrontID int) (*models.Storefront, error) {
	var storefront models.Storefront
	query := `SELECT * FROM b2c_stores WHERE id = $1 FOR SHARE`
	err := tx.GetContext(ctx, &storefront, query, storefrontID)
	if err != nil {
		return nil, err
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

func (s *OrderService) lockProductForUpdate(ctx context.Context, tx *sqlx.Tx, productID int64, variantID *int64) (*models.StorefrontProduct, *models.StorefrontProductVariant, error) {
	// Блокируем товар для обновления
	var product models.StorefrontProduct
	productQuery := `SELECT * FROM b2c_products WHERE id = $1 FOR UPDATE`
	err := tx.GetContext(ctx, &product, productQuery, productID)
	if err != nil {
		return nil, nil, err
	}

	// Если есть вариант, блокируем его тоже
	var variant *models.StorefrontProductVariant
	if variantID != nil {
		var v models.StorefrontProductVariant
		variantQuery := `SELECT * FROM b2c_product_variants WHERE id = $1 AND product_id = $2 FOR UPDATE`
		err := tx.GetContext(ctx, &v, variantQuery, *variantID, productID)
		if err != nil {
			return nil, nil, err
		}
		variant = &v
	}

	return &product, variant, nil
}

func (s *OrderService) createOrderInTransaction(ctx context.Context, tx *sqlx.Tx, order *models.StorefrontOrder) (*models.StorefrontOrder, error) {
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
		return nil, err
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

	return &createdOrder, nil
}

func (s *OrderService) createReservationTx(ctx context.Context, tx *sqlx.Tx, productID int64, variantID *int64, quantity int, orderID int64) (*models.InventoryReservation, error) {
	query := `
		INSERT INTO inventory_reservations (
			product_id, variant_id, quantity, order_id, status, expires_at
		) VALUES (
			$1, $2, $3, $4, 'active', NOW() + INTERVAL '30 minutes'
		) RETURNING id, product_id, variant_id, quantity, order_id, status, expires_at, created_at`

	var reservation models.InventoryReservation
	err := tx.QueryRowContext(ctx, query, productID, variantID, quantity, orderID).Scan(
		&reservation.ID,
		&reservation.ProductID,
		&reservation.VariantID,
		&reservation.Quantity,
		&reservation.OrderID,
		&reservation.Status,
		&reservation.ExpiresAt,
		&reservation.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &reservation, nil
}

func (s *OrderService) updateProductStockTx(ctx context.Context, tx *sqlx.Tx, productID int64, variantID *int64, newQuantity int) error {
	if variantID != nil {
		query := `UPDATE b2c_product_variants SET stock_quantity = $1, updated_at = NOW() WHERE id = $2 AND product_id = $3`
		_, err := tx.ExecContext(ctx, query, newQuantity, *variantID, productID)
		return err
	}

	query := `UPDATE b2c_products SET stock_quantity = $1, updated_at = NOW() WHERE id = $2`
	_, err := tx.ExecContext(ctx, query, newQuantity, productID)
	return err
}

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
	// Проверяем наличие OpenSearch репозитория
	if // TODO: s.productSearchRepo == nil {
		s.logger.Info("ProductSearchRepo is not configured, skipping stock update in OpenSearch")
		return
	}

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
		// Для товаров витрин используем префикс sp_ в ID документа
		err = // TODO: s.productSearchRepo.UpdateProductStock(ctx, int(item.ProductID), stockData)
		if err != nil {
			// Если товар не найден с обычным ID, пробуем с префиксом sp_ для товаров витрин
			s.logger.Info("Trying to update storefront product sp_%d in OpenSearch", item.ProductID)
			// Здесь нужно использовать специальный метод для обновления товаров витрин
			// Но пока просто логируем ошибку
			s.logger.Error("Failed to update product %d stock in OpenSearch: %v", item.ProductID, err)
			// Не прерываем процесс, продолжаем с другими товарами
		} else {
			s.logger.Info("Successfully updated stock for product %d in OpenSearch (new quantity: %d)",
				item.ProductID, stockQuantity)
		}
	}
}
