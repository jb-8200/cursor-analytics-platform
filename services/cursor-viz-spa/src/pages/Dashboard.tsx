import React from 'react';
import { useDashboard } from '../hooks/useDashboard';
import { VelocityHeatmap } from '../components/charts/VelocityHeatmap';
import { TeamRadarChart } from '../components/charts/TeamRadarChart';
import { DeveloperTable } from '../components/charts/DeveloperTable';

const Dashboard: React.FC = () => {
  const { data, loading, error } = useDashboard();

  if (loading) {
    return (
      <div data-route="dashboard">
        <div className="mb-6">
          <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
          <p className="mt-2 text-sm text-gray-600">
            Overview of AI coding assistant usage across your organization.
          </p>
        </div>
        <div className="flex items-center justify-center h-64">
          <div className="text-gray-500">Loading dashboard data...</div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div data-route="dashboard">
        <div className="mb-6">
          <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
          <p className="mt-2 text-sm text-gray-600">
            Overview of AI coding assistant usage across your organization.
          </p>
        </div>
        <div className="bg-red-50 border border-red-200 rounded-lg p-4">
          <p className="text-red-800">Error loading dashboard: {error.message}</p>
        </div>
      </div>
    );
  }

  const velocityData = data?.dailyTrend?.map(trend => ({
    date: new Date(trend.date),
    count: trend.aiLinesAdded || 0,
    level: trend.aiLinesAdded > 100 ? 4 : trend.aiLinesAdded > 50 ? 3 : trend.aiLinesAdded > 20 ? 2 : trend.aiLinesAdded > 0 ? 1 : 0,
  })) || [];

  const teamData = data?.teamComparison?.map(team => ({
    teamName: team.teamName,
    acceptanceRate: team.averageAcceptanceRate || 0,
    aiVelocity: team.aiVelocity || 0,
    suggestions: team.totalSuggestions || 0,
    chatInteractions: team.chatInteractions || 0,
  })) || [];

  const developerData = data?.teamComparison?.flatMap(team =>
    team.topPerformers?.map(dev => ({
      id: dev.id,
      name: dev.name,
      email: dev.email,
      team: team.teamName,
      seniority: dev.seniority,
      acceptanceRate: 0, // Not available in summary data
      suggestionsCount: 0,
      aiVelocity: 0,
    })) || []
  ) || [];

  return (
    <div data-route="dashboard">
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
        <p className="mt-2 text-sm text-gray-600">
          Overview of AI coding assistant usage across your organization.
        </p>
      </div>

      {/* KPI Cards */}
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4 mb-6">
        <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
          <div className="text-sm font-medium text-gray-600">Total Developers</div>
          <div className="mt-2 text-3xl font-bold text-gray-900">
            {data?.totalDevelopers || 0}
          </div>
        </div>
        <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
          <div className="text-sm font-medium text-gray-600">Active Developers</div>
          <div className="mt-2 text-3xl font-bold text-gray-900">
            {data?.activeDevelopers || 0}
          </div>
        </div>
        <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
          <div className="text-sm font-medium text-gray-600">Acceptance Rate</div>
          <div className="mt-2 text-3xl font-bold text-gray-900">
            {data?.overallAcceptanceRate?.toFixed(1) || '0.0'}%
          </div>
        </div>
        <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
          <div className="text-sm font-medium text-gray-600">AI Velocity Today</div>
          <div className="mt-2 text-3xl font-bold text-gray-900">
            {data?.aiVelocityToday || 0}
          </div>
        </div>
      </div>

      {/* Dashboard grid for charts */}
      <div className="dashboard-grid grid grid-cols-1 gap-6 lg:grid-cols-2">
        <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
          <h2 className="text-lg font-semibold text-gray-700 mb-4">
            Velocity Heatmap
          </h2>
          <VelocityHeatmap data={velocityData} />
        </div>

        <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
          <h2 className="text-lg font-semibold text-gray-700 mb-4">
            Team Radar
          </h2>
          <TeamRadarChart teams={teamData} />
        </div>

        <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-200 lg:col-span-2">
          <h2 className="text-lg font-semibold text-gray-700 mb-4">
            Developer Table
          </h2>
          <DeveloperTable developers={developerData} />
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
