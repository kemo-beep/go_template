package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AuditLog represents an audit log entry
type AuditLog struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserID       *uint     `json:"user_id" gorm:"index"`
	Action       string    `json:"action" gorm:"index"`
	Resource     string    `json:"resource"`
	ResourceID   *string   `json:"resource_id"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	RequestData  string    `json:"request_data" gorm:"type:text"`
	ResponseData string    `json:"response_data" gorm:"type:text"`
	Status       int       `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

// TableName returns the table name for AuditLog
func (AuditLog) TableName() string {
	return "audit_logs"
}

// AuditLogger handles audit logging
type AuditLogger struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewAuditLogger creates a new audit logger
func NewAuditLogger(db *gorm.DB, logger *zap.Logger) *AuditLogger {
	return &AuditLogger{
		db:     db,
		logger: logger,
	}
}

// AuditMiddleware creates audit logging middleware
func (al *AuditLogger) AuditMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Capture request body
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Create a custom response writer to capture response
		blw := &bodyLogWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = blw

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Get user ID from context
		var userID *uint
		if uid, exists := c.Get("user_id"); exists {
			if id, ok := uid.(uint); ok {
				userID = &id
			}
		}

		// Create audit log entry
		auditLog := AuditLog{
			UserID:       userID,
			Action:       c.Request.Method,
			Resource:     c.Request.URL.Path,
			ResourceID:   getResourceID(c),
			IPAddress:    c.ClientIP(),
			UserAgent:    c.Request.UserAgent(),
			RequestData:  string(requestBody),
			ResponseData: blw.body.String(),
			Status:       c.Writer.Status(),
			CreatedAt:    time.Now(),
		}

		// Log to database (async)
		go func() {
			if err := al.db.Create(&auditLog).Error; err != nil {
				al.logger.Error("Failed to create audit log", zap.Error(err))
			}
		}()

		// Log to application logs
		al.logger.Info("Audit log",
			zap.Uint("user_id", getUintValue(userID)),
			zap.String("action", auditLog.Action),
			zap.String("resource", auditLog.Resource),
			zap.String("ip", auditLog.IPAddress),
			zap.Int("status", auditLog.Status),
			zap.Duration("duration", duration),
		)
	}
}

// bodyLogWriter wraps gin.ResponseWriter to capture response body
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// getResourceID extracts resource ID from URL parameters
func getResourceID(c *gin.Context) *string {
	// Try to get ID from URL parameters
	if id := c.Param("id"); id != "" {
		return &id
	}
	if id := c.Param("tableName"); id != "" {
		return &id
	}
	if id := c.Param("userId"); id != "" {
		return &id
	}
	return nil
}

// getUintValue safely gets uint value from pointer
func getUintValue(ptr *uint) uint {
	if ptr != nil {
		return *ptr
	}
	return 0
}

// LogSecurityEvent logs security-related events
func (al *AuditLogger) LogSecurityEvent(eventType, description string, c *gin.Context, additionalData map[string]interface{}) {
	var userID *uint
	if uid, exists := c.Get("user_id"); exists {
		if id, ok := uid.(uint); ok {
			userID = &id
		}
	}

	// Create security event log
	securityLog := AuditLog{
		UserID:      userID,
		Action:      "SECURITY_EVENT",
		Resource:    eventType,
		ResourceID:  &description,
		IPAddress:   c.ClientIP(),
		UserAgent:   c.Request.UserAgent(),
		RequestData: marshalAdditionalData(additionalData),
		Status:      http.StatusForbidden,
		CreatedAt:   time.Now(),
	}

	// Log to database
	if err := al.db.Create(&securityLog).Error; err != nil {
		al.logger.Error("Failed to create security audit log", zap.Error(err))
	}

	// Log to application logs with WARN level
	al.logger.Warn("Security event",
		zap.String("event_type", eventType),
		zap.String("description", description),
		zap.Uint("user_id", getUintValue(userID)),
		zap.String("ip", c.ClientIP()),
		zap.Any("additional_data", additionalData),
	)
}

// LogAdminAction logs admin-specific actions
func (al *AuditLogger) LogAdminAction(action, resource string, c *gin.Context, additionalData map[string]interface{}) {
	userID, _ := c.Get("user_id")
	uid, _ := userID.(uint)

	adminLog := AuditLog{
		UserID:      &uid,
		Action:      "ADMIN_" + action,
		Resource:    resource,
		IPAddress:   c.ClientIP(),
		UserAgent:   c.Request.UserAgent(),
		RequestData: marshalAdditionalData(additionalData),
		Status:      c.Writer.Status(),
		CreatedAt:   time.Now(),
	}

	// Log to database
	if err := al.db.Create(&adminLog).Error; err != nil {
		al.logger.Error("Failed to create admin audit log", zap.Error(err))
	}

	// Log to application logs
	al.logger.Info("Admin action",
		zap.Uint("user_id", uid),
		zap.String("action", action),
		zap.String("resource", resource),
		zap.String("ip", c.ClientIP()),
		zap.Any("additional_data", additionalData),
	)
}

// marshalAdditionalData safely marshals additional data to JSON
func marshalAdditionalData(data map[string]interface{}) string {
	if data == nil {
		return ""
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}

	return string(jsonData)
}

// GetAuditLogs retrieves audit logs with filtering
func (al *AuditLogger) GetAuditLogs(userID *uint, action, resource string, limit, offset int) ([]AuditLog, int64, error) {
	var logs []AuditLog
	var total int64

	query := al.db.Model(&AuditLog{})

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if resource != "" {
		query = query.Where("resource LIKE ?", "%"+resource+"%")
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get logs with pagination
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// CleanupOldLogs removes audit logs older than specified days
func (al *AuditLogger) CleanupOldLogs(days int) error {
	cutoffDate := time.Now().AddDate(0, 0, -days)
	return al.db.Where("created_at < ?", cutoffDate).Delete(&AuditLog{}).Error
}
