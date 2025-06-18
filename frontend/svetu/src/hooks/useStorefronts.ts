import { useCallback } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import {
  fetchStorefronts,
  fetchStorefrontById,
  fetchMyStorefronts,
  createStorefront,
  updateStorefront,
  deleteStorefront,
  fetchStorefrontAnalytics,
  setFilters,
  clearFilters,
  setPagination,
  setCurrentStorefront,
  clearCurrentStorefront,
  clearError,
  selectStorefronts,
  selectCurrentStorefront,
  selectMyStorefronts,
  selectIsLoading,
  selectIsCreating,
  selectIsUpdating,
  selectIsDeleting,
  selectError,
  selectFilters,
  selectPagination,
  selectTotalCount,
  selectHasMore,
  selectAnalytics,
  selectIsLoadingAnalytics,
  selectVerifiedStorefronts,
  selectActiveStorefronts,
  type StorefrontFilters,
  type PaginationParams,
} from '@/store/slices/storefrontSlice';
import type { components } from '@/types/generated/api';

type Storefront =
  components['schemas']['backend_internal_domain_models.Storefront'];
type StorefrontCreateDTO =
  components['schemas']['backend_internal_domain_models.StorefrontCreateDTO'];
type StorefrontUpdateDTO =
  components['schemas']['backend_internal_domain_models.StorefrontUpdateDTO'];

/**
 * Хук для работы с витринами
 * Предоставляет удобный интерфейс для управления состоянием витрин
 */
