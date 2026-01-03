import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { BrowserRouter } from 'react-router-dom';
import Dashboard from '../Dashboard';

describe('Dashboard Page', () => {
  const renderWithRouter = (component: React.ReactElement) => {
    return render(<BrowserRouter>{component}</BrowserRouter>);
  };

  it('renders with dashboard route identifier', () => {
    const { container } = renderWithRouter(<Dashboard />);

    const routeElement = container.querySelector('[data-route="dashboard"]');
    expect(routeElement).toBeInTheDocument();
  });

  it('displays dashboard heading', () => {
    renderWithRouter(<Dashboard />);

    expect(screen.getByRole('heading', { name: /dashboard/i })).toBeInTheDocument();
  });

  it('has placeholder for charts grid', () => {
    const { container } = renderWithRouter(<Dashboard />);

    // Should have a container for future chart components
    expect(container.querySelector('.dashboard-grid')).toBeInTheDocument();
  });
});
