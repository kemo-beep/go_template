# Production-Ready Gin (Go) Mobile Backend Template

A comprehensive specification for building a production-ready Gin (Go) backend template for mobile apps with PostgreSQL, Docker Compose deployment, and Cloudflare R2 storage.

## 🏗️ Project Structure

```
/cmd
   /server         # main.go entrypoint
/config            # configuration (env, YAML, etc.)
/internal
   /api
      /v1          # versioned APIs
         /auth
         /users
         /files
   /db
      /migrations  # goose or migrate files
      /repository  # DB access layer (GORM or sqlc)
   /middleware     # auth, logging, recovery, rate limiting
   /services       # business logic
   /utils          # helpers (hashing, validation, etc.)
/pkg
   /logger         # zap or logrus wrapper
   /config         # viper or env loader
```

### VPS + Docker-Compose Friendly Structure

```
/cmd
   /server           # main.go entrypoint
/config              # configs (env, YAML)
/deploy
   docker-compose.yml
   Dockerfile
   init.sql          # optional DB init
/internal
   /api/v1
      /auth
      /users
      /files         # R2 storage endpoints
   /db
      /migrations    # Goose migration files
      /repository
   /middleware
   /services
   /storage          # R2 integration
/pkg
   /config           # viper/env loader
   /logger           # zap logger
```

## ⚙️ Tech Stack

### Core Framework
- **Framework**: Gin (fast, minimal, production-ready)
- **Database**: PostgreSQL (with Goose migrations)
- **ORM**: GORM (popular, feature-rich) or sqlc (strict type safety)

### Authentication & Security
- **Auth**: JWT (with refresh tokens) + OAuth2 if needed
- **Password Hashing**: bcrypt/argon2id
- **Security**: CORS, Helmet-style headers, rate limiting (ulule/limiter)

### Configuration & Logging
- **Config**: Viper (environment-based config)
- **Logging**: Zap (structured JSON logs)
- **Validation**: Go-playground/validator

### Testing & Development
- **Tests**: testify for unit tests
- **Hot Reload**: Air for development
- **Linting**: golangci-lint

### Containerization & Deployment
- **Containerization**: Docker + docker-compose for local dev
- **Deployment**: VPS with Docker Compose
- **Storage**: Cloudflare R2 (S3-compatible object storage via aws-sdk-go-v2)
- **Reverse Proxy**: Caddy / Nginx (optional but recommended)

## 🔑 Key Features

### Auth & Users
- Sign up / login with email + password (hashed with bcrypt/argon2)
- JWT-based auth with refresh tokens
- Role-based access control (basic admin/user)

### Database & Migrations
- PostgreSQL schema managed with Goose
- GORM repositories for CRUD
- Transaction support

### API Standards
- REST (JSON), easy to extend to GraphQL if needed
- Versioned endpoints (/api/v1/...)
- Swagger/OpenAPI auto-docs with swaggo/swag

### DX Tools
- Makefile with commands: `make run`, `make migrate`, `make test`
- Hot reload via Air
- Pre-configured linter (golangci-lint)

### Production Readiness
- Graceful shutdown (context + signals)
- Structured logging
- Health check endpoint (/healthz)
- CI/CD pipeline template (GitHub Actions or Gitea CI)
- Environment-based config (dev, staging, prod)

## 🔧 Configuration

### Environment Variables

```bash
# Application
APP_ENV=production
APP_PORT=8080

# Database
DB_HOST=db
DB_PORT=5433
DB_USER=app
DB_PASS=secret
DB_NAME=myapp

# Authentication
JWT_SECRET=super-secret-key

# Cloudflare R2
R2_ACCOUNT_ID=xxxx
R2_ACCESS_KEY=xxxx
R2_SECRET_KEY=xxxx
R2_BUCKET=myapp-bucket
R2_ENDPOINT=https://<account_id>.r2.cloudflarestorage.com
```

## 🚀 Docker Setup

### docker-compose.yml

```yaml
version: '3.9'

services:
  app:
    build: ..
    container_name: myapp_backend
    depends_on:
      - db
    env_file:
      - ../.env
    ports:
      - "8080:8080"
    volumes:
      - ../:/app
    restart: unless-stopped

  db:
    image: postgres:15-alpine
    container_name: myapp_db
    environment:
      POSTGRES_USER: app
      POSTGRES_PASSWORD: secret
      POSTGRES_DB: myapp
    volumes:
      - db_data:/var/lib/postgresql/data
    restart: unless-stopped

volumes:
  db_data:
```

### Dockerfile

```dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server ./cmd/server

# Runtime
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/server .
COPY .env .env
EXPOSE 8080
CMD ["./server"]
```

## 📝 Example main.go

