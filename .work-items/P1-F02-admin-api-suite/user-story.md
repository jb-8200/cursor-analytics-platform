# User Story: P1-F02 Admin API Suite

**Feature**: P1-F02 Admin API Suite
**Phase**: P1 (Foundation)
**Created**: January 10, 2026
**Dependencies**: P1-F01 (Environment Variables - will be completed as part of this feature)

---

## Problem Statement

As a **DevOps engineer and platform administrator**, I need a comprehensive Admin API Suite for cursor-sim so that I can:

1. **Configure simulation parameters dynamically** without restarting containers (especially for large datasets: 1200 developers, 400 days)
2. **Manage seed files** (upload/swap org/team structures) in JSON, YAML, or CSV format without rebuilding images
3. **Inspect current runtime configuration and seed structure** before making changes
4. **Monitor comprehensive statistics** including quality metrics, variance, and time series data
5. **Control initial startup** via environment variables for standard deployments

### Current Limitations

- Cannot reconfigure without restarting Docker container or Cloud Run service
- No way to dynamically change org/team structures (seed files are baked into image)
- Cannot inspect current configuration or seed data via API
- No visibility into generated data statistics (commits, quality metrics, variance)
- Environment variables only support MODE, SEED, PORT, DAYS, VELOCITY (missing DEVELOPERS, MAX_COMMITS, MONTHS)

### Business Impact

- **Slow iteration**: Must rebuild/restart containers to test different scenarios
- **Poor observability**: Cannot inspect what data was generated or current configuration
- **Limited flexibility**: Cannot swap seed files or org structures dynamically
- **Operational overhead**: GCP Cloud Run deployments require image rebuilds for parameter changes

---

## User Stories (EARS Format)

### Story 1: Environment Variable Configuration

**AS** a DevOps engineer
**I WANT** to configure all generation parameters via environment variables
**SO THAT** I can deploy cursor-sim to Docker or GCP Cloud Run with custom settings without code changes

**GIVEN** I have cursor-sim running in a container
**WHEN** I set environment variables `CURSOR_SIM_DEVELOPERS=1200`, `CURSOR_SIM_DAYS=400`, `CURSOR_SIM_VELOCITY=high`, `CURSOR_SIM_MAX_COMMITS=500`
**THEN** cursor-sim generates data for 1200 developers over 400 days with high velocity and max 500 commits per developer
**AND** I can verify this via the `/admin/config` endpoint

### Story 2: Runtime Reconfiguration (Override Mode)

**AS** a platform administrator
**I WANT** to regenerate simulation data with new parameters via API without restarting
**SO THAT** I can quickly test different scenarios (small team vs enterprise) without container restarts

**GIVEN** cursor-sim is running with default configuration (50 developers, 90 days)
**WHEN** I call `POST /admin/regenerate` with `{"mode":"override","developers":1200,"days":400,"velocity":"high","max_commits":500}`
**THEN** cursor-sim clears existing data and regenerates for 1200 developers over 400 days
**AND** returns statistics showing ~600,000 commits generated
**AND** I can verify via `GET /teams/members` that there are now 1200 developers

### Story 3: Incremental Data Addition (Append Mode)

**AS** a data analyst
**I WANT** to append additional historical data to existing simulation
**SO THAT** I can extend time series without losing existing data

**GIVEN** cursor-sim has 90 days of historical data
**WHEN** I call `POST /admin/regenerate` with `{"mode":"append","days":30,"velocity":"medium"}`
**THEN** cursor-sim adds 30 more days of data without clearing existing commits
**AND** I now have 120 days of cumulative historical data

### Story 4: Dynamic Seed Management

**AS** a platform administrator
**I WANT** to upload new seed files (org/team structures) dynamically
**SO THAT** I can test different organizational hierarchies without rebuilding Docker images

**GIVEN** cursor-sim is running with the default seed file
**WHEN** I call `POST /admin/seed` with a new org structure (3 divisions, 10 teams, 100 developers)
**THEN** cursor-sim loads the new seed data and reports the org structure
**AND** I can optionally trigger regeneration in the same request
**AND** future data generation uses the new org/team hierarchy

### Story 5: Configuration Inspection

