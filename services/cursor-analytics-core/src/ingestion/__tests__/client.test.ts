/**
 * Unit tests for CursorSimClient
 * Following TDD approach - these tests define expected behavior
 */

/* eslint-disable @typescript-eslint/require-await */
/* eslint-disable @typescript-eslint/no-unsafe-assignment */
/* eslint-disable @typescript-eslint/no-unsafe-member-access */

import { CursorSimClient } from '../client';
import type {
  TeamMember,
  CursorCommit,
  PaginatedResponse,
  TeamMembersResponse,
} from '../types';

// Mock fetch globally
global.fetch = jest.fn();

describe('CursorSimClient', () => {
  const mockConfig = {
    baseUrl: 'http://localhost:8080',
    apiKey: 'test-api-key',
    timeout: 5000,
    retryAttempts: 3, // Total attempts including initial (1 initial + 2 retries)
    retryDelayMs: 100,
  };

  let client: CursorSimClient;

  beforeEach(() => {
    client = new CursorSimClient(mockConfig);
    jest.clearAllMocks();
  });

  afterEach(() => {
    jest.restoreAllMocks();
  });

  describe('constructor', () => {
    it('should create client with provided config', () => {
      expect(client).toBeInstanceOf(CursorSimClient);
    });

    it('should use default timeout if not provided', () => {
      const clientWithDefaults = new CursorSimClient({
        baseUrl: 'http://localhost:8080',
        apiKey: 'test-key',
      });
      expect(clientWithDefaults).toBeInstanceOf(CursorSimClient);
    });
  });

  describe('getTeamMembers', () => {
    const mockMembers: TeamMember[] = [
      {
        user_id: 'user_001',
        name: 'Alice Developer',
        email: 'alice@example.com',
        seniority: 'senior',
        ai_preference: 0.8,
        active: true,
      },
      {
        user_id: 'user_002',
        name: 'Bob Engineer',
        email: 'bob@example.com',
        seniority: 'mid',
        ai_preference: 0.5,
        active: true,
      },
    ];

    it('should fetch team members successfully', async () => {
      const mockResponse: TeamMembersResponse = {
        data: mockMembers,
      };

      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => mockResponse,
      });

      const result = await client.getTeamMembers();

      expect(result).toEqual(mockMembers);
      expect(global.fetch).toHaveBeenCalledWith(
        'http://localhost:8080/teams/members',
        expect.objectContaining({
          headers: expect.objectContaining({
            Authorization: expect.stringContaining('Basic'),
          }),
        })
      );
    });

    it('should include correct Basic Auth header', async () => {
      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({ data: [] }),
      });

      await client.getTeamMembers();

      const authHeader = (global.fetch as jest.Mock).mock.calls[0][1].headers
        .Authorization;
      // Basic auth format: "Basic base64(apiKey:)"
      const expectedAuth = `Basic ${Buffer.from('test-api-key:').toString('base64')}`;
      expect(authHeader).toBe(expectedAuth);
    });

    it('should throw error on 401 Unauthorized', async () => {
      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: false,
        status: 401,
        statusText: 'Unauthorized',
        json: async () => ({ error: 'Invalid API key' }),
      });

      await expect(client.getTeamMembers()).rejects.toThrow('Unauthorized');
    });

    it('should throw error on network failure', async () => {
      (global.fetch as jest.Mock).mockRejectedValueOnce(
        new Error('Network error')
      );

      await expect(client.getTeamMembers()).rejects.toThrow();
    });
  });

  describe('getCommits', () => {
    const mockCommits: CursorCommit[] = [
      {
        commitHash: 'abc123',
        userId: 'user_001',
        userEmail: 'alice@example.com',
        userName: 'Alice Developer',
        repoName: 'acme/platform',
        branchName: 'main',
        isPrimaryBranch: true,
        totalLinesAdded: 100,
        totalLinesDeleted: 20,
        tabLinesAdded: 60,
        tabLinesDeleted: 10,
        composerLinesAdded: 30,
        composerLinesDeleted: 5,
        nonAiLinesAdded: 10,
        nonAiLinesDeleted: 5,
        message: 'feat: add authentication',
        commitTs: '2026-01-01T10:00:00Z',
        createdAt: '2026-01-01T10:00:00Z',
      },
    ];

    it('should fetch commits without parameters', async () => {
      const mockResponse: PaginatedResponse<CursorCommit> = {
        data: mockCommits,
        pagination: {
          page: 1,
          pageSize: 100,
          totalPages: 1,
          hasNextPage: false,
          hasPreviousPage: false,
        },
        params: { page: 1, pageSize: 100 },
      };

      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => mockResponse,
      });

      const result = await client.getCommits();

      expect(result.data).toEqual(mockCommits);
      expect(result.pagination.page).toBe(1);
      expect(global.fetch).toHaveBeenCalledWith(
        'http://localhost:8080/analytics/ai-code/commits',
        expect.any(Object)
      );
    });

    it('should include query parameters when provided', async () => {
      const mockResponse: PaginatedResponse<CursorCommit> = {
        data: mockCommits,
        pagination: {
          page: 2,
          pageSize: 50,
          totalPages: 5,
          hasNextPage: true,
          hasPreviousPage: true,
        },
        params: {
          from: '2026-01-01',
          to: '2026-01-07',
          page: 2,
          pageSize: 50,
        },
      };

      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => mockResponse,
      });

      const result = await client.getCommits({
        from: '2026-01-01',
        to: '2026-01-07',
        page: 2,
        page_size: 50,
      });

      expect(result.pagination.page).toBe(2);
      expect(result.pagination.pageSize).toBe(50);

      const callUrl = (global.fetch as jest.Mock).mock.calls[0][0];
      expect(callUrl).toContain('from=2026-01-01');
      expect(callUrl).toContain('to=2026-01-07');
      expect(callUrl).toContain('page=2');
      expect(callUrl).toContain('page_size=50');
    });

    it('should handle pagination correctly', async () => {
      const mockResponse: PaginatedResponse<CursorCommit> = {
        data: mockCommits,
        pagination: {
          page: 1,
          pageSize: 100,
          totalPages: 3,
          hasNextPage: true,
          hasPreviousPage: false,
        },
        params: { page: 1, pageSize: 100 },
      };

      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => mockResponse,
      });

      const result = await client.getCommits({ page: 1, page_size: 100 });

      expect(result.pagination.hasNextPage).toBe(true);
      expect(result.pagination.hasPreviousPage).toBe(false);
    });

    it('should filter by user_id when provided', async () => {
      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({
          data: mockCommits,
          pagination: {
            page: 1,
            pageSize: 100,
            totalPages: 1,
            hasNextPage: false,
            hasPreviousPage: false,
          },
          params: { user_id: 'alice@example.com' },
        }),
      });

      await client.getCommits({ user_id: 'alice@example.com' });

      const callUrl = (global.fetch as jest.Mock).mock.calls[0][0];
      // URL encoding: @ becomes %40
      expect(callUrl).toContain('user_id=alice%40example.com');
    });

    it('should filter by repo_name when provided', async () => {
      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({
          data: mockCommits,
          pagination: {
            page: 1,
            pageSize: 100,
            totalPages: 1,
            hasNextPage: false,
            hasPreviousPage: false,
          },
          params: { repo_name: 'acme/platform' },
        }),
      });

      await client.getCommits({ repo_name: 'acme/platform' });

      const callUrl = (global.fetch as jest.Mock).mock.calls[0][0];
      expect(callUrl).toContain('repo_name=acme%2Fplatform');
    });
  });

  describe('error handling', () => {
    it('should handle 429 Rate Limit errors', async () => {
      // Mock all retries to return 429
      (global.fetch as jest.Mock).mockResolvedValue({
        ok: false,
        status: 429,
        statusText: 'Too Many Requests',
        json: async () => ({ error: 'Rate limit exceeded' }),
      });

      await expect(client.getTeamMembers()).rejects.toThrow('Too Many Requests');
    });

    it('should handle 500 Internal Server Error', async () => {
      // Mock all retries to return 500
      (global.fetch as jest.Mock).mockResolvedValue({
        ok: false,
        status: 500,
        statusText: 'Internal Server Error',
        json: async () => ({ error: 'Server error' }),
      });

      await expect(client.getTeamMembers()).rejects.toThrow('Internal Server Error');
    });

    it('should handle malformed JSON response', async () => {
      (global.fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => {
          throw new Error('Invalid JSON');
        },
      });

      await expect(client.getTeamMembers()).rejects.toThrow();
    });
  });

  describe('retry logic', () => {
    beforeEach(() => {
      // Ensure clean slate for retry tests
      jest.clearAllMocks();
      jest.resetAllMocks();
    });

    it('should retry on network failure', async () => {
      // First two calls fail, third succeeds
      let callCount = 0;
      (global.fetch as jest.Mock).mockImplementation(() => {
        callCount++;
        if (callCount <= 2) {
          return Promise.reject(new Error('Network error'));
        }
        return Promise.resolve({
          ok: true,
          status: 200,
          json: () => Promise.resolve({ data: [] }),
        });
      });

      const result = await client.getTeamMembers();

      expect(result).toEqual([]);
      expect(global.fetch).toHaveBeenCalledTimes(3);
    });

    it('should give up after max retries', async () => {
      (global.fetch as jest.Mock).mockImplementation(() =>
        Promise.reject(new Error('Persistent network error'))
      );

      await expect(client.getTeamMembers()).rejects.toThrow(
        'Persistent network error'
      );

      // Initial attempt + 2 retries = 3 total
      expect(global.fetch).toHaveBeenCalledTimes(3);
    });

    it('should not retry on 4xx client errors', async () => {
      (global.fetch as jest.Mock).mockImplementation(() =>
        Promise.resolve({
          ok: false,
          status: 400,
          statusText: 'Bad Request',
          json: () => Promise.resolve({ error: 'Invalid parameters' }),
        })
      );

      await expect(client.getTeamMembers()).rejects.toThrow('Bad Request');

      // Should not retry on client errors
      expect(global.fetch).toHaveBeenCalledTimes(1);
    });

    it('should retry on 5xx server errors', async () => {
      (global.fetch as jest.Mock)
        .mockResolvedValueOnce({
          ok: false,
          status: 503,
          statusText: 'Service Unavailable',
          json: async () => ({ error: 'Temporary error' }),
        })
        .mockResolvedValueOnce({
          ok: true,
          status: 200,
          json: async () => ({ data: [] }),
        });

      const result = await client.getTeamMembers();

      expect(result).toEqual([]);
      expect(global.fetch).toHaveBeenCalledTimes(2);
    });
  });

  describe('timeout handling', () => {
    it('should timeout long-running requests', async () => {
      // Create client with shorter timeout for faster test
      const fastClient = new CursorSimClient({
        ...mockConfig,
        timeout: 100, // 100ms timeout
        retryAttempts: 1, // No retries for faster test
      });

      // Mock fetch that takes longer than timeout
      let abortCalled = false;
      (global.fetch as jest.Mock).mockImplementationOnce(
        (_url: string, options?: { signal?: AbortSignal }) => {
          return new Promise((_resolve, reject) => {
            // Listen for abort signal
            if (options?.signal) {
              options.signal.addEventListener('abort', () => {
                abortCalled = true;
                const error = new Error('The operation was aborted');
                error.name = 'AbortError';
                reject(error);
              });
            }
          });
        }
      );

      await expect(fastClient.getTeamMembers()).rejects.toThrow('timeout');
      expect(abortCalled).toBe(true);
    }, 5000); // Test timeout of 5 seconds to ensure it completes
  });
});
