'use client';

import React, { createContext, useContext, ReactNode } from 'react';
import { useVisibleCities } from '../hooks/useVisibleCities';
import type { City, District } from '../types';

interface VisibleCitiesContextType {
  visibleCities: City[];
  closestCity: { city: City; distance: number } | null;
  availableDistricts: District[];
  loading: boolean;
  error: string | null;
}

const VisibleCitiesContext = createContext<VisibleCitiesContextType | undefined>(undefined);

interface VisibleCitiesProviderProps {
  children: ReactNode;
}

export function VisibleCitiesProvider({ children }: VisibleCitiesProviderProps) {
  const visibleCitiesData = useVisibleCities();

  console.log('🌍 VisibleCitiesProvider render:', {
    availableDistricts: visibleCitiesData.availableDistricts.length,
    closestCity: visibleCitiesData.closestCity?.city.name,
    loading: visibleCitiesData.loading
  });

  return (
    <VisibleCitiesContext.Provider value={visibleCitiesData}>
      {children}
    </VisibleCitiesContext.Provider>
  );
}

export function useVisibleCitiesContext(): VisibleCitiesContextType {
  const context = useContext(VisibleCitiesContext);
  if (context === undefined) {
    throw new Error('useVisibleCitiesContext must be used within a VisibleCitiesProvider');
  }
  return context;
}
