import { describe, it, expect } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { MockedProvider } from '@apollo/client/testing';
import type { ReactNode } from 'react';
import { useDashboard } from './useDashboard';
import { GET_DASHBOARD_SUMMARY } from '../graphql/queries';
import type { DashboardSummary } from '../graphql/types';

/**
 * Tests for useDashboard hook
 */

describe('useDashboard', () => {
  const mockDashboardData: DashboardSummary = {
    totalDevelopers: 50,
    activeDevelopers: 42,
    overallAcceptanceRate: 68.5,
    totalSuggestionsToday: 1250,
    totalAcceptedToday: 856,
    aiVelocityToday: 42.3,
    teamComparison: [
      {
        teamName: 'Engineering',
        memberCount: 25,
        activeMemberCount: 22,
        averageAcceptanceRate: 72.1,
        totalSuggestions: 5000,
        aiVelocity: 45.2,
        chatInteractions: 320,
        topPerformers: [],
      },
    ],
    dailyTrend: [
      {
        date: '2026-01-01',
        suggestionsShown: 1200,
        suggestionsAccepted: 820,
        acceptanceRate: 68.3,
        aiLinesAdded: 4500,
        humanLinesAdded: 2100,
        chatInteractions: 45,
      },
    ],
  };

  describe('without date range', () => {
    it('should return loading state initially', () => {
      const mocks = [
        {
          request: {
            query: GET_DASHBOARD_SUMMARY,
            variables: {},
          },
          result: {
            data: { dashboardSummary: mockDashboardData },
          },
        },
      ];

      const wrapper = ({ children }: { children: ReactNode }) => {
        return (
          <MockedProvider mocks={mocks} addTypename={false}>
            {children}
          </MockedProvider>
        );
      };

      const { result } = renderHook(() => useDashboard(), { wrapper });

      expect(result.current.loading).toBe(true);
      expect(result.current.data).toBeUndefined();
      expect(result.current.error).toBeUndefined();
    });

    it('should return dashboard data on success', async () => {
      const mocks = [
        {
          request: {
            query: GET_DASHBOARD_SUMMARY,
            variables: {},
          },
          result: {
            data: { dashboardSummary: mockDashboardData },
          },
        },
      ];

      const wrapper = ({ children }: { children: ReactNode }) => {
        return (
          <MockedProvider mocks={mocks} addTypename={false}>
            {children}
          </MockedProvider>
        );
      };

      const { result } = renderHook(() => useDashboard(), { wrapper });

      await waitFor(() => {
        expect(result.current.loading).toBe(false);
      });

      expect(result.current.data).toEqual(mockDashboardData);
      expect(result.current.error).toBeUndefined();
    });

    it('should handle errors', async () => {
      const mocks = [
        {
          request: {
            query: GET_DASHBOARD_SUMMARY,
            variables: {},
          },
          error: new Error('Network error'),
        },
      ];

      const wrapper = ({ children }: { children: ReactNode }) => {
        return (
          <MockedProvider mocks={mocks} addTypename={false}>
            {children}
          </MockedProvider>
        );
      };

      const { result } = renderHook(() => useDashboard(), { wrapper });

      await waitFor(() => {
        expect(result.current.loading).toBe(false);
      });

      expect(result.current.error).toBeDefined();
      expect(result.current.data).toBeUndefined();
    });
  });

  describe('with date range', () => {
    it('should pass date range to query variables', async () => {
      const dateRange = {
        from: '2026-01-01',
        to: '2026-01-31',
      };

      const mocks = [
        {
          request: {
            query: GET_DASHBOARD_SUMMARY,
            variables: { range: dateRange },
          },
          result: {
            data: { dashboardSummary: mockDashboardData },
          },
        },
      ];

      const wrapper = ({ children }: { children: ReactNode }) => {
        return (
          <MockedProvider mocks={mocks} addTypename={false}>
            {children}
          </MockedProvider>
        );
      };

      const { result } = renderHook(() => useDashboard(dateRange), { wrapper });

      await waitFor(() => {
        expect(result.current.loading).toBe(false);
      });

      expect(result.current.data).toEqual(mockDashboardData);
    });
  });

  describe('refetch functionality', () => {
    it('should provide refetch function', async () => {
      const mocks = [
        {
          request: {
            query: GET_DASHBOARD_SUMMARY,
            variables: {},
          },
          result: {
            data: { dashboardSummary: mockDashboardData },
          },
        },
      ];

      const wrapper = ({ children }: { children: ReactNode }) => {
        return (
          <MockedProvider mocks={mocks} addTypename={false}>
            {children}
          </MockedProvider>
        );
      };

      const { result } = renderHook(() => useDashboard(), { wrapper });

      await waitFor(() => {
        expect(result.current.loading).toBe(false);
      });

      expect(result.current.refetch).toBeDefined();
      expect(typeof result.current.refetch).toBe('function');
    });
  });
});
