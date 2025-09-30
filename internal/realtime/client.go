package realtime

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period (must be less than pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512 * 1024 // 512 KB
)

// Client represents a WebSocket client
type Client struct {
	// The websocket connection
	conn *websocket.Conn

	// Hub this client belongs to
	hub *Hub

	// Buffered channel of outbound messages
	send chan []byte

	// User information
	UserID    uint
	Username  string
	UserAgent string

	// Room/channel this client is subscribed to
	Room string

	logger *zap.Logger
}

// ClientMessage represents an incoming message from client
type ClientMessage struct {
	Type    string                 `json:"type"`
	Event   string                 `json:"event"`
	Channel string                 `json:"channel,omitempty"`
	Payload map[string]interface{} `json:"payload,omitempty"`
}

// NewClient creates a new WebSocket client
func NewClient(conn *websocket.Conn, hub *Hub, userID uint, username string, room string, userAgent string, logger *zap.Logger) *Client {
	return &Client{
		conn:      conn,
		hub:       hub,
		send:      make(chan []byte, 256),
		UserID:    userID,
		Username:  username,
		Room:      room,
		UserAgent: userAgent,
		logger:    logger,
	}
}

// ReadPump pumps messages from the websocket connection to the hub
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Error("WebSocket error", zap.Error(err))
			}
			break
		}

		// Parse incoming message
		var clientMsg ClientMessage
		if err := json.Unmarshal(message, &clientMsg); err != nil {
			c.logger.Warn("Failed to parse client message", zap.Error(err))
			continue
		}

		// Handle different message types
		c.handleMessage(&clientMsg)
	}
}

// WritePump pumps messages from the hub to the websocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) handleMessage(msg *ClientMessage) {
	switch msg.Type {
	case "subscribe":
		c.handleSubscribe(msg)
	case "unsubscribe":
		c.handleUnsubscribe(msg)
	case "broadcast":
		c.handleBroadcast(msg)
	case "presence":
		c.handlePresence(msg)
	default:
		c.logger.Warn("Unknown message type", zap.String("type", msg.Type))
	}
}

func (c *Client) handleSubscribe(msg *ClientMessage) {
	if msg.Channel == "" {
		return
	}

	c.hub.mu.Lock()
	defer c.hub.mu.Unlock()

	// Add to new room
	if c.hub.rooms[msg.Channel] == nil {
		c.hub.rooms[msg.Channel] = make(map[*Client]bool)
	}
	c.hub.rooms[msg.Channel][c] = true

	// Remove from old room if different
	if c.Room != "" && c.Room != msg.Channel && c.hub.rooms[c.Room] != nil {
		delete(c.hub.rooms[c.Room], c)
		if len(c.hub.rooms[c.Room]) == 0 {
			delete(c.hub.rooms, c.Room)
		}
	}

	c.Room = msg.Channel

	c.logger.Info("Client subscribed to channel",
		zap.Uint("user_id", c.UserID),
		zap.String("channel", msg.Channel),
	)

	// Send confirmation
	response := map[string]interface{}{
		"type":    "subscribed",
		"channel": msg.Channel,
		"success": true,
	}
	data, _ := json.Marshal(response)
	c.send <- data
}

func (c *Client) handleUnsubscribe(msg *ClientMessage) {
	if msg.Channel == "" {
		return
	}

	c.hub.mu.Lock()
	defer c.hub.mu.Unlock()

	if c.hub.rooms[msg.Channel] != nil {
		delete(c.hub.rooms[msg.Channel], c)
		if len(c.hub.rooms[msg.Channel]) == 0 {
			delete(c.hub.rooms, msg.Channel)
		}
	}

	c.logger.Info("Client unsubscribed from channel",
		zap.Uint("user_id", c.UserID),
		zap.String("channel", msg.Channel),
	)

	// Send confirmation
	response := map[string]interface{}{
		"type":    "unsubscribed",
		"channel": msg.Channel,
		"success": true,
	}
	data, _ := json.Marshal(response)
	c.send <- data
}

func (c *Client) handleBroadcast(msg *ClientMessage) {
	// Create broadcast message
	broadcastMsg := &Message{
		Type:      "broadcast",
		Channel:   msg.Channel,
		Event:     msg.Event,
		Payload:   msg.Payload,
		UserID:    c.UserID,
		Timestamp: time.Now(),
	}

	c.hub.broadcast <- broadcastMsg
}

func (c *Client) handlePresence(msg *ClientMessage) {
	status, ok := msg.Payload["status"].(string)
	if !ok {
		return
	}

	c.hub.mu.Lock()
	if info, exists := c.hub.presence[c.UserID]; exists {
		info.Status = status
		info.LastSeen = time.Now()
	}
	c.hub.mu.Unlock()

	// Broadcast presence change
	c.hub.broadcastPresenceUpdate(c.UserID, status)
}
