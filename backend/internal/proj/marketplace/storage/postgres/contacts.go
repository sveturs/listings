package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"backend/internal/domain/models"
	"backend/internal/logger"
)

// ErrContactNotFound возвращается когда контакт не найден
var ErrContactNotFound = errors.New("contact not found")

// Добавить контакт
func (s *Storage) AddContact(ctx context.Context, contact *models.UserContact) error {
	query := `
		INSERT INTO user_contacts (
			user_id, contact_user_id, status, notes, added_from_chat_id
		) VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, contact_user_id) 
		DO UPDATE SET 
			status = EXCLUDED.status,
			notes = EXCLUDED.notes,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, created_at, updated_at
	`

	err := s.pool.QueryRow(ctx, query,
		contact.UserID,
		contact.ContactUserID,
		contact.Status,
		contact.Notes,
		contact.AddedFromChatID,
	).Scan(&contact.ID, &contact.CreatedAt, &contact.UpdatedAt)
	if err != nil {
		return fmt.Errorf("error adding contact: %w", err)
	}

	return nil
}

// Обновить статус контакта
func (s *Storage) UpdateContactStatus(ctx context.Context, userID, contactUserID int, status, notes string) error {
	query := `
		UPDATE user_contacts 
		SET status = $3, notes = $4, updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $1 AND contact_user_id = $2
	`

	result, err := s.pool.Exec(ctx, query, userID, contactUserID, status, notes)
	if err != nil {
		return fmt.Errorf("error updating contact status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("contact not found")
	}

	return nil
}

// Получить контакт
func (s *Storage) GetContact(ctx context.Context, userID, contactUserID int) (*models.UserContact, error) {
	query := `
		SELECT 
			uc.id, uc.user_id, uc.contact_user_id, uc.status, 
			uc.created_at, uc.updated_at, uc.notes, uc.added_from_chat_id,
			'' as contact_name, '' as contact_email, '' as contact_picture
		FROM user_contacts uc
		WHERE uc.user_id = $1 AND uc.contact_user_id = $2
	`

	contact := &models.UserContact{
		ContactUser: &models.User{},
	}

	var contactPicture sql.NullString

	err := s.pool.QueryRow(ctx, query, userID, contactUserID).Scan(
		&contact.ID,
		&contact.UserID,
		&contact.ContactUserID,
		&contact.Status,
		&contact.CreatedAt,
		&contact.UpdatedAt,
		&contact.Notes,
		&contact.AddedFromChatID,
		&contact.ContactUser.Name,
		&contact.ContactUser.Email,
		&contactPicture,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrContactNotFound // Контакт не найден
		}
		return nil, fmt.Errorf("error getting contact: %w", err)
	}

	contact.ContactUser.ID = contact.ContactUserID
	contact.ContactUser.PictureURL = contactPicture.String

	return contact, nil
}

