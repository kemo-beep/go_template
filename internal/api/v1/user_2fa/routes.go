package user_2fa

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go-mobile-backend-template/internal/middleware"
	authService "go-mobile-backend-template/internal/services/auth"
	"go-mobile-backend-template/pkg/config"
)

// RegisterRoutes registers user_2fa routes
func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, logger *zap.Logger, cfg *config.Config) {
	handler := NewHandler(db, logger)
	
	// Initialize JWT service for auth middleware
	jwtService := authService.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenExpireInt,
		cfg.JWT.RefreshTokenExpireInt,
	)

	// user_2fa routes (all protected)
	user2faRoutes := router.Group("/user_2fa")
	user2faRoutes.Use(middleware.AuthMiddleware(jwtService))
	{
		user2faRoutes.POST("", handler.CreateUser2fa)
		user2faRoutes.GET("", handler.GetAllUser2fas)
		user2faRoutes.GET("/:id", handler.GetUser2fa)
		user2faRoutes.PUT("/:id", handler.UpdateUser2fa)
		user2faRoutes.DELETE("/:id", handler.DeleteUser2fa)
	}
}
