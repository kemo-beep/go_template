package generator

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-mobile-backend-template/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CRUDHandlerGenerator generates CRUD handlers for tables
type CRUDHandlerGenerator struct {
	db     *gorm.DB
	logger *zap.Logger
	config *GeneratorConfig
}

// NewCRUDHandlerGenerator creates a new CRUD handler generator
func NewCRUDHandlerGenerator(db *gorm.DB, logger *zap.Logger, config *GeneratorConfig) *CRUDHandlerGenerator {
	return &CRUDHandlerGenerator{
		db:     db,
		logger: logger,
		config: config,
	}
}

// GenerateHandlers generates all CRUD handlers for a table
func (g *CRUDHandlerGenerator) GenerateHandlers(table *TableInfo) (map[string]gin.HandlerFunc, error) {
	handlers := make(map[string]gin.HandlerFunc)
	tableConfig := g.config.GetTableConfig(table.Name)

	// Generate handlers for each endpoint type
	for _, endpointType := range tableConfig.Endpoints {
		handler, err := g.generateHandler(table, endpointType, tableConfig)
		if err != nil {
			g.logger.Error("Failed to generate handler",
				zap.String("table", table.Name),
				zap.String("endpoint", endpointType),
				zap.Error(err))
			continue
		}

		handlerName := fmt.Sprintf("%s%s", g.toCamelCase(endpointType), g.toCamelCase(table.Name))
		handlers[handlerName] = handler
	}

	return handlers, nil
}

// generateHandler generates a specific handler for a table
func (g *CRUDHandlerGenerator) generateHandler(table *TableInfo, endpointType string, config *TableConfig) (gin.HandlerFunc, error) {
	switch endpointType {
	case "list":
		return g.generateListHandler(table, config), nil
	case "create":
		return g.generateCreateHandler(table, config), nil
	case "get":
		return g.generateGetHandler(table, config), nil
	case "update":
		return g.generateUpdateHandler(table, config), nil
	case "delete":
		return g.generateDeleteHandler(table, config), nil
	case "bulk":
		return g.generateBulkHandler(table, config), nil
	case "search":
		return g.generateSearchHandler(table, config), nil
	case "stats":
		return g.generateStatsHandler(table, config), nil
	case "export":
		return g.generateExportHandler(table, config), nil
	default:
		return nil, fmt.Errorf("unknown endpoint type: %s", endpointType)
	}
}

// generateListHandler generates a list handler
func (g *CRUDHandlerGenerator) generateListHandler(table *TableInfo, config *TableConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse pagination parameters
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", fmt.Sprintf("%d", config.Pagination.DefaultLimit)))

		if limit > config.Pagination.MaxLimit {
			limit = config.Pagination.MaxLimit
		}

		offset := (page - 1) * limit

		// Parse sorting parameters
		sort := c.Query("sort")
		order := c.Query("order")
		if order == "" {
			order = "asc"
		}

		// Build query
		query := g.db.Table(table.Name)

		// Apply filters
		if err := g.applyFilters(query, c, config); err != nil {
			g.logger.Error("Failed to apply filters", zap.Error(err))
			c.JSON(http.StatusBadRequest, utils.ErrorResponseData("Invalid filter parameters"))
			return
		}

		// Apply sorting
		if sort != "" {
			if config.Sorting != nil && g.isAllowedSortField(sort, config.Sorting.AllowedFields) {
				query = query.Order(fmt.Sprintf("%s %s", sort, order))
			}
		} else if config.Sorting != nil && config.Sorting.DefaultSort != "" {
			query = query.Order(config.Sorting.DefaultSort)
		}

		// Get total count
		var total int64
		if err := query.Count(&total).Error; err != nil {
			g.logger.Error("Failed to count records", zap.Error(err))
			c.JSON(http.StatusInternalServerError, utils.ErrorResponseData("Failed to count records"))
			return
		}

		// Apply pagination
		query = query.Offset(offset).Limit(limit)

		// Execute query
		var results []map[string]interface{}
		if err := query.Find(&results).Error; err != nil {
			g.logger.Error("Failed to fetch records", zap.Error(err))
			c.JSON(http.StatusInternalServerError, utils.ErrorResponseData("Failed to fetch records"))
			return
		}

		// Calculate pagination info
		totalPages := int((total + int64(limit) - 1) / int64(limit))
		hasNext := page < totalPages
		hasPrev := page > 1

		// Return response
		c.JSON(http.StatusOK, utils.SuccessResponseData("Records retrieved successfully", gin.H{
			"data": results,
			"pagination": gin.H{
				"page":        page,
				"limit":       limit,
				"total":       total,
				"total_pages": totalPages,
				"has_next":    hasNext,
				"has_prev":    hasPrev,
			},
		}))
	}
}

