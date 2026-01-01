# Testing Strategy: Cursor Usage Analytics Platform

**Version**: 1.0.0  
**Last Updated**: January 2026  

This document defines the testing approach for the Cursor Usage Analytics Platform, including the recommended testing stack, TDD workflows, and coverage requirements. The strategy is designed to support spec-driven development where tests are written from specifications before implementation code.

## Testing Philosophy

The project follows Test-Driven Development (TDD) as its primary development methodology. This means tests are written before implementation code, serving as executable specifications that define the expected behavior of each component. The testing pyramid guides our approach: many unit tests form the foundation, integration tests verify component interactions, and a smaller number of end-to-end tests validate complete user workflows.

Writing tests first provides several benefits that improve the overall quality of the codebase. Tests act as documentation that never becomes stale since failing tests indicate incorrect behavior immediately. The TDD cycle of red-green-refactor encourages simple, focused implementations that do exactly what the specification requires. Test-first development also leads to more testable designs because developers must consider how code will be tested before writing it.

## Testing Stack by Service

Each service uses testing tools appropriate for its technology stack while maintaining consistent patterns across the project.

### Service A: cursor-sim (Go)

The Go testing stack leverages the language's built-in testing capabilities while adding assertion libraries for cleaner test code.

**Core Tools:**

The standard `testing` package provides the foundation for all tests. Go's testing conventions require test files to be named with a `_test.go` suffix and test functions to start with `Test`. The built-in test runner supports parallel execution, benchmarks, and coverage reporting out of the box.

**Testify** (`github.com/stretchr/testify`) adds assertion functions that make test expectations more readable. Instead of writing manual comparison code with conditional failures, developers can write expressive assertions like `assert.Equal(t, expected, actual)` or `require.NoError(t, err)`. The `require` package variant stops test execution immediately on failure, useful for setup validations.

**Mockery** (`github.com/vektra/mockery`) generates mock implementations from interfaces automatically. This enables testing components in isolation by replacing dependencies with mocks that return predetermined values. Generated mocks track method calls, allowing assertions about how dependencies were used.

**Example Test Structure:**

```go
// internal/generator/developer_generator_test.go
package generator_test

import (
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/org/cursor-sim/internal/generator"
)

func TestGenerateDevelopers_Count(t *testing.T) {
    // Arrange: Set up test parameters
    count := 50
    seed := int64(12345)
    
    // Act: Execute the function under test
    developers := generator.GenerateDevelopers(count, seed)
    
    // Assert: Verify the results
    assert.Len(t, developers, count)
}

func TestGenerateDevelopers_UniqueIDs(t *testing.T) {
    developers := generator.GenerateDevelopers(100, 12345)
    
    seen := make(map[string]bool)
    for _, dev := range developers {
        require.False(t, seen[dev.ID], "Duplicate ID found: %s", dev.ID)
        seen[dev.ID] = true
    }
}

func TestGenerateDevelopers_SeniorityDistribution(t *testing.T) {
    developers := generator.GenerateDevelopers(1000, 12345)
    
    counts := make(map[string]int)
    for _, dev := range developers {
        counts[dev.Seniority]++
    }
    
    // Allow 10% variance from target distribution
    juniorPct := float64(counts["junior"]) / 1000.0
    assert.InDelta(t, 0.20, juniorPct, 0.10, "Junior should be ~20%")
    
    midPct := float64(counts["mid"]) / 1000.0
    assert.InDelta(t, 0.50, midPct, 0.10, "Mid should be ~50%")
    
    seniorPct := float64(counts["senior"]) / 1000.0
    assert.InDelta(t, 0.30, seniorPct, 0.10, "Senior should be ~30%")
}
```

**Running Tests:**

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific test
go test -run TestGenerateDevelopers ./internal/generator/

# Run tests in parallel
go test -parallel 4 ./...

# Run benchmarks
go test -bench=. ./...
```

### Service B: cursor-analytics-core (TypeScript)

The TypeScript testing stack combines Jest for test execution with specialized libraries for database and GraphQL testing.

**Core Tools:**

**Jest** serves as the primary test runner, providing a rich set of features including test discovery, parallel execution, mocking utilities, and snapshot testing. The `ts-jest` transformer enables running TypeScript tests directly without a separate compilation step.

**Supertest** (`supertest`) simplifies HTTP integration testing by providing a fluent API for making requests and asserting on responses. It works well with Express applications, allowing tests to exercise the full HTTP stack including middleware.

**GraphQL Tools Mock** (`@graphql-tools/mock`) creates mock GraphQL schemas automatically, enabling frontend tests to run against realistic mock data without a running backend.

**pg-mem** (`pg-mem`) provides an in-memory PostgreSQL implementation for testing database operations without a real database server. This speeds up tests significantly while maintaining SQL compatibility.

**Example Test Structure:**

```typescript
// src/services/metrics.test.ts
import { MetricsService } from './metrics';
import { createTestDatabase, seedTestData } from '../test/helpers';

