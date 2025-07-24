'use client';

import { useState, useEffect } from 'react';
import type { GridColumns } from '@/components/common/GridColumnsToggle';

const GRID_COLUMNS_KEY = 'marketplace-grid-columns';

export function useGridColumns(defaultColumns: GridColumns = 1) {
  const [columns, setColumns] = useState<GridColumns>(defaultColumns);

  // Load preference from localStorage on mount
  useEffect(() => {
    try {
      const savedColumns = localStorage.getItem(GRID_COLUMNS_KEY);
      if (savedColumns && [1, 2, 3].includes(Number(savedColumns))) {
        setColumns(Number(savedColumns) as GridColumns);
      }
    } catch (error) {
      // Ignore localStorage errors (e.g., in SSR or private mode)
      console.warn('Failed to load grid columns preference:', error);
    }
  }, []);

  // Save preference when it changes
  const updateColumns = (newColumns: GridColumns) => {
    setColumns(newColumns);
    try {
      localStorage.setItem(GRID_COLUMNS_KEY, newColumns.toString());
    } catch (error) {
      // Ignore localStorage errors
      console.warn('Failed to save grid columns preference:', error);
    }
  };

  return [columns, updateColumns] as const;
}
