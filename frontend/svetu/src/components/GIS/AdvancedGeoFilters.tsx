'use client';

import { useState, useEffect, useCallback, useMemo } from 'react';
import { useTranslations } from 'next-intl';
import {
  MapIcon,
  ClockIcon,
  BuildingOfficeIcon,
  UsersIcon,
} from '@heroicons/react/24/outline';
import debounce from 'lodash/debounce';

export interface AdvancedGeoFiltersProps {
  onFiltersChange: (filters: AdvancedGeoFilters) => void;
  currentLocation?: { lat: number; lng: number };
  className?: string;
}

export interface AdvancedGeoFilters {
  travelTime?: {
    centerLat: number;
    centerLng: number;
    maxMinutes: number;
    transportMode: 'walking' | 'driving' | 'cycling' | 'transit';
  };
  poiFilter?: {
    poiType: POIType;
    maxDistance: number;
    minCount?: number;
  };
  densityFilter?: {
    avoidCrowded: boolean;
    maxDensity?: number;
    minDensity?: number;
  };
}

export type POIType =
  | 'school'
  | 'hospital'
  | 'metro'
  | 'supermarket'
  | 'park'
  | 'bank'
  | 'pharmacy'
  | 'bus_stop';

interface POIOption {
  value: POIType;
  label: string;
  icon: string;
}

