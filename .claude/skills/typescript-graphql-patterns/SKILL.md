---
name: typescript-graphql-patterns
description: TypeScript and GraphQL best practices for cursor-analytics-core. Use when implementing GraphQL schemas, resolvers, TypeScript services, or data fetching. Covers type safety, error handling, and Apollo Server patterns. (project)
---

# TypeScript & GraphQL Patterns

Best practices for cursor-analytics-core development.

## TypeScript Configuration

### Strict Mode Required

```json
{
  "compilerOptions": {
    "strict": true,
    "noImplicitAny": true,
    "strictNullChecks": true,
    "noImplicitReturns": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "esModuleInterop": true
  }
}
```

## GraphQL Schema Design

### Type-First Development

Define GraphQL types that mirror cursor-sim models:

```graphql
type Developer {
  id: ID!
  name: String!
  email: String!
  seniority: Seniority!
  aiPreference: Float!
  active: Boolean!
}

enum Seniority {
  JUNIOR
  MID
  SENIOR
}

type Commit {
  sha: ID!
  message: String!
  author: Developer!
  timestamp: DateTime!
  aiLinesAdded: Int!
  humanLinesAdded: Int!
  aiRatio: Float!
}
```

### Query Design

Support flexible filtering and aggregation:

```graphql
type Query {
  # Team-level queries
  teamOverview(dateRange: DateRangeInput!): TeamOverview!
  teamUsage(dateRange: DateRangeInput!): TeamUsage!

  # Developer queries
  developers(
    filter: DeveloperFilter
    pagination: PaginationInput
  ): DeveloperConnection!

  developer(id: ID!): Developer

  # Analytics queries with aggregation
  aiUsageTrend(
    dateRange: DateRangeInput!
    granularity: Granularity!
  ): [TrendDataPoint!]!
}

input DateRangeInput {
  startDate: Date!
  endDate: Date!
}

enum Granularity {
  DAILY
  WEEKLY
  MONTHLY
}
```

## Resolver Patterns

### Type-Safe Resolvers

```typescript
import { Resolvers } from './generated/types';

const resolvers: Resolvers = {
  Query: {
    teamOverview: async (_, { dateRange }, { dataSources }) => {
      return dataSources.cursorSim.getTeamOverview(dateRange);
    },

    developers: async (_, { filter, pagination }, { dataSources }) => {
      const data = await dataSources.cursorSim.getDevelopers(filter);
      return paginateResults(data, pagination);
    },
  },

  Developer: {
    // Field resolver for computed fields
    totalAiLines: (parent) => {
      return parent.commits.reduce((sum, c) => sum + c.aiLinesAdded, 0);
    },
  },
};
```

### DataLoader for N+1 Prevention

```typescript
import DataLoader from 'dataloader';

const createLoaders = (cursorSimClient: CursorSimClient) => ({
  developerLoader: new DataLoader<string, Developer>(async (ids) => {
    const developers = await cursorSimClient.getDevelopersByIds(ids);
    return ids.map(id => developers.find(d => d.id === id) ?? null);
  }),
});
```

## Error Handling

### Custom GraphQL Errors

```typescript
import { GraphQLError } from 'graphql';

export class NotFoundError extends GraphQLError {
  constructor(resource: string, id: string) {
    super(`${resource} not found: ${id}`, {
      extensions: {
        code: 'NOT_FOUND',
        resource,
        id,
      },
    });
  }
}

export class ValidationError extends GraphQLError {
  constructor(message: string, field: string) {
    super(message, {
      extensions: {
        code: 'VALIDATION_ERROR',
        field,
      },
    });
  }
}
```

### Error Handling in Resolvers

```typescript
const resolvers: Resolvers = {
  Query: {
    developer: async (_, { id }, { dataSources }) => {
      const developer = await dataSources.cursorSim.getDeveloper(id);
      if (!developer) {
        throw new NotFoundError('Developer', id);
      }
      return developer;
    },
  },
};
```

## cursor-sim Client

### Type-Safe API Client

```typescript
interface CursorSimClient {
  getTeamOverview(dateRange: DateRange): Promise<TeamOverview>;
  getDevelopers(filter?: DeveloperFilter): Promise<Developer[]>;
  getCommits(repoName: string, dateRange: DateRange): Promise<Commit[]>;
}

class CursorSimApiClient implements CursorSimClient {
  constructor(
    private baseUrl: string,
    private apiKey: string,
  ) {}

  private async fetch<T>(endpoint: string, params?: Record<string, string>): Promise<T> {
    const url = new URL(endpoint, this.baseUrl);
    if (params) {
      Object.entries(params).forEach(([k, v]) => url.searchParams.set(k, v));
    }

    const response = await fetch(url.toString(), {
      headers: {
        'Authorization': `Basic ${Buffer.from(this.apiKey + ':').toString('base64')}`,
      },
    });

    if (!response.ok) {
      throw new Error(`cursor-sim API error: ${response.status}`);
    }

    return response.json();
  }

  async getTeamOverview(dateRange: DateRange): Promise<TeamOverview> {
    return this.fetch('/analytics/team/overview', {
      start_date: dateRange.startDate,
      end_date: dateRange.endDate,
    });
  }
}
```

## Testing Patterns

### Unit Tests for Resolvers

```typescript
import { createTestServer } from './test-utils';

describe('Developer Queries', () => {
  it('returns developer by ID', async () => {
    const server = createTestServer({
      mocks: {
        cursorSim: {
          getDeveloper: jest.fn().mockResolvedValue({
            id: 'dev-1',
            name: 'Alice',
          }),
        },
      },
    });

    const result = await server.executeOperation({
      query: `
        query GetDeveloper($id: ID!) {
          developer(id: $id) {
            id
            name
          }
        }
      `,
      variables: { id: 'dev-1' },
    });

    expect(result.data?.developer.name).toBe('Alice');
  });
});
```

### Integration Tests

```typescript
describe('cursor-sim Integration', () => {
  it('fetches real data from cursor-sim', async () => {
    const client = new CursorSimApiClient(
      process.env.CURSOR_SIM_URL!,
      process.env.CURSOR_SIM_API_KEY!,
    );

    const overview = await client.getTeamOverview({
      startDate: '2025-01-01',
      endDate: '2025-01-31',
    });

    expect(overview).toHaveProperty('total_ai_lines');
  });
});
```

## Code Generation

### Generate Types from Schema

```bash
# codegen.yml
generates:
  src/generated/types.ts:
    plugins:
      - typescript
      - typescript-resolvers
    config:
      contextType: ../context#Context
      mappers:
        Developer: ../models#DeveloperModel
```

Run: `npm run codegen`

## Performance

### Query Complexity Limits

```typescript
import { createComplexityPlugin } from 'graphql-query-complexity';

const complexityPlugin = createComplexityPlugin({
  maximumComplexity: 1000,
  variables: {},
  onComplete: (complexity) => {
    console.log('Query Complexity:', complexity);
  },
});
```

### Response Caching

```typescript
import responseCachePlugin from '@apollo/server-plugin-response-cache';

const server = new ApolloServer({
  plugins: [
    responseCachePlugin({
      sessionId: (requestContext) =>
        requestContext.request.http?.headers.get('authorization') ?? null,
    }),
  ],
});
```
