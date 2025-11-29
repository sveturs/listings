// Package domain defines core business entities and domain models for the listings microservice.
package domain

import (
	"errors"
	"fmt"
	"time"

	pb "github.com/vondi-global/listings/api/proto/listings/v1"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// OrderStatus represents the current state of an order
type OrderStatus string

const (
	OrderStatusUnspecified OrderStatus = "unspecified"
	OrderStatusPending     OrderStatus = "pending"    // Order created, awaiting payment
	OrderStatusConfirmed   OrderStatus = "confirmed"  // Payment successful, ready for processing
	OrderStatusAccepted    OrderStatus = "accepted"   // Seller accepted the order
	OrderStatusProcessing  OrderStatus = "processing" // Order being prepared (shipment created)
	OrderStatusShipped     OrderStatus = "shipped"    // Order shipped to customer
	OrderStatusDelivered   OrderStatus = "delivered"  // Order delivered successfully
	OrderStatusCancelled   OrderStatus = "cancelled"  // Order cancelled (by user or admin)
	OrderStatusRefunded    OrderStatus = "refunded"   // Payment refunded
	OrderStatusFailed      OrderStatus = "failed"     // Order processing failed
)

// PaymentStatus represents the current state of payment
type PaymentStatus string

const (
	PaymentStatusUnspecified PaymentStatus = "unspecified"
	PaymentStatusPending     PaymentStatus = "pending"     // Payment initiated, awaiting confirmation
	PaymentStatusProcessing  PaymentStatus = "processing"  // Payment being processed
	PaymentStatusCompleted   PaymentStatus = "completed"   // Payment successful
	PaymentStatusCODPending  PaymentStatus = "cod_pending" // Cash on Delivery - payment will be collected at delivery
	PaymentStatusFailed      PaymentStatus = "failed"      // Payment failed
	PaymentStatusRefunded    PaymentStatus = "refunded"    // Payment refunded to customer
)

// Address represents a flexible address structure stored as JSONB
type Address struct {
	Name       string  `json:"name,omitempty"`
	Street     string  `json:"street,omitempty"`
	City       string  `json:"city,omitempty"`
	PostalCode string  `json:"postal_code,omitempty"`
	Country    string  `json:"country,omitempty"`
	Phone      string  `json:"phone,omitempty"`
	Email      string  `json:"email,omitempty"`
	Additional *string `json:"additional,omitempty"` // Additional notes/instructions
}

// OrderFinancials contains all money-related fields
type OrderFinancials struct {
	Subtotal     float64 `json:"subtotal"`      // Sum of all items (before tax/shipping)
	Tax          float64 `json:"tax"`           // Tax amount
	ShippingCost float64 `json:"shipping_cost"` // Shipping cost
	Discount     float64 `json:"discount"`      // Discount amount (coupons, promotions)
	Total        float64 `json:"total"`         // Final amount to pay
	Commission   float64 `json:"commission"`    // Platform commission (for accounting)
	SellerAmount float64 `json:"seller_amount"` // Amount seller receives (total - commission)
	Currency     string  `json:"currency"`      // ISO 4217 code (USD, EUR, RSD)
}

// Order represents a customer order
type Order struct {
	// Identification
	ID             int64   `json:"id" db:"id"`
	OrderNumber    string  `json:"order_number" db:"order_number"`   // Unique order number (e.g., ORD-2025-001234)
	UserID         *int64  `json:"user_id,omitempty" db:"user_id"`   // NULL for guest orders
	StorefrontID   int64   `json:"storefront_id" db:"storefront_id"` // Storefront that fulfills the order
	StorefrontName *string `json:"storefront_name,omitempty" db:"-"` // Storefront name (joined, not in DB)

	// Order status
	Status OrderStatus `json:"status" db:"status"` // Order lifecycle status

	// Financial data
	Subtotal     float64 `json:"subtotal" db:"subtotal"`
	Tax          float64 `json:"tax" db:"tax"`
	Shipping     float64 `json:"shipping" db:"shipping"`
	Discount     float64 `json:"discount" db:"discount"`
	Total        float64 `json:"total" db:"total"`
	Commission   float64 `json:"commission" db:"commission"`
	SellerAmount float64 `json:"seller_amount" db:"seller_amount"`
	Currency     string  `json:"currency" db:"currency"`

	// Payment information
	PaymentMethod        *string       `json:"payment_method,omitempty" db:"payment_method"`                 // cash, card, bank_transfer, paypal, etc.
	PaymentStatus        PaymentStatus `json:"payment_status" db:"payment_status"`                           // Payment processing status
	PaymentTransactionID *string       `json:"payment_transaction_id,omitempty" db:"payment_transaction_id"` // External payment ID
	PaymentCompletedAt   *time.Time    `json:"payment_completed_at,omitempty" db:"payment_completed_at"`

	// Shipping information
	ShippingAddress  map[string]interface{} `json:"shipping_address,omitempty" db:"shipping_address"`   // JSONB
	BillingAddress   map[string]interface{} `json:"billing_address,omitempty" db:"billing_address"`     // JSONB
	ShippingMethod   *string                `json:"shipping_method,omitempty" db:"shipping_method"`     // standard, express, overnight
	ShippingProvider *string                `json:"shipping_provider,omitempty" db:"shipping_provider"` // Post Express, AKS, DHL, etc.
	TrackingNumber   *string                `json:"tracking_number,omitempty" db:"tracking_number"`     // Shipment tracking number
	ShipmentID       *int64                 `json:"shipment_id,omitempty" db:"shipment_id"`             // FK to Delivery Service shipment

	// Escrow (platform holds funds)
	EscrowReleaseDate *time.Time `json:"escrow_release_date,omitempty" db:"escrow_release_date"` // When funds released to seller
	EscrowDays        int32      `json:"escrow_days" db:"escrow_days"`                           // Days to hold in escrow (default: 3)

	// Customer contact
	CustomerName  *string `json:"customer_name,omitempty" db:"customer_name"`
	CustomerEmail *string `json:"customer_email,omitempty" db:"customer_email"`
	CustomerPhone *string `json:"customer_phone,omitempty" db:"customer_phone"`

	// Notes
	CustomerNotes *string `json:"customer_notes,omitempty" db:"notes"`      // Customer instructions
	AdminNotes    *string `json:"admin_notes,omitempty" db:"admin_notes"`   // Internal admin notes
	SellerNotes   *string `json:"seller_notes,omitempty" db:"seller_notes"` // Seller notes about the order

	// Timestamps
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	ConfirmedAt *time.Time `json:"confirmed_at,omitempty" db:"confirmed_at"`
	AcceptedAt  *time.Time `json:"accepted_at,omitempty" db:"accepted_at"` // When seller accepted
	ShippedAt   *time.Time `json:"shipped_at,omitempty" db:"shipped_at"`
	DeliveredAt *time.Time `json:"delivered_at,omitempty" db:"delivered_at"`
	CancelledAt *time.Time `json:"cancelled_at,omitempty" db:"cancelled_at"`

	// Shipping label
	LabelURL *string `json:"label_url,omitempty" db:"label_url"` // URL to shipping label PDF

	// Relations (loaded on demand)
	Items []*OrderItem `json:"items,omitempty" db:"-"`
}

// OrderItem represents a single line item in an order
type OrderItem struct {
	ID        int64  `json:"id" db:"id"`
	OrderID   int64  `json:"order_id" db:"order_id"`
	ListingID int64  `json:"listing_id" db:"listing_id"`           // FK to listings or products
	VariantID *int64 `json:"variant_id,omitempty" db:"variant_id"` // FK to listing_variants or product_variants

	// Snapshot data (immutable after order creation)
	ListingName string                 `json:"listing_name" db:"listing_name"`           // Product name at purchase time
	SKU         *string                `json:"sku,omitempty" db:"sku"`                   // SKU at purchase time
	VariantData map[string]interface{} `json:"variant_data,omitempty" db:"variant_data"` // Variant attributes snapshot
	Attributes  map[string]interface{} `json:"attributes,omitempty" db:"attributes"`     // Product attributes snapshot

	// Quantity and pricing
	Quantity  int32   `json:"quantity" db:"quantity"` // Quantity ordered
	UnitPrice float64 `json:"unit_price" db:"price"`  // Price per item
	Subtotal  float64 `json:"subtotal" db:"subtotal"` // quantity * unit_price
	Discount  float64 `json:"discount" db:"discount"` // Item-level discount
	Total     float64 `json:"total" db:"total"`       // subtotal - discount

	// Product snapshot
	ImageURL *string `json:"image_url,omitempty" db:"image_url"` // Primary image at purchase time

	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Validate validates the Order entity
func (o *Order) Validate() error {
	if o == nil {
		return errors.New("order cannot be nil")
	}

	if o.OrderNumber == "" {
		return errors.New("order_number is required")
	}

	if o.StorefrontID <= 0 {
		return errors.New("storefront_id must be greater than 0")
	}

	if o.Subtotal < 0 {
		return errors.New("subtotal cannot be negative")
	}

	if o.Total < 0 {
		return errors.New("total cannot be negative")
	}

	if o.SellerAmount < 0 {
		return errors.New("seller_amount cannot be negative")
	}

	if o.EscrowDays < 0 {
		return errors.New("escrow_days cannot be negative")
	}

	return nil
}

// ValidateAddresses validates shipping and billing addresses
func (o *Order) ValidateAddresses() error {
	if o.ShippingAddress == nil || len(o.ShippingAddress) == 0 {
		return errors.New("shipping_address is required")
	}

	// Billing address defaults to shipping if not provided
	// So no strict validation needed

	return nil
}

// CalculateFinancials calculates and populates financial fields
func (o *Order) CalculateFinancials(items []*OrderItem, taxRate float64, shippingCost float64, discount float64, commissionRate float64) error {
	if items == nil || len(items) == 0 {
		return errors.New("order must have at least one item")
	}

	// Calculate subtotal from items
	var subtotal float64
	for _, item := range items {
		subtotal += item.Total
	}

	o.Subtotal = subtotal
	o.Tax = subtotal * taxRate
	o.Shipping = shippingCost
	o.Discount = discount
	o.Total = subtotal + o.Tax + o.Shipping - o.Discount

	// Calculate commission and seller amount
	o.Commission = subtotal * commissionRate
	o.SellerAmount = o.Total - o.Commission

	return nil
}

// CanCancel checks if order can be cancelled
func (o *Order) CanCancel() bool {
	// Can only cancel pending, confirmed or accepted orders
	return o.Status == OrderStatusPending || o.Status == OrderStatusConfirmed || o.Status == OrderStatusAccepted
}

// CanAccept checks if order can be accepted by seller
func (o *Order) CanAccept() bool {
	return o.Status == OrderStatusConfirmed
}

// CanCreateShipment checks if shipment can be created for this order
func (o *Order) CanCreateShipment() bool {
	return o.Status == OrderStatusAccepted
}

// CanMarkShipped checks if order can be marked as shipped
func (o *Order) CanMarkShipped() bool {
	return o.Status == OrderStatusProcessing && o.TrackingNumber != nil
}

// CanUpdateStatus checks if order status can be updated to newStatus
func (o *Order) CanUpdateStatus(newStatus OrderStatus) bool {
	// Define valid state transitions
	// Flow: pending → confirmed → accepted → processing → shipped → delivered
	validTransitions := map[OrderStatus][]OrderStatus{
		OrderStatusPending: {
			OrderStatusConfirmed,
			OrderStatusCancelled,
			OrderStatusFailed,
		},
		OrderStatusConfirmed: {
			OrderStatusAccepted, // Seller accepts
			OrderStatusCancelled,
		},
		OrderStatusAccepted: {
			OrderStatusProcessing, // Shipment created
			OrderStatusCancelled,
		},
		OrderStatusProcessing: {
			OrderStatusShipped,
			OrderStatusCancelled,
		},
		OrderStatusShipped: {
			OrderStatusDelivered,
		},
		OrderStatusDelivered: {
			OrderStatusRefunded,
		},
	}

	allowedStatuses, exists := validTransitions[o.Status]
	if !exists {
		return false
	}

	for _, allowed := range allowedStatuses {
		if allowed == newStatus {
			return true
		}
	}

	return false
}

// GenerateOrderNumber generates a unique order number
// Format: ORD-YYYY-NNNNNN (e.g., ORD-2025-001234)
func GenerateOrderNumber(year int, sequence int64) string {
	return fmt.Sprintf("ORD-%d-%06d", year, sequence)
}

// ValidateOrderItem validates the OrderItem entity
func (i *OrderItem) Validate() error {
	if i == nil {
		return errors.New("order item cannot be nil")
	}

	if i.OrderID <= 0 {
		return errors.New("order_id must be greater than 0")
	}

	if i.ListingID <= 0 {
		return errors.New("listing_id must be greater than 0")
	}

	if i.ListingName == "" {
		return errors.New("listing_name is required")
	}

	if i.Quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	if i.UnitPrice < 0 {
		return errors.New("unit_price cannot be negative")
	}

	if i.Total < 0 {
		return errors.New("total cannot be negative")
	}

	return nil
}

// OrderStatusFromProto converts proto OrderStatus to domain OrderStatus
func OrderStatusFromProto(pbStatus pb.OrderStatus) OrderStatus {
	switch pbStatus {
	case pb.OrderStatus_ORDER_STATUS_PENDING:
		return OrderStatusPending
	case pb.OrderStatus_ORDER_STATUS_CONFIRMED:
		return OrderStatusConfirmed
	case pb.OrderStatus_ORDER_STATUS_ACCEPTED:
		return OrderStatusAccepted
	case pb.OrderStatus_ORDER_STATUS_PROCESSING:
		return OrderStatusProcessing
	case pb.OrderStatus_ORDER_STATUS_SHIPPED:
		return OrderStatusShipped
	case pb.OrderStatus_ORDER_STATUS_DELIVERED:
		return OrderStatusDelivered
	case pb.OrderStatus_ORDER_STATUS_CANCELLED:
		return OrderStatusCancelled
	case pb.OrderStatus_ORDER_STATUS_REFUNDED:
		return OrderStatusRefunded
	case pb.OrderStatus_ORDER_STATUS_FAILED:
		return OrderStatusFailed
	default:
		return OrderStatusUnspecified
	}
}

// ToProtoOrderStatus converts domain OrderStatus to proto OrderStatus
func (s OrderStatus) ToProtoOrderStatus() pb.OrderStatus {
	switch s {
	case OrderStatusPending:
		return pb.OrderStatus_ORDER_STATUS_PENDING
	case OrderStatusConfirmed:
		return pb.OrderStatus_ORDER_STATUS_CONFIRMED
	case OrderStatusAccepted:
		return pb.OrderStatus_ORDER_STATUS_ACCEPTED
	case OrderStatusProcessing:
		return pb.OrderStatus_ORDER_STATUS_PROCESSING
	case OrderStatusShipped:
		return pb.OrderStatus_ORDER_STATUS_SHIPPED
	case OrderStatusDelivered:
		return pb.OrderStatus_ORDER_STATUS_DELIVERED
	case OrderStatusCancelled:
		return pb.OrderStatus_ORDER_STATUS_CANCELLED
	case OrderStatusRefunded:
		return pb.OrderStatus_ORDER_STATUS_REFUNDED
	case OrderStatusFailed:
		return pb.OrderStatus_ORDER_STATUS_FAILED
	default:
		return pb.OrderStatus_ORDER_STATUS_UNSPECIFIED
	}
}

// PaymentStatusFromProto converts proto PaymentStatus to domain PaymentStatus
func PaymentStatusFromProto(pbStatus pb.PaymentStatus) PaymentStatus {
	switch pbStatus {
	case pb.PaymentStatus_PAYMENT_STATUS_PENDING:
		return PaymentStatusPending
	case pb.PaymentStatus_PAYMENT_STATUS_PROCESSING:
		return PaymentStatusProcessing
	case pb.PaymentStatus_PAYMENT_STATUS_COMPLETED:
		return PaymentStatusCompleted
	case pb.PaymentStatus_PAYMENT_STATUS_COD_PENDING:
		return PaymentStatusCODPending
	case pb.PaymentStatus_PAYMENT_STATUS_FAILED:
		return PaymentStatusFailed
	case pb.PaymentStatus_PAYMENT_STATUS_REFUNDED:
		return PaymentStatusRefunded
	default:
		return PaymentStatusUnspecified
	}
}

// ToProtoPaymentStatus converts domain PaymentStatus to proto PaymentStatus
func (s PaymentStatus) ToProtoPaymentStatus() pb.PaymentStatus {
	switch s {
	case PaymentStatusPending:
		return pb.PaymentStatus_PAYMENT_STATUS_PENDING
	case PaymentStatusProcessing:
		return pb.PaymentStatus_PAYMENT_STATUS_PROCESSING
	case PaymentStatusCompleted:
		return pb.PaymentStatus_PAYMENT_STATUS_COMPLETED
	case PaymentStatusCODPending:
		return pb.PaymentStatus_PAYMENT_STATUS_COD_PENDING
	case PaymentStatusFailed:
		return pb.PaymentStatus_PAYMENT_STATUS_FAILED
	case PaymentStatusRefunded:
		return pb.PaymentStatus_PAYMENT_STATUS_REFUNDED
	default:
		return pb.PaymentStatus_PAYMENT_STATUS_UNSPECIFIED
	}
}

// OrderFromProto converts proto Order to domain Order
func OrderFromProto(pb *pb.Order) *Order {
	if pb == nil {
		return nil
	}

	order := &Order{
		ID:            pb.Id,
		OrderNumber:   pb.OrderNumber,
		StorefrontID:  pb.StorefrontId,
		Status:        OrderStatusFromProto(pb.Status),
		PaymentStatus: PaymentStatusFromProto(pb.PaymentStatus),
		EscrowDays:    pb.EscrowDays,
		Currency:      pb.Financials.Currency,
	}

	if pb.UserId != nil {
		userID := *pb.UserId
		order.UserID = &userID
	}

	// Financial data
	if pb.Financials != nil {
		order.Subtotal = pb.Financials.Subtotal
		order.Tax = pb.Financials.Tax
		order.Shipping = pb.Financials.ShippingCost
		order.Discount = pb.Financials.Discount
		order.Total = pb.Financials.Total
		order.Commission = pb.Financials.Commission
		order.SellerAmount = pb.Financials.SellerAmount
	}

	// Payment info
	if pb.PaymentMethod != nil {
		order.PaymentMethod = pb.PaymentMethod
	}
	if pb.PaymentTransactionId != nil {
		order.PaymentTransactionID = pb.PaymentTransactionId
	}
	if pb.PaymentCompletedAt != nil {
		t := pb.PaymentCompletedAt.AsTime()
		order.PaymentCompletedAt = &t
	}

	// Addresses
	if pb.ShippingAddress != nil {
		order.ShippingAddress = pb.ShippingAddress.AsMap()
	}
	if pb.BillingAddress != nil {
		order.BillingAddress = pb.BillingAddress.AsMap()
	}

	// Shipping info
	if pb.ShippingMethod != nil {
		order.ShippingMethod = pb.ShippingMethod
	}
	if pb.ShippingProvider != nil {
		order.ShippingProvider = pb.ShippingProvider
	}
	if pb.TrackingNumber != nil {
		order.TrackingNumber = pb.TrackingNumber
	}
	if pb.ShipmentId != nil {
		shipmentID := *pb.ShipmentId
		order.ShipmentID = &shipmentID
	}

	// Escrow
	if pb.EscrowReleaseDate != nil {
		t := pb.EscrowReleaseDate.AsTime()
		order.EscrowReleaseDate = &t
	}

	// Customer contact
	if pb.CustomerName != nil {
		order.CustomerName = pb.CustomerName
	}
	if pb.CustomerEmail != nil {
		order.CustomerEmail = pb.CustomerEmail
	}
	if pb.CustomerPhone != nil {
		order.CustomerPhone = pb.CustomerPhone
	}

	// Notes
	if pb.CustomerNotes != nil {
		order.CustomerNotes = pb.CustomerNotes
	}
	if pb.AdminNotes != nil {
		order.AdminNotes = pb.AdminNotes
	}
	if pb.SellerNotes != nil {
		order.SellerNotes = pb.SellerNotes
	}

	// Timestamps
	if pb.CreatedAt != nil {
		order.CreatedAt = pb.CreatedAt.AsTime()
	}
	if pb.UpdatedAt != nil {
		order.UpdatedAt = pb.UpdatedAt.AsTime()
	}
	if pb.ConfirmedAt != nil {
		t := pb.ConfirmedAt.AsTime()
		order.ConfirmedAt = &t
	}
	if pb.AcceptedAt != nil {
		t := pb.AcceptedAt.AsTime()
		order.AcceptedAt = &t
	}
	if pb.ShippedAt != nil {
		t := pb.ShippedAt.AsTime()
		order.ShippedAt = &t
	}
	if pb.DeliveredAt != nil {
		t := pb.DeliveredAt.AsTime()
		order.DeliveredAt = &t
	}
	if pb.CancelledAt != nil {
		t := pb.CancelledAt.AsTime()
		order.CancelledAt = &t
	}

	// Shipping label
	if pb.LabelUrl != nil {
		order.LabelURL = pb.LabelUrl
	}

	// Convert items
	if pb.Items != nil && len(pb.Items) > 0 {
		order.Items = make([]*OrderItem, 0, len(pb.Items))
		for _, pbItem := range pb.Items {
			order.Items = append(order.Items, OrderItemFromProto(pbItem))
		}
	}

	return order
}

// ToProto converts domain Order to proto Order
func (o *Order) ToProto() *pb.Order {
	if o == nil {
		return nil
	}

	pbOrder := &pb.Order{
		Id:            o.ID,
		OrderNumber:   o.OrderNumber,
		StorefrontId:  o.StorefrontID,
		Status:        o.Status.ToProtoOrderStatus(),
		PaymentStatus: o.PaymentStatus.ToProtoPaymentStatus(),
		EscrowDays:    o.EscrowDays,
		CreatedAt:     timestamppb.New(o.CreatedAt),
		UpdatedAt:     timestamppb.New(o.UpdatedAt),
	}

	if o.UserID != nil {
		pbOrder.UserId = o.UserID
	}

	// Financial data
	pbOrder.Financials = &pb.OrderFinancials{
		Subtotal:     o.Subtotal,
		Tax:          o.Tax,
		ShippingCost: o.Shipping,
		Discount:     o.Discount,
		Total:        o.Total,
		Commission:   o.Commission,
		SellerAmount: o.SellerAmount,
		Currency:     o.Currency,
	}

	// Payment info
	if o.PaymentMethod != nil {
		pbOrder.PaymentMethod = o.PaymentMethod
	}
	if o.PaymentTransactionID != nil {
		pbOrder.PaymentTransactionId = o.PaymentTransactionID
	}
	if o.PaymentCompletedAt != nil {
		pbOrder.PaymentCompletedAt = timestamppb.New(*o.PaymentCompletedAt)
	}

	// Addresses
	if o.ShippingAddress != nil {
		shippingStruct, err := structpb.NewStruct(o.ShippingAddress)
		if err == nil {
			pbOrder.ShippingAddress = shippingStruct
		}
	}
	if o.BillingAddress != nil {
		billingStruct, err := structpb.NewStruct(o.BillingAddress)
		if err == nil {
			pbOrder.BillingAddress = billingStruct
		}
	}

	// Shipping info
	if o.ShippingMethod != nil {
		pbOrder.ShippingMethod = o.ShippingMethod
	}
	if o.ShippingProvider != nil {
		pbOrder.ShippingProvider = o.ShippingProvider
	}
	if o.TrackingNumber != nil {
		pbOrder.TrackingNumber = o.TrackingNumber
	}
	if o.ShipmentID != nil {
		pbOrder.ShipmentId = o.ShipmentID
	}

	// Escrow
	if o.EscrowReleaseDate != nil {
		pbOrder.EscrowReleaseDate = timestamppb.New(*o.EscrowReleaseDate)
	}

	// Customer contact
	if o.CustomerName != nil {
		pbOrder.CustomerName = o.CustomerName
	}
	if o.CustomerEmail != nil {
		pbOrder.CustomerEmail = o.CustomerEmail
	}
	if o.CustomerPhone != nil {
		pbOrder.CustomerPhone = o.CustomerPhone
	}

	// Notes
	if o.CustomerNotes != nil {
		pbOrder.CustomerNotes = o.CustomerNotes
	}
	if o.AdminNotes != nil {
		pbOrder.AdminNotes = o.AdminNotes
	}
	if o.SellerNotes != nil {
		pbOrder.SellerNotes = o.SellerNotes
	}

	// Timestamps
	if o.ConfirmedAt != nil {
		pbOrder.ConfirmedAt = timestamppb.New(*o.ConfirmedAt)
	}
	if o.AcceptedAt != nil {
		pbOrder.AcceptedAt = timestamppb.New(*o.AcceptedAt)
	}
	if o.ShippedAt != nil {
		pbOrder.ShippedAt = timestamppb.New(*o.ShippedAt)
	}
	if o.DeliveredAt != nil {
		pbOrder.DeliveredAt = timestamppb.New(*o.DeliveredAt)
	}
	if o.CancelledAt != nil {
		pbOrder.CancelledAt = timestamppb.New(*o.CancelledAt)
	}

	// Shipping label
	if o.LabelURL != nil {
		pbOrder.LabelUrl = o.LabelURL
	}

	// Convert items
	if o.Items != nil && len(o.Items) > 0 {
		pbOrder.Items = make([]*pb.OrderItem, 0, len(o.Items))
		for _, item := range o.Items {
			pbOrder.Items = append(pbOrder.Items, item.ToProto())
		}
	}

	return pbOrder
}

// OrderItemFromProto converts proto OrderItem to domain OrderItem
func OrderItemFromProto(pb *pb.OrderItem) *OrderItem {
	if pb == nil {
		return nil
	}

	item := &OrderItem{
		ID:          pb.Id,
		OrderID:     pb.OrderId,
		ListingID:   pb.ListingId,
		ListingName: pb.ListingName,
		Quantity:    pb.Quantity,
		UnitPrice:   pb.UnitPrice,
		Subtotal:    pb.Subtotal,
		Discount:    pb.Discount,
		Total:       pb.Total,
	}

	if pb.VariantId != nil {
		variantID := *pb.VariantId
		item.VariantID = &variantID
	}

	if pb.Sku != nil {
		item.SKU = pb.Sku
	}

	if pb.VariantData != nil {
		item.VariantData = pb.VariantData.AsMap()
	}

	if pb.Attributes != nil {
		item.Attributes = pb.Attributes.AsMap()
	}

	if pb.ImageUrl != nil {
		item.ImageURL = pb.ImageUrl
	}

	if pb.CreatedAt != nil {
		item.CreatedAt = pb.CreatedAt.AsTime()
	}

	return item
}

// ToProto converts domain OrderItem to proto OrderItem
func (i *OrderItem) ToProto() *pb.OrderItem {
	if i == nil {
		return nil
	}

	pbItem := &pb.OrderItem{
		Id:          i.ID,
		OrderId:     i.OrderID,
		ListingId:   i.ListingID,
		ListingName: i.ListingName,
		Quantity:    i.Quantity,
		UnitPrice:   i.UnitPrice,
		Subtotal:    i.Subtotal,
		Discount:    i.Discount,
		Total:       i.Total,
		CreatedAt:   timestamppb.New(i.CreatedAt),
	}

	if i.VariantID != nil {
		pbItem.VariantId = i.VariantID
	}

	if i.SKU != nil {
		pbItem.Sku = i.SKU
	}

	if i.VariantData != nil {
		variantStruct, err := structpb.NewStruct(i.VariantData)
		if err == nil {
			pbItem.VariantData = variantStruct
		}
	}

	if i.Attributes != nil {
		attrsStruct, err := structpb.NewStruct(i.Attributes)
		if err == nil {
			pbItem.Attributes = attrsStruct
		}
	}

	if i.ImageURL != nil {
		pbItem.ImageUrl = i.ImageURL
	}

	return pbItem
}
