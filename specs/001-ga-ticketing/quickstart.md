# Quick Start Guide: GA Ticketing System

**Branch**: `001-ga-ticketing` | **Date**: 2025-11-06

## Overview

This guide helps developers get the GA Ticketing System running locally for development and testing. The system provides a REST API for General Affairs service request management, approval workflows, and asset inventory management.

## Prerequisites

### System Requirements
- **Go**: 1.22 or higher
- **PostgreSQL**: 15.0 or higher
- **Git**: Latest stable version
- **Operating System**: Linux, macOS, or Windows (with WSL2)

### Development Tools (Recommended)
- **IDE/Editor**: VS Code, GoLand, or Vim with Go extensions
- **API Client**: Postman, Insomnia, or curl
- **Database Tool**: pgAdmin, DBeaver, or psql CLI
- **Git Client**: SourceTree, GitHub Desktop, or Git CLI

## Local Development Setup

### 1. Clone Repository and Setup Environment

```bash
# Clone the repository
git clone <repository-url>
cd ga-ticketing

# Switch to the feature branch
git checkout 001-ga-ticketing

# Install Go dependencies
go mod download

# Install development tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/swaggo/swag/cmd/swag@latest
```

### 2. Database Setup

```bash
# Install PostgreSQL if not already installed
# Ubuntu/Debian:
sudo apt update
sudo apt install postgresql postgresql-contrib

# macOS (with Homebrew):
brew install postgresql
brew services start postgresql

# Create database and user
sudo -u postgres psql
```

```sql
-- In PostgreSQL prompt
CREATE DATABASE ga_ticketing;
CREATE USER ga_user WITH PASSWORD 'ga_password';
GRANT ALL PRIVILEGES ON DATABASE ga_ticketing TO ga_user;
ALTER USER ga_user CREATEDB;
\q
```

### 3. Environment Configuration

Create environment configuration file:

```bash
# Create environment file
cat > .env << EOF
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=ga_ticketing
DB_USER=ga_user
DB_PASSWORD=ga_password
DB_SSLMODE=disable

# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_READ_TIMEOUT=30s
SERVER_WRITE_TIMEOUT=30s

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRY=24h
JWT_REFRESH_EXPIRY=168h

# OAuth2 Configuration (for production)
OAUTH2_CLIENT_ID=your-oauth2-client-id
OAUTH2_CLIENT_SECRET=your-oauth2-client-secret
OAUTH2_REDIRECT_URL=http://localhost:8080/auth/callback

# Redis Configuration (for caching)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Logging Configuration
LOG_LEVEL=info
LOG_FORMAT=json

# Application Configuration
APP_NAME=GA Ticketing System
APP_VERSION=1.0.0
APP_ENV=development
EOF
```

### 4. Database Migration

```bash
# Run database migrations
go run cmd/migrate/main.go up

# Or use the Makefile if available
make migrate-up
```

### 5. Start the Development Server

```bash
# Run the server
go run cmd/server/main.go

# Or use the Makefile
make run

# For development with hot reload (using air)
air
```

The server should start at `http://localhost:8080`

## Database Schema Overview

### Core Tables
- **users**: User accounts and role management
- **tickets**: Service request tickets with workflow states
- **assets**: Physical inventory and asset management
- **comments**: Ticket communication and updates
- **approvals**: Approval workflow records
- **status_history**: Ticket status change audit trail
- **inventory_logs**: Asset quantity change history

### Key Relationships
- Users create tickets (requester role)
- Admins manage and process tickets
- Approvers handle budget approvals
- Tickets can have multiple comments and status changes
- Assets track inventory quantities and locations

## API Testing

### 1. Authentication

Get JWT token for API calls:

```bash
# Login as user (using curl)
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "employee@company.com",
    "password": "password123"
  }'

# Save the token for subsequent requests
export AUTH_TOKEN="your-jwt-token-here"
```

### 2. Basic API Operations

```bash
# Create a ticket
curl -X POST http://localhost:8080/v1/tickets \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Request Office Supplies",
    "description": "Need notebooks and pens for team",
    "category": "office_supplies",
    "priority": "medium",
    "estimated_cost": 250000
  }'

# Get user's tickets
curl -X GET "http://localhost:8080/v1/tickets?page=1&limit=10" \
  -H "Authorization: Bearer $AUTH_TOKEN"

# Get specific ticket details
curl -X GET http://localhost:8080/v1/tickets/{ticket-id} \
  -H "Authorization: Bearer $AUTH_TOKEN"
```

### 3. Using Postman Collection

1. Import the Postman collection from `specs/001-ga-ticketing/contracts/postman-collection.json`
2. Set environment variables:
   - `baseUrl`: `http://localhost:8080`
   - `authToken`: Your JWT token from login
