# Technical Design: P4-F05 Insomnia Collections for External APIs

**Feature**: P4-F05 Insomnia Collections for External Data Sources
**Phase**: P4 (CLI Enhancements)
**Created**: January 10, 2026
**Status**: Planned

---

## Architecture Overview

This feature adds comprehensive Insomnia REST API collections and enhanced E2E test coverage for three external data source APIs that are already implemented in P4-F04:

1. **Harvey AI API** - Legal AI platform integration
2. **Microsoft 365 Copilot API** - GitHub Copilot metrics
3. **Qualtrics Survey API** - Survey response export

**Key Insight**: No code changes required to cursor-sim. This is purely documentation and testing enhancement.

---

## Current State Analysis

### Existing Implementation (P4-F04)

**Harvey AI API** (`/harvey/api/v1/*`)
- Handler: `internal/api/harvey/handlers.go`
- Endpoint: GET `/harvey/api/v1/history/usage`
- Query params: from, to, user, task, page, page_size
- Response: Paginated usage events with feedback and document tracking
- Conditional registration: Only if `seed.ExternalDataSources.Harvey.Enabled == true`

**Copilot Usage API** (`/microsoft/copilot/*`)
- Handler: `internal/api/microsoft/copilot_handlers.go`
- Endpoint: GET `/reports/getMicrosoft365CopilotUsageUserDetail(period='D30')`
- Periods: D7, D30, D90, D180
- Formats: JSON (default), CSV ($format=text/csv)
- Response: OData-style with @odata.context, value array, app breakdowns
- Conditional registration: Only if `seed.ExternalDataSources.Copilot.Enabled == true`

**Qualtrics Export API** (`/qualtrics/*`)
- Handler: `internal/api/qualtrics/handlers.go`
- Endpoints (3-step workflow):
  1. POST `/API/v3/surveys/{surveyId}/export-responses` - Start export
  2. GET `/API/v3/surveys/{surveyId}/export-responses/{progressId}` - Poll progress
  3. GET `/API/v3/surveys/{surveyId}/export-responses/{fileId}/file` - Download ZIP
- State machine: queued → inProgress → complete
- Conditional registration: Only if `seed.ExternalDataSources.Qualtrics.Enabled == true`

### Existing Test Coverage

File: `test/e2e/external_data_test.go` (14 tests)

**Harvey (4 tests)**:
- TestHarvey_E2E_UsageEndpoint
- TestHarvey_E2E_Pagination
- TestHarvey_E2E_DateFiltering
- TestHarvey_E2E_DisabledWhenNotConfigured

**Copilot (4 tests)**:
- TestCopilot_E2E_JSONResponse
- TestCopilot_E2E_CSVExport
- TestCopilot_E2E_AllPeriods
- TestCopilot_E2E_DisabledWhenNotConfigured

**Qualtrics (3 tests)**:
- TestQualtrics_E2E_FullExportFlow
- TestQualtrics_E2E_ProgressAdvancement
- TestQualtrics_E2E_DisabledWhenNotConfigured

**Cross-cutting (3 tests)**:
- TestExternalData_E2E_AuthenticationRequired
- TestExternalData_E2E_AllAPIsEnabled
- TestExternalData_E2E_ErrorCases

### Existing Insomnia Collection

File: `docs/insomnia/Insomnia_2026-01-09.yaml`

**Structure**:
- 9 existing folders (Health Check, Quality Analysis, GitHub, AI Code, Research, Analytics, etc.)
- Base URL environment variable: `{{ _.baseUrl }}`
- Authentication: Basic Auth with apiKey as username, empty password
- Date parameters: `{{ _.startDate }}`, `{{ _.endDate }}`

**Pattern to follow**:
```yaml
_type: request_group
parentId: wrk_xxx
name: "Folder Name"
authentication:
  type: basic
  username: "{{ _.apiKey }}"
  password: ""
```

---

