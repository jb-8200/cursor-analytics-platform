/**
 * Commit Resolvers
 *
 * GraphQL resolvers for querying commit data (usage events with eventType = 'cpp_suggestion_accepted').
 * In our schema, commits are represented by usage events that were accepted by developers.
 */

import { GraphQLContext } from '../context';
import { Developer } from './developer';

// Type definitions for resolver arguments
export interface CommitsArgs {
  userId?: string;
  team?: string;
  dateRange?: {
    from: string;
    to: string;
  };
  sortBy?: string;
  sortOrder?: string;
  limit?: number;
  offset?: number;
}

// Type for Commit (usage event)
export interface Commit {
  id: string;
  externalId: string;
  developerId: string;
  eventType: string;
  eventTimestamp: Date;
  linesAdded: number;
  linesDeleted: number;
  modelUsed: string | null;
  accepted: boolean | null;
  tokensInput: number;
  tokensOutput: number;
  createdAt: Date;
  developer?: Developer;
}

// Response types
export interface CommitConnection {
  nodes: Commit[];
  totalCount: number;
  pageInfo: {
    hasNextPage: boolean;
    hasPreviousPage: boolean;
    startCursor: string | null;
    endCursor: string | null;
  };
}

export const commitResolvers = {
  Query: {
    /**
     * List commits (usage events) with filtering and pagination
     * Returns accepted AI suggestions as commits
     */
    commits: async (
      _parent: unknown,
      args: CommitsArgs,
      context: GraphQLContext,
    ): Promise<CommitConnection> => {
      const {
        userId,
        team,
        dateRange,
        sortBy = 'timestamp',
        sortOrder = 'desc',
        limit = 50,
        offset = 0,
      } = args;

      // Build where clause
      interface WhereClause {
        eventType: string;
        developerId?: string;
        developer?: { team: string };
        eventTimestamp?: {
          gte: Date;
          lte: Date;
        };
      }

      const where: WhereClause = {
        eventType: 'cpp_suggestion_accepted', // Only fetch commits (accepted suggestions)
      };

      if (userId) {
        where.developerId = userId;
      }

      if (team) {
        where.developer = { team };
      }

      if (dateRange) {
        where.eventTimestamp = {
          gte: new Date(dateRange.from),
          lte: new Date(dateRange.to),
        };
      }

      // Build orderBy clause
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const orderBy: any = sortBy === 'author'
        ? { developer: { name: sortOrder } }
        : { eventTimestamp: sortOrder };

      // Execute queries
      const [nodes, totalCount] = await Promise.all([
        context.db.usageEvent.findMany({
          where,
          orderBy,
          take: limit,
          skip: offset,
          include: {
            developer: true,
          },
        }),
        context.db.usageEvent.count({ where }),
      ]);

      // Calculate pagination info
      const hasNextPage = offset + limit < totalCount;
      const hasPreviousPage = offset > 0;

      return {
        nodes: nodes as Commit[],
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

  Commit: {
    /**
     * Map eventTimestamp to timestamp for GraphQL schema
     */
    timestamp: (parent: Commit): Date => {
      return parent.eventTimestamp;
    },

    /**
     * Get the developer who created this commit
     */
    author: (parent: Commit, _args: unknown, _context: GraphQLContext): Developer => {
      // The developer is already included via the include option in the query
      // so we can just return it from the parent object
      return parent.developer!;
    },
  },
};
