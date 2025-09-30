package database

import (
	"fmt"

	"go-mobile-backend-template/internal/db/repository"
	"go-mobile-backend-template/pkg/config"
)

// RunMigrations runs database migrations
func RunMigrations(cfg config.Database) error {
	// Connect to database using GORM
	db, err := Connect(cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate the schema
	if err := db.AutoMigrate(
		&repository.User{},
		&repository.File{},
		&repository.RefreshToken{},
	); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
