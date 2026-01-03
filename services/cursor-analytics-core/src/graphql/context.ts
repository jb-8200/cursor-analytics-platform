/**
 * GraphQL Context
 *
 * Defines the context object available to all resolvers.
 * Provides access to:
 * - Database client (Prisma)
 * - REST client for cursor-sim
 * - Request metadata
 */

import { PrismaClient } from '../generated/prisma';
import { CursorSimClient } from '../ingestion/client';

export interface GraphQLContext {
  /**
   * Prisma database client
   */
  db: PrismaClient;

  /**
   * cursor-sim REST API client
   */
  simClient: CursorSimClient;

  /**
   * Request metadata
   */
  requestId?: string;
}

/**
 * Create context for a GraphQL request
 */
export function createContext(params: {
  db: PrismaClient;
  simClient: CursorSimClient;
}): GraphQLContext {
  return {
    db: params.db,
    simClient: params.simClient,
    requestId: generateRequestId(),
  };
}

/**
 * Generate a unique request ID for tracing
 */
function generateRequestId(): string {
  return `req_${Date.now()}_${Math.random().toString(36).substring(7)}`;
}
