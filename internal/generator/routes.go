package generator

import (
	"fmt"
	"strings"

	"go-mobile-backend-template/internal/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RouteGenerator generates routes for generated endpoints
type RouteGenerator struct {
	db         *gorm.DB
	logger     *zap.Logger
	config     *GeneratorConfig
	handlerGen *CRUDHandlerGenerator
}

// NewRouteGenerator creates a new route generator
func NewRouteGenerator(db *gorm.DB, logger *zap.Logger, config *GeneratorConfig) *RouteGenerator {
	return &RouteGenerator{
		db:         db,
		logger:     logger,
		config:     config,
		handlerGen: NewCRUDHandlerGenerator(db, logger, config),
	}
}

// GenerateRoutes generates all routes for discovered tables
func (g *RouteGenerator) GenerateRoutes(router *gin.Engine, tables []*TableInfo) error {
	g.logger.Info("Generating routes for tables", zap.Int("count", len(tables)))

	// Create API v1 group
	apiV1 := router.Group("/api/v1")

	// Apply global middleware
	apiV1.Use(
		middleware.Logger(g.logger),
		middleware.Recovery(g.logger),
		middleware.DevelopmentCORS(),
		middleware.SecurityHeaders(),
	)

	// Generate routes for each table
	for _, table := range tables {
		if !g.config.ShouldGenerateTable(table.Name) {
			g.logger.Info("Skipping table", zap.String("table", table.Name))
			continue
		}

		g.logger.Info("Generating routes for table", zap.String("table", table.Name))

		if err := g.generateTableRoutes(apiV1, table); err != nil {
			g.logger.Error("Failed to generate routes for table",
				zap.String("table", table.Name),
				zap.Error(err))
			continue
		}
	}

	g.logger.Info("Route generation completed")
	return nil
}

// generateTableRoutes generates routes for a specific table
func (g *RouteGenerator) generateTableRoutes(router *gin.RouterGroup, table *TableInfo) error {
	tableConfig := g.config.GetTableConfig(table.Name)
	tableName := strings.ToLower(table.Name)

	// Create table group
	tableGroup := router.Group("/" + tableName)

	// Apply table-specific middleware
	if err := g.applyTableMiddleware(tableGroup, table, tableConfig); err != nil {
		return fmt.Errorf("failed to apply middleware: %w", err)
	}

	// Generate handlers
	handlers, err := g.handlerGen.GenerateHandlers(table)
	if err != nil {
		return fmt.Errorf("failed to generate handlers: %w", err)
	}

	// Register routes
	for _, endpointType := range tableConfig.Endpoints {
		handlerName := fmt.Sprintf("%s%s", g.toCamelCase(endpointType), g.toCamelCase(table.Name))
		handler, exists := handlers[handlerName]
		if !exists {
			g.logger.Warn("Handler not found",
				zap.String("table", table.Name),
				zap.String("handler", handlerName))
			continue
		}

		if err := g.registerRoute(tableGroup, table, endpointType, handler, tableConfig); err != nil {
			g.logger.Error("Failed to register route",
				zap.String("table", table.Name),
				zap.String("endpoint", endpointType),
				zap.Error(err))
			continue
		}
	}

	// Generate relationship routes
	if err := g.generateRelationshipRoutes(tableGroup, table, tableConfig); err != nil {
		g.logger.Error("Failed to generate relationship routes",
			zap.String("table", table.Name),
			zap.Error(err))
	}

	return nil
}

// registerRoute registers a specific route
func (g *RouteGenerator) registerRoute(router *gin.RouterGroup, table *TableInfo, endpointType string, handler gin.HandlerFunc, config *TableConfig) error {
	tableName := strings.ToLower(table.Name)

	switch endpointType {
	case "list":
		router.GET("", handler)
		g.logger.Debug("Registered route",
			zap.String("method", "GET"),
			zap.String("path", "/api/v1/"+tableName),
			zap.String("handler", "List"+g.toCamelCase(table.Name)))

	case "create":
		router.POST("", handler)
		g.logger.Debug("Registered route",
			zap.String("method", "POST"),
			zap.String("path", "/api/v1/"+tableName),
			zap.String("handler", "Create"+g.toCamelCase(table.Name)))

	case "get":
		router.GET("/:id", handler)
		g.logger.Debug("Registered route",
			zap.String("method", "GET"),
			zap.String("path", "/api/v1/"+tableName+"/:id"),
			zap.String("handler", "Get"+g.toCamelCase(table.Name)))

	case "update":
		router.PUT("/:id", handler)
		router.PATCH("/:id", handler) // Also support PATCH for partial updates
		g.logger.Debug("Registered route",
			zap.String("method", "PUT/PATCH"),
			zap.String("path", "/api/v1/"+tableName+"/:id"),
			zap.String("handler", "Update"+g.toCamelCase(table.Name)))

	case "delete":
		router.DELETE("/:id", handler)
		g.logger.Debug("Registered route",
			zap.String("method", "DELETE"),
			zap.String("path", "/api/v1/"+tableName+"/:id"),
			zap.String("handler", "Delete"+g.toCamelCase(table.Name)))

	case "bulk":
		router.POST("/bulk", handler)
		g.logger.Debug("Registered route",
			zap.String("method", "POST"),
			zap.String("path", "/api/v1/"+tableName+"/bulk"),
			zap.String("handler", "Bulk"+g.toCamelCase(table.Name)))

	case "search":
		router.GET("/search", handler)
		g.logger.Debug("Registered route",
			zap.String("method", "GET"),
			zap.String("path", "/api/v1/"+tableName+"/search"),
			zap.String("handler", "Search"+g.toCamelCase(table.Name)))

	case "stats":
		router.GET("/stats", handler)
		g.logger.Debug("Registered route",
			zap.String("method", "GET"),
			zap.String("path", "/api/v1/"+tableName+"/stats"),
			zap.String("handler", "Stats"+g.toCamelCase(table.Name)))

	case "export":
		router.GET("/export", handler)
		g.logger.Debug("Registered route",
			zap.String("method", "GET"),
			zap.String("path", "/api/v1/"+tableName+"/export"),
			zap.String("handler", "Export"+g.toCamelCase(table.Name)))

	default:
		return fmt.Errorf("unknown endpoint type: %s", endpointType)
	}

	return nil
}

// generateRelationshipRoutes generates routes for table relationships
func (g *RouteGenerator) generateRelationshipRoutes(router *gin.RouterGroup, table *TableInfo, config *TableConfig) error {
	for _, relationship := range config.Relationships {
		relGroup := router.Group("/:id/" + relationship)

		// Apply relationship-specific middleware
		if err := g.applyRelationshipMiddleware(relGroup, table, relationship, config); err != nil {
			g.logger.Error("Failed to apply relationship middleware",
				zap.String("table", table.Name),
				zap.String("relationship", relationship),
				zap.Error(err))
			continue
		}

		// Generate relationship handlers
		handlers, err := g.generateRelationshipHandlers(table, relationship, config)
		if err != nil {
			g.logger.Error("Failed to generate relationship handlers",
				zap.String("table", table.Name),
				zap.String("relationship", relationship),
				zap.Error(err))
			continue
		}

		// Register relationship routes
		for method, handler := range handlers {
			relGroup.Handle(method, "", handler)
			g.logger.Debug("Registered relationship route",
				zap.String("method", method),
				zap.String("path", "/api/v1/"+strings.ToLower(table.Name)+"/:id/"+relationship),
				zap.String("table", table.Name),
				zap.String("relationship", relationship))
		}
	}

	return nil
}

// generateRelationshipHandlers generates handlers for table relationships
func (g *RouteGenerator) generateRelationshipHandlers(table *TableInfo, relationship string, config *TableConfig) (map[string]gin.HandlerFunc, error) {
	handlers := make(map[string]gin.HandlerFunc)

	// Generate list relationship handler
	handlers["GET"] = g.generateListRelationshipHandler(table, relationship, config)

	// Generate create relationship handler
	handlers["POST"] = g.generateCreateRelationshipHandler(table, relationship, config)

	// Generate delete relationship handler
	handlers["DELETE"] = g.generateDeleteRelationshipHandler(table, relationship, config)

	return handlers, nil
}

// generateListRelationshipHandler generates a handler to list related records
func (g *RouteGenerator) generateListRelationshipHandler(table *TableInfo, relationship string, config *TableConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// This would implement listing related records
		// For now, return a placeholder response
		c.JSON(200, gin.H{
			"message": fmt.Sprintf("List %s for %s", relationship, table.Name),
			"data":    []interface{}{},
		})
	}
}

