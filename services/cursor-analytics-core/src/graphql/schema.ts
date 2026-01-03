/**
 * GraphQL Schema Definition
 *
 * Defines the complete type system for the cursor-analytics-core GraphQL API.
 * This schema exposes data from cursor-sim through an aggregated view optimized
 * for frontend consumption by cursor-viz-spa.
 */

export const typeDefs = `#graphql
  """
  ISO 8601 DateTime scalar
  """
  scalar DateTime

  # ===========================
  # Query Root
  # ===========================

  type Query {
    """
    Get a single developer by ID
    """
    developer(id: ID!): Developer

    """
    List developers with optional filtering
    """
    developers(
      team: String
      seniority: String
      limit: Int = 50
      offset: Int = 0
      sortBy: String = "name"
      sortOrder: String = "asc"
    ): DeveloperConnection!

    """
    List commits (accepted AI suggestions) with optional filtering
    """
    commits(
      userId: ID
      team: String
      dateRange: DateRangeInput
      sortBy: String = "timestamp"
      sortOrder: String = "desc"
      limit: Int = 50
      offset: Int = 0
    ): CommitConnection!

    """
    Get team statistics
    """
    teamStats(teamName: String!): TeamStats

    """
    List all teams with statistics
    """
    teams: [TeamStats!]!

    """
    Dashboard summary - optimized for main dashboard view
    """
    dashboardSummary(
      range: DateRangeInput
      preset: DateRangePreset
    ): DashboardKPI!

    """
    Health check endpoint
    """
    health: HealthStatus!
  }

  # ===========================
  # Developer Types
  # ===========================

  """
  Developer profile information
  """
  type Developer {
    id: ID!
    externalId: String!
    name: String!
    email: String!
    team: String!
    seniority: String

    """
    Aggregated usage statistics for the developer
    """
    stats(range: DateRangeInput): UsageStats

    """
    Daily statistics breakdown
    """
    dailyStats(range: DateRangeInput): [DailyStats!]!
  }

  """
  Paginated developer list response
  """
  type DeveloperConnection {
    nodes: [Developer!]!
    totalCount: Int!
    pageInfo: PageInfo!
  }

  """
  Pagination metadata
  """
  type PageInfo {
    hasNextPage: Boolean!
    hasPreviousPage: Boolean!
    startCursor: String
    endCursor: String
  }

  # ===========================
  # Commit Types
  # ===========================

  """
  Commit (accepted AI suggestion) information
  """
  type Commit {
    id: ID!
    externalId: String!
    timestamp: DateTime!
    linesAdded: Int!
    linesDeleted: Int!
    modelUsed: String
    tokensInput: Int!
    tokensOutput: Int!
    author: Developer!
  }

  """
  Paginated commit list response
  """
  type CommitConnection {
    nodes: [Commit!]!
    totalCount: Int!
    pageInfo: PageInfo!
  }

  # ===========================
  # Statistics Types
  # ===========================

  """
  Aggregated usage statistics
  """
  type UsageStats {
    totalSuggestions: Int!
    acceptedSuggestions: Int!
    acceptanceRate: Float
    chatInteractions: Int!
    cmdKUsages: Int!
    totalLinesAdded: Int!
    totalLinesDeleted: Int!
    aiLinesAdded: Int!
    aiVelocity: Float
  }

  """
  Daily statistics for a developer
  """
  type DailyStats {
    date: DateTime!
    suggestionsShown: Int!
    suggestionsAccepted: Int!
    acceptanceRate: Float
    chatInteractions: Int!
    cmdKUsages: Int!
    linesAdded: Int!
    linesDeleted: Int!
    aiLinesAdded: Int!
  }

  """
  Team-level statistics
  """
  type TeamStats {
    teamName: String!
    memberCount: Int!
    activeMemberCount: Int!
    averageAcceptanceRate: Float
    totalSuggestions: Int!
    totalAccepted: Int!
    chatInteractions: Int!
    aiVelocity: Float
    topPerformer: Developer
  }

  """
  Dashboard KPIs - top-level metrics
  """
  type DashboardKPI {
    totalDevelopers: Int!
    activeDevelopers: Int!
    overallAcceptanceRate: Float
    totalSuggestionsToday: Int!
    totalAcceptedToday: Int!
    aiVelocityToday: Float
    teamComparison: [TeamStats!]!
    dailyTrend: [DailyStats!]!
  }

  # ===========================
  # Input Types
  # ===========================

  """
  Date range filter
  """
  input DateRangeInput {
    from: DateTime!
    to: DateTime!
  }

  """
  Common date range presets
  """
  enum DateRangePreset {
    TODAY
    THIS_WEEK
    THIS_MONTH
    LAST_7_DAYS
    LAST_30_DAYS
    LAST_90_DAYS
  }

  # ===========================
  # Health Check Types
  # ===========================

  """
  Service health status
  """
  type HealthStatus {
    status: String!
    database: String!
    simulator: String!
    lastIngestion: DateTime
    version: String!
  }
`;
