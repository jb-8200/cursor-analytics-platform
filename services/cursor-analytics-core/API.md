# cursor-analytics-core GraphQL API Documentation

**Version**: 0.0.1-p0
**Service**: cursor-analytics-core
**Port**: 4000
**Endpoint**: `http://localhost:4000/graphql`

## Overview

The cursor-analytics-core service provides a GraphQL API for querying Cursor AI usage analytics. It aggregates data from the cursor-sim service and exposes rich queries for dashboards and reports.

## Quick Start

### Starting the Server

```bash
cd services/cursor-analytics-core
npm install
npm run db:migrate
npm run db:seed
npm run dev
```

The GraphQL Playground will be available at `http://localhost:4000/graphql`.

### Example Query

```graphql
query GetDashboard {
  dashboardSummary(preset: LAST_7_DAYS) {
    totalDevelopers
    activeDevelopers
    overallAcceptanceRate
  }
}
```

## Authentication

Currently, the API does not require authentication. In production, implement Bearer token authentication via HTTP headers.

## Schema Reference

### Scalars

#### DateTime

ISO 8601 date-time string.

```
2026-01-03T10:30:00Z
```

### Types

#### Developer

Represents a developer in the organization.

| Field | Type | Description |
|-------|------|-------------|
| `id` | `ID!` | Unique identifier (UUID) |
| `externalId` | `String!` | External system ID |
| `name` | `String!` | Developer's full name |
| `email` | `String!` | Developer's email address |
| `team` | `String!` | Team name |
| `seniority` | `String` | Seniority level (junior, mid, senior) |
| `createdAt` | `DateTime!` | Account creation timestamp |
| `stats(range)` | `UsageStats` | Aggregated usage statistics |
| `dailyStats(range)` | `[DailyStats!]!` | Daily breakdown of usage |

#### UsageStats

Aggregated usage metrics for a developer or team.

| Field | Type | Description |
|-------|------|-------------|
| `totalSuggestions` | `Int!` | Total AI suggestions shown |
| `acceptedSuggestions` | `Int!` | Total AI suggestions accepted |
| `acceptanceRate` | `Float` | Acceptance rate percentage (0-100) |
| `chatInteractions` | `Int!` | Total chat messages sent |
| `cmdKUsages` | `Int!` | Total Cmd+K prompts used |
| `totalLinesAdded` | `Int!` | Total lines of code added |
| `totalLinesDeleted` | `Int!` | Total lines of code deleted |
| `aiLinesAdded` | `Int!` | Lines added via AI suggestions |
| `aiVelocity` | `Float` | AI velocity percentage (0-100) |

#### DailyStats

Daily usage statistics breakdown.

| Field | Type | Description |
|-------|------|-------------|
| `date` | `DateTime!` | Date of statistics |
| `suggestionsShown` | `Int!` | Suggestions shown on this day |
| `suggestionsAccepted` | `Int!` | Suggestions accepted on this day |
| `acceptanceRate` | `Float` | Daily acceptance rate |
| `chatInteractions` | `Int!` | Chat messages on this day |
| `cmdKUsages` | `Int!` | Cmd+K usages on this day |
| `linesAdded` | `Int!` | Lines added on this day |
| `linesDeleted` | `Int!` | Lines deleted on this day |
| `aiLinesAdded` | `Int!` | AI-generated lines on this day |

#### Commit

Represents an accepted AI suggestion (commit).

| Field | Type | Description |
|-------|------|-------------|
| `id` | `ID!` | Unique identifier (UUID) |
| `externalId` | `String!` | External event ID |
| `timestamp` | `DateTime!` | When the suggestion was accepted |
| `linesAdded` | `Int!` | Lines added in this commit |
| `linesDeleted` | `Int!` | Lines deleted in this commit |
| `modelUsed` | `String` | AI model used (e.g., claude-sonnet-4-5) |
| `tokensInput` | `Int!` | Input tokens consumed |
| `tokensOutput` | `Int!` | Output tokens generated |
| `author` | `Developer!` | Developer who accepted the suggestion |

#### TeamStats

Aggregated statistics for a team.

