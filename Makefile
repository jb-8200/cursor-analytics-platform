# Makefile for Cursor Analytics Platform
#
# This Makefile provides convenient commands for common development tasks.
# All commands are designed to work from the project root directory.
#
# Usage:
#   make help          - Show available commands
#   make setup         - Initial project setup
#   make dev           - Start all services in development mode
#   make test          - Run all tests
#   make test-coverage - Run tests with coverage reports

.PHONY: help setup dev stop clean test test-coverage lint build

# Default target shows help
.DEFAULT_GOAL := help

# =============================================================================
# Help
# =============================================================================

help: ## Show this help message
	@echo "Cursor Analytics Platform - Development Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# =============================================================================
# Setup
# =============================================================================

setup: ## Initial project setup (install dependencies)
	@echo "Setting up cursor-sim..."
	cd services/cursor-sim && go mod download
	@echo "Setting up cursor-analytics-core..."
	cd services/cursor-analytics-core && npm install
	@echo "Setting up cursor-viz-spa..."
	cd services/cursor-viz-spa && npm install
	@echo "Setup complete!"

setup-sim: ## Setup cursor-sim only
	cd services/cursor-sim && go mod download

setup-core: ## Setup cursor-analytics-core only
	cd services/cursor-analytics-core && npm install

setup-spa: ## Setup cursor-viz-spa only
	cd services/cursor-viz-spa && npm install

# =============================================================================
# Development
# =============================================================================

dev: ## Start all services with Docker Compose
	docker-compose up -d
	@echo ""
	@echo "Services starting..."
	@echo "  Simulator:  http://localhost:8080"
	@echo "  GraphQL:    http://localhost:4000/graphql"
	@echo "  Dashboard:  http://localhost:3000"
	@echo ""
	@echo "Run 'make logs' to view logs or 'make stop' to stop services."

dev-sim: ## Run cursor-sim locally (without Docker)
	cd services/cursor-sim && go run . --port=8080 --developers=50 --velocity=high

dev-core: ## Run cursor-analytics-core locally (requires PostgreSQL)
	cd services/cursor-analytics-core && npm run dev

dev-spa: ## Run cursor-viz-spa locally
	cd services/cursor-viz-spa && npm run dev

logs: ## Follow logs from all services
	docker-compose logs -f

logs-sim: ## Follow logs from cursor-sim
	docker-compose logs -f cursor-sim

logs-core: ## Follow logs from cursor-analytics-core
	docker-compose logs -f cursor-analytics-core

logs-spa: ## Follow logs from cursor-viz-spa
	docker-compose logs -f cursor-viz-spa

stop: ## Stop all services
	docker-compose down

clean: ## Stop services and remove volumes (reset all data)
	docker-compose down -v
	@echo "All containers and volumes removed."

restart: stop dev ## Restart all services

# =============================================================================
# Testing
# =============================================================================

test: ## Run all tests
	@echo "Testing cursor-sim..."
	cd services/cursor-sim && go test ./...
	@echo ""
	@echo "Testing cursor-analytics-core..."
	cd services/cursor-analytics-core && npm test
	@echo ""
	@echo "Testing cursor-viz-spa..."
	cd services/cursor-viz-spa && npm test

test-sim: ## Run cursor-sim tests
	cd services/cursor-sim && go test -v ./...

test-core: ## Run cursor-analytics-core tests
	cd services/cursor-analytics-core && npm test

test-spa: ## Run cursor-viz-spa tests
	cd services/cursor-viz-spa && npm test

test-coverage: ## Run all tests with coverage
	@echo "Testing cursor-sim with coverage..."
	cd services/cursor-sim && go test -coverprofile=coverage.out ./... && \
		go tool cover -func=coverage.out
	@echo ""
	@echo "Testing cursor-analytics-core with coverage..."
	cd services/cursor-analytics-core && npm test -- --coverage
	@echo ""
	@echo "Testing cursor-viz-spa with coverage..."
	cd services/cursor-viz-spa && npm test -- --coverage

test-watch: ## Run tests in watch mode (frontend only)
	cd services/cursor-viz-spa && npm test -- --watch

# =============================================================================
# Linting & Formatting
# =============================================================================

lint: ## Run linters on all services
	@echo "Linting cursor-sim..."
	cd services/cursor-sim && golangci-lint run
	@echo ""
	@echo "Linting cursor-analytics-core..."
	cd services/cursor-analytics-core && npm run lint
	@echo ""
	@echo "Linting cursor-viz-spa..."
	cd services/cursor-viz-spa && npm run lint

lint-fix: ## Run linters and fix auto-fixable issues
	@echo "Fixing cursor-sim..."
	cd services/cursor-sim && go fmt ./...
	@echo ""
	@echo "Fixing cursor-analytics-core..."
	cd services/cursor-analytics-core && npm run lint -- --fix
	@echo ""
	@echo "Fixing cursor-viz-spa..."
	cd services/cursor-viz-spa && npm run lint -- --fix

format: ## Format code in all services
	cd services/cursor-sim && go fmt ./...
	cd services/cursor-analytics-core && npm run format
	cd services/cursor-viz-spa && npm run format

