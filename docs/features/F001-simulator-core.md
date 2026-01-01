# Feature F001: Cursor API Simulator Core

**Feature ID:** F001  
**Service:** cursor-sim  
**Priority:** P0 (Critical Path)  
**Status:** Specification Complete

---

## 1. Overview

The Cursor API Simulator provides a mock implementation of the Cursor Business Admin API that generates statistically realistic usage telemetry. This feature encompasses the core simulation engine, developer profile generation, and event streaming capabilities.

### 1.1 Business Value

Engineering teams need realistic test data to develop and validate analytics tooling before gaining access to production Cursor Business API credentials. The simulator eliminates this dependency by providing a fully functional mock that produces data indistinguishable from real usage patterns.

### 1.2 Success Criteria

The simulator is considered successful when it can generate data for 100+ developers over 90 days that passes statistical tests for realistic distribution, when the API response format exactly matches documented Cursor Business API schemas, and when the system can handle concurrent requests without degradation.

---

## 2. Functional Requirements

### 2.1 CLI Configuration (FR-SIM-001)

The simulator must accept command-line flags to configure its behavior at startup. These flags allow users to customize the simulation without modifying code.

**Required Flags:**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--port` | int | 8080 | HTTP server port |
| `--developers` | int | 50 | Number of simulated developers |
| `--velocity` | string | "medium" | Event rate: low, medium, high |
| `--fluctuation` | float | 0.2 | Variance factor (0.0-1.0) |
| `--seed` | int | (random) | Random seed for reproducibility |
| `--history` | int | 30 | Days of historical data to generate |

**Acceptance Criteria:**
- AC1: Running `./cursor-sim --help` displays all flags with descriptions
- AC2: Invalid flag values produce clear error messages and exit code 1
- AC3: Using the same `--seed` value produces identical data across runs
- AC4: The simulator starts within 5 seconds for up to 100 developers

### 2.2 Developer Profile Generation (FR-SIM-002)

The simulator must generate diverse developer profiles that model real-world team compositions. Each profile should have characteristics that influence their usage patterns.

**Profile Schema:**

```go
type Developer struct {
    ID             uuid.UUID `json:"id"`
    Email          string    `json:"email"`
    Name           string    `json:"name"`
    Team           string    `json:"team"`
    Role           string    `json:"role"`
    JoinedAt       time.Time `json:"joinedAt"`
    // Internal simulation parameters (not exposed in API)
    AcceptanceBase float64   `json:"-"`
    ActivityMult   float64   `json:"-"`
    WorkStartHour  int       `json:"-"`
    WorkEndHour    int       `json:"-"`
}
```

**Team Distribution:**
- Backend: 25%
- Frontend: 25%
- Platform: 20%
- Mobile: 15%
- Data: 15%

**Role Distribution:**
- Senior: 20% (higher acceptance rates, more efficient usage)
- Mid: 50% (moderate rates, typical patterns)
- Junior: 30% (lower acceptance rates, more chat dependency)

**Acceptance Criteria:**
- AC1: Developer emails are unique and follow `{firstname}.{lastname}@example.com` format
- AC2: Team distribution matches specified percentages within ±5%
- AC3: Role distribution matches specified percentages within ±5%
- AC4: Senior developers have acceptance base rates of 85-95%
- AC5: Mid developers have acceptance base rates of 65-80%
- AC6: Junior developers have acceptance base rates of 40-60%

### 2.3 Event Generation Engine (FR-SIM-003)

The simulator must generate usage events following a Poisson distribution to model realistic request patterns. Events should respect working hours and show temporal patterns.

**Event Types:**

| Event Type | Description | Frequency Weight |
|------------|-------------|-----------------|
| `cpp_suggestion_shown` | Code completion displayed | 40% |
| `cpp_suggestion_accepted` | Completion accepted via Tab | Derived from shown × acceptance rate |
| `chat_message` | Chat conversation message | 25% |
| `cmd_k_prompt` | Inline edit prompt | 15% |
| `composer_request` | Composer/multi-file edit | 10% |
| `agent_request` | Background agent invocation | 5% |
| `bugbot_usage` | Automated bug detection | 5% |

**Velocity Settings:**
- Low: λ = 10 events/hour
- Medium: λ = 50 events/hour
- High: λ = 100 events/hour

**Temporal Patterns:**
- Working hours: Events concentrated between developer's work start and end hours
- Weekday bias: 20% fewer events on Fridays, 80% fewer on weekends
- Sprint effect: 20% more events during "sprint weeks" (weeks 1, 3, 5, 7... of the month)

**Acceptance Criteria:**
- AC1: Event counts per hour follow Poisson distribution within statistical bounds
- AC2: The ratio of `cpp_suggestion_accepted` to `cpp_suggestion_shown` matches developer's acceptance rate ±5%
- AC3: 90% of events fall within defined working hours
- AC4: Weekend event counts are significantly lower than weekday counts
- AC5: Fluctuation parameter introduces visible variance between developers with similar roles

### 2.4 REST API Endpoints (FR-SIM-004)

The simulator must expose REST endpoints that match Cursor Business API conventions. Response formats must be compatible with production API schemas.

**Endpoints:**

#### GET /v1/org/users

Lists all simulated developers.

**Response:**
```json
{
  "users": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "alice.smith@example.com",
      "name": "Alice Smith",
      "team": "Backend",
      "role": "Senior"
    }
  ],
  "total": 50
}
```

**Query Parameters:**
- `team` (string): Filter by team name
- `page` (int): Page number (1-indexed)
- `pageSize` (int): Results per page (default: 100, max: 1000)

#### GET /v1/stats/activity

Returns aggregated daily usage statistics.

**Query Parameters:**
- `from` (ISO 8601): Start date (required)
- `to` (ISO 8601): End date (required)
- `email` (string): Filter by developer email
- `team` (string): Filter by team
- `page` (int): Page number (1-indexed)
- `pageSize` (int): Results per page (default: 100, max: 1000)

**Response:**
```json
{
  "data": [
    {
      "email": "alice.smith@example.com",
      "date": "2025-01-01",
      "isActive": true,
      "totalTabsShown": 142,
      "totalTabsAccepted": 98,
      "totalLinesAdded": 487,
      "totalLinesDeleted": 123,
      "acceptedLinesAdded": 312,
      "composerRequests": 8,
      "chatRequests": 23,
      "agentRequests": 2,
      "cmdkUsages": 15,
      "mostUsedModel": "claude-sonnet-4"
    }
  ],
  "pagination": {
    "page": 1,
    "pageSize": 100,
    "totalPages": 3,
    "totalRecords": 250,
    "hasNextPage": true,
    "hasPreviousPage": false
  }
}
```

#### GET /v1/events

Returns raw event stream (for debugging/testing).

**Query Parameters:**
- `from` (ISO 8601): Start timestamp
- `to` (ISO 8601): End timestamp
- `type` (string): Filter by event type
- `developerId` (UUID): Filter by developer
- `limit` (int): Maximum events (default: 1000, max: 10000)

**Response:**
```json
{
  "events": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440001",
      "developerId": "550e8400-e29b-41d4-a716-446655440000",
      "type": "cpp_suggestion_shown",
      "timestamp": "2025-01-01T10:23:45Z",
      "metadata": {
        "model": "claude-sonnet-4",
        "linesShown": 3,
        "language": "typescript"
      }
    }
  ],
  "total": 5000,
  "returned": 1000
}
```

#### GET /health

Health check endpoint for Docker orchestration.

**Response:**
```json
{
  "status": "healthy",
  "version": "1.0.0",
  "developers": 50,
  "eventsGenerated": 125000,
  "uptime": "2h34m12s"
}
```

**Acceptance Criteria:**
- AC1: All endpoints return valid JSON with correct Content-Type header
- AC2: Invalid date ranges return 400 Bad Request with descriptive error
- AC3: Pagination works correctly with consistent ordering across pages
- AC4: Health endpoint returns within 100ms
- AC5: API handles 100 concurrent requests without errors

### 2.5 In-Memory Storage (FR-SIM-005)

The simulator must use in-memory SQLite for fast data access and SQL query capabilities without external dependencies.

**Requirements:**
- Generate all data on startup and store in memory
- Support efficient time-range queries on events
- Index appropriately for pagination and filtering
- Clear state on restart (volatile by design)

**Acceptance Criteria:**
- AC1: Startup time remains under 10 seconds for 100 developers with 90 days of history
- AC2: Query response time under 100ms for typical API requests
- AC3: Memory usage stays under 500MB for 100 developers with 90 days of history

---

## 3. Non-Functional Requirements

### 3.1 Performance

- Startup: Generate 90 days of data for 100 developers in under 10 seconds
- Throughput: Handle 100 requests/second for the activity endpoint
- Latency: P99 response time under 200ms for all endpoints

### 3.2 Reliability

- The simulator should run continuously without memory leaks
- Graceful shutdown on SIGTERM/SIGINT
- Structured logging for debugging

### 3.3 Maintainability

- Code coverage minimum: 80%
- All public functions documented
- Configuration externalized (no hardcoded values)

---

## 4. Technical Design Notes

### 4.1 Poisson Distribution Implementation

The Poisson distribution models the number of events occurring in a fixed interval. For a given λ (lambda, the expected number of events), we generate inter-arrival times using the exponential distribution.

```go
// Generate Poisson-distributed count for an interval
func (g *PoissonGenerator) CountForInterval(lambda float64, intervalHours float64) int {
    adjustedLambda := lambda * intervalHours
    // Using inverse transform sampling
    L := math.Exp(-adjustedLambda)
    k := 0
    p := 1.0
    for p > L {
        k++
        p *= g.rand.Float64()
    }
    return k - 1
}
```

### 4.2 Deterministic Seeding

All random number generators must be seeded from the CLI `--seed` parameter to ensure reproducible simulations. This is critical for testing and debugging.

```go
func NewSimulator(config Config) *Simulator {
    var seed int64
    if config.Seed != 0 {
        seed = config.Seed
    } else {
        seed = time.Now().UnixNano()
    }
    return &Simulator{
        rand: rand.New(rand.NewSource(seed)),
        // ...
    }
}
```

### 4.3 Working Hours Modeling

Events should cluster during working hours. Each developer has randomized work patterns:

```go
type WorkPattern struct {
    StartHour    int     // 7-10
    EndHour      int     // 16-20
    LunchHour    int     // 11-13
    LunchProb    float64 // Probability of reduced activity during lunch
    FridayFactor float64 // 0.7-0.9
}
```

---

## 5. Dependencies

### 5.1 External Libraries

| Library | Version | Purpose |
|---------|---------|---------|
| chi | v5.0+ | HTTP router |
| go-sqlite3 | v1.14+ | In-memory database |
| uuid | v1.3+ | UUID generation |
| cobra | v1.8+ | CLI framework |
| testify | v1.8+ | Testing assertions |
| faker | v4.0+ | Realistic data generation |

### 5.2 Internal Dependencies

None - this service is the data source and has no upstream dependencies.

---

## 6. Test Cases

### 6.1 Unit Tests

| ID | Test Case | Expected Result |
|----|-----------|-----------------|
| UT-SIM-001 | Poisson generator produces correct mean | Mean of 10000 samples within ±2% of λ |
| UT-SIM-002 | Developer profiles have unique emails | No duplicate emails in generated set |
| UT-SIM-003 | Role distribution matches specification | Counts within ±5% of expected |
| UT-SIM-004 | Acceptance rates by role are in expected ranges | Senior: 85-95%, Mid: 65-80%, Junior: 40-60% |
| UT-SIM-005 | Same seed produces same developers | Developer list identical for same seed |
| UT-SIM-006 | Event timestamps respect working hours | 90%+ events within work hours |

### 6.2 Integration Tests

| ID | Test Case | Expected Result |
|----|-----------|-----------------|
| IT-SIM-001 | GET /v1/org/users returns all developers | Status 200, correct count |
| IT-SIM-002 | GET /v1/org/users with team filter | Only specified team returned |
| IT-SIM-003 | GET /v1/stats/activity with date range | Events within range only |
| IT-SIM-004 | Pagination returns consistent results | Page N+1 continues where N ended |
| IT-SIM-005 | Invalid date range returns 400 | Error message explains issue |
| IT-SIM-006 | Health endpoint reflects current state | Correct developer and event counts |

### 6.3 Performance Tests

| ID | Test Case | Expected Result |
|----|-----------|-----------------|
| PT-SIM-001 | Startup with 100 developers, 90 days | Completes in under 10 seconds |
| PT-SIM-002 | 100 concurrent requests to /v1/stats/activity | All complete successfully, P99 < 200ms |
| PT-SIM-003 | Memory usage after 1 hour | No significant growth beyond initial |

---

## 7. Related User Stories

- [US-SIM-001](../user-stories/US-SIM-001-configure-simulation.md): Configure Simulation Parameters
- [US-SIM-002](../user-stories/US-SIM-002-generate-developers.md): Generate Realistic Developer Profiles
- [US-SIM-003](../user-stories/US-SIM-003-simulate-events.md): Simulate Usage Events
- [US-SIM-004](../user-stories/US-SIM-004-query-activity.md): Query Activity Data

---

## 8. Implementation Tasks

- [TASK-001](../tasks/TASK-001-sim-project-setup.md): Set up Go project structure
- [TASK-002](../tasks/TASK-002-sim-cli.md): Implement CLI with Cobra
- [TASK-003](../tasks/TASK-003-sim-developer-gen.md): Implement developer profile generation
- [TASK-004](../tasks/TASK-004-sim-poisson.md): Implement Poisson event generator
- [TASK-005](../tasks/TASK-005-sim-sqlite.md): Set up in-memory SQLite storage
- [TASK-006](../tasks/TASK-006-sim-api-users.md): Implement /v1/org/users endpoint
- [TASK-007](../tasks/TASK-007-sim-api-activity.md): Implement /v1/stats/activity endpoint
- [TASK-008](../tasks/TASK-008-sim-api-events.md): Implement /v1/events endpoint
- [TASK-009](../tasks/TASK-009-sim-docker.md): Create Dockerfile and health check
