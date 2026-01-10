# Task Breakdown: P4-F05 Insomnia Collections for External APIs

**Feature**: P4-F05 Insomnia Collections for External Data Sources
**Phase**: P4 (CLI Enhancements)
**Created**: January 10, 2026
**Status**: Planned
**Estimated Total**: 8-9 hours

---

## Implementation Strategy

This feature extends the existing Insomnia collection (`docs/insomnia/Insomnia_2026-01-09.yaml`) with comprehensive API testing for Harvey, Copilot, and Qualtrics endpoints. The work is organized into 3 sequential phases:

1. **Phase 1**: Create Insomnia collection folders and environment variables
2. **Phase 2**: Enhance E2E test coverage for edge cases
3. **Phase 3**: Update documentation (SPEC.md, usage guide)

**Key Insight**: No code changes to cursor-sim required. This is purely documentation and testing enhancement.

---

## Task Summary

### By Status

| Status | Count | Tasks |
|--------|-------|-------|
| PENDING | 8 | TASK-INS-01 through TASK-INS-08 |
| IN_PROGRESS | 0 | - |
| COMPLETE | 0 | - |

### By Phase

| Phase | Description | Tasks | Hours |
|-------|-------------|-------|-------|
| Phase 1 | Insomnia Collection Creation | 4 | 3-4h |
| Phase 2 | E2E Test Enhancement | 2 | 2-3h |
| Phase 3 | Documentation Updates | 2 | 1-2h |

---

## Phase 1: Insomnia Collection Creation (3-4 hours)

### TASK-INS-01: Create Harvey AI Folder (1h)

**Status**: PENDING
**Dependencies**: None
**Time Estimate**: 1.0 hours

**Goal**: Add Harvey AI API folder to Insomnia collection with complete endpoint configuration.

**Changes**:
- Add "Harvey AI" request_group to `docs/insomnia/Insomnia_2026-01-09.yaml`
- Create "GET Usage History" request with all query parameters
- Document response schema in request description
- Add example request with expected response

**Deliverables**:
- [ ] Harvey AI folder added with Basic Auth
- [ ] GET /harvey/api/v1/history/usage endpoint configured
- [ ] Query parameters: from, to, user, task, page, page_size
- [ ] Response schema documented with pagination structure
- [ ] Environment variables used: baseUrl, startDate, endDate, harveyUser, harveyTask

**Success Criteria**:
- [ ] Collection imports without errors
- [ ] Request executes successfully against running simulator
- [ ] Pagination controls work correctly
- [ ] Date and user/task filters function as expected

**YAML Structure**:
```yaml
_type: request_group
parentId: wrk_xxx
name: "Harvey AI"
authentication:
  type: basic
  username: "{{ _.apiKey }}"
  password: ""
requests:
  - _type: request
    name: "GET Usage History"
    method: GET
    url: "{{ _.baseUrl }}/harvey/api/v1/history/usage"
    parameters:
      - name: from
        value: "{{ _.startDate }}"
      - name: to
        value: "{{ _.endDate }}"
      - name: page
        value: "1"
      - name: page_size
        value: "50"
```

**Test Plan**:
1. Import collection into Insomnia
2. Configure environment with test API key
3. Execute GET Usage History request
4. Verify response matches documented schema
5. Test pagination by changing page parameter
6. Test date filtering with different ranges

---

### TASK-INS-02: Create Copilot Folder (1h)

**Status**: PENDING
**Dependencies**: None (can run parallel with TASK-INS-01)
**Time Estimate**: 1.0 hours

**Goal**: Add Microsoft 365 Copilot API folder with all period variants and CSV export.

**Changes**:
- Add "Microsoft 365 Copilot" request_group to Insomnia collection
- Create 4 period variant requests (D7, D30, D90, D180)
- Create CSV export request with $format parameter
- Document OData response structure

**Deliverables**:
- [ ] Copilot folder added with Basic Auth
- [ ] 4 period requests: D7, D30, D90, D180
- [ ] CSV export request with text/csv format parameter
- [ ] OData response schema documented
- [ ] App breakdown fields (Teams, Word, Outlook, Excel, PowerPoint, OneNote) documented

