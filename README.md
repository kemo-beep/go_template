# Go Mobile Backend Template

A production-ready Gin (Go) backend template for mobile apps with PostgreSQL, Docker Compose deployment, and Cloudflare R2 storage.

## 🚀 Features

- **Gin Framework** - Fast, minimal, production-ready HTTP web framework
- **PostgreSQL** - Robust relational database with GORM ORM
- **JWT Authentication** - Secure authentication with refresh tokens
- **Cloudflare R2** - S3-compatible object storage for files
- **Docker Support** - Complete containerization with Docker Compose
- **Database Migrations** - Schema management with Goose
- **Structured Logging** - JSON logging with Zap
- **API Documentation** - Auto-generated Swagger docs
- **Rate Limiting** - Built-in request throttling
- **Security Headers** - Production-ready security middleware
- **Hot Reload** - Development with Air for fast iteration

## 📁 Project Structure

```
go-mobile-backend-template/
├── cmd/server/           # Application entry point
├── config/               # Configuration files
├── internal/
│   ├── api/v1/          # API handlers (auth, users, files)
│   ├── db/              # Database models and migrations
│   ├── middleware/       # HTTP middleware
│   ├── services/         # Business logic
│   └── utils/           # Utility functions
├── pkg/                  # Shared packages
├── deploy/               # Production deployment files
├── tests/                # Test files
└── docs/                 # Documentation
```

## 🛠️ Tech Stack

- **Go 1.23+** - Programming language
- **Gin** - HTTP web framework
- **PostgreSQL 15** - Database
- **GORM** - ORM
- **Goose** - Database migrations
- **JWT** - Authentication
- **Cloudflare R2** - File storage
- **Zap** - Logging
- **Docker** - Containerization
- **Swagger** - API documentation

## 🚀 Quick Start

### Prerequisites

- Go 1.23+
- Docker & Docker Compose
- Make (optional)

### 1. Clone and Setup

```bash
git clone <repository-url>
cd go-mobile-backend-template
make setup
```

### 2. Configure Environment

```bash
cp env.example .env
# Edit .env with your configuration
```

### 3. Start Development

```bash
# Using Docker Compose
make docker-run

# Or using Air for hot reload
make dev
```

### 4. Run Migrations

```bash
make migrate-up
```

## 📚 API Endpoints

### Authentication
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh token
- `POST /api/v1/auth/logout` - User logout

### Users
- `GET /api/v1/users/me` - Get current user
- `PUT /api/v1/users/me` - Update user profile
- `DELETE /api/v1/users/me` - Delete user account

### Files
- `POST /api/v1/files/upload` - Upload file
- `GET /api/v1/files/:id` - Get file info
- `GET /api/v1/files/:id/download` - Download file
- `DELETE /api/v1/files/:id` - Delete file

### Health
- `GET /healthz` - Health check

## 🔧 Development

### Available Commands

```bash
make build          # Build the application
make test           # Run tests
make coverage       # Run tests with coverage
make dev            # Run with hot reload
make lint           # Lint code
make format         # Format code
make migrate-up     # Run migrations
make migrate-down   # Rollback migrations
make docker-run     # Run with Docker
make swagger        # Generate API docs
```

### Environment Variables

```bash
# Application
APP_ENV=development
APP_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5433
DB_USER=app
DB_PASS=secret
DB_NAME=myapp

# JWT
JWT_SECRET=your-secret-key

# Cloudflare R2
R2_ACCOUNT_ID=your-account-id
R2_ACCESS_KEY=your-access-key
R2_SECRET_KEY=your-secret-key
R2_BUCKET=your-bucket
R2_ENDPOINT=https://your-account-id.r2.cloudflarestorage.com
```

## 🐳 Docker Deployment

### Development

```bash
docker-compose up --build
```

### Production

```bash
docker-compose -f deploy/docker-compose.yml up --build
```

## 📖 API Documentation

Once the server is running, visit:
- Swagger UI: `http://localhost:8080/docs/`
- Health Check: `http://localhost:8080/healthz`

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Run specific test
go test ./internal/api/v1/auth/...
```

## 🔒 Security Features

- JWT-based authentication
- Password hashing with bcrypt
- CORS protection
- Rate limiting
- Security headers
- Input validation
- SQL injection protection

## 📝 License

MIT License - see LICENSE file for details

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📞 Support

For support and questions, please open an issue on GitHub.
