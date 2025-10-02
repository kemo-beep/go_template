package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// TypeScriptGenerator generates TypeScript types and API clients
type TypeScriptGenerator struct {
	db     *gorm.DB
	logger *zap.Logger
	config *GeneratorConfig
}

// NewTypeScriptGenerator creates a new TypeScript generator
func NewTypeScriptGenerator(db *gorm.DB, logger *zap.Logger, config *GeneratorConfig) *TypeScriptGenerator {
	return &TypeScriptGenerator{
		db:     db,
		logger: logger,
		config: config,
	}
}

// GenerateAll generates all TypeScript types and API clients
func (tg *TypeScriptGenerator) GenerateAll() error {
	if !tg.config.GenerateTypeScript {
		tg.logger.Info("TypeScript generation is disabled")
		return nil
	}

	tg.logger.Info("Starting TypeScript type generation...")

	// Discover all tables
	tableList, err := tg.discoverTables()
	if err != nil {
		return fmt.Errorf("failed to discover tables: %w", err)
	}

	// Generate types for each table
	for _, table := range tableList {
		if !tg.config.ShouldGenerateTable(table.Name) {
			continue
		}

		if err := tg.generateTableTypes(table); err != nil {
			tg.logger.Error("Failed to generate types for table",
				zap.String("table", table.Name),
				zap.Error(err))
			continue
		}
	}

	// Generate API client
	if err := tg.generateAPIClient(); err != nil {
		return fmt.Errorf("failed to generate API client: %w", err)
	}

	// Generate index file
	if err := tg.generateIndexFile(); err != nil {
		return fmt.Errorf("failed to generate index file: %w", err)
	}

	tg.logger.Info("TypeScript type generation completed successfully")
	return nil
}

// discoverTables discovers all tables in the database
func (tg *TypeScriptGenerator) discoverTables() ([]*TableInfo, error) {
	// Reuse the existing schema analyzer
	analyzer := NewSchemaAnalyzer(tg.db, tg.logger)
	return analyzer.DiscoverTables()
}

// generateTableTypes generates TypeScript types for a specific table
func (tg *TypeScriptGenerator) generateTableTypes(table *TableInfo) error {
	tg.logger.Info("Generating TypeScript types for table", zap.String("table", table.Name))

	// Create output directory
	outputDir := filepath.Join("frontend", "lib", "types", "generated")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate main type file
	if err := tg.generateMainTypeFile(table, outputDir); err != nil {
		return fmt.Errorf("failed to generate main type file: %w", err)
	}

	// Generate API types file
	if err := tg.generateAPITypesFile(table, outputDir); err != nil {
		return fmt.Errorf("failed to generate API types file: %w", err)
	}

	return nil
}

// generateMainTypeFile generates the main TypeScript type file for a table
func (tg *TypeScriptGenerator) generateMainTypeFile(table *TableInfo, outputDir string) error {
	filename := fmt.Sprintf("%s.ts", strings.ToLower(table.Name))
	filepath := filepath.Join(outputDir, filename)

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	tmpl := `// Auto-generated TypeScript types for {{.TableName}}
// Generated on: {{.GeneratedAt}}

export interface {{.StructName}} {
{{range .Columns}}
  {{.FieldName}}: {{.TypeScriptType}}{{if .IsOptional}}?{{end}};{{if .Comment}} // {{.Comment}}{{end}}
{{end}}
}

export interface {{.StructName}}CreateRequest {
{{range .CreateColumns}}
  {{.FieldName}}: {{.TypeScriptType}}{{if .IsOptional}}?{{end}};{{if .Comment}} // {{.Comment}}{{end}}
{{end}}
}

export interface {{.StructName}}UpdateRequest {
{{range .UpdateColumns}}
  {{.FieldName}}?: {{.TypeScriptType}};{{if .Comment}} // {{.Comment}}{{end}}
{{end}}
}

export interface {{.StructName}}Response {
{{range .Columns}}
  {{.FieldName}}: {{.TypeScriptType}}{{if .IsOptional}}?{{end}};{{if .Comment}} // {{.Comment}}{{end}}
{{end}}
}

export interface {{.StructName}}PaginationInfo {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
  has_next: boolean;
  has_prev: boolean;
}

export interface {{.StructName}}PaginationResponse {
  data: {{.StructName}}Response[];
  pagination: {{.StructName}}PaginationInfo;
}
`

	t := template.Must(template.New("types").Parse(tmpl))

	// Prepare template data
	columns := tg.generateColumnInfo(table.Columns)
	createColumns := tg.generateCreateColumns(table.Columns)
	updateColumns := tg.generateUpdateColumns(table.Columns)

	data := map[string]interface{}{
		"TableName":     table.Name,
		"StructName":    tg.toPascalCase(table.Name),
		"GeneratedAt":   "2025-10-02T10:00:00Z", // You could use time.Now()
		"Columns":       columns,
		"CreateColumns": createColumns,
		"UpdateColumns": updateColumns,
	}

	return t.Execute(file, data)
}

