# Task: TASK-010

## Set Up Node.js TypeScript Project for Aggregator

**Task ID:** TASK-010  
**Service:** cursor-analytics-core  
**Feature:** [F002 - Aggregator Ingestion](../features/F002-aggregator-ingestion.md)  
**User Story:** [US-AGG-001](../user-stories/US-AGG-001-sync-data.md)  
**Estimated Hours:** 4  
**Status:** Ready

---

## Objective

Initialize the Node.js TypeScript project with proper structure, dependencies, and tooling to support test-driven development of the Aggregator Service.

---

## Prerequisites

Node.js 20 or later must be installed. Docker and Docker Compose should be available for local PostgreSQL. The cursor-sim service should be implemented or mocked for integration testing.

---

## Implementation Steps

### Step 1: Initialize Node.js Project

Create the package.json with appropriate metadata and scripts.

```bash
cd services/cursor-analytics-core
npm init -y
```

Then update package.json with the following configuration:

```json
{
  "name": "cursor-analytics-core",
  "version": "1.0.0",
  "description": "Aggregator service for Cursor usage analytics",
  "main": "dist/index.js",
  "scripts": {
    "build": "tsc",
    "start": "node dist/index.js",
    "dev": "tsx watch src/index.ts",
    "test": "jest",
    "test:watch": "jest --watch",
    "test:coverage": "jest --coverage",
    "lint": "eslint src --ext .ts",
    "lint:fix": "eslint src --ext .ts --fix",
    "typecheck": "tsc --noEmit",
    "db:migrate": "prisma migrate dev",
    "db:generate": "prisma generate",
    "db:push": "prisma db push",
    "codegen": "graphql-codegen"
  },
  "engines": {
    "node": ">=20.0.0"
  }
}
```

### Step 2: Create Directory Structure

The project uses a layered architecture with clear separation between API, services, and data access.

```
services/cursor-analytics-core/
├── src/
│   ├── index.ts                 # Entry point
│   ├── server.ts                # Apollo Server setup
│   ├── config/
│   │   ├── index.ts             # Configuration loading
│   │   └── database.ts          # Database connection
│   ├── graphql/
│   │   ├── schema.graphql       # GraphQL type definitions
│   │   ├── resolvers/
│   │   │   ├── index.ts         # Resolver composition
│   │   │   ├── developer.ts     # Developer resolvers
│   │   │   ├── stats.ts         # Statistics resolvers
│   │   │   └── dashboard.ts     # Dashboard resolvers
│   │   └── loaders/
│   │       ├── index.ts         # DataLoader factory
│   │       ├── developer.ts     # Developer DataLoader
│   │       └── stats.ts         # Stats DataLoader
│   ├── services/
│   │   ├── ingestion/
│   │   │   ├── worker.ts        # Background sync worker
│   │   │   ├── fetcher.ts       # Simulator API client
│   │   │   └── transformer.ts   # Data transformation
│   │   ├── kpi/
│   │   │   ├── calculator.ts    # KPI calculation logic
│   │   │   └── aggregator.ts    # Team aggregations
│   │   └── sync-state.ts        # Sync state management
│   ├── models/
│   │   ├── developer.ts         # Developer type
│   │   ├── event.ts             # Event types
│   │   └── stats.ts             # Statistics types
│   └── utils/
│       ├── logger.ts            # Pino logger setup
│       └── errors.ts            # Custom error classes
├── prisma/
│   ├── schema.prisma            # Database schema
│   └── migrations/              # Migration files
├── tests/
│   ├── setup.ts                 # Jest setup
│   ├── mocks/                   # Mock data factories
│   ├── unit/                    # Unit tests
│   └── integration/             # Integration tests
├── package.json
├── tsconfig.json
├── jest.config.js
├── .eslintrc.js
├── codegen.ts                   # GraphQL codegen config
├── Dockerfile
└── README.md
```

### Step 3: Install Dependencies

Install all required dependencies in the correct categories.

```bash
# Runtime dependencies
npm install express @apollo/server graphql graphql-tag dataloader
npm install pg prisma @prisma/client
npm install node-cron zod pino pino-pretty
npm install axios

# Development dependencies
npm install -D typescript tsx @types/node @types/express
npm install -D jest ts-jest @types/jest
npm install -D eslint @typescript-eslint/parser @typescript-eslint/eslint-plugin
npm install -D @graphql-codegen/cli @graphql-codegen/typescript @graphql-codegen/typescript-resolvers
npm install -D pg-mem
```

### Step 4: Configure TypeScript

Create tsconfig.json with strict type checking enabled.

```json
{
  "compilerOptions": {
    "target": "ES2022",
    "module": "NodeNext",
    "moduleResolution": "NodeNext",
    "lib": ["ES2022"],
    "outDir": "./dist",
    "rootDir": "./src",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true,
    "resolveJsonModule": true,
    "declaration": true,
    "declarationMap": true,
    "sourceMap": true,
    "noImplicitReturns": true,
    "noFallthroughCasesInSwitch": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true
  },
  "include": ["src/**/*"],
  "exclude": ["node_modules", "dist", "tests"]
}
```

### Step 5: Configure Jest

Create jest.config.js for TypeScript testing with pg-mem for database tests.

```javascript
// jest.config.js
module.exports = {
  preset: 'ts-jest',
  testEnvironment: 'node',
  roots: ['<rootDir>/src', '<rootDir>/tests'],
  testMatch: ['**/*.test.ts'],
  transform: {
    '^.+\\.ts$': 'ts-jest',
  },
  collectCoverageFrom: [
    'src/**/*.ts',
    '!src/**/*.d.ts',
    '!src/index.ts',
  ],
  coverageThreshold: {
    global: {
      branches: 80,
      functions: 80,
      lines: 80,
      statements: 80,
    },
  },
  setupFilesAfterEnv: ['<rootDir>/tests/setup.ts'],
  moduleNameMapper: {
    '^@/(.*)$': '<rootDir>/src/$1',
  },
};
```

