// backend/internal/proj/marketplace/handler/chat.go

package handler

import (
	"backend/internal/config"
	"backend/internal/domain/models"
	"backend/internal/logger"
	globalService "backend/internal/proj/global/service"
	"backend/pkg/utils"
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"strconv"
	"sync"
	"time"
)

type ChatHandler struct {
	services globalService.ServicesInterface
	config   *config.Config
}

func NewChatHandler(services globalService.ServicesInterface, config *config.Config) *ChatHandler {
	return &ChatHandler{
		services: services,
		config:   config,
	}
}

// GetChats возвращает список чатов пользователя
// @Summary Get user's chats
// @Description Returns all chats where the user is a participant
// @Tags marketplace-chat
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=[]models.MarketplaceChat} "List of chats"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.getChatsError"
// @Security BearerAuth
// @Router /api/v1/marketplace/chat [get]
func (h *ChatHandler) GetChats(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	logger.Info().Int("userId", userID).Msg("GetChats called")

	chats, err := h.services.Chat().GetChats(c.Context(), userID)
	if err != nil {
		logger.Error().Err(err).Int("userId", userID).Msg("Error in GetChats")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.getChatsError")
	}

	logger.Info().Int("userId", userID).Int("chatsCount", len(chats)).Msg("GetChats successful")
	return utils.SuccessResponse(c, chats)
}

// GetMessages возвращает сообщения чата
// @Summary Get chat messages
// @Description Returns paginated messages from a chat
// @Tags marketplace-chat
// @Accept json
// @Produce json
// @Param chat_id query string false "Chat ID"
// @Param listing_id query string false "Listing ID"
// @Param receiver_id query string false "Receiver ID for direct messages"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} object{data=[]models.MarketplaceMessage,meta=object{page=int,limit=int,hasMore=bool,total=int}} "Chat messages"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidChatId or marketplace.invalidListingId or marketplace.chatParamsRequired"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.getMessagesError"
// @Security BearerAuth
// @Router /api/v1/marketplace/chat/messages [get]
func (h *ChatHandler) GetMessages(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	chatID := c.Query("chat_id")
	listingID := c.Query("listing_id")

	// Если есть chat_id, используем его
	if chatID != "" {
		chatIDInt, err := strconv.Atoi(chatID)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidChatId")
		}
		c.Context().SetUserValue("chat_id", chatIDInt)

		// Если указан chat_id, listing_id не обязателен
		// Получим listing_id из чата
		chat, err := h.services.Chat().GetChat(c.Context(), chatIDInt, userID)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.getChatError")
		}
		listingID = strconv.Itoa(chat.ListingID)
	}

	// Получаем receiver_id для прямых сообщений
	receiverID := c.Query("receiver_id")

	// Если нет ни chat_id, ни listing_id, ни receiver_id - ошибка
	if listingID == "" && chatID == "" && receiverID == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.chatParamsRequired")
	}

	listingIDInt := 0
	if listingID != "" {
		var err error
		listingIDInt, err = strconv.Atoi(listingID)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidListingId")
		}
	}

	// Убедимся, что пагинация не даст отрицательный offset
	page := utils.StringToInt(c.Query("page"), 1)
	if page < 1 {
		page = 1
	}

	limit := utils.StringToInt(c.Query("limit"), 20)
	if limit < 1 {
		limit = 20
	}

	offset := (page - 1) * limit

	logger.Info().
		Int("page", page).
		Int("limit", limit).
		Int("offset", offset).
		Str("chatId", c.Query("chat_id")).
		Int("listingId", listingIDInt).
		Int("userId", userID).
		Msg("GetMessages")

	// Создаем новый context.Context с chat_id
	ctx := context.Background()
	if chatID != "" {
		// Преобразуем строку в int для контекста
		if chatIDInt, err := strconv.Atoi(chatID); err == nil {
			ctx = context.WithValue(ctx, "chat_id", chatIDInt)
		}
	}

	messages, err := h.services.Chat().GetMessages(ctx, listingIDInt, userID, offset, limit)
	if err != nil {
		logger.Error().Err(err).Msg("Error fetching messages")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.getMessagesError")
	}

	logger.Info().Int("messagesCount", len(messages)).Msg("GetMessages: returned messages")

	// Загружаем вложения для каждого сообщения
	for i := range messages {
		if messages[i].HasAttachments {
			attachments, err := h.services.Storage().GetMessageAttachments(c.Context(), messages[i].ID)
			if err == nil && len(attachments) > 0 {
				// Конвертируем []*ChatAttachment в []ChatAttachment
				messages[i].Attachments = make([]models.ChatAttachment, len(attachments))
				for j, att := range attachments {
					messages[i].Attachments[j] = *att
				}
			}
		}
	}

	// Получаем общее количество сообщений, если есть chat_id
	var total int = -1 // По умолчанию -1 означает, что количество неизвестно
	if chatIDStr := c.Query("chat_id"); chatIDStr != "" {
		if _, err := strconv.Atoi(chatIDStr); err == nil {
			// TODO: добавить метод в сервис для получения общего количества
			// Пока что используем -1, что заставит фронтенд определять hasMore по количеству возвращенных сообщений
		}
	}

	// Возвращаем структурированный ответ
	response := fiber.Map{
		"messages": messages,
		"total":    total,
		"page":     page,
		"limit":    limit,
	}

	return utils.SuccessResponse(c, response)
}

