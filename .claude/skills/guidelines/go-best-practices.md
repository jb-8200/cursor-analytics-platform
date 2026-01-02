# Go Best Practices Skill

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

---

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
- **PascalCase**: `DefaultPort`, `MaxDevelopers`
- **Grouped**: Use `const ()` blocks

```go
const (
    DefaultPort      = 8080
    MaxDevelopers    = 10000
    MinDevelopers    = 1
    DefaultVelocity  = "high"
)
```

---

## Error Handling

### Pattern 1: Return Errors, Don't Panic

```go
// ✓ GOOD
func LoadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read config: %w", err)
    }

    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("failed to parse config: %w", err)
    }

    return &cfg, nil
}

// ✗ BAD
func LoadConfig(path string) *Config {
    data, err := os.ReadFile(path)
    if err != nil {
        panic(err)  // Don't panic in library code
    }
    // ...
}
```

### Pattern 2: Wrap Errors with Context

```go
// Use fmt.Errorf with %w to wrap
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

// Usage
func ValidateConfig(cfg *Config) error {
    if cfg.Developers < 1 {
        return &ValidationError{
            Field:   "developers",
            Message: "must be at least 1",
        }
    }
    return nil
}
```

---

## Struct Patterns

### Use Struct Tags Correctly

```go
type Developer struct {
    ID       string    `json:"id"`              // JSON serialization
    Email    string    `json:"email"`
    Region   string    `json:"region"`
    CreateAt time.Time `json:"created_at"`
}
```

### Constructor Pattern

```go
// Always provide New* constructors
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
    if !isValidEmail(d.Email) {
        return errors.New("invalid email format")
    }
    return nil
}
```

---

## Interface Design

### Keep Interfaces Small

```go
// ✓ GOOD - Single responsibility
type DeveloperStore interface {
    Save(ctx context.Context, dev *Developer) error
    Get(ctx context.Context, id string) (*Developer, error)
    List(ctx context.Context, filters Filters) ([]*Developer, error)
}

// ✗ BAD - Too many responsibilities
type Store interface {
    SaveDeveloper(...)
    GetDeveloper(...)
    SaveCommit(...)
    GetCommit(...)
    SaveChange(...)
    GetChange(...)
    // ... 20 more methods
}
```

### Accept Interfaces, Return Structs

```go
// ✓ GOOD
func ProcessDevelopers(store DeveloperStore, gen DeveloperGenerator) error {
    // ...
}

// ✗ BAD - Returns interface
func NewDeveloperStore() DeveloperStore {
    return &memoryStore{}
}

// ✓ GOOD - Returns concrete type
func NewMemoryStore() *MemoryStore {
    return &MemoryStore{
        data: make(map[string]*Developer),
    }
}
```

---

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

func (s *MemoryStore) Save(dev *Developer) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.data[dev.ID] = dev
    return nil
}

