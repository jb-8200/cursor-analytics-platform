import { graphql, HttpResponse } from 'msw';
import type {
  DashboardSummary,
  PaginatedDevelopers,
  TeamStats,
  DailyStats,
} from '../../graphql/types';

/**
 * Mock GraphQL Handlers for Testing
 *
 * These handlers provide stub data for GraphQL queries during testing.
 * They simulate responses from cursor-analytics-core (P5).
 */

// ============================================================================
// Mock Data
// ============================================================================

const mockDevelopers = [
  {
    id: 'dev-1',
    name: 'Alice Johnson',
    email: 'alice@example.com',
    team: 'Frontend',
    seniority: 'senior' as const,
    active: true,
    stats: {
      totalSuggestions: 1250,
      acceptedSuggestions: 875,
      acceptanceRate: 70.0,
      aiLinesAdded: 4500,
      aiLinesDeleted: 320,
      humanLinesAdded: 1200,
      humanLinesDeleted: 150,
      chatInteractions: 45,
      aiVelocity: 0.79,
    },
  },
  {
    id: 'dev-2',
    name: 'Bob Smith',
    email: 'bob@example.com',
    team: 'Backend',
    seniority: 'mid' as const,
    active: true,
    stats: {
      totalSuggestions: 980,
      acceptedSuggestions: 588,
      acceptanceRate: 60.0,
      aiLinesAdded: 3200,
      aiLinesDeleted: 250,
      humanLinesAdded: 1800,
      humanLinesDeleted: 200,
      chatInteractions: 32,
      aiVelocity: 0.64,
    },
  },
  {
    id: 'dev-3',
    name: 'Carol Davis',
    email: 'carol@example.com',
    team: 'Frontend',
    seniority: 'junior' as const,
    active: true,
    stats: {
      totalSuggestions: 650,
      acceptedSuggestions: 260,
      acceptanceRate: 40.0,
      aiLinesAdded: 1800,
      aiLinesDeleted: 120,
      humanLinesAdded: 2400,
      humanLinesDeleted: 180,
      chatInteractions: 28,
      aiVelocity: 0.43,
    },
  },
];

const mockDailyStats: DailyStats[] = Array.from({ length: 30 }, (_, i) => {
  const date = new Date();
  date.setDate(date.getDate() - (29 - i));
  return {
    date: date.toISOString().split('T')[0],
    suggestionsShown: Math.floor(Math.random() * 500) + 200,
    suggestionsAccepted: Math.floor(Math.random() * 300) + 100,
    acceptanceRate: Math.random() * 30 + 50, // 50-80%
    aiLinesAdded: Math.floor(Math.random() * 2000) + 500,
    humanLinesAdded: Math.floor(Math.random() * 1000) + 300,
    chatInteractions: Math.floor(Math.random() * 50) + 10,
  };
});

const mockTeamStats: TeamStats[] = [
  {
    teamName: 'Frontend',
    memberCount: 5,
    activeMemberCount: 4,
    averageAcceptanceRate: 65.5,
    totalSuggestions: 3200,
    aiVelocity: 0.72,
    chatInteractions: 125,
    topPerformers: [mockDevelopers[0], mockDevelopers[2]],
  },
  {
    teamName: 'Backend',
    memberCount: 4,
    activeMemberCount: 3,
    averageAcceptanceRate: 58.3,
    totalSuggestions: 2800,
    aiVelocity: 0.61,
    chatInteractions: 98,
    topPerformers: [mockDevelopers[1]],
  },
];

// ============================================================================
// GraphQL Handlers
// ============================================================================

export const graphqlHandlers = [
  // Dashboard Summary Query
  graphql.query('GetDashboardSummary', () => {
    const dashboardSummary: DashboardSummary = {
      totalDevelopers: 12,
      activeDevelopers: 9,
      overallAcceptanceRate: 62.4,
      totalSuggestionsToday: 245,
      totalAcceptedToday: 153,
      aiVelocityToday: 0.68,
      teamComparison: mockTeamStats,
      dailyTrend: mockDailyStats,
    };

    return HttpResponse.json({
      data: {
        dashboardSummary,
      },
    });
  }),

  // Developers Query
  graphql.query('GetDevelopers', ({ variables }) => {
    const { team, limit = 25, offset = 0, search } = variables as {
      team?: string;
      limit?: number;
      offset?: number;
      search?: string;
    };

    let filteredDevelopers = [...mockDevelopers];

    // Filter by team
    if (team) {
      filteredDevelopers = filteredDevelopers.filter((dev) => dev.team === team);
    }

    // Filter by search
    if (search) {
      const searchLower = search.toLowerCase();
      filteredDevelopers = filteredDevelopers.filter((dev) =>
        dev.name.toLowerCase().includes(searchLower)
      );
    }

    // Pagination
    const totalCount = filteredDevelopers.length;
    const nodes = filteredDevelopers.slice(offset, offset + limit);
    const hasNextPage = offset + limit < totalCount;
    const hasPreviousPage = offset > 0;

    const paginatedDevelopers: PaginatedDevelopers = {
      nodes,
      totalCount,
      pageInfo: {
        hasNextPage,
        hasPreviousPage,
        startCursor: offset > 0 ? `cursor-${offset}` : undefined,
        endCursor: hasNextPage ? `cursor-${offset + limit}` : undefined,
      },
    };

    return HttpResponse.json({
      data: {
        developers: paginatedDevelopers,
      },
    });
  }),

  // Single Developer Query
  graphql.query('GetDeveloper', ({ variables }) => {
    const { id } = variables as { id: string };
    const developer = mockDevelopers.find((dev) => dev.id === id);

    if (!developer) {
      return HttpResponse.json({
        errors: [
          {
            message: `Developer with id ${id} not found`,
            path: ['developer'],
            extensions: { code: 'NOT_FOUND' },
          },
        ],
        data: null,
      });
    }

    return HttpResponse.json({
      data: {
        developer,
      },
    });
  }),

  // Team Stats Query
  graphql.query('GetTeamStats', ({ variables }) => {
    const { teamName } = variables as { teamName: string };
    const teamStats = mockTeamStats.find((team) => team.teamName === teamName);

    if (!teamStats) {
      return HttpResponse.json({
        errors: [
          {
            message: `Team ${teamName} not found`,
            path: ['teamStats'],
            extensions: { code: 'NOT_FOUND' },
          },
        ],
        data: null,
      });
    }

    return HttpResponse.json({
      data: {
        teamStats,
      },
    });
  }),

  // Teams List Query
  graphql.query('GetTeams', () => {
    return HttpResponse.json({
      data: {
        teams: ['Frontend', 'Backend', 'DevOps', 'QA'],
      },
    });
  }),

  // Daily Stats Query
  graphql.query('GetDailyStats', () => {
    return HttpResponse.json({
      data: {
        dailyStats: mockDailyStats,
      },
    });
  }),

  // Health Check Query
  graphql.query('HealthCheck', () => {
    return HttpResponse.json({
      data: {
        health: {
          status: 'ok',
          timestamp: new Date().toISOString(),
        },
      },
    });
  }),
];