// GetUnreadCount возвращает количество непрочитанных сообщений
// @Summary Get unread messages count
// @Description Returns the number of unread messages for the user
// @Tags marketplace-chat
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=object{count=int}} "Unread count"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.getUnreadCountError"
// @Security BearerAuth
// @Router /api/v1/marketplace/messages/unread [get]
func (h *ChatHandler) GetUnreadCount(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	count, err := h.services.Chat().GetUnreadMessagesCount(c.Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.getUnreadCountError")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"count": count,
	})
}

// SendMessage отправляет сообщение в чат
// @Summary Send chat message
// @Description Sends a new message to a chat
// @Tags marketplace-chat
// @Accept json
// @Produce json
// @Param body body models.CreateMessageRequest true "Message data"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.MarketplaceMessage} "Sent message"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidRequest or marketplace.receiverRequired"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.sendMessageError"
// @Security BearerAuth
// @Router /api/v1/marketplace/messages [post]
func (h *ChatHandler) SendMessage(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	var req models.CreateMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidRequest")
	}

	// Валидация входных данных
	if req.ReceiverID == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.receiverRequired")
	}

	// Санитизация контента для защиты от XSS
	req.Content = utils.SanitizeText(req.Content)

	// Для прямых сообщений между контактами достаточно ReceiverID
	// ListingID или ChatID нужны только для чатов с привязкой к объявлению
	// Если нет ни ListingID, ни ChatID - это прямое сообщение контакту

	msg := &models.MarketplaceMessage{
		ChatID:     req.ChatID,
		ListingID:  req.ListingID,
		SenderID:   userID,
		ReceiverID: req.ReceiverID,
		Content:    req.Content,
		Sender:     &models.User{}, // Инициализируем структуры
		Receiver:   &models.User{},
	}

	if err := h.services.Chat().SendMessage(c.Context(), msg); err != nil {
		logger.Error().Err(err).Msg("Error sending message")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.sendMessageError")
	}

	return utils.SuccessResponse(c, msg)
}

// MarkAsRead отмечает сообщения как прочитанные
// @Summary Mark messages as read
// @Description Marks specified messages as read
// @Tags marketplace-chat
// @Accept json
// @Produce json
// @Param body body models.MarkAsReadRequest true "Message IDs to mark as read"
// @Success 200 {object} utils.SuccessResponseSwag{data=object{message=string}} "Messages marked as read"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidRequest"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.markAsReadError"
// @Security BearerAuth
// @Router /api/v1/marketplace/messages/read [post]
func (h *ChatHandler) MarkAsRead(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	var req models.MarkAsReadRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidRequest")
	}

	if err := h.services.Chat().MarkMessagesAsRead(c.Context(), req.MessageIDs, userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.markAsReadError")
	}

	return utils.SuccessResponse(c, fiber.Map{"message": "marketplace.messagesMarkedAsRead"})
}
// ArchiveChat архивирует чат
// @Summary Archive chat
// @Description Archives a chat for the user
// @Tags marketplace-chat
// @Accept json
// @Produce json
// @Param chat_id path int true "Chat ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=object{message=string}} "Chat archived"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidChatId"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.archiveChatError"
// @Security BearerAuth
// @Router /api/v1/marketplace/chats/{chat_id}/archive [post]
func (h *ChatHandler) ArchiveChat(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	chatID, err := c.ParamsInt("chat_id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidChatId")
	}

	err = h.services.Chat().ArchiveChat(c.Context(), chatID, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.archiveChatError")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "marketplace.chatArchived",
	})
}

