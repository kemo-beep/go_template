package generator

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SchemaWatcher monitors database schema changes and triggers API regeneration
type SchemaWatcher struct {
	db           *gorm.DB
	logger       *zap.Logger
	config       *GeneratorConfig
	lastChecksum string
	mu           sync.RWMutex
	stopChan     chan struct{}
	isRunning    bool
	onChange     func() error // Callback function when schema changes
}

// NewSchemaWatcher creates a new schema watcher
func NewSchemaWatcher(db *gorm.DB, logger *zap.Logger, config *GeneratorConfig) *SchemaWatcher {
	return &SchemaWatcher{
		db:       db,
		logger:   logger,
		config:   config,
		stopChan: make(chan struct{}),
	}
}

// Start begins watching for schema changes
func (sw *SchemaWatcher) Start(ctx context.Context, onChange func() error) error {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	if sw.isRunning {
		return fmt.Errorf("schema watcher is already running")
	}

	sw.onChange = onChange
	sw.isRunning = true

	// Get initial checksum
	checksum, err := sw.getSchemaChecksum()
	if err != nil {
		sw.logger.Error("Failed to get initial schema checksum", zap.Error(err))
		return err
	}
	sw.lastChecksum = checksum

	sw.logger.Info("Starting schema watcher", zap.String("interval", sw.config.AutoRegistration.WatchInterval.String()))

	go sw.watchLoop(ctx)
	return nil
}

// Stop stops the schema watcher
func (sw *SchemaWatcher) Stop() {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	if !sw.isRunning {
		return
	}

	close(sw.stopChan)
	sw.isRunning = false
	sw.logger.Info("Schema watcher stopped")
}

// watchLoop runs the main watching loop
func (sw *SchemaWatcher) watchLoop(ctx context.Context) {
	ticker := time.NewTicker(sw.config.AutoRegistration.WatchInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			sw.logger.Info("Schema watcher context cancelled")
			return
		case <-sw.stopChan:
			sw.logger.Info("Schema watcher stopped")
			return
		case <-ticker.C:
			if err := sw.checkForChanges(); err != nil {
				sw.logger.Error("Error checking for schema changes", zap.Error(err))
			}
		}
	}
}

// checkForChanges checks if the database schema has changed
func (sw *SchemaWatcher) checkForChanges() error {
	sw.mu.RLock()
	defer sw.mu.RUnlock()

	if !sw.isRunning {
		return nil
	}

	currentChecksum, err := sw.getSchemaChecksum()
	if err != nil {
		return fmt.Errorf("failed to get current schema checksum: %w", err)
	}

	if currentChecksum != sw.lastChecksum {
		sw.logger.Info("Schema change detected",
			zap.String("old_checksum", sw.lastChecksum),
			zap.String("new_checksum", currentChecksum))

		// Update the checksum
		sw.mu.Lock()
		sw.lastChecksum = currentChecksum
		sw.mu.Unlock()

		// Trigger regeneration
		if sw.onChange != nil {
			if err := sw.onChange(); err != nil {
				sw.logger.Error("Failed to regenerate APIs after schema change", zap.Error(err))
				return err
			}
			sw.logger.Info("APIs regenerated successfully after schema change")
		}
	}

	return nil
}

// getSchemaChecksum generates a checksum of the current database schema
func (sw *SchemaWatcher) getSchemaChecksum() (string, error) {
	// Get all table names
	var tables []string
	if err := sw.db.Raw(`
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_type = 'BASE TABLE'
		ORDER BY table_name
	`).Scan(&tables).Error; err != nil {
		return "", err
	}

	// Get schema information for each table
	var checksumData string
	for _, table := range tables {
		// Get column information
		var columns []struct {
			ColumnName    string  `db:"column_name"`
			DataType      string  `db:"data_type"`
			IsNullable    string  `db:"is_nullable"`
			ColumnDefault *string `db:"column_default"`
			MaxLength     *int    `db:"character_maximum_length"`
		}

		if err := sw.db.Raw(`
			SELECT 
				column_name,
				data_type,
				is_nullable,
				column_default,
				character_maximum_length
			FROM information_schema.columns 
			WHERE table_name = ? 
			ORDER BY ordinal_position
		`, table).Scan(&columns).Error; err != nil {
			return "", err
		}

		// Add table and column info to checksum data
		checksumData += table + ":"
		for _, col := range columns {
			checksumData += fmt.Sprintf("%s:%s:%s:", col.ColumnName, col.DataType, col.IsNullable)
			if col.ColumnDefault != nil {
				checksumData += *col.ColumnDefault + ":"
			}
			if col.MaxLength != nil {
				checksumData += fmt.Sprintf("%d:", *col.MaxLength)
			}
		}
		checksumData += ";"
	}

	// Generate simple checksum (in production, use proper hash)
	return fmt.Sprintf("%x", len(checksumData)), nil
}

// IsRunning returns whether the watcher is currently running
func (sw *SchemaWatcher) IsRunning() bool {
	sw.mu.RLock()
	defer sw.mu.RUnlock()
	return sw.isRunning
}

// GetLastChecksum returns the last known schema checksum
func (sw *SchemaWatcher) GetLastChecksum() string {
	sw.mu.RLock()
	defer sw.mu.RUnlock()
	return sw.lastChecksum
}