export const useStorefronts = () => {
  const dispatch = useAppDispatch();

  // Selectors
  const storefronts = useAppSelector(selectStorefronts);
  const currentStorefront = useAppSelector(selectCurrentStorefront);
  const myStorefronts = useAppSelector(selectMyStorefronts);
  const isLoading = useAppSelector(selectIsLoading);
  const isCreating = useAppSelector(selectIsCreating);
  const isUpdating = useAppSelector(selectIsUpdating);
  const isDeleting = useAppSelector(selectIsDeleting);
  const error = useAppSelector(selectError);
  const filters = useAppSelector(selectFilters);
  const pagination = useAppSelector(selectPagination);
  const totalCount = useAppSelector(selectTotalCount);
  const hasMore = useAppSelector(selectHasMore);
  const analytics = useAppSelector(selectAnalytics);
  const isLoadingAnalytics = useAppSelector(selectIsLoadingAnalytics);
  const verifiedStorefronts = useAppSelector(selectVerifiedStorefronts);
  const activeStorefronts = useAppSelector(selectActiveStorefronts);

  // Actions
  const loadStorefronts = useCallback(
    (options?: {
      filters?: StorefrontFilters;
      pagination?: PaginationParams;
    }) => {
      return dispatch(fetchStorefronts(options || {}));
    },
    [dispatch]
  );

  const loadStorefrontById = useCallback(
    (id: number) => {
      return dispatch(fetchStorefrontById(id));
    },
    [dispatch]
  );

  const loadMyStorefronts = useCallback(() => {
    return dispatch(fetchMyStorefronts());
  }, [dispatch]);

  const createNewStorefront = useCallback(
    (data: StorefrontCreateDTO) => {
      return dispatch(createStorefront(data));
    },
    [dispatch]
  );

  const updateExistingStorefront = useCallback(
    (id: number, data: StorefrontUpdateDTO) => {
      return dispatch(updateStorefront({ id, data }));
    },
    [dispatch]
  );

  const deleteExistingStorefront = useCallback(
    (id: number) => {
      return dispatch(deleteStorefront(id));
    },
    [dispatch]
  );

  const loadStorefrontAnalytics = useCallback(
    (id: number, from?: string, to?: string) => {
      return dispatch(fetchStorefrontAnalytics({ id, from, to }));
    },
    [dispatch]
  );

  const updateFilters = useCallback(
    (newFilters: Partial<StorefrontFilters>) => {
      dispatch(setFilters(newFilters));
    },
    [dispatch]
  );

  const resetFilters = useCallback(() => {
    dispatch(clearFilters());
  }, [dispatch]);

  const updatePagination = useCallback(
    (newPagination: Partial<PaginationParams>) => {
      dispatch(setPagination(newPagination));
    },
    [dispatch]
  );

  const selectStorefront = useCallback(
    (storefront: Storefront | null) => {
      dispatch(setCurrentStorefront(storefront));
    },
    [dispatch]
  );

  const clearSelectedStorefront = useCallback(() => {
    dispatch(clearCurrentStorefront());
  }, [dispatch]);

  const clearErrors = useCallback(() => {
    dispatch(clearError());
  }, [dispatch]);

  // Загрузка следующей страницы (для бесконечной прокрутки)
  const loadNextPage = useCallback(() => {
    if (!isLoading && hasMore) {
      const nextOffset = pagination.offset + pagination.limit;
      dispatch(setPagination({ offset: nextOffset }));
      return dispatch(
        fetchStorefronts({
          filters,
          pagination: { ...pagination, offset: nextOffset },
        })
      );
    }
  }, [dispatch, isLoading, hasMore, pagination, filters]);

  // Поиск витрин с применением фильтров
  const searchStorefronts = useCallback(
    (searchFilters: Partial<StorefrontFilters>) => {
      // Обновляем фильтры и сбрасываем пагинацию
      dispatch(setFilters(searchFilters));
      dispatch(setPagination({ offset: 0 }));

      // Выполняем поиск
      return dispatch(
        fetchStorefronts({
          filters: { ...filters, ...searchFilters },
          pagination: { ...pagination, offset: 0 },
        })
      );
    },
    [dispatch, filters, pagination]
  );

  // Поиск витрин по городу
  const searchByCity = useCallback(
    (city: string) => {
      return searchStorefronts({ city });
    },
    [searchStorefronts]
  );

  // Поиск витрин по текстовому запросу
  const searchByText = useCallback(
    (search: string) => {
      return searchStorefronts({ search });
    },
    [searchStorefronts]
  );

  // Поиск витрин по геолокации
  const searchByLocation = useCallback(
    (latitude: number, longitude: number, radiusKm?: number) => {
      return searchStorefronts({ latitude, longitude, radiusKm });
    },
    [searchStorefronts]
  );

  // Получение витрины по слагу из текущего списка
  const getStorefrontBySlug = useCallback(
    (slug: string) => {
      return storefronts.find((storefront) => storefront.slug === slug);
    },
    [storefronts]
  );

  // Получение витрин по городу из текущего списка
  const getStorefrontsByCity = useCallback(
    (city: string) => {
      return storefronts.filter((storefront) => storefront.city === city);
    },
    [storefronts]
  );

  // Проверка, является ли витрина моей
  const isMyStorefront = useCallback(
    (storefrontId: number) => {
      return myStorefronts.some((storefront) => storefront.id === storefrontId);
    },
    [myStorefronts]
  );

  return {
    // State
    storefronts,
    currentStorefront,
    myStorefronts,
    verifiedStorefronts,
    activeStorefronts,
    filters,
    pagination,
    totalCount,
    hasMore,
    analytics,

    // Loading states
    isLoading,
    isCreating,
    isUpdating,
    isDeleting,
    isLoadingAnalytics,
    error,

    // Actions
    loadStorefronts,
    loadStorefrontById,
    loadMyStorefronts,
    createNewStorefront,
    updateExistingStorefront,
    deleteExistingStorefront,
    loadStorefrontAnalytics,
    loadNextPage,

    // Filters and search
    updateFilters,
    resetFilters,
    searchStorefronts,
    searchByCity,
    searchByText,
    searchByLocation,

    // Pagination
    updatePagination,

    // Selection
    selectStorefront,
    clearSelectedStorefront,

    // Utility
    clearErrors,
    getStorefrontBySlug,
    getStorefrontsByCity,
    isMyStorefront,
  };
};

export default useStorefronts;
