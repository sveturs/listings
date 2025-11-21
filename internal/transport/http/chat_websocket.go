package http

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog"
	authservice "github.com/sveturs/auth/pkg/service"

	"github.com/sveturs/listings/internal/domain"
	ws "github.com/sveturs/listings/internal/websocket"
)

const (
	// MaxMessageSize defines maximum message size in bytes (512KB)
	MaxMessageSize = 512 * 1024
)

// ErrorMessage represents an error response to client
type ErrorMessage struct {
	Type      string `json:"type"`
	Error     string `json:"error"`
	Timestamp string `json:"timestamp"`
}

// ChatRepository interface for getting chat info
type ChatRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Chat, error)
}

// ChatWebSocketHandler handles WebSocket connections for chat
type ChatWebSocketHandler struct {
	hub              *ws.ChatHub
	authService      *authservice.AuthService
	chatRepo         ChatRepository
	securityMW       *ws.SecurityMiddleware
	logger           zerolog.Logger
	startTime        time.Time
}

// NewChatWebSocketHandler creates a new chat WebSocket handler
func NewChatWebSocketHandler(hub *ws.ChatHub, authService *authservice.AuthService, chatRepo ChatRepository, logger zerolog.Logger) *ChatWebSocketHandler {
	// Allowed origins from environment or config
	// For now, allow common development and production origins
	allowedOrigins := []string{
		"http://localhost:3001",
		"https://dev.svetu.rs",
		"https://svetu.rs",
	}

	return &ChatWebSocketHandler{
		hub:         hub,
		authService: authService,
		chatRepo:    chatRepo,
		securityMW:  ws.NewSecurityMiddleware(60, allowedOrigins, logger), // 60 messages per minute
		logger:      logger.With().Str("component", "chat_websocket_handler").Logger(),
		startTime:   time.Now(),
	}
}

// HandleChatWebSocket handles WebSocket upgrade and connection
func (h *ChatWebSocketHandler) HandleChatWebSocket(c *fiber.Ctx) error {
	// Check if request is WebSocket upgrade
	if !websocket.IsWebSocketUpgrade(c) {
		return fiber.ErrUpgradeRequired
	}

	// Check origin (graceful degradation - log but don't block)
	origin := c.Get("Origin")
	if origin != "" {
		if err := h.securityMW.CheckOrigin(origin); err != nil {
			h.logger.Warn().
				Str("origin", origin).
				Str("remote_ip", c.IP()).
				Msg("WebSocket connection from unallowed origin (allowed for now)")
			// Don't block for now - just log
		}
	}

	// Get JWT token from query parameter
	token := c.Query("token")
	if token == "" {
		h.logger.Warn().Msg("WebSocket connection attempt without token")
		h.securityMW.TrackError()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "token required in query parameter",
		})
	}

	// Validate token and get user ID
	claims, err := h.authService.ValidateToken(c.Context(), token)
	if err != nil {
		h.logger.Warn().
			Err(err).
			Str("remote_ip", c.IP()).
			Msg("invalid token for WebSocket connection")
		h.securityMW.TrackError()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid or expired token",
		})
	}

	userID := int64(claims.UserID)
	if userID == 0 {
		h.logger.Warn().
			Str("remote_ip", c.IP()).
			Msg("token validation succeeded but user_id is 0")
		h.securityMW.TrackError()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid user_id in token",
		})
	}

	h.logger.Info().
		Int64("user_id", userID).
		Str("remote_ip", c.IP()).
		Str("origin", origin).
		Msg("WebSocket upgrade authorized")

	// Upgrade to WebSocket
	return websocket.New(func(conn *websocket.Conn) {
		h.handleConnection(conn, userID)
	})(c)
}