describe('MetricsService', () => {
    let metricsService: MetricsService;
    let db: TestDatabase;
    
    beforeAll(async () => {
        db = await createTestDatabase();
        metricsService = new MetricsService(db);
    });
    
    beforeEach(async () => {
        await db.clear();
    });
    
    afterAll(async () => {
        await db.close();
    });
    
    describe('calculateAcceptanceRate', () => {
        it('should calculate rate as (accepted/shown)*100', async () => {
            // Arrange: Seed test data with known values
            await seedTestData(db, {
                developers: [{ id: 'dev-1', name: 'Test Developer' }],
                events: [
                    { developerId: 'dev-1', type: 'cpp_suggestion_shown' },
                    { developerId: 'dev-1', type: 'cpp_suggestion_shown' },
                    { developerId: 'dev-1', type: 'cpp_suggestion_shown' },
                    { developerId: 'dev-1', type: 'cpp_suggestion_shown' },
                    { developerId: 'dev-1', type: 'cpp_suggestion_accepted' },
                    { developerId: 'dev-1', type: 'cpp_suggestion_accepted' },
                    { developerId: 'dev-1', type: 'cpp_suggestion_accepted' },
                ]
            });
            
            // Act
            const rate = await metricsService.calculateAcceptanceRate('dev-1');
            
            // Assert
            expect(rate).toBe(75.0); // 3 accepted / 4 shown = 75%
        });
        
        it('should return null when no suggestions shown', async () => {
            await seedTestData(db, {
                developers: [{ id: 'dev-1', name: 'Test Developer' }],
                events: [] // No events
            });
            
            const rate = await metricsService.calculateAcceptanceRate('dev-1');
            
            expect(rate).toBeNull();
        });
        
        it('should respect date range boundaries', async () => {
            const lastWeek = new Date('2026-01-08');
            const thisWeek = new Date('2026-01-15');
            
            await seedTestData(db, {
                developers: [{ id: 'dev-1', name: 'Test Developer' }],
                events: [
                    // Last week events (should be excluded)
                    { developerId: 'dev-1', type: 'cpp_suggestion_shown', timestamp: lastWeek },
                    { developerId: 'dev-1', type: 'cpp_suggestion_accepted', timestamp: lastWeek },
                    // This week events (should be included)
                    { developerId: 'dev-1', type: 'cpp_suggestion_shown', timestamp: thisWeek },
                    { developerId: 'dev-1', type: 'cpp_suggestion_shown', timestamp: thisWeek },
                ]
            });
            
            const rate = await metricsService.calculateAcceptanceRate('dev-1', {
                from: new Date('2026-01-12'),
                to: new Date('2026-01-19')
            });
            
            expect(rate).toBe(0); // 0 accepted / 2 shown = 0%
        });
    });
});
```

**Running Tests:**

```bash
# Run all tests
npm test

# Run tests in watch mode
npm test -- --watch

# Run with coverage
npm test -- --coverage

# Run specific test file
npm test -- src/services/metrics.test.ts

# Run tests matching pattern
npm test -- --testNamePattern="acceptance"

# Update snapshots
npm test -- -u
```

### Service C: cursor-viz-spa (React)

The React testing stack emphasizes testing user interactions and component behavior from the user's perspective.

**Core Tools:**

**Vitest** provides a Jest-compatible test runner optimized for Vite projects. It offers faster startup times than Jest through native ES modules support and integrates seamlessly with Vite's configuration.

**React Testing Library** (`@testing-library/react`) encourages testing components through their public interface, meaning through DOM interactions rather than internal implementation details. Tests render components and interact with them the way users would.

**MSW** (Mock Service Worker) intercepts HTTP requests at the network level, allowing tests to run against mock API responses without modifying application code. This is particularly valuable for testing loading states, error handling, and data transformations.

**Example Test Structure:**

```tsx
// src/components/charts/VelocityHeatmap.test.tsx
import { render, screen, fireEvent } from '@testing-library/react';
import { VelocityHeatmap } from './VelocityHeatmap';
import { generateTestDailyStats } from '../../test/factories';

