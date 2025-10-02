package password_reset_tokens

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go-mobile-backend-template/internal/middleware"
	authService "go-mobile-backend-template/internal/services/auth"
	"go-mobile-backend-template/pkg/config"
)

// RegisterRoutes registers password_reset_tokens routes
func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, logger *zap.Logger, cfg *config.Config) {
	handler := NewHandler(db, logger)
	
	// Initialize JWT service for auth middleware
	jwtService := authService.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenExpireInt,
		cfg.JWT.RefreshTokenExpireInt,
	)

	// password_reset_tokens routes (all protected)
	passwordResetTokensRoutes := router.Group("/password_reset_tokens")
	passwordResetTokensRoutes.Use(middleware.AuthMiddleware(jwtService))
	{
		passwordResetTokensRoutes.POST("", handler.CreatePasswordresettokens)
		passwordResetTokensRoutes.GET("", handler.GetAllPasswordresettokenss)
		passwordResetTokensRoutes.GET("/:id", handler.GetPasswordresettokens)
		passwordResetTokensRoutes.PUT("/:id", handler.UpdatePasswordresettokens)
		passwordResetTokensRoutes.DELETE("/:id", handler.DeletePasswordresettokens)
	}
}
