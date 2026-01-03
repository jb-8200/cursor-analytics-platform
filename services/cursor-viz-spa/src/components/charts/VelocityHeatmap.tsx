/**
 * VelocityHeatmap Component
 *
 * A GitHub-style contribution graph showing AI code acceptance intensity over time.
 * Displays a grid with 7 rows (days of week) and N columns (weeks),
 * with color intensity mapping to suggestions accepted count.
 */

import React, { useMemo } from 'react';
import { DailyStats } from '../../graphql/types';

export interface VelocityHeatmapProps {
  data: DailyStats[];
  weeks?: number;
  colorScale?: string[];
  onCellClick?: (date: Date) => void;
}

// Default GitHub-style color scale
const DEFAULT_COLOR_SCALE = [
  '#ebedf0', // 0 (no activity)
  '#9be9a8', // 1-25th percentile
  '#40c463', // 25-50th percentile
  '#30a14e', // 50-75th percentile
  '#216e39', // 75-100th percentile
];

const DAYS_OF_WEEK = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
const MONTHS = [
  'Jan',
  'Feb',
  'Mar',
  'Apr',
  'May',
  'Jun',
  'Jul',
  'Aug',
  'Sep',
  'Oct',
  'Nov',
  'Dec',
];

const CELL_SIZE = 12;
const CELL_GAP = 2;
const DAY_LABEL_WIDTH = 30;
const MONTH_LABEL_HEIGHT = 20;

/**
 * Generate a complete grid of dates covering the specified number of weeks
 * If data is provided, use the earliest date as the start, otherwise use today
 */
function generateDateGrid(
  weeks: number,
  data: DailyStats[],
  endDate: Date = new Date()
): Date[] {
  const dates: Date[] = [];

  // If we have data, use the earliest date as the start
  let startDate: Date;
  if (data.length > 0) {
    const sortedDates = data.map((d) => new Date(d.date)).sort((a, b) => a.getTime() - b.getTime());
    startDate = sortedDates[0];
  } else {
    // Start from weeks ago
    startDate = new Date(endDate);
    startDate.setDate(startDate.getDate() - weeks * 7 + 1);
  }

  const totalDays = weeks * 7;
  for (let i = 0; i < totalDays; i++) {
    const date = new Date(startDate);
    date.setDate(date.getDate() + i);
    dates.push(date);
  }

  return dates;
}

/**
 * Map data to a lookup by date string
 */
function createDataLookup(data: DailyStats[]): Map<string, DailyStats> {
  const lookup = new Map<string, DailyStats>();
  data.forEach((stat) => {
    lookup.set(stat.date, stat);
  });
  return lookup;
}

/**
 * Get color for a cell based on acceptance count
 */
function getCellColor(
  acceptedCount: number,
  thresholds: number[],
  colorScale: string[]
): string {
  if (acceptedCount === 0) return colorScale[0];

  for (let i = thresholds.length - 1; i >= 0; i--) {
    if (acceptedCount >= thresholds[i]) {
      return colorScale[i + 1];
    }
  }

  return colorScale[1];
}

/**
 * Calculate percentile thresholds for color mapping
 */
function calculateThresholds(data: DailyStats[]): number[] {
  if (data.length === 0) return [0, 0, 0, 0];

  const counts = data
    .map((d) => d.suggestionsAccepted)
    .filter((c) => c > 0)
    .sort((a, b) => a - b);

  if (counts.length === 0) return [0, 0, 0, 0];

  const percentiles = [0.25, 0.5, 0.75, 1.0];
  return percentiles.map((p) => {
    const index = Math.floor(counts.length * p) - 1;
    return counts[Math.max(0, index)] || 0;
  });
}

