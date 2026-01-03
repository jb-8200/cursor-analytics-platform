/**
 * Dashboard Resolvers
 *
 * GraphQL resolvers for dashboard summary, team stats, and teams queries.
 */

import { GraphQLContext } from '../context';
import { MetricsService } from '../../services/metrics';
import type { Developer, UsageEvent } from '../../generated/prisma';

// Type definitions for resolver arguments
export interface DashboardSummaryArgs {
  range?: DateRangeInput;
  preset?: DateRangePreset;
}

export interface TeamStatsArgs {
  teamName: string;
}

export interface DateRangeInput {
  from: string;
  to: string;
}

export type DateRangePreset = 'TODAY' | 'LAST_7_DAYS' | 'LAST_30_DAYS' | 'LAST_90_DAYS' | 'THIS_WEEK' | 'THIS_MONTH';

// Response types
export interface DashboardKPI {
  totalDevelopers: number;
  activeDevelopers: number;
  overallAcceptanceRate: number | null;
  totalSuggestionsToday: number;
  totalAcceptedToday: number;
  aiVelocityToday: number | null;
  teamComparison: TeamStats[];
  dailyTrend: DailyStats[];
}

export interface TeamStats {
  teamName: string;
  memberCount: number;
  activeMemberCount: number;
  averageAcceptanceRate: number | null;
  totalSuggestions: number;
  totalAccepted: number;
  chatInteractions: number;
  aiVelocity: number | null;
  topPerformer: Developer | null;
}

export interface DailyStats {
  date: string;
  suggestionsShown: number;
  suggestionsAccepted: number;
  acceptanceRate: number | null;
  chatInteractions: number;
  cmdKUsages: number;
  linesAdded: number;
  linesDeleted: number;
  aiLinesAdded: number;
}

/**
 * Helper: Calculate team statistics for a given team
 */
async function calculateTeamStats(
  teamName: string,
  context: GraphQLContext,
  dateRange?: { from: Date; to: Date }
): Promise<TeamStats | null> {
  const metricsService = new MetricsService(context.db);

  // Get all team members
  const teamMembers = await context.db.developer.findMany({
    where: { team: teamName },
  });

  if (teamMembers.length === 0) {
    return null;
  }

  // Get active members (those with events in the date range)
  const activeMemberIds: string[] = [];
  for (const member of teamMembers) {
    const eventCount = await context.db.usageEvent.count({
      where: {
        developerId: member.id,
        ...(dateRange && {
          eventTimestamp: {
            gte: dateRange.from,
            lte: dateRange.to,
          },
        }),
      },
    });

    if (eventCount > 0) {
      activeMemberIds.push(member.id);
    }
  }

  // Get all events for the team
  const teamEvents = await context.db.usageEvent.findMany({
    where: {
      developerId: { in: teamMembers.map((m) => m.id) },
      ...(dateRange && {
        eventTimestamp: {
          gte: dateRange.from,
          lte: dateRange.to,
        },
      }),
    },
  });

  // Calculate aggregated stats
  const totalSuggestions = teamEvents.filter((e) => e.eventType === 'cpp_suggestion_shown').length;
  const totalAccepted = teamEvents.filter((e) => e.eventType === 'cpp_suggestion_accepted').length;
  const chatInteractions = teamEvents.filter((e) => e.eventType === 'chat_message').length;

  const totalLinesAdded = teamEvents.reduce((sum, e) => sum + e.linesAdded, 0);
  const totalAILines = teamEvents
    .filter((e) => e.accepted === true)
    .reduce((sum, e) => sum + e.linesAdded, 0);

  const averageAcceptanceRate = metricsService.calculateAcceptanceRate(totalAccepted, totalSuggestions);
  const aiVelocity = metricsService.calculateAIVelocity(totalAILines, totalLinesAdded);

  // Find top performer (by AI lines added)
  let topPerformer = null;
  if (activeMemberIds.length > 0) {
    const memberAILines = activeMemberIds.map((memberId) => {
      const memberEvents = teamEvents.filter((e) => e.developerId === memberId);
      const aiLines = memberEvents
        .filter((e) => e.accepted === true)
        .reduce((sum, e) => sum + e.linesAdded, 0);
      return { memberId, aiLines };
    });

    const topPerformerId = memberAILines.reduce((max, current) =>
      current.aiLines > max.aiLines ? current : max
    ).memberId;

    topPerformer = teamMembers.find((m) => m.id === topPerformerId) || null;
  }

  return {
    teamName,
    memberCount: teamMembers.length,
    activeMemberCount: activeMemberIds.length,
    averageAcceptanceRate,
    totalSuggestions,
    totalAccepted,
    chatInteractions,
    aiVelocity,
    topPerformer,
  };
}

/**
 * Helper: Calculate daily trend statistics
 */
