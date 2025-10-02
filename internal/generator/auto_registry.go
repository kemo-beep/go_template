package generator

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AutoRegistry manages automatic API registration and regeneration
type AutoRegistry struct {
	db             *gorm.DB
	logger         *zap.Logger
	config         *GeneratorConfig
	generator      *APIGenerator
	watcher        *SchemaWatcher
	router         *gin.RouterGroup
	registeredAPIs map[string]bool
	mu             sync.RWMutex
}

// NewAutoRegistry creates a new auto registry
func NewAutoRegistry(db *gorm.DB, logger *zap.Logger, config *GeneratorConfig) *AutoRegistry {
	return &AutoRegistry{
		db:             db,
		logger:         logger,
		config:         config,
		generator:      NewAPIGenerator(db, logger, config),
		registeredAPIs: make(map[string]bool),
	}
}

// Initialize sets up the auto registry with a router
func (ar *AutoRegistry) Initialize(router *gin.RouterGroup) error {
	ar.router = router
	ar.watcher = NewSchemaWatcher(ar.db, ar.logger, ar.config)

	// Generate initial APIs
	if err := ar.regenerateAPIs(); err != nil {
		return fmt.Errorf("failed to generate initial APIs: %w", err)
	}

	// Register all generated APIs
	if err := ar.registerAllAPIs(); err != nil {
		return fmt.Errorf("failed to register initial APIs: %w", err)
	}

	// Start schema watcher
	ctx := context.Background()
	if err := ar.watcher.Start(ctx, ar.onSchemaChange); err != nil {
		return fmt.Errorf("failed to start schema watcher: %w", err)
	}

	ar.logger.Info("Auto registry initialized successfully")
	return nil
}

// onSchemaChange is called when the database schema changes
func (ar *AutoRegistry) onSchemaChange() error {
	ar.logger.Info("Schema change detected, regenerating APIs...")

	// Regenerate APIs
	if err := ar.regenerateAPIs(); err != nil {
		return fmt.Errorf("failed to regenerate APIs: %w", err)
	}

	// Re-register APIs
	if err := ar.registerAllAPIs(); err != nil {
		return fmt.Errorf("failed to re-register APIs: %w", err)
	}

	// Generate TypeScript types for frontend
	if err := ar.generateTypeScriptTypes(); err != nil {
		ar.logger.Warn("Failed to generate TypeScript types", zap.Error(err))
	}

	ar.logger.Info("APIs regenerated and re-registered successfully")
	return nil
}

// regenerateAPIs regenerates all API files
func (ar *AutoRegistry) regenerateAPIs() error {
	ar.logger.Info("Regenerating API files...")

	// Generate all APIs
	if err := ar.generator.GenerateAll(); err != nil {
		return fmt.Errorf("failed to generate APIs: %w", err)
	}

	// Regenerate Swagger documentation
	if err := ar.regenerateSwaggerDocs(); err != nil {
		return fmt.Errorf("failed to regenerate Swagger docs: %w", err)
	}

	ar.logger.Info("API files regenerated successfully")
	return nil
}

// registerAllAPIs registers all generated APIs
func (ar *AutoRegistry) registerAllAPIs() error {
	ar.mu.Lock()
	defer ar.mu.Unlock()

	ar.logger.Info("Registering generated APIs...")

	// Clear existing registrations
	ar.registeredAPIs = make(map[string]bool)

	// Find all generated API directories
	apiDir := "internal/api/v1"
	entries, err := os.ReadDir(apiDir)
	if err != nil {
		return fmt.Errorf("failed to read API directory: %w", err)
	}

	registeredCount := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Skip non-generated directories
		if ar.isSystemDirectory(entry.Name()) {
			continue
		}

		// Check if this directory has a routes.go file
		routesFile := filepath.Join(apiDir, entry.Name(), "routes.go")
		if _, err := os.Stat(routesFile); os.IsNotExist(err) {
			continue
		}

		// Register this API
		if err := ar.registerAPI(entry.Name()); err != nil {
			ar.logger.Error("Failed to register API",
				zap.String("api", entry.Name()),
				zap.Error(err))
			continue
		}

		ar.registeredAPIs[entry.Name()] = true
		registeredCount++
	}

	ar.logger.Info("API registration completed",
		zap.Int("registered_apis", registeredCount),
		zap.Int("total_apis", len(ar.registeredAPIs)))

	return nil
}

