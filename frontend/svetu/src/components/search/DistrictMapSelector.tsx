'use client';

import React from 'react';
import type { Feature, Polygon } from 'geojson';
import type { MapBounds } from '@/components/GIS/types/gis';

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

interface SpatialSearchResult {
  id: string;
  title: string;
  description?: string;
  latitude: number;
  longitude: number;
  distance?: number;
  category?: string;
  price?: number;
  currency?: string;
  imageUrl?: string;
  first_image_url?: string;
  category_name?: string;
  address?: string;
  user_email?: string;
}

interface DistrictMapSelectorProps {
  onSearchResults?: (results: SpatialSearchResult[]) => void;
  onDistrictBoundsChange?: (
    bounds: [number, number, number, number] | null
  ) => void;
  onDistrictBoundaryChange?: (boundary: Feature<Polygon> | null) => void;
  onViewportChange?: (
    bounds: MapBounds,
    center: { lat: number; lng: number }
  ) => void;
  currentViewport?: {
    bounds: MapBounds;
    center: { lat: number; lng: number };
  } | null;
  className?: string;
}

export function DistrictMapSelector({
  onSearchResults,
  onDistrictBoundsChange,
  onDistrictBoundaryChange,
  onViewportChange,
  currentViewport,
  className = '',
}: DistrictMapSelectorProps) {
  // DISTRICT FUNCTIONALITY TEMPORARILY DISABLED
  return null;
}