import { useState, useRef, useEffect } from 'react';
import { format } from 'date-fns';
import { DateRange, DateRangePreset } from '../../hooks/useDateRange';

export interface DateRangePickerProps {
  value: DateRange;
  onChange: (range: DateRange) => void;
  className?: string;
}

const PRESETS: { label: string; value: DateRangePreset }[] = [
  { label: 'Last 7 Days', value: 'LAST_7_DAYS' },
  { label: 'Last 30 Days', value: 'LAST_30_DAYS' },
  { label: 'Last 90 Days', value: 'LAST_90_DAYS' },
  { label: 'Last 6 Months', value: 'LAST_6_MONTHS' },
  { label: 'Last 1 Year', value: 'LAST_1_YEAR' },
  { label: 'Custom', value: 'CUSTOM' },
];

function formatDateRange(range: DateRange): string {
  if (range.preset && range.preset !== 'CUSTOM') {
    const preset = PRESETS.find((p) => p.value === range.preset);
    return preset?.label || 'Select Range';
  }

  const fromStr = format(range.from, 'MMM d, yyyy');
  const toStr = format(range.to, 'MMM d, yyyy');
  return `${fromStr} - ${toStr}`;
}

export function DateRangePicker({ value, onChange, className = '' }: DateRangePickerProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [selectedPreset, setSelectedPreset] = useState<DateRangePreset>(
    value.preset || 'CUSTOM'
  );
  const [customFrom, setCustomFrom] = useState(format(value.from, 'yyyy-MM-dd'));
  const [customTo, setCustomTo] = useState(format(value.to, 'yyyy-MM-dd'));

  const dropdownRef = useRef<HTMLDivElement>(null);
  const buttonRef = useRef<HTMLButtonElement>(null);

  // Close dropdown when clicking outside
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target as Node) &&
        !buttonRef.current?.contains(event.target as Node)
      ) {
        setIsOpen(false);
      }
    }

    if (isOpen) {
      document.addEventListener('mousedown', handleClickOutside);
      return () => {
        document.removeEventListener('mousedown', handleClickOutside);
      };
    }
  }, [isOpen]);

  // Sync custom inputs with value prop
  useEffect(() => {
    setCustomFrom(format(value.from, 'yyyy-MM-dd'));
    setCustomTo(format(value.to, 'yyyy-MM-dd'));
    setSelectedPreset(value.preset || 'CUSTOM');
  }, [value]);

  const handlePresetSelect = (preset: DateRangePreset) => {
    setSelectedPreset(preset);

    if (preset !== 'CUSTOM') {
      // For non-custom presets, close dropdown and trigger onChange
      // The parent component (using useDateRange) will calculate the actual dates
      onChange({
        ...value,
        preset,
      });
      setIsOpen(false);
    }
    // For CUSTOM, keep dropdown open to show date inputs
  };

  const handleCustomFromChange = (dateStr: string) => {
    setCustomFrom(dateStr);
    const newFrom = new Date(dateStr);

    if (!isNaN(newFrom.getTime())) {
      onChange({
        from: newFrom,
        to: value.to,
        preset: 'CUSTOM',
      });
    }
  };

  const handleCustomToChange = (dateStr: string) => {
    setCustomTo(dateStr);
    const newTo = new Date(dateStr);

    if (!isNaN(newTo.getTime())) {
      onChange({
        from: value.from,
        to: newTo,
        preset: 'CUSTOM',
      });
    }
  };

  const isValidRange = value.from <= value.to;

  const handleKeyDown = (event: React.KeyboardEvent, preset?: DateRangePreset) => {
    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault();
      if (preset) {
        handlePresetSelect(preset);
      } else {
        setIsOpen(!isOpen);
      }
    }
  };

  return (
    <div className={`relative ${className}`}>
      <button
        ref={buttonRef}
        type="button"
        onClick={() => setIsOpen(!isOpen)}
        onKeyDown={(e) => handleKeyDown(e)}
        className="flex items-center justify-between w-full px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
        aria-haspopup="true"
        aria-expanded={isOpen}
      >
        <span>{formatDateRange(value)}</span>
        <svg
          className={`w-5 h-5 ml-2 transition-transform ${isOpen ? 'rotate-180' : ''}`}
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 20 20"
          fill="currentColor"
          aria-hidden="true"
        >
          <path
            fillRule="evenodd"
            d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z"
            clipRule="evenodd"
          />
        </svg>
      </button>

      {isOpen && (
        <div
          ref={dropdownRef}
          className="absolute z-10 mt-2 w-full bg-white border border-gray-300 rounded-md shadow-lg"
        >
          <div className="py-1" role="menu" aria-orientation="vertical">
            {PRESETS.map((preset) => (
              <button
                key={preset.value}
                type="button"
                onClick={() => handlePresetSelect(preset.value)}
                onKeyDown={(e) => handleKeyDown(e, preset.value)}
                className={`block w-full px-4 py-2 text-sm text-left hover:bg-gray-100 focus:outline-none focus:bg-gray-100 ${
                  selectedPreset === preset.value
                    ? 'bg-primary-50 text-primary-700 font-medium'
                    : 'text-gray-700'
                }`}
                role="menuitem"
              >
                {preset.label}
              </button>
            ))}
          </div>

          {selectedPreset === 'CUSTOM' && (
            <div className="px-4 py-3 border-t border-gray-200">
              <div className="space-y-3">
                <div>
                  <label htmlFor="date-from" className="block text-xs font-medium text-gray-700 mb-1">
                    From
                  </label>
                  <input
                    id="date-from"
                    type="date"
                    value={customFrom}
                    onChange={(e) => handleCustomFromChange(e.target.value)}
                    className="block w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                  />
                </div>
                <div>
                  <label htmlFor="date-to" className="block text-xs font-medium text-gray-700 mb-1">
                    To
                  </label>
                  <input
                    id="date-to"
                    type="date"
                    value={customTo}
                    onChange={(e) => handleCustomToChange(e.target.value)}
                    className="block w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                  />
                </div>

                {!isValidRange && (
                  <div className="text-xs text-red-600">
                    From date must be before To date
                  </div>
                )}

                <button
                  type="button"
                  onClick={() => setIsOpen(false)}
                  className="w-full px-3 py-2 text-sm font-medium text-white bg-primary-500 rounded-md hover:bg-primary-600 focus:outline-none focus:ring-2 focus:ring-primary-500"
                >
                  Apply
                </button>
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
