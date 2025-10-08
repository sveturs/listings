package service

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"

	"backend/internal/domain/models"
	orderRepo "backend/internal/proj/orders/repository"
)

// OrderServiceTx предоставляет методы для работы с заказами в транзакциях
type OrderServiceTx struct {
	*OrderService
	db *sqlx.DB
}

// NewOrderServiceTx создает новый экземпляр сервиса с поддержкой транзакций
func NewOrderServiceTx(orderService *OrderService, db *sqlx.DB) *OrderServiceTx {
	return &OrderServiceTx{
		OrderService: orderService,
		db:           db,
	}
}

// CreateOrderTx создает новый заказ в транзакции
func (s *OrderServiceTx) CreateOrderTx(ctx context.Context, req *models.CreateOrderRequest, userID int) (*models.StorefrontOrder, error) {
	var createdOrder *models.StorefrontOrder

	err := orderRepo.WithTx(ctx, s.db, func(tx *sqlx.Tx) error {
		s.logger.Info("Creating order in transaction (user_id: %d, storefront_id: %d)", userID, req.StorefrontID)

		// Проверяем существование витрины
		storefront, err := s.storefrontRepo.GetByID(ctx, req.StorefrontID)
		if err != nil {
			return fmt.Errorf("failed to get storefront: %w", err)
		}

		if !storefront.IsActive {
			return fmt.Errorf("storefront is not active")
		}

		// Получаем позиции заказа
		items, err := s.getOrderItems(ctx, req, userID)
		if err != nil {
			return err
		}

		if len(items) == 0 {
			return fmt.Errorf("no items in order")
		}

		// Создаем заказ
		order := s.prepareOrder(req, userID, storefront)

		// Создаем заказ в базе используя транзакцию
		createdOrder, err = s.createOrderInTx(ctx, tx, order)
		if err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		// Обрабатываем позиции заказа и резервируем товары
		orderItems, reservations, err := s.processOrderItemsInTx(ctx, tx, createdOrder, items)
		if err != nil {
			return err
		}

		// Добавляем позиции к заказу
		createdOrder.Items = orderItems

		// Рассчитываем суммы
		s.calculateOrderTotals(ctx, createdOrder, storefront)

		// Обновляем заказ с рассчитанными суммами в транзакции
		if err := s.updateOrderInTx(ctx, tx, createdOrder); err != nil {
			return fmt.Errorf("failed to update order totals: %w", err)
		}

		// Сохраняем информацию о резервированиях
		createdOrder.Metadata["reservations"] = reservations

		// Очищаем корзину если заказ был создан из неё
		if req.CartID != nil {
			if err := s.clearCartInTx(ctx, tx, *req.CartID); err != nil {
				s.logger.Error("Failed to clear cart %d: %v", *req.CartID, err)
				// Не прерываем транзакцию из-за ошибки очистки корзины
			}
		}

		s.logger.Info("Order created successfully (order_id: %d)", createdOrder.ID)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return createdOrder, nil
}

// getOrderItems получает позиции заказа из корзины или запроса
func (s *OrderServiceTx) getOrderItems(ctx context.Context, req *models.CreateOrderRequest, userID int) ([]models.OrderItemRequest, error) {
	var items []models.OrderItemRequest

	if req.CartID != nil {
		// Из корзины
		cart, err := s.cartRepo.GetByID(ctx, *req.CartID)
		if err != nil {
			return nil, fmt.Errorf("failed to get cart: %w", err)
		}

		if cart.UserID == nil || *cart.UserID != userID {
			return nil, fmt.Errorf("cart does not belong to user")
		}

		if cart.StorefrontID != req.StorefrontID {
			return nil, fmt.Errorf("cart belongs to different storefront")
		}

		for _, cartItem := range cart.Items {
			items = append(items, models.OrderItemRequest{
				ProductID: cartItem.ProductID,
				VariantID: cartItem.VariantID,
				Quantity:  cartItem.Quantity,
			})
		}
	} else {
		// Из переданных позиций
		items = req.Items
	}

	return items, nil
}

// prepareOrder подготавливает структуру заказа
func (s *OrderServiceTx) prepareOrder(req *models.CreateOrderRequest, userID int, storefront *models.Storefront) *models.StorefrontOrder {
	return &models.StorefrontOrder{
		StorefrontID:    req.StorefrontID,
		CustomerID:      userID,
		UserID:          userID, // Для совместимости
		Status:          models.OrderStatusPending,
		Currency:        "RSD", // TODO: получать из настроек витрины
		ShippingMethod:  &req.ShippingMethod,
		CustomerNotes:   &req.CustomerNotes,
		EscrowDays:      s.calculateEscrowDays(storefront),
		ShippingAddress: convertToJSONB(req.ShippingAddress),
		BillingAddress:  convertToJSONB(req.BillingAddress),
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

// processOrderItemsInTx обрабатывает позиции заказа в транзакции
func (s *OrderServiceTx) processOrderItemsInTx(
	ctx context.Context,
	tx *sqlx.Tx,
	order *models.StorefrontOrder,
	items []models.OrderItemRequest,
) ([]models.StorefrontOrderItem, []*models.InventoryReservation, error) {
	var orderItems []models.StorefrontOrderItem
	var reservations []*models.InventoryReservation

	for _, item := range items {
		// Получаем товар
		product, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get product %d: %w", item.ProductID, err)
		}

		if !product.IsActive {
			return nil, nil, fmt.Errorf("product %d is not active", item.ProductID)
		}

		// Получаем вариант если указан
		var variant *models.StorefrontProductVariant
		var price decimal.Decimal
		var stockQuantity int

		if item.VariantID != nil {
			variant, err = s.productRepo.GetVariantByID(ctx, *item.VariantID)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to get variant %d: %w", *item.VariantID, err)
			}
			if !variant.IsActive {
				return nil, nil, fmt.Errorf("variant %d is not active", *item.VariantID)
			}
			price = decimal.NewFromFloat(*variant.Price)
			stockQuantity = variant.StockQuantity
		} else {
			price = decimal.NewFromFloat(product.Price)
			stockQuantity = product.StockQuantity
		}

		// Проверяем наличие на складе
		if stockQuantity < item.Quantity {
			return nil, nil, fmt.Errorf("insufficient stock for product %d: requested %d, available %d",
				item.ProductID, item.Quantity, stockQuantity)
		}

		// Резервируем товар в рамках транзакции
		reservation, err := s.reserveStockInTx(ctx, tx, item.ProductID, item.VariantID, item.Quantity, order.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to reserve stock for product %d: %w", item.ProductID, err)
		}
		reservations = append(reservations, reservation)

		// Создаем позицию заказа
		orderItem := models.StorefrontOrderItem{
			OrderID:           order.ID,
			ProductID:         item.ProductID,
			VariantID:         item.VariantID,
			ProductName:       product.Name,
			ProductSKU:        product.SKU,
			Quantity:          item.Quantity,
			PricePerUnit:      price,
			TotalPrice:        price.Mul(decimal.NewFromInt(int64(item.Quantity))),
			ProductAttributes: product.Attributes,
		}

		if variant != nil {
			orderItem.VariantName = nil
		}

		// Сохраняем позицию заказа в транзакции
		if err := s.createOrderItemTx(ctx, tx, &orderItem); err != nil {
			return nil, nil, fmt.Errorf("failed to add order item: %w", err)
		}

		orderItems = append(orderItems, orderItem)
	}

	return orderItems, reservations, nil
}

// createOrderItemTx создает позицию заказа в транзакции
func (s *OrderServiceTx) createOrderItemTx(ctx context.Context, tx *sqlx.Tx, item *models.StorefrontOrderItem) error {
	query := `
		INSERT INTO storefront_order_items (
			order_id, product_id, variant_id, product_name, variant_name,
			product_sku, quantity, price_per_unit, total_price, product_attributes
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
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
		item.ProductAttributes,
	)

	return err
}

// createOrderInTx создает заказ в рамках транзакции
func (s *OrderServiceTx) createOrderInTx(ctx context.Context, tx *sqlx.Tx, order *models.StorefrontOrder) (*models.StorefrontOrder, error) {
	return s.createOrderInTransaction(ctx, tx, order)
}

// updateOrderInTx обновляет заказ в рамках транзакции
func (s *OrderServiceTx) updateOrderInTx(ctx context.Context, tx *sqlx.Tx, order *models.StorefrontOrder) error {
	return s.updateOrderTx(ctx, tx, order)
}

// clearCartInTx очищает корзину в рамках транзакции
func (s *OrderServiceTx) clearCartInTx(ctx context.Context, tx *sqlx.Tx, cartID int64) error {
	return s.clearCartTx(ctx, tx, cartID)
}

// reserveStockInTx резервирует товар в рамках транзакции
func (s *OrderServiceTx) reserveStockInTx(ctx context.Context, tx *sqlx.Tx, productID int64, variantID *int64, quantity int, orderID int64) (*models.InventoryReservation, error) {
	return s.createReservationTx(ctx, tx, productID, variantID, quantity, orderID)
}
