package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go-mobile-backend-template/internal/utils"
	"go-mobile-backend-template/pkg/cache"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Requests int                       // Number of requests allowed
	Window   time.Duration             // Time window
	KeyFunc  func(*gin.Context) string // Function to generate cache key
}

// RateLimiter handles rate limiting logic
type RateLimiter struct {
	redis  *cache.RedisClient
	logger *zap.Logger
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(redis *cache.RedisClient, logger *zap.Logger) *RateLimiter {
	return &RateLimiter{
		redis:  redis,
		logger: logger,
	}
}

// RateLimit middleware that limits requests per IP
func (rl *RateLimiter) RateLimit(config RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := config.KeyFunc(c)
		if key == "" {
			c.Next()
			return
		}

		// Create rate limit key
		rateLimitKey := fmt.Sprintf("rate_limit:%s", key)

		// Get current count
		current, err := rl.redis.Get(context.Background(), rateLimitKey)
		if err != nil {
			// If key doesn't exist, start counting
			current = "0"
		}

		count, err := strconv.Atoi(current)
		if err != nil {
			count = 0
		}

		// Check if limit exceeded
		if count >= config.Requests {
			rl.logger.Warn("Rate limit exceeded",
				zap.String("key", key),
				zap.Int("count", count),
				zap.Int("limit", config.Requests),
			)

			c.Header("X-RateLimit-Limit", strconv.Itoa(config.Requests))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(config.Window).Unix(), 10))

			utils.ErrorResponse(c, http.StatusTooManyRequests, "Rate limit exceeded")
			c.Abort()
			return
		}

		// Increment counter
		if count == 0 {
			// First request in window
			err = rl.redis.Set(context.Background(), rateLimitKey, "1", config.Window)
		} else {
			// Increment existing counter
			_, err = rl.redis.GetClient().Incr(context.Background(), rateLimitKey).Result()
		}

		if err != nil {
			rl.logger.Error("Failed to update rate limit counter", zap.Error(err))
			// Don't block request on Redis errors
		}

		// Set response headers
		remaining := config.Requests - count - 1
		if remaining < 0 {
			remaining = 0
		}

		c.Header("X-RateLimit-Limit", strconv.Itoa(config.Requests))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(config.Window).Unix(), 10))

		c.Next()
	}
}

// IPKeyFunc generates rate limit key based on client IP
func IPKeyFunc(c *gin.Context) string {
	// Get real IP (considering proxies)
	ip := c.ClientIP()
	if ip == "" {
		ip = c.RemoteIP()
	}
	return ip
}

// UserKeyFunc generates rate limit key based on user ID
func UserKeyFunc(c *gin.Context) string {
	userID, exists := c.Get("user_id")
	if !exists {
		return ""
	}
	return fmt.Sprintf("user:%v", userID)
}

// EndpointKeyFunc generates rate limit key based on endpoint and IP
func EndpointKeyFunc(c *gin.Context) string {
	ip := c.ClientIP()
	if ip == "" {
		ip = c.RemoteIP()
	}
	return fmt.Sprintf("endpoint:%s:%s", c.Request.Method, c.Request.URL.Path)
}

// Predefined rate limit configurations
var (
	// Strict rate limiting for auth endpoints
	AuthRateLimit = RateLimitConfig{
		Requests: 5,                // 5 requests
		Window:   15 * time.Minute, // per 15 minutes
		KeyFunc:  IPKeyFunc,
	}

	// Moderate rate limiting for API endpoints
	APIRateLimit = RateLimitConfig{
		Requests: 100,             // 100 requests
		Window:   1 * time.Minute, // per minute
		KeyFunc:  UserKeyFunc,
	}

	// Lenient rate limiting for public endpoints
	PublicRateLimit = RateLimitConfig{
		Requests: 200,             // 200 requests
		Window:   1 * time.Minute, // per minute
		KeyFunc:  IPKeyFunc,
	}

	// Very strict rate limiting for admin endpoints
	AdminRateLimit = RateLimitConfig{
		Requests: 50,              // 50 requests
		Window:   1 * time.Minute, // per minute
		KeyFunc:  UserKeyFunc,
	}

	// Strict rate limiting for database operations
	DatabaseRateLimit = RateLimitConfig{
		Requests: 20,              // 20 requests
		Window:   1 * time.Minute, // per minute
		KeyFunc:  UserKeyFunc,
	}
)
