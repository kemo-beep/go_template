package admin

import (
	"fmt"
	"net/http"
	"strings"

	"go-mobile-backend-template/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TableManagerHandler struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewTableManagerHandler(db *gorm.DB, logger *zap.Logger) *TableManagerHandler {
	return &TableManagerHandler{
		db:     db,
		logger: logger,
	}
}

// CreateTableRequest represents a request to create a new table
type CreateTableRequest struct {
	TableName string          `json:"table_name" binding:"required"`
	Columns   []ColumnRequest `json:"columns" binding:"required,min=1"`
}

// ColumnRequest represents a column definition
type ColumnRequest struct {
	Name         string  `json:"name" binding:"required"`
	Type         string  `json:"type" binding:"required"` // varchar, integer, boolean, timestamp, text, jsonb, etc.
	Length       *int    `json:"length,omitempty"`        // For varchar(length)
	NotNull      bool    `json:"not_null"`
	PrimaryKey   bool    `json:"primary_key"`
	Unique       bool    `json:"unique"`
	DefaultValue *string `json:"default_value,omitempty"`
	References   *string `json:"references,omitempty"` // For foreign keys: "table_name(column)"
}

// CreateTable godoc
// @Summary Create a new database table (Admin)
// @Description Create a new table with specified columns
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateTableRequest true "Table creation request"
// @Success 201 {object} map[string]interface{}
// @Router /admin/database/tables [post]
func (h *TableManagerHandler) CreateTable(c *gin.Context) {
	var req CreateTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponseData("Invalid request"))
		return
	}

	// Validate table name (prevent SQL injection)
	if !isValidIdentifier(req.TableName) {
		c.JSON(http.StatusBadRequest, utils.ErrorResponseData("Invalid table name"))
		return
	}

	// Build CREATE TABLE SQL
	var sqlBuilder strings.Builder
	sqlBuilder.WriteString(fmt.Sprintf("CREATE TABLE %s (\n", req.TableName))

	columnDefs := make([]string, 0, len(req.Columns))
	var primaryKeys []string

	for i, col := range req.Columns {
		if !isValidIdentifier(col.Name) {
			c.JSON(http.StatusBadRequest, utils.ErrorResponseData(fmt.Sprintf("Invalid column name: %s", col.Name)))
			return
		}

		var colDef strings.Builder
		colDef.WriteString(fmt.Sprintf("  %s", col.Name))

		// Add data type
		dataType := strings.ToUpper(col.Type)
		if col.Length != nil && (dataType == "VARCHAR" || dataType == "CHAR") {
			colDef.WriteString(fmt.Sprintf(" %s(%d)", dataType, *col.Length))
		} else {
			colDef.WriteString(fmt.Sprintf(" %s", dataType))
		}

		// Add constraints
		if col.NotNull {
			colDef.WriteString(" NOT NULL")
		}

		if col.Unique {
			colDef.WriteString(" UNIQUE")
		}

		if col.DefaultValue != nil {
			colDef.WriteString(fmt.Sprintf(" DEFAULT %s", *col.DefaultValue))
		}

		if col.PrimaryKey {
			primaryKeys = append(primaryKeys, col.Name)
		}

		columnDefs = append(columnDefs, colDef.String())

		// Add foreign key constraint if specified
		if col.References != nil {
			parts := strings.Split(*col.References, "(")
			if len(parts) == 2 {
				refTable := parts[0]
				refColumn := strings.TrimSuffix(parts[1], ")")
				fkConstraint := fmt.Sprintf("  CONSTRAINT fk_%s_%s FOREIGN KEY (%s) REFERENCES %s(%s)",
					req.TableName, col.Name, col.Name, refTable, refColumn)
				columnDefs = append(columnDefs, fkConstraint)
			}
		}

		if i < len(req.Columns)-1 || len(primaryKeys) > 0 {
			columnDefs[len(columnDefs)-1] += ","
		}
	}

	// Add primary key constraint
	if len(primaryKeys) > 0 {
		pkConstraint := fmt.Sprintf("  PRIMARY KEY (%s)", strings.Join(primaryKeys, ", "))
		columnDefs = append(columnDefs, pkConstraint)
	}

	sqlBuilder.WriteString(strings.Join(columnDefs, "\n"))
	sqlBuilder.WriteString("\n);")

	sql := sqlBuilder.String()

	// Execute SQL
	if err := h.db.Exec(sql).Error; err != nil {
		h.logger.Error("Failed to create table", zap.Error(err), zap.String("sql", sql))
		c.JSON(http.StatusInternalServerError, utils.ErrorResponseData("Failed to create table: "+err.Error()))
		return
	}

	h.logger.Info("Table created successfully", zap.String("table", req.TableName))

	c.JSON(http.StatusCreated, utils.SuccessResponseData("Table created successfully", gin.H{
		"table_name": req.TableName,
		"sql":        sql,
	}))
}

