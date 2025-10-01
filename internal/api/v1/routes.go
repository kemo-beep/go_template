package v1

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go-mobile-backend-template/internal/api/v1/admin"
	"go-mobile-backend-template/internal/api/v1/auth"
	"go-mobile-backend-template/internal/api/v1/files"
	"go-mobile-backend-template/internal/api/v1/migration"
	realtimeAPI "go-mobile-backend-template/internal/api/v1/realtime"
	"go-mobile-backend-template/internal/api/v1/users"
	"go-mobile-backend-template/internal/middleware"
	"go-mobile-backend-template/internal/realtime"
	authService "go-mobile-backend-template/internal/services/auth"
	migrationService "go-mobile-backend-template/internal/services/migration"
	"go-mobile-backend-template/pkg/config"
)

// RegisterRoutes registers all v1 API routes
func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, logger *zap.Logger, cfg *config.Config, hub *realtime.Hub) {
	// Initialize handlers
	authHandler := auth.NewHandler(db, logger, cfg)
	usersHandler := users.NewHandler(db, logger)
	filesHandler := files.NewHandler(db, logger, cfg)

	// Initialize JWT service for auth middleware
	jwtService := authService.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenExpireInt,
		cfg.JWT.RefreshTokenExpireInt,
	)

	// Public auth routes
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/refresh", authHandler.RefreshToken)
	}

	// Protected auth routes
	authProtected := router.Group("/auth")
	authProtected.Use(middleware.AuthMiddleware(jwtService))
	{
		authProtected.POST("/logout", authHandler.Logout)
	}

	// User routes (all protected)
	userRoutes := router.Group("/users")
	userRoutes.Use(middleware.AuthMiddleware(jwtService))
	{
		userRoutes.GET("/me", usersHandler.GetProfile)
		userRoutes.PUT("/me", usersHandler.UpdateProfile)
		userRoutes.DELETE("/me", usersHandler.DeleteAccount)
		userRoutes.POST("/me/change-password", usersHandler.ChangePassword)
	}

	// File routes (all protected)
	fileRoutes := router.Group("/files")
	fileRoutes.Use(middleware.AuthMiddleware(jwtService))
	{
		fileRoutes.POST("/upload", filesHandler.Upload)
		fileRoutes.GET("", filesHandler.ListFiles)
		fileRoutes.GET("/:id", filesHandler.GetFile)
		fileRoutes.GET("/:id/download", filesHandler.GetDownloadURL)
		fileRoutes.DELETE("/:id", filesHandler.DeleteFile)
	}

	// Admin routes (admin only)
	adminRoutes := router.Group("/admin")
	admin.RegisterRoutes(adminRoutes, db, logger, cfg)

	// Migration routes (admin only)
	migrationConfig := &migrationService.GoogleScriptsConfig{
		ScriptURL:     cfg.GoogleScripts.URL,
		AccessToken:   cfg.GoogleScripts.AccessToken,
		ProjectID:     cfg.GoogleScripts.ProjectID,
		MigrationsDir: "./internal/db/migrations", // Use the existing migrations directory
	}
	migration.SetupMigrationRoutes(router, db, migrationConfig, jwtService)

	// Real-time routes (WebSocket, presence, etc.)
	realtimeRoutes := router.Group("/realtime")
	realtimeAPI.RegisterRoutes(realtimeRoutes, hub, logger, cfg)
}