**Success Criteria**:
- [ ] All 5 requests execute successfully
- [ ] JSON responses include @odata.context
- [ ] CSV export returns proper CSV format with headers
- [ ] All periods return appropriate data ranges
- [ ] App-specific breakdowns present in responses

**YAML Structure**:
```yaml
_type: request_group
parentId: wrk_xxx
name: "Microsoft 365 Copilot"
authentication:
  type: basic
  username: "{{ _.apiKey }}"
  password: ""
requests:
  - name: "GET D7 Usage"
    url: "{{ _.baseUrl }}/reports/getMicrosoft365CopilotUsageUserDetail(period='D7')"
  - name: "GET D30 Usage"
    url: "{{ _.baseUrl }}/reports/getMicrosoft365CopilotUsageUserDetail(period='D30')"
  - name: "GET D90 Usage"
    url: "{{ _.baseUrl }}/reports/getMicrosoft365CopilotUsageUserDetail(period='D90')"
  - name: "GET D180 Usage"
    url: "{{ _.baseUrl }}/reports/getMicrosoft365CopilotUsageUserDetail(period='D180')"
  - name: "GET D30 CSV Export"
    url: "{{ _.baseUrl }}/reports/getMicrosoft365CopilotUsageUserDetail(period='D30')"
    parameters:
      - name: $format
        value: "text/csv"
```

**Test Plan**:
1. Import collection into Insomnia
2. Execute each period variant (D7, D30, D90, D180)
3. Verify response times scale appropriately
4. Execute CSV export request
5. Verify CSV format with proper headers and data rows
6. Validate OData structure in JSON responses

---

### TASK-INS-03: Create Qualtrics Folder (1.5h)

**Status**: PENDING
**Dependencies**: None (can run parallel with TASK-INS-01, TASK-INS-02)
**Time Estimate**: 1.5 hours

**Goal**: Add Qualtrics Survey Export API folder with complete 3-step workflow.

**Changes**:
- Add "Qualtrics Survey Export" request_group to Insomnia collection
- Create 3 workflow requests (Start Export, Check Progress, Download ZIP)
- Document multi-step workflow in folder description
- Add workflow instructions for copying IDs between steps

**Deliverables**:
- [ ] Qualtrics folder added with Basic Auth
- [ ] POST /API/v3/surveys/{surveyId}/export-responses endpoint
- [ ] GET /API/v3/surveys/{surveyId}/export-responses/{progressId} endpoint
- [ ] GET /API/v3/surveys/{surveyId}/export-responses/{fileId}/file endpoint
- [ ] Workflow documentation in folder description
- [ ] Environment variables: qualtricsSurveyId, qualtricsProgressId, qualtricsFileId

**Success Criteria**:
- [ ] Step 1 (Start Export) returns progressId
- [ ] Step 2 (Check Progress) shows state progression (queued → inProgress → complete)
- [ ] Step 3 (Download ZIP) returns ZIP file with proper Content-Type
- [ ] Workflow instructions are clear and easy to follow
- [ ] Environment variables can be manually updated between steps

**YAML Structure**:
```yaml
_type: request_group
parentId: wrk_xxx
name: "Qualtrics Survey Export"
description: |
  Three-step export workflow:
  1. POST to start export → receive progressId
  2. GET with progressId to poll status (queued → inProgress → complete)
  3. GET with fileId from step 2 to download ZIP

  Instructions:
  - Execute "1. Start Export"
  - Copy progressId from response to {{ _.qualtricsProgressId }}
  - Execute "2. Check Progress" repeatedly until status = "complete"
  - Copy fileId from response to {{ _.qualtricsFileId }}
  - Execute "3. Download ZIP"
authentication:
  type: basic
  username: "{{ _.apiKey }}"
  password: ""
requests:
  - name: "1. Start Export"
    method: POST
    url: "{{ _.baseUrl }}/API/v3/surveys/{{ _.qualtricsSurveyId }}/export-responses"
    body:
      mimeType: application/json
      text: '{"format": "json"}'
  - name: "2. Check Progress"
    method: GET
    url: "{{ _.baseUrl }}/API/v3/surveys/{{ _.qualtricsSurveyId }}/export-responses/{{ _.qualtricsProgressId }}"
  - name: "3. Download ZIP"
    method: GET
    url: "{{ _.baseUrl }}/API/v3/surveys/{{ _.qualtricsSurveyId }}/export-responses/{{ _.qualtricsFileId }}/file"
```

