'use client';

import { useEffect, useRef } from 'react';

const STORAGE_KEY = 'searchFilters';
const STORAGE_EXPIRY_KEY = 'searchFiltersExpiry';
const EXPIRY_TIME = 24 * 60 * 60 * 1000; // 24 hours

interface FilterPersistenceOptions {
  categoryId?: number;
  filters: Record<string, any>;
}

export const useFilterPersistence = ({
  categoryId,
  filters,
}: FilterPersistenceOptions) => {
  const isInitialized = useRef(false);

  // Save filters to localStorage
  const saveFilters = () => {
    try {
      const dataToSave = {
        categoryId,
        filters,
        timestamp: Date.now(),
      };
      localStorage.setItem(STORAGE_KEY, JSON.stringify(dataToSave));
      localStorage.setItem(
        STORAGE_EXPIRY_KEY,
        String(Date.now() + EXPIRY_TIME)
      );
    } catch (error) {
      console.warn('Failed to save filters to localStorage:', error);
    }
  };

  // Load filters from localStorage
  const loadFilters = (): {
    categoryId?: number;
    filters: Record<string, any>;
  } | null => {
    try {
      const expiryTime = localStorage.getItem(STORAGE_EXPIRY_KEY);

      // Check if filters have expired
      if (expiryTime && Date.now() > parseInt(expiryTime)) {
        clearFilters();
        return null;
      }

      const savedData = localStorage.getItem(STORAGE_KEY);
      if (savedData) {
        const parsed = JSON.parse(savedData);
        return {
          categoryId: parsed.categoryId,
          filters: parsed.filters || {},
        };
      }
    } catch (error) {
      console.warn('Failed to load filters from localStorage:', error);
    }
    return null;
  };

  // Clear filters from localStorage
  const clearFilters = () => {
    try {
      localStorage.removeItem(STORAGE_KEY);
      localStorage.removeItem(STORAGE_EXPIRY_KEY);
    } catch (error) {
      console.warn('Failed to clear filters from localStorage:', error);
    }
  };

  // Save filters when they change (debounced)
  useEffect(() => {
    if (!isInitialized.current) {
      isInitialized.current = true;
      return;
    }

    const timeoutId = setTimeout(() => {
      if (Object.keys(filters).length > 0 || categoryId) {
        saveFilters();
      } else {
        clearFilters();
      }
    }, 500); // Debounce for 500ms

    return () => clearTimeout(timeoutId);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [categoryId, JSON.stringify(filters)]);

  return {
    loadFilters,
    clearFilters,
  };
};

// Hook for saving search queries
const QUERY_HISTORY_KEY = 'searchQueryHistory';
const MAX_QUERY_HISTORY = 10;

export const useSearchHistory = () => {
  // Save search query to history
  const saveQuery = (query: string) => {
    if (!query.trim()) return;

    try {
      const history = getQueryHistory();

      // Remove duplicate if exists
      const filteredHistory = history.filter((q) => q !== query);

      // Add new query at the beginning
      const newHistory = [query, ...filteredHistory].slice(
        0,
        MAX_QUERY_HISTORY
      );

      localStorage.setItem(QUERY_HISTORY_KEY, JSON.stringify(newHistory));
    } catch (error) {
      console.warn('Failed to save search query:', error);
    }
  };

  // Get search query history
  const getQueryHistory = (): string[] => {
    try {
      const saved = localStorage.getItem(QUERY_HISTORY_KEY);
      return saved ? JSON.parse(saved) : [];
    } catch (error) {
      console.warn('Failed to load search history:', error);
      return [];
    }
  };

  // Clear search history
  const clearHistory = () => {
    try {
      localStorage.removeItem(QUERY_HISTORY_KEY);
    } catch (error) {
      console.warn('Failed to clear search history:', error);
    }
  };

  return {
    saveQuery,
    getQueryHistory,
    clearHistory,
  };
};
