---
name: react-vite-patterns
description: React and Vite best practices for cursor-viz-spa dashboard. Use when implementing React components, Apollo Client queries, Tailwind styling, or Vite configuration. Covers hooks, state management, and data visualization patterns.
---

# React & Vite Patterns

Best practices for cursor-viz-spa dashboard development.

## Vite Configuration

### Standard Setup

```typescript
// vite.config.ts
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import path from 'path';

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
      '@components': path.resolve(__dirname, './src/components'),
      '@hooks': path.resolve(__dirname, './src/hooks'),
      '@graphql': path.resolve(__dirname, './src/graphql'),
    },
  },
  server: {
    port: 3000,
    proxy: {
      '/graphql': 'http://localhost:4000',
    },
  },
});
```

## Component Patterns

### Functional Components with TypeScript

```tsx
interface DeveloperCardProps {
  developer: Developer;
  onSelect?: (id: string) => void;
  isSelected?: boolean;
}

export function DeveloperCard({
  developer,
  onSelect,
  isSelected = false
}: DeveloperCardProps) {
  const handleClick = () => {
    onSelect?.(developer.id);
  };

  return (
    <div
      className={cn(
        'rounded-lg border p-4 cursor-pointer transition-colors',
        isSelected ? 'border-blue-500 bg-blue-50' : 'border-gray-200 hover:border-gray-300'
      )}
      onClick={handleClick}
    >
      <h3 className="font-semibold">{developer.name}</h3>
      <p className="text-sm text-gray-600">{developer.email}</p>
      <div className="mt-2 flex items-center gap-2">
        <Badge variant={developer.seniority}>{developer.seniority}</Badge>
        <span className="text-sm">AI: {(developer.aiPreference * 100).toFixed(0)}%</span>
      </div>
    </div>
  );
}
```

### Chart Components

```tsx
import { LineChart, Line, XAxis, YAxis, Tooltip, ResponsiveContainer } from 'recharts';

interface TrendChartProps {
  data: TrendDataPoint[];
  title: string;
  loading?: boolean;
  error?: Error | null;
}

export function TrendChart({ data, title, loading, error }: TrendChartProps) {
  if (loading) {
    return <ChartSkeleton title={title} />;
  }

  if (error) {
    return <ChartError title={title} message={error.message} />;
  }

  if (data.length === 0) {
    return <ChartEmpty title={title} message="No data available" />;
  }

  return (
    <div className="rounded-lg border bg-white p-4">
      <h3 className="mb-4 font-semibold">{title}</h3>
      <ResponsiveContainer width="100%" height={300}>
        <LineChart data={data}>
          <XAxis dataKey="date" />
          <YAxis />
          <Tooltip />
          <Line
            type="monotone"
            dataKey="aiRatio"
            stroke="#3b82f6"
            strokeWidth={2}
          />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}
```

## Apollo Client Setup

### Client Configuration

```typescript
// src/graphql/client.ts
import { ApolloClient, InMemoryCache, HttpLink } from '@apollo/client';

const httpLink = new HttpLink({
  uri: import.meta.env.VITE_GRAPHQL_URL || '/graphql',
});

export const apolloClient = new ApolloClient({
  link: httpLink,
  cache: new InMemoryCache({
    typePolicies: {
      Query: {
        fields: {
          developers: {
            keyArgs: ['filter'],
            merge(existing, incoming, { args }) {
              if (!args?.pagination?.page || args.pagination.page === 1) {
                return incoming;
              }
              return {
                ...incoming,
                edges: [...(existing?.edges || []), ...incoming.edges],
              };
            },
          },
        },
      },
    },
  }),
  defaultOptions: {
    watchQuery: {
      fetchPolicy: 'cache-and-network',
    },
  },
});
```

### Typed Queries

```typescript
// src/graphql/queries/team.ts
import { gql } from '@apollo/client';

export const GET_TEAM_OVERVIEW = gql`
  query GetTeamOverview($dateRange: DateRangeInput!) {
    teamOverview(dateRange: $dateRange) {
      totalDevelopers
      activeDevelopers
      totalAiLines
      totalHumanLines
      averageAiRatio
    }
  }
`;

export const GET_AI_USAGE_TREND = gql`
  query GetAiUsageTrend($dateRange: DateRangeInput!, $granularity: Granularity!) {
    aiUsageTrend(dateRange: $dateRange, granularity: $granularity) {
      date
      aiLines
      humanLines
      aiRatio
    }
  }
`;
```

### Query Hooks

