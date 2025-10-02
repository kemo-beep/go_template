package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

// DefaultCORSConfig returns a secure default CORS configuration
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"http://localhost:3001",
			"http://localhost:3002",
			"https://yourdomain.com", // Replace with your production domain
		},
		AllowedMethods: []string{
			"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS",
		},
		AllowedHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-CSRF-Token",
			"X-Requested-With",
			"Cache-Control",
			"Pragma",
		},
		ExposedHeaders: []string{
			"X-RateLimit-Limit",
			"X-RateLimit-Remaining",
			"X-RateLimit-Reset",
		},
		AllowCredentials: true,
		MaxAge:           86400, // 24 hours
	}
}

// CORS middleware for cross-origin requests with security
func CORS(config CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		allowed := false
		if origin != "" {
			for _, allowedOrigin := range config.AllowedOrigins {
				if origin == allowedOrigin || allowedOrigin == "*" {
					allowed = true
					break
				}
			}
		}

		// Set CORS headers
		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		} else if len(config.AllowedOrigins) > 0 && config.AllowedOrigins[0] == "*" {
			c.Header("Access-Control-Allow-Origin", "*")
		}

		if config.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		if len(config.AllowedMethods) > 0 {
			c.Header("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
		}

		if len(config.AllowedHeaders) > 0 {
			c.Header("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))
		}

		if len(config.ExposedHeaders) > 0 {
			c.Header("Access-Control-Expose-Headers", strings.Join(config.ExposedHeaders, ", "))
		}

		if config.MaxAge > 0 {
			c.Header("Access-Control-Max-Age", string(rune(config.MaxAge)))
		}

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// SecureCORS returns a secure CORS configuration for production
func SecureCORS() gin.HandlerFunc {
	config := CORSConfig{
		AllowedOrigins: []string{
			"https://yourdomain.com", // Replace with your production domain
		},
		AllowedMethods: []string{
			"GET", "POST", "PUT", "DELETE", "PATCH",
		},
		AllowedHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-CSRF-Token",
		},
		ExposedHeaders: []string{
			"X-RateLimit-Limit",
			"X-RateLimit-Remaining",
			"X-RateLimit-Reset",
		},
		AllowCredentials: true,
		MaxAge:           3600, // 1 hour
	}

	return CORS(config)
}

// DevelopmentCORS returns a permissive CORS configuration for development
func DevelopmentCORS() gin.HandlerFunc {
	config := CORSConfig{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"http://localhost:3001",
			"http://localhost:3002",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:3001",
			"http://127.0.0.1:3002",
		},
		AllowedMethods: []string{
			"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS",
		},
		AllowedHeaders: []string{
			"*",
		},
		ExposedHeaders: []string{
			"*",
		},
		AllowCredentials: true,
		MaxAge:           86400, // 24 hours
	}

	return CORS(config)
}
