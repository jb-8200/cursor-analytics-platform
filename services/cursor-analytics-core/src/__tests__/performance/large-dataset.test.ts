/**
 * Performance Tests: Large Dataset Handling
 *
 * Tests performance with 10,000+ commits to ensure the system
 * can handle production-scale data volumes.
 */

import { ApolloServer } from '@apollo/server';
import { createApolloServer } from '../../graphql/server';
import { createContext, type GraphQLContext } from '../../graphql/context';
import { createTestDb, createTestSimClient, cleanupDb } from '../integration/setup';
import type { PrismaClient } from '../../generated/prisma';
import type { CursorSimClient } from '../../ingestion/client';

describe('Performance Tests: Large Dataset', () => {
  let db: PrismaClient;
  let simClient: CursorSimClient;
  let server: ApolloServer<GraphQLContext>;

  // Increase timeout for performance tests
  jest.setTimeout(30000);

  beforeAll(async () => {
    db = createTestDb();
    simClient = createTestSimClient();
    server = createApolloServer({ db, simClient });

    console.log('Seeding large dataset (10,000+ events)...');
    const startTime = Date.now();
    await seedLargeDataset(db);
    const seedTime = Date.now() - startTime;
    console.log(`Dataset seeded in ${seedTime}ms`);
  });

  afterAll(async () => {
    await cleanupDb(db);
  });

  /**
   * Seed database with 10,000+ events across 50 developers
   */
  async function seedLargeDataset(db: PrismaClient) {
    const teams = ['Engineering', 'Product', 'Data', 'DevOps', 'QA'];
    const seniorities = ['junior', 'mid', 'senior'];
    const models = ['claude-sonnet-4-5', 'claude-haiku-4', 'claude-opus-4-5'];

    // Create 50 developers
    const developers = [];
    for (let i = 0; i < 50; i++) {
      const dev = await db.developer.create({
        data: {
          externalId: `perf-dev-${i}`,
          name: `Developer ${i}`,
          email: `dev${i}@company.com`,
          team: teams[i % teams.length],
          seniority: seniorities[i % seniorities.length],
        },
      });
      developers.push(dev);
    }

    // Generate events for last 30 days
    const today = new Date();
    const eventBatches = [];

    for (let dayOffset = 0; dayOffset < 30; dayOffset++) {
      const date = new Date(today);
      date.setDate(date.getDate() - dayOffset);

      for (const dev of developers) {
        // Each developer generates ~7 events per day
        // 50 devs × 7 events × 30 days = 10,500 events

        // 3 suggestions shown
        for (let i = 0; i < 3; i++) {
          eventBatches.push({
            externalId: `perf-${dev.externalId}-day${dayOffset}-shown-${i}`,
            developerId: dev.id,
            eventType: 'cpp_suggestion_shown',
            eventTimestamp: date,
            linesAdded: 0,
            linesDeleted: 0,
            tokensInput: 100 + Math.floor(Math.random() * 100),
            tokensOutput: 150 + Math.floor(Math.random() * 150),
          });
        }

        // 2 accepted suggestions
        for (let i = 0; i < 2; i++) {
          eventBatches.push({
            externalId: `perf-${dev.externalId}-day${dayOffset}-accepted-${i}`,
            developerId: dev.id,
            eventType: 'cpp_suggestion_accepted',
            eventTimestamp: date,
            linesAdded: 5 + Math.floor(Math.random() * 20),
            linesDeleted: Math.floor(Math.random() * 5),
            accepted: true,
            tokensInput: 100 + Math.floor(Math.random() * 100),
            tokensOutput: 150 + Math.floor(Math.random() * 150),
            modelUsed: models[Math.floor(Math.random() * models.length)],
          });
        }

        // 1 chat message
        eventBatches.push({
          externalId: `perf-${dev.externalId}-day${dayOffset}-chat`,
          developerId: dev.id,
          eventType: 'chat_message',
          eventTimestamp: date,
          linesAdded: 0,
          linesDeleted: 0,
          tokensInput: 200 + Math.floor(Math.random() * 300),
          tokensOutput: 300 + Math.floor(Math.random() * 500),
        });

        // 1 cmd_k usage
        eventBatches.push({
          externalId: `perf-${dev.externalId}-day${dayOffset}-cmdk`,
          developerId: dev.id,
          eventType: 'cmd_k_prompt',
          eventTimestamp: date,
          linesAdded: 0,
          linesDeleted: 0,
          tokensInput: 150 + Math.floor(Math.random() * 200),
          tokensOutput: 250 + Math.floor(Math.random() * 400),
        });
      }
    }

    // Batch insert for performance (chunks of 500)
    const chunkSize = 500;
    for (let i = 0; i < eventBatches.length; i += chunkSize) {
      const chunk = eventBatches.slice(i, i + chunkSize);
      await db.usageEvent.createMany({
        data: chunk,
      });
    }

    console.log(`Created ${eventBatches.length} events for 50 developers over 30 days`);
  }

  describe('Query Performance with Large Dataset', () => {
    it('should fetch dashboard summary in < 2000ms', async () => {
      const startTime = Date.now();

      const result = await server.executeOperation(
        {
          query: `
            query DashboardPerf {
              dashboardSummary(preset: LAST_30_DAYS) {
                totalDevelopers
                activeDevelopers
                overallAcceptanceRate
                totalSuggestionsToday
                totalAcceptedToday
                teamComparison {
                  teamName
                  memberCount
                  averageAcceptanceRate
                  totalSuggestions
                  totalAccepted
                }
                dailyTrend {
                  date
                  suggestionsShown
                  suggestionsAccepted
                }
              }
            }
          `,
        },
        {
          contextValue: createContext({ db, simClient }),
        },
      );

      const queryTime = Date.now() - startTime;

      expect(result.body.kind).toBe('single');
      if (result.body.kind === 'single') {
        expect(result.body.singleResult.errors).toBeUndefined();

        const summary = result.body.singleResult.data?.dashboardSummary as any;
        expect(summary.totalDevelopers).toBe(50);
        expect(summary.teamComparison.length).toBeGreaterThan(0);
        expect(summary.dailyTrend.length).toBe(30);
      }

      console.log(`Dashboard query completed in ${queryTime}ms`);
      expect(queryTime).toBeLessThan(2000); // Should complete in under 2 seconds
    });

    it('should list 100 developers with stats in < 1000ms', async () => {
      const today = new Date();
      const sevenDaysAgo = new Date(today);
      sevenDaysAgo.setDate(sevenDaysAgo.getDate() - 7);

      const startTime = Date.now();

      const result = await server.executeOperation(
        {
          query: `
            query DevelopersPerf($range: DateRangeInput!) {
              developers(limit: 100) {
                nodes {
                  id
                  name
                  team
                  stats(range: $range) {
                    totalSuggestions
                    acceptedSuggestions
                    acceptanceRate
                    aiVelocity
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

      const queryTime = Date.now() - startTime;

      expect(result.body.kind).toBe('single');
      if (result.body.kind === 'single') {
        expect(result.body.singleResult.errors).toBeUndefined();

        const developers = result.body.singleResult.data?.developers as any;
        expect(developers.totalCount).toBe(50);
        expect(developers.nodes.length).toBe(50);

        // Verify stats are calculated for each developer
        for (const dev of developers.nodes) {
          expect(dev.stats).toBeDefined();
          expect(dev.stats.totalSuggestions).toBeGreaterThan(0);
        }
      }

      console.log(`Developers list query completed in ${queryTime}ms`);
      expect(queryTime).toBeLessThan(1000); // Should complete in under 1 second
    });

    it('should paginate through 1000 commits in < 500ms', async () => {
      const startTime = Date.now();

      const result = await server.executeOperation(
        {
          query: `
            query CommitsPerf {
              commits(limit: 100, sortBy: "timestamp", sortOrder: "desc") {
                nodes {
                  id
                  linesAdded
                  modelUsed
                  timestamp
                  author {
                    name
                    team
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

      const queryTime = Date.now() - startTime;

      expect(result.body.kind).toBe('single');
      if (result.body.kind === 'single') {
        expect(result.body.singleResult.errors).toBeUndefined();

        const commits = result.body.singleResult.data?.commits as any;
        expect(commits.totalCount).toBeGreaterThan(1000);
        expect(commits.nodes.length).toBe(100);

        // Verify each commit has author populated (no N+1 queries)
        for (const commit of commits.nodes) {
          expect(commit.author).toBeDefined();
          expect(commit.author.name).toBeTruthy();
        }
      }

      console.log(`Commits pagination query completed in ${queryTime}ms`);
      expect(queryTime).toBeLessThan(500); // Should complete in under 500ms
    });

    it('should aggregate team stats across all teams in < 800ms', async () => {
      const startTime = Date.now();

      const result = await server.executeOperation(
        {
          query: `
            query AllTeamsPerf {
              teams {
                teamName
                memberCount
                activeMemberCount
                averageAcceptanceRate
                totalSuggestions
                totalAccepted
                aiVelocity
                topPerformer {
                  name
                  email
                }
              }
            }
          `,
        },
        {
          contextValue: createContext({ db, simClient }),
        },
      );

      const queryTime = Date.now() - startTime;

      expect(result.body.kind).toBe('single');
      if (result.body.kind === 'single') {
        expect(result.body.singleResult.errors).toBeUndefined();

        const teams = result.body.singleResult.data?.teams as any;
        expect(teams.length).toBe(5); // Engineering, Product, Data, DevOps, QA

        for (const team of teams) {
          expect(team.memberCount).toBeGreaterThan(0);
          expect(team.topPerformer).toBeDefined();
        }
      }

      console.log(`All teams aggregation completed in ${queryTime}ms`);
      expect(queryTime).toBeLessThan(800); // Should complete in under 800ms
    });

    it('should filter commits by date range efficiently in < 300ms', async () => {
      const today = new Date();
      const threeDaysAgo = new Date(today);
      threeDaysAgo.setDate(threeDaysAgo.getDate() - 3);

      const startTime = Date.now();

      const result = await server.executeOperation(
        {
          query: `
            query FilteredCommitsPerf($range: DateRangeInput!) {
              commits(dateRange: $range, limit: 100) {
                nodes {
                  id
                  timestamp
                  linesAdded
                }
                totalCount
              }
            }
          `,
          variables: {
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

      const queryTime = Date.now() - startTime;

      expect(result.body.kind).toBe('single');
      if (result.body.kind === 'single') {
        expect(result.body.singleResult.errors).toBeUndefined();

        const commits = result.body.singleResult.data?.commits as any;
        expect(commits.totalCount).toBeGreaterThan(0);

        // Verify all commits are within date range
        for (const commit of commits.nodes) {
          const commitDate = new Date(commit.timestamp);
          expect(commitDate.getTime()).toBeGreaterThanOrEqual(threeDaysAgo.getTime());
          expect(commitDate.getTime()).toBeLessThanOrEqual(today.getTime());
        }
      }

      console.log(`Filtered commits query completed in ${queryTime}ms`);
      expect(queryTime).toBeLessThan(300); // Should complete in under 300ms
    });
  });

  describe('Memory and Resource Usage', () => {
    it('should handle multiple concurrent queries without degradation', async () => {
      const queries = Array(10)
        .fill(null)
        .map(() =>
          server.executeOperation(
            {
              query: `
                query ConcurrentTest {
                  developers(limit: 20) {
                    nodes {
                      name
                      stats {
                        totalSuggestions
                        acceptanceRate
                      }
                    }
                  }
                }
              `,
            },
            {
              contextValue: createContext({ db, simClient }),
            },
          ),
        );

      const startTime = Date.now();
      const results = await Promise.all(queries);
      const totalTime = Date.now() - startTime;

      // All queries should succeed
      for (const result of results) {
        expect(result.body.kind).toBe('single');
        if (result.body.kind === 'single') {
          expect(result.body.singleResult.errors).toBeUndefined();
        }
      }

      console.log(`10 concurrent queries completed in ${totalTime}ms`);
      expect(totalTime).toBeLessThan(3000); // Should complete in under 3 seconds
    });
  });

  describe('Pagination Performance', () => {
    it('should paginate through large result sets efficiently', async () => {
      const pageSize = 50;
      const totalPages = 5; // 250 total commits

      const timings = [];

      for (let page = 0; page < totalPages; page++) {
        const startTime = Date.now();

        const result = await server.executeOperation(
          {
            query: `
              query PaginationTest($offset: Int!, $limit: Int!) {
                commits(offset: $offset, limit: $limit) {
                  nodes {
                    id
                    linesAdded
                    author {
                      name
                    }
                  }
                  pageInfo {
                    hasNextPage
                  }
                }
              }
            `,
            variables: {
              offset: page * pageSize,
              limit: pageSize,
            },
          },
          {
            contextValue: createContext({ db, simClient }),
          },
        );

        const pageTime = Date.now() - startTime;
        timings.push(pageTime);

        expect(result.body.kind).toBe('single');
        if (result.body.kind === 'single') {
          expect(result.body.singleResult.errors).toBeUndefined();
        }
      }

      const avgTime = timings.reduce((a, b) => a + b, 0) / timings.length;
      console.log(`Average pagination time: ${avgTime.toFixed(2)}ms`);
      console.log(`Page timings: ${timings.join(', ')}ms`);

      // Average page load should be under 400ms
      expect(avgTime).toBeLessThan(400);

      // Performance should not degrade significantly across pages
      const firstPageTime = timings[0];
      const lastPageTime = timings[timings.length - 1];
      expect(lastPageTime).toBeLessThan(firstPageTime * 2); // No more than 2x slower
    });
  });
});
