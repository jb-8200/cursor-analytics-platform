# Cursor-Sim Service - Detailed Design Document

## Executive Summary

**cursor-sim** is a sophisticated Go-based simulator that generates synthetic Cursor analytics and code tracking data. It mimics real developer behavior across multiple regions, divisions, groups, and teams while providing realistic metrics aligned with Cursor's official Analytics API and AI Code Tracking API.

The simulator accepts JSON-based configuration, manages in-memory data persistence, exposes REST and GraphQL interfaces, and provides interactive CLI controls for data generation management.

---

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Core Features & Breakdown](#core-features--breakdown)
3. [Data Models](#data-models)
4. [Configuration Specification](#configuration-specification)
5. [API Specifications](#api-specifications)
6. [CLI Interface & Controls](#cli-interface--controls)
7. [Development Task Breakdown](#development-task-breakdown)
8. [Testing Strategy](#testing-strategy)

---

## Architecture Overview

### High-Level System Design

```
┌──────────────────────────────────────────────────────────────┐
│                    cursor-sim Service (Go)                   │
├──────────────────────────────────────────────────────────────┤
│                                                                │
│  ┌─────────────────┐      ┌──────────────────┐              │
│  │   JSON Config   │──►  │  CLI Parser       │              │
│  │   (Input)       │      │  & Validator     │              │
│  └─────────────────┘      └────────┬─────────┘              │
│                                     │                        │
│  ┌─────────────────────────────────▼────────────────────┐  │
│  │           Data Generator Engine                       │  │
│  │  ┌──────────────────────────────────────────────┐   │  │
│  │  │ • Developer Profiles (Regions, Teams, Skills)│   │  │
│  │  │ • Commits/Changes (TAB & COMPOSER)           │   │  │
│  │  │ • Analytics Events (DAU, Model Usage, etc.)  │   │  │
│  │  │ • Poisson-Distributed Event Timing           │   │  │
│  │  └──────────────────────────────────────────────┘   │  │
│  └──────────────────┬───────────────────────────────────┘  │
│                     │                                        │
│  ┌──────────────────▼──────────────────┐                   │
│  │    In-Memory Database               │                   │
│  │  ┌────────────────────────────────┐ │                   │
│  │  │ • Developers                   │ │                   │
│  │  │ • Commits & Changes            │ │                   │
│  │  │ • Analytics Metrics            │ │                   │
│  │  │ • Aggregated Statistics        │ │                   │
│  │  └────────────────────────────────┘ │                   │
│  └──────────────────┬───────────────────┘                   │
│                     │                                        │
│  ┌──────────────────┴──────────────────┐                   │
│  │                                      │                   │
│  ▼                                      ▼                   │
│ ┌─────────────────┐          ┌──────────────────┐          │
│ │  REST API       │          │  GraphQL API     │          │
│ │  (Cursor-like)  │          │  (Flexible Query)│          │
│ └─────────────────┘          └──────────────────┘          │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │           Interactive CLI Dashboard                   │  │
│  │  • Real-time Generation Stats                        │  │
│  │  • Signal Handlers (Ctrl+S, Ctrl+E, Ctrl+C)         │  │
│  │  • Status Display (Developers, PRs, Concurrency)    │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

### Component Responsibilities

| Component | Responsibility |
|-----------|-----------------|
| **Config Parser** | Validates and parses JSON input configuration |
| **Data Generator** | Creates realistic synthetic developer profiles and events |
| **In-Memory DB** | Stores and retrieves all generated data with thread-safe access |
| **REST API Server** | Exposes endpoints compatible with Cursor API format |
| **GraphQL API** | Provides flexible querying of stored metrics |
| **CLI Controller** | Manages interactive controls and displays stats |
| **Auth Simulator** | Implements fake Basic Authentication with fixed credentials |

---

## Core Features & Breakdown

### Feature 1: Configuration & Initialization
**Purpose**: Accept and validate JSON configuration to control simulation parameters

**Components**:
- JSON schema validation
- Configuration parser
- Parameter defaults and bounds checking

**Key Parameters**:
```json
{
  "auth": {
    "api_key": "string",
    "api_secret": "string",
    "team_id": "string"
  },
  "simulation": {
    "developers": 100,
    "regions": ["US", "EU", "APAC"],
    "divisions": ["AGS", "AT", "ST"],
    "groups": ["TMOBILE", "ATANT"],
    "teams": ["dev", "support"],
    "pr_velocity": "high|medium|low",
    "volatility": 0.0-1.0,
    "break_condition": {
      "type": "pr_count|time_duration|none",
      "value": 100000
    }
  },
  "api": {
    "port": 8080,
    "enable_graphql": true
  }
}
```

---

### Feature 2: Developer Profile Generation
**Purpose**: Create realistic synthetic developer profiles with regional and organizational context

**Characteristics**:
- Unique email addresses per developer
- Regional assignment (US/EU/APAC with realistic distributions)
- Division assignment (AGS/AT/ST)
- Group assignment (TMOBILE/ATANT)
- Team assignment (dev/support)
- Random skill sets and expertise levels
- Realistic client versions (IDE versions)
- Activation timestamps

**Data Model**:
```go
type Developer struct {
    ID               string
    Email            string
    Name             string
    Region           string    // US, EU, APAC
    Division         string    // AGS, AT, ST
    Group            string    // TMOBILE, ATANT
    Team             string    // dev, support
    Skills           []string
    ClientVersion    string
    CreatedAt        time.Time
    LastActiveAt     time.Time
}
```

---

### Feature 3: Commit & Change Data Generation
**Purpose**: Simulate realistic code commits with AI-assisted changes

**Models**:
```go
type Commit struct {
    Hash              string
    Timestamp         time.Time
    Message           string
    DeveloperID       string
    Repository        string
    Branch            string
    LinesDelta        int
    LinesFromTAB      int
    LinesFromComposer int
    LinesNonAI        int
    IngestionTime     time.Time
}

type Change struct {
    ChangeID          string
    CommitHash        string
    DeveloperID       string
    Timestamp         time.Time
    Source            string    // TAB or COMPOSER
    Model             string    // claude-3.5-sonnet, etc.
    FilesDelta        int
    LinesAdded        int
    LinesModified     int
    IngestionTime     time.Time
}
```

**Generation Rules**:
- Poisson-distributed event timing based on velocity
- Realistic line count deltas (5-500 lines per commit)
- TAB vs COMPOSER distribution based on developer profile
- AI model assignment (rotating through available models)
- Repository naming based on team/division context

---

### Feature 4: Analytics Metrics Generation
**Purpose**: Generate aggregated analytics aligned with Cursor Analytics API

**Metrics to Generate**:
- **Daily Active Users (DAU)**: Unique developers per day
- **Agent Edits**: Accepted AI-assisted code edits
- **Tab Usage**: TAB autocomplete metrics per developer
- **Model Usage**: Distribution of AI models used
- **Client Versions**: IDE version distribution
- **Top File Extensions**: .ts, .js, .py, .go, .jsx, .tsx, etc.
- **MCP Adoption**: Model Context Protocol tool usage
- **Commands Adoption**: Cursor command usage (e.g., cmd+k)
- **Plans Adoption**: Plan mode usage by model
- **Ask Mode Adoption**: Ask mode query patterns
- **Team Leaderboard**: Top contributors by metric

**Data Model**:
```go
type DailyMetrics struct {
    Date                 time.Time
    DAU                  int
    AgentEditsAccepted   int
    TabCompletions      int
    ComposerEdits       int
    ModelUsage          map[string]int
    TopFileExtensions   map[string]int
    CommandUsage        map[string]int
}

type DeveloperMetrics struct {
    DeveloperID         string
    Date                time.Time
    AgentEdits          int
    TabCompletions      int
    ComposerEdits       int
    ModelsUsed          map[string]int
    FileExtensions      map[string]int
}
```

---

### Feature 5: In-Memory Database
**Purpose**: Store and efficiently query all generated data with thread-safe access

**Storage Strategy**:
- Use `sync.Map` or `go-memdb` for concurrent access
- Index data by: DeveloperID, Date, Region, Division, Team
- Support range queries (date ranges)
- Maintain aggregation caches for performance

**Collections**:
- `developers`: All developer profiles
- `commits`: All generated commits
- `changes`: All generated changes
- `daily_metrics`: Aggregated daily statistics
- `developer_metrics`: Per-developer daily statistics

---

### Feature 6: REST API Server
**Purpose**: Expose endpoints compatible with Cursor's API format

**Authentication**: Basic Auth with fixed credentials from config

**Endpoints**:

#### AI Code Tracking API
```
GET /v1/analytics/ai-code/commits
  Query: startDate, endDate, page, pageSize, user
  Response: Paginated commit metrics (JSON)

GET /v1/analytics/ai-code/commits.csv
  Query: startDate, endDate, user
  Response: CSV stream of commit data

GET /v1/analytics/ai-code/changes
  Query: startDate, endDate, page, pageSize, user
  Response: Paginated change metrics (JSON)

GET /v1/analytics/ai-code/changes.csv
  Query: startDate, endDate, user
  Response: CSV stream of change data
```

#### Team Analytics API
```
GET /v1/analytics/team/agent-edits
  Query: startDate, endDate, page, pageSize
  Response: Team-wide agent edit metrics

GET /v1/analytics/team/tabs
  Query: startDate, endDate, page, pageSize
  Response: TAB usage metrics

GET /v1/analytics/team/dau
  Query: startDate, endDate
  Response: Daily active users

GET /v1/analytics/team/models
  Query: startDate, endDate
  Response: AI model usage distribution

GET /v1/analytics/team/client-versions
  Query: startDate, endDate
  Response: IDE version distribution

GET /v1/analytics/team/top-file-extensions
  Query: startDate, endDate, limit
  Response: Most-edited file types

GET /v1/analytics/team/leaderboard
  Query: startDate, endDate, metric, limit
  Response: Ranked developers by metric

GET /v1/health
  Response: Service health status
```

**Response Format**:
```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "pageSize": 100,
    "total": 5000
  },
  "params": {
    "startDate": "2024-01-01",
    "endDate": "2024-01-31"
  }
}
```

---

### Feature 7: GraphQL API
**Purpose**: Provide flexible querying with parameterized filters

**Schema Overview**:
```graphql
type Query {
  developers(
    region: [String]
    division: [String]
    team: [String]
    limit: Int
  ): [Developer]

  commits(
    startDate: DateTime
    endDate: DateTime
    developerId: String
    repository: String
  ): [Commit]

  dailyMetrics(
    startDate: DateTime
    endDate: DateTime
    region: [String]
  ): [DailyMetrics]

  developerMetrics(
    developerId: String
    startDate: DateTime
    endDate: DateTime
  ): [DeveloperMetrics]

  leaderboard(
    metric: String
    startDate: DateTime
    endDate: DateTime
    limit: Int
  ): [LeaderboardEntry]
}
```

---

### Feature 8: Interactive CLI Dashboard
**Purpose**: Provide real-time monitoring and control of data generation

**Display Elements**:
```
╔════════════════════════════════════════════════════════════╗
║         Cursor Analytics Simulator - Running                ║
╠════════════════════════════════════════════════════════════╣
║ Status: GENERATING                                          ║
║ Uptime: 00:05:32                                           ║
║                                                              ║
║ Developers:     142 total                                  ║
║   • US:         62 (44%)                                   ║
║   • EU:         54 (38%)                                   ║
║   • APAC:       26 (18%)                                   ║
║                                                              ║
║ Commits Generated:                                          ║
║   • Total:      3,847 (cumulative)                         ║
║   • Current:    12 commits/min                             ║
║   • From TAB:   2,341 (61%)                                ║
║   • From Composer: 1,506 (39%)                             ║
║                                                              ║
║ Generation Speed:                                           ║
║   • Velocity:   HIGH                                       ║
║   • Volatility: 0.35                                       ║
║                                                              ║
║ Break Condition:                                            ║
║   • Type:       PR Count                                   ║
║   • Progress:   3,847 / 100,000 (3.8%)                    ║
║                                                              ║
╠════════════════════════════════════════════════════════════╣
║ Controls: [S]oft Stop | [E]xport DB | [Q]uit               ║
║ Listening on http://localhost:8080                         ║
╚════════════════════════════════════════════════════════════╝
```

**Interactive Controls**:
- **Ctrl+S** (Soft Stop): Pause data generation, keep API active
- **Ctrl+E** (Export): Dump entire in-memory DB to JSON/binary file
- **Ctrl+C** (Quit): Graceful shutdown of entire application
- **Ctrl+Q** (Query): Show recent statistics snapshot

---

### Feature 9: Data Export & Persistence
**Purpose**: Allow dumping in-memory database for analysis and sharing

**Export Format Options**:
- **JSON**: Pretty-printed JSON file with all data
- **Binary**: Compressed binary format using gob encoding

**Export Structure**:
```json
{
  "export_timestamp": "2024-01-15T14:30:00Z",
  "simulation_config": { ... },
  "statistics": {
    "total_developers": 142,
    "total_commits": 3847,
    "total_changes": 8291,
    "generation_duration": "00:05:32"
  },
  "data": {
    "developers": [...],
    "commits": [...],
    "changes": [...],
    "daily_metrics": [...]
  }
}
```

---

## Data Models

### Core Structures

```go
// Developer represents a simulated developer
type Developer struct {
    ID              string                 `json:"id"`
    Email           string                 `json:"email"`
    Name            string                 `json:"name"`
    Region          string                 `json:"region"`      // US, EU, APAC
    Division        string                 `json:"division"`    // AGS, AT, ST
    Group           string                 `json:"group"`       // TMOBILE, ATANT
    Team            string                 `json:"team"`        // dev, support
    Skills          []string               `json:"skills"`
    ClientVersion   string                 `json:"client_version"`
    IsActive        bool                   `json:"is_active"`
    CreatedAt       time.Time              `json:"created_at"`
    LastActiveAt    time.Time              `json:"last_active_at"`
}

// Commit represents a single Git commit with AI-assisted metrics
type Commit struct {
    Hash              string    `json:"hash"`
    Timestamp         time.Time `json:"timestamp"`
    Message           string    `json:"message"`
    DeveloperID       string    `json:"developer_id"`
    Repository        string    `json:"repository"`
    Branch            string    `json:"branch"`
    LinesDelta        int       `json:"lines_delta"`
    LinesFromTAB      int       `json:"lines_from_tab"`
    LinesFromComposer int       `json:"lines_from_composer"`
    LinesNonAI        int       `json:"lines_non_ai"`
    IngestionTime     time.Time `json:"ingestion_time"`
}

// Change represents an individual AI-suggested change
type Change struct {
    ChangeID      string    `json:"change_id"`
    CommitHash    string    `json:"commit_hash"`
    DeveloperID   string    `json:"developer_id"`
    Timestamp     time.Time `json:"timestamp"`
    Source        string    `json:"source"` // TAB or COMPOSER
    Model         string    `json:"model"`
    FilesDelta    int       `json:"files_delta"`
    LinesAdded    int       `json:"lines_added"`
    LinesModified int       `json:"lines_modified"`
    IngestionTime time.Time `json:"ingestion_time"`
}

// DailyMetrics represents aggregated daily statistics
type DailyMetrics struct {
    Date                time.Time         `json:"date"`
    DAU                 int               `json:"dau"`
    AgentEditsAccepted  int               `json:"agent_edits_accepted"`
    TabCompletions      int               `json:"tab_completions"`
    ComposerEdits       int               `json:"composer_edits"`
    ModelUsage          map[string]int    `json:"model_usage"`
    TopFileExtensions   map[string]int    `json:"top_file_extensions"`
    CommandUsage        map[string]int    `json:"command_usage"`
}
```

---

## Configuration Specification

### JSON Input Schema

```json
{
  "auth": {
    "api_key": "cursor_test_key_12345",
    "api_secret": "cursor_test_secret_67890",
    "team_id": "team_default"
  },
  "simulation": {
    "developers": 100,
    "regions": ["US", "EU", "APAC"],
    "region_distribution": {
      "US": 0.5,
      "EU": 0.35,
      "APAC": 0.15
    },
    "divisions": ["AGS", "AT", "ST"],
    "division_distribution": {
      "AGS": 0.4,
      "AT": 0.35,
      "ST": 0.25
    },
    "groups": ["TMOBILE", "ATANT"],
    "group_distribution": {
      "TMOBILE": 0.6,
      "ATANT": 0.4
    },
    "teams": ["dev", "support"],
    "team_distribution": {
      "dev": 0.75,
      "support": 0.25
    },
    "pr_velocity": "high",
    "volatility": 0.35,
    "break_condition": {
      "type": "pr_count",
      "value": 100000
    }
  },
  "api": {
    "port": 8080,
    "host": "localhost",
    "enable_graphql": true
  },
  "output": {
    "export_format": "json",
    "export_path": "./exports"
  }
}
```

### Configuration Defaults

| Parameter | Default | Bounds |
|-----------|---------|--------|
| developers | 50 | 1-10000 |
| pr_velocity | "medium" | "low"\|"medium"\|"high" |
| volatility | 0.2 | 0.0-1.0 |
| port | 8080 | 1024-65535 |
| pageSize | 100 | 1-1000 |

---

## API Specifications

### REST API Response Envelopes

**Standard Response**:
```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "pageSize": 100,
    "total": 5000,
    "hasMore": true
  },
  "params": {
    "startDate": "2024-01-01",
    "endDate": "2024-01-31",
    "user": null
  },
  "timestamp": "2024-01-15T14:30:00Z"
}
```

**Error Response**:
```json
{
  "error": {
    "code": "INVALID_DATE_RANGE",
    "message": "Start date cannot be after end date",
    "details": {}
  }
}
```

### Authentication

**Mechanism**: HTTP Basic Authentication
```
Authorization: Basic {base64(api_key:api_secret)}
```

**Credential Validation**:
- Fixed credentials from config
- Return 401 Unauthorized if invalid
- Rate limiting: 100 requests/minute for team endpoints

---

## CLI Interface & Controls

### Startup Sequence

1. **Load Configuration**: Read JSON from file or stdin
2. **Validate Parameters**: Check bounds and consistency
3. **Initialize Database**: Create in-memory storage
4. **Generate Initial Data**: Create developer profiles
5. **Start API Server**: Bind to configured port
6. **Start Generation Loop**: Begin event generation
7. **Display Dashboard**: Show interactive status

### Signal Handlers

| Signal | Action |
|--------|--------|
| `Ctrl+S` (SIGTERM) | Soft stop - pause generation, keep API running |
| `Ctrl+E` (Custom) | Export database to file |
| `Ctrl+C` (SIGINT) | Graceful shutdown - stop API, save state, exit |
| `Ctrl+Q` (Custom) | Display quick stats snapshot |

### Status Display Refresh Rate
- Update every 1 second for live metrics
- Calculate rolling averages (1-min, 5-min, 15-min)

---

## Development Task Breakdown

### Folder Structure

```
services/cursor-sim/
├── DESIGN.md                          # This document
├── SPEC.md                            # Original spec
├── README.md                          # Quick start guide
├── Makefile                           # Build targets
├── go.mod                             # Go module definition
├── go.sum                             # Dependency lock
│
├── cmd/
│   └── simulator/
│       └── main.go                    # Entry point
│
├── internal/
│   ├── config/
│   │   ├── loader.go                  # Config file parsing
│   │   ├── validator.go               # Input validation
│   │   └── loader_test.go
│   │
│   ├── models/
│   │   ├── developer.go               # Developer data model
│   │   ├── commit.go                  # Commit data model
│   │   ├── change.go                  # Change data model
│   │   ├── metrics.go                 # Metrics data model
│   │   └── models_test.go
│   │
│   ├── generator/
│   │   ├── developer_generator.go     # Create developer profiles
│   │   ├── commit_generator.go        # Create commits
│   │   ├── change_generator.go        # Create changes
│   │   ├── metrics_calculator.go      # Calculate metrics
│   │   └── generator_test.go
│   │
│   ├── db/
│   │   ├── store.go                   # In-memory DB interface
│   │   ├── memory_store.go            # sync.Map implementation
│   │   ├── indexer.go                 # Indexing & querying
│   │   └── store_test.go
│   │
│   ├── api/
│   │   ├── server.go                  # HTTP server setup
│   │   ├── auth.go                    # Basic auth middleware
│   │   ├── handlers.go                # REST endpoint handlers
│   │   └── api_test.go
│   │
│   ├── graphql/
│   │   ├── schema.go                  # GraphQL schema definition
│   │   ├── resolvers.go               # Query resolvers
│   │   └── graphql_test.go
│   │
│   ├── cli/
│   │   ├── controller.go              # CLI interaction handler
│   │   ├── dashboard.go               # Terminal UI display
│   │   ├── signals.go                 # Signal handling
│   │   └── cli_test.go
│   │
│   └── export/
│       ├── exporter.go                # Export to JSON/binary
│       └── exporter_test.go
│
└── tests/
    ├── integration_test.go            # End-to-end tests
    ├── fixtures/
    │   └── config.sample.json         # Example config
    └── testdata/
        └── expected_exports/          # Expected export samples
```

---

## Testing Strategy

### Test Coverage Requirements
- **Minimum Coverage**: 80% across all packages
- **Unit Tests**: For each component in isolation
- **Integration Tests**: Full workflow testing
- **API Tests**: REST and GraphQL endpoint verification

### Test Categories

1. **Config Tests**: Validation, parsing, defaults
2. **Model Tests**: Data structure creation and validation
3. **Generator Tests**: Realistic data generation
4. **Database Tests**: CRUD operations, indexing, querying
5. **API Tests**: Authentication, response formats, pagination
6. **CLI Tests**: Signal handling, display formatting
7. **Integration Tests**: Full simulator workflow

### Example Test File Structure

```go
// internal/config/loader_test.go
package config

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestLoadValidConfig(t *testing.T) { ... }
func TestLoadConfigMissingFile(t *testing.T) { ... }
func TestValidateConfigBounds(t *testing.T) { ... }
func TestApplyConfigDefaults(t *testing.T) { ... }
```

---

## Success Criteria

### Phase Completion Checklist

- [ ] Configuration system fully functional with validation
- [ ] Developer generation creates realistic profiles
- [ ] Commit/change generation follows Poisson distribution
- [ ] In-memory database supports efficient queries
- [ ] REST API endpoints match Cursor API format exactly
- [ ] GraphQL API supports all required query types
- [ ] CLI dashboard displays real-time metrics
- [ ] Signal handlers (Ctrl+S, E, C) all functional
- [ ] Data export works for both JSON and binary formats
- [ ] Test coverage ≥80% across all packages
- [ ] Documentation complete and examples provided
- [ ] Docker integration verified with docker-compose

---

## Appendix: Cursor API Reference

### AI Code Tracking API Endpoints
- `GET /analytics/ai-code/commits` - JSON or CSV
- `GET /analytics/ai-code/changes` - JSON or CSV

**Supported Filters**:
- startDate / endDate (ISO 8601 or relative: "7d", "now")
- page / pageSize (default: 100, max: 1000)
- user (email or ID)

### Analytics API Endpoints
- Team-level: agent-edits, tabs, dau, models, client-versions, top-file-extensions, mcp, commands, plans, ask-mode, leaderboard
- By-user variants available for most endpoints

**Rate Limits**:
- Team endpoints: 100 req/min
- User endpoints: 50 req/min

### Response Structure
```json
{
  "data": [...],
  "pagination": {...},
  "params": {...}
}
```

---

## Next Steps

1. **Review this design** with stakeholders
2. **Approve folder structure** and task breakdown
3. **Begin Feature 1**: Configuration & Initialization
4. **Iterate through features** following atomic task lists
5. **Maintain 80%+ test coverage** throughout development