3. Use the collection to test various API endpoints

## Development Workflow

### 1. Code Organization

```
src/
├── domain/          # Domain entities and business logic
├── application/     # Use cases and application services
├── infrastructure/  # Database, external services, frameworks
├── interface/       # HTTP handlers, controllers
├── config/          # Configuration management
└── shared/          # Common utilities and types

tests/
├── unit/            # Unit tests
├── integration/     # Integration tests
├── contract/        # API contract tests
└── e2e/            # End-to-end tests
```

### 2. Running Tests

```bash
# Run all tests
go test ./...

# Run specific test suites
go test ./tests/unit/...
go test ./tests/integration/...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run performance benchmarks
go test -bench=. ./tests/benchmarks/...
```

### 3. Code Quality

```bash
# Run linter
golangci-lint run

# Format code
go fmt ./...

# Vet for potential issues
go vet ./...

# Generate API documentation
swag init -g cmd/server/main.go
```

## Common Development Tasks

### 1. Adding New API Endpoints

1. Define request/response structs in `interface/http/dto/`
2. Add validation tags and documentation
3. Implement use case in `application/services/`
4. Create HTTP handler in `interface/http/handlers/`
5. Add route in `interface/http/router/`
6. Write tests for the new endpoint
7. Update API documentation

### 2. Database Schema Changes

1. Create migration file in `migrations/`
2. Add/up migration with schema changes
3. Add/down migration for rollback
4. Update domain models if needed
5. Run migration: `go run cmd/migrate/main.go up`
6. Test both up and down migrations

### 3. Adding New Business Logic

1. Define domain entities or value objects
2. Implement domain services with business rules
3. Create application use cases that orchestrate domain services
4. Write comprehensive unit tests
5. Add integration tests with database

## Troubleshooting

### Common Issues

**Database Connection Failed**
```bash
# Check PostgreSQL status
sudo systemctl status postgresql

# Verify database exists
psql -h localhost -U ga_user -d ga_ticketing

# Check connection string in .env
echo $DB_HOST $DB_PORT $DB_NAME
```

**JWT Token Issues**
```bash
# Verify JWT secret is set
grep JWT_SECRET .env

# Test token generation
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@company.com","password":"password"}'
```

**Permission Denied Errors**
- Verify user role in database
- Check JWT token contains correct role
- Review endpoint authorization rules

**Build Errors**
```bash
# Clean module cache
go clean -modcache
go mod download

# Verify Go version
go version

# Check for missing dependencies
go mod tidy
```

### Performance Issues

**Slow Database Queries**
```bash
# Enable query logging
export LOG_LEVEL=debug

# Add database indexes as needed
# See data-model.md for recommended indexes

# Use connection pooling
# Check DB_MAX_CONNECTIONS in .env
```

**High Memory Usage**
```bash
# Check for memory leaks
go tool pprof http://localhost:8080/debug/pprof/heap

# Monitor connection pools
# Review database connection settings
```

## Production Deployment

### 1. Environment Variables for Production

```bash
# Required production environment variables
export APP_ENV=production
export LOG_LEVEL=warn
export DB_SSLMODE=require
export DB_MAX_CONNECTIONS=50
export REDIS_PASSWORD=your-redis-password
export JWT_SECRET=your-production-jwt-secret
```

### 2. Docker Deployment

```dockerfile
# Build application
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main cmd/server/main.go

# Runtime image
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
CMD ["./main"]
```

### 3. Health Checks

```bash
# Application health endpoint
curl http://localhost:8080/health

# Database connectivity check
curl http://localhost:8080/health/db

# Readiness check
curl http://localhost:8080/ready
```

## Additional Resources

- **API Documentation**: [OpenAPI Spec](./contracts/api.yaml)
- **Postman Collection**: [API Testing](./contracts/postman-collection.json)
- **Data Model**: [Domain Entities](./data-model.md)
- **Implementation Plan**: [Project Architecture](./plan.md)
- **Feature Specification**: [Requirements](./spec.md)

## Getting Help

1. Check the [Troubleshooting](#troubleshooting) section first
2. Review the [Data Model](./data-model.md) for entity relationships
3. Consult the [API Documentation](./contracts/api.yaml) for endpoint details
4. Reference the [Feature Specification](./spec.md) for business requirements
5. Join the development Slack channel or create a GitHub issue for support

## Contributing

1. Follow the code style and quality standards defined in the project constitution
2. Write comprehensive tests for all new features
3. Update documentation for any API changes
4. Ensure all tests pass before submitting pull requests
5. Follow the Git workflow and branch naming conventions