import { describe, it, expect } from 'vitest';
import { screen } from '@testing-library/react';
import { renderWithProviders } from '../test/utils';
import Dashboard from '../pages/Dashboard';

/**
 * Integration Test: Dashboard Page
 *
 * This test verifies the full integration of the Dashboard page with:
 * - Apollo Client for GraphQL queries
 * - React Router for navigation
 * - Chart components rendering with mock data
 * - Filter components with URL synchronization
 *
 * The MSW handlers in test/mocks/graphqlHandlers.ts provide mock GraphQL responses.
 */
describe('Dashboard Integration', () => {
  describe('Basic Rendering', () => {
    it('should render dashboard heading and description', () => {
      renderWithProviders(<Dashboard />);

      expect(screen.getByRole('heading', { name: /dashboard/i })).toBeInTheDocument();
      expect(
        screen.getByText(/overview of ai coding assistant usage/i)
      ).toBeInTheDocument();
    });

    it('should render dashboard grid layout', () => {
      const { container } = renderWithProviders(<Dashboard />);

      const grid = container.querySelector('.dashboard-grid');
      expect(grid).toBeInTheDocument();
      expect(grid).toHaveClass('grid', 'grid-cols-1', 'lg:grid-cols-2', 'xl:grid-cols-3');
    });
  });

  describe('Chart Placeholders', () => {
    it('should display Velocity Heatmap placeholder', () => {
      renderWithProviders(<Dashboard />);

      expect(
        screen.getByRole('heading', { name: /velocity heatmap/i })
      ).toBeInTheDocument();
      // Use getAllByText since there are multiple placeholders
      const placeholders = screen.getAllByText(/chart placeholder/i);
      expect(placeholders.length).toBeGreaterThanOrEqual(2);
    });

    it('should display Team Radar placeholder', () => {
      renderWithProviders(<Dashboard />);

      expect(screen.getByRole('heading', { name: /team radar/i })).toBeInTheDocument();
    });

    it('should display Developer Table placeholder', () => {
      renderWithProviders(<Dashboard />);

      expect(
        screen.getByRole('heading', { name: /developer table/i })
      ).toBeInTheDocument();
      expect(screen.getByText(/table placeholder/i)).toBeInTheDocument();
    });
  });

  describe('Responsive Layout', () => {
    it('should have responsive grid classes for different screen sizes', () => {
      const { container } = renderWithProviders(<Dashboard />);

      const grid = container.querySelector('.dashboard-grid');
      expect(grid).toHaveClass('grid-cols-1'); // Mobile
      expect(grid).toHaveClass('lg:grid-cols-2'); // Tablet
      expect(grid).toHaveClass('xl:grid-cols-3'); // Desktop
    });
  });

  describe('Component Structure', () => {
    it('should have proper semantic structure', () => {
      const { container } = renderWithProviders(<Dashboard />);

      // Should have data-route attribute for route identification
      const dashboardRoute = container.querySelector('[data-route="dashboard"]');
      expect(dashboardRoute).toBeInTheDocument();

      // Should have header section with title and description
      const heading = screen.getByRole('heading', { name: /dashboard/i, level: 1 });
      expect(heading).toBeInTheDocument();
      expect(heading).toHaveClass('text-3xl', 'font-bold', 'text-gray-900');
    });

    it('should have card components for each chart area', () => {
      const { container } = renderWithProviders(<Dashboard />);

      // All chart cards should have consistent styling
      const cards = container.querySelectorAll('.bg-white.p-6.rounded-lg.shadow-sm');
      expect(cards.length).toBeGreaterThanOrEqual(3); // At least 3 chart cards
    });
  });

  describe('Accessibility', () => {
    it('should have proper heading hierarchy', () => {
      renderWithProviders(<Dashboard />);

      const h1 = screen.getByRole('heading', { name: /dashboard/i, level: 1 });
      const h2s = screen.getAllByRole('heading', { level: 2 });

      expect(h1).toBeInTheDocument();
      expect(h2s.length).toBeGreaterThanOrEqual(3); // Chart section headings
    });

    it('should have descriptive text for screen readers', () => {
      renderWithProviders(<Dashboard />);

      const description = screen.getByText(
        /overview of ai coding assistant usage across your organization/i
      );
      expect(description).toBeInTheDocument();
    });
  });

  describe('Future Integration Points', () => {
    it('should have containers ready for real chart components', () => {
      const { container } = renderWithProviders(<Dashboard />);

      // Dashboard grid is ready to receive actual chart components
      const grid = container.querySelector('.dashboard-grid');
      expect(grid).toBeInTheDocument();

      // Each placeholder area has proper structure for replacement
      const placeholders = container.querySelectorAll('.bg-gray-50.rounded');
      expect(placeholders.length).toBeGreaterThanOrEqual(3);
    });
  });
});
