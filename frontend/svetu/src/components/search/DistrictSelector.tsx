'use client';

import React from 'react';

// Временные интерфейсы до исправления API типов
interface District {
  id: string;
  name: string;
  geometry?: any;
  boundary?: {
    coordinates: number[][][];
  };
  bounds?: [number, number, number, number];
  population?: number;
  area?: number;
  area_km2?: number;
}

interface Municipality {
  id: string;
  name: string;
  districts?: District[];
}

interface DistrictSelectorProps {
  selectedDistrictId?: string;
  selectedMunicipalityId?: string;
  onDistrictChange?: (districtId: string | null) => void;
  onMunicipalityChange?: (municipalityId: string | null) => void;
  className?: string;
}

export function DistrictSelector({
  selectedDistrictId,
  selectedMunicipalityId,
  onDistrictChange,
  onMunicipalityChange,
  className = '',
}: DistrictSelectorProps) {
  // DISTRICT FUNCTIONALITY TEMPORARILY DISABLED
  return null;
}