```typescript
// src/hooks/useTeamOverview.ts
import { useQuery } from '@apollo/client';
import { GET_TEAM_OVERVIEW } from '@graphql/queries/team';
import type { TeamOverview, DateRange } from '@/types';

interface UseTeamOverviewResult {
  data: TeamOverview | null;
  loading: boolean;
  error: Error | null;
  refetch: () => void;
}

export function useTeamOverview(dateRange: DateRange): UseTeamOverviewResult {
  const { data, loading, error, refetch } = useQuery(GET_TEAM_OVERVIEW, {
    variables: { dateRange },
    skip: !dateRange.startDate || !dateRange.endDate,
  });

  return {
    data: data?.teamOverview ?? null,
    loading,
    error: error ?? null,
    refetch,
  };
}
```

## Custom Hooks

### Filter State Hook

```typescript
// src/hooks/useFilters.ts
import { useState, useCallback, useMemo } from 'react';
import { useSearchParams } from 'react-router-dom';

interface FilterState {
  dateRange: DateRange;
  seniority: Seniority[];
  minAiRatio: number;
}

export function useFilters() {
  const [searchParams, setSearchParams] = useSearchParams();

  const filters = useMemo<FilterState>(() => ({
    dateRange: {
      startDate: searchParams.get('start') || getDefaultStartDate(),
      endDate: searchParams.get('end') || getDefaultEndDate(),
    },
    seniority: (searchParams.get('seniority')?.split(',') || []) as Seniority[],
    minAiRatio: Number(searchParams.get('minAi')) || 0,
  }), [searchParams]);

  const setFilters = useCallback((updates: Partial<FilterState>) => {
    setSearchParams((prev) => {
      const next = new URLSearchParams(prev);
      if (updates.dateRange) {
        next.set('start', updates.dateRange.startDate);
        next.set('end', updates.dateRange.endDate);
      }
      if (updates.seniority) {
        next.set('seniority', updates.seniority.join(','));
      }
      if (updates.minAiRatio !== undefined) {
        next.set('minAi', String(updates.minAiRatio));
      }
      return next;
    });
  }, [setSearchParams]);

  return { filters, setFilters };
}
```

### Debounced Input Hook

```typescript
// src/hooks/useDebounce.ts
import { useState, useEffect } from 'react';

export function useDebounce<T>(value: T, delay: number): T {
  const [debouncedValue, setDebouncedValue] = useState(value);

  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedValue(value);
    }, delay);

    return () => clearTimeout(timer);
  }, [value, delay]);

  return debouncedValue;
}
```

## Tailwind CSS Patterns

### Design Tokens

```javascript
// tailwind.config.js
module.exports = {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        // Chart colors for data visualization
        chart: {
          ai: '#3b82f6',      // Blue for AI code
          human: '#10b981',   // Green for human code
          mixed: '#8b5cf6',   // Purple for mixed
        },
        // Seniority badge colors
        seniority: {
          junior: '#fbbf24',
          mid: '#60a5fa',
          senior: '#34d399',
        },
      },
    },
  },
  plugins: [],
};
```

### Component Classes

```tsx
// Common pattern: cn() utility for conditional classes
import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

// Usage
<button
  className={cn(
    'px-4 py-2 rounded-md font-medium transition-colors',
    variant === 'primary' && 'bg-blue-600 text-white hover:bg-blue-700',
    variant === 'secondary' && 'bg-gray-100 text-gray-700 hover:bg-gray-200',
    disabled && 'opacity-50 cursor-not-allowed'
  )}
>
  {children}
</button>
```

## State Management

### Context for Global State

```tsx
// src/contexts/FilterContext.tsx
import { createContext, useContext, ReactNode } from 'react';
import { useFilters } from '@hooks/useFilters';

type FilterContextType = ReturnType<typeof useFilters>;

const FilterContext = createContext<FilterContextType | null>(null);

export function FilterProvider({ children }: { children: ReactNode }) {
  const value = useFilters();
  return (
    <FilterContext.Provider value={value}>
      {children}
    </FilterContext.Provider>
  );
}

export function useFilterContext() {
  const context = useContext(FilterContext);
  if (!context) {
    throw new Error('useFilterContext must be used within FilterProvider');
  }
  return context;
}
```

## Testing Patterns

### Component Tests with Testing Library

