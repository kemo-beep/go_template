# File & Folder Structure + Tech Stack

## 📁 Complete Project Structure

```
go-mobile-backend-template/
├── cmd/
│   └── server/
│       └── main.go                    # Application entry point
├── config/
│   ├── config.go                      # Configuration structs
│   ├── config.yaml                    # Default configuration
│   └── environments/
│       ├── development.yaml
│       ├── staging.yaml
│       └── production.yaml
├── deploy/
│   ├── docker-compose.yml             # Production Docker Compose
│   ├── docker-compose.dev.yml         # Development Docker Compose
│   ├── Dockerfile                     # Multi-stage Docker build
│   ├── nginx.conf                     # Nginx reverse proxy config
│   └── init.sql                       # Database initialization
├── internal/
│   ├── api/
│   │   └── v1/
│   │       ├── auth/
│   │       │   ├── handler.go         # Auth HTTP handlers
│   │       │   ├── service.go         # Auth business logic
│   │       │   ├── request.go         # Request/response structs
│   │       │   └── middleware.go      # Auth middleware
│   │       ├── users/
│   │       │   ├── handler.go         # User HTTP handlers
│   │       │   ├── service.go         # User business logic
│   │       │   └── request.go         # Request/response structs
│   │       ├── files/
│   │       │   ├── handler.go         # File HTTP handlers
│   │       │   ├── service.go         # File business logic
│   │       │   └── request.go         # Request/response structs
│   │       └── routes.go              # API route registration
│   ├── db/
│   │   ├── migrations/
│   │   │   ├── 000001_create_users_table.up.sql
│   │   │   ├── 000001_create_users_table.down.sql
│   │   │   ├── 000002_create_files_table.up.sql
│   │   │   ├── 000002_create_files_table.down.sql
│   │   │   ├── 000003_create_refresh_tokens_table.up.sql
│   │   │   └── 000003_create_refresh_tokens_table.down.sql
│   │   └── repository/
│   │       ├── user.go                # User repository
│   │       ├── file.go                # File repository
│   │       ├── refresh_token.go       # Refresh token repository
│   │       └── interfaces.go          # Repository interfaces
│   ├── middleware/
│   │   ├── auth.go                    # JWT authentication
│   │   ├── cors.go                    # CORS configuration
│   │   ├── logging.go                 # Request logging
│   │   ├── recovery.go                # Panic recovery
│   │   ├── rate_limit.go              # Rate limiting
│   │   ├── security.go                # Security headers
│   │   └── validation.go              # Request validation
│   ├── services/
│   │   ├── auth/
│   │   │   ├── jwt.go                 # JWT token management
│   │   │   ├── password.go            # Password hashing
│   │   │   └── service.go             # Auth service
│   │   ├── storage/
│   │   │   ├── r2.go                  # Cloudflare R2 client
│   │   │   ├── file.go                # File operations
│   │   │   └── service.go             # Storage service
│   │   └── user/
│   │       └── service.go             # User service
│   └── utils/
│       ├── crypto.go                  # Cryptographic utilities
│       ├── validation.go              # Validation helpers
│       ├── response.go                # HTTP response helpers
│       └── errors.go                  # Custom error types
├── pkg/
│   ├── config/
│   │   ├── config.go                  # Configuration loader
│   │   └── viper.go                   # Viper integration
│   ├── logger/
│   │   ├── logger.go                  # Logger interface
│   │   └── zap.go                     # Zap implementation
│   └── database/
│       ├── postgres.go                # PostgreSQL connection
│       └── migrations.go              # Migration runner
├── tests/
│   ├── integration/
│   │   ├── auth_test.go               # Auth integration tests
│   │   ├── users_test.go              # User integration tests
│   │   └── files_test.go              # File integration tests
│   ├── unit/
│   │   ├── services/
│   │   ├── middleware/
│   │   └── utils/
│   └── fixtures/
│       ├── users.json                 # Test user data
│       └── files.json                 # Test file data
├── scripts/
│   ├── setup.sh                       # Initial setup script
│   ├── migrate.sh                     # Migration runner
│   ├── test.sh                        # Test runner
│   └── deploy.sh                      # Deployment script
├── docs/
│   ├── api/
│   │   └── swagger.json               # Generated API docs
│   ├── deployment.md                  # Deployment guide
│   └── development.md                 # Development guide
├── .env.example                       # Environment variables template
├── .env                               # Local environment variables
├── .gitignore                         # Git ignore rules
├── .golangci.yml                      # Linter configuration
├── .air.toml                          # Air hot reload config
├── Makefile                           # Build and development commands
├── go.mod                             # Go module definition
├── go.sum                             # Go module checksums
├── docker-compose.yml                 # Development Docker Compose
├── Dockerfile                         # Application Dockerfile
└── README.md                          # Project documentation
```

## 🛠️ Tech Stack

### **Core Framework & Language**
- **Go 1.23+** - Programming language
- **Gin** - HTTP web framework
- **GORM** - Object-Relational Mapping
- **PostgreSQL 15** - Primary database

### **Authentication & Security**
- **JWT (golang-jwt/jwt/v5)** - JSON Web Tokens
- **bcrypt/argon2id** - Password hashing
- **CORS** - Cross-Origin Resource Sharing
- **Rate Limiting** - Request throttling
- **Helmet-style headers** - Security headers