describe('VelocityHeatmap', () => {
    const mockData = generateTestDailyStats({
        weeks: 52,
        baseValue: 50,
        variance: 20
    });
    
    it('should render correct number of cells', () => {
        render(<VelocityHeatmap data={mockData} />);
        
        // 52 weeks × 7 days = 364 cells
        const cells = screen.getAllByRole('cell');
        expect(cells).toHaveLength(364);
    });
    
    it('should apply color intensity based on value', () => {
        const dataWithKnownValues = [
            { date: '2026-01-01', acceptedCount: 0 },   // Should be lightest
            { date: '2026-01-02', acceptedCount: 50 },  // Should be medium
            { date: '2026-01-03', acceptedCount: 100 }, // Should be darkest
        ];
        
        render(<VelocityHeatmap data={dataWithKnownValues} />);
        
        const cells = screen.getAllByRole('cell');
        
        // Check that CSS classes or styles reflect intensity
        expect(cells[0]).toHaveClass('intensity-0');
        expect(cells[1]).toHaveClass('intensity-2');
        expect(cells[2]).toHaveClass('intensity-4');
    });
    
    it('should display tooltip on hover', async () => {
        render(<VelocityHeatmap data={mockData} />);
        
        const firstCell = screen.getAllByRole('cell')[0];
        fireEvent.mouseEnter(firstCell);
        
        // Tooltip should appear with date and count
        expect(await screen.findByRole('tooltip')).toBeInTheDocument();
        expect(screen.getByText(/January 1, 2026/)).toBeInTheDocument();
    });
    
    it('should show day-of-week labels', () => {
        render(<VelocityHeatmap data={mockData} />);
        
        expect(screen.getByText('Mon')).toBeInTheDocument();
        expect(screen.getByText('Wed')).toBeInTheDocument();
        expect(screen.getByText('Fri')).toBeInTheDocument();
    });
});
```

**Running Tests:**

```bash
# Run all tests
npm test

# Run tests in watch mode
npm test -- --watch

# Run with coverage
npm test -- --coverage

# Run with UI
npm test -- --ui

# Run specific test file
npm test -- VelocityHeatmap
```

## TDD Workflow

The Test-Driven Development workflow follows a specific cycle that ensures tests drive the implementation.

### The Red-Green-Refactor Cycle

The TDD cycle consists of three phases that repeat for each piece of functionality.

**Red Phase:** Write a failing test that describes the expected behavior. The test should fail because the functionality doesn't exist yet. This phase confirms that the test is actually testing something meaningful—a test that passes without implementation code may be testing the wrong thing.

**Green Phase:** Write the minimum amount of code necessary to make the test pass. The goal is not to write perfect code but to satisfy the test's requirements. Resist the urge to add additional functionality not covered by tests.

**Refactor Phase:** Improve the code while keeping tests green. This might involve extracting functions, renaming variables for clarity, or restructuring the code. The tests act as a safety net, ensuring that refactoring doesn't change behavior.

### Applying TDD to This Project

When implementing a feature from the specifications in this project, the workflow should proceed as follows.

**Step 1: Read the Specification**

Before writing any code, read the relevant specification document. For a feature like the acceptance rate calculation, consult `docs/USER_STORIES.md` to understand the user story US-CORE-003 with its acceptance criteria.

**Step 2: Derive Test Cases**

Transform each acceptance criterion into one or more test cases. The acceptance criteria use Given-When-Then format which maps directly to test structure. "Given" becomes the test setup (Arrange), "When" becomes the action (Act), and "Then" becomes the assertions (Assert).

**Step 3: Write Failing Tests**

Create test files and write test functions for each derived test case. Run the tests to confirm they fail for the expected reason (usually because the function or class under test doesn't exist yet).

**Step 4: Implement Minimally**

Write just enough code to make each test pass. Work through tests one at a time, running the test suite after each change to verify progress.

**Step 5: Refactor**

Once all tests pass, review the implementation for opportunities to improve code quality. Extract duplicated logic, improve naming, and ensure the code is readable. Run tests after each refactoring to catch any regressions.

**Step 6: Update Documentation**

If the implementation revealed any gaps or changes needed in the specification, update the relevant documentation. Specifications should remain accurate reflections of actual behavior.

### Example TDD Session

This example demonstrates implementing the acceptance rate calculation.

```typescript
// Step 3: Write failing test
// src/services/metrics.test.ts

