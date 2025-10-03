package wow_table

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go-mobile-backend-template/internal/middleware"
	authService "go-mobile-backend-template/internal/services/auth"
	"go-mobile-backend-template/pkg/config"
)

// RegisterRoutes registers wow_table routes
func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, logger *zap.Logger, cfg *config.Config) {
	handler := NewHandler(db, logger)
	
	// Initialize JWT service for auth middleware
	jwtService := authService.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenExpireInt,
		cfg.JWT.RefreshTokenExpireInt,
	)

	// wow_table routes (all protected)
	wowTableRoutes := router.Group("/wow_table")
	wowTableRoutes.Use(middleware.AuthMiddleware(jwtService))
	{
		wowTableRoutes.POST("", handler.CreateWowtable)
		wowTableRoutes.GET("", handler.GetAllWowtables)
		wowTableRoutes.GET("/:id", handler.GetWowtable)
		wowTableRoutes.PUT("/:id", handler.UpdateWowtable)
		wowTableRoutes.DELETE("/:id", handler.DeleteWowtable)
	}
}
