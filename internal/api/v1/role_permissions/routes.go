package role_permissions

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go-mobile-backend-template/internal/middleware"
	authService "go-mobile-backend-template/internal/services/auth"
	"go-mobile-backend-template/pkg/config"
)

// RegisterRoutes registers role_permissions routes
func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, logger *zap.Logger, cfg *config.Config) {
	handler := NewHandler(db, logger)
	
	// Initialize JWT service for auth middleware
	jwtService := authService.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenExpireInt,
		cfg.JWT.RefreshTokenExpireInt,
	)

	// role_permissions routes (all protected)
	rolePermissionsRoutes := router.Group("/role_permissions")
	rolePermissionsRoutes.Use(middleware.AuthMiddleware(jwtService))
	{
		rolePermissionsRoutes.POST("", handler.CreateRolepermissions)
		rolePermissionsRoutes.GET("", handler.GetAllRolepermissionss)
		rolePermissionsRoutes.GET("/:id", handler.GetRolepermissions)
		rolePermissionsRoutes.PUT("/:id", handler.UpdateRolepermissions)
		rolePermissionsRoutes.DELETE("/:id", handler.DeleteRolepermissions)
	}
}
