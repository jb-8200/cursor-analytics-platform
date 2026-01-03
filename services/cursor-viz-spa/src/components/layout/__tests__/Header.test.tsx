import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import Header from '../Header';

describe('Header', () => {
  it('renders application title', () => {
    render(<Header />);

    expect(screen.getByText(/cursor analytics/i)).toBeInTheDocument();
  });

  it('renders as a header landmark', () => {
    render(<Header />);

    expect(screen.getByRole('banner')).toBeInTheDocument();
  });

  it('has sticky positioning for scroll behavior', () => {
    render(<Header />);

    const header = screen.getByRole('banner');
    expect(header).toHaveClass('sticky', 'top-0');
  });

  it('includes placeholder for user menu', () => {
    render(<Header />);

    // User menu placeholder should be present
    const header = screen.getByRole('banner');
    expect(header).toContainHTML('user-menu');
  });
});
