package migration

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MigrationStatus represents the status of a migration
type MigrationStatus string

const (
	StatusPending    MigrationStatus = "pending"
	StatusRunning    MigrationStatus = "running"
	StatusCompleted  MigrationStatus = "completed"
	StatusFailed     MigrationStatus = "failed"
	StatusRolledBack MigrationStatus = "rolled_back"
)

// Migration represents a database migration
type Migration struct {
	ID           string          `json:"id" gorm:"primaryKey"`
	TableName    string          `json:"table_name" gorm:"not null"`
	SQLQuery     string          `json:"sql_query" gorm:"type:text"`
	RollbackSQL  string          `json:"rollback_sql" gorm:"type:text"`
	Status       MigrationStatus `json:"status" gorm:"default:'pending'"`
	ErrorMessage string          `json:"error_message,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	CompletedAt  *time.Time      `json:"completed_at,omitempty"`
	CreatedBy    string          `json:"created_by"`
}

// ColumnChange represents a change to a table column
type ColumnChange struct {
	Action       string `json:"action"` // "add", "modify", "drop", "rename"
	ColumnName   string `json:"column_name"`
	NewName      string `json:"new_name,omitempty"`
	Type         string `json:"type,omitempty"`
	Nullable     bool   `json:"nullable,omitempty"`
	DefaultValue string `json:"default_value,omitempty"`
	IsPrimaryKey bool   `json:"is_primary_key,omitempty"`
	IsForeignKey bool   `json:"is_foreign_key,omitempty"`
	References   string `json:"references,omitempty"`
}

// TableAlterRequest represents a request to alter a table
type TableAlterRequest struct {
	TableName   string         `json:"table_name"`
	Changes     []ColumnChange `json:"changes"`
	RequestedBy string         `json:"requested_by"`
}

// GoogleScriptsConfig holds configuration for Google Scripts integration
type GoogleScriptsConfig struct {
	ScriptURL     string
	AccessToken   string
	ProjectID     string
	MigrationsDir string
}

// GooseMigrationService handles migrations using goose
type GooseMigrationService struct {
	db            *gorm.DB
	migrationsDir string
}

// NewGooseMigrationService creates a new goose migration service
func NewGooseMigrationService(db *gorm.DB, migrationsDir string) *GooseMigrationService {
	return &GooseMigrationService{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

// CreateMigration creates a new migration using goose
func (s *GooseMigrationService) CreateMigration(ctx context.Context, req *TableAlterRequest) (*Migration, error) {
	// Create migration record
	migration := &Migration{
		ID:        uuid.New().String(),
		TableName: req.TableName,
		Status:    StatusPending,
		CreatedAt: time.Now(),
		CreatedBy: req.RequestedBy,
	}

	// Generate migration name with timestamp format like 20250925062029_modify_transaction_table.sql
	timestamp := time.Now().Format("20060102150405")
	migrationName := fmt.Sprintf("%s_modify_%s_table", timestamp, req.TableName)

	// Create migration files using goose
	upSQL, downSQL, err := s.generateMigrationFiles(migrationName, req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate migration files: %w", err)
	}

	migration.SQLQuery = upSQL
	migration.RollbackSQL = downSQL

	// Save migration record
	if err := s.db.Create(migration).Error; err != nil {
		return nil, fmt.Errorf("failed to save migration: %w", err)
	}

	return migration, nil
}

// ExecuteMigration executes a migration by running the SQL directly
func (s *GooseMigrationService) ExecuteMigration(ctx context.Context, migrationID string) error {
	// Get migration record
	var migration Migration
	if err := s.db.Where("id = ?", migrationID).First(&migration).Error; err != nil {
		return fmt.Errorf("migration not found: %w", err)
	}

	// Update status to running
	migration.Status = StatusRunning
	if err := s.db.Save(&migration).Error; err != nil {
		return fmt.Errorf("failed to update migration status: %w", err)
	}

	// Execute migration SQL directly
	if err := s.db.Exec(migration.SQLQuery).Error; err != nil {
		// Update status to failed
		migration.Status = StatusFailed
		migration.ErrorMessage = err.Error()
		s.db.Save(&migration)
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	// Update status to completed
	now := time.Now()
	migration.Status = StatusCompleted
	migration.CompletedAt = &now
	if err := s.db.Save(&migration).Error; err != nil {
		return fmt.Errorf("failed to update migration status: %w", err)
	}

	return nil
}

// RollbackMigration rolls back a migration by running the rollback SQL directly
func (s *GooseMigrationService) RollbackMigration(ctx context.Context, migrationID string) error {
	// Get migration record
	var migration Migration
	if err := s.db.Where("id = ?", migrationID).First(&migration).Error; err != nil {
		return fmt.Errorf("migration not found: %w", err)
	}

	// Update status to running
	migration.Status = StatusRunning
	if err := s.db.Save(&migration).Error; err != nil {
		return fmt.Errorf("failed to update migration status: %w", err)
	}

	// Execute rollback SQL directly
	if err := s.db.Exec(migration.RollbackSQL).Error; err != nil {
		// Update status to failed
		migration.Status = StatusFailed
		migration.ErrorMessage = err.Error()
		s.db.Save(&migration)
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	// Update status to rolled back
	now := time.Now()
	migration.Status = StatusRolledBack
	migration.CompletedAt = &now
	if err := s.db.Save(&migration).Error; err != nil {
		return fmt.Errorf("failed to update migration status: %w", err)
	}

	return nil
}

// generateMigrationFiles generates a single goose migration file with both up and down sections
func (s *GooseMigrationService) generateMigrationFiles(migrationName string, req *TableAlterRequest) (string, string, error) {
	// Ensure migrations directory exists
	if err := os.MkdirAll(s.migrationsDir, 0755); err != nil {
		return "", "", fmt.Errorf("failed to create migrations directory: %w", err)
	}

	// Use timestamp-based naming convention: YYYYMMDDHHMMSS_modify_tablename_table.sql
	migrationFileName := fmt.Sprintf("%s.sql", migrationName)

	// Generate up and down SQL
	upSQL := s.generateUpSQLContent(req)
	downSQL := s.generateDownSQLContent(req)

	// Combine into a single goose migration file
	migrationContent := fmt.Sprintf(`-- +goose Up
-- +goose StatementBegin
%s
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
%s
-- +goose StatementEnd
`, upSQL, downSQL)

	// Write the migration file
	migrationFile := filepath.Join(s.migrationsDir, migrationFileName)
	if err := os.WriteFile(migrationFile, []byte(migrationContent), 0644); err != nil {
		return "", "", fmt.Errorf("failed to write migration file: %w", err)
	}

	return upSQL, downSQL, nil
}

// getNextMigrationNumber gets the next migration number
func (s *GooseMigrationService) getNextMigrationNumber() (int, error) {
	// List existing migration files (both .sql and .up.sql for backwards compatibility)
	files, err := filepath.Glob(filepath.Join(s.migrationsDir, "*.sql"))
	if err != nil {
		return 1, err
	}

	maxNumber := 0
	for _, file := range files {
		// Extract number from filename like "000001_create_users_table.sql" or "000001_create_users_table.up.sql"
		filename := filepath.Base(file)
		if len(filename) >= 6 {
			numberStr := filename[:6]
			if number, err := strconv.Atoi(numberStr); err == nil {
				if number > maxNumber {
					maxNumber = number
				}
			}
		}
	}

	// Return the next number (maxNumber is already the actual number like 4, so next is 5)
	// Since we have 000001-000004, the next should be 5
	// But we need to return 5, not 5, since the existing files are 000001-000004
	if maxNumber < 4 {
		return 5, nil
	}
	return maxNumber + 1, nil
}

// generateUpSQLContent generates the up migration SQL content (without goose directives)
func (s *GooseMigrationService) generateUpSQLContent(req *TableAlterRequest) string {
	var sql string

	for _, change := range req.Changes {
		switch change.Action {
		case "add":
			sql += s.generateAddColumnSQL(req.TableName, change)
		case "modify":
			sql += s.generateModifyColumnSQL(req.TableName, change)
		case "drop":
			sql += s.generateDropColumnSQL(req.TableName, change)
		case "rename":
			sql += s.generateRenameColumnSQL(req.TableName, change)
		}
	}

	return sql
}

// generateDownSQLContent generates the down migration SQL content (without goose directives)
func (s *GooseMigrationService) generateDownSQLContent(req *TableAlterRequest) string {
	var sql string

	// Reverse the changes
	for i := len(req.Changes) - 1; i >= 0; i-- {
		change := req.Changes[i]
		switch change.Action {
		case "add":
			sql += s.generateDropColumnSQL(req.TableName, change)
		case "modify":
			// For modify, we'd need to know the original column definition
			// This is a simplified version
			sql += fmt.Sprintf("-- TODO: Restore original column definition for %s\n", change.ColumnName)
		case "drop":
			// For drop, we'd need to know the original column definition
			sql += fmt.Sprintf("-- TODO: Restore dropped column %s\n", change.ColumnName)
		case "rename":
			sql += s.generateRenameColumnSQL(req.TableName, ColumnChange{
				Action:     "rename",
				ColumnName: change.NewName,
				NewName:    change.ColumnName,
			})
		}
	}

	return sql
}

// generateAddColumnSQL generates SQL for adding a column
func (s *GooseMigrationService) generateAddColumnSQL(tableName string, change ColumnChange) string {
	sql := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", tableName, change.ColumnName, change.Type)

	if !change.Nullable {
		sql += " NOT NULL"
	}

	if change.DefaultValue != "" {
		sql += fmt.Sprintf(" DEFAULT %s", change.DefaultValue)
	}

	sql += ";\n"
	return sql
}

// generateModifyColumnSQL generates SQL for modifying a column
func (s *GooseMigrationService) generateModifyColumnSQL(tableName string, change ColumnChange) string {
	sql := fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s TYPE %s", tableName, change.ColumnName, change.Type)

	if !change.Nullable {
		sql += " NOT NULL"
	}

	if change.DefaultValue != "" {
		sql += fmt.Sprintf(" DEFAULT %s", change.DefaultValue)
	}

	sql += ";\n"
	return sql
}

// generateDropColumnSQL generates SQL for dropping a column
func (s *GooseMigrationService) generateDropColumnSQL(tableName string, change ColumnChange) string {
	return fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s;\n", tableName, change.ColumnName)
}

// generateRenameColumnSQL generates SQL for renaming a column
func (s *GooseMigrationService) generateRenameColumnSQL(tableName string, change ColumnChange) string {
	return fmt.Sprintf("ALTER TABLE %s RENAME COLUMN %s TO %s;\n", tableName, change.ColumnName, change.NewName)
}

// runGooseUp runs goose up command
func (s *GooseMigrationService) runGooseUp() error {
	cmd := exec.Command("goose", "-dir", s.migrationsDir, "postgres", s.getDSN(), "up")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("goose up failed: %s, output: %s", err, string(output))
	}
	return nil
}

// runGooseDown runs goose down command
func (s *GooseMigrationService) runGooseDown() error {
	cmd := exec.Command("goose", "-dir", s.migrationsDir, "postgres", s.getDSN(), "down")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("goose down failed: %s, output: %s", err, string(output))
	}
	return nil
}

// getDSN returns the database connection string
func (s *GooseMigrationService) getDSN() string {
	// For now, return a simple DSN - this should be improved to get actual connection details
	return "postgres://app:secret@localhost:5433/myapp?sslmode=disable"
}

// GetMigrations gets all migrations with pagination
func (s *GooseMigrationService) GetMigrations(ctx context.Context, limit, offset int) ([]*Migration, error) {
	var migrations []*Migration
	query := s.db.Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&migrations).Error; err != nil {
		return nil, fmt.Errorf("failed to get migrations: %w", err)
	}

	return migrations, nil
}

// GetMigration gets a specific migration by ID
func (s *GooseMigrationService) GetMigration(ctx context.Context, id string) (*Migration, error) {
	var migration Migration
	if err := s.db.Where("id = ?", id).First(&migration).Error; err != nil {
		return nil, fmt.Errorf("migration not found: %w", err)
	}
	return &migration, nil
}

// GetMigrationHistory gets migration history for a specific table
func (s *GooseMigrationService) GetMigrationHistory(ctx context.Context, tableName string, limit, offset int) ([]*Migration, error) {
	var migrations []*Migration
	query := s.db.Where("table_name = ?", tableName).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&migrations).Error; err != nil {
		return nil, fmt.Errorf("failed to get migration history: %w", err)
	}

	return migrations, nil
}

// GetMigrationFile gets migration file content
func (s *GooseMigrationService) GetMigrationFile(ctx context.Context, id string) (string, string, error) {
	var migration Migration
	if err := s.db.Where("id = ?", id).First(&migration).Error; err != nil {
		return "", "", fmt.Errorf("migration not found: %w", err)
	}
	return migration.SQLQuery, migration.RollbackSQL, nil
}

// ValidateMigration validates a migration
func (s *GooseMigrationService) ValidateMigration(ctx context.Context, id string) (bool, []string, []string, error) {
	var migration Migration
	if err := s.db.Where("id = ?", id).First(&migration).Error; err != nil {
		return false, nil, nil, fmt.Errorf("migration not found: %w", err)
	}

	// Simple validation - check if migration is in pending or running state
	if migration.Status == StatusPending || migration.Status == StatusRunning {
		return true, []string{}, []string{}, nil
	}

	return false, []string{}, []string{"Migration is not in a valid state for execution"}, nil
}

// GetMigrationStatus gets migration status
func (s *GooseMigrationService) GetMigrationStatus(ctx context.Context, id string) (*Migration, error) {
	return s.GetMigration(ctx, id)
}