**Test Plan**:
1. Import collection into Insomnia
2. Configure qualtricsSurveyId environment variable
3. Execute "1. Start Export" request
4. Manually copy progressId to environment
5. Execute "2. Check Progress" multiple times
6. Verify state progression (queued → inProgress → complete)
7. Copy fileId to environment
8. Execute "3. Download ZIP"
9. Verify ZIP file downloads successfully

---

### TASK-INS-04: Add Environment Variables (0.5h)

**Status**: PENDING
**Dependencies**: TASK-INS-01, TASK-INS-02, TASK-INS-03
**Time Estimate**: 0.5 hours

**Goal**: Extend "Local Development" environment with all variables needed for external API testing.

**Changes**:
- Add 6 new environment variables to existing environment
- Document variable usage and example values
- Ensure backward compatibility with existing variables

**Deliverables**:
- [ ] harveyUser variable added (optional filter)
- [ ] harveyTask variable added (optional filter)
- [ ] copilotPeriod variable added (default period)
- [ ] qualtricsSurveyId variable added (survey ID)
- [ ] qualtricsProgressId variable added (workflow step 1 → step 2)
- [ ] qualtricsFileId variable added (workflow step 2 → step 3)

**Success Criteria**:
- [ ] All 6 variables present in environment
- [ ] Variables have sensible default/example values
- [ ] Documentation explains each variable's purpose
- [ ] Existing environment variables unchanged

**YAML Structure**:
```yaml
_type: environment
name: "Local Development"
data:
  # Existing variables
  baseUrl: "http://localhost:8080"
  apiKey: "dev-key-001"
  startDate: "2026-01-01"
  endDate: "2026-01-10"

  # Harvey AI variables (NEW)
  harveyUser: "developer@example.com"  # Optional: Filter by specific user
  harveyTask: "contract_review"         # Optional: Filter by task type

  # Copilot variables (NEW)
  copilotPeriod: "D30"                  # Default period for Copilot requests

  # Qualtrics variables (NEW)
  qualtricsSurveyId: "SV_aitools_q1_2026"  # Survey ID for export
  qualtricsProgressId: ""                   # Populated from step 1 response
  qualtricsFileId: ""                       # Populated from step 2 response
```

**Test Plan**:
1. Open Insomnia and navigate to Environments
2. Verify all 6 new variables present
3. Execute Harvey request with user/task filters
4. Execute Copilot request with period variable
5. Execute Qualtrics workflow with survey ID variable
6. Verify variables interpolate correctly in requests

---

## Phase 2: E2E Test Enhancement (2-3 hours)

### TASK-INS-05: Verify Existing E2E Coverage (1h)

**Status**: PENDING
**Dependencies**: None
**Time Estimate**: 1.0 hours

**Goal**: Run and document all 14 existing E2E tests to establish baseline coverage.

**Changes**:
- Run `go test ./test/e2e/external_data_test.go -v`
- Document test results and any failures
- Identify coverage gaps by comparing tests to API capabilities
- Create test coverage matrix

**Deliverables**:
- [ ] All 14 existing tests executed
- [ ] Test results documented (pass/fail/skip)
- [ ] Coverage gaps identified for each API
- [ ] Test coverage matrix created

**Success Criteria**:
- [ ] All 14 tests pass without modifications
- [ ] Test execution time < 30 seconds total
- [ ] Coverage gaps clearly documented
- [ ] No flaky tests detected

**Existing Tests** (from `test/e2e/external_data_test.go`):

