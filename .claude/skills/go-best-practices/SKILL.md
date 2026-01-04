---
name: go-best-practices
description: Go coding standards for cursor-sim service. Use when writing Go code, implementing HTTP handlers, creating structs, handling errors, or writing Go tests. Covers naming conventions, error handling, concurrency, and testing patterns. (project)
---

# Go Best Practices

This skill provides Go-specific coding standards for the Cursor Analytics Platform.

## Code Organization

### Project Structure

```
services/cursor-sim/
├── cmd/
│   └── simulator/
│       └── main.go          # Entry point only
├── internal/                 # Private packages
│   ├── config/              # Configuration
│   ├── models/              # Domain types
│   ├── generator/           # Business logic
│   ├── api/                 # HTTP handlers
│   └── db/                  # Data access
└── pkg/                     # Public packages (if any)
```

**Rules:**
1. `cmd/` contains only `main.go` and initialization
2. `internal/` for all private implementation
3. `pkg/` only if code is meant to be imported by other services
4. One package per directory
5. Package names match directory names

## Naming Conventions

### Packages
- **Lowercase**, no underscores: `config`, `generator`, `api`
- **Singular nouns**: `model` not `models`
- **Short but descriptive**: `auth` not `authentication`

### Files
- **Snake case**: `developer_generator.go`
- **Test files**: `developer_generator_test.go`
- **Group by functionality**: Not by type

### Types
- **PascalCase** for exported: `Developer`, `CommitGenerator`
- **camelCase** for unexported: `developerCache`, `configLoader`

### Functions
- **PascalCase** for exported: `NewDeveloper()`, `GenerateCommit()`
- **camelCase** for unexported: `parseDate()`, `validateConfig()`

### Constants

```go
const (
    DefaultPort      = 8080
    MaxDevelopers    = 10000
    MinDevelopers    = 1
    DefaultVelocity  = "high"
)
```

## Error Handling

### Pattern 1: Return Errors, Don't Panic

```go
// GOOD
func LoadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read config: %w", err)
    }
    return &cfg, nil
}

// BAD - Don't panic in library code
func LoadConfig(path string) *Config {
    data, err := os.ReadFile(path)
    if err != nil {
        panic(err)
    }
}
```

### Pattern 2: Wrap Errors with Context

```go
if err := validateDeveloperCount(cfg.Developers); err != nil {
    return fmt.Errorf("config validation failed: %w", err)
}

// Check wrapped errors
if errors.Is(err, ErrInvalidConfig) {
    // Handle specific error
}
```

### Pattern 3: Custom Error Types

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}
```

## Struct Patterns

### Constructor Pattern

```go
func NewDeveloper(id, email, region string) *Developer {
    return &Developer{
        ID:        id,
        Email:     email,
        Region:    region,
        CreatedAt: time.Now().UTC(),
    }
}
```

### Validation Methods

```go
func (d *Developer) Validate() error {
    if d.ID == "" {
        return errors.New("developer ID is required")
    }
    return nil
}
```

## Interface Design

### Keep Interfaces Small

```go
// GOOD - Single responsibility
type DeveloperStore interface {
    Save(ctx context.Context, dev *Developer) error
    Get(ctx context.Context, id string) (*Developer, error)
    List(ctx context.Context, filters Filters) ([]*Developer, error)
}
```

### Accept Interfaces, Return Structs

```go
// GOOD - Returns concrete type
func NewMemoryStore() *MemoryStore {
    return &MemoryStore{
        data: make(map[string]*Developer),
    }
}
```

## Concurrency Patterns

### Use sync.WaitGroup for Goroutine Coordination

```go
func GenerateForAllDevelopers(developers []*Developer) {
    var wg sync.WaitGroup
    for _, dev := range developers {
        wg.Add(1)
        go func(d *Developer) {
            defer wg.Done()
            generateEvents(d)
        }(dev)
    }
    wg.Wait()
}
```

### Protect Shared State with Mutex

```go
type MemoryStore struct {
    mu   sync.RWMutex
    data map[string]*Developer
}

func (s *MemoryStore) Get(id string) (*Developer, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    dev, ok := s.data[id]
    if !ok {
        return nil, ErrNotFound
    }
    return dev, nil
}
```

## Testing Standards

### Table-Driven Tests

```go
func TestValidateConfig(t *testing.T) {
    tests := []struct {
        name    string
        config  Config
        wantErr bool
    }{
        {"valid config", Config{Developers: 50, Velocity: "high"}, false},
        {"negative developers", Config{Developers: -1}, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateConfig(&tt.config)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Use testify for Assertions

```go
func TestDeveloperCreation(t *testing.T) {
    dev := NewDeveloper("id123", "test@example.com", "US")
    assert.Equal(t, "id123", dev.ID)
    require.NoError(t, dev.Validate())
}
```

## HTTP Handler Pattern

```go
func (h *Handler) GetDevelopers(w http.ResponseWriter, r *http.Request) {
    // 1. Parse and validate inputs
    page := parseIntParam(r, "page", 1)

    // 2. Call business logic
    developers, total, err := h.store.ListDevelopers(r.Context(), page, pageSize)
    if err != nil {
        h.writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
        return
    }

    // 3. Build and write response
    h.writeJSON(w, http.StatusOK, response)
}
```

## Code Review Checklist

- [ ] All errors are handled (no `_` ignores)
- [ ] All exported functions have comments
- [ ] All tests pass (`go test ./...`)
- [ ] Code is formatted (`gofmt -s -w .`)
- [ ] No panics in library code
- [ ] HTTP handlers return proper status codes
