// Package service provides business logic layer for the listings microservice.
package service

import (
	"fmt"
	"math"

	"github.com/vondi-global/listings/internal/domain"
)

// FinancialConfig contains configuration for financial calculations
type FinancialConfig struct {
	// Tax rate (e.g., 0.20 for 20% VAT)
	TaxRate float64

	// Commission rate that platform takes (e.g., 0.10 for 10%)
	CommissionRate float64

	// Default currency (ISO 4217 code)
	DefaultCurrency string

	// Escrow hold period in days (funds held before releasing to seller)
	EscrowDays int32
}

// DefaultFinancialConfig returns default financial configuration
func DefaultFinancialConfig() *FinancialConfig {
	return &FinancialConfig{
		TaxRate:         0.20,  // 20% VAT
		CommissionRate:  0.10,  // 10% platform fee
		DefaultCurrency: "RSD", // Serbian Dinar
		EscrowDays:      3,     // 3 days escrow hold
	}
}

// OrderFinancials contains all calculated financial values for an order
type OrderFinancials struct {
	Subtotal     float64 // Sum of all items (before tax/shipping)
	Tax          float64 // Tax amount
	ShippingCost float64 // Shipping cost
	Discount     float64 // Discount amount (coupons, promotions)
	Total        float64 // Final amount to pay
	Commission   float64 // Platform commission
	SellerAmount float64 // Amount seller receives (total - commission)
	Currency     string  // ISO 4217 currency code
}

// CalculateOrderFinancials calculates all financial values for an order
func CalculateOrderFinancials(
	items []*domain.OrderItem,
	shippingCost float64,
	discountAmount float64,
	config *FinancialConfig,
) (*OrderFinancials, error) {
	if config == nil {
		config = DefaultFinancialConfig()
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("cannot calculate financials for empty order")
	}

	// Calculate subtotal from items
	var subtotal float64
	for _, item := range items {
		if item == nil {
			continue
		}
		subtotal += item.Total
	}

	// Calculate tax (applied to subtotal only, not shipping)
	tax := roundCurrency(subtotal * config.TaxRate)

	// Ensure discount doesn't exceed subtotal
	if discountAmount > subtotal {
		discountAmount = subtotal
	}

	// Calculate total (subtotal + tax + shipping - discount)
	total := roundCurrency(subtotal + tax + shippingCost - discountAmount)

	// Ensure total is not negative
	if total < 0 {
		total = 0
	}

	// Calculate commission (based on subtotal, not total)
	commission := roundCurrency(subtotal * config.CommissionRate)

	// Calculate seller amount (total - commission)
	sellerAmount := roundCurrency(total - commission)

	// Ensure seller amount is not negative
	if sellerAmount < 0 {
		sellerAmount = 0
	}

	return &OrderFinancials{
		Subtotal:     roundCurrency(subtotal),
		Tax:          tax,
		ShippingCost: roundCurrency(shippingCost),
		Discount:     roundCurrency(discountAmount),
		Total:        total,
		Commission:   commission,
		SellerAmount: sellerAmount,
		Currency:     config.DefaultCurrency,
	}, nil
}

// CalculateItemFinancials calculates financial values for a single order item
func CalculateItemFinancials(quantity int32, unitPrice float64, discountPercent float64) (subtotal, discount, total float64) {
	if quantity <= 0 || unitPrice < 0 {
		return 0, 0, 0
	}

	// Calculate subtotal
	subtotal = float64(quantity) * unitPrice

	// Calculate discount
	if discountPercent > 0 && discountPercent <= 100 {
		discount = subtotal * (discountPercent / 100.0)
	}

	// Calculate total
	total = subtotal - discount

	// Round to 2 decimal places
	subtotal = roundCurrency(subtotal)
	discount = roundCurrency(discount)
	total = roundCurrency(total)

	return subtotal, discount, total
}

// ValidateDiscount validates a discount amount/percentage
func ValidateDiscount(discountAmount float64, subtotal float64) error {
	if discountAmount < 0 {
		return fmt.Errorf("discount cannot be negative")
	}

	if discountAmount > subtotal {
		return fmt.Errorf("discount (%0.2f) cannot exceed subtotal (%0.2f)", discountAmount, subtotal)
	}

	return nil
}

// ValidateShippingCost validates a shipping cost
func ValidateShippingCost(shippingCost float64) error {
	if shippingCost < 0 {
		return fmt.Errorf("shipping cost cannot be negative")
	}

	// Reasonable upper limit (adjust based on business rules)
	const maxShippingCost = 100000.0 // 100,000 RSD
	if shippingCost > maxShippingCost {
		return fmt.Errorf("shipping cost (%0.2f) exceeds maximum allowed (%0.2f)", shippingCost, maxShippingCost)
	}

	return nil
}

