/**
 * CursorSimClient - REST client for cursor-sim API
 * Handles authentication, pagination, retries, and error handling
 */

import type {
  CursorSimClientConfig,
  TeamMember,
  TeamMembersResponse,
  CursorCommit,
  PaginatedResponse,
  CommitsQueryParams,
} from './types';

/**
 * HTTP client for cursor-sim REST API
 */
export class CursorSimClient {
  private readonly baseUrl: string;
  private readonly apiKey: string;
  private readonly timeout: number;
  private readonly retryAttempts: number;
  private readonly retryDelayMs: number;
  private readonly authHeader: string;

  constructor(config: CursorSimClientConfig) {
    this.baseUrl = config.baseUrl.replace(/\/$/, ''); // Remove trailing slash
    this.apiKey = config.apiKey;
    this.timeout = config.timeout ?? 30000;
    this.retryAttempts = config.retryAttempts ?? 3;
    this.retryDelayMs = config.retryDelayMs ?? 1000;

    // Pre-compute Basic Auth header: "Basic base64(apiKey:)"
    this.authHeader = `Basic ${Buffer.from(`${this.apiKey}:`).toString('base64')}`;
  }

  /**
   * Fetch team members from cursor-sim
   * GET /teams/members
   */
  async getTeamMembers(): Promise<TeamMember[]> {
    const response = await this.request<TeamMembersResponse>(
      '/teams/members'
    );
    return response.data;
  }

  /**
   * Fetch commits from cursor-sim with optional filtering and pagination
   * GET /analytics/ai-code/commits
   */
  async getCommits(
    params?: CommitsQueryParams
  ): Promise<PaginatedResponse<CursorCommit>> {
    const queryString = this.buildQueryString(params);
    const url = `/analytics/ai-code/commits${queryString ? `?${queryString}` : ''}`;
    return this.request<PaginatedResponse<CursorCommit>>(url);
  }

  /**
   * Build URL query string from parameters
   */
  private buildQueryString(params?: CommitsQueryParams): string {
    if (!params) return '';

    const entries = Object.entries(params)
      .filter(([_, value]) => value !== undefined && value !== null)
      .map(([key, value]) => `${key}=${encodeURIComponent(String(value))}`);

    return entries.join('&');
  }

  /**
   * Execute HTTP request with retry logic and error handling
   */
  private async request<T>(
    path: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseUrl}${path}`;
    const headers = {
      'Content-Type': 'application/json',
      Authorization: this.authHeader,
      ...options.headers,
    };

    let lastError: Error | null = null;

    for (let attempt = 0; attempt < this.retryAttempts; attempt++) {
      try {
        const controller = new AbortController();
        const timeoutId = setTimeout(() => controller.abort(), this.timeout);

        let response;
        try {
          response = await fetch(url, {
            ...options,
            headers,
            signal: controller.signal,
          });
          clearTimeout(timeoutId);
        } catch (fetchError) {
          clearTimeout(timeoutId);
          throw fetchError;
        }

        // Handle HTTP errors
        if (!response.ok) {
          const shouldRetry = this.shouldRetryStatus(response.status);
          const errorMessage = `${response.statusText} (${response.status})`;

          if (!shouldRetry) {
            // Don't retry on client errors (4xx except 429)
            // Throw immediately, bypassing the retry loop
            throw new ClientError(errorMessage);
          }

          // Retry on server errors (5xx) and 429
          if (attempt < this.retryAttempts - 1) {
            lastError = new Error(errorMessage);
            await this.delay(this.calculateRetryDelay(attempt));
            continue;
          }

          throw new Error(errorMessage);
        }

        // Parse JSON response
        const data = await response.json();
        return data as T;
      } catch (error) {
        // If it's a ClientError, rethrow immediately without retry
        if (error instanceof ClientError) {
          throw new Error(error.message);
        }

        lastError = error as Error;

        // Check if it's an abort error (timeout)
        if (error instanceof Error && error.name === 'AbortError') {
          throw new Error(`Request timeout after ${this.timeout}ms`);
        }

        // Retry on network errors if not the last attempt
        if (attempt < this.retryAttempts - 1) {
          await this.delay(this.calculateRetryDelay(attempt));
          continue;
        }
      }
    }

    // All retries exhausted
    throw lastError || new Error('Request failed after retries');
  }

  /**
   * Determine if HTTP status code should trigger a retry
   */
  private shouldRetryStatus(status: number): boolean {
    // Retry on server errors (5xx) and rate limiting (429)
    // Do not retry on client errors (4xx) except 429
    return status >= 500 || status === 429;
  }

  /**
   * Calculate exponential backoff delay
   */
  private calculateRetryDelay(attempt: number): number {
    // Exponential backoff: delay * (2 ^ attempt)
    // Capped at 30 seconds
    const delay = this.retryDelayMs * Math.pow(2, attempt);
    return Math.min(delay, 30000);
  }

  /**
   * Delay helper for retry logic
   */
  private delay(ms: number): Promise<void> {
    return new Promise((resolve) => setTimeout(resolve, ms));
  }
}

/**
 * Custom error class for non-retryable client errors (4xx)
 */
class ClientError extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'ClientError';
  }
}