| Field | Type | Description |
|-------|------|-------------|
| `teamName` | `String!` | Team name |
| `memberCount` | `Int!` | Total team members |
| `activeMemberCount` | `Int!` | Active members (with usage in period) |
| `averageAcceptanceRate` | `Float` | Weighted average acceptance rate |
| `totalSuggestions` | `Int!` | Total suggestions across team |
| `totalAccepted` | `Int!` | Total accepted suggestions |
| `chatInteractions` | `Int!` | Total chat interactions |
| `aiVelocity` | `Float` | Weighted average AI velocity |
| `topPerformer` | `Developer` | Developer with most AI lines |

#### DashboardKPI

High-level KPIs for the dashboard view.

| Field | Type | Description |
|-------|------|-------------|
| `totalDevelopers` | `Int!` | Total developers in organization |
| `activeDevelopers` | `Int!` | Active developers in period |
| `overallAcceptanceRate` | `Float` | Organization-wide acceptance rate |
| `totalSuggestionsToday` | `Int!` | Total suggestions shown today |
| `totalAcceptedToday` | `Int!` | Total suggestions accepted today |
| `aiVelocityToday` | `Float` | AI velocity for today |
| `teamComparison` | `[TeamStats!]!` | Team-by-team breakdown |
| `dailyTrend` | `[DailyStats!]!` | Daily trend data |

#### PageInfo

Pagination metadata.

| Field | Type | Description |
|-------|------|-------------|
| `hasNextPage` | `Boolean!` | Whether more results exist after this page |
| `hasPreviousPage` | `Boolean!` | Whether more results exist before this page |
| `startCursor` | `String` | Cursor for first item in page |
| `endCursor` | `String` | Cursor for last item in page |

### Inputs

#### DateRangeInput

Specifies a date range for filtering.

```graphql
input DateRangeInput {
  from: DateTime!
  to: DateTime!
}
```

**Example**:
```graphql
{
  from: "2026-01-01",
  to: "2026-01-31"
}
```

#### DateRangePreset

Enum for common date range presets.

| Value | Description |
|-------|-------------|
| `TODAY` | Current day only |
| `THIS_WEEK` | Current week (Sun-Sat) |
| `THIS_MONTH` | Current month |
| `LAST_7_DAYS` | Last 7 days including today |
| `LAST_30_DAYS` | Last 30 days including today |
| `LAST_90_DAYS` | Last 90 days including today |

## Queries

### health

Health check for all services.

**Returns**: `HealthStatus!`

**Example**:
```graphql
query HealthCheck {
  health {
    status
    database
    simulator
    lastIngestion
    version
  }
}
```

**Response**:
```json
{
  "data": {
    "health": {
      "status": "healthy",
      "database": "connected",
      "simulator": "reachable",
      "lastIngestion": "2026-01-03T10:30:00Z",
      "version": "0.0.1-p0"
    }
  }
}
```

### developer(id)

Fetch a single developer by ID.

**Arguments**:
- `id: ID!` - Developer UUID

**Returns**: `Developer`

**Example**:
```graphql
query GetDeveloper($id: ID!) {
  developer(id: $id) {
    id
    name
    email
    team
    seniority
    createdAt
  }
}
```

**Variables**:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### developers

List all developers with optional filtering and pagination.

**Arguments**:
- `team: String` - Filter by team name
- `seniority: String` - Filter by seniority level
- `limit: Int` - Max results per page (default: 50)
- `offset: Int` - Number of results to skip (default: 0)
- `sortBy: String` - Field to sort by (default: "name")
- `sortOrder: String` - Sort direction: "asc" or "desc" (default: "asc")

**Returns**: `DeveloperConnection!`

**Example 1: List all developers**:
```graphql
query ListDevelopers {
  developers(limit: 10) {
    nodes {
      id
      name
      email
      team
    }
    totalCount
    pageInfo {
      hasNextPage
      hasPreviousPage
    }
  }
}
```

**Example 2: Filter by team**:
```graphql
query FilterByTeam {
  developers(team: "Backend", limit: 20) {
    nodes {
      name
      email
      seniority
    }
    totalCount
  }
}
```

