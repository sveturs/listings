package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"backend/internal/proj/viber/models"
	"backend/internal/storage/postgres"
)

// ErrNoActiveSession indicates that there is no active session for the user
var ErrNoActiveSession = errors.New("no active session")

// SessionManager управляет сессиями Viber (24-часовые окна бесплатных сообщений)
type SessionManager struct {
	db *postgres.Database
}

// NewSessionManager создаёт новый менеджер сессий
func NewSessionManager(db *postgres.Database) *SessionManager {
	return &SessionManager{db: db}
}

// GetActiveSession возвращает активную сессию для пользователя
func (sm *SessionManager) GetActiveSession(ctx context.Context, viberID string) (*models.ViberSession, error) {
	// Сначала получаем пользователя
	user, err := sm.getOrCreateViberUser(ctx, viberID)
	if err != nil {
		return nil, err
	}

	// Ищем активную сессию
	query := `
		SELECT id, viber_user_id, started_at, last_message_at, expires_at,
		       message_count, context, active, created_at
		FROM viber_sessions
		WHERE viber_user_id = $1
		  AND active = true
		  AND expires_at > CURRENT_TIMESTAMP
		ORDER BY created_at DESC
		LIMIT 1
	`

	var session models.ViberSession
	err = sm.db.GetSQLXDB().GetContext(ctx, &session, query, user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoActiveSession
		}
		return nil, fmt.Errorf("failed to get active session: %w", err)
	}

	return &session, nil
}

// CreateSession создаёт новую сессию при получении сообщения от пользователя
func (sm *SessionManager) CreateSession(ctx context.Context, viberID string) (*models.ViberSession, error) {
	user, err := sm.getOrCreateViberUser(ctx, viberID)
	if err != nil {
		return nil, err
	}

	// Закрываем предыдущие сессии
	_, err = sm.db.ExecContext(ctx, `
		UPDATE viber_sessions
		SET active = false
		WHERE viber_user_id = $1 AND active = true
	`, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to close previous sessions: %w", err)
	}

	// Создаём новую сессию
	session := &models.ViberSession{
		ViberUserID:   user.ID,
		StartedAt:     time.Now(),
		LastMessageAt: time.Now(),
		ExpiresAt:     time.Now().Add(24 * time.Hour),
		MessageCount:  1,
		Context:       json.RawMessage("{}"),
		Active:        true,
	}

	query := `
		INSERT INTO viber_sessions (viber_user_id, started_at, last_message_at,
		                           expires_at, message_count, context, active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`

	err = sm.db.QueryRowContext(ctx, query,
		session.ViberUserID,
		session.StartedAt,
		session.LastMessageAt,
		session.ExpiresAt,
		session.MessageCount,
		session.Context,
		session.Active,
	).Scan(&session.ID, &session.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Обновляем последнюю сессию пользователя
	_, err = sm.db.ExecContext(ctx, `
		UPDATE viber_users
		SET last_session_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`, user.ID)

	return session, err
}

// UpdateSession обновляет активную сессию при новом сообщении
func (sm *SessionManager) UpdateSession(ctx context.Context, sessionID int) error {
	query := `
		UPDATE viber_sessions
		SET last_message_at = CURRENT_TIMESTAMP,
		    message_count = message_count + 1
		WHERE id = $1 AND active = true
	`

	result, err := sm.db.ExecContext(ctx, query, sessionID)
	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("session not found or not active")
	}

	return nil
}

// CloseExpiredSessions закрывает истёкшие сессии
func (sm *SessionManager) CloseExpiredSessions(ctx context.Context) error {
	query := `
		UPDATE viber_sessions
		SET active = false
		WHERE active = true AND expires_at < CURRENT_TIMESTAMP
	`

	result, err := sm.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to close expired sessions: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows > 0 {
		fmt.Printf("Closed %d expired sessions\n", rows)
	}

	return nil
}

