// Package service provides business logic layer for the listings microservice.
package service

import (
	"errors"
	"fmt"
)

// Common errors
var (
	// ErrNotFound indicates that a requested entity was not found
	ErrNotFound = errors.New("entity not found")

	// ErrInvalidInput indicates that input validation failed
	ErrInvalidInput = errors.New("invalid input")

	// ErrUnauthorized indicates that the user is not authorized to perform the action
	ErrUnauthorized = errors.New("unauthorized")

	// ErrConflict indicates a conflict with the current state
	ErrConflict = errors.New("conflict")

	// ErrInternal indicates an internal server error
	ErrInternal = errors.New("internal error")
)

// Cart-specific errors

// ErrCartEmpty indicates that the cart has no items
var ErrCartEmpty = errors.New("cart is empty")

// ErrCartNotFound indicates that the cart was not found
var ErrCartNotFound = errors.New("cart not found")

// ErrCartItemNotFound indicates that the cart item was not found
var ErrCartItemNotFound = errors.New("cart item not found")

// ErrStorefrontMismatch indicates that cart items belong to different storefronts
type ErrStorefrontMismatch struct {
	CartStorefrontID int64
	ItemStorefrontID int64
}

func (e ErrStorefrontMismatch) Error() string {
	return fmt.Sprintf("cart storefront mismatch: cart belongs to storefront %d, but item belongs to storefront %d",
		e.CartStorefrontID, e.ItemStorefrontID)
}

// ErrListingNotFound indicates that a listing/product was not found
type ErrListingNotFound struct {
	ListingID int64
}

func (e ErrListingNotFound) Error() string {
	return fmt.Sprintf("listing %d not found", e.ListingID)
}

// ErrListingInactive indicates that a listing/product is not active
type ErrListingInactive struct {
	ListingID int64
}

func (e ErrListingInactive) Error() string {
	return fmt.Sprintf("listing %d is not active", e.ListingID)
}

// PriceChangeItem represents a single item with price change
type PriceChangeItem struct {
	ListingID     int64
	ListingName   string
	OldPrice      float64
	NewPrice      float64
	PriceIncrease bool
}

// ErrPriceChanged indicates that one or more items have changed price
type ErrPriceChanged struct {
	Changes []PriceChangeItem
}

func (e ErrPriceChanged) Error() string {
	if len(e.Changes) == 0 {
		return "price changed for cart items"
	}
	return fmt.Sprintf("price changed for %d item(s)", len(e.Changes))
}

// Order-specific errors

// ErrOrderNotFound indicates that the order was not found
var ErrOrderNotFound = errors.New("order not found")

// ErrOrderAlreadyConfirmed indicates that the order is already confirmed
var ErrOrderAlreadyConfirmed = errors.New("order already confirmed")

// ErrOrderAlreadyCancelled indicates that the order is already cancelled
var ErrOrderAlreadyCancelled = errors.New("order already cancelled")

// ErrOrderCannotCancel indicates that the order cannot be cancelled in its current state
type ErrOrderCannotCancel struct {
	OrderID int64
	Status  string
}

func (e ErrOrderCannotCancel) Error() string {
	return fmt.Sprintf("order %d cannot be cancelled (current status: %s)", e.OrderID, e.Status)
}

// ErrOrderCannotUpdateStatus indicates that the order status cannot be updated
type ErrOrderCannotUpdateStatus struct {
	OrderID    int64
	FromStatus string
	ToStatus   string
}

func (e ErrOrderCannotUpdateStatus) Error() string {
	return fmt.Sprintf("order %d cannot transition from '%s' to '%s'", e.OrderID, e.FromStatus, e.ToStatus)
}

// ErrOrderInvalidStatus indicates that the order is in invalid status for the requested action
type ErrOrderInvalidStatus struct {
	OrderID        int64
	CurrentStatus  string
	ExpectedStatus string
	Action         string
}

func (e ErrOrderInvalidStatus) Error() string {
	return fmt.Sprintf("cannot %s order %d: current status '%s', expected '%s'",
		e.Action, e.OrderID, e.CurrentStatus, e.ExpectedStatus)
}

// ErrOrderMissingTrackingNumber indicates that the order doesn't have a tracking number
type ErrOrderMissingTrackingNumber struct {
	OrderID int64
}

