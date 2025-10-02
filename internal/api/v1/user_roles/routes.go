package user_roles

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go-mobile-backend-template/internal/middleware"
	authService "go-mobile-backend-template/internal/services/auth"
	"go-mobile-backend-template/pkg/config"
)

// RegisterRoutes registers user_roles routes
func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, logger *zap.Logger, cfg *config.Config) {
	handler := NewHandler(db, logger)
	
	// Initialize JWT service for auth middleware
	jwtService := authService.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenExpireInt,
		cfg.JWT.RefreshTokenExpireInt,
	)

	// user_roles routes (all protected)
	userRolesRoutes := router.Group("/user_roles")
	userRolesRoutes.Use(middleware.AuthMiddleware(jwtService))
	{
		userRolesRoutes.POST("", handler.CreateUserroles)
		userRolesRoutes.GET("", handler.GetAllUserroless)
		userRolesRoutes.GET("/:id", handler.GetUserroles)
		userRolesRoutes.PUT("/:id", handler.UpdateUserroles)
		userRolesRoutes.DELETE("/:id", handler.DeleteUserroles)
	}
}
