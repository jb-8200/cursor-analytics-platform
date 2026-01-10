# Testing Strategy: Cursor Usage Analytics Platform

**Version**: 1.0.0  
**Last Updated**: January 2026  

This document defines the testing approach for the Cursor Usage Analytics Platform, including the recommended testing stack, TDD workflows, and coverage requirements. The strategy is designed to support spec-driven development where tests are written from specifications before implementation code.

## Testing Philosophy

The project follows Test-Driven Development (TDD) as its primary development methodology. This means tests are written before implementation code, serving as executable specifications that define the expected behavior of each component. The testing pyramid guides our approach: many unit tests form the foundation, integration tests verify component interactions, and a smaller number of end-to-end tests validate complete user workflows.

Writing tests first provides several benefits that improve the overall quality of the codebase. Tests act as documentation that never becomes stale since failing tests indicate incorrect behavior immediately. The TDD cycle of red-green-refactor encourages simple, focused implementations that do exactly what the specification requires. Test-first development also leads to more testable designs because developers must consider how code will be tested before writing it.

### Data Contract and Source of Truth

**cursor-sim (P4) is the authoritative source of truth** for the analytics platform. All data flowing through the system must:

1. **Validate against cursor-sim's API contract** defined in `services/cursor-sim/SPEC.md`
2. **Preserve data fidelity** as it moves through the pipeline: cursor-sim → Data Tier → Dashboard
3. **Establish contracts at each layer**:
   - **API Contract**: cursor-sim response format and fields (source of truth)
   - **Data Tier Contract**: dbt transformations map API fields to analytics columns
   - **Dashboard Contract**: Streamlit queries consume only mart tables, never raw data

**Data contract hierarchy:**
```
cursor-sim API (camelCase, {items:[]})
    ↓ [api-loader validates against SPEC.md]
DuckDB raw schema (raw_commits, raw_pull_requests)
    ↓ [dbt staging transforms camelCase → snake_case]
DuckDB staging schema (stg_commits, stg_pull_requests)
    ↓ [dbt marts aggregate for analytics]
DuckDB mart schema (mart_velocity, mart_ai_impact, mart_quality, mart_review_costs)
    ↓ [dashboard queries consume marts only]
Streamlit Dashboard (KPI visualizations)
```

**Testing implications**:
- Data tier tests must verify: API format handling → DuckDB schema alignment → dbt transformation correctness
- Dashboard tests must verify: parameterized queries → no SQL injection → correct mart column usage
- E2E tests must trace: API response → DuckDB table → dbt transform → mart result → dashboard display

## Data Pipeline Testing Strategy

The data pipeline (P8-F01 data tier + P9-F01 dashboard) requires specialized testing to validate data contracts and transformations end-to-end.

### P8: Data Tier Testing (api-loader → dbt → DuckDB)

**Contract Validation Tests**: Verify each stage preserves data integrity according to contracts.

