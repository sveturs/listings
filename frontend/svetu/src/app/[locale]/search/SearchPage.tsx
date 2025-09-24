'use client';

import { useState, useEffect, useRef, useCallback } from 'react';
import { useTranslations, useLocale } from 'next-intl';
import { useSearchParams, useRouter } from 'next/navigation';
import { SearchBar } from '@/components/SearchBar';
import { UnifiedProductCard } from '@/components/common/UnifiedProductCard';
import { adaptMarketplaceItem } from '@/utils/product-adapters';
import ViewToggle from '@/components/common/ViewToggle';
import { SearchResultCard } from '@/components/search';
import { CategorySelector } from '@/components/search/CategorySelector';
import { DynamicFilters } from '@/components/search/DynamicFilters';
import { QuickFilterPresets } from '@/components/search/QuickFilterPresets';
import { MobileFilterDrawer } from '@/components/search/MobileFilterDrawer';
import { useViewPreference } from '@/hooks/useViewPreference';
import { useInfiniteScroll } from '@/hooks/useInfiniteScroll';
import { useBehaviorTracking } from '@/hooks/useBehaviorTracking';
import {
  useFilterPersistence,
  useSearchHistory,
} from '@/hooks/useFilterPersistence';
import InfiniteScrollTrigger from '@/components/common/InfiniteScrollTrigger';
import { ListingGridSkeleton } from '@/components/ui/skeletons';
import {
  UnifiedSearchService,
  UnifiedSearchResult,
  UnifiedSearchParams,
} from '@/services/unifiedSearch';
import { MarketplaceItem } from '@/types/marketplace';
import { PageTransition } from '@/components/ui/PageTransition';

interface SearchFilters {
  category_id?: string;
  category_ids?: number[];
  price_min?: number;
  price_max?: number;
  product_types?: string[];
  sort_by?: string;
  sort_order?: string;
  city?: string;
  condition?: string;
  distance?: number;
  // Автомобильные фильтры
  car_make?: string;
  car_model?: string;
  car_year_from?: number;
  car_year_to?: number;
  car_mileage_max?: number;
  car_fuel_type?: string;
  car_transmission?: string;
  car_body_type?: string[];
}