```go
package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/gin-gonic/gin"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"

    "myapp/internal/api/v1"
    "myapp/internal/middleware"
    "myapp/pkg/config"
    "myapp/pkg/logger"
)

func main() {
    // Load config
    cfg := config.Load()

    // Init logger
    log := logger.New(cfg.Env)

    // Connect DB
    db, err := gorm.Open(postgres.Open(cfg.DB.DSN()), &gorm.Config{})
    if err != nil {
        log.Fatal("failed to connect to db", err)
    }

    // Setup router
    r := gin.New()
    r.Use(gin.Recovery(), middleware.Logger(log), middleware.CORS())

    // Health check
    r.GET("/healthz", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    // API v1
    v1.RegisterRoutes(r.Group("/api/v1"), db, log)

    // Graceful shutdown
    srv := config.NewServer(r, cfg)
    go srv.Run()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Info("shutting down server")
    srv.Shutdown()
}
```

## ☁️ Cloudflare R2 Integration

### R2 Client Implementation

```go
package storage

import (
    "context"
    "fmt"
    "os"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/s3"
)

type R2Client struct {
    Client *s3.Client
    Bucket string
}

func NewR2Client() (*R2Client, error) {
    endpoint := os.Getenv("R2_ENDPOINT")
    accountID := os.Getenv("R2_ACCOUNT_ID")
    accessKey := os.Getenv("R2_ACCESS_KEY")
    secretKey := os.Getenv("R2_SECRET_KEY")
    bucket := os.Getenv("R2_BUCKET")

    customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, opts ...interface{}) (aws.Endpoint, error) {
        return aws.Endpoint{URL: endpoint}, nil
    })

    cfg, err := config.LoadDefaultConfig(context.TODO(),
        config.WithRegion("auto"),
        config.WithCredentialsProvider(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
            return aws.Credentials{
                AccessKeyID:     accessKey,
                SecretAccessKey: secretKey,
            }, nil
        })),
        config.WithEndpointResolverWithOptions(customResolver),
    )
    if err != nil {
        return nil, err
    }

    return &R2Client{
        Client: s3.NewFromConfig(cfg),
        Bucket: bucket,
    }, nil
}

func (r *R2Client) Upload(ctx context.Context, key string, body []byte) (string, error) {
    _, err := r.Client.PutObject(ctx, &s3.PutObjectInput{
        Bucket: &r.Bucket,
        Key:    &key,
        Body:   aws.ReadSeekCloser(bytes.NewReader(body)),
    })
    if err != nil {
        return "", err
    }
    return fmt.Sprintf("%s/%s", r.Bucket, key), nil
}
```

## 🌐 Example Endpoints

- `POST /api/v1/auth/register` → user registration
- `POST /api/v1/auth/login` → JWT login
- `GET /api/v1/users/me` → get profile (auth required)
- `POST /api/v1/files/upload` → upload file to R2

## 🛠️ Essential Add-Ons

### 1. Database & Schema Management
- **Goose** → run SQL or Go migrations inside Docker (`make migrate-up/down`)
- **Alternative**: Atlas → more declarative schema management
- **Why**: keeps schema consistent across environments (dev, staging, prod)

### 2. Authentication & Authorization
- JWT with refresh token flow
- Password hashing (bcrypt/argon2id)
- Optional: OAuth2/OpenID Connect for Google/Apple login
- Casbin (RBAC/ABAC) for role-based access control

### 3. Config & Secrets
- **Viper** for config (env vars + YAML)
- SOPS or Vault for secret management (if you don't want plain .env)

### 4. Observability
- **Logging** → Zap or Zerolog (structured, JSON)
- **Monitoring** → Prometheus + Grafana (expose /metrics endpoint with Gin middleware)
- **Error tracking** → Sentry (via github.com/getsentry/sentry-go)
- **Health checks** → /healthz + DB connection check

### 5. Testing
- Testify for unit/integration tests
- Spin up Postgres test container with testcontainers-go

### 6. Developer Experience
- Air → live reloading in dev
- golangci-lint → linting, style, static analysis
- Makefile with common commands:
  ```makefile
  run: docker-compose up --build
  migrate-up: goose -dir ./internal/db/migrations postgres "$(DB_DSN)" up
  test: go test ./...
  ```

### 7. API Documentation
- **Swaggo** → auto-generate Swagger docs from comments
- Expose at /docs

### 8. Security
- Rate limiting → ulule/limiter Gin middleware
- CORS config for mobile clients
- Helmet-style middleware for headers
- Dependency scanning (e.g. trivy)

### 9. File Storage (R2)
- Signed URLs for secure downloads
- Object lifecycle policies (clean old uploads)

### 10. Deployment Additions
- Nginx or Caddy reverse proxy with HTTPS (Let's Encrypt)
- Zero-downtime deploys → `docker-compose up -d --no-deps --build`
- Backups → nightly Postgres dumps (cron + pg_dump)

## 🔒 Production Considerations

- Run behind Nginx / Caddy with HTTPS (Let's Encrypt)
- Enable rate limiting middleware
- Store secrets with Vault or at least .env mounted securely
- Automate deploys with Gitea CI → Docker Compose up

## ⚡ Next Steps

This template provides:
- A ready-to-use template repo with Gin + PostgreSQL + JWT + migrations
- Docker + docker-compose for easy development setup
- Example routes (/auth/register, /auth/login, /users/me)
- Production-ready configuration for VPS deployment with Cloudflare R2 storage