**Example 3: With usage stats**:
```graphql
query DevelopersWithStats {
  developers(limit: 10) {
    nodes {
      name
      stats(range: { from: "2026-01-01", to: "2026-01-31" }) {
        totalSuggestions
        acceptedSuggestions
        acceptanceRate
        aiVelocity
      }
    }
  }
}
```

### commits

List commits (accepted AI suggestions) with filtering and pagination.

**Arguments**:
- `userId: ID` - Filter by developer ID
- `team: String` - Filter by team name
- `dateRange: DateRangeInput` - Filter by date range
- `sortBy: String` - Field to sort by (default: "timestamp")
- `sortOrder: String` - Sort direction (default: "desc")
- `limit: Int` - Max results per page (default: 50)
- `offset: Int` - Number of results to skip (default: 0)

**Returns**: `CommitConnection!`

**Example 1: Recent commits**:
```graphql
query RecentCommits {
  commits(limit: 20, sortBy: "timestamp", sortOrder: "desc") {
    nodes {
      id
      timestamp
      linesAdded
      linesDeleted
      modelUsed
      author {
        name
        team
      }
    }
    totalCount
  }
}
```

**Example 2: Commits by user in date range**:
```graphql
query UserCommits($userId: ID!, $range: DateRangeInput!) {
  commits(userId: $userId, dateRange: $range, limit: 50) {
    nodes {
      timestamp
      linesAdded
      modelUsed
    }
    totalCount
  }
}
```

**Variables**:
```json
{
  "userId": "550e8400-e29b-41d4-a716-446655440000",
  "range": {
    "from": "2026-01-01",
    "to": "2026-01-31"
  }
}
```

### teamStats(teamName)

Get statistics for a specific team.

**Arguments**:
- `teamName: String!` - Team name

**Returns**: `TeamStats`

**Example**:
```graphql
query GetTeamStats($teamName: String!) {
  teamStats(teamName: $teamName) {
    teamName
    memberCount
    activeMemberCount
    averageAcceptanceRate
    totalSuggestions
    totalAccepted
    aiVelocity
    topPerformer {
      name
      email
      seniority
    }
  }
}
```

**Variables**:
```json
{
  "teamName": "Backend"
}
```

### teams

List all teams with statistics.

**Returns**: `[TeamStats!]!`

**Example**:
```graphql
query AllTeams {
  teams {
    teamName
    memberCount
    activeMemberCount
    averageAcceptanceRate
    totalSuggestions
    totalAccepted
    topPerformer {
      name
    }
  }
}
```

### dashboardSummary

Get comprehensive dashboard KPIs and trends.

**Arguments**:
- `range: DateRangeInput` - Custom date range
- `preset: DateRangePreset` - Preset date range (overrides `range`)

**Returns**: `DashboardKPI!`

**Example 1: Last 7 days**:
```graphql
query DashboardLast7Days {
  dashboardSummary(preset: LAST_7_DAYS) {
    totalDevelopers
    activeDevelopers
    overallAcceptanceRate
    totalSuggestionsToday
    totalAcceptedToday
    aiVelocityToday
    teamComparison {
      teamName
      memberCount
      averageAcceptanceRate
      topPerformer {
        name
      }
    }
    dailyTrend {
      date
      suggestionsShown
      suggestionsAccepted
      acceptanceRate
    }
  }
}
```

**Example 2: Custom date range**:
```graphql
query DashboardCustomRange($range: DateRangeInput!) {
  dashboardSummary(range: $range) {
    totalDevelopers
    activeDevelopers
    overallAcceptanceRate
    teamComparison {
      teamName
      totalSuggestions
      totalAccepted
    }
  }
}
```

**Variables**:
```json
{
  "range": {
    "from": "2026-01-01",
    "to": "2026-01-31"
  }
}
```

## Complex Query Examples

### Developer Profile with Full Stats

