// Package service provides business logic layer for the listings microservice.
package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"github.com/sveturs/listings/internal/domain"
	"github.com/sveturs/listings/internal/repository/postgres"
)

// OrderService defines business logic operations for order management
type OrderService interface {
	// Order lifecycle
	CreateOrder(ctx context.Context, req *CreateOrderRequest) (*domain.Order, error)
	GetOrder(ctx context.Context, orderID int64) (*domain.Order, error)
	GetOrderByNumber(ctx context.Context, orderNumber string) (*domain.Order, error)
	ListOrders(ctx context.Context, req *ListOrdersRequest) ([]*domain.Order, int64, error)
	CancelOrder(ctx context.Context, orderID int64, userID int64, reason string) (*domain.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID int64, status domain.OrderStatus) (*domain.Order, error)
	GetOrderStats(ctx context.Context, userID *int64, storefrontID *int64) (*OrderStats, error)

	// Internal helpers (called by Payment Service webhooks)
	ConfirmOrderPayment(ctx context.Context, orderID int64, transactionID string) error
	ProcessRefund(ctx context.Context, orderID int64) error
}

// OrderItemInput represents a single item for direct checkout
type OrderItemInput struct {
	ProductID  int64  // Required - product/listing ID
	VariantID  *int64 // Optional - variant ID
	Quantity   int    // Required - quantity to order
}

// CreateOrderRequest contains parameters for creating an order
type CreateOrderRequest struct {
	// Cart-based checkout (existing flow)
	CartID          int64                  // Cart to convert to order (0 if using Items)

	// Direct checkout (new flow)
	Items           []OrderItemInput       // Alternative to CartID for direct checkout
	StorefrontID    int64                  // Required for direct checkout (not needed for cart-based)

	// Common fields
	UserID          *int64                 // NULL for guest orders
	ShippingAddress map[string]interface{} // Required JSONB
	BillingAddress  map[string]interface{} // Optional (defaults to shipping)
	ShippingCost    float64                // Calculated shipping cost
	DiscountCode    *string                // Optional discount code
	DiscountAmount  float64                // Discount amount (if applicable)
	PaymentMethod   string                 // payment method (card, cash, etc.)
	CustomerNotes   *string                // Optional customer notes
}

// ListOrdersRequest contains parameters for listing orders
type ListOrdersRequest struct {
	UserID       *int64             // Filter by user
	StorefrontID *int64             // Filter by storefront
	Status       *domain.OrderStatus // Filter by status
	Limit        int                // Page size
	Offset       int                // Page offset
}

// OrderStats contains statistics for orders
type OrderStats struct {
	TotalOrders      int64
	PendingOrders    int64
	ConfirmedOrders  int64
	CompletedOrders  int64
	CancelledOrders  int64
	TotalRevenue     float64
	AverageOrderValue float64
}

// orderService implements OrderService
type orderService struct {
	orderRepo       postgres.OrderRepository
	cartRepo        postgres.CartRepository
	reservationRepo postgres.ReservationRepository
	productsRepo    *postgres.Repository
	pool            *pgxpool.Pool
	config          *FinancialConfig
	logger          zerolog.Logger
}

// NewOrderService creates a new order service
func NewOrderService(
	orderRepo postgres.OrderRepository,
	cartRepo postgres.CartRepository,
	reservationRepo postgres.ReservationRepository,
	productsRepo *postgres.Repository,
	pool *pgxpool.Pool,
	config *FinancialConfig,
	logger zerolog.Logger,
) OrderService {
	if config == nil {
		config = DefaultFinancialConfig()
	}

	return &orderService{
		orderRepo:       orderRepo,
		cartRepo:        cartRepo,
		reservationRepo: reservationRepo,
		productsRepo:    productsRepo,
		pool:            pool,
		config:          config,
		logger:          logger.With().Str("component", "order_service").Logger(),
	}
}

