'use client';

import React, { useState, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { RadiusSearchControlConfig, GeoSearchResult } from './types/gis';
import {
  useRadiusSearch,
  formatRadius,
  normalizeRadius,
} from './hooks/useRadiusSearch';
import { useGeolocation } from './hooks/useGeolocation';
import { useGeoSearch } from './hooks/useGeoSearch';

interface RadiusSearchControlProps {
  config?: Partial<RadiusSearchControlConfig>;
  onSearch?: (
    results: any[],
    center: { latitude: number; longitude: number },
    radius: number
  ) => void;
  onRadiusChange?: (radius: number) => void;
  onCenterChange?: (center: { latitude: number; longitude: number }) => void;
  onShowCircleToggle?: (show: boolean) => void;
  className?: string;
  disabled?: boolean;
}

const defaultConfig: RadiusSearchControlConfig = {
  minRadius: 0.1,
  maxRadius: 50,
  defaultRadius: 5,
  step: 0.1,
  showMyLocation: true,
  showAddressInput: true,
  showRadiusCircle: true,
  enableGeolocation: true,
};

export default function RadiusSearchControl({
  config = {},
  onSearch,
  onRadiusChange,
  onCenterChange,
  onShowCircleToggle,
  className = '',
  disabled = false,
}: RadiusSearchControlProps) {
  const t = useTranslations('gis.radius_search');
  const mergedConfig = { ...defaultConfig, ...config };

  // Состояние компонента
  const [radius, setRadius] = useState(mergedConfig.defaultRadius);
  const [searchAddress, setSearchAddress] = useState('');
  const [showCircle, setShowCircle] = useState(mergedConfig.showRadiusCircle);
  const [isExpanded, setIsExpanded] = useState(false);
  const [addressSuggestions, setAddressSuggestions] = useState<
    GeoSearchResult[]
  >([]);
  const [showSuggestions, setShowSuggestions] = useState(false);

  // Хуки
  const {
    loading: searchLoading,
    error: searchError,
    total,
    searchCenter,
    searchByAddress,
    searchByCurrentLocation,
    clearResults,
  } = useRadiusSearch();

  const { loading: locationLoading, error: locationError } = useGeolocation();

  const {
    results: _geoResults,
    loading: geoLoading,
    search: geoSearch,
  } = useGeoSearch();

  // Обработчики событий
  const handleRadiusChange = useCallback(
    (newRadius: number) => {
      const validRadius = normalizeRadius(
        newRadius,
        mergedConfig.minRadius,
        mergedConfig.maxRadius
      );
      setRadius(validRadius);
      onRadiusChange?.(validRadius);
    },
    [mergedConfig.minRadius, mergedConfig.maxRadius, onRadiusChange]
  );

  const handleShowCircleToggle = useCallback(() => {
    const newShowCircle = !showCircle;
    setShowCircle(newShowCircle);
    onShowCircleToggle?.(newShowCircle);
  }, [showCircle, onShowCircleToggle]);

  const handleAddressSearch = useCallback(
    async (address: string) => {
      if (!address.trim()) return;

      try {
        const result = await searchByAddress(address, radius);
        if (result) {
          onSearch?.(result.items, result.center, result.radius);
          onCenterChange?.(result.center);
        }
      } catch (error) {
        console.error('Address search failed:', error);
      }
    },
    [radius, searchByAddress, onSearch, onCenterChange]
  );

  const handleCurrentLocationSearch = useCallback(async () => {
    if (!mergedConfig.enableGeolocation) return;

    try {
      const result = await searchByCurrentLocation(radius);
      if (result) {
        onSearch?.(result.items, result.center, result.radius);
        onCenterChange?.(result.center);
      }
    } catch (error) {
      console.error('Current location search failed:', error);
    }
  }, [
    radius,
    searchByCurrentLocation,
    mergedConfig.enableGeolocation,
    onSearch,
    onCenterChange,
  ]);

  // Автодополнение адресов
  const handleAddressInputChange = useCallback(
    async (value: string) => {
      setSearchAddress(value);

      if (value.length >= 3) {
        try {
          const suggestions = await geoSearch({
            query: value,
            limit: 5,
            language: 'ru',
          });
          setAddressSuggestions(suggestions);
          setShowSuggestions(true);
        } catch (error) {
          console.error('Address suggestions failed:', error);
          setAddressSuggestions([]);
        }
      } else {
        setAddressSuggestions([]);
        setShowSuggestions(false);
      }
    },
    [geoSearch]
  );

  const handleSuggestionSelect = useCallback(
    (suggestion: GeoSearchResult) => {
      setSearchAddress(suggestion.display_name);
      setShowSuggestions(false);
      handleAddressSearch(suggestion.display_name);
    },
    [handleAddressSearch]
  );

  const handleSubmit = useCallback(
    (e: React.FormEvent) => {
      e.preventDefault();
      if (searchAddress.trim()) {
        handleAddressSearch(searchAddress);
        setShowSuggestions(false);
      }
    },
    [searchAddress, handleAddressSearch]
  );

  // Отображение ошибок
  const hasError = searchError || locationError;
  const isLoading = searchLoading || locationLoading || geoLoading;

  return (
    <div
      className={`bg-base-100 shadow-lg rounded-lg border border-base-300 ${className}`}
    >
      {/* Заголовок с кнопкой разворачивания */}
      <div className="flex items-center justify-between p-4 border-b border-base-300">
        <div className="flex items-center space-x-2">
          <div className="w-2 h-2 bg-primary rounded-full"></div>
          <h3 className="font-semibold text-sm">{t('title')}</h3>
          {total > 0 && (
            <span className="badge badge-primary badge-sm">{total}</span>
          )}
        </div>
        <button
          type="button"
          onClick={() => setIsExpanded(!isExpanded)}
          className="btn btn-ghost btn-sm btn-square"
          disabled={disabled}
        >
          <svg
            className={`w-4 h-4 transition-transform ${isExpanded ? 'rotate-180' : ''}`}
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M19 9l-7 7-7-7"
            />
          </svg>
        </button>
      </div>

      {/* Основной контент */}
      {isExpanded && (
        <div className="p-4 space-y-4">
          {/* Слайдер радиуса */}
          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <label className="text-sm font-medium text-base-content">
                {t('radius')}
              </label>
              <span className="text-sm font-mono bg-base-200 px-2 py-1 rounded">
                {formatRadius(radius)}
              </span>
            </div>
            <input
              type="range"
              min={mergedConfig.minRadius}
              max={mergedConfig.maxRadius}
              step={mergedConfig.step}
              value={radius}
              onChange={(e) => handleRadiusChange(parseFloat(e.target.value))}
              className="range range-primary range-sm"
              disabled={disabled || isLoading}
            />
            <div className="flex justify-between text-xs text-base-content/60">
              <span>{formatRadius(mergedConfig.minRadius)}</span>
              <span>{formatRadius(mergedConfig.maxRadius)}</span>
            </div>
          </div>

          {/* Поиск по адресу */}
          {mergedConfig.showAddressInput && (
            <div className="space-y-2">
              <label className="text-sm font-medium text-base-content">
                {t('address')}
              </label>
              <form onSubmit={handleSubmit} className="relative">
                <div className="flex space-x-2">
                  <div className="relative flex-1">
                    <input
                      type="text"
                      value={searchAddress}
                      onChange={(e) => handleAddressInputChange(e.target.value)}
                      placeholder={t('address_placeholder')}
                      className="input input-bordered input-sm w-full pr-10"
                      disabled={disabled || isLoading}
                    />
                    {geoLoading && (
                      <div className="absolute right-3 top-1/2 transform -translate-y-1/2">
                        <div className="loading loading-spinner loading-xs"></div>
                      </div>
                    )}
                  </div>
                  <button
                    type="submit"
                    className="btn btn-primary btn-sm"
                    disabled={disabled || isLoading || !searchAddress.trim()}
                  >
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
                        d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                      />
                    </svg>
                  </button>
                </div>

                {/* Предложения адресов */}
                {showSuggestions && addressSuggestions.length > 0 && (
                  <div className="absolute top-full left-0 right-0 z-50 mt-1 bg-base-100 border border-base-300 rounded-lg shadow-lg max-h-60 overflow-y-auto">
                    {addressSuggestions.map((suggestion) => (
                      <button
                        key={suggestion.id}
                        type="button"
                        onClick={() => handleSuggestionSelect(suggestion)}
                        className="w-full text-left px-3 py-2 hover:bg-base-200 first:rounded-t-lg last:rounded-b-lg"
                      >
                        <div className="text-sm text-base-content">
                          {suggestion.display_name}
                        </div>
                      </button>
                    ))}
                  </div>
                )}
              </form>
            </div>
          )}

          {/* Кнопка "Мое местоположение" */}
          {mergedConfig.showMyLocation && mergedConfig.enableGeolocation && (
            <button
              type="button"
              onClick={handleCurrentLocationSearch}
              className="btn btn-outline btn-sm w-full"
              disabled={disabled || isLoading}
            >
              {locationLoading ? (
                <div className="loading loading-spinner loading-xs"></div>
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
              {t('use_my_location')}
            </button>
          )}

          {/* Переключатель показа круга */}
          <div className="flex items-center justify-between">
            <label className="text-sm font-medium text-base-content">
              {t('show_search_area')}
            </label>
            <input
              type="checkbox"
              checked={showCircle}
              onChange={handleShowCircleToggle}
              className="checkbox checkbox-primary checkbox-sm"
              disabled={disabled}
            />
          </div>

          {/* Индикатор загрузки */}
          {isLoading && (
            <div className="flex items-center justify-center space-x-2 p-2">
              <div className="loading loading-spinner loading-sm"></div>
              <span className="text-sm text-base-content/70">
                {t('searching')}
              </span>
            </div>
          )}

          {/* Ошибки */}
          {hasError && (
            <div className="alert alert-error alert-sm">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                className="stroke-current shrink-0 h-4 w-4"
                fill="none"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              <span className="text-xs">
                {searchError || locationError?.message}
              </span>
            </div>
          )}

          {/* Результаты поиска */}
          {total > 0 && searchCenter && (
            <div className="space-y-2">
              <div className="text-sm text-base-content/70">
                {t('results_found', {
                  count: total,
                  radius: formatRadius(radius),
                })}
              </div>
              <button
                type="button"
                onClick={clearResults}
                className="btn btn-ghost btn-xs"
              >
                {t('clear_results')}
              </button>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
