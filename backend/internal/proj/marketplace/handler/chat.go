// backend/internal/proj/marketplace/handler/chat.go

package handler

import (
	"backend/internal/config"
	"backend/internal/domain/models"
	globalService "backend/internal/proj/global/service"
	"backend/pkg/utils"
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"log"
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

// REST эндпоинты
func (h *ChatHandler) GetChats(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	log.Printf("GetChats called for userID: %d", userID)

	chats, err := h.services.Chat().GetChats(c.Context(), userID)
	if err != nil {
		log.Printf("Error in GetChats for userID %d: %v", userID, err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error fetching chats")
	}

	log.Printf("GetChats successful for userID %d, found %d chats", userID, len(chats))
	return utils.SuccessResponse(c, chats)
}

func (h *ChatHandler) GetMessages(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	chatID := c.Query("chat_id")
	listingID := c.Query("listing_id")

	// Если есть chat_id, используем его
	if chatID != "" {
		chatIDInt, err := strconv.Atoi(chatID)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid chat ID format")
		}
		c.Context().SetUserValue("chat_id", chatIDInt)

		// Если указан chat_id, listing_id не обязателен
		// Получим listing_id из чата
		chat, err := h.services.Chat().GetChat(c.Context(), chatIDInt, userID)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error fetching chat")
		}
		listingID = strconv.Itoa(chat.ListingID)
	}

	// Получаем receiver_id для прямых сообщений
	receiverID := c.Query("receiver_id")

	// Если нет ни chat_id, ни listing_id, ни receiver_id - ошибка
	if listingID == "" && chatID == "" && receiverID == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Either chat_id, listing_id or receiver_id is required")
	}

	listingIDInt := 0
	if listingID != "" {
		var err error
		listingIDInt, err = strconv.Atoi(listingID)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid listing ID format")
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

	log.Printf("GetMessages: page=%d, limit=%d, offset=%d, chatID=%s, listingID=%d, userID=%d",
		page, limit, offset, c.Query("chat_id"), listingIDInt, userID)

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
		log.Printf("Error fetching messages: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error fetching messages")
	}

	log.Printf("GetMessages: returned %d messages", len(messages))

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

func (h *ChatHandler) GetUnreadCount(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	count, err := h.services.Chat().GetUnreadMessagesCount(c.Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error getting unread count")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"count": count,
	})
}

func (h *ChatHandler) SendMessage(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	var req models.CreateMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Валидация входных данных
	if req.ReceiverID == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Missing receiver ID")
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
		log.Printf("Error sending message: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error sending message")
	}

	return utils.SuccessResponse(c, msg)
}

func (h *ChatHandler) MarkAsRead(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	var req models.MarkAsReadRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.services.Chat().MarkMessagesAsRead(c.Context(), req.MessageIDs, userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error marking messages as read")
	}

	return utils.SuccessResponse(c, fiber.Map{"message": "Messages marked as read"})
}
func (h *ChatHandler) ArchiveChat(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	chatID, err := c.ParamsInt("chat_id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный ID чата")
	}

	err = h.services.Chat().ArchiveChat(c.Context(), chatID, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка при архивировании чата")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Чат архивирован",
	})
}

