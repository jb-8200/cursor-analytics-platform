import { useState, useCallback } from 'react';
import { subDays, subMonths, subYears, format } from 'date-fns';

export type DateRangePreset =
  | 'LAST_7_DAYS'
  | 'LAST_30_DAYS'
  | 'LAST_90_DAYS'
  | 'LAST_6_MONTHS'
  | 'LAST_1_YEAR'
  | 'CUSTOM';

export interface DateRange {
  from: Date;
  to: Date;
  preset?: DateRangePreset;
}

export interface UseDateRangeReturn {
  dateRange: DateRange;
  setDateRange: (range: Partial<DateRange>) => void;
  setPreset: (preset: DateRangePreset) => void;
  formatRange: () => string;
}

const PRESET_LABELS: Record<DateRangePreset, string> = {
  LAST_7_DAYS: 'Last 7 Days',
  LAST_30_DAYS: 'Last 30 Days',
  LAST_90_DAYS: 'Last 90 Days',
  LAST_6_MONTHS: 'Last 6 Months',
  LAST_1_YEAR: 'Last 1 Year',
  CUSTOM: 'Custom Range',
};

function calculatePresetRange(preset: DateRangePreset, now: Date = new Date()): DateRange {
  const to = now;
  let from: Date;

  switch (preset) {
    case 'LAST_7_DAYS':
      from = subDays(to, 7);
      break;
    case 'LAST_30_DAYS':
      from = subDays(to, 30);
      break;
    case 'LAST_90_DAYS':
      from = subDays(to, 90);
      break;
    case 'LAST_6_MONTHS':
      from = subMonths(to, 6);
      break;
    case 'LAST_1_YEAR':
      from = subYears(to, 1);
      break;
    case 'CUSTOM':
      // For custom, return current dates without modification
      return { from: to, to, preset };
    default:
      // Default to last 30 days
      from = subDays(to, 30);
  }

  return { from, to, preset };
}

export function useDateRange(initialRange?: DateRange): UseDateRangeReturn {
  const [dateRange, setDateRangeState] = useState<DateRange>(() => {
    if (initialRange) {
      return initialRange;
    }
    return calculatePresetRange('LAST_30_DAYS');
  });

  const setDateRange = useCallback((range: Partial<DateRange>) => {
    setDateRangeState((prev) => ({
      ...prev,
      ...range,
    }));
  }, []);

  const setPreset = useCallback((preset: DateRangePreset) => {
    if (preset === 'CUSTOM') {
      // For custom, only update the preset flag, keep current dates
      setDateRangeState((prev) => ({
        ...prev,
        preset,
      }));
    } else {
      // For preset ranges, calculate new dates
      const newRange = calculatePresetRange(preset);
      setDateRangeState(newRange);
    }
  }, []);

  const formatRange = useCallback(() => {
    if (dateRange.preset && dateRange.preset !== 'CUSTOM') {
      return PRESET_LABELS[dateRange.preset];
    }

    // Format custom range
    const fromStr = format(dateRange.from, 'MMM d, yyyy');
    const toStr = format(dateRange.to, 'MMM d, yyyy');
    return `${fromStr} - ${toStr}`;
  }, [dateRange]);

  return {
    dateRange,
    setDateRange,
    setPreset,
    formatRange,
  };
}
