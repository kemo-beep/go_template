package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Environment   string        `mapstructure:"environment"`
	Server        Server        `mapstructure:"server"`
	Database      Database      `mapstructure:"database"`
	JWT           JWT           `mapstructure:"jwt"`
	R2            R2            `mapstructure:"r2"`
	GoogleScripts GoogleScripts `mapstructure:"google_scripts"`
	Logging       Logging       `mapstructure:"logging"`
}

// Server configuration
type Server struct {
	Port         string `mapstructure:"port"`
	Host         string `mapstructure:"host"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

// Database configuration
type Database struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

// JWT configuration
type JWT struct {
	Secret                string        `mapstructure:"secret"`
	AccessTokenExpire     time.Duration `mapstructure:"access_token_expire"`
	RefreshTokenExpire    time.Duration `mapstructure:"refresh_token_expire"`
	AccessTokenExpireInt  int           `mapstructure:"access_token_expire_int"`
	RefreshTokenExpireInt int           `mapstructure:"refresh_token_expire_int"`
}

// R2 (Cloudflare) configuration
type R2 struct {
	AccountID string `mapstructure:"account_id"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket"`
	Endpoint  string `mapstructure:"endpoint"`
}

// Google Scripts configuration
type GoogleScripts struct {
	URL           string `mapstructure:"url"`
	AccessToken   string `mapstructure:"access_token"`
	ProjectID     string `mapstructure:"project_id"`
	MigrationsDir string `mapstructure:"migrations_dir"`
}

// Logging configuration
type Logging struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// DSN returns the database connection string
func (d *Database) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode)
}

// Load loads configuration from file and environment variables
func Load() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	// Set default values
	setDefaults()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// Enable reading from environment variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("fatal error unmarshaling config: %w", err))
	}

	return &config
}

// setDefaults sets default configuration values
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.read_timeout", 30)
	viper.SetDefault("server.write_timeout", 30)

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.user", "app")
	viper.SetDefault("database.password", "secret")
	viper.SetDefault("database.name", "myapp")
	viper.SetDefault("database.ssl_mode", "disable")

	// JWT defaults
	viper.SetDefault("jwt.access_token_expire_int", 15)     // minutes
	viper.SetDefault("jwt.refresh_token_expire_int", 10080) // minutes (7 days)

	// Google Scripts defaults
	viper.SetDefault("google_scripts.url", "")
	viper.SetDefault("google_scripts.access_token", "")
	viper.SetDefault("google_scripts.project_id", "")
	viper.SetDefault("google_scripts.migrations_dir", "./migrations")

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
}
