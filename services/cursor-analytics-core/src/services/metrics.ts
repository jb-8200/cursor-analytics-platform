/**
 * Metrics Service
 *
 * Provides metric calculation functions for analytics
 */

import { PrismaClient } from '../generated/prisma';

export interface MemberStats {
  suggestionsShown: number;
  suggestionsAccepted: number;
}

export interface MemberAIVelocityStats {
  aiLinesAdded: number;
  totalLinesAdded: number;
}

export interface DateRange {
  from: Date;
  to: Date;
}

export type DateRangePreset = 'TODAY' | 'LAST_7_DAYS' | 'LAST_30_DAYS' | 'LAST_90_DAYS' | 'THIS_WEEK' | 'THIS_MONTH';

export class MetricsService {
  constructor(private db: PrismaClient) {}

  /**
   * Calculate acceptance rate as a percentage
   * @returns Percentage rounded to 2 decimal places, or null if no suggestions
   */
  calculateAcceptanceRate(accepted: number, shown: number): number | null {
    if (shown === 0) return null;
    return Math.round((accepted / shown) * 10000) / 100;
  }

  /**
   * Calculate AI velocity as a percentage
   * @returns Percentage rounded to 2 decimal places, or null if no lines
   */
  calculateAIVelocity(aiLines: number, totalLines: number): number | null {
    if (totalLines === 0) return null;
    return Math.round((aiLines / totalLines) * 10000) / 100;
  }

  /**
   * Calculate team acceptance rate using weighted average
   * @returns Percentage rounded to 2 decimal places, or null if no data
   */
  calculateTeamAcceptanceRate(memberStats: MemberStats[]): number | null {
    if (memberStats.length === 0) return null;

    const totalShown = memberStats.reduce((sum, stat) => sum + stat.suggestionsShown, 0);
    const totalAccepted = memberStats.reduce((sum, stat) => sum + stat.suggestionsAccepted, 0);

    return this.calculateAcceptanceRate(totalAccepted, totalShown);
  }

  /**
   * Calculate team AI velocity using weighted average
   * @returns Percentage rounded to 2 decimal places, or null if no data
   */
  calculateTeamAIVelocity(memberStats: MemberAIVelocityStats[]): number | null {
    if (memberStats.length === 0) return null;

    const totalAILines = memberStats.reduce((sum, stat) => sum + stat.aiLinesAdded, 0);
    const totalLines = memberStats.reduce((sum, stat) => sum + stat.totalLinesAdded, 0);

    return this.calculateAIVelocity(totalAILines, totalLines);
  }

  /**
   * Get active developers (those with events in the given date range)
   */
  async getActiveDevelopers(dateRange?: DateRange): Promise<any[]> {
    const developers = await this.db.developer.findMany();

    if (!dateRange) {
      const activeDevelopers = [];
      for (const dev of developers) {
        const eventCount = await this.db.usageEvent.count({
          where: { developerId: dev.id },
        });
        if (eventCount > 0) {
          activeDevelopers.push(dev);
        }
      }
      return activeDevelopers;
    }

    const activeDevelopers = [];
    for (const dev of developers) {
      const eventCount = await this.db.usageEvent.count({
        where: {
          developerId: dev.id,
          eventTimestamp: {
            gte: dateRange.from,
            lte: dateRange.to,
          },
        },
      });
      if (eventCount > 0) {
        activeDevelopers.push(dev);
      }
    }

    return activeDevelopers;
  }

  /**
   * Expand a date range preset into actual dates
   */
  expandDateRangePreset(preset: DateRangePreset): DateRange {
    const now = new Date();
    const today = new Date();
    today.setHours(0, 0, 0, 0);

    switch (preset) {
      case 'TODAY':
        return { from: today, to: now };

      case 'LAST_7_DAYS': {
        const from = new Date();
        from.setDate(from.getDate() - 7);
        from.setHours(0, 0, 0, 0);
        return { from, to: now };
      }

      case 'LAST_30_DAYS': {
        const from = new Date();
        from.setDate(from.getDate() - 30);
        from.setHours(0, 0, 0, 0);
        return { from, to: now };
      }

      case 'LAST_90_DAYS': {
        const from = new Date();
        from.setDate(from.getDate() - 90);
        from.setHours(0, 0, 0, 0);
        return { from, to: now };
      }

      case 'THIS_WEEK': {
        const dayOfWeek = today.getDay();
        const monday = new Date(today);
        monday.setDate(today.getDate() - dayOfWeek + (dayOfWeek === 0 ? -6 : 1));
        monday.setHours(0, 0, 0, 0);
        return { from: monday, to: now };
      }

      case 'THIS_MONTH': {
        const firstDayOfMonth = new Date(today.getFullYear(), today.getMonth(), 1);
        return { from: firstDayOfMonth, to: now };
      }

      default:
        throw new Error(`Invalid date range preset: ${preset}`);
    }
  }

  /**
   * Filter events by date range
   */
  filterEventsByDateRange(events: any[], dateRange?: DateRange): any[] {
    if (!dateRange) return events;

    return events.filter(
      (event) =>
        event.eventTimestamp >= dateRange.from && event.eventTimestamp <= dateRange.to
    );
  }
}
