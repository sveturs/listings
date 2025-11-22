// Package domain defines core business entities and domain models for the listings microservice.
package domain

import (
	"errors"
	"fmt"
	"time"

	pb "github.com/sveturs/listings/api/proto/listings/v1"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Cart represents a shopping cart entity (authenticated or anonymous)
type Cart struct {
	ID           int64     `json:"id" db:"id"`
	UserID       *int64    `json:"user_id,omitempty" db:"user_id"`       // NULL for anonymous carts
	SessionID    *string   `json:"session_id,omitempty" db:"session_id"` // NULL for authenticated carts
	StorefrontID int64     `json:"storefront_id" db:"storefront_id"`     // Required - cart belongs to one storefront
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`

	// Relations (loaded on demand)
	Items []*CartItem `json:"items,omitempty" db:"-"`
}

// CartItem represents a single item in shopping cart
type CartItem struct {
	ID            int64     `json:"id" db:"id"`
	CartID        int64     `json:"cart_id" db:"cart_id"`
	ListingID     int64     `json:"listing_id" db:"listing_id"`           // FK to listings
	VariantID     *int64    `json:"variant_id,omitempty" db:"variant_id"` // FK to listing_variants
	Quantity      int32     `json:"quantity" db:"quantity"`               // Quantity to purchase
	PriceSnapshot float64   `json:"price_snapshot" db:"price_snapshot"`   // Price at add-to-cart time
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`

	// Embedded data (for frontend display) - loaded separately
	ListingName    *string                `json:"listing_name,omitempty" db:"-"`    // Listing/product name
	ListingImage   *string                `json:"listing_image,omitempty" db:"-"`   // Primary image URL
	VariantData    map[string]interface{} `json:"variant_data,omitempty" db:"-"`    // Variant attributes (color, size, etc.)
	AvailableStock *int32                 `json:"available_stock,omitempty" db:"-"` // Current stock availability
	CurrentPrice   *float64               `json:"current_price,omitempty" db:"-"`   // Current price (may differ from price_snapshot)
}

// Validate validates the Cart entity
func (c *Cart) Validate() error {
	if c == nil {
		return errors.New("cart cannot be nil")
	}

	// Must have either user_id or session_id (not both, not neither)
	if (c.UserID == nil && c.SessionID == nil) || (c.UserID != nil && c.SessionID != nil) {
		return errors.New("cart must have either user_id or session_id (not both, not neither)")
	}

	if c.StorefrontID <= 0 {
		return errors.New("storefront_id must be greater than 0")
	}

	return nil
}

// ValidateCartItem validates the CartItem entity
func (i *CartItem) Validate() error {
	if i == nil {
		return errors.New("cart item cannot be nil")
	}

	if i.CartID <= 0 {
		return errors.New("cart_id must be greater than 0")
	}

	if i.ListingID <= 0 {
		return errors.New("listing_id must be greater than 0")
	}

	if i.Quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	if i.PriceSnapshot < 0 {
		return errors.New("price_snapshot cannot be negative")
	}

	return nil
}

// CalculateTotal calculates the total price for all items in cart
func (c *Cart) CalculateTotal() float64 {
	if c.Items == nil || len(c.Items) == 0 {
		return 0.0
	}

	var total float64
	for _, item := range c.Items {
		total += item.PriceSnapshot * float64(item.Quantity)
	}

	return total
}

// CalculateTotalItems calculates the total number of items in cart
func (c *Cart) CalculateTotalItems() int32 {
	if c.Items == nil || len(c.Items) == 0 {
		return 0
	}

	var totalItems int32
	for _, item := range c.Items {
		totalItems += item.Quantity
	}

	return totalItems
}

// AddItem adds a new item to the cart or updates quantity if already exists
func (c *Cart) AddItem(item *CartItem) error {
	if item == nil {
		return errors.New("item cannot be nil")
	}

	if err := item.Validate(); err != nil {
		return fmt.Errorf("invalid item: %w", err)
	}

	// Check if item already exists (same listing_id and variant_id)
	for _, existingItem := range c.Items {
		if existingItem.ListingID == item.ListingID {
			// Check if variant_id matches (both nil or both equal)
			variantMatch := (existingItem.VariantID == nil && item.VariantID == nil) ||
				(existingItem.VariantID != nil && item.VariantID != nil && *existingItem.VariantID == *item.VariantID)

			if variantMatch {
				// Update quantity
				existingItem.Quantity += item.Quantity
				existingItem.UpdatedAt = time.Now()
				return nil
			}
		}
	}

	// Add new item
	if c.Items == nil {
		c.Items = []*CartItem{}
	}
	item.CartID = c.ID
	c.Items = append(c.Items, item)

	return nil
}

// RemoveItem removes an item from the cart by item ID
func (c *Cart) RemoveItem(itemID int64) error {
	if c.Items == nil || len(c.Items) == 0 {
		return errors.New("cart is empty")
	}

	for i, item := range c.Items {
		if item.ID == itemID {
			// Remove item by slicing
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("item with ID %d not found in cart", itemID)
}

// UpdateQuantity updates the quantity of a cart item
func (c *Cart) UpdateQuantity(itemID int64, newQuantity int32) error {
	if newQuantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	if c.Items == nil || len(c.Items) == 0 {
		return errors.New("cart is empty")
	}

	for _, item := range c.Items {
		if item.ID == itemID {
			item.Quantity = newQuantity
			item.UpdatedAt = time.Now()
			return nil
		}
	}

	return fmt.Errorf("item with ID %d not found in cart", itemID)
}

// Clear removes all items from the cart
func (c *Cart) Clear() {
	c.Items = []*CartItem{}
}

// HasPriceChanges checks if any cart items have price changes
func (c *Cart) HasPriceChanges() bool {
	if c.Items == nil || len(c.Items) == 0 {
		return false
	}

	for _, item := range c.Items {
		if item.CurrentPrice != nil && *item.CurrentPrice != item.PriceSnapshot {
			return true
		}
	}

	return false
}

// GetPriceChanges returns a list of items with price changes
func (c *Cart) GetPriceChanges() []string {
	warnings := []string{}

	if c.Items == nil || len(c.Items) == 0 {
		return warnings
	}

	for _, item := range c.Items {
		if item.CurrentPrice != nil && *item.CurrentPrice != item.PriceSnapshot {
			warning := fmt.Sprintf("Price changed for %s: %.2f → %.2f",
				getStringOrDefault(item.ListingName, "item"),
				item.PriceSnapshot,
				*item.CurrentPrice)
			warnings = append(warnings, warning)
		}
	}

	return warnings
}

// CartFromProto converts proto Cart to domain Cart
func CartFromProto(pb *pb.Cart) *Cart {
	if pb == nil {
		return nil
	}

	cart := &Cart{
		ID:           pb.Id,
		StorefrontID: pb.StorefrontId,
	}

	if pb.UserId != nil {
		userID := *pb.UserId
		cart.UserID = &userID
	}

	if pb.SessionId != nil {
		sessionID := *pb.SessionId
		cart.SessionID = &sessionID
	}

	if pb.CreatedAt != nil {
		cart.CreatedAt = pb.CreatedAt.AsTime()
	}

	if pb.UpdatedAt != nil {
		cart.UpdatedAt = pb.UpdatedAt.AsTime()
	}

	// Convert items
	if pb.Items != nil && len(pb.Items) > 0 {
		cart.Items = make([]*CartItem, 0, len(pb.Items))
		for _, pbItem := range pb.Items {
			cart.Items = append(cart.Items, CartItemFromProto(pbItem))
		}
	}

	return cart
}

// ToProto converts domain Cart to proto Cart
func (c *Cart) ToProto() *pb.Cart {
	if c == nil {
		return nil
	}

	pbCart := &pb.Cart{
		Id:           c.ID,
		StorefrontId: c.StorefrontID,
		CreatedAt:    timestamppb.New(c.CreatedAt),
		UpdatedAt:    timestamppb.New(c.UpdatedAt),
	}

	if c.UserID != nil {
		pbCart.UserId = c.UserID
	}

	if c.SessionID != nil {
		pbCart.SessionId = c.SessionID
	}

	// Convert items
	if c.Items != nil && len(c.Items) > 0 {
		pbCart.Items = make([]*pb.CartItem, 0, len(c.Items))
		for _, item := range c.Items {
			pbCart.Items = append(pbCart.Items, item.ToProto())
		}
	}

	return pbCart
}

// CartItemFromProto converts proto CartItem to domain CartItem
func CartItemFromProto(pb *pb.CartItem) *CartItem {
	if pb == nil {
		return nil
	}

	item := &CartItem{
		ID:            pb.Id,
		CartID:        pb.CartId,
		ListingID:     pb.ListingId,
		Quantity:      pb.Quantity,
		PriceSnapshot: pb.PriceSnapshot,
	}

	if pb.VariantId != nil {
		variantID := *pb.VariantId
		item.VariantID = &variantID
	}

	if pb.CreatedAt != nil {
		item.CreatedAt = pb.CreatedAt.AsTime()
	}

	if pb.UpdatedAt != nil {
		item.UpdatedAt = pb.UpdatedAt.AsTime()
	}

	// Embedded data
	if pb.ListingName != nil {
		listingName := *pb.ListingName
		item.ListingName = &listingName
	}

	if pb.ListingImage != nil {
		listingImage := *pb.ListingImage
		item.ListingImage = &listingImage
	}

	if pb.VariantData != nil {
		item.VariantData = pb.VariantData.AsMap()
	}

	if pb.AvailableStock != nil {
		availableStock := *pb.AvailableStock
		item.AvailableStock = &availableStock
	}

	if pb.CurrentPrice != nil {
		currentPrice := *pb.CurrentPrice
		item.CurrentPrice = &currentPrice
	}

	return item
}

// ToProto converts domain CartItem to proto CartItem
func (i *CartItem) ToProto() *pb.CartItem {
	if i == nil {
		return nil
	}

	pbItem := &pb.CartItem{
		Id:            i.ID,
		CartId:        i.CartID,
		ListingId:     i.ListingID,
		Quantity:      i.Quantity,
		PriceSnapshot: i.PriceSnapshot,
		CreatedAt:     timestamppb.New(i.CreatedAt),
		UpdatedAt:     timestamppb.New(i.UpdatedAt),
	}

	if i.VariantID != nil {
		pbItem.VariantId = i.VariantID
	}

	// Embedded data
	if i.ListingName != nil {
		pbItem.ListingName = i.ListingName
	}

	if i.ListingImage != nil {
		pbItem.ListingImage = i.ListingImage
	}

	if i.VariantData != nil {
		variantStruct, err := structpb.NewStruct(i.VariantData)
		if err == nil {
			pbItem.VariantData = variantStruct
		}
	}

	if i.AvailableStock != nil {
		pbItem.AvailableStock = i.AvailableStock
	}

	if i.CurrentPrice != nil {
		pbItem.CurrentPrice = i.CurrentPrice
	}

	return pbItem
}

// Helper function to get string or default value
func getStringOrDefault(ptr *string, defaultVal string) string {
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}
