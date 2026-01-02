.PHONY: help build test lint fmt clean run docker-build docker-run

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the simulator binary
	@echo "Building cursor-sim..."
	@go build -o bin/cursor-sim ./cmd/simulator

test: ## Run tests
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests with coverage report
	@echo "Generating coverage report..."
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

lint: ## Run linter (optional)
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not installed (skipping lint)"; \
		echo "   Install with: brew install golangci-lint"; \
		echo "   or visit https://golangci-lint.run/usage/install/"; \
	fi

fmt: ## Format code
	@echo "Formatting code..."
	@gofmt -s -w .
	@go mod tidy

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@go clean

run: build ## Build and run the simulator
	@echo "Starting cursor-sim..."
	@./bin/cursor-sim

dev: ## Run in development mode
	@echo "Running in dev mode..."
	@go run ./cmd/simulator

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t cursor-sim:latest .

docker-run: docker-build ## Run in Docker
	@echo "Running in Docker..."
	@docker run -p 8080:8080 --rm cursor-sim:latest

all: fmt test build ## Run all checks and build
