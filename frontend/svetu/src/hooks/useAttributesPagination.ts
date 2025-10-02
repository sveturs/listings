import { useState, useCallback, useMemo } from 'react';
import { adminApi, Attribute } from '@/services/admin';
import { usePaginatedData, PaginatedDataResponse } from './usePaginatedData';

export interface AttributeFilters {
  search?: string;
  type?: string;
  categoryId?: number;
  groupId?: number;
}

export function useAttributesPagination(
  filters: AttributeFilters = {},
  pageSize: number = 50
) {
  const [localFilters, setLocalFilters] = useState<AttributeFilters>(filters);

  // Create fetch function with current filters
  const fetchAttributes = useCallback(
    async (
      page: number,
      size: number
    ): Promise<PaginatedDataResponse<Attribute>> => {
      const response = await adminApi.attributes.getAll(
        page,
        size,
        localFilters.search,
        localFilters.type
      );

      // Адаптируем ответ к ожидаемому формату
      return {
        data: response.data || [],
        total: response.total || 0,
        page: response.page || page,
        pageSize: response.page_size || size,
        totalPages:
          response.total_pages || Math.ceil((response.total || 0) / size),
      };
    },
    [localFilters.search, localFilters.type]
  );

  const [paginatedState, paginatedActions] = usePaginatedData<Attribute>(
    fetchAttributes,
    { pageSize }
  );

  // Filter data locally if needed (for category or group filtering)
  const filteredData = useMemo(() => {
    const data = paginatedState.data;

    if (localFilters.categoryId !== undefined) {
      // This would need to be implemented based on your data structure
      // For now, we'll assume the API handles this
    }

    if (localFilters.groupId !== undefined) {
      // This would need to be implemented based on your data structure
      // For now, we'll assume the API handles this
    }

    return data;
  }, [paginatedState.data, localFilters.categoryId, localFilters.groupId]);

  const updateFilters = useCallback(
    async (newFilters: Partial<AttributeFilters>) => {
      setLocalFilters((prev) => ({ ...prev, ...newFilters }));
      // Reset pagination when filters change
      await paginatedActions.reset();
    },
    [paginatedActions]
  );

  const clearFilters = useCallback(async () => {
    setLocalFilters({});
    await paginatedActions.reset();
  }, [paginatedActions]);

  return {
    // Data
    attributes: filteredData,
    loading: paginatedState.loading,
    error: paginatedState.error,
    hasMore: paginatedState.hasMore,
    totalItems: paginatedState.totalItems,

    // Filters
    filters: localFilters,
    updateFilters,
    clearFilters,

    // Actions
    loadMore: paginatedActions.loadMore,
    refresh: paginatedActions.refresh,
    reset: paginatedActions.reset,
  };
}
