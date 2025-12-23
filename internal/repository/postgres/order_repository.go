// Package postgres implements PostgreSQL repository layer for listings microservice.
package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"github.com/vondi-global/listings/internal/domain"
)

// UpdatePaymentInfoParams contains fields for updating payment information
type UpdatePaymentInfoParams struct {
	PaymentProvider       *string
	PaymentSessionID      *string
	PaymentIntentID       *string
	PaymentIdempotencyKey *string
	PaymentStatus         *string
	PaymentTransactionID  *string
}

// OrderRepository defines operations for order management
type OrderRepository interface {
	// Order operations
	Create(ctx context.Context, order *domain.Order) error
	GetByID(ctx context.Context, orderID int64) (*domain.Order, error)
	GetByOrderNumber(ctx context.Context, orderNumber string) (*domain.Order, error)
	ListByUser(ctx context.Context, userID int64, limit, offset int) ([]*domain.Order, int, error)
	ListByStorefront(ctx context.Context, storefrontID int64, limit, offset int) ([]*domain.Order, int, error)
	Update(ctx context.Context, order *domain.Order) error
	UpdateStatus(ctx context.Context, orderID int64, status domain.OrderStatus) error
	UpdatePaymentInfo(ctx context.Context, orderID int64, params UpdatePaymentInfoParams) error
	Delete(ctx context.Context, orderID int64) error

	// Order item operations
	CreateItems(ctx context.Context, orderID int64, items []*domain.OrderItem) error
	GetItems(ctx context.Context, orderID int64) ([]*domain.OrderItem, error)

	// Locking (for ACID transactions)
	LockListingsByIDs(ctx context.Context, listingIDs []int64) error

	// Transaction support
	WithTx(tx pgx.Tx) OrderRepository
}

// orderRepository implements OrderRepository using PostgreSQL
type orderRepository struct {
	db     dbOrTx
	logger zerolog.Logger
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(pool *pgxpool.Pool, logger zerolog.Logger) OrderRepository {
	return &orderRepository{
		db:     pool,
		logger: logger.With().Str("component", "order_repository").Logger(),
	}
}

// WithTx returns a new repository instance using the provided transaction
func (r *orderRepository) WithTx(tx pgx.Tx) OrderRepository {
	return &orderRepository{
		db:     tx,
		logger: r.logger,
	}
}

// Create creates a new order
func (r *orderRepository) Create(ctx context.Context, order *domain.Order) error {
	if err := order.Validate(); err != nil {
		return fmt.Errorf("invalid order: %w", err)
	}

	// Marshal JSONB addresses
	shippingAddressJSON, err := json.Marshal(order.ShippingAddress)
	if err != nil {
		return fmt.Errorf("failed to marshal shipping address: %w", err)
	}

	var billingAddressJSON []byte
	if order.BillingAddress != nil {
		billingAddressJSON, err = json.Marshal(order.BillingAddress)
		if err != nil {
			return fmt.Errorf("failed to marshal billing address: %w", err)
		}
	}

	query := `
		INSERT INTO orders (
			order_number, user_id, storefront_id, status, payment_status,
			subtotal, tax, shipping, discount, total, commission, seller_amount, currency,
			payment_method, payment_transaction_id, payment_completed_at,
			shipping_address, billing_address, shipping_method, shipping_provider, tracking_number, shipment_id,
			escrow_release_date, escrow_days,
			customer_name, customer_email, customer_phone,
			notes, admin_notes
		)
		VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10, $11, $12, $13,
			$14, $15, $16,
			$17, $18, $19, $20, $21, $22,
			$23, $24,
			$25, $26, $27,
			$28, $29
		)
		RETURNING id, created_at, updated_at
	`

	err = r.db.QueryRow(ctx, query,
		order.OrderNumber, order.UserID, order.StorefrontID, string(order.Status), string(order.PaymentStatus),
		order.Subtotal, order.Tax, order.Shipping, order.Discount, order.Total, order.Commission, order.SellerAmount, order.Currency,
		order.PaymentMethod, order.PaymentTransactionID, order.PaymentCompletedAt,
		shippingAddressJSON, billingAddressJSON, order.ShippingMethod, order.ShippingProvider, order.TrackingNumber, order.ShipmentID,
		order.EscrowReleaseDate, order.EscrowDays,
		order.CustomerName, order.CustomerEmail, order.CustomerPhone,
		order.CustomerNotes, order.AdminNotes,
	).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)

	if err != nil {
		r.logger.Error().Err(err).Str("order_number", order.OrderNumber).Msg("failed to create order")
		return fmt.Errorf("failed to create order: %w", err)
	}

	r.logger.Info().Int64("order_id", order.ID).Str("order_number", order.OrderNumber).Msg("order created")
	return nil
}

