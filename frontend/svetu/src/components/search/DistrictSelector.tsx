'use client';

// import React from 'react';

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

interface _Municipality {
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
  selectedDistrictId: _selectedDistrictId,
  selectedMunicipalityId: _selectedMunicipalityId,
  onDistrictChange: _onDistrictChange,
  onMunicipalityChange: _onMunicipalityChange,
  className: _className = '',
}: DistrictSelectorProps) {
  // DISTRICT FUNCTIONALITY TEMPORARILY DISABLED
  return null;
}
