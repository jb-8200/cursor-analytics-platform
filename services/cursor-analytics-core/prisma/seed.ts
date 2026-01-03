import { PrismaClient } from '../src/generated/prisma';

const prisma = new PrismaClient();

async function main() {
  console.log('Starting database seed...');

  // Clear existing data
  console.log('Clearing existing data...');
  await prisma.usageEvent.deleteMany();
  await prisma.developer.deleteMany();

  // Create test developers
  console.log('Creating test developers...');
  const developers = await Promise.all([
    prisma.developer.create({
      data: {
        externalId: 'dev-001',
        name: 'Alice Johnson',
        email: 'alice@example.com',
        team: 'Platform',
        seniority: 'senior',
      },
    }),
    prisma.developer.create({
      data: {
        externalId: 'dev-002',
        name: 'Bob Smith',
        email: 'bob@example.com',
        team: 'Platform',
        seniority: 'mid',
      },
    }),
    prisma.developer.create({
      data: {
        externalId: 'dev-003',
        name: 'Charlie Davis',
        email: 'charlie@example.com',
        team: 'Product',
        seniority: 'junior',
      },
    }),
    prisma.developer.create({
      data: {
        externalId: 'dev-004',
        name: 'Diana Martinez',
        email: 'diana@example.com',
        team: 'Product',
        seniority: 'senior',
      },
    }),
    prisma.developer.create({
      data: {
        externalId: 'dev-005',
        name: 'Ethan Wilson',
        email: 'ethan@example.com',
        team: 'Infrastructure',
        seniority: 'mid',
      },
    }),
  ]);

  console.log(`Created ${developers.length} developers`);

  // Create test usage events
  console.log('Creating test usage events...');
  const now = new Date();
  const dayMs = 24 * 60 * 60 * 1000;

  const events = [];

  // Generate events for the past 7 days
  for (let day = 0; day < 7; day++) {
    const eventDate = new Date(now.getTime() - day * dayMs);

    for (const developer of developers) {
      // Suggestions shown
      for (let i = 0; i < 10 + Math.floor(Math.random() * 20); i++) {
        events.push({
          externalId: `event-${developer.externalId}-${day}-shown-${i}`,
          developerId: developer.id,
          eventType: 'cpp_suggestion_shown',
          eventTimestamp: new Date(
            eventDate.getTime() + Math.random() * dayMs,
          ),
          linesAdded: 0,
          linesDeleted: 0,
          modelUsed: Math.random() > 0.5 ? 'claude-sonnet-4' : 'claude-opus-4',
          accepted: null,
          tokensInput: Math.floor(Math.random() * 1000) + 100,
          tokensOutput: Math.floor(Math.random() * 500) + 50,
        });
      }

      // Suggestions accepted (70-90% acceptance rate)
      const acceptanceRate = 0.7 + Math.random() * 0.2;
      const acceptedCount = Math.floor(10 * acceptanceRate);

      for (let i = 0; i < acceptedCount; i++) {
        events.push({
          externalId: `event-${developer.externalId}-${day}-accepted-${i}`,
          developerId: developer.id,
          eventType: 'cpp_suggestion_accepted',
          eventTimestamp: new Date(
            eventDate.getTime() + Math.random() * dayMs,
          ),
          linesAdded: Math.floor(Math.random() * 50) + 5,
          linesDeleted: Math.floor(Math.random() * 20),
          modelUsed: Math.random() > 0.5 ? 'claude-sonnet-4' : 'claude-opus-4',
          accepted: true,
          tokensInput: 0,
          tokensOutput: 0,
        });
      }

      // Chat interactions
      for (let i = 0; i < Math.floor(Math.random() * 5); i++) {
        events.push({
          externalId: `event-${developer.externalId}-${day}-chat-${i}`,
          developerId: developer.id,
          eventType: 'chat_message',
          eventTimestamp: new Date(
            eventDate.getTime() + Math.random() * dayMs,
          ),
          linesAdded: 0,
          linesDeleted: 0,
          modelUsed: 'claude-sonnet-4',
          accepted: null,
          tokensInput: Math.floor(Math.random() * 2000) + 500,
          tokensOutput: Math.floor(Math.random() * 1500) + 300,
        });
      }

      // Cmd+K usages
      for (let i = 0; i < Math.floor(Math.random() * 3); i++) {
        events.push({
          externalId: `event-${developer.externalId}-${day}-cmdk-${i}`,
          developerId: developer.id,
          eventType: 'cmd_k_prompt',
          eventTimestamp: new Date(
            eventDate.getTime() + Math.random() * dayMs,
          ),
          linesAdded: Math.floor(Math.random() * 100) + 10,
          linesDeleted: Math.floor(Math.random() * 50),
          modelUsed: 'claude-sonnet-4',
          accepted: true,
          tokensInput: Math.floor(Math.random() * 3000) + 1000,
          tokensOutput: Math.floor(Math.random() * 2000) + 500,
        });
      }
    }
  }

  // Batch insert events
  console.log(`Creating ${events.length} usage events...`);
  await prisma.usageEvent.createMany({
    data: events,
  });

  console.log('Seed completed successfully!');
  console.log(`- ${developers.length} developers created`);
  console.log(`- ${events.length} usage events created`);
  console.log('- Events span 7 days of simulated activity');
}

main()
  .catch((e) => {
    console.error('Error during seed:', e);
    process.exit(1);
  })
  .finally(async () => {
    await prisma.$disconnect();
  });
