package utils

import (
	"regexp"
)

// Email validation regex
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// IsValidEmail validates an email address
func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// IsValidPassword validates a password
func IsValidPassword(password string) bool {
	return len(password) >= 8
}

// IsValidName validates a name
func IsValidName(name string) bool {
	return len(name) >= 2 && len(name) <= 100
}

// ValidationErrors holds validation error messages
type ValidationErrors map[string]string

// NewValidationErrors creates a new ValidationErrors instance
func NewValidationErrors() ValidationErrors {
	return make(ValidationErrors)
}

// Add adds a validation error
func (v ValidationErrors) Add(field, message string) {
	v[field] = message
}

// HasErrors checks if there are any validation errors
func (v ValidationErrors) HasErrors() bool {
	return len(v) > 0
}
