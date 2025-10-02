package users

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go-mobile-backend-template/internal/middleware"
	authService "go-mobile-backend-template/internal/services/auth"
	"go-mobile-backend-template/pkg/config"
)

// RegisterRoutes registers users routes
func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, logger *zap.Logger, cfg *config.Config) {
	handler := NewHandler(db, logger)
	
	// Initialize JWT service for auth middleware
	jwtService := authService.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenExpireInt,
		cfg.JWT.RefreshTokenExpireInt,
	)

	// users routes (all protected)
	usersRoutes := router.Group("/users")
	usersRoutes.Use(middleware.AuthMiddleware(jwtService))
	{
		usersRoutes.POST("", handler.CreateUsers)
		usersRoutes.GET("", handler.GetAllUserss)
		usersRoutes.GET("/:id", handler.GetUsers)
		usersRoutes.PUT("/:id", handler.UpdateUsers)
		usersRoutes.DELETE("/:id", handler.DeleteUsers)
	}
}
