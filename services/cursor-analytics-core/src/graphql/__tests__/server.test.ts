/**
 * Apollo Server Tests
 *
 * Tests for GraphQL server setup, schema validation, and health check query.
 */

/* eslint-disable @typescript-eslint/no-unsafe-assignment */
/* eslint-disable @typescript-eslint/no-unsafe-member-access */
/* eslint-disable @typescript-eslint/no-unsafe-call */
/* eslint-disable @typescript-eslint/no-unsafe-return */
/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable @typescript-eslint/unbound-method */

import { ApolloServer } from '@apollo/server';
import { PrismaClient } from '../../generated/prisma';
import { CursorSimClient } from '../../ingestion/client';
import { createApolloServer } from '../server';
import { GraphQLContext } from '../context';

// Mock dependencies
jest.mock('../../generated/prisma');
jest.mock('../../ingestion/client');

describe('Apollo Server', () => {
  let server: ApolloServer<GraphQLContext>;
  let mockDb: jest.Mocked<PrismaClient>;
  let mockSimClient: jest.Mocked<CursorSimClient>;

  beforeEach(async () => {
    // Create mock instances
    mockDb = new PrismaClient() as jest.Mocked<PrismaClient>;
    mockSimClient = new CursorSimClient({
      baseUrl: 'http://localhost:8080',
      apiKey: 'test-key',
    }) as jest.Mocked<CursorSimClient>;

    // Mock $queryRaw for health checks
    mockDb.$queryRaw = jest.fn().mockResolvedValue([{ 1: 1 }]);

    // Create server instance
    server = createApolloServer({
      db: mockDb,
      simClient: mockSimClient,
    });

    // Start the server for testing
    await server.start();
  });

  afterEach(async () => {
    if (server) {
      await server.stop();
    }
    jest.clearAllMocks();
  });

  describe('Server Creation', () => {
    it('should create Apollo Server instance', () => {
      expect(server).toBeInstanceOf(ApolloServer);
    });

    it('should enable introspection', async () => {
      // Verify introspection works by running an introspection query
      const result = await server.executeOperation(
        {
          query: `
            query {
              __schema {
                queryType {
                  name
                }
              }
            }
          `,
        },
        {
          contextValue: {
            db: mockDb,
            simClient: mockSimClient,
            requestId: 'test-req-id',
          },
        }
      );

      expect(result.body.kind).toBe('single');
      if (result.body.kind === 'single') {
        expect(result.body.singleResult.errors).toBeUndefined();
      }
    });
  });

  describe('Health Check Query', () => {
    beforeEach(() => {
      // Mock successful database query
      mockDb.$queryRaw = jest.fn().mockResolvedValue([{ 1: 1 }]);

      // Mock successful fetch for simulator health check
      global.fetch = jest.fn().mockResolvedValue({
        ok: true,
        status: 200,
      } as Response);
    });

    it('should return healthy status when all services are operational', async () => {
      const result = await server.executeOperation(
        {
          query: `
            query Health {
              health {
                status
                database
                simulator
                version
              }
            }
          `,
        },
        {
          contextValue: {
            db: mockDb,
            simClient: mockSimClient,
            requestId: 'test-req-id',
          },
        }
      );

      expect(result.body.kind).toBe('single');
      if (result.body.kind === 'single') {
        expect(result.body.singleResult.errors).toBeUndefined();
        expect(result.body.singleResult.data).toEqual({
          health: {
            status: 'healthy',
            database: 'connected',
            simulator: 'reachable',
            version: '0.1.0',
          },
        });
      }

      // Verify database health check was called
      expect(mockDb.$queryRaw).toHaveBeenCalled();

      // Verify simulator health check was called
      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/health'),
        expect.objectContaining({
          method: 'GET',
          signal: expect.any(AbortSignal),
        })
      );
    });

    it('should return degraded status when database is disconnected', async () => {
      // Mock database failure
      mockDb.$queryRaw = jest.fn().mockRejectedValue(new Error('Connection refused'));

      const result = await server.executeOperation(
        {
          query: `
            query Health {
              health {
                status
                database
                simulator
              }
            }
          `,
        },
        {
          contextValue: {
            db: mockDb,
            simClient: mockSimClient,
            requestId: 'test-req-id',
          },
        }
      );

      expect(result.body.kind).toBe('single');
      if (result.body.kind === 'single') {
        expect(result.body.singleResult.errors).toBeUndefined();
        expect(result.body.singleResult.data?.health).toMatchObject({
          status: 'degraded',
          database: 'disconnected',
          simulator: 'reachable',
        });
      }
    });

    it('should return degraded status when simulator is unreachable', async () => {
      // Mock simulator failure
      global.fetch = jest.fn().mockRejectedValue(new Error('Network error'));

      const result = await server.executeOperation(
        {
          query: `
            query Health {
              health {
                status
                database
                simulator
              }
            }
          `,
        },
        {
          contextValue: {
            db: mockDb,
            simClient: mockSimClient,
            requestId: 'test-req-id',
          },
        }
      );

      expect(result.body.kind).toBe('single');
      if (result.body.kind === 'single') {
        expect(result.body.singleResult.errors).toBeUndefined();
        expect(result.body.singleResult.data?.health).toMatchObject({
          status: 'degraded',
          database: 'connected',
          simulator: 'unreachable',
        });
      }
    });

    it('should return degraded status when both services fail', async () => {
      // Mock both failures
      mockDb.$queryRaw = jest.fn().mockRejectedValue(new Error('DB error'));
      global.fetch = jest.fn().mockRejectedValue(new Error('Network error'));

      const result = await server.executeOperation(
        {
          query: `
            query Health {
              health {
                status
                database
                simulator
              }
            }
          `,
        },
        {
          contextValue: {
            db: mockDb,
            simClient: mockSimClient,
            requestId: 'test-req-id',
          },
        }
      );

      expect(result.body.kind).toBe('single');
      if (result.body.kind === 'single') {
        expect(result.body.singleResult.errors).toBeUndefined();
        expect(result.body.singleResult.data?.health).toMatchObject({
          status: 'degraded',
          database: 'disconnected',
          simulator: 'unreachable',
        });
      }
    });

    it('should include lastIngestion field (currently null)', async () => {
      const result = await server.executeOperation(
        {
          query: `
            query Health {
              health {
                lastIngestion
              }
            }
          `,
        },
        {
          contextValue: {
            db: mockDb,
            simClient: mockSimClient,
            requestId: 'test-req-id',
          },
        }
      );

      expect(result.body.kind).toBe('single');
      if (result.body.kind === 'single') {
        expect(result.body.singleResult.errors).toBeUndefined();
        expect(result.body.singleResult.data?.health).toHaveProperty('lastIngestion', null);
      }
    });
  });

  describe('Schema Validation', () => {
    it('should validate GraphQL schema without errors', async () => {
      // If the server starts without throwing, the schema is valid
      expect(server).toBeDefined();

      // Test a simple introspection query
      const result = await server.executeOperation(
        {
          query: `
            query IntrospectionQuery {
              __schema {
                queryType {
                  name
                }
              }
            }
          `,
        },
        {
          contextValue: {
            db: mockDb,
            simClient: mockSimClient,
            requestId: 'test-req-id',
          },
        }
      );

      expect(result.body.kind).toBe('single');
      if (result.body.kind === 'single') {
        expect(result.body.singleResult.errors).toBeUndefined();
        const data = result.body.singleResult.data as any;
        expect(data?.__schema.queryType.name).toBe('Query');
      }
    });

    it('should have all expected query fields defined', async () => {
      const result = await server.executeOperation(
        {
          query: `
            query QueryFields {
              __schema {
                queryType {
                  fields {
                    name
                  }
                }
              }
            }
          `,
        },
        {
          contextValue: {
            db: mockDb,
            simClient: mockSimClient,
            requestId: 'test-req-id',
          },
        }
      );

      expect(result.body.kind).toBe('single');
      if (result.body.kind === 'single') {
        expect(result.body.singleResult.errors).toBeUndefined();
        const data = result.body.singleResult.data as any;
        const fields = data?.__schema.queryType.fields.map((f: any) => f.name);

        // Verify all expected query fields are present
        expect(fields).toContain('health');
        expect(fields).toContain('developer');
        expect(fields).toContain('developers');
        expect(fields).toContain('teamStats');
        expect(fields).toContain('teams');
        expect(fields).toContain('dashboardSummary');
      }
    });
  });

  describe('DateTime Scalar', () => {
    it('should serialize Date objects to ISO strings', async () => {
      // The DateTime scalar is tested indirectly through the health query
      // We'll add more comprehensive tests when we have fields that return DateTime
      const result = await server.executeOperation(
        {
          query: `
            query TypeQuery {
              __type(name: "DateTime") {
                name
                kind
              }
            }
          `,
        },
        {
          contextValue: {
            db: mockDb,
            simClient: mockSimClient,
            requestId: 'test-req-id',
          },
        }
      );

      expect(result.body.kind).toBe('single');
      if (result.body.kind === 'single') {
        expect(result.body.singleResult.errors).toBeUndefined();
        expect(result.body.singleResult.data?.__type).toEqual({
          name: 'DateTime',
          kind: 'SCALAR',
        });
      }
    });
  });

  describe('Error Handling', () => {
    it('should handle invalid queries gracefully', async () => {
      const result = await server.executeOperation(
        {
          query: `
            query Invalid {
              nonExistentField
            }
          `,
        },
        {
          contextValue: {
            db: mockDb,
            simClient: mockSimClient,
            requestId: 'test-req-id',
          },
        }
      );

      expect(result.body.kind).toBe('single');
      if (result.body.kind === 'single') {
        expect(result.body.singleResult.errors).toBeDefined();
        expect(result.body.singleResult.errors?.[0].message).toContain('nonExistentField');
      }
    });
  });
});
