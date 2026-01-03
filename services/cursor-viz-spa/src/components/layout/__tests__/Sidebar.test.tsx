import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { BrowserRouter } from 'react-router-dom';
import Sidebar from '../Sidebar';

describe('Sidebar', () => {
  const renderWithRouter = (component: React.ReactElement) => {
    return render(<BrowserRouter>{component}</BrowserRouter>);
  };

  it('renders navigation landmark', () => {
    renderWithRouter(<Sidebar />);

    expect(screen.getByRole('navigation')).toBeInTheDocument();
  });

  it('contains navigation links for all routes', () => {
    renderWithRouter(<Sidebar />);

    expect(screen.getByRole('link', { name: /dashboard/i })).toBeInTheDocument();
    expect(screen.getByRole('link', { name: /teams/i })).toBeInTheDocument();
    expect(screen.getByRole('link', { name: /developers/i })).toBeInTheDocument();
  });

  it('has responsive classes for mobile collapse', () => {
    renderWithRouter(<Sidebar />);

    const nav = screen.getByRole('navigation');
    // Should have classes for mobile hiding and desktop showing
    expect(nav).toHaveClass('md:block');
  });

  it('links point to correct routes', () => {
    renderWithRouter(<Sidebar />);

    expect(screen.getByRole('link', { name: /dashboard/i })).toHaveAttribute('href', '/dashboard');
    expect(screen.getByRole('link', { name: /teams/i })).toHaveAttribute('href', '/teams');
    expect(screen.getByRole('link', { name: /developers/i })).toHaveAttribute('href', '/developers');
  });
});