## Technical Design

### Component 1: Insomnia Collection Extensions

#### Harvey AI Folder Structure

```yaml
- Folder: "Harvey AI" (parentId: wrk_xxx)
  - Authentication: Basic Auth
  - Request: "GET Usage History"
    - URL: {{ _.baseUrl }}/harvey/api/v1/history/usage
    - Params:
      - from: {{ _.startDate }}
      - to: {{ _.endDate }}
      - user: {{ _.harveyUser }} (optional)
      - task: {{ _.harveyTask }} (optional)
      - page: 1
      - page_size: 50
    - Description: "Query Harvey AI usage history with date filtering and pagination"
```

**Response Schema**:
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

#### Copilot Folder Structure

```yaml
- Folder: "Microsoft 365 Copilot" (parentId: wrk_xxx)
  - Authentication: Basic Auth
  - Request: "GET D7 Usage"
    - URL: {{ _.baseUrl }}/reports/getMicrosoft365CopilotUsageUserDetail(period='D7')
  - Request: "GET D30 Usage"
    - URL: {{ _.baseUrl }}/reports/getMicrosoft365CopilotUsageUserDetail(period='D30')
  - Request: "GET D90 Usage"
    - URL: {{ _.baseUrl }}/reports/getMicrosoft365CopilotUsageUserDetail(period='D90')
  - Request: "GET D180 Usage"
    - URL: {{ _.baseUrl }}/reports/getMicrosoft365CopilotUsageUserDetail(period='D180')
  - Request: "GET D30 CSV Export"
    - URL: {{ _.baseUrl }}/reports/getMicrosoft365CopilotUsageUserDetail(period='D30')
    - Params:
      - $format: text/csv
```

**Response Schema (JSON)**:
```json
{
  "@odata.context": "...",
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

#### Qualtrics Folder Structure

```yaml
- Folder: "Qualtrics Survey Export" (parentId: wrk_xxx)
  - Authentication: Basic Auth
  - Request: "1. Start Export"
    - Method: POST
    - URL: {{ _.baseUrl }}/API/v3/surveys/{{ _.qualtricsSurveyId }}/export-responses
    - Body: {"format": "json"}
  - Request: "2. Check Progress"
    - Method: GET
    - URL: {{ _.baseUrl }}/API/v3/surveys/{{ _.qualtricsSurveyId }}/export-responses/{{ _.qualtricsProgressId }}
  - Request: "3. Download ZIP"
    - Method: GET
    - URL: {{ _.baseUrl }}/API/v3/surveys/{{ _.qualtricsSurveyId }}/export-responses/{{ _.qualtricsFileId }}/file
```

**Workflow Description** (in folder):
```
Three-step export workflow:
1. POST to start export → receive progressId
2. GET with progressId to poll status (queued → inProgress → complete)
3. GET with fileId from step 2 to download ZIP

Copy progressId from step 1 response to {{ _.qualtricsProgressId }}
Copy fileId from step 2 response to {{ _.qualtricsFileId }}
```

#### Environment Variables to Add

Extend "Local Development" environment:

```yaml
harveyUser: "developer@example.com"  # Optional: Filter by specific user
harveyTask: "contract_review"         # Optional: Filter by task type
copilotPeriod: "D30"                  # Default period for Copilot
qualtricsSurveyId: "SV_aitools_q1_2026"  # Survey ID for export
qualtricsProgressId: ""               # Populated from step 1 response
qualtricsFileId: ""                   # Populated from step 2 response
```

---

### Component 2: E2E Test Enhancements

#### Test Coverage Gaps

**Harvey - Missing Scenarios**:
1. User-specific filtering (`?user=dev@example.com`)
2. Task-specific filtering (`?task=contract_review`)
3. Invalid date range (from > to)
4. Pagination edge cases (last page, page beyond total)

**Copilot - Missing Scenarios**:
1. Invalid period value (`period='D45'`)
2. Empty response handling (no users in period)
3. CSV format validation (proper headers, data rows)

**Qualtrics - Missing Scenarios**:
1. Export job timeout simulation
2. Invalid survey ID handling
3. Progress polling with immediate completion
4. File download verification (Content-Type, file size)

#### New Test Structure

Add to `test/e2e/external_data_test.go`:

```go
// Harvey enhancements
func TestHarvey_E2E_UserFiltering(t *testing.T) { ... }
func TestHarvey_E2E_TaskFiltering(t *testing.T) { ... }
func TestHarvey_E2E_InvalidDateRange(t *testing.T) { ... }