**Harvey (4 tests)**:
1. TestHarvey_E2E_UsageEndpoint - Basic usage query
2. TestHarvey_E2E_Pagination - Page/page_size parameters
3. TestHarvey_E2E_DateFiltering - from/to date ranges
4. TestHarvey_E2E_DisabledWhenNotConfigured - Conditional registration

**Copilot (4 tests)**:
1. TestCopilot_E2E_JSONResponse - Default JSON format
2. TestCopilot_E2E_CSVExport - CSV export with $format
3. TestCopilot_E2E_AllPeriods - D7, D30, D90, D180
4. TestCopilot_E2E_DisabledWhenNotConfigured - Conditional registration

**Qualtrics (3 tests)**:
1. TestQualtrics_E2E_FullExportFlow - Complete 3-step workflow
2. TestQualtrics_E2E_ProgressAdvancement - State machine progression
3. TestQualtrics_E2E_DisabledWhenNotConfigured - Conditional registration

**Cross-cutting (3 tests)**:
1. TestExternalData_E2E_AuthenticationRequired - Basic Auth validation
2. TestExternalData_E2E_AllAPIsEnabled - All 3 APIs working together
3. TestExternalData_E2E_ErrorCases - Common error scenarios

**Coverage Gaps Identified**:

**Harvey Missing**:
- User-specific filtering (`?user=dev@example.com`)
- Task-specific filtering (`?task=contract_review`)
- Invalid date range handling (from > to)
- Pagination edge cases (last page, beyond total)

**Copilot Missing**:
- Invalid period value (`period='D45'`)
- Empty response handling (no users in period)
- CSV format validation (headers, data rows)

**Qualtrics Missing**:
- Export job timeout simulation
- Invalid survey ID handling
- Progress polling with immediate completion
- File download verification (Content-Type, size)

**Test Plan**:
1. Navigate to `services/cursor-sim/test/e2e/`
2. Run `go test ./external_data_test.go -v -count=1`
3. Document test output (pass/fail, timing)
4. Review test code to understand coverage
5. Compare test coverage to handler implementations
6. Create coverage matrix showing tested vs untested scenarios

---

### TASK-INS-06: Add Missing Test Scenarios (2h)

**Status**: PENDING
**Dependencies**: TASK-INS-05
**Time Estimate**: 2.0 hours

**Goal**: Implement 9 additional E2E test scenarios to achieve comprehensive coverage.

**Changes**:
- Add 3 Harvey test scenarios to `test/e2e/external_data_test.go`
- Add 3 Copilot test scenarios
- Add 3 Qualtrics test scenarios
- Ensure all response schemas validated

**Deliverables**:
- [ ] 3 Harvey tests added (user filtering, task filtering, invalid date range)
- [ ] 3 Copilot tests added (invalid period, empty response, CSV validation)
- [ ] 3 Qualtrics tests added (timeout, invalid survey ID, file verification)
- [ ] All new tests passing
- [ ] Total test count: 23 (14 existing + 9 new)

**Success Criteria**:
- [ ] All 23 tests pass
- [ ] Test execution time < 60 seconds total
- [ ] No code changes to handlers required
- [ ] All edge cases covered

**New Test Functions**:

**Harvey Enhancements**:
```go
func TestHarvey_E2E_UserFiltering(t *testing.T) {
    // Test ?user=dev@example.com parameter
    // Verify only events for specified user returned
}

func TestHarvey_E2E_TaskFiltering(t *testing.T) {
    // Test ?task=contract_review parameter
    // Verify only events for specified task returned
}

func TestHarvey_E2E_InvalidDateRange(t *testing.T) {
    // Test from > to (invalid range)
    // Verify appropriate error response
}
```

**Copilot Enhancements**:
```go
func TestCopilot_E2E_InvalidPeriod(t *testing.T) {
    // Test period='D45' (invalid value)
    // Verify 400 Bad Request with error message
}

func TestCopilot_E2E_EmptyResponse(t *testing.T) {
    // Test scenario with no users in period
    // Verify empty value array with valid @odata.context
}

func TestCopilot_E2E_CSVValidation(t *testing.T) {
    // Test CSV export format
    // Verify headers present and data rows formatted correctly
}
```