// GetByID retrieves an order by its ID
func (r *orderRepository) GetByID(ctx context.Context, orderID int64) (*domain.Order, error) {
	query := `
		SELECT o.id, o.order_number, o.user_id, o.storefront_id, o.status, o.payment_status,
		       o.subtotal, o.tax, o.shipping, o.discount, o.total, o.commission, o.seller_amount, o.currency,
		       o.payment_method, o.payment_transaction_id, o.payment_completed_at,
		       o.payment_provider, o.payment_session_id, o.payment_intent_id, o.payment_idempotency_key,
		       o.shipping_address, o.billing_address, o.shipping_method_id, o.shipping_provider, o.tracking_number, o.shipment_id,
		       o.escrow_release_date, o.escrow_days,
		       o.customer_name, o.customer_email, o.customer_phone,
		       o.customer_notes, o.admin_notes, o.seller_notes,
		       o.created_at, o.updated_at, o.confirmed_at, o.accepted_at, o.shipped_at, o.delivered_at, o.cancelled_at,
		       o.label_url,
		       s.name as storefront_name
		FROM orders o
		LEFT JOIN storefronts s ON o.storefront_id = s.id
		WHERE o.id = $1
	`

	var order domain.Order
	var userID sql.NullInt64
	var statusStr, paymentStatusStr string
	var paymentMethod, paymentTransactionID sql.NullString
	var paymentProvider, paymentSessionID, paymentIntentID, paymentIdempotencyKey sql.NullString
	var paymentCompletedAt, escrowReleaseDate, confirmedAt, acceptedAt, shippedAt, deliveredAt, cancelledAt sql.NullTime
	var shippingAddressJSON, billingAddressJSON []byte
	var shippingMethod, shippingProvider, trackingNumber sql.NullString
	var shipmentID sql.NullInt64
	var customerName, customerEmail, customerPhone, customerNotes, adminNotes, sellerNotes sql.NullString
	var labelURL sql.NullString
	var storefrontName sql.NullString

	err := r.db.QueryRow(ctx, query, orderID).Scan(
		&order.ID, &order.OrderNumber, &userID, &order.StorefrontID, &statusStr, &paymentStatusStr,
		&order.Subtotal, &order.Tax, &order.Shipping, &order.Discount, &order.Total, &order.Commission, &order.SellerAmount, &order.Currency,
		&paymentMethod, &paymentTransactionID, &paymentCompletedAt,
		&paymentProvider, &paymentSessionID, &paymentIntentID, &paymentIdempotencyKey,
		&shippingAddressJSON, &billingAddressJSON, &shippingMethod, &shippingProvider, &trackingNumber, &shipmentID,
		&escrowReleaseDate, &order.EscrowDays,
		&customerName, &customerEmail, &customerPhone,
		&customerNotes, &adminNotes, &sellerNotes,
		&order.CreatedAt, &order.UpdatedAt, &confirmedAt, &acceptedAt, &shippedAt, &deliveredAt, &cancelledAt,
		&labelURL,
		&storefrontName,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		r.logger.Error().Err(err).Int64("order_id", orderID).Msg("failed to get order by ID")
		return nil, fmt.Errorf("failed to get order by ID: %w", err)
	}

	// Handle nullable fields
	if userID.Valid {
		order.UserID = &userID.Int64
	}
	order.Status = domain.OrderStatus(statusStr)
	order.PaymentStatus = domain.PaymentStatus(paymentStatusStr)

	if paymentMethod.Valid {
		order.PaymentMethod = &paymentMethod.String
	}
	if paymentTransactionID.Valid {
		order.PaymentTransactionID = &paymentTransactionID.String
	}
	if paymentCompletedAt.Valid {
		order.PaymentCompletedAt = &paymentCompletedAt.Time
	}
	if paymentProvider.Valid {
		order.PaymentProvider = &paymentProvider.String
	}
	if paymentSessionID.Valid {
		order.PaymentSessionID = &paymentSessionID.String
	}
	if paymentIntentID.Valid {
		order.PaymentIntentID = &paymentIntentID.String
	}
	if paymentIdempotencyKey.Valid {
		order.PaymentIdempotencyKey = &paymentIdempotencyKey.String
	}

	// Parse JSONB addresses
	if len(shippingAddressJSON) > 0 {
		if err := json.Unmarshal(shippingAddressJSON, &order.ShippingAddress); err != nil {
			r.logger.Error().Err(err).Msg("failed to unmarshal shipping address")
		}
	}
	if len(billingAddressJSON) > 0 {
		if err := json.Unmarshal(billingAddressJSON, &order.BillingAddress); err != nil {
			r.logger.Error().Err(err).Msg("failed to unmarshal billing address")
		}
	}

	if shippingMethod.Valid {
		order.ShippingMethod = &shippingMethod.String
	}
	if shippingProvider.Valid {
		order.ShippingProvider = &shippingProvider.String
	}
	if trackingNumber.Valid {
		order.TrackingNumber = &trackingNumber.String
	}
	if shipmentID.Valid {
		order.ShipmentID = &shipmentID.Int64
	}

	if escrowReleaseDate.Valid {
		order.EscrowReleaseDate = &escrowReleaseDate.Time
	}

	if customerName.Valid {
		order.CustomerName = &customerName.String
	}
	if customerEmail.Valid {
		order.CustomerEmail = &customerEmail.String
	}
	if customerPhone.Valid {
		order.CustomerPhone = &customerPhone.String
	}
	if customerNotes.Valid {
		order.CustomerNotes = &customerNotes.String
	}
	if adminNotes.Valid {
		order.AdminNotes = &adminNotes.String
	}
	if sellerNotes.Valid {
		order.SellerNotes = &sellerNotes.String
	}

	if confirmedAt.Valid {
		order.ConfirmedAt = &confirmedAt.Time
	}
	if acceptedAt.Valid {
		order.AcceptedAt = &acceptedAt.Time
	}
	if shippedAt.Valid {
		order.ShippedAt = &shippedAt.Time
	}
	if deliveredAt.Valid {
		order.DeliveredAt = &deliveredAt.Time
	}
	if cancelledAt.Valid {
		order.CancelledAt = &cancelledAt.Time
	}
	if labelURL.Valid {
		order.LabelURL = &labelURL.String
	}
	if storefrontName.Valid {
		order.StorefrontName = &storefrontName.String
	}

	// Load order items
	items, err := r.GetItems(ctx, order.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load order items: %w", err)
	}
	order.Items = items

	return &order, nil
}

