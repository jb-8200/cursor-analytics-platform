/**
 * Commit Resolver Tests
 *
 * Tests for the commits query resolver that fetches usage events
 * (which represent commits in our schema).
 */

import { commitResolvers } from '../commit';
import { GraphQLContext } from '../../context';
import { PrismaClient } from '@prisma/client';
import { CursorSimClient } from '../../../ingestion/client';

// Mock Prisma client
const mockUsageEvent = {
  findMany: jest.fn(),
  count: jest.fn(),
};

const mockDeveloper = {
  findUnique: jest.fn(),
};

// eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
const mockPrismaClient = {
  usageEvent: mockUsageEvent,
  developer: mockDeveloper,
} as unknown as PrismaClient;

// Mock CursorSimClient
const mockSimClient = {} as CursorSimClient;

// Helper to create context
function createContext(): GraphQLContext {
  return {
    // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
    db: mockPrismaClient,
    simClient: mockSimClient,
  };
}

describe('Commit Resolvers', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('Query.commits', () => {
    it('should return paginated commits', async () => {
      // Arrange
      const mockCommits = [
        {
          id: 'event-1',
          externalId: 'commit-1',
          developerId: 'dev-1',
          eventType: 'cpp_suggestion_accepted',
          eventTimestamp: new Date('2026-01-01T10:00:00Z'),
          linesAdded: 50,
          linesDeleted: 10,
          modelUsed: 'claude-sonnet-4',
          accepted: true,
          tokensInput: 100,
          tokensOutput: 200,
          createdAt: new Date('2026-01-01T10:00:00Z'),
        },
        {
          id: 'event-2',
          externalId: 'commit-2',
          developerId: 'dev-2',
          eventType: 'cpp_suggestion_accepted',
          eventTimestamp: new Date('2026-01-01T11:00:00Z'),
          linesAdded: 30,
          linesDeleted: 5,
          modelUsed: 'claude-sonnet-4',
          accepted: true,
          tokensInput: 80,
          tokensOutput: 150,
          createdAt: new Date('2026-01-01T11:00:00Z'),
        },
      ];

      mockUsageEvent.findMany.mockResolvedValue(mockCommits);
      mockUsageEvent.count.mockResolvedValue(2);

      // Act
      const result = await commitResolvers.Query.commits(
        {},
        { limit: 10, offset: 0 },
        createContext(),
      );

      // Assert
      expect(result.nodes).toHaveLength(2);
      expect(result.totalCount).toBe(2);
      expect(result.pageInfo.hasNextPage).toBe(false);
      expect(result.pageInfo.hasPreviousPage).toBe(false);
      expect(mockUsageEvent.findMany).toHaveBeenCalledWith({
        where: {
          eventType: 'cpp_suggestion_accepted',
        },
        orderBy: { eventTimestamp: 'desc' },
        take: 10,
        skip: 0,
        include: {
          developer: true,
        },
      });
    });

    it('should filter by userId', async () => {
      // Arrange
      const mockCommits = [
        {
          id: 'event-1',
          externalId: 'commit-1',
          developerId: 'dev-1',
          eventType: 'cpp_suggestion_accepted',
          eventTimestamp: new Date('2026-01-01T10:00:00Z'),
          linesAdded: 50,
          linesDeleted: 10,
          modelUsed: 'claude-sonnet-4',
          accepted: true,
          tokensInput: 100,
          tokensOutput: 200,
          createdAt: new Date('2026-01-01T10:00:00Z'),
          developer: {
            id: 'dev-1',
            externalId: 'alice',
            name: 'Alice',
            email: 'alice@example.com',
            team: 'backend',
            seniority: 'senior',
            createdAt: new Date('2026-01-01T00:00:00Z'),
            updatedAt: new Date('2026-01-01T00:00:00Z'),
          },
        },
      ];

      mockUsageEvent.findMany.mockResolvedValue(mockCommits);
      mockUsageEvent.count.mockResolvedValue(1);

      // Act
      const result = await commitResolvers.Query.commits(
        {},
        { userId: 'dev-1', limit: 10, offset: 0 },
        createContext(),
      );

      // Assert
      expect(result.nodes).toHaveLength(1);
      expect(result.nodes[0].developerId).toBe('dev-1');
      expect(mockUsageEvent.findMany).toHaveBeenCalledWith({
        where: {
          eventType: 'cpp_suggestion_accepted',
          developerId: 'dev-1',
        },
        orderBy: { eventTimestamp: 'desc' },
        take: 10,
        skip: 0,
        include: {
          developer: true,
        },
      });
    });

    it('should filter by team', async () => {
      // Arrange
      const mockCommits = [
        {
          id: 'event-1',
          externalId: 'commit-1',
          developerId: 'dev-1',
          eventType: 'cpp_suggestion_accepted',
          eventTimestamp: new Date('2026-01-01T10:00:00Z'),
          linesAdded: 50,
          linesDeleted: 10,
          modelUsed: 'claude-sonnet-4',
          accepted: true,
          tokensInput: 100,
          tokensOutput: 200,
          createdAt: new Date('2026-01-01T10:00:00Z'),
          developer: {
            id: 'dev-1',
            externalId: 'alice',
            name: 'Alice',
            email: 'alice@example.com',
            team: 'backend',
            seniority: 'senior',
            createdAt: new Date('2026-01-01T00:00:00Z'),
            updatedAt: new Date('2026-01-01T00:00:00Z'),
          },
        },
      ];

      mockUsageEvent.findMany.mockResolvedValue(mockCommits);
      mockUsageEvent.count.mockResolvedValue(1);

      // Act
      const result = await commitResolvers.Query.commits(
        {},
        { team: 'backend', limit: 10, offset: 0 },
        createContext(),
      );

      // Assert
      expect(result.nodes).toHaveLength(1);
      expect(result.nodes[0].developer?.team).toBe('backend');
      expect(mockUsageEvent.findMany).toHaveBeenCalledWith({
        where: {
          eventType: 'cpp_suggestion_accepted',
          developer: { team: 'backend' },
        },
        orderBy: { eventTimestamp: 'desc' },
        take: 10,
        skip: 0,
        include: {
          developer: true,
        },
      });
    });

    it('should filter by date range', async () => {
      // Arrange
      const mockCommits = [
        {
          id: 'event-1',
          externalId: 'commit-1',
          developerId: 'dev-1',
          eventType: 'cpp_suggestion_accepted',
          eventTimestamp: new Date('2026-01-15T10:00:00Z'),
          linesAdded: 50,
          linesDeleted: 10,
          modelUsed: 'claude-sonnet-4',
          accepted: true,
          tokensInput: 100,
          tokensOutput: 200,
          createdAt: new Date('2026-01-15T10:00:00Z'),
          developer: {
            id: 'dev-1',
            externalId: 'alice',
            name: 'Alice',
            email: 'alice@example.com',
            team: 'backend',
            seniority: 'senior',
            createdAt: new Date('2026-01-01T00:00:00Z'),
            updatedAt: new Date('2026-01-01T00:00:00Z'),
          },
        },
      ];

      mockUsageEvent.findMany.mockResolvedValue(mockCommits);
      mockUsageEvent.count.mockResolvedValue(1);

      // Act
      const result = await commitResolvers.Query.commits(
        {},
        {
          dateRange: {
            from: '2026-01-01T00:00:00Z',
            to: '2026-01-31T23:59:59Z',
          },
          limit: 10,
          offset: 0,
        },
        createContext(),
      );

      // Assert
      expect(result.nodes).toHaveLength(1);
      expect(mockUsageEvent.findMany).toHaveBeenCalledWith({
        where: {
          eventType: 'cpp_suggestion_accepted',
          eventTimestamp: {
            gte: new Date('2026-01-01T00:00:00Z'),
            lte: new Date('2026-01-31T23:59:59Z'),
          },
        },
        orderBy: { eventTimestamp: 'desc' },
        take: 10,
        skip: 0,
        include: {
          developer: true,
        },
      });
    });

    it('should support sorting by timestamp ascending', async () => {
      // Arrange
      mockUsageEvent.findMany.mockResolvedValue([]);
      mockUsageEvent.count.mockResolvedValue(0);

      // Act
      await commitResolvers.Query.commits(
        {},
        { sortBy: 'timestamp', sortOrder: 'asc', limit: 10, offset: 0 },
        createContext(),
      );

      // Assert
      expect(mockUsageEvent.findMany).toHaveBeenCalledWith({
        where: {
          eventType: 'cpp_suggestion_accepted',
        },
        orderBy: { eventTimestamp: 'asc' },
        take: 10,
        skip: 0,
        include: {
          developer: true,
        },
      });
    });

    it('should support sorting by author name', async () => {
      // Arrange
      mockUsageEvent.findMany.mockResolvedValue([]);
      mockUsageEvent.count.mockResolvedValue(0);

      // Act
      await commitResolvers.Query.commits(
        {},
        { sortBy: 'author', sortOrder: 'asc', limit: 10, offset: 0 },
        createContext(),
      );

      // Assert
      expect(mockUsageEvent.findMany).toHaveBeenCalledWith({
        where: {
          eventType: 'cpp_suggestion_accepted',
        },
        orderBy: { developer: { name: 'asc' } },
        take: 10,
        skip: 0,
        include: {
          developer: true,
        },
      });
    });

    it('should support pagination with hasNextPage', async () => {
      // Arrange
      const mockCommits = [
        {
          id: 'event-1',
          externalId: 'commit-1',
          developerId: 'dev-1',
          eventType: 'cpp_suggestion_accepted',
          eventTimestamp: new Date('2026-01-01T10:00:00Z'),
          linesAdded: 50,
          linesDeleted: 10,
          modelUsed: 'claude-sonnet-4',
          accepted: true,
          tokensInput: 100,
          tokensOutput: 200,
          createdAt: new Date('2026-01-01T10:00:00Z'),
          developer: {
            id: 'dev-1',
            externalId: 'alice',
            name: 'Alice',
            email: 'alice@example.com',
            team: 'backend',
            seniority: 'senior',
            createdAt: new Date('2026-01-01T00:00:00Z'),
            updatedAt: new Date('2026-01-01T00:00:00Z'),
          },
        },
      ];

      mockUsageEvent.findMany.mockResolvedValue(mockCommits);
      mockUsageEvent.count.mockResolvedValue(25);

      // Act
      const result = await commitResolvers.Query.commits(
        {},
        { limit: 10, offset: 0 },
        createContext(),
      );

      // Assert
      expect(result.pageInfo.hasNextPage).toBe(true);
      expect(result.pageInfo.hasPreviousPage).toBe(false);
    });

    it('should support pagination with hasPreviousPage', async () => {
      // Arrange
      const mockCommits = [
        {
          id: 'event-1',
          externalId: 'commit-1',
          developerId: 'dev-1',
          eventType: 'cpp_suggestion_accepted',
          eventTimestamp: new Date('2026-01-01T10:00:00Z'),
          linesAdded: 50,
          linesDeleted: 10,
          modelUsed: 'claude-sonnet-4',
          accepted: true,
          tokensInput: 100,
          tokensOutput: 200,
          createdAt: new Date('2026-01-01T10:00:00Z'),
          developer: {
            id: 'dev-1',
            externalId: 'alice',
            name: 'Alice',
            email: 'alice@example.com',
            team: 'backend',
            seniority: 'senior',
            createdAt: new Date('2026-01-01T00:00:00Z'),
            updatedAt: new Date('2026-01-01T00:00:00Z'),
          },
        },
      ];

      mockUsageEvent.findMany.mockResolvedValue(mockCommits);
      mockUsageEvent.count.mockResolvedValue(25);

      // Act
      const result = await commitResolvers.Query.commits(
        {},
        { limit: 10, offset: 10 },
        createContext(),
      );

      // Assert
      expect(result.pageInfo.hasNextPage).toBe(true);
      expect(result.pageInfo.hasPreviousPage).toBe(true);
    });

    it('should combine multiple filters', async () => {
      // Arrange
      mockUsageEvent.findMany.mockResolvedValue([]);
      mockUsageEvent.count.mockResolvedValue(0);

      // Act
      await commitResolvers.Query.commits(
        {},
        {
          userId: 'dev-1',
          team: 'backend',
          dateRange: {
            from: '2026-01-01T00:00:00Z',
            to: '2026-01-31T23:59:59Z',
          },
          limit: 10,
          offset: 0,
        },
        createContext(),
      );

      // Assert
      expect(mockUsageEvent.findMany).toHaveBeenCalledWith({
        where: {
          eventType: 'cpp_suggestion_accepted',
          developerId: 'dev-1',
          developer: { team: 'backend' },
          eventTimestamp: {
            gte: new Date('2026-01-01T00:00:00Z'),
            lte: new Date('2026-01-31T23:59:59Z'),
          },
        },
        orderBy: { eventTimestamp: 'desc' },
        take: 10,
        skip: 0,
        include: {
          developer: true,
        },
      });
    });

    it('should return empty array when no commits found', async () => {
      // Arrange
      mockUsageEvent.findMany.mockResolvedValue([]);
      mockUsageEvent.count.mockResolvedValue(0);

      // Act
      const result = await commitResolvers.Query.commits(
        {},
        { limit: 10, offset: 0 },
        createContext(),
      );

      // Assert
      expect(result.nodes).toEqual([]);
      expect(result.totalCount).toBe(0);
      expect(result.pageInfo.hasNextPage).toBe(false);
      expect(result.pageInfo.hasPreviousPage).toBe(false);
      expect(result.pageInfo.startCursor).toBeNull();
      expect(result.pageInfo.endCursor).toBeNull();
    });
  });

  describe('Commit.author field resolver', () => {
    it('should return the associated developer', () => {
      // Arrange
      const mockCommit = {
        id: 'event-1',
        externalId: 'commit-1',
        developerId: 'dev-1',
        eventType: 'cpp_suggestion_accepted',
        eventTimestamp: new Date('2026-01-01T10:00:00Z'),
        linesAdded: 50,
        linesDeleted: 10,
        modelUsed: 'claude-sonnet-4',
        accepted: true,
        tokensInput: 100,
        tokensOutput: 200,
        createdAt: new Date('2026-01-01T10:00:00Z'),
        developer: {
          id: 'dev-1',
          externalId: 'alice',
          name: 'Alice',
          email: 'alice@example.com',
          team: 'backend',
          seniority: 'senior',
          createdAt: new Date('2026-01-01T00:00:00Z'),
          updatedAt: new Date('2026-01-01T00:00:00Z'),
        },
      };

      // Act
      const result = commitResolvers.Commit.author(mockCommit, {}, createContext());

      // Assert
      expect(result).toEqual(mockCommit.developer);
    });
  });
});