// Copilot enhancements
func TestCopilot_E2E_InvalidPeriod(t *testing.T) { ... }
func TestCopilot_E2E_EmptyResponse(t *testing.T) { ... }
func TestCopilot_E2E_CSVValidation(t *testing.T) { ... }

// Qualtrics enhancements
func TestQualtrics_E2E_ExportTimeout(t *testing.T) { ... }
func TestQualtrics_E2E_InvalidSurveyID(t *testing.T) { ... }
func TestQualtrics_E2E_FileDownloadVerification(t *testing.T) { ... }
```

---

### Component 3: Documentation Updates

#### SPEC.md Additions

Add to `services/cursor-sim/SPEC.md` after Admin API section:

```markdown
#### External Data Sources (P4-F04)

| Method | Path | Auth | Status |
|--------|------|------|--------|
| GET | `/harvey/api/v1/history/usage` | Yes | ✅ Implemented |
| GET | `/reports/getMicrosoft365CopilotUsageUserDetail(period='...')` | Yes | ✅ Implemented |
| POST | `/API/v3/surveys/{surveyId}/export-responses` | Yes | ✅ Implemented |
| GET | `/API/v3/surveys/{surveyId}/export-responses/{progressId}` | Yes | ✅ Implemented |
| GET | `/API/v3/surveys/{surveyId}/export-responses/{fileId}/file` | Yes | ✅ Implemented |

