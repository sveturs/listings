package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pb "github.com/sveturs/listings/api/proto/listings/v1"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// =============================================================================
// Cart Validation Tests
// =============================================================================

func TestCart_Validate_Success(t *testing.T) {
	tests := []struct {
		name string
		cart *Cart
	}{
		{
			name: "valid cart with user_id",
			cart: &Cart{
				ID:           1,
				UserID:       ptr[int64](123),
				StorefrontID: 1,
			},
		},
		{
			name: "valid cart with session_id",
			cart: &Cart{
				ID:           1,
				SessionID:    ptr("session-123"),
				StorefrontID: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cart.Validate()
			assert.NoError(t, err)
		})
	}
}

func TestCart_Validate_Failures(t *testing.T) {
	tests := []struct {
		name    string
		cart    *Cart
		wantErr string
	}{
		{
			name:    "nil cart",
			cart:    nil,
			wantErr: "cart cannot be nil",
		},
		{
			name: "no user_id or session_id",
			cart: &Cart{
				ID:           1,
				StorefrontID: 1,
			},
			wantErr: "cart must have either user_id or session_id",
		},
		{
			name: "both user_id and session_id",
			cart: &Cart{
				ID:           1,
				UserID:       ptr[int64](123),
				SessionID:    ptr("session-123"),
				StorefrontID: 1,
			},
			wantErr: "cart must have either user_id or session_id",
		},
		{
			name: "invalid storefront_id",
			cart: &Cart{
				ID:           1,
				UserID:       ptr[int64](123),
				StorefrontID: 0,
			},
			wantErr: "storefront_id must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cart.Validate()
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

// =============================================================================
// CartItem Validation Tests
// =============================================================================

func TestCartItem_Validate_Success(t *testing.T) {
	item := &CartItem{
		ID:            1,
		CartID:        1,
		ListingID:     100,
		Quantity:      2,
		PriceSnapshot: 99.99,
	}

	err := item.Validate()
	assert.NoError(t, err)
}

func TestCartItem_Validate_Failures(t *testing.T) {
	tests := []struct {
		name    string
		item    *CartItem
		wantErr string
	}{
		{
			name:    "nil item",
			item:    nil,
			wantErr: "cart item cannot be nil",
		},
		{
			name: "invalid cart_id",
			item: &CartItem{
				ID:            1,
				CartID:        0,
				ListingID:     100,
				Quantity:      2,
				PriceSnapshot: 99.99,
			},
			wantErr: "cart_id must be greater than 0",
		},
		{
			name: "invalid listing_id",
			item: &CartItem{
				ID:            1,
				CartID:        1,
				ListingID:     0,
				Quantity:      2,
				PriceSnapshot: 99.99,
			},
			wantErr: "listing_id must be greater than 0",
		},
		{
			name: "invalid quantity",
			item: &CartItem{
				ID:            1,
				CartID:        1,
				ListingID:     100,
				Quantity:      0,
				PriceSnapshot: 99.99,
			},
			wantErr: "quantity must be greater than 0",
		},
		{
			name: "negative price_snapshot",
			item: &CartItem{
				ID:            1,
				CartID:        1,
				ListingID:     100,
				Quantity:      2,
				PriceSnapshot: -10.00,
			},
			wantErr: "price_snapshot cannot be negative",
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
// Cart Operations Tests
// =============================================================================

func TestCart_AddItem_NewItem(t *testing.T) {
	cart := &Cart{
		ID:           1,
		UserID:       ptr[int64](123),
		StorefrontID: 1,
		Items:        []*CartItem{},
	}

	item := &CartItem{
		CartID:        1,
		ListingID:     100,
		Quantity:      2,
		PriceSnapshot: 99.99,
	}

	err := cart.AddItem(item)
	require.NoError(t, err)
	assert.Len(t, cart.Items, 1)
	assert.Equal(t, int32(2), cart.Items[0].Quantity)
	assert.Equal(t, int64(100), cart.Items[0].ListingID)
}

func TestCart_AddItem_UpdateExisting(t *testing.T) {
	oldTime := time.Now().Add(-1 * time.Hour)
	existingItem := &CartItem{
		ID:            1,
		CartID:        1,
		ListingID:     100,
		Quantity:      2,
		PriceSnapshot: 99.99,
		UpdatedAt:     oldTime,
	}

	cart := &Cart{
		ID:           1,
		UserID:       ptr[int64](123),
		StorefrontID: 1,
		Items:        []*CartItem{existingItem},
	}

	// Small sleep to ensure time difference
	time.Sleep(10 * time.Millisecond)

	newItem := &CartItem{
		CartID:        1,
		ListingID:     100,
		Quantity:      3,
		PriceSnapshot: 99.99,
	}

	err := cart.AddItem(newItem)
	require.NoError(t, err)
	assert.Len(t, cart.Items, 1)
	assert.Equal(t, int32(5), cart.Items[0].Quantity) // 2 + 3
	assert.True(t, cart.Items[0].UpdatedAt.After(oldTime))
}

func TestCart_AddItem_DifferentVariant(t *testing.T) {
	variant1 := int64(10)
	variant2 := int64(20)

	existingItem := &CartItem{
		ID:            1,
		CartID:        1,
		ListingID:     100,
		VariantID:     &variant1,
		Quantity:      2,
		PriceSnapshot: 99.99,
	}

	cart := &Cart{
		ID:           1,
		UserID:       ptr[int64](123),
		StorefrontID: 1,
		Items:        []*CartItem{existingItem},
	}

	newItem := &CartItem{
		CartID:        1,
		ListingID:     100,
		VariantID:     &variant2,
		Quantity:      3,
		PriceSnapshot: 99.99,
	}

	err := cart.AddItem(newItem)
	require.NoError(t, err)
	assert.Len(t, cart.Items, 2) // Two separate items
}

func TestCart_AddItem_NilItem(t *testing.T) {
	cart := &Cart{
		ID:           1,
		UserID:       ptr[int64](123),
		StorefrontID: 1,
	}

	err := cart.AddItem(nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "item cannot be nil")
}

func TestCart_RemoveItem_Success(t *testing.T) {
	cart := &Cart{
		ID:           1,
		UserID:       ptr[int64](123),
		StorefrontID: 1,
		Items: []*CartItem{
			{ID: 1, ListingID: 100},
			{ID: 2, ListingID: 200},
		},
	}

	err := cart.RemoveItem(1)
	require.NoError(t, err)
	assert.Len(t, cart.Items, 1)
	assert.Equal(t, int64(2), cart.Items[0].ID)
}

func TestCart_RemoveItem_EmptyCart(t *testing.T) {
	cart := &Cart{
		ID:           1,
		UserID:       ptr[int64](123),
		StorefrontID: 1,
		Items:        []*CartItem{},
	}

	err := cart.RemoveItem(1)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cart is empty")
}

func TestCart_RemoveItem_NotFound(t *testing.T) {
	cart := &Cart{
		ID:           1,
		UserID:       ptr[int64](123),
		StorefrontID: 1,
		Items: []*CartItem{
			{ID: 1, ListingID: 100},
		},
	}

	err := cart.RemoveItem(999)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found in cart")
}

func TestCart_UpdateQuantity_Success(t *testing.T) {
	cart := &Cart{
		ID:           1,
		UserID:       ptr[int64](123),
		StorefrontID: 1,
		Items: []*CartItem{
			{ID: 1, ListingID: 100, Quantity: 2, UpdatedAt: time.Now().Add(-1 * time.Hour)},
		},
	}

	oldUpdateTime := cart.Items[0].UpdatedAt

	err := cart.UpdateQuantity(1, 5)
	require.NoError(t, err)
	assert.Equal(t, int32(5), cart.Items[0].Quantity)
	assert.True(t, cart.Items[0].UpdatedAt.After(oldUpdateTime))
}

func TestCart_UpdateQuantity_InvalidQuantity(t *testing.T) {
	cart := &Cart{
		ID:           1,
		UserID:       ptr[int64](123),
		StorefrontID: 1,
		Items: []*CartItem{
			{ID: 1, ListingID: 100, Quantity: 2},
		},
	}

	err := cart.UpdateQuantity(1, 0)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "quantity must be greater than 0")
}

func TestCart_UpdateQuantity_NotFound(t *testing.T) {
	cart := &Cart{
		ID:           1,
		UserID:       ptr[int64](123),
		StorefrontID: 1,
		Items: []*CartItem{
			{ID: 1, ListingID: 100, Quantity: 2},
		},
	}

	err := cart.UpdateQuantity(999, 5)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found in cart")
}

func TestCart_Clear(t *testing.T) {
	cart := &Cart{
		ID:           1,
		UserID:       ptr[int64](123),
		StorefrontID: 1,
		Items: []*CartItem{
			{ID: 1, ListingID: 100},
			{ID: 2, ListingID: 200},
		},
	}

	cart.Clear()
	assert.Len(t, cart.Items, 0)
}

// =============================================================================
// Cart Calculation Tests
// =============================================================================

func TestCart_CalculateTotal_EmptyCart(t *testing.T) {
	cart := &Cart{
		ID:           1,
		UserID:       ptr[int64](123),
		StorefrontID: 1,
		Items:        []*CartItem{},
	}

	total := cart.CalculateTotal()
	assert.Equal(t, 0.0, total)
}

func TestCart_CalculateTotal_WithItems(t *testing.T) {
	cart := &Cart{
		ID:           1,
		UserID:       ptr[int64](123),
		StorefrontID: 1,
		Items: []*CartItem{
			{Quantity: 2, PriceSnapshot: 10.00},
			{Quantity: 3, PriceSnapshot: 15.00},
		},
	}

	total := cart.CalculateTotal()
	assert.Equal(t, 65.0, total) // 2*10 + 3*15 = 65
}

func TestCart_CalculateTotalItems_EmptyCart(t *testing.T) {
	cart := &Cart{
		ID:           1,
		UserID:       ptr[int64](123),
		StorefrontID: 1,
		Items:        []*CartItem{},
	}

	total := cart.CalculateTotalItems()
	assert.Equal(t, int32(0), total)
}

func TestCart_CalculateTotalItems_WithItems(t *testing.T) {
	cart := &Cart{
		ID:           1,
		UserID:       ptr[int64](123),
		StorefrontID: 1,
		Items: []*CartItem{
			{Quantity: 2},
			{Quantity: 3},
			{Quantity: 5},
		},
	}

	total := cart.CalculateTotalItems()
	assert.Equal(t, int32(10), total) // 2+3+5 = 10
}

// =============================================================================
// Price Change Detection Tests
// =============================================================================

func TestCart_HasPriceChanges_NoChanges(t *testing.T) {
	cart := &Cart{
		ID:           1,
		UserID:       ptr[int64](123),
		StorefrontID: 1,
		Items: []*CartItem{
			{PriceSnapshot: 10.00, CurrentPrice: ptr(10.00)},
			{PriceSnapshot: 15.00, CurrentPrice: ptr(15.00)},
		},
	}

	assert.False(t, cart.HasPriceChanges())
}

func TestCart_HasPriceChanges_WithChanges(t *testing.T) {
	cart := &Cart{
		ID:           1,
		UserID:       ptr[int64](123),
		StorefrontID: 1,
		Items: []*CartItem{
			{PriceSnapshot: 10.00, CurrentPrice: ptr(10.00)},
			{PriceSnapshot: 15.00, CurrentPrice: ptr(12.00)}, // Price changed
		},
	}

	assert.True(t, cart.HasPriceChanges())
}

func TestCart_HasPriceChanges_NilCurrentPrice(t *testing.T) {
	cart := &Cart{
		ID:           1,
		UserID:       ptr[int64](123),
		StorefrontID: 1,
		Items: []*CartItem{
			{PriceSnapshot: 10.00, CurrentPrice: nil}, // No current price
		},
	}

	assert.False(t, cart.HasPriceChanges())
}

func TestCart_GetPriceChanges(t *testing.T) {
	listingName := "Test Product"
	cart := &Cart{
		ID:           1,
		UserID:       ptr[int64](123),
		StorefrontID: 1,
		Items: []*CartItem{
			{PriceSnapshot: 10.00, CurrentPrice: ptr(10.00), ListingName: &listingName},
			{PriceSnapshot: 15.00, CurrentPrice: ptr(12.00), ListingName: &listingName}, // Price decreased
			{PriceSnapshot: 20.00, CurrentPrice: ptr(25.00), ListingName: &listingName}, // Price increased
		},
	}

	warnings := cart.GetPriceChanges()
	require.Len(t, warnings, 2)
	assert.Contains(t, warnings[0], "15.00")
	assert.Contains(t, warnings[0], "12.00")
	assert.Contains(t, warnings[1], "20.00")
	assert.Contains(t, warnings[1], "25.00")
}

// =============================================================================
// Proto Conversion Tests
// =============================================================================

func TestCartFromProto_Success(t *testing.T) {
	now := time.Now()
	userID := int64(123)
	sessionID := "session-123"

	pbCart := &pb.Cart{
		Id:           1,
		UserId:       &userID,
		SessionId:    &sessionID,
		StorefrontId: 10,
		CreatedAt:    timestamppb.New(now),
		UpdatedAt:    timestamppb.New(now),
		Items: []*pb.CartItem{
			{
				Id:            1,
				CartId:        1,
				ListingId:     100,
				Quantity:      2,
				PriceSnapshot: 99.99,
			},
		},
	}

	cart := CartFromProto(pbCart)
	require.NotNil(t, cart)
	assert.Equal(t, int64(1), cart.ID)
	assert.Equal(t, int64(123), *cart.UserID)
	assert.Equal(t, "session-123", *cart.SessionID)
	assert.Equal(t, int64(10), cart.StorefrontID)
	assert.Len(t, cart.Items, 1)
}

func TestCartFromProto_Nil(t *testing.T) {
	cart := CartFromProto(nil)
	assert.Nil(t, cart)
}

func TestCart_ToProto_Success(t *testing.T) {
	now := time.Now()
	userID := int64(123)

	cart := &Cart{
		ID:           1,
		UserID:       &userID,
		StorefrontID: 10,
		CreatedAt:    now,
		UpdatedAt:    now,
		Items: []*CartItem{
			{
				ID:            1,
				CartID:        1,
				ListingID:     100,
				Quantity:      2,
				PriceSnapshot: 99.99,
				CreatedAt:     now,
				UpdatedAt:     now,
			},
		},
	}

	pbCart := cart.ToProto()
	require.NotNil(t, pbCart)
	assert.Equal(t, int64(1), pbCart.Id)
	assert.Equal(t, int64(123), *pbCart.UserId)
	assert.Equal(t, int64(10), pbCart.StorefrontId)
	assert.Len(t, pbCart.Items, 1)
	assert.Equal(t, int64(100), pbCart.Items[0].ListingId)
}

func TestCart_ToProto_Nil(t *testing.T) {
	var cart *Cart
	pbCart := cart.ToProto()
	assert.Nil(t, pbCart)
}

func TestCartItemFromProto_Success(t *testing.T) {
	now := time.Now()
	variantID := int64(10)
	listingName := "Test Product"
	listingImage := "http://example.com/image.jpg"
	availableStock := int32(100)
	currentPrice := 89.99

	variantData := map[string]interface{}{
		"color": "red",
		"size":  "L",
	}
	variantStruct, _ := structpb.NewStruct(variantData)

	pbItem := &pb.CartItem{
		Id:             1,
		CartId:         1,
		ListingId:      100,
		VariantId:      &variantID,
		Quantity:       2,
		PriceSnapshot:  99.99,
		CreatedAt:      timestamppb.New(now),
		UpdatedAt:      timestamppb.New(now),
		ListingName:    &listingName,
		ListingImage:   &listingImage,
		VariantData:    variantStruct,
		AvailableStock: &availableStock,
		CurrentPrice:   &currentPrice,
	}

	item := CartItemFromProto(pbItem)
	require.NotNil(t, item)
	assert.Equal(t, int64(1), item.ID)
	assert.Equal(t, int64(10), *item.VariantID)
	assert.Equal(t, "Test Product", *item.ListingName)
	assert.Equal(t, int32(100), *item.AvailableStock)
	assert.Equal(t, "red", item.VariantData["color"])
}

func TestCartItemFromProto_Nil(t *testing.T) {
	item := CartItemFromProto(nil)
	assert.Nil(t, item)
}

func TestCartItem_ToProto_Success(t *testing.T) {
	now := time.Now()
	variantID := int64(10)
	listingName := "Test Product"
	listingImage := "http://example.com/image.jpg"
	availableStock := int32(100)
	currentPrice := 89.99

	item := &CartItem{
		ID:             1,
		CartID:         1,
		ListingID:      100,
		VariantID:      &variantID,
		Quantity:       2,
		PriceSnapshot:  99.99,
		CreatedAt:      now,
		UpdatedAt:      now,
		ListingName:    &listingName,
		ListingImage:   &listingImage,
		VariantData:    map[string]interface{}{"color": "red"},
		AvailableStock: &availableStock,
		CurrentPrice:   &currentPrice,
	}

	pbItem := item.ToProto()
	require.NotNil(t, pbItem)
	assert.Equal(t, int64(1), pbItem.Id)
	assert.Equal(t, int64(10), *pbItem.VariantId)
	assert.Equal(t, "Test Product", *pbItem.ListingName)
	assert.Equal(t, int32(100), *pbItem.AvailableStock)
}

func TestCartItem_ToProto_Nil(t *testing.T) {
	var item *CartItem
	pbItem := item.ToProto()
	assert.Nil(t, pbItem)
}

// =============================================================================
// Helper Functions
// =============================================================================

func ptr[T any](v T) *T {
	return &v
}
