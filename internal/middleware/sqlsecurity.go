package middleware

import (
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SQLSecurity middleware provides SQL injection protection
type SQLSecurity struct {
	logger *zap.Logger
}

// NewSQLSecurity creates a new SQL security middleware
func NewSQLSecurity(logger *zap.Logger) *SQLSecurity {
	return &SQLSecurity{
		logger: logger,
	}
}

// DangerousSQLPatterns contains regex patterns for dangerous SQL operations
var DangerousSQLPatterns = []*regexp.Regexp{
	// Comment patterns
	regexp.MustCompile(`(?i)--`),
	regexp.MustCompile(`(?i)/\*.*\*/`),

	// Union attacks
	regexp.MustCompile(`(?i)union\s+select`),
	regexp.MustCompile(`(?i)union\s+all\s+select`),

	// Information schema attacks
	regexp.MustCompile(`(?i)information_schema`),
	regexp.MustCompile(`(?i)pg_tables`),
	regexp.MustCompile(`(?i)pg_database`),

	// System functions
	regexp.MustCompile(`(?i)version\(\)`),
	regexp.MustCompile(`(?i)database\(\)`),
	regexp.MustCompile(`(?i)user\(\)`),
	regexp.MustCompile(`(?i)current_user`),
	regexp.MustCompile(`(?i)session_user`),

	// Time-based attacks
	regexp.MustCompile(`(?i)sleep\(`),
	regexp.MustCompile(`(?i)waitfor\s+delay`),
	regexp.MustCompile(`(?i)benchmark\(`),

	// Stacked queries
	regexp.MustCompile(`(?i);\s*(drop|delete|insert|update|create|alter|exec|execute)`),

	// Boolean-based attacks
	regexp.MustCompile(`(?i)and\s+1\s*=\s*1`),
	regexp.MustCompile(`(?i)or\s+1\s*=\s*1`),
	regexp.MustCompile(`(?i)and\s+true`),
	regexp.MustCompile(`(?i)or\s+true`),

	// Error-based attacks
	regexp.MustCompile(`(?i)extractvalue\(`),
	regexp.MustCompile(`(?i)updatexml\(`),
	regexp.MustCompile(`(?i)exp\(`),
	regexp.MustCompile(`(?i)floor\(`),

	// File operations
	regexp.MustCompile(`(?i)load_file\(`),
	regexp.MustCompile(`(?i)into\s+outfile`),
	regexp.MustCompile(`(?i)into\s+dumpfile`),

	// Privilege escalation
	regexp.MustCompile(`(?i)grant\s+`),
	regexp.MustCompile(`(?i)revoke\s+`),
	regexp.MustCompile(`(?i)create\s+user`),
	regexp.MustCompile(`(?i)drop\s+user`),
}

// ValidateSQLInput validates input for SQL injection patterns
func (s *SQLSecurity) ValidateSQLInput(input string) (bool, string) {
	input = strings.TrimSpace(input)

	// Check for dangerous patterns
	for _, pattern := range DangerousSQLPatterns {
		if pattern.MatchString(input) {
			s.logger.Warn("SQL injection attempt detected",
				zap.String("pattern", pattern.String()),
				zap.String("input", input),
			)
			return false, "Potentially dangerous SQL pattern detected"
		}
	}

	// Check for suspicious character sequences
	suspiciousChars := []string{
		"';", "';--", "';/*", "';#",
		"\"", "\";", "\";--", "\";/*",
		"\\x", "\\u", "\\n", "\\r", "\\t",
		"<script", "javascript:", "vbscript:",
	}

	for _, char := range suspiciousChars {
		if strings.Contains(strings.ToLower(input), char) {
			s.logger.Warn("Suspicious character sequence detected",
				zap.String("sequence", char),
				zap.String("input", input),
			)
			return false, "Suspicious character sequence detected"
		}
	}

	return true, ""
}

// SQLInjectionProtection middleware that validates request parameters
func (s *SQLSecurity) SQLInjectionProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check query parameters
		for key, values := range c.Request.URL.Query() {
			for _, value := range values {
				if valid, reason := s.ValidateSQLInput(value); !valid {
					s.logger.Error("SQL injection attempt blocked",
						zap.String("parameter", key),
						zap.String("value", value),
						zap.String("reason", reason),
						zap.String("ip", c.ClientIP()),
						zap.String("user_agent", c.Request.UserAgent()),
					)

					c.JSON(400, gin.H{
						"success": false,
						"error":   "Invalid input detected",
					})
					c.Abort()
					return
				}
			}
		}

		// Check path parameters
		for _, param := range c.Params {
			if valid, reason := s.ValidateSQLInput(param.Value); !valid {
				s.logger.Error("SQL injection attempt blocked in path parameter",
					zap.String("parameter", param.Key),
					zap.String("value", param.Value),
					zap.String("reason", reason),
					zap.String("ip", c.ClientIP()),
				)

				c.JSON(400, gin.H{
					"success": false,
					"error":   "Invalid input detected",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// SanitizeTableName sanitizes table names to prevent SQL injection
func (s *SQLSecurity) SanitizeTableName(tableName string) (string, error) {
	// Only allow alphanumeric characters and underscores
	matched, err := regexp.MatchString(`^[a-zA-Z_][a-zA-Z0-9_]*$`, tableName)
	if err != nil || !matched {
		return "", err
	}

	// Check length
	if len(tableName) > 63 { // PostgreSQL identifier limit
		return "", err
	}

	return tableName, nil
}

// SanitizeColumnName sanitizes column names to prevent SQL injection
func (s *SQLSecurity) SanitizeColumnName(columnName string) (string, error) {
	// Only allow alphanumeric characters and underscores
	matched, err := regexp.MatchString(`^[a-zA-Z_][a-zA-Z0-9_]*$`, columnName)
	if err != nil || !matched {
		return "", err
	}

	// Check length
	if len(columnName) > 63 { // PostgreSQL identifier limit
		return "", err
	}

	return columnName, nil
}

// ValidateSQLQuery validates SQL queries for safety
func (s *SQLSecurity) ValidateSQLQuery(query string) (bool, string) {
	query = strings.TrimSpace(query)

	// Only allow SELECT queries for read operations
	upperQuery := strings.ToUpper(query)
	if !strings.HasPrefix(upperQuery, "SELECT") {
		return false, "Only SELECT queries are allowed"
	}

	// Check for dangerous keywords
	dangerousKeywords := []string{
		"DROP", "DELETE", "INSERT", "UPDATE", "CREATE", "ALTER",
		"EXEC", "EXECUTE", "SP_", "XP_", "OPENROWSET", "OPENDATASOURCE",
		"BULK", "BULKINSERT", "BACKUP", "RESTORE", "SHUTDOWN",
		"RECONFIGURE", "DBCC", "KILL", "DENY", "REVOKE",
	}

	for _, keyword := range dangerousKeywords {
		if strings.Contains(upperQuery, keyword) {
			s.logger.Warn("Dangerous SQL keyword detected",
				zap.String("keyword", keyword),
				zap.String("query", query),
			)
			return false, "Dangerous SQL keyword detected: " + keyword
		}
	}

	// Check for SQL injection patterns
	if valid, reason := s.ValidateSQLInput(query); !valid {
		return false, reason
	}

	return true, ""
}
