package middleware

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go-mobile-backend-template/pkg/cache"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// PerformanceMonitor handles performance monitoring
type PerformanceMonitor struct {
	redis  *cache.RedisClient
	logger *zap.Logger
}

// NewPerformanceMonitor creates a new performance monitor
func NewPerformanceMonitor(redis *cache.RedisClient, logger *zap.Logger) *PerformanceMonitor {
	return &PerformanceMonitor{
		redis:  redis,
		logger: logger,
	}
}

// PerformanceMetrics holds performance metrics
type PerformanceMetrics struct {
	Endpoint     string        `json:"endpoint"`
	Method       string        `json:"method"`
	Duration     time.Duration `json:"duration"`
	Status       int           `json:"status"`
	ResponseSize int64         `json:"response_size"`
	Timestamp    time.Time     `json:"timestamp"`
	UserID       *uint         `json:"user_id,omitempty"`
}

// PerformanceMiddleware creates performance monitoring middleware
func (pm *PerformanceMonitor) PerformanceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Calculate metrics
		duration := time.Since(start)
		status := c.Writer.Status()

		// Get response size
		responseSize := int64(c.Writer.Size())

		// Get user ID
		var userID *uint
		if uid, exists := c.Get("user_id"); exists {
			if id, ok := uid.(uint); ok {
				userID = &id
			}
		}

		// Create metrics
		metrics := PerformanceMetrics{
			Endpoint:     c.Request.URL.Path,
			Method:       c.Request.Method,
			Duration:     duration,
			Status:       status,
			ResponseSize: responseSize,
			Timestamp:    time.Now(),
			UserID:       userID,
		}

		// Log performance metrics
		pm.logger.Info("Performance metrics",
			zap.String("endpoint", metrics.Endpoint),
			zap.String("method", metrics.Method),
			zap.Duration("duration", metrics.Duration),
			zap.Int("status", metrics.Status),
			zap.Int64("response_size", metrics.ResponseSize),
		)

		// Store metrics in Redis for analysis
		go pm.storeMetrics(metrics)
	}
}

// storeMetrics stores performance metrics in Redis
func (pm *PerformanceMonitor) storeMetrics(metrics PerformanceMetrics) {
	ctx := context.Background()

	// Create key for this endpoint
	key := fmt.Sprintf("perf:%s:%s", metrics.Method, metrics.Endpoint)

	// Store metrics as JSON
	metricsJSON, err := json.Marshal(metrics)
	if err != nil {
		pm.logger.Error("Failed to marshal performance metrics", zap.Error(err))
		return
	}

	// Store with 1 hour expiration
	if err := pm.redis.Set(ctx, key, string(metricsJSON), time.Hour); err != nil {
		pm.logger.Error("Failed to store performance metrics", zap.Error(err))
	}
}

// CacheMiddleware provides response caching
func (pm *PerformanceMonitor) CacheMiddleware(ttl time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only cache GET requests
		if c.Request.Method != "GET" {
			c.Next()
			return
		}

		// Generate cache key
		cacheKey := pm.generateCacheKey(c)

		// Try to get from cache
		ctx := context.Background()
		cached, err := pm.redis.Get(ctx, cacheKey)
		if err == nil && cached != "" {
			// Return cached response
			var cachedResponse map[string]interface{}
			if err := json.Unmarshal([]byte(cached), &cachedResponse); err == nil {
				c.JSON(200, cachedResponse)
				c.Abort()
				return
			}
		}

		// Create custom response writer to capture response
		blw := &bodyLogWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = blw

		// Process request
		c.Next()

		// Cache successful responses
		if c.Writer.Status() == 200 {
			go pm.cacheResponse(cacheKey, blw.body.Bytes(), ttl)
		}
	}
}

// generateCacheKey generates a unique cache key for the request
func (pm *PerformanceMonitor) generateCacheKey(c *gin.Context) string {
	// Include method, path, and query parameters
	key := fmt.Sprintf("%s:%s:%s", c.Request.Method, c.Request.URL.Path, c.Request.URL.RawQuery)

	// Include user ID if authenticated
	if userID, exists := c.Get("user_id"); exists {
		key = fmt.Sprintf("%s:user:%v", key, userID)
	}

	// Hash the key to keep it reasonable length
	hash := md5.Sum([]byte(key))
	return fmt.Sprintf("cache:%x", hash)
}

