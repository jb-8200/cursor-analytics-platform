/**
 * Integration Tests: GraphQL Queries with Real Database
 *
 * Tests full GraphQL query execution with Prisma and PostgreSQL.
 */

import { ApolloServer } from '@apollo/server';
import { createApolloServer } from '../../graphql/server';
import { createContext, type GraphQLContext } from '../../graphql/context';
import { createTestDb, createTestSimClient, cleanupDb, seedTestData } from './setup';
import type { PrismaClient } from '../../generated/prisma';
import type { CursorSimClient } from '../../ingestion/client';

describe('GraphQL Integration Tests', () => {
  let db: PrismaClient;
  let simClient: CursorSimClient;
  let server: ApolloServer<GraphQLContext>;
  let testDevelopers: any;

  beforeAll(async () => {
    db = createTestDb();
    simClient = createTestSimClient();
    server = createApolloServer({ db, simClient });

    // Seed test data
    testDevelopers = await seedTestData(db);
  });

  afterAll(async () => {
    await cleanupDb(db);
  });

  describe('health query', () => {
    it('should return health status for all services', async () => {
      const result = await server.executeOperation(
        {
          query: `
            query HealthCheck {
              health {
                status
                database
                version
              }
            }
          `,
        },
        {
          contextValue: createContext({ db, simClient }),
        },
      );

      expect(result.body.kind).toBe('single');
      if (result.body.kind === 'single') {
        expect(result.body.singleResult.errors).toBeUndefined();
        expect(result.body.singleResult.data?.health).toMatchObject({
          status: 'healthy',
          database: 'connected',
          version: expect.any(String),
        });
      }
    });
  });

  describe('developer query', () => {
    it('should fetch single developer by ID', async () => {
      const result = await server.executeOperation(
        {
          query: `
            query GetDeveloper($id: ID!) {
              developer(id: $id) {
                id
                externalId
                name
                email
                team
                seniority
              }
            }
          `,
          variables: { id: testDevelopers.dev1.id },
        },
        {
          contextValue: createContext({ db, simClient }),
        },
      );

      expect(result.body.kind).toBe('single');
      if (result.body.kind === 'single') {
        expect(result.body.singleResult.errors).toBeUndefined();
        expect(result.body.singleResult.data?.developer).toMatchObject({
          id: testDevelopers.dev1.id,
          externalId: 'test-dev-1',
          name: 'Alice Developer',
          email: 'alice@example.com',
          team: 'Team A',
          seniority: 'senior',
        });
      }
    });

    it('should return null for non-existent developer', async () => {
      const result = await server.executeOperation(
        {
          query: `
            query GetDeveloper($id: ID!) {
              developer(id: $id) {
                id
                name
              }
            }
          `,
          variables: { id: '00000000-0000-0000-0000-000000000000' },
        },
        {
          contextValue: createContext({ db, simClient }),
        },
      );

      expect(result.body.kind).toBe('single');
      if (result.body.kind === 'single') {
        expect(result.body.singleResult.errors).toBeUndefined();
        expect(result.body.singleResult.data?.developer).toBeNull();
      }
    });
  });

  describe('developers query', () => {
    it('should list all developers with pagination', async () => {
      const result = await server.executeOperation(
        {
          query: `
            query ListDevelopers {
              developers(limit: 10) {
                nodes {
                  id
                  name
                  team
                }
                totalCount
                pageInfo {
                  hasNextPage
                  hasPreviousPage
                }
              }
            }
          `,
        },
        {
          contextValue: createContext({ db, simClient }),
        },
      );

      expect(result.body.kind).toBe('single');
      if (result.body.kind === 'single') {
        expect(result.body.singleResult.errors).toBeUndefined();
        const data = result.body.singleResult.data?.developers as any;
        expect(data.nodes).toHaveLength(3);
        expect(data.totalCount).toBe(3);
        expect(data.pageInfo.hasNextPage).toBe(false);
        expect(data.pageInfo.hasPreviousPage).toBe(false);
      }
    });

    it('should filter developers by team', async () => {
      const result = await server.executeOperation(
        {
          query: `
            query FilterByTeam($team: String!) {
              developers(team: $team) {
                nodes {
                  name
                  team
                }
                totalCount
              }
            }
          `,
          variables: { team: 'Team A' },
        },
        {
          contextValue: createContext({ db, simClient }),
        },
      );

      expect(result.body.kind).toBe('single');
      if (result.body.kind === 'single') {
        expect(result.body.singleResult.errors).toBeUndefined();
        const data = result.body.singleResult.data?.developers as any;
        expect(data.nodes).toHaveLength(2);
        expect(data.totalCount).toBe(2);
        expect(data.nodes.every((d: any) => d.team === 'Team A')).toBe(true);
      }
    });

    it('should filter developers by seniority', async () => {
      const result = await server.executeOperation(
        {
          query: `
            query FilterBySeniority($seniority: String!) {
              developers(seniority: $seniority) {
                nodes {
                  name
                  seniority
                }
                totalCount
              }
            }
          `,
          variables: { seniority: 'senior' },
        },
        {
          contextValue: createContext({ db, simClient }),
        },
      );

      expect(result.body.kind).toBe('single');
      if (result.body.kind === 'single') {
        expect(result.body.singleResult.errors).toBeUndefined();
        const data = result.body.singleResult.data?.developers as any;
        expect(data.nodes).toHaveLength(1);
        expect(data.nodes[0].name).toBe('Alice Developer');
      }
    });
  });

  describe('developer stats', () => {
    it('should calculate usage stats for developer', async () => {
      const today = new Date().toISOString().split('T')[0];
      const tomorrow = new Date();
      tomorrow.setDate(tomorrow.getDate() + 1);

      const result = await server.executeOperation(
        {
          query: `
            query DeveloperStats($id: ID!, $range: DateRangeInput) {
              developer(id: $id) {
                name
                stats(range: $range) {
                  totalSuggestions
                  acceptedSuggestions
                  acceptanceRate
                  chatInteractions
                  cmdKUsages
                  totalLinesAdded
                  aiLinesAdded
                  aiVelocity
                }
              }
            }
          `,
          variables: {
            id: testDevelopers.dev1.id,
            range: {
              from: today,
              to: tomorrow.toISOString().split('T')[0],
            },
          },
        },
        {
          contextValue: createContext({ db, simClient }),
        },
      );

      expect(result.body.kind).toBe('single');
      if (result.body.kind === 'single') {
        expect(result.body.singleResult.errors).toBeUndefined();
        const stats = (result.body.singleResult.data as any)?.developer.stats;
        expect(stats).toMatchObject({
          totalSuggestions: 10,
          acceptedSuggestions: 8,
          acceptanceRate: expect.any(Number),
          chatInteractions: 1,
          cmdKUsages: 0,
          totalLinesAdded: expect.any(Number),
          aiLinesAdded: expect.any(Number),
        });
        expect(stats.acceptanceRate).toBeGreaterThan(70); // 8/10 = 80%
        expect(stats.aiVelocity).toBeDefined();
      }
    });

    it('should return daily stats breakdown', async () => {
      const today = new Date().toISOString().split('T')[0];
      const tomorrow = new Date();
      tomorrow.setDate(tomorrow.getDate() + 1);

      const result = await server.executeOperation(
        {
          query: `
            query DeveloperDailyStats($id: ID!, $range: DateRangeInput) {
              developer(id: $id) {
                name
                dailyStats(range: $range) {
                  date
                  suggestionsShown
                  suggestionsAccepted
                  acceptanceRate
                  linesAdded
                  aiLinesAdded
                }
              }
            }
          `,
          variables: {
            id: testDevelopers.dev1.id,
            range: {
              from: today,
              to: tomorrow.toISOString().split('T')[0],
            },
          },
        },
        {
          contextValue: createContext({ db, simClient }),
        },
      );

      expect(result.body.kind).toBe('single');
      if (result.body.kind === 'single') {
        expect(result.body.singleResult.errors).toBeUndefined();
        const dailyStats = (result.body.singleResult.data as any)?.developer.dailyStats;
        expect(dailyStats).toHaveLength(1); // One day of data
        expect(dailyStats[0]).toMatchObject({
          suggestionsShown: 10,
          suggestionsAccepted: 8,
          linesAdded: expect.any(Number),
          aiLinesAdded: expect.any(Number),
        });
      }
    });
  });

  describe('commits query', () => {
    it('should list commits with pagination', async () => {
      const result = await server.executeOperation(
        {
          query: `
            query ListCommits {
              commits(limit: 10) {
                nodes {
                  id
                  externalId
                  linesAdded
                  linesDeleted
                  modelUsed
                  author {
                    name
                  }
                }
                totalCount
                pageInfo {
                  hasNextPage
                }
              }
            }
          `,
        },
        {
          contextValue: createContext({ db, simClient }),
        },
      );

      expect(result.body.kind).toBe('single');
      if (result.body.kind === 'single') {
        expect(result.body.singleResult.errors).toBeUndefined();
        const data = result.body.singleResult.data?.commits as any;
        expect(data.nodes.length).toBeGreaterThan(0);
        expect(data.totalCount).toBeGreaterThan(0);
        expect(data.nodes[0]).toHaveProperty('author');
      }
    });

    it('should filter commits by userId', async () => {
      const result = await server.executeOperation(
        {
          query: `
            query FilterCommits($userId: ID!) {
              commits(userId: $userId, limit: 20) {
                nodes {
                  externalId
                  author {
                    id
                    name
                  }
                }
                totalCount
              }
            }
          `,
          variables: { userId: testDevelopers.dev1.id },
        },
        {
          contextValue: createContext({ db, simClient }),
        },
      );

      expect(result.body.kind).toBe('single');
      if (result.body.kind === 'single') {
        expect(result.body.singleResult.errors).toBeUndefined();
        const data = result.body.singleResult.data?.commits as any;
        expect(data.nodes.every((c: any) => c.author.id === testDevelopers.dev1.id)).toBe(true);
      }
    });
  });

  describe('dashboard summary', () => {
    it('should return comprehensive dashboard KPIs', async () => {
      const result = await server.executeOperation(
        {
          query: `
            query DashboardSummary($preset: DateRangePreset) {
              dashboardSummary(preset: $preset) {
                totalDevelopers
                activeDevelopers
                overallAcceptanceRate
                totalSuggestionsToday
                totalAcceptedToday
                teamComparison {
                  teamName
                  memberCount
                  activeMemberCount
                  averageAcceptanceRate
                  topPerformer {
                    name
                  }
                }
              }
            }
          `,
          variables: { preset: 'LAST_7_DAYS' },
        },
        {
          contextValue: createContext({ db, simClient }),
        },
      );

      expect(result.body.kind).toBe('single');
      if (result.body.kind === 'single') {
        expect(result.body.singleResult.errors).toBeUndefined();
        const summary = result.body.singleResult.data?.dashboardSummary as any;
        expect(summary).toMatchObject({
          totalDevelopers: 3,
          activeDevelopers: expect.any(Number),
          overallAcceptanceRate: expect.any(Number),
          totalSuggestionsToday: expect.any(Number),
          totalAcceptedToday: expect.any(Number),
          teamComparison: expect.any(Array),
        });
        expect(summary.teamComparison.length).toBeGreaterThan(0);
      }
    });
  });

  describe('team stats', () => {
    it('should return statistics for specific team', async () => {
      const result = await server.executeOperation(
        {
          query: `
            query TeamStats($teamName: String!) {
              teamStats(teamName: $teamName) {
                teamName
                memberCount
                activeMemberCount
                averageAcceptanceRate
                totalSuggestions
                totalAccepted
                topPerformer {
                  name
                  email
                }
              }
            }
          `,
          variables: { teamName: 'Team A' },
        },
        {
          contextValue: createContext({ db, simClient }),
        },
      );

      expect(result.body.kind).toBe('single');
      if (result.body.kind === 'single') {
        expect(result.body.singleResult.errors).toBeUndefined();
        const stats = result.body.singleResult.data?.teamStats as any;
        expect(stats).toMatchObject({
          teamName: 'Team A',
          memberCount: 2,
          activeMemberCount: expect.any(Number),
          topPerformer: expect.objectContaining({
            name: expect.any(String),
          }),
        });
      }
    });

    it('should list all teams', async () => {
      const result = await server.executeOperation(
        {
          query: `
            query AllTeams {
              teams {
                teamName
                memberCount
                averageAcceptanceRate
              }
            }
          `,
        },
        {
          contextValue: createContext({ db, simClient }),
        },
      );

      expect(result.body.kind).toBe('single');
      if (result.body.kind === 'single') {
        expect(result.body.singleResult.errors).toBeUndefined();
        const teams = result.body.singleResult.data?.teams as any;
        expect(teams.length).toBe(2); // Team A and Team B
        expect(teams.map((t: any) => t.teamName).sort()).toEqual(['Team A', 'Team B']);
      }
    });
  });
});
