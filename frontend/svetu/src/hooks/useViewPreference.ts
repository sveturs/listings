'use client';

import { useState, useEffect } from 'react';
import type { ViewMode } from '@/components/common/ViewToggle';

const VIEW_PREFERENCE_KEY = 'marketplace-view-preference';

export function useViewPreference(defaultView: ViewMode = 'grid') {
  const [viewMode, setViewMode] = useState<ViewMode>(defaultView);

  // Load preference from localStorage on mount
  useEffect(() => {
    try {
      const savedView = localStorage.getItem(VIEW_PREFERENCE_KEY);
      if (savedView === 'grid' || savedView === 'list') {
        setViewMode(savedView);
      }
    } catch (error) {
      // Ignore localStorage errors (e.g., in SSR or private mode)
      console.warn('Failed to load view preference:', error);
    }
  }, []);

  // Save preference when it changes
  const updateViewMode = (newMode: ViewMode) => {
    setViewMode(newMode);
    try {
      localStorage.setItem(VIEW_PREFERENCE_KEY, newMode);
    } catch (error) {
      // Ignore localStorage errors
      console.warn('Failed to save view preference:', error);
    }
  };

  return [viewMode, updateViewMode] as const;
}
