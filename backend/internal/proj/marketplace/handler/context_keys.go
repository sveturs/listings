package handler

// ContextKey is a type for context keys to avoid collisions
type ContextKey string

const (
	// Context keys
	ContextKeyLocale    ContextKey = "locale"
	ContextKeyUserID    ContextKey = "user_id"
	ContextKeyIPAddress ContextKey = "ip_address"
)