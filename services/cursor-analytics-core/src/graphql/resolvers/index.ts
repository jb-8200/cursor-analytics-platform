/**
 * Resolver Index
 *
 * Exports all GraphQL resolvers for the application.
 * Resolvers are modularized by domain (developer, commit, dashboard, etc.)
 */

import { developerResolvers } from './developer';
import { commitResolvers } from './commit';
import { dashboardResolvers } from './dashboard';
import { GraphQLContext } from '../context';
import { GraphQLError } from 'graphql';

/**
 * Merge all resolvers into a single object
 * This includes Query, Developer, Commit field resolvers, DateTime scalar, etc.
 */
export const resolvers = {
  Query: {
    // Health check resolver
    health: async (_parent: unknown, _args: unknown, context: GraphQLContext) => {
      // Check database connectivity
      let dbStatus = 'disconnected';
      try {
        await context.db.$queryRaw`SELECT 1`;
        dbStatus = 'connected';
      } catch (error) {
        console.error('Database health check failed:', error);
        dbStatus = 'disconnected';
      }

      // Check cursor-sim connectivity
      let simStatus = 'unreachable';
      try {
        const response = await fetch(`${context.simClient['baseUrl']}/health`, {
          method: 'GET',
          signal: AbortSignal.timeout(5000), // 5 second timeout
        });
        simStatus = response.ok ? 'reachable' : 'unreachable';
      } catch (error) {
        console.error('cursor-sim health check failed:', error);
        simStatus = 'unreachable';
      }

      return {
        status: dbStatus === 'connected' && simStatus === 'reachable' ? 'healthy' : 'degraded',
        database: dbStatus,
        simulator: simStatus,
        lastIngestion: null, // Will be implemented in Step 04 (Ingestion Worker)
        version: '0.1.0',
      };
    },

    // Developer resolvers
    ...developerResolvers.Query,

    // Commit resolvers
    ...commitResolvers.Query,

    // Dashboard resolvers
    ...dashboardResolvers.Query,
  },

  // Developer field resolvers
  Developer: {
    ...developerResolvers.Developer,
  },

  // Commit field resolvers
  Commit: {
    ...commitResolvers.Commit,
  },

  // DateTime scalar resolver
  DateTime: {
    serialize: (value: Date | string) => {
      if (value instanceof Date) {
        return value.toISOString();
      }
      return value;
    },
    parseValue: (value: string) => {
      return new Date(value);
    },
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    parseLiteral: (ast: any) => {
      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      if (ast.kind === 'StringValue') {
        // eslint-disable-next-line @typescript-eslint/no-unsafe-argument, @typescript-eslint/no-unsafe-member-access
        return new Date(ast.value);
      }
      throw new GraphQLError('DateTime must be a string in ISO 8601 format');
    },
  },
};
