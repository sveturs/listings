'use client';

// import React from 'react';

interface District {
  id: string;
  name: string;
  country_code: string;
  population?: number;
  area_km2?: number;
}

interface Municipality {
  id: string;
  name: string;
  district_id: string;
  country_code: string;
  population?: number;
  area_km2?: number;
}

interface DistrictFilterProps {
  onDistrictSelect?: (district: District | null) => void;
  onMunicipalitySelect?: (municipality: Municipality | null) => void;
  className?: string;
  disabled?: boolean;
}

export default function DistrictFilter({
  onDistrictSelect: _onDistrictSelect,
  onMunicipalitySelect: _onMunicipalitySelect,
  className: _className = '',
  disabled: _disabled = false,
}: DistrictFilterProps) {
  // DISTRICT FUNCTIONALITY TEMPORARILY DISABLED
  return null;
}