// GetByOrderNumber retrieves an order by its order number
func (r *orderRepository) GetByOrderNumber(ctx context.Context, orderNumber string) (*domain.Order, error) {
	query := `SELECT id FROM orders WHERE order_number = $1`

	var orderID int64
	err := r.db.QueryRow(ctx, query, orderNumber).Scan(&orderID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		r.logger.Error().Err(err).Str("order_number", orderNumber).Msg("failed to get order by order number")
		return nil, fmt.Errorf("failed to get order by order number: %w", err)
	}

	return r.GetByID(ctx, orderID)
}

// ListByUser retrieves orders for a user with pagination
func (r *orderRepository) ListByUser(ctx context.Context, userID int64, limit, offset int) ([]*domain.Order, int, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM orders WHERE user_id = $1`
	var totalCount int
	err := r.db.QueryRow(ctx, countQuery, userID).Scan(&totalCount)
	if err != nil {
		r.logger.Error().Err(err).Int64("user_id", userID).Msg("failed to count user orders")
		return nil, 0, fmt.Errorf("failed to count user orders: %w", err)
	}

	// Get orders
	query := `
		SELECT id FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		r.logger.Error().Err(err).Int64("user_id", userID).Msg("failed to list user orders")
		return nil, 0, fmt.Errorf("failed to list user orders: %w", err)
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		var orderID int64
		if err := rows.Scan(&orderID); err != nil {
			return nil, 0, fmt.Errorf("failed to scan order ID: %w", err)
		}

		order, err := r.GetByID(ctx, orderID)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get order: %w", err)
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating order rows: %w", err)
	}

	return orders, totalCount, nil
}

