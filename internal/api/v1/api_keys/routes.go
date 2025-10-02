package api_keys

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go-mobile-backend-template/internal/middleware"
	authService "go-mobile-backend-template/internal/services/auth"
	"go-mobile-backend-template/pkg/config"
)

// RegisterRoutes registers api_keys routes
func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, logger *zap.Logger, cfg *config.Config) {
	handler := NewHandler(db, logger)
	
	// Initialize JWT service for auth middleware
	jwtService := authService.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenExpireInt,
		cfg.JWT.RefreshTokenExpireInt,
	)

	// api_keys routes (all protected)
	apiKeysRoutes := router.Group("/api_keys")
	apiKeysRoutes.Use(middleware.AuthMiddleware(jwtService))
	{
		apiKeysRoutes.POST("", handler.CreateApikeys)
		apiKeysRoutes.GET("", handler.GetAllApikeyss)
		apiKeysRoutes.GET("/:id", handler.GetApikeys)
		apiKeysRoutes.PUT("/:id", handler.UpdateApikeys)
		apiKeysRoutes.DELETE("/:id", handler.DeleteApikeys)
	}
}
