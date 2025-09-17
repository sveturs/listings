package tracking

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
)

// Hub управляет WebSocket соединениями для трекинга
type Hub struct {
	// Registered connections по delivery_id
	connections map[int]map[*websocket.Conn]bool

	// Buffered channel of inbound messages
	broadcast chan []byte

	// Register requests from connections
	register chan subscription

	// Unregister requests from connections
	unregister chan subscription

	// Mutex для защиты connections map
	mutex sync.RWMutex
}

// subscription представляет подписку на updates для delivery
type subscription struct {
	conn       *websocket.Conn
	deliveryID int
}

// LocationUpdateMessage сообщение об обновлении локации
type LocationUpdateMessage struct {
	Type        string  `json:"type"`
	DeliveryID  int     `json:"delivery_id"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Speed       float64 `json:"speed"`
	Heading     int     `json:"heading"`
	Timestamp   string  `json:"timestamp"`
	CourierName string  `json:"courier_name,omitempty"`
}

// StatusUpdateMessage сообщение об обновлении статуса
type StatusUpdateMessage struct {
	Type       string `json:"type"`
	DeliveryID int    `json:"delivery_id"`
	Status     string `json:"status"`
	Timestamp  string `json:"timestamp"`
	Message    string `json:"message,omitempty"`
}

// ETAUpdateMessage сообщение об обновлении ETA
type ETAUpdateMessage struct {
	Type       string `json:"type"`
	DeliveryID int    `json:"delivery_id"`
	ETA        string `json:"eta"`
	Timestamp  string `json:"timestamp"`
}

// NewHub создает новый WebSocket hub
func NewHub() *Hub {
	return &Hub{
		broadcast:   make(chan []byte),
		register:    make(chan subscription),
		unregister:  make(chan subscription),
		connections: make(map[int]map[*websocket.Conn]bool),
	}
}

// Run запускает hub
func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case sub := <-h.register:
			h.mutex.Lock()
			if h.connections[sub.deliveryID] == nil {
				h.connections[sub.deliveryID] = make(map[*websocket.Conn]bool)
			}
			h.connections[sub.deliveryID][sub.conn] = true
			h.mutex.Unlock()

			log.Printf("WebSocket client connected to delivery %d", sub.deliveryID)

		case sub := <-h.unregister:
			h.mutex.Lock()
			if connections, ok := h.connections[sub.deliveryID]; ok {
				if _, ok := connections[sub.conn]; ok {
					delete(connections, sub.conn)
					sub.conn.Close()

					// Если больше нет соединений для этой доставки, удаляем мапу
					if len(connections) == 0 {
						delete(h.connections, sub.deliveryID)
					}
				}
			}
			h.mutex.Unlock()

			log.Printf("WebSocket client disconnected from delivery %d", sub.deliveryID)

		case message := <-h.broadcast:
			// Определяем delivery_id из сообщения
			var baseMsg struct {
				DeliveryID int `json:"delivery_id"`
			}

			if err := json.Unmarshal(message, &baseMsg); err != nil {
				log.Printf("Failed to unmarshal broadcast message: %v", err)
				continue
			}

			h.sendToDeliveryClients(baseMsg.DeliveryID, message)
		}
	}
}

// RegisterConnection регистрирует новое WebSocket соединение
func (h *Hub) RegisterConnection(conn *websocket.Conn, deliveryID int) {
	h.register <- subscription{conn: conn, deliveryID: deliveryID}
}

// UnregisterConnection отменяет регистрацию WebSocket соединения
func (h *Hub) UnregisterConnection(conn *websocket.Conn, deliveryID int) {
	h.unregister <- subscription{conn: conn, deliveryID: deliveryID}
}

// BroadcastLocationUpdate отправляет обновление локации всем подключенным клиентам
func (h *Hub) BroadcastLocationUpdate(deliveryID int, latitude, longitude, speed float64, heading int) {
	message := LocationUpdateMessage{
		Type:       "location_update",
		DeliveryID: deliveryID,
		Latitude:   latitude,
		Longitude:  longitude,
		Speed:      speed,
		Heading:    heading,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal location update: %v", err)
		return
	}

	h.broadcast <- data
}

// BroadcastStatusUpdate отправляет обновление статуса
func (h *Hub) BroadcastStatusUpdate(deliveryID int, status, message string) {
	statusMessage := StatusUpdateMessage{
		Type:       "status_update",
		DeliveryID: deliveryID,
		Status:     status,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Message:    message,
	}

	data, err := json.Marshal(statusMessage)
	if err != nil {
		log.Printf("Failed to marshal status update: %v", err)
		return
	}

	h.broadcast <- data
}

// BroadcastETAUpdate отправляет обновление ETA
func (h *Hub) BroadcastETAUpdate(deliveryID int, eta time.Time) {
	message := ETAUpdateMessage{
		Type:       "eta_update",
		DeliveryID: deliveryID,
		ETA:        eta.UTC().Format(time.RFC3339),
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal ETA update: %v", err)
		return
	}

	h.broadcast <- data
}

// sendToDeliveryClients отправляет сообщение всем клиентам определенной доставки
func (h *Hub) sendToDeliveryClients(deliveryID int, message []byte) {
	h.mutex.RLock()
	connections, ok := h.connections[deliveryID]
	h.mutex.RUnlock()

	if !ok {
		return // Нет подключенных клиентов для этой доставки
	}

	for conn := range connections {
		select {
		case <-time.After(time.Second):
			log.Printf("Write timeout for client, closing connection")
			h.UnregisterConnection(conn, deliveryID)
		default:
			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("WebSocket write error: %v", err)
				h.UnregisterConnection(conn, deliveryID)
			}
		}
	}
}

// GetActiveConnections возвращает статистику активных соединений
func (h *Hub) GetActiveConnections() map[string]int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	stats := make(map[string]int)
	totalConnections := 0

	for deliveryID, connections := range h.connections {
		count := len(connections)
		stats[string(rune(deliveryID))] = count
		totalConnections += count
	}

	stats["total"] = totalConnections
	return stats
}

// HandleWebSocket обрабатывает WebSocket соединение для трекинга
func (h *Hub) HandleWebSocket(conn *websocket.Conn, deliveryID int) {
	// Регистрируем соединение
	h.RegisterConnection(conn, deliveryID)

	// Отправляем начальное подтверждение соединения
	initialMessage := map[string]interface{}{
		"type":        "connection_established",
		"delivery_id": deliveryID,
		"timestamp":   time.Now().UTC().Format(time.RFC3339),
	}

	if data, err := json.Marshal(initialMessage); err == nil {
		conn.WriteMessage(websocket.TextMessage, data)
	}

	// Обрабатываем сообщения от клиента
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		if messageType == websocket.CloseMessage {
			log.Printf("WebSocket close message received")
			break
		}

		// Обрабатываем входящие сообщения (например, ping)
		if messageType == websocket.TextMessage {
			var msg map[string]interface{}
			if err := json.Unmarshal(message, &msg); err == nil {
				if msgType, ok := msg["type"].(string); ok && msgType == "ping" {
					// Отвечаем pong
					pongMessage := map[string]interface{}{
						"type":      "pong",
						"timestamp": time.Now().UTC().Format(time.RFC3339),
					}
					if data, err := json.Marshal(pongMessage); err == nil {
						conn.WriteMessage(websocket.TextMessage, data)
					}
				}
			}
		}
	}

	// Отменяем регистрацию при закрытии соединения
	h.UnregisterConnection(conn, deliveryID)
}