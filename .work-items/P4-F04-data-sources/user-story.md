# User Story: External Data Source Simulators

**Feature ID**: P4-F04-data-sources
**Created**: January 9, 2026
**Status**: Planning
**Priority**: P4 (cursor-sim enhancements)

---

## Story (EARS Format)

### As-a / I-want / So-that

**As a** researcher or platform developer building AI adoption analytics
**I want** cursor-sim to simulate additional enterprise AI tool APIs (Harvey, Microsoft 365 Copilot, Qualtrics surveys)
**So that** I can develop and test analytics pipelines against realistic synthetic data without production API access

---

## Context

### Current State

1. **Limited API Coverage**
   - cursor-sim only simulates Cursor Business API and GitHub
   - No support for legal AI (Harvey), productivity AI (M365 Copilot), or survey data (Qualtrics)
   - Analytics platform cannot test multi-source correlation scenarios

2. **Missing Enterprise AI Tracking**
   - Harvey: Leading AI legal assistant with no test data available
   - M365 Copilot: Microsoft's productivity AI with complex usage reporting
   - Qualtrics: Survey platform for qualitative AI feedback data

3. **No Async API Patterns**
   - Current simulators are synchronous request/response
   - Qualtrics requires 3-step async state machine pattern
   - Missing realistic enterprise integration patterns

### Desired State

1. **Harvey API Simulator**
   - Generates realistic AI legal assistant usage events
   - Tracks tasks (Assist, Draft, Review, Research), feedback, and client matters
   - Correlates with developer activity patterns

2. **Microsoft 365 Copilot Usage API Simulator**
   - Implements Microsoft Graph API beta endpoint format
   - Supports period parameters (D7, D30, D90, D180, ALL)
   - Generates per-application usage data (Teams, Word, Excel, etc.)

3. **Qualtrics Survey Export API Simulator**
   - Implements 3-step async state machine (start, poll, download)
   - Generates realistic survey responses about AI tool satisfaction
   - Produces ZIP/CSV export format

---

## Requirements

### Functional Requirements (EARS)

#### Harvey API

##### FR-H01: Usage History Endpoint
**WHEN** client calls `GET /harvey/api/v1/history/usage`
**THEN** system returns array of usage events
**AND** each event includes event_id, user, task type, source, and feedback
**AND** events are filtered by date range if provided

##### FR-H02: Task Type Distribution
**WHEN** generating Harvey usage events
**THEN** system distributes tasks across types:
- Assist: 35% (general questions)
- Draft: 30% (document drafting)
- Review: 25% (contract review)
- Research: 10% (legal research)
**AND** distribution varies by user role

##### FR-H03: Feedback Sentiment
**WHEN** generating feedback for events
**THEN** system assigns sentiment (positive, negative, neutral)
**AND** sentiment correlates with task completion quality
**AND** negative feedback includes descriptive comments

#### Microsoft 365 Copilot API

##### FR-M01: Usage User Detail Endpoint
**WHEN** client calls `GET /reports/getMicrosoft365CopilotUsageUserDetail(period='{period}')`
**THEN** system returns user-level Copilot usage data
**AND** supports period values: D7, D30, D90, D180, ALL
**AND** returns JSON or CSV based on $format parameter

##### FR-M02: Per-Application Activity Dates
**WHEN** returning user usage detail
**THEN** response includes last activity date for each Copilot-enabled app:
- Microsoft Teams Copilot
- Word Copilot
- Excel Copilot
- PowerPoint Copilot
- Outlook Copilot
- OneNote Copilot
- Loop Copilot
- Copilot Chat

##### FR-M03: Pagination Support
**WHEN** JSON response contains more than 100 users
**THEN** response includes `@odata.nextLink` for pagination
**AND** pagination uses skiptoken pattern

##### FR-M04: CSV Redirect
**WHEN** $format=text/csv is requested
**THEN** system returns 302 redirect
**AND** Location header contains pre-authenticated download URL
**AND** download URL is valid for simulated timeframe

#### Qualtrics Survey Export API

