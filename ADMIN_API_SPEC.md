# Admin API Specification

This document outlines the required backend API endpoints to support the Admin Dashboard.

## 🔐 Authentication

All admin endpoints require:
- JWT authentication
- Admin role verification
- Audit logging

## 📋 Required Endpoints

### 1. User Management

#### List Users
```
GET /api/v1/admin/users
Response: [
  {
    "id": "uuid",
    "email": "user@example.com",
    "name": "John Doe",
    "is_active": true,
    "is_admin": false,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
]
```

#### Get User
```
GET /api/v1/admin/users/:id
Response: {
  "id": "uuid",
  "email": "user@example.com",
  "name": "John Doe",
  "is_active": true,
  "is_admin": false,
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

#### Create User
```
POST /api/v1/admin/users
Body: {
  "email": "user@example.com",
  "name": "John Doe",
  "password": "secure_password",
  "is_admin": false
}
Response: {
  "id": "uuid",
  "email": "user@example.com",
  "name": "John Doe"
}
```

#### Update User
```
PUT /api/v1/admin/users/:id
Body: {
  "name": "Jane Doe",
  "is_active": true,
  "is_admin": false
}
Response: {
  "id": "uuid",
  "email": "user@example.com",
  "name": "Jane Doe"
}
```

#### Delete User
```
DELETE /api/v1/admin/users/:id
Response: 204 No Content
```

#### Reset Password
```
POST /api/v1/admin/users/:id/reset-password
Body: {
  "password": "new_secure_password"
}
Response: { "success": true }
```

### 2. Database Explorer

#### List Tables
```
GET /api/v1/admin/database/tables
Response: [
  {
    "name": "users",
    "row_count": 100,
    "size": "1.5 MB"
  },
  {
    "name": "files",
    "row_count": 50,
    "size": "500 KB"
  }
]
```

#### Get Table Data
```
GET /api/v1/admin/database/tables/:table?limit=100&offset=0
Response: {
  "columns": ["id", "email", "name", "created_at"],
  "rows": [
    {
      "id": "uuid",
      "email": "user@example.com",
      "name": "John Doe",
      "created_at": "2025-01-01T00:00:00Z"
    }
  ],
  "total": 100,
  "limit": 100,
  "offset": 0
}
```

#### Execute Query
```
POST /api/v1/admin/database/query
Body: {
  "query": "SELECT * FROM users WHERE email LIKE '%@example.com' LIMIT 10"
}
Response: {
  "columns": ["id", "email", "name"],
  "rows": [...],
  "execution_time": "45ms",
  "row_count": 10
}
```

#### Get Schema
```
GET /api/v1/admin/database/schema
Response: [
  {
    "table": "users",
    "columns": [
      {
        "name": "id",
        "type": "uuid",
        "nullable": false,
        "primary_key": true
      },
      {
        "name": "email",
        "type": "varchar(255)",
        "nullable": false,
        "unique": true
      }
    ],
    "indexes": [...],
    "foreign_keys": [...]
  }
]
```

### 3. Storage (R2) Management

Uses existing `/api/v1/files` endpoints:
- `GET /api/v1/files` - List files
- `POST /api/v1/files/upload` - Upload file
- `DELETE /api/v1/files/:id` - Delete file
- `GET /api/v1/files/:id/download` - Get download URL

### 4. Logs & Monitoring

#### Get Logs
```
GET /api/v1/admin/logs?limit=100&level=error&since=2025-01-01T00:00:00Z
Response: [
  {
    "timestamp": "2025-01-01T00:00:00Z",
    "level": "error",
    "method": "POST",
    "path": "/api/v1/auth/login",
    "status": 401,
    "latency": "50ms",
    "client_ip": "192.168.1.1",
    "user_agent": "Mozilla/5.0...",
    "error": "invalid credentials"
  }
]
```

#### Get Metrics
```
GET /api/v1/admin/metrics?from=2025-01-01T00:00:00Z&to=2025-01-02T00:00:00Z&interval=1h
Response: [
  {
    "timestamp": "2025-01-01T00:00:00Z",
    "requests_per_minute": 100,
    "error_rate": 0.05,
    "avg_latency": 45.5,
    "p95_latency": 120.0,
    "p99_latency": 250.0
  }
]
```

### 5. Developer Tools

#### List Migrations
```
GET /api/v1/admin/migrations
Response: [
  {
    "version": "000001",
    "name": "create_users_table",
    "applied": true,
    "applied_at": "2025-01-01T00:00:00Z"
  },
  {
    "version": "000002",
    "name": "create_files_table",
    "applied": false,
    "applied_at": null
  }
]
```

#### Run Migration
```
POST /api/v1/admin/migrations/run
Body: {
  "direction": "up",  // or "down"
  "version": "000002" // optional, run specific version
}
Response: {
  "success": true,
  "message": "Migration applied successfully",
  "version": "000002"
}
```

#### List Feature Flags
```
GET /api/v1/admin/feature-flags
Response: [
  {
    "name": "new_ui",
    "enabled": true,
    "description": "Enable new UI features",
    "updated_at": "2025-01-01T00:00:00Z"
  }
]
```

#### Update Feature Flag
```
PUT /api/v1/admin/feature-flags/:name
Body: {
  "enabled": true
}
Response: {
  "name": "new_ui",
  "enabled": true
}
```

#### Run Background Job
```
POST /api/v1/admin/jobs/run
Body: {
  "job": "cleanup_old_files",
  "params": {
    "older_than_days": 30
  }
}
Response: {
  "job_id": "uuid",
  "status": "queued"
}
```

## 🔒 Implementation Guidelines

### 1. Admin Middleware

Create middleware to protect admin routes:

```go
func AdminOnly() gin.HandlerFunc {
    return func(c *gin.Context) {
        user := c.MustGet("user").(models.User)
        if !user.IsAdmin {
            c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

### 2. Audit Logging

Log all admin actions:

```go
func AuditLog() gin.HandlerFunc {
    return func(c *gin.Context) {
        user := c.MustGet("user").(models.User)
        // Log: user ID, action, timestamp, IP, etc.
        logger.Info("Admin action",
            zap.String("user_id", user.ID),
            zap.String("action", c.Request.Method + " " + c.Request.URL.Path),
            zap.String("ip", c.ClientIP()),
        )
        c.Next()
    }
}
```

### 3. Rate Limiting

Protect sensitive endpoints:

```go
adminGroup := router.Group("/api/v1/admin")
adminGroup.Use(middleware.Auth())
adminGroup.Use(middleware.AdminOnly())
adminGroup.Use(middleware.RateLimit(100, time.Minute))
adminGroup.Use(middleware.AuditLog())
```

### 4. SQL Query Safety

For database explorer, sanitize queries:

```go
func ExecuteQuery(c *gin.Context) {
    var req struct {
        Query string `json:"query"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // Only allow SELECT queries
    if !strings.HasPrefix(strings.ToUpper(strings.TrimSpace(req.Query)), "SELECT") {
        c.JSON(400, gin.H{"error": "only SELECT queries allowed"})
        return
    }
    
    // Execute with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    rows, err := db.QueryContext(ctx, req.Query)
    // ... handle results
}
```

### 5. Pagination

Always implement pagination:

```go
func GetTableData(c *gin.Context) {
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
    offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
    
    // Enforce max limit
    if limit > 1000 {
        limit = 1000
    }
    
    // Query with limit and offset
    // ...
}
```

## 📊 Metrics Collection

Implement metrics collection using:

1. **Prometheus**: For metrics storage
2. **Middleware**: To collect request metrics
3. **Aggregation**: Calculate RPM, error rates, latencies

Example metrics middleware:

```go
func MetricsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        duration := time.Since(start)
        
        // Record metrics
        metrics.RequestDuration.WithLabelValues(
            c.Request.Method,
            c.Request.URL.Path,
            strconv.Itoa(c.Writer.Status()),
        ).Observe(duration.Seconds())
        
        metrics.RequestTotal.WithLabelValues(
            c.Request.Method,
            c.Request.URL.Path,
            strconv.Itoa(c.Writer.Status()),
        ).Inc()
    }
}
```

## 🧪 Testing

Test all admin endpoints:

```go
func TestAdminUserList(t *testing.T) {
    // Setup
    router := setupTestRouter()
    w := httptest.NewRecorder()
    
    // Create admin user and get token
    token := createAdminToken()
    
    // Make request
    req, _ := http.NewRequest("GET", "/api/v1/admin/users", nil)
    req.Header.Set("Authorization", "Bearer " + token)
    router.ServeHTTP(w, req)
    
    // Assert
    assert.Equal(t, 200, w.Code)
}
```

## 📝 Next Steps

1. Create `internal/api/v1/admin` package
2. Implement handlers for each endpoint
3. Add admin middleware
4. Set up audit logging
5. Implement rate limiting
6. Add comprehensive tests
7. Document all endpoints with Swagger

## 🔗 Related Files

- Frontend: `/frontend/lib/api-client.ts`
- Backend: `/internal/api/v1/admin/` (to be created)
- Middleware: `/internal/middleware/admin.go` (to be created)