// ListByStorefront retrieves orders for a storefront with pagination
func (r *orderRepository) ListByStorefront(ctx context.Context, storefrontID int64, limit, offset int) ([]*domain.Order, int, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM orders WHERE storefront_id = $1`
	var totalCount int
	err := r.db.QueryRow(ctx, countQuery, storefrontID).Scan(&totalCount)
	if err != nil {
		r.logger.Error().Err(err).Int64("storefront_id", storefrontID).Msg("failed to count storefront orders")
		return nil, 0, fmt.Errorf("failed to count storefront orders: %w", err)
	}

	// Get orders
	query := `
		SELECT id FROM orders
		WHERE storefront_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, storefrontID, limit, offset)
	if err != nil {
		r.logger.Error().Err(err).Int64("storefront_id", storefrontID).Msg("failed to list storefront orders")
		return nil, 0, fmt.Errorf("failed to list storefront orders: %w", err)
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		var orderID int64
		if err := rows.Scan(&orderID); err != nil {
			return nil, 0, fmt.Errorf("failed to scan order ID: %w", err)
		}

		order, err := r.GetByID(ctx, orderID)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get order: %w", err)
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating order rows: %w", err)
	}

	return orders, totalCount, nil
}

// Update updates an existing order
func (r *orderRepository) Update(ctx context.Context, order *domain.Order) error {
	if err := order.Validate(); err != nil {
		return fmt.Errorf("invalid order: %w", err)
	}

	// Marshal JSONB addresses
	shippingAddressJSON, err := json.Marshal(order.ShippingAddress)
	if err != nil {
		return fmt.Errorf("failed to marshal shipping address: %w", err)
	}

	var billingAddressJSON []byte
	if order.BillingAddress != nil {
		billingAddressJSON, err = json.Marshal(order.BillingAddress)
		if err != nil {
			return fmt.Errorf("failed to marshal billing address: %w", err)
		}
	}

	query := `
		UPDATE orders SET
			status = $1, payment_status = $2,
			subtotal = $3, tax = $4, shipping = $5, discount = $6, total = $7, commission = $8, seller_amount = $9,
			payment_method = $10, payment_transaction_id = $11, payment_completed_at = $12,
			payment_provider = $13, payment_session_id = $14, payment_intent_id = $15, payment_idempotency_key = $16,
			shipping_address = $17, billing_address = $18, shipping_method = $19, shipping_provider = $20, tracking_number = $21, shipment_id = $22,
			escrow_release_date = $23, escrow_days = $24,
			customer_name = $25, customer_email = $26, customer_phone = $27,
			notes = $28, admin_notes = $29, seller_notes = $30,
			confirmed_at = $31, accepted_at = $32, shipped_at = $33, delivered_at = $34, cancelled_at = $35,
			label_url = $36
		WHERE id = $37
		RETURNING updated_at
	`

	err = r.db.QueryRow(ctx, query,
		string(order.Status), string(order.PaymentStatus),
		order.Subtotal, order.Tax, order.Shipping, order.Discount, order.Total, order.Commission, order.SellerAmount,
		order.PaymentMethod, order.PaymentTransactionID, order.PaymentCompletedAt,
		order.PaymentProvider, order.PaymentSessionID, order.PaymentIntentID, order.PaymentIdempotencyKey,
		shippingAddressJSON, billingAddressJSON, order.ShippingMethod, order.ShippingProvider, order.TrackingNumber, order.ShipmentID,
		order.EscrowReleaseDate, order.EscrowDays,
		order.CustomerName, order.CustomerEmail, order.CustomerPhone,
		order.CustomerNotes, order.AdminNotes, order.SellerNotes,
		order.ConfirmedAt, order.AcceptedAt, order.ShippedAt, order.DeliveredAt, order.CancelledAt,
		order.LabelURL,
		order.ID,
	).Scan(&order.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("order not found")
		}
		r.logger.Error().Err(err).Int64("order_id", order.ID).Msg("failed to update order")
		return fmt.Errorf("failed to update order: %w", err)
	}

	r.logger.Info().Int64("order_id", order.ID).Msg("order updated")
	return nil
}

