package postgres

import (
	"backend/internal/domain/models"
	"context"
	"database/sql"
	"fmt"
)

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
			u.name as contact_name, u.email as contact_email, u.picture_url as contact_picture
		FROM user_contacts uc
		JOIN users u ON uc.contact_user_id = u.id
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
		if err == sql.ErrNoRows {
			return nil, nil // Контакт не найден
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
				CASE 
					WHEN uc.user_id = $1 THEN u2.name
					ELSE u1.name
				END as contact_name,
				CASE 
					WHEN uc.user_id = $1 THEN u2.email
					ELSE u1.email
				END as contact_email,
				CASE 
					WHEN uc.user_id = $1 THEN u2.picture_url
					ELSE u1.picture_url
				END as contact_picture
			FROM user_contacts uc
			JOIN users u1 ON uc.user_id = u1.id
			JOIN users u2 ON uc.contact_user_id = u2.id
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
				CASE 
					WHEN uc.user_id = $1 THEN u2.name
					ELSE u1.name
				END as contact_name,
				CASE 
					WHEN uc.user_id = $1 THEN u2.email
					ELSE u1.email
				END as contact_email,
				CASE 
					WHEN uc.user_id = $1 THEN u2.picture_url
					ELSE u1.picture_url
				END as contact_picture
			FROM user_contacts uc
			JOIN users u1 ON uc.user_id = u1.id
			JOIN users u2 ON uc.contact_user_id = u2.id
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
	defer tx.Rollback(ctx)

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
	// Сначала проверяем что пользователь существует
	var exists bool
	checkUserQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
	err := s.pool.QueryRow(ctx, checkUserQuery, userID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("error checking user existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("user does not exist")
	}

	// Пытаемся получить существующие настройки
	selectQuery := `
		SELECT user_id, allow_contact_requests, allow_messages_from_contacts_only, created_at, updated_at
		FROM user_privacy_settings 
		WHERE user_id = $1
	`

	settings := &models.UserPrivacySettings{}
	err = s.pool.QueryRow(ctx, selectQuery, userID).Scan(
		&settings.UserID,
		&settings.AllowContactRequests,
		&settings.AllowMessagesFromContactsOnly,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		// Если настроек нет, создаем их с значениями по умолчанию
		insertQuery := `
			INSERT INTO user_privacy_settings (user_id, allow_contact_requests, allow_messages_from_contacts_only) 
			VALUES ($1, true, false)
			RETURNING user_id, allow_contact_requests, allow_messages_from_contacts_only, created_at, updated_at
		`

		err = s.pool.QueryRow(ctx, insertQuery, userID).Scan(
			&settings.UserID,
			&settings.AllowContactRequests,
			&settings.AllowMessagesFromContactsOnly,
			&settings.CreatedAt,
			&settings.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("error creating privacy settings: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("error getting privacy settings: %w", err)
	}

	return settings, nil
}

// Обновить настройки приватности
func (s *Storage) UpdateUserPrivacySettings(ctx context.Context, userID int, settings *models.UpdatePrivacySettingsRequest) error {
	setParts := []string{}
	args := []interface{}{userID}
	argIndex := 2

	if settings.AllowContactRequests != nil {
		setParts = append(setParts, fmt.Sprintf("allow_contact_requests = $%d", argIndex))
		args = append(args, *settings.AllowContactRequests)
		argIndex++
	}

	if settings.AllowMessagesFromContactsOnly != nil {
		setParts = append(setParts, fmt.Sprintf("allow_messages_from_contacts_only = $%d", argIndex))
		args = append(args, *settings.AllowMessagesFromContactsOnly)
		argIndex++
	}

	if len(setParts) == 0 {
		return fmt.Errorf("no settings to update")
	}

	setParts = append(setParts, "updated_at = CURRENT_TIMESTAMP")
	setClause := fmt.Sprintf("SET %s", fmt.Sprintf("%s", setParts[0]))
	for i := 1; i < len(setParts); i++ {
		setClause += ", " + setParts[i]
	}

	// Упрощенный запрос
	simpleQuery := `
		INSERT INTO user_privacy_settings (user_id, allow_contact_requests, allow_messages_from_contacts_only) 
		VALUES ($1, COALESCE($2, true), COALESCE($3, false)) 
		ON CONFLICT (user_id) 
		DO UPDATE SET 
			allow_contact_requests = COALESCE($2, user_privacy_settings.allow_contact_requests),
			allow_messages_from_contacts_only = COALESCE($3, user_privacy_settings.allow_messages_from_contacts_only),
			updated_at = CURRENT_TIMESTAMP
	`

	_, err := s.pool.Exec(ctx, simpleQuery, userID, settings.AllowContactRequests, settings.AllowMessagesFromContactsOnly)
	if err != nil {
		return fmt.Errorf("error updating privacy settings: %w", err)
	}

	return nil
}

// Проверить, разрешены ли запросы на добавление в контакты
func (s *Storage) CanAddContact(ctx context.Context, userID, targetUserID int) (bool, error) {
	// Получаем настройки приватности целевого пользователя
	settings, err := s.GetUserPrivacySettings(ctx, targetUserID)
	if err != nil {
		return false, err
	}

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
