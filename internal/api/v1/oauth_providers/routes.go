package oauth_providers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go-mobile-backend-template/internal/middleware"
	authService "go-mobile-backend-template/internal/services/auth"
	"go-mobile-backend-template/pkg/config"
)

// RegisterRoutes registers oauth_providers routes
func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, logger *zap.Logger, cfg *config.Config) {
	handler := NewHandler(db, logger)
	
	// Initialize JWT service for auth middleware
	jwtService := authService.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenExpireInt,
		cfg.JWT.RefreshTokenExpireInt,
	)

	// oauth_providers routes (all protected)
	oauthProvidersRoutes := router.Group("/oauth_providers")
	oauthProvidersRoutes.Use(middleware.AuthMiddleware(jwtService))
	{
		oauthProvidersRoutes.POST("", handler.CreateOauthproviders)
		oauthProvidersRoutes.GET("", handler.GetAllOauthproviderss)
		oauthProvidersRoutes.GET("/:id", handler.GetOauthproviders)
		oauthProvidersRoutes.PUT("/:id", handler.UpdateOauthproviders)
		oauthProvidersRoutes.DELETE("/:id", handler.DeleteOauthproviders)
	}
}
