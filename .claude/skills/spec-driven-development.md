# Spec-Driven Development Skill

This skill enables Claude to follow spec-driven development practices effectively within the Cursor Analytics Platform project.

## Overview

Spec-driven development is an approach where specifications and documentation drive the implementation process. Tests are derived from specifications, and code is written to satisfy those tests. This skill provides Claude with the knowledge and procedures to follow this methodology.

## Core Principles

### Specifications Come First

Every feature in this project has a specification that defines its behavior. Before implementing anything, Claude must locate and understand the relevant specification. Specifications live in several places within the project structure.

The service-level specifications in `services/{service}/SPEC.md` provide detailed technical requirements for each microservice. These include data models, API contracts, configuration options, and performance requirements.

The user stories in `docs/USER_STORIES.md` describe features from the user's perspective with Given-When-Then acceptance criteria. These criteria map directly to test cases.

The task breakdown in `docs/TASKS.md` lists specific implementation tasks with their dependencies and definitions of done.

### Tests Before Code

Following Test-Driven Development principles, tests must exist and fail before implementation code is written. The workflow proceeds in three phases that repeat for each piece of functionality.

During the Red phase, Claude writes a test based on the specification's acceptance criteria. The test should fail because the functionality doesn't exist yet. This confirms the test is actually testing something meaningful.

During the Green phase, Claude writes the minimum code necessary to make the test pass. The goal is satisfying the test's requirements, not writing perfect code. Additional functionality not covered by tests should be avoided.

During the Refactor phase, Claude improves the code while keeping tests green. This might involve extracting functions, improving names, or restructuring. The tests act as a safety net ensuring refactoring doesn't change behavior.

### Documentation Stays Current

When implementation reveals gaps or changes needed in specifications, the documentation must be updated. Specifications should remain accurate reflections of actual behavior. This bidirectional relationship between specs and code ensures the documentation never becomes stale.

## Procedures

### Starting a New Feature

When asked to implement a feature, Claude should follow this procedure.

First, locate the specification by checking the service SPEC.md, searching USER_STORIES.md for relevant stories, and finding the task in TASKS.md. Read all relevant documentation before proceeding.

Second, verify dependencies by checking that prerequisite tasks are complete. If dependencies are not met, either complete them first or inform the user that the feature is blocked.

Third, derive test cases from acceptance criteria. Each Given-When-Then scenario becomes at least one test case. Create the test file if it doesn't exist.

Fourth, write failing tests that encode the expected behavior. Run the tests to confirm they fail for the expected reason, typically because the function or class under test doesn't exist.

Fifth, implement minimally by writing just enough code to make each test pass. Work through tests one at a time, running the suite after each change.

Sixth, refactor while tests are green. Clean up the code, extract duplications, and improve readability. Run tests after each refactoring step.

Seventh, update documentation if the implementation revealed any specification gaps or if behavior differs from what was documented.

### Locating Relevant Specifications

When working on a feature, Claude should use these strategies to find specifications.

For simulator features, check `services/cursor-sim/SPEC.md` for technical details about data models, API endpoints, and generation logic. Search `docs/USER_STORIES.md` for stories starting with "US-SIM-".

For aggregator features, check `services/cursor-analytics-core/SPEC.md` for database schema, GraphQL types, and metric calculations. Search for stories starting with "US-CORE-".

For dashboard features, check `services/cursor-viz-spa/SPEC.md` for component specifications, props interfaces, and behavior descriptions. Search for stories starting with "US-VIZ-".

For cross-cutting features, check `docs/DESIGN.md` for architectural decisions and `docs/TESTING_STRATEGY.md` for testing approaches.

### Writing Tests from Specifications

Acceptance criteria in USER_STORIES.md use Given-When-Then format which maps directly to test structure.

The Given clause describes the test setup or preconditions. This becomes the Arrange section of the test where test data is created and dependencies are configured.

The When clause describes the action being tested. This becomes the Act section where the function under test is called or the interaction is performed.

The Then clause describes the expected outcome. This becomes the Assert section where results are verified against expectations.

Here is an example transformation from acceptance criteria to test code.

The acceptance criterion states: "Given a developer has 100 cpp_suggestion_shown events and 75 cpp_suggestion_accepted events, when I query their acceptance rate, then the result should be 75.0."

This becomes the following test structure:

```typescript
it('should calculate rate as (accepted/shown)*100', async () => {
    // Arrange (Given)
    await seedTestData(db, {
        developers: [{ id: 'dev-1' }],
        events: [
            ...Array(100).fill({ developerId: 'dev-1', type: 'cpp_suggestion_shown' }),
            ...Array(75).fill({ developerId: 'dev-1', type: 'cpp_suggestion_accepted' }),
        ]
    });
    
    // Act (When)
    const rate = await metricsService.calculateAcceptanceRate('dev-1');
    
    // Assert (Then)
    expect(rate).toBe(75.0);
});
```

### Handling Specification Gaps

Sometimes implementation reveals that specifications are incomplete or ambiguous. When this happens, Claude should make a reasonable assumption based on the specification's intent, document the assumption in a code comment, implement based on that assumption, and flag the gap by suggesting a specification update.

The suggested update should be specific and actionable. For example: "The SPEC.md should clarify whether acceptance rate returns 0 or null when a developer has no suggestions. I implemented null to avoid implying certainty about a rate that cannot be calculated."

## Quality Standards

### Test Coverage

This project requires 80% line coverage for all services. Critical paths like metric calculations require 95% coverage. Claude should write tests that exercise both happy paths and edge cases to meet these thresholds.

### Code Style

Each service has specific style requirements. Go code in cursor-sim must pass golangci-lint. TypeScript code must compile with strict mode enabled. React components should follow the functional component pattern with hooks.

### Documentation

Code comments should explain why, not what. The code itself should be readable enough that what it does is clear. Comments add value by explaining the reasoning behind decisions, especially when the code handles edge cases or implements business rules from the specification.
