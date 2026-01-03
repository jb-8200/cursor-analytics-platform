import { useQuery } from '@apollo/client';
import { GET_DEVELOPERS } from '../graphql/queries';
import type {
  GetDevelopersResponse,
  DeveloperFilters,
  PaginationInput,
} from '../graphql/types';

/**
 * Custom hook for fetching developers list
 *
 * @param filters - Optional filters (team, seniority)
 * @param pagination - Optional pagination parameters
 * @returns Developers list with loading and error states
 */
export function useDevelopers(filters?: DeveloperFilters, pagination?: PaginationInput) {
  const { data, loading, error, refetch, fetchMore } = useQuery<GetDevelopersResponse>(
    GET_DEVELOPERS,
    {
      variables: {
        ...( filters && { filters }),
        ...(pagination && { pagination }),
      },
      fetchPolicy: 'cache-and-network',
      notifyOnNetworkStatusChange: true,
    }
  );

  return {
    developers: data?.developers?.nodes || [],
    pageInfo: data?.developers?.pageInfo,
    loading,
    error,
    refetch,
    fetchMore,
  };
}

export type UseDevelopersReturn = ReturnType<typeof useDevelopers>;
