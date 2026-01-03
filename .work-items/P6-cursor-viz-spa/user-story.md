# User Story: cursor-viz-spa Dashboard

**Feature ID**: P6-cursor-viz-spa
**Phase**: P6 (Visualization SPA)
**Created**: January 3, 2026
**Status**: TODO

## Overview

As a **SDLC researcher** or **engineering manager**, I want a **React-based visualization dashboard** so that I can **interactively explore AI coding assistant analytics** from the cursor-analytics-core GraphQL API.

## Acceptance Criteria

```gherkin
Given cursor-analytics-core is running with populated data
When I access the cursor-viz-spa dashboard
Then I can view visualizations of:
  - Team AI code usage trends
  - Developer leaderboards
  - Model usage breakdowns
  - Code quality metrics
  - Research correlations
And I can filter by date range, team, and developer
And I can export visualizations as PNG/SVG
```

## Dependencies

- cursor-analytics-core (P5) must be implemented first
- GraphQL API schema from P5

---

**Next Steps**: Create technical design and task breakdown
