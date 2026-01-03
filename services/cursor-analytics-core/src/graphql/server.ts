/**
 * Apollo Server Configuration
 *
 * Sets up Apollo Server with:
 * - GraphQL schema (type definitions)
 * - Resolvers
 * - Context creation
 * - Error handling
 */

import { ApolloServer } from '@apollo/server';
import { typeDefs } from './schema';
import { GraphQLContext } from './context';
import { PrismaClient } from '../generated/prisma';
import { CursorSimClient } from '../ingestion/client';
import { GraphQLError } from 'graphql';

/**
 * Base resolvers - currently only health check
 * Other resolvers (developer, dashboard, etc.) will be added in later steps
 */
const resolvers = {
  Query: {
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

/**
 * Create and configure Apollo Server instance
 */
export function createApolloServer(_params: {
  db: PrismaClient;
  simClient: CursorSimClient;
}): ApolloServer<GraphQLContext> {
  return new ApolloServer<GraphQLContext>({
    typeDefs,
    resolvers,
    introspection: true, // Enable GraphQL Playground in all environments
    formatError: (formattedError) => {
      // Log errors for debugging
      console.error('GraphQL Error:', {
        message: formattedError.message,
        path: formattedError.path,
        extensions: formattedError.extensions,
      });

      // Return formatted error to client
      return formattedError;
    },
  });
}