// UpdateStatus updates the status of an order
func (r *orderRepository) UpdateStatus(ctx context.Context, orderID int64, status domain.OrderStatus) error {
	query := `
		UPDATE orders SET status = $1
		WHERE id = $2
		RETURNING updated_at
	`

	var updatedAt sql.NullTime
	err := r.db.QueryRow(ctx, query, string(status), orderID).Scan(&updatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("order not found")
		}
		r.logger.Error().Err(err).Int64("order_id", orderID).Msg("failed to update order status")
		return fmt.Errorf("failed to update order status: %w", err)
	}

	r.logger.Info().Int64("order_id", orderID).Str("status", string(status)).Msg("order status updated")
	return nil
}

// UpdatePaymentInfo updates payment information for an order
func (r *orderRepository) UpdatePaymentInfo(ctx context.Context, orderID int64, params UpdatePaymentInfoParams) error {
	// Build dynamic UPDATE query based on non-nil params
	setClauses := []string{}
	args := []interface{}{}
	argPos := 1

	if params.PaymentProvider != nil {
		setClauses = append(setClauses, fmt.Sprintf("payment_provider = $%d", argPos))
		args = append(args, *params.PaymentProvider)
		argPos++
	}

	if params.PaymentSessionID != nil {
		setClauses = append(setClauses, fmt.Sprintf("payment_session_id = $%d", argPos))
		args = append(args, *params.PaymentSessionID)
		argPos++
	}

	if params.PaymentIntentID != nil {
		setClauses = append(setClauses, fmt.Sprintf("payment_intent_id = $%d", argPos))
		args = append(args, *params.PaymentIntentID)
		argPos++
	}

	if params.PaymentIdempotencyKey != nil {
		setClauses = append(setClauses, fmt.Sprintf("payment_idempotency_key = $%d", argPos))
		args = append(args, *params.PaymentIdempotencyKey)
		argPos++
	}

	if params.PaymentStatus != nil {
		setClauses = append(setClauses, fmt.Sprintf("payment_status = $%d", argPos))
		args = append(args, *params.PaymentStatus)
		argPos++
	}

	if params.PaymentTransactionID != nil {
		setClauses = append(setClauses, fmt.Sprintf("payment_transaction_id = $%d", argPos))
		args = append(args, *params.PaymentTransactionID)
		argPos++
	}

	// If no fields to update, return error
	if len(setClauses) == 0 {
		return fmt.Errorf("no payment fields to update")
	}

	// Add order_id as last parameter
	args = append(args, orderID)
	whereClause := fmt.Sprintf("WHERE id = $%d", argPos)

	query := fmt.Sprintf(`
		UPDATE orders SET %s
		%s
		RETURNING updated_at
	`, strings.Join(setClauses, ", "), whereClause)

	var updatedAt sql.NullTime
	err := r.db.QueryRow(ctx, query, args...).Scan(&updatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("order not found")
		}
		r.logger.Error().Err(err).Int64("order_id", orderID).Msg("failed to update payment info")
		return fmt.Errorf("failed to update payment info: %w", err)
	}

	r.logger.Info().Int64("order_id", orderID).Msg("payment info updated")
	return nil
}

