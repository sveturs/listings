'use client';

import React, {
  useState,
  useEffect,
  useCallback,
  useRef,
  useMemo,
} from 'react';
import { useTranslations } from 'next-intl';
import { debounce } from 'lodash';
import {
  MagnifyingGlassIcon,
  ClockIcon,
  TrendingUpIcon,
  SparklesIcon,
  ArrowTrendingUpIcon,
  UserIcon,
} from '@heroicons/react/24/outline';

interface QuerySuggestion {
  query: string;
  type: 'recent' | 'popular' | 'trending' | 'personalized' | 'fuzzy';
  score: number;
  attributes?: string[];
  highlight?: string;
  metadata?: {
    frequency?: number;
    ctr?: number;
    conversions?: number;
    trend_score?: number;
    distance?: number;
    user_frequency?: number;
    last_used?: string;
  };
}

interface QuerySuggestionsProps {
  value: string;
  onChange: (value: string) => void;
  onSelect?: (suggestion: QuerySuggestion) => void;
  placeholder?: string;
  className?: string;
  autoFocus?: boolean;
  categoryId?: number;
  userId?: number;
  disabled?: boolean;
}

export default function QuerySuggestions({
  value,
  onChange,
  onSelect,
  placeholder,
  className = '',
  autoFocus = false,
  categoryId,
  userId,
  disabled = false,
}: QuerySuggestionsProps) {
  const t = useTranslations('search');
  const [suggestions, setSuggestions] = useState<QuerySuggestion[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [showSuggestions, setShowSuggestions] = useState(false);
  const [selectedIndex, setSelectedIndex] = useState(-1);
  const inputRef = useRef<HTMLInputElement>(null);
  const suggestionsRef = useRef<HTMLDivElement>(null);

  // Fetch suggestions from API
  const fetchSuggestions = useCallback(
    async (query: string) => {
      if (query.length < 2) {
        setSuggestions([]);
        return;
      }

      setIsLoading(true);

      try {
        const params = new URLSearchParams({
          q: query,
          limit: '10',
          enable_fuzzy: 'true',
          include_trending: 'true',
        });

        if (categoryId) {
          params.set('category_id', categoryId.toString());
        }

        if (userId) {
          params.set('user_id', userId.toString());
        }

        const response = await fetch(`/api/v1/suggestions?${params}`);

        if (response.ok) {
          const data = await response.json();
          setSuggestions(data.data || []);
        }
      } catch (error) {
        console.error('Failed to fetch suggestions:', error);
        setSuggestions([]);
      } finally {
        setIsLoading(false);
      }
    },
    [categoryId, userId]
  );

  // Debounced fetch
  const debouncedFetch = useMemo(
    () => debounce(fetchSuggestions, 200),
    [fetchSuggestions]
  );

  // Handle input change
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newValue = e.target.value;
    onChange(newValue);
    setSelectedIndex(-1);

    if (newValue.length >= 2) {
      setShowSuggestions(true);
      debouncedFetch(newValue);
    } else {
      setShowSuggestions(false);
      setSuggestions([]);
    }
  };

  // Handle keyboard navigation
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (!showSuggestions || suggestions.length === 0) return;

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
        if (selectedIndex >= 0 && selectedIndex < suggestions.length) {
          handleSuggestionClick(suggestions[selectedIndex]);
        }
        break;

      case 'Escape':
        e.preventDefault();
        setShowSuggestions(false);
        setSelectedIndex(-1);
        break;
    }
  };

  // Handle suggestion click
  const handleSuggestionClick = (suggestion: QuerySuggestion) => {
    onChange(suggestion.query);
    setShowSuggestions(false);
    setSelectedIndex(-1);

    if (onSelect) {
      onSelect(suggestion);
    }

    // Track suggestion click for learning
    trackSuggestionClick(suggestion);
  };

  // Track suggestion click
  const trackSuggestionClick = async (suggestion: QuerySuggestion) => {
    try {
      await fetch('/api/v1/suggestions/track', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          query: suggestion.query,
          type: 'click',
          user_id: userId,
          category_id: categoryId,
          suggestion_type: suggestion.type,
          position: suggestions.indexOf(suggestion),
        }),
      });
    } catch (error) {
      console.error('Failed to track suggestion click:', error);
    }
  };

  // Get icon for suggestion type
  const getSuggestionIcon = (type: string) => {
    switch (type) {
      case 'recent':
        return <ClockIcon className="w-4 h-4 text-base-content/50" />;
      case 'trending':
        return <TrendingUpIcon className="w-4 h-4 text-success" />;
      case 'personalized':
        return <UserIcon className="w-4 h-4 text-primary" />;
      case 'popular':
        return <ArrowTrendingUpIcon className="w-4 h-4 text-info" />;
      case 'fuzzy':
        return <SparklesIcon className="w-4 h-4 text-warning" />;
      default:
        return <MagnifyingGlassIcon className="w-4 h-4 text-base-content/50" />;
    }
  };

  // Get suggestion label
  const getSuggestionLabel = (type: string) => {
    switch (type) {
      case 'recent':
        return t('suggestions.recent');
      case 'trending':
        return t('suggestions.trending');
      case 'personalized':
        return t('suggestions.personalized');
      case 'popular':
        return t('suggestions.popular');
      case 'fuzzy':
        return t('suggestions.didYouMean');
      default:
        return '';
    }
  };

  // Highlight matching text
  const highlightText = (text: string, highlight?: string) => {
    if (!highlight) return text;

    // Parse highlight HTML
    const parser = new DOMParser();
    const doc = parser.parseFromString(highlight, 'text/html');
    const markElement = doc.querySelector('mark');

    if (!markElement) return text;

    const highlightedText = markElement.textContent || '';
    const parts = text.toLowerCase().split(highlightedText.toLowerCase());

    if (parts.length <= 1) return text;

    return (
      <>
        {parts[0]}
        <mark className="bg-primary/20 text-primary font-semibold">
          {text.substr(parts[0].length, highlightedText.length)}
        </mark>
        {text.substr(parts[0].length + highlightedText.length)}
      </>
    );
  };

  // Close suggestions on outside click
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        suggestionsRef.current &&
        !suggestionsRef.current.contains(event.target as Node) &&
        inputRef.current &&
        !inputRef.current.contains(event.target as Node)
      ) {
        setShowSuggestions(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  // Group suggestions by type
  const groupedSuggestions = useMemo(() => {
    const groups: Record<string, QuerySuggestion[]> = {};

    suggestions.forEach((suggestion) => {
      if (!groups[suggestion.type]) {
        groups[suggestion.type] = [];
      }
      groups[suggestion.type].push(suggestion);
    });

    // Order groups by priority
    const orderedTypes = [
      'personalized',
      'recent',
      'trending',
      'popular',
      'fuzzy',
    ];
    const orderedGroups: Array<[string, QuerySuggestion[]]> = [];

    orderedTypes.forEach((type) => {
      if (groups[type]) {
        orderedGroups.push([type, groups[type]]);
      }
    });

    return orderedGroups;
  }, [suggestions]);

  return (
    <div className="relative w-full">
      {/* Search Input */}
      <div className="relative">
        <input
          ref={inputRef}
          type="text"
          value={value}
          onChange={handleInputChange}
          onKeyDown={handleKeyDown}
          onFocus={() => value.length >= 2 && setShowSuggestions(true)}
          placeholder={placeholder || t('search.placeholder')}
          className={`input input-bordered w-full pr-10 ${className}`}
          autoFocus={autoFocus}
          disabled={disabled}
        />
        <div className="absolute right-3 top-1/2 -translate-y-1/2">
          {isLoading ? (
            <span className="loading loading-spinner loading-sm"></span>
          ) : (
            <MagnifyingGlassIcon className="w-5 h-5 text-base-content/50" />
          )}
        </div>
      </div>

      {/* Suggestions Dropdown */}
      {showSuggestions && suggestions.length > 0 && (
        <div
          ref={suggestionsRef}
          className="absolute top-full left-0 right-0 mt-1 bg-base-100 border border-base-200 rounded-lg shadow-lg z-50 max-h-96 overflow-y-auto"
        >
          {groupedSuggestions.map(([type, items], groupIndex) => (
            <div
              key={type}
              className={groupIndex > 0 ? 'border-t border-base-200' : ''}
            >
              {/* Group Header */}
              <div className="px-3 py-2 flex items-center gap-2 bg-base-200/50">
                {getSuggestionIcon(type)}
                <span className="text-xs font-medium text-base-content/70">
                  {getSuggestionLabel(type)}
                </span>
              </div>

              {/* Group Items */}
              {items.map((suggestion, itemIndex) => {
                const globalIndex = suggestions.indexOf(suggestion);
                const isSelected = selectedIndex === globalIndex;

                return (
                  <div
                    key={`${type}-${itemIndex}`}
                    className={`
                      px-4 py-2 cursor-pointer transition-colors
                      ${isSelected ? 'bg-primary/10' : 'hover:bg-base-200'}
                    `}
                    onClick={() => handleSuggestionClick(suggestion)}
                    onMouseEnter={() => setSelectedIndex(globalIndex)}
                  >
                    <div className="flex items-start justify-between">
                      <div className="flex-1">
                        {/* Query Text */}
                        <div className="font-medium">
                          {highlightText(
                            suggestion.query,
                            suggestion.highlight
                          )}
                        </div>

                        {/* Attributes */}
                        {suggestion.attributes &&
                          suggestion.attributes.length > 0 && (
                            <div className="flex flex-wrap gap-1 mt-1">
                              {suggestion.attributes
                                .slice(0, 3)
                                .map((attr, idx) => (
                                  <span
                                    key={idx}
                                    className="badge badge-xs badge-ghost"
                                  >
                                    {attr}
                                  </span>
                                ))}
                            </div>
                          )}
                      </div>

                      {/* Metadata */}
                      <div className="text-xs text-base-content/50 ml-2">
                        {suggestion.metadata?.frequency && (
                          <span>
                            {suggestion.metadata.frequency}{' '}
                            {t('suggestions.searches')}
                          </span>
                        )}
                        {suggestion.metadata?.trend_score && (
                          <span className="text-success">
                            +{Math.round(suggestion.metadata.trend_score)}%{' '}
                            {t('suggestions.trending')}
                          </span>
                        )}
                      </div>
                    </div>
                  </div>
                );
              })}
            </div>
          ))}

          {/* Footer */}
          <div className="px-3 py-2 bg-base-200/30 text-xs text-base-content/50 border-t border-base-200">
            {t('suggestions.hint')}
          </div>
        </div>
      )}

      {/* No Suggestions */}
      {showSuggestions &&
        !isLoading &&
        value.length >= 2 &&
        suggestions.length === 0 && (
          <div className="absolute top-full left-0 right-0 mt-1 bg-base-100 border border-base-200 rounded-lg shadow-lg z-50 p-4 text-center text-base-content/50">
            {t('suggestions.noResults')}
          </div>
        )}
    </div>
  );
}
