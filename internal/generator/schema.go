package generator

import (
	"database/sql"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// TableInfo represents information about a database table
type TableInfo struct {
	Name        string           `json:"name"`
	Schema      string           `json:"schema"`
	Columns     []ColumnInfo     `json:"columns"`
	Indexes     []IndexInfo      `json:"indexes"`
	ForeignKeys []ForeignKeyInfo `json:"foreign_keys"`
	Constraints []ConstraintInfo `json:"constraints"`
	Comment     string           `json:"comment"`
}

// ColumnInfo represents information about a table column
type ColumnInfo struct {
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	GoType       string         `json:"go_type"`
	TSType       string         `json:"ts_type"`
	IsNullable   bool           `json:"is_nullable"`
	IsPrimaryKey bool           `json:"is_primary_key"`
	IsUnique     bool           `json:"is_unique"`
	IsForeignKey bool           `json:"is_foreign_key"`
	DefaultValue *string        `json:"default_value,omitempty"`
	MaxLength    *int           `json:"max_length,omitempty"`
	Precision    *int           `json:"precision,omitempty"`
	Scale        *int           `json:"scale,omitempty"`
	Comment      string         `json:"comment"`
	References   *ForeignKeyRef `json:"references,omitempty"`
}

// IndexInfo represents information about a table index
type IndexInfo struct {
	Name    string   `json:"name"`
	Columns []string `json:"columns"`
	Unique  bool     `json:"unique"`
	Type    string   `json:"type"`
}

// ForeignKeyInfo represents foreign key information
type ForeignKeyInfo struct {
	Column     string `json:"column"`
	References string `json:"references"`
	RefTable   string `json:"ref_table"`
	RefColumn  string `json:"ref_column"`
	OnDelete   string `json:"on_delete"`
	OnUpdate   string `json:"on_update"`
}

// ConstraintInfo represents constraint information
type ConstraintInfo struct {
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	Columns   []string `json:"columns"`
	Check     string   `json:"check,omitempty"`
	Reference string   `json:"reference,omitempty"`
}

// ForeignKeyRef represents a foreign key reference
type ForeignKeyRef struct {
	Table  string `json:"table"`
	Column string `json:"column"`
}

// SchemaAnalyzer analyzes database schema and extracts table information
type SchemaAnalyzer struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewSchemaAnalyzer creates a new schema analyzer
func NewSchemaAnalyzer(db *gorm.DB, logger *zap.Logger) *SchemaAnalyzer {
	return &SchemaAnalyzer{
		db:     db,
		logger: logger,
	}
}

// DiscoverTables discovers all tables in the database
func (sa *SchemaAnalyzer) DiscoverTables() ([]*TableInfo, error) {
	var tables []*TableInfo

	// Get all tables from information_schema
	rows, err := sa.db.Raw(`
		SELECT table_name, table_schema, COALESCE(obj_description(c.oid), '') as table_comment
		FROM information_schema.tables t
		LEFT JOIN pg_class c ON c.relname = t.table_name
		WHERE t.table_schema = 'public' 
		AND t.table_type = 'BASE TABLE'
		ORDER BY table_name
	`).Rows()

	if err != nil {
		return nil, fmt.Errorf("failed to discover tables: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tableName, tableSchema, tableComment string
		if err := rows.Scan(&tableName, &tableSchema, &tableComment); err != nil {
			sa.logger.Error("Failed to scan table row", zap.Error(err))
			continue
		}

		// Skip system tables
		if sa.isSystemTable(tableName) {
			continue
		}

		tableInfo, err := sa.analyzeTable(tableName, tableSchema, tableComment)
		if err != nil {
			sa.logger.Error("Failed to analyze table",
				zap.String("table", tableName),
				zap.Error(err))
			continue
		}

		tables = append(tables, tableInfo)
	}

	sa.logger.Info("Discovered tables", zap.Int("count", len(tables)))
	return tables, nil
}

// analyzeTable analyzes a specific table and extracts its schema information
func (sa *SchemaAnalyzer) analyzeTable(tableName, schema, comment string) (*TableInfo, error) {
	tableInfo := &TableInfo{
		Name:    tableName,
		Schema:  schema,
		Comment: comment,
	}

	// Get columns
	columns, err := sa.getColumns(tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get columns for table %s: %w", tableName, err)
	}
	tableInfo.Columns = columns

	// Get indexes
	indexes, err := sa.getIndexes(tableName)
	if err != nil {
		sa.logger.Warn("Failed to get indexes", zap.String("table", tableName), zap.Error(err))
	}
	tableInfo.Indexes = indexes

	// Get foreign keys
	foreignKeys, err := sa.getForeignKeys(tableName)
	if err != nil {
		sa.logger.Warn("Failed to get foreign keys", zap.String("table", tableName), zap.Error(err))
	}
	tableInfo.ForeignKeys = foreignKeys

	// Get constraints
	constraints, err := sa.getConstraints(tableName)
	if err != nil {
		sa.logger.Warn("Failed to get constraints", zap.String("table", tableName), zap.Error(err))
	}
	tableInfo.Constraints = constraints

	return tableInfo, nil
}

// getColumns retrieves column information for a table
func (sa *SchemaAnalyzer) getColumns(tableName string) ([]ColumnInfo, error) {
	rows, err := sa.db.Raw(`
		SELECT 
			c.column_name,
			c.data_type,
			c.is_nullable,
			c.column_default,
			c.character_maximum_length,
			c.numeric_precision,
			c.numeric_scale,
			COALESCE(col_description(pgc.oid, c.ordinal_position), '') as column_comment,
			CASE WHEN pk.column_name IS NOT NULL THEN true ELSE false END as is_primary_key,
			CASE WHEN u.column_name IS NOT NULL THEN true ELSE false END as is_unique
		FROM information_schema.columns c
		LEFT JOIN pg_class pgc ON pgc.relname = c.table_name
		LEFT JOIN (
			SELECT ku.column_name
			FROM information_schema.table_constraints tc
			JOIN information_schema.key_column_usage ku ON tc.constraint_name = ku.constraint_name
			WHERE tc.table_name = ? AND tc.constraint_type = 'PRIMARY KEY'
		) pk ON c.column_name = pk.column_name
		LEFT JOIN (
			SELECT ku.column_name
			FROM information_schema.table_constraints tc
			JOIN information_schema.key_column_usage ku ON tc.constraint_name = ku.constraint_name
			WHERE tc.table_name = ? AND tc.constraint_type = 'UNIQUE'
		) u ON c.column_name = u.column_name
		WHERE c.table_name = ?
		ORDER BY c.ordinal_position
	`, tableName, tableName, tableName).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []ColumnInfo
	for rows.Next() {
		var col ColumnInfo
		var isNullable, isPrimaryKey, isUnique string
		var maxLength, precision, scale sql.NullInt64
		var defaultValue, comment sql.NullString

		err := rows.Scan(
			&col.Name,
			&col.Type,
			&isNullable,
			&defaultValue,
			&maxLength,
			&precision,
			&scale,
			&comment,
			&isPrimaryKey,
			&isUnique,
		)
		if err != nil {
			return nil, err
		}

		col.IsNullable = isNullable == "YES"
		col.IsPrimaryKey = isPrimaryKey == "true"
		col.IsUnique = isUnique == "true"

		if defaultValue.Valid {
			col.DefaultValue = &defaultValue.String
		}
		if maxLength.Valid {
			val := int(maxLength.Int64)
			col.MaxLength = &val
		}
		if precision.Valid {
			val := int(precision.Int64)
			col.Precision = &val
		}
		if scale.Valid {
			val := int(scale.Int64)
			col.Scale = &val
		}
		if comment.Valid {
			col.Comment = comment.String
		}

		// Map database types to Go and TypeScript types
		col.GoType = sa.mapToGoType(col.Type, col.IsNullable)
		col.TSType = sa.mapToTSType(col.Type, col.IsNullable)

		columns = append(columns, col)
	}

	return columns, nil
}

// getIndexes retrieves index information for a table
func (sa *SchemaAnalyzer) getIndexes(tableName string) ([]IndexInfo, error) {
	rows, err := sa.db.Raw(`
		SELECT 
			i.indexname,
			i.indexdef,
			CASE WHEN i.indexdef LIKE '%UNIQUE%' THEN true ELSE false END as is_unique
		FROM pg_indexes i
		WHERE i.tablename = ?
		ORDER BY i.indexname
	`, tableName).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indexes []IndexInfo
	for rows.Next() {
		var idx IndexInfo
		var indexDef string
		var isUnique bool

		err := rows.Scan(&idx.Name, &indexDef, &isUnique)
		if err != nil {
			return nil, err
		}

		idx.Unique = isUnique
		idx.Type = "btree" // Default type

		// Extract columns from index definition
		idx.Columns = sa.extractColumnsFromIndexDef(indexDef)

		indexes = append(indexes, idx)
	}

	return indexes, nil
}

// getForeignKeys retrieves foreign key information for a table
func (sa *SchemaAnalyzer) getForeignKeys(tableName string) ([]ForeignKeyInfo, error) {
	rows, err := sa.db.Raw(`
		SELECT 
			kcu.column_name,
			ccu.table_name AS foreign_table_name,
			ccu.column_name AS foreign_column_name,
			rc.delete_rule,
			rc.update_rule
		FROM information_schema.table_constraints AS tc
		JOIN information_schema.key_column_usage AS kcu
			ON tc.constraint_name = kcu.constraint_name
		JOIN information_schema.constraint_column_usage AS ccu
			ON ccu.constraint_name = tc.constraint_name
		JOIN information_schema.referential_constraints AS rc
			ON tc.constraint_name = rc.constraint_name
		WHERE tc.constraint_type = 'FOREIGN KEY' AND tc.table_name = ?
	`, tableName).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var foreignKeys []ForeignKeyInfo
	for rows.Next() {
		var fk ForeignKeyInfo
		err := rows.Scan(
			&fk.Column,
			&fk.RefTable,
			&fk.RefColumn,
			&fk.OnDelete,
			&fk.OnUpdate,
		)
		if err != nil {
			return nil, err
		}

		fk.References = fk.RefTable + "." + fk.RefColumn
		foreignKeys = append(foreignKeys, fk)
	}

	return foreignKeys, nil
}

// getConstraints retrieves constraint information for a table
func (sa *SchemaAnalyzer) getConstraints(tableName string) ([]ConstraintInfo, error) {
	rows, err := sa.db.Raw(`
		SELECT 
			tc.constraint_name,
			tc.constraint_type,
			string_agg(kcu.column_name, ', ' ORDER BY kcu.ordinal_position) as columns
		FROM information_schema.table_constraints tc
		LEFT JOIN information_schema.key_column_usage kcu
			ON tc.constraint_name = kcu.constraint_name
		WHERE tc.table_name = ?
		GROUP BY tc.constraint_name, tc.constraint_type
	`, tableName).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var constraints []ConstraintInfo
	for rows.Next() {
		var constraint ConstraintInfo
		var columns string

		err := rows.Scan(&constraint.Name, &constraint.Type, &columns)
		if err != nil {
			return nil, err
		}

		if columns != "" {
			constraint.Columns = strings.Split(columns, ", ")
		}

		constraints = append(constraints, constraint)
	}

	return constraints, nil
}

// mapToGoType maps database types to Go types
func (sa *SchemaAnalyzer) mapToGoType(dbType string, nullable bool) string {
	goType := sa.getGoType(dbType)
	if nullable && goType != "string" && goType != "[]byte" {
		return "*" + goType
	}
	return goType
}

// mapToTSType maps database types to TypeScript types
func (sa *SchemaAnalyzer) mapToTSType(dbType string, nullable bool) string {
	tsType := sa.getTSType(dbType)
	if nullable {
		return tsType + " | null"
	}
	return tsType
}

// getGoType returns the Go type for a database type
func (sa *SchemaAnalyzer) getGoType(dbType string) string {
	switch strings.ToLower(dbType) {
	case "bigint", "bigserial":
		return "int64"
	case "integer", "int", "int4", "serial":
		return "int"
	case "smallint", "int2", "smallserial":
		return "int16"
	case "boolean", "bool":
		return "bool"
	case "real", "float4":
		return "float32"
	case "double precision", "float8":
		return "float64"
	case "numeric", "decimal":
		return "float64"
	case "character varying", "varchar", "text", "char", "character":
		return "string"
	case "bytea":
		return "[]byte"
	case "date", "time", "timestamp", "timestamptz":
		return "time.Time"
	case "uuid":
		return "string"
	case "json", "jsonb":
		return "string"
	default:
		return "string"
	}
}

// getTSType returns the TypeScript type for a database type
func (sa *SchemaAnalyzer) getTSType(dbType string) string {
	switch strings.ToLower(dbType) {
	case "bigint", "bigserial", "integer", "int", "int4", "serial", "smallint", "int2", "smallserial":
		return "number"
	case "boolean", "bool":
		return "boolean"
	case "real", "float4", "double precision", "float8", "numeric", "decimal":
		return "number"
	case "character varying", "varchar", "text", "char", "character", "uuid":
		return "string"
	case "date", "time", "timestamp", "timestamptz":
		return "string"
	case "json", "jsonb":
		return "any"
	default:
		return "string"
	}
}

// extractColumnsFromIndexDef extracts column names from index definition
func (sa *SchemaAnalyzer) extractColumnsFromIndexDef(indexDef string) []string {
	// Simple extraction - in production, you'd want more sophisticated parsing
	// This is a basic implementation
	parts := strings.Split(indexDef, "(")
	if len(parts) < 2 {
		return []string{}
	}

	columnsPart := strings.TrimSuffix(parts[1], ")")
	columns := strings.Split(columnsPart, ",")

	var result []string
	for _, col := range columns {
		col = strings.TrimSpace(col)
		if col != "" {
			result = append(result, col)
		}
	}

	return result
}

// isSystemTable checks if a table is a system table
func (sa *SchemaAnalyzer) isSystemTable(tableName string) bool {
	systemTables := []string{
		"migrations",
		"goose_db_version",
		"audit_logs",
		"schema_migrations",
	}

	for _, sysTable := range systemTables {
		if strings.EqualFold(tableName, sysTable) {
			return true
		}
	}

	return false
}

// GetTableByName retrieves information for a specific table
func (sa *SchemaAnalyzer) GetTableByName(tableName string) (*TableInfo, error) {
	return sa.analyzeTable(tableName, "public", "")
}

// GetTablesByPattern retrieves tables matching a pattern
func (sa *SchemaAnalyzer) GetTablesByPattern(pattern string) ([]*TableInfo, error) {
	allTables, err := sa.DiscoverTables()
	if err != nil {
		return nil, err
	}

	var filtered []*TableInfo
	for _, table := range allTables {
		if strings.Contains(strings.ToLower(table.Name), strings.ToLower(pattern)) {
			filtered = append(filtered, table)
		}
	}

	return filtered, nil
}
