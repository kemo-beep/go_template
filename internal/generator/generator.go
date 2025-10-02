package generator

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// APIGenerator is the main generator that orchestrates the API generation process
type APIGenerator struct {
	db             *gorm.DB
	logger         *zap.Logger
	config         *GeneratorConfig
	schemaAnalyzer *SchemaAnalyzer
	tables         []*TableInfo
	endpoints      []*GeneratedEndpoint
}

// GeneratedEndpoint represents a generated API endpoint
type GeneratedEndpoint struct {
	Method      string               `json:"method"`
	Path        string               `json:"path"`
	Handler     string               `json:"handler"` // Handler function name
	Middleware  []string             `json:"middleware"`
	Validation  *ValidationRules     `json:"validation"`
	Pagination  *PaginationConfig    `json:"pagination"`
	Filtering   *FilteringConfig     `json:"filtering"`
	Sorting     *SortingConfig       `json:"sorting"`
	Joins       []JoinConfig         `json:"joins"`
	Cache       *CacheConfig         `json:"cache"`
	Security    *SecurityConfig      `json:"security"`
	Description string               `json:"description"`
	Tags        []string             `json:"tags"`
	Parameters  []ParameterInfo      `json:"parameters"`
	Responses   map[int]ResponseInfo `json:"responses"`
}

// ValidationRules represents validation rules for an endpoint
type ValidationRules struct {
	Required  []string
	MinLength map[string]int
	MaxLength map[string]int
	MinValue  map[string]float64
	MaxValue  map[string]float64
	Email     []string
	URL       []string
	UUID      []string
	Enum      map[string][]string
	Custom    map[string]string
}

// JoinConfig represents join configuration for relationships
type JoinConfig struct {
	Type       string                 `json:"type"` // "belongs_to", "has_many", "has_one", "many_to_many"
	Table      string                 `json:"table"`
	LocalKey   string                 `json:"local_key"`
	ForeignKey string                 `json:"foreign_key"`
	Alias      string                 `json:"alias"`
	Select     []string               `json:"select"`
	Where      map[string]interface{} `json:"where"`
	Order      string                 `json:"order"`
	Limit      int                    `json:"limit"`
}

// ParameterInfo represents API parameter information
type ParameterInfo struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
	Location    string `json:"location"` // "path", "query", "body", "header"
}

// ResponseInfo represents API response information
type ResponseInfo struct {
	Description string            `json:"description"`
	Schema      interface{}       `json:"schema"`
	Headers     map[string]string `json:"headers"`
}

// NewAPIGenerator creates a new API generator
func NewAPIGenerator(db *gorm.DB, logger *zap.Logger, config *GeneratorConfig) *APIGenerator {
	return &APIGenerator{
		db:             db,
		logger:         logger,
		config:         config,
		schemaAnalyzer: NewSchemaAnalyzer(db, logger),
		tables:         []*TableInfo{},
		endpoints:      []*GeneratedEndpoint{},
	}
}

// GenerateAll generates APIs for all discovered tables
func (g *APIGenerator) GenerateAll() error {
	g.logger.Info("Starting file-based API generation for all tables")

	// Discover all tables
	tableList, err := g.schemaAnalyzer.DiscoverTables()
	if err != nil {
		return fmt.Errorf("failed to discover tables: %w", err)
	}

	// Convert slice to map for easier access
	tables := make(map[string]*TableInfo)
	for _, table := range tableList {
		tables[table.Name] = table
	}

	g.tables = tableList

	// Create file generator
	fileGenerator := NewFileGenerator(g.db, g.logger, g.config)

	// Generate files for each table
	for tableName, tableInfo := range tables {
		if !g.config.ShouldGenerateTable(tableName) {
			g.logger.Info("Skipping table", zap.String("table", tableName))
			continue
		}

		g.logger.Info("Generating files for table", zap.String("table", tableName))

		if err := fileGenerator.GenerateFiles(tableName, tableInfo); err != nil {
			g.logger.Error("Failed to generate files", zap.String("table", tableName), zap.Error(err))
			continue
		}

		g.logger.Info("Generated files for table", zap.String("table", tableName))
	}

	// Generate main routes file that imports all generated routes
	if err := g.generateMainRoutesFile(tables); err != nil {
		g.logger.Error("Failed to generate main routes file", zap.Error(err))
	}

	g.logger.Info("File-based API generation completed", zap.Int("tables", len(tables)))
	return nil
}

