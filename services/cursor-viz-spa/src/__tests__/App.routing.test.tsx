import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { MemoryRouter } from 'react-router-dom';
import AppLayout from '../components/layout/AppLayout';
import Dashboard from '../pages/Dashboard';
import TeamList from '../pages/TeamList';
import DeveloperList from '../pages/DeveloperList';

// Test routes without Apollo Provider since we're just testing routing
const TestApp = () => (
  <AppLayout>
    <Routes>
      <Route path="/" element={<Navigate to="/dashboard" replace />} />
      <Route path="/dashboard" element={<Dashboard />} />
      <Route path="/teams" element={<TeamList />} />
      <Route path="/developers" element={<DeveloperList />} />
    </Routes>
  </AppLayout>
);

describe('App Routing', () => {
  it('redirects root path to /dashboard', async () => {
    const { container } = render(
      <MemoryRouter initialEntries={['/']}>
        <TestApp />
      </MemoryRouter>
    );

    // Should redirect to dashboard
    const routeElement = container.querySelector('[data-route="dashboard"]');
    expect(routeElement).toBeInTheDocument();
  });

  it('renders Dashboard route at /dashboard', () => {
    const { container } = render(
      <MemoryRouter initialEntries={['/dashboard']}>
        <TestApp />
      </MemoryRouter>
    );

    const routeElement = container.querySelector('[data-route="dashboard"]');
    expect(routeElement).toBeInTheDocument();
    // Check for the Dashboard page content
    expect(screen.getByText(/overview of ai coding assistant usage/i)).toBeInTheDocument();
  });

  it('renders Teams route at /teams', () => {
    const { container } = render(
      <MemoryRouter initialEntries={['/teams']}>
        <TestApp />
      </MemoryRouter>
    );

    const routeElement = container.querySelector('[data-route="teams"]');
    expect(routeElement).toBeInTheDocument();
  });

  it('renders Developers route at /developers', () => {
    const { container } = render(
      <MemoryRouter initialEntries={['/developers']}>
        <TestApp />
      </MemoryRouter>
    );

    const routeElement = container.querySelector('[data-route="developers"]');
    expect(routeElement).toBeInTheDocument();
  });
});
