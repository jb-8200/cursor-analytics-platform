import { useQuery } from '@apollo/client';
import { GET_DASHBOARD_SUMMARY } from '../graphql/queries';
import type {
  DateRangeInput,
  GetDashboardSummaryResponse,
} from '../graphql/types';

/**
 * Custom hook for fetching dashboard summary data
 *
 * @param range - Optional date range filter
 * @returns Dashboard summary data with loading and error states
 */
export function useDashboard(range?: DateRangeInput) {
  const { data, loading, error, refetch } = useQuery<GetDashboardSummaryResponse>(
    GET_DASHBOARD_SUMMARY,
    {
      variables: range ? { range } : {},
      fetchPolicy: 'cache-and-network',
      notifyOnNetworkStatusChange: true,
    }
  );

  return {
    data: data?.dashboardSummary,
    loading,
    error,
    refetch,
  };
}

export type UseDashboardReturn = ReturnType<typeof useDashboard>;