// generateAPITypesFile generates API-specific types for a table
func (tg *TypeScriptGenerator) generateAPITypesFile(table *TableInfo, outputDir string) error {
	filename := fmt.Sprintf("%s-api.ts", strings.ToLower(table.Name))
	filepath := filepath.Join(outputDir, filename)

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	tmpl := `// Auto-generated API types for {{.TableName}}
// Generated on: {{.GeneratedAt}}

import { {{.StructName}}, {{.StructName}}CreateRequest, {{.StructName}}UpdateRequest, {{.StructName}}Response, {{.StructName}}PaginationResponse } from './{{.LowerName}}';

export interface {{.StructName}}API {
  // Get all {{.LowerName}}s with pagination
  getAll(params?: {
    page?: number;
    limit?: number;
    search?: string;
    sort_by?: string;
    sort_order?: 'asc' | 'desc';
  }): Promise<{{.StructName}}PaginationResponse>;

  // Get {{.LowerName}} by ID
  getById(id: number): Promise<{{.StructName}}Response>;

  // Create new {{.LowerName}}
  create(data: {{.StructName}}CreateRequest): Promise<{{.StructName}}Response>;

  // Update {{.LowerName}} by ID
  update(id: number, data: {{.StructName}}UpdateRequest): Promise<{{.StructName}}Response>;

  // Delete {{.LowerName}} by ID
  delete(id: number): Promise<void>;

  // Bulk operations
  bulkCreate(data: {{.StructName}}CreateRequest[]): Promise<{{.StructName}}Response[]>;
  bulkUpdate(updates: { id: number; data: {{.StructName}}UpdateRequest }[]): Promise<{{.StructName}}Response[]>;
  bulkDelete(ids: number[]): Promise<void>;
}

export const {{.LowerName}}API: {{.StructName}}API = {
  async getAll(params = {}) {
    const searchParams = new URLSearchParams();
    if (params.page) searchParams.set('page', params.page.toString());
    if (params.limit) searchParams.set('limit', params.limit.toString());
    if (params.search) searchParams.set('search', params.search);
    if (params.sort_by) searchParams.set('sort_by', params.sort_by);
    if (params.sort_order) searchParams.set('sort_order', params.sort_order);

    const response = await fetch('/api/v1/{{.LowerName}}?' + searchParams, {
      method: 'GET',
      headers: {
        'Authorization': 'Bearer ' + localStorage.getItem('access_token'),
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error('Failed to fetch {{.LowerName}}s: ' + response.statusText);
    }

    return response.json();
  },

  async getById(id: number) {
    const response = await fetch('/api/v1/{{.LowerName}}/' + id, {
      method: 'GET',
      headers: {
        'Authorization': 'Bearer ' + localStorage.getItem('access_token'),
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error('Failed to fetch {{.LowerName}}: ' + response.statusText);
    }

    return response.json();
  },

  async create(data: {{.StructName}}CreateRequest) {
    const response = await fetch('/api/v1/{{.LowerName}}', {
      method: 'POST',
      headers: {
        'Authorization': 'Bearer ' + localStorage.getItem('access_token'),
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      throw new Error('Failed to create {{.LowerName}}: ' + response.statusText);
    }

    return response.json();
  },

  async update(id: number, data: {{.StructName}}UpdateRequest) {
    const response = await fetch('/api/v1/{{.LowerName}}/' + id, {
      method: 'PUT',
      headers: {
        'Authorization': 'Bearer ' + localStorage.getItem('access_token'),
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      throw new Error('Failed to update {{.LowerName}}: ' + response.statusText);
    }

    return response.json();
  },

  async delete(id: number) {
    const response = await fetch('/api/v1/{{.LowerName}}/' + id, {
      method: 'DELETE',
      headers: {
        'Authorization': 'Bearer ' + localStorage.getItem('access_token'),
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error('Failed to delete {{.LowerName}}: ' + response.statusText);
    }
  },

  async bulkCreate(data: {{.StructName}}CreateRequest[]) {
    const response = await fetch('/api/v1/{{.LowerName}}/bulk', {
      method: 'POST',
      headers: {
        'Authorization': 'Bearer ' + localStorage.getItem('access_token'),
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ data }),
    });

    if (!response.ok) {
      throw new Error('Failed to bulk create {{.LowerName}}s: ' + response.statusText);
    }

    return response.json();
  },

  async bulkUpdate(updates: { id: number; data: {{.StructName}}UpdateRequest }[]) {
    const response = await fetch('/api/v1/{{.LowerName}}/bulk', {
      method: 'PUT',
      headers: {
        'Authorization': 'Bearer ' + localStorage.getItem('access_token'),
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ updates }),
    });

    if (!response.ok) {
      throw new Error('Failed to bulk update {{.LowerName}}s: ' + response.statusText);
    }

    return response.json();
  },

  async bulkDelete(ids: number[]) {
    const response = await fetch('/api/v1/{{.LowerName}}/bulk', {
      method: 'DELETE',
      headers: {
        'Authorization': 'Bearer ' + localStorage.getItem('access_token'),
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ ids }),
    });

    if (!response.ok) {
      throw new Error('Failed to bulk delete {{.LowerName}}s: ' + response.statusText);
    }
  },
};
`

	t := template.Must(template.New("api").Parse(tmpl))

	data := map[string]interface{}{
		"TableName":   table.Name,
		"StructName":  tg.toPascalCase(table.Name),
		"LowerName":   strings.ToLower(table.Name),
		"GeneratedAt": "2025-10-02T10:00:00Z",
	}

	return t.Execute(file, data)
}

