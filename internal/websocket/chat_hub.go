package websocket

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog"

	"github.com/sveturs/listings/internal/domain"
)

// ChatHub manages WebSocket connections for chat
type ChatHub struct {
	// Connections mapped by user_id -> set of connections
	connections map[int64]map[*websocket.Conn]bool

	// Write mutexes for each connection to prevent concurrent writes
	connWriteMu map[*websocket.Conn]*sync.Mutex

	// Channels for managing connections and broadcasts
	broadcast  chan *BroadcastMessage
	register   chan *ClientConnection
	unregister chan *ClientConnection

	// Mutex for protecting connections map
	mutex sync.RWMutex

	// Presence tracking
	userLastSeen map[int64]time.Time // userId -> last seen timestamp
	lastSeenMu   sync.RWMutex

	// Typing tracking: chatId -> userId -> started typing time
	typingUsers map[int64]map[int64]time.Time
	typingMu    sync.RWMutex

	// Logger
	logger zerolog.Logger
}

// ClientConnection represents a WebSocket client connection
type ClientConnection struct {
	conn   *websocket.Conn
	userID int64
}

// WSMessage represents incoming message from client
type WSMessage struct {
	Type      string     `json:"type"` // ping, typing, mark_read, logout
	ChatID    *int64     `json:"chat_id,omitempty"`
	MessageID *int64     `json:"message_id,omitempty"`
	IsTyping  *bool      `json:"is_typing,omitempty"`
	Payload   *WSPayload `json:"payload,omitempty"` // For nested payload format
}

// WSPayload represents nested payload in WebSocket messages
type WSPayload struct {
	ChatID   *int64 `json:"chat_id,omitempty"`
	UserID   *int64 `json:"user_id,omitempty"`
	IsTyping *bool  `json:"is_typing,omitempty"`
}

// UserLastSeenInfo represents last seen info for a user
type UserLastSeenInfo struct {
	UserID   int64  `json:"user_id"`
	LastSeen string `json:"last_seen"` // ISO timestamp
}

// BroadcastMessage represents outgoing message to clients
type BroadcastMessage struct {
	Type          string             `json:"type"` // new_message, message_read, message_delivered, user_typing, user_online, user_offline, connected, pong, users_last_seen
	ChatID        *int64             `json:"chat_id,omitempty"`
	Message       *domain.Message    `json:"message,omitempty"`
	MessageID     *int64             `json:"message_id,omitempty"`
	MessageIDs    []int64            `json:"message_ids,omitempty"` // For batch message_read events
	ReadBy        *int64             `json:"read_by,omitempty"`
	UserID        *int64             `json:"user_id,omitempty"`
	IsTyping      *bool              `json:"is_typing,omitempty"`
	Status        string             `json:"status,omitempty"`    // online, offline
	LastSeen      string             `json:"last_seen,omitempty"` // ISO timestamp for offline
	DeliveredAt   string             `json:"delivered_at,omitempty"`
	OnlineUsers   []int64            `json:"online_users,omitempty"`    // For online_users_list event
	UsersLastSeen []UserLastSeenInfo `json:"users_last_seen,omitempty"` // For users_last_seen event
	TargetUserIDs []int64            `json:"-"`                         // Internal: explicit target users (not serialized)
	Timestamp     string             `json:"timestamp"`
}

// NewChatHub creates a new WebSocket hub for chat
func NewChatHub(logger zerolog.Logger) *ChatHub {
	return &ChatHub{
		broadcast:    make(chan *BroadcastMessage, 256),
		register:     make(chan *ClientConnection),
		unregister:   make(chan *ClientConnection),
		connections:  make(map[int64]map[*websocket.Conn]bool),
		connWriteMu:  make(map[*websocket.Conn]*sync.Mutex),
		userLastSeen: make(map[int64]time.Time),
		typingUsers:  make(map[int64]map[int64]time.Time),
		logger:       logger.With().Str("component", "chat_hub").Logger(),
	}
}