// Получить список контактов пользователя
func (s *Storage) GetUserContacts(ctx context.Context, userID int, status string, page, limit int) ([]models.UserContact, int, error) {
	offset := (page - 1) * limit

	// Ищем контакты где пользователь является либо отправителем, либо получателем
	// Для принятых контактов используем DISTINCT чтобы избежать дубликатов
	whereClause := "WHERE (uc.user_id = $1 OR uc.contact_user_id = $1)"
	args := []interface{}{userID}
	argIndex := 2

	if status != "" {
		whereClause += fmt.Sprintf(" AND uc.status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	// Для принятых контактов считаем только уникальные пары пользователей
	countQuery := ""
	if status == "accepted" || status == "" {
		// Для всех контактов и принятых используем DISTINCT
		countQuery = fmt.Sprintf(`
			SELECT COUNT(DISTINCT 
				CASE 
					WHEN uc.user_id = $1 THEN uc.contact_user_id
					ELSE uc.user_id
				END
			)
			FROM user_contacts uc 
			%s
		`, whereClause)
	} else {
		countQuery = fmt.Sprintf(`
			SELECT COUNT(*) 
			FROM user_contacts uc 
			%s
		`, whereClause)
	}

	var total int
	err := s.pool.QueryRow(ctx, countQuery, args[:argIndex-1]...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting contacts: %w", err)
	}

	// Получаем контакты с пагинацией
	// Для принятых контактов берем только одну запись из пары
	query := ""
	if status == "accepted" || status == "" {
		// Для всех контактов или только принятых используем DISTINCT для избежания дубликатов
		query = fmt.Sprintf(`
			SELECT DISTINCT ON (
				CASE 
					WHEN uc.user_id = $1 THEN uc.contact_user_id
					ELSE uc.user_id
				END
			)
				uc.id, uc.user_id, uc.contact_user_id, uc.status, 
				uc.created_at, uc.updated_at, uc.notes, uc.added_from_chat_id,
				'' as contact_name,
				'' as contact_email,
				'' as contact_picture
			FROM user_contacts uc
			%s
			ORDER BY 
				CASE 
					WHEN uc.user_id = $1 THEN uc.contact_user_id
					ELSE uc.user_id
				END,
				uc.updated_at DESC
			LIMIT $%d OFFSET $%d
		`, whereClause, argIndex, argIndex+1)
	} else {
		// Для других статусов используем обычный запрос
		query = fmt.Sprintf(`
			SELECT 
				uc.id, uc.user_id, uc.contact_user_id, uc.status, 
				uc.created_at, uc.updated_at, uc.notes, uc.added_from_chat_id,
				'' as contact_name,
				'' as contact_email,
				'' as contact_picture
			FROM user_contacts uc
			%s
			ORDER BY uc.updated_at DESC
			LIMIT $%d OFFSET $%d
		`, whereClause, argIndex, argIndex+1)
	}

	args = append(args, limit, offset)

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying contacts: %w", err)
	}
	defer rows.Close()

	var contacts []models.UserContact
	for rows.Next() {
		contact := models.UserContact{
			ContactUser: &models.User{},
		}

		var contactPicture sql.NullString

		err := rows.Scan(
			&contact.ID,
			&contact.UserID,
			&contact.ContactUserID,
			&contact.Status,
			&contact.CreatedAt,
			&contact.UpdatedAt,
			&contact.Notes,
			&contact.AddedFromChatID,
			&contact.ContactUser.Name,
			&contact.ContactUser.Email,
			&contactPicture,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning contact: %w", err)
		}

		// Устанавливаем ID "другого" пользователя
		if contact.UserID == userID {
			contact.ContactUser.ID = contact.ContactUserID
		} else {
			contact.ContactUser.ID = contact.UserID
		}
		contact.ContactUser.PictureURL = contactPicture.String

		contacts = append(contacts, contact)
	}

	return contacts, total, nil
}

// Удалить контакт
func (s *Storage) RemoveContact(ctx context.Context, userID, contactUserID int) error {
	// Начинаем транзакцию
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			// Игнорируем ошибку если транзакция уже была завершена
			_ = err // Explicitly ignore error
		}
	}()

	// Удаляем контакт в обоих направлениях
	query := `DELETE FROM user_contacts WHERE (user_id = $1 AND contact_user_id = $2) OR (user_id = $2 AND contact_user_id = $1)`

	result, err := tx.Exec(ctx, query, userID, contactUserID)
	if err != nil {
		return fmt.Errorf("error removing contact: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("contact not found")
	}

	// Коммитим транзакцию
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

// Получить/создать настройки приватности
func (s *Storage) GetUserPrivacySettings(ctx context.Context, userID int) (*models.UserPrivacySettings, error) {
	// User existence проверяется через JWT auth middleware
	// Если запрос дошел сюда, пользователь уже аутентифицирован в auth-service

	// Пытаемся получить существующие настройки (включая JSONB settings)
	selectQuery := `
		SELECT user_id, allow_contact_requests, allow_messages_from_contacts_only, COALESCE(settings, '{}'::jsonb), created_at, updated_at
		FROM user_privacy_settings
		WHERE user_id = $1
	`

	settings := &models.UserPrivacySettings{}
	var settingsJSON []byte
	err := s.pool.QueryRow(ctx, selectQuery, userID).Scan(
		&settings.UserID,
		&settings.AllowContactRequests,
		&settings.AllowMessagesFromContactsOnly,
		&settingsJSON,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)

	if err == nil {
		// Парсим JSONB settings
		if len(settingsJSON) > 0 {
			if err := json.Unmarshal(settingsJSON, &settings.Settings); err != nil {
				logger.Warn().Err(err).Msg("Failed to parse settings JSONB, using empty map")
				settings.Settings = make(map[string]interface{})
			}
		} else {
			settings.Settings = make(map[string]interface{})
		}
	}

	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
		// Если настроек нет, создаем их с значениями по умолчанию
		insertQuery := `
			INSERT INTO user_privacy_settings (user_id, allow_contact_requests, allow_messages_from_contacts_only)
			VALUES ($1, true, false)
			RETURNING user_id, allow_contact_requests, allow_messages_from_contacts_only, COALESCE(settings, '{}'::jsonb), created_at, updated_at
		`

		err = s.pool.QueryRow(ctx, insertQuery, userID).Scan(
			&settings.UserID,
			&settings.AllowContactRequests,
			&settings.AllowMessagesFromContactsOnly,
			&settingsJSON,
			&settings.CreatedAt,
			&settings.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error creating privacy settings: %w", err)
		}

		// Парсим JSONB settings для новой записи
		if len(settingsJSON) > 0 {
			if err := json.Unmarshal(settingsJSON, &settings.Settings); err != nil {
				logger.Warn().Err(err).Msg("Failed to parse settings JSONB for new record")
				settings.Settings = make(map[string]interface{})
			}
		} else {
			settings.Settings = make(map[string]interface{})
		}
	} else if err != nil {
		return nil, fmt.Errorf("error getting privacy settings: %w", err)
	}

	return settings, nil
}