```tsx
import { render, screen, fireEvent } from '@testing-library/react';
import { MockedProvider } from '@apollo/client/testing';
import { DeveloperCard } from './DeveloperCard';

const mockDeveloper = {
  id: 'dev-1',
  name: 'Alice',
  email: 'alice@example.com',
  seniority: 'SENIOR',
  aiPreference: 0.75,
};

describe('DeveloperCard', () => {
  it('renders developer information', () => {
    render(<DeveloperCard developer={mockDeveloper} />);

    expect(screen.getByText('Alice')).toBeInTheDocument();
    expect(screen.getByText('alice@example.com')).toBeInTheDocument();
    expect(screen.getByText('SENIOR')).toBeInTheDocument();
  });

  it('calls onSelect when clicked', () => {
    const onSelect = jest.fn();
    render(<DeveloperCard developer={mockDeveloper} onSelect={onSelect} />);

    fireEvent.click(screen.getByText('Alice'));
    expect(onSelect).toHaveBeenCalledWith('dev-1');
  });
});
```

### Query Tests

```tsx
import { render, screen, waitFor } from '@testing-library/react';
import { MockedProvider } from '@apollo/client/testing';
import { GET_TEAM_OVERVIEW } from '@graphql/queries/team';
import { Dashboard } from './Dashboard';

const mocks = [
  {
    request: {
      query: GET_TEAM_OVERVIEW,
      variables: { dateRange: { startDate: '2025-01-01', endDate: '2025-01-31' } },
    },
    result: {
      data: {
        teamOverview: {
          totalDevelopers: 10,
          activeDevelopers: 8,
          totalAiLines: 5000,
          totalHumanLines: 3000,
          averageAiRatio: 0.625,
        },
      },
    },
  },
];

describe('Dashboard', () => {
  it('displays team overview data', async () => {
    render(
      <MockedProvider mocks={mocks} addTypename={false}>
        <Dashboard />
      </MockedProvider>
    );

    await waitFor(() => {
      expect(screen.getByText('10 developers')).toBeInTheDocument();
    });
  });
});
```

## Error Boundaries

```tsx
// src/components/common/ErrorBoundary.tsx
import { Component, ReactNode } from 'react';

interface Props {
  children: ReactNode;
  fallback?: ReactNode;
}

interface State {
  hasError: boolean;
  error: Error | null;
}

export class ErrorBoundary extends Component<Props, State> {
  state: State = { hasError: false, error: null };

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  render() {
    if (this.state.hasError) {
      return this.props.fallback || (
        <div className="rounded-lg border border-red-200 bg-red-50 p-4">
          <h2 className="font-semibold text-red-800">Something went wrong</h2>
          <p className="text-sm text-red-600">{this.state.error?.message}</p>
        </div>
      );
    }

    return this.props.children;
  }
}
```

## Performance

### Memoization

```tsx
import { memo, useMemo, useCallback } from 'react';

// Memoize expensive computations
const sortedDevelopers = useMemo(() => {
  return [...developers].sort((a, b) => b.aiRatio - a.aiRatio);
}, [developers]);

// Memoize callbacks
const handleSelect = useCallback((id: string) => {
  setSelectedId(id);
}, []);

// Memoize components
export const DeveloperCard = memo(function DeveloperCard({
  developer,
  onSelect,
}: DeveloperCardProps) {
  // Component implementation
});
```

### Lazy Loading

```tsx
import { lazy, Suspense } from 'react';
import { Routes, Route } from 'react-router-dom';

const Dashboard = lazy(() => import('./pages/Dashboard'));
const Analytics = lazy(() => import('./pages/Analytics'));
const Research = lazy(() => import('./pages/Research'));

function App() {
  return (
    <Suspense fallback={<PageSkeleton />}>
      <Routes>
        <Route path="/" element={<Dashboard />} />
        <Route path="/analytics" element={<Analytics />} />
        <Route path="/research" element={<Research />} />
      </Routes>
    </Suspense>
  );
}
```

## Accessibility

### ARIA Patterns

```tsx
// Accessible chart with description
<div role="img" aria-label={`Line chart showing ${title}`}>
  <span className="sr-only">
    {data.map(d => `${d.date}: ${d.aiRatio}%`).join(', ')}
  </span>
  <ResponsiveContainer>
    <LineChart data={data}>
      {/* Chart implementation */}
    </LineChart>
  </ResponsiveContainer>
</div>

// Accessible filter controls
<fieldset>
  <legend className="font-semibold">Filter by Seniority</legend>
  {seniorityOptions.map(option => (
    <label key={option} className="flex items-center gap-2">
      <input
        type="checkbox"
        checked={filters.seniority.includes(option)}
        onChange={() => toggleSeniority(option)}
        aria-describedby={`${option}-description`}
      />
      <span>{option}</span>
    </label>
  ))}
</fieldset>
```