**Qualtrics Enhancements**:
```go
func TestQualtrics_E2E_ExportTimeout(t *testing.T) {
    // Test export job with extended inProgress state
    // Verify timeout handling
}

func TestQualtrics_E2E_InvalidSurveyID(t *testing.T) {
    // Test POST with non-existent survey ID
    // Verify 404 Not Found response
}

func TestQualtrics_E2E_FileDownloadVerification(t *testing.T) {
    // Test ZIP download
    // Verify Content-Type: application/zip
    // Verify file size > 0
}
```

**Test Plan**:
1. For each new test function:
   - Write test following existing patterns
   - Use httptest.NewRequest for setup
   - Assert response status code
   - Validate response body structure
   - Check error messages where applicable
2. Run individual test: `go test -run TestName -v`
3. Run all tests: `go test ./test/e2e/ -v`
4. Verify no regressions in existing tests
5. Document any handler issues discovered

---

## Phase 3: Documentation Updates (1-2 hours)

### TASK-INS-07: Update SPEC.md (1h)

**Status**: PENDING
**Dependencies**: TASK-INS-01, TASK-INS-02, TASK-INS-03
**Time Estimate**: 1.0 hours

**Goal**: Add comprehensive endpoint documentation for all 3 external data source APIs.

**Changes**:
- Add "External Data Sources (P4-F04)" section to `services/cursor-sim/SPEC.md`
- Document all 5 endpoints with full request/response examples
- Update endpoints summary table
- Add query parameters and validation rules

**Deliverables**:
- [ ] External Data Sources section added after Admin API
- [ ] 5 endpoints documented with examples
- [ ] Request/response schemas for all endpoints
- [ ] Query parameter documentation
- [ ] Error response examples
- [ ] Conditional registration documented

**Success Criteria**:
- [ ] All endpoints documented with consistent format
- [ ] Request examples include all parameters
- [ ] Response examples show actual data structure
- [ ] Error scenarios documented
- [ ] Pagination and OData structures explained

**SPEC.md Structure**:

