# BookClubApp Makefile

.PHONY: help build run test clean deps lint format

# Variables
BINARY_NAME=bookclubapp
MAIN_FILE=main.go
BUILD_DIR=build

# Default command
help: ## Show this help
	@echo "BookClubApp - Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Development
deps: ## Install dependencies
	@echo "📦 Installing dependencies..."
	go mod tidy
	go mod download

run: ## Run server in development mode
	@echo "🚀 Starting server..."
	go run $(MAIN_FILE)

build: ## Compile the project
	@echo "🔨 Compiling project..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "✅ Binary created at $(BUILD_DIR)/$(BINARY_NAME)"

# Cleanup
clean: ## Remove build files
	@echo "🧹 Cleaning build..."
	@rm -rf $(BUILD_DIR)
	@go clean

# Code quality
lint: ## Run linter
	@echo "🔍 Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

format: ## Format code
	@echo "✨ Formatting code..."
	go fmt ./...
	go vet ./...

# Tests
test: ## Run tests
	@echo "🧪 Running tests..."
	go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "🧪 Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report generated: coverage.html"

# Docker
docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	docker build -t bookclubapp:latest .

docker-run: ## Run Docker container
	@echo "🐳 Running Docker container..."
	docker run -p 8080:8080 bookclubapp:latest

# Development
dev: deps run ## Install dependencies and run server

# Installation
install: ## Install the project
	@echo "📥 Installing BookClubApp..."
	go install .

# Verification
check: format lint test ## Format, verify and test code

# Complete development
all: clean deps build test ## Run complete development pipeline