// generateMainRoutesFile generates a main routes file that imports all generated routes
func (g *APIGenerator) generateMainRoutesFile(tables map[string]*TableInfo) error {
	routesFile := "internal/api/v1/generated_routes.go"

	tmpl := `package v1

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go-mobile-backend-template/pkg/config"
{{range .Tables}}
	"go-mobile-backend-template/internal/api/v1/{{.}}"
{{end}}
)

// RegisterGeneratedRoutes registers all auto-generated API routes
func RegisterGeneratedRoutes(router *gin.RouterGroup, db *gorm.DB, logger *zap.Logger, cfg *config.Config) {
{{range .Tables}}
	{{.}}.RegisterRoutes(router, db, logger, cfg)
{{end}}
}
`

	t, err := template.New("main_routes").Parse(tmpl)
	if err != nil {
		return err
	}

	file, err := os.Create(routesFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get table names for template
	var tableNames []string
	for tableName := range tables {
		if g.config.ShouldGenerateTable(tableName) {
			tableNames = append(tableNames, tableName)
		}
	}

	return t.Execute(file, map[string]interface{}{
		"Tables": tableNames,
	})
}

// GenerateTableEndpoints generates all endpoints for a specific table
func (g *APIGenerator) GenerateTableEndpoints(table *TableInfo) ([]*GeneratedEndpoint, error) {
	var endpoints []*GeneratedEndpoint
	tableConfig := g.config.GetTableConfig(table.Name)

	// Generate CRUD endpoints
	for _, endpointType := range tableConfig.Endpoints {
		endpoint, err := g.generateEndpoint(table, endpointType, tableConfig)
		if err != nil {
			g.logger.Error("Failed to generate endpoint",
				zap.String("table", table.Name),
				zap.String("endpoint", endpointType),
				zap.Error(err))
			continue
		}

		endpoints = append(endpoints, endpoint)
	}

	// Generate relationship endpoints
	for _, relationship := range tableConfig.Relationships {
		relEndpoints, err := g.generateRelationshipEndpoints(table, relationship, tableConfig)
		if err != nil {
			g.logger.Error("Failed to generate relationship endpoints",
				zap.String("table", table.Name),
				zap.String("relationship", relationship),
				zap.Error(err))
			continue
		}

		endpoints = append(endpoints, relEndpoints...)
	}

	return endpoints, nil
}

// generateEndpoint generates a specific endpoint for a table
func (g *APIGenerator) generateEndpoint(table *TableInfo, endpointType string, config *TableConfig) (*GeneratedEndpoint, error) {
	basePath := fmt.Sprintf("/api/v1/%s", strings.ToLower(table.Name))

	switch endpointType {
	case "list":
		return g.generateListEndpoint(table, basePath, config)
	case "create":
		return g.generateCreateEndpoint(table, basePath, config)
	case "get":
		return g.generateGetEndpoint(table, basePath, config)
	case "update":
		return g.generateUpdateEndpoint(table, basePath, config)
	case "delete":
		return g.generateDeleteEndpoint(table, basePath, config)
	case "bulk":
		return g.generateBulkEndpoint(table, basePath, config)
	case "search":
		return g.generateSearchEndpoint(table, basePath, config)
	case "stats":
		return g.generateStatsEndpoint(table, basePath, config)
	case "export":
		return g.generateExportEndpoint(table, basePath, config)
	default:
		return nil, fmt.Errorf("unknown endpoint type: %s", endpointType)
	}
}

// generateListEndpoint generates a list endpoint
func (g *APIGenerator) generateListEndpoint(table *TableInfo, basePath string, config *TableConfig) (*GeneratedEndpoint, error) {
	return &GeneratedEndpoint{
		Method:      "GET",
		Path:        basePath,
		Handler:     fmt.Sprintf("List%s", g.toCamelCase(table.Name)),
		Middleware:  g.getMiddlewareForEndpoint("list", config),
		Validation:  g.generateValidationRules(table, "list", config),
		Pagination:  config.Pagination,
		Filtering:   config.Filtering,
		Sorting:     config.Sorting,
		Joins:       g.generateJoins(table, config),
		Cache:       config.Caching,
		Security:    config.Security,
		Description: fmt.Sprintf("List all %s records", strings.ToLower(table.Name)),
		Tags:        []string{table.Name},
		Parameters:  g.generateListParameters(table, config),
		Responses: map[int]ResponseInfo{
			200: {
				Description: "List of records",
				Schema:      g.generateListResponseSchema(table),
			},
			400: {
				Description: "Bad request",
			},
			500: {
				Description: "Internal server error",
			},
		},
	}, nil
}

// generateCreateEndpoint generates a create endpoint
func (g *APIGenerator) generateCreateEndpoint(table *TableInfo, basePath string, config *TableConfig) (*GeneratedEndpoint, error) {
	return &GeneratedEndpoint{
		Method:      "POST",
		Path:        basePath,
		Handler:     fmt.Sprintf("Create%s", g.toCamelCase(table.Name)),
		Middleware:  g.getMiddlewareForEndpoint("create", config),
		Validation:  g.generateValidationRules(table, "create", config),
		Cache:       config.Caching,
		Security:    config.Security,
		Description: fmt.Sprintf("Create a new %s record", strings.ToLower(table.Name)),
		Tags:        []string{table.Name},
		Parameters:  g.generateCreateParameters(table, config),
		Responses: map[int]ResponseInfo{
			201: {
				Description: "Record created successfully",
				Schema:      g.generateSingleResponseSchema(table),
			},
			400: {
				Description: "Bad request",
			},
			409: {
				Description: "Conflict",
			},
			500: {
				Description: "Internal server error",
			},
		},
	}, nil
}

// generateGetEndpoint generates a get endpoint
func (g *APIGenerator) generateGetEndpoint(table *TableInfo, basePath string, config *TableConfig) (*GeneratedEndpoint, error) {
	return &GeneratedEndpoint{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/:id", basePath),
		Handler:     fmt.Sprintf("Get%s", g.toCamelCase(table.Name)),
		Middleware:  g.getMiddlewareForEndpoint("get", config),
		Validation:  g.generateValidationRules(table, "get", config),
		Joins:       g.generateJoins(table, config),
		Cache:       config.Caching,
		Security:    config.Security,
		Description: fmt.Sprintf("Get a %s record by ID", strings.ToLower(table.Name)),
		Tags:        []string{table.Name},
		Parameters:  g.generateGetParameters(table, config),
		Responses: map[int]ResponseInfo{
			200: {
				Description: "Record found",
				Schema:      g.generateSingleResponseSchema(table),
			},
			404: {
				Description: "Record not found",
			},
			500: {
				Description: "Internal server error",
			},
		},
	}, nil
}

// generateUpdateEndpoint generates an update endpoint
func (g *APIGenerator) generateUpdateEndpoint(table *TableInfo, basePath string, config *TableConfig) (*GeneratedEndpoint, error) {
	return &GeneratedEndpoint{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/:id", basePath),
		Handler:     fmt.Sprintf("Update%s", g.toCamelCase(table.Name)),
		Middleware:  g.getMiddlewareForEndpoint("update", config),
		Validation:  g.generateValidationRules(table, "update", config),
		Cache:       config.Caching,
		Security:    config.Security,
		Description: fmt.Sprintf("Update a %s record", strings.ToLower(table.Name)),
		Tags:        []string{table.Name},
		Parameters:  g.generateUpdateParameters(table, config),
		Responses: map[int]ResponseInfo{
			200: {
				Description: "Record updated successfully",
				Schema:      g.generateSingleResponseSchema(table),
			},
			400: {
				Description: "Bad request",
			},
			404: {
				Description: "Record not found",
			},
			500: {
				Description: "Internal server error",
			},
		},
	}, nil
}

// generateDeleteEndpoint generates a delete endpoint
func (g *APIGenerator) generateDeleteEndpoint(table *TableInfo, basePath string, config *TableConfig) (*GeneratedEndpoint, error) {
	return &GeneratedEndpoint{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/:id", basePath),
		Handler:     fmt.Sprintf("Delete%s", g.toCamelCase(table.Name)),
		Middleware:  g.getMiddlewareForEndpoint("delete", config),
		Validation:  g.generateValidationRules(table, "delete", config),
		Cache:       config.Caching,
		Security:    config.Security,
		Description: fmt.Sprintf("Delete a %s record", strings.ToLower(table.Name)),
		Tags:        []string{table.Name},
		Parameters:  g.generateDeleteParameters(table, config),
		Responses: map[int]ResponseInfo{
			204: {
				Description: "Record deleted successfully",
			},
			404: {
				Description: "Record not found",
			},
			500: {
				Description: "Internal server error",
			},
		},
	}, nil
}

// generateBulkEndpoint generates a bulk operations endpoint
func (g *APIGenerator) generateBulkEndpoint(table *TableInfo, basePath string, config *TableConfig) (*GeneratedEndpoint, error) {
	return &GeneratedEndpoint{
		Method:      "POST",
		Path:        fmt.Sprintf("%s/bulk", basePath),
		Handler:     fmt.Sprintf("Bulk%s", g.toCamelCase(table.Name)),
		Middleware:  g.getMiddlewareForEndpoint("bulk", config),
		Validation:  g.generateValidationRules(table, "bulk", config),
		Cache:       config.Caching,
		Security:    config.Security,
		Description: fmt.Sprintf("Bulk operations for %s records", strings.ToLower(table.Name)),
		Tags:        []string{table.Name},
		Parameters:  g.generateBulkParameters(table, config),
		Responses: map[int]ResponseInfo{
			200: {
				Description: "Bulk operation completed",
				Schema:      g.generateBulkResponseSchema(table),
			},
			400: {
				Description: "Bad request",
			},
			500: {
				Description: "Internal server error",
			},
		},
	}, nil
}

// generateSearchEndpoint generates a search endpoint
func (g *APIGenerator) generateSearchEndpoint(table *TableInfo, basePath string, config *TableConfig) (*GeneratedEndpoint, error) {
	return &GeneratedEndpoint{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/search", basePath),
		Handler:     fmt.Sprintf("Search%s", g.toCamelCase(table.Name)),
		Middleware:  g.getMiddlewareForEndpoint("search", config),
		Validation:  g.generateValidationRules(table, "search", config),
		Pagination:  config.Pagination,
		Filtering:   config.Filtering,
		Sorting:     config.Sorting,
		Cache:       config.Caching,
		Security:    config.Security,
		Description: fmt.Sprintf("Search %s records", strings.ToLower(table.Name)),
		Tags:        []string{table.Name},
		Parameters:  g.generateSearchParameters(table, config),
		Responses: map[int]ResponseInfo{
			200: {
				Description: "Search results",
				Schema:      g.generateListResponseSchema(table),
			},
			400: {
				Description: "Bad request",
			},
			500: {
				Description: "Internal server error",
			},
		},
	}, nil
}

// generateStatsEndpoint generates a stats endpoint
func (g *APIGenerator) generateStatsEndpoint(table *TableInfo, basePath string, config *TableConfig) (*GeneratedEndpoint, error) {
	return &GeneratedEndpoint{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/stats", basePath),
		Handler:     fmt.Sprintf("Stats%s", g.toCamelCase(table.Name)),
		Middleware:  g.getMiddlewareForEndpoint("stats", config),
		Cache:       config.Caching,
		Security:    config.Security,
		Description: fmt.Sprintf("Get statistics for %s records", strings.ToLower(table.Name)),
		Tags:        []string{table.Name},
		Parameters:  g.generateStatsParameters(table, config),
		Responses: map[int]ResponseInfo{
			200: {
				Description: "Statistics",
				Schema:      g.generateStatsResponseSchema(table),
			},
			500: {
				Description: "Internal server error",
			},
		},
	}, nil
}

// generateExportEndpoint generates an export endpoint
func (g *APIGenerator) generateExportEndpoint(table *TableInfo, basePath string, config *TableConfig) (*GeneratedEndpoint, error) {
	return &GeneratedEndpoint{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/export", basePath),
		Handler:     fmt.Sprintf("Export%s", g.toCamelCase(table.Name)),
		Middleware:  g.getMiddlewareForEndpoint("export", config),
		Validation:  g.generateValidationRules(table, "export", config),
		Filtering:   config.Filtering,
		Sorting:     config.Sorting,
		Cache:       config.Caching,
		Security:    config.Security,
		Description: fmt.Sprintf("Export %s records", strings.ToLower(table.Name)),
		Tags:        []string{table.Name},
		Parameters:  g.generateExportParameters(table, config),
		Responses: map[int]ResponseInfo{
			200: {
				Description: "Export file",
				Headers: map[string]string{
					"Content-Type": "application/octet-stream",
				},
			},
			400: {
				Description: "Bad request",
			},
			500: {
				Description: "Internal server error",
			},
		},
	}, nil
}

// Helper methods for generating endpoint components

func (g *APIGenerator) getMiddlewareForEndpoint(endpointType string, config *TableConfig) []string {
	middleware := []string{"Logger", "Recovery", "CORS", "SecurityHeaders"}

	if config.Security != nil {
		if config.Security.RBAC != nil {
			middleware = append(middleware, "RBAC")
		}
		if config.Security.RateLimit != nil {
			middleware = append(middleware, "RateLimit")
		}
		if config.Security.AuditLog {
			middleware = append(middleware, "AuditLog")
		}
	}

	// Add validation middleware for write operations
	if endpointType == "create" || endpointType == "update" || endpointType == "bulk" {
		middleware = append(middleware, "Validation")
	}

	// Add caching middleware for read operations
	if endpointType == "list" || endpointType == "get" || endpointType == "search" {
		middleware = append(middleware, "Cache")
	}

	return middleware
}

func (g *APIGenerator) generateValidationRules(table *TableInfo, endpointType string, config *TableConfig) *ValidationRules {
	rules := &ValidationRules{
		Required:  []string{},
		MinLength: make(map[string]int),
		MaxLength: make(map[string]int),
		MinValue:  make(map[string]float64),
		MaxValue:  make(map[string]float64),
		Email:     []string{},
		URL:       []string{},
		UUID:      []string{},
		Enum:      make(map[string][]string),
		Custom:    make(map[string]string),
	}

	if config.Validation == nil {
		return rules
	}

	// Apply global validation rules
	rules.Required = config.Validation.Required
	rules.MinLength = config.Validation.MinLength
	rules.MaxLength = config.Validation.MaxLength
	rules.MinValue = config.Validation.MinValue
	rules.MaxValue = config.Validation.MaxValue
	rules.Email = config.Validation.Email
	rules.URL = config.Validation.URL
	rules.UUID = config.Validation.UUID
	rules.Enum = config.Validation.Enum

	// Apply custom rules
	for field, rule := range config.Validation.CustomRules {
		rules.Custom[field] = rule
	}

	return rules
}

func (g *APIGenerator) generateJoins(table *TableInfo, config *TableConfig) []JoinConfig {
	var joins []JoinConfig

	// Auto-detect relationships from foreign keys
	for _, fk := range table.ForeignKeys {
		join := JoinConfig{
			Type:       "belongs_to",
			Table:      fk.RefTable,
			LocalKey:   fk.Column,
			ForeignKey: fk.RefColumn,
			Alias:      fk.RefTable,
		}
		joins = append(joins, join)
	}

	return joins
}

func (g *APIGenerator) generateListParameters(table *TableInfo, config *TableConfig) []ParameterInfo {
	params := []ParameterInfo{
		{Name: "page", Type: "integer", Required: false, Description: "Page number", Location: "query"},
		{Name: "limit", Type: "integer", Required: false, Description: "Number of records per page", Location: "query"},
		{Name: "sort", Type: "string", Required: false, Description: "Sort field and direction", Location: "query"},
	}

	// Add filtering parameters
	if config.Filtering != nil {
		for _, field := range config.Filtering.AllowedFields {
			params = append(params, ParameterInfo{
				Name:        field,
				Type:        "string",
				Required:    false,
				Description: fmt.Sprintf("Filter by %s", field),
				Location:    "query",
			})
		}
	}

	return params
}

func (g *APIGenerator) generateCreateParameters(table *TableInfo, config *TableConfig) []ParameterInfo {
	return []ParameterInfo{
		{Name: "body", Type: "object", Required: true, Description: "Record data", Location: "body"},
	}
}

func (g *APIGenerator) generateGetParameters(table *TableInfo, config *TableConfig) []ParameterInfo {
	return []ParameterInfo{
		{Name: "id", Type: "string", Required: true, Description: "Record ID", Location: "path"},
	}
}

func (g *APIGenerator) generateUpdateParameters(table *TableInfo, config *TableConfig) []ParameterInfo {
	return []ParameterInfo{
		{Name: "id", Type: "string", Required: true, Description: "Record ID", Location: "path"},
		{Name: "body", Type: "object", Required: true, Description: "Updated record data", Location: "body"},
	}
}

func (g *APIGenerator) generateDeleteParameters(table *TableInfo, config *TableConfig) []ParameterInfo {
	return []ParameterInfo{
		{Name: "id", Type: "string", Required: true, Description: "Record ID", Location: "path"},
	}
}

func (g *APIGenerator) generateBulkParameters(table *TableInfo, config *TableConfig) []ParameterInfo {
	return []ParameterInfo{
		{Name: "body", Type: "object", Required: true, Description: "Bulk operation data", Location: "body"},
	}
}

func (g *APIGenerator) generateSearchParameters(table *TableInfo, config *TableConfig) []ParameterInfo {
	params := []ParameterInfo{
		{Name: "q", Type: "string", Required: true, Description: "Search query", Location: "query"},
		{Name: "page", Type: "integer", Required: false, Description: "Page number", Location: "query"},
		{Name: "limit", Type: "integer", Required: false, Description: "Number of records per page", Location: "query"},
	}

	return params
}

func (g *APIGenerator) generateStatsParameters(table *TableInfo, config *TableConfig) []ParameterInfo {
	return []ParameterInfo{}
}

func (g *APIGenerator) generateExportParameters(table *TableInfo, config *TableConfig) []ParameterInfo {
	return []ParameterInfo{
		{Name: "format", Type: "string", Required: false, Description: "Export format (csv, json, xlsx)", Location: "query"},
		{Name: "fields", Type: "string", Required: false, Description: "Comma-separated list of fields to export", Location: "query"},
	}
}

func (g *APIGenerator) generateListResponseSchema(table *TableInfo) interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"data": map[string]interface{}{
				"type":  "array",
				"items": g.generateRecordSchema(table),
			},
			"pagination": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"page":        map[string]interface{}{"type": "integer"},
					"limit":       map[string]interface{}{"type": "integer"},
					"total":       map[string]interface{}{"type": "integer"},
					"total_pages": map[string]interface{}{"type": "integer"},
					"has_next":    map[string]interface{}{"type": "boolean"},
					"has_prev":    map[string]interface{}{"type": "boolean"},
				},
			},
		},
	}
}

