# Feature 3: In-Memory Database - Development Tasks

## Feature Overview
Implement a thread-safe in-memory database that stores and efficiently retrieves all generated data. The database must support concurrent access, indexing, and range queries while maintaining data consistency.

---

## Task 3.1: Database Interface Definition

**User Story**:
> As a developer, I want a well-defined database interface, so that I can swap implementations if needed and test easily with mocks.

**Description**:
Define the complete database interface that all data access must go through.

**Acceptance Criteria**:
- [ ] `Store` interface defined with all required methods
- [ ] Methods for CRUD operations: Create, Read, Update, Delete
- [ ] Methods for querying: GetByID, GetByRange, GetAll with filters
- [ ] Methods for aggregation: CountByRegion, CountByTeam, etc.
- [ ] Methods return errors for invalid operations
- [ ] Thread-safe method signatures (no mutable parameters)
- [ ] Transactions/batch operations for consistency
- [ ] Clear documentation of expected behavior

**Success Criteria**:
- Interface is complete and logical
- All data types support required operations
- Error cases are well-defined
- Mock implementations are straightforward

**Files to Create**:
- `internal/db/store.go` - Store interface definition

**Example Interface**:
```go
type Store interface {
    // Developers
    SaveDeveloper(ctx context.Context, dev *Developer) error
    GetDeveloper(ctx context.Context, id string) (*Developer, error)
    ListDevelopers(ctx context.Context, filters DeveloperFilters, limit int, offset int) ([]*Developer, error)
    CountDevelopers(ctx context.Context) (int64, error)

    // Commits
    SaveCommit(ctx context.Context, commit *Commit) error
    GetCommitsInRange(ctx context.Context, startDate, endDate time.Time) ([]*Commit, error)

    // ... other methods
}
```

---

## Task 3.2: Memory Store Implementation

**User Story**:
> As a developer, I want a concrete in-memory store implementation using sync.Map, so that I can store and retrieve data efficiently with thread-safe access.

**Description**:
Implement the Store interface using Go's sync.Map for thread-safe storage.

**Acceptance Criteria**:
- [ ] `MemoryStore` struct defined with sync.Map instances for each data type
- [ ] All Store interface methods implemented
- [ ] Developers storage with primary key (ID) and secondary indices (Email, Region, Team)
- [ ] Commits storage with composite indices (DeveloperID + Date)
- [ ] Changes storage with composite indices
- [ ] Metrics storage with date-based indexing
- [ ] No locks causing performance issues (minimal contention)
- [ ] Context support for cancellation
- [ ] Proper error handling (not found, invalid operations)

**Success Criteria**:
- All CRUD operations work correctly
- Concurrent read/writes don't cause panics
- Data consistency maintained
- Queries return correct results

**Files to Create**:
- `internal/db/memory_store.go` - MemoryStore implementation
- `internal/db/memory_store_test.go` - Store tests

**Example Test**:
```go
func TestMemoryStoreCreate(t *testing.T) {
    store := NewMemoryStore()
    dev := NewDeveloper("test@example.com", "Test")
    err := store.SaveDeveloper(context.Background(), dev)
    assert.NoError(t, err)
}

func TestMemoryStoreRead(t *testing.T) {
    store := NewMemoryStore()
    dev := NewDeveloper("test@example.com", "Test")
    store.SaveDeveloper(context.Background(), dev)

    retrieved, err := store.GetDeveloper(context.Background(), dev.ID)
    assert.NoError(t, err)
    assert.Equal(t, dev.ID, retrieved.ID)
}

func TestMemoryStoreConcurrentWrites(t *testing.T) {
    store := NewMemoryStore()
    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(idx int) {
            defer wg.Done()
            dev := NewDeveloper(fmt.Sprintf("dev%d@example.com", idx), "Dev")
            store.SaveDeveloper(context.Background(), dev)
        }(i)
    }
    wg.Wait()
    count, _ := store.CountDevelopers(context.Background())
    assert.Equal(t, count, int64(100))
}
```

---

## Task 3.3: Indexing & Query Optimization

**User Story**:
> As a developer, I want efficient queries by region, team, and date range, so that the API can return results quickly even with large datasets.

**Description**:
Implement secondary indices for common query patterns to optimize performance.

**Acceptance Criteria**:
- [ ] Index by DeveloperID for quick lookups
- [ ] Index by Region (US, EU, APAC) for regional queries
- [ ] Index by Division for division-based queries
- [ ] Index by Team for team-based queries
- [ ] Date-range index for time-based queries
- [ ] Composite indices for multiple filters
- [ ] Indices automatically maintained on writes
- [ ] Query performance remains O(1) or O(log n) for indexed fields
- [ ] Unindexed queries still work but may be slower

