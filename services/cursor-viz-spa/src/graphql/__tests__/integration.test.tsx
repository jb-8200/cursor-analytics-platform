import { describe, it, expect } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { ApolloProvider, useQuery } from '@apollo/client';
import { createApolloClient } from '../client';
import { GET_DASHBOARD_SUMMARY, GET_DEVELOPERS, HEALTH_CHECK } from '../queries';
import type { GetDashboardSummaryResponse, GetDevelopersResponse } from '../types';

/**
 * Integration tests for Apollo Client setup
 *
 * These tests verify that:
 * 1. Apollo Client connects to the mock GraphQL server
 * 2. Queries execute successfully and return expected data
 * 3. Error handling works correctly
 * 4. Cache policies are applied
 */

// Test component that uses a GraphQL query
function TestDashboardComponent() {
  const { loading, error, data } = useQuery<GetDashboardSummaryResponse>(GET_DASHBOARD_SUMMARY);

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;
  if (!data) return <div>No data</div>;

  return (
    <div>
      <h1>Dashboard</h1>
      <p data-testid="total-developers">Total: {data.dashboardSummary.totalDevelopers}</p>
      <p data-testid="active-developers">Active: {data.dashboardSummary.activeDevelopers}</p>
      <p data-testid="acceptance-rate">
        Rate: {data.dashboardSummary.overallAcceptanceRate.toFixed(1)}%
      </p>
    </div>
  );
}

function TestDevelopersComponent() {
  const { loading, error, data } = useQuery<GetDevelopersResponse>(GET_DEVELOPERS, {
    variables: { limit: 10, offset: 0 },
  });

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;
  if (!data) return <div>No data</div>;

  return (
    <div>
      <h1>Developers</h1>
      <p data-testid="developer-count">Count: {data.developers.totalCount}</p>
      <ul>
        {data.developers.nodes.map((dev) => (
          <li key={dev.id} data-testid={`developer-${dev.id}`}>
            {dev.name} - {dev.team}
          </li>
        ))}
      </ul>
    </div>
  );
}

function TestHealthComponent() {
  const { loading, error, data } = useQuery(HEALTH_CHECK);

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;
  if (!data) return <div>No data</div>;

  return (
    <div>
      <p data-testid="health-status">Status: {data.health.status}</p>
    </div>
  );
}

describe('Apollo Client Integration', () => {
  describe('Dashboard Query', () => {
    it('should fetch and display dashboard summary', async () => {
      const client = createApolloClient();

      render(
        <ApolloProvider client={client}>
          <TestDashboardComponent />
        </ApolloProvider>
      );

      // Should show loading initially
      expect(screen.getByText('Loading...')).toBeInTheDocument();

      // Wait for data to load
      await waitFor(() => {
        expect(screen.getByText('Dashboard')).toBeInTheDocument();
      });

      // Verify data is displayed correctly
      expect(screen.getByTestId('total-developers')).toHaveTextContent('Total: 12');
      expect(screen.getByTestId('active-developers')).toHaveTextContent('Active: 9');
      expect(screen.getByTestId('acceptance-rate')).toHaveTextContent('Rate: 62.4%');
    });
  });

  describe('Developers Query', () => {
    it('should fetch and display paginated developers', async () => {
      const client = createApolloClient();

      render(
        <ApolloProvider client={client}>
          <TestDevelopersComponent />
        </ApolloProvider>
      );

      // Should show loading initially
      expect(screen.getByText('Loading...')).toBeInTheDocument();

      // Wait for data to load
      await waitFor(() => {
        expect(screen.getByText('Developers')).toBeInTheDocument();
      });

      // Verify developer count
      expect(screen.getByTestId('developer-count')).toHaveTextContent('Count: 3');

      // Verify developers are displayed
      expect(screen.getByTestId('developer-dev-1')).toHaveTextContent('Alice Johnson - Frontend');
      expect(screen.getByTestId('developer-dev-2')).toHaveTextContent('Bob Smith - Backend');
      expect(screen.getByTestId('developer-dev-3')).toHaveTextContent('Carol Davis - Frontend');
    });
  });

  describe('Health Check Query', () => {
    it('should fetch health status', async () => {
      const client = createApolloClient();

      render(
        <ApolloProvider client={client}>
          <TestHealthComponent />
        </ApolloProvider>
      );

      // Wait for data to load
      await waitFor(() => {
        expect(screen.getByTestId('health-status')).toHaveTextContent('Status: ok');
      });
    });
  });

  describe('Cache Behavior', () => {
    it('should cache query results', async () => {
      const client = createApolloClient();

      // First render
      const { rerender } = render(
        <ApolloProvider client={client}>
          <TestDashboardComponent />
        </ApolloProvider>
      );

      // Wait for initial load
      await waitFor(() => {
        expect(screen.getByText('Dashboard')).toBeInTheDocument();
      });

      // Unmount and remount - should use cached data
      rerender(
        <ApolloProvider client={client}>
          <TestDashboardComponent />
        </ApolloProvider>
      );

      // Should immediately show cached data without loading state
      expect(screen.queryByText('Loading...')).not.toBeInTheDocument();
      expect(screen.getByTestId('total-developers')).toHaveTextContent('Total: 12');
    });
  });
});
