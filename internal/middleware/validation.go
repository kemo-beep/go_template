package middleware

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"go-mobile-backend-template/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ValidationConfig holds validation configuration
type ValidationConfig struct {
	MaxRequestSize int64    // Maximum request size in bytes
	MaxFileSize    int64    // Maximum file size in bytes
	AllowedTypes   []string // Allowed content types
}

// DefaultValidationConfig returns default validation configuration
func DefaultValidationConfig() ValidationConfig {
	return ValidationConfig{
		MaxRequestSize: 10 * 1024 * 1024, // 10MB
		MaxFileSize:    50 * 1024 * 1024, // 50MB
		AllowedTypes: []string{
			"application/json",
			"application/x-www-form-urlencoded",
			"multipart/form-data",
			"text/plain",
		},
	}
}

// RequestSizeLimit middleware limits request body size
func RequestSizeLimit(config ValidationConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check content length header
		if contentLength := c.GetHeader("Content-Length"); contentLength != "" {
			if size, err := strconv.ParseInt(contentLength, 10, 64); err == nil {
				if size > config.MaxRequestSize {
					utils.ErrorResponse(c, http.StatusRequestEntityTooLarge, "Request too large")
					c.Abort()
					return
				}
			}
		}

		// Limit request body
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, config.MaxRequestSize)

		c.Next()
	}
}

// ContentTypeValidation middleware validates content type
func ContentTypeValidation(config ValidationConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip validation for GET requests without body
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" {
			c.Next()
			return
		}

		contentType := c.GetHeader("Content-Type")
		if contentType == "" {
			utils.ErrorResponse(c, http.StatusBadRequest, "Content-Type header required")
			c.Abort()
			return
		}

		// Check if content type is allowed
		allowed := false
		for _, allowedType := range config.AllowedTypes {
			if strings.HasPrefix(contentType, allowedType) {
				allowed = true
				break
			}
		}

		if !allowed {
			utils.ErrorResponse(c, http.StatusUnsupportedMediaType, "Unsupported content type")
			c.Abort()
			return
		}

		c.Next()
	}
}

// InputSanitizer sanitizes input data
type InputSanitizer struct {
	logger *zap.Logger
}

// NewInputSanitizer creates a new input sanitizer
func NewInputSanitizer(logger *zap.Logger) *InputSanitizer {
	return &InputSanitizer{
		logger: logger,
	}
}

// SanitizeString sanitizes a string input
func (s *InputSanitizer) SanitizeString(input string) string {
	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")

	// Trim whitespace
	input = strings.TrimSpace(input)

	// Remove control characters except newlines and tabs
	var result strings.Builder
	for _, r := range input {
		if r >= 32 || r == '\n' || r == '\t' || r == '\r' {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// SanitizeEmail sanitizes email input
func (s *InputSanitizer) SanitizeEmail(email string) string {
	email = s.SanitizeString(email)
	email = strings.ToLower(email)

	// Remove any whitespace
	email = strings.ReplaceAll(email, " ", "")

	return email
}

// ValidateEmail validates email format
func (s *InputSanitizer) ValidateEmail(email string) bool {
	// Basic email validation regex
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, email)
	return matched
}

// ValidatePassword validates password strength
func (s *InputSanitizer) ValidatePassword(password string) (bool, string) {
	if len(password) < 8 {
		return false, "Password must be at least 8 characters long"
	}

	if len(password) > 128 {
		return false, "Password must be less than 128 characters"
	}

	// Check for at least one uppercase letter
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	// Check for at least one lowercase letter
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	// Check for at least one digit
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	// Check for at least one special character
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)

	if !hasUpper {
		return false, "Password must contain at least one uppercase letter"
	}
	if !hasLower {
		return false, "Password must contain at least one lowercase letter"
	}
	if !hasDigit {
		return false, "Password must contain at least one digit"
	}
	if !hasSpecial {
		return false, "Password must contain at least one special character"
	}

	return true, ""
}

// ValidateUsername validates username format
func (s *InputSanitizer) ValidateUsername(username string) (bool, string) {
	username = s.SanitizeString(username)

	if len(username) < 3 {
		return false, "Username must be at least 3 characters long"
	}

	if len(username) > 30 {
		return false, "Username must be less than 30 characters"
	}

	// Only allow alphanumeric characters, underscores, and hyphens
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, username)
	if !matched {
		return false, "Username can only contain letters, numbers, underscores, and hyphens"
	}

	return true, ""
}

// InputValidation middleware that validates and sanitizes input
func (s *InputSanitizer) InputValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		// This would be implemented based on specific endpoint requirements
		// For now, we'll just log the request for monitoring
		s.logger.Info("Input validation middleware executed",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("ip", c.ClientIP()),
		)

		c.Next()
	}
}

// FileUploadValidation validates file uploads
func FileUploadValidation(maxSize int64, allowedTypes []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// This would be implemented in file upload handlers
		// For now, we'll just set the max file size
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)

		c.Next()
	}
}

// CSRFProtection provides basic CSRF protection
func CSRFProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip CSRF for GET, HEAD, OPTIONS
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// Check for CSRF token in header
		csrfToken := c.GetHeader("X-CSRF-Token")
		if csrfToken == "" {
			utils.ErrorResponse(c, http.StatusForbidden, "CSRF token required")
			c.Abort()
			return
		}

		// In a real implementation, you would validate the CSRF token
		// against a session or database value

		c.Next()
	}
}
