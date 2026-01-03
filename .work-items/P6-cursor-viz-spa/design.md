# Technical Design: cursor-viz-spa Dashboard

**Feature ID**: P6-cursor-viz-spa
**Phase**: P6 (Visualization SPA)
**Created**: January 3, 2026
**Status**: TODO

## Architecture

```
┌──────────────────────────────────────┐
│         cursor-viz-spa               │
│         (React + Vite)               │
│         Port: 3000                   │
│                                      │
│  ┌────────────────────────────────┐ │
│  │  React Components              │ │
│  │  - Dashboard Layout            │ │
│  │  - Chart Components            │ │
│  │  - Filter Controls             │ │
│  └────────────────────────────────┘ │
│                │                     │
│                ▼                     │
│  ┌────────────────────────────────┐ │
│  │  Apollo Client (GraphQL)       │ │
│  └────────────────────────────────┘ │
└──────────────────┬───────────────────┘
                   │ GraphQL Queries
                   ▼
    ┌──────────────────────────────┐
    │  cursor-analytics-core       │
    │  GraphQL API (Port: 4000)    │
    └──────────────────────────────┘
```

## Technology Stack

- **Framework**: React 18+ with TypeScript
- **Build Tool**: Vite
- **GraphQL Client**: Apollo Client
- **Charting**: Recharts or Chart.js
- **Styling**: Tailwind CSS
- **State Management**: React Context or Zustand

## Key Features

1. **AI Code Analytics Dashboard**
2. **Team Leaderboards**
3. **Model Usage Visualization**
4. **Time-Series Charts**
5. **Exportable Reports**

---

**Next Steps**: Define tasks and components
