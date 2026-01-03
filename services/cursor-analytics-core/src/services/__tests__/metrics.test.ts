/**
 * Metrics Service Tests
 *
 * Unit tests for metric calculation functions.
 */

import { MetricsService } from '../metrics';
import { PrismaClient } from '../../generated/prisma';

jest.mock('../../generated/prisma', () => ({
  PrismaClient: jest.fn(),
}));

describe('MetricsService', () => {
  let metricsService: MetricsService;
  let mockDb: jest.Mocked<PrismaClient>;

  beforeEach(() => {
    mockDb = {
      usageEvent: {
        findMany: jest.fn(),
        groupBy: jest.fn(),
        count: jest.fn(),
      },
      developer: {
        findMany: jest.fn(),
        count: jest.fn(),
      },
    } as unknown as jest.Mocked<PrismaClient>;

    metricsService = new MetricsService(mockDb);
  });

  describe('calculateAcceptanceRate', () => {
    it('should calculate acceptance rate correctly', () => {
      const result = metricsService.calculateAcceptanceRate(75, 100);
      expect(result).toBe(75.0);
    });

    it('should return null when total suggestions is 0', () => {
      const result = metricsService.calculateAcceptanceRate(0, 0);
      expect(result).toBeNull();
    });

    it('should round to 2 decimal places', () => {
      const result = metricsService.calculateAcceptanceRate(2, 3);
      expect(result).toBe(66.67);
    });

    it('should handle 100% acceptance', () => {
      const result = metricsService.calculateAcceptanceRate(50, 50);
      expect(result).toBe(100.0);
    });

    it('should handle 0% acceptance', () => {
      const result = metricsService.calculateAcceptanceRate(0, 100);
      expect(result).toBe(0.0);
    });
  });

  describe('calculateAIVelocity', () => {
    it('should calculate AI velocity correctly', () => {
      const result = metricsService.calculateAIVelocity(500, 2000);
      expect(result).toBe(25.0);
    });

    it('should return null when total lines is 0', () => {
      const result = metricsService.calculateAIVelocity(0, 0);
      expect(result).toBeNull();
    });

    it('should round to 2 decimal places', () => {
      const result = metricsService.calculateAIVelocity(1, 3);
      expect(result).toBe(33.33);
    });

    it('should handle 100% AI velocity', () => {
      const result = metricsService.calculateAIVelocity(1000, 1000);
      expect(result).toBe(100.0);
    });

    it('should handle 0% AI velocity', () => {
      const result = metricsService.calculateAIVelocity(0, 1000);
      expect(result).toBe(0.0);
    });
  });

  describe('calculateTeamAcceptanceRate', () => {
    it('should calculate weighted team acceptance rate', () => {
      const memberStats = [
        { suggestionsShown: 100, suggestionsAccepted: 75 },
        { suggestionsShown: 200, suggestionsAccepted: 150 },
        { suggestionsShown: 50, suggestionsAccepted: 40 },
      ];

      const result = metricsService.calculateTeamAcceptanceRate(memberStats);
      expect(result).toBe(75.71);
    });

    it('should return null when no suggestions shown', () => {
      const memberStats = [
        { suggestionsShown: 0, suggestionsAccepted: 0 },
        { suggestionsShown: 0, suggestionsAccepted: 0 },
      ];

      const result = metricsService.calculateTeamAcceptanceRate(memberStats);
      expect(result).toBeNull();
    });

    it('should handle empty member list', () => {
      const result = metricsService.calculateTeamAcceptanceRate([]);
      expect(result).toBeNull();
    });

    it('should use weighted average not simple average', () => {
      const memberStats = [
        { suggestionsShown: 1, suggestionsAccepted: 1 },
        { suggestionsShown: 99, suggestionsAccepted: 50 },
      ];

      const result = metricsService.calculateTeamAcceptanceRate(memberStats);
      expect(result).toBe(51.0);
    });
  });

  describe('calculateTeamAIVelocity', () => {
    it('should calculate weighted team AI velocity', () => {
      const memberStats = [
        { aiLinesAdded: 500, totalLinesAdded: 2000 },
        { aiLinesAdded: 300, totalLinesAdded: 1000 },
        { aiLinesAdded: 200, totalLinesAdded: 500 },
      ];

      const result = metricsService.calculateTeamAIVelocity(memberStats);
      expect(result).toBe(28.57);
    });

    it('should return null when no lines added', () => {
      const memberStats = [
        { aiLinesAdded: 0, totalLinesAdded: 0 },
        { aiLinesAdded: 0, totalLinesAdded: 0 },
      ];

      const result = metricsService.calculateTeamAIVelocity(memberStats);
      expect(result).toBeNull();
    });

    it('should handle empty member list', () => {
      const result = metricsService.calculateTeamAIVelocity([]);
      expect(result).toBeNull();
    });
  });

  describe('getActiveDevelopers', () => {
    it('should return developers with events in date range', async () => {
      const mockDevelopers = [
        { id: 'dev1', name: 'Alice' },
        { id: 'dev2', name: 'Bob' },
      ];

      (mockDb.developer.findMany as jest.Mock).mockResolvedValue(mockDevelopers);
      (mockDb.usageEvent.count as jest.Mock).mockResolvedValueOnce(10).mockResolvedValueOnce(5);

      const from = new Date('2026-01-01');
      const to = new Date('2026-01-07');

      const result = await metricsService.getActiveDevelopers({ from, to });

      expect(result).toEqual(mockDevelopers);
      expect(mockDb.developer.findMany).toHaveBeenCalled();
    });

    it('should filter out developers with no events', async () => {
      const mockDevelopers = [
        { id: 'dev1', name: 'Alice' },
        { id: 'dev2', name: 'Bob' },
        { id: 'dev3', name: 'Charlie' },
      ];

      (mockDb.developer.findMany as jest.Mock).mockResolvedValue(mockDevelopers);
      (mockDb.usageEvent.count as jest.Mock)
        .mockResolvedValueOnce(10)
        .mockResolvedValueOnce(5)
        .mockResolvedValueOnce(0);

      const from = new Date('2026-01-01');
      const to = new Date('2026-01-07');

      const result = await metricsService.getActiveDevelopers({ from, to });

      expect(result).toHaveLength(2);
      expect(result.map((d) => d.id)).toEqual(['dev1', 'dev2']);
    });
  });

  describe('expandDateRangePreset', () => {
    it('should expand TODAY preset correctly', () => {
      const result = metricsService.expandDateRangePreset('TODAY');

      const today = new Date();
      today.setHours(0, 0, 0, 0);

      expect(result.from.toDateString()).toBe(today.toDateString());
      expect(result.to.getTime()).toBeGreaterThanOrEqual(Date.now() - 1000);
    });

    it('should expand THIS_WEEK preset correctly', () => {
      const result = metricsService.expandDateRangePreset('THIS_WEEK');
      const now = new Date();
      const dayOfWeek = now.getDay();
      const expectedFrom = new Date(now);
      expectedFrom.setDate(now.getDate() - dayOfWeek);
      expectedFrom.setHours(0, 0, 0, 0);
      expect(result.from.toDateString()).toBe(expectedFrom.toDateString());
    });

    it('should expand THIS_MONTH preset correctly', () => {
      const result = metricsService.expandDateRangePreset('THIS_MONTH');
      const now = new Date();
      const expectedFrom = new Date(now.getFullYear(), now.getMonth(), 1);
      expect(result.from.toDateString()).toBe(expectedFrom.toDateString());
    });

    it('should expand LAST_7_DAYS preset correctly', () => {
      const result = metricsService.expandDateRangePreset('LAST_7_DAYS');

      const expectedFrom = new Date();
      expectedFrom.setDate(expectedFrom.getDate() - 7);
      expectedFrom.setHours(0, 0, 0, 0);

      expect(result.from.toDateString()).toBe(expectedFrom.toDateString());
    });

    it('should expand LAST_30_DAYS preset correctly', () => {
      const result = metricsService.expandDateRangePreset('LAST_30_DAYS');
      const expectedFrom = new Date();
      expectedFrom.setDate(expectedFrom.getDate() - 30);
      expectedFrom.setHours(0, 0, 0, 0);
      expect(result.from.toDateString()).toBe(expectedFrom.toDateString());
    });

    it('should expand LAST_90_DAYS preset correctly', () => {
      const result = metricsService.expandDateRangePreset('LAST_90_DAYS');
      const expectedFrom = new Date();
      expectedFrom.setDate(expectedFrom.getDate() - 90);
      expectedFrom.setHours(0, 0, 0, 0);
      expect(result.from.toDateString()).toBe(expectedFrom.toDateString());
    });

    it('should throw error for invalid preset', () => {
      expect(() => {
        metricsService.expandDateRangePreset('INVALID' as any);
      }).toThrow('Invalid date range preset: INVALID');
    });
  });

  describe('filterEventsByDateRange', () => {
    it('should filter events within date range', () => {
      const events = [
        { id: '1', eventTimestamp: new Date('2026-01-01') },
        { id: '2', eventTimestamp: new Date('2026-01-05') },
        { id: '3', eventTimestamp: new Date('2026-01-10') },
      ];

      const result = metricsService.filterEventsByDateRange(events, {
        from: new Date('2026-01-03'),
        to: new Date('2026-01-08'),
      });

      expect(result).toHaveLength(1);
      expect(result[0].id).toBe('2');
    });

    it('should return all events when no date range provided', () => {
      const events = [
        { id: '1', eventTimestamp: new Date('2026-01-01') },
        { id: '2', eventTimestamp: new Date('2026-01-05') },
      ];

      const result = metricsService.filterEventsByDateRange(events);

      expect(result).toHaveLength(2);
    });
  });
});
