.PHONY: build run test clean docker-build docker-run deps lint security

# Variables
APP_NAME=golang-rest-api
VERSION=1.0.0
BUILD_DIR=build
DOCKER_IMAGE=$(APP_NAME):$(VERSION)

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) .

# Run the application
run:
	@echo "Running $(APP_NAME)..."
	@go run main.go

# Run tests
test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Lint code
lint:
	@echo "Running linter..."
	@golangci-lint run

# Security scan
security:
	@echo "Running security scan..."
	@gosec ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

# Docker build
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .

# Docker run
docker-run:
	@echo "Running Docker container..."
	@docker-compose up -d

# Docker stop
docker-stop:
	@echo "Stopping Docker containers..."
	@docker-compose down

# Database migration
migrate:
	@echo "Running database migration..."
	@mysql -h127.0.0.1 -uroot -ppassword goblog < schema/schema.sql

# Development setup
dev-setup: deps
	@echo "Setting up development environment..."
	@cp .env.example .env
	@echo "Please edit .env file with your configuration"

# Production build
prod-build:
	@echo "Building for production..."
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w -s' -o $(BUILD_DIR)/$(APP_NAME) .

# Help
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  deps          - Install dependencies"
	@echo "  lint          - Run linter"
	@echo "  security      - Run security scan"
	@echo "  clean         - Clean build artifacts"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run with Docker Compose"
	@echo "  docker-stop   - Stop Docker containers"
	@echo "  migrate       - Run database migration"
	@echo "  dev-setup     - Setup development environment"
	@echo "  prod-build    - Build for production"
	@echo "  help          - Show this help"