// registerAPI registers a specific API
func (ar *AutoRegistry) registerAPI(apiName string) error {
	// This is a simplified version - in practice, you'd need to dynamically
	// import and call the RegisterRoutes function for each API
	ar.logger.Debug("Registering API", zap.String("api", apiName))

	// For now, we'll rely on the generated_routes.go file to handle registration
	// In a more sophisticated implementation, you'd dynamically load and call
	// the RegisterRoutes function for each API package
	return nil
}

// isSystemDirectory checks if a directory is a system directory (not generated)
func (ar *AutoRegistry) isSystemDirectory(name string) bool {
	systemDirs := map[string]bool{
		"admin":     true,
		"auth":      true,
		"files":     true,
		"migration": true,
		"realtime":  true,
		"users":     true,
	}
	return systemDirs[name]
}

// regenerateSwaggerDocs regenerates Swagger documentation
func (ar *AutoRegistry) regenerateSwaggerDocs() error {
	ar.logger.Info("Regenerating Swagger documentation...")

	// Run swag init command
	cmd := exec.Command("swag", "init", "-g", "cmd/server/main.go", "-o", "docs", "--parseInternal")
	cmd.Dir = "."

	output, err := cmd.CombinedOutput()
	if err != nil {
		ar.logger.Error("Failed to regenerate Swagger docs",
			zap.Error(err),
			zap.String("output", string(output)))
		return fmt.Errorf("swag init failed: %w", err)
	}

	ar.logger.Info("Swagger documentation regenerated successfully")
	return nil
}

// generateTypeScriptTypes generates TypeScript types for the frontend
func (ar *AutoRegistry) generateTypeScriptTypes() error {
	ar.logger.Info("Generating TypeScript types...")

	// Create TypeScript type generator
	tsGenerator := NewTypeScriptGenerator(ar.db, ar.logger, ar.config)

	// Generate types
	if err := tsGenerator.GenerateAll(); err != nil {
		return fmt.Errorf("failed to generate TypeScript types: %w", err)
	}

	ar.logger.Info("TypeScript types generated successfully")
	return nil
}

// GetRegisteredAPIs returns a list of currently registered APIs
func (ar *AutoRegistry) GetRegisteredAPIs() []string {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	apis := make([]string, 0, len(ar.registeredAPIs))
	for api := range ar.registeredAPIs {
		apis = append(apis, api)
	}
	return apis
}

// IsAPIRegistered checks if a specific API is registered
func (ar *AutoRegistry) IsAPIRegistered(apiName string) bool {
	ar.mu.RLock()
	defer ar.mu.RUnlock()
	return ar.registeredAPIs[apiName]
}

// Stop stops the auto registry
func (ar *AutoRegistry) Stop() {
	if ar.watcher != nil {
		ar.watcher.Stop()
	}
	ar.logger.Info("Auto registry stopped")
}

// GetStatus returns the current status of the auto registry
func (ar *AutoRegistry) GetStatus() map[string]interface{} {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	status := map[string]interface{}{
		"registered_apis": len(ar.registeredAPIs),
		"apis":            ar.registeredAPIs,
		"watcher_running": false,
		"last_checksum":   "",
	}

	if ar.watcher != nil {
		status["watcher_running"] = ar.watcher.IsRunning()
		status["last_checksum"] = ar.watcher.GetLastChecksum()
	}

	return status
}