func (e ErrOrderMissingTrackingNumber) Error() string {
	return fmt.Sprintf("order %d does not have a tracking number", e.OrderID)
}

// ErrInsufficientStock indicates that there is not enough stock for an item
type ErrInsufficientStock struct {
	ListingID      int64
	ListingName    string
	RequestedQty   int32
	AvailableStock int32
}

func (e ErrInsufficientStock) Error() string {
	if e.ListingName != "" {
		return fmt.Sprintf("insufficient stock for '%s': requested %d, available %d",
			e.ListingName, e.RequestedQty, e.AvailableStock)
	}
	return fmt.Sprintf("insufficient stock for listing %d: requested %d, available %d",
		e.ListingID, e.RequestedQty, e.AvailableStock)
}

// ErrInvalidAddress indicates that the shipping or billing address is invalid
var ErrInvalidAddress = errors.New("invalid address")

// ErrInvalidPaymentMethod indicates that the payment method is invalid
var ErrInvalidPaymentMethod = errors.New("invalid payment method")

// ErrPaymentFailed indicates that payment processing failed
type ErrPaymentFailed struct {
	Reason string
}

func (e ErrPaymentFailed) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("payment failed: %s", e.Reason)
	}
	return "payment failed"
}

// Inventory/Reservation-specific errors

// ErrReservationNotFound indicates that the reservation was not found
var ErrReservationNotFound = errors.New("reservation not found")

// ErrReservationExpired indicates that the reservation has expired
type ErrReservationExpired struct {
	ReservationID int64
}

func (e ErrReservationExpired) Error() string {
	return fmt.Sprintf("reservation %d has expired", e.ReservationID)
}

// ErrReservationCannotCommit indicates that the reservation cannot be committed
type ErrReservationCannotCommit struct {
	ReservationID int64
	Reason        string
}

func (e ErrReservationCannotCommit) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("reservation %d cannot be committed: %s", e.ReservationID, e.Reason)
	}
	return fmt.Sprintf("reservation %d cannot be committed", e.ReservationID)
}

// ErrReservationCannotRelease indicates that the reservation cannot be released
type ErrReservationCannotRelease struct {
	ReservationID int64
	Reason        string
}

func (e ErrReservationCannotRelease) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("reservation %d cannot be released: %s", e.ReservationID, e.Reason)
	}
	return fmt.Sprintf("reservation %d cannot be released", e.ReservationID)
}

// ErrStockNotAvailable indicates that stock is not available (locked by reservations)
type ErrStockNotAvailable struct {
	ListingID      int64
	RequestedQty   int32
	TotalStock     int32
	ReservedStock  int32
	AvailableStock int32
}

func (e ErrStockNotAvailable) Error() string {
	return fmt.Sprintf("stock not available for listing %d: requested %d, total stock %d, reserved %d, available %d",
		e.ListingID, e.RequestedQty, e.TotalStock, e.ReservedStock, e.AvailableStock)
}

// Chat-specific errors

// ErrChatNotFound indicates that the chat was not found
var ErrChatNotFound = errors.New("chat not found")

// ErrChatBlocked indicates that the chat is blocked (spam/abuse)
var ErrChatBlocked = errors.New("chat is blocked")

// ErrChatAlreadyExists indicates that a chat already exists for this context
type ErrChatAlreadyExists struct {
	ChatID int64
}

func (e ErrChatAlreadyExists) Error() string {
	return fmt.Sprintf("chat %d already exists for this context", e.ChatID)
}

// ErrChatWithSelf indicates that user is trying to chat with themselves
var ErrChatWithSelf = errors.New("cannot chat with yourself")

// ErrMessageNotFound indicates that the message was not found
var ErrMessageNotFound = errors.New("message not found")

// ErrMessageTooLong indicates that the message content exceeds maximum length
type ErrMessageTooLong struct {
	Length    int
	MaxLength int
}

func (e ErrMessageTooLong) Error() string {
	return fmt.Sprintf("message too long: %d characters (max: %d)", e.Length, e.MaxLength)
}

// ErrMessageEmpty indicates that the message content is empty
var ErrMessageEmpty = errors.New("message content cannot be empty")

// ErrAttachmentNotFound indicates that the attachment was not found
var ErrAttachmentNotFound = errors.New("attachment not found")

// ErrAttachmentTooLarge indicates that the attachment exceeds size limit
type ErrAttachmentTooLarge struct {
	FileType AttachmentFileType
	Size     int64
	MaxSize  int64
}