// generateCreateRelationshipHandler generates a handler to create relationships
func (g *RouteGenerator) generateCreateRelationshipHandler(table *TableInfo, relationship string, config *TableConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// This would implement creating relationships
		// For now, return a placeholder response
		c.JSON(201, gin.H{
			"message": fmt.Sprintf("Create %s for %s", relationship, table.Name),
			"data":    gin.H{},
		})
	}
}

// generateDeleteRelationshipHandler generates a handler to delete relationships
func (g *RouteGenerator) generateDeleteRelationshipHandler(table *TableInfo, relationship string, config *TableConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// This would implement deleting relationships
		// For now, return a placeholder response
		c.JSON(204, nil)
	}
}

// applyTableMiddleware applies middleware specific to a table
func (g *RouteGenerator) applyTableMiddleware(router *gin.RouterGroup, table *TableInfo, config *TableConfig) error {
	// Apply rate limiting
	if config.Security != nil && config.Security.RateLimit != nil {
		// This would apply rate limiting middleware
		// For now, we'll skip it as it requires Redis setup
	}

	// Apply RBAC middleware
	if config.Security != nil && config.Security.RBAC != nil {
		// This would apply RBAC middleware
		// For now, we'll skip it as it requires authentication setup
	}

	// Apply audit logging
	if config.Security != nil && config.Security.AuditLog {
		// This would apply audit logging middleware
		// For now, we'll skip it as it requires database setup
	}

	return nil
}

// applyRelationshipMiddleware applies middleware specific to relationships
func (g *RouteGenerator) applyRelationshipMiddleware(router *gin.RouterGroup, table *TableInfo, relationship string, config *TableConfig) error {
	// Apply the same middleware as the parent table
	return g.applyTableMiddleware(router, table, config)
}

// Utility methods

func (g *RouteGenerator) toCamelCase(str string) string {
	if str == "" {
		return ""
	}

	parts := strings.Split(str, "_")
	result := strings.Title(parts[0])

	for i := 1; i < len(parts); i++ {
		result += strings.Title(parts[i])
	}

	return result
}

// GetRouteInfo returns information about generated routes
func (g *RouteGenerator) GetRouteInfo() map[string]interface{} {
	routes := make(map[string]interface{})

	// This would return information about all generated routes
	// For now, return empty map

	return routes
}

// ValidateRoutes validates that all routes are properly configured
func (g *RouteGenerator) ValidateRoutes() error {
	// This would validate that all routes are properly configured
	// For now, return nil

	return nil
}
