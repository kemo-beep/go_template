package generator

import (
	"fmt"
	"strings"

	"go-mobile-backend-template/pkg/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// APIGeneratorMain is the main entry point for the API generator
type APIGeneratorMain struct {
	db             *gorm.DB
	logger         *zap.Logger
	config         *GeneratorConfig
	schemaAnalyzer *SchemaAnalyzer
	routeGen       *RouteGenerator
}

// NewAPIGeneratorMain creates a new main generator instance
func NewAPIGeneratorMain(db *gorm.DB, logger *zap.Logger, cfg *config.Config) *APIGeneratorMain {
	// Create generator config from app config
	genConfig := DefaultGeneratorConfig()

	// Override with app config if available
	// Note: We'll use the default config for now since we can't unmarshal directly
	// In a real implementation, you'd want to properly handle the config unmarshaling

	return &APIGeneratorMain{
		db:             db,
		logger:         logger,
		config:         genConfig,
		schemaAnalyzer: NewSchemaAnalyzer(db, logger),
		routeGen:       NewRouteGenerator(db, logger, genConfig),
	}
}

// GenerateAll generates all APIs using file-based approach
func (g *APIGeneratorMain) GenerateAll() error {
	if !g.config.Enabled {
		g.logger.Info("API generation is disabled")
		return nil
	}

	g.logger.Info("Starting file-based auto API generation")

	// Use the new file-based generator
	apiGenerator := NewAPIGenerator(g.db, g.logger, g.config)
	return apiGenerator.GenerateAll()
}

// GenerateForTable generates APIs for a specific table
func (g *APIGeneratorMain) GenerateForTable(tableName string) (*gin.Engine, error) {
	if !g.config.Enabled {
		g.logger.Info("API generation is disabled")
		return nil, nil
	}

	g.logger.Info("Generating APIs for table", zap.String("table", tableName))

	// Get table information
	table, err := g.schemaAnalyzer.GetTableByName(tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get table %s: %w", tableName, err)
	}

	// Create Gin router
	router := gin.New()

	// Generate routes for the table
	if err := g.routeGen.GenerateRoutes(router, []*TableInfo{table}); err != nil {
		return nil, fmt.Errorf("failed to generate routes for table %s: %w", tableName, err)
	}

	g.logger.Info("API generation completed for table", zap.String("table", tableName))
	return router, nil
}

// generateDocumentation generates API documentation
// TODO: Implement documentation generation
// func (g *APIGeneratorMain) generateDocumentation(tables []*TableInfo) error {
// 	if g.config.Global.Documentation == nil {
// 		return nil
// 	}
//
// 	g.logger.Info("Generating API documentation")
//
// 	// Create output directory
// 	docsDir := filepath.Join(g.config.OutputDir, "docs")
// 	if err := os.MkdirAll(docsDir, 0755); err != nil {
// 		return fmt.Errorf("failed to create docs directory: %w", err)
// 	}
//
// 	// Generate OpenAPI specification
// 	openAPIGen := NewOpenAPIGenerator(g.logger, g.config.Global.Documentation)
// 	if err := openAPIGen.GenerateOpenAPI(tables, docsDir); err != nil {
// 		return fmt.Errorf("failed to generate OpenAPI spec: %w", err)
// 	}
//
// 	g.logger.Info("API documentation generated", zap.String("dir", docsDir))
// 	return nil
// }

// generateTypes generates type definitions
// TODO: Implement type generation
// func (g *APIGeneratorMain) generateTypes(tables []*TableInfo) error {
// 	g.logger.Info("Generating type definitions")
//
// 	// Create output directory
// 	typesDir := filepath.Join(g.config.OutputDir, "types")
// 	if err := os.MkdirAll(typesDir, 0755); err != nil {
// 		return fmt.Errorf("failed to create types directory: %w", err)
// 	}
//
// 	// Generate Go types
// 	goGen := NewGoTypeGenerator(g.logger, g.config.PackageName)
// 	if err := goGen.GenerateGoTypes(tables, typesDir); err != nil {
// 		return fmt.Errorf("failed to generate Go types: %w", err)
// 	}
//
// 	// Generate TypeScript types
// 	tsGen := NewTypeScriptGenerator(g.logger)
// 	if err := tsGen.GenerateTypeScriptTypes(tables, typesDir); err != nil {
// 		return fmt.Errorf("failed to generate TypeScript types: %w", err)
// 	}
//
// 	g.logger.Info("Type definitions generated", zap.String("dir", typesDir))
// 	return nil
// }

// GetTableInfo returns information about discovered tables
func (g *APIGeneratorMain) GetTableInfo() ([]*TableInfo, error) {
	return g.schemaAnalyzer.DiscoverTables()
}

// GetGeneratedEndpoints returns information about generated endpoints
func (g *APIGeneratorMain) GetGeneratedEndpoints() ([]*GeneratedEndpoint, error) {
	tables, err := g.schemaAnalyzer.DiscoverTables()
	if err != nil {
		return nil, err
	}

	var endpoints []*GeneratedEndpoint
	for _, table := range tables {
		if !g.config.ShouldGenerateTable(table.Name) {
			continue
		}

		tableConfig := g.config.GetTableConfig(table.Name)
		for _, endpointType := range tableConfig.Endpoints {
			endpoint, err := g.generateEndpointInfo(table, endpointType, tableConfig)
			if err != nil {
				g.logger.Error("Failed to generate endpoint info",
					zap.String("table", table.Name),
					zap.String("endpoint", endpointType),
					zap.Error(err))
				continue
			}
			endpoints = append(endpoints, endpoint)
		}
	}

	return endpoints, nil
}

// generateEndpointInfo generates endpoint information
func (g *APIGeneratorMain) generateEndpointInfo(table *TableInfo, endpointType string, config *TableConfig) (*GeneratedEndpoint, error) {
	basePath := fmt.Sprintf("/api/v1/%s", strings.ToLower(table.Name))

	switch endpointType {
	case "list":
		return &GeneratedEndpoint{
			Method:      "GET",
			Path:        basePath,
			Handler:     fmt.Sprintf("List%s", g.toCamelCase(table.Name)),
			Description: fmt.Sprintf("List all %s records", strings.ToLower(table.Name)),
			Tags:        []string{table.Name},
		}, nil
	case "create":
		return &GeneratedEndpoint{
			Method:      "POST",
			Path:        basePath,
			Handler:     fmt.Sprintf("Create%s", g.toCamelCase(table.Name)),
			Description: fmt.Sprintf("Create a new %s record", strings.ToLower(table.Name)),
			Tags:        []string{table.Name},
		}, nil
	case "get":
		return &GeneratedEndpoint{
			Method:      "GET",
			Path:        fmt.Sprintf("%s/:id", basePath),
			Handler:     fmt.Sprintf("Get%s", g.toCamelCase(table.Name)),
			Description: fmt.Sprintf("Get a %s record by ID", strings.ToLower(table.Name)),
			Tags:        []string{table.Name},
		}, nil
	case "update":
		return &GeneratedEndpoint{
			Method:      "PUT",
			Path:        fmt.Sprintf("%s/:id", basePath),
			Handler:     fmt.Sprintf("Update%s", g.toCamelCase(table.Name)),
			Description: fmt.Sprintf("Update a %s record", strings.ToLower(table.Name)),
			Tags:        []string{table.Name},
		}, nil
	case "delete":
		return &GeneratedEndpoint{
			Method:      "DELETE",
			Path:        fmt.Sprintf("%s/:id", basePath),
			Handler:     fmt.Sprintf("Delete%s", g.toCamelCase(table.Name)),
			Description: fmt.Sprintf("Delete a %s record", strings.ToLower(table.Name)),
			Tags:        []string{table.Name},
		}, nil
	case "bulk":
		return &GeneratedEndpoint{
			Method:      "POST",
			Path:        fmt.Sprintf("%s/bulk", basePath),
			Handler:     fmt.Sprintf("Bulk%s", g.toCamelCase(table.Name)),
			Description: fmt.Sprintf("Bulk operations for %s records", strings.ToLower(table.Name)),
			Tags:        []string{table.Name},
		}, nil
	case "search":
		return &GeneratedEndpoint{
			Method:      "GET",
			Path:        fmt.Sprintf("%s/search", basePath),
			Handler:     fmt.Sprintf("Search%s", g.toCamelCase(table.Name)),
			Description: fmt.Sprintf("Search %s records", strings.ToLower(table.Name)),
			Tags:        []string{table.Name},
		}, nil
	case "stats":
		return &GeneratedEndpoint{
			Method:      "GET",
			Path:        fmt.Sprintf("%s/stats", basePath),
			Handler:     fmt.Sprintf("Stats%s", g.toCamelCase(table.Name)),
			Description: fmt.Sprintf("Get statistics for %s records", strings.ToLower(table.Name)),
			Tags:        []string{table.Name},
		}, nil
	case "export":
		return &GeneratedEndpoint{
			Method:      "GET",
			Path:        fmt.Sprintf("%s/export", basePath),
			Handler:     fmt.Sprintf("Export%s", g.toCamelCase(table.Name)),
			Description: fmt.Sprintf("Export %s records", strings.ToLower(table.Name)),
			Tags:        []string{table.Name},
		}, nil
	default:
		return nil, fmt.Errorf("unknown endpoint type: %s", endpointType)
	}
}

// Utility methods

func (g *APIGeneratorMain) toCamelCase(str string) string {
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

// ValidateConfig validates the generator configuration
func (g *APIGeneratorMain) ValidateConfig() error {
	return g.config.ValidateConfig()
}

// GetConfig returns the generator configuration
func (g *APIGeneratorMain) GetConfig() *GeneratorConfig {
	return g.config
}

// UpdateConfig updates the generator configuration
func (g *APIGeneratorMain) UpdateConfig(config *GeneratorConfig) error {
	if err := config.ValidateConfig(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	g.config = config
	g.routeGen = NewRouteGenerator(g.db, g.logger, config)

	return nil
}