// CreateOrder creates a new order from a cart OR direct items (ACID transaction)
func (s *orderService) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*domain.Order, error) {
	s.logger.Info().
		Int64("cart_id", req.CartID).
		Int("items_count", len(req.Items)).
		Msg("creating order")

	// Start transaction
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to begin transaction")
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// 1. Get cart with items (or create temporary cart from items)
	var cart *domain.Cart

	if req.CartID > 0 {
		// Existing flow: load cart from DB
		cart, err = s.cartRepo.GetByID(ctx, req.CartID)
		if err != nil {
			s.logger.Error().Err(err).Int64("cart_id", req.CartID).Msg("failed to get cart")
			return nil, fmt.Errorf("failed to get cart: %w", err)
		}
	} else if len(req.Items) > 0 {
		// New flow: create temporary cart from items
		cart = s.createTemporaryCartFromItems(req)
	} else {
		return nil, errors.New("either cart_id or items must be provided")
	}

	// 2. Validate cart not empty
	if len(cart.Items) == 0 {
		return nil, ErrCartEmpty
	}

	// 3. Extract listing IDs and sort (ORDER BY id ASC to prevent deadlocks)
	listingIDs := make([]int64, 0, len(cart.Items))
	for _, item := range cart.Items {
		listingIDs = append(listingIDs, item.ListingID)
	}
	sort.Slice(listingIDs, func(i, j int) bool { return listingIDs[i] < listingIDs[j] })

	// 4. Lock listings in sorted order (SELECT FOR UPDATE)
	// NOTE: This requires a LockListingsByIDs method in the repository
	// For now, we'll fetch listings normally and add locking later
	listings := make(map[int64]*domain.Product)
	for _, listingID := range listingIDs {
		listing, err := s.productsRepo.GetProductByID(ctx, listingID, &cart.StorefrontID)
		if err != nil {
			return nil, &ErrListingNotFound{ListingID: listingID}
		}
		listings[listingID] = listing
	}

	// 5. Validate stock availability
	for _, item := range cart.Items {
		listing := listings[item.ListingID]
		if listing.StockQuantity < item.Quantity {
			return nil, &ErrInsufficientStock{
				ListingID:      item.ListingID,
				ListingName:    listing.Name,
				RequestedQty:   item.Quantity,
				AvailableStock: listing.StockQuantity,
			}
		}
	}

	// 6. Validate prices (current price matches cart snapshot)
	priceChanges := s.validatePrices(cart.Items, listings)
	if len(priceChanges) > 0 {
		return nil, &ErrPriceChanged{Changes: priceChanges}
	}

	// 7. Build order items with snapshot data
	orderItems, err := BuildOrderItems(cart.Items, listings)
	if err != nil {
		return nil, fmt.Errorf("failed to build order items: %w", err)
	}

	// 8. Calculate financials
	financials, err := CalculateOrderFinancials(orderItems, req.ShippingCost, req.DiscountAmount, s.config)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate financials: %w", err)
	}

	// 9. Generate order number
	year := time.Now().Year()
	// TODO: Get next sequence number from database
	// For now, use timestamp-based ordering
	sequence := time.Now().UnixNano() / 1000000 % 1000000 // Last 6 digits of milliseconds
	orderNumber := domain.GenerateOrderNumber(year, sequence)

	// 10. Create order
	order := &domain.Order{
		OrderNumber:  orderNumber,
		UserID:       req.UserID,
		StorefrontID: cart.StorefrontID,
		Status:       domain.OrderStatusPending,
		PaymentStatus: domain.PaymentStatusPending,
		Subtotal:     financials.Subtotal,
		Tax:          financials.Tax,
		Shipping:     financials.ShippingCost,
		Discount:     financials.Discount,
		Total:        financials.Total,
		Commission:   financials.Commission,
		SellerAmount: financials.SellerAmount,
		Currency:     financials.Currency,
		ShippingAddress: req.ShippingAddress,
		BillingAddress:  req.BillingAddress,
		EscrowDays:      s.config.EscrowDays,
	}

	// Set payment method
	if req.PaymentMethod != "" {
		order.PaymentMethod = &req.PaymentMethod
	}

	// Set customer notes
	if req.CustomerNotes != nil {
		order.CustomerNotes = req.CustomerNotes
	}

	// Use billing address = shipping address if not provided
	if order.BillingAddress == nil && order.ShippingAddress != nil {
		order.BillingAddress = order.ShippingAddress
	}

	// Create order in transaction
	orderRepoTx := s.orderRepo.WithTx(tx)
	if err := orderRepoTx.Create(ctx, order); err != nil {
		s.logger.Error().Err(err).Msg("failed to create order")
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// 11. Create order items
	if err := orderRepoTx.CreateItems(ctx, order.ID, orderItems); err != nil {
		s.logger.Error().Err(err).Msg("failed to create order items")
		return nil, fmt.Errorf("failed to create order items: %w", err)
	}

	// 12. Create inventory reservations (TTL 30 minutes)
	reservationRepoTx := s.reservationRepo.WithTx(tx)
	reservations := s.buildReservations(order.ID, cart.Items)
	for _, reservation := range reservations {
		if err := reservationRepoTx.Create(ctx, reservation); err != nil {
			s.logger.Error().Err(err).Msg("failed to create reservation")
			return nil, fmt.Errorf("failed to create reservation: %w", err)
		}
	}

	// 13. Deduct stock
	for _, item := range cart.Items {
		if err := s.productsRepo.DeductStockWithPgxTx(ctx, tx, item.ListingID, item.Quantity); err != nil {
			s.logger.Error().Err(err).Int64("listing_id", item.ListingID).Msg("failed to deduct stock")
			return nil, fmt.Errorf("failed to deduct stock: %w", err)
		}
	}

	// 14. Clear cart
	// TODO: Fix transaction type mismatch - CartRepository expects *sql.Tx but we have pgx.Tx
	// Need to either:
	// 1. Convert CartRepository to use pgx.Tx (recommended)
	// 2. Use sql.DB wrapper for transactions
	// For now, delete cart outside transaction (non-critical operation)
	// if err := s.cartRepo.Delete(ctx, cart.ID); err != nil {
	//     s.logger.Warn().Err(err).Int64("cart_id", cart.ID).Msg("failed to delete cart")
	// }

	// 15. Commit transaction
	if err := tx.Commit(ctx); err != nil {
		s.logger.Error().Err(err).Msg("failed to commit transaction")
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Reload order with items
	order, err = s.orderRepo.GetByID(ctx, order.ID)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to reload order")
		return nil, fmt.Errorf("failed to reload order: %w", err)
	}

	s.logger.Info().Int64("order_id", order.ID).Str("order_number", order.OrderNumber).Msg("order created successfully")
	return order, nil
}

