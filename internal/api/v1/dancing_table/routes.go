package dancing_table

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go-mobile-backend-template/internal/middleware"
	authService "go-mobile-backend-template/internal/services/auth"
	"go-mobile-backend-template/pkg/config"
)

// RegisterRoutes registers dancing_table routes
func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, logger *zap.Logger, cfg *config.Config) {
	handler := NewHandler(db, logger)
	
	// Initialize JWT service for auth middleware
	jwtService := authService.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenExpireInt,
		cfg.JWT.RefreshTokenExpireInt,
	)

	// dancing_table routes (all protected)
	dancingTableRoutes := router.Group("/dancing_table")
	dancingTableRoutes.Use(middleware.AuthMiddleware(jwtService))
	{
		dancingTableRoutes.POST("", handler.CreateDancingtable)
		dancingTableRoutes.GET("", handler.GetAllDancingtables)
		dancingTableRoutes.GET("/:id", handler.GetDancingtable)
		dancingTableRoutes.PUT("/:id", handler.UpdateDancingtable)
		dancingTableRoutes.DELETE("/:id", handler.DeleteDancingtable)
	}
}
