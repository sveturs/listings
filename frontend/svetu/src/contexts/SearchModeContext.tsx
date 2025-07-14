'use client';

import React, { createContext, useContext, useState, ReactNode } from 'react';

export type SearchMode = 'radius' | 'district' | 'none';

interface SearchModeContextType {
  searchMode: SearchMode;
  setSearchMode: (mode: SearchMode) => void;
  isDistrictSearchActive: boolean;
  isRadiusSearchActive: boolean;
}

const SearchModeContext = createContext<SearchModeContextType | undefined>(
  undefined
);

interface SearchModeProviderProps {
  children: ReactNode;
}

export function SearchModeProvider({ children }: SearchModeProviderProps) {
  const [searchMode, setSearchMode] = useState<SearchMode>('none');

  const value: SearchModeContextType = {
    searchMode,
    setSearchMode,
    isDistrictSearchActive: searchMode === 'district',
    isRadiusSearchActive: searchMode === 'radius',
  };

  return (
    <SearchModeContext.Provider value={value}>
      {children}
    </SearchModeContext.Provider>
  );
}

export function useSearchMode() {
  const context = useContext(SearchModeContext);
  if (context === undefined) {
    throw new Error('useSearchMode must be used within a SearchModeProvider');
  }
  return context;
}