// handleConnection handles a WebSocket connection lifecycle
func (h *ChatWebSocketHandler) handleConnection(conn *websocket.Conn, userID int64) {
	// Track connection
	connStart := time.Now()
	ctx := context.Background()
	defer h.securityMW.TrackConnection(ctx)()

	// Register connection
	h.hub.RegisterConnection(conn, userID)
	defer func() {
		duration := time.Since(connStart)
		h.hub.UnregisterConnection(conn, userID)
		h.logger.Info().
			Int64("user_id", userID).
			Dur("duration", duration).
			Msg("WebSocket connection closed")
	}()

	h.logger.Info().
		Int64("user_id", userID).
		Str("remote_addr", conn.RemoteAddr().String()).
		Msg("WebSocket connection established")

	// Send initial connected message (using SafeWriteMessage to avoid concurrent write panic)
	connectedMsg := &ws.BroadcastMessage{
		Type:      "connected",
		UserID:    &userID,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	data, _ := json.Marshal(connectedMsg)
	if err := h.hub.SafeWriteMessage(conn, websocket.TextMessage, data, userID); err != nil {
		h.logger.Error().
			Err(err).
			Int64("user_id", userID).
			Msg("Failed to send connected message")
		h.securityMW.TrackError()
		return
	}
	h.securityMW.TrackMessage(false) // Sent message
	h.logger.Debug().Int64("user_id", userID).Msg("Sent connected message successfully")

	// Send list of currently online users to the new client
	h.hub.SendOnlineUsersList(conn, userID)

	// Send last seen data for offline users
	h.hub.SendUsersLastSeen(conn, userID)

	// Set up pong handler
	conn.SetPongHandler(func(string) error {
		// Reset read deadline on pong
		return conn.SetReadDeadline(time.Now().Add(90 * time.Second))
	})

	// Set initial read deadline
	_ = conn.SetReadDeadline(time.Now().Add(90 * time.Second))

	// Start ping ticker (reduced to 20 seconds for better connection health)
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	// Channel to signal done
	done := make(chan struct{})

	// Goroutine for sending pings (using SafeWriteMessage to avoid concurrent write panic)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := h.hub.SafeWriteMessage(conn, websocket.PingMessage, []byte{}, userID); err != nil {
					h.logger.Debug().Err(err).Int64("user_id", userID).Msg("failed to send ping")
					close(done)
					return
				}
			case <-done:
				return
			}
		}
	}()

	// Main message loop
	h.logger.Debug().Int64("user_id", userID).Msg("Entering main read loop, waiting for messages...")
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				h.logger.Warn().
					Err(err).
					Int64("user_id", userID).
					Msg("WebSocket unexpected close error")
				h.securityMW.TrackError()
			} else {
				h.logger.Debug().
					Err(err).
					Int64("user_id", userID).
					Msg("WebSocket read error")
			}
			close(done)
			break
		}

		h.securityMW.TrackMessage(true) // Received message

		// Validate message size
		if len(message) > MaxMessageSize {
			h.logger.Warn().
				Int("message_size", len(message)).
				Int("max_size", MaxMessageSize).
				Int64("user_id", userID).
				Msg("Message size exceeds limit")
			h.securityMW.TrackError()
			// Send error to client (using SafeWriteMessage to avoid concurrent write panic)
			errMsg := &ErrorMessage{
				Type:      "error",
				Error:     fmt.Sprintf("Message too large (max %d bytes)", MaxMessageSize),
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			}
			errData, _ := json.Marshal(errMsg)
			_ = h.hub.SafeWriteMessage(conn, websocket.TextMessage, errData, userID)
			continue
		}

		h.logger.Debug().
			Int("messageType", messageType).
			Int("messageSize", len(message)).
			Int64("user_id", userID).
			Msg("Received message from client")

		// Handle close message
		if messageType == websocket.CloseMessage {
			h.logger.Info().Int64("user_id", userID).Msg("WebSocket close message received")
			close(done)
			break
		}

		// Handle text messages
		if messageType == websocket.TextMessage {
			h.handleMessage(conn, userID, message)
		}

		// Reset read deadline after each message
		_ = conn.SetReadDeadline(time.Now().Add(90 * time.Second))
	}
}

