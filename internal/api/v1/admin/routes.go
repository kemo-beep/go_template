package admin

import (
	"go-mobile-backend-template/internal/middleware"
	authService "go-mobile-backend-template/internal/services/auth"
	"go-mobile-backend-template/pkg/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RegisterRoutes registers admin routes
func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, logger *zap.Logger, cfg *config.Config) {
	// Initialize JWT service for auth middleware
	jwtService := authService.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenExpireInt,
		cfg.JWT.RefreshTokenExpireInt,
	)

	// Apply auth middleware
	router.Use(middleware.AuthMiddleware(jwtService))

	// Apply admin middleware - all admin routes require admin role
	router.Use(middleware.RequireAdmin())

	userHandler := NewUserHandler(db, logger)
	roleHandler := NewRoleHandler(db, logger)
	dbHandler := NewDatabaseHandler(db, logger)
	tableHandler := NewTableManagerHandler(db, logger)
	tableDataHandler := NewTableDataHandler(db)

	// User management routes
	users := router.Group("/users")
	{
		users.GET("", userHandler.ListUsers)
		users.GET("/:id", userHandler.GetUser)
		users.PUT("/:id", userHandler.UpdateUser)
		users.DELETE("/:id", userHandler.DeleteUser)
		users.POST("/:id/roles", userHandler.AssignRole)
		users.DELETE("/:id/roles/:roleId", userHandler.RemoveRole)
	}

	// Role management routes
	roles := router.Group("/roles")
	{
		roles.GET("", roleHandler.ListRoles)
		roles.POST("", roleHandler.CreateRole)
		roles.GET("/:id", roleHandler.GetRole)
		roles.POST("/:id/permissions", roleHandler.AssignPermissions)
	}

	// Permission routes
	router.GET("/permissions", roleHandler.ListPermissions)

	// Database management routes
	database := router.Group("/database")
	{
		// Read operations
		database.GET("/tables", dbHandler.ListTables)
		database.GET("/tables/:tableName/schema", dbHandler.GetTableSchema)
		database.GET("/tables/:tableName/data", dbHandler.GetTableData)
		database.POST("/query", dbHandler.ExecuteQuery)
		database.GET("/stats", dbHandler.GetDatabaseStats)

		// Table management operations (CREATE, ALTER, DROP)
		database.POST("/tables", tableHandler.CreateTable)
		database.DELETE("/tables/:tableName", tableHandler.DropTable)
		database.PUT("/tables/:tableName/rename", tableHandler.RenameTable)
		database.POST("/tables/:tableName/columns", tableHandler.AddColumn)
		database.DELETE("/tables/:tableName/columns/:columnName", tableHandler.DropColumn)

		// Table data operations (INSERT, UPDATE, DELETE rows)
		database.POST("/tables/:tableName/rows", tableDataHandler.InsertTableRow)
		database.PUT("/tables/:tableName/rows/:pkValue", tableDataHandler.UpdateTableRow)
		database.DELETE("/tables/:tableName/rows/:pkValue", tableDataHandler.DeleteTableRow)
	}
}
