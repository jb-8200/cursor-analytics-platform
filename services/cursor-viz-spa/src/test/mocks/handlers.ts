import { graphqlHandlers } from './graphqlHandlers';

/**
 * MSW handlers for mocking API requests
 *
 * GraphQL handlers are defined in graphqlHandlers.ts and provide
 * mock responses for all GraphQL queries used in the application.
 */

export const handlers = [...graphqlHandlers];
