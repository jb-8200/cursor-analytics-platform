import { describe, it, expect } from 'vitest';
import { renderWithProviders } from '../utils';
import { useQuery, gql } from '@apollo/client';

// Use the HealthCheck query that's already mocked in graphqlHandlers
const HEALTH_CHECK_QUERY = gql`
  query HealthCheck {
    health {
      status
      timestamp
    }
  }
`;

// Simple test component that uses Apollo
function TestComponent() {
  const { data, loading, error } = useQuery(HEALTH_CHECK_QUERY);

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;
  if (data) return <div>Status: {data.health.status}</div>;
  return <div>No data</div>;
}

describe('Test Utilities', () => {
  describe('renderWithProviders', () => {
    it('should wrap component with Apollo Provider', async () => {
      const { findByText } = renderWithProviders(<TestComponent />);

      // Should show loading initially
      expect(await findByText(/loading/i)).toBeInTheDocument();

      // Should eventually show mocked data
      expect(await findByText(/status: ok/i)).toBeInTheDocument();
    });

    it('should provide Router context for navigation', () => {
      const { container } = renderWithProviders(
        <div>
          <a href="/dashboard">Dashboard</a>
        </div>
      );

      expect(container.querySelector('a[href="/dashboard"]')).toBeInTheDocument();
    });
  });
});
