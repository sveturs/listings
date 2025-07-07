package contextkeys

// ContextKey is the type for context keys
type ContextKey string

const (
	// ChatIDKey is the context key for chat ID
	ChatIDKey ContextKey = "chat_id"
	// UserIDKey is the context key for user ID
	UserIDKey ContextKey = "user_id"
	// RequestIDKey is the context key for request ID
	RequestIDKey ContextKey = "request_id"
	// TransactionKey is the context key for database transaction
	TransactionKey ContextKey = "transaction"
	// ListingExistsKey is the context key for listing exists flag
	ListingExistsKey ContextKey = "listing_exists"
)