const VelocityHeatmap: React.FC<VelocityHeatmapProps> = ({
  data,
  weeks = 52,
  colorScale = DEFAULT_COLOR_SCALE,
  onCellClick,
}) => {
  // Generate grid and data lookup
  const { dateGrid, dataLookup, thresholds } = useMemo(() => {
    const grid = generateDateGrid(weeks, data);
    const lookup = createDataLookup(data);
    const thresh = calculateThresholds(data);
    return { dateGrid: grid, dataLookup: lookup, thresholds: thresh };
  }, [data, weeks]);

  // Calculate dimensions
  const weeksInGrid = weeks;
  const width = DAY_LABEL_WIDTH + weeksInGrid * (CELL_SIZE + CELL_GAP);
  const height = MONTH_LABEL_HEIGHT + 7 * (CELL_SIZE + CELL_GAP);

  // Group dates by week
  const weekGroups = useMemo(() => {
    const groups: Date[][] = [];
    for (let i = 0; i < dateGrid.length; i += 7) {
      groups.push(dateGrid.slice(i, i + 7));
    }
    return groups;
  }, [dateGrid]);

  // Get month labels
  const monthLabels = useMemo(() => {
    const labels: { month: string; x: number }[] = [];
    let currentMonth = -1;

    weekGroups.forEach((week, weekIndex) => {
      const firstDayOfWeek = week[0];
      const month = firstDayOfWeek.getMonth();

      if (month !== currentMonth) {
        currentMonth = month;
        labels.push({
          month: MONTHS[month],
          x: DAY_LABEL_WIDTH + weekIndex * (CELL_SIZE + CELL_GAP),
        });
      }
    });

    return labels;
  }, [weekGroups]);

  const handleCellClick = (date: Date) => {
    if (onCellClick) {
      onCellClick(date);
    }
  };

  return (
    <div className="velocity-heatmap">
      <svg
        width={width}
        height={height}
        role="img"
        aria-label="AI code acceptance activity heatmap"
        className="font-sans"
      >
        {/* Month Labels */}
        <g className="month-labels">
          {monthLabels.map((label, index) => (
            <text
              key={`month-${index}`}
              x={label.x}
              y={12}
              fontSize="10"
              fill="currentColor"
              className="text-gray-600"
            >
              {label.month}
            </text>
          ))}
        </g>

        {/* Day Labels */}
        <g className="day-labels">
          {[1, 3, 5].map((dayIndex) => (
            <text
              key={`day-${dayIndex}`}
              x={0}
              y={MONTH_LABEL_HEIGHT + dayIndex * (CELL_SIZE + CELL_GAP) + CELL_SIZE}
              fontSize="9"
              fill="currentColor"
              className="text-gray-600"
            >
              {DAYS_OF_WEEK[dayIndex]}
            </text>
          ))}
        </g>

        {/* Heatmap Grid */}
        <g className="cells" transform={`translate(${DAY_LABEL_WIDTH}, ${MONTH_LABEL_HEIGHT})`}>
          {weekGroups.map((week, weekIndex) =>
            week.map((date, dayIndex) => {
              const dateStr = date.toISOString().split('T')[0];
              const stats = dataLookup.get(dateStr);
              const acceptedCount = stats?.suggestionsAccepted || 0;
              const color = getCellColor(acceptedCount, thresholds, colorScale);
              const cellIndex = weekIndex * 7 + dayIndex;

              const tooltip = stats
                ? `${dateStr}: ${acceptedCount} suggestions accepted`
                : `${dateStr}: No activity`;

              return (
                <g key={`cell-${weekIndex}-${dayIndex}`}>
                  <title>{tooltip}</title>
                  <rect
                    data-testid={`heatmap-cell-${cellIndex}`}
                    x={weekIndex * (CELL_SIZE + CELL_GAP)}
                    y={dayIndex * (CELL_SIZE + CELL_GAP)}
                    width={CELL_SIZE}
                    height={CELL_SIZE}
                    fill={color}
                    stroke="#1b1f230a"
                    strokeWidth="1"
                    rx="2"
                    aria-label={tooltip}
                    className="cursor-pointer hover:stroke-gray-400 transition-all"
                    onClick={() => handleCellClick(date)}
                    tabIndex={weekIndex === 0 && dayIndex === 0 ? 0 : -1}
                  />
                </g>
              );
            })
          )}
        </g>
      </svg>

      {/* Legend */}
      <div className="mt-2 flex items-center justify-end gap-1 text-xs text-gray-600">
        <span>Less</span>
        {colorScale.map((color, index) => (
          <div
            key={`legend-${index}`}
            className="w-3 h-3 border border-gray-200 rounded-sm"
            style={{ backgroundColor: color }}
            title={`Level ${index}`}
          />
        ))}
        <span>More</span>
      </div>
    </div>
  );
};

export default VelocityHeatmap;
