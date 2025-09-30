package common

// ContextKey тип для ключей контекста
type ContextKey string

const (
	// ContextKeyIsAdmin ключ для флага администратора в контексте
	ContextKeyIsAdmin ContextKey = "is_admin"
	// ContextKeyUserID ключ для ID пользователя в контексте
	ContextKeyUserID ContextKey = "user_id"
	// ContextKeyHardDelete ключ для флага жесткого удаления в контексте
	ContextKeyHardDelete ContextKey = "hard_delete"
)