// Run starts the hub event loop
func (h *ChatHub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			h.logger.Info().Msg("chat hub shutting down")
			h.closeAllConnections()
			return

		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

// RegisterConnection registers a new WebSocket connection
func (h *ChatHub) RegisterConnection(conn *websocket.Conn, userID int64) {
	// Pre-create write mutex synchronously so it's available immediately for SendOnlineUsersList
	h.mutex.Lock()
	h.connWriteMu[conn] = &sync.Mutex{}
	h.mutex.Unlock()

	h.register <- &ClientConnection{conn: conn, userID: userID}
}

// UnregisterConnection unregisters a WebSocket connection
func (h *ChatHub) UnregisterConnection(conn *websocket.Conn, userID int64) {
	h.unregister <- &ClientConnection{conn: conn, userID: userID}
}

// BroadcastNewMessage broadcasts a new message to all participants in the chat
func (h *ChatHub) BroadcastNewMessage(chatID int64, message *domain.Message) {
	h.broadcast <- &BroadcastMessage{
		Type:      "new_message",
		ChatID:    &chatID,
		Message:   message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// BroadcastMessageRead broadcasts a message read event
func (h *ChatHub) BroadcastMessageRead(chatID, messageID, userID int64) {
	h.broadcast <- &BroadcastMessage{
		Type:      "message_read",
		ChatID:    &chatID,
		MessageID: &messageID,
		ReadBy:    &userID,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// BroadcastTyping broadcasts a typing indicator (deprecated - use BroadcastTypingToUser)
func (h *ChatHub) BroadcastTyping(chatID, userID int64, isTyping bool) {
	h.broadcast <- &BroadcastMessage{
		Type:      "user_typing", // Must match frontend expectation
		ChatID:    &chatID,
		UserID:    &userID,
		IsTyping:  &isTyping,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// BroadcastTypingToUser broadcasts a typing indicator to a specific user
func (h *ChatHub) BroadcastTypingToUser(chatID, typerID, targetUserID int64, isTyping bool) {
	h.broadcast <- &BroadcastMessage{
		Type:          "user_typing",
		ChatID:        &chatID,
		UserID:        &typerID,
		IsTyping:      &isTyping,
		TargetUserIDs: []int64{targetUserID},
		Timestamp:     time.Now().UTC().Format(time.RFC3339),
	}
}

// registerClient handles client registration
func (h *ChatHub) registerClient(client *ClientConnection) {
	h.mutex.Lock()

	// Check if this is the first connection for this user
	isFirstConnection := h.connections[client.userID] == nil || len(h.connections[client.userID]) == 0

	if h.connections[client.userID] == nil {
		h.connections[client.userID] = make(map[*websocket.Conn]bool)
	}
	h.connections[client.userID][client.conn] = true

	// Note: write mutex is pre-created in RegisterConnection for immediate availability

	h.logger.Info().
		Int64("user_id", client.userID).
		Int("total_user_connections", len(h.connections[client.userID])).
		Bool("is_first_connection", isFirstConnection).
		Msg("WebSocket client connected")

	h.mutex.Unlock()

	// Broadcast user online status if this is their first connection
	if isFirstConnection {
		h.BroadcastUserOnline(client.userID)
	}
}

// unregisterClient handles client disconnection
func (h *ChatHub) unregisterClient(client *ClientConnection) {
	h.mutex.Lock()

	isLastConnection := false

	if connections, ok := h.connections[client.userID]; ok {
		if _, ok := connections[client.conn]; ok {
			delete(connections, client.conn)
			// Remove write mutex for this connection
			delete(h.connWriteMu, client.conn)
			_ = client.conn.Close()

			// If no more connections for this user, remove the map
			if len(connections) == 0 {
				delete(h.connections, client.userID)
				isLastConnection = true
			}
		}
	}

	h.logger.Info().
		Int64("user_id", client.userID).
		Bool("is_last_connection", isLastConnection).
		Msg("WebSocket client disconnected")

	h.mutex.Unlock()

	// Broadcast user offline status if this was their last connection
	if isLastConnection {
		h.BroadcastUserOffline(client.userID)
	}
}

// broadcastMessage sends a message to relevant users
func (h *ChatHub) broadcastMessage(msg *BroadcastMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to marshal broadcast message")
		return
	}

	// Determine target users based on message type
	var targetUserIDs []int64
	broadcastToAll := false

	switch msg.Type {
	case "new_message":
		// Send to both sender and receiver
		if msg.Message != nil {
			targetUserIDs = []int64{msg.Message.SenderID, msg.Message.ReceiverID}
		}
	case "message_read":
		// Send to both users in the chat
		// We need to get chat participants - for now, broadcast to all
		// In production, you might want to pass participant IDs
		if msg.ReadBy != nil {
			// Send to the reader and potentially the sender
			// Since we don't have sender ID here, we'll broadcast to the reader
			targetUserIDs = []int64{*msg.ReadBy}
		}
	case "message_delivered":
		// Send to sender that message was delivered
		// For now, broadcast to all in case we don't know the sender
		broadcastToAll = true
	case "typing", "user_typing":
		// Use TargetUserIDs if provided, otherwise broadcast to all
		if len(msg.TargetUserIDs) > 0 {
			targetUserIDs = msg.TargetUserIDs
		} else {
			broadcastToAll = true
		}
	case "user_online", "user_offline":
		// Broadcast presence to all connected users
		broadcastToAll = true
	}

	// Send to target users
	// Collect connections to send to, then release read lock before writing
	type connTarget struct {
		conn   *websocket.Conn
		userID int64
	}
	var targets []connTarget

	h.mutex.RLock()
	if broadcastToAll {
		// Send to all connected users
		for userID, connections := range h.connections {
			for conn := range connections {
				targets = append(targets, connTarget{conn: conn, userID: userID})
			}
		}
	} else {
		// Send to specific users
		for _, userID := range targetUserIDs {
			if connections, ok := h.connections[userID]; ok {
				for conn := range connections {
					targets = append(targets, connTarget{conn: conn, userID: userID})
				}
			}
		}
	}
	totalConnectedUsers := len(h.connections)
	h.mutex.RUnlock()

	h.logger.Debug().
		Str("type", msg.Type).
		Int("target_users", len(targetUserIDs)).
		Int("target_connections", len(targets)).
		Int("total_connected_users", totalConnectedUsers).
		Msg("broadcasting message")

	// Now send to all targets without holding the read lock
	// This prevents race condition when UnregisterConnection tries to acquire write lock
	var failedConnections []connTarget
	for _, target := range targets {
		if err := h.sendToConnection(target.conn, data, target.userID); err != nil {
			h.logger.Warn().
				Err(err).
				Int64("user_id", target.userID).
				Msg("failed to send message to connection")
			failedConnections = append(failedConnections, target)
		} else {
			h.logger.Debug().
				Int64("user_id", target.userID).
				Msg("message sent successfully to user")
		}
	}

	// Unregister failed connections
	for _, target := range failedConnections {
		h.UnregisterConnection(target.conn, target.userID)
	}
}

// sendToConnection sends data to a specific connection with timeout
func (h *ChatHub) sendToConnection(conn *websocket.Conn, data []byte, userID int64) error {
	// Get write mutex for this connection
	h.mutex.RLock()
	mu, ok := h.connWriteMu[conn]
	h.mutex.RUnlock()

	if !ok {
		h.logger.Warn().Int64("user_id", userID).Msg("no write mutex found for connection")
		return nil // Connection might be already unregistered
	}

	// Lock the connection for writing
	mu.Lock()
	defer mu.Unlock()

	// Set write deadline - 10 seconds to handle slow networks
	deadline := time.Now().Add(10 * time.Second)
	if err := conn.SetWriteDeadline(deadline); err != nil {
		h.logger.Warn().Err(err).Int64("user_id", userID).Msg("failed to set write deadline")
		return err
	}

	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		h.logger.Warn().Err(err).Int64("user_id", userID).Msg("failed to write message")
		return err
	}

	return nil
}

// closeAllConnections closes all active connections
func (h *ChatHub) closeAllConnections() {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for userID, connections := range h.connections {
		for conn := range connections {
			_ = conn.Close()
		}
		delete(h.connections, userID)
	}

	h.logger.Info().Msg("all connections closed")
}

// GetActiveConnections returns statistics about active connections
func (h *ChatHub) GetActiveConnections() map[string]interface{} {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	totalConnections := 0
	userCount := len(h.connections)

	for _, connections := range h.connections {
		totalConnections += len(connections)
	}

	return map[string]interface{}{
		"total_connections": totalConnections,
		"unique_users":      userCount,
	}
}

// StartTypingCleanup starts a background goroutine to clean up stale typing indicators
func (h *ChatHub) StartTypingCleanup(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			h.logger.Info().Msg("typing cleanup stopped")
			return
		case <-ticker.C:
			h.cleanupStaleTyping()
		}
	}
}

// cleanupStaleTyping removes typing indicators older than 5 seconds
func (h *ChatHub) cleanupStaleTyping() {
	h.typingMu.Lock()
	defer h.typingMu.Unlock()

	now := time.Now()
	staleThreshold := 5 * time.Second

	for chatID, users := range h.typingUsers {
		for userID, startedAt := range users {
			if now.Sub(startedAt) > staleThreshold {
				delete(users, userID)

				// Broadcast typing stopped
				go h.BroadcastTyping(chatID, userID, false)

				h.logger.Debug().
					Int64("chat_id", chatID).
					Int64("user_id", userID).
					Msg("cleaned up stale typing indicator")
			}
		}

		if len(users) == 0 {
			delete(h.typingUsers, chatID)
		}
	}
}

// SetUserTyping updates typing status for a user in a chat
func (h *ChatHub) SetUserTyping(chatID, userID int64, isTyping bool) {
	h.typingMu.Lock()
	defer h.typingMu.Unlock()

	if isTyping {
		if h.typingUsers[chatID] == nil {
			h.typingUsers[chatID] = make(map[int64]time.Time)
		}
		h.typingUsers[chatID][userID] = time.Now()
	} else {
		if users, ok := h.typingUsers[chatID]; ok {
			delete(users, userID)
			if len(users) == 0 {
				delete(h.typingUsers, chatID)
			}
		}
	}

	// Broadcast to other participants
	h.BroadcastTyping(chatID, userID, isTyping)
}

// BroadcastUserOnline broadcasts user online status to all connected users
func (h *ChatHub) BroadcastUserOnline(userID int64) {
	h.lastSeenMu.Lock()
	delete(h.userLastSeen, userID) // Remove last seen when user comes online
	h.lastSeenMu.Unlock()

	h.broadcast <- &BroadcastMessage{
		Type:      "user_online",
		UserID:    &userID,
		Status:    "online",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	h.logger.Debug().Int64("user_id", userID).Msg("broadcasting user online")
}

// BroadcastUserOffline broadcasts user offline status with last seen
func (h *ChatHub) BroadcastUserOffline(userID int64) {
	now := time.Now().UTC()

	h.lastSeenMu.Lock()
	h.userLastSeen[userID] = now
	h.lastSeenMu.Unlock()

	h.broadcast <- &BroadcastMessage{
		Type:      "user_offline",
		UserID:    &userID,
		Status:    "offline",
		LastSeen:  now.Format(time.RFC3339),
		Timestamp: now.Format(time.RFC3339),
	}

	h.logger.Debug().Int64("user_id", userID).Msg("broadcasting user offline")
}

// BroadcastMessageDelivered broadcasts message delivered status
func (h *ChatHub) BroadcastMessageDelivered(chatID, messageID int64, deliveredAt time.Time) {
	h.broadcast <- &BroadcastMessage{
		Type:        "message_delivered",
		ChatID:      &chatID,
		MessageID:   &messageID,
		DeliveredAt: deliveredAt.UTC().Format(time.RFC3339),
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
}

// BroadcastMessagesRead broadcasts batch message read status
func (h *ChatHub) BroadcastMessagesRead(chatID int64, messageIDs []int64, readerID int64, readAt time.Time) {
	h.broadcast <- &BroadcastMessage{
		Type:       "message_read",
		ChatID:     &chatID,
		MessageIDs: messageIDs,
		ReadBy:     &readerID,
		Timestamp:  readAt.UTC().Format(time.RFC3339),
	}
}

// IsUserOnline checks if a user has active connections
func (h *ChatHub) IsUserOnline(userID int64) bool {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	_, ok := h.connections[userID]
	return ok
}

// GetUserLastSeen returns the last seen time for a user
func (h *ChatHub) GetUserLastSeen(userID int64) *time.Time {
	h.lastSeenMu.RLock()
	defer h.lastSeenMu.RUnlock()

	if lastSeen, ok := h.userLastSeen[userID]; ok {
		return &lastSeen
	}
	return nil
}

// GetOnlineUsers returns list of online user IDs
func (h *ChatHub) GetOnlineUsers() []int64 {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	users := make([]int64, 0, len(h.connections))
	for userID := range h.connections {
		users = append(users, userID)
	}
	return users
}

// SendOnlineUsersList sends the current list of online users to a specific connection
// This should be called AFTER sending the connected message to avoid concurrent writes
func (h *ChatHub) SendOnlineUsersList(conn *websocket.Conn, userID int64) {
	// Get online users (excluding the connecting user)
	h.mutex.RLock()
	onlineUsers := make([]int64, 0, len(h.connections))
	for uid := range h.connections {
		if uid != userID {
			onlineUsers = append(onlineUsers, uid)
		}
	}
	h.mutex.RUnlock()

	msg := &BroadcastMessage{
		Type:        "online_users_list",
		OnlineUsers: onlineUsers,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to marshal online_users_list message")
		return
	}

	if err := h.sendToConnection(conn, data, userID); err != nil {
		h.logger.Warn().
			Err(err).
			Int64("user_id", userID).
			Msg("failed to send online_users_list to client")
	} else {
		h.logger.Debug().
			Int64("user_id", userID).
			Int("online_count", len(onlineUsers)).
			Msg("sent online_users_list to new client")
	}
}

// GetTypingUsersInChat returns list of users currently typing in a chat
func (h *ChatHub) GetTypingUsersInChat(chatID int64) []int64 {
	h.typingMu.RLock()
	defer h.typingMu.RUnlock()

	if users, ok := h.typingUsers[chatID]; ok {
		result := make([]int64, 0, len(users))
		for userID := range users {
			result = append(result, userID)
		}
		return result
	}
	return nil
}

// SafeWriteMessage safely writes a message to a connection using the connection's write mutex
// This should be used by external code (like chat_websocket.go) to avoid concurrent write panics
func (h *ChatHub) SafeWriteMessage(conn *websocket.Conn, messageType int, data []byte, userID int64) error {
	// Get write mutex for this connection
	h.mutex.RLock()
	mu, ok := h.connWriteMu[conn]
	h.mutex.RUnlock()

	if !ok {
		h.logger.Warn().Int64("user_id", userID).Msg("SafeWriteMessage: no write mutex found for connection")
		return nil // Connection might be already unregistered
	}

	// Lock the connection for writing
	mu.Lock()
	defer mu.Unlock()

	// Set write deadline
	deadline := time.Now().Add(10 * time.Second)
	if err := conn.SetWriteDeadline(deadline); err != nil {
		h.logger.Warn().Err(err).Int64("user_id", userID).Msg("SafeWriteMessage: failed to set write deadline")
		return err
	}

	if err := conn.WriteMessage(messageType, data); err != nil {
		h.logger.Warn().Err(err).Int64("user_id", userID).Msg("SafeWriteMessage: failed to write message")
		return err
	}

	return nil
}

// SendUsersLastSeen sends last seen data for offline users to a specific connection
func (h *ChatHub) SendUsersLastSeen(conn *websocket.Conn, userID int64) {
	h.lastSeenMu.RLock()
	usersLastSeen := make([]UserLastSeenInfo, 0, len(h.userLastSeen))
	for uid, lastSeen := range h.userLastSeen {
		usersLastSeen = append(usersLastSeen, UserLastSeenInfo{
			UserID:   uid,
			LastSeen: lastSeen.Format(time.RFC3339),
		})
	}
	h.lastSeenMu.RUnlock()

	if len(usersLastSeen) == 0 {
		return // No last seen data to send
	}

	msg := &BroadcastMessage{
		Type:          "users_last_seen",
		UsersLastSeen: usersLastSeen,
		Timestamp:     time.Now().UTC().Format(time.RFC3339),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to marshal users_last_seen message")
		return
	}

	if err := h.sendToConnection(conn, data, userID); err != nil {
		h.logger.Warn().
			Err(err).
			Int64("user_id", userID).
			Msg("failed to send users_last_seen to client")
	} else {
		h.logger.Debug().
			Int64("user_id", userID).
			Int("count", len(usersLastSeen)).
			Msg("sent users_last_seen to new client")
	}
}
