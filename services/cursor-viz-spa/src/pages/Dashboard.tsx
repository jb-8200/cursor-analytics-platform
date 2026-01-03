import React from 'react';

const Dashboard: React.FC = () => {
  return (
    <div data-route="dashboard">
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
        <p className="mt-2 text-sm text-gray-600">
          Overview of AI coding assistant usage across your organization.
        </p>
      </div>

      {/* Dashboard grid for charts (placeholder) */}
      <div className="dashboard-grid grid grid-cols-1 gap-6 lg:grid-cols-2 xl:grid-cols-3">
        {/* Chart components will be added in TASK04 */}
        <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
          <h2 className="text-lg font-semibold text-gray-700 mb-4">
            Velocity Heatmap
          </h2>
          <div className="h-48 bg-gray-50 rounded flex items-center justify-center text-gray-400">
            Chart placeholder
          </div>
        </div>

        <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
          <h2 className="text-lg font-semibold text-gray-700 mb-4">
            Team Radar
          </h2>
          <div className="h-48 bg-gray-50 rounded flex items-center justify-center text-gray-400">
            Chart placeholder
          </div>
        </div>

        <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
          <h2 className="text-lg font-semibold text-gray-700 mb-4">
            Developer Table
          </h2>
          <div className="h-48 bg-gray-50 rounded flex items-center justify-center text-gray-400">
            Table placeholder
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