// Delete deletes an order by ID
func (r *orderRepository) Delete(ctx context.Context, orderID int64) error {
	query := `DELETE FROM orders WHERE id = $1`

	result, err := r.db.Exec(ctx, query, orderID)
	if err != nil {
		r.logger.Error().Err(err).Int64("order_id", orderID).Msg("failed to delete order")
		return fmt.Errorf("failed to delete order: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("order not found")
	}

	r.logger.Info().Int64("order_id", orderID).Msg("order deleted")
	return nil
}

// CreateItems creates order items in batch
func (r *orderRepository) CreateItems(ctx context.Context, orderID int64, items []*domain.OrderItem) error {
	if len(items) == 0 {
		return fmt.Errorf("no items to create")
	}

	query := `
		INSERT INTO order_items (
			order_id, listing_id, variant_id, variant_uuid, stock_reservation_id,
			listing_name, sku,
			variant_data, attributes, quantity, price, subtotal, discount, total, image_url
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id, created_at
	`

	batch := &pgx.Batch{}
	for _, item := range items {
		if err := item.Validate(); err != nil {
			return fmt.Errorf("invalid order item: %w", err)
		}

		// Marshal JSONB data
		variantDataJSON, err := json.Marshal(item.VariantData)
		if err != nil {
			return fmt.Errorf("failed to marshal variant data: %w", err)
		}
		attributesJSON, err := json.Marshal(item.Attributes)
		if err != nil {
			return fmt.Errorf("failed to marshal attributes: %w", err)
		}

		batch.Queue(query,
			orderID, item.ListingID, item.VariantID, item.VariantUUID, item.StockReservationID,
			item.ListingName, item.SKU,
			variantDataJSON, attributesJSON, item.Quantity, item.UnitPrice, item.Subtotal, item.Discount, item.Total, item.ImageURL,
		)
	}

	br := r.db.SendBatch(ctx, batch)
	defer br.Close()

	for i := 0; i < len(items); i++ {
		err := br.QueryRow().Scan(&items[i].ID, &items[i].CreatedAt)
		if err != nil {
			r.logger.Error().Err(err).Int("item_index", i).Msg("failed to create order item")
			return fmt.Errorf("failed to create order item: %w", err)
		}
	}

	r.logger.Info().Int64("order_id", orderID).Int("items_count", len(items)).Msg("order items created")
	return nil
}

// GetItems retrieves all items for an order
func (r *orderRepository) GetItems(ctx context.Context, orderID int64) ([]*domain.OrderItem, error) {
	query := `
		SELECT id, order_id, listing_id, variant_id, stock_reservation_id,
		       listing_name, sku,
		       variant_data, attributes, quantity, price, subtotal, discount, total, image_url, created_at
		FROM order_items
		WHERE order_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		r.logger.Error().Err(err).Int64("order_id", orderID).Msg("failed to get order items")
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}
	defer rows.Close()

	var items []*domain.OrderItem
	for rows.Next() {
		var item domain.OrderItem
		var variantID sql.NullInt64
		var sku, imageURL, stockReservationID sql.NullString
		var variantDataJSON, attributesJSON []byte

		err := rows.Scan(
			&item.ID, &item.OrderID, &item.ListingID, &variantID, &stockReservationID,
			&item.ListingName, &sku,
			&variantDataJSON, &attributesJSON, &item.Quantity, &item.UnitPrice, &item.Subtotal, &item.Discount, &item.Total, &imageURL, &item.CreatedAt,
		)
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to scan order item")
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}

		// Handle nullable fields
		if variantID.Valid {
			item.VariantID = &variantID.Int64
		}
		if stockReservationID.Valid {
			item.StockReservationID = &stockReservationID.String
		}
		if sku.Valid {
			item.SKU = &sku.String
		}
		if imageURL.Valid {
			item.ImageURL = &imageURL.String
		}

		// Parse JSONB data
		if len(variantDataJSON) > 0 {
			if err := json.Unmarshal(variantDataJSON, &item.VariantData); err != nil {
				r.logger.Error().Err(err).Msg("failed to unmarshal variant data")
			}
		}
		if len(attributesJSON) > 0 {
			if err := json.Unmarshal(attributesJSON, &item.Attributes); err != nil {
				r.logger.Error().Err(err).Msg("failed to unmarshal attributes")
			}
		}

		items = append(items, &item)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error().Err(err).Msg("error iterating order item rows")
		return nil, fmt.Errorf("error iterating order item rows: %w", err)
	}

	return items, nil
}

// LockListingsByIDs locks listings by IDs in ascending order to prevent deadlocks
func (r *orderRepository) LockListingsByIDs(ctx context.Context, listingIDs []int64) error {
	if len(listingIDs) == 0 {
		return nil
	}

	// Lock in ascending ID order to prevent deadlocks
	query := `
		SELECT id FROM listings
		WHERE id = ANY($1)
		ORDER BY id ASC
		FOR UPDATE
	`

	rows, err := r.db.Query(ctx, query, listingIDs)
	if err != nil {
		r.logger.Error().Err(err).Interface("listing_ids", listingIDs).Msg("failed to lock listings")
		return fmt.Errorf("failed to lock listings: %w", err)
	}
	defer rows.Close()

	lockedCount := 0
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("failed to scan locked listing ID: %w", err)
		}
		lockedCount++
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("error iterating locked listings: %w", err)
	}

	if lockedCount != len(listingIDs) {
		return fmt.Errorf("failed to lock all listings: requested %d, locked %d", len(listingIDs), lockedCount)
	}

	r.logger.Debug().Interface("listing_ids", listingIDs).Msg("listings locked for transaction")
	return nil
}
