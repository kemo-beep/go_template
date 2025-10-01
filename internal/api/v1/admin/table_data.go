package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TableDataHandler struct {
	db *gorm.DB
}

func NewTableDataHandler(db *gorm.DB) *TableDataHandler {
	return &TableDataHandler{db: db}
}

// InsertTableRow inserts a new row into a table
// @Summary Insert a row into a table
// @Description Insert a new row into the specified table
// @Tags Database
// @Accept json
// @Produce json
// @Param tableName path string true "Table name"
// @Param row body map[string]interface{} true "Row data"
// @Success 201 {object} map[string]interface{}
// @Router /admin/database/tables/{tableName}/rows [post]
func (h *TableDataHandler) InsertTableRow(c *gin.Context) {
	tableName := c.Param("tableName")

	var rowData map[string]interface{}
	if err := c.ShouldBindJSON(&rowData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid request body"})
		return
	}

	// Build column names and placeholders
	columns := []string{}
	values := []interface{}{}

	for col, val := range rowData {
		columns = append(columns, col)
		values = append(values, val)
	}

	// Build INSERT query
	query := fmt.Sprintf("INSERT INTO %s (", tableName)
	for i, col := range columns {
		if i > 0 {
			query += ", "
		}
		query += col
	}
	query += ") VALUES ("
	for i := range columns {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("$%d", i+1)
	}
	query += ") RETURNING *"

	// Execute query
	result := make(map[string]interface{})
	if err := h.db.Raw(query, values...).Scan(&result).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": fmt.Sprintf("Failed to insert row: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Row inserted successfully",
		"data":    result,
	})
}

// UpdateTableRow updates a row in a table
// @Summary Update a row in a table
// @Description Update an existing row in the specified table
// @Tags Database
// @Accept json
// @Produce json
// @Param tableName path string true "Table name"
// @Param pkValue path string true "Primary key value"
// @Param row body map[string]interface{} true "Row data"
// @Success 200 {object} map[string]interface{}
// @Router /admin/database/tables/{tableName}/rows/{pkValue} [put]
func (h *TableDataHandler) UpdateTableRow(c *gin.Context) {
	tableName := c.Param("tableName")
	pkValue := c.Param("pkValue")

	var rowData map[string]interface{}
	if err := c.ShouldBindJSON(&rowData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid request body"})
		return
	}

	// Get primary key column name
	var pkColumn string
	err := h.db.Raw(`
		SELECT a.attname
		FROM pg_index i
		JOIN pg_attribute a ON a.attrelid = i.indrelid AND a.attnum = ANY(i.indkey)
		WHERE i.indrelid = ?::regclass AND i.indisprimary
	`, tableName).Scan(&pkColumn).Error

	if err != nil || pkColumn == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Could not determine primary key"})
		return
	}

	// Build UPDATE query
	query := fmt.Sprintf("UPDATE %s SET ", tableName)
	values := []interface{}{}
	paramIndex := 1

	first := true
	for col, val := range rowData {
		if col == pkColumn {
			continue // Don't update primary key
		}
		if !first {
			query += ", "
		}
		query += fmt.Sprintf("%s = $%d", col, paramIndex)
		values = append(values, val)
		paramIndex++
		first = false
	}

	query += fmt.Sprintf(" WHERE %s = $%d RETURNING *", pkColumn, paramIndex)
	values = append(values, pkValue)

	// Execute query
	result := make(map[string]interface{})
	if err := h.db.Raw(query, values...).Scan(&result).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": fmt.Sprintf("Failed to update row: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Row updated successfully",
		"data":    result,
	})
}

// DeleteTableRow deletes a row from a table
// @Summary Delete a row from a table
// @Description Delete an existing row from the specified table
// @Tags Database
// @Accept json
// @Produce json
// @Param tableName path string true "Table name"
// @Param pkValue path string true "Primary key value"
// @Success 200 {object} map[string]interface{}
// @Router /admin/database/tables/{tableName}/rows/{pkValue} [delete]
func (h *TableDataHandler) DeleteTableRow(c *gin.Context) {
	tableName := c.Param("tableName")
	pkValue := c.Param("pkValue")

	// Get primary key column name
	var pkColumn string
	err := h.db.Raw(`
		SELECT a.attname
		FROM pg_index i
		JOIN pg_attribute a ON a.attrelid = i.indrelid AND a.attnum = ANY(i.indkey)
		WHERE i.indrelid = ?::regclass AND i.indisprimary
	`, tableName).Scan(&pkColumn).Error

	if err != nil || pkColumn == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Could not determine primary key"})
		return
	}

	// Build DELETE query
	query := fmt.Sprintf("DELETE FROM %s WHERE %s = $1", tableName, pkColumn)

	// Execute query
	result := h.db.Exec(query, pkValue)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": fmt.Sprintf("Failed to delete row: %v", result.Error)})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Row not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Row deleted successfully",
	})
}
