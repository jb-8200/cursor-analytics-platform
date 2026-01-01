# Task: TASK-001

## Set Up Go Project Structure for Simulator

**Task ID:** TASK-001  
**Service:** cursor-sim  
**Feature:** [F001 - Simulator Core](../features/F001-simulator-core.md)  
**User Story:** [US-SIM-001](../user-stories/US-SIM-001-configure-simulation.md)  
**Estimated Hours:** 4  
**Status:** Ready

---

## Objective

Initialize the Go project with proper structure, dependencies, and tooling to support test-driven development of the Cursor API Simulator.

---

## Prerequisites

Before starting this task, ensure you have Go 1.22 or later installed, Docker and Docker Compose are available on your machine, and you have access to the project repository.

---

## Specification Reference

This task implements the foundation for Feature F001. Review the complete feature specification for context on how this setup will be used.

---

## Implementation Steps

### Step 1: Initialize Go Module

Create the go.mod file with the appropriate module path. The module name should reflect the project structure and be suitable for import by other services if needed in the future.

```bash
cd services/cursor-sim
go mod init github.com/cursor-analytics/cursor-sim
```

### Step 2: Create Directory Structure

The project follows a standard Go project layout with internal packages for non-exported code and cmd for the main entry point.

```
services/cursor-sim/
├── cmd/
│   └── cursor-sim/
│       └── main.go           # Entry point
├── internal/
│   ├── api/
│   │   ├── handlers.go       # HTTP handlers
│   │   ├── handlers_test.go
│   │   ├── router.go         # Chi router setup
│   │   └── middleware.go     # Logging, recovery
│   ├── generator/
│   │   ├── developer.go      # Developer profile generation
│   │   ├── developer_test.go
│   │   ├── event.go          # Event generation
│   │   ├── event_test.go
│   │   ├── poisson.go        # Poisson distribution
│   │   └── poisson_test.go
│   ├── storage/
│   │   ├── sqlite.go         # In-memory SQLite
│   │   └── sqlite_test.go
│   └── config/
│       ├── config.go         # Configuration struct
│       └── config_test.go
├── pkg/
│   └── models/
│       ├── developer.go      # Exported types
│       └── event.go
├── go.mod
├── go.sum
├── Dockerfile
├── Makefile
└── README.md
```

### Step 3: Add Dependencies

Install the required dependencies using go get. These are specified in the feature document.

```bash
go get github.com/go-chi/chi/v5
go get github.com/mattn/go-sqlite3
go get github.com/google/uuid
go get github.com/spf13/cobra
go get github.com/spf13/viper
go get github.com/stretchr/testify
go get github.com/go-faker/faker/v4
```

### Step 4: Create Main Entry Point

Create a minimal main.go that initializes and runs the application. This should be a thin wrapper that delegates to internal packages.

```go
// cmd/cursor-sim/main.go
package main

import (
	"os"

	"github.com/cursor-analytics/cursor-sim/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
```

### Step 5: Set Up Makefile

Create a Makefile with common development commands. This standardizes the developer experience across the team.

```makefile
# Makefile
.PHONY: build test lint run clean

BINARY_NAME=cursor-sim
BUILD_DIR=./build

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/cursor-sim

test:
	go test -v -race -cover ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./...

run:
	go run ./cmd/cursor-sim --developers=50 --velocity=medium

clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

docker-build:
	docker build -t cursor-sim:latest .

docker-run:
	docker run -p 8080:8080 cursor-sim:latest
```

### Step 6: Create Initial Dockerfile

Create a multi-stage Dockerfile for efficient builds. The first stage compiles the Go binary, and the second stage creates a minimal runtime image.

```dockerfile
# Dockerfile
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod files first for layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=1 go build -o cursor-sim ./cmd/cursor-sim

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Install runtime dependencies for SQLite
RUN apk add --no-cache ca-certificates

# Copy binary from builder
COPY --from=builder /app/cursor-sim .

# Expose default port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["./cursor-sim"]
```

### Step 7: Write Initial Tests

Create placeholder test files that will fail, establishing the TDD cycle. These tests define the expected behavior before implementation.

```go
// internal/config/config_test.go
package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/cursor-analytics/cursor-sim/internal/config"
)

func TestDefaultConfig(t *testing.T) {
	cfg := config.Default()
	
	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, 50, cfg.Developers)
	assert.Equal(t, "medium", cfg.Velocity)
	assert.Equal(t, 0.2, cfg.Fluctuation)
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  config.Config
		wantErr bool
	}{
		{
			name:    "valid config",
			config:  config.Config{Port: 8080, Developers: 50, Velocity: "medium"},
			wantErr: false,
		},
		{
			name:    "negative developers",
			config:  config.Config{Port: 8080, Developers: -1, Velocity: "medium"},
			wantErr: true,
		},
		{
			name:    "invalid velocity",
			config:  config.Config{Port: 8080, Developers: 50, Velocity: "turbo"},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
```

### Step 8: Verify Setup

Run the following commands to verify the setup is correct. All commands should execute without error, though tests will fail as expected until implementation.

```bash
# Verify module is valid
go mod tidy

# Verify build works (will fail, but should parse)
go build ./... 2>&1 | head -20

# Verify test framework works
go test ./... -v

# Verify linting works
golangci-lint run ./...
```

---

## Acceptance Criteria Mapping

This task contributes to the following acceptance criteria from US-SIM-001:

| Criteria | Task Contribution |
|----------|-------------------|
| AC1: Help flag | Project structure supports CLI implementation |
| AC7: Startup time | Build configuration optimized for performance |

---

## Test-Driven Development Checklist

Before marking this task complete, ensure the TDD process was followed.

- [ ] Wrote test file structure before implementation
- [ ] Tests currently fail (as expected for a setup task)
- [ ] Test file naming follows `*_test.go` convention
- [ ] Test functions follow `Test*` naming convention
- [ ] Table-driven tests used where appropriate

---

## Definition of Done

- [ ] Go module initialized with correct path
- [ ] Directory structure matches specification
- [ ] All dependencies installed and verified
- [ ] Makefile contains all required targets
- [ ] Dockerfile builds successfully
- [ ] Initial test files created
- [ ] `go mod tidy` runs without errors
- [ ] `golangci-lint` runs without errors (may have warnings)
- [ ] README.md created with basic project information
- [ ] Code committed and pushed to feature branch

---

## Notes for AI Code Assistants

When implementing this task, follow these guidelines.

First, create the directory structure before writing any code. Use `mkdir -p` commands to create all directories at once.

Second, always start with test files. Write the test that describes the expected behavior, then write the minimal code to make it compile (but fail).

Third, follow Go conventions for package naming—use short, lowercase names without underscores.

Fourth, use CGO_ENABLED=1 for the build because go-sqlite3 requires cgo.

Fifth, the Makefile should be the primary interface for developers. All common tasks should be available as make targets.