// CalculateRefundAmount calculates the refund amount for a cancelled order
func CalculateRefundAmount(order *domain.Order) float64 {
	if order == nil {
		return 0
	}

	// Full refund of the total amount paid
	return roundCurrency(order.Total)
}

// CalculatePartialRefund calculates a partial refund for specific items
func CalculatePartialRefund(items []*domain.OrderItem, config *FinancialConfig) (*OrderFinancials, error) {
	if config == nil {
		config = DefaultFinancialConfig()
	}

	// For partial refunds, we refund item subtotal + proportional tax
	// But no commission or shipping refund
	var subtotal float64
	for _, item := range items {
		if item != nil {
			subtotal += item.Total
		}
	}

	tax := roundCurrency(subtotal * config.TaxRate)
	refundAmount := roundCurrency(subtotal + tax)

	return &OrderFinancials{
		Subtotal:     roundCurrency(subtotal),
		Tax:          tax,
		ShippingCost: 0, // No shipping refund for partial returns
		Discount:     0,
		Total:        refundAmount,
		Commission:   0, // Commission not refunded
		SellerAmount: 0,
		Currency:     config.DefaultCurrency,
	}, nil
}

// BuildOrderItems converts cart items to order items with snapshot data
func BuildOrderItems(cartItems []*domain.CartItem, listings map[int64]*domain.Product) ([]*domain.OrderItem, error) {
	if len(cartItems) == 0 {
		return nil, fmt.Errorf("cart items cannot be empty")
	}

	orderItems := make([]*domain.OrderItem, 0, len(cartItems))

	for _, cartItem := range cartItems {
		if cartItem == nil {
			continue
		}

		// Find corresponding listing
		listing, exists := listings[cartItem.ListingID]
		if !exists {
			return nil, &ErrListingNotFound{ListingID: cartItem.ListingID}
		}

		// Calculate item financials (no item-level discount for now)
		subtotal, discount, total := CalculateItemFinancials(cartItem.Quantity, listing.Price, 0)

		// Create order item with snapshot data
		orderItem := &domain.OrderItem{
			ListingID:   cartItem.ListingID,
			VariantID:   cartItem.VariantID,
			ListingName: listing.Name,
			Quantity:    cartItem.Quantity,
			UnitPrice:   listing.Price,
			Subtotal:    subtotal,
			Discount:    discount,
			Total:       total,
		}

		// Snapshot SKU if available
		if listing.SKU != nil {
			orderItem.SKU = listing.SKU
		}

		// Snapshot variant data if available
		if cartItem.VariantData != nil {
			orderItem.VariantData = cartItem.VariantData
		}

		// Snapshot product attributes
		if listing.Attributes != nil {
			orderItem.Attributes = listing.Attributes
		}

		// Snapshot image URL (get primary image or first available)
		if len(listing.Images) > 0 {
			// Find primary image first
			for _, img := range listing.Images {
				if img.IsPrimary {
					orderItem.ImageURL = &img.URL
					break
				}
			}
			// If no primary image found, use first image
			if orderItem.ImageURL == nil {
				orderItem.ImageURL = &listing.Images[0].URL
			}
		}

		orderItems = append(orderItems, orderItem)
	}

	return orderItems, nil
}

// CalculateCartTotal calculates the total price for cart items
func CalculateCartTotal(items []*domain.CartItem) float64 {
	if len(items) == 0 {
		return 0.0
	}

	var total float64
	for _, item := range items {
		if item != nil {
			total += item.PriceSnapshot * float64(item.Quantity)
		}
	}

	return roundCurrency(total)
}

// roundCurrency rounds a float64 to 2 decimal places (standard for currency)
func roundCurrency(value float64) float64 {
	return math.Round(value*100) / 100
}

// FormatCurrency formats a currency value with 2 decimal places
func FormatCurrency(value float64, currency string) string {
	return fmt.Sprintf("%.2f %s", value, currency)
}

// ValidateCurrency validates a currency code (basic validation)
func ValidateCurrency(currency string) error {
	// List of supported currencies (extend as needed)
	supportedCurrencies := map[string]bool{
		"RSD": true, // Serbian Dinar
		"EUR": true, // Euro
		"USD": true, // US Dollar
	}

	if !supportedCurrencies[currency] {
		return fmt.Errorf("unsupported currency: %s", currency)
	}

	return nil
}