// AddColumnRequest represents a request to add a column to an existing table
type AddColumnRequest struct {
	Column ColumnRequest `json:"column" binding:"required"`
}

// AddColumn godoc
// @Summary Add column to table (Admin)
// @Description Add a new column to an existing table
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param tableName path string true "Table name"
// @Param request body AddColumnRequest true "Add column request"
// @Success 200 {object} map[string]interface{}
// @Router /admin/database/tables/{tableName}/columns [post]
func (h *TableManagerHandler) AddColumn(c *gin.Context) {
	tableName := c.Param("tableName")
	if !isValidIdentifier(tableName) {
		c.JSON(http.StatusBadRequest, utils.ErrorResponseData("Invalid table name"))
		return
	}

	var req AddColumnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponseData("Invalid request"))
		return
	}

	col := req.Column
	if !isValidIdentifier(col.Name) {
		c.JSON(http.StatusBadRequest, utils.ErrorResponseData("Invalid column name"))
		return
	}

	// Build ALTER TABLE SQL
	sql := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s",
		tableName, col.Name, strings.ToUpper(col.Type))

	if col.Length != nil && (strings.ToUpper(col.Type) == "VARCHAR" || strings.ToUpper(col.Type) == "CHAR") {
		sql = fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s(%d)",
			tableName, col.Name, strings.ToUpper(col.Type), *col.Length)
	}

	if col.NotNull {
		sql += " NOT NULL"
	}

	if col.DefaultValue != nil {
		sql += fmt.Sprintf(" DEFAULT %s", *col.DefaultValue)
	}

	// Execute SQL
	if err := h.db.Exec(sql).Error; err != nil {
		h.logger.Error("Failed to add column", zap.Error(err), zap.String("sql", sql))
		c.JSON(http.StatusInternalServerError, utils.ErrorResponseData("Failed to add column: "+err.Error()))
		return
	}

	h.logger.Info("Column added successfully",
		zap.String("table", tableName),
		zap.String("column", col.Name))

	c.JSON(http.StatusOK, utils.SuccessResponseData("Column added successfully", gin.H{
		"table":  tableName,
		"column": col.Name,
		"sql":    sql,
	}))
}

// DropColumn godoc
// @Summary Drop column from table (Admin)
// @Description Remove a column from an existing table
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param tableName path string true "Table name"
// @Param columnName path string true "Column name"
// @Success 200 {object} map[string]interface{}
// @Router /admin/database/tables/{tableName}/columns/{columnName} [delete]
func (h *TableManagerHandler) DropColumn(c *gin.Context) {
	tableName := c.Param("tableName")
	columnName := c.Param("columnName")

	if !isValidIdentifier(tableName) || !isValidIdentifier(columnName) {
		c.JSON(http.StatusBadRequest, utils.ErrorResponseData("Invalid table or column name"))
		return
	}

	sql := fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", tableName, columnName)

	if err := h.db.Exec(sql).Error; err != nil {
		h.logger.Error("Failed to drop column", zap.Error(err))
		c.JSON(http.StatusInternalServerError, utils.ErrorResponseData("Failed to drop column: "+err.Error()))
		return
	}

	h.logger.Info("Column dropped successfully",
		zap.String("table", tableName),
		zap.String("column", columnName))

	c.JSON(http.StatusOK, utils.SuccessResponseData("Column dropped successfully", gin.H{
		"table":  tableName,
		"column": columnName,
	}))
}

