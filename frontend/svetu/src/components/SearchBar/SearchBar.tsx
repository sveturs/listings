'use client';

import { useState, useEffect, useRef, useCallback } from 'react';
import { useTranslations, useLocale } from 'next-intl';
import { useRouter } from 'next/navigation';
import { useDebounce } from '@/hooks/useDebounce';
import {
  UnifiedSearchService,
  SearchSuggestion,
} from '@/services/unifiedSearch';
import { SearchIcon, CloseIcon, TrendingIcon } from './icons';

interface SearchBarProps {
  placeholder?: string;
  initialQuery?: string;
  onSearch?: (query: string) => void;
  className?: string;
  variant?: 'default' | 'hero' | 'minimal';
  showTrending?: boolean;
}

export default function SearchBar({
  placeholder,
  initialQuery = '',
  onSearch,
  className = '',
  variant = 'default',
  showTrending = false,
}: SearchBarProps) {
  const t = useTranslations('search');
  const locale = useLocale();
  const router = useRouter();
  const [query, setQuery] = useState(initialQuery);
  const [suggestions, setSuggestions] = useState<SearchSuggestion[]>([]);
  const [isLoadingSuggestions, setIsLoadingSuggestions] = useState(false);
  const [showSuggestions, setShowSuggestions] = useState(false);
  const [selectedIndex, setSelectedIndex] = useState(-1);
  const [searchHistory, setSearchHistory] = useState<string[]>([]);
  const [trendingSearches, setTrendingSearches] = useState<string[]>([]);
  const debouncedQuery = useDebounce(query, 300);
  const inputRef = useRef<HTMLInputElement>(null);
  const suggestionRef = useRef<HTMLDivElement>(null);

  // Загрузка истории и трендов при монтировании
  useEffect(() => {
    UnifiedSearchService.getSearchHistory().then(setSearchHistory);
    if (showTrending) {
      UnifiedSearchService.getTrendingSearches().then(setTrendingSearches);
    }
  }, [showTrending]);

  // Загрузка предложений
  useEffect(() => {
    if (debouncedQuery.length >= 2) {
      loadSuggestions(debouncedQuery);
    } else if (
      debouncedQuery.length === 0 &&
      (searchHistory.length > 0 || trendingSearches.length > 0)
    ) {
      // Показываем историю и тренды когда поле пустое
      setShowSuggestions(true);
      setSuggestions([]);
    } else {
      setSuggestions([]);
      setShowSuggestions(false);
    }
  }, [debouncedQuery, searchHistory, trendingSearches]);

  // Закрытие предложений при клике вне
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        suggestionRef.current &&
        !suggestionRef.current.contains(event.target as Node) &&
        !inputRef.current?.contains(event.target as Node)
      ) {
        setShowSuggestions(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const loadSuggestions = async (searchQuery: string) => {
    setIsLoadingSuggestions(true);
    try {
      const results = await UnifiedSearchService.getSuggestions(searchQuery);
      setSuggestions(results);
      setShowSuggestions(true);
      setSelectedIndex(-1);
    } catch (error) {
      console.error('Error loading suggestions:', error);
      setSuggestions([]);
    } finally {
      setIsLoadingSuggestions(false);
    }
  };

  const handleSearch = useCallback(
    (searchQuery: string = query) => {
      const trimmedQuery = searchQuery.trim();
      if (!trimmedQuery) return;

      setShowSuggestions(false);
      UnifiedSearchService.saveToHistory(trimmedQuery);

      if (onSearch) {
        onSearch(trimmedQuery);
      } else {
        // Строим URL с правильной локалью
        const searchUrl = `/${locale}/search?q=${encodeURIComponent(trimmedQuery)}`;
        router.push(searchUrl);
      }
    },
    [query, onSearch, router, locale]
  );

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (
      !showSuggestions ||
      (suggestions.length === 0 && searchHistory.length === 0)
    ) {
      if (e.key === 'Enter') {
        handleSearch();
      }
      return;
    }

    const totalItems = suggestions.length + searchHistory.length;

    switch (e.key) {
      case 'ArrowDown':
        e.preventDefault();
        setSelectedIndex((prev) => (prev < totalItems - 1 ? prev + 1 : prev));
        break;
      case 'ArrowUp':
        e.preventDefault();
        setSelectedIndex((prev) => (prev > 0 ? prev - 1 : -1));
        break;
      case 'Enter':
        e.preventDefault();
        if (selectedIndex >= 0) {
          if (selectedIndex < suggestions.length) {
            handleSuggestionClick(suggestions[selectedIndex]);
          } else {
            const historyIndex = selectedIndex - suggestions.length;
            handleSearch(searchHistory[historyIndex]);
          }
        } else {
          handleSearch();
        }
        break;
      case 'Escape':
        setShowSuggestions(false);
        setSelectedIndex(-1);
        inputRef.current?.blur();
        break;
    }
  };

  const handleSuggestionClick = (suggestion: SearchSuggestion) => {
    setShowSuggestions(false);

    // Если это категория, переходим на страницу категории
    if (suggestion.type === 'category' && suggestion.category) {
      router.push(
        `/${locale}/category/${suggestion.category.slug || suggestion.category.id}`
      );
      return;
    }

    // Если это товар, можно перейти на страницу товара (если есть id)
    if (suggestion.type === 'product' && suggestion.product_id) {
      router.push(`/${locale}/listing/${suggestion.product_id}`);
      return;
    }

    // Иначе выполняем обычный поиск
    setQuery(suggestion.text);
    handleSearch(suggestion.text);
  };

  const handleClear = () => {
    setQuery('');
    setSuggestions([]);
    setShowSuggestions(false);
    inputRef.current?.focus();
  };

  const renderSuggestionIcon = (type: string) => {
    switch (type) {
      case 'category':
        return (
          <svg
            className="w-4 h-4 text-primary"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.99 1.99 0 013 12V7a4 4 0 014-4z"
            />
          </svg>
        );
      case 'product':
        return (
          <svg
            className="w-4 h-4 text-secondary"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"
            />
          </svg>
        );
      default:
        return <SearchIcon className="w-4 h-4 text-base-content/60" />;
    }
  };

  const getSearchBarClasses = () => {
    const base = 'relative';
    switch (variant) {
      case 'hero':
        return `${base} w-full max-w-2xl mx-auto`;
      case 'minimal':
        return `${base} w-full`;
      default:
        return `${base} ${className}`;
    }
  };

  const getInputClasses = () => {
    const base = 'input input-bordered w-full pr-10';
    switch (variant) {
      case 'hero':
        return `${base} input-lg text-lg shadow-xl hover:shadow-2xl transition-shadow duration-300 border-base-300/50 focus:border-primary`;
      case 'minimal':
        return `${base} input-sm`;
      default:
        return base;
    }
  };

  return (
    <div className={getSearchBarClasses()}>
      <div className="form-control">
        <div className="relative">
          <input
            ref={inputRef}
            type="text"
            placeholder={placeholder || t('placeholder')}
            className={getInputClasses()}
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            onKeyDown={handleKeyDown}
            onFocus={() => {
              if (
                suggestions.length > 0 ||
                searchHistory.length > 0 ||
                trendingSearches.length > 0
              ) {
                setShowSuggestions(true);
              }
            }}
          />

          {/* Иконка поиска или загрузки */}
          <div className="absolute right-3 top-1/2 -translate-y-1/2 flex items-center gap-2">
            {query && (
              <button
                type="button"
                className="btn btn-ghost btn-xs btn-circle"
                onClick={handleClear}
              >
                <CloseIcon className="w-4 h-4" />
              </button>
            )}

            <button
              type="button"
              className={`btn btn-ghost btn-xs btn-circle ${!query.trim() ? 'btn-disabled' : ''}`}
              onClick={() => handleSearch()}
              disabled={!query.trim()}
            >
              {isLoadingSuggestions ? (
                <span className="loading loading-spinner loading-xs"></span>
              ) : (
                <SearchIcon className="w-4 h-4" />
              )}
            </button>
          </div>
        </div>
      </div>

      {/* Предложения автодополнения */}
      {showSuggestions && (
        <div
          ref={suggestionRef}
          className="absolute top-full left-0 right-0 bg-base-100/95 backdrop-blur-md border border-base-300/50 rounded-xl mt-2 shadow-2xl z-[100] max-h-96 overflow-y-auto"
        >
          {/* Результаты поиска */}
          {suggestions.length > 0 && (
            <div className="p-2">
              {suggestions.map((suggestion, index) => (
                <div
                  key={`suggestion-${index}`}
                  className={`flex items-center gap-3 px-3 py-2 rounded-lg cursor-pointer transition-all duration-150 ${
                    index === selectedIndex
                      ? 'bg-base-200'
                      : 'hover:bg-base-200'
                  } active:scale-[0.98]`}
                  onClick={() => handleSuggestionClick(suggestion)}
                >
                  {renderSuggestionIcon(suggestion.type)}
                  <div className="flex-1 min-w-0">
                    <div className="text-sm font-medium text-base-content truncate">
                      {suggestion.text}
                    </div>
                    {suggestion.category && (
                      <div className="text-xs text-base-content/60 truncate">
                        {t('inCategory')}: {suggestion.category.name}
                      </div>
                    )}
                  </div>
                  {suggestion.type === 'category' && (
                    <div className="flex items-center gap-2">
                      <div className="badge badge-primary badge-sm">
                        {t('category')}
                      </div>
                      <svg
                        className="w-3 h-3 text-base-content/40"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M9 5l7 7-7 7"
                        />
                      </svg>
                    </div>
                  )}
                  {suggestion.type === 'product' && (
                    <svg
                      className="w-3 h-3 text-base-content/40"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M9 5l7 7-7 7"
                      />
                    </svg>
                  )}
                </div>
              ))}
            </div>
          )}

          {/* История поиска */}
          {suggestions.length === 0 && searchHistory.length > 0 && (
            <div className="p-2">
              <div className="px-3 py-2 text-xs font-semibold text-base-content/60 uppercase">
                История поиска
              </div>
              {searchHistory.slice(0, 5).map((historyItem, index) => (
                <div
                  key={`history-${index}`}
                  className={`flex items-center gap-3 px-3 py-2 rounded-lg cursor-pointer transition-colors ${
                    suggestions.length + index === selectedIndex
                      ? 'bg-base-200'
                      : 'hover:bg-base-200'
                  }`}
                  onClick={() => {
                    setQuery(historyItem);
                    handleSearch(historyItem);
                  }}
                >
                  <svg
                    className="w-4 h-4 text-base-content/40"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                  </svg>
                  <span className="text-sm flex-1 truncate">{historyItem}</span>
                </div>
              ))}
              <button
                className="w-full text-left px-3 py-2 text-xs text-primary hover:bg-base-200 rounded-lg transition-colors"
                onClick={() => {
                  UnifiedSearchService.clearHistory();
                  setSearchHistory([]);
                }}
              >
                Очистить историю
              </button>
            </div>
          )}

          {/* Популярные запросы */}
          {suggestions.length === 0 &&
            searchHistory.length === 0 &&
            trendingSearches.length > 0 &&
            showTrending && (
              <div className="p-2">
                <div className="px-3 py-2 text-xs font-semibold text-base-content/60 uppercase flex items-center gap-2">
                  <TrendingIcon className="w-4 h-4" />
                  Популярные запросы
                </div>
                {trendingSearches.map((trending, index) => (
                  <div
                    key={`trending-${index}`}
                    className="flex items-center gap-3 px-3 py-2 rounded-lg cursor-pointer hover:bg-base-200 transition-colors"
                    onClick={() => {
                      setQuery(trending);
                      handleSearch(trending);
                    }}
                  >
                    <span className="text-sm flex-1">{trending}</span>
                  </div>
                ))}
              </div>
            )}
        </div>
      )}
    </div>
  );
}