// UploadAttachments загружает вложения для сообщения
func (h *ChatHandler) UploadAttachments(c *fiber.Ctx) error {
	log.Printf("UploadAttachments called")
	userID := c.Locals("user_id").(int)
	messageID, err := c.ParamsInt("id")
	if err != nil {
		log.Printf("Error parsing message ID: %v", err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid message ID")
	}
	log.Printf("UploadAttachments: userID=%d, messageID=%d", userID, messageID)

	// Получаем сообщение для проверки прав
	message, err := h.services.Storage().GetMessageByID(c.Context(), messageID)
	if err != nil {
		log.Printf("Error getting message by ID %d: %v", messageID, err)
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Message not found")
	}

	// Проверяем, что пользователь является отправителем сообщения
	if message.SenderID != userID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "You can only upload attachments to your own messages")
	}

	// Получаем файлы из запроса
	form, err := c.MultipartForm()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Error parsing multipart form")
	}

	files := form.File["files"]
	if len(files) == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "No files provided")
	}

	// Ограничение на количество файлов
	if len(files) > 10 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Maximum 10 files allowed per message")
	}

	// Загружаем файлы через сервис
	attachments, err := h.services.ChatAttachment().UploadAttachments(c.Context(), messageID, files)
	if err != nil {
		log.Printf("Error uploading attachments: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	// Отправляем обновленное сообщение через WebSocket
	if len(attachments) > 0 {
		log.Printf("UploadAttachments: uploading %d attachments for message %d", len(attachments), messageID)
		// Получаем обновленное сообщение с вложениями
		updatedMessage, err := h.services.Storage().GetMessageByID(c.Context(), messageID)
		if err == nil {
			log.Printf("UploadAttachments: got message from DB, senderID=%d, receiverID=%d", updatedMessage.SenderID, updatedMessage.ReceiverID)
			// Конвертируем вложения для сообщения
			updatedMessage.Attachments = make([]models.ChatAttachment, len(attachments))
			for i, att := range attachments {
				updatedMessage.Attachments[i] = *att
			}
			updatedMessage.HasAttachments = true
			updatedMessage.AttachmentsCount = len(attachments)

			log.Printf("UploadAttachments: broadcasting message with %d attachments", len(updatedMessage.Attachments))
			// Отправляем обновленное сообщение через WebSocket
			h.services.Chat().BroadcastMessage(updatedMessage)
		} else {
			log.Printf("UploadAttachments: error getting message by ID: %v", err)
		}
	}

	return utils.SuccessResponse(c, attachments)
}

// GetAttachment получает информацию о вложении
func (h *ChatHandler) GetAttachment(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	attachmentID, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid attachment ID")
	}

	// Получаем вложение
	attachment, err := h.services.ChatAttachment().GetAttachment(c.Context(), attachmentID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Attachment not found")
	}

	// Проверяем доступ к вложению через сообщение
	message, err := h.services.Storage().GetMessageByID(c.Context(), attachment.MessageID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Message not found")
	}

	// Пользователь должен быть участником чата
	if message.SenderID != userID && message.ReceiverID != userID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Access denied")
	}

	return utils.SuccessResponse(c, attachment)
}

// DeleteAttachment удаляет вложение
func (h *ChatHandler) DeleteAttachment(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	attachmentID, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid attachment ID")
	}

	if err := h.services.ChatAttachment().DeleteAttachment(c.Context(), attachmentID, userID); err != nil {
		if err.Error() == "permission denied" {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "You can only delete your own attachments")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error deleting attachment")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Attachment deleted successfully",
	})
}

// HandleWebSocketWithAuth - WebSocket хендлер с переданным userID
func (h *ChatHandler) HandleWebSocketWithAuth(c *websocket.Conn, userID int) {
	if c == nil {
		return // Защита от nil pointer
	}

	// Проверяем, что userID валидный
	if userID == 0 {
		log.Printf("WebSocket: Invalid user_id: %d, closing connection", userID)
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
			log.Printf("SECURITY: WebSocket invalid origin %s for user %d, closing connection", origin, userID)
			c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseUnsupportedData, "Invalid origin"))
			c.Close()
			return
		}
	}

	log.Printf("WebSocket: User %d connected from origin %s", userID, origin)

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
		log.Printf("WebSocket: No user_id found, closing connection")
		c.Close()
		return
	}

	userID, ok := userIDRaw.(int)
	if !ok || userID == 0 {
		log.Printf("WebSocket: Invalid user_id: %v, closing connection", userIDRaw)
		c.Close()
		return
	}

	log.Printf("WebSocket: User %d connected", userID)

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
					log.Printf("error reading message: %v", err)
				}
				return
			}

			if messageType == websocket.TextMessage {
				// Пытаемся распарсить сообщение как JSON
				var rawMsg map[string]interface{}
				if err := json.Unmarshal(message, &rawMsg); err != nil {
					log.Printf("Error unmarshaling message: %v", err)
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
					log.Printf("Error unmarshaling WebSocket message: %v", err)
					continue
				}

				// Валидация входных данных
				if msg.ReceiverID == 0 {
					log.Printf("Error: ReceiverID is 0 in WebSocket message")
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
					log.Printf("Error sending message via WebSocket: %v", err)
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
