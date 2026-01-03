import { describe, it, expect } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { MockedProvider } from '@apollo/client/testing';
import type { ReactNode } from 'react';
import { useDevelopers } from './useDevelopers';
import { GET_DEVELOPERS } from '../graphql/queries';

/**
 * Tests for useDevelopers hook
 */

describe('useDevelopers', () => {
  const mockDevelopersData = {
    developers: {
      nodes: [
        {
          id: 'dev1',
          name: 'Alice Developer',
          email: 'alice@example.com',
          team: 'Engineering',
          seniority: 'senior',
        },
        {
          id: 'dev2',
          name: 'Bob Developer',
          email: 'bob@example.com',
          team: 'Product',
          seniority: 'mid',
        },
      ],
      pageInfo: {
        hasNextPage: false,
        hasPreviousPage: false,
        startCursor: 'dev1',
        endCursor: 'dev2',
      },
    },
  };

  describe('without filters', () => {
    it('should return loading state initially', () => {
      const mocks = [
        {
          request: {
            query: GET_DEVELOPERS,
            variables: {},
          },
          result: {
            data: mockDevelopersData,
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

      const { result } = renderHook(() => useDevelopers(), { wrapper });

      expect(result.current.loading).toBe(true);
      expect(result.current.developers).toEqual([]);
    });

    it('should return developers list on success', async () => {
      const mocks = [
        {
          request: {
            query: GET_DEVELOPERS,
            variables: {},
          },
          result: {
            data: mockDevelopersData,
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

      const { result } = renderHook(() => useDevelopers(), { wrapper });

      await waitFor(() => {
        expect(result.current.loading).toBe(false);
      });

      expect(result.current.developers).toHaveLength(2);
      expect(result.current.developers[0].id).toBe('dev1');
      expect(result.current.pageInfo).toBeDefined();
    });
  });

  describe('with filters', () => {
    it('should pass filters to query variables', async () => {
      const filters = { team: 'Engineering', seniority: 'senior' };

      const mocks = [
        {
          request: {
            query: GET_DEVELOPERS,
            variables: { filters },
          },
          result: {
            data: {
              developers: {
                nodes: [mockDevelopersData.developers.nodes[0]],
                pageInfo: mockDevelopersData.developers.pageInfo,
              },
            },
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

      const { result } = renderHook(() => useDevelopers(filters), { wrapper });

      await waitFor(() => {
        expect(result.current.loading).toBe(false);
      });

      expect(result.current.developers).toHaveLength(1);
      expect(result.current.developers[0].team).toBe('Engineering');
    });
  });

  describe('with pagination', () => {
    it('should pass pagination to query variables', async () => {
      const pagination = { first: 10, after: 'cursor1' };

      const mocks = [
        {
          request: {
            query: GET_DEVELOPERS,
            variables: { pagination },
          },
          result: {
            data: mockDevelopersData,
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

      const { result } = renderHook(() => useDevelopers(pagination), { wrapper });

      await waitFor(() => {
        expect(result.current.loading).toBe(false);
      });

      expect(result.current.developers).toHaveLength(2);
      expect(result.current.pageInfo).toBeDefined();
    });
  });

  describe('refetch and fetchMore', () => {
    it('should provide refetch and fetchMore functions', async () => {
      const mocks = [
        {
          request: {
            query: GET_DEVELOPERS,
            variables: {},
          },
          result: {
            data: mockDevelopersData,
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

      const { result } = renderHook(() => useDevelopers(), { wrapper });

      await waitFor(() => {
        expect(result.current.loading).toBe(false);
      });

      expect(result.current.refetch).toBeDefined();
      expect(typeof result.current.refetch).toBe('function');
      expect(result.current.fetchMore).toBeDefined();
      expect(typeof result.current.fetchMore).toBe('function');
    });
  });
});