```python
# tools/api-loader/tests/test_api_contract.py
# Test 1: Verify API response format matches cursor-sim SPEC.md

def test_api_response_format_dual_support():
    """Cursor-sim may return {items:[...]} or raw array [...] - both must be handled."""
    extractor = BaseAPIExtractor()

    # Format 1: Paginated response (production format)
    paginated = {
        "items": [{"commitHash": "abc123", "userEmail": "dev@example.com"}],
        "totalCount": 100,
        "page": 1,
        "pageSize": 50
    }
    result1 = extractor.fetch_cursor_style_paginated('/commits', paginated)
    assert len(result1) == 1
    assert result1[0]['commitHash'] == 'abc123'

    # Format 2: Raw array response (fallback format)
    raw_array = [{"commitHash": "abc123", "userEmail": "dev@example.com"}]
    result2 = extractor.fetch_cursor_style_paginated('/commits', raw_array)
    assert len(result2) == 1
    assert result2[0]['commitHash'] == 'abc123'


def test_api_column_mapping_contract():
    """Verify API response fields match cursor-sim SPEC.md field names."""
    # Contract: cursor-sim returns camelCase fields
    api_response = {
        "items": [{
            "commitHash": "abc123",
            "userEmail": "dev@example.com",
            "tabLinesAdded": 45,
            "composerLinesAdded": 12,
            "commitTs": "2026-01-10T10:30:00Z"
        }]
    }

    # api-loader must not transform API response
    # (transformation happens in dbt staging layer)
    extractor = BaseAPIExtractor()
    result = extractor.fetch_cursor_style_paginated('/commits', api_response)

    # Verify API fields are preserved as-is
    commit = result[0]
    assert 'commitHash' in commit, "API camelCase fields must be preserved"
    assert 'userEmail' in commit, "API camelCase fields must be preserved"
    assert 'tabLinesAdded' in commit, "API camelCase fields must be preserved"


def test_duckdb_raw_schema_loading():
    """Verify API data loads into raw_commits without transformation."""
    api_data = [{
        "commitHash": "abc123",
        "userEmail": "dev@example.com",
        "tabLinesAdded": 45,
        "composerLinesAdded": 12,
        "commitTs": "2026-01-10T10:30:00Z"
    }]

    loader = DuckDBLoader(db_path=":memory:")
    loader.load_raw_commits(api_data)

    # Raw table should have API fields as-is
    df = loader.conn.execute("SELECT * FROM raw_commits LIMIT 1").df()
    assert 'commitHash' in df.columns, "Raw schema should preserve API camelCase"
    assert df.iloc[0]['commitHash'] == 'abc123'


def test_dbt_staging_column_mapping():
    """Verify dbt staging layer maps camelCase → snake_case correctly."""
    # This test assumes dbt has been run and staging tables exist
    db = DuckDBLoader(db_path=":memory:")

    # Query the staging layer (dbt-transformed data)
    df = db.conn.execute("""
        SELECT
            commit_hash,
            user_email,
            tab_lines_added,
            composer_lines_added,
            committed_at
        FROM main_staging.stg_commits
        LIMIT 1
    """).df()

    # Verify snake_case transformation
    expected_columns = {
        'commit_hash': 'commitHash',
        'user_email': 'userEmail',
        'tab_lines_added': 'tabLinesAdded',
        'composer_lines_added': 'composerLinesAdded',
        'committed_at': 'commitTs'
    }

    for snake_col, camel_col in expected_columns.items():
        assert snake_col in df.columns, f"Staging table must have {snake_col}"


def test_dbt_mart_aggregation():
    """Verify dbt mart layer produces correct aggregates for analytics."""
    db = DuckDBLoader(db_path=":memory:")

    # Query mart_velocity table (dbt mart output)
    df = db.conn.execute("""
        SELECT
            week,
            repo_name,
            active_developers,
            total_prs,
            avg_total_cycle_time,
            avg_ai_ratio
        FROM main_mart.mart_velocity
        WHERE week >= CURRENT_DATE - INTERVAL '4' WEEK
        ORDER BY week DESC
    """).df()

    # Verify all expected columns exist
    required_columns = {
        'week', 'repo_name', 'active_developers',
        'total_prs', 'avg_total_cycle_time', 'avg_ai_ratio'
    }
    assert required_columns.issubset(set(df.columns)), \
        f"Missing columns: {required_columns - set(df.columns)}"

    # Verify no null averages (indicates aggregation worked)
    assert df['avg_total_cycle_time'].notna().all(), \
        "avg_total_cycle_time should not be null after aggregation"
    assert df['avg_ai_ratio'].notna().all(), \
        "avg_ai_ratio should not be null after aggregation"
```

### P9: Dashboard Testing (parameterized queries, SQL injection prevention)

**Security and Parameterization Tests**: Verify dashboard is protected against SQL injection.