[Full endpoint documentation with request/response examples]
```

#### Insomnia Usage Guide

Create `docs/insomnia/README.md`:

**Contents**:
1. **Import Instructions**: How to import `Insomnia_2026-01-09.yaml`
2. **Environment Configuration**: Setting up Local Development environment
3. **Workflow Examples**:
   - Harvey: Querying usage with filters
   - Copilot: Testing different periods and formats
   - Qualtrics: Complete 3-step export workflow
4. **Troubleshooting**: Common issues and solutions

---

## Implementation Plan

### Phase 1: Insomnia Collection Creation (3-4 hours)

**TASK-INS-01**: Create Harvey AI Folder (1h)
- Add folder with Basic Auth
- Create single endpoint with all query parameters
- Document response schema
- Test with seed data

**TASK-INS-02**: Create Copilot Folder (1h)
- Add folder with Basic Auth
- Create 4 period variants + CSV export
- Document OData response structure
- Test all periods

**TASK-INS-03**: Create Qualtrics Folder (1.5h)
- Add folder with Basic Auth
- Create 3 workflow endpoints
- Document multi-step process in folder description
- Test complete workflow

**TASK-INS-04**: Add Environment Variables (0.5h)
- Extend "Local Development" environment
- Add 6 new variables
- Document variable usage

### Phase 2: E2E Test Enhancement (2-3 hours)

**TASK-INS-05**: Verify Existing E2E Coverage (1h)
- Run all 14 existing tests
- Document any failures
- Identify coverage gaps

**TASK-INS-06**: Add Missing Test Scenarios (2h)
- Implement Harvey filtering tests
- Implement Copilot validation tests
- Implement Qualtrics error tests
- Ensure all schemas validated

### Phase 3: Documentation Updates (1-2 hours)

**TASK-INS-07**: Update SPEC.md (1h)
- Add External Data Sources section
- Document all 5 endpoints
- Include request/response examples
- Update endpoints table

**TASK-INS-08**: Create Insomnia Usage Guide (1h)
- Write docs/insomnia/README.md
- Document import process
- Provide workflow examples
- Add troubleshooting section

---

## Design Decisions

### Decision 1: No Admin API Changes

**Context**: External data sources might need admin endpoints for configuration

**Decision**: Use existing seed file configuration, no new admin endpoints

**Rationale**:
- Seed data already supports `external_data_sources` configuration
- Router conditionally registers endpoints based on seed config
- Admin API (P1-F02) focuses on runtime reconfiguration
- Seed upload (P1-F02-11) allows changing external source config
- Adding admin endpoints would duplicate existing functionality

**Alternatives Considered**:
- Add GET /admin/external-sources (rejected: seed config already exposed via GET /admin/config)
- Add POST /admin/external-sources/enable (rejected: can use POST /admin/seed instead)

### Decision 2: Extend Existing Insomnia Collection

**Context**: Could create separate collection or extend existing one

**Decision**: Extend `Insomnia_2026-01-09.yaml` with 3 new folders

**Rationale**:
- Maintains single source of truth
- Reuses existing environment variables (baseUrl, apiKey, dates)
- Consistent authentication pattern across all endpoints
- Easier onboarding (one import, all APIs available)

**Alternatives Considered**:
- Separate collections per API (rejected: management overhead, duplicated environments)
- Nested folder structure (rejected: existing collection uses flat structure)

### Decision 3: Focus on E2E Test Gaps

**Context**: 14 E2E tests already exist, could rewrite or enhance

**Decision**: Keep existing tests, add targeted scenarios for gaps

**Rationale**:
- Existing tests already validate core functionality
- Gaps are specific edge cases and error scenarios
- Adding tests is less risky than modifying existing ones
- Incremental improvement aligns with SDD methodology

**Alternatives Considered**:
- Rewrite all E2E tests (rejected: high risk, no added value)
- Add unit tests instead (rejected: external APIs need E2E validation)

---

## Risk Assessment

### Technical Risks

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Insomnia YAML format incompatibility | Low | Medium | Follow existing pattern, validate with import test |
| E2E tests flaky due to timing | Low | Medium | Use proper test server setup, avoid hardcoded delays |
| Documentation drift from code | Low | Low | Reference SPEC.md as source of truth |

### Schedule Risks

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| YAML editing more complex than expected | Low | Low | Existing collection provides clear pattern |
| E2E test scenarios require code changes | Very Low | Medium | Verify APIs already support all test scenarios |

---

## Success Criteria

- [ ] All 3 external API folders added to Insomnia collection
- [ ] 6 environment variables added and documented
- [ ] Collection successfully imports into Insomnia without errors
- [ ] All endpoints execute successfully with valid responses
- [ ] 9+ new E2E test scenarios added (3 per API)
- [ ] All E2E tests pass (existing 14 + new 9 = 23 total)
- [ ] SPEC.md has complete endpoint documentation
- [ ] docs/insomnia/README.md provides clear usage guide
- [ ] Zero onboarding questions from team after documentation release

---

## Future Enhancements (Out of Scope)

- Postman collection generation
- Swagger/OpenAPI spec generation
- Performance testing scenarios in Insomnia
- Automated collection generation from SPEC.md
- GraphQL playground for analytics-core (separate P5 work)

---

## References

- Existing Insomnia collection: `docs/insomnia/Insomnia_2026-01-09.yaml`
- E2E test patterns: `test/e2e/external_data_test.go`, `test/e2e/github_test.go`
- Handler implementations: `internal/api/harvey/`, `internal/api/microsoft/`, `internal/api/qualtrics/`
- Router registration: `internal/server/router.go` (lines 89-169)
- Seed configuration: `internal/seed/seed.go` (ExternalDataSources struct)
