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

// FileGenerator handles file-based code generation
type FileGenerator struct {
	db     *gorm.DB
	logger *zap.Logger
	config *GeneratorConfig
}

// NewFileGenerator creates a new file generator
func NewFileGenerator(db *gorm.DB, logger *zap.Logger, config *GeneratorConfig) *FileGenerator {
	return &FileGenerator{
		db:     db,
		logger: logger,
		config: config,
	}
}

// GenerateFiles generates all files for a table
func (fg *FileGenerator) GenerateFiles(tableName string, tableInfo *TableInfo) error {
	// Create directory structure
	apiDir := filepath.Join("internal", "api", "v1", tableName)
	modelDir := filepath.Join("internal", "db", "repository", "generated")

	if err := os.MkdirAll(apiDir, 0755); err != nil {
		return fmt.Errorf("failed to create API directory: %w", err)
	}

	if err := os.MkdirAll(modelDir, 0755); err != nil {
		return fmt.Errorf("failed to create model directory: %w", err)
	}

	// Generate files
	if err := fg.generateModel(tableName, tableInfo, modelDir); err != nil {
		return fmt.Errorf("failed to generate model: %w", err)
	}

	if err := fg.generateRepository(tableName, tableInfo, modelDir); err != nil {
		return fmt.Errorf("failed to generate repository: %w", err)
	}

	if err := fg.generateRequestStructs(tableName, tableInfo, apiDir); err != nil {
		return fmt.Errorf("failed to generate request structs: %w", err)
	}

	if err := fg.generateHandler(tableName, tableInfo, apiDir); err != nil {
		return fmt.Errorf("failed to generate handler: %w", err)
	}

	if err := fg.generateRoutes(tableName, tableInfo, apiDir); err != nil {
		return fmt.Errorf("failed to generate routes: %w", err)
	}

	return nil
}