// UploadAttachments загружает вложения для сообщения
// @Summary Upload message attachments
// @Description Uploads attachments for a chat message
// @Tags marketplace-chat
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Message ID"
// @Param files formData file true "Files to upload" 
// @Success 200 {object} utils.SuccessResponseSwag{data=[]models.ChatAttachment} "Uploaded attachments"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidMessageId or marketplace.noFilesUploaded or marketplace.tooManyFiles"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 404 {object} utils.ErrorResponseSwag "marketplace.messageNotFound"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.uploadAttachmentsError"
// @Security BearerAuth
// @Router /api/v1/marketplace/messages/{id}/attachments [post]
func (h *ChatHandler) UploadAttachments(c *fiber.Ctx) error {
	logger.Info().Msg("UploadAttachments called")
	userID := c.Locals("user_id").(int)
	messageID, err := c.ParamsInt("id")
	if err != nil {
		logger.Error().Err(err).Msg("Error parsing message ID")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidMessageId")
	}
	logger.Info().Int("userId", userID).Int("messageId", messageID).Msg("UploadAttachments")

	// Получаем сообщение для проверки прав
	message, err := h.services.Storage().GetMessageByID(c.Context(), messageID)
	if err != nil {
		logger.Error().Err(err).Int("messageId", messageID).Msg("Error getting message by ID")
		return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.messageNotFound")
	}

	// Проверяем, что пользователь является отправителем сообщения
	if message.SenderID != userID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "marketplace.attachmentsForbidden")
	}

	// Получаем файлы из запроса
	form, err := c.MultipartForm()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidFormData")
	}

	files := form.File["files"]
	if len(files) == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.noFilesUploaded")
	}

	// Ограничение на количество файлов
	if len(files) > 10 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.tooManyFiles")
	}

	// Загружаем файлы через сервис
	attachments, err := h.services.ChatAttachment().UploadAttachments(c.Context(), messageID, files)
	if err != nil {
		logger.Error().Err(err).Msg("Error uploading attachments")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.uploadAttachmentsError")
	}

	// Отправляем обновленное сообщение через WebSocket
	if len(attachments) > 0 {
		logger.Info().Int("attachmentsCount", len(attachments)).Int("messageId", messageID).Msg("UploadAttachments: uploading attachments")
		// Получаем обновленное сообщение с вложениями
		updatedMessage, err := h.services.Storage().GetMessageByID(c.Context(), messageID)
		if err == nil {
			logger.Info().Int("senderId", updatedMessage.SenderID).Int("receiverId", updatedMessage.ReceiverID).Msg("UploadAttachments: got message from DB")
			// Конвертируем вложения для сообщения
			updatedMessage.Attachments = make([]models.ChatAttachment, len(attachments))
			for i, att := range attachments {
				updatedMessage.Attachments[i] = *att
			}
			updatedMessage.HasAttachments = true
			updatedMessage.AttachmentsCount = len(attachments)

			logger.Info().Int("attachmentsCount", len(updatedMessage.Attachments)).Msg("UploadAttachments: broadcasting message")
			// Отправляем обновленное сообщение через WebSocket
			h.services.Chat().BroadcastMessage(updatedMessage)
		} else {
			logger.Error().Err(err).Msg("UploadAttachments: error getting message by ID")
		}
	}

	return utils.SuccessResponse(c, attachments)
}

// GetAttachment получает информацию о вложении
// @Summary Get attachment info
// @Description Returns information about a specific attachment
// @Tags marketplace-chat
// @Accept json
// @Produce json
// @Param id path int true "Attachment ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.ChatAttachment} "Attachment info"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidAttachmentId"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 403 {object} utils.ErrorResponseSwag "marketplace.accessDenied"
// @Failure 404 {object} utils.ErrorResponseSwag "marketplace.attachmentNotFound or marketplace.messageNotFound"
// @Security BearerAuth
// @Router /api/v1/marketplace/attachments/{id} [get]
func (h *ChatHandler) GetAttachment(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	attachmentID, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidAttachmentId")
	}

	// Получаем вложение
	attachment, err := h.services.ChatAttachment().GetAttachment(c.Context(), attachmentID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.attachmentNotFound")
	}

	// Проверяем доступ к вложению через сообщение
	message, err := h.services.Storage().GetMessageByID(c.Context(), attachment.MessageID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.messageNotFound")
	}

	// Пользователь должен быть участником чата
	if message.SenderID != userID && message.ReceiverID != userID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "marketplace.accessDenied")
	}

	return utils.SuccessResponse(c, attachment)
}

