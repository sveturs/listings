import { useState, useCallback, useEffect } from 'react';
import { useSearchParams } from 'next/navigation';
import { AdvancedGeoFilters } from '@/components/GIS/AdvancedGeoFilters';

export function useAdvancedGeoFilters() {
  const searchParams = useSearchParams();
  const [filters, setFilters] = useState<AdvancedGeoFilters>({});
  const [isLoading, _setIsLoading] = useState(false);

  // Загрузка фильтров из URL при инициализации
  useEffect(() => {
    const loadedFilters: AdvancedGeoFilters = {};

    // Travel Time Filter
    const travelTime = searchParams?.get('travel_time');
    if (travelTime) {
      try {
        const [lat, lng, minutes, mode] = travelTime.split(',');
        loadedFilters.travelTime = {
          centerLat: parseFloat(lat),
          centerLng: parseFloat(lng),
          maxMinutes: parseInt(minutes),
          transportMode: mode as any,
        };
      } catch (e) {
        console.error('Failed to parse travel_time filter:', e);
      }
    }

    // POI Filter
    const poiFilter = searchParams?.get('poi');
    if (poiFilter) {
      try {
        const [type, distance, minCount] = poiFilter.split(',');
        loadedFilters.poiFilter = {
          poiType: type as any,
          maxDistance: parseInt(distance),
          ...(minCount && { minCount: parseInt(minCount) }),
        };
      } catch (e) {
        console.error('Failed to parse poi filter:', e);
      }
    }

    // Density Filter
    const densityFilter = searchParams?.get('density');
    if (densityFilter) {
      try {
        const parts = densityFilter.split(',');
        loadedFilters.densityFilter = {
          avoidCrowded: parts[0] === 'avoid',
          ...(parts[1] && { maxDensity: parseFloat(parts[1]) }),
          ...(parts[2] && { minDensity: parseFloat(parts[2]) }),
        };
      } catch (e) {
        console.error('Failed to parse density filter:', e);
      }
    }

    setFilters(loadedFilters);
  }, [searchParams]);

  // Обновление URL при изменении фильтров
  const updateURL = useCallback(
    (newFilters: AdvancedGeoFilters) => {
      const params = new URLSearchParams(searchParams?.toString() || '');

      // Travel Time
      if (newFilters.travelTime) {
        const { centerLat, centerLng, maxMinutes, transportMode } =
          newFilters.travelTime;
        params.set(
          'travel_time',
          `${centerLat},${centerLng},${maxMinutes},${transportMode}`
        );
      } else {
        params.delete('travel_time');
      }

      // POI
      if (newFilters.poiFilter) {
        const { poiType, maxDistance, minCount } = newFilters.poiFilter;
        const value = minCount
          ? `${poiType},${maxDistance},${minCount}`
          : `${poiType},${maxDistance}`;
        params.set('poi', value);
      } else {
        params.delete('poi');
      }

      // Density
      if (newFilters.densityFilter) {
        const { avoidCrowded, maxDensity, minDensity } =
          newFilters.densityFilter;
        const parts = [avoidCrowded ? 'avoid' : 'normal'];
        if (maxDensity) parts.push(maxDensity.toString());
        if (minDensity) parts.push(minDensity.toString());
        params.set('density', parts.join(','));
      } else {
        params.delete('density');
      }

      // Обновляем URL без перезагрузки страницы
      const newURL = `${window.location.pathname}?${params.toString()}`;
      window.history.pushState({}, '', newURL);
    },
    [searchParams]
  );

  const handleFiltersChange = useCallback(
    (newFilters: AdvancedGeoFilters) => {
      setFilters(newFilters);
      updateURL(newFilters);

      // Отправляем событие для обновления результатов поиска
      window.dispatchEvent(
        new CustomEvent('advancedGeoFiltersChanged', {
          detail: newFilters,
        })
      );
    },
    [updateURL]
  );

  const clearFilters = useCallback(() => {
    setFilters({});
    updateURL({});
    localStorage.removeItem('advancedGeoFilters');

    window.dispatchEvent(
      new CustomEvent('advancedGeoFiltersChanged', {
        detail: {},
      })
    );
  }, [updateURL]);

  return {
    filters,
    isLoading,
    setFilters: handleFiltersChange,
    clearFilters,
  };
}
