package roles

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go-mobile-backend-template/internal/middleware"
	authService "go-mobile-backend-template/internal/services/auth"
	"go-mobile-backend-template/pkg/config"
)

// RegisterRoutes registers roles routes
func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, logger *zap.Logger, cfg *config.Config) {
	handler := NewHandler(db, logger)
	
	// Initialize JWT service for auth middleware
	jwtService := authService.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenExpireInt,
		cfg.JWT.RefreshTokenExpireInt,
	)

	// roles routes (all protected)
	rolesRoutes := router.Group("/roles")
	rolesRoutes.Use(middleware.AuthMiddleware(jwtService))
	{
		rolesRoutes.POST("", handler.CreateRoles)
		rolesRoutes.GET("", handler.GetAllRoless)
		rolesRoutes.GET("/:id", handler.GetRoles)
		rolesRoutes.PUT("/:id", handler.UpdateRoles)
		rolesRoutes.DELETE("/:id", handler.DeleteRoles)
	}
}
