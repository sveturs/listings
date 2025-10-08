'use client';

import React, { useState, useEffect, useRef, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { Search, X, Clock, TrendingUp, MapPin, Car } from 'lucide-react';
import { useDebounce } from '@/hooks/useDebounce';

interface SearchSuggestion {
  type: 'make' | 'model' | 'location' | 'recent' | 'popular';
  value: string;
  label: string;
  count?: number;
  icon?: React.ReactNode;
}

interface AutocompleteSearchProps {
  onSearch: (query: string, filters?: any) => void;
  placeholder?: string;
  className?: string;
}

const RECENT_SEARCHES_KEY = 'car-recent-searches';
const MAX_RECENT_SEARCHES = 5;
const POPULAR_SEARCHES = [
  { make: 'toyota', label: 'Toyota', count: 150 },
  { make: 'volkswagen', label: 'Volkswagen', count: 120 },
  { make: 'bmw', label: 'BMW', count: 100 },
  { make: 'mercedes-benz', label: 'Mercedes-Benz', count: 95 },
  { make: 'audi', label: 'Audi', count: 85 },
];

export default function AutocompleteSearch({
  onSearch,
  placeholder,
  className = '',
}: AutocompleteSearchProps) {
  const t = useTranslations('cars');
  const [query, setQuery] = useState('');
  const [suggestions, setSuggestions] = useState<SearchSuggestion[]>([]);
  const [showSuggestions, setShowSuggestions] = useState(false);
  const [selectedIndex, setSelectedIndex] = useState(-1);
  const [recentSearches, setRecentSearches] = useState<string[]>([]);
  const debouncedQuery = useDebounce(query, 300);
  const inputRef = useRef<HTMLInputElement>(null);
  const suggestionsRef = useRef<HTMLDivElement>(null);

  // Load recent searches from localStorage
  useEffect(() => {
    const stored = localStorage.getItem(RECENT_SEARCHES_KEY);
    if (stored) {
      try {
        setRecentSearches(JSON.parse(stored));
      } catch (e) {
        console.error('Failed to load recent searches:', e);
      }
    }
  }, []);

  // Save search to recent
  const saveToRecent = useCallback(
    (searchQuery: string) => {
      if (!searchQuery.trim()) return;

      const updated = [
        searchQuery,
        ...recentSearches.filter((s) => s !== searchQuery),
      ].slice(0, MAX_RECENT_SEARCHES);

      setRecentSearches(updated);
      localStorage.setItem(RECENT_SEARCHES_KEY, JSON.stringify(updated));
    },
    [recentSearches]
  );

  // Clear recent searches
  const clearRecentSearches = useCallback(() => {
    setRecentSearches([]);
    localStorage.removeItem(RECENT_SEARCHES_KEY);
  }, []);

  // Fetch suggestions based on query
  useEffect(() => {
    const fetchSuggestions = async () => {
      if (!debouncedQuery.trim()) {
        // Show recent and popular when no query
        const recentSuggestions: SearchSuggestion[] = recentSearches.map(
          (s) => ({
            type: 'recent',
            value: s,
            label: s,
            icon: <Clock className="w-4 h-4" />,
          })
        );

        const popularSuggestions: SearchSuggestion[] = POPULAR_SEARCHES.map(
          (p) => ({
            type: 'popular',
            value: p.make,
            label: p.label,
            count: p.count,
            icon: <TrendingUp className="w-4 h-4" />,
          })
        );

        setSuggestions([...recentSuggestions, ...popularSuggestions]);
        return;
      }

      try {
        // Fetch car makes/models matching query
        const makesResponse = await fetch('/api/v1/c2c/cars/makes');
        const makesData = await makesResponse.json();

        const matchingMakes: SearchSuggestion[] =
          makesData.data
            ?.filter((make: any) =>
              make.name.toLowerCase().includes(debouncedQuery.toLowerCase())
            )
            .slice(0, 3)
            .map((make: any) => ({
              type: 'make',
              value: make.slug,
              label: make.name,
              count: make.count,
              icon: <Car className="w-4 h-4" />,
            })) || [];

        // Add location suggestions (mock for now)
        const locations: SearchSuggestion[] = [
          'Belgrade',
          'Novi Sad',
          'Niš',
          'Kragujevac',
        ]
          .filter((city) =>
            city.toLowerCase().includes(debouncedQuery.toLowerCase())
          )
          .slice(0, 2)
          .map((city) => ({
            type: 'location',
            value: city.toLowerCase(),
            label: city,
            icon: <MapPin className="w-4 h-4" />,
          }));

        setSuggestions([...matchingMakes, ...locations]);
      } catch (error) {
        console.error('Failed to fetch suggestions:', error);
        setSuggestions([]);
      }
    };

    fetchSuggestions();
  }, [debouncedQuery, recentSearches]);

  // Handle keyboard navigation
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (!showSuggestions || suggestions.length === 0) {
      if (e.key === 'Enter' && query.trim()) {
        handleSearch(query);
      }
      return;
    }

    switch (e.key) {
      case 'ArrowDown':
        e.preventDefault();
        setSelectedIndex((prev) =>
          prev < suggestions.length - 1 ? prev + 1 : 0
        );
        break;
      case 'ArrowUp':
        e.preventDefault();
        setSelectedIndex((prev) =>
          prev > 0 ? prev - 1 : suggestions.length - 1
        );
        break;
      case 'Enter':
        e.preventDefault();
        if (selectedIndex >= 0) {
          handleSuggestionClick(suggestions[selectedIndex]);
        } else if (query.trim()) {
          handleSearch(query);
        }
        break;
      case 'Escape':
        setShowSuggestions(false);
        setSelectedIndex(-1);
        break;
    }
  };

  // Handle search
  const handleSearch = (searchQuery: string) => {
    if (!searchQuery.trim()) return;

    saveToRecent(searchQuery);
    onSearch(searchQuery);
    setShowSuggestions(false);
    setQuery('');
  };

  // Handle suggestion click
  const handleSuggestionClick = (suggestion: SearchSuggestion) => {
    switch (suggestion.type) {
      case 'make':
        onSearch('', { make: suggestion.value });
        saveToRecent(suggestion.label);
        break;
      case 'location':
        onSearch('', { location: suggestion.value });
        saveToRecent(suggestion.label);
        break;
      case 'recent':
      case 'popular':
        setQuery(suggestion.value);
        handleSearch(suggestion.value);
        break;
      default:
        handleSearch(suggestion.value);
    }
    setShowSuggestions(false);
    setQuery('');
  };

  // Click outside handler
  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (
        inputRef.current &&
        !inputRef.current.contains(e.target as Node) &&
        suggestionsRef.current &&
        !suggestionsRef.current.contains(e.target as Node)
      ) {
        setShowSuggestions(false);
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
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          onFocus={() => setShowSuggestions(true)}
          onKeyDown={handleKeyDown}
          placeholder={placeholder || t('search.placeholder')}
          className="w-full pl-10 pr-10 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
        />
        <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
        {query && (
          <button
            onClick={() => {
              setQuery('');
              inputRef.current?.focus();
            }}
            className="absolute right-3 top-1/2 transform -translate-y-1/2 p-1 hover:bg-gray-100 rounded"
          >
            <X className="w-4 h-4 text-gray-400" />
          </button>
        )}
      </div>

      {/* Suggestions dropdown */}
      {showSuggestions && suggestions.length > 0 && (
        <div
          ref={suggestionsRef}
          className="absolute top-full left-0 right-0 mt-1 bg-white border border-gray-200 rounded-lg shadow-lg z-50 max-h-96 overflow-y-auto"
        >
          {/* Recent searches section */}
          {!query && recentSearches.length > 0 && (
            <div className="border-b border-gray-100">
              <div className="px-4 py-2 flex justify-between items-center">
                <span className="text-sm text-gray-500 font-medium">
                  {t('search.recent')}
                </span>
                <button
                  onClick={clearRecentSearches}
                  className="text-xs text-primary hover:underline"
                >
                  {t('search.clearRecent')}
                </button>
              </div>
            </div>
          )}

          {/* Popular searches section */}
          {!query && suggestions.some((s) => s.type === 'popular') && (
            <div className="px-4 py-2 border-b border-gray-100">
              <span className="text-sm text-gray-500 font-medium">
                {t('search.popular')}
              </span>
            </div>
          )}

          {/* Suggestions list */}
          <div className="py-1">
            {suggestions.map((suggestion, index) => (
              <button
                key={`${suggestion.type}-${suggestion.value}-${index}`}
                onClick={() => handleSuggestionClick(suggestion)}
                onMouseEnter={() => setSelectedIndex(index)}
                className={`w-full px-4 py-2 text-left flex items-center justify-between hover:bg-gray-50 ${
                  selectedIndex === index ? 'bg-gray-50' : ''
                }`}
              >
                <div className="flex items-center gap-3">
                  <span className="text-gray-400">{suggestion.icon}</span>
                  <span className="text-gray-700">{suggestion.label}</span>
                </div>
                {suggestion.count && (
                  <span className="text-sm text-gray-400">
                    {suggestion.count}
                  </span>
                )}
              </button>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
