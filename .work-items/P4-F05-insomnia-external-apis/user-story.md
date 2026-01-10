# User Story: P4-F05 Insomnia Collections for External APIs

**Feature**: P4-F05 Insomnia Collections for External Data Sources
**Phase**: P4 (CLI Enhancements)
**Created**: January 10, 2026
**Status**: Planned

---

## Overview

As a **developer or QA engineer**, I want **comprehensive Insomnia REST collections for Harvey, Copilot, and Qualtrics APIs** so that I can **easily explore, test, and validate external data source integrations without writing custom scripts**.

---

## User Stories (EARS Format)

### Story 1: Harvey AI API Testing

**As a** QA engineer
**I want** an Insomnia collection for Harvey AI endpoints
**So that** I can validate legal AI usage tracking and feedback mechanisms

**Given** the simulator is running with Harvey enabled in seed data
**When** I import the Insomnia collection and execute Harvey requests
**Then** I can query usage history with date ranges, user filters, and task filters
**And** responses match the documented schema with proper pagination

### Story 2: Copilot Usage API Testing

**As a** developer
**I want** an Insomnia collection for Microsoft 365 Copilot endpoints
**So that** I can test GitHub Copilot metrics across different time periods

**Given** the simulator is running with Copilot enabled in seed data
**When** I execute Copilot usage requests for different periods (D7, D30, D90, D180)
**Then** I receive usage data in both JSON and CSV formats
**And** app-specific breakdowns (Teams, Word, Outlook, Excel) are included

### Story 3: Qualtrics Survey Export Workflow

**As a** integration engineer
**I want** an Insomnia collection for Qualtrics export endpoints
**So that** I can test the multi-step survey response export workflow

**Given** the simulator is running with Qualtrics enabled in seed data
**When** I execute the three-step workflow (start export → poll progress → download ZIP)
**Then** the export job progresses through states correctly
**And** I can download the final response ZIP file

### Story 4: Comprehensive E2E Test Coverage

**As a** CI/CD pipeline maintainer
**I want** comprehensive E2E tests for all external APIs
**So that** automated testing can verify data generation and retrieval

**Given** external data sources are enabled in test seed data
**When** E2E tests execute against all endpoints
**Then** all tests pass with proper data validation
**And** error cases are properly handled

### Story 5: Documentation and Onboarding

**As a** new team member
**I want** clear documentation on using Insomnia collections
**So that** I can quickly get started testing external APIs

**Given** I have imported the Insomnia collection
**When** I read the usage guide in docs/insomnia/README.md
**Then** I understand how to configure environments and execute workflows
**And** I can troubleshoot common issues independently

---

## Acceptance Criteria

### Insomnia Collections

- [ ] Harvey AI folder created with all endpoints
  - [ ] GET /harvey/api/v1/history/usage with full parameter support
  - [ ] Environment variables for dates, user, task filters
  - [ ] Example requests with expected responses

- [ ] Copilot folder created with period variants
  - [ ] 4 period requests (D7, D30, D90, D180)
  - [ ] JSON and CSV export examples
  - [ ] OData response structure documented

- [ ] Qualtrics folder created with 3-step workflow
  - [ ] POST /API/v3/surveys/{surveyId}/export-responses
  - [ ] GET /API/v3/surveys/{surveyId}/export-responses/{progressId}
  - [ ] GET /API/v3/surveys/{surveyId}/export-responses/{fileId}/file
  - [ ] Workflow sequence documented in folder description

- [ ] Environment variables extended
  - [ ] harveyUser, harveyTask, copilotPeriod, qualtricsSurveyId, qualtricsProgressId, qualtricsFileId

### E2E Test Coverage

- [ ] Existing 14 tests verified and passing
- [ ] Additional test scenarios added:
  - [ ] Harvey: user filtering, task filtering, error cases
  - [ ] Copilot: invalid period, empty responses
  - [ ] Qualtrics: timeout, invalid survey ID
- [ ] All response schemas validated against SPEC.md

### Documentation

- [ ] SPEC.md updated with endpoint documentation
  - [ ] Request/response examples for all endpoints
  - [ ] Query parameters and validation rules
  - [ ] Error responses documented

- [ ] docs/insomnia/README.md created
  - [ ] Import instructions
  - [ ] Environment configuration guide
  - [ ] Workflow examples for each API
  - [ ] Troubleshooting section

---

## Business Value

1. **Reduced Testing Time**: Manual API testing becomes point-and-click instead of curl scripts
2. **Better Onboarding**: New team members can explore APIs interactively
3. **QA Efficiency**: Comprehensive test scenarios catch integration issues early
4. **Documentation**: Living examples complement static SPEC.md documentation

---

## Dependencies

- **Upstream**: P4-F04 (External Data Sources) - ✅ COMPLETE
- **Downstream**: None (documentation/testing enhancement)
- **Blocking**: None

---

## Out of Scope

- Creating new admin endpoints (seed config already supports external sources)
- Modifying existing external API implementations
- Adding new external data sources
- Performance testing or load testing

---

## Success Metrics

- All 8 tasks completed within 8-9 hours estimated time
- 100% of external API endpoints documented in Insomnia
- E2E test coverage increased to cover all error scenarios
- Zero onboarding questions about external API testing after documentation added

---

## Notes

- This is primarily a documentation and testing enhancement
- No code changes to cursor-sim service required
- APIs are already proven functional with 14 existing E2E tests
- Risk: LOW (declarative YAML, no behavioral changes)