**AS** a platform administrator
**I WANT** to inspect the current runtime configuration and seed structure
**SO THAT** I can review settings before making changes and verify deployment configuration

**GIVEN** cursor-sim is running
**WHEN** I call `GET /admin/config`
**THEN** I receive:
- Current generation parameters (days, velocity, developers, max_commits)
- Active seed structure (org, divisions, teams, developers count)
- Developer breakdowns by seniority, region, team
- External data sources configuration (Harvey, Copilot, Qualtrics)
- Server info (port, version, uptime)

### Story 6: Comprehensive Statistics

**AS** a data analyst
**I WANT** to retrieve comprehensive statistics about generated data
**SO THAT** I can verify data quality, variance, and distribution

**GIVEN** cursor-sim has generated data
**WHEN** I call `GET /admin/stats`
**THEN** I receive:
- Generation stats (total commits, PRs, reviews, issues, data size)
- Developer stats (by seniority, region, team, activity level)
- Quality metrics (revert rate, hotfix rate, code survival, review thoroughness)
- Variance metrics (commits std dev, PR size std dev, cycle time std dev)
- Performance metrics (memory usage, generation time)
- Optional time series data (commits per day, PRs per day, cycle times)

---

## Acceptance Criteria

### AC1: Environment Variable Support
- [ ] `CURSOR_SIM_DEVELOPERS` environment variable sets developer count
- [ ] `CURSOR_SIM_MONTHS` environment variable converts to days (months × 30)
- [ ] `CURSOR_SIM_MAX_COMMITS` environment variable limits commits per developer
- [ ] Environment variables override flag defaults but are overridden by explicit CLI flags
- [ ] Docker Compose reads from `.env` file automatically
- [ ] GCP Cloud Run deployment script includes new environment variables
- [ ] Tests verify env var precedence: CLI flags > Env vars > Defaults
- [ ] Interactive mode (`-interactive`) does not conflict with environment variables

### AC2: Admin Regenerate API
- [ ] `POST /admin/regenerate` endpoint accepts `mode`, `days`, `velocity`, `developers`, `max_commits`
- [ ] Override mode clears all existing data before regenerating
- [ ] Append mode adds new data to existing dataset
- [ ] Returns detailed statistics: commits_added, prs_added, reviews_added, total counts, duration
- [ ] Validates parameters: mode (append/override), days (1-3650), velocity (low/medium/high), developers (0-10000), max_commits (0-100000)
- [ ] Thread-safe operations using existing storage mutex locks
- [ ] Returns 400 Bad Request for invalid parameters
- [ ] Returns 401 Unauthorized if API key missing

### AC3: Seed Management API
- [ ] `POST /admin/seed` endpoint accepts `data`, `format`, `regenerate`, `regenerate_config`
- [ ] Supports JSON, YAML, and CSV formats
- [ ] Validates seed data before accepting (required fields, org structure)
- [ ] Reports org/division/team structure after upload
- [ ] Optional auto-regeneration with regenerate_config parameter
- [ ] `GET /admin/seed/presets` returns predefined seed configurations
- [ ] Presets include: small-team, medium-team, enterprise, multi-region
- [ ] Thread-safe seed swapping

### AC4: Configuration Inspection API
- [ ] `GET /admin/config` returns current generation parameters
- [ ] Returns active seed structure (developers, repos, org hierarchy)
- [ ] Returns developer breakdowns by seniority, region, team
- [ ] Returns external data sources (Harvey, Copilot, Qualtrics)
- [ ] Returns server info (port, version, uptime)
- [ ] No mutation operations (GET only)

### AC5: Statistics API
- [ ] `GET /admin/stats` returns generation statistics
- [ ] Returns developer statistics by seniority, region, team, activity
- [ ] Returns quality metrics (revert rate, hotfix rate, code survival, review thoroughness)
- [ ] Returns variance metrics (std dev for commits, PR size, cycle time)
- [ ] Returns performance metrics (memory usage, generation time, data size)
- [ ] Optional `?include_timeseries=true` query parameter adds time series data
- [ ] Returns org structure (teams, divisions, repositories)