```graphql
query DeveloperProfile($id: ID!, $range: DateRangeInput!) {
  developer(id: $id) {
    id
    name
    email
    team
    seniority
    createdAt

    # Overall stats for the period
    stats(range: $range) {
      totalSuggestions
      acceptedSuggestions
      acceptanceRate
      chatInteractions
      cmdKUsages
      totalLinesAdded
      aiLinesAdded
      aiVelocity
    }

    # Daily breakdown for trend chart
    dailyStats(range: $range) {
      date
      suggestionsShown
      suggestionsAccepted
      acceptanceRate
      linesAdded
      aiLinesAdded
    }
  }
}
```

### Team Comparison Dashboard

```graphql
query TeamComparison {
  teams {
    teamName
    memberCount
    activeMemberCount
    averageAcceptanceRate
    totalSuggestions
    totalAccepted
    aiVelocity
    chatInteractions
    topPerformer {
      name
      email
      seniority
      stats {
        totalSuggestions
        acceptanceRate
        aiVelocity
      }
    }
  }
}
```

### Full Activity Feed

```graphql
query ActivityFeed($limit: Int!, $offset: Int!) {
  commits(limit: $limit, offset: $offset, sortBy: "timestamp", sortOrder: "desc") {
    nodes {
      id
      timestamp
      linesAdded
      linesDeleted
      modelUsed
      tokensInput
      tokensOutput
      author {
        name
        email
        team
        seniority
      }
    }
    totalCount
    pageInfo {
      hasNextPage
      hasPreviousPage
    }
  }
}
```

### Filtered Analytics Report

```graphql
query FilteredReport($team: String!, $range: DateRangeInput!) {
  # Team overview
  teamStats(teamName: $team) {
    memberCount
    averageAcceptanceRate
    aiVelocity
  }

  # Team members with stats
  developers(team: $team) {
    nodes {
      name
      seniority
      stats(range: $range) {
        acceptanceRate
        aiVelocity
        totalSuggestions
      }
    }
  }

  # Recent commits from team
  commits(team: $team, dateRange: $range, limit: 100) {
    nodes {
      timestamp
      linesAdded
      modelUsed
      author {
        name
      }
    }
    totalCount
  }
}
```

## Error Handling

Errors are returned in standard GraphQL format:

```json
{
  "errors": [
    {
      "message": "Developer not found",
      "extensions": {
        "code": "NOT_FOUND",
        "id": "invalid-id"
      }
    }
  ]
}
```

### Error Codes

| Code | Description |
|------|-------------|
| `NOT_FOUND` | Requested resource doesn't exist |
| `VALIDATION_ERROR` | Invalid input parameters |
| `DATABASE_ERROR` | Database operation failed |
| `SIMULATOR_ERROR` | Failed to fetch from cursor-sim |

## Performance Considerations

### Query Complexity

Queries are analyzed for complexity to prevent expensive operations. The complexity limit is configurable via `QUERY_COMPLEXITY_LIMIT` environment variable (default: 1000).

### Pagination

Always use pagination for large result sets:
- Default `limit`: 50
- Maximum `limit`: 100
- Use `offset` for simple pagination
- Use cursors for more efficient pagination

### Caching

- Individual requests are not cached by default
- Consider implementing Redis caching for frequently accessed queries
- `dashboardSummary` is a good candidate for caching with 30-60 second TTL

### Best Practices

1. **Request only needed fields** - GraphQL allows precise field selection
2. **Use date range filters** - Reduce dataset size for better performance
3. **Paginate large results** - Always specify `limit` for list queries
4. **Batch related queries** - Use GraphQL's query batching capability
5. **Monitor query performance** - Track execution times in production

## Environment Variables

```bash
# Database
DATABASE_URL=postgresql://user:password@host:5432/cursor_analytics

# cursor-sim
SIMULATOR_URL=http://localhost:8080
SIMULATOR_API_KEY=cursor-sim-dev-key

# Server
PORT=4000
NODE_ENV=development

# Performance
QUERY_COMPLEXITY_LIMIT=1000
```

## Testing

Run the test suite:

```bash
npm test                # Unit tests
npm run test:coverage   # With coverage report
```

Integration tests are in `src/__tests__/integration/` and require a PostgreSQL database.

## Support

For issues or questions:
- GitHub Issues: [cursor-analytics-platform/issues](https://github.com/your-org/cursor-analytics-platform/issues)
- Documentation: See `SPEC.md` for technical details
