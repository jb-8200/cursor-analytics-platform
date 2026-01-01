import { ApolloServer } from '@apollo/server';
import { startStandaloneServer } from '@apollo/server/standalone';

// P0 scaffolding - minimal GraphQL schema
const typeDefs = `#graphql
  type Query {
    health: Health!
    placeholder: String!
  }

  type Health {
    status: String!
    timestamp: String!
    service: String!
    version: String!
  }
`;

// P0 scaffolding - basic resolvers
const resolvers = {
  Query: {
    health: () => ({
      status: 'healthy',
      timestamp: new Date().toISOString(),
      service: 'cursor-analytics-core',
      version: '0.0.1-p0',
    }),
    placeholder: () => 'GraphQL server is running (P0 scaffolding - not yet implemented)',
  },
};

async function startServer() {
  const server = new ApolloServer({
    typeDefs,
    resolvers,
  });

  const { url } = await startStandaloneServer(server, {
    listen: { port: Number(process.env.PORT) || 4000 },
  });

  console.log(`🚀 cursor-analytics-core ready at ${url}`);
  console.log(`📊 GraphQL Playground: ${url}`);
  console.log(`Status: P0 scaffolding - schema not yet implemented`);
}

startServer().catch((err) => {
  console.error('Failed to start server:', err);
  process.exit(1);
});
