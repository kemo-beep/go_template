package realtime

import (
	"go-mobile-backend-template/internal/middleware"
	"go-mobile-backend-template/internal/realtime"
	authService "go-mobile-backend-template/internal/services/auth"
	"go-mobile-backend-template/pkg/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RegisterRoutes registers real-time routes
func RegisterRoutes(router *gin.RouterGroup, hub *realtime.Hub, logger *zap.Logger, cfg *config.Config) {
	// Initialize JWT service for auth middleware
	jwtService := authService.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenExpireInt,
		cfg.JWT.RefreshTokenExpireInt,
	)

	handler := NewHandler(hub, logger, jwtService)

	// WebSocket endpoint (token-based auth via query parameter)
	router.GET("/ws", handler.HandleWebSocket)

	// REST endpoints for real-time features
	authorized := router.Group("")
	authorized.Use(middleware.AuthMiddleware(jwtService))
	{
		// Presence
		authorized.GET("/presence", handler.GetPresence)

		// Broadcasting (admin only recommended)
		authorized.POST("/broadcast", middleware.RequireAdmin(), handler.BroadcastMessage)

		// Stats
		authorized.GET("/stats", handler.GetStats)
	}
}
