package migration

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"go-mobile-backend-template/internal/services/migration"
	"go-mobile-backend-template/internal/utils"
)

type MigrationHandler struct {
	migrationService interface {
		CreateMigration(ctx context.Context, req *migration.TableAlterRequest) (*migration.Migration, error)
		ExecuteMigration(ctx context.Context, migrationID string) error
		RollbackMigration(ctx context.Context, migrationID string) error
		GetMigrations(ctx context.Context, limit, offset int) ([]*migration.Migration, error)
		GetMigration(ctx context.Context, id string) (*migration.Migration, error)
		GetMigrationHistory(ctx context.Context, tableName string, limit, offset int) ([]*migration.Migration, error)
		GetMigrationFile(ctx context.Context, id string) (string, string, error)
		ValidateMigration(ctx context.Context, id string) (bool, []string, []string, error)
		GetMigrationStatus(ctx context.Context, id string) (*migration.Migration, error)
	}
}

func NewMigrationHandler(db *gorm.DB, config *migration.GoogleScriptsConfig) *MigrationHandler {
	// Use goose-based migration service instead of Google Scripts
	gooseService := migration.NewGooseMigrationService(db, config.MigrationsDir)
	return &MigrationHandler{
		migrationService: gooseService,
	}
}

// CreateMigrationRequest represents the request to create a migration
type CreateMigrationRequest struct {
	TableName   string                   `json:"table_name" binding:"required"`
	Changes     []migration.ColumnChange `json:"changes" binding:"required"`
	RequestedBy string                   `json:"requested_by" binding:"required"`
}

// CreateMigration creates a new migration
func (h *MigrationHandler) CreateMigration(c *gin.Context) {
	var req CreateMigrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request")
		return
	}

	// Create the migration
	migration, err := h.migrationService.CreateMigration(c.Request.Context(), &migration.TableAlterRequest{
		TableName:   req.TableName,
		Changes:     req.Changes,
		RequestedBy: req.RequestedBy,
	})
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create migration")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Migration created successfully", migration)
}

// ExecuteMigration executes a migration
func (h *MigrationHandler) ExecuteMigration(c *gin.Context) {
	migrationID := c.Param("id")
	if migrationID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Migration ID is required")
		return
	}

	// Validate UUID
	if _, err := uuid.Parse(migrationID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid migration ID")
		return
	}

	// Execute the migration
	if err := h.migrationService.ExecuteMigration(c.Request.Context(), migrationID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to execute migration")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Migration executed successfully", nil)
}

// RollbackMigration rolls back a migration
func (h *MigrationHandler) RollbackMigration(c *gin.Context) {
	migrationID := c.Param("id")
	if migrationID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Migration ID is required")
		return
	}

	// Validate UUID
	if _, err := uuid.Parse(migrationID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid migration ID")
		return
	}

	// Rollback the migration
	if err := h.migrationService.RollbackMigration(c.Request.Context(), migrationID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to rollback migration")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Migration rolled back successfully", nil)
}

// GetMigrations retrieves migrations with pagination
func (h *MigrationHandler) GetMigrations(c *gin.Context) {
	// Parse pagination parameters
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// Get migrations
	migrations, err := h.migrationService.GetMigrations(c.Request.Context(), limit, offset)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve migrations")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Migrations retrieved successfully", map[string]interface{}{
		"migrations": migrations,
		"limit":      limit,
		"offset":     offset,
	})
}

// GetMigration retrieves a specific migration
func (h *MigrationHandler) GetMigration(c *gin.Context) {
	migrationID := c.Param("id")
	if migrationID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Migration ID is required")
		return
	}

	// Validate UUID
	if _, err := uuid.Parse(migrationID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid migration ID")
		return
	}

	// Get the migration
	migration, err := h.migrationService.GetMigration(c.Request.Context(), migrationID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Migration not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Migration retrieved successfully", migration)
}

// GetMigrationStatus retrieves the status of a migration
func (h *MigrationHandler) GetMigrationStatus(c *gin.Context) {
	migrationID := c.Param("id")
	if migrationID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Migration ID is required")
		return
	}

	// Validate UUID
	if _, err := uuid.Parse(migrationID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid migration ID")
		return
	}

	// Get the migration
	migration, err := h.migrationService.GetMigration(c.Request.Context(), migrationID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Migration not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Migration status retrieved successfully", map[string]interface{}{
		"id":            migration.ID,
		"status":        migration.Status,
		"error_message": migration.ErrorMessage,
		"created_at":    migration.CreatedAt,
		"completed_at":  migration.CompletedAt,
	})
}

// GetMigrationHistory retrieves migration history for a specific table
func (h *MigrationHandler) GetMigrationHistory(c *gin.Context) {
	tableName := c.Query("table_name")
	if tableName == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Table name is required")
		return
	}

	// Parse pagination parameters
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// Get migrations for the table
	migrations, err := h.migrationService.GetMigrations(c.Request.Context(), limit, offset)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve migration history")
		return
	}

	// Filter by table name
	var tableMigrations []*migration.Migration
	for _, m := range migrations {
		if m.TableName == tableName {
			tableMigrations = append(tableMigrations, m)
		}
	}

	utils.SuccessResponse(c, http.StatusOK, "Migration history retrieved successfully", map[string]interface{}{
		"migrations": tableMigrations,
		"table_name": tableName,
		"limit":      limit,
		"offset":     offset,
	})
}

// GetMigrationFile retrieves the migration file content
func (h *MigrationHandler) GetMigrationFile(c *gin.Context) {
	migrationID := c.Param("id")
	if migrationID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Migration ID is required")
		return
	}

	// Validate UUID
	if _, err := uuid.Parse(migrationID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid migration ID")
		return
	}

	// Get the migration
	migration, err := h.migrationService.GetMigration(c.Request.Context(), migrationID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Migration not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Migration file retrieved successfully", map[string]interface{}{
		"id":           migration.ID,
		"table_name":   migration.TableName,
		"sql_query":    migration.SQLQuery,
		"rollback_sql": migration.RollbackSQL,
		"status":       migration.Status,
	})
}

// ValidateMigration validates a migration before execution
func (h *MigrationHandler) ValidateMigration(c *gin.Context) {
	migrationID := c.Param("id")
	if migrationID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Migration ID is required")
		return
	}

	// Validate UUID
	if _, err := uuid.Parse(migrationID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid migration ID")
		return
	}

	// Get the migration
	migration, err := h.migrationService.GetMigration(c.Request.Context(), migrationID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Migration not found")
		return
	}

	// Basic validation
	validationResult := map[string]interface{}{
		"valid":        true,
		"warnings":     []string{},
		"errors":       []string{},
		"migration_id": migration.ID,
		"table_name":   migration.TableName,
		"status":       migration.Status,
	}

	// Check if migration is in pending status
	if migration.Status != "pending" {
		validationResult["valid"] = false
		validationResult["errors"] = append(validationResult["errors"].([]string),
			"Migration is not in pending status")
	}

	// Check if SQL query is not empty
	if migration.SQLQuery == "" {
		validationResult["valid"] = false
		validationResult["errors"] = append(validationResult["errors"].([]string),
			"Migration SQL query is empty")
	}

	utils.SuccessResponse(c, http.StatusOK, "Migration validation completed", validationResult)
}
