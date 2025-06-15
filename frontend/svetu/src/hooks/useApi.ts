import { useState, useEffect, useCallback } from 'react';
import { ApiResponse } from '@/services/api-client';

interface UseApiOptions {
  immediate?: boolean; // Выполнить запрос сразу
  onSuccess?: (data: any) => void;
  onError?: (error: any) => void;
}

interface UseApiResult<T> {
  data: T | null;
  error: any | null;
  loading: boolean;
  execute: (...args: any[]) => Promise<void>;
  reset: () => void;
}

export function useApi<T = any>(
  apiCall: (...args: any[]) => Promise<ApiResponse<T>>,
  options: UseApiOptions = {}
): UseApiResult<T> {
  const [data, setData] = useState<T | null>(null);
  const [error, setError] = useState<any | null>(null);
  const [loading, setLoading] = useState(false);

  const execute = useCallback(
    async (...args: any[]) => {
      setLoading(true);
      setError(null);

      try {
        const response = await apiCall(...args);

        if (response.error) {
          setError(response.error);
          options.onError?.(response.error);
        } else {
          setData(response.data || null);
          options.onSuccess?.(response.data);
        }
      } catch (err) {
        const error = { message: 'Unexpected error', details: err };
        setError(error);
        options.onError?.(error);
      } finally {
        setLoading(false);
      }
    },
    [apiCall, options]
  );

  const reset = useCallback(() => {
    setData(null);
    setError(null);
    setLoading(false);
  }, []);

  useEffect(() => {
    if (options.immediate) {
      execute();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return { data, error, loading, execute, reset };
}

// Пример использования:
// const { data, loading, error, execute } = useApi(
//   () => marketplaceApi.getListings({ page: 1 }),
//   { immediate: true }
// );