export default function AdvancedGeoFilters({
  onFiltersChange,
  currentLocation,
  className = '',
}: AdvancedGeoFiltersProps) {
  const commonT = useTranslations('common');
  const t = (key: string) => commonT(`map.${key}`);

  // Travel Time Filter State
  const [travelTimeEnabled, setTravelTimeEnabled] = useState(false);
  const [travelMinutes, setTravelMinutes] = useState(15);
  const [transportMode, setTransportMode] = useState<
    'walking' | 'driving' | 'cycling' | 'transit'
  >('walking');

  // POI Filter State
  const [poiEnabled, setPOIEnabled] = useState(false);
  const [selectedPOI, setSelectedPOI] = useState<POIType>('metro');
  const [poiDistance, setPOIDistance] = useState(1000);
  const [poiMinCount, setPOIMinCount] = useState(0);
  const [poiResults, setPOIResults] = useState<any[]>([]);

  // Density Filter State
  const [densityEnabled, setDensityEnabled] = useState(false);
  const [avoidCrowded, setAvoidCrowded] = useState(false);
  const [maxDensity, setMaxDensity] = useState<number | undefined>();

  const poiOptions: POIOption[] = [
    { value: 'school', label: t('gis.poi.school'), icon: '🏫' },
    { value: 'hospital', label: t('gis.poi.hospital'), icon: '🏥' },
    { value: 'metro', label: t('gis.poi.metro'), icon: '🚇' },
    { value: 'supermarket', label: t('gis.poi.supermarket'), icon: '🛒' },
    { value: 'park', label: t('gis.poi.park'), icon: '🌳' },
    { value: 'bank', label: t('gis.poi.bank'), icon: '🏦' },
    { value: 'pharmacy', label: t('gis.poi.pharmacy'), icon: '💊' },
    { value: 'bus_stop', label: t('gis.poi.busStop'), icon: '🚌' },
  ];

  const transportModes = [
    { value: 'walking', label: t('gis.transport.walking'), icon: '🚶' },
    { value: 'driving', label: t('gis.transport.driving'), icon: '🚗' },
    { value: 'cycling', label: t('gis.transport.cycling'), icon: '🚴' },
    { value: 'transit', label: t('gis.transport.transit'), icon: '🚌' },
  ];

  // Debounced filter update
  const debouncedUpdate = useMemo(
    () =>
      debounce(() => {
        const filters: AdvancedGeoFilters = {};

        if (travelTimeEnabled && currentLocation) {
          filters.travelTime = {
            centerLat: currentLocation.lat,
            centerLng: currentLocation.lng,
            maxMinutes: travelMinutes,
            transportMode,
          };
        }

        if (poiEnabled) {
          filters.poiFilter = {
            poiType: selectedPOI,
            maxDistance: poiDistance,
            ...(poiMinCount > 0 && { minCount: poiMinCount }),
          };
        }

        if (densityEnabled) {
          filters.densityFilter = {
            avoidCrowded,
            ...(maxDensity && { maxDensity }),
          };
        }

        onFiltersChange(filters);
      }, 500),
    [
      onFiltersChange,
      travelTimeEnabled,
      currentLocation,
      travelMinutes,
      transportMode,
      poiEnabled,
      selectedPOI,
      poiDistance,
      poiMinCount,
      densityEnabled,
      avoidCrowded,
      maxDensity,
    ]
  );

  useEffect(() => {
    debouncedUpdate();
  }, [debouncedUpdate]);

  const searchPOIs = useCallback(async () => {
    if (!currentLocation) return;

    try {
      const response = await fetch(
        `/api/v1/gis/advanced/poi/search?` +
          new URLSearchParams({
            poi_type: selectedPOI,
            lat: currentLocation.lat.toString(),
            lng: currentLocation.lng.toString(),
            radius: '2',
          }),
        {
          headers: {
            'Content-Type': 'application/json',
          },
        }
      );

      if (response.ok) {
        const data = await response.json();
        setPOIResults(data.data || []);
      }
    } catch (error) {
      console.error('Failed to search POIs:', error);
    }
  }, [currentLocation, selectedPOI]);

  // Search POIs when enabled
  useEffect(() => {
    if (poiEnabled && currentLocation) {
      searchPOIs();
    }
  }, [poiEnabled, selectedPOI, currentLocation, searchPOIs]);

  // Save filters to localStorage
  const saveFilters = useCallback(() => {
    const savedFilters = {
      travelTimeEnabled,
      travelMinutes,
      transportMode,
      poiEnabled,
      selectedPOI,
      poiDistance,
      poiMinCount,
      densityEnabled,
      avoidCrowded,
      maxDensity,
    };
    localStorage.setItem('advancedGeoFilters', JSON.stringify(savedFilters));
  }, [
    travelTimeEnabled,
    travelMinutes,
    transportMode,
    poiEnabled,
    selectedPOI,
    poiDistance,
    poiMinCount,
    densityEnabled,
    avoidCrowded,
    maxDensity,
  ]);

  // Load filters from localStorage
  useEffect(() => {
    const saved = localStorage.getItem('advancedGeoFilters');
    if (saved) {
      try {
        const filters = JSON.parse(saved);
        setTravelTimeEnabled(filters.travelTimeEnabled || false);
        setTravelMinutes(filters.travelMinutes || 15);
        setTransportMode(filters.transportMode || 'walking');
        setPOIEnabled(filters.poiEnabled || false);
        setSelectedPOI(filters.selectedPOI || 'metro');
        setPOIDistance(filters.poiDistance || 1000);
        setPOIMinCount(filters.poiMinCount || 0);
        setDensityEnabled(filters.densityEnabled || false);
        setAvoidCrowded(filters.avoidCrowded || false);
        setMaxDensity(filters.maxDensity);
      } catch (e) {
        console.error('Failed to load saved filters:', e);
      }
    }
  }, []);

  useEffect(() => {
    saveFilters();
  }, [
    travelTimeEnabled,
    travelMinutes,
    transportMode,
    poiEnabled,
    selectedPOI,
    poiDistance,
    poiMinCount,
    densityEnabled,
    avoidCrowded,
    maxDensity,
    saveFilters,
  ]);

  return (
    <div className={`bg-base-200 rounded-lg p-4 space-y-4 ${className}`}>
      <h3 className="text-lg font-semibold flex items-center gap-2">
        <MapIcon className="h-5 w-5" />
        {t('gis.advancedFilters.title')}
      </h3>

      {/* Travel Time Filter */}
      <div className="border-t pt-4">
        <label className="flex items-center gap-2 cursor-pointer">
          <input
            type="checkbox"
            className="checkbox checkbox-primary"
            checked={travelTimeEnabled}
            onChange={(e) => setTravelTimeEnabled(e.target.checked)}
          />
          <ClockIcon className="h-5 w-5" />
          <span className="font-medium">
            {t('gis.advancedFilters.travelTime')}
          </span>
        </label>

        {travelTimeEnabled && (
          <div className="mt-3 space-y-3 pl-7">
            <div>
              <label className="text-sm">
                {t('gis.advancedFilters.maxTravelTime')}
              </label>
              <div className="flex items-center gap-2">
                <input
                  type="range"
                  min="5"
                  max="60"
                  value={travelMinutes}
                  onChange={(e) => setTravelMinutes(Number(e.target.value))}
                  className="range range-primary flex-1"
                />
                <span className="text-sm font-medium w-16">
                  {travelMinutes} {t('common.minutes')}
                </span>
              </div>
            </div>

            <div>
              <label className="text-sm">
                {t('gis.advancedFilters.transportMode')}
              </label>
              <div className="flex flex-wrap gap-2 mt-1">
                {transportModes.map((mode) => (
                  <button
                    key={mode.value}
                    onClick={() => setTransportMode(mode.value as any)}
                    className={`btn btn-sm ${transportMode === mode.value ? 'btn-primary' : 'btn-ghost'}`}
                  >
                    <span className="mr-1">{mode.icon}</span>
                    {mode.label}
                  </button>
                ))}
              </div>
            </div>
          </div>
        )}
      </div>

      {/* POI Filter */}
      <div className="border-t pt-4">
        <label className="flex items-center gap-2 cursor-pointer">
          <input
            type="checkbox"
            className="checkbox checkbox-primary"
            checked={poiEnabled}
            onChange={(e) => setPOIEnabled(e.target.checked)}
          />
          <BuildingOfficeIcon className="h-5 w-5" />
          <span className="font-medium">
            {t('gis.advancedFilters.nearPOI')}
          </span>
        </label>

        {poiEnabled && (
          <div className="mt-3 space-y-3 pl-7">
            <div>
              <label className="text-sm">
                {t('gis.advancedFilters.poiType')}
              </label>
              <select
                value={selectedPOI}
                onChange={(e) => setSelectedPOI(e.target.value as POIType)}
                className="select select-bordered select-sm w-full mt-1"
              >
                {poiOptions.map((poi) => (
                  <option key={poi.value} value={poi.value}>
                    {poi.icon} {poi.label}
                  </option>
                ))}
              </select>
            </div>

            <div>
              <label className="text-sm">
                {t('gis.advancedFilters.maxDistance')}
              </label>
              <div className="flex items-center gap-2">
                <input
                  type="range"
                  min="100"
                  max="5000"
                  step="100"
                  value={poiDistance}
                  onChange={(e) => setPOIDistance(Number(e.target.value))}
                  className="range range-primary flex-1"
                />
                <span className="text-sm font-medium w-20">
                  {poiDistance >= 1000
                    ? `${(poiDistance / 1000).toFixed(1)} ${t('common.km')}`
                    : `${poiDistance} ${t('common.m')}`}
                </span>
              </div>
            </div>

            {poiResults.length > 0 && (
              <div className="text-sm text-base-content/70">
                {t('gis.advancedFilters.foundPOIs', {
                  count: poiResults.length,
                })}
              </div>
            )}
          </div>
        )}
      </div>

      {/* Density Filter */}
      <div className="border-t pt-4">
        <label className="flex items-center gap-2 cursor-pointer">
          <input
            type="checkbox"
            className="checkbox checkbox-primary"
            checked={densityEnabled}
            onChange={(e) => setDensityEnabled(e.target.checked)}
          />
          <UsersIcon className="h-5 w-5" />
          <span className="font-medium">
            {t('gis.advancedFilters.density')}
          </span>
        </label>

        {densityEnabled && (
          <div className="mt-3 space-y-3 pl-7">
            <label className="flex items-center gap-2 cursor-pointer">
              <input
                type="checkbox"
                className="toggle toggle-primary toggle-sm"
                checked={avoidCrowded}
                onChange={(e) => setAvoidCrowded(e.target.checked)}
              />
              <span className="text-sm">
                {t('gis.advancedFilters.avoidCrowded')}
              </span>
            </label>

            {!avoidCrowded && (
              <div>
                <label className="text-sm">
                  {t('gis.advancedFilters.maxDensity')}
                </label>
                <input
                  type="number"
                  value={maxDensity || ''}
                  onChange={(e) =>
                    setMaxDensity(
                      e.target.value ? Number(e.target.value) : undefined
                    )
                  }
                  placeholder={t('gis.advancedFilters.densityPlaceholder')}
                  className="input input-bordered input-sm w-full mt-1"
                />
              </div>
            )}
          </div>
        )}
      </div>

      {/* Reset Filters */}
      <div className="flex justify-end pt-2">
        <button
          onClick={() => {
            setTravelTimeEnabled(false);
            setPOIEnabled(false);
            setDensityEnabled(false);
            localStorage.removeItem('advancedGeoFilters');
          }}
          className="btn btn-ghost btn-sm"
        >
          {t('common.resetFilters')}
        </button>
      </div>
    </div>
  );
}