describe('MetricsService', () => {
    describe('calculateAcceptanceRate', () => {
        it('should calculate rate as (accepted/shown)*100', async () => {
            // This test will fail because MetricsService doesn't exist
            const service = new MetricsService(db);
            const rate = await service.calculateAcceptanceRate('dev-1');
            expect(rate).toBe(75.0);
        });
    });
});
```

Running this test produces an error because `MetricsService` is not defined.

```typescript
// Step 4: Implement minimally
// src/services/metrics.ts

export class MetricsService {
    constructor(private db: Database) {}
    
    async calculateAcceptanceRate(developerId: string): Promise<number | null> {
        const result = await this.db.query(`
            SELECT 
                COUNT(*) FILTER (WHERE event_type = 'cpp_suggestion_shown') as shown,
                COUNT(*) FILTER (WHERE event_type = 'cpp_suggestion_accepted') as accepted
            FROM usage_events
            WHERE developer_id = $1
        `, [developerId]);
        
        if (result.shown === 0) {
            return null;
        }
        
        return (result.accepted / result.shown) * 100;
    }
}
```

The test now passes. Continue with additional tests for edge cases and date range filtering.

## Coverage Requirements

The project maintains minimum coverage thresholds to ensure adequate test coverage.

**Overall Coverage Target:** 80% line coverage across all services.

**Per-Service Targets:**

| Service | Line Coverage | Branch Coverage | Function Coverage |
|---------|---------------|-----------------|-------------------|
| cursor-sim | 80% | 70% | 90% |
| cursor-analytics-core | 80% | 75% | 90% |
| cursor-viz-spa | 70% | 65% | 85% |

The frontend has slightly lower targets because UI components may have visual behavior that is difficult to test automatically.

**Critical Path Targets:**

Certain code paths require higher coverage due to their importance. The metrics calculation functions must maintain 95% coverage. The data ingestion worker must maintain 90% coverage. GraphQL resolvers must maintain 85% coverage.

**Enforcement:**

Coverage thresholds are enforced in CI pipelines. Pull requests that decrease coverage below thresholds are blocked from merging. Coverage reports are generated for every PR and published as artifacts.

## Integration Testing

Integration tests verify that components work together correctly.

### API Integration Tests

Integration tests for the REST API (cursor-sim) and GraphQL API (cursor-analytics-core) exercise the full HTTP stack.

```go
// Integration test for cursor-sim API
func TestAPI_Integration(t *testing.T) {
    // Start the server in test mode
    srv := NewTestServer(t)
    defer srv.Close()
    
    // Test the full flow
    t.Run("create and query developers", func(t *testing.T) {
        // GET /v1/org/users should return developers
        resp, err := http.Get(srv.URL + "/v1/org/users")
        require.NoError(t, err)
        defer resp.Body.Close()
        
        assert.Equal(t, http.StatusOK, resp.StatusCode)
        
        var developers []Developer
        err = json.NewDecoder(resp.Body).Decode(&developers)
        require.NoError(t, err)
        assert.NotEmpty(t, developers)
    });
}
```

### Database Integration Tests

Database tests verify that queries and migrations work correctly with PostgreSQL.

```typescript
// Integration test for database operations
describe('Database Integration', () => {
    let db: Database;
    
    beforeAll(async () => {
        db = await createTestDatabase();
        await runMigrations(db);
    });
    
    it('should store and retrieve events', async () => {
        const event = {
            developerId: 'dev-1',
            eventType: 'cpp_suggestion_shown',
            timestamp: new Date()
        };
        
        await db.events.create(event);
        const retrieved = await db.events.findById(event.id);
        
        expect(retrieved).toMatchObject(event);
    });
});
```

## End-to-End Testing

End-to-end tests validate complete user workflows across all three services.

**Tool:** Playwright provides cross-browser testing with excellent support for modern web applications.

**Scope:** E2E tests focus on critical user journeys rather than exhaustive feature coverage.

```typescript
// e2e/dashboard.spec.ts
import { test, expect } from '@playwright/test';

