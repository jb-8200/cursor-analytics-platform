import { render, RenderOptions } from '@testing-library/react';
import { ReactElement } from 'react';

/**
 * Custom render function that wraps components with necessary providers
 * Can be extended to include Apollo Provider, Router, etc. as needed
 */
export function renderWithProviders(
  ui: ReactElement,
  options?: Omit<RenderOptions, 'wrapper'>
) {
  // Add providers here as they are implemented
  // const AllProviders = ({ children }: { children: React.ReactNode }) => {
  //   return (
  //     <ApolloProvider client={mockClient}>
  //       <BrowserRouter>{children}</BrowserRouter>
  //     </ApolloProvider>
  //   );
  // };

  return render(ui, { ...options });
}

export * from '@testing-library/react';
export { renderWithProviders as render };
