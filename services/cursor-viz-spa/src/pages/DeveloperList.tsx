import React from 'react';

const DeveloperList: React.FC = () => {
  return (
    <div data-route="developers">
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-900">Developers</h1>
        <p className="mt-2 text-sm text-gray-600">
          Individual developer metrics and performance.
        </p>
      </div>

      <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-200">
        <p className="text-gray-500">Developer list will be implemented in future tasks.</p>
      </div>
    </div>
  );
};

export default DeveloperList;
