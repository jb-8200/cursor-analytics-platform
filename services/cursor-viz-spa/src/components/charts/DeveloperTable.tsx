/**
 * DeveloperTable Component
 *
 * A sortable, filterable table displaying individual developer metrics.
 * Features:
 * - Column sorting (click headers)
 * - Search filtering (debounced)
 * - Pagination
 * - Row click navigation
 * - Highlights low acceptance rates
 */

import React, { useState, useCallback, useMemo } from 'react';
import { Developer } from '../../graphql/types';

export interface DeveloperTableProps {
  data: Developer[];
  onSort?: (column: string, direction: 'asc' | 'desc') => void;
  onSearch?: (term: string) => void;
  onPageChange?: (page: number) => void;
  onRowClick?: (developer: Developer) => void;
  pageSize?: number;
  currentPage?: number;
  highlightThreshold?: number; // Acceptance rate below this is highlighted
}

type SortDirection = 'asc' | 'desc';

interface SortState {
  column: string;
  direction: SortDirection;
}

const DeveloperTable: React.FC<DeveloperTableProps> = ({
  data,
  onSort,
  onSearch,
  onPageChange,
  onRowClick,
  pageSize = 25,
  currentPage = 1,
  highlightThreshold = 20,
}) => {
  const [sortState, setSortState] = useState<SortState>({
    column: 'name',
    direction: 'asc',
  });
  const [searchTerm, setSearchTerm] = useState('');

  // Handle column sort
  const handleSort = useCallback(
    (column: string) => {
      const newDirection: SortDirection =
        sortState.column === column && sortState.direction === 'asc'
          ? 'desc'
          : 'asc';

      setSortState({ column, direction: newDirection });

      if (onSort) {
        onSort(column, newDirection);
      }
    },
    [sortState, onSort]
  );

  // Handle search input change
  const handleSearchChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const term = e.target.value;
      setSearchTerm(term);

      if (onSearch) {
        onSearch(term);
      }
    },
    [onSearch]
  );

  // Clear search
  const handleClearSearch = useCallback(() => {
    setSearchTerm('');
    if (onSearch) {
      onSearch('');
    }
  }, [onSearch]);

  // Sort and filter data
  const processedData = useMemo(() => {
    let result = [...data];

    // Filter by search term
    if (searchTerm) {
      const term = searchTerm.toLowerCase();
      result = result.filter((dev) =>
        dev.name.toLowerCase().includes(term)
      );
    }

    // Sort
    result.sort((a, b) => {
      let aVal: any;
      let bVal: any;

      switch (sortState.column) {
        case 'name':
          aVal = a.name;
          bVal = b.name;
          break;
        case 'team':
          aVal = a.team;
          bVal = b.team;
          break;
        case 'suggestions':
          aVal = a.stats?.totalSuggestions || 0;
          bVal = b.stats?.totalSuggestions || 0;
          break;
        case 'accepted':
          aVal = a.stats?.acceptedSuggestions || 0;
          bVal = b.stats?.acceptedSuggestions || 0;
          break;
        case 'rate':
          aVal = a.stats?.acceptanceRate || 0;
          bVal = b.stats?.acceptanceRate || 0;
          break;
        case 'aiLines':
          aVal = a.stats?.aiLinesAdded || 0;
          bVal = b.stats?.aiLinesAdded || 0;
          break;
        default:
          aVal = a.name;
          bVal = b.name;
      }

      if (typeof aVal === 'string') {
        return sortState.direction === 'asc'
          ? aVal.localeCompare(bVal)
          : bVal.localeCompare(aVal);
      } else {
        return sortState.direction === 'asc' ? aVal - bVal : bVal - aVal;
      }
    });

    return result;
  }, [data, searchTerm, sortState]);

  // Paginate data
  const paginatedData = useMemo(() => {
    const startIndex = (currentPage - 1) * pageSize;
    const endIndex = startIndex + pageSize;
    return processedData.slice(startIndex, endIndex);
  }, [processedData, currentPage, pageSize]);

  const totalPages = Math.ceil(processedData.length / pageSize);
  const startIndex = (currentPage - 1) * pageSize + 1;
  const endIndex = Math.min(currentPage * pageSize, processedData.length);

  // Render sort icon
  const renderSortIcon = (column: string) => {
    if (sortState.column !== column) {
      return (
        <svg
          className="w-4 h-4 text-gray-400"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"
          />
        </svg>
      );
    }

    return (
      <svg
        className="w-4 h-4 text-blue-600"
        data-testid="sort-icon"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
      >
        {sortState.direction === 'asc' ? (
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M5 15l7-7 7 7"
          />
        ) : (
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M19 9l-7 7-7-7"
          />
        )}
      </svg>
    );
  };

  if (data.length === 0) {
    return (
      <div className="text-center py-12 text-gray-500">
        No developers found
      </div>
    );
  }

  return (
    <div className="developer-table w-full">
      {/* Search Bar */}
      <div className="mb-4 flex items-center gap-2">
        <input
          type="text"
          placeholder="Search developers..."
          value={searchTerm}
          onChange={handleSearchChange}
          className="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        />
        {searchTerm && (
          <button
            onClick={handleClearSearch}
            className="px-4 py-2 text-gray-600 hover:text-gray-900"
            aria-label="Clear search"
          >
            Clear
          </button>
        )}
      </div>

      {/* Table */}
      <div className="overflow-x-auto border border-gray-200 rounded-lg">
        <table className="w-full divide-y divide-gray-200" role="table">
          <thead className="bg-gray-50">
            <tr>
              <th
                className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
                onClick={() => handleSort('name')}
                aria-sort={
                  sortState.column === 'name'
                    ? sortState.direction === 'asc'
                      ? 'ascending'
                      : 'descending'
                    : 'none'
                }
              >
                <div className="flex items-center gap-2">
                  Name
                  {renderSortIcon('name')}
                </div>
              </th>
              <th
                className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
                onClick={() => handleSort('team')}
                aria-sort={
                  sortState.column === 'team'
                    ? sortState.direction === 'asc'
                      ? 'ascending'
                      : 'descending'
                    : 'none'
                }
              >
                <div className="flex items-center gap-2">
                  Team
                  {renderSortIcon('team')}
                </div>
              </th>
              <th
                className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
                onClick={() => handleSort('suggestions')}
                aria-sort={
                  sortState.column === 'suggestions'
                    ? sortState.direction === 'asc'
                      ? 'ascending'
                      : 'descending'
                    : 'none'
                }
              >
                <div className="flex items-center gap-2">
                  Suggestions
                  {renderSortIcon('suggestions')}
                </div>
              </th>
              <th
                className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
                onClick={() => handleSort('accepted')}
                aria-sort={
                  sortState.column === 'accepted'
                    ? sortState.direction === 'asc'
                      ? 'ascending'
                      : 'descending'
                    : 'none'
                }
              >
                <div className="flex items-center gap-2">
                  Accepted
                  {renderSortIcon('accepted')}
                </div>
              </th>
              <th
                className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
                onClick={() => handleSort('rate')}
                aria-sort={
                  sortState.column === 'rate'
                    ? sortState.direction === 'asc'
                      ? 'ascending'
                      : 'descending'
                    : 'none'
                }
              >
                <div className="flex items-center gap-2">
                  Rate
                  {renderSortIcon('rate')}
                </div>
              </th>
              <th
                className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
                onClick={() => handleSort('aiLines')}
                aria-sort={
                  sortState.column === 'aiLines'
                    ? sortState.direction === 'asc'
                      ? 'ascending'
                      : 'descending'
                    : 'none'
                }
              >
                <div className="flex items-center gap-2">
                  AI Lines
                  {renderSortIcon('aiLines')}
                </div>
              </th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {paginatedData.map((developer) => {
              const acceptanceRate = (developer.stats?.acceptanceRate || 0) * 100;
              const isLowAcceptance = acceptanceRate < highlightThreshold;

              return (
                <tr
                  key={developer.id}
                  className={`hover:bg-gray-50 transition-colors ${
                    isLowAcceptance ? 'bg-red-50' : ''
                  } ${onRowClick ? 'cursor-pointer' : ''}`}
                  onClick={() => onRowClick && onRowClick(developer)}
                >
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="text-sm font-medium text-gray-900 truncate max-w-xs">
                      {developer.name}
                    </div>
                    <div className="text-sm text-gray-500">{developer.email}</div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {developer.team}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {developer.stats?.totalSuggestions ?? 'N/A'}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {developer.stats?.acceptedSuggestions ?? 'N/A'}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm">
                    <span
                      className={`font-medium ${
                        isLowAcceptance ? 'text-red-600' : 'text-gray-900'
                      }`}
                    >
                      {developer.stats
                        ? `${acceptanceRate.toFixed(1)}%`
                        : 'N/A'}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {developer.stats?.aiLinesAdded ?? 'N/A'}
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>

      {/* Pagination */}
      <div className="mt-4 flex items-center justify-between">
        <div className="text-sm text-gray-700">
          Showing {startIndex}-{endIndex} of {processedData.length}
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={() => onPageChange && onPageChange(currentPage - 1)}
            disabled={currentPage === 1}
            className="px-4 py-2 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Previous
          </button>
          <span className="text-sm text-gray-700">
            Page {currentPage} of {totalPages}
          </span>
          <button
            onClick={() => onPageChange && onPageChange(currentPage + 1)}
            disabled={currentPage >= totalPages}
            className="px-4 py-2 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Next
          </button>
        </div>
      </div>
    </div>
  );
};

export default DeveloperTable;