// generateAPIClient generates a unified API client
func (tg *TypeScriptGenerator) generateAPIClient() error {
	tg.logger.Info("Generating unified API client...")

	// Create output directory
	outputDir := filepath.Join("frontend", "lib", "api")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	filepath := filepath.Join(outputDir, "generated-client.ts")

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	tmpl := `// Auto-generated unified API client
// Generated on: {{.GeneratedAt}}

{{range .Tables}}
import { {{.StructName}}API } from '../types/generated/{{.LowerName}}-api';
{{end}}

export class GeneratedAPIClient {
{{range .Tables}}
  public {{.LowerName}}: {{.StructName}}API;
{{end}}

  constructor() {
{{range .Tables}}
    this.{{.LowerName}} = {{.LowerName}}API;
{{end}}
  }
}

// Export individual APIs for convenience
{{range .Tables}}
export { {{.LowerName}}API } from '../types/generated/{{.LowerName}}-api';
{{end}}

// Export all types
{{range .Tables}}
export * from '../types/generated/{{.LowerName}}';
{{end}}

// Create a default instance
export const api = new GeneratedAPIClient();
export default api;
`

	t := template.Must(template.New("client").Parse(tmpl))

	// Get all tables
	tables, err := tg.discoverTables()
	if err != nil {
		return fmt.Errorf("failed to discover tables: %w", err)
	}

	var tableData []map[string]interface{}
	for _, table := range tables {
		if !tg.config.ShouldGenerateTable(table.Name) {
			continue
		}
		tableData = append(tableData, map[string]interface{}{
			"TableName":  table.Name,
			"StructName": tg.toPascalCase(table.Name),
			"LowerName":  strings.ToLower(table.Name),
		})
	}

	data := map[string]interface{}{
		"GeneratedAt": "2025-10-02T10:00:00Z",
		"Tables":      tableData,
	}

	return t.Execute(file, data)
}

