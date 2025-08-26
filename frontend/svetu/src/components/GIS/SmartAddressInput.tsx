'use client';

import React, { useState, useRef, useEffect, useCallback } from 'react';
import {
  useAddressGeocoding,
  AddressGeocodingResult,
} from '@/hooks/useAddressGeocoding';
import { toast } from '@/utils/toast';
import { useTranslations } from 'next-intl';

export interface SmartAddressInputProps {
  value: string;
  onChange: (value: string, result?: AddressGeocodingResult) => void;
  onLocationSelect?: (location: {
    lat: number;
    lng: number;
    address: string;
    confidence: number;
  }) => void;
  placeholder?: string;
  disabled?: boolean;
  required?: boolean;
  showCurrentLocation?: boolean;
  country?: string[];
  language?: string;
  error?: string;
  className?: string;
}

export default function SmartAddressInput({
  value,
  onChange,
  onLocationSelect,
  placeholder = 'Начните вводить адрес...',
  disabled = false,
  required = false,
  showCurrentLocation = true,
  country = ['rs', 'hr', 'ba', 'me'],
  language = 'ru',
  error,
  className = '',
}: SmartAddressInputProps) {
  const t = useTranslations('common');
  const [isOpen, setIsOpen] = useState(false);
  const [highlightedIndex, setHighlightedIndex] = useState(-1);
  const [isGettingLocation, setIsGettingLocation] = useState(false);

  const inputRef = useRef<HTMLInputElement>(null);
  const dropdownRef = useRef<HTMLDivElement>(null);

  const {
    suggestions,
    loading,
    error: geocodingError,
    search,
    reverseGeocode,
    clearSuggestions,
    clearError,
  } = useAddressGeocoding({
    country,
    language,
    minLength: 2,
    maxSuggestions: 5,
  });

  // Закрытие dropdown при клике вне компонента
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target as Node) &&
        !inputRef.current?.contains(event.target as Node)
      ) {
        setIsOpen(false);
        setHighlightedIndex(-1);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  // Обработка ввода
  const handleInputChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const newValue = e.target.value;
      onChange(newValue);

      if (newValue.trim()) {
        search(newValue);
        setIsOpen(true);
        setHighlightedIndex(-1);
      } else {
        clearSuggestions();
        setIsOpen(false);
      }

      clearError();
    },
    [onChange, search, clearSuggestions, clearError]
  );

  // Выбор предложения
  const selectSuggestion = useCallback(
    (suggestion: AddressGeocodingResult) => {
      onChange(suggestion.address_components.formatted, suggestion);

      if (onLocationSelect) {
        onLocationSelect({
          lat: suggestion.location.lat,
          lng: suggestion.location.lng,
          address: suggestion.address_components.formatted,
          confidence: suggestion.confidence,
        });
      }

      setIsOpen(false);
      setHighlightedIndex(-1);
      inputRef.current?.blur();
    },
    [onChange, onLocationSelect]
  );

  // Обработка клавиш
  const handleKeyDown = useCallback(
    (e: React.KeyboardEvent) => {
      if (!isOpen || suggestions.length === 0) {
        if (e.key === 'ArrowDown') {
          e.preventDefault();
          setIsOpen(true);
        }
        return;
      }

      switch (e.key) {
        case 'ArrowDown':
          e.preventDefault();
          setHighlightedIndex((prev) =>
            prev < suggestions.length - 1 ? prev + 1 : 0
          );
          break;

        case 'ArrowUp':
          e.preventDefault();
          setHighlightedIndex((prev) =>
            prev > 0 ? prev - 1 : suggestions.length - 1
          );
          break;

        case 'Enter':
          e.preventDefault();
          if (highlightedIndex >= 0 && highlightedIndex < suggestions.length) {
            selectSuggestion(suggestions[highlightedIndex]);
          }
          break;

        case 'Escape':
          e.preventDefault();
          setIsOpen(false);
          setHighlightedIndex(-1);
          inputRef.current?.blur();
          break;
      }
    },
    [isOpen, suggestions, highlightedIndex, selectSuggestion]
  );

  // Получение текущего местоположения
  const getCurrentLocation = useCallback(async () => {
    if (!navigator.geolocation) {
      toast.error('Геолокация не поддерживается вашим браузером');
      return;
    }

    setIsGettingLocation(true);

    navigator.geolocation.getCurrentPosition(
      async (position) => {
        try {
          const { latitude, longitude } = position.coords;
          const result = await reverseGeocode(latitude, longitude);

          if (result) {
            onChange(result.address_components.formatted, result);

            if (onLocationSelect) {
              onLocationSelect({
                lat: latitude,
                lng: longitude,
                address: result.address_components.formatted,
                confidence: result.confidence,
              });
            }
          } else {
            toast.warning(
              'Не удалось определить адрес по вашему местоположению'
            );
          }
        } catch (err) {
          console.error('Error getting current location:', err);
          toast.error('Ошибка при определении местоположения');
        } finally {
          setIsGettingLocation(false);
        }
      },
      (error) => {
        console.error('Geolocation error:', error);
        let message = 'Не удалось получить ваше местоположение.';

        switch (error.code) {
          case error.PERMISSION_DENIED:
            message =
              'Доступ к геолокации запрещен. Разрешите доступ в настройках браузера.';
            break;
          case error.POSITION_UNAVAILABLE:
            message = 'Информация о местоположении недоступна.';
            break;
          case error.TIMEOUT:
            message = 'Время ожидания определения местоположения истекло.';
            break;
        }

        toast.error(message);
        setIsGettingLocation(false);
      },
      {
        enableHighAccuracy: true,
        timeout: 10000,
        maximumAge: 300000, // 5 минут
      }
    );
  }, [reverseGeocode, onChange, onLocationSelect]);

  // Определение уровня доверия для цвета
  const getConfidenceColor = (confidence: number) => {
    if (confidence >= 0.8) return 'text-success';
    if (confidence >= 0.6) return 'text-warning';
    return 'text-error';
  };

  // Показать ли dropdown
  const shouldShowDropdown = isOpen && (suggestions.length > 0 || loading);

  return (
    <div className={`relative w-full ${className}`}>
      {/* Основное поле ввода */}
      <div className="relative">
        <input
          ref={inputRef}
          type="text"
          value={value}
          onChange={handleInputChange}
          onKeyDown={handleKeyDown}
          onFocus={() => {
            if (suggestions.length > 0) {
              setIsOpen(true);
            }
          }}
          placeholder={placeholder}
          disabled={disabled}
          required={required}
          className={`
            input input-bordered w-full pr-12
            ${error || geocodingError ? 'input-error' : ''}
            ${loading ? 'opacity-75' : ''}
          `}
          autoComplete="off"
        />

        {/* Кнопка текущего местоположения */}
        {showCurrentLocation && (
          <button
            type="button"
            onClick={getCurrentLocation}
            disabled={disabled || isGettingLocation}
            className="absolute right-2 top-1/2 -translate-y-1/2 btn btn-ghost btn-sm btn-circle"
            title="Использовать мое местоположение"
          >
            {isGettingLocation ? (
              <span className="loading loading-spinner loading-xs"></span>
            ) : (
              <svg
                className="w-4 h-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"
                />
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"
                />
              </svg>
            )}
          </button>
        )}
      </div>

      {/* Dropdown с предложениями */}
      {shouldShowDropdown && (
        <div
          ref={dropdownRef}
          className="absolute z-50 w-full mt-1 bg-base-100 border border-base-300 rounded-lg shadow-lg max-h-64 overflow-y-auto"
        >
          {/* Индикатор загрузки */}
          {loading && (
            <div className="p-3 text-center">
              <span className="loading loading-spinner loading-sm mr-2"></span>
              <span className="text-sm text-base-content/70">
                Поиск адресов...
              </span>
            </div>
          )}

          {/* Список предложений */}
          {suggestions.map((suggestion, index) => (
            <button
              key={suggestion.id}
              type="button"
              onClick={() => selectSuggestion(suggestion)}
              className={`
                w-full p-3 text-left border-b border-base-200 last:border-b-0
                hover:bg-base-200 focus:bg-base-200 focus:outline-none
                ${index === highlightedIndex ? 'bg-base-200' : ''}
              `}
            >
              <div className="flex items-start justify-between">
                <div className="flex-1 min-w-0">
                  {/* Основной текст адреса */}
                  <div className="font-medium text-sm text-base-content truncate">
                    {suggestion.text}
                  </div>

                  {/* Полный адрес */}
                  <div className="text-xs text-base-content/70 mt-1 truncate">
                    {suggestion.place_name}
                  </div>

                  {/* Компоненты адреса */}
                  {suggestion.address_components.city && (
                    <div className="text-xs text-base-content/50 mt-1">
                      {[
                        suggestion.address_components.city,
                        suggestion.address_components.country,
                      ]
                        .filter(Boolean)
                        .join(', ')}
                    </div>
                  )}
                </div>

                {/* Показатель доверия */}
                <div className="ml-2 flex items-center">
                  <div
                    className={`text-xs font-medium ${getConfidenceColor(suggestion.confidence)}`}
                  >
                    {Math.round(suggestion.confidence * 100)}%
                  </div>
                </div>
              </div>
            </button>
          ))}

          {/* Сообщение об отсутствии результатов */}
          {!loading && suggestions.length === 0 && value.trim().length >= 2 && (
            <div className="p-3 text-center text-sm text-base-content/50">
              Адреса не найдены. Попробуйте изменить запрос.
            </div>
          )}
        </div>
      )}

      {/* Сообщения об ошибках */}
      {(error || geocodingError) && (
        <div className="mt-1 text-sm text-error">{error || geocodingError}</div>
      )}

      {/* Вспомогательный текст */}
      {!error &&
        !geocodingError &&
        value.trim() &&
        !loading &&
        suggestions.length === 0 && (
          <div className="mt-1 text-xs text-base-content/50">
            {t('location_picker.min_chars_hint')}
          </div>
        )}
    </div>
  );
}
