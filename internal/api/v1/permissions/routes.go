package permissions

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go-mobile-backend-template/internal/middleware"
	authService "go-mobile-backend-template/internal/services/auth"
	"go-mobile-backend-template/pkg/config"
)

// RegisterRoutes registers permissions routes
func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, logger *zap.Logger, cfg *config.Config) {
	handler := NewHandler(db, logger)
	
	// Initialize JWT service for auth middleware
	jwtService := authService.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenExpireInt,
		cfg.JWT.RefreshTokenExpireInt,
	)

	// permissions routes (all protected)
	permissionsRoutes := router.Group("/permissions")
	permissionsRoutes.Use(middleware.AuthMiddleware(jwtService))
	{
		permissionsRoutes.POST("", handler.CreatePermissions)
		permissionsRoutes.GET("", handler.GetAllPermissionss)
		permissionsRoutes.GET("/:id", handler.GetPermissions)
		permissionsRoutes.PUT("/:id", handler.UpdatePermissions)
		permissionsRoutes.DELETE("/:id", handler.DeletePermissions)
	}
}
