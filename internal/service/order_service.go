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

	// Seller shipment workflow
	AcceptOrder(ctx context.Context, orderID int64, sellerID int64, sellerNotes string) (*domain.Order, error)
	CreateOrderShipment(ctx context.Context, req *CreateShipmentRequest) (*CreateShipmentResult, error)
	MarkOrderShipped(ctx context.Context, orderID int64, sellerID int64, sellerNotes string) (*domain.Order, error)
	GetOrderTracking(ctx context.Context, orderID int64, userID int64) (*TrackingInfo, error)

	// Internal helpers (called by Payment Service webhooks)
	ConfirmOrderPayment(ctx context.Context, orderID int64, transactionID string) error
	ProcessRefund(ctx context.Context, orderID int64) error

	// Configuration
	SetChatService(chatService ChatService)
	SetDeliveryClient(client DeliveryClient)
}

// OrderItemInput represents a single item for direct checkout
type OrderItemInput struct {
	ProductID int64  // Required - product/listing ID
	VariantID *int64 // Optional - variant ID
	Quantity  int    // Required - quantity to order
}

// CreateOrderRequest contains parameters for creating an order
type CreateOrderRequest struct {
	// Cart-based checkout (existing flow)
	CartID int64 // Cart to convert to order (0 if using Items)

	// Direct checkout (new flow)
	Items        []OrderItemInput // Alternative to CartID for direct checkout
	StorefrontID int64            // Required for direct checkout (not needed for cart-based)

	// Common fields
	UserID             *int64                 // NULL for guest orders
	ShippingAddress    map[string]interface{} // Required JSONB
	BillingAddress     map[string]interface{} // Optional (defaults to shipping)
	ShippingCost       float64                // Calculated shipping cost
	DiscountCode       *string                // Optional discount code
	DiscountAmount     float64                // Discount amount (if applicable)
	PaymentMethod      string                 // payment method (card, cash, etc.)
	CustomerNotes      *string                // Optional customer notes
	AcceptPriceChanges bool                   // If true, skip price validation (for direct checkout)

	// Customer contact info (for seller to see on sales page)
	CustomerName  *string // Customer full name
	CustomerEmail *string // Customer email
	CustomerPhone *string // Customer phone
}

// ListOrdersRequest contains parameters for listing orders
type ListOrdersRequest struct {
	UserID       *int64              // Filter by user
	StorefrontID *int64              // Filter by storefront
	Status       *domain.OrderStatus // Filter by status
	Limit        int                 // Page size
	Offset       int                 // Page offset
}

// OrderStats contains statistics for orders
type OrderStats struct {
	TotalOrders       int64
	PendingOrders     int64
	ConfirmedOrders   int64
	CompletedOrders   int64
	CancelledOrders   int64
	TotalRevenue      float64
	AverageOrderValue float64
}

// CreateShipmentRequest contains parameters for creating a shipment
type CreateShipmentRequest struct {
	OrderID        int64
	SellerID       int64
	ProviderCode   string      // post_express, bex_express, aks, d_express, city_express
	PackageInfo    PackageInfo // Package dimensions and weight
	UseCOD         bool        // Cash on delivery
	CODAmount      float64     // COD amount (if UseCOD = true)
	UseInsurance   bool        // Insurance for package
	InsuranceValue float64     // Declared value for insurance
}

// PackageInfo contains package dimensions and weight
type PackageInfo struct {
	WeightKg    float64 // Weight in kg
	LengthCm    float64 // Length in cm
	WidthCm     float64 // Width in cm
	HeightCm    float64 // Height in cm
	IsFragile   bool    // Fragile goods flag
	Description string  // Package contents description
}

// CreateShipmentResult contains the result of creating a shipment
type CreateShipmentResult struct {
	Order             *domain.Order
	ShipmentID        int64
	TrackingNumber    string
	Provider          string
	Status            string
	DeliveryCost      float64
	EstimatedDelivery string // RFC3339 formatted date
	LabelURL          string // URL to shipping label PDF
}

// TrackingInfo contains tracking information for an order
type TrackingInfo struct {
	TrackingNumber    string
	Provider          string
	Status            string
	EstimatedDelivery string          // RFC3339 formatted date
	Events            []TrackingEvent // Timeline of events
}

// TrackingEvent represents a single tracking event
type TrackingEvent struct {
	Status      string    // pending, picked_up, in_transit, delivered, etc.
	Location    string    // Location description
	Description string    // Event description
	Timestamp   time.Time // Event time
}

