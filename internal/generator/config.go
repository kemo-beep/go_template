package generator

import (
	"fmt"
	"time"
)

// GeneratorConfig holds configuration for the API generator
type GeneratorConfig struct {
	Enabled             bool                    `yaml:"enabled"`
	AutoScan            bool                    `yaml:"auto_scan"`
	OutputDir           string                  `yaml:"output_dir"`
	PackageName         string                  `yaml:"package_name"`
	AutoRegistration    *AutoRegistrationConfig `yaml:"auto_registration"`
	GenerateTypeScript  bool                    `yaml:"generate_typescript"`
	TypeScriptOutput    string                  `yaml:"typescript_output_dir"`
	TypeScriptAPIClient bool                    `yaml:"typescript_api_client"`
	Tables              map[string]*TableConfig `yaml:"tables"`
	Global              *GlobalConfig           `yaml:"global"`
}

// AutoRegistrationConfig holds configuration for auto-registration
type AutoRegistrationConfig struct {
	Enabled       bool          `yaml:"enabled"`
	WatchInterval time.Duration `yaml:"watch_interval"`
	AutoRestart   bool          `yaml:"auto_restart"`
	HotReload     bool          `yaml:"hot_reload"`
}

// TableConfig holds configuration for a specific table
type TableConfig struct {
	Enabled       bool                   `yaml:"enabled"`
	Endpoints     []string               `yaml:"endpoints"`
	Relationships []string               `yaml:"relationships"`
	Security      *SecurityConfig        `yaml:"security"`
	Validation    *ValidationConfig      `yaml:"validation"`
	Caching       *CacheConfig           `yaml:"caching"`
	Pagination    *PaginationConfig      `yaml:"pagination"`
	Filtering     *FilteringConfig       `yaml:"filtering"`
	Sorting       *SortingConfig         `yaml:"sorting"`
	Custom        map[string]interface{} `yaml:"custom"`
}

// GlobalConfig holds global generator configuration
type GlobalConfig struct {
	Security      *SecurityConfig      `yaml:"security"`
	Validation    *ValidationConfig    `yaml:"validation"`
	Caching       *CacheConfig         `yaml:"caching"`
	Pagination    *PaginationConfig    `yaml:"pagination"`
	Filtering     *FilteringConfig     `yaml:"filtering"`
	Sorting       *SortingConfig       `yaml:"sorting"`
	Documentation *DocumentationConfig `yaml:"documentation"`
}

// SecurityConfig holds security configuration
type SecurityConfig struct {
	RBAC       *RBACConfig      `yaml:"rbac"`
	RateLimit  *RateLimitConfig `yaml:"rate_limit"`
	AuditLog   bool             `yaml:"audit_log"`
	SoftDelete bool             `yaml:"soft_delete"`
	Timestamps bool             `yaml:"timestamps"`
	CSRF       bool             `yaml:"csrf"`
}

// RBACConfig holds RBAC configuration
type RBACConfig struct {
	Resource    string              `yaml:"resource"`
	Permissions map[string][]string `yaml:"permissions"`
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Requests int           `yaml:"requests"`
	Window   time.Duration `yaml:"window"`
	Burst    int           `yaml:"burst"`
}

// ValidationConfig holds validation configuration
type ValidationConfig struct {
	Strict      bool                `yaml:"strict"`
	CustomRules map[string]string   `yaml:"custom_rules"`
	Required    []string            `yaml:"required"`
	Optional    []string            `yaml:"optional"`
	MinLength   map[string]int      `yaml:"min_length"`
	MaxLength   map[string]int      `yaml:"max_length"`
	MinValue    map[string]float64  `yaml:"min_value"`
	MaxValue    map[string]float64  `yaml:"max_value"`
	Email       []string            `yaml:"email"`
	URL         []string            `yaml:"url"`
	UUID        []string            `yaml:"uuid"`
	Enum        map[string][]string `yaml:"enum"`
}

// CacheConfig holds caching configuration
type CacheConfig struct {
	TTL          time.Duration `yaml:"ttl"`
	KeyPattern   string        `yaml:"key_pattern"`
	InvalidateOn []string      `yaml:"invalidate_on"`
	SkipCache    []string      `yaml:"skip_cache"`
	Strategy     string        `yaml:"strategy"` // "memory", "redis"
}

// PaginationConfig holds pagination configuration
type PaginationConfig struct {
	DefaultLimit int    `yaml:"default_limit"`
	MaxLimit     int    `yaml:"max_limit"`
	PageParam    string `yaml:"page_param"`
	LimitParam   string `yaml:"limit_param"`
	SortParam    string `yaml:"sort_param"`
	OrderParam   string `yaml:"order_param"`
	EnableCursor bool   `yaml:"enable_cursor"`
	CursorParam  string `yaml:"cursor_param"`
}

// FilteringConfig holds filtering configuration
type FilteringConfig struct {
	AllowedFields []string          `yaml:"allowed_fields"`
	Operators     map[string]string `yaml:"operators"`
	DateRanges    []string          `yaml:"date_ranges"`
	TextSearch    []string          `yaml:"text_search"`
	CustomFilters map[string]string `yaml:"custom_filters"`
}

// SortingConfig holds sorting configuration
type SortingConfig struct {
	AllowedFields []string `yaml:"allowed_fields"`
	DefaultSort   string   `yaml:"default_sort"`
	MultiSort     bool     `yaml:"multi_sort"`
	CaseSensitive bool     `yaml:"case_sensitive"`
}

// DocumentationConfig holds documentation configuration
type DocumentationConfig struct {
	Enabled     bool         `yaml:"enabled"`
	Title       string       `yaml:"title"`
	Version     string       `yaml:"version"`
	Description string       `yaml:"description"`
	BaseURL     string       `yaml:"base_url"`
	Contact     *ContactInfo `yaml:"contact"`
}

