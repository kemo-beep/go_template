package main

import (
	"flag"
	"fmt"
	"log"

	"go-mobile-backend-template/internal/db"
	"go-mobile-backend-template/internal/generator"
	"go-mobile-backend-template/pkg/config"

	"go.uber.org/zap"
)

func main() {
	var (
		tableName = flag.String("table", "", "Generate APIs for specific table only")
		outputDir = flag.String("output", "./generated", "Output directory for generated files")
		verbose   = flag.Bool("verbose", false, "Enable verbose logging")
	)
	flag.Parse()

	// Load configuration
	cfg := config.Load()

	// Setup logger
	var logger *zap.Logger
	var err error

	if *verbose {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		log.Fatal("Failed to create logger:", err)
	}
	defer logger.Sync()

	// Connect to database
	dbConn, err := db.Connect(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Create generator
	gen := generator.NewAPIGeneratorMain(dbConn, logger, cfg)

	// Override output directory if specified
	if *outputDir != "" {
		genConfig := gen.GetConfig()
		genConfig.OutputDir = *outputDir
		if err := gen.UpdateConfig(genConfig); err != nil {
			logger.Fatal("Failed to update config", zap.Error(err))
		}
	}

	// Generate APIs
	if *tableName != "" {
		// Generate for specific table
		logger.Info("Generating APIs for specific table", zap.String("table", *tableName))

		router, err := gen.GenerateForTable(*tableName)
		if err != nil {
			logger.Fatal("Failed to generate APIs for table",
				zap.String("table", *tableName),
				zap.Error(err))
		}

		if router != nil {
			logger.Info("API generation completed for table",
				zap.String("table", *tableName),
				zap.String("output", *outputDir))
		}
	} else {
		// Generate for all tables using file-based approach
		logger.Info("Generating APIs for all tables using file-based approach")

		err := gen.GenerateAll()
		if err != nil {
			logger.Fatal("Failed to generate APIs", zap.Error(err))
		}

		logger.Info("File-based API generation completed", zap.String("output", *outputDir))
	}

	// Get generated endpoints info
	endpoints, err := gen.GetGeneratedEndpoints()
	if err != nil {
		logger.Error("Failed to get generated endpoints info", zap.Error(err))
	} else {
		logger.Info("Generated endpoints", zap.Int("count", len(endpoints)))

		// Print endpoint summary
		fmt.Println("\n📋 Generated API Endpoints:")
		fmt.Println("================================")

		currentTable := ""
		for _, endpoint := range endpoints {
			if endpoint.Tags[0] != currentTable {
				currentTable = endpoint.Tags[0]
				fmt.Printf("\n🏷️  %s:\n", currentTable)
			}
			fmt.Printf("  %-6s %s\n", endpoint.Method, endpoint.Path)
		}
	}

	// Get table info
	tables, err := gen.GetTableInfo()
	if err != nil {
		logger.Error("Failed to get table info", zap.Error(err))
	} else {
		fmt.Printf("\n📊 Discovered Tables: %d\n", len(tables))
		fmt.Println("================================")
		for _, table := range tables {
			fmt.Printf("  • %s (%d columns)\n", table.Name, len(table.Columns))
		}
	}

	fmt.Println("\n✅ Auto API generation completed successfully!")
	fmt.Printf("📁 Generated files saved to: %s\n", *outputDir)
	fmt.Println("\n🚀 You can now use the generated APIs in your application!")
}
