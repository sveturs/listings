import { useState, useCallback, useEffect, useMemo } from 'react';
import { useDebounce } from '@/hooks/useDebounce';
import { MapBounds } from '../types/gis';

// Типы для работы с городами и видимостью
export interface City {
  id: string;
  name: string;
  slug: string;
  country_code: string;
  center_point: {
    lat: number;
    lng: number;
  } | null;
  boundary?: any;
  population?: number;
  area_km2?: number;
  postal_codes?: string[];
  has_districts: boolean;
  priority: number;
  created_at: string;
  updated_at: string;
}

export interface CityWithDistance {
  city: City;
  distance: number; // расстояние в метрах от центра карты
}

export interface VisibleCitiesResponse {
  visible_cities: CityWithDistance[];
  closest_city: CityWithDistance | null;
}

export interface VisibleCitiesRequest {
  bounds: {
    north: number;
    south: number;
    east: number;
    west: number;
  };
  center: {
    lat: number;
    lng: number;
  };
}

export interface District {
  id: string;
  name: string;
  city_id?: string;
  country_code: string;
  center_point: {
    lat: number;
    lng: number;
  } | null;
  boundary?: any;
  population?: number;
  area_km2?: number;
  postal_codes?: string[];
  created_at: string;
  updated_at: string;
}

interface UseVisibleCitiesResult {
  // Основные данные
  visibleCities: CityWithDistance[];
  closestCity: CityWithDistance | null;
  availableDistricts: District[];

  // Состояние загрузки
  loading: boolean;
  error: string | null;

  // Методы
  updateViewport: (
    bounds: MapBounds,
    center: { lat: number; lng: number }
  ) => void;
  refreshCities: () => Promise<void>;
  getDistrictsForCity: (cityId: string) => Promise<District[]>;
  getDistrictsForCities: (cityIds: string[]) => Promise<District[]>;

  // Утилиты
  hasDistrictsInViewport: boolean;
  shouldShowDistrictSearch: boolean;
}

/**
 * Хук для определения видимых городов в viewport карты
 * и управления контекстно-зависимым поиском по районам
 */
