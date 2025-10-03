'use client';

import React, { createContext, useContext, ReactNode } from 'react';
import { useVisibleCities } from '../hooks/useVisibleCities';
import type { CityWithDistance, District } from '../hooks/useVisibleCities';

interface VisibleCitiesContextType {
  visibleCities: CityWithDistance[];
  closestCity: CityWithDistance | null;
  availableDistricts: District[];
  loading: boolean;
  error: string | null;
}

const VisibleCitiesContext = createContext<
  VisibleCitiesContextType | undefined
>(undefined);

interface VisibleCitiesProviderProps {
  children: ReactNode;
}

export function VisibleCitiesProvider({
  children,
}: VisibleCitiesProviderProps) {
  const visibleCitiesData = useVisibleCities();

  return (
    <VisibleCitiesContext.Provider value={visibleCitiesData}>
      {children}
    </VisibleCitiesContext.Provider>
  );
}

export function useVisibleCitiesContext(): VisibleCitiesContextType {
  const context = useContext(VisibleCitiesContext);
  if (context === undefined) {
    throw new Error(
      'useVisibleCitiesContext must be used within a VisibleCitiesProvider'
    );
  }
  return context;
}
