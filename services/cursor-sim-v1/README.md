# cursor-sim - Cursor API Simulator

A Go-based simulator that generates synthetic usage data mimicking the Cursor Business API.

## Project Structure

```
cursor-sim/
├── cmd/
│   └── simulator/       # Main application entry point
│       └── main.go
├── internal/           # Private packages
│   ├── api/           # HTTP handlers (to be implemented)
│   ├── config/        # Configuration management ✓
│   ├── generator/     # Data generation logic (to be implemented)
│   ├── models/        # Domain types ✓
│   └── storage/       # In-memory storage (to be implemented)
├── Dockerfile         # Multi-stage Docker build
├── Makefile          # Build commands
├── .golangci.yml     # Linter configuration
└── go.mod            # Go module definition
```

## Quick Start

### Build

```bash
make build
```

### Run Tests

```bash
make test
```

### Run with Coverage

```bash
make test-coverage
```

### Format Code

```bash
make fmt
```

### Run Linter

```bash
make lint
```

### Build and Run

```bash
make run
```

### Run in Development Mode

```bash
make dev
```

### Docker Build

```bash
make docker-build
```

### Docker Run

```bash
make docker-run
```

## Development Status

**Phase**: P1 - Core Functionality
**Current Task**: TASK-SIM-002 ✓ Complete

### Completed
- ✓ TASK-SIM-001: Go project structure with internal packages
- ✓ TASK-SIM-002: CLI flag parsing with validation
  - Command-line flag parsing with standard library
  - JSON configuration file support
  - Comprehensive validation with helpful error messages
  - --help output with examples
  - Tests (80.4% coverage for config package)

### Next Steps
- TASK-SIM-003: Implement Developer Profile Generator
- TASK-SIM-004: Implement Event Generation Engine
- TASK-SIM-005: Implement In-Memory Storage
- TASK-SIM-006: Implement REST API Handlers
- TASK-SIM-007: Wire Up Main Application

## Configuration

### Command-Line Flags

Run with custom configuration:
```bash
./bin/cursor-sim --port 9000 --developers 100 --velocity high --fluctuation 0.3 --seed 42
```

View all options:
```bash
./bin/cursor-sim --help
```

### Configuration File

You can also use a JSON configuration file (see `config.example.json`):
```json
{
  "Port": 8080,
  "Developers": 100,
  "Velocity": "high",
  "Fluctuation": 0.3,
  "Seed": 42
}
```

Load from file:
```bash
# Future: ./bin/cursor-sim --config config.json
```

### Default Values

- Port: 8080 (range: 1024-65535)
- Developers: 50 (range: 1-10000)
- Velocity: "medium" (options: low|medium|high)
- Fluctuation: 0.2 (range: 0.0-1.0)
- Seed: 12345

## Testing

Run all tests:
```bash
go test ./...
```

Run with coverage:
```bash
go test -cover ./...
```

Generate coverage report:
```bash
make test-coverage
open coverage.html
```

## Code Quality

This project follows Go best practices:
- Standard project layout with `cmd/` and `internal/`
- Table-driven tests
- Error wrapping with context
- Exported function documentation
- Linting with golangci-lint

See `.claude/skills/go-best-practices.md` for complete guidelines.
