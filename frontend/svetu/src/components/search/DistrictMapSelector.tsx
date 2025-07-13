'use client';

import React, { useEffect, useState } from 'react';
import { useTranslations } from 'next-intl';
import { DistrictSelector } from './DistrictSelector';
import type { components } from '@/types/generated/api';

type District = components['schemas']['types.District'];
type SpatialSearchResult = components['schemas']['types.SpatialSearchResult'];

interface DistrictMapSelectorProps {
  onSearchResults?: (results: SpatialSearchResult[]) => void;
  onDistrictBoundsChange?: (
    bounds: [number, number, number, number] | null
  ) => void;
  className?: string;
}

export function DistrictMapSelector({
  onSearchResults,
  onDistrictBoundsChange,
  className = '',
}: DistrictMapSelectorProps) {
  const t = useTranslations('search');
  const [selectedDistrictId, setSelectedDistrictId] = useState<string | null>(
    null
  );
  const [selectedMunicipalityId, setSelectedMunicipalityId] = useState<
    string | null
  >(null);
  const [isSearching, setIsSearching] = useState(false);
  const [currentDistrict, setCurrentDistrict] = useState<District | null>(null);

  // Загрузка информации о выбранном районе
  useEffect(() => {
    if (!selectedDistrictId) {
      setCurrentDistrict(null);
      onDistrictBoundsChange?.(null);
      return;
    }

    const fetchDistrictDetails = async () => {
      try {
        const response = await fetch(
          `/api/v1/gis/districts/${selectedDistrictId}`
        );
        if (!response.ok) {
          throw new Error('Failed to fetch district details');
        }
        const data = await response.json();
        const district = data.data as District;
        setCurrentDistrict(district);

        // Если у района есть границы, передаем их для отображения на карте
        if (district.boundary?.coordinates?.[0]) {
          const coords = district.boundary.coordinates[0];
          let minLng = 180,
            maxLng = -180,
            minLat = 90,
            maxLat = -90;

          coords.forEach((coord) => {
            minLng = Math.min(minLng, coord[0]);
            maxLng = Math.max(maxLng, coord[0]);
            minLat = Math.min(minLat, coord[1]);
            maxLat = Math.max(maxLat, coord[1]);
          });

          onDistrictBoundsChange?.([minLng, minLat, maxLng, maxLat]);
        }
      } catch (err) {
        console.error('Error fetching district details:', err);
      }
    };

    fetchDistrictDetails();
  }, [selectedDistrictId, onDistrictBoundsChange]);

  // Поиск объявлений при изменении выбора
  useEffect(() => {
    const searchListings = async () => {
      if (!selectedDistrictId && !selectedMunicipalityId) {
        onSearchResults?.([]);
        return;
      }

      setIsSearching(true);
      try {
        let url = '';
        if (selectedMunicipalityId) {
          url = `/api/v1/gis/search/by-municipality/${selectedMunicipalityId}`;
        } else if (selectedDistrictId) {
          url = `/api/v1/gis/search/by-district/${selectedDistrictId}`;
        }

        const response = await fetch(url);
        if (!response.ok) {
          throw new Error('Failed to search listings');
        }

        const data = await response.json();
        onSearchResults?.(data.data || []);
      } catch (err) {
        console.error('Error searching listings:', err);
        onSearchResults?.([]);
      } finally {
        setIsSearching(false);
      }
    };

    searchListings();
  }, [selectedDistrictId, selectedMunicipalityId, onSearchResults]);

  const handleDistrictChange = (districtId: string | null) => {
    setSelectedDistrictId(districtId);
    setSelectedMunicipalityId(null);
  };

  const handleMunicipalityChange = (municipalityId: string | null) => {
    setSelectedMunicipalityId(municipalityId);
  };

  return (
    <div className={`card bg-base-100 shadow-lg ${className}`}>
      <div className="card-body">
        <h3 className="card-title text-lg">{t('searchByDistrict')}</h3>

        <DistrictSelector
          selectedDistrictId={selectedDistrictId || undefined}
          selectedMunicipalityId={selectedMunicipalityId || undefined}
          onDistrictChange={handleDistrictChange}
          onMunicipalityChange={handleMunicipalityChange}
        />

        {/* Информация о выбранном районе */}
        {currentDistrict && (
          <div className="stats stats-vertical lg:stats-horizontal shadow mt-4">
            {currentDistrict.population && (
              <div className="stat">
                <div className="stat-title">{t('population')}</div>
                <div className="stat-value text-base">
                  {currentDistrict.population.toLocaleString()}
                </div>
              </div>
            )}
            {currentDistrict.area_km2 && (
              <div className="stat">
                <div className="stat-title">{t('area')}</div>
                <div className="stat-value text-base">
                  {currentDistrict.area_km2} km²
                </div>
              </div>
            )}
          </div>
        )}

        {/* Индикатор поиска */}
        {isSearching && (
          <div className="flex items-center justify-center mt-4">
            <span className="loading loading-spinner loading-md"></span>
            <span className="ml-2">{t('searching')}</span>
          </div>
        )}
      </div>
    </div>
  );
}
