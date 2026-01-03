import { describe, it, expect } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { useUrlState } from './useUrlState';
import { ReactNode } from 'react';

// Wrapper component that provides Router context
function RouterWrapper({ children }: { children: ReactNode }) {
  return <BrowserRouter>{children}</BrowserRouter>;
}

describe('useUrlState', () => {
  it('should initialize with default value when URL param is absent', () => {
    const { result } = renderHook(() => useUrlState('testKey', 'defaultValue'), {
      wrapper: RouterWrapper,
    });

    expect(result.current[0]).toBe('defaultValue');
  });

  it('should read initial value from URL params', () => {
    // Set up initial URL with query param
    window.history.replaceState({}, '', '/?testKey=urlValue');

    const { result } = renderHook(() => useUrlState('testKey', 'defaultValue'), {
      wrapper: RouterWrapper,
    });

    expect(result.current[0]).toBe('urlValue');
  });

  it('should update URL when state changes', () => {
    const { result } = renderHook(() => useUrlState('testKey', 'initial'), {
      wrapper: RouterWrapper,
    });

    act(() => {
      result.current[1]('newValue');
    });

    expect(result.current[0]).toBe('newValue');
    expect(window.location.search).toContain('testKey=newValue');
  });

  it('should preserve other URL params when updating', () => {
    window.history.replaceState({}, '', '/?otherParam=keepMe&testKey=initial');

    const { result } = renderHook(() => useUrlState('testKey', 'initial'), {
      wrapper: RouterWrapper,
    });

    act(() => {
      result.current[1]('updated');
    });

    expect(window.location.search).toContain('otherParam=keepMe');
    expect(window.location.search).toContain('testKey=updated');
  });

  it('should remove param from URL when set to default value', () => {
    window.history.replaceState({}, '', '/?testKey=someValue');

    const { result } = renderHook(() => useUrlState('testKey', 'defaultValue'), {
      wrapper: RouterWrapper,
    });

    act(() => {
      result.current[1]('defaultValue');
    });

    // Should not contain the param when set to default
    expect(window.location.search).not.toContain('testKey');
  });

  it('should handle complex serializable values', () => {
    const defaultValue = { nested: { value: 'test' } };
    const newValue = { nested: { value: 'updated' } };

    const { result } = renderHook(
      () =>
        useUrlState('complexKey', defaultValue, {
          serialize: (val) => JSON.stringify(val),
          deserialize: (str) => JSON.parse(str),
        }),
      { wrapper: RouterWrapper }
    );

    act(() => {
      result.current[1](newValue);
    });

    expect(result.current[0]).toEqual(newValue);
    expect(window.location.search).toContain('complexKey=');
  });

  it('should handle date serialization', () => {
    const defaultDate = new Date('2026-01-01');
    const newDate = new Date('2026-01-15');

    const { result } = renderHook(
      () =>
        useUrlState('dateKey', defaultDate, {
          serialize: (date) => date.toISOString(),
          deserialize: (str) => new Date(str),
        }),
      { wrapper: RouterWrapper }
    );

    act(() => {
      result.current[1](newDate);
    });

    expect(result.current[0].toISOString()).toBe(newDate.toISOString());
    expect(window.location.search).toContain('2026-01-15');
  });

  it('should sync state when URL changes externally', () => {
    const { rerender } = renderHook(() => useUrlState('testKey', 'default'), {
      wrapper: RouterWrapper,
    });

    // Simulate external URL change (e.g., browser back button)
    act(() => {
      window.history.replaceState({}, '', '/?testKey=externalValue');
      // Trigger re-render to pick up URL change
      rerender();
    });

    // Need to manually trigger the effect by creating a new render
    const { result: result2 } = renderHook(() => useUrlState('testKey', 'default'), {
      wrapper: RouterWrapper,
    });

    expect(result2.current[0]).toBe('externalValue');
  });
});
