# cursor-viz-spa

React-based visualization dashboard for Cursor AI analytics.

## Overview

Interactive web dashboard that consumes the cursor-analytics-core GraphQL API to display AI coding assistant usage analytics, team performance metrics, and code quality insights.

## Tech Stack

- **Framework**: React 18+ with TypeScript
- **Build Tool**: Vite
- **GraphQL Client**: Apollo Client
- **Charting**: Recharts
- **Styling**: Tailwind CSS
- **Testing**: Vitest + Testing Library
- **Port**: 3000

## Getting Started

### Prerequisites

- Node.js 18+
- npm or yarn
- cursor-analytics-core running on port 4000

### Installation

```bash
npm install
```

### Development

```bash
# Start dev server
npm run dev

# Run tests
npm run test

# Run tests with UI
npm run test:ui

# Run tests with coverage
npm run test:coverage

# Type check
npm run type-check

# Lint code
npm run lint

# Format code
npm run format
```

### Build

```bash
npm run build
npm run preview
```

## Project Structure

```
src/
├── components/
│   ├── layout/       # Header, Sidebar, etc.
│   ├── charts/       # Visualization components
│   ├── filters/      # Filter controls
│   └── common/       # Shared components
├── pages/            # Route pages
├── hooks/            # Custom React hooks
├── graphql/          # GraphQL queries & client
├── utils/            # Utility functions
├── test/             # Test setup & mocks
└── __tests__/        # Test files
```

## Configuration

Create a `.env` file (see `.env.example`):

```bash
VITE_GRAPHQL_URL=http://localhost:4000/graphql
VITE_APP_TITLE=Cursor Analytics Dashboard
VITE_DEFAULT_DATE_RANGE=LAST_30_DAYS
```

## GraphQL Code Generation

Generate TypeScript types from GraphQL schema:

```bash
npm run codegen
```

This requires cursor-analytics-core to be running on port 4000.

## Testing

- Unit tests: Component and hook testing
- Integration tests: Full page rendering with mocked API
- Coverage target: 80%

## Documentation

- **Specification**: `SPEC.md`
- **User Story**: `.work-items/P6-cursor-viz-spa/user-story.md`
- **Design**: `.work-items/P6-cursor-viz-spa/design.md`
- **Tasks**: `.work-items/P6-cursor-viz-spa/task.md`

## Dependencies

See `SPEC.md` for full dependency rationale and configuration details.