**Success Criteria**:
- All indices maintain consistency with data
- Queries use appropriate indices
- Range queries return correct results
- Index updates don't introduce race conditions

**Files to Create**:
- `internal/db/indexer.go` - Index management
- `internal/db/indexer_test.go` - Index tests

**Example Test**:
```go
func TestIndexByRegion(t *testing.T) {
    store := NewMemoryStore()

    // Add developers from different regions
    devUS := NewDeveloper("us@example.com", "US Dev")
    devUS.Region = "US"
    devEU := NewDeveloper("eu@example.com", "EU Dev")
    devEU.Region = "EU"

    store.SaveDeveloper(context.Background(), devUS)
    store.SaveDeveloper(context.Background(), devEU)

    // Query by region
    usDevs, _ := store.ListDevelopers(context.Background(),
        DeveloperFilters{Region: "US"}, 100, 0)
    assert.Len(t, usDevs, 1)
    assert.Equal(t, usDevs[0].ID, devUS.ID)
}

func TestDateRangeQuery(t *testing.T) {
    store := NewMemoryStore()

    // Add commits on different dates
    commit1 := NewCommit("sha1", "dev-id", "repo")
    commit1.Timestamp = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

    commit2 := NewCommit("sha2", "dev-id", "repo")
    commit2.Timestamp = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

    store.SaveCommit(context.Background(), commit1)
    store.SaveCommit(context.Background(), commit2)

    // Query range
    start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
    end := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
    commits, _ := store.GetCommitsInRange(context.Background(), start, end)
    assert.Len(t, commits, 1)
}
```

---

## Task 3.4: Aggregation Functions

**User Story**:
> As an analyst, I want aggregated metrics computed on-demand, so that I can get team-wide statistics without manual calculation.

**Description**:
Implement methods to compute aggregated statistics from stored data.

**Acceptance Criteria**:
- [ ] `GetDailyMetrics(date time.Time)` computes metrics for a day
- [ ] `GetDeveloperMetrics(devID string, date time.Time)` for per-dev metrics
- [ ] `GetRegionalStats(region string, startDate, endDate time.Time)` for regional analysis
- [ ] `GetTeamStats(team string, startDate, endDate time.Time)` for team analysis
- [ ] Leaderboard functions ranking by various metrics
- [ ] All aggregations are computed from stored data (no stored aggregates)
- [ ] Aggregations are reasonably fast even with large datasets
- [ ] Caching layer optional for frequently accessed aggregations

**Success Criteria**:
- Aggregations return correct results
- Performance is acceptable (< 1s for typical queries)
- Different aggregation methods produce consistent results

**Files to Create**:
- `internal/db/aggregations.go` - Aggregation logic
- `internal/db/aggregations_test.go` - Aggregation tests

**Example Test**:
```go
func TestComputeDailyMetrics(t *testing.T) {
    store := NewMemoryStore()

    // Add commits for a day
    date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
    for i := 0; i < 10; i++ {
        commit := NewCommit(fmt.Sprintf("sha%d", i), fmt.Sprintf("dev%d", i), "repo")
        commit.Timestamp = date.Add(time.Duration(i) * time.Hour)
        store.SaveCommit(context.Background(), commit)
    }

    metrics := store.GetDailyMetrics(context.Background(), date)
    assert.Equal(t, metrics.Date, date)
    assert.Equal(t, metrics.Commits, 10)
}

func TestLeaderboard(t *testing.T) {
    store := NewMemoryStore()

    // Add varied data
    // ...

    leaderboard, _ := store.GetLeaderboard(context.Background(),
        "agent_edits", 10)
    assert.True(t, len(leaderboard) > 0)
    // Verify sorted in descending order
    for i := 1; i < len(leaderboard); i++ {
        assert.True(t, leaderboard[i-1].Metric >= leaderboard[i].Metric)
    }
}
```

---

## Task 3.5: Batch Operations & Transactions

**User Story**:
> As a developer, I want to perform batch operations atomically, so that data consistency is maintained during bulk inserts.

**Description**:
Implement batch and transaction-like operations for consistent multi-record operations.

**Acceptance Criteria**:
- [ ] `BatchInsertDevelopers(ctx context.Context, devs []*Developer)` inserts multiple
- [ ] `BatchInsertCommits(ctx context.Context, commits []*Commit)` inserts multiple
- [ ] All-or-nothing semantics (if any fails, all rollback)
- [ ] Atomic from external view (no partial data visible)
- [ ] Reasonable performance for large batches (10k+ records)
- [ ] Error handling reports which records failed

**Success Criteria**:
- Batch operations complete successfully
- All records inserted or none are
- Performance scales reasonably with batch size

**Files to Use**:
- `internal/db/memory_store.go` - Add batch methods

