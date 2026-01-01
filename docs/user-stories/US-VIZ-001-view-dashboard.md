# User Story: US-VIZ-001

## View Dashboard Summary

**Story ID:** US-VIZ-001  
**Feature:** [F003 - Dashboard Views](../features/F003-dashboard-views.md)  
**Priority:** High  
**Story Points:** 5

---

## Story

**As an** engineering manager,  
**I want to** see a summary dashboard of AI coding assistant usage,  
**So that** I can quickly understand how my team is adopting and benefiting from Cursor.

---

## Description

The dashboard is the primary interface for understanding team-wide AI adoption. At a glance, an engineering manager needs to know whether the team is actively using the AI assistant, how effective that usage is (measured by acceptance rate), and who the top performers are. This information helps identify training needs, justify tool investments, and celebrate successes.

The dashboard should load quickly and present information hierarchically—the most important metrics at the top, followed by increasingly detailed views. The manager should not need to click around to get the basic picture; it should be immediately visible upon page load.

Real-time updates are important because the manager may leave the dashboard open during standup meetings to monitor activity. The data should refresh automatically without requiring a manual page reload.

---

## Acceptance Criteria

### AC1: Dashboard Loads Within 2 Seconds
**Given** the user navigates to the dashboard URL  
**When** the page loads  
**Then** all KPI cards and charts are visible within 2 seconds on a standard development machine

### AC2: KPI Cards Display Current Metrics
**Given** the dashboard is loaded  
**When** viewing the KPI card row  
**Then** I see cards for Active Developers, Average Acceptance Rate, AI Lines This Week, and Top Performer

### AC3: Active Developers Count is Accurate
**Given** 42 out of 50 developers have activity in the selected date range  
**When** viewing the Active Developers card  
**Then** the card displays "42" with a subtitle indicating "out of 50"

### AC4: Acceptance Rate Uses Correct Formula
**Given** developers have varying acceptance rates  
**When** viewing the Average Acceptance Rate card  
**Then** the displayed value equals the mean of all active developers' individual rates

### AC5: Top Performer is Correctly Identified
**Given** Alice has the highest acceptance rate at 95.2%  
**When** viewing the Top Performer card  
**Then** Alice's name is displayed with her acceptance rate

### AC6: Dashboard Refreshes Automatically
**Given** the dashboard is open  
**When** 30 seconds have elapsed  
**Then** the data refreshes and the "Last updated" timestamp changes

### AC7: Loading States Are User-Friendly
**Given** data is being fetched  
**When** the request is in progress  
**Then** skeleton loaders appear in place of actual content (not spinners or blank space)

---

## Technical Notes

The dashboard uses TanStack Query for data fetching with a staleTime of 25 seconds and refetchInterval of 30 seconds. This ensures data is cached but still refreshes automatically.

KPI calculations happen on the server (in the aggregator) to ensure consistency. The frontend simply displays the pre-calculated values.

Consider using React.lazy for chart components to improve initial load time. The charts can load asynchronously while KPI cards display immediately.

---

## Definition of Done

- [ ] All acceptance criteria pass
- [ ] Component tests verify each KPI card renders correctly
- [ ] Integration test confirms dashboard loads with mock GraphQL data
- [ ] E2E test verifies full user journey on real services
- [ ] Accessibility audit passes (keyboard navigation, screen reader support)
- [ ] Code reviewed and merged to main branch

---

## Related Tasks

- [TASK-020](../tasks/TASK-020-viz-layout.md): Implement dashboard layout
- [TASK-021](../tasks/TASK-021-viz-kpi-cards.md): Implement KPI summary cards
- [TASK-025](../tasks/TASK-025-viz-polling.md): Implement real-time polling

---

## Mockup Reference

```
┌──────────────────────────────────────────────────────────────────────┐
│  Cursor Analytics          [Last 30 Days ▼]    Updated: 2 min ago   │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌───────────┐│
│  │  👥 Active   │  │  ✓ Avg       │  │  📈 AI Lines │  │  🏆 Top   ││
│  │  Developers  │  │  Accept Rate │  │  This Week   │  │  Performer││
│  │              │  │              │  │              │  │           ││
│  │     42       │  │    72.3%     │  │   15,234     │  │ Alice S.  ││
│  │   of 50      │  │              │  │   ↑ 12%      │  │  95.2%    ││
│  └──────────────┘  └──────────────┘  └──────────────┘  └───────────┘│
│                                                                      │
└──────────────────────────────────────────────────────────────────────┘
```
