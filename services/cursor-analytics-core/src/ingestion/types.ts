/**
 * Type definitions for cursor-sim API responses
 * Based on cursor-sim SPEC.md
 */

/**
 * Team member response from /teams/members
 */
export interface TeamMember {
  user_id: string;
  name: string;
  email: string;
  seniority: 'junior' | 'mid' | 'senior';
  ai_preference: number; // 0.0-1.0
  active: boolean;
}

/**
 * Commit response from /analytics/ai-code/commits
 */
export interface CursorCommit {
  commitHash: string;
  userId: string;
  userEmail: string;
  userName: string;
  repoName: string;
  branchName: string;
  isPrimaryBranch: boolean;
  totalLinesAdded: number;
  totalLinesDeleted: number;
  tabLinesAdded: number;
  tabLinesDeleted: number;
  composerLinesAdded: number;
  composerLinesDeleted: number;
  nonAiLinesAdded: number;
  nonAiLinesDeleted: number;
  message: string;
  commitTs: string; // ISO 8601
  createdAt: string; // ISO 8601
}

/**
 * Pagination metadata in API responses
 */
export interface PaginationMetadata {
  page: number;
  pageSize: number;
  totalPages: number;
  hasNextPage: boolean;
  hasPreviousPage: boolean;
}

/**
 * Query parameters echoed in API responses
 */
export interface ResponseParams {
  from?: string;
  to?: string;
  page?: number;
  pageSize?: number;
  user_id?: string;
  repo_name?: string;
}

/**
 * Generic paginated response wrapper
 */
export interface PaginatedResponse<T> {
  data: T[];
  pagination: PaginationMetadata;
  params: ResponseParams;
}

/**
 * Team members response (non-paginated)
 */
export interface TeamMembersResponse {
  data: TeamMember[];
}

/**
 * Query parameters for commits endpoint
 */
export interface CommitsQueryParams {
  from?: string; // YYYY-MM-DD
  to?: string; // YYYY-MM-DD
  page?: number;
  page_size?: number;
  user_id?: string;
  repo_name?: string;
}

/**
 * Configuration for cursor-sim client
 */
export interface CursorSimClientConfig {
  baseUrl: string; // e.g., "http://localhost:8080"
  apiKey: string; // Used as Basic Auth username
  timeout?: number; // Request timeout in milliseconds (default: 30000)
  retryAttempts?: number; // Max retry attempts (default: 3)
  retryDelayMs?: number; // Initial retry delay (default: 1000)
}

/**
 * Error response from cursor-sim API
 */
export interface CursorSimError {
  error: string;
  message: string;
  status?: number;
}