// DeleteAttachment удаляет вложение
// @Summary Delete attachment
// @Description Deletes an attachment (only by the message sender)
// @Tags marketplace-chat
// @Accept json
// @Produce json
// @Param id path int true "Attachment ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=object{message=string}} "Attachment deleted"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidAttachmentId"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 403 {object} utils.ErrorResponseSwag "marketplace.deleteAttachmentForbidden"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.deleteAttachmentError"
// @Security BearerAuth
// @Router /api/v1/marketplace/attachments/{id} [delete]
func (h *ChatHandler) DeleteAttachment(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	attachmentID, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidAttachmentId")
	}

	if err := h.services.ChatAttachment().DeleteAttachment(c.Context(), attachmentID, userID); err != nil {
		if err.Error() == "permission denied" {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "marketplace.deleteAttachmentForbidden")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.deleteAttachmentError")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "marketplace.attachmentDeleted",
	})
}

// HandleWebSocketWithAuth - WebSocket хендлер с переданным userID
func (h *ChatHandler) HandleWebSocketWithAuth(c *websocket.Conn, userID int) {
	if c == nil {
		return // Защита от nil pointer
	}

	// Проверяем, что userID валидный
	if userID == 0 {
		logger.Warn().Int("userId", userID).Msg("WebSocket: Invalid user_id, closing connection")
		c.Close()
		return
	}

	// Проверка Origin для защиты от CSRF
	origin := c.Headers("Origin")
	if origin != "" && h.config.Environment == "production" {
		allowedOrigins := []string{
			h.config.FrontendURL,
			"https://svetu.rs",
			"https://www.svetu.rs",
		}

		validOrigin := false
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				validOrigin = true
				break
			}
		}

		if !validOrigin {
			logger.Warn().Str("origin", origin).Int("userId", userID).Msg("SECURITY: WebSocket invalid origin, closing connection")
			c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseUnsupportedData, "Invalid origin"))
			c.Close()
			return
		}
	}

	logger.Info().Int("userId", userID).Str("origin", origin).Msg("WebSocket: User connected")

	// Вызываем основной обработчик
	h.handleWebSocketConnection(c, userID)
}

// WebSocket хендлер (для обратной совместимости)
func (h *ChatHandler) HandleWebSocket(c *websocket.Conn) {
	if c == nil {
		return // Защита от nil pointer
	}

	// Получаем userID, переданный из middleware
	userIDRaw := c.Locals("user_id")
	if userIDRaw == nil {
		logger.Warn().Msg("WebSocket: No user_id found, closing connection")
		c.Close()
		return
	}

	userID, ok := userIDRaw.(int)
	if !ok || userID == 0 {
		logger.Warn().Interface("userIdRaw", userIDRaw).Msg("WebSocket: Invalid user_id, closing connection")
		c.Close()
		return
	}

	logger.Info().Int("userId", userID).Msg("WebSocket: User connected")

	// Вызываем основной обработчик
	h.handleWebSocketConnection(c, userID)
}

