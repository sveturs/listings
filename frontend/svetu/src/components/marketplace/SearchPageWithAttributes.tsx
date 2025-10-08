'use client';

import React, { useState, useEffect, useCallback, Suspense, lazy } from 'react';
import { useTranslations } from 'next-intl';
import { useSearchWithAttributes } from '@/hooks/useSearchWithAttributes';
import MarketplaceList from './MarketplaceList';
import {
  UnifiedSearchItem,
  UnifiedSearchService,
} from '@/services/unifiedSearch';
import { countActiveAttributeFilters } from '@/utils/urlAttributeSync';
import { useSearchParams } from 'next/navigation';
import {
  MagnifyingGlassIcon,
  AdjustmentsHorizontalIcon,
  XMarkIcon,
  FunnelIcon,
} from '@heroicons/react/24/outline';
import type { components } from '@/types/generated/api';

// Lazy load SmartAttributeFilters for better performance
const SmartAttributeFilters = lazy(
  () => import('@/components/shared/SmartAttributeFilters')
);

type UnifiedAttribute = components['schemas']['models.UnifiedAttribute'];
type Category = components['schemas']['models.MarketplaceCategory'];

interface SearchPageWithAttributesProps {
  initialData: {
    items: UnifiedSearchItem[];
    total: number;
    page: number;
    limit: number;
    has_more: boolean;
  } | null;
  categories?: Category[];
  locale: string;
  error?: Error | null;
}

