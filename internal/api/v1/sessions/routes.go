package sessions

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go-mobile-backend-template/internal/middleware"
	authService "go-mobile-backend-template/internal/services/auth"
	"go-mobile-backend-template/pkg/config"
)

// RegisterRoutes registers sessions routes
func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, logger *zap.Logger, cfg *config.Config) {
	handler := NewHandler(db, logger)
	
	// Initialize JWT service for auth middleware
	jwtService := authService.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenExpireInt,
		cfg.JWT.RefreshTokenExpireInt,
	)

	// sessions routes (all protected)
	sessionsRoutes := router.Group("/sessions")
	sessionsRoutes.Use(middleware.AuthMiddleware(jwtService))
	{
		sessionsRoutes.POST("", handler.CreateSessions)
		sessionsRoutes.GET("", handler.GetAllSessionss)
		sessionsRoutes.GET("/:id", handler.GetSessions)
		sessionsRoutes.PUT("/:id", handler.UpdateSessions)
		sessionsRoutes.DELETE("/:id", handler.DeleteSessions)
	}
}