// ContactInfo holds contact information for documentation
type ContactInfo struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
	URL   string `yaml:"url"`
}

// DefaultGeneratorConfig returns a default generator configuration
func DefaultGeneratorConfig() *GeneratorConfig {
	return &GeneratorConfig{
		Enabled:     true,
		AutoScan:    true,
		OutputDir:   "./generated",
		PackageName: "generated",
		Tables:      make(map[string]*TableConfig),
		Global: &GlobalConfig{
			Security: &SecurityConfig{
				AuditLog:   true,
				SoftDelete: true,
				Timestamps: true,
				CSRF:       false,
				RateLimit: &RateLimitConfig{
					Requests: 100,
					Window:   time.Minute,
					Burst:    10,
				},
			},
			Validation: &ValidationConfig{
				Strict: true,
			},
			Caching: &CacheConfig{
				TTL:        5 * time.Minute,
				KeyPattern: "{table}:{operation}:{params}",
				Strategy:   "redis",
			},
			Pagination: &PaginationConfig{
				DefaultLimit: 20,
				MaxLimit:     100,
				PageParam:    "page",
				LimitParam:   "limit",
				SortParam:    "sort",
				OrderParam:   "order",
				EnableCursor: false,
			},
			Filtering: &FilteringConfig{
				Operators: map[string]string{
					"eq":    "=",
					"ne":    "!=",
					"gt":    ">",
					"gte":   ">=",
					"lt":    "<",
					"lte":   "<=",
					"like":  "LIKE",
					"ilike": "ILIKE",
					"in":    "IN",
					"nin":   "NOT IN",
					"null":  "IS NULL",
					"nnull": "IS NOT NULL",
				},
			},
			Sorting: &SortingConfig{
				MultiSort:     true,
				CaseSensitive: false,
			},
			Documentation: &DocumentationConfig{
				Enabled:     true,
				Title:       "Auto Generated API",
				Version:     "1.0.0",
				Description: "Automatically generated API endpoints",
				BaseURL:     "http://localhost:8080/api/v1",
			},
		},
	}
}

// GetTableConfig returns configuration for a specific table
func (gc *GeneratorConfig) GetTableConfig(tableName string) *TableConfig {
	if config, exists := gc.Tables[tableName]; exists {
		return config
	}

	// Return default configuration for the table
	return &TableConfig{
		Enabled:       true,
		Endpoints:     []string{"list", "create", "get", "update", "delete"},
		Relationships: []string{},
		Security:      gc.Global.Security,
		Validation:    gc.Global.Validation,
		Caching:       gc.Global.Caching,
		Pagination:    gc.Global.Pagination,
		Filtering:     gc.Global.Filtering,
		Sorting:       gc.Global.Sorting,
		Custom:        make(map[string]interface{}),
	}
}

// ShouldGenerateTable checks if a table should be generated
func (gc *GeneratorConfig) ShouldGenerateTable(tableName string) bool {
	if !gc.Enabled {
		return false
	}

	config := gc.GetTableConfig(tableName)
	return config.Enabled
}

// GetEndpointsForTable returns the endpoints to generate for a table
func (gc *GeneratorConfig) GetEndpointsForTable(tableName string) []string {
	config := gc.GetTableConfig(tableName)
	return config.Endpoints
}

// GetRelationshipsForTable returns the relationships to generate for a table
func (gc *GeneratorConfig) GetRelationshipsForTable(tableName string) []string {
	config := gc.GetTableConfig(tableName)
	return config.Relationships
}

// MergeTableConfig merges table-specific config with global config
func (gc *GeneratorConfig) MergeTableConfig(tableName string, tableConfig *TableConfig) *TableConfig {
	merged := &TableConfig{
		Enabled:       tableConfig.Enabled,
		Endpoints:     tableConfig.Endpoints,
		Relationships: tableConfig.Relationships,
		Custom:        tableConfig.Custom,
	}

	// Merge security config
	if tableConfig.Security != nil {
		merged.Security = tableConfig.Security
	} else {
		merged.Security = gc.Global.Security
	}

	// Merge validation config
	if tableConfig.Validation != nil {
		merged.Validation = tableConfig.Validation
	} else {
		merged.Security = gc.Global.Security
	}

	// Merge caching config
	if tableConfig.Caching != nil {
		merged.Caching = tableConfig.Caching
	} else {
		merged.Caching = gc.Global.Caching
	}

	// Merge pagination config
	if tableConfig.Pagination != nil {
		merged.Pagination = tableConfig.Pagination
	} else {
		merged.Pagination = gc.Global.Pagination
	}

	// Merge filtering config
	if tableConfig.Filtering != nil {
		merged.Filtering = tableConfig.Filtering
	} else {
		merged.Filtering = gc.Global.Filtering
	}

	// Merge sorting config
	if tableConfig.Sorting != nil {
		merged.Sorting = tableConfig.Sorting
	} else {
		merged.Sorting = gc.Global.Sorting
	}

	return merged
}

// ValidateConfig validates the generator configuration
func (gc *GeneratorConfig) ValidateConfig() error {
	if gc.OutputDir == "" {
		return fmt.Errorf("output directory is required")
	}

	if gc.PackageName == "" {
		return fmt.Errorf("package name is required")
	}

	if gc.Global.Pagination.MaxLimit <= 0 {
		return fmt.Errorf("max limit must be greater than 0")
	}

	if gc.Global.Pagination.DefaultLimit > gc.Global.Pagination.MaxLimit {
		return fmt.Errorf("default limit cannot be greater than max limit")
	}

	return nil
}
