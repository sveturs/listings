package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pb "github.com/vondi-global/listings/api/proto/listings/v1"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// =============================================================================
// Order Validation Tests
// =============================================================================

func TestOrder_Validate_Success(t *testing.T) {
	order := &Order{
		OrderNumber:  "ORD-2025-000001",
		StorefrontID: 1,
		Subtotal:     100.00,
		Tax:          10.00,
		Shipping:     5.00,
		Discount:     0.00,
		Total:        115.00,
		Commission:   5.00,
		SellerAmount: 110.00,
		EscrowDays:   3,
	}

	err := order.Validate()
	assert.NoError(t, err)
}

func TestOrder_Validate_Failures(t *testing.T) {
	tests := []struct {
		name    string
		order   *Order
		wantErr string
	}{
		{
			name:    "nil order",
			order:   nil,
			wantErr: "order cannot be nil",
		},
		{
			name: "empty order_number",
			order: &Order{
				StorefrontID: 1,
				Total:        100.00,
			},
			wantErr: "order_number is required",
		},
		{
			name: "invalid storefront_id",
			order: &Order{
				OrderNumber:  "ORD-2025-000001",
				StorefrontID: 0,
				Total:        100.00,
			},
			wantErr: "storefront_id must be greater than 0",
		},
		{
			name: "negative subtotal",
			order: &Order{
				OrderNumber:  "ORD-2025-000001",
				StorefrontID: 1,
				Subtotal:     -10.00,
				Total:        100.00,
			},
			wantErr: "subtotal cannot be negative",
		},
		{
			name: "negative total",
			order: &Order{
				OrderNumber:  "ORD-2025-000001",
				StorefrontID: 1,
				Subtotal:     100.00,
				Total:        -10.00,
			},
			wantErr: "total cannot be negative",
		},
		{
			name: "negative seller_amount",
			order: &Order{
				OrderNumber:  "ORD-2025-000001",
				StorefrontID: 1,
				Subtotal:     100.00,
				Total:        100.00,
				SellerAmount: -10.00,
			},
			wantErr: "seller_amount cannot be negative",
		},
		{
			name: "negative escrow_days",
			order: &Order{
				OrderNumber:  "ORD-2025-000001",
				StorefrontID: 1,
				Total:        100.00,
				EscrowDays:   -1,
			},
			wantErr: "escrow_days cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.order.Validate()
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestOrder_ValidateAddresses_Success(t *testing.T) {
	order := &Order{
		ShippingAddress: map[string]interface{}{
			"street": "123 Main St",
			"city":   "New York",
		},
	}

	err := order.ValidateAddresses()
	assert.NoError(t, err)
}

func TestOrder_ValidateAddresses_Failures(t *testing.T) {
	tests := []struct {
		name    string
		order   *Order
		wantErr string
	}{
		{
			name: "nil shipping_address",
			order: &Order{
				ShippingAddress: nil,
			},
			wantErr: "shipping_address is required",
		},
		{
			name: "empty shipping_address",
			order: &Order{
				ShippingAddress: map[string]interface{}{},
			},
			wantErr: "shipping_address is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.order.ValidateAddresses()
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

// =============================================================================
// OrderItem Validation Tests
// =============================================================================

func TestOrderItem_Validate_Success(t *testing.T) {
	item := &OrderItem{
		OrderID:     1,
		ListingID:   100,
		ListingName: "Test Product",
		Quantity:    2,
		UnitPrice:   99.99,
		Subtotal:    199.98,
		Discount:    0.00,
		Total:       199.98,
	}

	err := item.Validate()
	assert.NoError(t, err)
}

func TestOrderItem_Validate_Failures(t *testing.T) {
	tests := []struct {
		name    string
		item    *OrderItem
		wantErr string
	}{
		{
			name:    "nil item",
			item:    nil,
			wantErr: "order item cannot be nil",
		},
		{
			name: "invalid order_id",
			item: &OrderItem{
				OrderID:     0,
				ListingID:   100,
				ListingName: "Test",
				Quantity:    1,
				UnitPrice:   10.00,
				Total:       10.00,
			},
			wantErr: "order_id must be greater than 0",
		},
		{
			name: "invalid listing_id",
			item: &OrderItem{
				OrderID:     1,
				ListingID:   0,
				ListingName: "Test",
				Quantity:    1,
				UnitPrice:   10.00,
				Total:       10.00,
			},
			wantErr: "listing_id must be greater than 0",
		},
		{
			name: "empty listing_name",
			item: &OrderItem{
				OrderID:   1,
				ListingID: 100,
				Quantity:  1,
				UnitPrice: 10.00,
				Total:     10.00,
			},
			wantErr: "listing_name is required",
		},
		{
			name: "invalid quantity",
			item: &OrderItem{
				OrderID:     1,
				ListingID:   100,
				ListingName: "Test",
				Quantity:    0,
				UnitPrice:   10.00,
				Total:       10.00,
			},
			wantErr: "quantity must be greater than 0",
		},
		{
			name: "negative unit_price",
			item: &OrderItem{
				OrderID:     1,
				ListingID:   100,
				ListingName: "Test",
				Quantity:    1,
				UnitPrice:   -10.00,
				Total:       10.00,
			},
			wantErr: "unit_price cannot be negative",
		},
		{
			name: "negative total",
			item: &OrderItem{
				OrderID:     1,
				ListingID:   100,
				ListingName: "Test",
				Quantity:    1,
				UnitPrice:   10.00,
				Total:       -10.00,
			},
			wantErr: "total cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.item.Validate()
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

// =============================================================================
// Order Financial Calculation Tests
// =============================================================================

func TestOrder_CalculateFinancials_Success(t *testing.T) {
	order := &Order{}
	items := []*OrderItem{
		{
			UnitPrice: 100.00,
			Quantity:  2,
			Discount:  10.00,
			Total:     190.00, // (100*2) - 10
		},
		{
			UnitPrice: 50.00,
			Quantity:  1,
			Discount:  0.00,
			Total:     50.00,
		},
	}

	err := order.CalculateFinancials(items, 0.1, 15.00, 20.00, 0.05)
	require.NoError(t, err)

	assert.Equal(t, 240.00, order.Subtotal) // 190 + 50
	assert.Equal(t, 24.00, order.Tax)       // 240 * 0.1
	assert.Equal(t, 15.00, order.Shipping)
	assert.Equal(t, 20.00, order.Discount)
	assert.Equal(t, 259.00, order.Total)        // 240 + 24 + 15 - 20
	assert.Equal(t, 12.00, order.Commission)    // 240 * 0.05
	assert.Equal(t, 247.00, order.SellerAmount) // 259 - 12
}

func TestOrder_CalculateFinancials_NoItems(t *testing.T) {
	order := &Order{}

	err := order.CalculateFinancials(nil, 0.1, 15.00, 20.00, 0.05)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "order must have at least one item")
}

// =============================================================================
// Order State Machine Tests
// =============================================================================

func TestOrder_CanCancel_PendingStatus(t *testing.T) {
	order := &Order{Status: OrderStatusPending}
	assert.True(t, order.CanCancel())
}

func TestOrder_CanCancel_ConfirmedStatus(t *testing.T) {
	order := &Order{Status: OrderStatusConfirmed}
	assert.True(t, order.CanCancel())
}

func TestOrder_CanCancel_ShippedStatus(t *testing.T) {
	order := &Order{Status: OrderStatusShipped}
	assert.False(t, order.CanCancel())
}

func TestOrder_CanUpdateStatus_ValidTransitions(t *testing.T) {
	tests := []struct {
		name      string
		current   OrderStatus
		new       OrderStatus
		canUpdate bool
	}{
		// From Pending
		{"pending -> confirmed", OrderStatusPending, OrderStatusConfirmed, true},
		{"pending -> cancelled", OrderStatusPending, OrderStatusCancelled, true},
		{"pending -> failed", OrderStatusPending, OrderStatusFailed, true},
		{"pending -> processing", OrderStatusPending, OrderStatusProcessing, false},

		// From Confirmed
		{"confirmed -> accepted", OrderStatusConfirmed, OrderStatusAccepted, true},
		{"confirmed -> processing", OrderStatusConfirmed, OrderStatusProcessing, false},
		{"confirmed -> cancelled", OrderStatusConfirmed, OrderStatusCancelled, true},
		{"confirmed -> shipped", OrderStatusConfirmed, OrderStatusShipped, false},

		// From Accepted
		{"accepted -> processing", OrderStatusAccepted, OrderStatusProcessing, true},
		{"accepted -> cancelled", OrderStatusAccepted, OrderStatusCancelled, true},

		// From Processing
		{"processing -> shipped", OrderStatusProcessing, OrderStatusShipped, true},
		{"processing -> cancelled", OrderStatusProcessing, OrderStatusCancelled, true},
		{"processing -> delivered", OrderStatusProcessing, OrderStatusDelivered, false},

		// From Shipped
		{"shipped -> delivered", OrderStatusShipped, OrderStatusDelivered, true},
		{"shipped -> cancelled", OrderStatusShipped, OrderStatusCancelled, false},

		// From Delivered
		{"delivered -> refunded", OrderStatusDelivered, OrderStatusRefunded, true},
		{"delivered -> cancelled", OrderStatusDelivered, OrderStatusCancelled, false},

		// Terminal states
		{"cancelled -> pending", OrderStatusCancelled, OrderStatusPending, false},
		{"refunded -> confirmed", OrderStatusRefunded, OrderStatusConfirmed, false},
		{"failed -> pending", OrderStatusFailed, OrderStatusPending, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := &Order{Status: tt.current}
			canUpdate := order.CanUpdateStatus(tt.new)
			assert.Equal(t, tt.canUpdate, canUpdate)
		})
	}
}

// =============================================================================
// Order Number Generation Tests
// =============================================================================

func TestGenerateOrderNumber(t *testing.T) {
	tests := []struct {
		year     int
		sequence int64
		expected string
	}{
		{2025, 1, "ORD-2025-000001"},
		{2025, 123, "ORD-2025-000123"},
		{2025, 999999, "ORD-2025-999999"},
		{2026, 1, "ORD-2026-000001"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := GenerateOrderNumber(tt.year, tt.sequence)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// =============================================================================
// Order Status Proto Conversion Tests
// =============================================================================

func TestOrderStatusFromProto(t *testing.T) {
	tests := []struct {
		pbStatus pb.OrderStatus
		expected OrderStatus
	}{
		{pb.OrderStatus_ORDER_STATUS_PENDING, OrderStatusPending},
		{pb.OrderStatus_ORDER_STATUS_CONFIRMED, OrderStatusConfirmed},
		{pb.OrderStatus_ORDER_STATUS_PROCESSING, OrderStatusProcessing},
		{pb.OrderStatus_ORDER_STATUS_SHIPPED, OrderStatusShipped},
		{pb.OrderStatus_ORDER_STATUS_DELIVERED, OrderStatusDelivered},
		{pb.OrderStatus_ORDER_STATUS_CANCELLED, OrderStatusCancelled},
		{pb.OrderStatus_ORDER_STATUS_REFUNDED, OrderStatusRefunded},
		{pb.OrderStatus_ORDER_STATUS_FAILED, OrderStatusFailed},
		{pb.OrderStatus_ORDER_STATUS_UNSPECIFIED, OrderStatusUnspecified},
	}

	for _, tt := range tests {
		t.Run(string(tt.expected), func(t *testing.T) {
			result := OrderStatusFromProto(tt.pbStatus)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestOrderStatus_ToProtoOrderStatus(t *testing.T) {
	tests := []struct {
		status   OrderStatus
		expected pb.OrderStatus
	}{
		{OrderStatusPending, pb.OrderStatus_ORDER_STATUS_PENDING},
		{OrderStatusConfirmed, pb.OrderStatus_ORDER_STATUS_CONFIRMED},
		{OrderStatusProcessing, pb.OrderStatus_ORDER_STATUS_PROCESSING},
		{OrderStatusShipped, pb.OrderStatus_ORDER_STATUS_SHIPPED},
		{OrderStatusDelivered, pb.OrderStatus_ORDER_STATUS_DELIVERED},
		{OrderStatusCancelled, pb.OrderStatus_ORDER_STATUS_CANCELLED},
		{OrderStatusRefunded, pb.OrderStatus_ORDER_STATUS_REFUNDED},
		{OrderStatusFailed, pb.OrderStatus_ORDER_STATUS_FAILED},
		{OrderStatusUnspecified, pb.OrderStatus_ORDER_STATUS_UNSPECIFIED},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			result := tt.status.ToProtoOrderStatus()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// =============================================================================
// Payment Status Proto Conversion Tests
// =============================================================================

func TestPaymentStatusFromProto(t *testing.T) {
	tests := []struct {
		pbStatus pb.PaymentStatus
		expected PaymentStatus
	}{
		{pb.PaymentStatus_PAYMENT_STATUS_PENDING, PaymentStatusPending},
		{pb.PaymentStatus_PAYMENT_STATUS_PROCESSING, PaymentStatusProcessing},
		{pb.PaymentStatus_PAYMENT_STATUS_COMPLETED, PaymentStatusCompleted},
		{pb.PaymentStatus_PAYMENT_STATUS_FAILED, PaymentStatusFailed},
		{pb.PaymentStatus_PAYMENT_STATUS_REFUNDED, PaymentStatusRefunded},
		{pb.PaymentStatus_PAYMENT_STATUS_UNSPECIFIED, PaymentStatusUnspecified},
	}

	for _, tt := range tests {
		t.Run(string(tt.expected), func(t *testing.T) {
			result := PaymentStatusFromProto(tt.pbStatus)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPaymentStatus_ToProtoPaymentStatus(t *testing.T) {
	tests := []struct {
		status   PaymentStatus
		expected pb.PaymentStatus
	}{
		{PaymentStatusPending, pb.PaymentStatus_PAYMENT_STATUS_PENDING},
		{PaymentStatusProcessing, pb.PaymentStatus_PAYMENT_STATUS_PROCESSING},
		{PaymentStatusCompleted, pb.PaymentStatus_PAYMENT_STATUS_COMPLETED},
		{PaymentStatusFailed, pb.PaymentStatus_PAYMENT_STATUS_FAILED},
		{PaymentStatusRefunded, pb.PaymentStatus_PAYMENT_STATUS_REFUNDED},
		{PaymentStatusUnspecified, pb.PaymentStatus_PAYMENT_STATUS_UNSPECIFIED},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			result := tt.status.ToProtoPaymentStatus()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// =============================================================================
// Order Proto Conversion Tests
// =============================================================================

func TestOrderFromProto_Success(t *testing.T) {
	now := time.Now()
	userID := int64(123)
	paymentMethod := "card"
	transactionID := "txn-123"

	shippingAddress := map[string]interface{}{
		"street": "123 Main St",
		"city":   "New York",
	}
	shippingStruct, _ := structpb.NewStruct(shippingAddress)

	pbOrder := &pb.Order{
		Id:                   1,
		OrderNumber:          "ORD-2025-000001",
		UserId:               &userID,
		StorefrontId:         10,
		Status:               pb.OrderStatus_ORDER_STATUS_CONFIRMED,
		PaymentStatus:        pb.PaymentStatus_PAYMENT_STATUS_COMPLETED,
		PaymentMethod:        &paymentMethod,
		PaymentTransactionId: &transactionID,
		ShippingAddress:      shippingStruct,
		EscrowDays:           3,
		CreatedAt:            timestamppb.New(now),
		UpdatedAt:            timestamppb.New(now),
		Financials: &pb.OrderFinancials{
			Subtotal:     100.00,
			Tax:          10.00,
			ShippingCost: 5.00,
			Discount:     0.00,
			Total:        115.00,
			Commission:   5.00,
			SellerAmount: 110.00,
			Currency:     "USD",
		},
		Items: []*pb.OrderItem{
			{
				Id:          1,
				OrderId:     1,
				ListingId:   100,
				ListingName: "Test Product",
				Quantity:    2,
				UnitPrice:   50.00,
				Subtotal:    100.00,
				Total:       100.00,
			},
		},
	}

	order := OrderFromProto(pbOrder)
	require.NotNil(t, order)
	assert.Equal(t, int64(1), order.ID)
	assert.Equal(t, "ORD-2025-000001", order.OrderNumber)
	assert.Equal(t, int64(123), *order.UserID)
	assert.Equal(t, OrderStatusConfirmed, order.Status)
	assert.Equal(t, PaymentStatusCompleted, order.PaymentStatus)
	assert.Equal(t, 100.00, order.Subtotal)
	assert.Equal(t, "USD", order.Currency)
	assert.Len(t, order.Items, 1)
}

func TestOrderFromProto_Nil(t *testing.T) {
	order := OrderFromProto(nil)
	assert.Nil(t, order)
}

func TestOrder_ToProto_Success(t *testing.T) {
	now := time.Now()
	userID := int64(123)
	paymentMethod := "card"

	order := &Order{
		ID:            1,
		OrderNumber:   "ORD-2025-000001",
		UserID:        &userID,
		StorefrontID:  10,
		Status:        OrderStatusConfirmed,
		PaymentStatus: PaymentStatusCompleted,
		PaymentMethod: &paymentMethod,
		Subtotal:      100.00,
		Tax:           10.00,
		Shipping:      5.00,
		Total:         115.00,
		Commission:    5.00,
		SellerAmount:  110.00,
		Currency:      "USD",
		EscrowDays:    3,
		CreatedAt:     now,
		UpdatedAt:     now,
		ShippingAddress: map[string]interface{}{
			"street": "123 Main St",
		},
		Items: []*OrderItem{
			{
				ID:          1,
				OrderID:     1,
				ListingID:   100,
				ListingName: "Test Product",
				Quantity:    2,
				UnitPrice:   50.00,
				Subtotal:    100.00,
				Total:       100.00,
				CreatedAt:   now,
			},
		},
	}

	pbOrder := order.ToProto()
	require.NotNil(t, pbOrder)
	assert.Equal(t, int64(1), pbOrder.Id)
	assert.Equal(t, "ORD-2025-000001", pbOrder.OrderNumber)
	assert.Equal(t, int64(123), *pbOrder.UserId)
	assert.Equal(t, pb.OrderStatus_ORDER_STATUS_CONFIRMED, pbOrder.Status)
	assert.Equal(t, pb.PaymentStatus_PAYMENT_STATUS_COMPLETED, pbOrder.PaymentStatus)
	assert.NotNil(t, pbOrder.Financials)
	assert.Equal(t, 100.00, pbOrder.Financials.Subtotal)
	assert.Len(t, pbOrder.Items, 1)
}

func TestOrder_ToProto_Nil(t *testing.T) {
	var order *Order
	pbOrder := order.ToProto()
	assert.Nil(t, pbOrder)
}

func TestOrderItemFromProto_Success(t *testing.T) {
	now := time.Now()
	variantID := int64(10)
	sku := "SKU-123"
	imageURL := "http://example.com/image.jpg"

	variantData := map[string]interface{}{"color": "red"}
	variantStruct, _ := structpb.NewStruct(variantData)

	attributes := map[string]interface{}{"brand": "Nike"}
	attrsStruct, _ := structpb.NewStruct(attributes)

	pbItem := &pb.OrderItem{
		Id:          1,
		OrderId:     1,
		ListingId:   100,
		VariantId:   &variantID,
		ListingName: "Test Product",
		Sku:         &sku,
		VariantData: variantStruct,
		Attributes:  attrsStruct,
		Quantity:    2,
		UnitPrice:   50.00,
		Subtotal:    100.00,
		Discount:    0.00,
		Total:       100.00,
		ImageUrl:    &imageURL,
		CreatedAt:   timestamppb.New(now),
	}

	item := OrderItemFromProto(pbItem)
	require.NotNil(t, item)
	assert.Equal(t, int64(1), item.ID)
	assert.Equal(t, int64(10), *item.VariantID)
	assert.Equal(t, "SKU-123", *item.SKU)
	assert.Equal(t, "red", item.VariantData["color"])
	assert.Equal(t, "Nike", item.Attributes["brand"])
}

func TestOrderItemFromProto_Nil(t *testing.T) {
	item := OrderItemFromProto(nil)
	assert.Nil(t, item)
}

func TestOrderItem_ToProto_Success(t *testing.T) {
	now := time.Now()
	variantID := int64(10)
	sku := "SKU-123"
	imageURL := "http://example.com/image.jpg"

	item := &OrderItem{
		ID:          1,
		OrderID:     1,
		ListingID:   100,
		VariantID:   &variantID,
		ListingName: "Test Product",
		SKU:         &sku,
		VariantData: map[string]interface{}{"color": "red"},
		Attributes:  map[string]interface{}{"brand": "Nike"},
		Quantity:    2,
		UnitPrice:   50.00,
		Subtotal:    100.00,
		Discount:    0.00,
		Total:       100.00,
		ImageURL:    &imageURL,
		CreatedAt:   now,
	}

	pbItem := item.ToProto()
	require.NotNil(t, pbItem)
	assert.Equal(t, int64(1), pbItem.Id)
	assert.Equal(t, int64(10), *pbItem.VariantId)
	assert.Equal(t, "SKU-123", *pbItem.Sku)
	assert.NotNil(t, pbItem.VariantData)
	assert.NotNil(t, pbItem.Attributes)
}

func TestOrderItem_ToProto_Nil(t *testing.T) {
	var item *OrderItem
	pbItem := item.ToProto()
	assert.Nil(t, pbItem)
}
