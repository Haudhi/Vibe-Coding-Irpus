.PHONY: help build run test test-coverage lint clean docker-build docker-run migrate-up migrate-down

# Application variables
APP_NAME=ga-ticketing
VERSION=1.0.0
BUILD_DIR=build
SERVER_BIN=$(BUILD_DIR)/server
MIGRATE_BIN=$(BUILD_DIR)/migrate

# Go variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Docker variables
DOCKER_IMAGE=$(APP_NAME):$(VERSION)
DOCKER_COMPOSE_FILE=docker-compose.yml

# Database variables
DB_URL=postgres://ga_user:ga_password@localhost:5432/ga_ticketing?sslmode=disable
TEST_DB_URL=postgres://ga_user:ga_password@localhost:5432/ga_ticketing_test?sslmode=disable

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Setup and dependencies
setup: ## Setup development environment
	@echo "Setting up development environment..."
	$(GOMOD) download
	$(GOMOD) tidy
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/cosmtrek/air@latest
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Build targets
build: ## Build the application
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(SERVER_BIN) -v cmd/server/main.go
	$(GOBUILD) -o $(MIGRATE_BIN) -v cmd/migrate/main.go

build-all: ## Build application for multiple platforms
	@echo "Building $(APP_NAME) for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 -v cmd/server/main.go
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 -v cmd/server/main.go
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe -v cmd/server/main.go

# Development targets
run: ## Run the application
	@echo "Running $(APP_NAME)..."
	$(GOCMD) run cmd/server/main.go

run-dev: ## Run the application with hot reload
	@echo "Running $(APP_NAME) with hot reload..."
	air

# Test targets
test: ## Run all tests
	@echo "Running tests..."
	$(GOTEST) -v ./...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-unit: ## Run unit tests only
	@echo "Running unit tests..."
	$(GOTEST) -v ./tests/unit/...

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	$(GOTEST) -v -tags=integration ./tests/integration/...

test-contract: ## Run contract tests
	@echo "Running contract tests..."
	$(GOTEST) -v -tags=contract ./tests/contract/...

# Quality targets
lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run

fmt: ## Format code
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	$(GOCMD) vet ./...

quality: fmt vet lint test-coverage ## Run all quality checks

# Database targets
migrate-up: ## Run database migrations
	@echo "Running database migrations..."
	migrate -path migrations/ -database "$(DB_URL)" up

migrate-down: ## Rollback database migrations
	@echo "Rolling back database migrations..."
	migrate -path migrations/ -database "$(DB_URL)" down

migrate-create: ## Create new migration (usage: make migrate-create NAME=migration_name)
	@if [ -z "$(NAME)" ]; then echo "Error: NAME parameter is required. Usage: make migrate-create NAME=migration_name"; exit 1; fi
	@echo "Creating migration: $(NAME)..."
	migrate create -ext sql -dir migrations/ -seq $(NAME)

# Docker targets
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

docker-run: ## Run application with Docker
	@echo "Running $(APP_NAME) with Docker..."
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE)

docker-compose-up: ## Start development environment with Docker Compose
	@echo "Starting development environment..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

docker-compose-down: ## Stop development environment
	@echo "Stopping development environment..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

# Documentation targets
docs: ## Generate API documentation
	@echo "Generating API documentation..."
	swag init -g cmd/server/main.go -o docs/api

# Database test setup
test-db-setup: ## Setup test database
	@echo "Setting up test database..."
	createdb -U ga_user ga_ticketing_test || true
	migrate -path migrations/ -database "$(TEST_DB_URL)" up

test-db-teardown: ## Teardown test database
	@echo "Tearing down test database..."
	dropdb -U ga_user ga_ticketing_test || true

# CI/CD targets
ci: quality ## Run CI pipeline locally

pre-commit: fmt vet lint test ## Run pre-commit checks

# Cleanup targets
clean: ## Clean build artifacts
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	rm -rf tmp/
	rm -f .air.toml

clean-docker: ## Clean Docker images and containers
	@echo "Cleaning Docker resources..."
	docker system prune -f
	docker volume prune -f

# Development shortcuts
dev-setup: setup test-db-setup ## Complete development setup
	@echo "Development environment ready!"

dev-reset: clean test-db-teardown test-db-setup ## Reset development environment
	@echo "Development environment reset!"

# Production targets
release: clean quality docker-build ## Build release version
	@echo "Release build completed!"

deploy-staging: ## Deploy to staging
	@echo "Deploying to staging..."
	# Add staging deployment commands here

deploy-production: ## Deploy to production
	@echo "Deploying to production..."
	# Add production deployment commands here

# Performance targets
benchmark: ## Run performance benchmarks
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

profile: ## Run application profiling
	@echo "Starting application with profiling..."
	$(GOCMD) run -cpuprofile=cpu.prof -memprofile=mem.prof cmd/server/main.go