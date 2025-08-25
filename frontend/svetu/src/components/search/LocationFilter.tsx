'use client';

import { useState, useEffect, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { LocationService, City } from '@/services/location';
import { MapPin } from 'lucide-react';

interface LocationFilterProps {
  location: string;
  radius: number;
  onLocationChange: (location: string, radius: number) => void;
}

export default function LocationFilter({
  location,
  radius,
  onLocationChange,
}: LocationFilterProps) {
  const t = useTranslations('search');
  const [inputValue, setInputValue] = useState(location);
  const [suggestions, setSuggestions] = useState<City[]>([]);
  const [showSuggestions, setShowSuggestions] = useState(false);
  const [loading, setLoading] = useState(false);
  const [searchTimeout, setSearchTimeout] = useState<NodeJS.Timeout | null>(
    null
  );

  useEffect(() => {
    setInputValue(location);
  }, [location]);

  const searchCities = useCallback(async (query: string) => {
    if (!query.trim()) {
      setSuggestions([]);
      return;
    }

    setLoading(true);
    try {
      const cities = await LocationService.searchCities(query);
      setSuggestions(cities);
      setShowSuggestions(true);
    } catch (error) {
      console.error('Failed to search cities:', error);
      setSuggestions([]);
    } finally {
      setLoading(false);
    }
  }, []);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setInputValue(value);

    // Clear previous timeout
    if (searchTimeout) {
      clearTimeout(searchTimeout);
    }

    if (!value.trim()) {
      onLocationChange('', radius);
      setSuggestions([]);
      setShowSuggestions(false);
      return;
    }

    // Debounce search
    const timeout = setTimeout(() => {
      searchCities(value);
    }, 300);
    setSearchTimeout(timeout);
  };

  const handleCitySelect = (city: City) => {
    setInputValue(city.name);
    onLocationChange(city.name, radius);
    setShowSuggestions(false);
    setSuggestions([]);
  };

  const handleRadiusChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newRadius = parseInt(e.target.value, 10);
    onLocationChange(inputValue, newRadius);
  };

  const handleFocus = () => {
    if (suggestions.length > 0) {
      setShowSuggestions(true);
    }
  };

  const handleBlur = () => {
    // Delay to allow click on suggestion
    setTimeout(() => {
      setShowSuggestions(false);
    }, 200);
  };

  return (
    <div className="form-control">
      <label className="label" htmlFor="location-input">
        <span className="label-text font-medium flex items-center gap-2">
          <MapPin className="w-4 h-4 text-primary" />
          {t('location')}
        </span>
      </label>

      <div className="relative">
        <input
          id="location-input"
          type="text"
          className="input input-bordered w-full"
          placeholder={t('enterCity')}
          value={inputValue}
          onChange={handleInputChange}
          onFocus={handleFocus}
          onBlur={handleBlur}
        />

        {loading && (
          <div className="absolute right-2 top-1/2 -translate-y-1/2">
            <span className="loading loading-spinner loading-sm"></span>
          </div>
        )}

        {showSuggestions && suggestions.length > 0 && (
          <div className="absolute z-10 w-full mt-1 bg-base-100 border border-base-300 rounded-lg shadow-lg max-h-60 overflow-y-auto">
            {suggestions.map((city) => (
              <button
                key={city.id}
                className="w-full text-left px-4 py-2 hover:bg-base-200 focus:bg-base-200 transition-colors"
                onClick={() => handleCitySelect(city)}
                type="button"
              >
                <div className="font-medium">{city.name}</div>
                {city.region && (
                  <div className="text-sm text-base-content/60">
                    {city.region}, {city.country}
                  </div>
                )}
              </button>
            ))}
          </div>
        )}
      </div>

      {inputValue && (
        <div className="mt-4">
          <label className="label">
            <span className="label-text text-sm">
              {t('searchRadius')}:{' '}
              <span className="font-semibold">{radius} km</span>
            </span>
          </label>
          <input
            type="range"
            min="1"
            max="100"
            value={radius}
            onChange={handleRadiusChange}
            className="range range-primary range-sm"
            aria-label={t('searchRadius')}
          />
          <div className="w-full flex justify-between text-xs px-2 mt-1">
            <span>1 km</span>
            <span>25 km</span>
            <span>50 km</span>
            <span>75 km</span>
            <span>100 km</span>
          </div>
        </div>
      )}
    </div>
  );
}