// getOrCreateViberUser получает или создаёт пользователя Viber
func (sm *SessionManager) getOrCreateViberUser(ctx context.Context, viberID string) (*models.ViberUser, error) {
	var user models.ViberUser

	// Пытаемся найти существующего пользователя
	query := `
		SELECT id, viber_id, user_id, name, avatar_url, language,
		       country_code, api_version, subscribed, subscribed_at,
		       last_session_at, conversation_started_at, created_at, updated_at
		FROM viber_users
		WHERE viber_id = $1
	`

	err := sm.db.GetSQLXDB().GetContext(ctx, &user, query, viberID)
	if err == nil {
		return &user, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to get viber user: %w", err)
	}

	// Создаём нового пользователя
	insertQuery := `
		INSERT INTO viber_users (viber_id, language, subscribed)
		VALUES ($1, $2, $3)
		RETURNING id, viber_id, user_id, name, avatar_url, language,
		          country_code, api_version, subscribed, subscribed_at,
		          last_session_at, conversation_started_at, created_at, updated_at
	`

	err = sm.db.GetSQLXDB().GetContext(ctx, &user, insertQuery, viberID, "sr", false)
	if err != nil {
		return nil, fmt.Errorf("failed to create viber user: %w", err)
	}

	return &user, nil
}

// SaveUserInfo сохраняет информацию о пользователе из webhook
func (sm *SessionManager) SaveUserInfo(ctx context.Context, sender *models.ViberSender) error {
	query := `
		UPDATE viber_users
		SET name = $2,
		    avatar_url = $3,
		    country_code = $4,
		    language = $5,
		    api_version = $6,
		    updated_at = CURRENT_TIMESTAMP
		WHERE viber_id = $1
	`

	_, err := sm.db.ExecContext(ctx, query,
		sender.ID,
		sender.Name,
		sender.Avatar,
		sender.Country,
		sender.Language,
		sender.APIVersion,
	)

	return err
}

// SetUserSubscribed отмечает пользователя как подписанного
func (sm *SessionManager) SetUserSubscribed(ctx context.Context, viberID string, subscribed bool) error {
	query := `
		UPDATE viber_users
		SET subscribed = $2,
		    subscribed_at = CASE WHEN $2 THEN CURRENT_TIMESTAMP ELSE subscribed_at END,
		    updated_at = CURRENT_TIMESTAMP
		WHERE viber_id = $1
	`

	_, err := sm.db.ExecContext(ctx, query, viberID, subscribed)
	return err
}

// SetConversationStarted отмечает начало разговора
func (sm *SessionManager) SetConversationStarted(ctx context.Context, viberID string) error {
	query := `
		UPDATE viber_users
		SET conversation_started_at = CURRENT_TIMESTAMP,
		    updated_at = CURRENT_TIMESTAMP
		WHERE viber_id = $1
	`

	_, err := sm.db.ExecContext(ctx, query, viberID)
	return err
}

// LinkToUser связывает Viber аккаунт с пользователем системы
func (sm *SessionManager) LinkToUser(ctx context.Context, viberID string, userID int) error {
	query := `
		UPDATE viber_users
		SET user_id = $2,
		    updated_at = CURRENT_TIMESTAMP
		WHERE viber_id = $1
	`

	_, err := sm.db.ExecContext(ctx, query, viberID, userID)
	return err
}

// GetSessionStats возвращает статистику по сессиям
func (sm *SessionManager) GetSessionStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Активные сессии
	var activeSessions int
	err := sm.db.GetSQLXDB().GetContext(ctx, &activeSessions, `
		SELECT COUNT(*) FROM viber_sessions
		WHERE active = true AND expires_at > CURRENT_TIMESTAMP
	`)
	if err != nil {
		return nil, err
	}
	stats["active_sessions"] = activeSessions

	// Всего пользователей
	var totalUsers int
	err = sm.db.GetSQLXDB().GetContext(ctx, &totalUsers, `SELECT COUNT(*) FROM viber_users`)
	if err != nil {
		return nil, err
	}
	stats["total_users"] = totalUsers

	// Подписанные пользователи
	var subscribedUsers int
	err = sm.db.GetSQLXDB().GetContext(ctx, &subscribedUsers, `
		SELECT COUNT(*) FROM viber_users WHERE subscribed = true
	`)
	if err != nil {
		return nil, err
	}
	stats["subscribed_users"] = subscribedUsers

	// Сообщения за последние 24 часа
	var messages24h int
	err = sm.db.GetSQLXDB().GetContext(ctx, &messages24h, `
		SELECT COUNT(*) FROM viber_messages
		WHERE created_at > CURRENT_TIMESTAMP - INTERVAL '24 hours'
	`)
	if err != nil {
		return nil, err
	}
	stats["messages_24h"] = messages24h

	return stats, nil
}
