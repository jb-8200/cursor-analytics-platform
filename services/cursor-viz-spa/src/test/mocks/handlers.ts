import { http, HttpResponse } from 'msw';

/**
 * MSW handlers for mocking GraphQL API
 * Will be expanded as GraphQL schema is implemented in P5
 */

const GRAPHQL_URL =
  import.meta.env.VITE_GRAPHQL_URL || 'http://localhost:4000/graphql';

export const handlers = [
  // Mock GraphQL endpoint
  http.post(GRAPHQL_URL, () => {
    return HttpResponse.json({
      data: {
        // Add mock responses as queries are defined
      },
    });
  }),
];