```python
# services/streamlit-dashboard/tests/test_query_security.py

def test_velocity_query_parameterized_injection_attempt():
    """Verify parameterized queries prevent SQL injection."""
    from queries.velocity import get_velocity_data

    # Malicious input that would work if concatenated, should fail safely with params
    malicious_repo = "test'; DROP TABLE main_mart.mart_velocity; --"

    # This must NOT execute the DROP TABLE command
    df = get_velocity_data(repo_name=malicious_repo, days=30)

    # If parameterized correctly, query should execute with malicious string as literal value
    # and return either empty dataframe or error (no dropped tables)
    assert isinstance(df, pd.DataFrame), "Query should return DataFrame, not error"

    # Verify table still exists by running another query
    from db.connector import query
    verify_df = query("SELECT COUNT(*) as table_count FROM main_mart.mart_velocity")
    assert len(verify_df) > 0, "Table should still exist after malicious input attempt"


def test_schema_naming_requires_main_prefix():
    """Verify queries use DuckDB required main_mart.* schema prefix."""
    from queries.velocity import get_velocity_data

    # This should work (correct schema naming)
    df = get_velocity_data(repo_name="test-repo", days=30)
    assert isinstance(df, pd.DataFrame), "Should succeed with main_mart.* schema"

    # Verify error if schema name is wrong by checking internal query
    # (dashboard should never construct mart.* queries)
    from db.connector import query

    try:
        # This should fail - incorrect schema name
        bad_df = query("SELECT * FROM mart.mart_velocity LIMIT 1")
        assert False, "Should have raised error with incorrect mart.* schema"
    except Exception as e:
        assert "Catalog Error" in str(e) or "not found" in str(e).lower(), \
            "Should fail with catalog/not found error for incorrect schema"


def test_interval_syntax_with_days_parameter():
    """Verify INTERVAL syntax handles days parameter safely."""
    from queries.velocity import get_velocity_data

    # Days is integer and validated, f-string interpolation is acceptable
    test_cases = [
        (7, "last 7 days"),
        (30, "last 30 days"),
        (90, "last 90 days"),
    ]

    for days, description in test_cases:
        df = get_velocity_data(days=days)
        assert isinstance(df, pd.DataFrame), f"Should handle {description}"

        # Verify only data from date range is returned
        if len(df) > 0:
            cutoff = pd.Timestamp.now(tz='UTC').normalize() - pd.Timedelta(days=days)
            df['week'] = pd.to_datetime(df['week'])
            assert (df['week'] >= cutoff).all(), \
                f"All rows should be >= {days} days ago"


def test_filter_parameters_isolation():
    """Verify sidebar filters don't leak into SQL."""
    from components.sidebar import get_filter_params
    from queries.velocity import get_velocity_data

    # Sidebar returns raw values for parameterization downstream
    repo_name, date_range, days = get_filter_params()

    # get_velocity_data must parameterize these values
    df = get_velocity_data(repo_name=repo_name, days=days)
    assert isinstance(df, pd.DataFrame), "Filtered query should succeed"
```

### Data Contract Validation Test Examples

**Validate end-to-end data flow** from API through to dashboard:

