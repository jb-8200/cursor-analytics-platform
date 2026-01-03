import { useCallback, useEffect, useState } from 'react';
import { useSearchParams } from 'react-router-dom';

export interface UseUrlStateOptions<T> {
  serialize?: (value: T) => string;
  deserialize?: (value: string) => T;
}

export type UseUrlStateReturn<T> = [T, (value: T) => void];

/**
 * A hook that synchronizes state with URL query parameters.
 *
 * @param key - The URL parameter key
 * @param defaultValue - The default value when the parameter is not present
 * @param options - Optional serialization/deserialization functions
 *
 * @example
 * // Simple string state
 * const [search, setSearch] = useUrlState('search', '');
 *
 * @example
 * // Complex object state
 * const [filters, setFilters] = useUrlState('filters', { team: '', seniority: '' }, {
 *   serialize: (val) => JSON.stringify(val),
 *   deserialize: (str) => JSON.parse(str),
 * });
 *
 * @example
 * // Date state
 * const [date, setDate] = useUrlState('date', new Date(), {
 *   serialize: (d) => d.toISOString(),
 *   deserialize: (str) => new Date(str),
 * });
 */
export function useUrlState<T>(
  key: string,
  defaultValue: T,
  options?: UseUrlStateOptions<T>
): UseUrlStateReturn<T> {
  const [searchParams, setSearchParams] = useSearchParams();

  const serialize = options?.serialize || ((value: T) => String(value));
  const deserialize = options?.deserialize || ((value: string) => value as unknown as T);

  // Initialize state from URL or default
  const [state, setState] = useState<T>(() => {
    const urlValue = searchParams.get(key);
    if (urlValue !== null) {
      try {
        return deserialize(urlValue);
      } catch (error) {
        console.warn(`Failed to deserialize URL param "${key}":`, error);
        return defaultValue;
      }
    }
    return defaultValue;
  });

  // Sync state with URL when URL changes externally (e.g., browser back/forward)
  useEffect(() => {
    const urlValue = searchParams.get(key);
    if (urlValue !== null) {
      try {
        const deserialized = deserialize(urlValue);
        setState(deserialized);
      } catch (error) {
        console.warn(`Failed to deserialize URL param "${key}":`, error);
      }
    } else {
      // If param is removed from URL, reset to default
      setState(defaultValue);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [searchParams, key]);

  // Update URL when state changes
  const updateState = useCallback(
    (newValue: T) => {
      setState(newValue);

      // Update URL params
      setSearchParams((prevParams) => {
        const newParams = new URLSearchParams(prevParams);

        // Compare with default value
        const isDefaultValue =
          JSON.stringify(newValue) === JSON.stringify(defaultValue);

        if (isDefaultValue) {
          // Remove param when set to default value to keep URL clean
          newParams.delete(key);
        } else {
          // Set param to serialized value
          try {
            const serialized = serialize(newValue);
            newParams.set(key, serialized);
          } catch (error) {
            console.warn(`Failed to serialize value for URL param "${key}":`, error);
          }
        }

        return newParams;
      });
    },
    [key, defaultValue, serialize, setSearchParams]
  );

  return [state, updateState];
}
