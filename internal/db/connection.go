package db

import (
	"fmt"
	"time"

	"go-mobile-backend-template/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect establishes a database connection
func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.Database.DSN()

	// Configure GORM logger
	var gormLogger logger.Interface
	switch cfg.Logging.Level {
	case "debug":
		gormLogger = logger.Default.LogMode(logger.Info)
	case "info":
		gormLogger = logger.Default.LogMode(logger.Warn)
	default:
		gormLogger = logger.Default.LogMode(logger.Error)
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB for connection pool configuration
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// Close closes the database connection
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
