# cursor-sim v2.0.0

> High-fidelity Cursor Business API simulator for testing and development

cursor-sim generates synthetic developer usage data that exactly matches the Cursor Business API schema. Perfect for testing analytics pipelines, generating research datasets, or developing against the Cursor API without production access.

## Table of Contents

- [Quick Start](#quick-start)
- [Installation](#installation)
- [Usage](#usage)
  - [Interactive Mode](#interactive-mode-recommended)
  - [Non-Interactive Mode](#non-interactive-mode)
  - [Configuration Options](#configuration-options)
- [API Reference](#api-reference)
- [Seed Files](#seed-files)
- [Common Use Cases](#common-use-cases)
- [Docker Deployment](#docker-deployment)
- [Advanced Topics](#advanced-topics)

---

## Quick Start

### 1. Build the binary

```bash
cd services/cursor-sim
go build -o bin/cursor-sim ./cmd/simulator
```

### 2. Run with interactive configuration (easiest!)

```bash
./bin/cursor-sim -mode runtime -seed testdata/valid_seed.json -interactive
```

You'll be prompted to configure:
- Number of developers (1-100)
- Time period in months (1-24)
- Maximum commits per developer (100-2000)

### 3. Test the API

```bash
# Health check (no auth required)
curl http://localhost:8080/health

# Get team members (requires API key)
curl -u cursor-sim-dev-key: http://localhost:8080/teams/members

# Get commit history
curl -u cursor-sim-dev-key: http://localhost:8080/analytics/ai-code/commits
```

---

## Installation

### Prerequisites

- **Go 1.22+** ([install](https://go.dev/doc/install))
- **Git** (for cloning the repository)

### Build from Source

```bash
# Clone the repository
git clone <repository-url>
cd cursor-analytics-platform/services/cursor-sim

# Download dependencies
go mod download

# Build the binary
go build -o bin/cursor-sim ./cmd/simulator

# Verify installation
./bin/cursor-sim --help
```

### Run Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test ./... -cover

# Run specific package tests
go test ./internal/config/...
```

---

## Usage

cursor-sim has two operation modes:

### Interactive Mode (Recommended)

**Best for:** Experimentation, quick testing, demos

The interactive mode prompts you for configuration values with helpful defaults:

```bash
./bin/cursor-sim -mode runtime -seed testdata/valid_seed.json -interactive
```

**Example session:**

```
cursor-sim v2.0.0

Cursor Simulator - Interactive Configuration
Press Enter to use default values
Number of developers (default: 10): 5
Period in months (default: 6): 3
Maximum commits per developer (default: 500): 250

Configuration Summary:
  Developers: 5
  Period: 3 months (90 days)
  Max commits: 250

2026/01/08 12:00:00 Using interactive configuration: 5 developers, 90 days, max 250 commits
2026/01/08 12:00:00 Starting cursor-sim v2.0.0...
2026/01/08 12:00:00 Loading seed data from testdata/valid_seed.json...
2026/01/08 12:00:00 Loaded 2 developers from seed file, replicated to 5 developers
2026/01/08 12:00:00 Generating 90 days of commit history...
2026/01/08 12:00:00 Generated 245 commits across 5 developers
2026/01/08 12:00:00 HTTP server listening on port 8080
```

**Interactive Configuration:**
- **Developers**: How many developers to simulate (replicates from seed file)
- **Period**: How many months of history to generate
- **Max Commits**: Limit commits per developer (prevents runaway generation)

**Defaults:**
- 10 developers
- 6 months (180 days)
- 500 commits per developer

### Non-Interactive Mode

**Best for:** Automation, CI/CD, scripting

Specify all configuration via flags:

```bash
./bin/cursor-sim \
  -mode runtime \
  -seed testdata/valid_seed.json \
  -developers 5 \
  -months 3 \
  -max-commits 250 \
  -velocity high \
  -port 8080
```

**Equivalent to interactive example above, but scriptable!**

### Configuration Options

#### Core Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-mode` | string | `runtime` | Operation mode: `runtime` (generate data) or `replay` (from corpus) |
| `-seed` | string | _(required)_ | Path to seed.json file with team/developer configuration |
| `-port` | int | `8080` | HTTP server port |
| `-days` | int | `90` | Days of history to generate (overridden by `-months` or interactive) |
| `-velocity` | string | `medium` | Event generation rate: `low`, `medium`, `high` |

#### Interactive Mode (P4-F02)

| Flag | Description |
|------|-------------|
| `-interactive` | Enable interactive configuration prompts |

#### Non-Interactive Generation Parameters (P4-F02)

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-developers` | int | `0` | Number of developers (0 = use all from seed) |
| `-months` | int | `0` | Period in months (converted to days automatically) |
| `-max-commits` | int | `0` | Maximum commits per developer (0 = unlimited) |

**Note:** Cannot mix `-interactive` with `-developers`, `-months`, or `-max-commits`.

#### Velocity Settings

| Velocity | Commits/Day/Developer | Use Case |
|----------|----------------------|----------|
| `low` | ~5 | Small teams, proof of concepts |
| `medium` | ~15 | Realistic team activity |
| `high` | ~30 | High-activity teams, load testing |

#### Environment Variables

All flags can be set via environment variables. **CLI flags take precedence over environment variables.**

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `CURSOR_SIM_MODE` | string | `runtime` | Operation mode: runtime, replay, or preview |
| `CURSOR_SIM_SEED` | string | _(required)_ | Path to seed.json file |
| `CURSOR_SIM_PORT` | int | `8080` | HTTP server port |
| `CURSOR_SIM_DAYS` | int | `90` | Days of history to generate |
| `CURSOR_SIM_VELOCITY` | string | `medium` | Event rate: low, medium, high |
| `CURSOR_SIM_DEVELOPERS` | int | `0` | Number of developers (0 = use seed count) |
| `CURSOR_SIM_MONTHS` | int | `0` | Period in months (converted to days = months × 30) |
| `CURSOR_SIM_MAX_COMMITS` | int | `0` | Maximum commits per developer (0 = unlimited) |

**Examples:**

```bash
# Example 1: Basic configuration via environment variables
export CURSOR_SIM_SEED=testdata/valid_seed.json
export CURSOR_SIM_PORT=8080
export CURSOR_SIM_DAYS=90
export CURSOR_SIM_VELOCITY=medium
./bin/cursor-sim

# Example 2: Large-scale simulation (1200 developers, 400 days, high velocity)
export CURSOR_SIM_DEVELOPERS=1200
export CURSOR_SIM_DAYS=400
export CURSOR_SIM_VELOCITY=high
export CURSOR_SIM_MAX_COMMITS=500
./bin/cursor-sim -seed testdata/valid_seed.json

# Example 3: Using months instead of days
export CURSOR_SIM_MONTHS=12  # Automatically converts to 360 days
./bin/cursor-sim -seed testdata/valid_seed.json -developers 100

# Example 4: Docker Compose with .env file
cat > .env << EOF
CURSOR_SIM_DEVELOPERS=1200
CURSOR_SIM_DAYS=400
CURSOR_SIM_VELOCITY=high
CURSOR_SIM_MAX_COMMITS=500
EOF
docker-compose up -d cursor-sim
```

**Precedence:** CLI flags > Environment variables > Defaults

---

## API Reference

cursor-sim provides 29 REST endpoints matching the Cursor Business API schema.

### Authentication

All endpoints except `/health` require **Basic Auth**:
- **Username:** API key (default: `cursor-sim-dev-key`)
- **Password:** _(leave empty)_

```bash
curl -u cursor-sim-dev-key: http://localhost:8080/endpoint
```

### Key Endpoints

#### Health Check

```bash
GET /health
```

No authentication required. Returns service status.

**Response:**
```json
{
  "status": "healthy",
  "service": "cursor-sim",
  "version": "2.0.0"
}
```

#### Team Members

```bash
GET /teams/members
```

List all developers in the simulated team.

**Response:**
```json
{
  "data": [
    {
      "id": "user_001",
      "email": "alice@example.com",
      "name": "Alice Johnson",
      "role": "engineer"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 100,
    "total": 5
  }
}
```

#### AI Code Commits

```bash
GET /analytics/ai-code/commits
```

Query commit history with filtering and pagination.

**Query Parameters:**
- `from` (string): Start date YYYY-MM-DD
- `to` (string): End date YYYY-MM-DD
- `page` (int): Page number (default 1)
- `page_size` (int): Items per page (default 100, max 500)
- `user_id` (string): Filter by developer email
- `repo_name` (string): Filter by repository

**Example:**
```bash
curl -u cursor-sim-dev-key: \
  "http://localhost:8080/analytics/ai-code/commits?from=2026-01-01&to=2026-01-31&page_size=50"
```

**Response:**
```json
{
  "data": [
    {
      "commit_id": "abc123",
      "user_id": "alice@example.com",
      "repo_name": "backend-api",
      "timestamp": "2026-01-15T10:30:00Z",
      "ai_assisted": true,
      "lines_added": 45,
      "lines_deleted": 12
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 50,
    "total": 245
  }
}
```

#### Team Analytics

```bash
GET /analytics/team/commits           # Commit activity over time
GET /analytics/team/users             # User activity metrics
GET /analytics/team/repos             # Repository statistics
GET /analytics/team/models            # AI model usage
GET /analytics/team/client-versions   # Cursor client versions
GET /analytics/team/file-extensions   # Programming languages used
GET /analytics/team/mcp               # MCP server usage
GET /analytics/team/commands          # Command palette usage
```

See [SPEC.md](./SPEC.md) for complete endpoint documentation.

---

## Seed Files

Seed files define your simulated team's structure and behavior patterns.

### Minimal Seed File

```json
{
  "team": {
    "id": "team_001",
    "name": "Engineering Team",
    "plan": "business"
  },
  "developers": [
    {
      "id": "user_001",
      "email": "alice@example.com",
      "name": "Alice Johnson",
      "role": "engineer",
      "timezone": "America/New_York"
    },
    {
      "id": "user_002",
      "email": "bob@example.com",
      "name": "Bob Smith",
      "role": "engineer",
      "timezone": "America/Los_Angeles"
    }
  ],
  "repositories": [
    {
      "id": "repo_001",
      "name": "backend-api",
      "language": "go"
    },
    {
      "id": "repo_002",
      "name": "frontend-app",
      "language": "typescript"
    }
  ],
  "models": ["claude-3-5-sonnet-20241022", "claude-3-haiku-20240307"],
  "file_extensions": [".go", ".ts", ".tsx", ".py", ".md"]
}
```

### Developer Replication

If you specify more developers than in your seed file, cursor-sim **replicates** developers with unique IDs:

```bash
# Seed file has 2 developers, but you want 10
./bin/cursor-sim -seed seed.json -developers 10 -months 6
```

Generated developers:
- `alice@example.com` → `alice_clone1@example.com`, `alice_clone2@example.com`, ...
- `bob@example.com` → `bob_clone1@example.com`, `bob_clone2@example.com`, ...

### Example Seed Files

- **testdata/valid_seed.json** - Minimal 2-developer team
- **testdata/large_team_seed.json** - Enterprise team (if exists)

---

## Common Use Cases

### 1. Quick Testing (Minimal Data)

Generate 1 week of data for rapid iteration:

```bash
./bin/cursor-sim -seed testdata/valid_seed.json -days 7 -velocity low
```

### 2. Realistic Team Simulation

Simulate a 5-person team over 6 months:

```bash
./bin/cursor-sim -seed testdata/valid_seed.json -interactive
# Enter: 5 developers, 6 months, 500 commits
```

### 3. Load Testing

Generate high-volume data for performance testing:

```bash
./bin/cursor-sim -seed testdata/valid_seed.json -developers 50 -months 12 -velocity high
```

### 4. Analytics Pipeline Testing

Run cursor-sim as a test fixture for analytics-core:

```bash
# Terminal 1: Start cursor-sim
./bin/cursor-sim -seed testdata/valid_seed.json -interactive -port 8080

# Terminal 2: Point analytics-core at cursor-sim
cd ../cursor-analytics-core
CURSOR_API_URL=http://localhost:8080 npm run dev
```

### 5. Demo/Presentation

Generate clean, predictable data for demos:

```bash
./bin/cursor-sim \
  -seed testdata/valid_seed.json \
  -developers 3 \
  -months 3 \
  -max-commits 100 \
  -velocity medium
```

### 6. Research Export

Generate data and export for analysis:

```bash
# Start simulator
./bin/cursor-sim -seed testdata/valid_seed.json -developers 20 -months 12

# Export commits as CSV
curl -u cursor-sim-dev-key: http://localhost:8080/analytics/ai-code/commits.csv > commits.csv

# Analyze in Excel, Python, R, etc.
```

---

## Docker Deployment

### Quick Start with Docker

**Automated script (recommended):**

```bash
# From project root
./tools/docker-local.sh

# With custom configuration
DAYS=30 VELOCITY=high ./tools/docker-local.sh
```

**Manual Docker commands:**

```bash
# Build Docker image
docker build -t cursor-sim:latest services/cursor-sim

# Run with Docker
docker run --rm -p 8080:8080 \
  -e CURSOR_SIM_MODE=runtime \
  -e CURSOR_SIM_SEED=/app/seed.json \
  cursor-sim:latest
```

### GCP Cloud Run Deployment

Deploy cursor-sim to Google Cloud Run with automated SSL and custom domains:

```bash
# Deploy to staging (default .run.app URL)
./tools/deploy-cursor-sim.sh staging

# Deploy with custom domain
CUSTOM_DOMAIN=dox-a3.jishutech.io ./tools/deploy-cursor-sim.sh staging

# Deploy to production
CUSTOM_DOMAIN=api.yourdomain.com ./tools/deploy-cursor-sim.sh production
```

**Features:**
- Automatic SSL certificate provisioning
- Custom domain mapping support
- Zero-downtime deployments
- Auto-scaling (scale-to-zero for staging)
- Managed infrastructure

**Current Deployment:**
- **Staging URL**: https://cursor-sim-7m3ityidxa-uc.a.run.app
- **Custom Domain**: https://dox-a3.jishutech.io
- **Region**: us-central1
- **Project**: cursor-sim

### Complete Docker Documentation

For comprehensive Docker deployment information, see **[DOCKER.md](./DOCKER.md)**:

- Multi-stage build architecture
- Security features and performance metrics
- Custom seed file strategies
- Container management commands
- Troubleshooting Docker-specific issues
- Production deployment (Cloud Run, Kubernetes, Docker Compose)
- Configuration presets

**Also see:**
- [docs/cursor-sim-cloud-run.md](../../docs/cursor-sim-cloud-run.md) - Complete GCP Cloud Run deployment guide with custom domains

---

## Advanced Topics

### Replay Mode (Deferred)

Replay mode loads pre-generated event corpus files:

```bash
./bin/cursor-sim -mode replay -corpus data/events.parquet
```

**Status:** Deferred to Phase 3D (not yet implemented)

### Custom Generators

cursor-sim generates:
- **Commits** - Git commits with AI assistance metadata
- **Model Usage** - Claude model selection events
- **Client Versions** - Cursor app version distribution
- **File Extensions** - Programming language usage
- **Features** - MCP servers, commands, plans, ask mode

All generators use **Poisson processes** for realistic temporal distribution.

### Performance Tuning

**Memory usage** scales with:
- Number of developers
- Number of days
- Velocity setting

**Typical usage:**
- 10 developers × 6 months × medium velocity ≈ **50MB RAM**
- 100 developers × 12 months × high velocity ≈ **500MB RAM**

**Commit generation time:**
- 10 developers × 6 months ≈ **1-2 seconds**
- 100 developers × 12 months ≈ **5-10 seconds**

### Concurrent Requests

cursor-sim uses `sync.Map` for thread-safe storage. All endpoints support concurrent requests.

### Data Consistency

Generated data maintains referential integrity:
- All commits reference valid developers
- All commits reference valid repositories
- Model usage aligns with team plan
- Timestamps are chronologically ordered

---

## Troubleshooting

### "validation failed: seed path is required"

**Solution:** Provide `-seed` flag:
```bash
./bin/cursor-sim -mode runtime -seed testdata/valid_seed.json
```

### "failed to load seed data: no such file"

**Solution:** Check seed file path:
```bash
ls -la testdata/valid_seed.json
./bin/cursor-sim -seed testdata/valid_seed.json
```

### "bind: address already in use"

**Solution:** Use different port:
```bash
./bin/cursor-sim -seed testdata/valid_seed.json -port 8081
```

### API returns 401 Unauthorized

**Solution:** Include API key in request:
```bash
curl -u cursor-sim-dev-key: http://localhost:8080/teams/members
```

### Generated data is empty

**Solution:** Check that generators ran successfully in logs. If using `-max-commits`, ensure value isn't too restrictive.

---

## Platform Integration

cursor-sim is part of the Cursor Analytics Platform with two analytics paths:

### Path 1: GraphQL Analytics (Original)

```
┌─────────────┐       ┌──────────────────┐       ┌─────────────┐
│ cursor-sim  │──────▶│ analytics-core   │──────▶│  viz-spa    │
│   (P4)      │ REST  │     (P5)         │GraphQL│    (P6)     │
│  Port 8080  │       │  Port 4000       │       │  Port 3000  │
└─────────────┘       └──────────────────┘       └─────────────┘
                             │
                             ▼
                      ┌─────────────┐
                      │ PostgreSQL  │
                      │  Port 5432  │
                      └─────────────┘
```

### Path 2: dbt + Streamlit Analytics (New)

```
┌─────────────┐       ┌───────────────────┐
│ cursor-sim  │──────▶│ streamlit-dash    │
│   (P4)      │ REST  │     (P9)          │
│  Port 8080  │       │  Port 8501        │
└─────────────┘       │                   │
                      │ ┌───────────────┐ │
                      │ │ dbt (P8)      │ │
                      │ │ Transforms    │ │
                      │ └───────────────┘ │
                      │        ▼          │
                      │ ┌───────────────┐ │
                      │ │ DuckDB        │ │
                      │ │ Analytics DB  │ │
                      │ └───────────────┘ │
                      └───────────────────┘
```

### Related Services

- **cursor-analytics-core (P5)** - GraphQL aggregation layer (PostgreSQL backend)
- **cursor-viz-spa (P6)** - React dashboard for visualizations
- **streamlit-dashboard (P9)** - Python analytics dashboard with embedded dbt pipeline (DuckDB/Snowflake backend)

### Running with Docker Compose

The complete platform (cursor-sim + streamlit-dashboard) can be started with:

```bash
# Start cursor-sim + streamlit-dashboard
docker-compose up -d cursor-sim streamlit-dashboard

# Access services
# - cursor-sim API: http://localhost:8080
# - Streamlit Dashboard: http://localhost:8501

# View logs
docker-compose logs -f cursor-sim streamlit-dashboard

# Stop services
docker-compose down cursor-sim streamlit-dashboard
```

For the complete 4-service stack (including GraphQL path):

```bash
# Start all services
docker-compose up -d

# Services available:
# - cursor-sim: http://localhost:8080
# - cursor-analytics-core: http://localhost:4000/graphql
# - cursor-viz-spa: http://localhost:3000
# - streamlit-dashboard: http://localhost:8501
```

### Related Documentation

- **Technical Specification:** [SPEC.md](./SPEC.md)
- **Platform Architecture:** [docs/DESIGN.md](../../docs/DESIGN.md)
- **API Contract:** `.claude/skills/api-contract/SKILL.md`
- **Cloud Deployment:** [docs/cursor-sim-cloud-run.md](../../docs/cursor-sim-cloud-run.md)

---

## Development

### Running Tests

```bash
# All tests
go test ./...

# Specific package
go test ./internal/config/...

# With coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Code Structure

```
services/cursor-sim/
├── cmd/simulator/          # Main application entry point
├── internal/
│   ├── api/               # HTTP handlers and routing
│   ├── config/            # Configuration and CLI parsing
│   ├── generator/         # Data generation engines
│   ├── models/            # Data structures
│   ├── seed/              # Seed file loading and validation
│   ├── server/            # HTTP server setup
│   └── storage/           # In-memory data storage
├── testdata/              # Test fixtures and seed files
├── Dockerfile             # Container build
├── SPEC.md               # Technical specification
└── README.md             # This file
```

### Contributing

1. Follow [SDD workflow](../../docs/spec-driven-design.md)
2. Write tests first (TDD)
3. Run `go fmt` before committing
4. Update SPEC.md if changing API
5. Follow [Go best practices](.claude/rules/03-coding-standards.md)

---

## Version History

- **v2.0.0** (Jan 2026) - Interactive mode, developer replication, commit limits
- **v1.0.0** (Dec 2025) - Initial release with 29 endpoints

---

## License

Internal development tool - not for external distribution.

---

## Support

- **Issues:** Report bugs in project issue tracker
- **Questions:** See [SPEC.md](./SPEC.md) for technical details
- **Architecture:** See [docs/DESIGN.md](../../docs/DESIGN.md) for platform overview
