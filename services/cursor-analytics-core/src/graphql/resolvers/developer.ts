/**
 * Developer Resolvers
 *
 * GraphQL resolvers for querying developer data and calculating usage statistics.
 */

import { GraphQLContext } from '../context';

// Type definitions for resolver arguments
export interface DeveloperArgs {
  id: string;
}

export interface DevelopersArgs {
  team?: string;
  seniority?: string;
  limit?: number;
  offset?: number;
  sortBy?: string;
  sortOrder?: string;
}

export interface DateRangeInput {
  from: string;
  to: string;
}

export interface StatsArgs {
  range?: DateRangeInput;
}

export interface DailyStatsArgs {
  range?: DateRangeInput;
}

// Type for Developer from Prisma
export interface Developer {
  id: string;
  externalId: string;
  name: string;
  email: string;
  team: string;
  seniority: string | null;
  createdAt: Date;
  updatedAt: Date;
}

// Response types
export interface UsageStats {
  totalSuggestions: number;
  acceptedSuggestions: number;
  acceptanceRate: number | null;
  chatInteractions: number;
  cmdKUsages: number;
  totalLinesAdded: number;
  totalLinesDeleted: number;
  aiLinesAdded: number;
  aiVelocity: number | null;
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

export interface DeveloperConnection {
  nodes: Developer[];
  totalCount: number;
  pageInfo: {
    hasNextPage: boolean;
    hasPreviousPage: boolean;
    startCursor: string | null;
    endCursor: string | null;
  };
}

/**
 * Calculate acceptance rate
 * Returns null if no suggestions to avoid division by zero
 */
function calculateAcceptanceRate(accepted: number, shown: number): number | null {
  if (shown === 0) return null;
  return Math.round((accepted / shown) * 10000) / 100; // Round to 2 decimal places
}

/**
 * Calculate AI velocity
 * Returns null if no lines added to avoid division by zero
 */
function calculateAIVelocity(aiLines: number, totalLines: number): number | null {
  if (totalLines === 0) return null;
  return Math.round((aiLines / totalLines) * 10000) / 100;
}

export const developerResolvers = {
  Query: {
    /**
     * Get a single developer by ID
     */
    developer: async (
      _parent: unknown,
      args: DeveloperArgs,
      context: GraphQLContext,
    ): Promise<Developer | null> => {
      return context.db.developer.findUnique({
        where: { id: args.id },
      });
    },

    /**
     * List developers with filtering and pagination
     */
    developers: async (
      _parent: unknown,
      args: DevelopersArgs,
      context: GraphQLContext,
    ): Promise<DeveloperConnection> => {
      const { team, seniority, limit = 50, offset = 0, sortBy = 'name', sortOrder = 'asc' } = args;

      // Build where clause
      const where: {
        team?: string;
        seniority?: string;
      } = {};
      if (team) where.team = team;
      if (seniority) where.seniority = seniority;

      // Build orderBy clause
      const orderBy: Record<string, string> = {};
      orderBy[sortBy] = sortOrder;

      // Execute queries
      const [nodes, totalCount] = await Promise.all([
        context.db.developer.findMany({
          where,
          orderBy,
          take: limit,
          skip: offset,
        }),
        context.db.developer.count({ where }),
      ]);

      // Calculate pagination info
      const hasNextPage = offset + limit < totalCount;
      const hasPreviousPage = offset > 0;

      return {
        nodes,
        totalCount,
        pageInfo: {
          hasNextPage,
          hasPreviousPage,
          startCursor: nodes.length > 0 ? nodes[0].id : null,
          endCursor: nodes.length > 0 ? nodes[nodes.length - 1].id : null,
        },
      };
    },
  },

  Developer: {
    /**
     * Calculate aggregated usage statistics for a developer
     */
    stats: async (
      parent: Developer,
      args: StatsArgs,
      context: GraphQLContext,
    ): Promise<UsageStats> => {
      // Build where clause
      const where: {
        developerId: string;
        eventTimestamp?: {
          gte: Date;
          lte: Date;
        };
      } = {
        developerId: parent.id,
      };

      // Add date range filter if provided
      if (args.range) {
        where.eventTimestamp = {
          gte: new Date(args.range.from),
          lte: new Date(args.range.to),
        };
      }

      // Aggregate stats by event type
      const eventStats = await context.db.usageEvent.groupBy({
        by: ['eventType'],
        where,
        _count: { id: true },
        _sum: {
          linesAdded: true,
          linesDeleted: true,
        },
      });

      // Also get total lines (including human edits)
      const totalLinesStats = await context.db.usageEvent.groupBy({
        by: ['developerId'],
        where: {
          developerId: parent.id,
          ...(args.range
            ? {
                eventTimestamp: {
                  gte: new Date(args.range.from),
                  lte: new Date(args.range.to),
                },
              }
            : {}),
        },
        _sum: {
          linesAdded: true,
          linesDeleted: true,
        },
      });

      // Extract values from grouped results
      const suggestionShown =
        eventStats.find((s) => s.eventType === 'cpp_suggestion_shown')?._count.id ?? 0;
      const suggestionAccepted =
        eventStats.find((s) => s.eventType === 'cpp_suggestion_accepted')?._count.id ?? 0;
      const chatMessages = eventStats.find((s) => s.eventType === 'chat_message')?._count.id ?? 0;
      const cmdKUsages = eventStats.find((s) => s.eventType === 'cmd_k_prompt')?._count.id ?? 0;

      const aiLinesAdded =
        eventStats.find((s) => s.eventType === 'cpp_suggestion_accepted')?._sum.linesAdded ?? 0;
      const totalLinesAdded = totalLinesStats[0]?._sum.linesAdded ?? 0;
      const totalLinesDeleted = totalLinesStats[0]?._sum.linesDeleted ?? 0;

      return {
        totalSuggestions: suggestionShown,
        acceptedSuggestions: suggestionAccepted,
        acceptanceRate: calculateAcceptanceRate(suggestionAccepted, suggestionShown),
        chatInteractions: chatMessages,
        cmdKUsages,
        totalLinesAdded,
        totalLinesDeleted,
        aiLinesAdded,
        aiVelocity: calculateAIVelocity(aiLinesAdded, totalLinesAdded),
      };
    },

    /**
     * Get daily statistics breakdown for a developer
     */
    dailyStats: async (
      parent: Developer,
      args: DailyStatsArgs,
      context: GraphQLContext,
    ): Promise<DailyStats[]> => {
      // Build where clause
      const where: {
        developerId: string;
        eventTimestamp?: {
          gte: Date;
          lte: Date;
        };
      } = {
        developerId: parent.id,
      };

      // Add date range filter if provided
      if (args.range) {
        where.eventTimestamp = {
          gte: new Date(args.range.from),
          lte: new Date(args.range.to),
        };
      }

      // Fetch all usage events for the developer
      const events = await context.db.usageEvent.findMany({
        where,
        orderBy: { eventTimestamp: 'asc' },
      });

      // Group events by date
      const dailyMap = new Map<
        string,
        {
          suggestionsShown: number;
          suggestionsAccepted: number;
          chatInteractions: number;
          cmdKUsages: number;
          linesAdded: number;
          linesDeleted: number;
          aiLinesAdded: number;
        }
      >();

      for (const event of events) {
        const dateKey = event.eventTimestamp.toISOString().split('T')[0];

        if (!dailyMap.has(dateKey)) {
          dailyMap.set(dateKey, {
            suggestionsShown: 0,
            suggestionsAccepted: 0,
            chatInteractions: 0,
            cmdKUsages: 0,
            linesAdded: 0,
            linesDeleted: 0,
            aiLinesAdded: 0,
          });
        }

        const stats = dailyMap.get(dateKey)!;

        if (event.eventType === 'cpp_suggestion_shown') {
          stats.suggestionsShown++;
        } else if (event.eventType === 'cpp_suggestion_accepted') {
          stats.suggestionsAccepted++;
          stats.aiLinesAdded += event.linesAdded;
        } else if (event.eventType === 'chat_message') {
          stats.chatInteractions++;
        } else if (event.eventType === 'cmd_k_prompt') {
          stats.cmdKUsages++;
        }

        stats.linesAdded += event.linesAdded;
        stats.linesDeleted += event.linesDeleted;
      }

      // Convert map to array and calculate acceptance rates
      const dailyStatsArray: DailyStats[] = [];
      for (const [date, stats] of dailyMap.entries()) {
        dailyStatsArray.push({
          date,
          suggestionsShown: stats.suggestionsShown,
          suggestionsAccepted: stats.suggestionsAccepted,
          acceptanceRate: calculateAcceptanceRate(
            stats.suggestionsAccepted,
            stats.suggestionsShown,
          ),
          chatInteractions: stats.chatInteractions,
          cmdKUsages: stats.cmdKUsages,
          linesAdded: stats.linesAdded,
          linesDeleted: stats.linesDeleted,
          aiLinesAdded: stats.aiLinesAdded,
        });
      }

      return dailyStatsArray;
    },
  },
};