// handleMessage processes incoming text messages from client
func (h *ChatWebSocketHandler) handleMessage(conn *websocket.Conn, userID int64, message []byte) {
	// Check rate limit for non-ping messages
	var msg ws.WSMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		h.logger.Warn().
			Err(err).
			Int64("user_id", userID).
			Msg("failed to parse WebSocket message")
		h.securityMW.TrackError()
		return
	}

	// Apply rate limiting only for significant actions (not pings)
	if msg.Type != "ping" && msg.Type != "pong" {
		if err := h.securityMW.CheckRateLimit(userID); err != nil {
			h.logger.Warn().
				Int64("user_id", userID).
				Str("message_type", msg.Type).
				Msg("Rate limit exceeded for WebSocket message")
			h.securityMW.TrackError()

			// Send rate limit error to client (using SafeWriteMessage to avoid concurrent write panic)
			errMsg := &ErrorMessage{
				Type:      "error",
				Error:     "Rate limit exceeded. Please slow down.",
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			}
			errData, _ := json.Marshal(errMsg)
			_ = h.hub.SafeWriteMessage(conn, websocket.TextMessage, errData, userID)
			return
		}
	}

	switch msg.Type {
	case "ping":
		// Reset read deadline on text ping (frontend heartbeat)
		_ = conn.SetReadDeadline(time.Now().Add(90 * time.Second))

		// Respond with pong (using SafeWriteMessage to avoid concurrent write panic)
		pongMsg := &ws.BroadcastMessage{
			Type:      "pong",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
		data, _ := json.Marshal(pongMsg)
		if err := h.hub.SafeWriteMessage(conn, websocket.TextMessage, data, userID); err != nil {
			h.logger.Warn().
				Err(err).
				Int64("user_id", userID).
				Msg("failed to send pong")
			h.securityMW.TrackError()
			return
		}
		h.securityMW.TrackMessage(false) // Sent pong

	case "typing":
		// Broadcast typing indicator - support both direct and payload formats
		var chatID *int64
		var isTyping *bool

		// Check payload format first (from frontend sendTypingIndicator)
		if msg.Payload != nil {
			chatID = msg.Payload.ChatID
			isTyping = msg.Payload.IsTyping
		}
		// Fallback to direct format
		if chatID == nil {
			chatID = msg.ChatID
		}
		if isTyping == nil {
			isTyping = msg.IsTyping
		}

		if chatID != nil && isTyping != nil {
			// Get chat to find the other participant
			if h.chatRepo != nil {
				chat, err := h.chatRepo.GetByID(context.Background(), *chatID)
				if err == nil && chat != nil {
					// Send to the other participant
					var targetUserID int64
					if chat.BuyerID == userID {
						targetUserID = chat.SellerID
					} else {
						targetUserID = chat.BuyerID
					}
					h.hub.BroadcastTypingToUser(*chatID, userID, targetUserID, *isTyping)
				} else {
					// Fallback to old method (broadcasts to all)
					h.hub.BroadcastTyping(*chatID, userID, *isTyping)
				}
			} else {
				h.hub.BroadcastTyping(*chatID, userID, *isTyping)
			}
			h.logger.Debug().
				Int64("user_id", userID).
				Int64("chat_id", *chatID).
				Bool("is_typing", *isTyping).
				Msg("typing indicator processed")
		}

	case "logout":
		// Explicit logout - immediately mark user as offline
		h.logger.Info().
			Int64("user_id", userID).
			Msg("explicit logout message received, closing connection")
		// The defer in handleConnection will call UnregisterConnection which broadcasts user_offline
		// Just close the connection gracefully
		return

	case "mark_read":
		// Mark message as read
		// Note: This should trigger backend logic, not just broadcast
		// For now, we'll just log it - actual mark_read should go through gRPC/HTTP API
		h.logger.Debug().
			Int64("user_id", userID).
			Interface("message_id", msg.MessageID).
			Msg("mark_read received (should use API endpoint instead)")

	case "get_user_status":
		// Get status for a specific user
		var targetUserID *int64
		if msg.Payload != nil {
			targetUserID = msg.Payload.UserID
		}

		if targetUserID != nil {
			// Check if user is online
			isOnline := h.hub.IsUserOnline(*targetUserID)
			var lastSeen string

			if !isOnline {
				// Get last seen time
				if ts := h.hub.GetUserLastSeen(*targetUserID); ts != nil {
					lastSeen = ts.UTC().Format(time.RFC3339)
				}
			}

			// Determine status
			status := "offline"
			if isOnline {
				status = "online"
			}

			// Send response (using SafeWriteMessage to avoid concurrent write panic)
			response := &ws.BroadcastMessage{
				Type:      "user_status",
				UserID:    targetUserID,
				Status:    status,
				LastSeen:  lastSeen,
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			}
			data, _ := json.Marshal(response)
			_ = h.hub.SafeWriteMessage(conn, websocket.TextMessage, data, userID)

			h.logger.Debug().
				Int64("user_id", userID).
				Int64("target_user_id", *targetUserID).
				Str("status", status).
				Msg("get_user_status processed")
		}

	default:
		h.logger.Debug().
			Int64("user_id", userID).
			Str("type", msg.Type).
			Msg("unknown message type")
	}
}

// RegisterWebSocketRoute registers the WebSocket endpoint
func (h *ChatWebSocketHandler) RegisterWebSocketRoute(app *fiber.App) {
	app.Get("/ws/chat", h.HandleChatWebSocket)
	app.Get("/ws/health", h.HealthCheck)
	h.logger.Info().Msg("WebSocket routes registered: GET /ws/chat, GET /ws/health")
}

// HealthCheck returns WebSocket service health and metrics
func (h *ChatWebSocketHandler) HealthCheck(c *fiber.Ctx) error {
	stats := h.securityMW.GetStats()
	uptime := time.Since(h.startTime)

	return c.JSON(fiber.Map{
		"status": "ok",
		"service": "chat-websocket",
		"uptime_seconds": int64(uptime.Seconds()),
		"uptime_human": uptime.String(),
		"metrics": stats,
		"config": fiber.Map{
			"max_message_size": MaxMessageSize,
			"rate_limit_per_minute": 60,
		},
	})
}

// Stop gracefully stops the handler and cleanup resources
func (h *ChatWebSocketHandler) Stop() {
	h.logger.Info().Msg("Stopping WebSocket handler...")
	h.securityMW.Stop()
	h.logger.Info().Msg("WebSocket handler stopped")
}
