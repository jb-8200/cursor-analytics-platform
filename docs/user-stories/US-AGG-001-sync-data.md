# User Story: US-AGG-001

## Sync Data from Simulator

**Story ID:** US-AGG-001  
**Feature:** [F002 - Aggregator Ingestion](../features/F002-aggregator-ingestion.md)  
**Priority:** High  
**Story Points:** 5

---

## Story

**As a** system administrator,  
**I want to** have the aggregator automatically sync data from the simulator,  
**So that** the dashboard always displays current information without manual intervention.

---

## Description

The aggregator service needs to periodically fetch usage data from the simulator and persist it to the PostgreSQL database. This background synchronization process is the foundation of the entire analytics pipeline—without reliable data ingestion, all downstream visualizations would be empty or stale.

The sync process must be resilient to temporary failures. If the simulator is briefly unavailable (for example, during a restart), the aggregator should wait and retry rather than crashing. When the simulator returns, the aggregator should pick up where it left off and backfill any missed time periods.

Additionally, the sync must be idempotent. Running the same sync twice for the same time period should not create duplicate records in the database. This is crucial because network hiccups might cause the aggregator to retry a request that actually succeeded.

---

## Acceptance Criteria

### AC1: Automatic Sync on Startup
**Given** the aggregator service starts  
**When** it establishes a database connection  
**Then** it immediately triggers a full sync from the simulator

### AC2: Periodic Polling
**Given** the aggregator is running  
**When** 60 seconds have elapsed since the last sync  
**Then** a new sync cycle begins automatically

### AC3: Configurable Interval
**Given** the environment variable `POLL_INTERVAL_MS=30000` is set  
**When** the aggregator starts  
**Then** syncs occur every 30 seconds instead of the default 60

### AC4: Graceful Failure Handling
**Given** the simulator is unreachable  
**When** a sync attempt fails  
**Then** the aggregator logs the error, waits with exponential backoff, and retries

### AC5: Idempotent Ingestion
**Given** data for 2025-01-01 has already been synced  
**When** the same data is synced again  
**Then** no duplicate records are created (upsert behavior)

### AC6: Backfill Support
**Given** the simulator was unavailable for 2 hours  
**When** it becomes available again  
**Then** the aggregator syncs all missed data from the gap period

### AC7: Progress Logging
**Given** a sync is in progress  
**When** records are being processed  
**Then** structured logs show progress (e.g., "Synced 50/100 developers")

---

## Technical Notes

The ingestion worker should run in a separate async context from the GraphQL server so that slow syncs do not block API requests. Consider using a dedicated worker thread or process.

The backoff strategy should be exponential: 5s, 10s, 20s, 40s, 80s, capping at 5 minutes. After 10 consecutive failures, the worker should enter a "degraded" state that is visible in the health check.

For upserts, use PostgreSQL's `ON CONFLICT DO UPDATE` clause with the unique constraint on `(developer_id, date)` for daily_stats.

---

## Definition of Done

- [ ] All acceptance criteria pass
- [ ] Unit tests mock the simulator and verify sync logic
- [ ] Integration tests use real PostgreSQL (via pg-mem) to verify upserts
- [ ] Logs include correlation IDs for tracing sync cycles
- [ ] Health endpoint reflects sync status
- [ ] Code reviewed and merged to main branch

---

## Related Tasks

- [TASK-012](../tasks/TASK-012-agg-ingestion.md): Implement ingestion worker
