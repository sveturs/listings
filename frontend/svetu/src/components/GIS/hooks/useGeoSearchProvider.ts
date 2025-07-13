import { useGeoSearch, useGeocoding } from './useGeoSearch';
import { useMapboxGeoSearch, useMapboxGeocoding } from './useMapboxGeoSearch';

/**
 * Провайдер для выбора между Nominatim и MapBox геокодированием
 * Позволяет легко переключаться между API без изменения компонентов
 */

// Флаг для переключения между провайдерами
const USE_MAPBOX_GEOCODING =
  process.env.NEXT_PUBLIC_USE_MAPBOX_GEOCODING === 'true';

/**
 * Хук для поиска мест
 * Автоматически выбирает между Nominatim и MapBox на основе конфигурации
 */
export const useGeoSearchProvider = () => {
  const nominatimSearch = useGeoSearch();
  const mapboxSearch = useMapboxGeoSearch();

  // Возвращаем MapBox или Nominatim в зависимости от конфигурации
  if (USE_MAPBOX_GEOCODING) {
    console.log('[GeoSearchProvider] Using MapBox Geocoding API');
    return mapboxSearch;
  } else {
    console.log('[GeoSearchProvider] Using Nominatim API');
    return nominatimSearch;
  }
};

/**
 * Хук для геокодирования
 * Автоматически выбирает между Nominatim и MapBox на основе конфигурации
 */
export const useGeocodingProvider = () => {
  const nominatimGeocoding = useGeocoding();
  const mapboxGeocoding = useMapboxGeocoding();

  // Возвращаем MapBox или Nominatim в зависимости от конфигурации
  if (USE_MAPBOX_GEOCODING) {
    console.log('[GeocodingProvider] Using MapBox Geocoding API');
    return mapboxGeocoding;
  } else {
    console.log('[GeocodingProvider] Using Nominatim API');
    return nominatimGeocoding;
  }
};

/**
 * Утилита для проверки текущего провайдера
 */
export const getCurrentGeocodingProvider = (): 'mapbox' | 'nominatim' => {
  return USE_MAPBOX_GEOCODING ? 'mapbox' : 'nominatim';
};

/**
 * Утилита для получения информации о возможностях текущего провайдера
 */
export const getGeocodingProviderInfo = () => {
  if (USE_MAPBOX_GEOCODING) {
    return {
      provider: 'mapbox' as const,
      name: 'MapBox Geocoding API',
      features: {
        accuracy: 'high',
        performance: 'high',
        categories: true,
        languageSupport: 'excellent',
        rateLimits: 'generous',
        CORS: 'native',
      },
      limitations: {
        requiresToken: true,
        commercial: true,
      },
    };
  } else {
    return {
      provider: 'nominatim' as const,
      name: 'OpenStreetMap Nominatim',
      features: {
        accuracy: 'medium',
        performance: 'medium',
        categories: false,
        languageSupport: 'good',
        rateLimits: 'strict',
        CORS: 'proxy-required',
      },
      limitations: {
        requiresToken: false,
        commercial: false,
      },
    };
  }
};

// Экспортируем для тестирования
export { USE_MAPBOX_GEOCODING };