##### FR-Q01: Start Export Endpoint
**WHEN** client calls `POST /API/v3/surveys/{surveyId}/export-responses`
**THEN** system creates export job
**AND** returns progressId, status="inProgress", percentComplete=0

##### FR-Q02: Progress Polling Endpoint
**WHEN** client calls `GET /API/v3/surveys/{surveyId}/export-responses/{progressId}`
**THEN** system returns current progress (0-100%)
**AND** status transitions: "inProgress" -> "complete"
**AND** on completion, returns fileId for download

##### FR-Q03: File Download Endpoint
**WHEN** client calls `GET /API/v3/surveys/{surveyId}/export-responses/{fileId}/file`
**THEN** system returns ZIP file containing CSV
**AND** CSV contains survey responses with AI satisfaction questions

##### FR-Q04: State Machine Behavior
**WHEN** export job progresses
**THEN** system simulates realistic timing:
- First poll: 10-20% progress
- Subsequent polls: +20-30% progress each
- Total completion time: 5-15 seconds simulated (configurable)

---

### Non-Functional Requirements

#### NFR1: Consistency with Existing Patterns
- Follow cursor-sim model/generator/handler architecture
- Use same seed-based generation approach
- Integrate with existing storage layer

#### NFR2: Seed Extension
- Extend seed schema to include Harvey users, M365 tenants, survey configurations
- Support correlation with existing developer profiles

#### NFR3: Test Coverage
- 90%+ coverage for new models
- 85%+ coverage for generators
- E2E tests for all endpoints

#### NFR4: Performance
- Generation time < 100ms per 1000 events
- Memory footprint < 10MB per API source
- API response time < 50ms

---

## Acceptance Criteria (Given-When-Then)

### AC1: Harvey Usage Events
**GIVEN** seed file with 5 users having Harvey access
**WHEN** I query `GET /harvey/api/v1/history/usage?from=2026-01-01&to=2026-01-31`
**THEN** I receive 200-500 usage events
**AND** events span all 4 task types
**AND** feedback sentiment is realistically distributed
**AND** user emails match seed file

### AC2: Microsoft Copilot Usage Report
**GIVEN** seed file with 10 M365 users
**WHEN** I query `GET /reports/getMicrosoft365CopilotUsageUserDetail(period='D30')`
**THEN** I receive JSON with 10 user records
**AND** each record has last activity dates for all 8 apps
**AND** activity dates are within the last 30 days
**AND** reportRefreshDate matches current date

### AC3: Microsoft Copilot CSV Export
**GIVEN** seed file with M365 users
**WHEN** I query with `$format=text/csv`
**THEN** I receive 302 redirect
**AND** following the redirect returns valid CSV
**AND** CSV headers match Microsoft Graph schema

### AC4: Qualtrics Export Lifecycle
**GIVEN** survey configuration in seed file
**WHEN** I execute the 3-step export flow:
1. POST to start export
2. Poll GET until complete
3. GET file download
**THEN** each step returns expected response
**AND** progress increases monotonically
**AND** final ZIP contains valid CSV survey responses

### AC5: Qualtrics Survey Content
**GIVEN** completed survey export
**WHEN** I extract the CSV from ZIP
**THEN** responses include AI satisfaction questions:
- Overall AI tool satisfaction (1-5 scale)
- Specific tool ratings (Harvey, Copilot, Cursor)
- Free-text feedback
- Respondent demographics

### AC6: Cross-Source Correlation
**GIVEN** seed file with users having access to all three tools
**WHEN** I query all three APIs for the same time period
**THEN** user emails are consistent across sources
**AND** activity patterns show realistic correlation
**AND** high Harvey usage correlates with legal-team membership

---

## Out of Scope

- **OAuth authentication simulation**: Basic auth sufficient for simulator
- **Real file upload to cloud storage**: Local file serving only
- **Rate limiting per Microsoft specs**: Simplified throttling
- **Full Qualtrics survey builder**: Fixed survey templates only
- **Real-time streaming exports**: Batch export only

---

## Success Metrics

- All 3 APIs functional and passing E2E tests
- Seed schema documented with examples
- Integration with existing cursor-sim router
- Ready for downstream analytics-core consumption

