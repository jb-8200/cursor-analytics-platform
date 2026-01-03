import { ApolloClient, InMemoryCache, HttpLink, from } from '@apollo/client';
import { onError } from '@apollo/client/link/error';

/**
 * Create and configure Apollo Client for GraphQL communication
 * with cursor-analytics-core.
 *
 * Configuration includes:
 * - HTTP link to GraphQL endpoint
 * - In-memory cache with type policies
 * - Error handling link
 * - Default query options
 */
export function createApolloClient(): ApolloClient<unknown> {
  // Get GraphQL URL from environment variables
  const graphqlUrl = import.meta.env.VITE_GRAPHQL_URL || 'http://localhost:4000/graphql';

  // HTTP link for GraphQL requests
  const httpLink = new HttpLink({
    uri: graphqlUrl,
  });

  // Error handling link
  const errorLink = onError(({ graphQLErrors, networkError }) => {
    if (graphQLErrors) {
      graphQLErrors.forEach(({ message, locations, path }) => {
        console.error(
          `[GraphQL error]: Message: ${message}, Location: ${JSON.stringify(locations)}, Path: ${path}`
        );
      });
    }

    if (networkError) {
      console.error(`[Network error]: ${networkError}`);
    }
  });

  // Configure in-memory cache with type policies
  const cache = new InMemoryCache({
    typePolicies: {
      // Developer type: use 'id' as the key field
      Developer: {
        keyFields: ['id'],
      },
      // DailyStats type: use 'date' as the key field
      DailyStats: {
        keyFields: ['date'],
      },
      // Query type: configure field policies for pagination
      Query: {
        fields: {
          developers: {
            // Merge strategy for paginated developers
            keyArgs: ['team'], // Separate cache entries by team filter
            merge(existing, incoming, { args }) {
              if (!existing || args?.offset === 0) {
                return incoming;
              }
              return {
                ...incoming,
                nodes: [...(existing.nodes || []), ...(incoming.nodes || [])],
              };
            },
          },
        },
      },
    },
  });

  // Combine links
  const link = from([errorLink, httpLink]);

  // Create Apollo Client
  const client = new ApolloClient({
    link,
    cache,
    defaultOptions: {
      query: {
        errorPolicy: 'all',
        fetchPolicy: 'cache-first',
      },
      watchQuery: {
        errorPolicy: 'all',
        fetchPolicy: 'cache-and-network',
      },
      mutate: {
        errorPolicy: 'all',
      },
    },
  });

  return client;
}

// Export singleton instance for use in the app
export const apolloClient = createApolloClient();

// Default export for convenience
export default apolloClient;