```bash
#!/bin/bash
# scripts/test-data-contract-e2e.sh
# Validates complete data pipeline: API → DuckDB raw → dbt staging/marts → Dashboard

set -e

echo "Data Contract E2E Validation"
echo "=============================="

# 1. Verify cursor-sim API returns expected format
echo "1. Checking cursor-sim API contract..."
response=$(curl -s http://localhost:8080/v1/org/users | head -c 200)
if [[ $response == *"items"* ]]; then
    echo "   ✓ API returns {items:[...]} format (matches SPEC.md)"
elif [[ $response == *"[{"* ]]; then
    echo "   ✓ API returns raw array format (fallback supported)"
else
    echo "   ✗ API response format unknown"
    exit 1
fi

# 2. Verify DuckDB raw_commits table has API fields
echo "2. Checking DuckDB raw schema..."
docker exec streamlit-dashboard duckdb /data/analytics.duckdb << 'SQL' > /tmp/raw_columns.txt
.columns raw_commits
SQL

if grep -q "commitHash" /tmp/raw_columns.txt; then
    echo "   ✓ Raw schema preserves camelCase API fields"
else
    echo "   ✗ Raw schema missing API fields"
    exit 1
fi

# 3. Verify dbt staging has snake_case transformation
echo "3. Checking dbt staging transformation..."
docker exec streamlit-dashboard duckdb /data/analytics.duckdb << 'SQL' > /tmp/stg_columns.txt
.columns main_staging.stg_commits
SQL

if grep -q "commit_hash" /tmp/stg_columns.txt && ! grep -q "commitHash" /tmp/stg_columns.txt; then
    echo "   ✓ Staging transforms camelCase → snake_case"
else
    echo "   ✗ Staging transformation failed"
    exit 1
fi

# 4. Verify mart table aggregations work
echo "4. Checking mart aggregations..."
docker exec streamlit-dashboard duckdb /data/analytics.duckdb << 'SQL' > /tmp/mart_check.txt
SELECT COUNT(*) as row_count,
       COUNT(avg_total_cycle_time) as avg_cycle_non_null,
       COUNT(avg_ai_ratio) as avg_ai_non_null
FROM main_mart.mart_velocity
LIMIT 1
SQL

rows=$(cat /tmp/mart_check.txt | tail -1)
if [[ $rows == *"0"* ]]; then
    echo "   ⚠ Mart table is empty (expected if data not yet loaded)"
else
    echo "   ✓ Mart aggregations populated"
fi

# 5. Verify dashboard queries use parameterized syntax
echo "5. Checking dashboard query security..."
grep -r "WHERE repo_name = \$repo" services/streamlit-dashboard/queries/ > /dev/null && \
    echo "   ✓ Dashboard uses parameterized queries" || \
    { echo "   ✗ Dashboard has non-parameterized queries"; exit 1; }

# 6. Test SQL injection prevention
echo "6. Testing SQL injection prevention..."
docker exec streamlit-dashboard python << 'PYTHON'
from db.connector import query
import pandas as pd

# This would DROP TABLE if query was concatenated (not parameterized)
malicious_input = "test'; DROP TABLE main_mart.mart_velocity; --"

try:
    df = query(
        "SELECT * FROM main_mart.mart_velocity WHERE repo_name = $repo LIMIT 1",
        {"repo": malicious_input}
    )
    print("   ✓ Malicious input safely parameterized")
except Exception as e:
    print(f"   ✗ Query failed: {e}")
    exit(1)
PYTHON

echo ""
echo "=============================="
echo "✓ All data contract validations passed!"
```

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

## Reference Documentation

### Source of Truth Hierarchy

When writing tests, always reference these documents in priority order:

| Priority | Document | Purpose |
|----------|----------|---------|
| **1** | `services/cursor-sim/SPEC.md` | API contract: endpoints, response formats, field names/types |
| **2** | `tools/api-loader/` design docs | Data extraction strategy: API response handling, format normalization |
| **3** | `dbt/models/` + `docs/` | Transformation contracts: staging layer (camelCase→snake_case), mart aggregations |
| **4** | `services/streamlit-dashboard/` design docs | Consumer contracts: parameterized queries, schema naming (main_mart.*), column availability |

### Key Testing Checkpoints

When validating data contracts across the pipeline:

**API Contract** (cursor-sim SPEC.md)
- Response format: `{items: [...], totalCount, page, pageSize}` or raw array `[...]`
- Field naming: **camelCase** (commitHash, userEmail, tabLinesAdded, composerLinesAdded, commitTs)
- Pagination: cursor-sim implements cursor-based pagination

**Data Tier Contract** (api-loader → dbt → DuckDB)
- Raw layer: Preserves API fields as-is (camelCase)
- Staging layer: Transforms to snake_case (commit_hash, user_email, tab_lines_added, etc.)
- Mart layer: Aggregations produce analytics columns (avg_total_cycle_time, avg_ai_ratio, revert_rate, etc.)

**Dashboard Contract** (Streamlit → DuckDB)
- Query syntax: Parameterized with `$param` placeholders (DuckDB requirement)
- Schema naming: `main_mart.mart_*` not `mart.*` (DuckDB prefix requirement)
- INTERVAL syntax: Use f-string for days parameter: `CURRENT_DATE - INTERVAL '{days}' DAY`
- SQL injection: All user inputs from sidebar → parameterized in queries

### Testing checklist

✓ Does my test verify against the source of truth (cursor-sim SPEC.md)?
✓ Am I testing the contract boundary (input → transformation → output)?
✓ Does my test handle both API response formats ({items:[]} and raw array)?
✓ Do my dashboard tests verify parameterized queries (no f-string SQL)?
✓ Am I using correct schema naming (main_mart.mart_* not mart.*)?
✓ Do I handle the INTERVAL syntax correctly (f-string for days, not parameter)?
```