// generateCreateHandler generates a create handler
func (g *CRUDHandlerGenerator) generateCreateHandler(table *TableInfo, config *TableConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse request body
		var data map[string]interface{}
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponseData("Invalid request body"))
			return
		}

		// Validate data
		if err := g.validateData(data, table, config, "create"); err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponseData(err.Error()))
			return
		}

		// Add timestamps if enabled
		if config.Security != nil && config.Security.Timestamps {
			now := time.Now()
			data["created_at"] = now
			data["updated_at"] = now
		}

		// Insert record
		if err := g.db.Table(table.Name).Create(&data).Error; err != nil {
			g.logger.Error("Failed to create record", zap.Error(err))
			c.JSON(http.StatusInternalServerError, utils.ErrorResponseData("Failed to create record"))
			return
		}

		c.JSON(http.StatusCreated, utils.SuccessResponseData("Record created successfully", gin.H{
			"data": data,
		}))
	}
}

// generateGetHandler generates a get handler
func (g *CRUDHandlerGenerator) generateGetHandler(table *TableInfo, config *TableConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, utils.ErrorResponseData("ID parameter is required"))
			return
		}

		// Build query
		query := g.db.Table(table.Name)

		// Apply joins if configured
		if err := g.applyJoins(query, table, config); err != nil {
			g.logger.Error("Failed to apply joins", zap.Error(err))
		}

		// Find record by ID
		var result map[string]interface{}
		if err := query.Where("id = ?", id).First(&result).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, utils.ErrorResponseData("Record not found"))
				return
			}
			g.logger.Error("Failed to fetch record", zap.Error(err))
			c.JSON(http.StatusInternalServerError, utils.ErrorResponseData("Failed to fetch record"))
			return
		}

		c.JSON(http.StatusOK, utils.SuccessResponseData("Record retrieved successfully", gin.H{
			"data": result,
		}))
	}
}

// generateUpdateHandler generates an update handler
func (g *CRUDHandlerGenerator) generateUpdateHandler(table *TableInfo, config *TableConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, utils.ErrorResponseData("ID parameter is required"))
			return
		}

		// Parse request body
		var data map[string]interface{}
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponseData("Invalid request body"))
			return
		}

		// Validate data
		if err := g.validateData(data, table, config, "update"); err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponseData(err.Error()))
			return
		}

		// Add updated timestamp if enabled
		if config.Security != nil && config.Security.Timestamps {
			data["updated_at"] = time.Now()
		}

		// Update record
		result := g.db.Table(table.Name).Where("id = ?", id).Updates(data)
		if result.Error != nil {
			g.logger.Error("Failed to update record", zap.Error(result.Error))
			c.JSON(http.StatusInternalServerError, utils.ErrorResponseData("Failed to update record"))
			return
		}

		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, utils.ErrorResponseData("Record not found"))
			return
		}

		// Fetch updated record
		var updatedRecord map[string]interface{}
		if err := g.db.Table(table.Name).Where("id = ?", id).First(&updatedRecord).Error; err != nil {
			g.logger.Error("Failed to fetch updated record", zap.Error(err))
		}

		c.JSON(http.StatusOK, utils.SuccessResponseData("Record updated successfully", gin.H{
			"data": updatedRecord,
		}))
	}
}