test.describe('Dashboard', () => {
    test('should display data after full system startup', async ({ page }) => {
        // Navigate to dashboard
        await page.goto('http://localhost:3000');
        
        // Wait for loading to complete
        await expect(page.locator('.loading-skeleton')).toBeHidden();
        
        // Verify KPIs are displayed
        await expect(page.locator('[data-testid="total-developers"]')).toContainText(/\d+/);
        await expect(page.locator('[data-testid="acceptance-rate"]')).toContainText(/%/);
        
        // Verify charts are rendered
        await expect(page.locator('[data-testid="velocity-heatmap"]')).toBeVisible();
        await expect(page.locator('[data-testid="developer-table"]')).toBeVisible();
    });
    
    test('should filter data by date range', async ({ page }) => {
        await page.goto('http://localhost:3000');
        
        // Select "Last 7 Days" preset
        await page.click('[data-testid="date-range-picker"]');
        await page.click('text=Last 7 Days');
        
        // Verify data updates (charts should re-render)
        await expect(page.locator('[data-testid="loading-indicator"]')).toBeVisible();
        await expect(page.locator('[data-testid="loading-indicator"]')).toBeHidden();
    });
});
```

## CI/CD Integration

Tests are integrated into the continuous integration pipeline.

**Pipeline Stages:**

The CI pipeline runs in stages to provide fast feedback for common issues before running longer tests.

The lint stage runs first, checking code formatting and style rules. This catches simple issues quickly.

The unit test stage runs all unit tests in parallel across services. Tests run with coverage collection enabled.

The integration test stage spins up Docker containers for dependent services and runs integration tests.

The E2E test stage (on main branch only) runs the full system and executes Playwright tests.

**Configuration:**

```yaml
# .github/workflows/test.yml
name: Test

on: [push, pull_request]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Lint cursor-sim
        run: cd services/cursor-sim && golangci-lint run
      - name: Lint cursor-analytics-core
        run: cd services/cursor-analytics-core && npm run lint
      - name: Lint cursor-viz-spa
        run: cd services/cursor-viz-spa && npm run lint

  unit-test:
    runs-on: ubuntu-latest
    needs: lint
    strategy:
      matrix:
        service: [cursor-sim, cursor-analytics-core, cursor-viz-spa]
    steps:
      - uses: actions/checkout@v4
      - name: Run tests
        run: cd services/${{ matrix.service }} && make test-coverage
      - name: Upload coverage
        uses: codecov/codecov-action@v3

  integration-test:
    runs-on: ubuntu-latest
    needs: unit-test
    steps:
      - uses: actions/checkout@v4
      - name: Start services
        run: docker-compose up -d
      - name: Wait for healthy
        run: ./scripts/wait-for-healthy.sh
      - name: Run integration tests
        run: make test-integration
```

## Test Utilities and Helpers

The project provides shared testing utilities to reduce boilerplate.

**Factories:** Factory functions create test data with sensible defaults.

```typescript
// test/factories.ts
export function createTestDeveloper(overrides?: Partial<Developer>): Developer {
    return {
        id: `dev-${Math.random().toString(36).substr(2, 9)}`,
        name: 'Test Developer',
        email: 'test@example.com',
        team: 'Team Alpha',
        seniority: 'mid',
        ...overrides
    };
}

export function createTestEvent(overrides?: Partial<UsageEvent>): UsageEvent {
    return {
        id: `evt-${Math.random().toString(36).substr(2, 9)}`,
        developerId: 'dev-1',
        eventType: 'cpp_suggestion_shown',
        timestamp: new Date(),
        metadata: {},
        ...overrides
    };
}
```

**Database Helpers:** Utilities for setting up and tearing down test databases.

```typescript
// test/database.ts
export async function createTestDatabase(): Promise<TestDatabase> {
    const db = newDb();
    db.public.registerFunction({
        name: 'gen_random_uuid',
        implementation: () => crypto.randomUUID()
    });
    return db;
}

export async function seedTestData(db: TestDatabase, data: TestData): Promise<void> {
    for (const developer of data.developers || []) {
        await db.developers.create(developer);
    }
    for (const event of data.events || []) {
        await db.events.create(event);
    }
}
```

**Mock Builders:** Utilities for creating mock GraphQL responses.

```typescript
// test/mocks.ts
export function mockDashboardSummary(overrides?: Partial<DashboardKPI>): DashboardKPI {
    return {
        totalDevelopers: 50,
        activeDevelopers: 45,
        overallAcceptanceRate: 72.5,
        totalSuggestionsToday: 1250,
        totalAcceptedToday: 906,
        teamComparison: [],
        dailyTrend: [],
        ...overrides
    };
}
```
