'use client';

import { useState, useCallback, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { MagnifyingGlassIcon, XMarkIcon } from '@heroicons/react/24/outline';
import { useDebounce } from '@/hooks/useDebounce';
import {
  translationAdminApi,
  type Translation,
} from '@/services/translationAdminApi';

interface TranslationSearchProps {
  onResultSelect?: (translation: Translation) => void;
  className?: string;
}

export default function TranslationSearch({
  onResultSelect,
  className = '',
}: TranslationSearchProps) {
  const t = useTranslations('admin');

  const [query, setQuery] = useState('');
  const [results, setResults] = useState<Translation[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [showResults, setShowResults] = useState(false);
  const [selectedIndex, setSelectedIndex] = useState(-1);

  const debouncedQuery = useDebounce(query, 300);

  // Поиск переводов
  const searchTranslations = useCallback(async (searchQuery: string) => {
    if (!searchQuery || searchQuery.length < 2) {
      setResults([]);
      return;
    }

    setIsLoading(true);
    try {
      const response = await translationAdminApi.searchTranslations({
        query: searchQuery,
        limit: 20,
      });

      if (response.success) {
        setResults(response.data || []);
      }
    } catch (error) {
      console.error('Search error:', error);
      setResults([]);
    } finally {
      setIsLoading(false);
    }
  }, []);

  // Эффект для выполнения поиска
  useEffect(() => {
    searchTranslations(debouncedQuery);
  }, [debouncedQuery, searchTranslations]);

  // Обработка выбора результата
  const handleResultClick = (translation: Translation) => {
    onResultSelect?.(translation);
    setQuery('');
    setResults([]);
    setShowResults(false);
    setSelectedIndex(-1);
  };

  // Обработка клавиатурной навигации
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (!showResults || results.length === 0) return;

    switch (e.key) {
      case 'ArrowDown':
        e.preventDefault();
        setSelectedIndex((prev) => (prev < results.length - 1 ? prev + 1 : 0));
        break;

      case 'ArrowUp':
        e.preventDefault();
        setSelectedIndex((prev) => (prev > 0 ? prev - 1 : results.length - 1));
        break;

      case 'Enter':
        e.preventDefault();
        if (selectedIndex >= 0 && selectedIndex < results.length) {
          handleResultClick(results[selectedIndex]);
        }
        break;

      case 'Escape':
        setShowResults(false);
        setSelectedIndex(-1);
        break;
    }
  };

  // Подсветка совпадений в тексте
  const highlightMatch = (text: string, searchQuery: string) => {
    if (!searchQuery) return text;

    const parts = text.split(new RegExp(`(${searchQuery})`, 'gi'));
    return (
      <>
        {parts.map((part, index) =>
          part.toLowerCase() === searchQuery.toLowerCase() ? (
            <mark key={index} className="bg-yellow-200 text-base-content">
              {part}
            </mark>
          ) : (
            <span key={index}>{part}</span>
          )
        )}
      </>
    );
  };

  // Определение типа элемента
  const getItemType = (key: string): string => {
    if (key.startsWith('category.')) return t('types.category');
    if (key.startsWith('attribute.')) return t('types.attribute');
    if (key.startsWith('page.')) return t('types.page');
    if (key.startsWith('notification.')) return t('types.notification');
    if (key.startsWith('email.')) return t('types.email');
    return t('types.other');
  };

  return (
    <div className={`relative ${className}`}>
      {/* Поисковая строка */}
      <div className="relative">
        <input
          type="text"
          value={query}
          onChange={(e) => {
            setQuery(e.target.value);
            setShowResults(true);
          }}
          onFocus={() => setShowResults(true)}
          onKeyDown={handleKeyDown}
          placeholder={t('placeholder')}
          className="input input-bordered w-full pl-10 pr-10"
        />

        <MagnifyingGlassIcon className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-base-content/50" />

        {query && (
          <button
            onClick={() => {
              setQuery('');
              setResults([]);
              setShowResults(false);
            }}
            className="absolute right-3 top-1/2 -translate-y-1/2"
          >
            <XMarkIcon className="h-5 w-5 text-base-content/50 hover:text-base-content" />
          </button>
        )}
      </div>

      {/* Результаты поиска */}
      {showResults && (query.length >= 2 || results.length > 0) && (
        <div className="absolute z-50 w-full mt-1 bg-base-100 border border-base-300 rounded-lg shadow-lg max-h-96 overflow-y-auto">
          {isLoading ? (
            <div className="p-4 text-center">
              <span className="loading loading-spinner loading-sm"></span>
              <span className="ml-2">{t('searching')}</span>
            </div>
          ) : results.length > 0 ? (
            <ul className="py-2">
              {results.map((translation, index) => (
                <li key={translation.id}>
                  <button
                    onClick={() => handleResultClick(translation)}
                    className={`w-full px-4 py-3 text-left hover:bg-base-200 transition-colors ${
                      index === selectedIndex ? 'bg-base-200' : ''
                    }`}
                  >
                    {/* Ключ и тип */}
                    <div className="flex items-center justify-between mb-1">
                      <span className="font-mono text-sm text-primary">
                        {highlightMatch(
                          `${translation.entity_type}.${translation.entity_id}.${translation.field_name}`,
                          debouncedQuery
                        )}
                      </span>
                      <span className="badge badge-ghost badge-sm">
                        {getItemType(translation.entity_type)}
                      </span>
                    </div>

                    {/* Переводы */}
                    <div className="space-y-1">
                      <div className="text-sm">
                        <span className="font-semibold text-base-content/70">
                          {translation.language.toUpperCase()}:
                        </span>{' '}
                        <span className="text-base-content/90">
                          {highlightMatch(
                            translation.translated_text,
                            debouncedQuery
                          )}
                        </span>
                      </div>
                    </div>

                    {/* Дополнительная информация */}
                    {translation.is_verified && (
                      <div className="mt-1">
                        <span className="badge badge-success badge-xs">
                          Verified
                        </span>
                      </div>
                    )}
                  </button>
                </li>
              ))}
            </ul>
          ) : query.length >= 2 ? (
            <div className="p-4 text-center text-base-content/60">
              {t('noResults')}
            </div>
          ) : null}
        </div>
      )}
    </div>
  );
}