export default function SearchPage() {
  const t = useTranslations('search');
  const tCommon = useTranslations('common');
  const locale = useLocale();
  const router = useRouter();
  const searchParams = useSearchParams();

  // Парсим ВСЕ параметры из URL
  const parseSearchParams = () => {
    const params = searchParams || new URLSearchParams();
    return {
      query: params.get('q') || '',
      fuzzy: params.get('fuzzy') !== 'false',
      lat: params.get('lat') ? parseFloat(params.get('lat')!) : undefined,
      lon: params.get('lon') ? parseFloat(params.get('lon')!) : undefined,
      distance: params.get('distance')
        ? Number(params.get('distance'))
        : undefined,
      category_id: params.get('category') || undefined,
      category_ids: (() => {
        const categoriesParam = params.get('categories');
        if (!categoriesParam) return undefined;
        const ids = categoriesParam
          .split(',')
          .map(Number)
          .filter((n) => !isNaN(n) && n > 0);
        return ids.length > 0 ? ids : undefined;
      })(),
      price_min: params.get('price_min')
        ? Number(params.get('price_min'))
        : undefined,
      price_max: params.get('price_max')
        ? Number(params.get('price_max'))
        : undefined,
      product_types: params.get('types')
        ? params.get('types')!.split(',')
        : ['marketplace', 'storefront'],
      sort_by: params.get('sort') || 'relevance',
      sort_order: params.get('order') || 'desc',
      city: params.get('city') || undefined,
      condition: params.get('condition') || undefined,
      // Автомобильные параметры
      car_make: params.get('car_make') || undefined,
      car_model: params.get('car_model') || undefined,
      car_year_from: params.get('car_year_from')
        ? Number(params.get('car_year_from'))
        : undefined,
      car_year_to: params.get('car_year_to')
        ? Number(params.get('car_year_to'))
        : undefined,
      car_mileage_max: params.get('car_mileage_max')
        ? Number(params.get('car_mileage_max'))
        : undefined,
      car_fuel_type: params.get('car_fuel_type') || undefined,
      car_transmission: params.get('car_transmission') || undefined,
      car_body_type: params.get('car_body_type')
        ? params.get('car_body_type')!.split(',')
        : undefined,
    };
  };

  const initialParams = parseSearchParams();

  const [query, setQuery] = useState(initialParams.query);
  const [fuzzy, setFuzzy] = useState(initialParams.fuzzy);
  const [results, setResults] = useState<UnifiedSearchResult | null>(null);
  const [loading, setLoading] = useState(false);
  const [filtersLoading, setFiltersLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [filters, setFilters] = useState<SearchFilters>({
    category_id: initialParams.category_id,
    category_ids: initialParams.category_ids,
    price_min: initialParams.price_min,
    price_max: initialParams.price_max,
    product_types: initialParams.product_types,
    sort_by: initialParams.sort_by,
    sort_order: initialParams.sort_order,
    city: initialParams.city,
    condition: initialParams.condition,
    distance: initialParams.distance,
    // Автомобильные фильтры
    car_make: initialParams.car_make,
    car_model: initialParams.car_model,
    car_year_from: initialParams.car_year_from,
    car_year_to: initialParams.car_year_to,
    car_mileage_max: initialParams.car_mileage_max,
    car_fuel_type: initialParams.car_fuel_type,
    car_transmission: initialParams.car_transmission,
    car_body_type: initialParams.car_body_type,
  });
  const [page, setPage] = useState(1);
  const [allItems, setAllItems] = useState<any[]>([]);
  const [hasInitialSearchRun, setHasInitialSearchRun] = useState(false);
  const [showFilters, setShowFilters] = useState(false);
  const [mobileFiltersOpen, setMobileFiltersOpen] = useState(false);
  const [viewMode, setViewMode] = useViewPreference('grid');
  const [selectedCategoryId, setSelectedCategoryId] = useState<
    number | undefined
  >(
    initialParams.category_id ? parseInt(initialParams.category_id) : undefined
  );
  const [dynamicFilters, setDynamicFilters] = useState<Record<string, any>>({});
  const [geoParams, setGeoParams] = useState({
    latitude: initialParams.lat,
    longitude: initialParams.lon,
    distance: initialParams.distance,
  });

  // Для трекинга времени поиска
  const searchStartTimeRef = useRef<number>(0);

  // Behavior tracking
  const {
    trackSearchPerformed,
    trackSearchFilterApplied,
    trackSearchSortChanged,
  } = useBehaviorTracking();

  // Filter persistence
  const { loadFilters } = useFilterPersistence({
    categoryId: selectedCategoryId,
    filters: { ...filters, ...dynamicFilters },
  });

  // Search history
  const { saveQuery } = useSearchHistory();

  // Load saved filters on mount (only if URL doesn't have filters)
  useEffect(() => {
    if (!initialParams.category_id && Object.keys(initialParams).length <= 2) {
      const saved = loadFilters();
      if (saved) {
        if (saved.categoryId) {
          setSelectedCategoryId(saved.categoryId);
          setFilters((prev) => ({
            ...prev,
            category_id: saved.categoryId?.toString(),
          }));
        }
        if (saved.filters) {
          const { category_id: _categoryId, ...otherFilters } = saved.filters;
          setFilters((prev) => ({ ...prev, ...otherFilters }));
          setDynamicFilters(otherFilters);
        }
      }
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const handleLoadMore = () => {
    if (results && results.has_more) {
      setPage((prev) => prev + 1);
    }
  };

  const loadMoreRef = useInfiniteScroll({
    loading,
    hasMore: results?.has_more || false,
    onLoadMore: handleLoadMore,
  });

  // Функция для обновления URL с новыми параметрами
  const updateURL = (newParams: Partial<typeof initialParams>) => {
    const url = new URL(window.location.href);
    const currentParams = parseSearchParams();
    const mergedParams = { ...currentParams, ...newParams };

    // Очищаем все параметры
    url.search = '';

    // Добавляем только непустые параметры
    if (mergedParams.query) url.searchParams.set('q', mergedParams.query);
    if (!mergedParams.fuzzy) url.searchParams.set('fuzzy', 'false');
    if (mergedParams.lat !== undefined)
      url.searchParams.set('lat', mergedParams.lat.toString());
    if (mergedParams.lon !== undefined)
      url.searchParams.set('lon', mergedParams.lon.toString());
    if (mergedParams.distance !== undefined)
      url.searchParams.set('distance', mergedParams.distance.toString());
    if (mergedParams.category_id)
      url.searchParams.set('category', mergedParams.category_id);
    if (mergedParams.category_ids && mergedParams.category_ids.length > 0) {
      url.searchParams.set('categories', mergedParams.category_ids.join(','));
    }
    if (mergedParams.price_min !== undefined)
      url.searchParams.set('price_min', mergedParams.price_min.toString());
    if (mergedParams.price_max !== undefined)
      url.searchParams.set('price_max', mergedParams.price_max.toString());
    if (
      mergedParams.product_types &&
      mergedParams.product_types.join(',') !== 'marketplace,storefront'
    ) {
      url.searchParams.set('types', mergedParams.product_types.join(','));
    }
    if (mergedParams.sort_by && mergedParams.sort_by !== 'relevance')
      url.searchParams.set('sort', mergedParams.sort_by);
    if (mergedParams.sort_order && mergedParams.sort_order !== 'desc')
      url.searchParams.set('order', mergedParams.sort_order);
    if (mergedParams.city) url.searchParams.set('city', mergedParams.city);
    if (mergedParams.condition)
      url.searchParams.set('condition', mergedParams.condition);

    window.history.replaceState({}, '', url.toString());
  };

  // Handle URL query changes (this handles both initial load and subsequent changes)
  useEffect(() => {
    const params = parseSearchParams();

    // Всегда выполняем поиск при первой загрузке или изменении параметров
    if (!hasInitialSearchRun) {
      setQuery(params.query);
      setFuzzy(params.fuzzy);
      setGeoParams({
        latitude: params.lat,
        longitude: params.lon,
        distance: params.distance,
      });
      setFilters({
        category_id: params.category_id,
        category_ids: params.category_ids,
        price_min: params.price_min,
        price_max: params.price_max,
        product_types: params.product_types,
        sort_by: params.sort_by,
        sort_order: params.sort_order,
        city: params.city,
        condition: params.condition,
        distance: params.distance,
      });
      setPage(1);
      setAllItems([]);
      performSearch(
        params.query,
        1,
        {
          category_id: params.category_id,
          category_ids: params.category_ids,
          price_min: params.price_min,
          price_max: params.price_max,
          product_types: params.product_types,
          sort_by: params.sort_by,
          sort_order: params.sort_order,
          city: params.city,
          condition: params.condition,
          distance: params.distance,
          // Автомобильные параметры
          car_make: params.car_make,
          car_model: params.car_model,
          car_year_from: params.car_year_from,
          car_year_to: params.car_year_to,
          car_mileage_max: params.car_mileage_max,
          car_fuel_type: params.car_fuel_type,
          car_transmission: params.car_transmission,
          car_body_type: params.car_body_type,
        },
        params.fuzzy
      );
      setHasInitialSearchRun(true);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [searchParams]);

  // Load more pages
  useEffect(() => {
    if (page > 1) {
      performSearch(query, page, filters, fuzzy);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page]);

  // Track previous filters to detect real changes
  const prevFiltersRef = useRef<SearchFilters>(filters);

  // Create stable strings for array dependencies
  const categoryIdsKey = filters.category_ids?.join(',') || '';
  const productTypesKey = filters.product_types?.join(',') || '';

  // Handle filter changes (skip first render and only update on real changes)
  useEffect(() => {
    // Check if filters actually changed
    const filtersChanged =
      JSON.stringify(prevFiltersRef.current) !== JSON.stringify(filters);

    if (filtersChanged && hasInitialSearchRun) {
      prevFiltersRef.current = filters;
      setPage(1);
      setAllItems([]);
      performSearch(query, 1, filters, fuzzy);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [
    filters.category_id,
    categoryIdsKey,
    filters.price_min,
    filters.price_max,
    productTypesKey,
    filters.sort_by,
    filters.sort_order,
    filters.city,
    filters.condition,
    filters.distance,
    filters.car_make,
    filters.car_model,
    filters.car_year_from,
    filters.car_year_to,
    filters.car_mileage_max,
    filters.car_fuel_type,
    filters.car_transmission,
    filters.car_body_type?.join(','),
  ]);

  const performSearch = async (
    searchQuery: string,
    currentPage: number,
    currentFilters: SearchFilters,
    useFuzzy: boolean = true
  ) => {
    // Теперь разрешаем пустой запрос для показа всех товаров
    // if (!searchQuery.trim() && !currentFilters.category_id) return;

    // Запоминаем время начала поиска для трекинга
    if (currentPage === 1) {
      searchStartTimeRef.current = Date.now();
      // Сохраняем запрос в историю если он не пустой
      if (searchQuery.trim()) {
        saveQuery(searchQuery);
      }
    }

    setLoading(true);
    setError(null);

    try {
      // Если нет запроса, используем пустую строку
      const effectiveQuery = searchQuery || '';

      const params: UnifiedSearchParams = {
        query: effectiveQuery,
        page: currentPage,
        limit: 20,
        product_types: currentFilters.product_types as (
          | 'marketplace'
          | 'storefront'
        )[],
        sort_by: currentFilters.sort_by as any,
        sort_order: currentFilters.sort_order as any,
        // Передаем массив категорий если есть
        category_ids: currentFilters.category_ids,
        // Передаем единичную категорию только если нет массива
        category_id:
          !currentFilters.category_ids ||
          currentFilters.category_ids.length === 0
            ? currentFilters.category_id
            : undefined,
        price_min: currentFilters.price_min,
        price_max: currentFilters.price_max,
        city: currentFilters.city,
        fuzzy: useFuzzy,
        latitude: geoParams.latitude,
        longitude: geoParams.longitude,
        distance:
          currentFilters.distance?.toString() || geoParams.distance?.toString(),
        // Автомобильные параметры
        car_make: currentFilters.car_make,
        car_model: currentFilters.car_model,
        car_year_from: currentFilters.car_year_from,
        car_year_to: currentFilters.car_year_to,
        car_mileage_max: currentFilters.car_mileage_max,
        car_fuel_type: currentFilters.car_fuel_type,
        car_transmission: currentFilters.car_transmission,
        car_body_type: currentFilters.car_body_type,
      };

      // Debug logging
      console.log('SearchPage - currentFilters:', currentFilters);
      console.log('SearchPage - params being sent:', params);

      const data = await UnifiedSearchService.search(params);
      setResults(data);

      if (currentPage === 1) {
        setAllItems(data.items);

        // Трекинг выполненного поиска (только для первой страницы)
        try {
          // Определяем тип элемента для трекинга на основе фильтров
          const itemType =
            currentFilters.product_types?.length === 1
              ? (currentFilters.product_types[0] as
                  | 'marketplace'
                  | 'storefront')
              : data.items.length > 0
                ? (data.items[0].product_type as 'marketplace' | 'storefront')
                : 'marketplace';

          await trackSearchPerformed({
            search_query: searchQuery,
            search_filters: {
              ...currentFilters,
              fuzzy: useFuzzy,
            },
            search_sort: currentFilters.sort_by,
            results_count: data.total,
            search_duration_ms: Date.now() - searchStartTimeRef.current,
            item_type: itemType,
          });
        } catch (trackingError) {
          console.error('Failed to track search:', trackingError);
        }
      } else {
        setAllItems((prev) => [...prev, ...data.items]);
      }
    } catch (err) {
      console.error('Search error:', err);

      // Fallback: показываем пустые результаты вместо ошибки
      setResults({
        items: [],
        total: 0,
        page: currentPage,
        limit: 20,
        total_pages: 0,
        has_more: false,
        took_ms: 0,
      });

      if (currentPage === 1) {
        setAllItems([]);

        // Трекинг неудачного поиска
        try {
          // Определяем тип элемента для трекинга на основе фильтров
          const itemType =
            currentFilters.product_types?.length === 1
              ? (currentFilters.product_types[0] as
                  | 'marketplace'
                  | 'storefront')
              : 'marketplace';

          await trackSearchPerformed({
            search_query: searchQuery,
            search_filters: {
              ...currentFilters,
              fuzzy: useFuzzy,
            },
            search_sort: currentFilters.sort_by,
            results_count: 0,
            search_duration_ms: Date.now() - searchStartTimeRef.current,
            item_type: itemType,
          });
        } catch (trackingError) {
          console.error('Failed to track failed search:', trackingError);
        }
      }

      setError(null); // Не показываем ошибку пользователю
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = (newQuery: string, newFuzzy?: boolean) => {
    const searchFuzzy = newFuzzy !== undefined ? newFuzzy : fuzzy;
    setQuery(newQuery);
    if (newFuzzy !== undefined) {
      setFuzzy(newFuzzy);
    }
    setPage(1);
    setAllItems([]);
    performSearch(newQuery, 1, filters, searchFuzzy);

    // Обновляем URL со всеми параметрами
    updateURL({
      query: newQuery,
      fuzzy: searchFuzzy,
      ...filters,
    });
  };

  const handleCategorySelect = (categoryId: number | undefined) => {
    setSelectedCategoryId(categoryId);
    setDynamicFilters({});
    handleFilterChange({ category_id: categoryId?.toString() });
  };

  const handleDynamicFiltersChange = useCallback(
    (newFilters: Record<string, any>) => {
      setDynamicFilters((prevDynamicFilters) => {
        // Only update if filters actually changed
        const hasChanges = Object.keys(newFilters).some(
          (key) => prevDynamicFilters[key] !== newFilters[key]
        );

        if (!hasChanges) return prevDynamicFilters;

        // Update filters state as well
        setFilters((prevFilters) => ({ ...prevFilters, ...newFilters }));

        return newFilters;
      });
    },
    []
  );

  const handleFilterChange = async (newFilters: Partial<SearchFilters>) => {
    const prevFilters = filters;
    const updatedFilters = { ...filters, ...newFilters };
    setFilters(updatedFilters);
    setFiltersLoading(true);

    // Трекинг изменения фильтров
    if (query) {
      try {
        // Проверяем какие фильтры изменились
        for (const [key, value] of Object.entries(newFilters)) {
          if (prevFilters[key as keyof SearchFilters] !== value) {
            await trackSearchFilterApplied({
              search_query: query,
              filter_type: key,
              filter_value: JSON.stringify(value),
              results_count_before: results?.total || 0,
              results_count_after: 0, // Будет обновлено после выполнения поиска
            });
          }
        }
      } catch (error) {
        console.error('Failed to track filter change:', error);
      }

      // Выполняем новый поиск с обновленными фильтрами
      await performSearch(query, 1, updatedFilters, fuzzy);
    }

    // Обновляем URL со всеми параметрами
    updateURL({
      query,
      fuzzy,
      ...updatedFilters,
      lat: geoParams.latitude,
      lon: geoParams.longitude,
    });

    setFiltersLoading(false);
  };

  const convertToMarketplaceItem = (item: any): MarketplaceItem => {
    return {
      id: item.product_id,
      title: item.name,
      description: item.description,
      price: item.price,
      location: item.location?.city || '',
      city: item.location?.city || '',
      country: item.location?.country || '',
      images:
        item.images && item.images.length > 0
          ? item.images.map((img: any) => ({
              id: 0,
              public_url: img.url,
              is_main: img.is_main,
            }))
          : [],
      category: item.category,
      user_id: item.user_id || 0,
      condition: 'good',
      status: 'active',
      views_count: 0,
      created_at: item.created_at || new Date().toISOString(),
      updated_at: new Date().toISOString(),
      // Добавляем информацию о типе товара и storefront
      product_type: item.product_type,
      storefront_id: item.storefront?.id,
      // Добавляем данные о скидках
      has_discount: item.has_discount || false,
      old_price: item.old_price,
      discount_percentage: item.discount_percentage,
      currency: item.currency || 'РСД',
    };
  };

  const activeFiltersCount = () => {
    let count = 0;
    if (filters.category_id) count++;
    if (filters.price_min || filters.price_max) count++;
    if (filters.city) count++;
    if (filters.product_types?.length !== 2) count++;
    return count;
  };

  return (
    <PageTransition mode="fade">
      <div className="min-h-screen bg-base-100">
        {/* Компактный хедер с поиском */}
        <div className="bg-base-100 border-b border-base-200 sticky top-0 z-50">
          <div className="container mx-auto px-4 py-4">
            <div className="max-w-2xl mx-auto">
              <SearchBar
                initialQuery={query}
                onSearch={handleSearch}
                variant="minimal"
                showTrending={false}
                fuzzy={fuzzy}
                onFuzzyChange={setFuzzy}
              />

              {/* Быстрые фильтры */}
              <div className="flex flex-wrap gap-2 mt-3">
                <button
                  className={`btn btn-xs sm:btn-sm ${filters.product_types?.includes('marketplace') ? 'btn-primary' : 'btn-ghost'}`}
                  onClick={() => {
                    const types = filters.product_types || [];
                    if (types.includes('marketplace') && types.length > 1) {
                      handleFilterChange({
                        product_types: types.filter((t) => t !== 'marketplace'),
                      });
                    } else if (!types.includes('marketplace')) {
                      handleFilterChange({
                        product_types: [...types, 'marketplace'],
                      });
                    }
                  }}
                >
                  {t('listings')}
                </button>
                <button
                  className={`btn btn-xs sm:btn-sm lg:btn-md ${filters.product_types?.includes('storefront') ? 'btn-primary shadow-lg' : 'btn-ghost hover:btn-primary hover:btn-outline'} transition-all duration-200`}
                  aria-pressed={filters.product_types?.includes('storefront')}
                  aria-label="Filter by storefront items"
                  onClick={() => {
                    const types = filters.product_types || [];
                    if (types.includes('storefront') && types.length > 1) {
                      handleFilterChange({
                        product_types: types.filter((t) => t !== 'storefront'),
                      });
                    } else if (!types.includes('storefront')) {
                      handleFilterChange({
                        product_types: [...types, 'storefront'],
                      });
                    }
                  }}
                >
                  {t('storeProducts')}
                </button>
                <button
                  className={`btn btn-xs sm:btn-sm lg:btn-md ${filters.sort_by === 'price' ? 'btn-primary shadow-lg' : 'btn-ghost hover:btn-primary hover:btn-outline'} transition-all duration-200`}
                  onClick={async () => {
                    const newSortBy =
                      filters.sort_by === 'price' ? 'relevance' : 'price';
                    const previousSort = filters.sort_by;

                    // Трекинг изменения сортировки
                    if (query) {
                      try {
                        await trackSearchSortChanged({
                          search_query: query,
                          sort_type: newSortBy,
                          previous_sort: previousSort,
                          results_count: results?.total || 0,
                        });
                      } catch (error) {
                        console.error('Failed to track sort change:', error);
                      }
                    }

                    handleFilterChange({
                      sort_by: newSortBy,
                    });
                  }}
                >
                  <svg
                    className="w-4 h-4 mr-1"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                  </svg>
                  {t('byPrice')}
                </button>
                <button
                  className={`btn btn-xs sm:btn-sm lg:btn-md ${filters.sort_by === 'date' ? 'btn-primary shadow-lg' : 'btn-ghost hover:btn-primary hover:btn-outline'} transition-all duration-200`}
                  onClick={async () => {
                    const newSortBy =
                      filters.sort_by === 'date' ? 'relevance' : 'date';
                    const previousSort = filters.sort_by;

                    // Трекинг изменения сортировки
                    if (query) {
                      try {
                        await trackSearchSortChanged({
                          search_query: query,
                          sort_type: newSortBy,
                          previous_sort: previousSort,
                          results_count: results?.total || 0,
                        });
                      } catch (error) {
                        console.error('Failed to track sort change:', error);
                      }
                    }

                    handleFilterChange({
                      sort_by: newSortBy,
                    });
                  }}
                >
                  <svg
                    className="w-4 h-4 mr-1"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                    />
                  </svg>
                  {t('byDate')}
                </button>
              </div>
            </div>
          </div>
        </div>

        {/* Результаты поиска */}
        <div className="container mx-auto px-4 py-6">
          {/* Статистика поиска в стиле аналитики */}
          {(query || filters.category_id) && results && (
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
              {/* Найдено результатов */}
              <div className="stat bg-base-100 rounded-xl shadow-md">
                <div className="stat-figure text-primary">
                  <svg
                    className="w-8 h-8"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                    />
                  </svg>
                </div>
                <div className="stat-title text-sm">{t('found')}</div>
                <div className="stat-value text-2xl">{results.total || 0}</div>
                <div className="stat-desc">{t('results')}</div>
              </div>

              {/* Время поиска */}
              <div className="stat bg-base-100 rounded-xl shadow-md">
                <div className="stat-figure text-secondary">
                  <svg
                    className="w-8 h-8"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                  </svg>
                </div>
                <div className="stat-title text-sm">{t('speed')}</div>
                <div className="stat-value text-2xl">
                  {results.took_ms || 0}
                </div>
                <div className="stat-desc">{t('milliseconds')}</div>
              </div>

              {/* Категории */}
              <div className="stat bg-base-100 rounded-xl shadow-md">
                <div className="stat-figure text-accent">
                  <svg
                    className="w-8 h-8"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z"
                    />
                  </svg>
                </div>
                <div className="stat-title text-sm">{t('time')}</div>
                <div className="stat-value text-2xl">
                  {results.took_ms || 0}
                </div>
                <div className="stat-desc">{t('ms')}</div>
              </div>

              {/* Фильтры */}
              <div className="stat bg-base-100 rounded-xl shadow-md">
                <div className="stat-figure text-info">
                  <svg
                    className="w-8 h-8"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4"
                    />
                  </svg>
                </div>
                <div className="stat-title text-sm">{t('activeFilters')}</div>
                <div className="stat-value text-2xl">
                  {activeFiltersCount()}
                </div>
                <div className="stat-desc">{t('applied')}</div>
              </div>
            </div>
          )}

          {(query || filters.category_id) && (
            <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between mb-6 gap-4">
              <div>
                <h2 className="text-2xl font-bold">
                  {query
                    ? `${t('resultsFor')} "${query}"`
                    : t('categoryResults')}
                </h2>
              </div>

              {/* Кнопка фильтров для мобильных */}
              <button
                className="btn btn-outline btn-xs sm:btn-sm lg:hidden"
                onClick={() => setShowFilters(!showFilters)}
              >
                <svg
                  className="w-4 h-4 mr-2"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4"
                  />
                </svg>
                {t('filters')}
                {activeFiltersCount() > 0 && (
                  <span className="badge badge-primary badge-sm ml-2">
                    {activeFiltersCount()}
                  </span>
                )}
              </button>
            </div>
          )}

          <div className="flex flex-col lg:flex-row gap-6">
            {/* Боковая панель с фильтрами - современный дизайн */}
            <aside
              className={`lg:w-80 flex-shrink-0 ${showFilters ? 'block' : 'hidden lg:block'}`}
            >
              <div className="space-y-6">
                {/* Карточка выбора категории */}
                <div className="card bg-base-100 shadow-md">
                  <div className="card-body">
                    <CategorySelector
                      selectedCategoryId={selectedCategoryId}
                      onCategorySelect={handleCategorySelect}
                    />
                  </div>
                </div>

                {/* Карточка быстрых пресетов */}
                <div className="card bg-base-100 shadow-md">
                  <div className="card-body">
                    <QuickFilterPresets
                      onPresetSelect={(presetFilters, categoryId) => {
                        if (categoryId) {
                          handleCategorySelect(categoryId);
                        }
                        handleFilterChange({ ...filters, ...presetFilters });
                      }}
                    />
                  </div>
                </div>

                {/* Карточка фильтров */}
                <div className="card bg-base-100 shadow-md">
                  <div className="card-body">
                    <div className="flex items-center justify-between mb-4">
                      <h3 className="card-title text-lg flex items-center gap-2">
                        <svg
                          className="w-5 h-5 text-primary"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4"
                          />
                        </svg>
                        {t('filters')}
                      </h3>
                      {activeFiltersCount() > 0 && (
                        <button
                          className="btn btn-ghost btn-xs text-error"
                          onClick={() => {
                            const resetFilters = {
                              product_types: ['marketplace', 'storefront'],
                              sort_by: 'relevance',
                              sort_order: 'desc',
                              category_id: undefined,
                              category_ids: undefined,
                              price_min: undefined,
                              price_max: undefined,
                              city: undefined,
                              condition: undefined,
                              distance: undefined,
                            };
                            setFilters(resetFilters);
                            // Очищаем URL
                            updateURL({
                              query,
                              fuzzy,
                              ...resetFilters,
                              lat: geoParams.latitude,
                              lon: geoParams.longitude,
                            });
                          }}
                        >
                          <svg
                            className="w-4 h-4 mr-1"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth={2}
                              d="M6 18L18 6M6 6l12 12"
                            />
                          </svg>
                          {t('reset')}
                        </button>
                      )}
                    </div>

                    <div className="space-y-6">
                      {/* Индикатор загрузки фильтров */}
                      {filtersLoading && (
                        <div className="text-center py-2">
                          <span className="loading loading-spinner loading-sm text-primary"></span>
                          <span className="ml-2 text-sm text-base-content/60">
                            {t('updatingResults')}
                          </span>
                        </div>
                      )}

                      {/* Динамические фильтры */}
                      <DynamicFilters
                        categoryId={selectedCategoryId}
                        onFiltersChange={handleDynamicFiltersChange}
                        activeFilters={dynamicFilters}
                      />
                      <div>
                        <label className="label">
                          <span className="label-text font-medium">
                            {t('productTypes')}
                          </span>
                        </label>
                        <div className="grid grid-cols-2 gap-3">
                          <div
                            className={`card card-compact cursor-pointer transition-all ${
                              filters.product_types?.includes('marketplace')
                                ? 'ring-2 ring-primary bg-primary/5'
                                : 'bg-base-200 hover:bg-base-300'
                            }`}
                            onClick={() => {
                              const types = filters.product_types || [];
                              if (
                                types.includes('marketplace') &&
                                types.length > 1
                              ) {
                                handleFilterChange({
                                  product_types: types.filter(
                                    (t) => t !== 'marketplace'
                                  ),
                                });
                              } else if (!types.includes('marketplace')) {
                                handleFilterChange({
                                  product_types: [...types, 'marketplace'],
                                });
                              }
                            }}
                          >
                            <div className="card-body items-center text-center">
                              <svg
                                className="w-6 h-6 text-primary mb-1"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                              >
                                <path
                                  strokeLinecap="round"
                                  strokeLinejoin="round"
                                  strokeWidth={2}
                                  d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
                                />
                              </svg>
                              <span className="text-xs font-medium">
                                {t('private')}
                              </span>
                            </div>
                          </div>
                          <div
                            className={`card card-compact cursor-pointer transition-all ${
                              filters.product_types?.includes('storefront')
                                ? 'ring-2 ring-primary bg-primary/5'
                                : 'bg-base-200 hover:bg-base-300'
                            }`}
                            onClick={() => {
                              const types = filters.product_types || [];
                              if (
                                types.includes('storefront') &&
                                types.length > 1
                              ) {
                                handleFilterChange({
                                  product_types: types.filter(
                                    (t) => t !== 'storefront'
                                  ),
                                });
                              } else if (!types.includes('storefront')) {
                                handleFilterChange({
                                  product_types: [...types, 'storefront'],
                                });
                              }
                            }}
                          >
                            <div className="card-body items-center text-center">
                              <svg
                                className="w-6 h-6 text-secondary mb-1"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                              >
                                <path
                                  strokeLinecap="round"
                                  strokeLinejoin="round"
                                  strokeWidth={2}
                                  d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
                                />
                              </svg>
                              <span className="text-xs font-medium">
                                {t('stores')}
                              </span>
                            </div>
                          </div>
                        </div>
                      </div>
                      <div>
                        <label className="label">
                          <span className="label-text font-medium flex items-center gap-2">
                            <svg
                              className="w-4 h-4"
                              fill="none"
                              stroke="currentColor"
                              viewBox="0 0 24 24"
                            >
                              <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth={2}
                                d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                              />
                            </svg>
                            {t('priceRange')}
                          </span>
                        </label>
                        <div className="flex gap-2">
                          <div className="form-control flex-1">
                            <label className="input-group">
                              <span className="bg-base-200">{t('from')}</span>
                              <input
                                type="number"
                                className="input input-bordered w-full"
                                value={filters.price_min || ''}
                                onChange={(e) =>
                                  handleFilterChange({
                                    price_min: e.target.value
                                      ? Number(e.target.value)
                                      : undefined,
                                  })
                                }
                              />
                            </label>
                          </div>
                          <div className="form-control flex-1">
                            <label className="input-group">
                              <span className="bg-base-200">{t('to')}</span>
                              <input
                                type="number"
                                className="input input-bordered w-full"
                                value={filters.price_max || ''}
                                onChange={(e) =>
                                  handleFilterChange({
                                    price_max: e.target.value
                                      ? Number(e.target.value)
                                      : undefined,
                                  })
                                }
                              />
                            </label>
                          </div>
                        </div>
                      </div>

                      <div className="divider my-4"></div>

                      {/* Состояние товара */}
                      <div>
                        <label className="label">
                          <span className="label-text font-medium flex items-center gap-2">
                            <svg
                              className="w-4 h-4"
                              fill="none"
                              stroke="currentColor"
                              viewBox="0 0 24 24"
                            >
                              <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth={2}
                                d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                              />
                            </svg>
                            {t('condition')}
                          </span>
                        </label>
                        <select
                          className="select select-bordered w-full bg-base-100"
                          value={filters.condition || ''}
                          onChange={(e) =>
                            handleFilterChange({
                              condition: e.target.value || undefined,
                            })
                          }
                        >
                          <option value="">{t('allConditions')}</option>
                          <option value="new">{t('conditionNew')}</option>
                          <option value="like_new">
                            {t('conditionLikeNew')}
                          </option>
                          <option value="good">{t('conditionGood')}</option>
                          <option value="used">{t('conditionUsed')}</option>
                        </select>
                      </div>

                      <div className="divider my-4"></div>

                      {/* Местоположение */}
                      <div>
                        <label className="label">
                          <span className="label-text font-medium flex items-center gap-2">
                            <svg
                              className="w-4 h-4"
                              fill="none"
                              stroke="currentColor"
                              viewBox="0 0 24 24"
                            >
                              <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth={2}
                                d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"
                              />
                              <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth={2}
                                d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"
                              />
                            </svg>
                            {t('location')}
                          </span>
                        </label>
                        <input
                          type="text"
                          className="input input-bordered w-full"
                          placeholder={t('enterCity')}
                          value={filters.city || ''}
                          onChange={(e) =>
                            handleFilterChange({
                              city: e.target.value || undefined,
                            })
                          }
                        />

                        {/* Радиус поиска */}
                        <label className="label mt-2">
                          <span className="label-text text-sm">
                            {t('searchRadius')}
                          </span>
                        </label>
                        <div className="flex items-center gap-2">
                          <input
                            type="range"
                            min="0"
                            max="100"
                            value={filters.distance || 0}
                            className="range range-primary range-sm flex-1"
                            onChange={(e) =>
                              handleFilterChange({
                                distance: Number(e.target.value) || undefined,
                              })
                            }
                          />
                          <span className="text-sm font-medium w-12">
                            {filters.distance || 0} km
                          </span>
                        </div>
                      </div>

                      <div className="divider my-4"></div>

                      {/* Сортировка - современный селект */}
                      <div>
                        <label className="label">
                          <span className="label-text font-medium flex items-center gap-2">
                            <svg
                              className="w-4 h-4"
                              fill="none"
                              stroke="currentColor"
                              viewBox="0 0 24 24"
                            >
                              <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth={2}
                                d="M3 4h13M3 8h9m-9 4h6m4 0l4-4m0 0l4 4m-4-4v12"
                              />
                            </svg>
                            {t('sortBy')}
                          </span>
                        </label>
                        <select
                          className="select select-bordered w-full bg-base-100"
                          value={filters.sort_by}
                          onChange={async (e) => {
                            const newSortBy = e.target.value;
                            const previousSort = filters.sort_by;

                            // Трекинг изменения сортировки
                            if (query) {
                              try {
                                await trackSearchSortChanged({
                                  search_query: query,
                                  sort_type: newSortBy,
                                  previous_sort: previousSort,
                                  results_count: results?.total || 0,
                                });
                              } catch (error) {
                                console.error(
                                  'Failed to track sort change:',
                                  error
                                );
                              }
                            }

                            handleFilterChange({ sort_by: newSortBy });
                          }}
                        >
                          <option value="relevance">{t('relevance')}</option>
                          <option value="price">{t('price')}</option>
                          <option value="date">{t('date')}</option>
                          <option value="popularity">{t('popularity')}</option>
                        </select>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </aside>

            {/* Основной контент */}
            <main className="flex-1">
              {loading && allItems.length === 0 && (
                <div className="py-4">
                  <ListingGridSkeleton count={8} viewMode={viewMode} />
                </div>
              )}

              {error && (
                <div role="alert" className="alert alert-error mb-6">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    className="stroke-current shrink-0 h-6 w-6"
                    fill="none"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                  </svg>
                  <span>{error}</span>
                </div>
              )}

              {!loading &&
                allItems.length === 0 &&
                (query || filters.category_id) && (
                  <div className="card bg-base-100 shadow-xl mx-auto max-w-2xl">
                    <div className="card-body text-center py-16">
                      <div className="inline-flex items-center justify-center w-32 h-32 rounded-full bg-gradient-to-br from-primary/20 to-secondary/20 mb-8 mx-auto">
                        <svg
                          className="w-16 h-16 text-primary"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                          />
                        </svg>
                      </div>
                      <h3 className="text-3xl font-bold mb-4">
                        {t('noResults')}
                      </h3>
                      <p className="text-base-content/60 mb-8 max-w-md mx-auto">
                        {t('noResultsDescription')}
                      </p>

                      <div className="divider">
                        {locale === 'ru' ? 'Попробуйте' : 'Try'}
                      </div>

                      <div className="flex flex-wrap gap-2 justify-center mb-8">
                        <button className="badge badge-lg badge-outline hover:badge-primary cursor-pointer">
                          {locale === 'ru' ? 'Электроника' : 'Electronics'}
                        </button>
                        <button className="badge badge-lg badge-outline hover:badge-primary cursor-pointer">
                          {locale === 'ru' ? 'Одежда' : 'Clothing'}
                        </button>
                        <button className="badge badge-lg badge-outline hover:badge-primary cursor-pointer">
                          {locale === 'ru' ? 'Дом и сад' : 'Home & Garden'}
                        </button>
                        <button className="badge badge-lg badge-outline hover:badge-primary cursor-pointer">
                          {locale === 'ru' ? 'Авто' : 'Auto'}
                        </button>
                      </div>

                      <div className="flex flex-col sm:flex-row gap-3 justify-center">
                        <button
                          className="btn btn-outline btn-primary"
                          onClick={() => {
                            const resetFilters = {
                              product_types: ['marketplace', 'storefront'],
                              sort_by: 'relevance',
                              sort_order: 'desc',
                              category_id: undefined,
                              category_ids: undefined,
                              price_min: undefined,
                              price_max: undefined,
                              city: undefined,
                              condition: undefined,
                              distance: undefined,
                            };
                            setFilters(resetFilters);
                            // Очищаем URL
                            updateURL({
                              query,
                              fuzzy,
                              ...resetFilters,
                              lat: geoParams.latitude,
                              lon: geoParams.longitude,
                            });
                          }}
                        >
                          <svg
                            className="w-4 h-4 mr-2"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth={2}
                              d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
                            />
                          </svg>
                          {t('clearFilters')}
                        </button>
                        <button
                          className="btn btn-primary"
                          onClick={() => {
                            setQuery('');
                            setAllItems([]);
                            setResults(null);
                            router.push(`/${locale}`);
                          }}
                        >
                          <svg
                            className="w-4 h-4 mr-2"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth={2}
                              d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"
                            />
                          </svg>
                          {locale === 'ru' ? 'На главную' : 'Home'}
                        </button>
                      </div>
                    </div>
                  </div>
                )}

              {allItems.length > 0 && (
                <>
                  <div className="flex justify-end mb-4">
                    <ViewToggle
                      currentView={viewMode}
                      onViewChange={setViewMode}
                    />
                  </div>
                  <div
                    className={
                      viewMode === 'grid'
                        ? 'grid gap-6 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4'
                        : 'space-y-4'
                    }
                  >
                    {allItems.map((item, index) => (
                      <SearchResultCard
                        key={`${item.id}-${index}`}
                        searchQuery={query}
                        itemId={item.product_id?.toString() || item.id}
                        position={index + 1}
                        totalResults={allItems.length}
                        searchStartTime={searchStartTimeRef.current}
                        productType={
                          item.product_type as 'marketplace' | 'storefront'
                        }
                      >
                        <UnifiedProductCard
                          product={adaptMarketplaceItem(
                            convertToMarketplaceItem(item)
                          )}
                          locale={locale}
                          viewMode={viewMode}
                          index={index}
                        />
                      </SearchResultCard>
                    ))}
                  </div>

                  {/* Показываем скелетоны при подгрузке */}
                  {loading && allItems.length > 0 && (
                    <div className="mt-6">
                      <ListingGridSkeleton count={4} viewMode={viewMode} />
                    </div>
                  )}

                  <InfiniteScrollTrigger
                    ref={loadMoreRef}
                    loading={loading}
                    hasMore={results?.has_more || false}
                    onLoadMore={handleLoadMore}
                    loadMoreText={tCommon('loadMore')}
                  />
                </>
              )}
            </main>
          </div>
        </div>

        {/* Mobile Filter Drawer */}
        <MobileFilterDrawer
          isOpen={mobileFiltersOpen}
          onClose={() => setMobileFiltersOpen(false)}
          selectedCategoryId={selectedCategoryId}
          onCategorySelect={handleCategorySelect}
          filters={{ ...filters, ...dynamicFilters }}
          onFiltersChange={(newFilters) => {
            handleFilterChange(newFilters);
            setDynamicFilters(newFilters);
          }}
          activeFiltersCount={activeFiltersCount()}
        />

        {/* Floating Filter Button for Mobile */}
        <button
          onClick={() => setMobileFiltersOpen(true)}
          className={`
            fixed bottom-6 right-6 z-30
            btn btn-circle btn-primary btn-lg shadow-xl
            lg:hidden
            ${mobileFiltersOpen ? 'scale-0' : 'scale-100'}
            transition-transform duration-200
          `}
        >
          <svg
            className="w-6 h-6"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4"
            />
          </svg>
          {activeFiltersCount() > 0 && (
            <span className="badge badge-secondary badge-sm absolute -top-2 -right-2">
              {activeFiltersCount()}
            </span>
          )}
        </button>
      </div>
    </PageTransition>
  );
}
