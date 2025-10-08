'use client';

import { useState, useCallback, useRef, useEffect } from 'react';
// import { useTranslations } from 'next-intl';
import type { GeocodeResult, AutocompleteItem } from '@/types/b2c';

interface AddressSearchProps {
  onLocationSelect: (location: GeocodeResult) => void;
  placeholder?: string;
  className?: string;
  showCurrentLocation?: boolean;
  disabled?: boolean;
}

const AddressSearch: React.FC<AddressSearchProps> = ({
  onLocationSelect,
  placeholder = 'Введите адрес...',
  className = '',
  showCurrentLocation = true,
  disabled = false,
}) => {
  // const t = useTranslations('storefronts');
  const [query, setQuery] = useState('');
  const [suggestions, setSuggestions] = useState<AutocompleteItem[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isOpen, setIsOpen] = useState(false);
  const [selectedIndex, setSelectedIndex] = useState(-1);
  const inputRef = useRef<HTMLInputElement>(null);
  const dropdownRef = useRef<HTMLDivElement>(null);
  const debounceTimeout = useRef<NodeJS.Timeout | undefined>(undefined);

  // Получение текущего местоположения
  const getCurrentLocation = useCallback(() => {
    if (!navigator.geolocation) {
      alert('Геолокация не поддерживается вашим браузером');
      return;
    }

    setIsLoading(true);
    navigator.geolocation.getCurrentPosition(
      async (position) => {
        try {
          const { latitude, longitude } = position.coords;

          // Обратное геокодирование для получения адреса
          const response = await fetch(
            `https://nominatim.openstreetmap.org/reverse?format=json&lat=${latitude}&lon=${longitude}&addressdetails=1`
          );

          if (response.ok) {
            const data = await response.json();
            const geocodeResult: GeocodeResult = {
              address: data.display_name,
              latitude,
              longitude,
              city:
                data.address?.city ||
                data.address?.town ||
                data.address?.village,
              country: data.address?.country,
              postalCode: data.address?.postcode,
            };

            onLocationSelect(geocodeResult);
            setQuery(geocodeResult.address);
            setIsOpen(false);
          }
        } catch (error) {
          console.error('Ошибка при получении адреса:', error);
          alert('Не удалось получить адрес для текущего местоположения');
        } finally {
          setIsLoading(false);
        }
      },
      (error) => {
        setIsLoading(false);
        console.error('Ошибка геолокации:', error);
        alert('Не удалось получить ваше местоположение');
      },
      {
        enableHighAccuracy: true,
        timeout: 10000,
        maximumAge: 300000, // 5 минут
      }
    );
  }, [onLocationSelect]);

  // Поиск адресов с помощью Nominatim (OpenStreetMap)
  const searchAddresses = useCallback(async (searchQuery: string) => {
    if (searchQuery.length < 3) {
      setSuggestions([]);
      return;
    }

    setIsLoading(true);

    try {
      // Используем Nominatim для поиска адресов
      const response = await fetch(
        `https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(searchQuery)}&limit=5&addressdetails=1&countrycodes=rs,ba,hr,me,mk,si` // Ограничиваем поиск балканскими странами
      );

      if (response.ok) {
        const data = await response.json();
        const items: AutocompleteItem[] = data.map((item: any) => ({
          id: item.place_id,
          label: item.display_name,
          description: `${item.type} in ${item.address?.city || item.address?.town || item.address?.country}`,
          coordinates: {
            lat: parseFloat(item.lat),
            lng: parseFloat(item.lon),
          },
        }));

        setSuggestions(items);
      }
    } catch (error) {
      console.error('Ошибка поиска адресов:', error);
      setSuggestions([]);
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Debounced search
  const debouncedSearch = useCallback(
    (searchQuery: string) => {
      if (debounceTimeout.current) {
        clearTimeout(debounceTimeout.current);
      }

      debounceTimeout.current = setTimeout(() => {
        searchAddresses(searchQuery);
      }, 300);
    },
    [searchAddresses]
  );

  const handleInputChange = useCallback(
    (value: string) => {
      setQuery(value);
      setSelectedIndex(-1);

      if (value.trim()) {
        debouncedSearch(value);
        setIsOpen(true);
      } else {
        setSuggestions([]);
        setIsOpen(false);
      }
    },
    [debouncedSearch]
  );

  const handleSuggestionSelect = useCallback(
    (suggestion: AutocompleteItem) => {
      if (suggestion.coordinates) {
        const geocodeResult: GeocodeResult = {
          address: suggestion.label,
          latitude: suggestion.coordinates.lat,
          longitude: suggestion.coordinates.lng,
        };

        onLocationSelect(geocodeResult);
        setQuery(suggestion.label);
        setIsOpen(false);
        setSuggestions([]);
      }
    },
    [onLocationSelect]
  );

  const handleKeyDown = useCallback(
    (e: React.KeyboardEvent) => {
      if (!isOpen || suggestions.length === 0) return;

      switch (e.key) {
        case 'ArrowDown':
          e.preventDefault();
          setSelectedIndex((prev) =>
            prev < suggestions.length - 1 ? prev + 1 : prev
          );
          break;
        case 'ArrowUp':
          e.preventDefault();
          setSelectedIndex((prev) => (prev > 0 ? prev - 1 : prev));
          break;
        case 'Enter':
          e.preventDefault();
          if (selectedIndex >= 0) {
            handleSuggestionSelect(suggestions[selectedIndex]);
          }
          break;
        case 'Escape':
          setIsOpen(false);
          setSelectedIndex(-1);
          break;
      }
    },
    [isOpen, suggestions, selectedIndex, handleSuggestionSelect]
  );

  // Закрытие при клике вне компонента
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target as Node) &&
        inputRef.current &&
        !inputRef.current.contains(event.target as Node)
      ) {
        setIsOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  return (
    <div className={`relative ${className}`}>
      <div className="relative">
        <input
          ref={inputRef}
          type="text"
          placeholder={placeholder}
          value={query}
          onChange={(e) => handleInputChange(e.target.value)}
          onKeyDown={handleKeyDown}
          disabled={disabled}
          className="input input-bordered w-full pr-20"
        />

        <div className="absolute right-2 top-1/2 transform -translate-y-1/2 flex items-center gap-1">
          {isLoading && (
            <span className="loading loading-spinner loading-sm"></span>
          )}

          {showCurrentLocation && (
            <button
              type="button"
              onClick={getCurrentLocation}
              disabled={disabled || isLoading}
              className="btn btn-ghost btn-xs"
              title="Использовать текущее местоположение"
            >
              📍
            </button>
          )}
        </div>
      </div>

      {/* Dropdown с предложениями */}
      {isOpen && suggestions.length > 0 && (
        <div
          ref={dropdownRef}
          className="absolute z-50 w-full mt-1 bg-white border border-gray-200 rounded-lg shadow-lg max-h-60 overflow-y-auto"
        >
          {suggestions.map((suggestion, index) => (
            <button
              key={suggestion.id}
              type="button"
              onClick={() => handleSuggestionSelect(suggestion)}
              className={`w-full text-left px-4 py-3 hover:bg-gray-50 border-b border-gray-100 last:border-b-0 ${
                index === selectedIndex ? 'bg-blue-50' : ''
              }`}
            >
              <div className="font-medium text-sm truncate">
                {suggestion.label}
              </div>
              {suggestion.description && (
                <div className="text-xs text-gray-500 truncate mt-1">
                  {suggestion.description}
                </div>
              )}
            </button>
          ))}
        </div>
      )}
    </div>
  );
};

export default AddressSearch;
