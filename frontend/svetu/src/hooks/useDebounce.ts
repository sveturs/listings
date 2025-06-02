import { useEffect, useState } from 'react';

/**
 * Hook для debounce значений
 * @param value - значение которое нужно debounce
 * @param delay - задержка в миллисекундах (по умолчанию 300ms)
 * @returns дебаунсированное значение
 */
export function useDebounce<T>(value: T, delay: number = 300): T {
  const [debouncedValue, setDebouncedValue] = useState<T>(value);

  useEffect(() => {
    // Для null/undefined значений обновляем сразу без задержки
    if (value === null || value === undefined) {
      setDebouncedValue(value);
      return;
    }

    // Для пустых строк также обновляем сразу для лучшего UX
    if (typeof value === 'string' && value === '') {
      setDebouncedValue(value);
      return;
    }

    // Для falsy значений (0, false, etc.) обновляем без задержки
    if (!value && value !== 0) {
      setDebouncedValue(value);
      return;
    }

    const handler = setTimeout(() => {
      setDebouncedValue(value);
    }, delay);

    return () => {
      clearTimeout(handler);
    };
  }, [value, delay]);

  return debouncedValue;
}

/**
 * Hook для debounce callback функций
 * @param callback - функция для debounce
 * @param delay - задержка в миллисекундах (по умолчанию 300ms)
 * @returns дебаунсированная функция
 */
export function useDebouncedCallback<T extends (...args: unknown[]) => unknown>(
  callback: T,
  delay: number = 300
): T {
  const [debounceTimer, setDebounceTimer] = useState<NodeJS.Timeout | null>(
    null
  );

  const debouncedCallback = ((...args: Parameters<T>) => {
    if (debounceTimer) {
      clearTimeout(debounceTimer);
    }

    const newTimer = setTimeout(() => {
      callback(...args);
    }, delay);

    setDebounceTimer(newTimer);
  }) as T;

  useEffect(() => {
    return () => {
      if (debounceTimer) {
        clearTimeout(debounceTimer);
      }
    };
  }, [debounceTimer]);

  return debouncedCallback;
}
