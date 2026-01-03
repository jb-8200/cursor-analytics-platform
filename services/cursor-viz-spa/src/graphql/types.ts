/**
 * Stub GraphQL types for cursor-viz-spa
 *
 * These types are placeholders until the cursor-analytics-core (P5)
 * GraphQL schema is fully implemented. They are based on the expected
 * schema defined in the design documents and api-contract.
 *
 * Once P5 is complete, these will be replaced by generated types
 * from GraphQL Code Generator.
 */

// ============================================================================
// Base Types
// ============================================================================

export interface Developer {
  id: string;
  name: string;
  email: string;
  team: string;
  seniority: 'junior' | 'mid' | 'senior';
  active: boolean;
  stats?: DeveloperStats;
}

export interface DeveloperStats {
  totalSuggestions: number;
  acceptedSuggestions: number;
  acceptanceRate: number;
  aiLinesAdded: number;
  totalLinesAdded: number;
  totalLinesDeleted: number;
  chatInteractions: number;
  cmdKUsages: number;
  aiVelocity: number;
}

export interface DailyStats {
  date: string; // ISO date string (YYYY-MM-DD)
  suggestionsShown: number;
  suggestionsAccepted: number;
  acceptanceRate: number;
  aiLinesAdded: number;
  linesAdded: number;
  chatInteractions: number;
  cmdKUsages: number;
}

export interface TeamStats {
  teamName: string;
  memberCount: number;
  activeMemberCount: number;
  averageAcceptanceRate: number;
  totalSuggestions: number;
  aiVelocity: number;
  chatInteractions: number;
  topPerformer?: Developer;
}

// ============================================================================
// Dashboard Types
// ============================================================================

export interface DashboardSummary {
  totalDevelopers: number;
  activeDevelopers: number;
  overallAcceptanceRate: number;
  totalSuggestionsToday: number;
  totalAcceptedToday: number;
  aiVelocityToday: number;
  teamComparison: TeamStats[];
  dailyTrend: DailyStats[];
}

// ============================================================================
// Pagination Types
// ============================================================================

export interface PageInfo {
  hasNextPage: boolean;
  hasPreviousPage: boolean;
  startCursor?: string;
  endCursor?: string;
}

export interface PaginatedDevelopers {
  nodes: Developer[];
  totalCount: number;
  pageInfo: PageInfo;
}

// ============================================================================
// Input Types
// ============================================================================

export interface DateRangeInput {
  from: string; // ISO date string
  to: string; // ISO date string
}

export interface DeveloperQueryInput {
  team?: string;
  limit?: number;
  offset?: number;
  search?: string;
}

// ============================================================================
// Query Response Types
// ============================================================================

export interface GetDashboardSummaryResponse {
  dashboardSummary: DashboardSummary;
}

export interface GetDevelopersResponse {
  developers: PaginatedDevelopers;
}

export interface GetTeamStatsResponse {
  teamStats: TeamStats;
}

export interface GetDailyStatsResponse {
  dailyStats: DailyStats[];
}

// ============================================================================
// Error Types
// ============================================================================

export interface GraphQLError {
  message: string;
  path?: string[];
  extensions?: {
    code?: string;
    [key: string]: unknown;
  };
}

export interface ErrorResponse {
  errors: GraphQLError[];
  data?: null;
}
