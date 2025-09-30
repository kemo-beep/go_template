package realtime

import (
	"fmt"
	"net/http"

	"go-mobile-backend-template/internal/realtime"
	authService "go-mobile-backend-template/internal/services/auth"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// In production, implement proper origin checking
		return true
	},
}

type Handler struct {
	hub        *realtime.Hub
	logger     *zap.Logger
	jwtService *authService.JWTService
}

func NewHandler(hub *realtime.Hub, logger *zap.Logger, jwtService *authService.JWTService) *Handler {
	return &Handler{
		hub:        hub,
		logger:     logger,
		jwtService: jwtService,
	}
}

// HandleWebSocket godoc
// @Summary WebSocket connection endpoint
// @Description Establish WebSocket connection for real-time features
// @Tags realtime
// @Accept json
// @Produce json
// @Param token query string true "JWT token for authentication"
// @Param channel query string false "Channel to subscribe to"
// @Success 101 {string} string "Switching Protocols"
// @Router /realtime/ws [get]
func (h *Handler) HandleWebSocket(c *gin.Context) {
	// Get token from query parameter (WebSocket can't send headers)
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Token required",
		})
		return
	}

	// Validate token and get user info
	userID, username, err := h.validateToken(token)
	if err != nil {
		h.logger.Error("Invalid token for WebSocket connection", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Invalid token",
		})
		return
	}

	// Get optional channel/room parameter
	channel := c.Query("channel")

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Error("Failed to upgrade to WebSocket", zap.Error(err))
		return
	}

	// Create client
	client := realtime.NewClient(
		conn,
		h.hub,
		userID,
		username,
		channel,
		c.Request.UserAgent(),
		h.logger,
	)

	// Register client with hub
	h.hub.RegisterClient(client)

	// Start client goroutines
	go client.WritePump()
	go client.ReadPump()
}

// GetPresence godoc
// @Summary Get presence information
// @Description Get online presence for all users or a specific user
// @Tags realtime
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id query int false "User ID to get presence for"
// @Success 200 {object} map[string]interface{}
// @Router /realtime/presence [get]
func (h *Handler) GetPresence(c *gin.Context) {
	userIDStr := c.Query("user_id")

	if userIDStr != "" {
		// Get specific user presence
		var userID uint
		if _, err := fmt.Sscanf(userIDStr, "%d", &userID); err == nil {
			if info := h.hub.GetPresence(userID); info != nil {
				c.JSON(http.StatusOK, gin.H{
					"success": true,
					"data":    info,
				})
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "User not found",
		})
		return
	}

	// Get all presence
	presence := h.hub.GetAllPresence()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    presence,
		"count":   len(presence),
		"online":  h.hub.GetOnlineCount(),
	})
}

// BroadcastMessage godoc
// @Summary Broadcast a message
// @Description Send a message to a channel or all users
// @Tags realtime
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body BroadcastRequest true "Broadcast request"
// @Success 200 {object} map[string]interface{}
// @Router /realtime/broadcast [post]
func (h *Handler) BroadcastMessage(c *gin.Context) {
	var req BroadcastRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request",
		})
		return
	}

	if req.Channel != "" {
		h.hub.BroadcastToChannel(req.Channel, req.Event, req.Payload)
	} else {
		h.hub.BroadcastToAll(req.Event, req.Payload)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Message broadcasted",
	})
}

// GetStats godoc
// @Summary Get real-time statistics
// @Description Get statistics about connected clients and channels
// @Tags realtime
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /realtime/stats [get]
func (h *Handler) GetStats(c *gin.Context) {
	stats := h.hub.GetStats()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

type BroadcastRequest struct {
	Channel string                 `json:"channel,omitempty"`
	Event   string                 `json:"event" binding:"required"`
	Payload map[string]interface{} `json:"payload"`
}

// validateToken validates a JWT token and returns user info
func (h *Handler) validateToken(token string) (uint, string, error) {
	claims, err := h.jwtService.ValidateToken(token)
	if err != nil {
		return 0, "", err
	}

	return claims.UserID, claims.Email, nil
}
