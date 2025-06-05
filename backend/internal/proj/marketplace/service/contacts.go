package service

import (
	"backend/internal/domain/models"
	"backend/internal/storage"
	"context"
	"fmt"
)

type ContactsService struct {
	storage storage.Storage
}

func NewContactsService(storage storage.Storage) *ContactsService {
	return &ContactsService{
		storage: storage,
	}
}

// Добавить контакт
func (s *ContactsService) AddContact(ctx context.Context, userID int, req *models.AddContactRequest) (*models.UserContact, error) {
	// Проверяем, можно ли добавить этого пользователя в контакты
	canAdd, err := s.storage.CanAddContact(ctx, userID, req.ContactUserID)
	if err != nil {
		return nil, fmt.Errorf("error checking if can add contact: %w", err)
	}

	if !canAdd {
		return nil, fmt.Errorf("user does not allow contact requests or has blocked you")
	}

	// Проверяем, не пытается ли пользователь добавить себя
	if userID == req.ContactUserID {
		return nil, fmt.Errorf("cannot add yourself as contact")
	}

	// Проверяем, существует ли уже связь
	existingContact, err := s.storage.GetContact(ctx, userID, req.ContactUserID)
	if err != nil {
		// Игнорируем ошибку "контакт не найден" - это нормально
		existingContact = nil
	}

	if existingContact != nil {
		return nil, fmt.Errorf("contact already exists")
	}

	status := models.ContactStatusPending

	// Если уже есть обратная связь и она принята, сразу принимаем
	reverseContact, err := s.storage.GetContact(ctx, req.ContactUserID, userID)
	if err != nil {
		// Игнорируем ошибку "контакт не найден" - это нормально
		reverseContact = nil
	}

	if reverseContact != nil && reverseContact.Status == models.ContactStatusAccepted {
		status = models.ContactStatusAccepted
	}

	contact := &models.UserContact{
		UserID:          userID,
		ContactUserID:   req.ContactUserID,
		Status:          status,
		Notes:           req.Notes,
		AddedFromChatID: req.AddedFromChatID,
	}

	err = s.storage.AddContact(ctx, contact)
	if err != nil {
		return nil, err
	}

	// Если статус принят, также принимаем обратную связь
	if status == models.ContactStatusAccepted && reverseContact != nil {
		err = s.storage.UpdateContactStatus(ctx, req.ContactUserID, userID, models.ContactStatusAccepted, reverseContact.Notes)
		if err != nil {
			// Логируем, но не прерываем процесс
			fmt.Printf("Warning: failed to update reverse contact status: %v\n", err)
		}
	}

	// Загружаем полную информацию о контакте
	return s.storage.GetContact(ctx, userID, req.ContactUserID)
}

// Обновить статус контакта
func (s *ContactsService) UpdateContactStatus(ctx context.Context, userID int, contactUserID int, req *models.UpdateContactRequest) error {
	// ВАЖНО: userID - это текущий пользователь, который выполняет действие
	// contactUserID - это ID пользователя, чей запрос мы принимаем/отклоняем

	// Проверяем, существует ли запрос ОТ contactUserID К userID
	// То есть текущий пользователь (userID) принимает/отклоняет запрос от contactUserID
	existingContact, err := s.storage.GetContact(ctx, contactUserID, userID)
	if err != nil {
		return err
	}

	if existingContact == nil {
		return fmt.Errorf("contact request not found")
	}

	// Проверяем, что запрос в статусе pending
	if existingContact.Status != models.ContactStatusPending {
		return fmt.Errorf("can only update pending requests")
	}

	// Обновляем статус запроса от contactUserID к userID
	err = s.storage.UpdateContactStatus(ctx, contactUserID, userID, req.Status, req.Notes)
	if err != nil {
		return err
	}

	// Если статус принят, создаём взаимную связь
	if req.Status == models.ContactStatusAccepted {
		// Проверяем, существует ли уже обратная связь
		reverseContact, err := s.storage.GetContact(ctx, userID, contactUserID)
		if err != nil {
			// Игнорируем ошибку - обратный контакт может не существовать
			reverseContact = nil
		}

		if reverseContact == nil {
			// Создаём обратную связь со статусом accepted
			reverseContactData := &models.UserContact{
				UserID:        userID,
				ContactUserID: contactUserID,
				Status:        models.ContactStatusAccepted,
				Notes:         req.Notes,
			}
			err = s.storage.AddContact(ctx, reverseContactData)
			if err != nil {
				// Логируем, но не прерываем процесс
				fmt.Printf("Warning: failed to create reverse contact: %v\n", err)
			}
		} else if reverseContact.Status != models.ContactStatusAccepted {
			// Если обратная связь существует, но не принята, обновляем её статус
			err = s.storage.UpdateContactStatus(ctx, userID, contactUserID, models.ContactStatusAccepted, req.Notes)
			if err != nil {
				// Логируем, но не прерываем процесс
				fmt.Printf("Warning: failed to update reverse contact status: %v\n", err)
			}
		}
	}

	return nil
}

// Получить список контактов
func (s *ContactsService) GetContacts(ctx context.Context, userID int, status string, page, limit int) (*models.ContactsListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	contacts, total, err := s.storage.GetUserContacts(ctx, userID, status, page, limit)
	if err != nil {
		return nil, err
	}

	return &models.ContactsListResponse{
		Contacts: contacts,
		Total:    total,
		Page:     page,
		Limit:    limit,
	}, nil
}

// Удалить контакт
func (s *ContactsService) RemoveContact(ctx context.Context, userID, contactUserID int) error {
	return s.storage.RemoveContact(ctx, userID, contactUserID)
}

// Получить настройки приватности
func (s *ContactsService) GetPrivacySettings(ctx context.Context, userID int) (*models.UserPrivacySettings, error) {
	return s.storage.GetUserPrivacySettings(ctx, userID)
}

// Обновить настройки приватности
func (s *ContactsService) UpdatePrivacySettings(ctx context.Context, userID int, req *models.UpdatePrivacySettingsRequest) (*models.UserPrivacySettings, error) {
	err := s.storage.UpdateUserPrivacySettings(ctx, userID, req)
	if err != nil {
		return nil, err
	}

	return s.storage.GetUserPrivacySettings(ctx, userID)
}

// Проверить, являются ли пользователи контактами
func (s *ContactsService) AreContacts(ctx context.Context, userID1, userID2 int) (bool, error) {
	contact1, err := s.storage.GetContact(ctx, userID1, userID2)
	if err != nil {
		return false, err
	}

	contact2, err := s.storage.GetContact(ctx, userID2, userID1)
	if err != nil {
		return false, err
	}

	return contact1 != nil && contact1.Status == models.ContactStatusAccepted &&
		contact2 != nil && contact2.Status == models.ContactStatusAccepted, nil
}