func (s *MemoryStore) Get(id string) (*Developer, error) {
    s.mu.RLock()  // Read lock
    defer s.mu.RUnlock()

    dev, ok := s.data[id]
    if !ok {
        return nil, ErrNotFound
    }
    return dev, nil
}
```

### Use Context for Cancellation

```go
func FetchData(ctx context.Context) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    case <-time.After(5 * time.Second):
        // Do work
        return nil
    }
}
```

---

## Testing Standards

### Table-Driven Tests

```go
func TestValidateConfig(t *testing.T) {
    tests := []struct {
        name    string
        config  Config
        wantErr bool
    }{
        {
            name:    "valid config",
            config:  Config{Developers: 50, Velocity: "high"},
            wantErr: false,
        },
        {
            name:    "negative developers",
            config:  Config{Developers: -1},
            wantErr: true,
        },
        {
            name:    "invalid velocity",
            config:  Config{Developers: 50, Velocity: "invalid"},
            wantErr: true,
        },
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
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestDeveloperCreation(t *testing.T) {
    dev := NewDeveloper("id123", "test@example.com", "US")

    assert.Equal(t, "id123", dev.ID)
    assert.Equal(t, "US", dev.Region)
    assert.NotZero(t, dev.CreatedAt)

    // Use require for critical checks (stops test on failure)
    require.NoError(t, dev.Validate())
}
```

### Test Helpers

```go
// Put test helpers in _test.go files
func newTestDeveloper(t *testing.T) *Developer {
    t.Helper()  // Mark as helper
    return &Developer{
        ID:     "test-id",
        Email:  "test@example.com",
        Region: "US",
    }
}
```

---

## JSON Handling

### Use json.Marshal/Unmarshal

```go
// Marshaling
data, err := json.Marshal(developer)
if err != nil {
    return fmt.Errorf("failed to marshal: %w", err)
}

// Unmarshaling
var dev Developer
if err := json.Unmarshal(data, &dev); err != nil {
    return fmt.Errorf("failed to unmarshal: %w", err)
}
```

### Custom JSON Marshaling (When Needed)

```go
type Developer struct {
    ID        string
    CreatedAt time.Time
}

func (d Developer) MarshalJSON() ([]byte, error) {
    type Alias Developer
    return json.Marshal(&struct {
        CreatedAt string `json:"created_at"`
        *Alias
    }{
        CreatedAt: d.CreatedAt.Format(time.RFC3339),
        Alias:     (*Alias)(&d),
    })
}
```

---

## HTTP Handler Pattern

### Standard Handler Structure

```go
func (h *Handler) GetDevelopers(w http.ResponseWriter, r *http.Request) {
    // 1. Parse and validate inputs
    page := parseIntParam(r, "page", 1)
    pageSize := parseIntParam(r, "pageSize", 100)

    if pageSize > 1000 {
        h.writeError(w, http.StatusBadRequest, "INVALID_PARAMETER", "pageSize max is 1000")
        return
    }

    // 2. Call business logic
    developers, total, err := h.store.ListDevelopers(r.Context(), page, pageSize)
    if err != nil {
        h.writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
        return
    }

    // 3. Build response
    response := map[string]interface{}{
        "data": developers,
        "pagination": map[string]interface{}{
            "page":     page,
            "pageSize": pageSize,
            "total":    total,
        },
    }

    // 4. Write response
    h.writeJSON(w, http.StatusOK, response)
}
```

### Response Helpers

```go
func (h *Handler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}

func (h *Handler) writeError(w http.ResponseWriter, status int, code, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "error": map[string]string{
            "code":    code,
            "message": message,
        },
    })
}
```

---

## Common Anti-Patterns to Avoid

### 1. Don't Use init() for Complex Logic

```go
// ✗ BAD
var config *Config

func init() {
    var err error
    config, err = LoadConfig("config.json")
    if err != nil {
        panic(err)  // No way to handle errors in init
    }
}

// ✓ GOOD
func main() {
    config, err := LoadConfig("config.json")
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    // ...
}
```

### 2. Don't Ignore Errors

```go
// ✗ BAD
data, _ := json.Marshal(obj)

// ✓ GOOD
data, err := json.Marshal(obj)
if err != nil {
    return fmt.Errorf("marshal failed: %w", err)
}
```

### 3. Don't Use Naked Returns

```go
// ✗ BAD
func loadConfig(path string) (config *Config, err error) {
    // ... some code ...
    return  // Unclear what's being returned
}

// ✓ GOOD
func loadConfig(path string) (*Config, error) {
    // ... some code ...
    return config, err  // Explicit
}
```

---

## Code Review Checklist

Before submitting code for review:

- [ ] All errors are handled (no `_` ignores)
- [ ] All exported functions have comments
- [ ] All tests pass (`go test ./...`)
- [ ] Code is formatted (`gofmt -s -w .`)
- [ ] Linter passes (`golangci-lint run`)
- [ ] No panics in library code
- [ ] Concurrent access uses proper synchronization
- [ ] HTTP handlers return proper status codes
- [ ] JSON responses match Cursor API format
