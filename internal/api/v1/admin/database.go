package admin

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-mobile-backend-template/internal/middleware"
	"go-mobile-backend-template/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DatabaseHandler struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewDatabaseHandler(db *gorm.DB, logger *zap.Logger) *DatabaseHandler {
	return &DatabaseHandler{
		db:     db,
		logger: logger,
	}
}

// TableInfo represents information about a database table
type TableInfo struct {
	Name     string `json:"name"`
	RowCount int64  `json:"row_count"`
	Size     string `json:"size,omitempty"`
}

// ColumnInfo represents information about a table column
type ColumnInfo struct {
	Name         string  `json:"name"`
	Type         string  `json:"type"`
	Nullable     bool    `json:"nullable"`
	DefaultValue *string `json:"default_value"`
	IsPrimaryKey bool    `json:"is_primary_key"`
}

// ListTables godoc
// @Summary List all database tables (Admin)
// @Description Get list of all tables in the database
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Router /admin/database/tables [get]
func (h *DatabaseHandler) ListTables(c *gin.Context) {
	var tables []TableInfo

	// Get all table names
	rows, err := h.db.Raw(`
		SELECT 
			relname as name,
			n_live_tup as row_count
		FROM pg_stat_user_tables
		ORDER BY relname
	`).Rows()
	if err != nil {
		h.logger.Error("Failed to list tables", zap.Error(err))
		c.JSON(http.StatusInternalServerError, utils.ErrorResponseData("Failed to list tables"))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var table TableInfo
		var rowCount sql.NullInt64

		if err := rows.Scan(&table.Name, &rowCount); err != nil {
			h.logger.Error("Failed to scan table", zap.Error(err))
			continue
		}

		if rowCount.Valid {
			table.RowCount = rowCount.Int64
		}

		tables = append(tables, table)
	}

	c.JSON(http.StatusOK, utils.SuccessResponseData("Tables fetched successfully", tables))
}

// GetTableSchema godoc
// @Summary Get table schema (Admin)
// @Description Get column information for a specific table
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param tableName path string true "Table name"
// @Success 200 {object} utils.Response
// @Router /admin/database/tables/{tableName}/schema [get]
func (h *DatabaseHandler) GetTableSchema(c *gin.Context) {
	tableName := c.Param("tableName")
	h.logger.Info("Getting table schema", zap.String("table", tableName))

	var columns []ColumnInfo

	rows, err := h.db.Raw(`
		SELECT 
			column_name,
			data_type,
			is_nullable = 'YES' as nullable,
			column_default,
			false as is_primary_key
		FROM information_schema.columns
		WHERE table_name = ? AND table_schema = 'public'
		ORDER BY ordinal_position
	`, tableName).Rows()

	if err != nil {
		h.logger.Error("Failed to get table schema", zap.Error(err), zap.String("table", tableName))
		c.JSON(http.StatusInternalServerError, utils.ErrorResponseData("Failed to get table schema"))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var col ColumnInfo
		var defaultVal sql.NullString

		if err := rows.Scan(&col.Name, &col.Type, &col.Nullable, &defaultVal, &col.IsPrimaryKey); err != nil {
			h.logger.Error("Failed to scan column", zap.Error(err))
			continue
		}

		if defaultVal.Valid {
			col.DefaultValue = &defaultVal.String
		}

		columns = append(columns, col)
	}

	c.JSON(http.StatusOK, utils.SuccessResponseData("Schema fetched successfully", gin.H{
		"table":   tableName,
		"columns": columns,
	}))
}

// GetTableData godoc
// @Summary Get table data (Admin)
// @Description Get paginated data from a specific table
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param tableName path string true "Table name"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} utils.Response
// @Router /admin/database/tables/{tableName}/data [get]
func (h *DatabaseHandler) GetTableData(c *gin.Context) {
	tableName := c.Param("tableName")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	// Validate table name to prevent SQL injection
	var exists bool
	err := h.db.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = ?)", tableName).Scan(&exists).Error
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, utils.ErrorResponseData("Table not found"))
		return
	}

	// Get total count using parameterized query
	var total int64
	err = h.db.Raw("SELECT COUNT(*) FROM ?", tableName).Scan(&total).Error
	if err != nil {
		h.logger.Error("Failed to get table count", zap.Error(err), zap.String("table", tableName))
		c.JSON(http.StatusInternalServerError, utils.ErrorResponseData("Failed to get table count"))
		return
	}

	// Get data using parameterized query
	rows, err := h.db.Raw("SELECT * FROM ? LIMIT ? OFFSET ?", tableName, limit, offset).Rows()
	if err != nil {
		h.logger.Error("Failed to get table data", zap.Error(err), zap.String("table", tableName))
		c.JSON(http.StatusInternalServerError, utils.ErrorResponseData("Failed to get table data"))
		return
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		h.logger.Error("Failed to get columns", zap.Error(err))
		c.JSON(http.StatusInternalServerError, utils.ErrorResponseData("Failed to get columns"))
		return
	}

	// Prepare to scan rows
	var data []map[string]interface{}

	for rows.Next() {
		// Create a slice of interface{}'s to represent each column
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			h.logger.Error("Failed to scan row", zap.Error(err))
			continue
		}

		// Create map for this row
		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]

			// Convert []byte to string for display
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}

		data = append(data, row)
	}

	c.JSON(http.StatusOK, utils.SuccessResponseData("Data fetched successfully", gin.H{
		"table":   tableName,
		"columns": columns,
		"data":    data,
		"total":   total,
		"page":    page,
		"limit":   limit,
	}))
}

