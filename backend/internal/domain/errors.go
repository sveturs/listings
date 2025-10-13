// Package domain provides domain-level error types for the application.
//
// This package defines typed errors that replace string-based error checking
// throughout the codebase. Use errors.Is() and errors.As() instead of err.Error().
package domain

import "errors"

// Common domain errors that can be used across all modules
var (
	// General errors
	ErrNotFound         = errors.New("not found")
	ErrAlreadyExists    = errors.New("already exists")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
	ErrPermissionDenied = errors.New("permission denied")
	ErrValidationFailed = errors.New("validation failed")
	ErrInvalidInput     = errors.New("invalid input")
	ErrDuplicateKey     = errors.New("duplicate key")

	// Database errors
	ErrNoRows            = errors.New("no rows in result set")
	ErrTransactionFailed = errors.New("transaction failed")

	// Business logic errors
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrNotActive         = errors.New("not active")
	ErrNotAvailable      = errors.New("not available")
	ErrInvalidState      = errors.New("invalid state")

	// Authentication/Authorization errors
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("expired token")
	ErrInvalidSignature   = errors.New("invalid signature")
	ErrInvalidCredentials = errors.New("invalid credentials")

	// Marketplace-specific errors
	ErrListingNotFound     = errors.New("listing not found")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrAlreadyInFavorites  = errors.New("already in favorites")
	ErrNotInFavorites      = errors.New("not in favorites")

	// Subscription errors
	ErrActiveSubscriptionExists = errors.New("user already has active subscription")
	ErrNoActiveSubscription     = errors.New("no active subscription found")

	// Contact errors
	ErrContactAlreadyExists        = errors.New("contact already exists")
	ErrCannotAddYourself           = errors.New("cannot add yourself as contact")
	ErrUserNotAllowContactRequests = errors.New("user does not allow contact requests or has blocked you")

	// Search errors
	ErrSavedSearchNotFound  = errors.New("saved search not found")
	ErrSynonymAlreadyExists = errors.New("synonym already exists")
	ErrSynonymNotFound      = errors.New("synonym not found")

	// Translation errors
	ErrTranslationNotFound = errors.New("translation not found")

	// Order errors
	ErrOrderNotFound      = errors.New("order not found")
	ErrOrderAccessDenied  = errors.New("unauthorized: not a party of this order")
	ErrInvalidOrderStatus = errors.New("invalid order status")
	ErrCannotCancelOrder  = errors.New("order cannot be canceled")

	// Logistics errors
	ErrShipmentNotFound = errors.New("shipment not found")
	ErrProblemNotFound  = errors.New("problem not found")

	// VIN errors
	ErrInvalidVINLength   = errors.New("invalid VIN length")
	ErrVINDecoderDisabled = errors.New("VIN decoder is disabled")

	// OAuth errors
	ErrInvalidOAuthState       = errors.New("invalid state")
	ErrOAuthCodeExchangeFailed = errors.New("failed to exchange code")

	// TLS/Network errors
	ErrTLSHandshake = errors.New("TLS handshake failed")

	// Connection errors
	ErrConnectionBusy = errors.New("conn busy")
)

// IsNotFoundError checks if an error is any kind of "not found" error
func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrNotFound) ||
		errors.Is(err, ErrListingNotFound) ||
		errors.Is(err, ErrOrderNotFound) ||
		errors.Is(err, ErrSavedSearchNotFound) ||
		errors.Is(err, ErrSynonymNotFound) ||
		errors.Is(err, ErrTranslationNotFound) ||
		errors.Is(err, ErrShipmentNotFound) ||
		errors.Is(err, ErrProblemNotFound) ||
		errors.Is(err, ErrNoActiveSubscription) ||
		errors.Is(err, ErrNoRows)
}

// IsPermissionError checks if an error is any kind of permission/authorization error
func IsPermissionError(err error) bool {
	return errors.Is(err, ErrUnauthorized) ||
		errors.Is(err, ErrForbidden) ||
		errors.Is(err, ErrPermissionDenied) ||
		errors.Is(err, ErrOrderAccessDenied)
}

// IsValidationError checks if an error is any kind of validation error
func IsValidationError(err error) bool {
	return errors.Is(err, ErrValidationFailed) ||
		errors.Is(err, ErrInvalidInput) ||
		errors.Is(err, ErrInvalidVINLength) ||
		errors.Is(err, ErrInvalidOAuthState) ||
		errors.Is(err, ErrInvalidToken) ||
		errors.Is(err, ErrInvalidSignature) ||
		errors.Is(err, ErrInvalidCredentials)
}

// IsDuplicateError checks if an error is any kind of duplicate/already exists error
func IsDuplicateError(err error) bool {
	return errors.Is(err, ErrAlreadyExists) ||
		errors.Is(err, ErrDuplicateKey) ||
		errors.Is(err, ErrContactAlreadyExists) ||
		errors.Is(err, ErrAlreadyInFavorites) ||
		errors.Is(err, ErrSynonymAlreadyExists) ||
		errors.Is(err, ErrActiveSubscriptionExists)
}