func (g *APIGenerator) generateSingleResponseSchema(table *TableInfo) interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"data": g.generateRecordSchema(table),
		},
	}
}

func (g *APIGenerator) generateRecordSchema(table *TableInfo) interface{} {
	properties := make(map[string]interface{})

	for _, column := range table.Columns {
		properties[column.Name] = map[string]interface{}{
			"type":        column.TSType,
			"description": column.Comment,
		}
	}

	return map[string]interface{}{
		"type":       "object",
		"properties": properties,
	}
}

func (g *APIGenerator) generateBulkResponseSchema(table *TableInfo) interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"created": map[string]interface{}{"type": "integer"},
			"updated": map[string]interface{}{"type": "integer"},
			"deleted": map[string]interface{}{"type": "integer"},
			"errors":  map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
		},
	}
}

func (g *APIGenerator) generateStatsResponseSchema(table *TableInfo) interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"total":    map[string]interface{}{"type": "integer"},
			"active":   map[string]interface{}{"type": "integer"},
			"inactive": map[string]interface{}{"type": "integer"},
		},
	}
}

func (g *APIGenerator) generateRelationshipEndpoints(table *TableInfo, relationship string, config *TableConfig) ([]*GeneratedEndpoint, error) {
	// This would generate relationship-specific endpoints
	// For now, return empty slice
	return []*GeneratedEndpoint{}, nil
}

// Utility methods

func (g *APIGenerator) toCamelCase(str string) string {
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

// GetGeneratedEndpoints returns all generated endpoints
func (g *APIGenerator) GetGeneratedEndpoints() []*GeneratedEndpoint {
	return g.endpoints
}

// GetTables returns all discovered tables
func (g *APIGenerator) GetTables() []*TableInfo {
	return g.tables
}

// GetEndpointsForTable returns endpoints for a specific table
func (g *APIGenerator) GetEndpointsForTable(tableName string) []*GeneratedEndpoint {
	var endpoints []*GeneratedEndpoint

	for _, endpoint := range g.endpoints {
		if strings.Contains(endpoint.Path, "/"+strings.ToLower(tableName)) {
			endpoints = append(endpoints, endpoint)
		}
	}

	return endpoints
}
