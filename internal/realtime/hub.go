package realtime

import (
	"encoding/json"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Inbound messages from clients
	broadcast chan *Message

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Rooms/channels for topic-based broadcasting
	rooms map[string]map[*Client]bool

	// Presence tracking
	presence map[uint]*PresenceInfo

	// Mutex for concurrent access
	mu sync.RWMutex

	logger *zap.Logger
}

// Message represents a real-time message
type Message struct {
	Type      string                 `json:"type"`
	Channel   string                 `json:"channel,omitempty"`
	Event     string                 `json:"event"`
	Payload   map[string]interface{} `json:"payload"`
	UserID    uint                   `json:"user_id,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// PresenceInfo tracks user presence
type PresenceInfo struct {
	UserID     uint      `json:"user_id"`
	Username   string    `json:"username"`
	Status     string    `json:"status"` // online, away, busy
	LastSeen   time.Time `json:"last_seen"`
	DeviceInfo string    `json:"device_info,omitempty"`
}

// NewHub creates a new Hub
func NewHub(logger *zap.Logger) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan *Message, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		rooms:      make(map[string]map[*Client]bool),
		presence:   make(map[uint]*PresenceInfo),
		logger:     logger,
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client] = true

	// Subscribe to default room if specified
	if client.Room != "" {
		if h.rooms[client.Room] == nil {
			h.rooms[client.Room] = make(map[*Client]bool)
		}
		h.rooms[client.Room][client] = true
	}

	// Update presence
	if client.UserID > 0 {
		h.presence[client.UserID] = &PresenceInfo{
			UserID:     client.UserID,
			Username:   client.Username,
			Status:     "online",
			LastSeen:   time.Now(),
			DeviceInfo: client.UserAgent,
		}

		// Broadcast presence update
		h.broadcastPresenceUpdate(client.UserID, "online")
	}

	h.logger.Info("Client registered",
		zap.Uint("user_id", client.UserID),
		zap.String("room", client.Room),
		zap.Int("total_clients", len(h.clients)),
	)
}

func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)

		// Remove from room
		if client.Room != "" && h.rooms[client.Room] != nil {
			delete(h.rooms[client.Room], client)
			if len(h.rooms[client.Room]) == 0 {
				delete(h.rooms, client.Room)
			}
		}

		// Update presence
		if client.UserID > 0 {
			if info, ok := h.presence[client.UserID]; ok {
				info.Status = "offline"
				info.LastSeen = time.Now()
				h.broadcastPresenceUpdate(client.UserID, "offline")
			}
		}

		h.logger.Info("Client unregistered",
			zap.Uint("user_id", client.UserID),
			zap.String("room", client.Room),
			zap.Int("total_clients", len(h.clients)),
		)
	}
}

func (h *Hub) broadcastMessage(message *Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	messageData, err := json.Marshal(message)
	if err != nil {
		h.logger.Error("Failed to marshal message", zap.Error(err))
		return
	}

	// Broadcast to specific channel/room
	if message.Channel != "" {
		if clients, ok := h.rooms[message.Channel]; ok {
			for client := range clients {
				select {
				case client.send <- messageData:
				default:
					// Client send buffer is full, skip
					h.logger.Warn("Client send buffer full",
						zap.Uint("user_id", client.UserID),
						zap.String("channel", message.Channel),
					)
				}
			}
		}
		return
	}

	// Broadcast to all clients
	for client := range h.clients {
		select {
		case client.send <- messageData:
		default:
			h.logger.Warn("Client send buffer full", zap.Uint("user_id", client.UserID))
		}
	}
}

func (h *Hub) broadcastPresenceUpdate(userID uint, status string) {
	message := &Message{
		Type:   "presence",
		Event:  "status_change",
		UserID: userID,
		Payload: map[string]interface{}{
			"user_id": userID,
			"status":  status,
		},
		Timestamp: time.Now(),
	}

	messageData, _ := json.Marshal(message)

	for client := range h.clients {
		select {
		case client.send <- messageData:
		default:
		}
	}
}

// BroadcastToChannel sends a message to a specific channel
func (h *Hub) BroadcastToChannel(channel string, event string, payload map[string]interface{}) {
	message := &Message{
		Type:      "broadcast",
		Channel:   channel,
		Event:     event,
		Payload:   payload,
		Timestamp: time.Now(),
	}

	h.broadcast <- message
}

// BroadcastToAll sends a message to all connected clients
func (h *Hub) BroadcastToAll(event string, payload map[string]interface{}) {
	message := &Message{
		Type:      "broadcast",
		Event:     event,
		Payload:   payload,
		Timestamp: time.Now(),
	}

	h.broadcast <- message
}

// GetPresence returns presence information for a user
func (h *Hub) GetPresence(userID uint) *PresenceInfo {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if info, ok := h.presence[userID]; ok {
		return info
	}
	return nil
}

// GetAllPresence returns all presence information
func (h *Hub) GetAllPresence() map[uint]*PresenceInfo {
	h.mu.RLock()
	defer h.mu.RUnlock()

	presence := make(map[uint]*PresenceInfo)
	for k, v := range h.presence {
		presence[k] = v
	}
	return presence
}

// GetOnlineCount returns the number of online users
func (h *Hub) GetOnlineCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	count := 0
	for _, info := range h.presence {
		if info.Status == "online" {
			count++
		}
	}
	return count
}

// GetRoomClients returns the number of clients in a room
func (h *Hub) GetRoomClients(room string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.rooms[room]; ok {
		return len(clients)
	}
	return 0
}

// GetStats returns statistics about the hub
func (h *Hub) GetStats() map[string]interface{} {
	h.mu.RLock()
	defer h.mu.RUnlock()

	stats := map[string]interface{}{
		"total_clients": len(h.clients),
		"online_users":  h.GetOnlineCount(),
		"total_rooms":   len(h.rooms),
		"rooms":         make(map[string]int),
	}

	// Add room statistics
	for room, clients := range h.rooms {
		stats["rooms"].(map[string]int)[room] = len(clients)
	}

	return stats
}

// RegisterClient registers a new client with the hub
func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}