// generateDeleteHandler generates a delete handler
func (g *CRUDHandlerGenerator) generateDeleteHandler(table *TableInfo, config *TableConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, utils.ErrorResponseData("ID parameter is required"))
			return
		}

		// Check if soft delete is enabled
		if config.Security != nil && config.Security.SoftDelete {
			// Soft delete - update deleted_at timestamp
			result := g.db.Table(table.Name).Where("id = ?", id).Update("deleted_at", time.Now())
			if result.Error != nil {
				g.logger.Error("Failed to soft delete record", zap.Error(result.Error))
				c.JSON(http.StatusInternalServerError, utils.ErrorResponseData("Failed to delete record"))
				return
			}

			if result.RowsAffected == 0 {
				c.JSON(http.StatusNotFound, utils.ErrorResponseData("Record not found"))
				return
			}
		} else {
			// Hard delete
			result := g.db.Table(table.Name).Where("id = ?", id).Delete(nil)
			if result.Error != nil {
				g.logger.Error("Failed to delete record", zap.Error(result.Error))
				c.JSON(http.StatusInternalServerError, utils.ErrorResponseData("Failed to delete record"))
				return
			}

			if result.RowsAffected == 0 {
				c.JSON(http.StatusNotFound, utils.ErrorResponseData("Record not found"))
				return
			}
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

// generateBulkHandler generates a bulk operations handler
func (g *CRUDHandlerGenerator) generateBulkHandler(table *TableInfo, config *TableConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			Operation string                   `json:"operation" binding:"required"`
			Data      []map[string]interface{} `json:"data"`
			Where     map[string]interface{}   `json:"where"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponseData("Invalid request body"))
			return
		}

		// Start transaction
		tx := g.db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		var created, updated, deleted int
		var errors []string

		switch request.Operation {
		case "create":
			for _, record := range request.Data {
				if err := g.validateData(record, table, config, "create"); err != nil {
					errors = append(errors, fmt.Sprintf("Validation error: %s", err.Error()))
					continue
				}

				if config.Security != nil && config.Security.Timestamps {
					now := time.Now()
					record["created_at"] = now
					record["updated_at"] = now
				}

				if err := tx.Table(table.Name).Create(&record).Error; err != nil {
					errors = append(errors, fmt.Sprintf("Failed to create record: %s", err.Error()))
					continue
				}
				created++
			}

		case "update":
			for _, record := range request.Data {
				if err := g.validateData(record, table, config, "update"); err != nil {
					errors = append(errors, fmt.Sprintf("Validation error: %s", err.Error()))
					continue
				}

				if config.Security != nil && config.Security.Timestamps {
					record["updated_at"] = time.Now()
				}

				if err := tx.Table(table.Name).Where("id = ?", record["id"]).Updates(record).Error; err != nil {
					errors = append(errors, fmt.Sprintf("Failed to update record: %s", err.Error()))
					continue
				}
				updated++
			}

		case "delete":
			if request.Where != nil {
				// Delete by conditions
				result := tx.Table(table.Name).Where(request.Where).Delete(nil)
				if result.Error != nil {
					errors = append(errors, fmt.Sprintf("Failed to delete records: %s", result.Error.Error()))
				} else {
					deleted = int(result.RowsAffected)
				}
			} else {
				// Delete by IDs
				for _, record := range request.Data {
					if id, exists := record["id"]; exists {
						result := tx.Table(table.Name).Where("id = ?", id).Delete(nil)
						if result.Error != nil {
							errors = append(errors, fmt.Sprintf("Failed to delete record %v: %s", id, result.Error.Error()))
							continue
						}
						deleted++
					}
				}
			}

		default:
			c.JSON(http.StatusBadRequest, utils.ErrorResponseData("Invalid operation"))
			return
		}

		// Commit transaction
		if err := tx.Commit().Error; err != nil {
			c.JSON(http.StatusInternalServerError, utils.ErrorResponseData("Failed to commit transaction"))
			return
		}

		c.JSON(http.StatusOK, utils.SuccessResponseData("Bulk operation completed", gin.H{
			"created": created,
			"updated": updated,
			"deleted": deleted,
			"errors":  errors,
		}))
	}
}

// generateSearchHandler generates a search handler
func (g *CRUDHandlerGenerator) generateSearchHandler(table *TableInfo, config *TableConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			c.JSON(http.StatusBadRequest, utils.ErrorResponseData("Search query is required"))
			return
		}

		// Parse pagination parameters
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", fmt.Sprintf("%d", config.Pagination.DefaultLimit)))

		if limit > config.Pagination.MaxLimit {
			limit = config.Pagination.MaxLimit
		}

		offset := (page - 1) * limit

		// Build search query
		dbQuery := g.db.Table(table.Name)

		// Add search conditions for text fields
		if config.Filtering != nil && len(config.Filtering.TextSearch) > 0 {
			var conditions []string
			var args []interface{}

			for _, field := range config.Filtering.TextSearch {
				conditions = append(conditions, fmt.Sprintf("%s ILIKE ?", field))
				args = append(args, "%"+query+"%")
			}

			if len(conditions) > 0 {
				dbQuery = dbQuery.Where(strings.Join(conditions, " OR "), args...)
			}
		}

		// Get total count
		var total int64
		if err := dbQuery.Count(&total).Error; err != nil {
			g.logger.Error("Failed to count search results", zap.Error(err))
			c.JSON(http.StatusInternalServerError, utils.ErrorResponseData("Failed to search records"))
			return
		}

		// Apply pagination
		dbQuery = dbQuery.Offset(offset).Limit(limit)

		// Execute search
		var results []map[string]interface{}
		if err := dbQuery.Find(&results).Error; err != nil {
			g.logger.Error("Failed to execute search", zap.Error(err))
			c.JSON(http.StatusInternalServerError, utils.ErrorResponseData("Failed to search records"))
			return
		}

		// Calculate pagination info
		totalPages := int((total + int64(limit) - 1) / int64(limit))
		hasNext := page < totalPages
		hasPrev := page > 1

		c.JSON(http.StatusOK, utils.SuccessResponseData("Search completed", gin.H{
			"data": results,
			"pagination": gin.H{
				"page":        page,
				"limit":       limit,
				"total":       total,
				"total_pages": totalPages,
				"has_next":    hasNext,
				"has_prev":    hasPrev,
			},
		}))
	}
}

// generateStatsHandler generates a stats handler
func (g *CRUDHandlerGenerator) generateStatsHandler(table *TableInfo, config *TableConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get total count
		var total int64
		if err := g.db.Table(table.Name).Count(&total).Error; err != nil {
			g.logger.Error("Failed to get total count", zap.Error(err))
			c.JSON(http.StatusInternalServerError, utils.ErrorResponseData("Failed to get statistics"))
			return
		}

		// Get active count (if status field exists)
		var active int64
		if g.hasColumn(table, "status") {
			if err := g.db.Table(table.Name).Where("status = ?", "active").Count(&active).Error; err != nil {
				g.logger.Warn("Failed to get active count", zap.Error(err))
			}
		}

		// Get inactive count
		inactive := total - active

		c.JSON(http.StatusOK, utils.SuccessResponseData("Statistics retrieved", gin.H{
			"total":    total,
			"active":   active,
			"inactive": inactive,
		}))
	}
}

// generateExportHandler generates an export handler
func (g *CRUDHandlerGenerator) generateExportHandler(table *TableInfo, config *TableConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		format := c.DefaultQuery("format", "csv")
		fields := c.Query("fields")

		// Build query
		query := g.db.Table(table.Name)

		// Apply filters
		if err := g.applyFilters(query, c, config); err != nil {
			g.logger.Error("Failed to apply filters", zap.Error(err))
			c.JSON(http.StatusBadRequest, utils.ErrorResponseData("Invalid filter parameters"))
			return
		}

		// Select specific fields if requested
		if fields != "" {
			fieldList := strings.Split(fields, ",")
			query = query.Select(fieldList)
		}

		// Execute query
		var results []map[string]interface{}
		if err := query.Find(&results).Error; err != nil {
			g.logger.Error("Failed to fetch records for export", zap.Error(err))
			c.JSON(http.StatusInternalServerError, utils.ErrorResponseData("Failed to export records"))
			return
		}

		// Set appropriate headers
		switch format {
		case "csv":
			c.Header("Content-Type", "text/csv")
			c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", table.Name))
		case "json":
			c.Header("Content-Type", "application/json")
			c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.json", table.Name))
		case "xlsx":
			c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
			c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.xlsx", table.Name))
		default:
			c.JSON(http.StatusBadRequest, utils.ErrorResponseData("Unsupported export format"))
			return
		}

		// For now, return JSON. In production, you'd implement proper CSV/Excel export
		c.JSON(http.StatusOK, gin.H{
			"data":   results,
			"format": format,
			"count":  len(results),
		})
	}
}

// Helper methods

func (g *CRUDHandlerGenerator) applyFilters(query *gorm.DB, c *gin.Context, config *TableConfig) error {
	if config.Filtering == nil {
		return nil
	}

	for _, field := range config.Filtering.AllowedFields {
		value := c.Query(field)
		if value == "" {
			continue
		}

		// Apply filter based on field type
		// For now, use ILIKE for all fields - in production you'd check field types
		query = query.Where(fmt.Sprintf("%s ILIKE ?", field), "%"+value+"%")
	}

	return nil
}

func (g *CRUDHandlerGenerator) applyJoins(query *gorm.DB, table *TableInfo, config *TableConfig) error {
	// This would apply joins based on foreign keys
	// For now, return nil
	return nil
}

func (g *CRUDHandlerGenerator) validateData(data map[string]interface{}, table *TableInfo, config *TableConfig, operation string) error {
	if config.Validation == nil {
		return nil
	}

	// Check required fields
	for _, field := range config.Validation.Required {
		if _, exists := data[field]; !exists {
			return fmt.Errorf("field %s is required", field)
		}
	}

	// Validate field lengths
	for field, minLength := range config.Validation.MinLength {
		if value, exists := data[field]; exists {
			if str, ok := value.(string); ok && len(str) < minLength {
				return fmt.Errorf("field %s must be at least %d characters long", field, minLength)
			}
		}
	}

	for field, maxLength := range config.Validation.MaxLength {
		if value, exists := data[field]; exists {
			if str, ok := value.(string); ok && len(str) > maxLength {
				return fmt.Errorf("field %s must be at most %d characters long", field, maxLength)
			}
		}
	}

	return nil
}

func (g *CRUDHandlerGenerator) isAllowedSortField(field string, allowedFields []string) bool {
	return g.contains(allowedFields, field)
}

func (g *CRUDHandlerGenerator) hasColumn(table *TableInfo, columnName string) bool {
	for _, column := range table.Columns {
		if column.Name == columnName {
			return true
		}
	}
	return false
}

func (g *CRUDHandlerGenerator) isNumericField(table *TableInfo, fieldName string) bool {
	for _, column := range table.Columns {
		if column.Name == fieldName {
			return strings.Contains(column.Type, "int") || strings.Contains(column.Type, "numeric") || strings.Contains(column.Type, "decimal")
		}
	}
	return false
}

func (g *CRUDHandlerGenerator) contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (g *CRUDHandlerGenerator) toCamelCase(str string) string {
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
