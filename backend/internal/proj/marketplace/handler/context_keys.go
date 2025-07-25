package handler

import "backend/internal/common"

// ContextKey is a type for context keys to avoid collisions
type ContextKey = common.ContextKey

const (
	// Context keys
	ContextKeyLocale    = common.ContextKeyLocale
	ContextKeyUserID    = common.ContextKeyUserID
	ContextKeyIPAddress = common.ContextKeyIPAddress
)