export const useVisibleCities = (): UseVisibleCitiesResult => {
  // Состояние
  const [visibleCities, setVisibleCities] = useState<CityWithDistance[]>([]);
  const [closestCity, setClosestCity] = useState<CityWithDistance | null>(null);
  const [availableDistricts, setAvailableDistricts] = useState<District[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Viewport состояние (дебаунсированное для оптимизации)
  const [currentViewport, setCurrentViewport] = useState<{
    bounds: MapBounds;
    center: { lat: number; lng: number };
  } | null>(null);

  const debouncedViewport = useDebounce(currentViewport, 500);

  /**
   * Получение видимых городов из API
   */
  const fetchVisibleCities = useCallback(
    async (
      bounds: MapBounds,
      center: { lat: number; lng: number }
    ): Promise<VisibleCitiesResponse | null> => {
      try {
        const request: VisibleCitiesRequest = {
          bounds: {
            north: bounds.north,
            south: bounds.south,
            east: bounds.east,
            west: bounds.west,
          },
          center: {
            lat: center.lat,
            lng: center.lng,
          },
        };

        const response = await fetch('/api/v1/gis/cities/visible', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(request),
        });

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();

        if (!data.success || !data.data) {
          throw new Error(data.error || 'Failed to fetch visible cities');
        }

        return data.data as VisibleCitiesResponse;
      } catch (err) {
        console.error('Error fetching visible cities:', err);
        throw err;
      }
    },
    []
  );

  /**
   * Получение районов для конкретного города
   */
  const getDistrictsForCity = useCallback(
    async (cityId: string): Promise<District[]> => {
      try {
        const response = await fetch(
          `/api/v1/gis/districts?city_id=${encodeURIComponent(cityId)}`
        );

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();

        if (!data.success || !data.data) {
          throw new Error(data.error || 'Failed to fetch districts');
        }

        return data.data as District[];
      } catch (err) {
        console.error('Error fetching districts for city:', err);
        throw err;
      }
    },
    []
  );

  /**
   * Загружает районы для нескольких городов одновременно
   */
  const getDistrictsForCities = useCallback(
    async (cityIds: string[]): Promise<District[]> => {
      if (cityIds.length === 0) return [];

      try {
        const params = new URLSearchParams();
        cityIds.forEach((cityId) => {
          params.append('city_ids[]', cityId);
        });

        const response = await fetch(
          `/api/v1/gis/districts?${params.toString()}`
        );

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();

        if (!data.success || !data.data) {
          throw new Error(data.error || 'Failed to fetch districts');
        }

        return data.data as District[];
      } catch (err) {
        console.error('Error fetching districts for cities:', err);
        throw err;
      }
    },
    []
  );

  /**
   * Обновление видимых городов при изменении viewport
   */
  useEffect(() => {
    if (!debouncedViewport) return;

    const updateVisibleCities = async () => {
      setLoading(true);
      setError(null);

      try {
        const result = await fetchVisibleCities(
          debouncedViewport.bounds,
          debouncedViewport.center
        );

        if (result) {
          setVisibleCities(result.visible_cities);
          setClosestCity(result.closest_city);

          // Если есть ближайший город с районами, загружаем их
          if (result.closest_city?.city.has_districts) {
            const districts = await getDistrictsForCity(
              result.closest_city.city.id
            );
            setAvailableDistricts(districts);
          } else {
            setAvailableDistricts([]);
          }
        }
      } catch (err) {
        const errorMessage =
          err instanceof Error ? err.message : 'Unknown error';
        setError(errorMessage);
        setVisibleCities([]);
        setClosestCity(null);
        setAvailableDistricts([]);
      } finally {
        setLoading(false);
      }
    };

    updateVisibleCities();
  }, [debouncedViewport, fetchVisibleCities, getDistrictsForCity]);

  /**
   * Обновление viewport (вызывается из компонента карты)
   */
  const updateViewport = useCallback(
    (bounds: MapBounds, center: { lat: number; lng: number }) => {
      setCurrentViewport({ bounds, center });
    },
    []
  );

  /**
   * Принудительное обновление городов
   */
  const refreshCities = useCallback(async () => {
    if (!currentViewport) return;

    setLoading(true);
    setError(null);

    try {
      const result = await fetchVisibleCities(
        currentViewport.bounds,
        currentViewport.center
      );

      if (result) {
        setVisibleCities(result.visible_cities);
        setClosestCity(result.closest_city);

        if (result.closest_city?.city.has_districts) {
          const districts = await getDistrictsForCity(
            result.closest_city.city.id
          );
          setAvailableDistricts(districts);
        } else {
          setAvailableDistricts([]);
        }
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Unknown error';
      setError(errorMessage);
    } finally {
      setLoading(false);
    }
  }, [currentViewport, fetchVisibleCities, getDistrictsForCity]);

  /**
   * Вычисляемые свойства
   */
  const hasDistrictsInViewport = useMemo(() => {
    return visibleCities.some(
      (cityWithDistance) => cityWithDistance.city.has_districts
    );
  }, [visibleCities]);

  const shouldShowDistrictSearch = useMemo(() => {
    return hasDistrictsInViewport && availableDistricts.length > 0;
  }, [hasDistrictsInViewport, availableDistricts.length]);

  return {
    // Основные данные
    visibleCities,
    closestCity,
    availableDistricts,

    // Состояние загрузки
    loading,
    error,

    // Методы
    updateViewport,
    refreshCities,
    getDistrictsForCity,
    getDistrictsForCities,

    // Утилиты
    hasDistrictsInViewport,
    shouldShowDistrictSearch,
  };
};

/**
 * Утилитные функции для работы с городами
 */

/**
 * Вычисляет расстояние между двумя точками в метрах (формула Хаверсина)
 */
export const calculateDistance = (
  point1: { lat: number; lng: number },
  point2: { lat: number; lng: number }
): number => {
  const R = 6371000; // радиус Земли в метрах

  const lat1Rad = (point1.lat * Math.PI) / 180;
  const lat2Rad = (point2.lat * Math.PI) / 180;
  const deltaLatRad = ((point2.lat - point1.lat) * Math.PI) / 180;
  const deltaLngRad = ((point2.lng - point1.lng) * Math.PI) / 180;

  const a =
    Math.sin(deltaLatRad / 2) * Math.sin(deltaLatRad / 2) +
    Math.cos(lat1Rad) *
      Math.cos(lat2Rad) *
      Math.sin(deltaLngRad / 2) *
      Math.sin(deltaLngRad / 2);

  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));

  return R * c;
};

/**
 * Форматирует расстояние для отображения
 */
export const formatDistance = (distanceInMeters: number): string => {
  if (distanceInMeters < 1000) {
    return `${Math.round(distanceInMeters)} м`;
  } else {
    return `${(distanceInMeters / 1000).toFixed(1)} км`;
  }
};

/**
 * Определяет, находится ли точка внутри bounds
 */
export const isPointInBounds = (
  point: { lat: number; lng: number },
  bounds: MapBounds
): boolean => {
  return (
    point.lat >= bounds.south &&
    point.lat <= bounds.north &&
    point.lng >= bounds.west &&
    point.lng <= bounds.east
  );
};

/**
 * Создает viewport bounds с отступом
 */
export const expandBounds = (
  bounds: MapBounds,
  paddingPercent: number = 0.1
): MapBounds => {
  const latPadding = (bounds.north - bounds.south) * paddingPercent;
  const lngPadding = (bounds.east - bounds.west) * paddingPercent;

  return {
    north: Math.min(90, bounds.north + latPadding),
    south: Math.max(-90, bounds.south - latPadding),
    east: Math.min(180, bounds.east + lngPadding),
    west: Math.max(-180, bounds.west - lngPadding),
  };
};