function calculateDailyTrend(events: UsageEvent[], dateRange: { from: Date; to: Date }): DailyStats[] {
  const metricsService = new MetricsService(null as unknown as import('../../generated/prisma').PrismaClient); // Just for calculation helpers

  // Group events by date
  const eventsByDate = new Map<string, UsageEvent[]>();

  for (const event of events) {
    const dateKey = event.eventTimestamp.toISOString().split('T')[0];
    if (!eventsByDate.has(dateKey)) {
      eventsByDate.set(dateKey, []);
    }
    eventsByDate.get(dateKey)!.push(event);
  }

  // Generate daily stats for each day in range
  const dailyStats: DailyStats[] = [];
  const currentDate = new Date(dateRange.from);

  while (currentDate <= dateRange.to) {
    const dateKey = currentDate.toISOString().split('T')[0];
    const dayEvents = eventsByDate.get(dateKey) ?? [];

    const suggestionsShown = dayEvents.filter((e) => e.eventType === 'cpp_suggestion_shown').length;
    const suggestionsAccepted = dayEvents.filter((e) => e.eventType === 'cpp_suggestion_accepted').length;
    const chatInteractions = dayEvents.filter((e) => e.eventType === 'chat_message').length;
    const cmdKUsages = dayEvents.filter((e) => e.eventType === 'cmd_k_prompt').length;

    const linesAdded = dayEvents.reduce((sum, e) => sum + e.linesAdded, 0);
    const linesDeleted = dayEvents.reduce((sum, e) => sum + e.linesDeleted, 0);
    const aiLinesAdded = dayEvents
      .filter((e) => e.accepted === true)
      .reduce((sum, e) => sum + e.linesAdded, 0);

    const acceptanceRate = metricsService.calculateAcceptanceRate(suggestionsAccepted, suggestionsShown);

    dailyStats.push({
      date: currentDate.toISOString(),
      suggestionsShown,
      suggestionsAccepted,
      acceptanceRate,
      chatInteractions,
      cmdKUsages,
      linesAdded,
      linesDeleted,
      aiLinesAdded,
    });

    currentDate.setDate(currentDate.getDate() + 1);
  }

  return dailyStats;
}

/**
 * Dashboard resolvers
 */
export const dashboardResolvers = {
  Query: {
    /**
     * Get dashboard summary with KPIs
     */
    dashboardSummary: async (
      _parent: unknown,
      args: DashboardSummaryArgs,
      context: GraphQLContext
    ): Promise<DashboardKPI> => {
      const metricsService = new MetricsService(context.db);

      // Determine date range
      let dateRange: { from: Date; to: Date };
      if (args.range) {
        dateRange = {
          from: new Date(args.range.from),
          to: new Date(args.range.to),
        };
      } else if (args.preset) {
        dateRange = metricsService.expandDateRangePreset(args.preset);
      } else {
        // Default to LAST_7_DAYS
        dateRange = metricsService.expandDateRangePreset('LAST_7_DAYS');
      }

      // Get total developers
      const totalDevelopers = await context.db.developer.count();

      // Get active developers
      const activeDevelopers = await metricsService.getActiveDevelopers(dateRange);

      // Get all events in range
      const allEvents = await context.db.usageEvent.findMany({
        where: {
          eventTimestamp: {
            gte: dateRange.from,
            lte: dateRange.to,
          },
        },
      });

      // Calculate overall stats
      const totalSuggestionsToday = allEvents.filter((e) => e.eventType === 'cpp_suggestion_shown').length;
      const totalAcceptedToday = allEvents.filter((e) => e.eventType === 'cpp_suggestion_accepted').length;

      const overallAcceptanceRate = metricsService.calculateAcceptanceRate(
        totalAcceptedToday,
        totalSuggestionsToday
      );

      const totalLinesAdded = allEvents.reduce((sum, e) => sum + e.linesAdded, 0);
      const totalAILines = allEvents
        .filter((e) => e.accepted === true)
        .reduce((sum, e) => sum + e.linesAdded, 0);
      const aiVelocityToday = metricsService.calculateAIVelocity(totalAILines, totalLinesAdded);

      // Get team comparison
      const allDevelopers = await context.db.developer.findMany();
      const teamNames = [...new Set(allDevelopers.map((d) => d.team))];

      const teamComparison: TeamStats[] = [];
      for (const teamName of teamNames) {
        const teamStats = await calculateTeamStats(teamName, context, dateRange);
        if (teamStats) {
          teamComparison.push(teamStats);
        }
      }

      // Calculate daily trend
      const dailyTrend = calculateDailyTrend(allEvents, dateRange);

      return {
        totalDevelopers,
        activeDevelopers: activeDevelopers.length,
        overallAcceptanceRate,
        totalSuggestionsToday,
        totalAcceptedToday,
        aiVelocityToday,
        teamComparison,
        dailyTrend,
      };
    },

    /**
     * Get statistics for a specific team
     */
    teamStats: async (
      _parent: unknown,
      args: TeamStatsArgs,
      context: GraphQLContext
    ): Promise<TeamStats | null> => {
      return calculateTeamStats(args.teamName, context);
    },

    /**
     * Get statistics for all teams
     */
    teams: async (
      _parent: unknown,
      _args: unknown,
      context: GraphQLContext
    ): Promise<TeamStats[]> => {
      const allDevelopers = await context.db.developer.findMany();
      const teamNames = [...new Set(allDevelopers.map((d) => d.team))];

      const teams: TeamStats[] = [];
      for (const teamName of teamNames) {
        const teamStats = await calculateTeamStats(teamName, context);
        if (teamStats) {
          teams.push(teamStats);
        }
      }

      return teams;
    },
  },
};
