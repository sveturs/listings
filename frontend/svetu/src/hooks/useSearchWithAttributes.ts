import { useState, useCallback, useEffect, useMemo } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { debounce } from 'lodash';
import type { components } from '@/types/generated/api';

type _UnifiedAttribute =
  components['schemas']['backend_internal_domain_models.UnifiedAttribute'];
type _AttributeValue =
  components['schemas']['backend_internal_domain_models.UnifiedAttributeValue'];

interface SearchFilters {
  query?: string;
  categoryId?: number;
  minPrice?: number;
  maxPrice?: number;
  attributes?: Record<string, string[]>;
  sortBy?: 'price_asc' | 'price_desc' | 'date_desc' | 'relevance';
  page?: number;
  limit?: number;
}

interface UseSearchWithAttributesReturn {
  filters: SearchFilters;
  attributeFilters: Record<string, string[]>;
  isLoading: boolean;
  updateFilter: (key: keyof SearchFilters, value: any) => void;
  updateAttributeFilter: (attributeId: string, values: string[]) => void;
  clearAttributeFilters: () => void;
  applyFilters: () => void;
  syncWithURL: boolean;
}

const DEBOUNCE_DELAY = 300;
const DEFAULT_LIMIT = 20;

export function useSearchWithAttributes(
  initialFilters?: Partial<SearchFilters>,
  options?: {
    syncWithURL?: boolean;
    autoApply?: boolean;
    debounceDelay?: number;
  }
): UseSearchWithAttributesReturn {
  const router = useRouter();
  const searchParams = useSearchParams();

  const {
    syncWithURL = true,
    autoApply = true,
    debounceDelay = DEBOUNCE_DELAY,
  } = options || {};

  // Parse initial filters from URL if syncWithURL is enabled
  const getInitialFilters = useCallback((): SearchFilters => {
    if (!syncWithURL) {
      return {
        query: initialFilters?.query || '',
        categoryId: initialFilters?.categoryId,
        minPrice: initialFilters?.minPrice,
        maxPrice: initialFilters?.maxPrice,
        attributes: initialFilters?.attributes || {},
        sortBy: initialFilters?.sortBy || 'relevance',
        page: initialFilters?.page || 1,
        limit: initialFilters?.limit || DEFAULT_LIMIT,
      };
    }

    const query = searchParams.get('q') || initialFilters?.query || '';
    const categoryId = searchParams.get('category')
      ? parseInt(searchParams.get('category')!)
      : initialFilters?.categoryId;
    const minPrice = searchParams.get('min_price')
      ? parseFloat(searchParams.get('min_price')!)
      : initialFilters?.minPrice;
    const maxPrice = searchParams.get('max_price')
      ? parseFloat(searchParams.get('max_price')!)
      : initialFilters?.maxPrice;
    const sortBy =
      (searchParams.get('sort') as SearchFilters['sortBy']) ||
      initialFilters?.sortBy ||
      'relevance';
    const page = searchParams.get('page')
      ? parseInt(searchParams.get('page')!)
      : initialFilters?.page || 1;
    const limit = searchParams.get('limit')
      ? parseInt(searchParams.get('limit')!)
      : initialFilters?.limit || DEFAULT_LIMIT;

    // Parse attribute filters from URL
    const attributes: Record<string, string[]> = {};
    searchParams.forEach((value, key) => {
      if (key.startsWith('attr_')) {
        const attrId = key.replace('attr_', '');
        attributes[attrId] = value.split(',');
      }
    });

    return {
      query,
      categoryId,
      minPrice,
      maxPrice,
      attributes:
        Object.keys(attributes).length > 0
          ? attributes
          : initialFilters?.attributes || {},
      sortBy,
      page,
      limit,
    };
  }, [searchParams, initialFilters, syncWithURL]);

  const [filters, setFilters] = useState<SearchFilters>(getInitialFilters);
  const [attributeFilters, setAttributeFilters] = useState<
    Record<string, string[]>
  >(filters.attributes || {});
  const [isLoading, setIsLoading] = useState(false);

  // Update URL when filters change
  const updateURL = useCallback(
    (newFilters: SearchFilters) => {
      if (!syncWithURL) return;

      const params = new URLSearchParams();

      if (newFilters.query) params.set('q', newFilters.query);
      if (newFilters.categoryId)
        params.set('category', newFilters.categoryId.toString());
      if (newFilters.minPrice)
        params.set('min_price', newFilters.minPrice.toString());
      if (newFilters.maxPrice)
        params.set('max_price', newFilters.maxPrice.toString());
      if (newFilters.sortBy && newFilters.sortBy !== 'relevance') {
        params.set('sort', newFilters.sortBy);
      }
      if (newFilters.page && newFilters.page > 1)
        params.set('page', newFilters.page.toString());
      if (newFilters.limit && newFilters.limit !== DEFAULT_LIMIT) {
        params.set('limit', newFilters.limit.toString());
      }

      // Add attribute filters to URL
      if (newFilters.attributes) {
        Object.entries(newFilters.attributes).forEach(([key, values]) => {
          if (values.length > 0) {
            params.set(`attr_${key}`, values.join(','));
          }
        });
      }

      const newURL = `${window.location.pathname}?${params.toString()}`;
      router.push(newURL, { scroll: false });
    },
    [router, syncWithURL]
  );

  // Debounced search function
  const debouncedSearch = useMemo(
    () =>
      debounce((newFilters: SearchFilters) => {
        setIsLoading(false);
        updateURL(newFilters);
      }, debounceDelay),
    [debounceDelay, updateURL]
  );

  // Update single filter
  const updateFilter = useCallback(
    (key: keyof SearchFilters, value: any) => {
      setFilters((prev) => {
        const newFilters = { ...prev, [key]: value };

        // Reset page when filters change
        if (key !== 'page' && key !== 'limit') {
          newFilters.page = 1;
        }

        if (autoApply) {
          setIsLoading(true);
          debouncedSearch(newFilters);
        }

        return newFilters;
      });
    },
    [autoApply, debouncedSearch]
  );

  // Update attribute filter
  const updateAttributeFilter = useCallback(
    (attributeId: string, values: string[]) => {
      setAttributeFilters((prev) => {
        const newAttrs = { ...prev };

        if (values.length === 0) {
          delete newAttrs[attributeId];
        } else {
          newAttrs[attributeId] = values;
        }

        const newFilters = { ...filters, attributes: newAttrs, page: 1 };
        setFilters(newFilters);

        if (autoApply) {
          setIsLoading(true);
          debouncedSearch(newFilters);
        }

        return newAttrs;
      });
    },
    [filters, autoApply, debouncedSearch]
  );

  // Clear all attribute filters
  const clearAttributeFilters = useCallback(() => {
    setAttributeFilters({});
    const newFilters = { ...filters, attributes: {}, page: 1 };
    setFilters(newFilters);

    if (autoApply) {
      setIsLoading(true);
      debouncedSearch(newFilters);
    }
  }, [filters, autoApply, debouncedSearch]);

  // Apply filters manually (when autoApply is false)
  const applyFilters = useCallback(() => {
    setIsLoading(true);
    const newFilters = { ...filters, attributes: attributeFilters };
    setFilters(newFilters);
    debouncedSearch(newFilters);
  }, [filters, attributeFilters, debouncedSearch]);

  // Sync with URL changes (browser back/forward)
  useEffect(() => {
    if (syncWithURL) {
      const newFilters = getInitialFilters();
      setFilters(newFilters);
      setAttributeFilters(newFilters.attributes || {});
    }
  }, [searchParams, syncWithURL, getInitialFilters]);

  return {
    filters,
    attributeFilters,
    isLoading,
    updateFilter,
    updateAttributeFilter,
    clearAttributeFilters,
    applyFilters,
    syncWithURL,
  };
}
