/**
 * TeamRadarChart Component
 *
 * A multi-axis radar chart comparing teams across different metrics.
 * Displays 2-5 teams simultaneously as overlapping polygons with
 * metrics normalized to a 0-100 scale for comparability.
 */

import React, { useMemo } from 'react';
import {
  Radar,
  RadarChart,
  PolarGrid,
  PolarAngleAxis,
  PolarRadiusAxis,
  ResponsiveContainer,
  Legend,
  Tooltip,
} from 'recharts';
import { TeamStats } from '../../graphql/types';

export interface MetricConfig {
  key: keyof TeamStats;
  label: string;
  max?: number;
}

export interface TeamRadarChartProps {
  data: TeamStats[];
  selectedTeams: string[];
  onTeamSelect?: (teams: string[]) => void;
  metrics?: MetricConfig[];
  width?: number;
  height?: number;
}

// Default metrics configuration
const DEFAULT_METRICS: MetricConfig[] = [
  { key: 'chatInteractions', label: 'Chat Usage', max: 1000 },
  { key: 'totalSuggestions', label: 'Completions', max: 2000 },
  { key: 'averageAcceptanceRate', label: 'Acceptance', max: 1 },
  { key: 'aiVelocity', label: 'AI Velocity', max: 1 },
];

// Team colors for up to 5 teams
const TEAM_COLORS = [
  '#3b82f6', // blue-500
  '#10b981', // green-500
  '#f59e0b', // amber-500
  '#ef4444', // red-500
  '#8b5cf6', // violet-500
];

/**
 * Normalize a value to 0-100 scale based on max value
 */
function normalizeValue(value: number | undefined | null, max: number): number {
  if (value === undefined || value === null || max === 0) return 0;

  // For percentage values (0-1), convert to 0-100
  if (max === 1) {
    return value * 100;
  }

  // For other values, normalize to 100 scale
  return Math.min((value / max) * 100, 100);
}

/**
 * Transform team data for radar chart
 */
function transformDataForRadar(
  teams: TeamStats[],
  selectedTeams: string[],
  metrics: MetricConfig[]
): any[] {
  // Get only selected teams that exist in data
  const selectedTeamData = teams.filter((team) =>
    selectedTeams.includes(team.teamName)
  );

  if (selectedTeamData.length === 0) {
    return [];
  }

  // Transform to radar chart format: one object per metric with team values
  return metrics.map((metric) => {
    const dataPoint: any = {
      metric: metric.label,
      fullMark: 100, // All metrics normalized to 100
    };

    selectedTeamData.forEach((team) => {
      const rawValue = team[metric.key];
      const value = typeof rawValue === 'number' ? rawValue : 0;
      dataPoint[team.teamName] = normalizeValue(value, metric.max || 100);
    });

    return dataPoint;
  });
}

/**
 * Custom tooltip to show actual values
 */
const CustomTooltip = ({ active, payload }: any) => {
  if (!active || !payload || payload.length === 0) return null;

  return (
    <div className="bg-white border border-gray-200 rounded-lg shadow-lg p-3 text-sm">
      <p className="font-semibold text-gray-900 mb-2">{payload[0].payload.metric}</p>
      {payload.map((entry: any, index: number) => (
        <div key={`item-${index}`} className="flex items-center gap-2 mb-1">
          <div
            className="w-3 h-3 rounded-full"
            style={{ backgroundColor: entry.color }}
          />
          <span className="text-gray-700">{entry.name}:</span>
          <span className="font-medium text-gray-900">
            {entry.value.toFixed(1)}
          </span>
        </div>
      ))}
    </div>
  );
};

const TeamRadarChart: React.FC<TeamRadarChartProps> = ({
  data,
  selectedTeams,
  onTeamSelect,
  metrics = DEFAULT_METRICS,
  width,
  height,
}) => {
  // Transform data for chart
  const chartData = useMemo(
    () => transformDataForRadar(data, selectedTeams, metrics),
    [data, selectedTeams, metrics]
  );

  // Get selected team data for rendering
  const selectedTeamData = useMemo(
    () => data.filter((team) => selectedTeams.includes(team.teamName)),
    [data, selectedTeams]
  );

  if (chartData.length === 0) {
    return (
      <div className="team-radar-chart flex items-center justify-center h-64 text-gray-500">
        {selectedTeams.length === 0
          ? 'Select teams to compare'
          : 'No data available for selected teams'}
      </div>
    );
  }

  return (
    <div className="team-radar-chart" role="img" aria-label="Team comparison radar chart">
      <ResponsiveContainer width={width || '100%'} height={height || 400}>
        <RadarChart data={chartData} margin={{ top: 20, right: 30, bottom: 20, left: 30 }}>
          <PolarGrid stroke="#e5e7eb" />
          <PolarAngleAxis
            dataKey="metric"
            tick={{ fill: '#6b7280', fontSize: 12 }}
            tickLine={false}
          />
          <PolarRadiusAxis
            angle={90}
            domain={[0, 100]}
            tick={{ fill: '#9ca3af', fontSize: 10 }}
            tickCount={5}
          />

          {/* Render a Radar for each selected team */}
          {selectedTeamData.map((team, index) => (
            <Radar
              key={team.teamName}
              name={team.teamName}
              dataKey={team.teamName}
              stroke={TEAM_COLORS[index % TEAM_COLORS.length]}
              fill={TEAM_COLORS[index % TEAM_COLORS.length]}
              fillOpacity={0.25}
              strokeWidth={2}
            />
          ))}

          <Legend
            wrapperStyle={{
              paddingTop: '20px',
            }}
            iconType="circle"
          />
          <Tooltip content={<CustomTooltip />} />
        </RadarChart>
      </ResponsiveContainer>

      {/* Team Selection Controls (optional) */}
      {onTeamSelect && data.length > 0 && (
        <div className="mt-4 p-4 bg-gray-50 rounded-lg">
          <p className="text-sm font-medium text-gray-700 mb-2">Select Teams (2-5):</p>
          <div className="flex flex-wrap gap-2">
            {data.map((team) => {
              const isSelected = selectedTeams.includes(team.teamName);
              return (
                <button
                  key={team.teamName}
                  onClick={() => {
                    if (isSelected) {
                      onTeamSelect(selectedTeams.filter((t) => t !== team.teamName));
                    } else if (selectedTeams.length < 5) {
                      onTeamSelect([...selectedTeams, team.teamName]);
                    }
                  }}
                  disabled={!isSelected && selectedTeams.length >= 5}
                  className={`
                    px-3 py-1.5 text-sm rounded-md transition-colors
                    ${
                      isSelected
                        ? 'bg-blue-500 text-white hover:bg-blue-600'
                        : 'bg-white text-gray-700 border border-gray-300 hover:bg-gray-50'
                    }
                    disabled:opacity-50 disabled:cursor-not-allowed
                  `}
                  aria-pressed={isSelected}
                >
                  {team.teamName}
                </button>
              );
            })}
          </div>
        </div>
      )}
    </div>
  );
};

export default TeamRadarChart;
