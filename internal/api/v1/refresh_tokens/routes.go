package refresh_tokens

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go-mobile-backend-template/internal/middleware"
	authService "go-mobile-backend-template/internal/services/auth"
	"go-mobile-backend-template/pkg/config"
)

// RegisterRoutes registers refresh_tokens routes
func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, logger *zap.Logger, cfg *config.Config) {
	handler := NewHandler(db, logger)
	
	// Initialize JWT service for auth middleware
	jwtService := authService.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenExpireInt,
		cfg.JWT.RefreshTokenExpireInt,
	)

	// refresh_tokens routes (all protected)
	refreshTokensRoutes := router.Group("/refresh_tokens")
	refreshTokensRoutes.Use(middleware.AuthMiddleware(jwtService))
	{
		refreshTokensRoutes.POST("", handler.CreateRefreshtokens)
		refreshTokensRoutes.GET("", handler.GetAllRefreshtokenss)
		refreshTokensRoutes.GET("/:id", handler.GetRefreshtokens)
		refreshTokensRoutes.PUT("/:id", handler.UpdateRefreshtokens)
		refreshTokensRoutes.DELETE("/:id", handler.DeleteRefreshtokens)
	}
}
