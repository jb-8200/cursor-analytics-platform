# Development Process (TDD Workflow)

**Trigger**: Implementing features, writing code, fixing bugs

---

## The TDD Cycle

### RED Phase: Write Failing Test

**Before writing any implementation code:**

1. **Identify acceptance criterion** from user story
2. **Translate to test case**:
   - Given → Arrange (setup)
   - When → Act (execute)
   - Then → Assert (verify)
3. **Run test** - confirm it fails
4. **Verify failure reason** - should fail because code doesn't exist

Example:

```go
func TestAcceptanceRate_CalculatesCorrectly(t *testing.T) {
    // Arrange (Given developer with 100 suggestions, 75 accepted)
    dev := NewDeveloper("dev-1")
    dev.SuggestionsShown = 100
    dev.SuggestionsAccepted = 75

    // Act (When I calculate acceptance rate)
    rate := dev.AcceptanceRate()

    // Assert (Then result is 75.0%)
    assert.Equal(t, 75.0, rate)
}
```

Run: `go test ./... -run TestAcceptanceRate` → FAIL (method doesn't exist)

### GREEN Phase: Minimal Implementation

**Write just enough code to pass the test:**

1. **Implement minimally** - no extra features
2. **Run test** - confirm it passes
3. **No refactoring yet** - just make it work

```go
func (d *Developer) AcceptanceRate() float64 {
    if d.SuggestionsShown == 0 {
        return 0
    }
    return float64(d.SuggestionsAccepted) / float64(d.SuggestionsShown) * 100
}
```

Run: `go test ./... -run TestAcceptanceRate` → PASS

### REFACTOR Phase: Clean Up

**Improve code while keeping tests green:**

1. **Extract duplication**
2. **Improve naming**
3. **Simplify logic**
4. **Run tests after each change**

DO NOT:
- Add new functionality
- Change behavior
- Skip running tests

---

## Test Patterns by Type

### Unit Tests

Test individual functions/methods in isolation:

```go
func TestCalculator_Add(t *testing.T) {
    calc := NewCalculator()
    result := calc.Add(2, 3)
    assert.Equal(t, 5, result)
}
```

### Table-Driven Tests

Test multiple scenarios efficiently:

```go
func TestAcceptanceRate(t *testing.T) {
    tests := []struct {
        name     string
        shown    int
        accepted int
        expected float64
    }{
        {"all accepted", 100, 100, 100.0},
        {"half accepted", 100, 50, 50.0},
        {"none shown", 0, 0, 0.0},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            dev := NewDeveloper("test")
            dev.SuggestionsShown = tt.shown
            dev.SuggestionsAccepted = tt.accepted
            assert.Equal(t, tt.expected, dev.AcceptanceRate())
        })
    }
}
```

### Integration Tests

Test component interactions:

```go
func TestHandler_ReturnsCorrectFormat(t *testing.T) {
    store := storage.NewMemory()
    store.AddCommit(testCommit)

    handler := NewHandler(store)
    req := httptest.NewRequest("GET", "/commits", nil)
    rec := httptest.NewRecorder()

    handler.ServeHTTP(rec, req)

    assert.Equal(t, 200, rec.Code)
    assert.Contains(t, rec.Body.String(), "data")
}
```

### E2E Tests

Test full user flows:

```go
func TestE2E_CommitWorkflow(t *testing.T) {
    // Start server
    srv := startTestServer(t)
    defer srv.Close()

    // Make request
    resp, _ := http.Get(srv.URL + "/commits")

    // Verify response
    assert.Equal(t, 200, resp.StatusCode)
}
```

---

## Running Tests

### Go (cursor-sim)

```bash
# All tests
go test ./...

# Specific package
go test ./internal/models

# Specific test
go test ./... -run TestAcceptanceRate

# With coverage
go test ./... -cover

# Verbose
go test ./... -v
```

### TypeScript (future services)

```bash
# All tests
npm test

# Watch mode
npm test -- --watch

# With coverage
npm test -- --coverage
```

---

## When Tests Fail

1. **Read the error message carefully**
2. **Understand expected vs actual**
3. **Fix the issue** (code or test?)
4. **Run tests again**
5. **Don't proceed until green**

Common issues:
- Test has wrong expectation → Fix test
- Implementation has bug → Fix code
- Missing setup/teardown → Add fixtures
- Flaky test → Add synchronization

---

## Coverage Requirements

| Category | Minimum | Target |
|----------|---------|--------|
| Core logic | 90% | 95% |
| API handlers | 80% | 90% |
| Utilities | 70% | 80% |
| Overall | 80% | 85% |

Check coverage:
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## After TDD Cycle

**Every completed task requires:**

1. ✅ All tests pass
2. ✅ Coverage meets threshold
3. ✅ Git commit
4. ✅ task.md updated
5. ✅ DEVELOPMENT.md updated

See `sdd-checklist` skill for enforcement.
