'use client';

import React, { useState, useEffect, useRef, useCallback } from 'react';
import {
  FiSearch,
  FiTrendingUp,
  FiClock,
  FiX,
  FiFolder,
  FiShoppingBag,
} from 'react-icons/fi';
import { motion, AnimatePresence } from 'framer-motion';
import { useDebounce } from '@/hooks/useDebounce';
import api from '@/services/api';
import { useRouter } from 'next/navigation';

interface SearchAutocompleteProps {
  placeholder?: string;
  selectedCategory?: string | number;
  locale: string;
  onSearch?: (query: string) => void;
  className?: string;
}

interface Suggestion {
  type: 'query' | 'product' | 'category';
  value: string;
  label: string;
  count?: number;
  category_id?: number;
  product_id?: number;
  icon?: string;
  metadata?: {
    price?: number;
    image?: string;
    category?: string;
    source_type?: string;
    storefront_id?: number;
    storefront?: string;
    storefront_slug?: string;
    parent_id?: number | null;
    last_searched?: string;
  };
}

export const SearchAutocomplete: React.FC<SearchAutocompleteProps> = ({
  placeholder = 'Поиск среди 2 млн товаров...',
  selectedCategory = 'all',
  locale,
  onSearch,
  className = '',
}) => {
  const [query, setQuery] = useState('');
  const [suggestions, setSuggestions] = useState<Suggestion[]>([]);
  const [recentSearches, setRecentSearches] = useState<string[]>([]);
  const [trendingSearches] = useState<string[]>([
    'iPhone 15',
    'PS5',
    'Квартира центр',
    'MacBook',
    'Электросамокат',
  ]);
  const [isLoading, setIsLoading] = useState(false);
  const [showSuggestions, setShowSuggestions] = useState(false);
  const [selectedIndex, setSelectedIndex] = useState(-1);

  const debouncedQuery = useDebounce(query, 300);
  const searchRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const router = useRouter();

  // Загрузка недавних поисков из localStorage
  useEffect(() => {
    const stored = localStorage.getItem('recentSearches');
    if (stored) {
      setRecentSearches(JSON.parse(stored).slice(0, 5));
    }
  }, []);

  // Закрытие подсказок при клике вне компонента
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        searchRef.current &&
        !searchRef.current.contains(event.target as Node)
      ) {
        setShowSuggestions(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  // Загрузка подсказок при изменении запроса
  useEffect(() => {
    const fetchSuggestions = async () => {
      if (!debouncedQuery || debouncedQuery.length < 2) {
        setSuggestions([]);
        return;
      }

      setIsLoading(true);
      try {
        // Запрос на получение расширенных подсказок с типами
        const response = await api.get('/api/v1/search/suggestions', {
          params: {
            q: debouncedQuery,
            limit: 10,
            types: 'queries,categories,products',
            category:
              selectedCategory !== 'all' ? String(selectedCategory) : undefined,
            lang: locale,
          },
        });

        if (response.data?.data && Array.isArray(response.data.data)) {
          // Backend возвращает массив расширенных подсказок
          setSuggestions(response.data.data);
        } else {
          // Если API не возвращает подсказки, используем поиск
          const searchResponse = await api.get('/api/v1/search', {
            params: {
              query: debouncedQuery,
              size: 5,
              page: 1,
              language: locale,
              status: 'active',
            },
          });

          if (searchResponse.data?.items) {
            const productSuggestions: Suggestion[] = searchResponse.data.items
              .slice(0, 5)
              .map((item: any) => ({
                type: 'product' as const,
                value: item.name || item.title,
                label: item.name || item.title,
                product_id: item.id,
                metadata: {
                  price: item.price,
                  image: item.images?.[0]?.url || item.images?.[0]?.public_url,
                  category: item.category?.name,
                  source_type: 'marketplace',
                },
              }));

            // Добавляем поисковый запрос как первую подсказку
            setSuggestions([
              {
                type: 'query' as const,
                value: debouncedQuery,
                label: debouncedQuery,
                count: searchResponse.data.total,
              },
              ...productSuggestions,
            ]);
          }
        }
      } catch (error) {
        console.error('Failed to fetch suggestions:', error);
        // В случае ошибки показываем только поисковый запрос
        setSuggestions([
          {
            type: 'query' as const,
            value: debouncedQuery,
            label: debouncedQuery,
          },
        ]);
      } finally {
        setIsLoading(false);
      }
    };

    fetchSuggestions();
  }, [debouncedQuery, selectedCategory, locale]);

  // Обработка клавиатурной навигации
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (!showSuggestions) return;

    switch (e.key) {
      case 'ArrowDown':
        e.preventDefault();
        setSelectedIndex((prev) =>
          prev < suggestions.length - 1 ? prev + 1 : prev
        );
        break;
      case 'ArrowUp':
        e.preventDefault();
        setSelectedIndex((prev) => (prev > 0 ? prev - 1 : -1));
        break;
      case 'Enter':
        e.preventDefault();
        if (selectedIndex >= 0 && suggestions[selectedIndex]) {
          handleSuggestionClick(suggestions[selectedIndex]);
        } else {
          handleSearch();
        }
        break;
      case 'Escape':
        setShowSuggestions(false);
        setSelectedIndex(-1);
        break;
    }
  };

  // Сохранение поиска в историю
  const saveToHistory = useCallback(
    (searchQuery: string) => {
      const updated = [
        searchQuery,
        ...recentSearches.filter((s) => s !== searchQuery),
      ].slice(0, 5);
      setRecentSearches(updated);
      localStorage.setItem('recentSearches', JSON.stringify(updated));
    },
    [recentSearches]
  );

  // Обработка клика по подсказке
  const handleSuggestionClick = (suggestion: Suggestion) => {
    if (suggestion.type === 'product') {
      // Переход к товару
      if (
        suggestion.metadata?.source_type === 'storefront' &&
        suggestion.metadata?.storefront_slug
      ) {
        // Товар из витрины
        router.push(
          `/${locale}/storefronts/${suggestion.metadata.storefront_slug}/products/${suggestion.product_id}`
        );
      } else {
        // Обычное объявление - правильный путь /marketplace/[id]
        router.push(`/${locale}/marketplace/${suggestion.product_id}`);
      }
    } else if (suggestion.type === 'category') {
      // Переход к категории
      router.push(`/${locale}/search?category=${suggestion.category_id}`);
    } else {
      // Обычный поиск (тип 'query')
      setQuery(suggestion.value);
      handleSearch(suggestion.value);
    }
    setShowSuggestions(false);
    setSelectedIndex(-1);
  };

  // Выполнение поиска
  const handleSearch = useCallback(
    (searchQuery?: string) => {
      const finalQuery = searchQuery || query;
      if (!finalQuery.trim()) return;

      saveToHistory(finalQuery);

      if (onSearch) {
        onSearch(finalQuery);
      } else {
        // Переход на страницу поиска
        const params = new URLSearchParams();
        params.append('q', finalQuery);
        if (selectedCategory !== 'all') {
          params.append('category', String(selectedCategory));
        }
        router.push(`/${locale}/search?${params.toString()}`);
      }

      setShowSuggestions(false);
    },
    [query, selectedCategory, locale, onSearch, router, saveToHistory]
  );

  // Удаление из истории поиска
  const removeFromHistory = (searchQuery: string) => {
    const updated = recentSearches.filter((s) => s !== searchQuery);
    setRecentSearches(updated);
    localStorage.setItem('recentSearches', JSON.stringify(updated));
  };

  // Очистка поля поиска
  const clearSearch = () => {
    setQuery('');
    setSuggestions([]);
    inputRef.current?.focus();
  };

  return (
    <div ref={searchRef} className={`relative ${className}`}>
      <div className="relative">
        <input
          ref={inputRef}
          type="text"
          value={query}
          onChange={(e) => {
            setQuery(e.target.value);
            setShowSuggestions(true);
            setSelectedIndex(-1);
          }}
          onFocus={() => setShowSuggestions(true)}
          onKeyDown={handleKeyDown}
          placeholder={placeholder}
          className="input input-bordered w-full pr-20"
        />

        {/* Иконки в поле ввода */}
        <div className="absolute right-2 top-1/2 -translate-y-1/2 flex items-center gap-1">
          {query && (
            <button
              onClick={clearSearch}
              className="btn btn-ghost btn-circle btn-xs"
              aria-label="Очистить"
            >
              <FiX className="w-4 h-4" />
            </button>
          )}
          <button
            onClick={() => handleSearch()}
            className="btn btn-ghost btn-circle btn-xs"
            aria-label="Поиск"
          >
            <FiSearch className="w-4 h-4" />
          </button>
        </div>
      </div>

      {/* Выпадающий список подсказок */}
      <AnimatePresence>
        {showSuggestions && (
          <motion.div
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -10 }}
            transition={{ duration: 0.2 }}
            className="absolute top-full left-0 right-0 z-50 mt-1 bg-base-100 rounded-lg shadow-xl border border-base-300 max-h-[500px] overflow-y-auto"
          >
            {/* Индикатор загрузки */}
            {isLoading && (
              <div className="flex items-center justify-center p-4">
                <span className="loading loading-spinner loading-sm"></span>
                <span className="ml-2 text-sm">Поиск...</span>
              </div>
            )}

            {/* Подсказки на основе запроса */}
            {!isLoading && suggestions.length > 0 && (
              <div className="p-2">
                <div className="text-xs font-semibold text-base-content/60 px-3 py-1">
                  Результаты поиска
                </div>
                {suggestions.map((suggestion, index) => (
                  <button
                    key={`${suggestion.type}-${suggestion.value}-${index}`}
                    onClick={() => handleSuggestionClick(suggestion)}
                    onMouseEnter={() => setSelectedIndex(index)}
                    className={`w-full text-left px-3 py-2 rounded-lg hover:bg-base-200 transition-colors flex items-center gap-3 ${
                      selectedIndex === index ? 'bg-base-200' : ''
                    }`}
                  >
                    {/* Иконка или изображение */}
                    {suggestion.type === 'product' &&
                    suggestion.metadata?.image ? (
                      <img
                        src={`http://localhost:3000${suggestion.metadata.image}`}
                        alt=""
                        className="w-10 h-10 object-cover rounded"
                      />
                    ) : (
                      <div className="w-10 h-10 flex items-center justify-center bg-base-200 rounded">
                        {suggestion.type === 'query' && (
                          <FiSearch className="w-5 h-5" />
                        )}
                        {suggestion.type === 'category' && (
                          <FiFolder className="w-5 h-5 text-primary" />
                        )}
                        {suggestion.type === 'product' && (
                          <FiShoppingBag className="w-5 h-5 text-success" />
                        )}
                      </div>
                    )}

                    <div className="flex-1 min-w-0">
                      <div className="font-medium truncate flex items-center gap-2">
                        {suggestion.label}
                        {suggestion.count !== undefined &&
                          suggestion.count > 0 && (
                            <span className="text-xs text-base-content/60">
                              ({suggestion.count})
                            </span>
                          )}
                      </div>

                      {/* Дополнительная информация */}
                      {suggestion.type === 'product' &&
                        suggestion.metadata?.category && (
                          <div className="text-xs text-base-content/60">
                            в {suggestion.metadata.category}
                          </div>
                        )}
                      {suggestion.type === 'category' && (
                        <div className="text-xs text-success">Категория</div>
                      )}
                      {suggestion.metadata?.storefront && (
                        <div className="text-xs text-info">
                          Магазин: {suggestion.metadata.storefront}
                        </div>
                      )}
                    </div>

                    {/* Цена для товаров */}
                    {suggestion.type === 'product' &&
                      suggestion.metadata?.price && (
                        <div className="text-sm font-bold text-primary">
                          {suggestion.metadata.price} РСД
                        </div>
                      )}
                  </button>
                ))}
              </div>
            )}

            {/* Недавние поиски (когда нет запроса) */}
            {!isLoading && !query && recentSearches.length > 0 && (
              <div className="p-2">
                <div className="text-xs font-semibold text-base-content/60 px-3 py-1 flex items-center gap-1">
                  <FiClock className="w-3 h-3" />
                  Недавние поиски
                </div>
                {recentSearches.map((search) => (
                  <div
                    key={search}
                    className="flex items-center gap-2 px-3 py-2 rounded-lg hover:bg-base-200 transition-colors"
                  >
                    <button
                      onClick={() => {
                        setQuery(search);
                        handleSearch(search);
                      }}
                      className="flex-1 text-left"
                    >
                      {search}
                    </button>
                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        removeFromHistory(search);
                      }}
                      className="btn btn-ghost btn-circle btn-xs"
                    >
                      <FiX className="w-3 h-3" />
                    </button>
                  </div>
                ))}
              </div>
            )}

            {/* Популярные запросы (когда нет запроса) */}
            {!isLoading && !query && (
              <div className="p-2 border-t border-base-300">
                <div className="text-xs font-semibold text-base-content/60 px-3 py-1 flex items-center gap-1">
                  <FiTrendingUp className="w-3 h-3" />
                  Популярные запросы
                </div>
                <div className="flex flex-wrap gap-2 p-3">
                  {trendingSearches.map((search) => (
                    <button
                      key={search}
                      onClick={() => {
                        setQuery(search);
                        handleSearch(search);
                      }}
                      className="btn btn-sm btn-outline"
                    >
                      {search}
                    </button>
                  ))}
                </div>
              </div>
            )}
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
};