### AC6: Documentation
- [ ] `.env.example` documents all environment variables with descriptions
- [ ] `docker-compose.yml` includes all supported environment variables
- [ ] `services/cursor-sim/SPEC.md` documents all Admin API endpoints
- [ ] `services/cursor-sim/README.md` includes Admin API usage examples
- [ ] `tools/deploy-cursor-sim.sh` supports new environment variables
- [ ] API response schemas documented in SPEC.md

### AC7: Testing
- [ ] Unit tests for environment variable parsing (6 test cases)
- [ ] Unit tests for regenerate handler (append/override modes, validation)
- [ ] Unit tests for seed upload handler (JSON/YAML/CSV formats)
- [ ] Unit tests for config and stats handlers
- [ ] Storage layer tests for ClearAllData() and GetStats()
- [ ] E2E tests for all Admin API endpoints
- [ ] Test coverage: 80%+ for all new code

### AC8: Non-Breaking Changes
- [ ] Existing CLI behavior unchanged
- [ ] Existing API endpoints unchanged
- [ ] Backward compatible with existing deployments
- [ ] No changes to default behavior when env vars not set

---

## Out of Scope

- Web UI for admin operations (API-only for now)
- Authentication/authorization beyond existing API key
- Role-based access control
- Audit logging of admin operations
- Rate limiting for admin endpoints
- Seed file versioning or rollback
- Scheduled/cron-based regeneration
- Webhooks or notifications

---

## Technical Constraints

1. **Thread Safety**: All storage operations must use existing mutex locks
2. **API Authentication**: Use existing Basic Auth with API key
3. **Validation**: Strict parameter validation to prevent runaway generation
4. **Backward Compatibility**: No breaking changes to existing behavior
5. **Memory Limits**: Large datasets (1200 developers × 400 days) may require memory profiling
6. **GCP Cloud Run**: Environment variables must work in Cloud Run deployment

---

## Success Metrics

1. **Environment Variables**: Can deploy cursor-sim to Docker/GCP Cloud Run with custom parameters via `.env` file or Cloud Run env vars
2. **Dynamic Reconfiguration**: Can change from 50 developers to 1200 developers without container restart
3. **Seed Management**: Can upload new org structures (JSON/YAML/CSV) and regenerate data in single API call
4. **Observability**: Can inspect current config, seed structure, and statistics via GET endpoints
5. **Performance**: Regeneration for 1200 developers × 400 days completes in <30 seconds
6. **Test Coverage**: 80%+ coverage for all new code
7. **Documentation**: Complete API documentation in SPEC.md with curl examples

---

## Example Workflows

### Workflow 1: Docker Compose with .env file
```bash
# Create .env
cat > .env << EOF
CURSOR_SIM_DEVELOPERS=1200
CURSOR_SIM_DAYS=400
CURSOR_SIM_VELOCITY=high
CURSOR_SIM_MAX_COMMITS=500
EOF

# Start
docker-compose up -d cursor-sim

# Verify
curl -u cursor-sim-dev-key: http://localhost:8080/admin/config | jq
```

### Workflow 2: Runtime Reconfiguration
```bash
# Start with defaults
docker-compose up -d cursor-sim

# Reconfigure without restart
curl -X POST -u cursor-sim-dev-key: \
  -H "Content-Type: application/json" \
  -d '{"mode":"override","days":400,"velocity":"high","developers":1200,"max_commits":500}' \
  http://localhost:8080/admin/regenerate

# Verify
curl -u cursor-sim-dev-key: http://localhost:8080/admin/stats | jq
```

### Workflow 3: Upload Seed and Regenerate
```bash
# Upload new seed + regenerate in one call
curl -X POST -u cursor-sim-dev-key: \
  -H "Content-Type: application/json" \
  -d '{
    "format": "json",
    "data": "{...seed JSON...}",
    "regenerate": true,
    "regenerate_config": {"mode":"override","days":180,"velocity":"high"}
  }' \
  http://localhost:8080/admin/seed
```

---

**Estimated Effort**: 20.5 hours (2.5-3 days)
**Priority**: High (enables large-scale testing and operational flexibility)
**Risk**: Medium (thread safety, memory usage for large datasets)
