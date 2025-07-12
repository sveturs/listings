import { useState, useCallback, useRef, useEffect } from 'react';

// Типы для геокодирования
export interface AddressGeocodingResult {
  id: string;
  text: string;
  place_name: string;
  location: {
    lat: number;
    lng: number;
  };
  address_components: {
    country?: string;
    country_code?: string;
    city?: string;
    district?: string;
    street?: string;
    house_number?: string;
    postal_code?: string;
    formatted: string;
  };
  confidence: number;
  place_types: string[];
}

export interface UseAddressGeocodingOptions {
  debounceMs?: number;
  minLength?: number;
  maxSuggestions?: number;
  country?: string[];
  language?: string;
}

export interface UseAddressGeocodingReturn {
  suggestions: AddressGeocodingResult[];
  loading: boolean;
  error: string | null;
  search: (query: string) => void;
  reverseGeocode: (lat: number, lng: number) => Promise<AddressGeocodingResult | null>;
  validateAddress: (address: string) => Promise<{
    success: boolean;
    location?: { lat: number; lng: number };
    formatted_address?: string;
    confidence?: number;
    warnings?: string[];
  }>;
  clearSuggestions: () => void;
  clearError: () => void;
}

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3000';

export function useAddressGeocoding(options: UseAddressGeocodingOptions = {}): UseAddressGeocodingReturn {
  const {
    debounceMs = 300,
    minLength = 2,
    maxSuggestions = 5,
    country = ['rs', 'hr', 'ba', 'me'],
    language = 'en'
  } = options;

  const [suggestions, setSuggestions] = useState<AddressGeocodingResult[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  const debounceRef = useRef<NodeJS.Timeout>();
  const abortControllerRef = useRef<AbortController>();

  // Очистка при размонтировании
  useEffect(() => {
    return () => {
      if (debounceRef.current) {
        clearTimeout(debounceRef.current);
      }
      if (abortControllerRef.current) {
        abortControllerRef.current.abort();
      }
    };
  }, []);

  const search = useCallback((query: string) => {
    // Очищаем предыдущие таймеры и запросы
    if (debounceRef.current) {
      clearTimeout(debounceRef.current);
    }
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
    }

    // Проверяем минимальную длину
    if (query.length < minLength) {
      setSuggestions([]);
      setError(null);
      return;
    }

    // Устанавливаем debounce
    debounceRef.current = setTimeout(async () => {
      setLoading(true);
      setError(null);

      try {
        // Создаем новый AbortController для этого запроса
        abortControllerRef.current = new AbortController();

        const searchParams = new URLSearchParams({
          q: query,
          limit: maxSuggestions.toString(),
          language,
        });

        if (country.length > 0) {
          searchParams.set('country_code', country.join(','));
        }

        const response = await fetch(
          `${API_BASE_URL}/api/v1/gis/geocode/suggestions?${searchParams}`,
          {
            signal: abortControllerRef.current.signal,
            headers: {
              'Content-Type': 'application/json',
            },
          }
        );

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        
        if (data.success && Array.isArray(data.data)) {
          // Конвертируем данные в нужный формат
          const convertedSuggestions: AddressGeocodingResult[] = data.data.map((item: any) => ({
            id: item.id || `${item.location.lat}-${item.location.lng}`,
            text: item.text || item.place_name,
            place_name: item.place_name || item.text,
            location: {
              lat: item.location.lat,
              lng: item.location.lng,
            },
            address_components: {
              country: item.address_components?.country,
              country_code: item.address_components?.country_code,
              city: item.address_components?.city,
              district: item.address_components?.district,
              street: item.address_components?.street,
              house_number: item.address_components?.house_number,
              postal_code: item.address_components?.postal_code,
              formatted: item.address_components?.formatted || item.place_name,
            },
            confidence: item.confidence || 0.5,
            place_types: item.place_types || [],
          }));

          setSuggestions(convertedSuggestions);
        } else {
          setSuggestions([]);
        }
      } catch (err: any) {
        if (err.name === 'AbortError') {
          // Запрос был отменен - это нормально
          return;
        }
        
        console.error('Address geocoding error:', err);
        setError('Failed to search addresses. Please try again.');
        setSuggestions([]);
      } finally {
        setLoading(false);
      }
    }, debounceMs);
  }, [debounceMs, minLength, maxSuggestions, country, language]);

  const reverseGeocode = useCallback(async (lat: number, lng: number): Promise<AddressGeocodingResult | null> => {
    try {
      setLoading(true);
      setError(null);

      const searchParams = new URLSearchParams({
        lat: lat.toString(),
        lng: lng.toString(),
        language,
      });

      const response = await fetch(
        `${API_BASE_URL}/api/v1/gis/geocode/reverse?${searchParams}`,
        {
          headers: {
            'Content-Type': 'application/json',
          },
        }
      );

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      
      if (data.success && data.data) {
        return {
          id: data.data.id || `${lat}-${lng}`,
          text: data.data.text || data.data.place_name,
          place_name: data.data.place_name || data.data.text,
          location: { lat, lng },
          address_components: {
            country: data.data.address_components?.country,
            country_code: data.data.address_components?.country_code,
            city: data.data.address_components?.city,
            district: data.data.address_components?.district,
            street: data.data.address_components?.street,
            house_number: data.data.address_components?.house_number,
            postal_code: data.data.address_components?.postal_code,
            formatted: data.data.address_components?.formatted || data.data.place_name,
          },
          confidence: data.data.confidence || 0.5,
          place_types: data.data.place_types || [],
        };
      }

      return null;
    } catch (err: any) {
      console.error('Reverse geocoding error:', err);
      setError('Failed to get address from location.');
      return null;
    } finally {
      setLoading(false);
    }
  }, [language]);

  const validateAddress = useCallback(async (address: string) => {
    try {
      setLoading(true);
      setError(null);

      const response = await fetch(
        `${API_BASE_URL}/api/v1/gis/geocode/validate`,
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            address,
            language,
            country_code: country[0], // Используем первую страну как основную
          }),
        }
      );

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      
      if (data.success && data.data) {
        return {
          success: data.data.success,
          location: data.data.location,
          formatted_address: data.data.formatted_address,
          confidence: data.data.confidence,
          warnings: data.data.warnings,
        };
      }

      return { success: false };
    } catch (err: any) {
      console.error('Address validation error:', err);
      setError('Failed to validate address.');
      return { success: false };
    } finally {
      setLoading(false);
    }
  }, [language, country]);

  const clearSuggestions = useCallback(() => {
    setSuggestions([]);
  }, []);

  const clearError = useCallback(() => {
    setError(null);
  }, []);

  return {
    suggestions,
    loading,
    error,
    search,
    reverseGeocode,
    validateAddress,
    clearSuggestions,
    clearError,
  };
}