type AttachmentFileType string

const (
	AttachmentFileTypeImage    AttachmentFileType = "image"
	AttachmentFileTypeVideo    AttachmentFileType = "video"
	AttachmentFileTypeDocument AttachmentFileType = "document"
)

func (e ErrAttachmentTooLarge) Error() string {
	return fmt.Sprintf("%s file too large: %d bytes (max: %d bytes)", e.FileType, e.Size, e.MaxSize)
}

// ErrInvalidFileType indicates that the file type is not supported
type ErrInvalidFileType struct {
	ContentType string
}

func (e ErrInvalidFileType) Error() string {
	return fmt.Sprintf("invalid file type: %s", e.ContentType)
}

// ErrNotParticipant indicates that user is not a participant in the chat
var ErrNotParticipant = errors.New("user is not a participant in this chat")

// ErrNotReceiver indicates that user is not the receiver of the message
var ErrNotReceiver = errors.New("user is not the receiver of this message")

// Helper functions

// IsNotFoundError checks if the error is a "not found" error
func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	// Check for wrapped errors
	if errors.Is(err, ErrNotFound) ||
		errors.Is(err, ErrCartNotFound) ||
		errors.Is(err, ErrCartItemNotFound) ||
		errors.Is(err, ErrOrderNotFound) ||
		errors.Is(err, ErrReservationNotFound) ||
		errors.Is(err, ErrChatNotFound) ||
		errors.Is(err, ErrMessageNotFound) ||
		errors.Is(err, ErrAttachmentNotFound) {
		return true
	}

	// Check for typed errors
	var listingNotFound *ErrListingNotFound
	if errors.As(err, &listingNotFound) {
		return true
	}

	return false
}

// IsConflictError checks if the error is a conflict error
func IsConflictError(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, ErrConflict) ||
		errors.Is(err, ErrOrderAlreadyConfirmed) ||
		errors.Is(err, ErrOrderAlreadyCancelled) {
		return true
	}

	// Check for typed errors
	var priceChanged *ErrPriceChanged
	var storefrontMismatch *ErrStorefrontMismatch
	var insufficientStock *ErrInsufficientStock
	var orderCannotCancel *ErrOrderCannotCancel
	var orderCannotUpdateStatus *ErrOrderCannotUpdateStatus
	var orderInvalidStatus *ErrOrderInvalidStatus
	var orderMissingTrackingNumber *ErrOrderMissingTrackingNumber

	return errors.As(err, &priceChanged) ||
		errors.As(err, &storefrontMismatch) ||
		errors.As(err, &insufficientStock) ||
		errors.As(err, &orderCannotCancel) ||
		errors.As(err, &orderCannotUpdateStatus) ||
		errors.As(err, &orderInvalidStatus) ||
		errors.As(err, &orderMissingTrackingNumber)
}

// IsValidationError checks if the error is a validation error
func IsValidationError(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, ErrInvalidInput) ||
		errors.Is(err, ErrCartEmpty) ||
		errors.Is(err, ErrInvalidAddress) ||
		errors.Is(err, ErrInvalidPaymentMethod)
}

// Shipment-specific errors

// ErrDeliveryClientNotConfigured indicates that delivery client is not set
var ErrDeliveryClientNotConfigured = errors.New("delivery client not configured")

// ErrInvalidDeliveryProvider indicates that the delivery provider code is invalid
type ErrInvalidDeliveryProvider struct {
	ProviderCode string
}

func (e ErrInvalidDeliveryProvider) Error() string {
	return fmt.Sprintf("invalid delivery provider: %s", e.ProviderCode)
}

// ErrShipmentCreationFailed indicates that shipment creation failed
type ErrShipmentCreationFailed struct {
	OrderID int64
	Reason  string
}

func (e ErrShipmentCreationFailed) Error() string {
	return fmt.Sprintf("failed to create shipment for order %d: %s", e.OrderID, e.Reason)
}

// ErrTrackingInfoNotAvailable indicates that tracking info is not available
type ErrTrackingInfoNotAvailable struct {
	OrderID int64
	Reason  string
}

func (e ErrTrackingInfoNotAvailable) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("tracking info not available for order %d: %s", e.OrderID, e.Reason)
	}
	return fmt.Sprintf("tracking info not available for order %d", e.OrderID)
}