// handleWebSocketConnection - основная логика WebSocket соединения
func (h *ChatHandler) handleWebSocketConnection(c *websocket.Conn, userID int) {
	// Создаем mutex для синхронизации записи в WebSocket
	var writeMu sync.Mutex

	// Функция для безопасной записи в WebSocket
	writeMessage := func(messageType int, data []byte) error {
		writeMu.Lock()
		defer writeMu.Unlock()
		return c.WriteMessage(messageType, data)
	}

	// Подписываемся на сообщения
	msgChan := h.services.Chat().SubscribeToMessages(userID)
	defer h.services.Chat().UnsubscribeFromMessages(userID)

	// Подписываемся на обновления статуса
	statusChan := h.services.Chat().SubscribeToStatusUpdates(userID)
	defer h.services.Chat().UnsubscribeFromStatusUpdates(userID)

	// Устанавливаем пользователя онлайн
	h.services.Chat().SetUserOnline(userID)
	defer h.services.Chat().SetUserOffline(userID)

	// Создаем контекст, который будет отменен при закрытии соединения
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Горутина для чтения сообщений от клиента
	go func() {
		defer cancel()
		for {
			messageType, message, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					logger.Error().Err(err).Msg("error reading message")
				}
				return
			}

			if messageType == websocket.TextMessage {
				// Пытаемся распарсить сообщение как JSON
				var rawMsg map[string]interface{}
				if err := json.Unmarshal(message, &rawMsg); err != nil {
					logger.Error().Err(err).Msg("Error unmarshaling message")
					continue
				}

				// Проверяем тип сообщения
				msgType, ok := rawMsg["type"].(string)
				if ok {
					switch msgType {
					case "ping":
						// Отвечаем на ping pong-сообщением
						pongMsg := map[string]interface{}{
							"type":      "pong",
							"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
						}
						if pongBytes, err := json.Marshal(pongMsg); err == nil {
							writeMessage(websocket.TextMessage, pongBytes)
						}
						continue

					case "get_online_users":
						// Отправляем список онлайн пользователей
						onlineUsers := h.services.Chat().GetOnlineUsers()
						response := map[string]interface{}{
							"type": "online_users_list",
							"payload": map[string]interface{}{
								"users": onlineUsers,
							},
						}
						if respBytes, err := json.Marshal(response); err == nil {
							writeMessage(websocket.TextMessage, respBytes)
						}
						continue

					case "user_online":
						// Обновляем статус пользователя (уже установлен при подключении)
						h.services.Chat().SetUserOnline(userID)
						continue

					case "heartbeat":
						// Обновляем статус онлайн при получении heartbeat
						h.services.Chat().SetUserOnline(userID)
						continue
					}
				}

				// Обрабатываем обычное сообщение
				var msg models.MarketplaceMessage
				if err := json.Unmarshal(message, &msg); err != nil {
					logger.Error().Err(err).Msg("Error unmarshaling WebSocket message")
					continue
				}

				// Валидация входных данных
				if msg.ReceiverID == 0 {
					logger.Error().Msg("Error: ReceiverID is 0 in WebSocket message")
					errMsg := fiber.Map{
						"error": "ReceiverID is required",
					}
					if errBytes, err := json.Marshal(errMsg); err == nil {
						writeMessage(websocket.TextMessage, errBytes)
					}
					continue
				}

				msg.SenderID = userID
				if err := h.services.Chat().SendMessage(ctx, &msg); err != nil {
					logger.Error().Err(err).Msg("Error sending message via WebSocket")
					errMsg := fiber.Map{
						"error":      err.Error(),
						"chat_id":    msg.ChatID,
						"listing_id": msg.ListingID,
					}
					if errBytes, err := json.Marshal(errMsg); err == nil {
						writeMessage(websocket.TextMessage, errBytes)
					}
				}
			}
		}
	}()

	// Основной цикл для отправки сообщений клиенту
	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				return
			}
			// Отправляем только сообщения, относящиеся к этому пользователю
			if msg.SenderID == userID || msg.ReceiverID == userID {
				// Загружаем вложения для сообщения, если они есть
				if msg.HasAttachments {
					attachments, err := h.services.Storage().GetMessageAttachments(ctx, msg.ID)
					if err == nil && len(attachments) > 0 {
						// Конвертируем []*ChatAttachment в []ChatAttachment
						msg.Attachments = make([]models.ChatAttachment, len(attachments))
						for j, att := range attachments {
							msg.Attachments[j] = *att
						}
					}
				}

				// Оборачиваем сообщение в формат с типом
				wrappedMsg := map[string]interface{}{
					"type":    "new_message",
					"payload": msg,
				}
				if data, err := json.Marshal(wrappedMsg); err == nil {
					if err := writeMessage(websocket.TextMessage, data); err != nil {
						return
					}
				}
			}
		case statusMsg, ok := <-statusChan:
			if !ok {
				return
			}
			// Отправляем обновление статуса
			if data, err := json.Marshal(statusMsg); err == nil {
				if err := writeMessage(websocket.TextMessage, data); err != nil {
					return
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

// UploadAttachments загружает вложения для сообщения
