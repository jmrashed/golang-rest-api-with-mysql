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

# Release management
prepare-release:
	@if [ -z "$(VERSION)" ]; then echo "Usage: make prepare-release VERSION=x.y.z"; exit 1; fi
	@bash scripts/prepare-release.sh $(VERSION)

create-release:
	@if [ -z "$(VERSION)" ]; then echo "Usage: make create-release VERSION=x.y.z"; exit 1; fi
	@bash scripts/create-release.sh $(VERSION)

# Build release binaries
release-build:
	@echo "Building release binaries..."
	@mkdir -p dist
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o dist/$(APP_NAME)-linux-amd64 .
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -o dist/$(APP_NAME)-windows-amd64.exe .
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -o dist/$(APP_NAME)-darwin-amd64 .

# Help
help:
	@echo "Available commands:"
	@echo "  build           - Build the application"
	@echo "  run             - Run the application"
	@echo "  test            - Run tests"
	@echo "  test-coverage   - Run tests with coverage"
	@echo "  deps            - Install dependencies"
	@echo "  lint            - Run linter"
	@echo "  security        - Run security scan"
	@echo "  clean           - Clean build artifacts"
	@echo "  docker-build    - Build Docker image"
	@echo "  docker-run      - Run with Docker Compose"
	@echo "  docker-stop     - Stop Docker containers"
	@echo "  migrate         - Run database migration"
	@echo "  dev-setup       - Setup development environment"
	@echo "  prod-build      - Build for production"
	@echo "  prepare-release - Prepare release (VERSION=x.y.z)"
	@echo "  create-release  - Create and push release tag (VERSION=x.y.z)"
	@echo "  release-build   - Build release binaries"
	@echo "  help            - Show this help"