### Step 6: Create Prisma Schema

Initialize Prisma and create the database schema.

```prisma
// prisma/schema.prisma
generator client {
  provider = "prisma-client-js"
}

datasource db {
  provider = "postgresql"
  url      = env("DATABASE_URL")
}

model Developer {
  id        String   @id @default(uuid())
  email     String   @unique
  name      String
  team      String
  role      String
  createdAt DateTime @default(now()) @map("created_at")
  
  dailyStats DailyStats[]
  
  @@map("developers")
}

model DailyStats {
  id                 String   @id @default(uuid())
  developerId        String   @map("developer_id")
  date               DateTime @db.Date
  totalTabsShown     Int      @default(0) @map("total_tabs_shown")
  totalTabsAccepted  Int      @default(0) @map("total_tabs_accepted")
  linesAdded         Int      @default(0) @map("lines_added")
  linesDeleted       Int      @default(0) @map("lines_deleted")
  acceptedLinesAdded Int      @default(0) @map("accepted_lines_added")
  chatRequests       Int      @default(0) @map("chat_requests")
  composerRequests   Int      @default(0) @map("composer_requests")
  agentRequests      Int      @default(0) @map("agent_requests")
  cmdkUsages         Int      @default(0) @map("cmdk_usages")
  mostUsedModel      String?  @map("most_used_model")
  createdAt          DateTime @default(now()) @map("created_at")
  
  developer Developer @relation(fields: [developerId], references: [id])
  
  @@unique([developerId, date])
  @@index([date])
  @@map("daily_stats")
}

model TeamStats {
  id                String   @id @default(uuid())
  teamName          String   @map("team_name")
  date              DateTime @db.Date
  avgAcceptanceRate Float    @map("avg_acceptance_rate")
  totalAiLines      Int      @map("total_ai_lines")
  totalChatRequests Int      @map("total_chat_requests")
  activeDevelopers  Int      @map("active_developers")
  aiVelocityScore   Float    @map("ai_velocity_score")
  
  @@unique([teamName, date])
  @@map("team_stats")
}

model SyncState {
  id                    String   @id @default(uuid())
  lastSyncTime          DateTime @map("last_sync_time")
  lastSuccessfulSync    DateTime @map("last_successful_sync")
  consecutiveFailures   Int      @default(0) @map("consecutive_failures")
  totalRecordsSynced    Int      @default(0) @map("total_records_synced")
  
  @@map("sync_state")
}
```

### Step 7: Create Initial Test Files

Write failing tests that define the expected behavior of core services.

```typescript
// tests/unit/kpi/calculator.test.ts
import { kpiCalculator } from '@/services/kpi/calculator';

describe('KPI Calculator', () => {
  describe('calculateAcceptanceRate', () => {
    it('should calculate percentage of accepted suggestions', () => {
      const result = kpiCalculator.calculateAcceptanceRate(100, 75);
      expect(result).toBe(75);
    });

    it('should return 0 when no suggestions shown', () => {
      const result = kpiCalculator.calculateAcceptanceRate(0, 0);
      expect(result).toBe(0);
    });

    it('should handle decimal results correctly', () => {
      const result = kpiCalculator.calculateAcceptanceRate(3, 1);
      expect(result).toBeCloseTo(33.33, 1);
    });
  });

  describe('calculateAIVelocity', () => {
    it('should calculate AI contribution percentage', () => {
      const result = kpiCalculator.calculateAIVelocity(50, 100);
      expect(result).toBe(50);
    });

    it('should return 0 when no lines written', () => {
      const result = kpiCalculator.calculateAIVelocity(0, 0);
      expect(result).toBe(0);
    });
  });

  describe('calculateChatDependency', () => {
    it('should calculate ratio of chat to code events', () => {
      const result = kpiCalculator.calculateChatDependency(10, 100);
      expect(result).toBeCloseTo(0.099, 2);
    });

    it('should handle zero code events without division error', () => {
      const result = kpiCalculator.calculateChatDependency(10, 0);
      expect(result).toBe(10);
    });
  });
});
```

### Step 8: Create Dockerfile

Create a multi-stage Dockerfile for production builds.

```dockerfile
# Dockerfile
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./
COPY prisma ./prisma/

# Install dependencies
RUN npm ci

# Generate Prisma client
RUN npx prisma generate

# Copy source code
COPY . .

# Build TypeScript
RUN npm run build

# Production stage
FROM node:20-alpine

WORKDIR /app

# Copy built files and dependencies
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package*.json ./
COPY --from=builder /app/prisma ./prisma

# Expose GraphQL port
EXPOSE 4000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:4000/health || exit 1

CMD ["node", "dist/index.js"]
```

---

## Test-Driven Development Checklist

- [ ] Test files created before implementation
- [ ] Tests fail initially (as expected)
- [ ] Tests use descriptive names
- [ ] Mock data factories created
- [ ] pg-mem configured for database tests

---

## Definition of Done

- [ ] package.json configured with all scripts
- [ ] TypeScript compiles without errors
- [ ] Prisma schema created and generates client
- [ ] Jest runs tests (failing is expected)
- [ ] ESLint configured and runs without errors
- [ ] Dockerfile builds successfully
- [ ] README.md created
- [ ] Code committed to feature branch