// ExecuteQuery godoc
// @Summary Execute SQL query (Admin)
// @Description Execute a read-only SQL query
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ExecuteQueryRequest true "Query request"
// @Success 200 {object} utils.Response
// @Router /admin/database/query [post]
func (h *DatabaseHandler) ExecuteQuery(c *gin.Context) {
	var req ExecuteQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponseData("Invalid request"))
		return
	}

	// Enhanced SQL query validation
	if valid, reason := h.validateSQLQuery(req.Query); !valid {
		h.logger.Warn("Invalid SQL query attempted",
			zap.String("query", req.Query),
			zap.String("reason", reason),
			zap.String("ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, utils.ErrorResponseData("Invalid query: "+reason))
		return
	}

	// Execute query with timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	rows, err := h.db.WithContext(ctx).Raw(req.Query).Rows()
	if err != nil {
		h.logger.Error("Failed to execute query", zap.Error(err), zap.String("query", req.Query))
		c.JSON(http.StatusInternalServerError, utils.ErrorResponseData("Failed to execute query: "+err.Error()))
		return
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		h.logger.Error("Failed to get columns", zap.Error(err))
		c.JSON(http.StatusInternalServerError, utils.ErrorResponseData("Failed to get columns"))
		return
	}

	// Scan results
	var results []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			h.logger.Error("Failed to scan row", zap.Error(err))
			continue
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}

		results = append(results, row)
	}

	c.JSON(http.StatusOK, utils.SuccessResponseData("Query executed successfully", gin.H{
		"columns": columns,
		"data":    results,
		"count":   len(results),
	}))
}

type ExecuteQueryRequest struct {
	Query string `json:"query" binding:"required"`
}

// validateSQLQuery validates SQL queries for safety
func (h *DatabaseHandler) validateSQLQuery(query string) (bool, string) {
	query = strings.TrimSpace(query)

	// Only allow SELECT queries for read operations
	upperQuery := strings.ToUpper(query)
	if !strings.HasPrefix(upperQuery, "SELECT") {
		return false, "Only SELECT queries are allowed"
	}

	// Check for dangerous keywords
	dangerousKeywords := []string{
		"DROP", "DELETE", "INSERT", "UPDATE", "CREATE", "ALTER",
		"EXEC", "EXECUTE", "SP_", "XP_", "OPENROWSET", "OPENDATASOURCE",
		"BULK", "BULKINSERT", "BACKUP", "RESTORE", "SHUTDOWN",
		"RECONFIGURE", "DBCC", "KILL", "DENY", "REVOKE",
	}

	for _, keyword := range dangerousKeywords {
		if strings.Contains(upperQuery, keyword) {
			return false, "Dangerous SQL keyword detected: " + keyword
		}
	}

	// Check for SQL injection patterns
	sqlSecurity := middleware.NewSQLSecurity(h.logger)
	if valid, reason := sqlSecurity.ValidateSQLInput(query); !valid {
		return false, reason
	}

	return true, ""
}

// GetDatabaseStats godoc
// @Summary Get database statistics (Admin)
// @Description Get overall database statistics
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.Response
// @Router /admin/database/stats [get]
func (h *DatabaseHandler) GetDatabaseStats(c *gin.Context) {
	var stats struct {
		DatabaseSize string `json:"database_size"`
		TableCount   int64  `json:"table_count"`
		TotalRows    int64  `json:"total_rows"`
	}

	// Get database size
	h.db.Raw("SELECT pg_size_pretty(pg_database_size(current_database()))").Scan(&stats.DatabaseSize)

	// Get table count
	h.db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'").Scan(&stats.TableCount)

	// Get total rows across all tables
	h.db.Raw("SELECT SUM(n_live_tup) FROM pg_stat_user_tables").Scan(&stats.TotalRows)

	c.JSON(http.StatusOK, utils.SuccessResponseData("Stats fetched successfully", stats))
}
