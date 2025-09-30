package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// SecurityHeaders middleware adds security headers
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")

		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Enable XSS protection
		c.Header("X-XSS-Protection", "1; mode=block")

		// Strict Transport Security (HTTPS only)
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// Content Security Policy
		// Relaxed CSP for Swagger UI, stricter for other endpoints
		if strings.HasPrefix(c.Request.URL.Path, "/docs/") || c.Request.URL.Path == "/swagger-ui" {
			c.Header("Content-Security-Policy", "default-src 'self' https://unpkg.com; style-src 'self' 'unsafe-inline' https://unpkg.com; script-src 'self' 'unsafe-inline' 'unsafe-eval' https://unpkg.com; img-src 'self' data: https:")
		} else {
			c.Header("Content-Security-Policy", "default-src 'self'")
		}

		// Referrer Policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions Policy
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		c.Next()
	}
}