// Обновить настройки приватности
func (s *Storage) UpdateUserPrivacySettings(ctx context.Context, userID int, settings *models.UpdatePrivacySettingsRequest) error {
	// Проверяем, есть ли что обновлять
	if settings.AllowContactRequests == nil && settings.AllowMessagesFromContactsOnly == nil {
		return fmt.Errorf("no settings to update")
	}

	// Упрощенный запрос
	query := `
		INSERT INTO user_privacy_settings (user_id, allow_contact_requests, allow_messages_from_contacts_only)
		VALUES ($1, COALESCE($2, true), COALESCE($3, false))
		ON CONFLICT (user_id)
		DO UPDATE SET
			allow_contact_requests = COALESCE($2, user_privacy_settings.allow_contact_requests),
			allow_messages_from_contacts_only = COALESCE($3, user_privacy_settings.allow_messages_from_contacts_only),
			updated_at = CURRENT_TIMESTAMP
	`

	_, err := s.pool.Exec(ctx, query, userID, settings.AllowContactRequests, settings.AllowMessagesFromContactsOnly)
	if err != nil {
		return fmt.Errorf("error updating privacy settings: %w", err)
	}

	return nil
}

// UpdateChatSettings обновляет настройки чата в JSONB поле settings
func (s *Storage) UpdateChatSettings(ctx context.Context, userID int, settings *models.ChatUserSettings) error {
	// Сначала убеждаемся что запись существует (создаст если нет)
	_, err := s.GetUserPrivacySettings(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get/create privacy settings: %w", err)
	}

	// Обновляем JSONB поле settings используя jsonb_set
	query := `
		UPDATE user_privacy_settings
		SET settings = jsonb_set(
			jsonb_set(
				jsonb_set(
					jsonb_set(
						COALESCE(settings, '{}'::jsonb),
						'{auto_translate_chat}',
						to_jsonb($2::boolean)
					),
					'{preferred_language}',
					to_jsonb($3::text)
				),
				'{show_original_language_badge}',
				to_jsonb($4::boolean)
			),
			'{chat_tone_moderation}',
			to_jsonb($5::boolean)
		),
		updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $1
	`

	_, err = s.pool.Exec(ctx, query,
		userID,
		settings.AutoTranslate,
		settings.PreferredLanguage,
		settings.ShowLanguageBadge,
		settings.ModerateTone,
	)
	if err != nil {
		return fmt.Errorf("failed to update chat settings: %w", err)
	}

	return nil
}

