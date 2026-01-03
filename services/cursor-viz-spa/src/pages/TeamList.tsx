import React from 'react';

const TeamList: React.FC = () => {
  return (
    <div data-route="teams">
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-900">Teams</h1>
        <p className="mt-2 text-sm text-gray-600">
          View and compare team analytics.
        </p>
      </div>

      <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
        <p className="text-gray-500">Team list will be implemented in future tasks.</p>
      </div>
    </div>
  );
};

export default TeamList;
