/**
 * cursor-analytics-core Entry Point
 *
 * Starts the GraphQL server with full schema and resolvers.
 */

import { startStandaloneServer } from '@apollo/server/standalone';
import { PrismaClient } from './generated/prisma';
import { CursorSimClient } from './ingestion/client';
import { createApolloServer } from './graphql/server';
import { createContext } from './graphql/context';
import { config } from './config';

async function startServer() {
  // Initialize database client
  const db = new PrismaClient({
    log: config.nodeEnv === 'development' ? ['query', 'error', 'warn'] : ['error'],
  });

  // Initialize cursor-sim REST client
  const simClient = new CursorSimClient({
    baseUrl: config.simulatorUrl,
    apiKey: process.env.SIMULATOR_API_KEY || 'cursor-sim-dev-key',
    timeout: 10000,
    retryAttempts: 3,
    retryDelayMs: 1000,
  });

  // Create Apollo Server
  const server = createApolloServer({
    db,
    simClient,
  });

  // Start server
  const { url } = await startStandaloneServer(server, {
    listen: { port: config.port },
    context: () => createContext({ db, simClient }),
  });

  console.log(`🚀 cursor-analytics-core ready at ${url}`);
  console.log(`📊 GraphQL Playground: ${url}`);
  console.log(`✅ Step 05 (GraphQL Schema) - COMPLETE`);
  console.log(`
Available queries:
  - health          Health check for all services
  - developer       (Step 06 - not yet implemented)
  - developers      (Step 06 - not yet implemented)
  - teamStats       (Step 09 - not yet implemented)
  - teams           (Step 09 - not yet implemented)
  - dashboardSummary (Step 09 - not yet implemented)
  `);

  // Handle graceful shutdown
  process.on('SIGTERM', () => {
    console.log('SIGTERM received, shutting down gracefully...');
    void (async () => {
      await server.stop();
      await db.$disconnect();
      process.exit(0);
    })();
  });

  process.on('SIGINT', () => {
    console.log('SIGINT received, shutting down gracefully...');
    void (async () => {
      await server.stop();
      await db.$disconnect();
      process.exit(0);
    })();
  });
}

startServer().catch((err) => {
  console.error('Failed to start server:', err);
  process.exit(1);
});
