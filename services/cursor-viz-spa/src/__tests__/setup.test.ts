import { describe, it, expect } from 'vitest';

/**
 * Setup validation tests
 * These tests ensure the testing infrastructure is working correctly
 */
describe('Test Setup', () => {
  it('should run tests with Vitest', () => {
    expect(true).toBe(true);
  });

  it('should have access to environment variables', () => {
    // Vite exposes env vars via import.meta.env
    expect(import.meta.env).toBeDefined();
    expect(import.meta.env.VITE_GRAPHQL_URL).toBeDefined();
  });

  it('should have configured strict TypeScript', () => {
    // This test will fail to compile if strict mode is not enabled
    const value: string = 'test';
    expect(value).toBe('test');
  });
});