// DropTable godoc
// @Summary Drop table (Admin)
// @Description Delete a table from the database
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param tableName path string true "Table name"
// @Param cascade query boolean false "Drop with CASCADE"
// @Success 200 {object} map[string]interface{}
// @Router /admin/database/tables/{tableName} [delete]
func (h *TableManagerHandler) DropTable(c *gin.Context) {
	tableName := c.Param("tableName")
	cascade := c.Query("cascade") == "true"

	if !isValidIdentifier(tableName) {
		c.JSON(http.StatusBadRequest, utils.ErrorResponseData("Invalid table name"))
		return
	}

	sql := fmt.Sprintf("DROP TABLE %s", tableName)
	if cascade {
		sql += " CASCADE"
	}

	if err := h.db.Exec(sql).Error; err != nil {
		h.logger.Error("Failed to drop table", zap.Error(err))
		c.JSON(http.StatusInternalServerError, utils.ErrorResponseData("Failed to drop table: "+err.Error()))
		return
	}

	h.logger.Info("Table dropped successfully", zap.String("table", tableName))

	c.JSON(http.StatusOK, utils.SuccessResponseData("Table dropped successfully", gin.H{
		"table": tableName,
	}))
}

// RenameTableRequest represents a request to rename a table
type RenameTableRequest struct {
	NewName string `json:"new_name" binding:"required"`
}

// RenameTable godoc
// @Summary Rename table (Admin)
// @Description Rename an existing table
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param tableName path string true "Current table name"
// @Param request body RenameTableRequest true "Rename request"
// @Success 200 {object} map[string]interface{}
// @Router /admin/database/tables/{tableName}/rename [put]
func (h *TableManagerHandler) RenameTable(c *gin.Context) {
	tableName := c.Param("tableName")

	var req RenameTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponseData("Invalid request"))
		return
	}

	if !isValidIdentifier(tableName) || !isValidIdentifier(req.NewName) {
		c.JSON(http.StatusBadRequest, utils.ErrorResponseData("Invalid table name"))
		return
	}

	sql := fmt.Sprintf("ALTER TABLE %s RENAME TO %s", tableName, req.NewName)

	if err := h.db.Exec(sql).Error; err != nil {
		h.logger.Error("Failed to rename table", zap.Error(err))
		c.JSON(http.StatusInternalServerError, utils.ErrorResponseData("Failed to rename table: "+err.Error()))
		return
	}

	h.logger.Info("Table renamed successfully",
		zap.String("old_name", tableName),
		zap.String("new_name", req.NewName))

	c.JSON(http.StatusOK, utils.SuccessResponseData("Table renamed successfully", gin.H{
		"old_name": tableName,
		"new_name": req.NewName,
	}))
}

// isValidIdentifier checks if a string is a valid SQL identifier
func isValidIdentifier(name string) bool {
	if len(name) == 0 || len(name) > 63 {
		return false
	}

	// Must start with letter or underscore
	first := name[0]
	if !(first >= 'a' && first <= 'z') &&
		!(first >= 'A' && first <= 'Z') &&
		first != '_' {
		return false
	}

	// Rest must be alphanumeric or underscore
	for _, ch := range name[1:] {
		if !(ch >= 'a' && ch <= 'z') &&
			!(ch >= 'A' && ch <= 'Z') &&
			!(ch >= '0' && ch <= '9') &&
			ch != '_' {
			return false
		}
	}

	return true
}
