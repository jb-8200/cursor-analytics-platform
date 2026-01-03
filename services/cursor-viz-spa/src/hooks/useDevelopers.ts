import { useQuery } from '@apollo/client';
import { GET_DEVELOPERS } from '../graphql/queries';
import type {
  GetDevelopersResponse,
  DeveloperQueryInput,
} from '../graphql/types';

/**
 * Custom hook for fetching developers list
 *
 * @param queryInput - Optional query parameters (team, limit, offset, search)
 * @returns Developers list with loading and error states
 */
export function useDevelopers(queryInput?: DeveloperQueryInput) {
  const { data, loading, error, refetch, fetchMore } = useQuery<GetDevelopersResponse>(
    GET_DEVELOPERS,
    {
      variables: queryInput || {},
      fetchPolicy: 'cache-and-network',
      notifyOnNetworkStatusChange: true,
    }
  );

  return {
    developers: data?.developers?.nodes || [],
    pageInfo: data?.developers?.pageInfo,
    totalCount: data?.developers?.totalCount || 0,
    loading,
    error,
    refetch,
    fetchMore,
  };
}

export type UseDevelopersReturn = ReturnType<typeof useDevelopers>;