// Проверить, разрешены ли запросы на добавление в контакты
func (s *Storage) CanAddContact(ctx context.Context, userID, targetUserID int) (bool, error) {
	logger.Debug().
		Int("userID", userID).
		Int("targetUserID", targetUserID).
		Msg("[Storage] CanAddContact")

	// Получаем настройки приватности целевого пользователя
	settings, err := s.GetUserPrivacySettings(ctx, targetUserID)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("[Storage] GetUserPrivacySettings error")
		return false, err
	}

	logger.Debug().
		Int("targetUserID", targetUserID).
		Bool("AllowContactRequests", settings.AllowContactRequests).
		Msg("[Storage] Privacy settings for user")

	// Проверяем, разрешены ли запросы на добавление
	if !settings.AllowContactRequests {
		return false, nil
	}

	// Проверяем, не заблокирован ли инициатор
	existingContact, err := s.GetContact(ctx, targetUserID, userID)
	if err != nil {
		// Если метод GetContact вернул nil, nil - это нормально (контакт не найден)
		// Если GetContact вернул ошибку, логируем её но продолжаем работу
		// так как отсутствие контакта не должно блокировать добавление
		existingContact = nil
	}

	if existingContact != nil && existingContact.Status == models.ContactStatusBlocked {
		return false, nil
	}

	return true, nil
}

// GetIncomingContactRequests получает входящие запросы в контакты
func (s *Storage) GetIncomingContactRequests(ctx context.Context, userID int, page, limit int) ([]models.UserContact, int, error) {
	offset := (page - 1) * limit

	// Получаем только входящие запросы со статусом pending
	// Входящие = где текущий пользователь является contact_user_id
	countQuery := `
		SELECT COUNT(*)
		FROM user_contacts uc
		WHERE uc.contact_user_id = $1 AND uc.status = 'pending'
	`

	var total int
	err := s.pool.QueryRow(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting incoming requests: %w", err)
	}

	// Получаем запросы с информацией об отправителе
	query := `
		SELECT
			uc.id, uc.user_id, uc.contact_user_id, uc.status,
			uc.created_at, uc.updated_at, uc.notes, uc.added_from_chat_id,
			'' as sender_name,
			'' as sender_email,
			'' as sender_picture
		FROM user_contacts uc
		WHERE uc.contact_user_id = $1 AND uc.status = 'pending'
		ORDER BY uc.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := s.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying incoming requests: %w", err)
	}
	defer rows.Close()

	var contacts []models.UserContact
	for rows.Next() {
		contact := models.UserContact{
			User: &models.User{},
		}

		var senderPicture sql.NullString

		err := rows.Scan(
			&contact.ID,
			&contact.UserID,
			&contact.ContactUserID,
			&contact.Status,
			&contact.CreatedAt,
			&contact.UpdatedAt,
			&contact.Notes,
			&contact.AddedFromChatID,
			&contact.User.Name,
			&contact.User.Email,
			&senderPicture,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning incoming request: %w", err)
		}

		contact.User.ID = contact.UserID
		contact.User.PictureURL = senderPicture.String

		contacts = append(contacts, contact)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating incoming requests: %w", err)
	}

	return contacts, total, nil
}

// AreContacts проверяет, являются ли пользователи контактами
func (s *Storage) AreContacts(ctx context.Context, userID1, userID2 int) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM user_contacts
			WHERE user_id = $1 AND contact_user_id = $2 AND status = 'accepted'
		)
	`

	var exists bool
	err := s.pool.QueryRow(ctx, query, userID1, userID2).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
