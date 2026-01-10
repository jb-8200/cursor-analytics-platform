# cursor-sim REST API Collection for Insomnia

Comprehensive REST API testing collection for cursor-sim simulator, including Admin APIs, Analytics APIs, and External Data Sources (Harvey AI, Microsoft 365 Copilot, Qualtrics).

**Compatible with**: Insomnia v2024.5.0+
**Collection Files**:
- `Insomnia_2026-01-09.yaml` - Complete API collection (all endpoints)
- `Admin_APIs_2026-01-10.yaml` - Admin APIs only (quick import)

---

## Quick Start

### 1. Import Collection

**Option A: Full Collection**
1. Open Insomnia
2. Click **Dashboard** → **Import**
3. Select `docs/insomnia/Insomnia_2026-01-09.yaml`
4. Collection "cursor-sim API" appears in left sidebar

**Option B: Admin APIs Only**
1. Open Insomnia
2. Click **Dashboard** → **Import**
3. Select `docs/insomnia/Admin_APIs_2026-01-10.yaml`
4. Collection "cursor-sim Admin APIs" appears in left sidebar

### 2. Configure Environment

1. Click **Manage Environments** (bottom left corner)
2. Select or create "Local Development" environment
3. Edit environment with these variables:

```json
{
  "baseUrl": "http://localhost:8080",
  "apiKey": "cursor-sim-dev-key",
  "startDate": "2025-12-02",
  "endDate": "2026-01-06",
  "repoOwner": "acme-corp",
  "repoName": "payment-service",
  "prNumber": "1",
  "commitSha": "abc123",
  "harveyUser": "user@example.com",
  "harveyTask": "legal-review",
  "copilotPeriod": "D30",
  "qualtricsSurveyId": "SV_XXXXXXXXXXXXX",
  "qualtricsProgressId": "PROGRESS_ID_123456",
  "qualtricsFileId": "FILE_ID_123456"
}
```

### 3. Test Connection

1. Open **Health Check** folder
2. Click on **GET /health** request
3. Press **Send** button
4. Response: `{"status":"ok"}`

---

## Collection Structure

### Core API Folders

| Folder | Purpose | Endpoints |
|--------|---------|-----------|
| **Health Check** | Service monitoring | 1 |
| **Admin Configuration** | Runtime settings | 1 |
| **Admin Statistics** | Data metrics | 2 |
| **Admin Seed Management** | Seed file upload | 3 |
| **Admin Data Regeneration** | Data generation control | 3 |
| **Quick Workflows** | Common admin tasks | 3 |

### External Data Sources (NEW)

| Folder | Purpose | Endpoints |
|--------|---------|-----------|
| **Harvey AI** | Legal AI usage tracking | 1 |
| **Microsoft 365 Copilot** | Copilot usage metrics | 5 (periods + CSV) |
| **Qualtrics Survey Export** | Survey response export | 3 (start → poll → download) |

### Analytics & Research

| Folder | Purpose | Endpoints |
|--------|---------|-----------|
| **Team Management** | Developer lists | 1 |
| **AI Code Tracking** | AI-assisted development metrics | 2 |
| **Team Analytics** | Team-wide metrics | 11 |
| **By-User Analytics** | Per-developer metrics | 9 |
| **GitHub Analytics** | PR, review, issue tracking | 5 |

---

## Workflow Examples

### Workflow 1: Scale Up Server for Load Testing

**Goal**: Generate large dataset (1200 developers, 400 days)

**Steps**:
1. Navigate to **Admin Data Regeneration** folder
2. Open **Regenerate - Override Mode (Scale Up)**
3. Click **Send** (will take 8-15 seconds)
4. Wait for completion status
5. Verify with **Get Current Configuration**
6. Monitor with **Get Statistics (With Time Series)**

**Use Case**: Load testing, performance validation

---

### Workflow 2: Test External Data Source - Harvey AI

**Goal**: Query Harvey AI legal document analysis usage

**Steps**:
1. Ensure cursor-sim is running with seed that has Harvey enabled
2. Navigate to **Harvey AI** folder
3. Click **GET Usage History**
4. Verify parameters in URL:
   - `from`: Start date (auto-populated from environment)
   - `to`: End date (auto-populated from environment)
   - `user`: Optional (delete if not filtering)
   - `task`: Optional (delete if not filtering)
5. Click **Send**
6. Response includes usage events with pagination

**Response Example**:
```json
{
  "data": [
    {
      "id": "event_001",
      "timestamp": "2026-01-10T15:30:00Z",
      "user_email": "dev@example.com",
      "task_type": "legal_review",
      "document_name": "contract.pdf",
      "duration_seconds": 180
    }
  ],
  "pagination": {
    "page": 1,
    "pageSize": 50,
    "totalCount": 156,
    "totalPages": 4
  }
}
```

**Use Case**: Tracking legal AI usage across organization

---

### Workflow 3: Export Copilot Usage Metrics

**Goal**: Export Microsoft 365 Copilot usage in multiple formats

