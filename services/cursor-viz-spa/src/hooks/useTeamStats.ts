import { useQuery } from '@apollo/client';
import { GET_TEAM_STATS } from '../graphql/queries';
import type {
  GetTeamStatsResponse,
  DateRangeInput,
} from '../graphql/types';

/**
 * Custom hook for fetching team statistics
 *
 * @param teamName - Optional team name filter
 * @param range - Optional date range filter
 * @returns Team statistics with loading and error states
 */
export function useTeamStats(teamName?: string, range?: DateRangeInput) {
  const { data, loading, error, refetch } = useQuery<GetTeamStatsResponse>(
    GET_TEAM_STATS,
    {
      variables: {
        ...(teamName && { teamName }),
        ...(range && { range }),
      },
      fetchPolicy: 'cache-and-network',
      notifyOnNetworkStatusChange: true,
    }
  );

  return {
    teamStats: data?.teamStats,
    loading,
    error,
    refetch,
  };
}

export type UseTeamStatsReturn = ReturnType<typeof useTeamStats>;
