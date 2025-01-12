package handler

import (
    "backend/internal/domain/models"
    globalService "backend/internal/proj/global/service"
    "backend/pkg/utils"
    "context"
	"log"
    "encoding/json"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/websocket/v2"
)

type ChatHandler struct {
    services globalService.ServicesInterface
}

func NewChatHandler(services globalService.ServicesInterface) *ChatHandler {
    return &ChatHandler{
        services: services,
    }
}

// REST эндпоинты
func (h *ChatHandler) GetChats(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    
    chats, err := h.services.Chat().GetChats(c.Context(), userID)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error fetching chats")
    }
    
    return utils.SuccessResponse(c, chats)
}

func (h *ChatHandler) GetMessages(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    listingID, err := c.ParamsInt("listing_id")
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid listing ID")
    }
    
    page := utils.StringToInt(c.Query("page"), 1)
    limit := utils.StringToInt(c.Query("limit"), 20)
    
    messages, err := h.services.Chat().GetMessages(c.Context(), listingID, userID, page, limit)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error fetching messages")
    }
    
    return utils.SuccessResponse(c, messages)
}

func (h *ChatHandler) SendMessage(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    
    var req models.CreateMessageRequest
    if err := c.BodyParser(&req); err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
    }
    
    msg := &models.MarketplaceMessage{
        ListingID:  req.ListingID,
        SenderID:   userID,
        ReceiverID: req.ReceiverID,
        Content:    req.Content,
    }
    
    if err := h.services.Chat().SendMessage(c.Context(), msg); err != nil {
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
// WebSocket хендлер
func (h *ChatHandler) HandleWebSocket(c *websocket.Conn) {
    if c == nil {
        return // Защита от nil pointer
    }

    userID := c.Locals("user_id").(int)
    
    // Подписываемся на сообщения
    msgChan := h.services.Chat().SubscribeToMessages(userID)
    defer h.services.Chat().UnsubscribeFromMessages(userID)

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
                var msg models.MarketplaceMessage
                if err := json.Unmarshal(message, &msg); err != nil {
                    continue
                }
                
                msg.SenderID = userID
                if err := h.services.Chat().SendMessage(ctx, &msg); err != nil {
                    errMsg := fiber.Map{"error": err.Error()}
                    if errBytes, err := json.Marshal(errMsg); err == nil {
                        c.WriteMessage(websocket.TextMessage, errBytes)
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
                if data, err := json.Marshal(msg); err == nil {
                    if err := c.WriteMessage(websocket.TextMessage, data); err != nil {
                        return
                    }
                }
            }
        case <-ctx.Done():
            return
        }
    }
}