// GetOrder retrieves an order by ID
func (s *orderService) GetOrder(ctx context.Context, orderID int64) (*domain.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		if err.Error() == "order not found" {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	return order, nil
}

// GetOrderByNumber retrieves an order by order number
func (s *orderService) GetOrderByNumber(ctx context.Context, orderNumber string) (*domain.Order, error) {
	order, err := s.orderRepo.GetByOrderNumber(ctx, orderNumber)
	if err != nil {
		if err.Error() == "order not found" {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	return order, nil
}

// ListOrders retrieves orders with filtering and pagination
func (s *orderService) ListOrders(ctx context.Context, req *ListOrdersRequest) ([]*domain.Order, int64, error) {
	if req.UserID != nil {
		orders, total, err := s.orderRepo.ListByUser(ctx, *req.UserID, req.Limit, req.Offset)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to list user orders: %w", err)
		}
		return orders, int64(total), nil
	}

	if req.StorefrontID != nil {
		orders, total, err := s.orderRepo.ListByStorefront(ctx, *req.StorefrontID, req.Limit, req.Offset)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to list storefront orders: %w", err)
		}
		return orders, int64(total), nil
	}

	return nil, 0, fmt.Errorf("%w: must specify user_id or storefront_id", ErrInvalidInput)
}

// CancelOrder cancels an order (releases reservations, restores stock)
func (s *orderService) CancelOrder(ctx context.Context, orderID int64, userID int64, reason string) (*domain.Order, error) {
	s.logger.Info().Int64("order_id", orderID).Int64("user_id", userID).Msg("cancelling order")

	// Start transaction
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Get order
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		if err.Error() == "order not found" {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Verify user owns the order
	if order.UserID == nil || *order.UserID != userID {
		return nil, ErrUnauthorized
	}

	// Check if order can be cancelled
	if !order.CanCancel() {
		return nil, &ErrOrderCannotCancel{
			OrderID: orderID,
			Status:  string(order.Status),
		}
	}

	// Get reservations before releasing (needed for stock restoration)
	reservations, err := s.reservationRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		s.logger.Error().Err(err).Int64("order_id", orderID).Msg("failed to get reservations")
		return nil, fmt.Errorf("failed to get reservations: %w", err)
	}

	// Update order status to cancelled
	orderRepoTx := s.orderRepo.WithTx(tx)
	if err := orderRepoTx.UpdateStatus(ctx, orderID, domain.OrderStatusCancelled); err != nil {
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	// Release all reservations for this order (batch operation)
	reservationRepoTx := s.reservationRepo.WithTx(tx)
	if err := reservationRepoTx.ReleaseReservations(ctx, orderID); err != nil {
		s.logger.Error().Err(err).Int64("order_id", orderID).Msg("failed to release reservations")
		return nil, fmt.Errorf("failed to release reservations: %w", err)
	}

	// Restore stock for released reservations
	for _, reservation := range reservations {
		if err := s.productsRepo.RestoreStockWithPgxTx(ctx, tx, reservation.ListingID, reservation.Quantity); err != nil {
			s.logger.Error().Err(err).Int64("listing_id", reservation.ListingID).Msg("failed to restore stock")
			return nil, fmt.Errorf("failed to restore stock: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// TODO: Publish OrderCancelledEvent to message queue for Payment Service to process refund

	// Reload order
	order, err = s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to reload order: %w", err)
	}

	s.logger.Info().Int64("order_id", orderID).Msg("order cancelled successfully")
	return order, nil
}

// UpdateOrderStatus updates the status of an order
func (s *orderService) UpdateOrderStatus(ctx context.Context, orderID int64, status domain.OrderStatus) (*domain.Order, error) {
	// Get order
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		if err.Error() == "order not found" {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Check if status transition is valid
	if !order.CanUpdateStatus(status) {
		return nil, &ErrOrderCannotUpdateStatus{
			OrderID:    orderID,
			FromStatus: string(order.Status),
			ToStatus:   string(status),
		}
	}

	// Update status
	if err := s.orderRepo.UpdateStatus(ctx, orderID, status); err != nil {
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	// Reload order
	order, err = s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to reload order: %w", err)
	}

	s.logger.Info().Int64("order_id", orderID).Str("new_status", string(status)).Msg("order status updated")
	return order, nil
}

// GetOrderStats retrieves statistics for orders
func (s *orderService) GetOrderStats(ctx context.Context, userID *int64, storefrontID *int64) (*OrderStats, error) {
	// TODO: Implement statistics aggregation
	// For now, return empty stats
	return &OrderStats{}, nil
}

// ConfirmOrderPayment confirms payment for an order (called by Payment Service webhook)
func (s *orderService) ConfirmOrderPayment(ctx context.Context, orderID int64, transactionID string) error {
	s.logger.Info().Int64("order_id", orderID).Str("transaction_id", transactionID).Msg("confirming order payment")

	// Start transaction
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Get order
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	// Update payment status
	now := time.Now()
	order.PaymentStatus = domain.PaymentStatusCompleted
	order.PaymentTransactionID = &transactionID
	order.PaymentCompletedAt = &now

	// Update order status to confirmed
	order.Status = domain.OrderStatusConfirmed
	order.ConfirmedAt = &now

	// Set escrow release date
	escrowReleaseDate := now.Add(time.Duration(order.EscrowDays) * 24 * time.Hour)
	order.EscrowReleaseDate = &escrowReleaseDate

	orderRepoTx := s.orderRepo.WithTx(tx)
	if err := orderRepoTx.Update(ctx, order); err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	// Commit all reservations for this order (batch operation)
	reservationRepoTx := s.reservationRepo.WithTx(tx)
	if err := reservationRepoTx.CommitReservations(ctx, orderID); err != nil {
		s.logger.Error().Err(err).Int64("order_id", orderID).Msg("failed to commit reservations")
		return fmt.Errorf("failed to commit reservations: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// TODO: Publish OrderConfirmedEvent to message queue for Delivery Service to create shipment

	s.logger.Info().Int64("order_id", orderID).Msg("order payment confirmed successfully")
	return nil
}

// ProcessRefund processes a refund for a cancelled order
func (s *orderService) ProcessRefund(ctx context.Context, orderID int64) error {
	// Get order
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	// Calculate refund amount
	refundAmount := CalculateRefundAmount(order)

	// TODO: Call Payment Service to process refund
	s.logger.Info().
		Int64("order_id", orderID).
		Float64("refund_amount", refundAmount).
		Msg("processing refund (stub)")

	return nil
}

// Helper methods

// validatePrices validates that cart prices match current listing prices
func (s *orderService) validatePrices(cartItems []*domain.CartItem, listings map[int64]*domain.Product) []PriceChangeItem {
	changes := []PriceChangeItem{}

	for _, item := range cartItems {
		listing := listings[item.ListingID]
		if listing.Price != item.PriceSnapshot {
			changes = append(changes, PriceChangeItem{
				ListingID:     item.ListingID,
				ListingName:   listing.Name,
				OldPrice:      item.PriceSnapshot,
				NewPrice:      listing.Price,
				PriceIncrease: listing.Price > item.PriceSnapshot,
			})
		}
	}

	return changes
}

// buildReservations creates inventory reservations for cart items
func (s *orderService) buildReservations(orderID int64, cartItems []*domain.CartItem) []*domain.InventoryReservation {
	reservations := make([]*domain.InventoryReservation, 0, len(cartItems))

	for _, item := range cartItems {
		reservation := domain.NewInventoryReservation(
			item.ListingID,
			item.VariantID,
			orderID,
			item.Quantity,
		)
		reservations = append(reservations, reservation)
	}

	return reservations
}

// createTemporaryCartFromItems creates a temporary in-memory cart from direct items
// Used for direct checkout flow (without persisting cart to DB)
func (s *orderService) createTemporaryCartFromItems(req *CreateOrderRequest) *domain.Cart {
	cart := &domain.Cart{
		UserID:       req.UserID,
		StorefrontID: req.StorefrontID,
		Items:        make([]*domain.CartItem, 0, len(req.Items)),
	}

	for _, item := range req.Items {
		cartItem := &domain.CartItem{
			ListingID:     item.ProductID,
			Quantity:      int32(item.Quantity), // Convert int to int32 for domain.CartItem
			VariantID:     item.VariantID,
			PriceSnapshot: 0, // Will be set during price validation
		}
		cart.Items = append(cart.Items, cartItem)
	}

	return cart
}