// generateIndexFile generates an index file that exports everything
func (tg *TypeScriptGenerator) generateIndexFile() error {
	tg.logger.Info("Generating index file...")

	filepath := filepath.Join("frontend", "lib", "types", "generated", "index.ts")

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	tmpl := `// Auto-generated index file
// Generated on: {{.GeneratedAt}}

{{range .Tables}}
export * from './{{.LowerName}}';
export * from './{{.LowerName}}-api';
{{end}}
`

	t := template.Must(template.New("index").Parse(tmpl))

	// Get all tables
	tables, err := tg.discoverTables()
	if err != nil {
		return fmt.Errorf("failed to discover tables: %w", err)
	}

	var tableData []map[string]interface{}
	for _, table := range tables {
		if !tg.config.ShouldGenerateTable(table.Name) {
			continue
		}
		tableData = append(tableData, map[string]interface{}{
			"LowerName": strings.ToLower(table.Name),
		})
	}

	data := map[string]interface{}{
		"GeneratedAt": "2025-10-02T10:00:00Z",
		"Tables":      tableData,
	}

	return t.Execute(file, data)
}

// Helper methods (reuse from file_generator.go)
func (tg *TypeScriptGenerator) toPascalCase(s string) string {
	words := strings.Split(s, "_")
	for i, word := range words {
		words[i] = strings.Title(word)
	}
	return strings.Join(words, "")
}

func (tg *TypeScriptGenerator) generateColumnInfo(columns []ColumnInfo) []map[string]interface{} {
	var result []map[string]interface{}
	for _, col := range columns {
		result = append(result, map[string]interface{}{
			"FieldName":      tg.toPascalCase(col.Name),
			"TypeScriptType": tg.getTypeScriptType(col),
			"IsOptional":     col.IsNullable,
			"Comment":        col.Comment,
		})
	}
	return result
}

func (tg *TypeScriptGenerator) generateCreateColumns(columns []ColumnInfo) []map[string]interface{} {
	var result []map[string]interface{}
	for _, col := range columns {
		// Skip auto-generated fields for create requests
		if col.IsPrimaryKey || col.Name == "created_at" || col.Name == "updated_at" {
			continue
		}
		result = append(result, map[string]interface{}{
			"FieldName":      tg.toPascalCase(col.Name),
			"TypeScriptType": tg.getTypeScriptType(col),
			"IsOptional":     col.IsNullable,
			"Comment":        col.Comment,
		})
	}
	return result
}

func (tg *TypeScriptGenerator) generateUpdateColumns(columns []ColumnInfo) []map[string]interface{} {
	var result []map[string]interface{}
	for _, col := range columns {
		// Skip auto-generated fields for update requests
		if col.IsPrimaryKey || col.Name == "created_at" || col.Name == "updated_at" {
			continue
		}
		result = append(result, map[string]interface{}{
			"FieldName":      tg.toPascalCase(col.Name),
			"TypeScriptType": tg.getTypeScriptType(col),
			"Comment":        col.Comment,
		})
	}
	return result
}

func (tg *TypeScriptGenerator) getTypeScriptType(col ColumnInfo) string {
	switch col.Type {
	case "integer", "bigint", "smallint":
		return "number"
	case "real", "double precision", "numeric", "decimal":
		return "number"
	case "boolean":
		return "boolean"
	case "timestamp", "timestamptz", "date", "time":
		return "string" // ISO string format
	case "json", "jsonb":
		return "any"
	case "uuid":
		return "string"
	default:
		return "string"
	}
}
