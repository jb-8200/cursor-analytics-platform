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
import { resolvers } from './resolvers';

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
