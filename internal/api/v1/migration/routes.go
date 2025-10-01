package migration

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go-mobile-backend-template/internal/services/auth"
	"go-mobile-backend-template/internal/services/migration"
)

// SetupMigrationRoutes sets up migration-related routes
func SetupMigrationRoutes(router *gin.RouterGroup, db *gorm.DB, config *migration.GoogleScriptsConfig, jwtService *auth.JWTService) {
	handler := NewMigrationHandler(db, config)

	// Migration management routes
	migrationGroup := router.Group("/migrations")
	// Temporarily disabled auth for testing
	// migrationGroup.Use(middleware.AuthMiddleware(jwtService))
	// migrationGroup.Use(middleware.AdminMiddleware()) // Only admins can manage migrations

	{
		// Create a new migration
		migrationGroup.POST("", handler.CreateMigration)

		// Get all migrations with pagination
		migrationGroup.GET("", handler.GetMigrations)

		// Get migration history for a specific table
		migrationGroup.GET("/history", handler.GetMigrationHistory)

		// Get a specific migration
		migrationGroup.GET("/:id", handler.GetMigration)

		// Get migration file content
		migrationGroup.GET("/:id/file", handler.GetMigrationFile)

		// Validate a migration
		migrationGroup.GET("/:id/validate", handler.ValidateMigration)

		// Execute a migration
		migrationGroup.POST("/:id/execute", handler.ExecuteMigration)

		// Rollback a migration
		migrationGroup.POST("/:id/rollback", handler.RollbackMigration)
	}

	// Public migration status endpoint (for monitoring)
	statusGroup := router.Group("/migration-status")
	// Temporarily disabled auth for testing
	// statusGroup.Use(middleware.AuthMiddleware(jwtService))
	{
		// Get migration status
		statusGroup.GET("/:id", handler.GetMigrationStatus)
	}
}
