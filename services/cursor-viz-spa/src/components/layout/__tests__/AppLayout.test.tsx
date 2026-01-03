import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { BrowserRouter } from 'react-router-dom';
import AppLayout from '../AppLayout';

describe('AppLayout', () => {
  const renderWithRouter = (children: React.ReactNode) => {
    return render(<BrowserRouter>{children}</BrowserRouter>);
  };

  it('renders header, sidebar, and main content area', () => {
    renderWithRouter(
      <AppLayout>
        <div>Test Content</div>
      </AppLayout>
    );

    expect(screen.getByRole('banner')).toBeInTheDocument(); // header
    expect(screen.getByRole('navigation')).toBeInTheDocument(); // sidebar
    expect(screen.getByRole('main')).toBeInTheDocument(); // main content
    expect(screen.getByText('Test Content')).toBeInTheDocument();
  });

  it('applies responsive layout classes', () => {
    const { container } = renderWithRouter(
      <AppLayout>
        <div>Content</div>
      </AppLayout>
    );

    // Check the top-level div has the flex and min-h-screen classes
    const layout = container.querySelector('.flex.min-h-screen');
    expect(layout).toBeInTheDocument();
  });

  it('renders children in main content area', () => {
    renderWithRouter(
      <AppLayout>
        <div data-testid="child-content">Child Content</div>
      </AppLayout>
    );

    const main = screen.getByRole('main');
    expect(main).toContainElement(screen.getByTestId('child-content'));
  });
});
