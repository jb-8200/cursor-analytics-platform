/**
 * Developer Resolver Tests
 *
 * Tests for GraphQL resolver functions that query developer data.
 * Uses mocked Prisma Client to avoid database dependencies.
 */

/* eslint-disable @typescript-eslint/unbound-method */

import { developerResolvers } from '../developer';
import { GraphQLContext } from '../../context';
import { PrismaClient } from '../../../generated/prisma';
import { CursorSimClient } from '../../../ingestion/client';

// Mock Prisma Client
const mockPrismaClient = {
  developer: {
    findUnique: jest.fn(),
    findMany: jest.fn(),
    count: jest.fn(),
  },
  usageEvent: {
    groupBy: jest.fn(),
    findMany: jest.fn(),
  },
} as unknown as PrismaClient;

// Mock CursorSimClient
const mockSimClient = {} as CursorSimClient;

// Create mock context
function createMockContext(): GraphQLContext {
  return {
    db: mockPrismaClient,
    simClient: mockSimClient,
    requestId: 'test-request-id',
  };
}

describe('Developer Resolvers', () => {
  beforeEach(() => {
    // Reset all mocks before each test
    jest.clearAllMocks();
  });

  describe('Query.developer', () => {
    it('should return a developer by ID', async () => {
      const mockDeveloper = {
        id: 'dev-uuid-1',
        externalId: 'ext-dev-1',
        name: 'Alice Developer',
        email: 'alice@example.com',
        team: 'backend',
        seniority: 'senior',
        createdAt: new Date('2025-01-01T00:00:00Z'),
        updatedAt: new Date('2025-01-01T00:00:00Z'),
      };

      (mockPrismaClient.developer.findUnique as jest.Mock).mockResolvedValue(mockDeveloper);

      const context = createMockContext();
      const result = await developerResolvers.Query.developer(null, { id: 'dev-uuid-1' }, context);

      expect(result).toEqual(mockDeveloper);
      expect(mockPrismaClient.developer.findUnique).toHaveBeenCalledWith({
        where: { id: 'dev-uuid-1' },
      });
    });

    it('should return null if developer not found', async () => {
      (mockPrismaClient.developer.findUnique as jest.Mock).mockResolvedValue(null);

      const context = createMockContext();
      const result = await developerResolvers.Query.developer(null, { id: 'nonexistent' }, context);

      expect(result).toBeNull();
    });
  });

  describe('Query.developers', () => {
    it('should return paginated list of developers', async () => {
      const mockDevelopers = [
        {
          id: 'dev-1',
          externalId: 'ext-1',
          name: 'Alice',
          email: 'alice@example.com',
          team: 'backend',
          seniority: 'senior',
          createdAt: new Date('2025-01-01'),
          updatedAt: new Date('2025-01-01'),
        },
        {
          id: 'dev-2',
          externalId: 'ext-2',
          name: 'Bob',
          email: 'bob@example.com',
          team: 'frontend',
          seniority: 'mid',
          createdAt: new Date('2025-01-02'),
          updatedAt: new Date('2025-01-02'),
        },
      ];

      (mockPrismaClient.developer.findMany as jest.Mock).mockResolvedValue(mockDevelopers);
      (mockPrismaClient.developer.count as jest.Mock).mockResolvedValue(2);

      const context = createMockContext();
      const result = await developerResolvers.Query.developers(
        null,
        { limit: 50, offset: 0, sortBy: 'name', sortOrder: 'asc' },
        context,
      );

      expect(result.nodes).toEqual(mockDevelopers);
      expect(result.totalCount).toBe(2);
      expect(result.pageInfo.hasNextPage).toBe(false);
      expect(result.pageInfo.hasPreviousPage).toBe(false);
    });

    it('should filter developers by team', async () => {
      const mockDevelopers = [
        {
          id: 'dev-1',
          externalId: 'ext-1',
          name: 'Alice',
          email: 'alice@example.com',
          team: 'backend',
          seniority: 'senior',
          createdAt: new Date('2025-01-01'),
          updatedAt: new Date('2025-01-01'),
        },
      ];

      (mockPrismaClient.developer.findMany as jest.Mock).mockResolvedValue(mockDevelopers);
      (mockPrismaClient.developer.count as jest.Mock).mockResolvedValue(1);

      const context = createMockContext();
      const result = await developerResolvers.Query.developers(
        null,
        { team: 'backend', limit: 50, offset: 0, sortBy: 'name', sortOrder: 'asc' },
        context,
      );

      expect(result.nodes).toEqual(mockDevelopers);
      expect(mockPrismaClient.developer.findMany).toHaveBeenCalledWith({
        where: { team: 'backend' },
        orderBy: { name: 'asc' },
        take: 50,
        skip: 0,
      });
    });

    it('should filter developers by seniority', async () => {
      const mockDevelopers = [
        {
          id: 'dev-1',
          externalId: 'ext-1',
          name: 'Alice',
          email: 'alice@example.com',
          team: 'backend',
          seniority: 'senior',
          createdAt: new Date('2025-01-01'),
          updatedAt: new Date('2025-01-01'),
        },
      ];

      (mockPrismaClient.developer.findMany as jest.Mock).mockResolvedValue(mockDevelopers);
      (mockPrismaClient.developer.count as jest.Mock).mockResolvedValue(1);

      const context = createMockContext();
      await developerResolvers.Query.developers(
        null,
        { seniority: 'senior', limit: 50, offset: 0, sortBy: 'name', sortOrder: 'asc' },
        context,
      );

      expect(mockPrismaClient.developer.findMany).toHaveBeenCalledWith({
        where: { seniority: 'senior' },
        orderBy: { name: 'asc' },
        take: 50,
        skip: 0,
      });
    });

    it('should handle pagination correctly', async () => {
      const mockDevelopers = Array.from({ length: 10 }, (_, i) => ({
        id: `dev-${i + 21}`,
        externalId: `ext-${i + 21}`,
        name: `Developer ${i + 21}`,
        email: `dev${i + 21}@example.com`,
        team: 'engineering',
        seniority: 'mid',
        createdAt: new Date('2025-01-01'),
        updatedAt: new Date('2025-01-01'),
      }));

      (mockPrismaClient.developer.findMany as jest.Mock).mockResolvedValue(mockDevelopers);
      (mockPrismaClient.developer.count as jest.Mock).mockResolvedValue(100);

      const context = createMockContext();
      const result = await developerResolvers.Query.developers(
        null,
        { limit: 10, offset: 20, sortBy: 'name', sortOrder: 'asc' },
        context,
      );

      expect(result.totalCount).toBe(100);
      expect(result.pageInfo.hasNextPage).toBe(true); // offset 20 + limit 10 = 30 < 100
      expect(result.pageInfo.hasPreviousPage).toBe(true); // offset 20 > 0
    });

    it('should support sorting by different fields', async () => {
      const mockDevelopers = [
        {
          id: 'dev-1',
          externalId: 'ext-1',
          name: 'Alice',
          email: 'alice@example.com',
          team: 'backend',
          seniority: 'senior',
          createdAt: new Date('2025-01-01'),
          updatedAt: new Date('2025-01-01'),
        },
      ];

      (mockPrismaClient.developer.findMany as jest.Mock).mockResolvedValue(mockDevelopers);
      (mockPrismaClient.developer.count as jest.Mock).mockResolvedValue(1);

      const context = createMockContext();
      await developerResolvers.Query.developers(
        null,
        { limit: 50, offset: 0, sortBy: 'createdAt', sortOrder: 'desc' },
        context,
      );

      expect(mockPrismaClient.developer.findMany).toHaveBeenCalledWith({
        where: {},
        orderBy: { createdAt: 'desc' },
        take: 50,
        skip: 0,
      });
    });
  });

  describe('Developer.stats', () => {
    it('should calculate usage stats for a developer', async () => {
      const mockDeveloper = {
        id: 'dev-1',
        externalId: 'ext-1',
        name: 'Alice',
        email: 'alice@example.com',
        team: 'backend',
        seniority: 'senior',
        createdAt: new Date('2025-01-01'),
        updatedAt: new Date('2025-01-01'),
      };

      const mockGroupBy = [
        {
          eventType: 'cpp_suggestion_shown',
          _count: { id: 100 },
          _sum: { linesAdded: 0, linesDeleted: 0 },
        },
        {
          eventType: 'cpp_suggestion_accepted',
          _count: { id: 75 },
          _sum: { linesAdded: 1500, linesDeleted: 0 },
        },
        {
          eventType: 'chat_message',
          _count: { id: 20 },
          _sum: { linesAdded: 0, linesDeleted: 0 },
        },
        {
          eventType: 'cmd_k_prompt',
          _count: { id: 10 },
          _sum: { linesAdded: 0, linesDeleted: 0 },
        },
      ];

      (mockPrismaClient.usageEvent.groupBy as jest.Mock).mockResolvedValue(mockGroupBy);

      const context = createMockContext();
      const result = await developerResolvers.Developer.stats(mockDeveloper, {}, context);

      expect(result.totalSuggestions).toBe(100);
      expect(result.acceptedSuggestions).toBe(75);
      expect(result.acceptanceRate).toBe(75.0); // 75/100 * 100
      expect(result.chatInteractions).toBe(20);
      expect(result.cmdKUsages).toBe(10);
      expect(result.aiLinesAdded).toBe(1500);
    });

    it('should filter stats by date range', async () => {
      const mockDeveloper = {
        id: 'dev-1',
        externalId: 'ext-1',
        name: 'Alice',
        email: 'alice@example.com',
        team: 'backend',
        seniority: 'senior',
        createdAt: new Date('2025-01-01'),
        updatedAt: new Date('2025-01-01'),
      };

      (mockPrismaClient.usageEvent.groupBy as jest.Mock).mockResolvedValue([]);

      const context = createMockContext();
      const dateRange = {
        from: '2025-01-01T00:00:00Z',
        to: '2025-01-31T23:59:59Z',
      };

      await developerResolvers.Developer.stats(mockDeveloper, { range: dateRange }, context);

      expect(mockPrismaClient.usageEvent.groupBy).toHaveBeenCalledWith({
        by: ['eventType'],
        where: {
          developerId: 'dev-1',
          eventTimestamp: {
            gte: new Date('2025-01-01T00:00:00Z'),
            lte: new Date('2025-01-31T23:59:59Z'),
          },
        },
        _count: { id: true },
        _sum: {
          linesAdded: true,
          linesDeleted: true,
        },
      });
    });

    it('should return null acceptance rate when no suggestions', async () => {
      const mockDeveloper = {
        id: 'dev-1',
        externalId: 'ext-1',
        name: 'Alice',
        email: 'alice@example.com',
        team: 'backend',
        seniority: 'senior',
        createdAt: new Date('2025-01-01'),
        updatedAt: new Date('2025-01-01'),
      };

      (mockPrismaClient.usageEvent.groupBy as jest.Mock).mockResolvedValue([]);

      const context = createMockContext();
      const result = await developerResolvers.Developer.stats(mockDeveloper, {}, context);

      expect(result.totalSuggestions).toBe(0);
      expect(result.acceptedSuggestions).toBe(0);
      expect(result.acceptanceRate).toBeNull();
    });

    it('should calculate total lines including human edits', async () => {
      const mockDeveloper = {
        id: 'dev-1',
        externalId: 'ext-1',
        name: 'Alice',
        email: 'alice@example.com',
        team: 'backend',
        seniority: 'senior',
        createdAt: new Date('2025-01-01'),
        updatedAt: new Date('2025-01-01'),
      };

      const mockGroupBy = [
        {
          eventType: 'cpp_suggestion_accepted',
          _count: { id: 50 },
          _sum: { linesAdded: 1000, linesDeleted: 100 },
        },
      ];

      // Additional query for total lines
      const mockTotalLines = [
        {
          _sum: { linesAdded: 1500, linesDeleted: 200 },
        },
      ];

      (mockPrismaClient.usageEvent.groupBy as jest.Mock)
        .mockResolvedValueOnce(mockGroupBy)
        .mockResolvedValueOnce(mockTotalLines);

      const context = createMockContext();
      const result = await developerResolvers.Developer.stats(mockDeveloper, {}, context);

      expect(result.totalLinesAdded).toBe(1500);
      expect(result.totalLinesDeleted).toBe(200);
      expect(result.aiLinesAdded).toBe(1000);
    });
  });

  describe('Developer.dailyStats', () => {
    it('should return daily statistics grouped by date', async () => {
      const mockDeveloper = {
        id: 'dev-1',
        externalId: 'ext-1',
        name: 'Alice',
        email: 'alice@example.com',
        team: 'backend',
        seniority: 'senior',
        createdAt: new Date('2025-01-01'),
        updatedAt: new Date('2025-01-01'),
      };

      // Create individual event records (not aggregated)
      const mockEvents = [
        // Day 1: 20 shown, 15 accepted
        ...Array.from({ length: 20 }, () => ({
          eventTimestamp: new Date('2025-01-01T10:00:00Z'),
          eventType: 'cpp_suggestion_shown',
          linesAdded: 0,
          linesDeleted: 0,
        })),
        ...Array.from({ length: 15 }, () => ({
          eventTimestamp: new Date('2025-01-01T10:00:00Z'),
          eventType: 'cpp_suggestion_accepted',
          linesAdded: 20, // 20 lines per acceptance = 300 total
          linesDeleted: 0,
        })),
        // Day 2: 25 shown, 20 accepted
        ...Array.from({ length: 25 }, () => ({
          eventTimestamp: new Date('2025-01-02T10:00:00Z'),
          eventType: 'cpp_suggestion_shown',
          linesAdded: 0,
          linesDeleted: 0,
        })),
        ...Array.from({ length: 20 }, () => ({
          eventTimestamp: new Date('2025-01-02T10:00:00Z'),
          eventType: 'cpp_suggestion_accepted',
          linesAdded: 20, // 20 lines per acceptance = 400 total
          linesDeleted: 0,
        })),
      ];

      (mockPrismaClient.usageEvent.findMany as jest.Mock).mockResolvedValue(mockEvents);

      const context = createMockContext();
      const result = await developerResolvers.Developer.dailyStats(mockDeveloper, {}, context);

      expect(result).toHaveLength(2);
      expect(result[0].date).toBe('2025-01-01');
      expect(result[0].suggestionsShown).toBe(20);
      expect(result[0].suggestionsAccepted).toBe(15);
      expect(result[0].acceptanceRate).toBe(75.0);
      expect(result[0].linesAdded).toBe(300);
    });

    it('should filter daily stats by date range', async () => {
      const mockDeveloper = {
        id: 'dev-1',
        externalId: 'ext-1',
        name: 'Alice',
        email: 'alice@example.com',
        team: 'backend',
        seniority: 'senior',
        createdAt: new Date('2025-01-01'),
        updatedAt: new Date('2025-01-01'),
      };

      (mockPrismaClient.usageEvent.findMany as jest.Mock).mockResolvedValue([]);

      const context = createMockContext();
      const dateRange = {
        from: '2025-01-01T00:00:00Z',
        to: '2025-01-07T23:59:59Z',
      };

      await developerResolvers.Developer.dailyStats(mockDeveloper, { range: dateRange }, context);

      expect(mockPrismaClient.usageEvent.findMany).toHaveBeenCalledWith({
        where: {
          developerId: 'dev-1',
          eventTimestamp: {
            gte: new Date('2025-01-01T00:00:00Z'),
            lte: new Date('2025-01-07T23:59:59Z'),
          },
        },
        orderBy: { eventTimestamp: 'asc' },
      });
    });
  });
});
