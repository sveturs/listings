import { useState, useCallback } from 'react';
import {
  GeoSearchResult,
  GeoSearchParams,
  GISApiResponse,
  NearbySearchParams,
} from '../types/gis';

interface UseGeoSearchResult {
  results: GeoSearchResult[];
  loading: boolean;
  error: string | null;
  search: (params: GeoSearchParams) => Promise<GeoSearchResult[]>;
  searchNearby: (params: NearbySearchParams) => Promise<any[]>;
  clearResults: () => void;
}

export const useGeoSearch = (): UseGeoSearchResult => {
  const [results, setResults] = useState<GeoSearchResult[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const search = useCallback(
    async (params: GeoSearchParams): Promise<GeoSearchResult[]> => {
      setLoading(true);
      setError(null);

      try {
        const queryParams = new URLSearchParams({
          q: params.query,
          format: 'json',
          addressdetails: '1',
          limit: (params.limit || 10).toString(),
          ...(params.countrycodes && { countrycodes: params.countrycodes }),
          ...(params.bounded && { bounded: '1' }),
          ...(params.viewbox && { viewbox: params.viewbox }),
          ...(params.language && { 'accept-language': params.language }),
        });

        // Используем Nominatim API для поиска
        const response = await fetch(
          `https://nominatim.openstreetmap.org/search?${queryParams.toString()}`
        );

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();

        const searchResults: GeoSearchResult[] = data.map((item: any) => ({
          id: item.place_id?.toString() || `${item.lat}-${item.lon}`,
          display_name: item.display_name,
          lat: item.lat,
          lon: item.lon,
          boundingbox: item.boundingbox,
          type: item.type,
          class: item.class,
          importance: item.importance || 0,
        }));

        setResults(searchResults);
        return searchResults;
      } catch (err) {
        const errorMessage =
          err instanceof Error ? err.message : 'geo_search.unknown_error';
        setError(errorMessage);
        console.error('Geo search error:', err);
        return [];
      } finally {
        setLoading(false);
      }
    },
    []
  );

  const searchNearby = useCallback(
    async (params: NearbySearchParams): Promise<any[]> => {
      setLoading(true);
      setError(null);

      try {
        // Используем наш backend API для поиска ближайших объектов
        const response = await fetch('/api/gis/nearby', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(params),
        });

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const result: GISApiResponse = await response.json();

        if (!result.success) {
          throw new Error(result.error || 'geo_search.api_error');
        }

        return result.data || [];
      } catch (err) {
        const errorMessage =
          err instanceof Error ? err.message : 'geo_search.nearby_error';
        setError(errorMessage);
        console.error('Nearby search error:', err);
        return [];
      } finally {
        setLoading(false);
      }
    },
    []
  );

  const clearResults = useCallback(() => {
    setResults([]);
    setError(null);
  }, []);

  return {
    results,
    loading,
    error,
    search,
    searchNearby,
    clearResults,
  };
};

// Хук для работы с маршрутами
export const useRouteCalculation = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const calculateRoute = useCallback(
    async (
      origin: [number, number],
      destination: [number, number],
      mode: 'driving' | 'walking' | 'cycling' = 'driving'
    ) => {
      setLoading(true);
      setError(null);

      try {
        const response = await fetch('/api/gis/route', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            origin,
            destination,
            mode,
          }),
        });

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const result: GISApiResponse = await response.json();

        if (!result.success) {
          throw new Error(result.error || 'route.calculation_error');
        }

        return result.data;
      } catch (err) {
        const errorMessage =
          err instanceof Error ? err.message : 'route.unknown_error';
        setError(errorMessage);
        console.error('Route calculation error:', err);
        return null;
      } finally {
        setLoading(false);
      }
    },
    []
  );

  return {
    loading,
    error,
    calculateRoute,
  };
};

// Хук для работы с геокодированием
export const useGeocoding = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const geocode = useCallback(
    async (address: string): Promise<GeoSearchResult | null> => {
      setLoading(true);
      setError(null);

      try {
        const response = await fetch(
          `https://nominatim.openstreetmap.org/search?q=${encodeURIComponent(address)}&format=json&limit=1`
        );

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();

        if (data.length === 0) {
          return null;
        }

        const result: GeoSearchResult = {
          id: data[0].place_id?.toString() || `${data[0].lat}-${data[0].lon}`,
          display_name: data[0].display_name,
          lat: data[0].lat,
          lon: data[0].lon,
          boundingbox: data[0].boundingbox,
          type: data[0].type,
          class: data[0].class,
          importance: data[0].importance || 0,
        };

        return result;
      } catch (err) {
        const errorMessage =
          err instanceof Error ? err.message : 'geocoding.error';
        setError(errorMessage);
        console.error('Geocoding error:', err);
        return null;
      } finally {
        setLoading(false);
      }
    },
    []
  );

  const reverseGeocode = useCallback(
    async (lat: number, lon: number): Promise<GeoSearchResult | null> => {
      setLoading(true);
      setError(null);

      try {
        const response = await fetch(
          `https://nominatim.openstreetmap.org/reverse?lat=${lat}&lon=${lon}&format=json`
        );

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();

        if (!data.display_name) {
          return null;
        }

        const result: GeoSearchResult = {
          id: data.place_id?.toString() || `${lat}-${lon}`,
          display_name: data.display_name,
          lat: lat.toString(),
          lon: lon.toString(),
          boundingbox: data.boundingbox || [
            lat.toString(),
            lat.toString(),
            lon.toString(),
            lon.toString(),
          ],
          type: data.type || 'unknown',
          class: data.class || 'unknown',
          importance: data.importance || 0,
        };

        return result;
      } catch (err) {
        const errorMessage =
          err instanceof Error ? err.message : 'reverse_geocoding.error';
        setError(errorMessage);
        console.error('Reverse geocoding error:', err);
        return null;
      } finally {
        setLoading(false);
      }
    },
    []
  );

  return {
    loading,
    error,
    geocode,
    reverseGeocode,
  };
};

// Утилитные функции для работы с результатами поиска
export const formatSearchResult = (result: GeoSearchResult): string => {
  const parts = result.display_name.split(', ');
  if (parts.length > 3) {
    return `${parts[0]}, ${parts[1]}, ${parts[parts.length - 1]}`;
  }
  return result.display_name;
};

export const getResultIcon = (result: GeoSearchResult): string => {
  switch (result.class) {
    case 'amenity':
      return '🏪';
    case 'highway':
      return '🛣️';
    case 'building':
      return '🏢';
    case 'place':
      return '📍';
    default:
      return '📍';
  }
};

export const sortResultsByImportance = (
  results: GeoSearchResult[]
): GeoSearchResult[] => {
  return [...results].sort((a, b) => b.importance - a.importance);
};
