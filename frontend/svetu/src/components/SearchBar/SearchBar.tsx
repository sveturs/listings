'use client';

import { useState, useEffect, useRef, useCallback } from 'react';
import { useTranslations, useLocale } from 'next-intl';
import { useRouter } from 'next/navigation';
import { useDebounce } from '@/hooks/useDebounce';
import { useBehaviorTracking } from '@/hooks/useBehaviorTracking';
import {
  UnifiedSearchService,
  SearchSuggestion,
  EnhancedSuggestion,
} from '@/services/unifiedSearch';
import { SearchIcon, CloseIcon } from './icons';
import SearchAutocomplete from './SearchAutocomplete';

interface SearchBarProps {
  placeholder?: string;
  initialQuery?: string;
  onSearch?: (query: string, fuzzy?: boolean) => void;
  className?: string;
  variant?: 'default' | 'hero' | 'minimal';
  showTrending?: boolean;
  fuzzy?: boolean;
  onFuzzyChange?: (fuzzy: boolean) => void;
  geoLocation?: {
    lat: number;
    lon: number;
    radius: number; // в метрах
  };
}

export default function SearchBar({
  placeholder,
  initialQuery = '',
  onSearch,
  className = '',
  variant = 'default',
  showTrending = false,
  fuzzy: initialFuzzy = true,
  onFuzzyChange,
  geoLocation,
}: SearchBarProps) {
  const t = useTranslations('search');
  const locale = useLocale();
  const router = useRouter();
  const [query, setQuery] = useState(initialQuery);
  const [fuzzy, setFuzzy] = useState(initialFuzzy);
  const [suggestions, setSuggestions] = useState<
    (SearchSuggestion | EnhancedSuggestion)[]
  >([]);
  const [isLoadingSuggestions, setIsLoadingSuggestions] = useState(false);
  const [showSuggestions, setShowSuggestions] = useState(false);
  const [selectedIndex, setSelectedIndex] = useState(-1);
  const [searchHistory, setSearchHistory] = useState<string[]>([]);
  const [trendingSearches, setTrendingSearches] = useState<string[]>([]);
  const debouncedQuery = useDebounce(query, 300);
  const inputRef = useRef<HTMLInputElement>(null);
  const suggestionRef = useRef<HTMLDivElement>(null);

  // Behavior tracking
  const {
    trackSearchPerformed,
    trackSearchFilterApplied,
    trackResultClicked,
    startSearch,
  } = useBehaviorTracking();

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
    } else {
      setSuggestions([]);
      // Не показываем предложения автоматически при загрузке
      if (debouncedQuery.length > 0) {
        setShowSuggestions(false);
      }
    }
  }, [debouncedQuery]);

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
    async (searchQuery: string = query) => {
      const trimmedQuery = searchQuery.trim();
      if (!trimmedQuery) return;

      const searchStartTime = Date.now();
      setShowSuggestions(false);
      UnifiedSearchService.saveToHistory(trimmedQuery);

      // Подготавливаем фильтры для трекинга
      const searchFilters = {
        fuzzy,
        product_types: ['marketplace', 'storefront'], // дефолтные типы
      };

      // Запускаем трекинг поиска
      startSearch(trimmedQuery, searchFilters);

      try {
        // Выполняем поиск для получения количества результатов (если доступно)
        let resultsCount = 0;
        if (!onSearch) {
          // Если это прямой поиск (не через callback), пытаемся получить предварительные результаты
          try {
            const searchResults = await UnifiedSearchService.search({
              query: trimmedQuery,
              fuzzy,
              page: 1,
              limit: 1, // Минимальный запрос для получения total
            });
            resultsCount = searchResults.total;
          } catch (error) {
            console.warn('Failed to get results count for tracking:', error);
          }
        }

        // Трекинг выполненного поиска
        await trackSearchPerformed({
          search_query: trimmedQuery,
          search_filters: searchFilters,
          results_count: resultsCount,
          search_duration_ms: Date.now() - searchStartTime,
        });

        if (onSearch) {
          onSearch(trimmedQuery, fuzzy);
        } else {
          // Строим URL с правильной локалью и параметром fuzzy
          const searchUrl = new URL(
            `/${locale}/search`,
            window.location.origin
          );
          searchUrl.searchParams.set('q', trimmedQuery);
          searchUrl.searchParams.set('fuzzy', fuzzy.toString());

          // Добавляем геопараметры если они есть
          if (geoLocation) {
            searchUrl.searchParams.set('lat', geoLocation.lat.toString());
            searchUrl.searchParams.set('lon', geoLocation.lon.toString());
            const radiusKm = Math.round(geoLocation.radius / 1000);
            searchUrl.searchParams.set('distance', `${radiusKm}km`);
          }

          router.push(searchUrl.toString());
        }
      } catch (error) {
        console.error('Search tracking error:', error);
        // Продолжаем выполнение поиска даже при ошибке трекинга
        if (onSearch) {
          onSearch(trimmedQuery, fuzzy);
        } else {
          const searchUrl = new URL(
            `/${locale}/search`,
            window.location.origin
          );
          searchUrl.searchParams.set('q', trimmedQuery);
          searchUrl.searchParams.set('fuzzy', fuzzy.toString());

          // Добавляем геопараметры если они есть
          if (geoLocation) {
            searchUrl.searchParams.set('lat', geoLocation.lat.toString());
            searchUrl.searchParams.set('lon', geoLocation.lon.toString());
            const radiusKm = Math.round(geoLocation.radius / 1000);
            searchUrl.searchParams.set('distance', `${radiusKm}km`);
          }

          router.push(searchUrl.toString());
        }
      }
    },
    [
      query,
      fuzzy,
      onSearch,
      router,
      locale,
      trackSearchPerformed,
      startSearch,
      geoLocation,
    ]
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

  const handleSuggestionClick = async (
    suggestion: SearchSuggestion | EnhancedSuggestion
  ) => {
    setShowSuggestions(false);

    // Трекинг клика по результату поиска (предложению)
    try {
      if (query) {
        // Определяем позицию предложения в списке
        const position =
          suggestions.findIndex((s) => s.text === suggestion.text) + 1;

        // Определяем тип элемента на основе metadata
        let itemType = 'marketplace';
        if (
          suggestion.metadata &&
          'source_type' in suggestion.metadata &&
          suggestion.metadata.source_type === 'storefront'
        ) {
          itemType = 'storefront';
        }

        const itemId =
          'product' in suggestion && suggestion.product
            ? suggestion.product.id.toString()
            : 'product_id' in suggestion && suggestion.product_id
              ? suggestion.product_id.toString()
              : suggestion.text;

        await trackResultClicked({
          search_query: query,
          clicked_item_id: itemId,
          click_position: position,
          total_results: suggestions.length,
          click_time_from_search_ms: Date.now() - Date.now(), // приблизительное время
          item_type: itemType as any,
        });
      }
    } catch (error) {
      console.error('Failed to track suggestion click:', error);
    }

    // Если это категория, переходим на страницу категории
    if (suggestion.type === 'category' && suggestion.category) {
      router.push(
        `/${locale}/category/${suggestion.category.slug || suggestion.category.id}`
      );
      return;
    }

    // Если это товар (может быть EnhancedSuggestion или SearchSuggestion)
    const productId =
      'product' in suggestion && suggestion.product
        ? suggestion.product.id
        : 'product_id' in suggestion
          ? suggestion.product_id
          : undefined;

    if (suggestion.type === 'product' && productId) {
      // Проверяем тип товара через metadata
      if (
        suggestion.metadata &&
        'source_type' in suggestion.metadata &&
        suggestion.metadata.source_type === 'storefront' &&
        'storefront_id' in suggestion.metadata &&
        suggestion.metadata.storefront_id
      ) {
        // Для товаров витрин нужен правильный URL
        // Получаем slug витрины из metadata или делаем запрос
        const storefrontSlug =
          'storefront_slug' in suggestion.metadata
            ? suggestion.metadata.storefront_slug
            : undefined;
        if (storefrontSlug) {
          router.push(`/${locale}/b2c/${storefrontSlug}/products/${productId}`);
        } else {
          // Если нет slug, переходим на страницу поиска с фильтром
          router.push(
            `/${locale}/search?q=${encodeURIComponent(suggestion.text)}&product_types=storefront`
          );
        }
      } else {
        // Для товаров маркетплейса
        router.push(`/${locale}/c2c/${productId}`);
      }
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

  const handleFuzzyToggle = async () => {
    const newFuzzy = !fuzzy;
    setFuzzy(newFuzzy);

    // Трекинг изменения фильтра fuzzy поиска
    if (query) {
      try {
        await trackSearchFilterApplied({
          search_query: query,
          filter_type: 'fuzzy_search',
          filter_value: newFuzzy.toString(),
          results_count_before: 0, // TODO: можно получить из текущих результатов
          results_count_after: 0, // TODO: можно получить после применения фильтра
        });
      } catch (error) {
        console.error('Failed to track fuzzy filter change:', error);
      }
    }

    if (onFuzzyChange) {
      onFuzzyChange(newFuzzy);
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
              aria-label={isLoadingSuggestions ? t('searching') : t('searchButton')}
            >
              {isLoadingSuggestions ? (
                <span className="loading loading-spinner loading-xs"></span>
              ) : (
                <SearchIcon className="w-4 h-4" />
              )}
            </button>
          </div>
        </div>

        {/* Fuzzy search toggle for hero and minimal variants */}
        {(variant === 'hero' || variant === 'minimal') && (
          <div
            className={`flex items-center ${variant === 'hero' ? 'justify-center mt-3' : 'mt-2'} gap-2`}
          >
            <label className="label cursor-pointer flex items-center gap-2 p-0">
              <input
                type="checkbox"
                className="checkbox checkbox-xs checkbox-primary"
                checked={fuzzy}
                onChange={handleFuzzyToggle}
              />
              <span className="label-text text-xs">
                {t('fuzzySearch')}
                {variant === 'hero' && (
                  <span className="text-xs text-base-content/60 ml-1">
                    ({t('fuzzySearchHint')})
                  </span>
                )}
              </span>
              {variant === 'minimal' && (
                <div
                  className="tooltip tooltip-bottom"
                  data-tip={t('fuzzySearchTooltip')}
                >
                  <svg
                    className="w-4 h-4 text-base-content/60"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                  </svg>
                </div>
              )}
            </label>
          </div>
        )}
      </div>

      {/* Автодополнение с новым компонентом */}
      <SearchAutocomplete
        suggestions={suggestions}
        searchHistory={searchHistory}
        trendingSearches={trendingSearches}
        showSuggestions={showSuggestions}
        selectedIndex={selectedIndex}
        query={query}
        isLoading={isLoadingSuggestions}
        onSelect={(value) => {
          setQuery(value);
          handleSearch(value);
        }}
        onCategorySelect={(categoryId) => {
          router.push(`/${locale}/category/${categoryId}`);
        }}
        onProductSelect={(productId) => {
          router.push(`/${locale}/c2c/${productId}`);
        }}
      />
    </div>
  );
}
