# Insomnia Admin API Guide

This guide explains how to use the Insomnia Admin API collection for cursor-sim runtime management.

## Quick Start

### 1. Import the Collection

1. Open Insomnia
2. Click **Dashboard** → **Import**
3. Select the file: `Admin_APIs_2026-01-10.yaml`
4. The collection "cursor-sim Admin APIs" will be imported

### 2. Set Environment Variables

In Insomnia:

1. Click **Manage Environments** (bottom left)
2. Create or edit your environment with:

```json
{
  "baseUrl": "http://localhost:8080",
  "apiKey": "cursor-sim-dev-key"
}
```

### 3. Start Using Requests

All requests are organized in folders. Click any request and press **Send**

---

## Collection Structure

### Admin Configuration (Folder)

**Purpose:** Inspect current runtime settings

**Request:**
- **Get Current Configuration**
  - Method: GET
  - Path: `/admin/config`
  - Returns: Generation parameters, seed structure, external sources, server info

### Admin Statistics (Folder)

**Purpose:** Monitor data quality and generation progress

**Requests:**
- **Get Statistics (Basic)**
  - Method: GET
  - Path: `/admin/stats`
  - Returns: Commits, PRs, reviews, quality metrics, variance

- **Get Statistics (With Time Series)**
  - Method: GET
  - Path: `/admin/stats?include_timeseries=true`
  - Returns: All stats + commits per day, PRs per day arrays

### Admin Seed Management (Folder)

**Purpose:** Upload and manage team structures

**Requests:**
- **Get Seed Presets**
  - Method: GET
  - Path: `/admin/seed/presets`
  - Returns: Available presets (small-team, medium-team, enterprise, multi-region)

- **Upload Seed (JSON)**
  - Method: POST
  - Path: `/admin/seed`
  - Body: JSON seed data
  - Use when: You have a complete seed structure in JSON format

- **Upload Seed (CSV) + Regenerate**
  - Method: POST
  - Path: `/admin/seed`
  - Body: CSV data (user_id, email, name)
  - Auto-regenerates with override mode
  - Use when: You have developer list as CSV and want immediate generation

### Admin Data Regeneration (Folder)

**Purpose:** Scale data up/down or add more history

**Requests:**
- **Regenerate - Override Mode (Scale Up)**
  - Method: POST
  - Path: `/admin/regenerate`
  - Mode: override
  - Example: 1200 developers, 400 days, high velocity
  - Use when: Complete data replacement for large dataset

- **Regenerate - Override Mode (Small)**
  - Method: POST
  - Path: `/admin/regenerate`
  - Mode: override
  - Example: 50 developers, 30 days, low velocity
  - Use when: Fast testing with minimal data

- **Regenerate - Append Mode**
  - Method: POST
  - Path: `/admin/regenerate`
  - Mode: append
  - Example: Add 30 more days
  - Use when: Extending history without losing existing data

### Quick Workflows (Folder)

**Purpose:** Step-by-step workflow examples

**Requests:**
1. **[1] Choose Preset** - Get available seed presets
2. **[2] Verify Configuration** - Check current config after changes
3. **[3] Monitor Results** - Get statistics with time series

---

## Common Use Cases

### Scenario 1: Scale Up for Load Testing

**Goal:** Generate data for 1200 developers over 400 days

**Steps:**

1. Send: **Regenerate - Override Mode (Scale Up)**
   - Automatically configured for 1200 dev, 400 days, high velocity
   - Wait for response (will take 8-15 seconds)

2. Verify: **Get Current Configuration**
   - Check that `generation.developers = 1200` and `generation.days = 400`

3. Monitor: **Get Statistics (With Time Series)**
   - Check `generation.total_commits` and other metrics

### Scenario 2: Quick Test with Small Dataset

**Goal:** Generate small dataset for fast testing

**Steps:**

1. Send: **Regenerate - Override Mode (Small)**
   - Configured for 50 dev, 30 days, low velocity
   - Completes in 1-2 seconds

2. Verify: **Get Current Configuration**

3. Check: **Get Statistics (Basic)**

### Scenario 3: Change Team Structure

**Goal:** Upload new seed file with different team layout

**Steps:**

1. Get options: **Get Seed Presets**
   - See available preset configurations

2. Upload seed:
   - Either: **Upload Seed (JSON)** - if you have complete seed JSON
   - Or: **Upload Seed (CSV) + Regenerate** - if you have CSV dev list
   - This will load the new seed

3. (Optional) Regenerate: **Regenerate - Override Mode**
   - Use if you need to regenerate with new seed

4. Verify: **Get Current Configuration**
   - Check new seed structure loaded correctly

### Scenario 4: Add More Historical Data

**Goal:** Extend from 90 days to 120 days without losing data

**Steps:**

1. Send: **Regenerate - Append Mode**
   - Adds 30 more days to existing data
   - No data loss

2. Monitor: **Get Statistics (With Time Series)**
   - Verify commits_per_day now has more data points