```markdown
#### External Data Sources (P4-F04)

cursor-sim provides simulated endpoints for three external data source APIs. These endpoints are conditionally registered based on seed configuration.

**Endpoints Summary**:

| Method | Path | Auth | Description | Status |
|--------|------|------|-------------|--------|
| GET | `/harvey/api/v1/history/usage` | Yes | Harvey AI usage history | ✅ Implemented |
| GET | `/reports/getMicrosoft365CopilotUsageUserDetail(period='...')` | Yes | Copilot usage by period | ✅ Implemented |
| POST | `/API/v3/surveys/{surveyId}/export-responses` | Yes | Start Qualtrics export | ✅ Implemented |
| GET | `/API/v3/surveys/{surveyId}/export-responses/{progressId}` | Yes | Check export progress | ✅ Implemented |
| GET | `/API/v3/surveys/{surveyId}/export-responses/{fileId}/file` | Yes | Download export ZIP | ✅ Implemented |

##### Harvey AI API

**GET /harvey/api/v1/history/usage**

Query Harvey AI usage history with date filtering, pagination, and optional user/task filters.

**Query Parameters**:
- `from` (optional): Start date (RFC3339 format)
- `to` (optional): End date (RFC3339 format)
- `user` (optional): Filter by user email
- `task` (optional): Filter by task type
- `page` (optional): Page number (default: 1)
- `page_size` (optional): Results per page (default: 50, max: 100)

**Response**:
```json
{
  "events": [
    {
      "event_id": "evt_xxx",
      "user_email": "user@example.com",
      "timestamp": "2026-01-10T10:00:00Z",
      "task_type": "contract_review",
      "document_id": "doc_xxx",
      "feedback_score": 4,
      "time_saved_minutes": 15
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 50,
    "total": 1250,
    "has_next": true
  }
}
```

[Continue with Copilot and Qualtrics sections...]
```

**Test Plan**:
1. Open `services/cursor-sim/SPEC.md`
2. Add External Data Sources section after Admin API
3. Document each endpoint following existing format
4. Include request examples with all parameters
5. Include response examples with realistic data
6. Update "Last Updated" date
7. Verify markdown formatting with preview
8. Cross-reference with handler implementations

---

### TASK-INS-08: Create Insomnia Usage Guide (1h)

**Status**: PENDING
**Dependencies**: TASK-INS-04
**Time Estimate**: 1.0 hours

**Goal**: Create comprehensive user guide for Insomnia collection usage.

**Changes**:
- Create `docs/insomnia/README.md`
- Document import process
- Provide environment configuration guide
- Include workflow examples for each API
- Add troubleshooting section

**Deliverables**:
- [ ] docs/insomnia/README.md created
- [ ] Import instructions with screenshots (or detailed steps)
- [ ] Environment configuration guide
- [ ] 3 workflow examples (Harvey, Copilot, Qualtrics)
- [ ] Troubleshooting section with common issues

**Success Criteria**:
- [ ] New team member can import collection without assistance
- [ ] All workflows have step-by-step instructions
- [ ] Troubleshooting covers authentication, connectivity, data issues
- [ ] Links to SPEC.md for detailed API documentation

**README.md Structure**:

```markdown
# Insomnia REST Collection - cursor-sim

Comprehensive API testing collection for cursor-sim simulator.

## Quick Start

### 1. Import Collection

1. Open Insomnia
2. Click **Create** → **Import From File**
3. Select `docs/insomnia/Insomnia_2026-01-09.yaml`
4. Collection "cursor-sim API" will appear in sidebar

### 2. Configure Environment

1. Select **Local Development** environment from dropdown
2. Update environment variables:
   - `baseUrl`: http://localhost:8080 (or your simulator URL)
   - `apiKey`: Your API key (from seed data)
   - `startDate`: 2026-01-01 (or desired start date)
   - `endDate`: 2026-01-10 (or desired end date)

### 3. Test Connection

1. Open **Health Check** folder
2. Execute **GET /health** request
3. Verify response: `{"status": "healthy"}`

## API Collections

The collection includes 12 folders covering all cursor-sim APIs:

1. **Health Check** - Service health monitoring
2. **Quality Analysis** - Code quality metrics
3. **GitHub Analytics** - PR, review, issue tracking
4. **AI Code Tracking** - AI-assisted coding metrics
5. **Research Framework** - Experiment management
6. **Analytics Endpoints** - Team velocity, quality metrics
7. **Admin API Suite** - Configuration, regeneration, seed management
8. **Admin Statistics** - Usage statistics and monitoring
9. **Harvey AI** - Legal AI usage tracking ← NEW
10. **Microsoft 365 Copilot** - Copilot usage metrics ← NEW
11. **Qualtrics Survey Export** - Survey response export workflow ← NEW
12. **Authentication** - API key testing

## Workflow Examples

### Harvey AI: Querying Usage History

**Scenario**: Query Harvey AI usage for a specific user over the past 30 days.

**Steps**:
1. Navigate to **Harvey AI** folder
2. Open **GET Usage History** request
3. Update query parameters:
   - `from`: 2025-12-10
   - `to`: 2026-01-10
   - `user`: developer@example.com
   - `page`: 1
   - `page_size`: 50
4. Execute request
5. Review response:
   - `events` array contains usage events
   - `pagination` object shows total results and navigation

**Response Example**:
```json
{
  "events": [
    {
      "event_id": "evt_001",
      "user_email": "developer@example.com",
      "timestamp": "2026-01-09T14:30:00Z",
      "task_type": "contract_review",
      "document_id": "doc_12345",
      "feedback_score": 4,
      "time_saved_minutes": 15
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 50,
    "total": 127,
    "has_next": true
  }
}
```

### Copilot: Testing Different Periods

**Scenario**: Compare Copilot usage across different time periods.

**Steps**:
1. Navigate to **Microsoft 365 Copilot** folder
2. Execute requests in sequence:
   - **GET D7 Usage** (past 7 days)
   - **GET D30 Usage** (past 30 days)
   - **GET D90 Usage** (past 90 days)
   - **GET D180 Usage** (past 180 days)
3. Compare user counts and activity levels
4. For CSV export:
   - Execute **GET D30 CSV Export**
   - Verify response is CSV format with headers

**Response Example (JSON)**:
```json
{
  "@odata.context": "https://graph.microsoft.com/v1.0/$metadata#reports/...",
  "value": [
    {
      "userPrincipalName": "user@example.com",
      "lastActivityDate": "2026-01-10",
      "copilotActivity": {
        "teams": 50,
        "word": 30,
        "outlook": 20,
        "excel": 15,
        "powerpoint": 10,
        "onenote": 5
      }
    }
  ]
}
```

### Qualtrics: Complete Export Workflow

**Scenario**: Export survey responses as JSON and download ZIP file.

**Steps**:

**Step 1: Start Export**
1. Navigate to **Qualtrics Survey Export** folder
2. Update environment variable `qualtricsSurveyId` to your survey ID (e.g., "SV_aitools_q1_2026")
3. Execute **1. Start Export** request
4. Copy `progressId` from response
5. Update environment variable `qualtricsProgressId` with copied value

**Step 2: Check Progress**
1. Execute **2. Check Progress** request
2. Check `status` field in response:
   - `queued`: Export queued, wait and retry
   - `inProgress`: Export processing, wait and retry
   - `complete`: Export ready, proceed to step 3
3. Copy `fileId` from response
4. Update environment variable `qualtricsFileId` with copied value

**Step 3: Download ZIP**
1. Execute **3. Download ZIP** request
2. Verify response:
   - Content-Type: application/zip
   - File size > 0 bytes
3. Save response body as .zip file
4. Extract and review JSON export

**Response Examples**:

Step 1 Response:
```json
{
  "progressId": "prog_abc123",
  "status": "queued"
}
```

Step 2 Response:
```json
{
  "progressId": "prog_abc123",
  "status": "complete",
  "fileId": "file_xyz789"
}
```

Step 3 Response: Binary ZIP file

## Troubleshooting

### Authentication Errors

**Problem**: 401 Unauthorized response

**Solutions**:
- Verify `apiKey` environment variable matches seed data
- Check Basic Auth configuration (username: apiKey, password: empty)
- Ensure API key hasn't expired or been rotated

### Connection Refused

**Problem**: Failed to connect to localhost:8080

**Solutions**:
- Verify cursor-sim is running: `curl http://localhost:8080/health`
- Check `baseUrl` environment variable
- Ensure no firewall blocking port 8080

### Empty Responses

**Problem**: API returns empty arrays or no data

**Solutions**:
- Verify seed data includes external data sources:
  ```json
  "external_data_sources": {
    "harvey": {"enabled": true},
    "copilot": {"enabled": true},
    "qualtrics": {"enabled": true}
  }
  ```
- Check date range parameters (data may not exist for that period)
- Verify data generation completed successfully

### Invalid Parameters

**Problem**: 400 Bad Request with "invalid parameter" error

**Solutions**:
- Check query parameter format (dates must be RFC3339)
- Verify period values (Copilot: D7, D30, D90, D180 only)
- Ensure pagination values are positive integers
- Check survey ID format (Qualtrics)

## Additional Resources

- **API Documentation**: [cursor-sim SPEC.md](../../services/cursor-sim/SPEC.md)
- **E2E Tests**: [external_data_test.go](../../services/cursor-sim/test/e2e/external_data_test.go)
- **Handler Implementations**:
  - Harvey: `internal/api/harvey/handlers.go`
  - Copilot: `internal/api/microsoft/copilot_handlers.go`
  - Qualtrics: `internal/api/qualtrics/handlers.go`

## Support

For issues or questions:
1. Check SPEC.md for API details
2. Review E2E tests for example usage
3. Verify seed configuration includes external data sources
4. Check simulator logs for error details
```

**Test Plan**:
1. Follow import instructions step-by-step
2. Configure environment variables
3. Execute each workflow example
4. Verify all steps work as documented
5. Test troubleshooting scenarios
6. Ask colleague to review for clarity

---

## Dependencies and Sequencing

### Can Run in Parallel:
- TASK-INS-01, TASK-INS-02, TASK-INS-03 (Insomnia folder creation)
- TASK-INS-05 (E2E verification) can start immediately

### Must Run Sequentially:
- TASK-INS-04 depends on completion of TASK-INS-01, TASK-INS-02, TASK-INS-03
- TASK-INS-06 depends on TASK-INS-05
- TASK-INS-07 depends on TASK-INS-01, TASK-INS-02, TASK-INS-03
- TASK-INS-08 depends on TASK-INS-04

### Suggested Execution Order:

**Parallel Track 1** (Insomnia Collection):
1. TASK-INS-01 (Harvey folder)
2. TASK-INS-02 (Copilot folder)
3. TASK-INS-03 (Qualtrics folder)
4. TASK-INS-04 (Environment variables)
5. TASK-INS-08 (Usage guide)

**Parallel Track 2** (Testing):
1. TASK-INS-05 (Verify existing coverage)
2. TASK-INS-06 (Add missing scenarios)

**Final** (Documentation):
1. TASK-INS-07 (Update SPEC.md)

**Total Time**: 6-9 hours with parallelization vs 8-9 hours sequential

---

## Progress Tracking

Update this section as tasks complete:

| Task | Estimated | Actual | Status | Completion Date |
|------|-----------|--------|--------|-----------------|
| TASK-INS-01 | 1.0h | - | PENDING | - |
| TASK-INS-02 | 1.0h | - | PENDING | - |
| TASK-INS-03 | 1.5h | - | PENDING | - |
| TASK-INS-04 | 0.5h | - | PENDING | - |
| TASK-INS-05 | 1.0h | - | PENDING | - |
| TASK-INS-06 | 2.0h | - | PENDING | - |
| TASK-INS-07 | 1.0h | - | PENDING | - |
| TASK-INS-08 | 1.0h | - | PENDING | - |
| **TOTAL** | **8-9h** | **-** | **0% Complete** | - |

---

## Completion Checklist

Before marking P4-F05 as COMPLETE, verify:

### Insomnia Collection
- [ ] Harvey AI folder with 1 endpoint (GET Usage History)
- [ ] Copilot folder with 5 endpoints (4 periods + CSV)
- [ ] Qualtrics folder with 3 endpoints (3-step workflow)
- [ ] 6 environment variables added to Local Development
- [ ] Collection imports without errors
- [ ] All requests execute successfully

### E2E Tests
- [ ] All 14 existing tests passing
- [ ] 9 new test scenarios added (3 per API)
- [ ] Total 23 tests passing
- [ ] Test execution time < 60 seconds
- [ ] No code changes required to handlers

### Documentation
- [ ] SPEC.md updated with External Data Sources section
- [ ] All 5 endpoints documented with examples
- [ ] docs/insomnia/README.md created
- [ ] Import instructions clear and complete
- [ ] 3 workflow examples documented
- [ ] Troubleshooting section comprehensive

### Quality Gates
- [ ] No regressions in existing tests
- [ ] No new linter warnings
- [ ] All commits follow SDD workflow
- [ ] DEVELOPMENT.md updated
- [ ] Zero onboarding questions from team (success metric)

---

## Risk Mitigation

| Risk | Mitigation |
|------|------------|
| Insomnia YAML format incompatible | Follow existing pattern exactly, validate with import test |
| E2E tests flaky | Use proper test server setup, avoid hardcoded delays |
| Documentation drift | Reference SPEC.md as source of truth, update both |
| Missing test scenarios | Comprehensive coverage analysis in TASK-INS-05 |
| Environment variable conflicts | Use descriptive names, document purpose |

---

## Notes

- **No code changes to cursor-sim**: APIs already fully functional from P4-F04
- **No admin API changes**: Seed configuration already supports external sources
- **Low risk**: Declarative YAML and test additions only
- **High value**: Reduces manual testing time, improves onboarding
- **Documentation-heavy**: Focus on clarity and examples
