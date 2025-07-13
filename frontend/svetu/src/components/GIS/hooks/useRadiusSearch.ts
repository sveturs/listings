import { useState, useCallback } from 'react';
import {
  RadiusSearchParams,
  RadiusSearchResponse,
  RadiusSearchResult,
  GeoLocation as _GeoLocation,
} from '../types/gis';
import { useGeolocation } from './useGeolocation';
import { useGeocoding } from './useGeoSearch';

interface UseRadiusSearchResult {
  results: RadiusSearchResult[];
  loading: boolean;
  error: string | null;
  total: number;
  searchCenter: { latitude: number; longitude: number } | null;
  searchRadius: number;
  search: (params: RadiusSearchParams) => Promise<RadiusSearchResponse | null>;
  searchByAddress: (
    address: string,
    radius: number,
    category?: string
  ) => Promise<RadiusSearchResponse | null>;
  searchByCurrentLocation: (
    radius: number,
    category?: string
  ) => Promise<RadiusSearchResponse | null>;
  clearResults: () => void;
}

// Утилита для создания таймаута
const fetchWithTimeout = async (
  url: string,
  options: RequestInit = {},
  timeout = 10000
) => {
  const controller = new AbortController();
  const id = setTimeout(() => controller.abort(), timeout);

  try {
    const response = await fetch(url, {
      ...options,
      signal: controller.signal,
    });
    clearTimeout(id);
    return response;
  } catch (error) {
    clearTimeout(id);
    if (error instanceof Error && error.name === 'AbortError') {
      throw new Error('Request timeout');
    }
    throw error;
  }
};

export const useRadiusSearch = (): UseRadiusSearchResult => {
  const [results, setResults] = useState<RadiusSearchResult[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [total, setTotal] = useState(0);
  const [searchCenter, setSearchCenter] = useState<{
    latitude: number;
    longitude: number;
  } | null>(null);
  const [searchRadius, setSearchRadius] = useState(0);

  const { getCurrentPosition } = useGeolocation();
  const { geocode } = useGeocoding();

  // Основная функция поиска
  const search = useCallback(
    async (
      params: RadiusSearchParams
    ): Promise<RadiusSearchResponse | null> => {
      setLoading(true);
      setError(null);

      try {
        const response = await fetchWithTimeout(
          '/api/v1/gis/search/radius',
          {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify(params),
          },
          15000
        );

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const result = await response.json();

        if (result.error) {
          throw new Error(result.error);
        }

        const searchResponse: RadiusSearchResponse = {
          items: result.data?.items || [],
          total: result.data?.total || 0,
          center: {
            latitude: params.latitude,
            longitude: params.longitude,
          },
          radius: params.radius,
        };

        setResults(searchResponse.items);
        setTotal(searchResponse.total);
        setSearchCenter(searchResponse.center);
        setSearchRadius(searchResponse.radius);

        return searchResponse;
      } catch (err) {
        let errorMessage = 'radius_search.unknown_error';

        if (err instanceof TypeError && err.message === 'Failed to fetch') {
          errorMessage = 'radius_search.network_error';
          console.error(
            'Network error: Unable to reach radius search service. Possible network problem.'
          );
        } else if (err instanceof Error && err.message === 'Request timeout') {
          errorMessage = 'radius_search.timeout_error';
          console.error(
            'Request timeout: Radius search service took too long to respond.'
          );
        } else if (err instanceof Error) {
          errorMessage = err.message;
        }

        setError(errorMessage);
        console.error('Radius search error:', err);
        return null;
      } finally {
        setLoading(false);
      }
    },
    []
  );

  // Поиск по адресу
  const searchByAddress = useCallback(
    async (
      address: string,
      radius: number,
      category?: string
    ): Promise<RadiusSearchResponse | null> => {
      try {
        setLoading(true);
        setError(null);

        // Сначала геокодируем адрес
        const geocodeResult = await geocode(address);

        if (!geocodeResult) {
          throw new Error('radius_search.address_not_found');
        }

        const params: RadiusSearchParams = {
          latitude: parseFloat(geocodeResult.lat),
          longitude: parseFloat(geocodeResult.lon),
          radius,
          category,
          limit: 50,
        };

        return await search(params);
      } catch (err) {
        const errorMessage =
          err instanceof Error ? err.message : 'radius_search.address_error';
        setError(errorMessage);
        console.error('Search by address error:', err);
        return null;
      }
    },
    [geocode, search]
  );

  // Поиск по текущему местоположению
  const searchByCurrentLocation = useCallback(
    async (
      radius: number,
      category?: string
    ): Promise<RadiusSearchResponse | null> => {
      try {
        setLoading(true);
        setError(null);

        // Получаем текущее местоположение
        const location = await getCurrentPosition();

        const params: RadiusSearchParams = {
          latitude: location.latitude,
          longitude: location.longitude,
          radius,
          category,
          limit: 50,
        };

        return await search(params);
      } catch (err) {
        const errorMessage =
          err instanceof Error ? err.message : 'radius_search.location_error';
        setError(errorMessage);
        console.error('Search by current location error:', err);
        return null;
      }
    },
    [getCurrentPosition, search]
  );

  const clearResults = useCallback(() => {
    setResults([]);
    setTotal(0);
    setSearchCenter(null);
    setSearchRadius(0);
    setError(null);
  }, []);

  return {
    results,
    loading,
    error,
    total,
    searchCenter,
    searchRadius,
    search,
    searchByAddress,
    searchByCurrentLocation,
    clearResults,
  };
};

// Утилитные функции
export const formatRadius = (radius: number): string => {
  if (radius < 1) {
    return `${Math.round(radius * 1000)} м`;
  }
  return `${radius.toFixed(1)} км`;
};

export const validateRadius = (
  radius: number,
  minRadius: number,
  maxRadius: number
): boolean => {
  return radius >= minRadius && radius <= maxRadius;
};

export const normalizeRadius = (
  radius: number,
  minRadius: number,
  maxRadius: number
): number => {
  return Math.max(minRadius, Math.min(maxRadius, radius));
};

// Функция для создания круга радиуса на карте
export const createRadiusCircle = (
  center: { latitude: number; longitude: number },
  radius: number // в км
) => {
  const points = 64;
  const coordinates = [];
  const earthRadius = 6371; // радиус Земли в км

  for (let i = 0; i < points; i++) {
    const angle = (i * 360) / points;
    const radian = (angle * Math.PI) / 180;

    const lat1 = (center.latitude * Math.PI) / 180;
    const lon1 = (center.longitude * Math.PI) / 180;

    const lat2 = Math.asin(
      Math.sin(lat1) * Math.cos(radius / earthRadius) +
        Math.cos(lat1) * Math.sin(radius / earthRadius) * Math.cos(radian)
    );

    const lon2 =
      lon1 +
      Math.atan2(
        Math.sin(radian) * Math.sin(radius / earthRadius) * Math.cos(lat1),
        Math.cos(radius / earthRadius) - Math.sin(lat1) * Math.sin(lat2)
      );

    coordinates.push([(lon2 * 180) / Math.PI, (lat2 * 180) / Math.PI]);
  }

  // Замыкаем круг
  coordinates.push(coordinates[0]);

  return {
    type: 'Feature' as const,
    geometry: {
      type: 'Polygon' as const,
      coordinates: [coordinates],
    },
    properties: {
      radius,
      center,
    },
  };
};
