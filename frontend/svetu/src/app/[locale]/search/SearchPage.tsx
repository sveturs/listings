'use client';

import { useState, useEffect, useRef } from 'react';
import { useTranslations, useLocale } from 'next-intl';
import { useSearchParams, useRouter } from 'next/navigation';
import { SearchBar } from '@/components/SearchBar';
import MarketplaceCard from '@/components/MarketplaceCard';
import ViewToggle from '@/components/common/ViewToggle';
import { SearchResultCard } from '@/components/search';
import { useViewPreference } from '@/hooks/useViewPreference';
import { useInfiniteScroll } from '@/hooks/useInfiniteScroll';
import { useBehaviorTracking } from '@/hooks/useBehaviorTracking';
import InfiniteScrollTrigger from '@/components/common/InfiniteScrollTrigger';
import {
  UnifiedSearchService,
  UnifiedSearchResult,
  UnifiedSearchParams,
} from '@/services/unifiedSearch';
import { MarketplaceItem } from '@/types/marketplace';

interface SearchFilters {
  category_id?: string;
  price_min?: number;
  price_max?: number;
  product_types?: string[];
  sort_by?: string;
  sort_order?: string;
  city?: string;
}

export default function SearchPage() {
  const t = useTranslations('search');
  const tCommon = useTranslations('common');
  const locale = useLocale();
  const router = useRouter();
  const searchParams = useSearchParams();
  const initialQuery = searchParams.get('q') || '';
  const initialFuzzy = searchParams.get('fuzzy') !== 'false'; // По умолчанию true

  const [query, setQuery] = useState(initialQuery);
  const [fuzzy, setFuzzy] = useState(initialFuzzy);
  const [results, setResults] = useState<UnifiedSearchResult | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [filters, setFilters] = useState<SearchFilters>({
    product_types: ['marketplace', 'storefront'],
    sort_by: 'relevance',
    sort_order: 'desc',
  });
  const [page, setPage] = useState(1);
  const [allItems, setAllItems] = useState<any[]>([]);
  const [showFilters, setShowFilters] = useState(false);
  const [viewMode, setViewMode] = useViewPreference('grid');

  // Для трекинга времени поиска
  const searchStartTimeRef = useRef<number>(0);

  // Behavior tracking
  const {
    trackSearchPerformed,
    trackSearchFilterApplied,
    trackSearchSortChanged,
  } = useBehaviorTracking();

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

  // Handle URL query changes (this handles both initial load and subsequent changes)
  useEffect(() => {
    const searchQuery = searchParams.get('q');
    const searchFuzzy = searchParams.get('fuzzy') !== 'false';
    if (searchQuery) {
      // Only perform search if query actually changed
      if (searchQuery !== query || searchFuzzy !== fuzzy) {
        setQuery(searchQuery);
        setFuzzy(searchFuzzy);
        setPage(1);
        setAllItems([]);
        performSearch(searchQuery, 1, filters, searchFuzzy);
      }
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [searchParams]);

  // Load more pages
  useEffect(() => {
    if (query && page > 1) {
      performSearch(query, page, filters, fuzzy);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page]);

  // Handle filter changes (skip first render)
  useEffect(() => {
    // Skip effect on mount
    const isMount = allItems.length === 0 && !results;
    if (query && !isMount) {
      setPage(1);
      setAllItems([]);
      performSearch(query, 1, filters, fuzzy);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [filters]);

  const performSearch = async (
    searchQuery: string,
    currentPage: number,
    currentFilters: SearchFilters,
    useFuzzy: boolean = true
  ) => {
    if (!searchQuery.trim()) return;

    // Запоминаем время начала поиска для трекинга
    if (currentPage === 1) {
      searchStartTimeRef.current = Date.now();
    }

    setLoading(true);
    setError(null);

    try {
      const params: UnifiedSearchParams = {
        query: searchQuery,
        page: currentPage,
        limit: 20,
        product_types: currentFilters.product_types as (
          | 'marketplace'
          | 'storefront'
        )[],
        sort_by: currentFilters.sort_by as any,
        sort_order: currentFilters.sort_order as any,
        category_id: currentFilters.category_id,
        price_min: currentFilters.price_min,
        price_max: currentFilters.price_max,
        city: currentFilters.city,
        fuzzy: useFuzzy,
      };

      const data = await UnifiedSearchService.search(params);
      setResults(data);

      if (currentPage === 1) {
        setAllItems(data.items);

        // Трекинг выполненного поиска (только для первой страницы)
        try {
          await trackSearchPerformed({
            search_query: searchQuery,
            search_filters: {
              ...currentFilters,
              fuzzy: useFuzzy,
            },
            search_sort: currentFilters.sort_by,
            results_count: data.total,
            search_duration_ms: Date.now() - searchStartTimeRef.current,
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
          await trackSearchPerformed({
            search_query: searchQuery,
            search_filters: {
              ...currentFilters,
              fuzzy: useFuzzy,
            },
            search_sort: currentFilters.sort_by,
            results_count: 0,
            search_duration_ms: Date.now() - searchStartTimeRef.current,
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

    // Обновляем URL
    const url = new URL(window.location.href);
    url.searchParams.set('q', newQuery);
    url.searchParams.set('fuzzy', searchFuzzy.toString());
    window.history.replaceState({}, '', url.toString());
  };

  const handleFilterChange = async (newFilters: Partial<SearchFilters>) => {
    const prevFilters = filters;
    setFilters((prev) => ({ ...prev, ...newFilters }));

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
    }
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
              listing_id: item.product_id,
              file_path: '',
              file_name: img.alt_text || '',
              file_size: 0,
              content_type: 'image/jpeg',
              is_main: img.is_main,
              storage_type: 'minio',
              storage_bucket: '',
              public_url: img.public_url || img.url,
              created_at: new Date().toISOString(),
            }))
          : [],
      category: item.category,
      user_id: item.user_id || 0,
      condition: 'good',
      status: 'active',
      views_count: 0,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
      // Добавляем информацию о типе товара и storefront
      product_type: item.product_type,
      storefront_id: item.storefront?.id,
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
                {t('search.listings')}
              </button>
              <button
                className={`btn btn-xs sm:btn-sm lg:btn-md ${filters.product_types?.includes('storefront') ? 'btn-primary shadow-lg' : 'btn-ghost hover:btn-primary hover:btn-outline'} transition-all duration-200`}
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
                {t('search.storeProducts')}
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
                {t('search.byPrice')}
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
                {t('search.byDate')}
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Результаты поиска */}
      <div className="container mx-auto px-4 py-6">
        {/* Статистика поиска в стиле аналитики */}
        {query && results && (
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
              <div className="stat-title text-sm">{t('search.found')}</div>
              <div className="stat-value text-2xl">{results.total || 0}</div>
              <div className="stat-desc">{t('search.results')}</div>
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
              <div className="stat-title text-sm">{t('search.speed')}</div>
              <div className="stat-value text-2xl">{results.took_ms || 0}</div>
              <div className="stat-desc">{t('search.milliseconds')}</div>
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
              <div className="stat-title text-sm">{t('search.time')}</div>
              <div className="stat-value text-2xl">{results.took_ms || 0}</div>
              <div className="stat-desc">{t('search.ms')}</div>
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
              <div className="stat-title text-sm">
                {t('search.activeFilters')}
              </div>
              <div className="stat-value text-2xl">{activeFiltersCount()}</div>
              <div className="stat-desc">{t('search.applied')}</div>
            </div>
          </div>
        )}

        {query && (
          <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between mb-6 gap-4">
            <div>
              <h2 className="text-2xl font-bold">
                {t('search.resultsFor')} &quot;{query}&quot;
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
              {t('search.filters')}
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
                          setFilters({
                            product_types: ['marketplace', 'storefront'],
                            sort_by: 'relevance',
                            sort_order: 'desc',
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
                        {t('search.reset')}
                      </button>
                    )}
                  </div>

                  <div className="space-y-6">
                    {/* Тип товаров - карточки вместо чекбоксов */}
                    <div>
                      <label className="label">
                        <span className="label-text font-medium">
                          {t('search.productTypes')}
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
                              {t('search.private')}
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
                              {t('search.stores')}
                            </span>
                          </div>
                        </div>
                      </div>
                    </div>

                    <div className="divider my-4"></div>

                    {/* Диапазон цен - улучшенный дизайн */}
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
                            <span className="bg-base-200">
                              {t('search.from')}
                            </span>
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
                            <span className="bg-base-200">
                              {t('search.to')}
                            </span>
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
              <div className="flex justify-center py-16">
                <span className="loading loading-spinner loading-lg text-primary"></span>
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

            {!loading && allItems.length === 0 && query && (
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
                    {t('search.noResults')}
                  </h3>
                  <p className="text-base-content/60 mb-8 max-w-md mx-auto">
                    {t('search.noResultsDescription')}
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
                        setFilters({
                          product_types: ['marketplace', 'storefront'],
                          sort_by: 'relevance',
                          sort_order: 'desc',
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
                      {t('search.clearFilters')}
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
                    >
                      <MarketplaceCard
                        item={convertToMarketplaceItem(item)}
                        locale={locale}
                        viewMode={viewMode}
                      />
                    </SearchResultCard>
                  ))}
                </div>

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
    </div>
  );
}