// cacheResponse caches the response
func (pm *PerformanceMonitor) cacheResponse(key string, response []byte, ttl time.Duration) {
	ctx := context.Background()

	// Create cache entry
	cacheEntry := map[string]interface{}{
		"data":      string(response),
		"timestamp": time.Now().Unix(),
	}

	cacheJSON, err := json.Marshal(cacheEntry)
	if err != nil {
		pm.logger.Error("Failed to marshal cache entry", zap.Error(err))
		return
	}

	if err := pm.redis.Set(ctx, key, string(cacheJSON), ttl); err != nil {
		pm.logger.Error("Failed to cache response", zap.Error(err))
	}
}

// QueryCacheMiddleware provides database query caching
func (pm *PerformanceMonitor) QueryCacheMiddleware(ttl time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only cache GET requests to database endpoints
		if c.Request.Method != "GET" || !pm.isDatabaseEndpoint(c.Request.URL.Path) {
			c.Next()
			return
		}

		// Generate cache key for database query
		cacheKey := pm.generateQueryCacheKey(c)

		// Try to get from cache
		ctx := context.Background()
		cached, err := pm.redis.Get(ctx, cacheKey)
		if err == nil && cached != "" {
			// Return cached response
			var cachedResponse map[string]interface{}
			if err := json.Unmarshal([]byte(cached), &cachedResponse); err == nil {
				c.JSON(200, cachedResponse)
				c.Abort()
				return
			}
		}

		// Create custom response writer
		blw := &bodyLogWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = blw

		// Process request
		c.Next()

		// Cache successful database responses
		if c.Writer.Status() == 200 {
			go pm.cacheResponse(cacheKey, blw.body.Bytes(), ttl)
		}
	}
}

// generateQueryCacheKey generates cache key for database queries
func (pm *PerformanceMonitor) generateQueryCacheKey(c *gin.Context) string {
	// Include path, query params, and user context
	key := fmt.Sprintf("query:%s:%s", c.Request.URL.Path, c.Request.URL.RawQuery)

	// Include user ID for user-specific queries
	if userID, exists := c.Get("user_id"); exists {
		key = fmt.Sprintf("%s:user:%v", key, userID)
	}

	// Hash the key
	hash := md5.Sum([]byte(key))
	return fmt.Sprintf("query_cache:%x", hash)
}

// isDatabaseEndpoint checks if the endpoint is a database operation
func (pm *PerformanceMonitor) isDatabaseEndpoint(path string) bool {
	databasePaths := []string{
		"/api/v1/admin/database",
		"/api/v1/users",
		"/api/v1/files",
	}

	for _, dbPath := range databasePaths {
		if strings.HasPrefix(path, dbPath) {
			return true
		}
	}

	return false
}

// GetPerformanceStats retrieves performance statistics
func (pm *PerformanceMonitor) GetPerformanceStats(endpoint string, hours int) (map[string]interface{}, error) {
	// Get metrics for the endpoint
	// key := fmt.Sprintf("perf:*:%s", endpoint)

	// This is a simplified version - in production you'd use Redis SCAN
	// or a proper time-series database like InfluxDB

	stats := map[string]interface{}{
		"endpoint": endpoint,
		"hours":    hours,
		"metrics":  []PerformanceMetrics{},
	}

	return stats, nil
}

// ClearCache clears cache for a specific pattern
func (pm *PerformanceMonitor) ClearCache(pattern string) error {
	// In production, you'd use Redis SCAN to find and delete keys
	// For now, we'll just log the request
	pm.logger.Info("Cache clear requested", zap.String("pattern", pattern))

	return nil
}

// HealthCheck provides health check endpoint
func (pm *PerformanceMonitor) HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()

		// Check Redis connection
		_, err := pm.redis.Get(ctx, "health_check")
		if err != nil {
			c.JSON(503, gin.H{
				"status":    "unhealthy",
				"redis":     "disconnected",
				"timestamp": time.Now().Unix(),
			})
			return
		}

		c.JSON(200, gin.H{
			"status":    "healthy",
			"redis":     "connected",
			"timestamp": time.Now().Unix(),
		})
	}
}
