package files

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go-mobile-backend-template/internal/middleware"
	authService "go-mobile-backend-template/internal/services/auth"
	"go-mobile-backend-template/pkg/config"
)

// RegisterRoutes registers files routes
func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, logger *zap.Logger, cfg *config.Config) {
	handler := NewHandler(db, logger)
	
	// Initialize JWT service for auth middleware
	jwtService := authService.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenExpireInt,
		cfg.JWT.RefreshTokenExpireInt,
	)

	// files routes (all protected)
	filesRoutes := router.Group("/files")
	filesRoutes.Use(middleware.AuthMiddleware(jwtService))
	{
		filesRoutes.POST("", handler.CreateFiles)
		filesRoutes.GET("", handler.GetAllFiless)
		filesRoutes.GET("/:id", handler.GetFiles)
		filesRoutes.PUT("/:id", handler.UpdateFiles)
		filesRoutes.DELETE("/:id", handler.DeleteFiles)
	}
}
