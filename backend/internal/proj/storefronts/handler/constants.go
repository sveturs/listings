package handler

type contextKey string

const (
	// Boolean values
	boolValueTrue = "true"

	// Context keys
	isAdminKey    contextKey = "is_admin"
	userIDKey     contextKey = "user_id"
	hardDeleteKey contextKey = "hard_delete"
)