**Example Test**:
```go
func TestBatchInsertDevelopers(t *testing.T) {
    store := NewMemoryStore()

    devs := make([]*Developer, 1000)
    for i := 0; i < 1000; i++ {
        devs[i] = NewDeveloper(fmt.Sprintf("dev%d@example.com", i), "Dev")
    }

    err := store.BatchInsertDevelopers(context.Background(), devs)
    assert.NoError(t, err)

    count, _ := store.CountDevelopers(context.Background())
    assert.Equal(t, count, int64(1000))
}
```

---

## Task 3.6: Data Consistency & Recovery

**User Story**:
> As an operator, I want data consistency guarantees, so that I can trust the in-memory database.

**Description**:
Implement checks and utilities to maintain and verify data consistency.

**Acceptance Criteria**:
- [ ] `VerifyConsistency()` function checks all data for consistency
- [ ] Detects orphaned records (commit without developer, etc.)
- [ ] Detects duplicate IDs
- [ ] Reports consistency violations clearly
- [ ] Can repair simple inconsistencies (optional)
- [ ] Tests verify consistency after operations

**Success Criteria**:
- Consistency checks detect problems
- Reports are clear and actionable
- Repair operations work correctly

**Files to Create**:
- `internal/db/consistency.go` - Consistency checks
- `internal/db/consistency_test.go` - Consistency tests

---

## Task 3.7: Memory Management & Limits

**User Story**:
> As an operator, I want to understand memory usage, so that I can plan for large simulations.

**Description**:
Implement monitoring and optional limits on in-memory data.

**Acceptance Criteria**:
- [ ] `GetMemoryStats()` returns rough memory usage estimate
- [ ] Tracks count of each data type (developers, commits, changes)
- [ ] Optional max size limits to prevent unbounded growth
- [ ] Clear error when limit exceeded
- [ ] Statistics accessible via API

**Success Criteria**:
- Memory stats are available
- Limits prevent runaway growth if configured
- Estimates are reasonably accurate

**Files to Use**:
- `internal/db/memory_store.go` - Add memory tracking

---

## Task 3.8: Snapshot & Reload

**User Story**:
> As a developer, I want to snapshot the in-memory database state, so that I can restart with previous data.

**Description**:
Implement methods to save and restore database snapshots.

**Acceptance Criteria**:
- [ ] `CreateSnapshot(ctx context.Context)` creates snapshot of all data
- [ ] Snapshot is atomic (represents consistent point in time)
- [ ] `RestoreSnapshot(ctx context.Context, snapshot *Snapshot)` loads snapshot
- [ ] Snapshot can be serialized and saved to disk
- [ ] Restore validates data integrity

**Success Criteria**:
- Snapshots capture complete state
- Restore brings database to saved state
- Serialization/deserialization works correctly

**Files to Create**:
- `internal/db/snapshot.go` - Snapshot logic
- `internal/db/snapshot_test.go` - Snapshot tests

---

## Feature 3 Integration Test

**File**: `tests/integration_db_test.go`

```go
func TestDatabaseWorkflow(t *testing.T) {
    store := NewMemoryStore()

    // Add developers
    devs := make([]*Developer, 100)
    for i := 0; i < 100; i++ {
        dev := NewDeveloper(fmt.Sprintf("dev%d@example.com", i), "Dev")
        dev.Region = []string{"US", "EU", "APAC"}[i%3]
        devs[i] = dev
    }
    store.BatchInsertDevelopers(context.Background(), devs)

    // Add commits
    for i := 0; i < 500; i++ {
        commit := NewCommit(fmt.Sprintf("sha%d", i), devs[i%100].ID, "repo")
        store.SaveCommit(context.Background(), commit)
    }

    // Verify consistency
    errs := store.VerifyConsistency(context.Background())
    assert.Empty(t, errs)

    // Query various ways
    usDevs, _ := store.ListDevelopers(context.Background(),
        DeveloperFilters{Region: "US"}, 100, 0)
    assert.True(t, len(usDevs) > 0)

    // Get metrics
    metrics := store.GetDailyMetrics(context.Background(), time.Now())
    assert.Equal(t, metrics.DAU, 100)
}
```

---

## Feature 3 Completion Checklist

- [ ] Task 3.1: Store interface designed and documented
- [ ] Task 3.2: MemoryStore fully implemented
- [ ] Task 3.3: Indices working and queries fast
- [ ] Task 3.4: Aggregation functions implemented
- [ ] Task 3.5: Batch operations working
- [ ] Task 3.6: Consistency checks in place
- [ ] Task 3.7: Memory tracking implemented
- [ ] Task 3.8: Snapshots working
- [ ] Integration test passes
- [ ] Test coverage ≥85%
- [ ] Performance acceptable for 100k+ records
- [ ] No race conditions under concurrent access

