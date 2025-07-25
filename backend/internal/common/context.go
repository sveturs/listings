// Package common contains shared constants and types used across the application
package common

// ContextKey is a type for context keys to avoid collisions
type ContextKey string

const (
	// ContextKeyLocale is the key for storing locale in context
	ContextKeyLocale ContextKey = "locale"
	// ContextKeyUserID is the key for storing user ID in context
	ContextKeyUserID ContextKey = "user_id"
	// ContextKeyIPAddress is the key for storing IP address in context
	ContextKeyIPAddress ContextKey = "ip_address"
)
