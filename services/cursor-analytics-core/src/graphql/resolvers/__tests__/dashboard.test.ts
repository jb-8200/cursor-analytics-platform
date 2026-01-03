/**
 * Dashboard Resolvers Tests
 *
 * Unit tests for dashboard summary, team stats, and teams queries.
 */

import { dashboardResolvers, TeamStats } from '../dashboard';
import { GraphQLContext } from '../../context';
import { PrismaClient } from '../../../generated/prisma';
import { CursorSimClient } from '../../../ingestion/client';

// Mock dependencies
jest.mock('../../../generated/prisma', () => ({
  PrismaClient: jest.fn(),
}));

jest.mock('../../../ingestion/client');

describe('Dashboard Resolvers', () => {
  let mockDb: jest.Mocked<PrismaClient>;
  let mockSimClient: jest.Mocked<CursorSimClient>;
  let context: GraphQLContext;

  beforeEach(() => {
    mockDb = {
      developer: {
        findMany: jest.fn(),
        count: jest.fn(),
      },
      usageEvent: {
        findMany: jest.fn(),
        count: jest.fn(),
        groupBy: jest.fn(),
      },
    } as unknown as jest.Mocked<PrismaClient>;

    mockSimClient = {} as jest.Mocked<CursorSimClient>;

    context = {
      db: mockDb,
      simClient: mockSimClient,
    };
  });

  describe('Query.dashboardSummary', () => {
    it('should return dashboard KPIs with TODAY preset', async () => {
      // Mock developers
      (mockDb.developer.count as jest.Mock).mockResolvedValue(10);

      const mockDevelopers = [
        { id: 'dev1', name: 'Alice', team: 'TeamA', externalId: 'alice', email: 'alice@example.com' },
        { id: 'dev2', name: 'Bob', team: 'TeamB', externalId: 'bob', email: 'bob@example.com' },
      ];
      (mockDb.developer.findMany as jest.Mock).mockResolvedValue(mockDevelopers);

      // Mock events for active developers
      (mockDb.usageEvent.count as jest.Mock)
        .mockResolvedValueOnce(10) // dev1 has events
        .mockResolvedValueOnce(5); // dev2 has events

      // Mock usage events for TODAY
      const mockEvents = [
        {
          id: 'e1',
          developerId: 'dev1',
          eventType: 'cpp_suggestion_shown',
          eventTimestamp: new Date(),
          linesAdded: 10,
          aiLinesAdded: 8,
        },
        {
          id: 'e2',
          developerId: 'dev1',
          eventType: 'cpp_suggestion_accepted',
          eventTimestamp: new Date(),
          linesAdded: 10,
          aiLinesAdded: 8,
        },
      ];
      (mockDb.usageEvent.findMany as jest.Mock).mockResolvedValue(mockEvents);

      const result = await dashboardResolvers.Query.dashboardSummary(
        {},
        { preset: 'TODAY' },
        context
      );

      expect(result).toBeDefined();
      expect(result.totalDevelopers).toBe(10);
      expect(result.activeDevelopers).toBe(2);
      expect(result.totalSuggestionsToday).toBeGreaterThanOrEqual(0);
      expect(result.totalAcceptedToday).toBeGreaterThanOrEqual(0);
    });

    it('should use date range when provided', async () => {
      (mockDb.developer.count as jest.Mock).mockResolvedValue(5);
      (mockDb.developer.findMany as jest.Mock).mockResolvedValue([]);
      (mockDb.usageEvent.findMany as jest.Mock).mockResolvedValue([]);

      const from = new Date('2026-01-01');
      const to = new Date('2026-01-07');

      const result = await dashboardResolvers.Query.dashboardSummary(
        {},
        { range: { from: from.toISOString(), to: to.toISOString() } },
        context
      );

      expect(result).toBeDefined();
      // eslint-disable-next-line @typescript-eslint/unbound-method
      expect(mockDb.usageEvent.findMany).toHaveBeenCalledWith(
        // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
        expect.objectContaining({
          // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
          where: expect.objectContaining({
            // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
            eventTimestamp: {
              gte: expect.any(Date),
              lte: expect.any(Date),
            },
          }),
        })
      );
    });

    it('should default to LAST_7_DAYS when no range/preset provided', async () => {
      (mockDb.developer.count as jest.Mock).mockResolvedValue(3);
      (mockDb.developer.findMany as jest.Mock).mockResolvedValue([]);
      (mockDb.usageEvent.findMany as jest.Mock).mockResolvedValue([]);

      const result = await dashboardResolvers.Query.dashboardSummary(
        {},
        {},
        context
      );

      expect(result).toBeDefined();
      expect(result.totalDevelopers).toBe(3);
    });

    it('should calculate acceptance rate correctly', async () => {
      (mockDb.developer.count as jest.Mock).mockResolvedValue(1);
      (mockDb.developer.findMany as jest.Mock).mockResolvedValue([
        { id: 'dev1', name: 'Alice', team: 'TeamA' },
      ]);
      (mockDb.usageEvent.count as jest.Mock).mockResolvedValue(1);

      const mockEvents = [
        { id: 'e1', eventType: 'cpp_suggestion_shown', eventTimestamp: new Date() },
        { id: 'e2', eventType: 'cpp_suggestion_shown', eventTimestamp: new Date() },
        { id: 'e3', eventType: 'cpp_suggestion_accepted', eventTimestamp: new Date() },
      ];
      (mockDb.usageEvent.findMany as jest.Mock).mockResolvedValue(mockEvents);

      const result = await dashboardResolvers.Query.dashboardSummary(
        {},
        { preset: 'TODAY' },
        context
      );

      expect(result.totalSuggestionsToday).toBe(2);
      expect(result.totalAcceptedToday).toBe(1);
      expect(result.overallAcceptanceRate).toBe(50.0);
    });

    it('should return null acceptance rate when no suggestions', async () => {
      (mockDb.developer.count as jest.Mock).mockResolvedValue(1);
      (mockDb.developer.findMany as jest.Mock).mockResolvedValue([]);
      (mockDb.usageEvent.findMany as jest.Mock).mockResolvedValue([]);

      const result = await dashboardResolvers.Query.dashboardSummary(
        {},
        { preset: 'TODAY' },
        context
      );

      expect(result.overallAcceptanceRate).toBeNull();
    });

    it('should include team comparison', async () => {
      (mockDb.developer.count as jest.Mock).mockResolvedValue(2);
      (mockDb.developer.findMany as jest.Mock).mockResolvedValue([
        { id: 'dev1', name: 'Alice', team: 'TeamA', externalId: 'alice', email: 'alice@example.com' },
        { id: 'dev2', name: 'Bob', team: 'TeamB', externalId: 'bob', email: 'bob@example.com' },
      ]);
      (mockDb.usageEvent.count as jest.Mock).mockResolvedValue(1);
      (mockDb.usageEvent.findMany as jest.Mock).mockResolvedValue([]);

      const result = await dashboardResolvers.Query.dashboardSummary(
        {},
        { preset: 'TODAY' },
        context
      );

      expect(result.teamComparison).toBeDefined();
      expect(Array.isArray(result.teamComparison)).toBe(true);
    });

    it('should include daily trend', async () => {
      (mockDb.developer.count as jest.Mock).mockResolvedValue(1);
      (mockDb.developer.findMany as jest.Mock).mockResolvedValue([]);
      (mockDb.usageEvent.findMany as jest.Mock).mockResolvedValue([]);

      const result = await dashboardResolvers.Query.dashboardSummary(
        {},
        { preset: 'LAST_7_DAYS' },
        context
      );

      expect(result.dailyTrend).toBeDefined();
      expect(Array.isArray(result.dailyTrend)).toBe(true);
    });
  });

  describe('Query.teamStats', () => {
    it('should return stats for a specific team', async () => {
      const mockTeamMembers = [
        { id: 'dev1', name: 'Alice', team: 'TeamA', externalId: 'alice', email: 'alice@example.com' },
        { id: 'dev2', name: 'Bob', team: 'TeamA', externalId: 'bob', email: 'bob@example.com' },
      ];

      (mockDb.developer.findMany as jest.Mock).mockResolvedValue(mockTeamMembers);
      (mockDb.usageEvent.count as jest.Mock)
        .mockResolvedValueOnce(10)
        .mockResolvedValueOnce(5);

      const mockEvents = [
        { id: 'e1', developerId: 'dev1', eventType: 'cpp_suggestion_shown', linesAdded: 10, aiLinesAdded: 8 },
        { id: 'e2', developerId: 'dev2', eventType: 'cpp_suggestion_accepted', linesAdded: 5, aiLinesAdded: 4 },
      ];
      (mockDb.usageEvent.findMany as jest.Mock).mockResolvedValue(mockEvents);

      const result = await dashboardResolvers.Query.teamStats(
        {},
        { teamName: 'TeamA' },
        context
      );

      expect(result).toBeDefined();
      expect(result!.teamName).toBe('TeamA');
      expect(result!.memberCount).toBe(2);
      expect(result!.activeMemberCount).toBe(2);
    });

    it('should return null for non-existent team', async () => {
      (mockDb.developer.findMany as jest.Mock).mockResolvedValue([]);

      const result = await dashboardResolvers.Query.teamStats(
        {},
        { teamName: 'NonExistentTeam' },
        context
      );

      expect(result).toBeNull();
    });

    it('should calculate weighted team acceptance rate', async () => {
      const mockTeamMembers = [
        { id: 'dev1', name: 'Alice', team: 'TeamA' },
        { id: 'dev2', name: 'Bob', team: 'TeamA' },
      ];

      (mockDb.developer.findMany as jest.Mock).mockResolvedValue(mockTeamMembers);
      (mockDb.usageEvent.count as jest.Mock).mockResolvedValue(1);

      const mockEvents = [
        { id: 'e1', developerId: 'dev1', eventType: 'cpp_suggestion_shown' },
        { id: 'e2', developerId: 'dev1', eventType: 'cpp_suggestion_shown' },
        { id: 'e3', developerId: 'dev1', eventType: 'cpp_suggestion_accepted' },
        { id: 'e4', developerId: 'dev2', eventType: 'cpp_suggestion_shown' },
        { id: 'e5', developerId: 'dev2', eventType: 'cpp_suggestion_shown' },
        { id: 'e6', developerId: 'dev2', eventType: 'cpp_suggestion_accepted' },
      ];
      (mockDb.usageEvent.findMany as jest.Mock).mockResolvedValue(mockEvents);

      const result = await dashboardResolvers.Query.teamStats(
        {},
        { teamName: 'TeamA' },
        context
      );

      expect(result!.totalSuggestions).toBe(4);
      expect(result!.totalAccepted).toBe(2);
      expect(result!.averageAcceptanceRate).toBe(50.0);
    });

    it('should identify top performer', async () => {
      const mockTeamMembers = [
        { id: 'dev1', name: 'Alice', team: 'TeamA', externalId: 'alice', email: 'alice@example.com' },
        { id: 'dev2', name: 'Bob', team: 'TeamA', externalId: 'bob', email: 'bob@example.com' },
      ];

      (mockDb.developer.findMany as jest.Mock).mockResolvedValue(mockTeamMembers);
      (mockDb.usageEvent.count as jest.Mock).mockResolvedValue(1);

      const mockEvents = [
        { id: 'e1', developerId: 'dev1', eventType: 'cpp_suggestion_accepted', linesAdded: 100, aiLinesAdded: 80 },
        { id: 'e2', developerId: 'dev2', eventType: 'cpp_suggestion_accepted', linesAdded: 50, aiLinesAdded: 40 },
      ];
      (mockDb.usageEvent.findMany as jest.Mock).mockResolvedValue(mockEvents);

      const result = await dashboardResolvers.Query.teamStats(
        {},
        { teamName: 'TeamA' },
        context
      );

      expect(result!.topPerformer).toBeDefined();
      expect(result!.topPerformer?.id).toBe('dev1');
    });
  });

  describe('Query.teams', () => {
    it('should return stats for all teams', async () => {
      const mockDevelopers = [
        { id: 'dev1', name: 'Alice', team: 'TeamA', externalId: 'alice', email: 'alice@example.com' },
        { id: 'dev2', name: 'Bob', team: 'TeamB', externalId: 'bob', email: 'bob@example.com' },
        { id: 'dev3', name: 'Charlie', team: 'TeamA', externalId: 'charlie', email: 'charlie@example.com' },
      ];

      (mockDb.developer.findMany as jest.Mock).mockResolvedValue(mockDevelopers);
      (mockDb.usageEvent.count as jest.Mock).mockResolvedValue(1);
      (mockDb.usageEvent.findMany as jest.Mock).mockResolvedValue([]);

      const result = await dashboardResolvers.Query.teams({}, {}, context);

      expect(result).toBeDefined();
      expect(Array.isArray(result)).toBe(true);
      expect(result.length).toBe(2); // TeamA and TeamB
      expect(result.map((t: TeamStats) => t.teamName).sort()).toEqual(['TeamA', 'TeamB']);
    });

    it('should return empty array when no developers', async () => {
      (mockDb.developer.findMany as jest.Mock).mockResolvedValue([]);

      const result = await dashboardResolvers.Query.teams({}, {}, context);

      expect(result).toEqual([]);
    });

    it('should include member counts for each team', async () => {
      const allDevelopers = [
        { id: 'dev1', name: 'Alice', team: 'TeamA', externalId: 'alice', email: 'alice@example.com' },
        { id: 'dev2', name: 'Bob', team: 'TeamA', externalId: 'bob', email: 'bob@example.com' },
        { id: 'dev3', name: 'Charlie', team: 'TeamB', externalId: 'charlie', email: 'charlie@example.com' },
      ];

      // Mock for the initial call to get all developers
      (mockDb.developer.findMany as jest.Mock)
        .mockResolvedValueOnce(allDevelopers)
        // Mock for TeamA members
        .mockResolvedValueOnce(allDevelopers.filter((d) => d.team === 'TeamA'))
        // Mock for TeamB members
        .mockResolvedValueOnce(allDevelopers.filter((d) => d.team === 'TeamB'));

      (mockDb.usageEvent.count as jest.Mock).mockResolvedValue(1);
      (mockDb.usageEvent.findMany as jest.Mock).mockResolvedValue([]);

      const result = await dashboardResolvers.Query.teams({}, {}, context);

      const teamA = result.find((t: TeamStats) => t.teamName === 'TeamA');
      const teamB = result.find((t: TeamStats) => t.teamName === 'TeamB');

      expect(teamA?.memberCount).toBe(2);
      expect(teamB?.memberCount).toBe(1);
    });
  });
});
