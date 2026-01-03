/**
 * E2E Test: Full Pipeline from Database to GraphQL Queries
 *
 * Tests the complete data flow:
 * 1. Seed database with realistic data
 * 2. Execute complex GraphQL queries
 * 3. Verify aggregations and calculations
 */

import { ApolloServer } from '@apollo/server';
import { createApolloServer } from '../../graphql/server';
import { createContext, type GraphQLContext } from '../../graphql/context';
import { createTestDb, createTestSimClient, cleanupDb } from '../integration/setup';
import type { PrismaClient } from '../../generated/prisma';
import type { CursorSimClient } from '../../ingestion/client';

describe('E2E: Full Pipeline Test', () => {
  let db: PrismaClient;
  let simClient: CursorSimClient;
  let server: ApolloServer<GraphQLContext>;

  beforeAll(async () => {
    db = createTestDb();
    simClient = createTestSimClient();
    server = createApolloServer({ db, simClient });

    // Seed realistic multi-day dataset
    await seedRealisticData(db);
  });

  afterAll(async () => {
    await cleanupDb(db);
  });

  /**
   * Seed database with realistic multi-day, multi-developer data
   */
  async function seedRealisticData(db: PrismaClient) {
    // Create 5 developers across 2 teams
    const developers = await Promise.all([
      db.developer.create({
        data: {
          externalId: 'e2e-dev-1',
          name: 'Senior Alice',
          email: 'alice@company.com',
          team: 'Backend',
          seniority: 'senior',
        },
      }),
      db.developer.create({
        data: {
          externalId: 'e2e-dev-2',
          name: 'Mid Bob',
          email: 'bob@company.com',
          team: 'Backend',
          seniority: 'mid',
        },
      }),
      db.developer.create({
        data: {
          externalId: 'e2e-dev-3',
          name: 'Junior Charlie',
          email: 'charlie@company.com',
          team: 'Backend',
          seniority: 'junior',
        },
      }),
      db.developer.create({
        data: {
          externalId: 'e2e-dev-4',
          name: 'Senior Diana',
          email: 'diana@company.com',
          team: 'Frontend',
          seniority: 'senior',
        },
      }),
      db.developer.create({
        data: {
          externalId: 'e2e-dev-5',
          name: 'Mid Eve',
          email: 'eve@company.com',
          team: 'Frontend',
          seniority: 'mid',
        },
      }),
    ]);

    // Generate events for last 7 days
    const today = new Date();
    for (let dayOffset = 0; dayOffset < 7; dayOffset++) {
      const date = new Date(today);
      date.setDate(date.getDate() - dayOffset);

      for (const [idx, dev] of developers.entries()) {
        // Each developer has different AI usage patterns
        const aiPreference = [0.9, 0.7, 0.5, 0.8, 0.6][idx]; // Senior devs use AI more
        const suggestionsPerDay = Math.floor(20 * aiPreference);
        const acceptanceRate = aiPreference;

        // Generate suggestions shown
        for (let i = 0; i < suggestionsPerDay; i++) {
          await db.usageEvent.create({
            data: {
              externalId: `${dev.externalId}-day${dayOffset}-shown-${i}`,
              developerId: dev.id,
              eventType: 'cpp_suggestion_shown',
              eventTimestamp: date,
              linesAdded: 0,
              linesDeleted: 0,
              tokensInput: 100,
              tokensOutput: 150,
            },
          });
        }

        // Generate accepted suggestions
        const acceptedCount = Math.floor(suggestionsPerDay * acceptanceRate);
        for (let i = 0; i < acceptedCount; i++) {
          await db.usageEvent.create({
            data: {
              externalId: `${dev.externalId}-day${dayOffset}-accepted-${i}`,
              developerId: dev.id,
              eventType: 'cpp_suggestion_accepted',
              eventTimestamp: date,
              linesAdded: 5 + Math.floor(Math.random() * 15),
              linesDeleted: Math.floor(Math.random() * 3),
              accepted: true,
              tokensInput: 100,
              tokensOutput: 150,
              modelUsed: idx % 2 === 0 ? 'claude-sonnet-4-5' : 'claude-haiku-4',
            },
          });
        }

        // Add some chat interactions (2-5 per day)
        const chatCount = 2 + Math.floor(Math.random() * 4);
        for (let i = 0; i < chatCount; i++) {
          await db.usageEvent.create({
            data: {
              externalId: `${dev.externalId}-day${dayOffset}-chat-${i}`,
              developerId: dev.id,
              eventType: 'chat_message',
              eventTimestamp: date,
              linesAdded: 0,
              linesDeleted: 0,
              tokensInput: 200 + Math.floor(Math.random() * 300),
              tokensOutput: 300 + Math.floor(Math.random() * 500),
            },
          });
        }

        // Add cmd_k usages (1-3 per day)
        const cmdKCount = 1 + Math.floor(Math.random() * 3);
        for (let i = 0; i < cmdKCount; i++) {
          await db.usageEvent.create({
            data: {
              externalId: `${dev.externalId}-day${dayOffset}-cmdk-${i}`,
              developerId: dev.id,
              eventType: 'cmd_k_prompt',
              eventTimestamp: date,
              linesAdded: 0,
              linesDeleted: 0,
              tokensInput: 150 + Math.floor(Math.random() * 200),
              tokensOutput: 250 + Math.floor(Math.random() * 400),
            },
          });
        }
      }
    }
  }

  describe('Complex Multi-Query Pipeline', () => {
    it('should execute dashboard summary with team breakdown', async () => {
      const result = await server.executeOperation(
        {
          query: `
            query ComplexDashboard {
              dashboardSummary(preset: LAST_7_DAYS) {
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
                  totalSuggestions
                  totalAccepted
                  topPerformer {
                    name
                    email
                    seniority
                  }
                }
                dailyTrend {
                  date
                  suggestionsShown
                  suggestionsAccepted
                  acceptanceRate
                  linesAdded
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

        const summary = result.body.singleResult.data?.dashboardSummary as any;

        // Verify high-level KPIs
        expect(summary.totalDevelopers).toBe(5);
        expect(summary.activeDevelopers).toBeGreaterThan(0);
        expect(summary.overallAcceptanceRate).toBeGreaterThan(0);

        // Verify team comparison
        expect(summary.teamComparison).toHaveLength(2); // Backend and Frontend
        const backendTeam = summary.teamComparison.find((t: any) => t.teamName === 'Backend');
        expect(backendTeam).toBeDefined();
        expect(backendTeam.memberCount).toBe(3);
        expect(backendTeam.topPerformer).toBeDefined();
        expect(backendTeam.topPerformer.name).toBeTruthy();

        const frontendTeam = summary.teamComparison.find((t: any) => t.teamName === 'Frontend');
        expect(frontendTeam).toBeDefined();
        expect(frontendTeam.memberCount).toBe(2);

        // Verify daily trend has 7 days
        expect(summary.dailyTrend).toHaveLength(7);
        expect(summary.dailyTrend[0]).toMatchObject({
          date: expect.any(String),
          suggestionsShown: expect.any(Number),
          suggestionsAccepted: expect.any(Number),
          linesAdded: expect.any(Number),
        });
      }
    });

    it('should fetch developers with nested stats and commits', async () => {
      const today = new Date();
      const sevenDaysAgo = new Date(today);
      sevenDaysAgo.setDate(sevenDaysAgo.getDate() - 7);

      const result = await server.executeOperation(
        {
          query: `
            query DevelopersWithStats($range: DateRangeInput!) {
              developers(team: "Backend", limit: 10) {
                nodes {
                  name
                  team
                  seniority
                  stats(range: $range) {
                    totalSuggestions
                    acceptedSuggestions
                    acceptanceRate
                    chatInteractions
                    cmdKUsages
                    aiVelocity
                  }
                  dailyStats(range: $range) {
                    date
                    suggestionsShown
                    suggestionsAccepted
                    acceptanceRate
                  }
                }
                totalCount
              }
            }
          `,
          variables: {
            range: {
              from: sevenDaysAgo.toISOString().split('T')[0],
              to: today.toISOString().split('T')[0],
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

        const developers = result.body.singleResult.data?.developers as any;

        expect(developers.totalCount).toBe(3); // Backend team
        expect(developers.nodes).toHaveLength(3);

        // Verify each developer has computed stats
        for (const dev of developers.nodes) {
          expect(dev.stats).toBeDefined();
          expect(dev.stats.totalSuggestions).toBeGreaterThan(0);
          expect(dev.stats.acceptedSuggestions).toBeGreaterThan(0);
          expect(dev.stats.acceptanceRate).toBeGreaterThan(0);
          expect(dev.stats.chatInteractions).toBeGreaterThan(0);

          // Verify daily stats breakdown
          expect(dev.dailyStats).toHaveLength(7);
          expect(dev.dailyStats[0]).toMatchObject({
            date: expect.any(String),
            suggestionsShown: expect.any(Number),
            suggestionsAccepted: expect.any(Number),
          });
        }
      }
    });

    it('should aggregate commits across teams', async () => {
      const result = await server.executeOperation(
        {
          query: `
            query AllCommits {
              commits(limit: 100, sortBy: "timestamp", sortOrder: "desc") {
                nodes {
                  id
                  linesAdded
                  linesDeleted
                  modelUsed
                  timestamp
                  author {
                    name
                    team
                    seniority
                  }
                }
                totalCount
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

        const commits = result.body.singleResult.data?.commits as any;

        expect(commits.totalCount).toBeGreaterThan(0);
        expect(commits.nodes.length).toBeGreaterThan(0);

        // Verify commits are sorted by timestamp descending
        const timestamps = commits.nodes.map((c: any) => new Date(c.timestamp).getTime());
        for (let i = 1; i < timestamps.length; i++) {
          expect(timestamps[i - 1]).toBeGreaterThanOrEqual(timestamps[i]);
        }

        // Verify each commit has author details
        for (const commit of commits.nodes) {
          expect(commit.author).toBeDefined();
          expect(commit.author.name).toBeTruthy();
          expect(commit.author.team).toMatch(/Backend|Frontend/);
          expect(commit.linesAdded).toBeGreaterThan(0);
        }
      }
    });

    it('should compare team performance metrics', async () => {
      const result = await server.executeOperation(
        {
          query: `
            query CompareTeams {
              backend: teamStats(teamName: "Backend") {
                teamName
                memberCount
                averageAcceptanceRate
                totalSuggestions
                totalAccepted
                aiVelocity
                topPerformer {
                  name
                  seniority
                }
              }
              frontend: teamStats(teamName: "Frontend") {
                teamName
                memberCount
                averageAcceptanceRate
                totalSuggestions
                totalAccepted
                aiVelocity
                topPerformer {
                  name
                  seniority
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

        const { backend, frontend } = result.body.singleResult.data as any;

        // Verify Backend team
        expect(backend).toMatchObject({
          teamName: 'Backend',
          memberCount: 3,
          averageAcceptanceRate: expect.any(Number),
          totalSuggestions: expect.any(Number),
          totalAccepted: expect.any(Number),
        });
        expect(backend.topPerformer).toBeDefined();

        // Verify Frontend team
        expect(frontend).toMatchObject({
          teamName: 'Frontend',
          memberCount: 2,
          averageAcceptanceRate: expect.any(Number),
          totalSuggestions: expect.any(Number),
          totalAccepted: expect.any(Number),
        });
        expect(frontend.topPerformer).toBeDefined();

        // Backend should have more total suggestions (3 vs 2 members)
        expect(backend.totalSuggestions).toBeGreaterThan(frontend.totalSuggestions);
      }
    });

    it('should filter and paginate commits by user and date range', async () => {
      const today = new Date();
      const threeDaysAgo = new Date(today);
      threeDaysAgo.setDate(threeDaysAgo.getDate() - 3);

      // First get a developer ID
      const devResult = await server.executeOperation(
        {
          query: `
            query GetDev {
              developers(limit: 1) {
                nodes {
                  id
                  name
                }
              }
            }
          `,
        },
        {
          contextValue: createContext({ db, simClient }),
        },
      );

      expect(devResult.body.kind).toBe('single');
      if (devResult.body.kind !== 'single') return;

      const devId = (devResult.body.singleResult.data as any)?.developers.nodes[0].id;

      // Now query commits for that developer
      const result = await server.executeOperation(
        {
          query: `
            query FilteredCommits($userId: ID!, $range: DateRangeInput!) {
              commits(userId: $userId, dateRange: $range, limit: 50) {
                nodes {
                  id
                  timestamp
                  linesAdded
                  author {
                    id
                    name
                  }
                }
                totalCount
                pageInfo {
                  hasNextPage
                  hasPreviousPage
                }
              }
            }
          `,
          variables: {
            userId: devId,
            range: {
              from: threeDaysAgo.toISOString().split('T')[0],
              to: today.toISOString().split('T')[0],
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

        const commits = result.body.singleResult.data?.commits as any;

        expect(commits.totalCount).toBeGreaterThan(0);

        // Verify all commits belong to the filtered user
        for (const commit of commits.nodes) {
          expect(commit.author.id).toBe(devId);

          // Verify timestamp is within range
          const commitDate = new Date(commit.timestamp);
          expect(commitDate.getTime()).toBeGreaterThanOrEqual(threeDaysAgo.getTime());
          expect(commitDate.getTime()).toBeLessThanOrEqual(today.getTime());
        }
      }
    });
  });

  describe('Data Consistency Validation', () => {
    it('should maintain referential integrity across queries', async () => {
      // Get developer count from developers query
      const devsResult = await server.executeOperation(
        {
          query: `
            query CountDevs {
              developers(limit: 100) {
                totalCount
              }
            }
          `,
        },
        {
          contextValue: createContext({ db, simClient }),
        },
      );

      expect(devsResult.body.kind).toBe('single');
      if (devsResult.body.kind !== 'single') return;

      const devCount = (devsResult.body.singleResult.data as any)?.developers.totalCount;

      // Get developer count from dashboard
      const dashResult = await server.executeOperation(
        {
          query: `
            query DashCount {
              dashboardSummary {
                totalDevelopers
              }
            }
          `,
        },
        {
          contextValue: createContext({ db, simClient }),
        },
      );

      expect(dashResult.body.kind).toBe('single');
      if (dashResult.body.kind !== 'single') return;

      const dashCount = (dashResult.body.singleResult.data as any)?.dashboardSummary.totalDevelopers;

      // Both queries should return same developer count
      expect(devCount).toBe(dashCount);
    });

    it('should calculate weighted team averages correctly', async () => {
      const result = await server.executeOperation(
        {
          query: `
            query VerifyWeightedAvg {
              teamStats(teamName: "Backend") {
                memberCount
                averageAcceptanceRate
                totalSuggestions
                totalAccepted
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

        const stats = result.body.singleResult.data?.teamStats as any;

        // Verify weighted average calculation
        if (stats.totalSuggestions > 0) {
          const expectedRate = (stats.totalAccepted / stats.totalSuggestions) * 100;
          expect(Math.abs(stats.averageAcceptanceRate - expectedRate)).toBeLessThan(0.01);
        }
      }
    });
  });
});
