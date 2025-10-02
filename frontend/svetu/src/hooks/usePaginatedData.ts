import { useState, useEffect, useCallback } from 'react';

export interface PaginatedDataState<T> {
  data: T[];
  loading: boolean;
  hasMore: boolean;
  error: string | null;
  page: number;
  totalPages: number;
  totalItems: number;
}

export interface PaginatedDataActions {
  loadMore: () => Promise<void>;
  refresh: () => Promise<void>;
  reset: () => Promise<void>;
}

export interface PaginatedDataConfig {
  pageSize?: number;
  initialPage?: number;
}

export interface PaginatedDataResponse<T> {
  data: T[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

export function usePaginatedData<T>(
  fetchFunction: (
    page: number,
    pageSize: number
  ) => Promise<PaginatedDataResponse<T>>,
  config: PaginatedDataConfig = {}
): [PaginatedDataState<T>, PaginatedDataActions] {
  const { pageSize = 20, initialPage = 1 } = config;

  const [state, setState] = useState<PaginatedDataState<T>>({
    data: [],
    loading: false,
    hasMore: true,
    error: null,
    page: initialPage,
    totalPages: 0,
    totalItems: 0,
  });

  const loadData = useCallback(
    async (page: number, append: boolean = false) => {
      setState((prev) => ({ ...prev, loading: true, error: null }));

      try {
        const response = await fetchFunction(page, pageSize);

        setState((prev) => ({
          ...prev,
          data: append ? [...prev.data, ...response.data] : response.data,
          loading: false,
          hasMore: page < response.totalPages,
          page: response.page,
          totalPages: response.totalPages,
          totalItems: response.total,
        }));
      } catch (error) {
        setState((prev) => ({
          ...prev,
          loading: false,
          error: error instanceof Error ? error.message : 'Unknown error',
        }));
      }
    },
    [fetchFunction, pageSize]
  );

  const loadMore = useCallback(async () => {
    if (state.loading || !state.hasMore) return;
    await loadData(state.page + 1, true);
  }, [loadData, state.loading, state.hasMore, state.page]);

  const refresh = useCallback(async () => {
    await loadData(initialPage, false);
  }, [loadData, initialPage]);

  const reset = useCallback(async () => {
    setState({
      data: [],
      loading: false,
      hasMore: true,
      error: null,
      page: initialPage,
      totalPages: 0,
      totalItems: 0,
    });
    // Reload data after reset
    await loadData(initialPage, false);
  }, [initialPage, loadData]);

  useEffect(() => {
    loadData(initialPage, false);
  }, [loadData, initialPage]);

  return [
    state,
    {
      loadMore,
      refresh,
      reset,
    },
  ];
}