### **Database & Migrations**
- **PostgreSQL** - Primary database
- **Goose** - Database migrations
- **GORM** - ORM with migrations support
- **Connection pooling** - Database connection management

### **Configuration & Environment**
- **Viper** - Configuration management
- **Environment variables** - 12-factor app compliance
- **YAML configs** - Environment-specific settings
- **Validation** - Configuration validation

### **Logging & Monitoring**
- **Zap** - Structured logging
- **JSON logs** - Machine-readable format
- **Request ID** - Request tracing
- **Health checks** - Application monitoring
- **Prometheus** - Metrics collection (optional)

### **File Storage**
- **Cloudflare R2** - Object storage
- **AWS SDK v2** - S3-compatible API
- **Signed URLs** - Secure file access
- **File validation** - Type and size checks

### **API & Documentation**
- **REST API** - RESTful endpoints
- **Swagger/OpenAPI** - API documentation
- **JSON** - Request/response format
- **Validation** - Request validation
- **Versioning** - API versioning (/api/v1/)

### **Testing**
- **Testify** - Testing framework
- **Testcontainers** - Integration testing
- **Mocking** - Service mocking
- **Coverage** - Test coverage reporting

### **Development Tools**
- **Air** - Hot reloading
- **golangci-lint** - Linting and static analysis
- **Makefile** - Build automation
- **Pre-commit hooks** - Code quality
- **Go modules** - Dependency management

### **Containerization & Deployment**
- **Docker** - Containerization
- **Docker Compose** - Multi-container orchestration
- **Multi-stage builds** - Optimized images
- **Alpine Linux** - Minimal base image
- **Nginx** - Reverse proxy (optional)

### **CI/CD & Quality**
- **GitHub Actions** - Continuous Integration
- **Code scanning** - Security scanning
- **Dependency scanning** - Vulnerability detection
- **Automated testing** - CI/CD pipeline

## 📦 Key Dependencies

### **Core Dependencies**
```go
// Web Framework
github.com/gin-gonic/gin v1.9.1

// Database
gorm.io/gorm v1.25.5
gorm.io/driver/postgres v1.5.4
github.com/pressly/goose/v3 v3.17.0

// Configuration
github.com/spf13/viper v1.17.0

// Logging
go.uber.org/zap v1.26.0

// Authentication
github.com/golang-jwt/jwt/v5 v5.2.0
golang.org/x/crypto v0.17.0

// Validation
github.com/go-playground/validator/v10 v10.16.0

// File Storage
github.com/aws/aws-sdk-go-v2 v1.24.0
github.com/aws/aws-sdk-go-v2/service/s3 v1.47.0

// Rate Limiting
github.com/ulule/limiter/v3 v3.11.2

// API Documentation
github.com/swaggo/gin-swagger v1.6.0
github.com/swaggo/swag v1.16.2

// Testing
github.com/stretchr/testify v1.8.4
github.com/testcontainers/testcontainers-go v0.26.0
```

### **Development Dependencies**
```go
// Hot Reload
github.com/cosmtrek/air v1.49.0

// Linting
github.com/golangci/golangci-lint v1.55.2

// Code Generation
github.com/swaggo/swag v1.16.2
```

## 🚀 Quick Start Commands

### **Development**
```bash
# Setup
make setup

# Run development server
make dev

# Run tests
make test

# Run migrations
make migrate-up

# Lint code
make lint
```

### **Docker**
```bash
# Development
docker-compose up --build

# Production
docker-compose -f deploy/docker-compose.yml up --build
```

### **Database**
```bash
# Create migration
make migrate-create name=add_user_table

# Run migrations
make migrate-up

# Rollback migration
make migrate-down
```

## 📋 Environment Variables

### **Required Variables**
```bash
# Application
APP_ENV=development
APP_PORT=8080
APP_HOST=0.0.0.0

# Database
DB_HOST=localhost
DB_PORT=5433
DB_USER=app
DB_PASS=secret
DB_NAME=myapp
DB_SSLMODE=disable

# JWT
JWT_SECRET=your-super-secret-key
JWT_ACCESS_TOKEN_EXPIRE=15m
JWT_REFRESH_TOKEN_EXPIRE=7d

# Cloudflare R2
R2_ACCOUNT_ID=your-account-id
R2_ACCESS_KEY=your-access-key
R2_SECRET_KEY=your-secret-key
R2_BUCKET=your-bucket-name
R2_ENDPOINT=https://your-account-id.r2.cloudflarestorage.com

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

## 🎯 Key Features by Directory

### **`/cmd/server`** - Application Entry Point
- Main application bootstrap
- Graceful shutdown handling
- Configuration loading
- Database connection
- Router setup

### **`/internal/api/v1`** - API Layer
- HTTP handlers for all endpoints
- Request/response validation
- Business logic services
- Route registration

### **`/internal/db`** - Database Layer
- Database migrations
- Repository pattern implementation
- Database connection management
- Transaction handling

### **`/internal/middleware`** - Middleware
- Authentication middleware
- CORS configuration
- Request logging
- Rate limiting
- Security headers

### **`/internal/services`** - Business Logic
- Authentication services
- File storage services
- User management services
- JWT token management

### **`/pkg`** - Shared Packages
- Configuration management
- Logging utilities
- Database utilities
- Reusable components

This structure provides a **clean, scalable, and maintainable** foundation for any mobile backend application! 🚀
