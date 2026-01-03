/**
 * Integration Test Setup
 *
 * Provides utilities for integration tests with real database operations.
 */

import { PrismaClient } from '../../generated/prisma';
import { CursorSimClient } from '../../ingestion/client';

/**
 * Create a fresh test database client
 */
export function createTestDb(): PrismaClient {
  return new PrismaClient({
    datasources: {
      db: {
        url: process.env.DATABASE_URL || 'postgresql://localhost:5432/cursor_analytics_test',
      },
    },
  });
}

/**
 * Create a test cursor-sim client (mocked or pointing to local simulator)
 */
export function createTestSimClient(): CursorSimClient {
  return new CursorSimClient({
    baseUrl: process.env.SIMULATOR_URL || 'http://localhost:8080',
    apiKey: 'test-api-key',
    timeout: 5000,
    retryAttempts: 1,
    retryDelayMs: 100,
  });
}

/**
 * Clean up database after tests
 */
export async function cleanupDb(db: PrismaClient): Promise<void> {
  // Delete in order to respect foreign key constraints
  await db.usageEvent.deleteMany();
  await db.developer.deleteMany();
  await db.$disconnect();
}

/**
 * Seed database with test data
 */
export async function seedTestData(db: PrismaClient) {
  // Create developers
  const dev1 = await db.developer.create({
    data: {
      externalId: 'test-dev-1',
      name: 'Alice Developer',
      email: 'alice@example.com',
      team: 'Team A',
      seniority: 'senior',
    },
  });

  const dev2 = await db.developer.create({
    data: {
      externalId: 'test-dev-2',
      name: 'Bob Developer',
      email: 'bob@example.com',
      team: 'Team A',
      seniority: 'mid',
    },
  });

  const dev3 = await db.developer.create({
    data: {
      externalId: 'test-dev-3',
      name: 'Charlie Developer',
      email: 'charlie@example.com',
      team: 'Team B',
      seniority: 'junior',
    },
  });

  // Create usage events for today
  const today = new Date();
  const yesterday = new Date(today);
  yesterday.setDate(yesterday.getDate() - 1);

  // Alice - high AI usage (10 shown, 8 accepted)
  for (let i = 0; i < 10; i++) {
    await db.usageEvent.create({
      data: {
        externalId: `alice-shown-${i}`,
        developerId: dev1.id,
        eventType: 'cpp_suggestion_shown',
        eventTimestamp: today,
        linesAdded: 0,
        linesDeleted: 0,
        tokensInput: 100,
        tokensOutput: 200,
      },
    });
  }

  for (let i = 0; i < 8; i++) {
    await db.usageEvent.create({
      data: {
        externalId: `alice-accepted-${i}`,
        developerId: dev1.id,
        eventType: 'cpp_suggestion_accepted',
        eventTimestamp: today,
        linesAdded: 10 + i,
        linesDeleted: 2,
        accepted: true,
        tokensInput: 100,
        tokensOutput: 200,
        modelUsed: 'claude-sonnet-4-5',
      },
    });
  }

  // Bob - moderate AI usage (5 shown, 3 accepted)
  for (let i = 0; i < 5; i++) {
    await db.usageEvent.create({
      data: {
        externalId: `bob-shown-${i}`,
        developerId: dev2.id,
        eventType: 'cpp_suggestion_shown',
        eventTimestamp: today,
        linesAdded: 0,
        linesDeleted: 0,
        tokensInput: 50,
        tokensOutput: 100,
      },
    });
  }

  for (let i = 0; i < 3; i++) {
    await db.usageEvent.create({
      data: {
        externalId: `bob-accepted-${i}`,
        developerId: dev2.id,
        eventType: 'cpp_suggestion_accepted',
        eventTimestamp: today,
        linesAdded: 5 + i,
        linesDeleted: 1,
        accepted: true,
        tokensInput: 50,
        tokensOutput: 100,
        modelUsed: 'claude-sonnet-4-5',
      },
    });
  }

  // Charlie - low AI usage yesterday (2 shown, 1 accepted)
  for (let i = 0; i < 2; i++) {
    await db.usageEvent.create({
      data: {
        externalId: `charlie-shown-${i}`,
        developerId: dev3.id,
        eventType: 'cpp_suggestion_shown',
        eventTimestamp: yesterday,
        linesAdded: 0,
        linesDeleted: 0,
        tokensInput: 25,
        tokensOutput: 50,
      },
    });
  }

  await db.usageEvent.create({
    data: {
      externalId: 'charlie-accepted-0',
      developerId: dev3.id,
      eventType: 'cpp_suggestion_accepted',
      eventTimestamp: yesterday,
      linesAdded: 3,
      linesDeleted: 0,
      accepted: true,
      tokensInput: 25,
      tokensOutput: 50,
      modelUsed: 'claude-haiku-4',
    },
  });

  // Add some chat interactions
  await db.usageEvent.create({
    data: {
      externalId: 'alice-chat-1',
      developerId: dev1.id,
      eventType: 'chat_message',
      eventTimestamp: today,
      linesAdded: 0,
      linesDeleted: 0,
      tokensInput: 200,
      tokensOutput: 300,
    },
  });

  await db.usageEvent.create({
    data: {
      externalId: 'bob-cmd-k-1',
      developerId: dev2.id,
      eventType: 'cmd_k_prompt',
      eventTimestamp: today,
      linesAdded: 0,
      linesDeleted: 0,
      tokensInput: 150,
      tokensOutput: 250,
    },
  });

  return { dev1, dev2, dev3 };
}
