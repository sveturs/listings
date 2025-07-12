import { useState, useCallback } from 'react';
import {
  mapboxGeocoding,
  type NearbySearchResult,
} from '../utils/mapboxGeocoding';
import {
  GeoSearchResult,
  GeoSearchParams,
  NearbySearchParams,
} from '../types/gis';

interface UseMapboxGeoSearchResult {
  results: GeoSearchResult[];
  loading: boolean;
  error: string | null;
  search: (params: GeoSearchParams) => Promise<GeoSearchResult[]>;
  searchNearby: (params: NearbySearchParams) => Promise<NearbySearchResult[]>;
  clearResults: () => void;
}

/**
 * Хук для работы с MapBox Geocoding API
 * Заменяет старый useGeoSearch с более точными результатами и лучшей производительностью
 */
export const useMapboxGeoSearch = (): UseMapboxGeoSearchResult => {
  const [results, setResults] = useState<GeoSearchResult[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const search = useCallback(
    async (params: GeoSearchParams): Promise<GeoSearchResult[]> => {
      setLoading(true);
      setError(null);

      try {
        // Настройки для MapBox Geocoding
        const options: Parameters<typeof mapboxGeocoding.geocode>[1] = {
          limit: params.limit || 10,
          language: params.language ? [params.language] : ['en', 'ru', 'sr'],
          country: params.countrycodes
            ? params.countrycodes.split(',')
            : undefined,
        };

        // Добавляем proximity если есть viewbox (центральная точка)
        if (params.viewbox) {
          const bbox = params.viewbox.split(',').map(Number);
          if (bbox.length === 4) {
            const centerLon = (bbox[0] + bbox[2]) / 2;
            const centerLat = (bbox[1] + bbox[3]) / 2;
            options.proximity = [centerLon, centerLat];
            options.bbox = [bbox[0], bbox[1], bbox[2], bbox[3]];
          }
        }

        const mapboxResults = await mapboxGeocoding.geocode(
          params.query,
          options
        );

        // Конвертируем в формат GeoSearchResult для совместимости
        const searchResults: GeoSearchResult[] = mapboxResults.map(
          (item, index) => ({
            id: `mapbox-${index}-${item.lat}-${item.lon}`,
            display_name: item.display_name,
            lat: item.lat,
            lon: item.lon,
            boundingbox: item.boundingbox || ['0', '0', '0', '0'],
            type: item.type || 'unknown',
            class: item.class || 'unknown',
            importance: item.importance || 0,
          })
        );

        setResults(searchResults);
        return searchResults;
      } catch (err) {
        let errorMessage = 'geo_search.unknown_error';

        if (err instanceof Error) {
          if (err.message.includes('MAPBOX_ACCESS_TOKEN')) {
            errorMessage = 'geo_search.config_error';
          } else if (
            err.message.includes('network') ||
            err.message.includes('fetch')
          ) {
            errorMessage = 'geo_search.network_error';
          } else {
            errorMessage = err.message;
          }
        }

        setError(errorMessage);
        console.error('MapBox geo search error:', err);
        return [];
      } finally {
        setLoading(false);
      }
    },
    []
  );

  const searchNearby = useCallback(
    async (params: NearbySearchParams): Promise<NearbySearchResult[]> => {
      setLoading(true);
      setError(null);

      try {
        const nearbyResults = await mapboxGeocoding.searchNearby(
          params.latitude,
          params.longitude,
          params.category || '',
          {
            radius: params.radius,
            limit: params.limit || 10,
            types: ['poi', 'address'],
          }
        );

        return nearbyResults;
      } catch (err) {
        const errorMessage =
          err instanceof Error ? err.message : 'geo_search.nearby_error';
        setError(errorMessage);
        console.error('MapBox nearby search error:', err);
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

/**
 * Хук для работы с MapBox геокодированием
 * Заменяет старый useGeocoding с MapBox API
 */
export const useMapboxGeocoding = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const geocode = useCallback(
    async (address: string): Promise<GeoSearchResult | null> => {
      setLoading(true);
      setError(null);

      try {
        const results = await mapboxGeocoding.geocode(address, {
          limit: 1,
          language: ['en', 'ru', 'sr'],
        });

        if (results.length === 0) {
          return null;
        }

        const result = results[0];
        return {
          id: `mapbox-geocode-${result.lat}-${result.lon}`,
          display_name: result.display_name,
          lat: result.lat,
          lon: result.lon,
          boundingbox: result.boundingbox || ['0', '0', '0', '0'],
          type: result.type || 'unknown',
          class: result.class || 'unknown',
          importance: result.importance || 0,
        };
      } catch (err) {
        let errorMessage = 'geocoding.error';

        if (err instanceof Error) {
          if (err.message.includes('MAPBOX_ACCESS_TOKEN')) {
            errorMessage = 'geocoding.config_error';
          } else if (
            err.message.includes('network') ||
            err.message.includes('fetch')
          ) {
            errorMessage = 'geocoding.network_error';
          } else {
            errorMessage = err.message;
          }
        }

        setError(errorMessage);
        console.error('MapBox geocoding error:', err);
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
        const results = await mapboxGeocoding.reverseGeocode(lat, lon, {
          language: ['en', 'ru', 'sr'],
          types: ['address', 'poi', 'place'],
        });

        if (results.length === 0) {
          return null;
        }

        const result = results[0];
        return {
          id: `mapbox-reverse-${lat}-${lon}`,
          display_name: result.display_name,
          lat: result.lat,
          lon: result.lon,
          boundingbox: result.boundingbox || ['0', '0', '0', '0'],
          type: result.type || 'unknown',
          class: result.class || 'unknown',
          importance: result.importance || 0,
        };
      } catch (err) {
        let errorMessage = 'reverse_geocoding.error';

        if (err instanceof Error) {
          if (err.message.includes('MAPBOX_ACCESS_TOKEN')) {
            errorMessage = 'reverse_geocoding.config_error';
          } else if (
            err.message.includes('network') ||
            err.message.includes('fetch')
          ) {
            errorMessage = 'reverse_geocoding.network_error';
          } else {
            errorMessage = err.message;
          }
        }

        setError(errorMessage);
        console.error('MapBox reverse geocoding error:', err);
        return null;
      } finally {
        setLoading(false);
      }
    },
    []
  );

  const getCitySuggestions = useCallback(
    async (query: string, country?: string[]): Promise<GeoSearchResult[]> => {
      setLoading(true);
      setError(null);

      try {
        const results = await mapboxGeocoding.getCitySuggestions(query, {
          country,
          language: ['en', 'ru', 'sr'],
          limit: 5,
        });

        return results.map((item, index) => ({
          id: `mapbox-city-${index}-${item.lat}-${item.lon}`,
          display_name: item.display_name,
          lat: item.lat,
          lon: item.lon,
          boundingbox: item.boundingbox || ['0', '0', '0', '0'],
          type: item.type || 'unknown',
          class: item.class || 'unknown',
          importance: item.importance || 0,
        }));
      } catch (err) {
        let errorMessage = 'city_suggestions.error';

        if (err instanceof Error) {
          if (err.message.includes('MAPBOX_ACCESS_TOKEN')) {
            errorMessage = 'city_suggestions.config_error';
          } else if (
            err.message.includes('network') ||
            err.message.includes('fetch')
          ) {
            errorMessage = 'city_suggestions.network_error';
          } else {
            errorMessage = err.message;
          }
        }

        setError(errorMessage);
        console.error('MapBox city suggestions error:', err);
        return [];
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
    getCitySuggestions,
  };
};

// Утилитные функции остаются теми же
export const formatSearchResult = (result: GeoSearchResult): string => {
  const parts = result.display_name.split(', ');
  if (parts.length > 3) {
    return `${parts[0]}, ${parts[1]}, ${parts[parts.length - 1]}`;
  }
  return result.display_name;
};

export const getResultIcon = (result: GeoSearchResult): string => {
  switch (result.class) {
    case 'poi':
      return '📍';
    case 'address':
      return '🏠';
    case 'place':
      return '🏙️';
    case 'locality':
      return '🌆';
    case 'administrative':
      return '🏛️';
    case 'amenity':
      return '🏪';
    case 'highway':
      return '🛣️';
    case 'building':
      return '🏢';
    default:
      return '📍';
  }
};

export const sortResultsByImportance = (
  results: GeoSearchResult[]
): GeoSearchResult[] => {
  return [...results].sort((a, b) => b.importance - a.importance);
};
