import { render, RenderOptions } from '@testing-library/react';
import { ReactElement, ReactNode } from 'react';
import { ApolloProvider } from '@apollo/client';
import { MemoryRouter, MemoryRouterProps } from 'react-router-dom';
import { createApolloClient } from '../graphql/client';

/**
 * Options for renderWithProviders function
 */
export interface TestRenderOptions extends Omit<RenderOptions, 'wrapper'> {
  /**
   * Initial route entries for MemoryRouter
   * @default ['/']
   */
  routerProps?: MemoryRouterProps;
}

/**
 * Custom render function that wraps components with necessary providers
 * Includes Apollo Provider with test client and React Router's MemoryRouter
 */
export function renderWithProviders(
  ui: ReactElement,
  options?: TestRenderOptions
) {
  const { routerProps, ...renderOptions } = options || {};

  // Create a test Apollo client instance
  // MSW will intercept the GraphQL requests
  const testClient = createApolloClient();

  function AllProviders({ children }: { children: ReactNode }) {
    return (
      <ApolloProvider client={testClient}>
        <MemoryRouter {...routerProps}>{children}</MemoryRouter>
      </ApolloProvider>
    );
  }

  return render(ui, { wrapper: AllProviders, ...renderOptions });
}

export * from '@testing-library/react';
export { renderWithProviders as render };
