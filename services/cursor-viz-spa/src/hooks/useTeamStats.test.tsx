import { describe, it, expect } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { MockedProvider } from '@apollo/client/testing';
import type { ReactNode } from 'react';
import { useTeamStats } from './useTeamStats';
import { GET_TEAM_STATS } from '../graphql/queries';

/**
 * Tests for useTeamStats hook
 */

describe('useTeamStats', () => {
  const mockTeamStatsData = {
    teamStats: {
      teamName: 'Engineering',
      memberCount: 25,
      activeMemberCount: 22,
      averageAcceptanceRate: 72.1,
      totalSuggestions: 5000,
      aiVelocity: 45.2,
      chatInteractions: 320,
      topPerformers: [
        {
          id: 'dev1',
          name: 'Alice',
          acceptanceRate: 85.5,
          aiVelocity: 52.3,
        },
      ],
    },
  };

  describe('without filters', () => {
    it('should return loading state initially', () => {
      const mocks = [
        {
          request: {
            query: GET_TEAM_STATS,
            variables: {},
          },
          result: {
            data: mockTeamStatsData,
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

      const { result } = renderHook(() => useTeamStats(), { wrapper });

      expect(result.current.loading).toBe(true);
      expect(result.current.teamStats).toBeUndefined();
    });

    it('should return team stats on success', async () => {
      const mocks = [
        {
          request: {
            query: GET_TEAM_STATS,
            variables: {},
          },
          result: {
            data: mockTeamStatsData,
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

      const { result } = renderHook(() => useTeamStats(), { wrapper });

      await waitFor(() => {
        expect(result.current.loading).toBe(false);
      });

      expect(result.current.teamStats).toBeDefined();
      expect(result.current.teamStats?.teamName).toBe('Engineering');
      expect(result.current.teamStats?.memberCount).toBe(25);
    });
  });

  describe('with team name filter', () => {
    it('should pass team name to query variables', async () => {
      const teamName = 'Engineering';

      const mocks = [
        {
          request: {
            query: GET_TEAM_STATS,
            variables: { teamName },
          },
          result: {
            data: mockTeamStatsData,
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

      const { result } = renderHook(() => useTeamStats(teamName), { wrapper });

      await waitFor(() => {
        expect(result.current.loading).toBe(false);
      });

      expect(result.current.teamStats?.teamName).toBe('Engineering');
    });
  });

  describe('with date range', () => {
    it('should pass date range to query variables', async () => {
      const teamName = 'Engineering';
      const dateRange = { from: '2026-01-01', to: '2026-01-31' };

      const mocks = [
        {
          request: {
            query: GET_TEAM_STATS,
            variables: { teamName, range: dateRange },
          },
          result: {
            data: mockTeamStatsData,
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

      const { result } = renderHook(() => useTeamStats(teamName, dateRange), { wrapper });

      await waitFor(() => {
        expect(result.current.loading).toBe(false);
      });

      expect(result.current.teamStats).toBeDefined();
    });
  });

  describe('refetch functionality', () => {
    it('should provide refetch function', async () => {
      const mocks = [
        {
          request: {
            query: GET_TEAM_STATS,
            variables: {},
          },
          result: {
            data: mockTeamStatsData,
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

      const { result } = renderHook(() => useTeamStats(), { wrapper });

      await waitFor(() => {
        expect(result.current.loading).toBe(false);
      });

      expect(result.current.refetch).toBeDefined();
      expect(typeof result.current.refetch).toBe('function');
    });
  });
});
