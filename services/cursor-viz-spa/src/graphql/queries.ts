import { gql } from '@apollo/client';

/**
 * GraphQL Query Stubs for cursor-viz-spa
 *
 * These are placeholder queries based on the expected schema from
 * cursor-analytics-core (P5). They will be updated once P5's GraphQL
 * schema is finalized.
 *
 * Note: These queries use the stub types defined in ./types.ts
 */

// ============================================================================
// Dashboard Queries
// ============================================================================

/**
 * Get dashboard summary with KPIs and trends
 *
 * Variables:
 * - range: DateRangeInput (optional) - Filter by date range
 *
 * Returns:
 * - DashboardSummary with KPIs, team comparison, and daily trends
 */
export const GET_DASHBOARD_SUMMARY = gql`
  query GetDashboardSummary($range: DateRangeInput) {
    dashboardSummary(range: $range) {
      totalDevelopers
      activeDevelopers
      overallAcceptanceRate
      totalSuggestionsToday
      totalAcceptedToday
      aiVelocityToday
      teamComparison {
        teamName
        memberCount
        activeMemberCount
        averageAcceptanceRate
        totalSuggestions
        aiVelocity
        chatInteractions
        topPerformers {
          id
          name
          email
          seniority
        }
      }
      dailyTrend {
        date
        suggestionsShown
        suggestionsAccepted
        acceptanceRate
        aiLinesAdded
        humanLinesAdded
        chatInteractions
      }
    }
  }
`;

// ============================================================================
// Developer Queries
// ============================================================================

/**
 * Get paginated list of developers with stats
 *
 * Variables:
 * - team: String (optional) - Filter by team name
 * - limit: Int (optional) - Page size (default: 25)
 * - offset: Int (optional) - Offset for pagination (default: 0)
 * - search: String (optional) - Search by developer name
 *
 * Returns:
 * - PaginatedDevelopers with nodes, totalCount, and pageInfo
 */
export const GET_DEVELOPERS = gql`
  query GetDevelopers($team: String, $limit: Int, $offset: Int, $search: String) {
    developers(team: $team, limit: $limit, offset: $offset, search: $search) {
      nodes {
        id
        name
        email
        team
        seniority
        active
        stats {
          totalSuggestions
          acceptedSuggestions
          acceptanceRate
          aiLinesAdded
          aiLinesDeleted
          humanLinesAdded
          humanLinesDeleted
          chatInteractions
          aiVelocity
        }
      }
      totalCount
      pageInfo {
        hasNextPage
        hasPreviousPage
        startCursor
        endCursor
      }
    }
  }
`;

/**
 * Get single developer by ID with detailed stats
 *
 * Variables:
 * - id: ID! - Developer ID
 *
 * Returns:
 * - Developer with complete stats
 */
export const GET_DEVELOPER = gql`
  query GetDeveloper($id: ID!) {
    developer(id: $id) {
      id
      name
      email
      team
      seniority
      active
      stats {
        totalSuggestions
        acceptedSuggestions
        acceptanceRate
        aiLinesAdded
        aiLinesDeleted
        humanLinesAdded
        humanLinesDeleted
        chatInteractions
        aiVelocity
      }
    }
  }
`;

// ============================================================================
// Team Queries
// ============================================================================

/**
 * Get team statistics
 *
 * Variables:
 * - teamName: String! - Team name
 * - range: DateRangeInput (optional) - Filter by date range
 *
 * Returns:
 * - TeamStats with aggregated metrics
 */
export const GET_TEAM_STATS = gql`
  query GetTeamStats($teamName: String!, $range: DateRangeInput) {
    teamStats(teamName: $teamName, range: $range) {
      teamName
      memberCount
      activeMemberCount
      averageAcceptanceRate
      totalSuggestions
      aiVelocity
      chatInteractions
      topPerformers {
        id
        name
        email
        seniority
        stats {
          acceptanceRate
          aiVelocity
        }
      }
    }
  }
`;

/**
 * Get all teams list
 *
 * Returns:
 * - Array of team names
 */
export const GET_TEAMS = gql`
  query GetTeams {
    teams
  }
`;

// ============================================================================
// Daily Stats Queries
// ============================================================================

/**
 * Get daily statistics for heatmap visualization
 *
 * Variables:
 * - range: DateRangeInput! - Date range for stats
 * - developerId: ID (optional) - Filter by developer
 * - teamName: String (optional) - Filter by team
 *
 * Returns:
 * - Array of DailyStats
 */
export const GET_DAILY_STATS = gql`
  query GetDailyStats($range: DateRangeInput!, $developerId: ID, $teamName: String) {
    dailyStats(range: $range, developerId: $developerId, teamName: $teamName) {
      date
      suggestionsShown
      suggestionsAccepted
      acceptanceRate
      aiLinesAdded
      humanLinesAdded
      chatInteractions
    }
  }
`;

// ============================================================================
// Health Check Query
// ============================================================================

/**
 * Health check query for GraphQL endpoint
 *
 * Returns:
 * - Health status
 */
export const HEALTH_CHECK = gql`
  query HealthCheck {
    health {
      status
      timestamp
    }
  }
`;