# =============================================================================
# Building
# =============================================================================

build: ## Build all services
	@echo "Building cursor-sim..."
	cd services/cursor-sim && go build -o bin/cursor-sim .
	@echo ""
	@echo "Building cursor-analytics-core..."
	cd services/cursor-analytics-core && npm run build
	@echo ""
	@echo "Building cursor-viz-spa..."
	cd services/cursor-viz-spa && npm run build

build-docker: ## Build Docker images for all services
	docker-compose build

# =============================================================================
# Database
# =============================================================================

db-migrate: ## Run database migrations
	cd services/cursor-analytics-core && npm run migrate

db-migrate-rollback: ## Rollback last database migration
	cd services/cursor-analytics-core && npm run migrate:rollback

db-reset: ## Reset database (drop, create, migrate)
	cd services/cursor-analytics-core && npm run db:reset

db-shell: ## Open PostgreSQL shell
	docker-compose exec postgres psql -U cursor -d cursor_analytics

# =============================================================================
# Code Generation
# =============================================================================

codegen: ## Generate code (GraphQL types, etc.)
	cd services/cursor-analytics-core && npm run codegen
	cd services/cursor-viz-spa && npm run codegen

# =============================================================================
# Documentation
# =============================================================================

docs-serve: ## Serve documentation locally
	@echo "Documentation is in the docs/ directory"
	@echo "Key files:"
	@echo "  docs/DESIGN.md        - System architecture"
	@echo "  docs/FEATURES.md      - Feature breakdown"
	@echo "  docs/USER_STORIES.md  - User stories"
	@echo "  docs/TASKS.md         - Implementation tasks"

docs-validate: ## Validate documentation completeness
	@echo "Checking documentation..."
	@test -f docs/DESIGN.md || (echo "Missing: docs/DESIGN.md" && exit 1)
	@test -f docs/FEATURES.md || (echo "Missing: docs/FEATURES.md" && exit 1)
	@test -f docs/USER_STORIES.md || (echo "Missing: docs/USER_STORIES.md" && exit 1)
	@test -f docs/TASKS.md || (echo "Missing: docs/TASKS.md" && exit 1)
	@test -f services/cursor-sim/SPEC.md || (echo "Missing: cursor-sim SPEC.md" && exit 1)
	@test -f services/cursor-analytics-core/SPEC.md || (echo "Missing: cursor-analytics-core SPEC.md" && exit 1)
	@test -f services/cursor-viz-spa/SPEC.md || (echo "Missing: cursor-viz-spa SPEC.md" && exit 1)
	@echo "All documentation files present!"

# =============================================================================
# Health Checks
# =============================================================================

health: ## Check health of all running services
	@echo "Checking cursor-sim..."
	@curl -s http://localhost:8080/health | jq . || echo "Not responding"
	@echo ""
	@echo "Checking cursor-analytics-core..."
	@curl -s http://localhost:4000/health | jq . || echo "Not responding"
	@echo ""
	@echo "Checking cursor-viz-spa..."
	@curl -s -o /dev/null -w "%{http_code}" http://localhost:3000 || echo "Not responding"

# =============================================================================
# Utilities
# =============================================================================

ps: ## Show running containers
	docker-compose ps

stats: ## Show container resource usage
	docker stats --no-stream

shell-sim: ## Open shell in cursor-sim container
	docker-compose exec cursor-sim sh

shell-core: ## Open shell in cursor-analytics-core container
	docker-compose exec cursor-analytics-core sh

shell-spa: ## Open shell in cursor-viz-spa container
	docker-compose exec cursor-viz-spa sh

# =============================================================================
# Data Pipeline (P8 ETL)
# =============================================================================

pipeline: ## Run full ETL pipeline (extract, load, transform)
	./tools/run_pipeline.sh

extract: ## Extract data from cursor-sim API to Parquet
	python tools/api-loader/loader.py \
		--url $(or $(CURSOR_SIM_URL),http://localhost:8080) \
		--output ./data/raw \
		--api-key $(or $(API_KEY),cursor-sim-dev-key) \
		--start-date $(or $(START_DATE),90d)

load: ## Load Parquet files into DuckDB
	python tools/api-loader/duckdb_loader.py \
		--parquet-dir ./data/raw \
		--db-path ./data/analytics.duckdb

dbt-deps: ## Install dbt dependencies
	cd dbt && dbt deps

dbt-build: ## Run dbt models (build all)
	cd dbt && dbt build --target dev

dbt-test: ## Run dbt tests only
	cd dbt && dbt test --target dev

dbt-run: ## Run dbt models without tests
	cd dbt && dbt run --target dev

dbt-docs: ## Generate and serve dbt documentation
	cd dbt && dbt docs generate && dbt docs serve

ci-local: extract load dbt-build ## Run full pipeline for local CI verification

clean-data: ## Remove all generated data files
	rm -rf data/raw/*.parquet
	rm -f data/analytics.duckdb
	@echo "Data files cleaned"

query-duckdb: ## Open DuckDB CLI with analytics database
	duckdb data/analytics.duckdb