**Steps**:

**3a. Get JSON Response**
1. Navigate to **Microsoft 365 Copilot** folder
2. Open **D30 Usage (JSON)**
3. Click **Send**
4. Verify `@odata.context` and `value` fields
5. Each user record includes: reportRefreshDate, userPrincipalName, displayName, metrics

**3b. Export as CSV**
1. Open **D30 Usage (CSV)**
2. Click **Send**
3. Response headers include `Content-Disposition: attachment`
4. CSV format: headers + data rows
5. Can download directly

**3c. Compare Multiple Periods**
1. Try different period values: D7, D30, D90, D180
2. Each period returns different dataset sizes
3. reportPeriod field indicates the period

**Use Case**: Reporting Copilot adoption, usage trends

---

### Workflow 4: Complete Qualtrics Survey Export

**Goal**: Export survey responses through 3-step workflow

**Step 1: Start Export Job**
1. Navigate to **Qualtrics Survey Export** folder
2. Open **1. Start Export**
3. Update `qualtricsSurveyId` in environment (if needed)
4. Click **Send**
5. Response includes `progressId` (save this)

**Step 2: Poll Progress**
1. Open **2. Check Progress**
2. Check environment variables:
   - `qualtricsSurveyId`: From step 1 response
   - `qualtricsProgressId`: From step 1 response
3. Click **Send** repeatedly (5-10 times, every 1 second)
4. Watch `percentComplete` increase: 0% → 20% → 40% → 60% → 80% → 100%
5. When status = "complete", note the `fileId`

**Step 3: Download ZIP File**
1. Open **3. Download ZIP**
2. Update `qualtricsFileId` with fileId from step 2
3. Click **Send**
4. Response headers: `Content-Type: application/zip`
5. Binary ZIP file contains survey_responses.csv

**CSV Contents**:
```csv
ResponseID,RespondentEmail,OverallAISatisfaction,RecommendedFeatures,AdditionalComments
RESP_001,user1@company.com,9,Code generation,Great tool
RESP_002,user2@company.com,8,Debugging,Works well
```

**Use Case**: Collecting feedback on AI tools

---

## Configuration Reference

### Environment Variables

All variables use Insomnia interpolation syntax: `{{ _.variableName }}`