export default function SearchPageWithAttributes({
  initialData,
  categories = [],
  locale,
}: SearchPageWithAttributesProps) {
  const t = useTranslations('marketplace');
  const searchParams = useSearchParams();

  // State management
  const [isLoading, setIsLoading] = useState(false);
  const [searchResults, setSearchResults] = useState(initialData);
  const [showFilters, setShowFilters] = useState(false);
  const [availableAttributes, setAvailableAttributes] = useState<
    UnifiedAttribute[]
  >([]);
  const [mobileFiltersOpen, setMobileFiltersOpen] = useState(false);

  // Initialize search with attributes hook
  const {
    filters,
    attributeFilters,
    updateFilter,
    updateAttributeFilter,
    clearAttributeFilters,
  } = useSearchWithAttributes(
    {
      query: searchParams.get('q') || '',
      categoryId: searchParams.get('category')
        ? parseInt(searchParams.get('category') || '0')
        : undefined,
    },
    {
      syncWithURL: true,
      autoApply: true,
      debounceDelay: 300,
    }
  );

  // Count active filters
  const activeFilterCount = countActiveAttributeFilters(searchParams);

  // Load attributes for selected category
  useEffect(() => {
    const loadCategoryAttributes = async () => {
      if (!filters.categoryId) {
        setAvailableAttributes([]);
        return;
      }

      try {
        const response = await fetch(
          `/api/v1/unified-attributes/category/${filters.categoryId}`
        );
        if (response.ok) {
          const data = await response.json();
          setAvailableAttributes(data.data || []);
        }
      } catch (err) {
        console.error('Failed to load category attributes:', err);
      }
    };

    loadCategoryAttributes();
  }, [filters.categoryId]);

  // Perform search with filters
  const performSearch = useCallback(async () => {
    setIsLoading(true);

    try {
      // Build search parameters with attributes
      const searchParams: any = {
        query: filters.query || '',
        page: filters.page || 1,
        limit: filters.limit || 20,
        sort: filters.sortBy || 'relevance',
      };

      if (filters.categoryId) {
        searchParams.category_id = filters.categoryId;
      }

      if (filters.minPrice) {
        searchParams.min_price = filters.minPrice;
      }

      if (filters.maxPrice) {
        searchParams.max_price = filters.maxPrice;
      }

      // Add attribute filters
      if (filters.attributes && Object.keys(filters.attributes).length > 0) {
        searchParams.attributes = filters.attributes;
      }

      // Call unified search service
      const results = await UnifiedSearchService.search(searchParams);
      setSearchResults(results);
    } catch (err) {
      console.error('Search failed:', err);
      setSearchResults(null);
    } finally {
      setIsLoading(false);
    }
  }, [filters]);

  // Trigger search when filters change
  useEffect(() => {
    performSearch();
  }, [filters, performSearch]);

  // Clear all filters
  const handleClearAllFilters = () => {
    clearAttributeFilters();
    updateFilter('query', '');
    updateFilter('categoryId', undefined);
    updateFilter('minPrice', undefined);
    updateFilter('maxPrice', undefined);
  };

  return (
    <div className="min-h-screen bg-base-100">
      {/* Search Header */}
      <div className="sticky top-0 z-40 bg-base-100 border-b border-base-200 shadow-sm">
        <div className="container mx-auto px-4 py-4">
          <div className="flex flex-col md:flex-row gap-4">
            {/* Search Input */}
            <div className="flex-1">
              <div className="relative">
                <input
                  type="text"
                  placeholder={t('search.placeholder')}
                  value={filters.query || ''}
                  onChange={(e) => updateFilter('query', e.target.value)}
                  className="input input-bordered w-full pr-10"
                />
                <MagnifyingGlassIcon className="absolute right-3 top-1/2 -translate-y-1/2 w-5 h-5 text-base-content/50" />
              </div>
            </div>

            {/* Category Selector */}
            <select
              value={filters.categoryId || ''}
              onChange={(e) =>
                updateFilter(
                  'categoryId',
                  e.target.value ? parseInt(e.target.value) : undefined
                )
              }
              className="select select-bordered w-full md:w-auto"
            >
              <option value="">{t('search.allCategories')}</option>
              {categories.map((cat) => (
                <option key={cat.id} value={cat.id}>
                  {cat.name}
                </option>
              ))}
            </select>

            {/* Filters Toggle */}
            <button
              onClick={() => setShowFilters(!showFilters)}
              className={`btn ${showFilters ? 'btn-primary' : 'btn-outline'} gap-2`}
            >
              <AdjustmentsHorizontalIcon className="w-5 h-5" />
              <span className="hidden sm:inline">{t('search.filters')}</span>
              {activeFilterCount > 0 && (
                <span className="badge badge-sm">{activeFilterCount}</span>
              )}
            </button>

            {/* Mobile Filters Button */}
            <button
              onClick={() => setMobileFiltersOpen(true)}
              className="btn btn-outline md:hidden gap-2"
            >
              <FunnelIcon className="w-5 h-5" />
              {activeFilterCount > 0 && (
                <span className="badge badge-sm badge-primary">
                  {activeFilterCount}
                </span>
              )}
            </button>
          </div>

          {/* Active Filter Pills */}
          {activeFilterCount > 0 && (
            <div className="flex flex-wrap gap-2 mt-4">
              <span className="text-sm text-base-content/70">
                {t('search.activeFilters')}:
              </span>
              {Object.entries(attributeFilters).map(([attrId, values]) =>
                values.map((value) => (
                  <div
                    key={`${attrId}-${value}`}
                    className="badge badge-lg gap-1"
                  >
                    <span>{value}</span>
                    <button
                      onClick={() => {
                        const newValues = values.filter((v) => v !== value);
                        updateAttributeFilter(attrId, newValues as any);
                      }}
                      className="btn btn-ghost btn-xs btn-circle"
                    >
                      <XMarkIcon className="w-3 h-3" />
                    </button>
                  </div>
                ))
              )}
              <button
                onClick={handleClearAllFilters}
                className="btn btn-ghost btn-xs"
              >
                {t('search.clearAll')}
              </button>
            </div>
          )}
        </div>
      </div>

      <div className="container mx-auto px-4 py-6">
        <div className="flex flex-col lg:flex-row gap-6">
          {/* Desktop Filters Sidebar */}
          {showFilters && (
            <aside className="hidden lg:block lg:w-80 space-y-4">
              {/* Price Range */}
              <div className="card bg-base-200">
                <div className="card-body">
                  <h3 className="card-title text-sm">
                    {t('filters.priceRange')}
                  </h3>
                  <div className="flex gap-2">
                    <input
                      type="number"
                      placeholder={t('filters.min')}
                      value={filters.minPrice || ''}
                      onChange={(e) =>
                        updateFilter(
                          'minPrice',
                          e.target.value
                            ? parseFloat(e.target.value)
                            : undefined
                        )
                      }
                      className="input input-sm input-bordered w-full"
                    />
                    <span className="self-center">-</span>
                    <input
                      type="number"
                      placeholder={t('filters.max')}
                      value={filters.maxPrice || ''}
                      onChange={(e) =>
                        updateFilter(
                          'maxPrice',
                          e.target.value
                            ? parseFloat(e.target.value)
                            : undefined
                        )
                      }
                      className="input input-sm input-bordered w-full"
                    />
                  </div>
                </div>
              </div>

              {/* Smart Attribute Filters */}
              {availableAttributes.length > 0 && (
                <Suspense
                  fallback={
                    <div className="card bg-base-200 animate-pulse">
                      <div className="card-body">
                        <div className="h-4 bg-base-300 rounded w-1/3 mb-2"></div>
                        <div className="space-y-2">
                          <div className="h-8 bg-base-300 rounded"></div>
                          <div className="h-8 bg-base-300 rounded"></div>
                        </div>
                      </div>
                    </div>
                  }
                >
                  <SmartAttributeFilters
                    categoryId={filters.categoryId}
                    initialFilters={{}}
                    onFiltersChange={(filters) => {
                      Object.entries(filters).forEach(([key, value]) => {
                        if (value && value.text_value) {
                          updateAttributeFilter(key, [value.text_value]);
                        }
                      });
                    }}
                    className=""
                  />
                </Suspense>
              )}
            </aside>
          )}

          {/* Main Content */}
          <div className="flex-1">
            {/* Results Header */}
            <div className="flex justify-between items-center mb-4">
              <div>
                {searchResults && (
                  <p className="text-base-content/70">
                    {t('search.resultsCount', { count: searchResults.total })}
                  </p>
                )}
              </div>

              {/* Sort Options */}
              <select
                value={filters.sortBy || 'relevance'}
                onChange={(e) => updateFilter('sortBy', e.target.value)}
                className="select select-bordered select-sm"
              >
                <option value="relevance">{t('sort.relevance')}</option>
                <option value="price_asc">{t('sort.priceAsc')}</option>
                <option value="price_desc">{t('sort.priceDesc')}</option>
                <option value="date_desc">{t('sort.dateDesc')}</option>
              </select>
            </div>

            {/* Results List */}
            {isLoading ? (
              <div className="flex justify-center py-12">
                <span className="loading loading-spinner loading-lg"></span>
              </div>
            ) : searchResults && searchResults.items.length > 0 ? (
              <MarketplaceList initialData={searchResults} locale={locale} />
            ) : (
              <div className="text-center py-12">
                <p className="text-base-content/70">{t('search.noResults')}</p>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Mobile Filters Modal */}
      {mobileFiltersOpen && (
        <div className="fixed inset-0 z-50 lg:hidden">
          <div
            className="fixed inset-0 bg-black/50"
            onClick={() => setMobileFiltersOpen(false)}
          />
          <div className="fixed inset-y-0 right-0 w-full max-w-sm bg-base-100 shadow-xl overflow-y-auto">
            <div className="sticky top-0 bg-base-100 border-b border-base-200 p-4">
              <div className="flex justify-between items-center">
                <h2 className="text-lg font-semibold">{t('filters.title')}</h2>
                <button
                  onClick={() => setMobileFiltersOpen(false)}
                  className="btn btn-ghost btn-circle btn-sm"
                >
                  <XMarkIcon className="w-5 h-5" />
                </button>
              </div>
            </div>

            <div className="p-4 space-y-4">
              {/* Mobile filters content */}
              {/* Price Range */}
              <div>
                <h3 className="font-medium mb-2">{t('filters.priceRange')}</h3>
                <div className="flex gap-2">
                  <input
                    type="number"
                    placeholder={t('filters.min')}
                    value={filters.minPrice || ''}
                    onChange={(e) =>
                      updateFilter(
                        'minPrice',
                        e.target.value ? parseFloat(e.target.value) : undefined
                      )
                    }
                    className="input input-bordered w-full"
                  />
                  <input
                    type="number"
                    placeholder={t('filters.max')}
                    value={filters.maxPrice || ''}
                    onChange={(e) =>
                      updateFilter(
                        'maxPrice',
                        e.target.value ? parseFloat(e.target.value) : undefined
                      )
                    }
                    className="input input-bordered w-full"
                  />
                </div>
              </div>

              {/* Attributes */}
              {availableAttributes.length > 0 && (
                <Suspense
                  fallback={<div className="loading loading-spinner" />}
                >
                  <SmartAttributeFilters
                    categoryId={filters.categoryId}
                    initialFilters={{}}
                    onFiltersChange={(filters) => {
                      Object.entries(filters).forEach(([key, value]) => {
                        if (value && value.text_value) {
                          updateAttributeFilter(key, [value.text_value]);
                        }
                      });
                    }}
                    className=""
                  />
                </Suspense>
              )}

              {/* Apply Button */}
              <div className="sticky bottom-0 bg-base-100 pt-4 border-t border-base-200">
                <button
                  onClick={() => setMobileFiltersOpen(false)}
                  className="btn btn-primary btn-block"
                >
                  {t('filters.apply')}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