// generateModel generates the GORM model file
func (fg *FileGenerator) generateModel(tableName string, tableInfo *TableInfo, modelDir string) error {
	modelFile := filepath.Join(modelDir, fmt.Sprintf("%s.go", tableName))

	tmpl := `package generated
{{if .HasTimeFields}}
import "time"
{{end}}

// {{.StructName}} represents the {{.TableName}} table
type {{.StructName}} struct {
{{range .Columns}}
	{{.FieldName}} {{.GoType}} ` + "`" + `gorm:"{{.GormTag}}" json:"{{.JSONTag}}"` + "`" + `{{if .Comment}} // {{.Comment}}{{end}}
{{end}}
}

// TableName returns the table name for {{.StructName}}
func ({{.StructName}}) TableName() string {
	return "{{.TableName}}"
}
`

	t, err := template.New("model").Parse(tmpl)
	if err != nil {
		return err
	}

	file, err := os.Create(modelFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Check if any column uses time.Time
	columns := fg.generateColumnInfo(tableInfo.Columns)
	hasTimeFields := false
	for _, col := range columns {
		if col["GoType"] == "time.Time" {
			hasTimeFields = true
			break
		}
	}

	return t.Execute(file, map[string]interface{}{
		"StructName":    fg.toPascalCase(tableName),
		"TableName":     tableName,
		"Columns":       columns,
		"HasTimeFields": hasTimeFields,
	})
}

// generateRepository generates the repository file
func (fg *FileGenerator) generateRepository(tableName string, tableInfo *TableInfo, modelDir string) error {
	repoFile := filepath.Join(modelDir, fmt.Sprintf("%s_repository.go", tableName))

	tmpl := `package generated

import (
	"context"
	"gorm.io/gorm"
)

// {{.StructName}}Repository interface for {{.TableName}} operations
type {{.StructName}}Repository interface {
	Create(ctx context.Context, {{.LowerName}} *{{.StructName}}) error
	GetByID(ctx context.Context, id uint) (*{{.StructName}}, error)
	GetAll(ctx context.Context, limit, offset int) ([]{{.StructName}}, int64, error)
	Update(ctx context.Context, {{.LowerName}} *{{.StructName}}) error
	Delete(ctx context.Context, id uint) error
}

// {{.LowerName}}Repository implements {{.StructName}}Repository
type {{.LowerName}}Repository struct {
	db *gorm.DB
}

// New{{.StructName}}Repository creates a new {{.StructName}}Repository
func New{{.StructName}}Repository(db *gorm.DB) {{.StructName}}Repository {
	return &{{.LowerName}}Repository{db: db}
}

// Create creates a new {{.LowerName}}
func (r *{{.LowerName}}Repository) Create(ctx context.Context, {{.LowerName}} *{{.StructName}}) error {
	return r.db.WithContext(ctx).Create({{.LowerName}}).Error
}

// GetByID gets a {{.LowerName}} by ID
func (r *{{.LowerName}}Repository) GetByID(ctx context.Context, id uint) (*{{.StructName}}, error) {
	var {{.LowerName}} {{.StructName}}
	err := r.db.WithContext(ctx).First(&{{.LowerName}}, id).Error
	if err != nil {
		return nil, err
	}
	return &{{.LowerName}}, nil
}

// GetAll gets all {{.LowerName}}s with pagination
func (r *{{.LowerName}}Repository) GetAll(ctx context.Context, limit, offset int) ([]{{.StructName}}, int64, error) {
	var {{.LowerName}}s []{{.StructName}}
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&{{.StructName}}{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get data with pagination
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&{{.LowerName}}s).Error
	return {{.LowerName}}s, total, err
}

// Update updates a {{.LowerName}}
func (r *{{.LowerName}}Repository) Update(ctx context.Context, {{.LowerName}} *{{.StructName}}) error {
	return r.db.WithContext(ctx).Save({{.LowerName}}).Error
}

// Delete deletes a {{.LowerName}} by ID
func (r *{{.LowerName}}Repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&{{.StructName}}{}, id).Error
}
`

	t, err := template.New("repository").Parse(tmpl)
	if err != nil {
		return err
	}

	file, err := os.Create(repoFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return t.Execute(file, map[string]interface{}{
		"StructName": fg.toPascalCase(tableName),
		"LowerName":  fg.toCamelCase(tableName),
		"TableName":  tableName,
	})
}

// generateRequestStructs generates request/response structs
func (fg *FileGenerator) generateRequestStructs(tableName string, tableInfo *TableInfo, apiDir string) error {
	requestFile := filepath.Join(apiDir, "request.go")

	tmpl := `package {{.PackageName}}
{{if .HasTimeFields}}
import "time"
{{end}}

// {{.StructName}}Response represents {{.TableName}} response
type {{.StructName}}Response struct {
{{range .Columns}}
	{{.FieldName}} {{.GoType}} ` + "`" + `json:"{{.JSONTag}}"` + "`" + `{{if .Comment}} // {{.Comment}}{{end}}
{{end}}
}

// {{.StructName}}CreateRequest represents create {{.TableName}} request
type {{.StructName}}CreateRequest struct {
{{range .CreateColumns}}
	{{.FieldName}} {{.GoType}} ` + "`" + `json:"{{.JSONTag}}" binding:"{{.BindingTag}}"` + "`" + `{{if .Comment}} // {{.Comment}}{{end}}
{{end}}
}

// {{.StructName}}UpdateRequest represents update {{.TableName}} request
type {{.StructName}}UpdateRequest struct {
{{range .UpdateColumns}}
	{{.FieldName}} {{.GoType}} ` + "`" + `json:"{{.JSONTag}}" binding:"{{.BindingTag}}"` + "`" + `{{if .Comment}} // {{.Comment}}{{end}}
{{end}}
}

// PaginationResponse represents pagination response
type PaginationResponse struct {
	Data       []{{.StructName}}Response ` + "`" + `json:"data"` + "`" + `
	Pagination PaginationInfo            ` + "`" + `json:"pagination"` + "`" + `
}

// PaginationInfo represents pagination information
type PaginationInfo struct {
	Page       int   ` + "`" + `json:"page"` + "`" + `
	Limit      int   ` + "`" + `json:"limit"` + "`" + `
	Total      int64 ` + "`" + `json:"total"` + "`" + `
	TotalPages int   ` + "`" + `json:"total_pages"` + "`" + `
	HasNext    bool  ` + "`" + `json:"has_next"` + "`" + `
	HasPrev    bool  ` + "`" + `json:"has_prev"` + "`" + `
}
`

	t, err := template.New("request").Parse(tmpl)
	if err != nil {
		return err
	}

	file, err := os.Create(requestFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Check if any column uses time.Time
	columns := fg.generateColumnInfo(tableInfo.Columns)
	hasTimeFields := false
	for _, col := range columns {
		if col["GoType"] == "time.Time" {
			hasTimeFields = true
			break
		}
	}

	return t.Execute(file, map[string]interface{}{
		"PackageName":   tableName,
		"StructName":    fg.toPascalCase(tableName),
		"TableName":     tableName,
		"Columns":       columns,
		"CreateColumns": fg.generateCreateColumns(tableInfo.Columns),
		"UpdateColumns": fg.generateUpdateColumns(tableInfo.Columns),
		"HasTimeFields": hasTimeFields,
	})
}

// generateHandler generates the handler file with Swagger annotations
func (fg *FileGenerator) generateHandler(tableName string, tableInfo *TableInfo, apiDir string) error {
	handlerFile := filepath.Join(apiDir, "handler.go")

	tmpl := `package {{.PackageName}}

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go-mobile-backend-template/internal/db/repository/generated"
	"go-mobile-backend-template/internal/utils"
)

// Handler handles {{.TableName}} requests
type Handler struct {
	{{.LowerName}}Repo generated.{{.StructName}}Repository
	logger             *zap.Logger
}

// NewHandler creates a new {{.TableName}} handler
func NewHandler(db *gorm.DB, logger *zap.Logger) *Handler {
	return &Handler{
		{{.LowerName}}Repo: generated.New{{.StructName}}Repository(db),
		logger:             logger,
	}
}

// Create{{.StructName}} creates a new {{.TableName}}
// @Summary Create {{.TableName}}
// @Description Create a new {{.TableName}} record
// @Tags {{.TableName}}
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body {{.StructName}}CreateRequest true "Create {{.TableName}} request"
// @Success 201 {object} {{.StructName}}Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /{{.TableName}} [post]
func (h *Handler) Create{{.StructName}}(c *gin.Context) {
	var req {{.StructName}}CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	{{.LowerName}} := &generated.{{.StructName}}{
{{range .CreateColumns}}
		{{.FieldName}}: req.{{.FieldName}},
{{end}}
	}

	ctx := context.Background()
	if err := h.{{.LowerName}}Repo.Create(ctx, {{.LowerName}}); err != nil {
		h.logger.Error("Failed to create {{.TableName}}", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create {{.TableName}}")
		return
	}

	response := {{.StructName}}Response{
{{range .Columns}}
		{{.FieldName}}: {{$.LowerName}}.{{.FieldName}},
{{end}}
	}

	utils.SuccessResponse(c, http.StatusCreated, "{{.TableName}} created successfully", response)
}

// Get{{.StructName}} gets a {{.TableName}} by ID
// @Summary Get {{.TableName}}
// @Description Get a {{.TableName}} by ID
// @Tags {{.TableName}}
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "{{.TableName}} ID"
// @Success 200 {object} {{.StructName}}Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /{{.TableName}}/{id} [get]
func (h *Handler) Get{{.StructName}}(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	ctx := context.Background()
	{{.LowerName}}, err := h.{{.LowerName}}Repo.GetByID(ctx, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "{{.TableName}} not found")
			return
		}
		h.logger.Error("Failed to get {{.TableName}}", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get {{.TableName}}")
		return
	}

	response := {{.StructName}}Response{
{{range .Columns}}
		{{.FieldName}}: {{$.LowerName}}.{{.FieldName}},
{{end}}
	}

	utils.SuccessResponse(c, http.StatusOK, "{{.TableName}} retrieved successfully", response)
}

// GetAll{{.StructName}}s gets all {{.TableName}}s with pagination
// @Summary Get all {{.TableName}}s
// @Description Get all {{.TableName}}s with pagination
// @Tags {{.TableName}}
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(20)
// @Success 200 {object} PaginationResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /{{.TableName}} [get]
func (h *Handler) GetAll{{.StructName}}s(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	ctx := context.Background()
	{{.LowerName}}s, total, err := h.{{.LowerName}}Repo.GetAll(ctx, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get {{.TableName}}s", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get {{.TableName}}s")
		return
	}

	var responses []{{.StructName}}Response
	for _, {{.LowerName}} := range {{.LowerName}}s {
		responses = append(responses, {{.StructName}}Response{
{{range .Columns}}
			{{.FieldName}}: {{$.LowerName}}.{{.FieldName}},
{{end}}
		})
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	pagination := PaginationResponse{
		Data: responses,
		Pagination: PaginationInfo{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
	}

	utils.SuccessResponse(c, http.StatusOK, "{{.TableName}}s retrieved successfully", pagination)
}

// Update{{.StructName}} updates a {{.TableName}}
// @Summary Update {{.TableName}}
// @Description Update a {{.TableName}} by ID
// @Tags {{.TableName}}
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "{{.TableName}} ID"
// @Param request body {{.StructName}}UpdateRequest true "Update {{.TableName}} request"
// @Success 200 {object} {{.StructName}}Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /{{.TableName}}/{id} [put]
func (h *Handler) Update{{.StructName}}(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req {{.StructName}}UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx := context.Background()
	{{.LowerName}}, err := h.{{.LowerName}}Repo.GetByID(ctx, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "{{.TableName}} not found")
			return
		}
		h.logger.Error("Failed to get {{.TableName}}", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get {{.TableName}}")
		return
	}

{{range .UpdateColumns}}
	{{if ne .FieldName "ID"}}
	{{$.LowerName}}.{{.FieldName}} = req.{{.FieldName}}
	{{end}}
{{end}}

	if err := h.{{.LowerName}}Repo.Update(ctx, {{.LowerName}}); err != nil {
		h.logger.Error("Failed to update {{.TableName}}", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update {{.TableName}}")
		return
	}

	response := {{.StructName}}Response{
{{range .Columns}}
		{{.FieldName}}: {{$.LowerName}}.{{.FieldName}},
{{end}}
	}

	utils.SuccessResponse(c, http.StatusOK, "{{.TableName}} updated successfully", response)
}

// Delete{{.StructName}} deletes a {{.TableName}}
// @Summary Delete {{.TableName}}
// @Description Delete a {{.TableName}} by ID
// @Tags {{.TableName}}
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "{{.TableName}} ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /{{.TableName}}/{id} [delete]
func (h *Handler) Delete{{.StructName}}(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	ctx := context.Background()
	if err := h.{{.LowerName}}Repo.Delete(ctx, uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "{{.TableName}} not found")
			return
		}
		h.logger.Error("Failed to delete {{.TableName}}", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete {{.TableName}}")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "{{.TableName}} deleted successfully", nil)
}
`

	t, err := template.New("handler").Parse(tmpl)
	if err != nil {
		return err
	}

	file, err := os.Create(handlerFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return t.Execute(file, map[string]interface{}{
		"PackageName":   tableName,
		"StructName":    fg.toPascalCase(tableName),
		"LowerName":     fg.toCamelCase(tableName),
		"TableName":     tableName,
		"Columns":       fg.generateColumnInfo(tableInfo.Columns),
		"CreateColumns": fg.generateCreateColumns(tableInfo.Columns),
		"UpdateColumns": fg.generateUpdateColumns(tableInfo.Columns),
	})
}

// generateRoutes generates the routes file
func (fg *FileGenerator) generateRoutes(tableName string, tableInfo *TableInfo, apiDir string) error {
	routesFile := filepath.Join(apiDir, "routes.go")

	tmpl := `package {{.PackageName}}

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go-mobile-backend-template/internal/middleware"
	authService "go-mobile-backend-template/internal/services/auth"
	"go-mobile-backend-template/pkg/config"
)

// RegisterRoutes registers {{.TableName}} routes
func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, logger *zap.Logger, cfg *config.Config) {
	handler := NewHandler(db, logger)
	
	// Initialize JWT service for auth middleware
	jwtService := authService.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenExpireInt,
		cfg.JWT.RefreshTokenExpireInt,
	)

	// {{.TableName}} routes (all protected)
	{{.LowerName}}Routes := router.Group("/{{.TableName}}")
	{{.LowerName}}Routes.Use(middleware.AuthMiddleware(jwtService))
	{
		{{.LowerName}}Routes.POST("", handler.Create{{.StructName}})
		{{.LowerName}}Routes.GET("", handler.GetAll{{.StructName}}s)
		{{.LowerName}}Routes.GET("/:id", handler.Get{{.StructName}})
		{{.LowerName}}Routes.PUT("/:id", handler.Update{{.StructName}})
		{{.LowerName}}Routes.DELETE("/:id", handler.Delete{{.StructName}})
	}
}
`

	t, err := template.New("routes").Parse(tmpl)
	if err != nil {
		return err
	}

	file, err := os.Create(routesFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return t.Execute(file, map[string]interface{}{
		"PackageName": tableName,
		"StructName":  fg.toPascalCase(tableName),
		"LowerName":   fg.toCamelCase(tableName),
		"TableName":   tableName,
	})
}

// Helper functions
func (fg *FileGenerator) toPascalCase(s string) string {
	return strings.Title(strings.ReplaceAll(s, "_", ""))
}

func (fg *FileGenerator) toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	if len(parts) == 0 {
		return ""
	}

	result := strings.ToLower(parts[0])
	for _, part := range parts[1:] {
		result += strings.Title(part)
	}
	return result
}

func (fg *FileGenerator) generateColumnInfo(columns []ColumnInfo) []map[string]interface{} {
	var result []map[string]interface{}

	for _, col := range columns {
		result = append(result, map[string]interface{}{
			"FieldName": fg.toPascalCase(col.Name),
			"GoType":    fg.getGoType(col.Type),
			"GormTag":   fg.getGormTag(col),
			"JSONTag":   fg.getJSONTag(col.Name),
			"Comment":   col.Comment,
		})
	}

	return result
}

func (fg *FileGenerator) generateCreateColumns(columns []ColumnInfo) []map[string]interface{} {
	var result []map[string]interface{}

	for _, col := range columns {
		// Skip auto-generated fields for create
		if col.Name == "id" || col.Name == "created_at" || col.Name == "updated_at" {
			continue
		}

		result = append(result, map[string]interface{}{
			"FieldName":  fg.toPascalCase(col.Name),
			"GoType":     fg.getGoType(col.Type),
			"JSONTag":    fg.getJSONTag(col.Name),
			"BindingTag": fg.getBindingTag(col),
			"Comment":    col.Comment,
		})
	}

	return result
}

func (fg *FileGenerator) generateUpdateColumns(columns []ColumnInfo) []map[string]interface{} {
	var result []map[string]interface{}

	for _, col := range columns {
		// Skip auto-generated fields for update
		if col.Name == "id" || col.Name == "created_at" || col.Name == "updated_at" {
			continue
		}

		result = append(result, map[string]interface{}{
			"FieldName":  fg.toPascalCase(col.Name),
			"GoType":     fg.getGoType(col.Type),
			"JSONTag":    fg.getJSONTag(col.Name),
			"BindingTag": fg.getUpdateBindingTag(col),
			"Comment":    col.Comment,
		})
	}

	return result
}

func (fg *FileGenerator) getGoType(dbType string) string {
	switch {
	case strings.Contains(dbType, "varchar"), strings.Contains(dbType, "text"), strings.Contains(dbType, "char"):
		return "string"
	case strings.Contains(dbType, "int"):
		return "uint"
	case strings.Contains(dbType, "bigint"):
		return "uint64"
	case strings.Contains(dbType, "boolean"), strings.Contains(dbType, "bool"):
		return "bool"
	case strings.Contains(dbType, "timestamp"), strings.Contains(dbType, "datetime"):
		return "time.Time"
	case strings.Contains(dbType, "decimal"), strings.Contains(dbType, "numeric"), strings.Contains(dbType, "float"):
		return "float64"
	default:
		return "string"
	}
}

func (fg *FileGenerator) getGormTag(col ColumnInfo) string {
	tags := []string{}

	if col.IsPrimaryKey {
		tags = append(tags, "primaryKey")
	}

	// Check if it's an auto-increment field (usually ID fields)
	if col.IsPrimaryKey && (col.Type == "serial" || col.Type == "bigserial" || strings.Contains(col.Type, "auto_increment")) {
		tags = append(tags, "autoIncrement")
	}

	if !col.IsNullable {
		tags = append(tags, "not null")
	}

	if col.DefaultValue != nil && *col.DefaultValue != "" {
		tags = append(tags, fmt.Sprintf("default:%s", *col.DefaultValue))
	}

	return strings.Join(tags, ";")
}

func (fg *FileGenerator) getJSONTag(fieldName string) string {
	return strings.ToLower(fieldName)
}

func (fg *FileGenerator) getBindingTag(col ColumnInfo) string {
	tags := []string{}

	if !col.IsNullable {
		tags = append(tags, "required")
	}

	if col.MaxLength != nil && *col.MaxLength > 0 {
		tags = append(tags, fmt.Sprintf("max=%d", *col.MaxLength))
	}

	return strings.Join(tags, ",")
}

func (fg *FileGenerator) getUpdateBindingTag(col ColumnInfo) string {
	// For updates, make fields optional
	tags := []string{"omitempty"}

	if col.MaxLength != nil && *col.MaxLength > 0 {
		tags = append(tags, fmt.Sprintf("max=%d", *col.MaxLength))
	}

	return strings.Join(tags, ",")
}
