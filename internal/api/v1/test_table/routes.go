package test_table

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go-mobile-backend-template/internal/middleware"
	authService "go-mobile-backend-template/internal/services/auth"
	"go-mobile-backend-template/pkg/config"
)

// RegisterRoutes registers test_table routes
func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, logger *zap.Logger, cfg *config.Config) {
	handler := NewHandler(db, logger)
	
	// Initialize JWT service for auth middleware
	jwtService := authService.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenExpireInt,
		cfg.JWT.RefreshTokenExpireInt,
	)

	// test_table routes (all protected)
	testTableRoutes := router.Group("/test_table")
	testTableRoutes.Use(middleware.AuthMiddleware(jwtService))
	{
		testTableRoutes.POST("", handler.CreateTesttable)
		testTableRoutes.GET("", handler.GetAllTesttables)
		testTableRoutes.GET("/:id", handler.GetTesttable)
		testTableRoutes.PUT("/:id", handler.UpdateTesttable)
		testTableRoutes.DELETE("/:id", handler.DeleteTesttable)
	}
}
