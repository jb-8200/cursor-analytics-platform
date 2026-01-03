/**
 * GraphQL Context Tests
 */

/* eslint-disable @typescript-eslint/no-unsafe-assignment */
/* eslint-disable @typescript-eslint/await-thenable */

import { createContext } from '../context';
import { PrismaClient } from '../../generated/prisma';
import { CursorSimClient } from '../../ingestion/client';

// Mock dependencies
jest.mock('../../generated/prisma');
jest.mock('../../ingestion/client');

describe('GraphQL Context', () => {
  let mockDb: jest.Mocked<PrismaClient>;
  let mockSimClient: jest.Mocked<CursorSimClient>;

  beforeEach(() => {
    mockDb = new PrismaClient() as jest.Mocked<PrismaClient>;
    mockSimClient = new CursorSimClient({
      baseUrl: 'http://localhost:8080',
      apiKey: 'test-key',
    }) as jest.Mocked<CursorSimClient>;
  });

  describe('createContext', () => {
    it('should create context with db and simClient', async () => {
      const context = await createContext({
        db: mockDb,
        simClient: mockSimClient,
      });

      expect(context).toEqual({
        db: mockDb,
        simClient: mockSimClient,
        requestId: expect.stringMatching(/^req_\d+_[a-z0-9]+$/),
      });
    });

    it('should generate unique request IDs', async () => {
      const context1 = await createContext({
        db: mockDb,
        simClient: mockSimClient,
      });

      const context2 = await createContext({
        db: mockDb,
        simClient: mockSimClient,
      });

      expect(context1.requestId).toBeDefined();
      expect(context2.requestId).toBeDefined();
      expect(context1.requestId).not.toBe(context2.requestId);
    });

    it('should include all required properties', async () => {
      const context = await createContext({
        db: mockDb,
        simClient: mockSimClient,
      });

      expect(context).toHaveProperty('db');
      expect(context).toHaveProperty('simClient');
      expect(context).toHaveProperty('requestId');
    });
  });
});
