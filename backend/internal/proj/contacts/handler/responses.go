package handler

// ContactStatusUpdateResponse структура ответа при обновлении статуса контакта
// @Description Ответ при изменении статуса контакта
type ContactStatusUpdateResponse struct {
	// Сообщение об успешном обновлении
	Message string `json:"message" example:"contacts.statusUpdated"`
}

// ContactRemoveResponse структура ответа при удалении контакта
// @Description Ответ при удалении контакта из списка
type ContactRemoveResponse struct {
	// Сообщение об успешном удалении
	Message string `json:"message" example:"contacts.removed"`
}

// ContactStatusCheckResponse структура ответа проверки статуса контакта
// @Description Информация о статусе связи между пользователями
type ContactStatusCheckResponse struct {
	// Являются ли пользователи контактами
	AreContacts bool `json:"are_contacts" example:"true"`
	// ID текущего пользователя
	UserID int `json:"user_id" example:"123"`
	// ID проверяемого контакта
	ContactID int `json:"contact_id" example:"456"`
}
