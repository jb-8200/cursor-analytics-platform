import { setupServer } from 'msw/node';
import { handlers } from './handlers';

/**
 * MSW server for Node.js test environment
 * Automatically intercepts HTTP requests during tests
 */
export const server = setupServer(...handlers);