---

## Response Examples

### Get Current Configuration Response

```json
{
  "generation": {
    "days": 90,
    "velocity": "medium",
    "developers": 50,
    "max_commits": 1000
  },
  "seed": {
    "version": "1.0",
    "developers": 50,
    "repositories": 5,
    "organizations": ["acme-corp"],
    "divisions": ["Engineering"],
    "teams": ["Backend", "Frontend"]
  },
  "external_sources": {
    "harvey": {"enabled": true},
    "copilot": {"enabled": true},
    "qualtrics": {"enabled": true}
  },
  "server": {
    "port": 8080,
    "version": "2.0.0",
    "uptime": "5m30s"
  }
}
```

### Regenerate Response

```json
{
  "status": "success",
  "mode": "override",
  "data_cleaned": true,
  "commits_added": 600000,
  "prs_added": 60000,
  "reviews_added": 120000,
  "issues_added": 10000,
  "total_commits": 600000,
  "total_prs": 60000,
  "total_developers": 1200,
  "duration": "8.5s",
  "config": {
    "days": 400,
    "velocity": "high",
    "developers": 1200,
    "max_commits": 500
  }
}
```

### Get Statistics Response

```json
{
  "generation": {
    "total_commits": 4500,
    "total_prs": 450,
    "total_reviews": 900,
    "total_issues": 150,
    "total_developers": 100,
    "data_size": "5.2 MB"
  },
  "developers": {
    "by_seniority": {"junior": 20, "mid": 50, "senior": 30},
    "by_region": {"US": 50, "EU": 30, "APAC": 20},
    "by_team": {"Backend": 40, "Frontend": 35, "DevOps": 25}
  },
  "quality": {
    "avg_revert_rate": 0.02,
    "avg_hotfix_rate": 0.08,
    "avg_code_survival_30d": 0.85,
    "avg_review_thoroughness": 0.75
  },
  "variance": {
    "commits_std_dev": 15.2,
    "pr_size_std_dev": 75.5
  },
  "performance": {
    "last_generation_time": "2.34s",
    "memory_usage": "125 MB"
  },
  "time_series": {
    "commits_per_day": [15, 18, 12, 20, 16],
    "prs_per_day": [3, 2, 4, 3, 2]
  }
}
```

---

## Tips & Tricks

### 1. Customize Request Bodies

All POST requests have editable bodies. Click the **Body** tab and modify:

- `days`: Adjust number of days (1-3650)
- `velocity`: Change to "low", "medium", or "high"
- `developers`: Set number of developers (0 = use seed count)
- `max_commits`: Limit commits per developer (0 = unlimited)

### 2. Use Environment Variables

Reference variables with `{{ _.variableName }}`

Current variables:
- `{{ _.baseUrl }}` - Base URL (default: http://localhost:8080)
- `{{ _.apiKey }}` - API key (default: cursor-sim-dev-key)

### 3. Inspect Response Details

Click **Preview** tab to see formatted JSON responses

### 4. Test Multiple Scenarios Sequentially

1. Save responses by clicking **Save Response** for comparison
2. Run different scenarios one after another
3. Use **Get Statistics** between runs to see changes

### 5. Monitor Long Operations

For large datasets (1000+ developers, 400+ days):
1. Start regeneration (**Regenerate - Override Mode (Scale Up)**)
2. Wait for response
3. Check **Get Statistics** multiple times to see progress
4. System will indicate when complete

---

## Authentication

All requests use **Basic Auth** with:
- **Username:** cursor-sim-dev-key
- **Password:** (empty)

This is pre-configured in all requests. No additional setup needed.

---

## Common Issues

### Issue: 401 Unauthorized

**Cause:** Wrong API key or missing authentication

**Fix:** Verify environment variable `apiKey` is set to `cursor-sim-dev-key`

### Issue: Connection Refused

**Cause:** cursor-sim not running on configured port

**Fix:**
1. Verify cursor-sim is running: `curl http://localhost:8080/health`
2. Check `baseUrl` environment variable matches running instance
3. Update if running on different port (e.g., `http://localhost:9000`)

### Issue: Timeout on Large Datasets

**Cause:** Regenerating 1000+ developers takes time

**Fix:**
1. Wait for operation to complete (up to 15 seconds)
2. Use smaller dataset for testing (small-team preset)
3. Check server logs for progress

---

## Additional Resources

- **Full API Documentation:** See `services/cursor-sim/SPEC.md` (Admin API section)
- **README Guide:** See `services/cursor-sim/README.md` (Admin API Suite section)
- **Code Reference:** See `services/cursor-sim/internal/api/cursor/admin_*.go`

---

## Support

For issues or questions:

1. Check the Admin API section in `SPEC.md`
2. Review the README examples
3. Check cursor-sim logs: `docker logs cursor-sim`
4. Verify health endpoint: `curl -u cursor-sim-dev-key: http://localhost:8080/health`