// DeliveryClient defines the interface for delivery microservice client
type DeliveryClient interface {
	CreateShipment(ctx context.Context, req *DeliveryCreateShipmentRequest) (*DeliveryShipment, error)
	TrackShipment(ctx context.Context, trackingNumber string) (*DeliveryTrackingInfo, error)
	CalculateRate(ctx context.Context, req *DeliveryCalculateRateRequest) (*DeliveryRateInfo, error)
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
	chatService     ChatService    // For sending order notifications
	deliveryClient  DeliveryClient // For delivery microservice integration
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

// SetChatService sets the chat service for sending order notifications
// This allows delayed initialization to avoid circular dependencies
func (s *orderService) SetChatService(chatService ChatService) {
	s.chatService = chatService
}

// SetDeliveryClient sets the delivery microservice client
// This allows delayed initialization to avoid circular dependencies
func (s *orderService) SetDeliveryClient(client DeliveryClient) {
	s.deliveryClient = client
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
	// Skip price validation for direct checkout when AcceptPriceChanges is true
	// This is necessary because direct checkout doesn't have price snapshots in cart
	if !req.AcceptPriceChanges {
		priceChanges := s.validatePrices(cart.Items, listings)
		if len(priceChanges) > 0 {
			return nil, &ErrPriceChanged{Changes: priceChanges}
		}
	} else {
		// For direct checkout with AcceptPriceChanges=true,
		// update cart item price snapshots to current prices
		for _, item := range cart.Items {
			if listing, ok := listings[item.ListingID]; ok {
				item.PriceSnapshot = listing.Price
			}
		}
	}

	// 7. Build temporary order items for financial calculations
	// NOTE: These items don't have order_id yet - that will be set after order creation
	tempOrderItems, err := BuildOrderItems(cart.Items, listings)
	if err != nil {
		return nil, fmt.Errorf("failed to build order items: %w", err)
	}

	// 8. Calculate financials
	financials, err := CalculateOrderFinancials(tempOrderItems, req.ShippingCost, req.DiscountAmount, s.config)
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
		OrderNumber:     orderNumber,
		UserID:          req.UserID,
		StorefrontID:    cart.StorefrontID,
		Status:          domain.OrderStatusPending,
		PaymentStatus:   domain.PaymentStatusPending,
		Subtotal:        financials.Subtotal,
		Tax:             financials.Tax,
		Shipping:        financials.ShippingCost,
		Discount:        financials.Discount,
		Total:           financials.Total,
		Commission:      financials.Commission,
		SellerAmount:    financials.SellerAmount,
		Currency:        financials.Currency,
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

	// Set customer contact info (for seller display on sales page)
	// First try from explicit fields, then fallback to shipping_address
	if req.CustomerName != nil && *req.CustomerName != "" {
		order.CustomerName = req.CustomerName
	} else if req.ShippingAddress != nil {
		if fullName, ok := req.ShippingAddress["full_name"].(string); ok && fullName != "" {
			order.CustomerName = &fullName
		}
	}

	if req.CustomerEmail != nil && *req.CustomerEmail != "" {
		order.CustomerEmail = req.CustomerEmail
	} else if req.ShippingAddress != nil {
		if email, ok := req.ShippingAddress["email"].(string); ok && email != "" {
			order.CustomerEmail = &email
		}
	}

	if req.CustomerPhone != nil && *req.CustomerPhone != "" {
		order.CustomerPhone = req.CustomerPhone
	} else if req.ShippingAddress != nil {
		if phone, ok := req.ShippingAddress["phone"].(string); ok && phone != "" {
			order.CustomerPhone = &phone
		}
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

	// 11. Build final order items and set order_id (now available from database)
	finalOrderItems, err := BuildOrderItems(cart.Items, listings)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to build order items")
		return nil, fmt.Errorf("failed to build order items: %w", err)
	}

	// Set order_id on all items (order.ID is now populated from DB auto-increment)
	for _, item := range finalOrderItems {
		item.OrderID = order.ID
	}

	// Create order items in database
	if err := orderRepoTx.CreateItems(ctx, order.ID, finalOrderItems); err != nil {
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

	// Send system notification to storefront owner about new order
	s.notifyStorefrontOwnerAboutOrder(ctx, order)

	// Auto-confirm order for cash-on-delivery (COD) orders
	// COD orders should immediately move to "confirmed" status since payment
	// will be collected upon delivery, not upfront.
	// Payment status is set to "cod_pending" - NOT "completed" because payment hasn't happened yet!
	if order.PaymentMethod != nil && isCashOnDeliveryMethod(*order.PaymentMethod) {
		if err := s.confirmCODOrder(ctx, order); err != nil {
			s.logger.Error().Err(err).
				Int64("order_id", order.ID).
				Str("payment_method", *order.PaymentMethod).
				Msg("failed to auto-confirm COD order")
			// Don't fail the order creation, just log the error
			// The order is still valid, seller can manually confirm
		} else {
			// Reload order with updated status
			order, err = s.orderRepo.GetByID(ctx, order.ID)
			if err != nil {
				s.logger.Warn().Err(err).Int64("order_id", order.ID).Msg("failed to reload order after COD confirmation")
			}
			s.logger.Info().
				Int64("order_id", order.ID).
				Str("order_number", order.OrderNumber).
				Str("payment_method", *order.PaymentMethod).
				Str("payment_status", string(order.PaymentStatus)).
				Msg("COD order auto-confirmed successfully")
		}
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

// notifyStorefrontOwnerAboutOrder sends a system message to the storefront owner
// notifying them about a new order. This is a non-blocking async operation.
func (s *orderService) notifyStorefrontOwnerAboutOrder(ctx context.Context, order *domain.Order) {
	if s.chatService == nil {
		s.logger.Warn().Msg("chat service not configured, skipping order notification")
		return
	}

	// Get storefront to find owner
	storefront, err := s.productsRepo.GetStorefrontByID(ctx, order.StorefrontID, nil)
	if err != nil {
		s.logger.Error().Err(err).Int64("storefront_id", order.StorefrontID).Msg("failed to get storefront for notification")
		return
	}

	// Build notification message (multilingual - will use buyer's language if available)
	itemCount := len(order.Items)
	message := fmt.Sprintf(
		"🛒 New order #%s!\n\n"+
			"Items: %d\n"+
			"Total: %.2f %s\n\n"+
			"Please prepare the order for shipping.",
		order.OrderNumber,
		itemCount,
		order.Total,
		order.Currency,
	)

	// Send system message asynchronously to not block order creation
	go func() {
		// Create new context since parent context may be cancelled
		notifyCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		req := &SendSystemMessageRequest{
			ReceiverID:       storefront.UserID,
			Content:          message,
			OriginalLanguage: "en", // Default to English, could be improved with user preferences
		}

		_, err := s.chatService.SendSystemMessage(notifyCtx, req)
		if err != nil {
			s.logger.Error().Err(err).
				Int64("storefront_id", order.StorefrontID).
				Int64("owner_id", storefront.UserID).
				Str("order_number", order.OrderNumber).
				Msg("failed to send order notification to storefront owner")
		} else {
			s.logger.Info().
				Int64("storefront_id", order.StorefrontID).
				Int64("owner_id", storefront.UserID).
				Str("order_number", order.OrderNumber).
				Msg("order notification sent to storefront owner")
		}
	}()
}

// ============================================================================
// SELLER SHIPMENT WORKFLOW METHODS
// ============================================================================

// AcceptOrder handles seller accepting an order
// Flow: confirmed -> accepted
// Validates: order status, seller is storefront owner
func (s *orderService) AcceptOrder(ctx context.Context, orderID int64, sellerID int64, sellerNotes string) (*domain.Order, error) {
	s.logger.Info().
		Int64("order_id", orderID).
		Int64("seller_id", sellerID).
		Msg("accepting order")

	// Get order
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		if err.Error() == "order not found" {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Verify seller is storefront owner
	storefront, err := s.productsRepo.GetStorefrontByID(ctx, order.StorefrontID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get storefront: %w", err)
	}
	if storefront.UserID != sellerID {
		s.logger.Warn().
			Int64("order_id", orderID).
			Int64("seller_id", sellerID).
			Int64("storefront_owner_id", storefront.UserID).
			Msg("seller is not storefront owner")
		return nil, ErrUnauthorized
	}

	// Check if order can be accepted
	if !order.CanAccept() {
		return nil, &ErrOrderInvalidStatus{
			OrderID:        orderID,
			CurrentStatus:  string(order.Status),
			ExpectedStatus: string(domain.OrderStatusConfirmed),
			Action:         "accept",
		}
	}

	// Update order
	now := time.Now()
	order.Status = domain.OrderStatusAccepted
	order.AcceptedAt = &now
	if sellerNotes != "" {
		order.SellerNotes = &sellerNotes
	}

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	// Reload order with items
	order, err = s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to reload order: %w", err)
	}

	// Send notification to buyer about order acceptance
	s.notifyBuyerAboutOrderAccepted(ctx, order, storefront)

	s.logger.Info().Int64("order_id", orderID).Msg("order accepted successfully")
	return order, nil
}

// CreateOrderShipment creates shipment via Delivery Service
// Flow: accepted -> processing
// Validates: order status, seller is storefront owner, package info
func (s *orderService) CreateOrderShipment(ctx context.Context, req *CreateShipmentRequest) (*CreateShipmentResult, error) {
	s.logger.Info().
		Int64("order_id", req.OrderID).
		Int64("seller_id", req.SellerID).
		Str("provider_code", req.ProviderCode).
		Msg("creating order shipment")

	// Get order
	order, err := s.orderRepo.GetByID(ctx, req.OrderID)
	if err != nil {
		if err.Error() == "order not found" {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Verify seller is storefront owner
	storefront, err := s.productsRepo.GetStorefrontByID(ctx, order.StorefrontID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get storefront: %w", err)
	}
	if storefront.UserID != req.SellerID {
		s.logger.Warn().
			Int64("order_id", req.OrderID).
			Int64("seller_id", req.SellerID).
			Int64("storefront_owner_id", storefront.UserID).
			Msg("seller is not storefront owner")
		return nil, ErrUnauthorized
	}

	// Check if order can have shipment created
	if !order.CanCreateShipment() {
		return nil, &ErrOrderInvalidStatus{
			OrderID:        req.OrderID,
			CurrentStatus:  string(order.Status),
			ExpectedStatus: string(domain.OrderStatusAccepted),
			Action:         "create shipment",
		}
	}

	// Parse delivery provider from code
	provider, err := s.parseDeliveryProvider(req.ProviderCode)
	if err != nil {
		return nil, &ErrInvalidDeliveryProvider{ProviderCode: req.ProviderCode}
	}

	// Prepare shipment data
	var shipmentID int64
	var trackingNumber, labelURL, estimatedDelivery string
	var deliveryCost float64

	// Call Delivery Service if configured
	if s.deliveryClient != nil {
		// Build addresses from order data
		fromAddress := s.buildSellerAddress(storefront)
		toAddress := s.buildBuyerAddress(order)

		// Build package info
		pkg := s.buildDeliveryPackage(&req.PackageInfo, order, req.UseCOD, req.CODAmount, req.UseInsurance, req.InsuranceValue)

		// Create shipment via Delivery Service
		deliveryReq := &DeliveryCreateShipmentRequest{
			Provider:    provider,
			FromAddress: fromAddress,
			ToAddress:   toAddress,
			Package:     pkg,
			UserID:      fmt.Sprintf("%d", req.SellerID),
		}

		shipment, err := s.deliveryClient.CreateShipment(ctx, deliveryReq)
		if err != nil {
			s.logger.Error().Err(err).
				Int64("order_id", req.OrderID).
				Str("provider_code", req.ProviderCode).
				Msg("failed to create shipment via delivery service")
			return nil, &ErrShipmentCreationFailed{
				OrderID: req.OrderID,
				Reason:  err.Error(),
			}
		}

		// Parse shipment ID from string to int64
		shipmentID, _ = parseShipmentID(shipment.ID)
		trackingNumber = shipment.TrackingNumber
		labelURL = "" // Will be provided when label is generated
		if !shipment.EstimatedDelivery.IsZero() {
			estimatedDelivery = shipment.EstimatedDelivery.Format(time.RFC3339)
		} else {
			estimatedDelivery = time.Now().Add(72 * time.Hour).Format(time.RFC3339)
		}
		deliveryCost = parseCost(shipment.Cost)

		s.logger.Info().
			Str("shipment_id", shipment.ID).
			Str("tracking_number", trackingNumber).
			Msg("shipment created via delivery service")
	} else {
		// Fallback: generate mock shipment data if delivery client not configured
		s.logger.Warn().Msg("delivery client not configured, using mock shipment data")
		shipmentID = time.Now().UnixNano() / 1000000 % 10000000
		trackingNumber = fmt.Sprintf("TRK%d%06d", time.Now().Year(), shipmentID%1000000)
		labelURL = fmt.Sprintf("https://delivery.svetu.rs/labels/%s.pdf", trackingNumber)
		estimatedDelivery = time.Now().Add(72 * time.Hour).Format(time.RFC3339)
		deliveryCost = order.Shipping
	}

	// Update order with shipment info
	order.Status = domain.OrderStatusProcessing
	order.TrackingNumber = &trackingNumber
	order.ShippingProvider = &req.ProviderCode
	order.ShipmentID = &shipmentID
	if labelURL != "" {
		order.LabelURL = &labelURL
	}

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	// Reload order with items
	order, err = s.orderRepo.GetByID(ctx, req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to reload order: %w", err)
	}

	// Send notification to buyer about shipment creation
	s.notifyBuyerAboutShipmentCreated(ctx, order, storefront, trackingNumber)

	s.logger.Info().
		Int64("order_id", req.OrderID).
		Str("tracking_number", trackingNumber).
		Msg("order shipment created successfully")

	return &CreateShipmentResult{
		Order:             order,
		ShipmentID:        shipmentID,
		TrackingNumber:    trackingNumber,
		Provider:          req.ProviderCode,
		Status:            "label_created",
		DeliveryCost:      deliveryCost,
		EstimatedDelivery: estimatedDelivery,
		LabelURL:          labelURL,
	}, nil
}

// MarkOrderShipped marks order as shipped
// Flow: processing -> shipped
// Validates: order status, tracking number exists
func (s *orderService) MarkOrderShipped(ctx context.Context, orderID int64, sellerID int64, sellerNotes string) (*domain.Order, error) {
	s.logger.Info().
		Int64("order_id", orderID).
		Int64("seller_id", sellerID).
		Msg("marking order as shipped")

	// Get order
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		if err.Error() == "order not found" {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Verify seller is storefront owner
	storefront, err := s.productsRepo.GetStorefrontByID(ctx, order.StorefrontID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get storefront: %w", err)
	}
	if storefront.UserID != sellerID {
		s.logger.Warn().
			Int64("order_id", orderID).
			Int64("seller_id", sellerID).
			Int64("storefront_owner_id", storefront.UserID).
			Msg("seller is not storefront owner")
		return nil, ErrUnauthorized
	}

	// Check if order can be marked as shipped
	if !order.CanMarkShipped() {
		// Provide more specific error message
		if order.Status != domain.OrderStatusProcessing {
			return nil, &ErrOrderInvalidStatus{
				OrderID:        orderID,
				CurrentStatus:  string(order.Status),
				ExpectedStatus: string(domain.OrderStatusProcessing),
				Action:         "mark shipped",
			}
		}
		if order.TrackingNumber == nil {
			return nil, &ErrOrderMissingTrackingNumber{OrderID: orderID}
		}
	}

	// Update order
	now := time.Now()
	order.Status = domain.OrderStatusShipped
	order.ShippedAt = &now
	if sellerNotes != "" {
		if order.SellerNotes != nil {
			combinedNotes := *order.SellerNotes + "\n" + sellerNotes
			order.SellerNotes = &combinedNotes
		} else {
			order.SellerNotes = &sellerNotes
		}
	}

	if err := s.orderRepo.Update(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	// Reload order with items
	order, err = s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to reload order: %w", err)
	}

	// Send notification to buyer about shipment dispatch
	s.notifyBuyerAboutOrderShipped(ctx, order, storefront)

	s.logger.Info().Int64("order_id", orderID).Msg("order marked as shipped successfully")
	return order, nil
}

// GetOrderTracking gets tracking info from Delivery Service
// Returns: tracking events timeline
func (s *orderService) GetOrderTracking(ctx context.Context, orderID int64, userID int64) (*TrackingInfo, error) {
	s.logger.Debug().
		Int64("order_id", orderID).
		Int64("user_id", userID).
		Msg("getting order tracking")

	// Get order
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		if err.Error() == "order not found" {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Verify user has access (either buyer or seller)
	hasAccess := false
	if order.UserID != nil && *order.UserID == userID {
		hasAccess = true
	} else {
		// Check if user is storefront owner
		storefront, err := s.productsRepo.GetStorefrontByID(ctx, order.StorefrontID, nil)
		if err == nil && storefront.UserID == userID {
			hasAccess = true
		}
	}

	if !hasAccess {
		s.logger.Warn().
			Int64("order_id", orderID).
			Int64("user_id", userID).
			Msg("user does not have access to order tracking")
		return nil, ErrUnauthorized
	}

	// Check if order has tracking info
	if order.TrackingNumber == nil {
		return nil, &ErrOrderMissingTrackingNumber{OrderID: orderID}
	}

	// Try to get tracking from Delivery Service
	if s.deliveryClient != nil {
		deliveryInfo, err := s.deliveryClient.TrackShipment(ctx, *order.TrackingNumber)
		if err == nil && deliveryInfo != nil {
			// Convert delivery tracking info to service tracking info
			return s.convertDeliveryTrackingToService(order, deliveryInfo), nil
		}
		// Log error but continue with fallback
		s.logger.Warn().Err(err).
			Int64("order_id", orderID).
			Str("tracking_number", *order.TrackingNumber).
			Msg("failed to get tracking from delivery service, using local data")
	}

	// Fallback: build tracking data from order timestamps
	tracking := &TrackingInfo{
		TrackingNumber: *order.TrackingNumber,
		Provider:       "delivery_service",
		Status:         string(order.Status),
		Events:         []TrackingEvent{},
	}

	if order.ShippingProvider != nil {
		tracking.Provider = *order.ShippingProvider
	}

	// Build events based on order timestamps
	if order.CreatedAt.Unix() > 0 {
		tracking.Events = append(tracking.Events, TrackingEvent{
			Status:      "order_created",
			Location:    "System",
			Description: "Order created",
			Timestamp:   order.CreatedAt,
		})
	}

	if order.ConfirmedAt != nil {
		tracking.Events = append(tracking.Events, TrackingEvent{
			Status:      "payment_confirmed",
			Location:    "System",
			Description: "Payment confirmed",
			Timestamp:   *order.ConfirmedAt,
		})
	}

	if order.AcceptedAt != nil {
		tracking.Events = append(tracking.Events, TrackingEvent{
			Status:      "accepted",
			Location:    "Seller",
			Description: "Order accepted by seller",
			Timestamp:   *order.AcceptedAt,
		})
	}

	if order.Status == domain.OrderStatusProcessing || order.Status == domain.OrderStatusShipped || order.Status == domain.OrderStatusDelivered {
		// Add shipment created event (estimate based on accepted_at + 1 hour)
		shipmentCreatedTime := order.UpdatedAt
		if order.AcceptedAt != nil {
			shipmentCreatedTime = order.AcceptedAt.Add(time.Hour)
		}
		tracking.Events = append(tracking.Events, TrackingEvent{
			Status:      "shipment_created",
			Location:    "Seller",
			Description: "Shipment label created",
			Timestamp:   shipmentCreatedTime,
		})
	}

	if order.ShippedAt != nil {
		tracking.Events = append(tracking.Events, TrackingEvent{
			Status:      "shipped",
			Location:    "Courier",
			Description: "Package handed to courier",
			Timestamp:   *order.ShippedAt,
		})
		// Estimate delivery 2-3 days after shipping
		tracking.EstimatedDelivery = order.ShippedAt.Add(72 * time.Hour).Format(time.RFC3339)
	}

	if order.DeliveredAt != nil {
		tracking.Events = append(tracking.Events, TrackingEvent{
			Status:      "delivered",
			Location:    "Destination",
			Description: "Package delivered",
			Timestamp:   *order.DeliveredAt,
		})
		tracking.Status = "delivered"
	}

	return tracking, nil
}

// isCashOnDeliveryMethod checks if payment method is cash-on-delivery (COD)
// COD orders should be auto-confirmed since payment is collected upon delivery
func isCashOnDeliveryMethod(method string) bool {
	switch method {
	case "cash_on_delivery", "cod", "cash", "pouzecem", "pouzećem":
		return true
	default:
		return false
	}
}

// confirmCODOrder confirms a Cash-on-Delivery order without marking payment as completed.
// COD orders move to "confirmed" status but payment_status stays as "cod_pending"
// because actual payment will happen at delivery time.
func (s *orderService) confirmCODOrder(ctx context.Context, order *domain.Order) error {
	s.logger.Info().
		Int64("order_id", order.ID).
		Str("order_number", order.OrderNumber).
		Msg("confirming COD order")

	// Start transaction
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Update order status to confirmed, but payment_status to cod_pending (NOT completed!)
	now := time.Now()
	order.Status = domain.OrderStatusConfirmed
	order.PaymentStatus = domain.PaymentStatusCODPending
	order.ConfirmedAt = &now

	// Set escrow release date (for when payment is actually collected at delivery)
	escrowReleaseDate := now.Add(time.Duration(order.EscrowDays) * 24 * time.Hour)
	order.EscrowReleaseDate = &escrowReleaseDate

	// COD transaction ID for tracking
	codTransactionID := "COD-" + order.OrderNumber
	order.PaymentTransactionID = &codTransactionID

	orderRepoTx := s.orderRepo.WithTx(tx)
	if err := orderRepoTx.Update(ctx, order); err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	// Commit reservations (stock already deducted, now mark as committed)
	reservationRepoTx := s.reservationRepo.WithTx(tx)
	if err := reservationRepoTx.CommitReservations(ctx, order.ID); err != nil {
		s.logger.Warn().Err(err).Int64("order_id", order.ID).Msg("failed to commit reservations for COD order")
		// Don't fail - reservations are not critical for COD
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info().
		Int64("order_id", order.ID).
		Str("payment_status", string(order.PaymentStatus)).
		Str("order_status", string(order.Status)).
		Msg("COD order confirmed successfully")

	return nil
}

// ============================================================================
// SHIPMENT WORKFLOW HELPER METHODS
// ============================================================================

// parseDeliveryProvider converts provider code string to DeliveryProvider enum
func (s *orderService) parseDeliveryProvider(code string) (DeliveryProvider, error) {
	switch code {
	case "post_express":
		return DeliveryProviderPostExpress, nil
	case "bex_express":
		return DeliveryProviderBexExpress, nil
	case "aks", "aks_express":
		return DeliveryProviderAksExpress, nil
	case "d_express":
		return DeliveryProviderDExpress, nil
	case "city_express":
		return DeliveryProviderCityExpress, nil
	default:
		return DeliveryProviderUnspecified, fmt.Errorf("unknown provider: %s", code)
	}
}

// buildSellerAddress builds DeliveryAddress from storefront data
func (s *orderService) buildSellerAddress(storefront *domain.Storefront) *DeliveryAddress {
	addr := &DeliveryAddress{
		Country: "RS", // Default to Serbia
	}

	// Set contact name from storefront
	if storefront.Name != "" {
		addr.ContactName = storefront.Name
	}

	// Get contact phone from storefront
	if storefront.Phone != nil {
		addr.ContactPhone = *storefront.Phone
	}

	// Get address fields from storefront
	if storefront.Address != nil {
		addr.Street = *storefront.Address
	}
	if storefront.City != nil {
		addr.City = *storefront.City
	}
	if storefront.PostalCode != nil {
		addr.PostalCode = *storefront.PostalCode
	}
	if storefront.Country != nil && *storefront.Country != "" {
		addr.Country = *storefront.Country
	}

	return addr
}

// buildBuyerAddress builds DeliveryAddress from order shipping address
func (s *orderService) buildBuyerAddress(order *domain.Order) *DeliveryAddress {
	addr := &DeliveryAddress{
		Country: "RS", // Default to Serbia
	}

	if order.ShippingAddress == nil {
		return addr
	}

	// Extract address components from shipping_address JSONB
	if street, ok := order.ShippingAddress["street"].(string); ok {
		addr.Street = street
	}
	if city, ok := order.ShippingAddress["city"].(string); ok {
		addr.City = city
	}
	if postalCode, ok := order.ShippingAddress["postal_code"].(string); ok {
		addr.PostalCode = postalCode
	}
	if country, ok := order.ShippingAddress["country"].(string); ok && country != "" {
		addr.Country = country
	}
	if state, ok := order.ShippingAddress["state"].(string); ok {
		addr.State = state
	}

	// Get contact info
	if order.CustomerName != nil {
		addr.ContactName = *order.CustomerName
	} else if fullName, ok := order.ShippingAddress["full_name"].(string); ok {
		addr.ContactName = fullName
	}

	if order.CustomerPhone != nil {
		addr.ContactPhone = *order.CustomerPhone
	} else if phone, ok := order.ShippingAddress["phone"].(string); ok {
		addr.ContactPhone = phone
	}

	return addr
}

// buildDeliveryPackage builds DeliveryPackage from package info and order
func (s *orderService) buildDeliveryPackage(info *PackageInfo, order *domain.Order, useCOD bool, codAmount float64, useInsurance bool, insuranceValue float64) *DeliveryPackage {
	pkg := &DeliveryPackage{
		Description: info.Description,
	}

	// Format dimensions as strings (delivery service expects strings)
	if info.WeightKg > 0 {
		pkg.Weight = fmt.Sprintf("%.2f", info.WeightKg)
	} else {
		pkg.Weight = "1.0" // Default 1kg
	}

	if info.LengthCm > 0 {
		pkg.Length = fmt.Sprintf("%.0f", info.LengthCm)
	}
	if info.WidthCm > 0 {
		pkg.Width = fmt.Sprintf("%.0f", info.WidthCm)
	}
	if info.HeightCm > 0 {
		pkg.Height = fmt.Sprintf("%.0f", info.HeightCm)
	}

	// Set declared value for insurance
	if useInsurance && insuranceValue > 0 {
		pkg.DeclaredValue = fmt.Sprintf("%.2f", insuranceValue)
	} else if useCOD && codAmount > 0 {
		pkg.DeclaredValue = fmt.Sprintf("%.2f", codAmount)
	} else {
		pkg.DeclaredValue = fmt.Sprintf("%.2f", order.Total)
	}

	// Add description if not provided
	if pkg.Description == "" {
		itemCount := len(order.Items)
		if itemCount > 0 {
			pkg.Description = fmt.Sprintf("Order #%s (%d items)", order.OrderNumber, itemCount)
		} else {
			pkg.Description = fmt.Sprintf("Order #%s", order.OrderNumber)
		}
	}

	return pkg
}

// convertDeliveryTrackingToService converts delivery tracking info to service tracking info
func (s *orderService) convertDeliveryTrackingToService(order *domain.Order, deliveryInfo *DeliveryTrackingInfo) *TrackingInfo {
	tracking := &TrackingInfo{
		TrackingNumber: *order.TrackingNumber,
		Status:         string(order.Status),
		Events:         []TrackingEvent{},
	}

	// Set provider from shipment info
	if deliveryInfo.Shipment != nil {
		tracking.Provider = deliveryProviderToString(deliveryInfo.Shipment.Provider)
		if !deliveryInfo.Shipment.EstimatedDelivery.IsZero() {
			tracking.EstimatedDelivery = deliveryInfo.Shipment.EstimatedDelivery.Format(time.RFC3339)
		}
	} else if order.ShippingProvider != nil {
		tracking.Provider = *order.ShippingProvider
	}

	// Convert delivery events to service events
	for _, e := range deliveryInfo.Events {
		tracking.Events = append(tracking.Events, TrackingEvent{
			Status:      e.Status,
			Location:    e.Location,
			Description: e.Description,
			Timestamp:   e.Timestamp,
		})
	}

	// Add order-level events at the beginning
	orderEvents := []TrackingEvent{}

	if order.CreatedAt.Unix() > 0 {
		orderEvents = append(orderEvents, TrackingEvent{
			Status:      "order_created",
			Location:    "System",
			Description: "Order created",
			Timestamp:   order.CreatedAt,
		})
	}

	if order.ConfirmedAt != nil {
		orderEvents = append(orderEvents, TrackingEvent{
			Status:      "payment_confirmed",
			Location:    "System",
			Description: "Payment confirmed",
			Timestamp:   *order.ConfirmedAt,
		})
	}

	if order.AcceptedAt != nil {
		orderEvents = append(orderEvents, TrackingEvent{
			Status:      "accepted",
			Location:    "Seller",
			Description: "Order accepted by seller",
			Timestamp:   *order.AcceptedAt,
		})
	}

	// Prepend order events to delivery events
	tracking.Events = append(orderEvents, tracking.Events...)

	return tracking
}

// deliveryProviderToString converts DeliveryProvider enum to human-readable string
func deliveryProviderToString(p DeliveryProvider) string {
	switch p {
	case DeliveryProviderPostExpress:
		return "Post Express"
	case DeliveryProviderBexExpress:
		return "BEX Express"
	case DeliveryProviderAksExpress:
		return "AKS Express"
	case DeliveryProviderDExpress:
		return "D Express"
	case DeliveryProviderCityExpress:
		return "City Express"
	default:
		return "Unknown"
	}
}

// parseShipmentID parses shipment ID from string to int64
func parseShipmentID(id string) (int64, error) {
	var shipmentID int64
	_, err := fmt.Sscanf(id, "%d", &shipmentID)
	return shipmentID, err
}

// parseCost parses cost from string to float64
func parseCost(cost string) float64 {
	var c float64
	fmt.Sscanf(cost, "%f", &c)
	return c
}

// ============================================================================
// BUYER NOTIFICATION METHODS
// ============================================================================

// notifyBuyerAboutOrderAccepted sends notification to buyer about order acceptance
func (s *orderService) notifyBuyerAboutOrderAccepted(ctx context.Context, order *domain.Order, storefront *domain.Storefront) {
	if s.chatService == nil || order.UserID == nil {
		return
	}

	storefrontName := storefront.Name
	if storefrontName == "" && order.StorefrontName != nil {
		storefrontName = *order.StorefrontName
	}

	message := fmt.Sprintf(
		"Your order #%s has been accepted by %s.\n\n"+
			"The seller is preparing your items for shipment. "+
			"You will receive tracking information once the package is shipped.",
		order.OrderNumber,
		storefrontName,
	)

	s.sendBuyerNotification(ctx, *order.UserID, message, order.OrderNumber)
}

// notifyBuyerAboutShipmentCreated sends notification to buyer about shipment creation
func (s *orderService) notifyBuyerAboutShipmentCreated(ctx context.Context, order *domain.Order, storefront *domain.Storefront, trackingNumber string) {
	if s.chatService == nil || order.UserID == nil {
		return
	}

	storefrontName := storefront.Name
	if storefrontName == "" && order.StorefrontName != nil {
		storefrontName = *order.StorefrontName
	}

	provider := "the courier"
	if order.ShippingProvider != nil {
		provider = *order.ShippingProvider
	}

	message := fmt.Sprintf(
		"A shipping label has been created for your order #%s.\n\n"+
			"Tracking Number: %s\n"+
			"Carrier: %s\n\n"+
			"The seller will hand over your package to the courier soon.",
		order.OrderNumber,
		trackingNumber,
		provider,
	)

	s.sendBuyerNotification(ctx, *order.UserID, message, order.OrderNumber)
}

// notifyBuyerAboutOrderShipped sends notification to buyer about order shipment
func (s *orderService) notifyBuyerAboutOrderShipped(ctx context.Context, order *domain.Order, storefront *domain.Storefront) {
	if s.chatService == nil || order.UserID == nil {
		return
	}

	storefrontName := storefront.Name
	if storefrontName == "" && order.StorefrontName != nil {
		storefrontName = *order.StorefrontName
	}

	trackingInfo := ""
	if order.TrackingNumber != nil {
		provider := "courier"
		if order.ShippingProvider != nil {
			provider = *order.ShippingProvider
		}
		trackingInfo = fmt.Sprintf("\nTracking Number: %s\nCarrier: %s", *order.TrackingNumber, provider)
	}

	message := fmt.Sprintf(
		"Great news! Your order #%s has been shipped by %s.\n%s\n\n"+
			"You can track your package status in the order details page. "+
			"Expected delivery: 2-3 business days.",
		order.OrderNumber,
		storefrontName,
		trackingInfo,
	)

	s.sendBuyerNotification(ctx, *order.UserID, message, order.OrderNumber)
}

// sendBuyerNotification sends a system message to the buyer
func (s *orderService) sendBuyerNotification(ctx context.Context, buyerID int64, message, orderNumber string) {
	go func() {
		// Create new context since parent context may be cancelled
		notifyCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		req := &SendSystemMessageRequest{
			ReceiverID:       buyerID,
			Content:          message,
			OriginalLanguage: "en",
		}

		_, err := s.chatService.SendSystemMessage(notifyCtx, req)
		if err != nil {
			s.logger.Error().Err(err).
				Int64("buyer_id", buyerID).
				Str("order_number", orderNumber).
				Msg("failed to send order notification to buyer")
		} else {
			s.logger.Debug().
				Int64("buyer_id", buyerID).
				Str("order_number", orderNumber).
				Msg("order notification sent to buyer")
		}
	}()
}
