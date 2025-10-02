package email_verification_tokens

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go-mobile-backend-template/internal/middleware"
	authService "go-mobile-backend-template/internal/services/auth"
	"go-mobile-backend-template/pkg/config"
)

// RegisterRoutes registers email_verification_tokens routes
func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, logger *zap.Logger, cfg *config.Config) {
	handler := NewHandler(db, logger)
	
	// Initialize JWT service for auth middleware
	jwtService := authService.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenExpireInt,
		cfg.JWT.RefreshTokenExpireInt,
	)

	// email_verification_tokens routes (all protected)
	emailVerificationTokensRoutes := router.Group("/email_verification_tokens")
	emailVerificationTokensRoutes.Use(middleware.AuthMiddleware(jwtService))
	{
		emailVerificationTokensRoutes.POST("", handler.CreateEmailverificationtokens)
		emailVerificationTokensRoutes.GET("", handler.GetAllEmailverificationtokenss)
		emailVerificationTokensRoutes.GET("/:id", handler.GetEmailverificationtokens)
		emailVerificationTokensRoutes.PUT("/:id", handler.UpdateEmailverificationtokens)
		emailVerificationTokensRoutes.DELETE("/:id", handler.DeleteEmailverificationtokens)
	}
}
