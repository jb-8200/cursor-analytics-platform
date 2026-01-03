import { describe, it, expect, beforeEach, vi } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useDateRange } from './useDateRange';

describe('useDateRange', () => {
  beforeEach(() => {
    vi.useFakeTimers();
    vi.setSystemTime(new Date('2026-01-15T12:00:00Z'));
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  describe('initialization', () => {
    it('should initialize with last 30 days as default', () => {
      const { result } = renderHook(() => useDateRange());

      expect(result.current.dateRange.preset).toBe('LAST_30_DAYS');
      expect(result.current.dateRange.from).toBeInstanceOf(Date);
      expect(result.current.dateRange.to).toBeInstanceOf(Date);
    });

    it('should initialize with custom initial range', () => {
      const customRange = {
        from: new Date('2026-01-01'),
        to: new Date('2026-01-10'),
        preset: 'CUSTOM' as const,
      };

      const { result } = renderHook(() => useDateRange(customRange));

      expect(result.current.dateRange.preset).toBe('CUSTOM');
      expect(result.current.dateRange.from).toEqual(customRange.from);
      expect(result.current.dateRange.to).toEqual(customRange.to);
    });
  });

  describe('setDateRange', () => {
    it('should update date range', () => {
      const { result } = renderHook(() => useDateRange());

      const newRange = {
        from: new Date('2026-01-01'),
        to: new Date('2026-01-15'),
        preset: 'CUSTOM' as const,
      };

      act(() => {
        result.current.setDateRange(newRange);
      });

      expect(result.current.dateRange).toEqual(newRange);
    });

    it('should allow partial updates', () => {
      const { result } = renderHook(() => useDateRange());

      const originalTo = result.current.dateRange.to;

      act(() => {
        result.current.setDateRange({
          from: new Date('2026-01-01'),
        });
      });

      expect(result.current.dateRange.from).toEqual(new Date('2026-01-01'));
      expect(result.current.dateRange.to).toEqual(originalTo);
    });
  });

  describe('setPreset', () => {
    it('should set LAST_7_DAYS preset correctly', () => {
      const { result } = renderHook(() => useDateRange());

      act(() => {
        result.current.setPreset('LAST_7_DAYS');
      });

      expect(result.current.dateRange.preset).toBe('LAST_7_DAYS');

      // Should be 7 days before today
      const expectedFrom = new Date('2026-01-08T12:00:00Z');
      const expectedTo = new Date('2026-01-15T12:00:00Z');

      expect(result.current.dateRange.from.toISOString()).toBe(expectedFrom.toISOString());
      expect(result.current.dateRange.to.toISOString()).toBe(expectedTo.toISOString());
    });

    it('should set LAST_30_DAYS preset correctly', () => {
      const { result } = renderHook(() => useDateRange());

      act(() => {
        result.current.setPreset('LAST_30_DAYS');
      });

      expect(result.current.dateRange.preset).toBe('LAST_30_DAYS');

      // Should be 30 days before today
      const expectedFrom = new Date('2025-12-16T12:00:00Z');
      const expectedTo = new Date('2026-01-15T12:00:00Z');

      expect(result.current.dateRange.from.toISOString()).toBe(expectedFrom.toISOString());
      expect(result.current.dateRange.to.toISOString()).toBe(expectedTo.toISOString());
    });

    it('should set LAST_90_DAYS preset correctly', () => {
      const { result } = renderHook(() => useDateRange());

      act(() => {
        result.current.setPreset('LAST_90_DAYS');
      });

      expect(result.current.dateRange.preset).toBe('LAST_90_DAYS');

      // Should be 90 days before today
      const diffInDays = Math.floor(
        (result.current.dateRange.to.getTime() - result.current.dateRange.from.getTime()) /
          (1000 * 60 * 60 * 24)
      );

      expect(diffInDays).toBe(90);
      expect(result.current.dateRange.to.toISOString()).toBe(new Date('2026-01-15T12:00:00Z').toISOString());
    });

    it('should set LAST_6_MONTHS preset correctly', () => {
      const { result } = renderHook(() => useDateRange());

      act(() => {
        result.current.setPreset('LAST_6_MONTHS');
      });

      expect(result.current.dateRange.preset).toBe('LAST_6_MONTHS');

      // Should be 6 months before today (check month difference)
      const fromMonth = result.current.dateRange.from.getMonth();
      const toMonth = result.current.dateRange.to.getMonth();
      const monthDiff = (toMonth - fromMonth + 12) % 12;

      expect(monthDiff).toBe(6);
      expect(result.current.dateRange.to.toISOString()).toBe(new Date('2026-01-15T12:00:00Z').toISOString());
    });

    it('should set LAST_1_YEAR preset correctly', () => {
      const { result } = renderHook(() => useDateRange());

      act(() => {
        result.current.setPreset('LAST_1_YEAR');
      });

      expect(result.current.dateRange.preset).toBe('LAST_1_YEAR');

      // Should be 1 year before today
      const expectedFrom = new Date('2025-01-15T12:00:00Z');
      const expectedTo = new Date('2026-01-15T12:00:00Z');

      expect(result.current.dateRange.from.toISOString()).toBe(expectedFrom.toISOString());
      expect(result.current.dateRange.to.toISOString()).toBe(expectedTo.toISOString());
    });

    it('should set CUSTOM preset without changing dates', () => {
      const { result } = renderHook(() => useDateRange());

      const originalFrom = result.current.dateRange.from;
      const originalTo = result.current.dateRange.to;

      act(() => {
        result.current.setPreset('CUSTOM');
      });

      expect(result.current.dateRange.preset).toBe('CUSTOM');
      expect(result.current.dateRange.from).toEqual(originalFrom);
      expect(result.current.dateRange.to).toEqual(originalTo);
    });
  });

  describe('formatRange', () => {
    it('should format preset ranges correctly', () => {
      const { result } = renderHook(() => useDateRange());

      act(() => {
        result.current.setPreset('LAST_7_DAYS');
      });

      expect(result.current.formatRange()).toBe('Last 7 Days');

      act(() => {
        result.current.setPreset('LAST_30_DAYS');
      });

      expect(result.current.formatRange()).toBe('Last 30 Days');
    });

    it('should format custom ranges correctly', () => {
      const { result } = renderHook(() => useDateRange());

      act(() => {
        result.current.setDateRange({
          from: new Date('2026-01-01'),
          to: new Date('2026-01-15'),
          preset: 'CUSTOM',
        });
      });

      expect(result.current.formatRange()).toBe('Jan 1, 2026 - Jan 15, 2026');
    });
  });
});