**Base Configuration**:
- `baseUrl` - Simulator URL (http://localhost:8080)
- `apiKey` - Basic Auth username (cursor-sim-dev-key)

**Date Range**:
- `startDate` - Query start (YYYY-MM-DD format)
- `endDate` - Query end (YYYY-MM-DD format)

**Repository/PR Context**:
- `repoOwner` - GitHub organization
- `repoName` - Repository name
- `prNumber` - Pull request number
- `commitSha` - Commit hash

**External Data Source Parameters**:
- `harveyUser` - Harvey user email for filtering
- `harveyTask` - Harvey task type (legal_review, contract_analysis, etc.)
- `copilotPeriod` - Copilot report period (D7, D30, D90, D180)
- `qualtricsSurveyId` - Qualtrics survey ID (SV_...)
- `qualtricsProgressId` - Export progress ID (ES_...)
- `qualtricsFileId` - Downloaded file ID (FILE_...)

### Authentication

All requests use **Basic Auth**:
- Username: Value of `apiKey` environment variable
- Password: (empty)

Pre-configured in all request headers. No additional setup required.

---

## Common Workflows

### Workflow A: Start Small, Scale Up

1. **Start**: Load small dataset (50 devs, 30 days)
   - Open **Regenerate - Override Mode (Small)**
   - Complete in 1-2 seconds

2. **Verify**: Check configuration
   - Open **Get Current Configuration**
   - Confirm: developers=50, days=30

3. **Test**: Query endpoints
   - Try **AI Code Tracking** endpoints
   - Try **GitHub Analytics** endpoints

4. **Scale**: Gradually increase
   - **Regenerate - Override Mode (Scale Up)**
   - 1200 developers, 400 days
   - Takes 8-15 seconds

### Workflow B: Monitor Data Quality

1. **Regenerate** data with desired parameters
2. **Check Statistics**:
   - Basic: Total counts and quality metrics
   - With Time Series: Trends over time
3. **Verify Consistency**:
   - Compare across time ranges
   - Check for anomalies

### Workflow C: Test External Data APIs

1. **Verify Seed**: Ensure external sources enabled
   - Run **Get Current Configuration**
   - Check `external_sources` section

2. **Test Each API**:
   - Harvey: **GET Usage History**
   - Copilot: **D30 Usage (JSON)**
   - Qualtrics: Full 3-step export workflow

3. **Validate Responses**:
   - Check status codes (200 OK)
   - Verify JSON structure
   - Download CSV/ZIP files

---

## Troubleshooting

### Connection Issues

**Problem**: "Connection refused" on all requests

**Solution**:
1. Verify cursor-sim is running: `curl http://localhost:8080/health`
2. Check `baseUrl` in environment matches running instance
3. If running on different port: Update `baseUrl` to `http://localhost:PORT`
4. Firewall: Ensure port 8080 is accessible

**Problem**: Request timeout

**Solution**:
1. Large datasets may take 8-15 seconds
2. Increase Insomnia timeout: Settings → Request → Timeout (seconds)
3. For regeneration: D30 = ~5s, D90 = ~10s, D400 = ~15s

---

### Authentication Issues

**Problem**: "401 Unauthorized" on all requests

**Solution**:
1. Verify `apiKey` in environment (should be `cursor-sim-dev-key`)
2. All requests include Basic Auth header
3. Check cursor-sim running with correct API key
4. Re-select environment to refresh headers

**Problem**: Different API key needed

**Solution**:
1. Update environment `apiKey` variable
2. All requests use this value automatically
3. Re-send request to apply changes

---

### External Data Source Issues

**Problem**: External API returns 404

**Solution**:
1. Verify external source is enabled in seed:
   - Run **Get Current Configuration**
   - Check `external_sources` → `{api}` → `enabled: true`
2. If disabled, regenerate with seed that enables it
3. Restart cursor-sim after seed changes

**Problem**: Harvey returns empty data

**Solution**:
1. Remove optional filters: `user`, `task`
2. Expand date range: Use wider date range if available
3. Check generation parameters in stats

**Problem**: Copilot returns "invalid period"

**Solution**:
1. Only valid periods: D7, D30, D90, D180
2. Don't use other values like D45, D60
3. Verify period in URL is correct

**Problem**: Qualtrics export stuck in "inProgress"

**Solution**:
1. Wait longer (exports can take 5-10 seconds)
2. Keep polling progress endpoint
3. Status will eventually change to "complete" or "failed"
4. If still stuck after 30 seconds: Check server logs

---

### Data Issues

**Problem**: Data seems inconsistent

**Solution**:
1. Check statistics: **Get Statistics (With Time Series)**
2. Verify generation parameters match expectation
3. Check date range is correct
4. Regenerate with specific parameters if needed

**Problem**: CSV export has no data

**Solution**:
1. Verify CSV format request: `$format=text/csv`
2. Check HTTP response code (200 OK)
3. Save CSV file and inspect locally
4. Empty CSV usually means empty dataset

---

## Tips & Tricks

### 1. Save Responses for Comparison

1. Click **Save Response** after each request
2. Insomnia stores response alongside request
3. Switch between saved versions to compare
4. Useful for: Before/after regeneration, debugging changes

### 2. Use Request Chains

Some workflows benefit from ordering:
1. **Qualtrics workflow**: Must follow order (start → poll → download)
2. **Admin workflow**: Often: Regenerate → Verify → Stats
3. Organize requests in folders to group related calls

### 3. Customize Request Bodies

For POST requests (seed upload, regeneration):
1. Click **Body** tab
2. Edit JSON directly
3. Parameters are documented in each request

### 4. Export Collection

To backup or share:
1. Right-click collection name
2. Select **Export**
3. Save as `.yaml` or `.json`
4. Share or version control

### 5. Monitor Performance

For load testing:
1. Open **Get Statistics**
2. Note `last_generation_time`
3. For large datasets: May see 10-15s generation time
4. Memory usage in stats helps capacity planning

---

## Links & References

- **Full API Documentation**: `services/cursor-sim/SPEC.md`
- **Admin API Guide**: `docs/insomnia/ADMIN_API_GUIDE.md`
- **External Data Sources**: See SPEC.md "External Data Sources API" section
- **Seed Configuration**: `services/cursor-sim/testdata/valid_seed.json`

---

## Support

### Getting Help

1. Check **Troubleshooting** section above (covers 90% of issues)
2. Review relevant API documentation in SPEC.md
3. Check cursor-sim server logs: `docker logs cursor-sim`
4. Verify health endpoint: `curl -u cursor-sim-dev-key: http://localhost:8080/health`

### Common Questions

**Q: Can I use this collection with production cursor-sim?**
A: Yes, but change `apiKey` to production key. Recommend separate environment for prod.

**Q: How do I add custom requests?**
A: Right-click folder → Add Request, configure method/path/auth, save.

**Q: Can I automate testing with Insomnia?**
A: Yes, use Insomnia CLI or export to OpenAPI for integration tests.

**Q: What's the maximum dataset size?**
A: Tested up to 10,000 developers, 3650 days. Performance degrades beyond this.

---

## Version History

| Date | Version | Changes |
|------|---------|---------|
| 2026-01-10 | 2.0 | Added External Data Sources (Harvey, Copilot, Qualtrics) |
| 2026-01-10 | 1.5 | Admin API Suite complete (config, stats, regenerate, seed) |
| 2026-01-02 | 1.0 | Initial collection with core Analytics APIs |

---

**Last Updated**: January 10, 2026
**Collection Status**: Complete
**Total Endpoints**: 36 (Admin: 5, Analytics: 29, External: 5, Health: 1)