---

## User Scenarios

### Scenario 1: Legal AI Analytics Development
```bash
# Generate Harvey usage data
./bin/cursor-sim -mode runtime -seed testdata/enterprise_seed.yaml

# Query Harvey API
curl -u api-key: "http://localhost:8080/harvey/api/v1/history/usage?from=2026-01-01&to=2026-01-31"

# Response: 250 usage events
{
  "data": [
    {
      "event_id": 103230489,
      "message_ID": "ab12a1ab-abcd-1a12-1234-1234ab123456",
      "Time": "2026-01-15T09:30:00Z",
      "User": "attorney@lawfirm.com",
      "Task": "Review",
      "Client Matter #": 2024.789,
      "Source": "Files",
      "Number of documents": 3,
      "Feedback Sentiment": "positive"
    },
    ...
  ]
}
```

### Scenario 2: M365 Copilot Adoption Report
```bash
# Query Copilot usage
curl -u api-key: "http://localhost:8080/reports/getMicrosoft365CopilotUsageUserDetail(period='D30')?$format=application/json"

# Response: User-level adoption metrics
{
  "@odata.nextLink": null,
  "value": [
    {
      "reportRefreshDate": "2026-01-09",
      "reportPeriod": 30,
      "userPrincipalName": "dev@company.com",
      "displayName": "Jane Developer",
      "lastActivityDate": "2026-01-08",
      "microsoftTeamsCopilotLastActivityDate": "2026-01-08",
      "wordCopilotLastActivityDate": "2026-01-05",
      "excelCopilotLastActivityDate": null,
      "powerPointCopilotLastActivityDate": "2025-12-20",
      "outlookCopilotLastActivityDate": "2026-01-07",
      "oneNoteCopilotLastActivityDate": null,
      "loopCopilotLastActivityDate": null,
      "copilotChatLastActivityDate": "2026-01-09"
    }
  ]
}
```

### Scenario 3: Survey Export Workflow
```bash
# Step 1: Start export
curl -X POST -u api-key: "http://localhost:8080/API/v3/surveys/SV_abc123/export-responses"
# Response: {"progressId": "ES_xyz789", "percentComplete": 0, "status": "inProgress"}

# Step 2: Poll progress (repeat until complete)
curl -u api-key: "http://localhost:8080/API/v3/surveys/SV_abc123/export-responses/ES_xyz789"
# Response: {"percentComplete": 45, "status": "inProgress"}
# ... later ...
# Response: {"percentComplete": 100, "status": "complete", "fileId": "FILE_abc"}

# Step 3: Download file
curl -u api-key: "http://localhost:8080/API/v3/surveys/SV_abc123/export-responses/FILE_abc/file" -o responses.zip

# Unzip and analyze
unzip responses.zip
# Contains: survey_responses.csv
```

---

## Dependencies

- **Completed**: P4-F03 (TUI enhancements with event architecture)
- **Completed**: P3-F01 (Research framework with export patterns)
- **Required**: Seed schema extension design
- **Required**: Router integration points

---

## Risks & Mitigations

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Seed schema complexity | Medium | Medium | Incremental extension with defaults |
| Qualtrics state machine bugs | Medium | High | Comprehensive state transition tests |
| Microsoft API format changes | Low | Medium | Document beta API caveat in SPEC.md |
| ZIP file generation complexity | Low | Medium | Use standard library archive/zip |
| Cross-source correlation logic | Medium | Medium | Seed-based user linking |

---

## Timeline Estimate

| Phase | Tasks | Hours |
|-------|-------|-------|
| Harvey API | Models + Generator + Handler + Tests | 6.0h |
| Microsoft Copilot API | Models + Generator + Handler + Tests | 6.5h |
| Qualtrics API | Models + State Machine + Handler + Tests | 8.0h |
| Integration | Router + Seed Extension + E2E | 4.0h |
| Documentation | SPEC.md + Examples | 1.5h |
| **TOTAL** | **15 tasks** | **26.0h** |

---

**Next Step**: Review design.md for technical architecture and API contracts.
